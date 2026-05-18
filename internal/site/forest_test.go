package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Forest-region render tests: node shape (badge/label/markers),
// long description, related PRs, escaping, the byte-identity
// field-less node pin, and the empty-state inside the forest
// region. Relocated from render_test.go by the m2 t1
// region-ownership split — this is t2's owned test surface.

func TestRenderLongDescriptionInline(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed", LongDescription: "First paragraph.\n\nSecond paragraph."},
	}, nil)
	html := renderTree(t, roots)
	if !strings.Contains(html, `<div class="long-desc">First paragraph.`) {
		t.Errorf("long description not rendered inline; html:\n%s", html)
	}
	if !strings.Contains(html, "Second paragraph.</div>") {
		t.Errorf("multi-paragraph body not fully rendered; html:\n%s", html)
	}
}

func TestRenderRelatedPRsAbsoluteURLIsAnchor(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed", RelatedPRs: []string{
			"https://github.com/o/r/pull/7",
		}},
	}, nil)
	html := renderTree(t, roots)
	want := `<a href="https://github.com/o/r/pull/7">https://github.com/o/r/pull/7</a>`
	if !strings.Contains(html, want) {
		t.Errorf("absolute-URL related PR not rendered as anchor; want %q in:\n%s", want, html)
	}
	if !strings.Contains(html, `<ul class="related-prs">`) {
		t.Errorf("related-prs list not emitted; html:\n%s", html)
	}
}

func TestRenderRelatedPRsNonURLIsPlainText(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed", RelatedPRs: []string{"#123"}},
	}, nil)
	html := renderTree(t, roots)
	// The non-URL entry shows as escaped plain text inside the li,
	// with no anchor and no broken in-page href.
	if !strings.Contains(html, "<li>#123</li>") {
		t.Errorf("non-URL related PR should render as plain text; html:\n%s", html)
	}
	if strings.Contains(html, `href="#123"`) {
		t.Errorf("non-URL related PR must not produce a broken in-page anchor; html:\n%s", html)
	}
	// Match an actual anchor tag, not the bare "<a" substring (the
	// shell's <aside> region tag legitimately contains "<a").
	if strings.Contains(html, "<a ") || strings.Contains(html, "<a>") {
		t.Errorf("non-URL related PR must not be wrapped in an anchor; html:\n%s", html)
	}
}

// TestRenderNestedBoxesReplaceBulletList asserts the m2 t2 C1/C6
// structural contract: a node with children renders as a native
// <details> box whose summary carries the node header, the child
// renders as a nested box inside the body, and the old
// <ul>/<li> descendant nesting is gone. Semantic/structural, not
// a byte-identity pin (C6 — the prior byte pin is replaced, not
// re-ratcheted, because t2 deliberately rewrites the node shape).
func TestRenderNestedBoxesReplaceBulletList(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed"},
		{Slug: "alpha-m1", Status: "Landed"},
	}, nil)
	html := renderTree(t, roots)

	// Parent with children is a collapsible box; child is a nested
	// leaf box, not an <li>. Expand state (open) is a separate
	// concern (TestRenderDefaultOpenByActiveWork); this is the
	// structural contract only.
	if !strings.Contains(html, `<details class="box box-root"`) {
		t.Errorf("parent node not rendered as a collapsible box; html:\n%s", html)
	}
	if !strings.Contains(html, `<summary><span class="box-header">`) {
		t.Errorf("collapsible box missing summary header; html:\n%s", html)
	}
	if !strings.Contains(html, `<div class="box box-milestone box-leaf">`) {
		t.Errorf("child node not rendered as a nested leaf box; html:\n%s", html)
	}
	if !strings.Contains(html, `<span class="label" title="alpha-m1">m1</span>`) {
		t.Errorf("child node label lost; html:\n%s", html)
	}
	// The flat descendant-nesting <ul>/<li> is fully removed (no
	// related PRs here, so any <li> would be the old nesting).
	if strings.Contains(html, "<li>") {
		t.Errorf("descendant <ul>/<li> nesting not removed; html:\n%s", html)
	}
}

