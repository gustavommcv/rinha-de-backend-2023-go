package routes

import (
	"net/http"

	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/handlers"
)

func NewPeopleRouter(handler handlers.PeopleHandler) *http.ServeMux {
	peopleRouter := http.NewServeMux()

	peopleRouter.HandleFunc("GET /contagem-pessoas", handler.Count)
	peopleRouter.HandleFunc("POST /pessoas", handler.Create)
	peopleRouter.HandleFunc("GET /pessoas", handler.Search)
	peopleRouter.HandleFunc("GET /pessoas/{id}", handler.FindById)

	return peopleRouter
}
