package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/repositories"
)

type PersonRequestDTO struct {
	Surname   string    `json:"apelido"`
	Name      string    `json:"nome"`
	Birthdate time.Time `json:"nascimento"`
}

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
	var personRequest PersonRequestDTO

	// TODO
	// MARSHAL
	err := json.NewDecoder(r.Body).Decode(&personRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	ctx := context.WithValue(r.Context(), PersonRequest, personRequest)
	r.WithContext(ctx)

	fmt.Fprintf(w, "Creating person...\n%s", personRequest.Surname)
}

func (p *PeopleHandler) FindById(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "ID not supplied", http.StatusBadRequest)
		return
	}
	id := pathParts[2]

	fmt.Fprintf(w, "Get person id: %v", id)
}

func (p *PeopleHandler) Search(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("t")

	fmt.Fprintf(w, "Termo de busca: %s", term)
}
