package httpapi

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/your-user/go-service-template/internal/task"
)

// Pinger is the small database dependency needed by the health endpoint.
// Keeping it as an interface makes the router straightforward to test.
type Pinger interface {
	PingContext(context.Context) error
}

func NewRouter(tasks *task.Service, database Pinger) http.Handler {
	r := chi.NewRouter()
	r.Use(recoverer, requestLogger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := database.PingContext(r.Context()); err != nil {
			writeError(w, http.StatusServiceUnavailable, "database unavailable")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "running"})
	})

	h := task.NewHandler(tasks)
	r.Route("/api/v1/tasks", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
	return r
}
