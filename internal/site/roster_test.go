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
// name-then-slug rule across the rendered .roster-label
// surface: the reported name is the label when present, the
// slug when not, and the wst-<uuid> actor never surfaces in any
// .roster-label element in either case. p3 F9 K3 introduces an
// explicit "Actor id:" facts-block field inside the disclosed
// body — the actor id is allowed there as a deliberately
// labeled facts-block surface, just not in the entry's label
// (the identity rendering rule t4 and p2 locked).
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
		for _, content := range rosterLabelContents(rosterHTML) {
			if strings.Contains(content, uuid) {
				t.Errorf("roster label must never carry the wst-<uuid> actor (%q in %q); roster html:\n%s",
					uuid, content, rosterHTML)
			}
		}
		// The K3 facts block IS allowed to surface the actor id;
		// confirm each uuid appears there (positive guard that the
		// data is reaching the disclosure, not the negative guard
		// the rule is about).
		if !strings.Contains(rosterHTML, `<dd class="roster-fact-actor">`+uuid+`</dd>`) {
			t.Errorf("K3 known-facts header must surface actor id %q; roster html:\n%s", uuid, rosterHTML)
		}
	}
}

// rosterLabelContents extracts the text content of every
// <span class="roster-label...">...</span> element in the
// rendered roster HTML. Used by the F9 + p2 identity-rendering
// invariant tests to assert the wst-<uuid> actor never surfaces
// as an entry's identity label, while still allowing it to
// appear inside the K3 known-facts header's Actor id field.
func rosterLabelContents(html string) []string {
	const open = `<span class="roster-label`
	const close = `</span>`
	var out []string
	rest := html
	for {
		i := strings.Index(rest, open)
		if i < 0 {
			return out
		}
		// Skip to the closing '>' of the opening tag.
		gt := strings.Index(rest[i:], ">")
		if gt < 0 {
			return out
		}
		start := i + gt + 1
		j := strings.Index(rest[start:], close)
		if j < 0 {
			return out
		}
		out = append(out, rest[start:start+j])
		rest = rest[start+j+len(close):]
	}
}

