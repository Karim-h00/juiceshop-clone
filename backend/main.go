package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/karim-h00/juiceshop-clone/internal/database"
	_ "github.com/lib/pq"
)

type config struct {
	queries *database.Queries
	secret  string
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
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

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	cfg := config{
		queries: dbQueries,
		secret:  secret,
	}

	ServeMux := http.NewServeMux()

	ServeMux.HandleFunc("POST /api/login", cfg.handlerLogin)
	ServeMux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	ServeMux.Handle("PUT /api/users", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerUpdateUser)))
	ServeMux.Handle("POST /api/user/password", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerUpdatePassword)))

	ServeMux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	ServeMux.Handle("POST /api/logout", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerLogout)))

	ServeMux.HandleFunc("GET /api", cfg.handlerGetJuice)
	ServeMux.HandleFunc("GET /api/juice/{juiceName}", cfg.handlerGetJuiceByName)

	ServeMux.Handle("GET /api/admin/test", cfg.middlewareCheckAdmin(http.HandlerFunc(cfg.handlerAdminTest)))
	ServeMux.Handle("POST /api/admin/juice", cfg.middlewareCheckAdmin(http.HandlerFunc(cfg.handlerAddJuice)))
	ServeMux.Handle("PUT /api/admin/juice/{juiceID}", cfg.middlewareCheckAdmin(http.HandlerFunc(cfg.handlerUpdateJuice)))
	ServeMux.Handle("DELETE /api/admin/juice/{juiceID}", cfg.middlewareCheckAdmin(http.HandlerFunc(cfg.handlerDeleteJuice)))

	ServeMux.Handle("POST /api/order", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerOrderJuice)))
	ServeMux.Handle("GET /api/order", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerGetUserOrderHistory)))
	ServeMux.Handle("GET /api/order/{orderID}", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerGetOrderByID)))
	ServeMux.Handle("DELETE /api/order/{orderID}", cfg.middlewareCheckAdmin(http.HandlerFunc(cfg.handlerDeleteOrder)))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: ServeMux,
	}

	log.Printf("Serving on: http://localhost:%s/api/\n", port)
	log.Fatal(server.ListenAndServe())
}
