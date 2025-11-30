package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/Racuwcka/service-courier/internal/http-server/handlers/courier/add"
)

type CourierHandlers struct {
	add *add.Handler
}

func New(add *add.Handler) *CourierHandlers {
	return &CourierHandlers{
		add: add,
	}
}

func (h *CourierHandlers) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, "hello")
	})

	r.Post("/", h.add.Create())

	return r
}
