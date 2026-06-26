package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
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
			respondWithError(w, http.StatusNotFound, "juice not found")
			return
		}
		cfg.logger.Error("get juiceID for reviews", "error", err)
		respondWithError(w, http.StatusInternalServerError, "failed to fetch reviews")
		return
	}
	reviewData, err := cfg.queries.GetJuiceReviews(r.Context(), database.GetJuiceReviewsParams{
		JuiceID: juiceID,
		Limit:   20,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithJSON(w, http.StatusOK, []reviewResponse{})
			return
		}
		cfg.logger.Error("get reviews", "error", err)
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
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "juice is unavailable")
			return
		}
		cfg.logger.Error("get juiceID to add reviews", "user_id", userID, "error", err)
		respondWithError(w, http.StatusInternalServerError, "failed to fetch reviews")
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := reviewParams{}
	err = decoder.Decode(&params)
	if err != nil {
		cfg.logger.Warn("add review", "user_id", userID, "reason", "failed to decode body", "error", err, "ip", getClientIP(r))
		respondWithError(w, 400, "error decoding params")
		return
	}

	if params.Rating < 1 || params.Rating > 5 {
		cfg.logger.Warn("add review", "user_id", userID, "reason", "rating invalid", "ip", getClientIP(r))
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
		if strings.Contains(err.Error(), "unique_user_juice_review") {
			cfg.logger.Warn("add review", "user_id", userID, "reason", "user already added a rating", "error", err, "ip", getClientIP(r))
			respondWithError(w, 409, "You have already reviewed this juice")
			return
		}
		cfg.logger.Error("add reviews", "user_id", userID, "error", err)
		respondWithError(w, 500, "Error creating review")
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
		cfg.logger.Warn("delete juice", "user_id", userID, "reason", "failed to parse review id", "error", err, "ip", getClientIP(r))
		respondWithError(w, 400, "invalid review ID")
		return
	}

	review, err := cfg.queries.GetReviewByID(r.Context(), parsedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Review not found")
			return
		}
		cfg.logger.Error("delete review", "user_id", userID, "error", err)
		respondWithError(w, http.StatusInternalServerError, "failed to fetch review")
		return
	}

	if review.UserID != userID && role != "admin" {
		cfg.logger.Warn("delete review", "user", userID, "reason", "invalid user", "ip", getClientIP(r))
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	err = cfg.queries.DeleteReview(r.Context(), parsedID)
	if err != nil {
		cfg.logger.Error("delete review", "user_id", userID, "error", err, "ip", getClientIP(r))
		respondWithError(w, http.StatusInternalServerError, "failed to fetch reviews")
		return
	}

	if role == "admin" && review.UserID != userID {
		cfg.queries.AddLog(r.Context(), database.AddLogParams{
			UserID:     uuid.NullUUID{UUID: userID, Valid: true},
			Action:     "delete",
			TargetType: "review",
			TargetID:   uuid.NullUUID{UUID: parsedID, Valid: true},
			TargetName: sql.NullString{String: review.UserID.String(), Valid: true},
		})
	}
	cfg.logger.Info("delete review", "user_id", userID, "ip", getClientIP(r))
	w.WriteHeader(204)
}
