package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/database"
	"github.com/joho/godotenv"
)

var greeting string

func main() {
	godotenv.Load()
	ctx := context.Background()

	config := database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
	}

	pool, err := database.NewPool(ctx, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create postgres pool: %v", err)
		os.Exit(1)
	}
	defer pool.Close()

	query := "SELECT 'Hello, World!'"
	err = pool.QueryRow(ctx, query).Scan(&greeting)

	if err != nil {
		fmt.Fprintf(os.Stderr, "query row failed: %v\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", greeting)
}
