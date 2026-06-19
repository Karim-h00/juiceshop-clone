package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/karim-h00/juiceshop-clone/internal/database"
)

func (cfg *config) handlerGetReviews(w http.ResponseWriter, r *http.Request) {

	type reviewResponse struct {
		ID        uuid.UUID `json:"id"`
		UserID    uuid.UUID `json:"user_id"`
		JuiceID   uuid.UUID `json:"juice_id"`
		Rating    int32     `json:"rating"`
		Comment   *string   `json:"comment"`
		CreatedAt time.Time `json:"created_at"`
		Username  string    `json:"username"`
	}
	slug := r.PathValue("slug")
	name := slugToName(slug)

	juiceID, err := cfg.queries.GetJuiceID(r.Context(), name)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithJSON(w, http.StatusOK, []interface{}{})
			return
		}
		respondWithError(w, http.StatusInternalServerError, "failed to fetch reviews")
		return
	}
	reviewData, err := cfg.queries.GetJuiceReviews(r.Context(), juiceID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithJSON(w, http.StatusOK, []interface{}{})
			return
		}
		respondWithError(w, http.StatusInternalServerError, "failed to fetch reviews")
		return
	}
	responses := make([]reviewResponse, len(reviewData))
	for i, r := range reviewData {
		var comment *string
		if r.Comment.Valid {
			comment = &r.Comment.String
		}
		responses[i] = reviewResponse{
			ID:        r.ID,
			UserID:    r.UserID,
			JuiceID:   r.JuiceID,
			Rating:    r.Rating,
			Comment:   comment,
			CreatedAt: r.CreatedAt,
			Username:  r.Username,
		}
	}
	respondWithJSON(w, http.StatusOK, responses)
}

func (cfg *config) handlerAddReview(w http.ResponseWriter, r *http.Request) {
	type reviewParams struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}

	userID := r.Context().Value(contextKeyUserID).(uuid.UUID)
	slug := r.PathValue("slug")
	name := slugToName(slug)

	juiceID, err := cfg.queries.GetJuiceID(r.Context(), name)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithJSON(w, http.StatusOK, []interface{}{})
			return
		}
		respondWithError(w, http.StatusInternalServerError, "failed to fetch reviews")
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
		JuiceID: juiceID,
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
