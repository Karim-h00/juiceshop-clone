package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/karim-h00/juiceshop-clone/internal/database"
)

type config struct {
	queries *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)

	dbQueries := database.New(db)
	cfg := config{
		queries: dbQueries,
	}

	ServeMux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: ServeMux,
	}
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("something went wrong: %v", err)
	}
}
