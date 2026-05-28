package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *config) handlerGetOrderByID(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(contextKeyUserID).(uuid.UUID)
	orderID := r.PathValue("orderID")
	parsedID, err := uuid.Parse(orderID)
	if err != nil {
		respondWithError(w, 400, "failed to parse ID")
		return
	}

	order, err := cfg.queries.GetOrderByOrderID(r.Context(), userID)
	if err != nil {
		respondWithError(w, 404, "order not found")
		return
	}
	if order.ID != userID {
		respondWithError(w, 403, "forbidden")
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
