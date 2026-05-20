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
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"sort"
	"strings"

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

	meta, err := loadSessionMetadata(ctx, s.db, active)
	if err != nil {
		slog.Error("load session metadata", "err", err)
		http.Error(w, "failed to load session metadata", http.StatusInternalServerError)
		return
	}

	// Populate the forest's display fields on each active work-
	// instance. Both forest and roster read these from the same
	// per-request metadata resolution above; loadActiveWorkInstances
	// doesn't carry display data, so the assignment happens here.
	for slug, wis := range active {
		for _, wi := range wis {
			wi.Slug = slug
			wi.Name = meta[wi.ID].Name
		}
	}

	roots := buildTree(docs, active)
	roster := buildRoster(docs, active, meta)

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
		`SELECT id, slug, actor FROM work_instances WHERE state = ?`,
		models.StateActive,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[string][]*ActiveWorkInstance{}
	for rows.Next() {
		var id, slug, actor string
		if err := rows.Scan(&id, &slug, &actor); err != nil {
			return nil, err
		}
		out[slug] = append(out[slug], &ActiveWorkInstance{ID: id, Actor: actor})
	}
	return out, rows.Err()
}

// sessionMeta is one active work-instance's resolved reported
// metadata: the conventionally-read display Name (empty when the
// session reported no JSON-object `name`) and the
// deliberately-unstructured raw-JSON Detail (empty when the
// session reported no metadata at all — that entry renders as a
// plain row with no disclosure, not an empty <details>).
type sessionMeta struct {
	Name   string
	Detail string
}

