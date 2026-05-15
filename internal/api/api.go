// Package api implements the workstream-tracker HTTP API. See
// design/v0.1-design.md Section 4 for the endpoint shapes.
package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MountRoutes registers the workstream-tracker API endpoints on r.
//
// v0.0 stubs return 501 Not Implemented; the actual API logic
// (slug validation, event log append, current-state fold) lands
// in v0.1.
func MountRoutes(r chi.Router) {
	r.Post("/", registerWorkInstance)
	r.Post("/{id}/events", recordEvent)
}

func registerWorkInstance(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "POST /work-instances: not implemented in v0.0", http.StatusNotImplemented)
}

func recordEvent(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "POST /work-instances/{id}/events: not implemented in v0.0", http.StatusNotImplemented)
}