// TestRenderFieldlessNodeEmitsNoDetailMarkup is the preserved
// additive guard (C5): a node carrying neither a long description
// nor related PRs emits no detail markup. The semantic
// replacement for the retired byte-identity pin.
func TestRenderFieldlessNodeEmitsNoDetailMarkup(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed"},
		{Slug: "alpha-m1", Status: "Landed"},
	}, nil)
	html := renderTree(t, roots)

	// Guard the emitted detail markup, not the CSS rule names (the
	// stylesheet legitimately defines `.long-desc`/`.related-prs`).
	if strings.Contains(html, `<div class="long-desc">`) ||
		strings.Contains(html, `<ul class="related-prs">`) {
		t.Errorf("field-less render emitted detail markup; html:\n%s", html)
	}
}

// TestRenderLeafAndStubBox asserts the C5 leaf/stub contract: a
// node with no children renders as a valid box with NO disclosure
// control (no <details>/<summary> for that node), and a bare
// `slug` + `Status: In draft` stub still renders as a leaf box
// with just its badge and label.
func TestRenderLeafAndStubBox(t *testing.T) {
	// Single root, no children: a leaf box, no disclosure control.
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed"},
	}, nil)
	html := renderTree(t, roots)
	if !strings.Contains(html, `<div class="box box-root box-leaf">`) {
		t.Errorf("childless root not rendered as a leaf box; html:\n%s", html)
	}
	if strings.Contains(html, "<details") || strings.Contains(html, "<summary") {
		t.Errorf("leaf box must have no disclosure control; html:\n%s", html)
	}

	// A `slug` + `Status: In draft` stub renders as a valid leaf
	// box with its badge and label, nothing dropped or errored.
	stub := renderTree(t, buildTree([]parsedDoc{
		{Slug: "beta", Status: "In draft"},
	}, nil))
	if !strings.Contains(stub, `<div class="box box-root box-leaf">`) {
		t.Errorf("stub not rendered as a valid leaf box; html:\n%s", stub)
	}
	if !strings.Contains(stub, `<span class="badge status-in-draft">In draft</span>`) ||
		!strings.Contains(stub, `<span class="label" title="beta">beta</span>`) {
		t.Errorf("stub box lost its badge or label; html:\n%s", stub)
	}
}

// TestRenderActorMarkersOnBox asserts the cross-cutting invariant
// (v0.1 actor tags on nodes must not regress): every box still
// renders the active-work actor markers in its header.
func TestRenderActorMarkersOnBox(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "In progress"},
		{Slug: "alpha-m1", Status: "Proposed"},
	}, map[string][]*ActiveWorkInstance{
		"alpha":    {{Actor: "agent-1"}},
		"alpha-m1": {{Actor: "agent-2"}},
	})
	html := renderTree(t, roots)
	if !strings.Contains(html, `<span class="actor-marker">agent-1</span>`) {
		t.Errorf("actor marker on parent box lost; html:\n%s", html)
	}
	if !strings.Contains(html, `<span class="actor-marker">agent-2</span>`) {
		t.Errorf("actor marker on child box lost; html:\n%s", html)
	}
}

// TestRenderEmptyStateInForestRegion confirms the no-roots
// empty-state renders inside the forest region while the roster
// placeholder still renders — the empty forest is not a blank
// region and the roster is not lost.
func TestRenderEmptyStateInForestRegion(t *testing.T) {
	html := renderTree(t, nil)

	empty := `<p class="empty">No plan-tree roots found at <code>docs/plans</code>.</p>`
	if !strings.Contains(html, empty) {
		t.Errorf("empty-state not rendered; html:\n%s", html)
	}
	// Empty-state sits inside the forest region (between the
	// forest open tag and the roster open tag).
	forestOpen := strings.Index(html, `<main class="forest">`)
	rosterOpen := strings.Index(html, `<aside class="roster">`)
	emptyAt := strings.Index(html, empty)
	if forestOpen < 0 || rosterOpen < 0 || emptyAt < 0 ||
		!(forestOpen < emptyAt && emptyAt < rosterOpen) {
		t.Errorf("empty-state must render inside the forest region (forest=%d empty=%d roster=%d)",
			forestOpen, emptyAt, rosterOpen)
	}
	// The roster region is not lost when the forest is empty. (t1
	// asserted the literal placeholder string here; m2 t4 p1
	// replaced the placeholder with the real roster, so this is
	// reconciled to the roster's deliberate empty state — a
	// plan-flagged cross-region reconciliation, see
	// t4-p1-bare-roster.md "Files to touch".)
	if !strings.Contains(html, `<p class="roster-empty">No active sessions.`) {
		t.Errorf("roster region must still render its deliberate state when the forest is empty; html:\n%s", html)
	}
}

