package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Monitor MonitorHandler
}

func New(h Handlers) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		registerMonitors(r, h)
	})

	return r
}
