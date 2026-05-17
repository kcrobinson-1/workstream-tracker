// Package site renders the workstream-tracker website. See
// design/v0.1-design.md Section 7 for the visualization scope.
//
// Each request walks the plan-tree directory (plansPath) under
// the per-root-folder layout convention in
// spec/planning-doc-location.md, parses the YAML frontmatter of
// each markdown file (slug + Status), joins the result with
// active work-instance state from the DB, and renders the forest
// as HTML.
package site

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kcrobinson-1/workstream-tracker/internal/models"
)

// Server holds the dependencies the site handlers need.
type Server struct {
	db        *sql.DB
	plansPath string
}

// New constructs a Server backed by the given database and
// reading plan-tree docs from plansPath.
func New(db *sql.DB, plansPath string) *Server {
	return &Server{db: db, plansPath: plansPath}
}

// Router returns a chi.Router serving the workstream-tracker
// website.
func (s *Server) Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/", s.index)
	return r
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	docs, err := walkPlans(s.plansPath)
	if err != nil {
		slog.Error("walk plans", "path", s.plansPath, "err", err)
		http.Error(w, "failed to read plan tree", http.StatusInternalServerError)
		return
	}

	active, err := loadActiveWorkInstances(ctx, s.db)
	if err != nil {
		slog.Error("load work-instances", "err", err)
		http.Error(w, "failed to load work-instances", http.StatusInternalServerError)
		return
	}

	roots := buildTree(docs, active)

	// Best-effort gh pr list auto-discovery (t4 P2). A discovery
	// failure is non-fatal: log once at warn and render the
	// unaugmented tree (frontmatter PRs only). No caching /
	// goroutine / file-watch — walk-on-every-request invariant.
	if prs, err := discoverPRsByTitle(ctx); err != nil {
		slog.Warn("gh pr discovery", "err", err)
	} else {
		augmentRelatedPRs(roots, prs)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := renderIndex(w, indexData{Roots: roots, PlansPath: s.plansPath}); err != nil {
		slog.Error("render index", "err", err)
	}
}

// loadActiveWorkInstances returns a slug → active-work-instances
// map, keyed by the slug each work-instance is attached to. Only
// work-instances in state 'active' are included.
func loadActiveWorkInstances(ctx context.Context, db *sql.DB) (map[string][]*ActiveWorkInstance, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT slug, actor FROM work_instances WHERE state = ?`,
		models.StateActive,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[string][]*ActiveWorkInstance{}
	for rows.Next() {
		var slug, actor string
		if err := rows.Scan(&slug, &actor); err != nil {
			return nil, err
		}
		out[slug] = append(out[slug], &ActiveWorkInstance{Actor: actor})
	}
	return out, rows.Err()
}
