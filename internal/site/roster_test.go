package site

import (
	"bytes"
	"strings"
	"testing"
)

// Roster-region tests (m2 t4 p1): the bare bound/unbound session
// roster that replaces t1's placeholder. Covers the buildRoster
// classification (incl. a parent-gapped-but-walked doc still
// bound, and unbound listed not dropped), deterministic order,
// slug-as-label with no wst-<uuid> rendered, the deliberate
// empty state, and the no-affordance guard. This is t4's owned
// test surface.

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
		"alpha":           {{Actor: "wst-aaaa"}},
		"alpha-m1":        {{Actor: "wst-bbbb"}},
		"typoed-not-real": {{Actor: "wst-cccc"}},
	}
	entries := buildRoster(docs, active)

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

// TestBuildRosterParentGappedDocStillBound pins scoping SD1: the
// bound set is the parsed-doc slug set, NOT the post-buildTree
// tree. "alpha-m1-t1" has no parsed parent "alpha-m1", so
// buildTree drops it from the rendered forest — but it is still a
// walked plan doc, so the roster must classify it bound.
func TestBuildRosterParentGappedDocStillBound(t *testing.T) {
	docs := []parsedDoc{
		{Slug: "alpha", Status: "Proposed"},
		{Slug: "alpha-m1-t1", Status: "In draft"}, // parent alpha-m1 absent
	}
	// Confirm the premise: buildTree drops the gapped node.
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
		"alpha-m1-t1": {{Actor: "wst-dddd"}},
	})
	if len(entries) != 1 || !entries[0].Bound {
		t.Errorf("parent-gapped-but-walked doc must be bound: %+v", entries)
	}
}

func TestBuildRosterDeterministicOrder(t *testing.T) {
	docs := []parsedDoc{{Slug: "alpha"}}
	active := map[string][]*ActiveWorkInstance{
		"beta":  {{Actor: "wst-2"}, {Actor: "wst-1"}},
		"alpha": {{Actor: "wst-9"}},
	}
	want := []RosterEntry{
		{Slug: "alpha", Actor: "wst-9", Bound: true},
		{Slug: "beta", Actor: "wst-1", Bound: false},
		{Slug: "beta", Actor: "wst-2", Bound: false},
	}
	// Run several times: map iteration order is randomized, so a
	// non-deterministic sort would flake here.
	for i := 0; i < 8; i++ {
		got := buildRoster(docs, active)
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
		`<span class="roster-slug">alpha</span>`,
		`<span class="roster-slug">typoed-not-real</span>`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("roster missing %q; html:\n%s", want, html)
		}
	}
	// Bound vs. unbound is visibly distinguishable (a per-entry
	// class + tag), not merely a struct field.
	if !strings.Contains(html, `class="roster-entry roster-bound"`) {
		t.Errorf("bound entry not distinguishably marked; html:\n%s", html)
	}
	if !strings.Contains(html, `class="roster-entry roster-unbound"`) {
		t.Errorf("unbound entry not distinguishably marked; html:\n%s", html)
	}
	// The entries render inside the roster region.
	rosterOpen := strings.Index(html, `<aside class="roster">`)
	listAt := strings.Index(html, `<ul class="roster-list">`)
	if rosterOpen < 0 || listAt < 0 || rosterOpen > listAt {
		t.Errorf("roster list must render inside the roster region (roster=%d list=%d)", rosterOpen, listAt)
	}
}

// TestRenderRosterNoUUIDLabel pins scoping SD2: the wst-<uuid>
// actor is never rendered in the roster region in p1; the label
// is the slug.
func TestRenderRosterNoUUIDLabel(t *testing.T) {
	roots := buildTree([]parsedDoc{{Slug: "alpha"}}, nil)
	html := renderRoster(t, roots, []RosterEntry{
		{Slug: "alpha", Actor: "wst-deadbeef-1234-5678", Bound: true},
	})
	// Scope the assertion to the roster region (the forest region
	// legitimately renders actor markers — that is the intentional,
	// backlog-tracked inconsistency, not p1's surface).
	rosterStart := strings.Index(html, `<aside class="roster">`)
	if rosterStart < 0 {
		t.Fatalf("no roster region; html:\n%s", html)
	}
	rosterHTML := html[rosterStart:]
	if strings.Contains(rosterHTML, "wst-deadbeef-1234-5678") {
		t.Errorf("roster must not render the wst-<uuid> actor; roster html:\n%s", rosterHTML)
	}
	if !strings.Contains(rosterHTML, `<span class="roster-slug">alpha</span>`) {
		t.Errorf("roster entry label must be the slug; roster html:\n%s", rosterHTML)
	}
}

// TestRenderRosterSessionReportingNothingStillLists: a session
// that reported only the registration minimum (slug + actor)
// still appears — enrichment is additive (p2), absence is not a
// drop.
func TestRenderRosterSessionReportingNothingStillLists(t *testing.T) {
	roots := buildTree([]parsedDoc{{Slug: "alpha"}}, nil)
	html := renderRoster(t, roots, []RosterEntry{
		{Slug: "alpha", Actor: "wst-bare", Bound: true},
	})
	if !strings.Contains(html, `<span class="roster-slug">alpha</span>`) {
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

// TestRenderRosterNoAffordance: the roster lists only — it must
// not pre-empt the deferred triage action with a
// promote/dismiss/attach control (milestone Cross-Task Risk).
func TestRenderRosterNoAffordance(t *testing.T) {
	roots := buildTree([]parsedDoc{{Slug: "alpha"}}, nil)
	html := renderRoster(t, roots, []RosterEntry{
		{Slug: "alpha", Actor: "wst-x", Bound: true},
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
