package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Monitor MonitorHandler // позже конкретизируем
}

func New(h Handlers) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		registerMonitors(r, h)
	})

	return r
}
