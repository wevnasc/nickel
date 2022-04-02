package routers

import (
	"nickel/http/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func ListenEntityRouters(r *chi.Mux, h *handlers.EntryHandlers) {
	r.Use(middleware.Logger)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Get("/", h.List())
	r.Post("/", h.Create())
}
