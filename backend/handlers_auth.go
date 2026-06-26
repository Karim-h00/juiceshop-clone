package main

import (
	"database/sql"
	"encoding/json"
	"errors"
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
		cfg.logger.Warn("login failed", "reason", "invalid request body", "ip", getClientIP(r))
		respondWithError(w, 400, "Error decoding params")
		return
	}
	const defaultExpiry = time.Hour

	user_data, err := cfg.queries.GetPasswordByEmail(r.Context(), params.Email)
	if err != nil {
		cfg.logger.Warn("login failed", "reason", "invalid email", "ip", getClientIP(r))
		respondWithError(w, 401, "Wrong email or password")
		return
	}
	_, err = auth.CheckPasswordHash(params.Password, user_data.HashedPassword)
	if err != nil {
		cfg.logger.Warn("login failed", "reason", "invalid password", "ip", getClientIP(r))
		respondWithError(w, 401, "Wrong email or password")
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
			cfg.logger.Error("failed to delete existing refresh token", "user_id", user_data.ID, "err", err)
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
			cfg.logger.Error("failed to create refresh token", "user_id", user_data.ID, "err", err)
			respondWithError(w, 500, "Could not create refresh token")
			return
		}
	}

	now := time.Now().UTC()
	user_token, err := auth.MakeJWT(user_data.ID, cfg.secret, defaultExpiry, user_data.Role)
	if err != nil {
		cfg.logger.Error("failed to create jwt", "user_id", user_data.ID, "err", err)
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

	if user_data.Role == "admin" {
		err = cfg.queries.AddLog(r.Context(), database.AddLogParams{
			UserID:     uuid.NullUUID{UUID: user_data.ID, Valid: true},
			Action:     "admin_login",
			TargetType: "login",
			TargetID:   uuid.NullUUID{UUID: user_data.ID, Valid: true},
			TargetName: sql.NullString{String: user_data.Username, Valid: true},
			CreatedAt:  now,
		})
		if err != nil {
			cfg.logger.Error("add audit log", "error", err, "ip", getClientIP(r))
		}
	}

	cfg.logger.Info("login success", "user_id", user_data.ID, "ip", getClientIP(r))

	respondWithJSON(w, 200, struct {
		ID       uuid.UUID `json:"id"`
		Email    string    `json:"email"`
		Username string    `json:"username"`
		Token    string    `json:"token"`
	}{
		ID:       user_data.ID,
		Email:    user_data.Email,
		Username: user_data.Username,
		Token:    user_token,
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
		cfg.logger.Error("logout failed", "reason", "failed to revoke refresh token", "err", err)
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
		cfg.logger.Warn("refresh failed", "reason", "invalid or expired refresh token", "ip", getClientIP(r))
		respondWithError(w, 401, "Invalid or expired refresh token")
		return
	}

	newToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour, user.Role)
	if err != nil {
		cfg.logger.Error("failed to create jwt", "user_id", user.ID, "err", err)
		respondWithError(w, 500, "Could not create token")
		return
	}

	cfg.logger.Info("refresh success", "user id", user.ID, "ip", getClientIP(r))
	respondWithJSON(w, 200, struct {
		Token string `json:"token"`
	}{Token: newToken})
}

func (cfg *config) handlerMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(contextKeyUserID).(uuid.UUID)

	user, err := cfg.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			cfg.logger.Warn("get me failed", "reason", "user not found", "user id", userID, "ip", getClientIP(r))
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		cfg.logger.Error("get me failed", "reason", "server error", "user id", userID, "ip", getClientIP(r))
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondWithJSON(w, 200, struct {
		ID       uuid.UUID `json:"id"`
		Email    string    `json:"email"`
		Username string    `json:"username"`
		Role     string    `json:"role"`
	}{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	})
}
