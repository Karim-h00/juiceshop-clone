package main

import (
	"encoding/json"
	"net/http"

	"github.com/karim-h00/juiceshop-clone/internal/auth"
	"github.com/karim-h00/juiceshop-clone/internal/database"
)

type User_parameters struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *config) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := User_parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error decoding params")
		return
	}
	if params.Password == "" {
		respondWithError(w, 400, "password is required")
	}
	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Error creating user")
		return
	}

	user, err := cfg.queries.CreateUser(r.Context(), database.CreateUserParams{
		Username:       params.Username,
		Email:          params.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		respondWithError(w, 500, "Error creating user")
		return
	}
	respondWithJSON(w, 201, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Username:  user.Username,
	})
}