// TestRenderDefaultOpenByActiveWork asserts the C3 default expand
// state: a collapsible box with an active work-instance anywhere
// in its subtree renders `open`, an idle subtree renders closed,
// and an active leaf forces its ancestor boxes open ("or any
// descendant") so the leaf is visible. Leaf boxes have no
// disclosure control and no expand state.
func TestRenderDefaultOpenByActiveWork(t *testing.T) {
	docs := []parsedDoc{
		// active root: active leaf m1 forces alpha open.
		{Slug: "alpha", Status: "In progress"},
		{Slug: "alpha-m1", Status: "Proposed"},
		// idle root: no active work-instance anywhere.
		{Slug: "beta", Status: "Proposed"},
		{Slug: "beta-m1", Status: "Proposed"},
	}
	roots := buildTree(docs, map[string][]*ActiveWorkInstance{
		"alpha-m1": {{Actor: "agent-1"}},
	})
	html := renderTree(t, roots)

	// alpha is collapsible (has child) and has an active
	// descendant -> open.
	if !strings.Contains(html, `<details class="box box-root" open>`) {
		t.Errorf("active-subtree root box should render open; html:\n%s", html)
	}
	// beta is collapsible but idle -> closed (no open attr).
	if !strings.Contains(html, `<details class="box box-root">`) {
		t.Errorf("idle root box should render closed (no open attr); html:\n%s", html)
	}
	// The active leaf alpha-m1 has no disclosure control at all.
	if !strings.Contains(html, `<div class="box box-milestone box-leaf">`) {
		t.Errorf("leaf box must have no expand state; html:\n%s", html)
	}
}

// declaredCellCount counts the declared (non-Drafting) progress
// cells in the rendered HTML. The Drafting cell carries the
// extra progress-cell-drafting class, so a declared cell is the
// bare `<span class="progress-cell">` opening tag.
func declaredCellCount(html string) int {
	return strings.Count(html, `<span class="progress-cell">`)
}

func draftingCellCount(html string) int {
	return strings.Count(html, `<span class="progress-cell progress-cell-drafting">Drafting</span>`)
}

// TestRenderProgressRowDeclaredStages asserts m2 t3 C2: a doc
// declaring N stages renders the reserved Drafting cell followed
// by one cell per declared stage, in document order (N + 1 cells).
func TestRenderProgressRowDeclaredStages(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed", ProgressStages: []string{"Spec", "Parser", "Render"}},
	}, nil)
	html := renderTree(t, roots)

	if got := draftingCellCount(html); got != 1 {
		t.Errorf("Drafting cell count = %d, want exactly 1; html:\n%s", got, html)
	}
	if got := declaredCellCount(html); got != 3 {
		t.Errorf("declared cell count = %d, want 3 (N declared ⇒ N+1 cells); html:\n%s", got, html)
	}
	// Reserved Drafting cell precedes the declared cells, and the
	// declared cells appear in document order.
	iDrafting := strings.Index(html, `progress-cell-drafting`)
	iSpec := strings.Index(html, `<span class="progress-cell">Spec</span>`)
	iParser := strings.Index(html, `<span class="progress-cell">Parser</span>`)
	iRender := strings.Index(html, `<span class="progress-cell">Render</span>`)
	if !(iDrafting >= 0 && iDrafting < iSpec && iSpec < iParser && iParser < iRender) {
		t.Errorf("progress cells out of order: drafting=%d Spec=%d Parser=%d Render=%d; html:\n%s",
			iDrafting, iSpec, iParser, iRender, html)
	}
}

// TestRenderProgressRowFieldlessOnlyDrafting asserts m2 t3 C3: a
// field-omitting doc renders exactly the one reserved Drafting
// cell as an intentional observed state — no declared cells, no
// errored or empty row.
func TestRenderProgressRowFieldlessOnlyDrafting(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed"},
	}, nil)
	html := renderTree(t, roots)

	if !strings.Contains(html, `<div class="progress-row">`) {
		t.Errorf("field-omitting doc must still render the progress row; html:\n%s", html)
	}
	if got := draftingCellCount(html); got != 1 {
		t.Errorf("Drafting cell count = %d, want exactly 1; html:\n%s", got, html)
	}
	if got := declaredCellCount(html); got != 0 {
		t.Errorf("declared cell count = %d, want 0 for a field-omitting doc; html:\n%s", got, html)
	}
}

