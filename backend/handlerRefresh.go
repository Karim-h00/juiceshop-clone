package main

import (
	"net/http"
	"time"

	"github.com/karim-h00/juiceshop-clone/internal/auth"
)

func (cfg *config) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	user, err := cfg.queries.GetUserFromRefreshToken(r.Context(), refreshToken)
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
