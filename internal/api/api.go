// Package api implements the workstream-tracker HTTP API. See
// design/v0.1-design.md Section 4 for the endpoint shapes.
package api

import (
	"database/sql"
	"encoding/json"

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
func (s *Server) MountRoutes(r chi.Router) {
	r.Post("/", s.registerWorkInstance)
	r.Post("/{id}/events", s.recordEvent)
}

// RegisterRequest is the body of POST /work-instances.
//
// ParentPath distinguishes the two flows:
//   - nil (field omitted from JSON): root case. RootSlug names a
//     new root; node_type must be "epic" or "task".
//   - non-nil (field present, possibly empty string): descendant
//     case. ParentPath is the slug-suffix of the parent (empty
//     when the parent is the root). The server generates the
//     descendant slug.
type RegisterRequest struct {
	RootSlug   string          `json:"root_slug"`
	ParentPath *string         `json:"parent_path,omitempty"`
	NodeType   string          `json:"node_type"`
	Actor      string          `json:"actor"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
}

// RegisterResponse is the body returned by a successful POST
// /work-instances.
type RegisterResponse struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
}

// EventRequest is the body of POST /work-instances/{id}/events.
//
// State is empty for heartbeats and one of "completed" or
// "abandoned" for state transitions.
type EventRequest struct {
	State    string          `json:"state,omitempty"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
}

// EventResponse is the body returned by a successful POST to the
// events endpoint.
type EventResponse struct {
	ID string `json:"id"`
}
