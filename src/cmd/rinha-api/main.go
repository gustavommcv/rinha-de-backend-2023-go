package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/database"
	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/handlers"
	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/repositories"
	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/routes"
	"github.com/joho/godotenv"
)

const (
	DEFAULT_PORT = ":8080"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	ctx := context.Background()

	pool, err := database.NewPool(ctx, generateConfig())
	checkError(err, "failed to create postgres pool")
	defer pool.Close()

	personRepository := repositories.NewUserRepository(pool)
	peopleHandler := handlers.NewPeopleHandler(*personRepository)
	indexRouter := routes.NewIndexRouter(*peopleHandler)

	fmt.Println("Running at ", DEFAULT_PORT)
	log.Fatal(http.ListenAndServe(DEFAULT_PORT, indexRouter))
}

func checkError(err error, message string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error: %v\n", message, err)
		os.Exit(1)
	}
}

func generateConfig() database.Config {
	config := database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
	}

	return config
}
