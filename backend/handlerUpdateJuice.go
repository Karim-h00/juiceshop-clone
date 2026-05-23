package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/karim-h00/juiceshop-clone/internal/auth"
	"github.com/karim-h00/juiceshop-clone/internal/database"
)

type update_juice_params struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
}

func (cfg *config) handlerUpdateJuice(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	userID, tokenRole, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 400, "Could not make session")
		return
	}
	role, err := cfg.queries.GetUserRole(r.Context(), userID)
	if err != nil {
		respondWithError(w, 500, "Error getting user role")
		return
	}
	if role != tokenRole {
		respondWithError(w, 403, "mismatched token")
		return
	}
	if role != "admin" {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	juiceID := r.PathValue("juiceID")
	parsedID, err := uuid.Parse(juiceID)
	if err != nil {
		respondWithError(w, 400, "failed to parse ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := update_juice_params{}

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error decoding params")
		return
	}

	if params.Name == "" {
		respondWithError(w, 400, "Name is required")
		return
	}
	if params.Price <= 0 {
		respondWithError(w, 400, "price must be positive")
		return
	}
	if params.Stock < 0 {
		respondWithError(w, 400, "stock must not be negative")
		return
	}
	juice, err := cfg.queries.UpdateJuice(r.Context(), database.UpdateJuiceParams{
		ID:          parsedID,
		Name:        params.Name,
		Description: params.Description,
		Price:       int32(params.Price),
		Stock:       int32(params.Stock),
	})
	if err != nil {
		respondWithError(w, 500, "Error updating juice")
		return
	}
	respondWithJSON(w, 200, juice)
}