// loadSessionMetadata resolves, per active work-instance, the
// reported metadata per the t4 task plan "Metadata read policy":
// the register event's metadata is the identity baseline; the
// latest later event's (heartbeat / state-transition) metadata is
// overlaid key-by-key (later keys win; keys absent from the later
// event keep the register value). It is a per-request read of a
// distinct source (the event log) scoped to exactly the active
// work-instance ids and keyed by the indexed
// events.work_instance_id — not a cache and not a second walk of
// already-walked data (the walk-on-every-request invariant). The
// baseline+overlay fold runs in Go; the query is only an indexed,
// active-id-scoped row read (scoping SD1/SD3).
func loadSessionMetadata(ctx context.Context, db *sql.DB, active map[string][]*ActiveWorkInstance) (map[string]sessionMeta, error) {
	var ids []string
	for _, wis := range active {
		for _, wi := range wis {
			ids = append(ids, wi.ID)
		}
	}
	if len(ids) == 0 {
		return map[string]sessionMeta{}, nil
	}

	placeholders := strings.Repeat("?,", len(ids)-1) + "?"
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}

	rows, err := db.QueryContext(ctx,
		`SELECT work_instance_id, type, metadata, received_at
		 FROM events
		 WHERE work_instance_id IN (`+placeholders+`)
		   AND metadata IS NOT NULL`,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	baseline := map[string]json.RawMessage{}
	latest := map[string]json.RawMessage{}
	latestAt := map[string]int64{}
	for rows.Next() {
		var wid, etype string
		var meta []byte
		var receivedAt int64
		if err := rows.Scan(&wid, &etype, &meta, &receivedAt); err != nil {
			return nil, err
		}
		if etype == string(models.EventRegister) {
			baseline[wid] = json.RawMessage(meta)
			continue
		}
		if receivedAt >= latestAt[wid] {
			latest[wid] = json.RawMessage(meta)
			latestAt[wid] = receivedAt
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	out := make(map[string]sessionMeta, len(ids))
	for _, id := range ids {
		if name, detail, ok := resolveSessionMeta(baseline[id], latest[id]); ok {
			out[id] = sessionMeta{Name: name, Detail: detail}
		}
	}
	return out, nil
}

// resolveSessionMeta folds a work-instance's register-baseline
// metadata and its latest-later-event metadata into the
// represented view. When either is a JSON object the result is
// the key-by-key overlay (later keys win; register-only keys
// survive) and `name`, if a JSON string in the result, is the
// display label. Non-object metadata has no defined key-by-key
// merge under the schema-loose posture: it is rendered verbatim
// (latest non-empty blob, else the register blob) with no name —
// not coerced into a schema (scoping SD3). ok is false only when
// the session reported no metadata at all.
func resolveSessionMeta(base, later json.RawMessage) (name, detail string, ok bool) {
	baseObj := asObject(base)
	laterObj := asObject(later)

	if baseObj != nil || laterObj != nil {
		merged := map[string]json.RawMessage{}
		for k, v := range baseObj {
			merged[k] = v
		}
		for k, v := range laterObj {
			merged[k] = v
		}
		if raw, ok2 := merged["name"]; ok2 {
			var s string
			if json.Unmarshal(raw, &s) == nil {
				name = s
			}
		}
		pretty, err := json.MarshalIndent(merged, "", "  ")
		if err != nil {
			return "", "", false
		}
		return name, string(pretty), true
	}

	// No object metadata anywhere: render whatever was reported
	// verbatim (schema-loose), preferring the later blob.
	raw := later
	if len(raw) == 0 {
		raw = base
	}
	if len(raw) == 0 {
		return "", "", false
	}
	var buf bytes.Buffer
	if json.Indent(&buf, raw, "", "  ") == nil {
		return "", buf.String(), true
	}
	return "", string(raw), true
}

// asObject returns raw parsed as a JSON object, or nil when raw
// is empty, JSON null, or not an object. A non-object blob
// "contributes no keys" to the key-by-key overlay (scoping SD3).
func asObject(raw json.RawMessage) map[string]json.RawMessage {
	if len(raw) == 0 {
		return nil
	}
	var obj map[string]json.RawMessage
	if json.Unmarshal(raw, &obj) != nil {
		return nil
	}
	return obj
}

// RosterEntry is one active work-instance in the session roster,
// classified by plan-tree membership.
//
// ID is work_instances.id (the generated PRIMARY KEY), the key
// the event log is joined on (events.work_instance_id ->
// work_instances.id). p1 deliberately did not carry it; p2's
// event-log join surfaces it per the recorded Cross-PR
// coordination handoff. It is identity/join data, never rendered.
//
// Actor is NOT a work-instance identity and is never rendered in
// the roster (scoping SD2): it is the deterministic secondary
// sort key so the walk-on-every-request page renders a stable
// (slug, actor) order when a slug carries several active
// work-instances.
//
// Name is the reported display label (the conventionally-read
// `name` metadata key) — empty when the session reported none,
// in which case the slug is the label (never the wst-<uuid>
// actor; the task-level name-then-slug rule). Detail is the
// deliberately-unstructured raw-JSON view of the session's
// resolved reported metadata — empty when the session reported
// nothing, in which case the entry renders as a plain row with
// no disclosure (enrichment is additive; absence is not a drop).
type RosterEntry struct {
	ID     string
	Slug   string
	Actor  string
	Bound  bool
	Name   string
	Detail string
}

// buildRoster classifies every active work-instance as bound
// (its slug matches a walked plan doc) or unbound (its slug is
// absent from the walked tree), returning entries in
// deterministic (slug, actor) order.
//
// It is a pure function over the three values Server.index
// already loaded once per request — the walker's parsed-doc
// slice, the active-work-instance map, and the per-request
// resolved session metadata (keyed by work_instances.id) — so
// the roster adds no second walk or query of its own (the
// walk-on-every-request invariant; scoping SD1). The bound set is
// the parsed-doc slug set, the same slice buildTree joins
// against; classifying against it (not the post-buildTree tree)
// keeps a doc the tree later drops for a parent gap classified
// bound, matching the task-level "bound = slug matches a walked
// plan doc" contract. Unbound rows are listed, not dropped —
// bypassing buildTree's join-site drop.
func buildRoster(docs []parsedDoc, active map[string][]*ActiveWorkInstance, meta map[string]sessionMeta) []RosterEntry {
	bound := make(map[string]bool, len(docs))
	for _, d := range docs {
		bound[d.Slug] = true
	}

	var entries []RosterEntry
	for slug, wis := range active {
		for _, wi := range wis {
			m := meta[wi.ID]
			entries = append(entries, RosterEntry{
				ID:     wi.ID,
				Slug:   slug,
				Actor:  wi.Actor,
				Bound:  bound[slug],
				Name:   m.Name,
				Detail: m.Detail,
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
