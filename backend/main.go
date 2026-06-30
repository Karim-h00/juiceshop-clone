package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/karim-h00/juiceshop-clone/internal/database"
	"github.com/karim-h00/juiceshop-clone/internal/ratelimit"
	_ "github.com/lib/pq"
	"golang.org/x/time/rate"
)

type config struct {
	db            *sql.DB
	queries       *database.Queries
	secret        string
	assetsRoot    string
	port          string
	baseURL       string
	logger        *slog.Logger
	globalLimiter *ratelimit.Limiter
	authLimiter   *ratelimit.Limiter
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("dbURL must be set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	assetsRoot := os.Getenv("ASSETS_ROOT")
	if assetsRoot == "" {
		log.Fatal("ASSETS_ROOT environment variable is not set")
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Fatal("BASE_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	dbQueries := database.New(db)
	cfg := config{
		db:            db,
		queries:       dbQueries,
		secret:        secret,
		assetsRoot:    assetsRoot,
		port:          port,
		baseURL:       baseURL,
		logger:        logger,
		globalLimiter: ratelimit.New(rate.Every(time.Second), 20, 30*time.Minute),
		authLimiter:   ratelimit.New(rate.Every(10*time.Second), 5, 30*time.Minute),
	}

	ServeMux := http.NewServeMux()

	assetsHandler := http.StripPrefix("/assets", http.FileServer(http.Dir(assetsRoot)))
	ServeMux.Handle("/assets/", assetsHandler)

	ServeMux.HandleFunc("GET /api/healthz", handlerReadiness)
	ServeMux.Handle("POST /api/login", cfg.RateLimitMiddleware(cfg.authLimiter)(http.HandlerFunc(cfg.handlerLogin)))
	ServeMux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	ServeMux.Handle("PUT /api/users", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerUpdateUser)))
	ServeMux.Handle("POST /api/user/password", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerUpdatePassword)))
	ServeMux.HandleFunc("POST /api/check-password", cfg.handlerCheckPwn)
	ServeMux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	ServeMux.HandleFunc("POST /api/logout", cfg.handlerLogout)
	ServeMux.Handle("GET /api/me", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerMe)))

	ServeMux.HandleFunc("GET /api", cfg.handlerGetJuice)
	ServeMux.HandleFunc("GET /api/juice/{slug}", cfg.handlerGetJuiceByName)

	ServeMux.Handle("POST /api/order", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerOrderJuice)))
	ServeMux.Handle("GET /api/order", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerGetUserOrderHistory)))
	ServeMux.Handle("GET /api/order/{orderID}", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerGetOrderByID)))

	ServeMux.HandleFunc("GET /api/juice/{slug}/reviews", cfg.handlerGetReviews)
	ServeMux.Handle("POST /api/juice/{slug}/reviews", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerAddReview)))
	ServeMux.Handle("DELETE /api/review/{reviewID}", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerDeleteReview)))

	ServeMux.Handle("POST /api/admin/juice", cfg.middlewareCheckAdmin(http.HandlerFunc(cfg.handlerAddJuice)))
	ServeMux.Handle("PUT /api/admin/juice/{juiceID}", cfg.middlewareCheckAdmin(http.HandlerFunc(cfg.handlerUpdateJuice)))
	ServeMux.Handle("PUT /api/admin/juice/{juiceID}/image", cfg.middlewareCheckAdmin(http.HandlerFunc(cfg.handlerUpdateJuiceImage)))
	ServeMux.Handle("DELETE /api/admin/juice/{juiceID}", cfg.middlewareCheckAdmin(http.HandlerFunc(cfg.handlerDeleteJuice)))

	ServeMux.Handle("GET /api/admin/users", cfg.middlewareCheckAdmin(http.HandlerFunc(cfg.handlerGetAllUsers)))
	ServeMux.Handle("PATCH /api/admin/users/{userID}/role", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerAdminUpdate)))
	ServeMux.Handle("DELETE /api/admin/users/{userID}", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerDeleteUser)))

	ServeMux.Handle("DELETE /api/admin/order/{orderID}", cfg.middlewareCheckAdmin(http.HandlerFunc(cfg.handlerDeleteOrder)))
	ServeMux.Handle("GET /api/admin/orders", cfg.middlewareCheckAdmin(http.HandlerFunc(cfg.handlerAdminGetAllOrders)))

	ServeMux.Handle("GET /api/admin/audits", cfg.middlewareCheckAdmin(http.HandlerFunc(cfg.handlerGetAuditLogs)))

	handler := middlewareCORS(cfg.RateLimitMiddleware(cfg.globalLimiter)(ServeMux))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	logger.Info("server starting", "addr", "http://localhost:"+port+"/api/")
	log.Fatal(server.ListenAndServe())
}