// TestRenderProgressRowStubOnlyDrafting asserts m2 t3 C3 for the
// already-supported `slug` + `Status: In draft` stub: it renders
// exactly the one Drafting cell, Status-independently (C5 — the
// gating is field presence, never the Status token).
func TestRenderProgressRowStubOnlyDrafting(t *testing.T) {
	stub := renderTree(t, buildTree([]parsedDoc{
		{Slug: "beta", Status: "In draft"},
	}, nil))

	if got := draftingCellCount(stub); got != 1 {
		t.Errorf("stub Drafting cell count = %d, want exactly 1; html:\n%s", got, stub)
	}
	if got := declaredCellCount(stub); got != 0 {
		t.Errorf("stub declared cell count = %d, want 0; html:\n%s", got, stub)
	}
	// The stub still renders its badge and label (preserved t2
	// surface) alongside the new progress row.
	if !strings.Contains(stub, `<span class="badge status-in-draft">In draft</span>`) ||
		!strings.Contains(stub, `<span class="label" title="beta">beta</span>`) {
		t.Errorf("stub lost its badge or label; html:\n%s", stub)
	}
}

// TestRenderProgressRowMalformedNeverDropsNode ties the parser's
// tolerance (m2 t3 C1) to the render: a doc whose progress_stages
// is partially malformed (a non-string element) never errors or
// drops the node — the node still renders, with the Drafting cell
// and only the well-formed declared entries.
func TestRenderProgressRowMalformedNeverDropsNode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "epic-a", "README.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	content := "---\nslug: epic-a\nStatus: Proposed\nprogress_stages:\n  - Spec\n  - 7\n  - Render\n---\n# x\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	docs, err := walkPlans(dir)
	if err != nil {
		t.Fatalf("walkPlans: %v", err)
	}
	if len(docs) != 1 {
		t.Fatalf("walkPlans dropped the malformed doc: got %v", docs)
	}
	html := renderTree(t, buildTree(docs, nil))

	if !strings.Contains(html, `<span class="label" title="epic-a">epic-a</span>`) {
		t.Errorf("node with partially-malformed progress_stages was dropped; html:\n%s", html)
	}
	if got := draftingCellCount(html); got != 1 {
		t.Errorf("Drafting cell count = %d, want exactly 1; html:\n%s", got, html)
	}
	if got := declaredCellCount(html); got != 2 {
		t.Errorf("declared cell count = %d, want 2 (non-string dropped); html:\n%s", got, html)
	}
}

// TestRenderProgressRowCoexistsWithPreservedSurfaces asserts m2 t3
// C4: the progress row renders alongside — not in place of — the
// preserved t2 surfaces (Status badge, actor markers, long
// description, related-PR list).
func TestRenderProgressRowCoexistsWithPreservedSurfaces(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "In progress",
			LongDescription: "The long body.",
			RelatedPRs:      []string{"https://github.com/o/r/pull/9"},
			ProgressStages:  []string{"Spec", "Render"}},
	}, map[string][]*ActiveWorkInstance{
		"alpha": {{Actor: "agent-1"}},
	})
	html := renderTree(t, roots)

	for _, want := range []string{
		`<span class="badge status-in-progress">In progress</span>`,
		`<span class="actor-marker">agent-1</span>`,
		`<div class="long-desc">The long body.</div>`,
		`<a href="https://github.com/o/r/pull/9">https://github.com/o/r/pull/9</a>`,
		`<div class="progress-row">`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("preserved/added surface missing %q; html:\n%s", want, html)
		}
	}
}

func TestRenderURLAndTextEscaping(t *testing.T) {
	// An entry that is a non-URL string containing HTML-special
	// characters must be escaped by html/template, never emitted
	// raw and never linkified.
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed",
			LongDescription: "<script>alert(1)</script>",
			RelatedPRs:      []string{`<b>not a url</b>`}},
	}, nil)
	html := renderTree(t, roots)
	if strings.Contains(html, "<script>alert(1)</script>") {
		t.Errorf("long description not escaped; html:\n%s", html)
	}
	if strings.Contains(html, "<b>not a url</b>") {
		t.Errorf("non-URL related PR not escaped; html:\n%s", html)
	}
}
