package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/karim-h00/juiceshop-clone/internal/database"
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

func (cfg *config) handlerAddReview(w http.ResponseWriter, r *http.Request) {
	type reviewParams struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}

	userID := r.Context().Value(contextKeyUserID).(uuid.UUID)
	juiceID := r.PathValue("juiceID")
	parsedJuiceID, err := uuid.Parse(juiceID)
	if err != nil {
		respondWithError(w, 400, "invalid juice ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := reviewParams{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "error decoding params")
		return
	}

	if params.Rating < 1 || params.Rating > 5 {
		respondWithError(w, 400, "rating must be between 1 and 5")
		return
	}

	review, err := cfg.queries.AddReview(r.Context(), database.AddReviewParams{
		UserID:  userID,
		JuiceID: parsedJuiceID,
		Rating:  int32(params.Rating),
		Comment: sql.NullString{String: params.Comment, Valid: params.Comment != ""},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not add review")
		return
	}
	respondWithJSON(w, 201, review)
}

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
