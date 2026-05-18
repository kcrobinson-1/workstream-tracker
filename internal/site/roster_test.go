package site

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// Roster-region tests. p1 covered the bare bound/unbound roster
// (classification, deterministic order, slug label, empty state,
// no-affordance guard). p2 adds: reported-name label with slug
// fallback (never the wst-<uuid> actor), the register-baseline /
// latest-later key-by-key overlay (resolveSessionMeta), and the
// expandable raw-JSON detail present for a metadata-bearing entry
// / absent (plain row) for a no-metadata entry. t4's owned test
// surface.

// renderRoster renders a forest + roster and returns the HTML.
// Separate from render_test.go's renderTree helper so the shared
// helper's signature (which the forest/shell tests depend on)
// stays unchanged.
func renderRoster(t *testing.T, roots []*PlanNode, roster []RosterEntry) string {
	t.Helper()
	var buf bytes.Buffer
	if err := renderIndex(&buf, indexData{Roots: roots, PlansPath: "docs/plans", Roster: roster}); err != nil {
		t.Fatalf("renderIndex: %v", err)
	}
	return buf.String()
}

func TestBuildRosterClassifiesBoundAndUnbound(t *testing.T) {
	docs := []parsedDoc{
		{Slug: "alpha", Status: "Proposed"},
		{Slug: "alpha-m1", Status: "In progress"},
	}
	active := map[string][]*ActiveWorkInstance{
		"alpha":           {{ID: "w1", Actor: "wst-aaaa"}},
		"alpha-m1":        {{ID: "w2", Actor: "wst-bbbb"}},
		"typoed-not-real": {{ID: "w3", Actor: "wst-cccc"}},
	}
	entries := buildRoster(docs, active, nil)

	got := map[string]bool{}
	for _, e := range entries {
		got[e.Slug] = e.Bound
	}
	if len(entries) != 3 {
		t.Fatalf("got %d entries, want 3: %+v", len(entries), entries)
	}
	if !got["alpha"] || !got["alpha-m1"] {
		t.Errorf("alpha / alpha-m1 should be bound: %+v", entries)
	}
	if got["typoed-not-real"] {
		t.Errorf("typoed-not-real should be unbound: %+v", entries)
	}
}

// TestBuildRosterParentGappedDocStillBound pins scoping SD1 (p1):
// the bound set is the parsed-doc slug set, NOT the post-buildTree
// tree.
func TestBuildRosterParentGappedDocStillBound(t *testing.T) {
	docs := []parsedDoc{
		{Slug: "alpha", Status: "Proposed"},
		{Slug: "alpha-m1-t1", Status: "In draft"}, // parent alpha-m1 absent
	}
	roots := buildTree(docs, nil)
	var walk func(n *PlanNode) bool
	walk = func(n *PlanNode) bool {
		if n.Slug == "alpha-m1-t1" {
			return true
		}
		for _, c := range n.Children {
			if walk(c) {
				return true
			}
		}
		return false
	}
	for _, r := range roots {
		if walk(r) {
			t.Fatalf("premise broken: buildTree did NOT drop the parent-gapped node")
		}
	}

	entries := buildRoster(docs, map[string][]*ActiveWorkInstance{
		"alpha-m1-t1": {{ID: "w9", Actor: "wst-dddd"}},
	}, nil)
	if len(entries) != 1 || !entries[0].Bound {
		t.Errorf("parent-gapped-but-walked doc must be bound: %+v", entries)
	}
}

func TestBuildRosterDeterministicOrder(t *testing.T) {
	docs := []parsedDoc{{Slug: "alpha"}}
	active := map[string][]*ActiveWorkInstance{
		"beta":  {{ID: "b2", Actor: "wst-2"}, {ID: "b1", Actor: "wst-1"}},
		"alpha": {{ID: "a9", Actor: "wst-9"}},
	}
	want := []RosterEntry{
		{ID: "a9", Slug: "alpha", Actor: "wst-9", Bound: true},
		{ID: "b1", Slug: "beta", Actor: "wst-1", Bound: false},
		{ID: "b2", Slug: "beta", Actor: "wst-2", Bound: false},
	}
	// Run several times: map iteration order is randomized, so a
	// non-deterministic sort would flake here.
	for i := 0; i < 8; i++ {
		got := buildRoster(docs, active, nil)
		if len(got) != len(want) {
			t.Fatalf("got %d entries, want %d", len(got), len(want))
		}
		for j := range want {
			if got[j] != want[j] {
				t.Fatalf("entry %d = %+v, want %+v", j, got[j], want[j])
			}
		}
	}
}

