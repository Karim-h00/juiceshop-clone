package main

import (
	"database/sql"
	"fmt"
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
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	cfg := config{
		queries: dbQueries,
	}

	ServeMux := http.NewServeMux()
	ServeMux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	ServeMux.HandleFunc("POST /api/users", cfg.handlerCreateUser)

	server := &http.Server{
		Addr:    ":8080",
		Handler: ServeMux,
	}
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("something went wrong: %v", err)
	}
}
