package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/karim-h00/juiceshop-clone/internal/database"
)

func (cfg *config) handlerOrderJuice(w http.ResponseWriter, r *http.Request) {
	type order_params struct {
		Items []struct {
			JuiceID  string `json:"juice_id"`
			Quantity int    `json:"quantity"`
		} `json:"items"`
	}

	userID := r.Context().Value(contextKeyUserID).(uuid.UUID)

	decoder := json.NewDecoder(r.Body)
	params := order_params{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error decoding params")
		return
	}
	if len(params.Items) == 0 {
		respondWithError(w, 400, "Order must contain at least one item")
		return
	}

	juiceIDs := make([]uuid.UUID, len(params.Items))
	for i, item := range params.Items {
		if item.Quantity <= 0 {
			respondWithError(w, 400, "Quantity must be positive")
			return
		}
		parsed, err := uuid.Parse(item.JuiceID)
		if err != nil {
			respondWithError(w, 400, "Invalid juice ID")
			return
		}
		juiceIDs[i] = parsed
	}
	juices, err := cfg.queries.GetJuicesByIDs(r.Context(), juiceIDs)
	if err != nil {
		respondWithError(w, 500, "Could not fetch juices")
		return
	}

	type juiceInfo struct {
		name     string
		price    int32
		quantity int32
	}
	stockMap := make(map[uuid.UUID]juiceInfo, len(juices))
	for _, j := range juices {
		stockMap[j.ID] = juiceInfo{name: j.Name, price: j.Price, quantity: j.Stock}
	}

	var computedTotal int32
	for i, item := range params.Items {
		info, ok := stockMap[juiceIDs[i]]
		if !ok {
			respondWithError(w, 400, "Juice not found: "+item.JuiceID)
			return
		}
		if int32(item.Quantity) > info.quantity {
			respondWithError(w, 400, "Insufficient stock for: "+info.name)
			return
		}
		computedTotal += info.price * int32(item.Quantity)
	}

	order, err := cfg.queries.CreateOrder(r.Context(), database.CreateOrderParams{
		UserID: userID,
		Total:  computedTotal,
	})
	if err != nil {
		respondWithError(w, 500, "Could not create order")
		return
	}

	for i, item := range params.Items {
		_, err := cfg.queries.CreateOrderItem(r.Context(), database.CreateOrderItemParams{
			OrderID:  order.ID,
			JuiceID:  juiceIDs[i],
			Quantity: int32(item.Quantity),
		})
		if err != nil {
			respondWithError(w, 500, "Could not create order item")
			return
		}
		err = cfg.queries.DecrementJuiceStock(r.Context(), database.DecrementJuiceStockParams{
			ID:    juiceIDs[i],
			Stock: int32(item.Quantity),
		})
		if err != nil {
			respondWithError(w, 500, "Could not update stock")
			return
		}
	}

	type orderItemResponse struct {
		JuiceID  string `json:"juice_id"`
		Quantity int    `json:"quantity"`
		Price    int32  `json:"price"`
	}
	type orderResponse struct {
		OrderID   string              `json:"order_id"`
		UserID    string              `json:"user_id"`
		Total     int32               `json:"total"`
		CreatedAt time.Time           `json:"created_at"`
		Items     []orderItemResponse `json:"items"`
	}
	itemsResp := make([]orderItemResponse, len(params.Items))
	for i, item := range params.Items {
		id := juiceIDs[i]
		itemsResp[i] = orderItemResponse{
			JuiceID:  item.JuiceID,
			Quantity: item.Quantity,
			Price:    stockMap[id].price,
		}
	}

	respondWithJSON(w, 201, orderResponse{
		OrderID:   order.ID.String(),
		UserID:    userID.String(),
		Total:     computedTotal,
		CreatedAt: order.CreatedAt,
		Items:     itemsResp,
	})
}

func (cfg *config) handlerGetUserOrderHistory(w http.ResponseWriter, r *http.Request) {
	type orderItems struct {
		Name     string `json:"name"`
		Quantity int32  `json:"quantity"`
	}
	type orderResponse struct {
		OrderID   string       `json:"order_id"`
		Total     int32        `json:"total"`
		CreatedAt time.Time    `json:"created_at"`
		Items     []orderItems `json:"items"`
	}

	userID := r.Context().Value(contextKeyUserID).(uuid.UUID)

	orders, err := cfg.queries.GetOrdersByUserID(r.Context(), userID)
	if err != nil {
		respondWithError(w, 500, "Could not fetch orders")
		return
	}

	result := make([]orderResponse, len(orders))
	for i, order := range orders {
		items, err := cfg.queries.GetOrderItemsByOrderID(r.Context(), order.ID)
		if err != nil {
			respondWithError(w, 500, "Could not fetch order items")
			return
		}
		itemsResp := make([]orderItems, len(items))
		for j, item := range items {
			itemsResp[j] = orderItems{
				Name:     item.Name,
				Quantity: item.Quantity,
			}
		}
		result[i] = orderResponse{
			OrderID:   order.ID.String(),
			Total:     order.Total,
			CreatedAt: order.CreatedAt,
			Items:     itemsResp,
		}
	}
	respondWithJSON(w, 200, result)
}

func (cfg *config) handlerGetOrderByID(w http.ResponseWriter, r *http.Request) {
	type orderItems struct {
		Name     string `json:"name"`
		Quantity int32  `json:"quantity"`
	}
	type orderResponse struct {
		OrderID   string       `json:"order_id"`
		Total     int32        `json:"total"`
		CreatedAt time.Time    `json:"created_at"`
		Items     []orderItems `json:"items"`
	}

	userID := r.Context().Value(contextKeyUserID).(uuid.UUID)
	orderID := r.PathValue("orderID")
	parsedID, err := uuid.Parse(orderID)
	if err != nil {
		respondWithError(w, 400, "failed to parse ID")
		return
	}

	order, err := cfg.queries.GetOrderByOrderID(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, 404, "order not found")
		return
	}
	if order.UserID != userID {
		respondWithError(w, 403, "forbidden")
		return
	}

	items, err := cfg.queries.GetOrderItemsByOrderID(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, 500, "could not fetch order items")
		return
	}

	itemsResp := make([]orderItems, len(items))
	for j, item := range items {
		itemsResp[j] = orderItems{
			Name:     item.Name,
			Quantity: item.Quantity,
		}
	}

	respondWithJSON(w, 200, orderResponse{
		OrderID:   order.ID.String(),
		Total:     order.Total,
		CreatedAt: order.CreatedAt,
		Items:     itemsResp,
	})
}

func (cfg *config) handlerDeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("orderID")
	parsedID, err := uuid.Parse(orderID)
	if err != nil {
		respondWithError(w, 400, "failed to parse ID")
		return
	}

	err = cfg.queries.DeleteOrderByOrderID(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, 500, "Couldn't delete order")
		return
	}
	w.WriteHeader(204)
}

func (cfg *config) handlerAdminGetAllOrders(w http.ResponseWriter, r *http.Request) {

	page := 1

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil {
			respondWithError(w, 400, "invalid page number")
			return
		}
		page = parsedPage
	}
	offset := (page - 1) * 10

	orderData, err := cfg.queries.GetAllOrders(r.Context(), int32(offset))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get orders")
		return
	}

	respondWithJSON(w, http.StatusOK, orderData)
}
