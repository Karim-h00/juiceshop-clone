package main

import (
	"encoding/json"
	"net/http"

	"github.com/karim-h00/juiceshop-clone/internal/auth"
	"github.com/karim-h00/juiceshop-clone/internal/database"
)

type user_update_params struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (cfg *config) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	userID, _, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 400, "Could not make session")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := user_update_params{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error decoding params")
		return
	}
	user, err := cfg.queries.UpdateUser(r.Context(), database.UpdateUserParams{
		Username: params.Username,
		Email:    params.Email,
		ID:       userID,
	})
	if err != nil {
		respondWithError(w, 500, "Error creating user")
		return
	}
	respondWithJSON(w, 200, user)
}
