package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/karim-h00/juiceshop-clone/internal/auth"
	"github.com/karim-h00/juiceshop-clone/internal/database"
)

func (cfg *config) handlerOrderJuice(w http.ResponseWriter, r *http.Request) {
	type order_params struct {
		Items []struct {
			JuiceID  string `json:"juice_id"`
			Quantity int    `json:"quantity"`
		} `json:"items"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	userID, _, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 400, "Could not make session")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := order_params{}
	err = decoder.Decode(&params)
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
