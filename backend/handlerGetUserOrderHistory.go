package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
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

func (cfg *config) handlerGetUserOrderHistory(w http.ResponseWriter, r *http.Request) {

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
