package main

import (
	"context"
	"net/http"

	"github.com/karim-h00/juiceshop-clone/internal/auth"
	"github.com/karim-h00/juiceshop-clone/internal/ratelimit"
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
			respondWithError(w, 401, "Unauthorized")
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
			respondWithError(w, 401, "Unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUserID, userID)
		ctx = context.WithValue(ctx, contextKeyRole, tokenRole)
		ctx = context.WithValue(ctx, contextKeyToken, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func middlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (cfg *config) RateLimitMiddleware(limiter *ratelimit.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)

			if !limiter.Allow(ip) {
				cfg.logger.Warn("rate limit exceeded", "ip", ip, "path", r.URL.Path)
				respondWithError(w, http.StatusTooManyRequests, "too many requests")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
