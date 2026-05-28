package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/karim-h00/juiceshop-clone/internal/database"
)

func (cfg *config) handlerLogout(w http.ResponseWriter, r *http.Request) {

	token := r.Context().Value(contextKeyUserID).(string)

	err := cfg.queries.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		Token:     token,
		UpdatedAt: time.Now().UTC(),
		RevokedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke refresh token")
		return
	}
	respondWithJSON(w, http.StatusOK, "Logged out")
}
