package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *config) handlerDeleteReview(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(contextKeyUserID).(uuid.UUID)
	role := r.Context().Value(contextKeyRole).(string)

	reviewID := r.PathValue("reviewID")
	parsedID, err := uuid.Parse(reviewID)
	if err != nil {
		respondWithError(w, 400, "invalid review ID")
		return
	}

	review, err := cfg.queries.GetReviewByID(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "review not found")
		return
	}

	if review.UserID != userID && role != "admin" {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	err = cfg.queries.DeleteReview(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete review")
		return
	}

	w.WriteHeader(204)
}
