package main

import (
	"net/http"

	"github.com/google/uuid"
)

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
