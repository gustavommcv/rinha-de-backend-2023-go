package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/repositories"
)

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

func (p *PeopleHandler) Create(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Creating person...")
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
