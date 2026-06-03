package main

import (
	"context"
	"net/http"

	"github.com/karim-h00/juiceshop-clone/internal/auth"
)

func (cfg *config) middlewareCheckAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		next.ServeHTTP(w, r)
	})
}

type contextKey string

const (
	contextKeyUserID contextKey = "userID"
	contextKeyRole   contextKey = "role"
	contextKeyToken  contextKey = "token"
)

func (cfg *config) middlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx := context.WithValue(r.Context(), contextKeyUserID, userID)
		ctx = context.WithValue(ctx, contextKeyRole, tokenRole)
		ctx = context.WithValue(ctx, contextKeyToken, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
