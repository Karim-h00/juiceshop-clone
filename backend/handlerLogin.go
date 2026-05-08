package main

import (
	"encoding/json"
	"net/http"

	"github.com/karim-h00/juiceshop-clone/internal/auth"
)

type Login_parameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *config) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := User_parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error decoding params")
		return
	}

	user_data, err := cfg.queries.GetPasswordByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Wrong email or password!")
		return
	}
	_, err = auth.CheckPasswordHash(params.Password, user_data.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "wrong email or password")
		return
	}

	respondWithJSON(w, 200, User{
		ID:        user_data.ID,
		Email:     user_data.Email,
		Username:  user_data.Username,
		CreatedAt: user_data.CreatedAt,
		UpdatedAt: user_data.UpdatedAt,
	})
}
