package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/entities"
	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/repositories"
)

const PersonRequest = "person-request"

type PeopleHandler struct {
	repository repositories.PersonRepository
}

func NewPeopleHandler(repository repositories.PersonRepository) *PeopleHandler {
	return &PeopleHandler{
		repository: repository,
	}
}

func (p *PeopleHandler) Count(w http.ResponseWriter, _ *http.Request) {
	count, err := p.repository.GetPeopleCount(context.Background())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error counting users: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", count)
}

func (p *PeopleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var personRequest entities.PersonRequestDTO

	err := json.NewDecoder(r.Body).Decode(&personRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validations
	// Surname
	if personRequest.Surname == "" {
		http.Error(w, "Surname required", http.StatusUnprocessableEntity)
		return
	}
	if len(personRequest.Surname) > 32 {
		http.Error(w, "Surname max 32 characters", http.StatusUnprocessableEntity)
		return
	}
	// TODO
	// Tem que ser unique

	// Name
	if personRequest.Name == "" {
		http.Error(w, "Name required", http.StatusUnprocessableEntity)
		return
	}
	if len(personRequest.Name) > 100 {
		http.Error(w, "Name max 100 characters", http.StatusUnprocessableEntity)
		return
	}

	//Birthdate
	_, err = time.Parse("2006-01-02", personRequest.Birthdate)
	if err != nil {
		http.Error(w, "Birthdate format invalid - AAAA-MM-DD", http.StatusUnprocessableEntity)
		return
	}

	// Stack
	stackErr := slices.Contains(personRequest.Stack, "")
	if stackErr == true {
		http.Error(w, "A language cant be null", http.StatusUnprocessableEntity)
		return
	}

	ctx := context.WithValue(r.Context(), PersonRequest, personRequest)
	r.WithContext(ctx)

	personResponse, err := p.repository.CreatePerson(context.Background(), personRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/pessoas/%s", personResponse.Id))
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Creating person...\n%s", *personResponse)
}

func (p *PeopleHandler) FindById(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "ID not supplied", http.StatusBadRequest)
		return
	}
	id := pathParts[2]

	if err := uuid.Validate(id); err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	response, err := p.repository.FindById(context.Background(), uuid.MustParse(id))
	if err != nil {
		http.Error(w, "Get id error", http.StatusBadRequest)
		return
	}

	if response.Id == "" {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	user, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%v", string(user))
}

func (p *PeopleHandler) Search(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("t")

	if searchTerm == "" {
		http.Error(w, "Term not provided", http.StatusBadRequest)
		return
	}

	response, err := p.repository.Search(context.Background(), searchTerm)
	if err != nil {
		http.Error(w, "Search error", http.StatusBadRequest)
		return
	}

	people, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", string(people))
}
