package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *config) handlerGetReviews(w http.ResponseWriter, r *http.Request) {
	juiceID := r.PathValue("juiceID")
	parsedJuiceID, err := uuid.Parse(juiceID)
	if err != nil {
		respondWithError(w, 400, "invalid juice ID")
		return
	}

	reviewData, err := cfg.queries.GetJuiceReviews(r.Context(), parsedJuiceID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "no available data")
		return
	}
	respondWithJSON(w, http.StatusOK, reviewData)
}
