package routes

import (
	"net/http"

	"github.com/gustavommcv/rinha-de-backend-2023-go/src/internal/handlers"
)

func NewIndexRouter(handler handlers.PeopleHandler) *http.ServeMux {
	peopleRouter := NewPeopleRouter(handler)
	indexRouter := http.NewServeMux()

	indexRouter.Handle("/", peopleRouter)

	return indexRouter
}
