// Package api implements the workstream-tracker HTTP API. See
// design/v0.1-design.md Section 4 for the endpoint shapes.
package api

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Server holds the dependencies the API handlers need (currently
// just the database). Constructed once at startup and used to
// register routes via MountRoutes.
type Server struct {
	db *sql.DB
}

// New constructs a Server backed by the given database.
func New(db *sql.DB) *Server {
	return &Server{db: db}
}

// MountRoutes registers the workstream-tracker API endpoints on r.
//
// v0.0 stubs return 501 Not Implemented; the actual API logic
// (slug validation, event log append, current-state fold) lands
// in the next milestone.
func (s *Server) MountRoutes(r chi.Router) {
	r.Post("/", s.registerWorkInstance)
	r.Post("/{id}/events", s.recordEvent)
}

func (s *Server) registerWorkInstance(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "POST /work-instances: not implemented in v0.0", http.StatusNotImplemented)
}

func (s *Server) recordEvent(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "POST /work-instances/{id}/events: not implemented in v0.0", http.StatusNotImplemented)
}
