package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/karim-h00/juiceshop-clone/internal/auth"
	"github.com/karim-h00/juiceshop-clone/internal/database"
)

func (cfg *config) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type Login_parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := Login_parameters{}
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

	var refreshTokenStr string
	existingToken, err := cfg.queries.GetRefreshTokenByUserID(r.Context(), user_data.ID)
	if err == nil && !existingToken.RevokedAt.Valid && time.Now().UTC().Before(existingToken.ExpiresAt) {
		refreshTokenStr = existingToken.Token
	} else {
		cfg.queries.DeleteRefreshTokenByUserID(r.Context(), user_data.ID)
		refreshTokenStr, err = auth.MakeRefreshToken()
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
	}

	user_token, err := auth.MakeJWT(user_data.ID, cfg.secret, defaultExpiry, user_data.Role)
	if err != nil {
		respondWithError(w, 500, "Could not create session")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenStr,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   60 * 24 * 60 * 60,
	})

	respondWithJSON(w, 200, User{
		ID:        user_data.ID,
		Email:     user_data.Email,
		Username:  user_data.Username,
		CreatedAt: user_data.CreatedAt,
		UpdatedAt: user_data.UpdatedAt,
		Token:     user_token,
	})
}

func (cfg *config) handlerLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.queries.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		Token:     cookie.Value,
		UpdatedAt: time.Now().UTC(),
		RevokedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke refresh token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *config) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	user, err := cfg.queries.GetUserFromRefreshToken(r.Context(), cookie.Value)
	if err != nil {
		respondWithError(w, 401, "Invalid or expired refresh token")
		return
	}

	newToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour, user.Role)
	if err != nil {
		respondWithError(w, 500, "Could not create token")
		return
	}

	respondWithJSON(w, 200, struct {
		Token string `json:"token"`
	}{Token: newToken})
}

func (cfg *config) handlerMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(contextKeyUserID).(uuid.UUID)

	user, err := cfg.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, 404, "User not found")
		return
	}

	respondWithJSON(w, 200, User{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}
