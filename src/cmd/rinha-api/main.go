package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/database"
	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/repositories"
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
	checkError(err, "failed to create postgres pool")
	defer pool.Close()

	personRepository := repositories.NewUserRepository(pool)

	http.HandleFunc("GET /contagem-pessoas", func(w http.ResponseWriter, r *http.Request) {
		countHandler(w, r, personRepository)
	})
	http.HandleFunc("POST /pessoas", createPersonHandler)
	http.HandleFunc("GET /pessoas/{id}", getPersonHandler)
	http.HandleFunc("GET /pessoas", getPeopleHandler)

	fmt.Println("Running at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func countHandler(w http.ResponseWriter, _ *http.Request, repository *repositories.PersonRepository) {
	count, err := repository.GetPeopleCount(context.Background())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error counting users: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(269)
	fmt.Fprintf(w, "%d", count)
}

func createPersonHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Creating person...")
}

func getPersonHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "ID not supplied", http.StatusBadRequest)
		return
	}
	id := pathParts[2]

	fmt.Fprintf(w, "Get person id: %v", id)
}

func getPeopleHandler(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("t")

	fmt.Fprintf(w, "Termo de busca: %s", term)
}

func checkError(err error, message string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error: %v\n", message, err)
		os.Exit(1)
	}
}