// TestBuildRosterJoinsMetadataByID: buildRoster attaches resolved
// metadata by work_instances.id (the event-join key), not by slug
// or actor.
func TestBuildRosterJoinsMetadataByID(t *testing.T) {
	docs := []parsedDoc{{Slug: "alpha"}}
	active := map[string][]*ActiveWorkInstance{
		"alpha": {{ID: "wid-1", Actor: "wst-x"}},
	}
	meta := map[string]sessionMeta{
		"wid-1": {Name: "Alpha session", Detail: `{"name":"Alpha session"}`},
	}
	entries := buildRoster(docs, active, meta)
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %+v", entries)
	}
	if entries[0].Name != "Alpha session" || entries[0].Detail == "" {
		t.Errorf("metadata not joined by id: %+v", entries[0])
	}
}

func TestRenderRosterListsBoundAndUnbound(t *testing.T) {
	roots := buildTree([]parsedDoc{{Slug: "alpha", Status: "Proposed"}}, nil)
	roster := []RosterEntry{
		{Slug: "alpha", Actor: "wst-aaaa", Bound: true},
		{Slug: "typoed-not-real", Actor: "wst-bbbb", Bound: false},
	}
	html := renderRoster(t, roots, roster)

	for _, want := range []string{
		`<section class="roster-panel">`,
		`<h2 class="roster-title">Sessions</h2>`,
		`<ul class="roster-list">`,
		`>alpha</span>`,
		`>typoed-not-real</span>`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("roster missing %q; html:\n%s", want, html)
		}
	}
	if !strings.Contains(html, `class="roster-entry roster-bound"`) {
		t.Errorf("bound entry not distinguishably marked; html:\n%s", html)
	}
	if !strings.Contains(html, `class="roster-entry roster-unbound"`) {
		t.Errorf("unbound entry not distinguishably marked; html:\n%s", html)
	}
	rosterOpen := strings.Index(html, `<aside class="roster">`)
	listAt := strings.Index(html, `<ul class="roster-list">`)
	if rosterOpen < 0 || listAt < 0 || rosterOpen > listAt {
		t.Errorf("roster list must render inside the roster region (roster=%d list=%d)", rosterOpen, listAt)
	}
}

// TestRenderRosterNameLabelAndSlugFallback pins the task-level
// name-then-slug rule: the reported name is the label when
// present, the slug when not, and the wst-<uuid> actor is never
// rendered in the roster region in either case.
func TestRenderRosterNameLabelAndSlugFallback(t *testing.T) {
	roots := buildTree([]parsedDoc{{Slug: "alpha"}}, nil)
	html := renderRoster(t, roots, []RosterEntry{
		{Slug: "alpha", Actor: "wst-deadbeef-1234", Bound: true, Name: "Refactor the roster"},
		{Slug: "beta-typo", Actor: "wst-cafe-5678", Bound: false},
	})
	rosterStart := strings.Index(html, `<aside class="roster">`)
	if rosterStart < 0 {
		t.Fatalf("no roster region; html:\n%s", html)
	}
	rosterHTML := html[rosterStart:]
	if !strings.Contains(rosterHTML, `>Refactor the roster</span>`) {
		t.Errorf("named session must show its reported name; roster html:\n%s", rosterHTML)
	}
	if !strings.Contains(rosterHTML, `>beta-typo</span>`) {
		t.Errorf("no-name session must fall back to the slug; roster html:\n%s", rosterHTML)
	}
	for _, uuid := range []string{"wst-deadbeef-1234", "wst-cafe-5678"} {
		if strings.Contains(rosterHTML, uuid) {
			t.Errorf("roster must never render the wst-<uuid> actor (%q); roster html:\n%s", uuid, rosterHTML)
		}
	}
}

// TestRenderRosterExpandableDetailVsPlainRow pins scoping SD5: an
// entry with reported metadata is expandable (native <details>),
// an entry with none is a plain row with no disclosure (additive
// — absence is not a drop).
func TestRenderRosterExpandableDetailVsPlainRow(t *testing.T) {
	roots := buildTree([]parsedDoc{{Slug: "alpha"}}, nil)
	html := renderRoster(t, roots, []RosterEntry{
		{Slug: "alpha", Actor: "wst-a", Bound: true, Name: "Has detail",
			Detail: `{"name":"Has detail","pr":"#42"}`},
		{Slug: "bare-slug", Actor: "wst-b", Bound: false},
	})
	rosterStart := strings.Index(html, `<aside class="roster">`)
	rosterHTML := html[rosterStart:]

	if !strings.Contains(rosterHTML, `<details class="roster-disclosure"><summary`) {
		t.Errorf("metadata-bearing entry must be expandable via <details>; html:\n%s", rosterHTML)
	}
	if !strings.Contains(rosterHTML, `#42`) || !strings.Contains(rosterHTML, `<pre class="roster-detail">`) {
		t.Errorf("expanded detail must show the raw reported JSON; html:\n%s", rosterHTML)
	}
	// The no-metadata entry renders, but as a plain row: its label
	// is present and it is NOT wrapped in a <details>.
	bareAt := strings.Index(rosterHTML, `>bare-slug</span>`)
	if bareAt < 0 {
		t.Fatalf("no-metadata session must still list; html:\n%s", rosterHTML)
	}
	// Exactly one <details> in the roster (the metadata-bearing one).
	if n := strings.Count(rosterHTML, "<details"); n != 1 {
		t.Errorf("want exactly one <details> (the metadata entry), got %d; html:\n%s", n, rosterHTML)
	}
}

