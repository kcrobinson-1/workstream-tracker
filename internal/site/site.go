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
	"sort"

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
	roster := buildRoster(docs, active)

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
	if err := renderIndex(w, indexData{Roots: roots, PlansPath: s.plansPath, Roster: roster}); err != nil {
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

// RosterEntry is one active work-instance in the session roster,
// classified by plan-tree membership. p1 (bare roster) carries
// only the slug — the entry's display label, since p1 has no
// reported-name source — the actor (the identity key and the
// deterministic secondary sort key; never rendered as a label in
// p1, per scoping SD2), and the bound flag. p2 enriches the same
// loader output with the event-log-joined reported name and
// detail; it does not need to reshape this struct's identity.
type RosterEntry struct {
	Slug  string
	Actor string
	Bound bool
}

// buildRoster classifies every active work-instance as bound
// (its slug matches a walked plan doc) or unbound (its slug is
// absent from the walked tree), returning entries in
// deterministic (slug, actor) order.
//
// It is a pure function over the two values Server.index already
// loaded once per request — the walker's parsed-doc slice and
// the active-work-instance map — so the roster adds no second
// walk or query (the walk-on-every-request invariant; scoping
// SD1). The bound set is the parsed-doc slug set, the same slice
// buildTree joins against; classifying against it (not the
// post-buildTree tree) keeps a doc the tree later drops for a
// parent gap classified bound, matching the task-level "bound =
// slug matches a walked plan doc" contract. Unbound rows are
// listed, not dropped — bypassing buildTree's join-site drop.
func buildRoster(docs []parsedDoc, active map[string][]*ActiveWorkInstance) []RosterEntry {
	bound := make(map[string]bool, len(docs))
	for _, d := range docs {
		bound[d.Slug] = true
	}

	var entries []RosterEntry
	for slug, wis := range active {
		for _, wi := range wis {
			entries = append(entries, RosterEntry{
				Slug:  slug,
				Actor: wi.Actor,
				Bound: bound[slug],
			})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Slug != entries[j].Slug {
			return entries[i].Slug < entries[j].Slug
		}
		return entries[i].Actor < entries[j].Actor
	})
	return entries
}
