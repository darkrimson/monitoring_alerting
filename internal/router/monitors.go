package router

import (
	"net/http"

	"github.com/darkrimson/monitoring_alerting/internal/handler"
	"github.com/go-chi/chi/v5"
)

type MonitorHandler interface {
	Create(http.ResponseWriter, *http.Request)
	List(http.ResponseWriter, *http.Request)
	GetByID(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

func registerMonitors(r chi.Router, h Handlers) {
	mh := h.Monitor.(*handler.MonitorHandler)

	r.Route("/monitors", func(r chi.Router) {
		r.Post("/", mh.Create)
		r.Get("/", mh.List)
		r.Get("/{id}", mh.GetByID)
		r.Put("/{id}", mh.Update)
		r.Delete("/{id}", mh.Delete)
	})
}
