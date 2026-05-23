package main

import (
	"encoding/json"
	"net/http"

	"github.com/karim-h00/juiceshop-clone/internal/auth"
	"github.com/karim-h00/juiceshop-clone/internal/database"
)

type juice_params struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
}

func (cfg *config) handlerAddJuice(w http.ResponseWriter, r *http.Request) {
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

	decoder := json.NewDecoder(r.Body)
	params := juice_params{}

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
	juice, err := cfg.queries.AddJuice(r.Context(), database.AddJuiceParams{
		Name:        params.Name,
		Description: params.Description,
		Price:       int32(params.Price),
		Stock:       int32(params.Stock),
	})
	if err != nil {
		respondWithError(w, 500, "Error creating juice")
		return
	}
	respondWithJSON(w, 201, juice)
}
