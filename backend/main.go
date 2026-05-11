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
	Role         string    `json:"role"`
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
		log.Fatal("JWT_SECTER environment variable is not set")
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

	ServeMux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	ServeMux.HandleFunc("POST /api/login", cfg.handlerLogin)

	ServeMux.HandleFunc("POST /api/users", cfg.handlerCreateUser)

	ServeMux.HandleFunc("GET /api", cfg.handlerGetJuice)
	ServeMux.HandleFunc("GET /api/juice/{juiceID}", cfg.handlerGetChirpByID)
	ServeMux.HandleFunc("GET /api/admin/test", cfg.hanlderAdminTest)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: ServeMux,
	}

	log.Printf("Serving on: http://localhost:%s/api/\n", port)
	log.Fatal(server.ListenAndServe())
}