// TestRenderRosterSessionReportingNothingStillLists: a session
// that reported only the registration minimum still appears.
func TestRenderRosterSessionReportingNothingStillLists(t *testing.T) {
	roots := buildTree([]parsedDoc{{Slug: "alpha"}}, nil)
	html := renderRoster(t, roots, []RosterEntry{
		{Slug: "alpha", Actor: "wst-bare", Bound: true},
	})
	if !strings.Contains(html, `>alpha</span>`) {
		t.Errorf("a session reporting nothing must still list; html:\n%s", html)
	}
}

func TestRenderRosterEmptyState(t *testing.T) {
	roots := buildTree([]parsedDoc{{Slug: "alpha", Status: "Proposed"}}, nil)
	html := renderRoster(t, roots, nil)

	if !strings.Contains(html, `<section class="roster-panel">`) ||
		!strings.Contains(html, `<h2 class="roster-title">Sessions</h2>`) {
		t.Errorf("empty roster must still render the framed panel + heading; html:\n%s", html)
	}
	if !strings.Contains(html, `<p class="roster-empty">No active sessions.`) {
		t.Errorf("empty roster must render the deliberate observed state; html:\n%s", html)
	}
	if strings.Contains(html, `<ul class="roster-list">`) {
		t.Errorf("empty roster must not render an (empty) list element; html:\n%s", html)
	}
}

// TestRenderRosterNoAffordance: the roster lists only — no
// promote/dismiss/attach control (milestone Cross-Task Risk). The
// <details> disclosure is a read affordance, not a triage button.
func TestRenderRosterNoAffordance(t *testing.T) {
	roots := buildTree([]parsedDoc{{Slug: "alpha"}}, nil)
	html := renderRoster(t, roots, []RosterEntry{
		{Slug: "alpha", Actor: "wst-x", Bound: true, Name: "n", Detail: `{"name":"n"}`},
		{Slug: "ghost", Actor: "wst-y", Bound: false},
	})
	rosterStart := strings.Index(html, `<aside class="roster">`)
	rosterHTML := strings.ToLower(html[rosterStart:])
	for _, banned := range []string{"promote", "dismiss", "<button", "<form", "<input"} {
		if strings.Contains(rosterHTML, banned) {
			t.Errorf("roster must not render a triage affordance (%q); roster html:\n%s", banned, html[rosterStart:])
		}
	}
}

// TestResolveSessionMeta pins the task-level "Metadata read
// policy": register metadata is the identity baseline; the latest
// later event's metadata is overlaid key-by-key (later keys win;
// register-only keys survive); a non-object blob is rendered
// verbatim with no name (schema-loose, scoping SD3); no metadata
// at all yields ok=false (the entry is a plain row).
func TestResolveSessionMeta(t *testing.T) {
	// Register-only baseline.
	name, detail, ok := resolveSessionMeta(json.RawMessage(`{"name":"Reg"}`), nil)
	if !ok || name != "Reg" || !strings.Contains(detail, `"Reg"`) {
		t.Fatalf("register-only: name=%q detail=%q ok=%v", name, detail, ok)
	}

	// Later event overlays a key; register-only key survives.
	name, detail, ok = resolveSessionMeta(
		json.RawMessage(`{"name":"Reg","branch":"main"}`),
		json.RawMessage(`{"name":"Heartbeat"}`),
	)
	if !ok || name != "Heartbeat" {
		t.Fatalf("overlay: later key must win, got name=%q", name)
	}
	if !strings.Contains(detail, `"branch"`) || !strings.Contains(detail, `"main"`) {
		t.Fatalf("overlay: register-only key must survive, detail=%q", detail)
	}

	// No name reported anywhere → empty name, still has detail.
	name, _, ok = resolveSessionMeta(json.RawMessage(`{"pr":"#1"}`), nil)
	if !ok || name != "" {
		t.Fatalf("no-name object: name=%q ok=%v", name, ok)
	}

	// Non-object metadata: rendered verbatim, no name, no schema.
	name, detail, ok = resolveSessionMeta(json.RawMessage(`"just a string"`), nil)
	if !ok || name != "" || !strings.Contains(detail, "just a string") {
		t.Fatalf("non-object verbatim: name=%q detail=%q ok=%v", name, detail, ok)
	}

	// No metadata at all → not ok (entry renders as a plain row).
	if _, _, ok := resolveSessionMeta(nil, nil); ok {
		t.Fatalf("no metadata must yield ok=false (plain row)")
	}
}
