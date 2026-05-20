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

	juiceIDs := make([]uuid.UUID, len(params.Items))
	for i, item := range params.Items {
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
	priceMap := make(map[uuid.UUID]int32, len(juices))
	for _, j := range juices {
		priceMap[j.ID] = j.Price
	}

	var computedTotal int32
	for _, item := range params.Items {
		id := uuid.MustParse(item.JuiceID)
		price, ok := priceMap[id]
		if !ok {
			respondWithError(w, 400, "Juice not found: "+item.JuiceID)
			return
		}
		computedTotal += price * int32(item.Quantity)
	}

	order, err := cfg.queries.CreateOrder(r.Context(), database.CreateOrderParams{
		UserID: userID,
		Total:  computedTotal,
	})
	if err != nil {
		respondWithError(w, 500, "Could not create order")
		return
	}

	for _, item := range params.Items {
		_, err := cfg.queries.CreateOrderItem(r.Context(), database.CreateOrderItemParams{
			OrderID:  order.ID,
			JuiceID:  uuid.MustParse(item.JuiceID),
			Quantity: int32(item.Quantity),
		})
		if err != nil {
			respondWithError(w, 500, "Could not create order item")
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
		id := uuid.MustParse(item.JuiceID)
		itemsResp[i] = orderItemResponse{
			JuiceID:  item.JuiceID,
			Quantity: item.Quantity,
			Price:    priceMap[id],
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