// TestRenderRosterEveryEntryOpensToK3Disclosure pins p3 F9
// (parent C3; supersedes t4's "no metadata ⇒ plain row" under
// parent D3): every roster entry — bound or unbound, with or
// without reported metadata — opens to the same K3-shape
// disclosure. The metadata-bearing entry shows the raw-JSON
// block; the no-metadata entry shows the no-metadata sentinel
// in place of the block; both share the K3 known-facts header.
// The same outer disclosure structure across every entry is
// what the contract locks.
func TestRenderRosterEveryEntryOpensToK3Disclosure(t *testing.T) {
	roots := buildTree([]parsedDoc{{Slug: "alpha"}}, nil)
	html := renderRoster(t, roots, []RosterEntry{
		{Slug: "alpha", Actor: "wst-a", Bound: true, Name: "Has detail",
			Detail:       `{"name":"Has detail","pr":"#42"}`,
			RegisteredAt: 1700000000 * 1e9, LastEventAt: 1700000600 * 1e9},
		{Slug: "bare-slug", Actor: "wst-b", Bound: false,
			RegisteredAt: 1700000100 * 1e9},
	})
	rosterStart := strings.Index(html, `<aside class="roster">`)
	rosterHTML := html[rosterStart:]

	// Every entry opens to a <details> — superseded D3: no plain
	// row branch for no-metadata entries.
	if n := strings.Count(rosterHTML, `<details class="roster-disclosure">`); n != 2 {
		t.Errorf("every entry must open to a roster-disclosure <details>; got %d (want 2); html:\n%s",
			n, rosterHTML)
	}
	// Every entry's body carries the K3 known-facts header.
	if n := strings.Count(rosterHTML, `<dl class="roster-facts">`); n != 2 {
		t.Errorf("every entry must render the K3 known-facts header; got %d (want 2); html:\n%s",
			n, rosterHTML)
	}
	// Metadata-bearing entry: raw-JSON block present, no
	// no-metadata sentinel inside its body.
	if !strings.Contains(rosterHTML, `<pre class="roster-detail">`) ||
		!strings.Contains(rosterHTML, `#42`) {
		t.Errorf("metadata-bearing entry must show the raw-JSON block; html:\n%s", rosterHTML)
	}
	// No-metadata entry: sentinel present in place of the
	// raw-JSON block.
	if !strings.Contains(rosterHTML, `<p class="roster-empty-meta">(no reported metadata)</p>`) {
		t.Errorf("no-metadata entry must show the no-metadata sentinel; html:\n%s", rosterHTML)
	}
	// Same outer disclosure structure for every entry: counting
	// the per-entry <div class="roster-body"> wrappers.
	if n := strings.Count(rosterHTML, `<div class="roster-body">`); n != 2 {
		t.Errorf("every entry must wrap its disclosed body identically; got %d (want 2); html:\n%s",
			n, rosterHTML)
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

// TestRenderRosterFourObservableStatesAllOpen asserts p3 F9 across
// the four observable conditions the parent task plan's Validation
// Gate names: (a) name-bearing bound, (b) no-name bound, (c)
// name-bearing unbound, (d) no-name unbound. Every state renders
// the same outer disclosure (roster-disclosure <details> + body
// + K3 known-facts header + raw-JSON or sentinel branch). No
// wst-<uuid> surfaces in any .roster-label across any state.
func TestRenderRosterFourObservableStatesAllOpen(t *testing.T) {
	roots := buildTree([]parsedDoc{{Slug: "alpha"}}, nil)
	html := renderRoster(t, roots, []RosterEntry{
		{Slug: "alpha", Actor: "wst-aaaa", Bound: true, Name: "DemoBoundNamed",
			Detail: `{"name":"DemoBoundNamed"}`, RegisteredAt: 1700000000 * 1e9, LastEventAt: 1700000600 * 1e9},
		{Slug: "alpha", Actor: "wst-bbbb", Bound: true,
			RegisteredAt: 1700000100 * 1e9},
		{Slug: "unbound-named", Actor: "wst-cccc", Bound: false, Name: "DemoUnboundNamed",
			Detail: `{"name":"DemoUnboundNamed"}`, RegisteredAt: 1700000200 * 1e9, LastEventAt: 1700000700 * 1e9},
		{Slug: "unbound-bare", Actor: "wst-dddd", Bound: false,
			RegisteredAt: 1700000300 * 1e9},
	})
	rosterStart := strings.Index(html, `<aside class="roster">`)
	rosterHTML := html[rosterStart:]

	// Four entries, each opens to a roster-disclosure.
	if n := strings.Count(rosterHTML, `<details class="roster-disclosure">`); n != 4 {
		t.Errorf("every observable state must open; got %d disclosures (want 4); html:\n%s",
			n, rosterHTML)
	}
	if n := strings.Count(rosterHTML, `<dl class="roster-facts">`); n != 4 {
		t.Errorf("every observable state must render the K3 known-facts header; got %d (want 4); html:\n%s",
			n, rosterHTML)
	}
	// Two metadata-bearing entries → two raw-JSON blocks; two
	// no-metadata entries → two sentinels.
	if n := strings.Count(rosterHTML, `<pre class="roster-detail">`); n != 2 {
		t.Errorf("metadata-bearing entries must show raw-JSON block; got %d (want 2); html:\n%s",
			n, rosterHTML)
	}
	if n := strings.Count(rosterHTML, `<p class="roster-empty-meta">(no reported metadata)</p>`); n != 2 {
		t.Errorf("no-metadata entries must show the sentinel; got %d (want 2); html:\n%s",
			n, rosterHTML)
	}
	// No wst-<uuid> text in any .roster-label across all four
	// states (the identity-rendering rule from t4 + p2; F9's K3
	// header surfaces actor id only inside the K3 facts block,
	// not in the entry label).
	for _, content := range rosterLabelContents(rosterHTML) {
		if strings.Contains(content, "wst-") {
			t.Errorf("no .roster-label may carry wst-<uuid>; saw %q in:\n%s",
				content, rosterHTML)
		}
	}
}

// TestRenderRosterK3TimestampsRender asserts p3 F9 OD6: the K3
// known-facts header renders the registered-at and last-event
// timestamps from the per-request loader's already-tracked
// values. The "Last event" field falls back to the
// registered-at timestamp when no later event has been seen
// (LastEventAt == 0). Zero registered-at renders as the
// em-dash placeholder.
func TestRenderRosterK3TimestampsRender(t *testing.T) {
	roots := buildTree([]parsedDoc{{Slug: "alpha"}}, nil)
	html := renderRoster(t, roots, []RosterEntry{
		// Has both timestamps — both render. Values are
		// Unix-epoch nanoseconds matching events.received_at's
		// storage shape (now.UnixNano() in handlers.go); the
		// seconds-shaped epoch numbers are scaled by 1e9 here
		// so the expected formatted strings stay readable.
		{Slug: "alpha", Actor: "wst-1", Bound: true, Name: "Both",
			RegisteredAt: 1700000000 * 1e9, LastEventAt: 1700000600 * 1e9},
		// Only register seen (LastEventAt == 0) — last-event
		// falls back to registered-at.
		{Slug: "alpha", Actor: "wst-2", Bound: true, Name: "RegisterOnly",
			RegisteredAt: 1700000100 * 1e9},
		// Neither (both zero) — placeholder for both.
		{Slug: "alpha", Actor: "wst-3", Bound: true, Name: "NoTimestamps"},
	})
	rosterStart := strings.Index(html, `<aside class="roster">`)
	rosterHTML := html[rosterStart:]

	// Both-timestamps entry: each timestamp surfaces in its own
	// <dd> cell. The exact rendered string format is the
	// formatEventTime helper's UTC RFC-3339-without-T shape;
	// asserting on the prefix "2023-" (the year Unix-epoch
	// 1700000000 falls in) keeps the test format-tolerant.
	if !strings.Contains(rosterHTML, `<dd>2023-11-14 22:13:20 UTC</dd>`) {
		t.Errorf("registered-at must render formatted; html:\n%s", rosterHTML)
	}
	if !strings.Contains(rosterHTML, `<dd>2023-11-14 22:23:20 UTC</dd>`) {
		t.Errorf("last-event must render formatted when distinct from registered; html:\n%s", rosterHTML)
	}
	// Register-only entry: last-event falls back to register
	// timestamp. Both <dd> cells for THAT entry render the
	// register's formatted string; asserting the count >= 2 of
	// the register timestamp tests the fallback without
	// over-constraining the rest of the rendered output.
	registerOnlyAt := strings.Index(rosterHTML, `<span class="roster-label">RegisterOnly`)
	if registerOnlyAt < 0 {
		t.Fatalf("RegisterOnly entry missing; html:\n%s", rosterHTML)
	}
	// Find the bounded HTML of this one entry.
	registerOnlyEnd := strings.Index(rosterHTML[registerOnlyAt:], `</li>`)
	registerOnlyHTML := rosterHTML[registerOnlyAt : registerOnlyAt+registerOnlyEnd]
	if n := strings.Count(registerOnlyHTML, `<dd>2023-11-14 22:15:00 UTC</dd>`); n != 2 {
		t.Errorf("register-only entry must repeat registered timestamp for last-event; got %d (want 2); entry html:\n%s",
			n, registerOnlyHTML)
	}
	// No-timestamps entry: em-dash placeholder for both fields.
	noTimestampsAt := strings.Index(rosterHTML, `<span class="roster-label">NoTimestamps`)
	if noTimestampsAt < 0 {
		t.Fatalf("NoTimestamps entry missing; html:\n%s", rosterHTML)
	}
	noTimestampsEnd := strings.Index(rosterHTML[noTimestampsAt:], `</li>`)
	noTimestampsHTML := rosterHTML[noTimestampsAt : noTimestampsAt+noTimestampsEnd]
	if n := strings.Count(noTimestampsHTML, `<dd>—</dd>`); n != 2 {
		t.Errorf("zero-timestamp entry must render em-dash for both fields; got %d (want 2); entry html:\n%s",
			n, noTimestampsHTML)
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
