package router

import (
	"github.com/go-chi/chi/v5"
)

func registerMonitors(r chi.Router, h Handlers) {

	r.Route("/monitors", func(r chi.Router) {
		r.Post("/", h.Monitor.Create)
		r.Get("/", h.Monitor.List)
		r.Get("/{id}", h.Monitor.GetByID)
		r.Put("/{id}", h.Monitor.Update)
		r.Delete("/{id}", h.Monitor.Delete)
	})
}
