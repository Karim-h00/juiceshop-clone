package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/karim-h00/juiceshop-clone/internal/auth"
	"github.com/karim-h00/juiceshop-clone/internal/database"
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
	const defaultExpiry = time.Hour

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
	refreshTokenStr, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token")
		return
	}

	_, err = cfg.queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshTokenStr,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user_data.ID,
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 60),
	})
	if err != nil {
		respondWithError(w, 500, "Could not create refresh token")
		return
	}

	user_token, err := auth.MakeJWT(user_data.ID, cfg.secret, defaultExpiry, user_data.Role)
	if err != nil {
		respondWithError(w, 500, "Could not create session")
		return
	}

	respondWithJSON(w, 200, User{
		ID:           user_data.ID,
		Email:        user_data.Email,
		Username:     user_data.Username,
		CreatedAt:    user_data.CreatedAt,
		UpdatedAt:    user_data.UpdatedAt,
		Token:        user_token,
		RefreshToken: refreshTokenStr,
	})
}
