package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/karim-h00/juiceshop-clone/internal/database"
)

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
		cfg.logger.Warn("order juice", "reason", "failed to decode params", "user_id", userID, "ip", r.RemoteAddr)
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
		cfg.logger.Error("order juice", "reason", "failed to fetch juices", "user_id", userID, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Could not fetch juices")
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

	tx, err := cfg.db.BeginTx(r.Context(), nil)
	if err != nil {
		cfg.logger.Error("order juice", "reason", "failed to start transaction", "user_id", userID, "error", err)
		respondWithError(w, 500, "Could not start transaction")
		return
	}
	defer tx.Rollback()

	qtx := cfg.queries.WithTx(tx)

	order, err := qtx.CreateOrder(r.Context(), database.CreateOrderParams{
		UserID: userID,
		Total:  computedTotal,
	})
	if err != nil {
		cfg.logger.Error("order juice", "reason", "failed to create order", "user_id", userID, "error", err)
		respondWithError(w, 500, "Could not create order")
		return
	}

	for i, item := range params.Items {
		_, err := qtx.CreateOrderItem(r.Context(), database.CreateOrderItemParams{
			OrderID:  order.ID,
			JuiceID:  juiceIDs[i],
			Quantity: int32(item.Quantity),
		})
		if err != nil {
			cfg.logger.Error("order juice", "reason", "failed to create order item", "user_id", userID, "order_id", order.ID, "error", err)
			respondWithError(w, http.StatusInternalServerError, "Could not create order item")
			return
		}
		err = qtx.DecrementJuiceStock(r.Context(), database.DecrementJuiceStockParams{
			ID:    juiceIDs[i],
			Stock: int32(item.Quantity),
		})
		if err != nil {
			cfg.logger.Error("order juice", "reason", "failed to decrement stock", "user_id", userID, "order_id", order.ID, "error", err)
			respondWithError(w, http.StatusInternalServerError, "Could not update stock")
			return
		}
	}
	if err = tx.Commit(); err != nil {
		cfg.logger.Error("order juice", "reason", "failed to commit transaction", "user_id", userID, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Could not complete order")
		return
	}
	now := time.Now().UTC()
	err = cfg.queries.AddLog(r.Context(), database.AddLogParams{
		UserID:     uuid.NullUUID{UUID: userID, Valid: true},
		Action:     "create_order",
		TargetType: "order",
		TargetID:   uuid.NullUUID{UUID: order.ID, Valid: true},
		TargetName: sql.NullString{},
		CreatedAt:  now,
	})
	if err != nil {
		cfg.logger.Error("add audit log", "error", err)
	}

	type orderItemResponse struct {
		JuiceID  string `json:"juice_id"`
		Quantity int    `json:"quantity"`
		Price    int32  `json:"price"`
	}
	type orderCreatedResponse struct {
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

	cfg.logger.Info("order juice successful", "user_id", userID, "order_id", order.ID, "total", computedTotal, "ip", r.RemoteAddr)
	respondWithJSON(w, 201, orderCreatedResponse{
		OrderID:   order.ID.String(),
		UserID:    userID.String(),
		Total:     computedTotal,
		CreatedAt: order.CreatedAt,
		Items:     itemsResp,
	})
}

func (cfg *config) handlerGetUserOrderHistory(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(contextKeyUserID).(uuid.UUID)

	orders, err := cfg.queries.GetOrdersByUserID(r.Context(), database.GetOrdersByUserIDParams{
		UserID: userID,
		Limit:  5,
		Offset: 0,
	})
	if err != nil {
		cfg.logger.Error("get order history", "user_id", userID, "error", err)
		respondWithError(w, 500, "Could not fetch orders")
		return
	}

	result := make([]orderResponse, len(orders))
	for i, order := range orders {
		items, err := cfg.queries.GetOrderItemsByOrderID(r.Context(), order.ID)
		if err != nil {
			cfg.logger.Error("get order history(items)", "user_id", userID, "error", err)
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

	userID := r.Context().Value(contextKeyUserID).(uuid.UUID)
	role := r.Context().Value(contextKeyRole).(string)
	orderID := r.PathValue("orderID")
	parsedID, err := uuid.Parse(orderID)
	if err != nil {
		cfg.logger.Warn("invalid order id", "order_id", orderID, "ip", r.RemoteAddr)
		respondWithError(w, 400, "failed to parse ID")
		return
	}

	order, err := cfg.queries.GetOrderByOrderID(r.Context(), parsedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			cfg.logger.Warn("get order by id", "reason", "order id not found", "order_id", orderID, "ip", r.RemoteAddr)
			respondWithError(w, 404, "order not found")
			return
		}
		cfg.logger.Error("get order by id", "error", err)
		respondWithError(w, http.StatusInternalServerError, "order not found")
		return
	}
	if role != "admin" && order.UserID != userID {
		cfg.logger.Warn("get order by id", "reason", "user unauthorized to view order", "order_id", orderID, "ip", r.RemoteAddr)
		respondWithError(w, 403, "forbidden")
		return
	}

	items, err := cfg.queries.GetOrderItemsByOrderID(r.Context(), parsedID)
	if err != nil {
		cfg.logger.Error("get order by id", "error", err, "ip", r.RemoteAddr)
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
	adminID := r.Context().Value(contextKeyUserID).(uuid.UUID)
	orderID := r.PathValue("orderID")
	parsedID, err := uuid.Parse(orderID)
	if err != nil {
		cfg.logger.Warn("invalid order id", "order_id", orderID, "ip", r.RemoteAddr)
		respondWithError(w, 400, "failed to parse ID")
		return
	}

	now := time.Now().UTC()
	err = cfg.queries.DeleteOrderByOrderID(r.Context(), parsedID)
	if err != nil {
		cfg.logger.Error("delete order", "error", err)
		respondWithError(w, 500, "Couldn't delete order")
		return
	}
	err = cfg.queries.AddLog(r.Context(), database.AddLogParams{
		UserID:     uuid.NullUUID{UUID: adminID, Valid: true},
		Action:     "delete_order",
		TargetType: "order",
		TargetID:   uuid.NullUUID{UUID: parsedID, Valid: true},
		TargetName: sql.NullString{String: "", Valid: false},
		CreatedAt:  now,
	})
	if err != nil {
		cfg.logger.Error("add audit log", "error", err)
	}
	cfg.logger.Info("delete order", "admin_id", adminID, "order_id", orderID, "ip", r.RemoteAddr)
	w.WriteHeader(204)
}

func (cfg *config) handlerAdminGetAllOrders(w http.ResponseWriter, r *http.Request) {

	page := 1

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil {
			cfg.logger.Warn("invalid page number", "page", pageStr, "ip", r.RemoteAddr)
			respondWithError(w, 400, "invalid page number")
			return
		}
		page = parsedPage
	}
	offset := (page - 1) * 10

	orderData, err := cfg.queries.GetAllOrders(r.Context(), int32(offset))
	if err != nil {
		cfg.logger.Error("get all orders", "page", page, "error", err, "ip", r.RemoteAddr)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get orders")
		return
	}
	cfg.logger.Info("get all orders", "page", page)
	respondWithJSON(w, http.StatusOK, orderData)
}
