package site

import (
	"fmt"
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

// TestRenderLongDescriptionInsideBodyDisclosure asserts p3 F7
// (parent C2): the long description does not render in the
// header-only first view; it sits inside a nested <details> body
// disclosure that the reader opens explicitly. When opened, the
// body renders as markdown via goldmark (paragraphs wrapped in
// <p>) per F7 OD3 = I2b.
func TestRenderLongDescriptionInsideBodyDisclosure(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed", LongDescription: "First paragraph.\n\nSecond paragraph."},
	}, nil)
	html := renderTree(t, roots)

	// The body sits inside a body-disclosure <details>, not free
	// in the node body. The .long-desc wrapper still exists for
	// styling but lives under .body-disclosure.
	bodyDiscAt := strings.Index(html, `<details class="body-disclosure">`)
	if bodyDiscAt < 0 {
		t.Fatalf("body-disclosure <details> wrapper missing; html:\n%s", html)
	}
	longDescAt := strings.Index(html, `<div class="long-desc">`)
	if longDescAt < 0 || longDescAt < bodyDiscAt {
		t.Errorf("long-desc must render inside body-disclosure; bodyDisc=%d longDesc=%d html:\n%s",
			bodyDiscAt, longDescAt, html)
	}
	// goldmark wraps each paragraph in <p>; the markdown render
	// is what surfaces the body content (not the raw plain-text
	// fall-through path the pre-p3 inline render used).
	if !strings.Contains(html, "<p>First paragraph.</p>") {
		t.Errorf("first paragraph not rendered as markdown <p>; html:\n%s", html)
	}
	if !strings.Contains(html, "<p>Second paragraph.</p>") {
		t.Errorf("multi-paragraph body not fully rendered through markdown; html:\n%s", html)
	}
}

// TestRenderLongDescriptionLineCapped asserts a long description
// past the line cap is truncated even when the box is expanded:
// the head lines render, the tail is dropped, and a visible
// truncation marker is appended (the forest is a context-and-goal
// overview, not a full plan-file viewer).
func TestRenderLongDescriptionLineCapped(t *testing.T) {
	var b strings.Builder
	for i := 0; i < maxLongDescLines+50; i++ {
		fmt.Fprintf(&b, "line %d\n", i)
	}
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed", LongDescription: b.String()},
	}, nil)
	html := renderTree(t, roots)

	if !strings.Contains(html, "line 0") {
		t.Errorf("head of long description should render; html:\n%s", html)
	}
	if !strings.Contains(html, "… (truncated — see the plan doc for the full text)") {
		t.Errorf("truncation marker missing; html:\n%s", html)
	}
	// A line well past the cap must not appear.
	if strings.Contains(html, fmt.Sprintf("line %d", maxLongDescLines+40)) {
		t.Errorf("tail past the line cap should be dropped; html:\n%s", html)
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
	if !strings.Contains(html, `<details class="box box-milestone box-leaf">`) {
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

// TestRenderSummaryCarriesExplicitMarkerTreatment asserts the m2
// post-m2-ux-correction p1 F8 contract: every node's <summary>
// carries an explicit marker treatment so the disclosure marker is
// baseline-aligned and does not overlap the box border (the UA-
// default marker is suppressed and an in-flow triangle element is
// rendered inside the header flex row). Semantic-not-byte-exact
// per the m2 t3 C7 precedent: the rule set carries the declarations,
// and a triangle element is present inside the summary; the visual
// result is observed via the manual walkthrough, not asserted here.
func TestRenderSummaryCarriesExplicitMarkerTreatment(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed"},
		{Slug: "alpha-m1", Status: "Landed"},
	}, nil)
	html := renderTree(t, roots)

	// The rule set suppresses the UA-default marker.
	if !strings.Contains(html, "list-style: none") {
		t.Errorf("summary rule set must suppress the UA-default marker via list-style: none; html:\n%s", html)
	}
	// An inline triangle element renders inside the summary's
	// header flex row, baseline-aligned with the label.
	if !strings.Contains(html, `<span class="triangle"`) {
		t.Errorf("summary missing the inline triangle marker element; html:\n%s", html)
	}
	// The triangle sits inside the box-header so the existing
	// baseline-aligned flex row aligns it with the label.
	iHeader := strings.Index(html, `<span class="box-header">`)
	iTriangle := strings.Index(html, `<span class="triangle"`)
	iLabel := strings.Index(html, `<span class="label-group">`)
	if !(iHeader >= 0 && iHeader < iTriangle && iTriangle < iLabel) {
		t.Errorf("triangle must sit inside box-header, before the label-group: header=%d triangle=%d label=%d; html:\n%s",
			iHeader, iTriangle, iLabel, html)
	}
	// The open-state rotation rule exists so the marker reads as
	// "open" vs. "closed" in both <details> states.
	if !strings.Contains(html, "details[open] > summary .triangle") {
		t.Errorf("missing details[open] rule that rotates the triangle for the open state; html:\n%s", html)
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

// TestRenderLeafAndStubBox asserts the leaf/stub contract: a node
// with no children renders as a collapsible <details> box (with
// the box-leaf class for styling), closed by default when there is
// no active work so its body text is hidden until expanded, and a
// bare `slug` + `Status: In draft` stub still renders as a leaf
// box with its badge and label.
func TestRenderLeafAndStubBox(t *testing.T) {
	// Single root, no children: a collapsible leaf box, closed by
	// default (no active work), so no body text shows until opened.
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed"},
	}, nil)
	html := renderTree(t, roots)
	if !strings.Contains(html, `<details class="box box-root box-leaf">`) {
		t.Errorf("childless root not rendered as a collapsible leaf box; html:\n%s", html)
	}
	if strings.Contains(html, `box box-root box-leaf" open`) {
		t.Errorf("idle leaf box must be closed by default (no open attr); html:\n%s", html)
	}
	if !strings.Contains(html, `<summary><span class="box-header">`) {
		t.Errorf("leaf box must carry its header in a <summary>; html:\n%s", html)
	}

	// A `slug` + `Status: In draft` stub renders as a valid leaf
	// box with its badge and label, nothing dropped or errored.
	stub := renderTree(t, buildTree([]parsedDoc{
		{Slug: "beta", Status: "In draft"},
	}, nil))
	if !strings.Contains(stub, `<details class="box box-root box-leaf">`) {
		t.Errorf("stub not rendered as a valid leaf box; html:\n%s", stub)
	}
	if !strings.Contains(stub, `<span class="badge status-in-draft">In draft</span>`) ||
		!strings.Contains(stub, `<span class="label" title="beta">beta</span>`) {
		t.Errorf("stub box lost its badge or label; html:\n%s", stub)
	}
}

// TestRenderActorMarkersOnBox asserts the cross-cutting invariant
// (v0.1 actor tags on nodes must not regress): every box still
// renders an active-work actor-marker span in its header per
// attached work-instance. The rendered text is the resolved
// display name (Name with Slug fallback); the post-m2-ux-correction
// p2 F4 contract changed the marker's text source but not its
// presence or attachment.
func TestRenderActorMarkersOnBox(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "In progress"},
		{Slug: "alpha-m1", Status: "Proposed"},
	}, map[string][]*ActiveWorkInstance{
		"alpha":    {{Actor: "wst-1", Name: "agent-1", Slug: "alpha"}},
		"alpha-m1": {{Actor: "wst-2", Name: "agent-2", Slug: "alpha-m1"}},
	})
	html := renderTree(t, roots)
	if !strings.Contains(html, `<span class="actor-marker">agent-1</span>`) {
		t.Errorf("actor marker on parent box lost; html:\n%s", html)
	}
	if !strings.Contains(html, `<span class="actor-marker">agent-2</span>`) {
		t.Errorf("actor marker on child box lost; html:\n%s", html)
	}
}

// TestRenderForestActorMarkerNameBearing asserts post-m2-ux-correction
// p2 C1 / OD4.a for a name-bearing session: the actor-marker text
// is the reported display name, and the raw wst-<uuid> actor never
// surfaces alongside it.
func TestRenderForestActorMarkerNameBearing(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "In progress"},
	}, map[string][]*ActiveWorkInstance{
		"alpha": {{Actor: "wst-abc-123", Name: "Demo bound", Slug: "alpha"}},
	})
	html := renderTree(t, roots)
	if !strings.Contains(html, `<span class="actor-marker">Demo bound</span>`) {
		t.Errorf("name-bearing session must render its reported name; html:\n%s", html)
	}
	if strings.Contains(html, "wst-") {
		t.Errorf("forest actor-marker must not surface the wst- prefix; html:\n%s", html)
	}
}

// TestRenderForestActorMarkerSlugFallback asserts post-m2-ux-correction
// p2 C1 / OD4.a for a no-name session: the actor-marker falls back
// to the work-instance's registered slug (never the wst-<uuid>
// actor).
func TestRenderForestActorMarkerSlugFallback(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "In progress"},
	}, map[string][]*ActiveWorkInstance{
		"alpha": {{Actor: "wst-abc-123", Slug: "alpha"}},
	})
	html := renderTree(t, roots)
	if !strings.Contains(html, `<span class="actor-marker">alpha</span>`) {
		t.Errorf("no-name session must render the slug fallback; html:\n%s", html)
	}
	if strings.Contains(html, "wst-") {
		t.Errorf("forest actor-marker must not surface the wst- prefix; html:\n%s", html)
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
// descendant") so the leaf is visible. Leaf boxes are collapsible
// too: an idle leaf is closed, an active leaf renders open.
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
	// The active leaf alpha-m1 is a collapsible box rendered open
	// (its own active work-instance makes it active-in-subtree).
	if !strings.Contains(html, `<details class="box box-milestone box-leaf" open>`) {
		t.Errorf("active leaf box should render open; html:\n%s", html)
	}
	// The idle leaf beta-m1 is collapsible but closed by default.
	if !strings.Contains(html, `<details class="box box-milestone box-leaf">`) {
		t.Errorf("idle leaf box should render closed (no open attr); html:\n%s", html)
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

// TestRenderDefaultProgressRowPerStatusBucket asserts p3 F3a
// (parent C2; supersedes m2 t3 D5 under parent D1 and the
// stub-children stub-render contract under parent D2): a doc
// declaring no progress_stages renders the default D / P / I / V
// row, Status-driven, in one of three shape buckets:
//   - Landed ⇒ all four cells "landed" (filled green).
//   - In draft ⇒ D filled "in-draft" (amber); P / I / V "empty"
//     (dashed border). The stub case (slug + Status: In draft)
//     reads this branch — supersedes the previous stub-only
//     Drafting cell render.
//   - Everything else (In progress, Proposed, Validating,
//     Deferred, any unrecognized Status) ⇒ all four "neutral"
//     (filled grey). Reads the statusClass fallback to "unknown"
//     for the unrecognized case.
//
// Across all three buckets the per-cell DOM is preserved as
// every-cell-is-its-own-element per C-INV-1 (cell-anchor) — the
// F3b future per-cell attachment surface is what depends on this.
func TestRenderDefaultProgressRowPerStatusBucket(t *testing.T) {
	cases := []struct {
		name        string
		status      string
		wantClasses [4]string
	}{
		{"Landed", "Landed", [4]string{"landed", "landed", "landed", "landed"}},
		{"In draft", "In draft", [4]string{"in-draft", "empty", "empty", "empty"}},
		{"In progress", "In progress", [4]string{"neutral", "neutral", "neutral", "neutral"}},
		{"Proposed", "Proposed", [4]string{"neutral", "neutral", "neutral", "neutral"}},
		{"Validating", "Validating", [4]string{"neutral", "neutral", "neutral", "neutral"}},
		{"Deferred", "Deferred", [4]string{"neutral", "neutral", "neutral", "neutral"}},
		{"Deferred with reason", "Deferred — out of scope", [4]string{"neutral", "neutral", "neutral", "neutral"}},
		{"unrecognized falls through to neutral", "Frobnicating", [4]string{"neutral", "neutral", "neutral", "neutral"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			roots := buildTree([]parsedDoc{
				{Slug: "alpha", Status: tc.status},
			}, nil)
			html := renderTree(t, roots)

			// The default row replaces the prior "exactly one
			// reserved Drafting cell" render for a field-omitting
			// doc (parent D1 / D2 supersession): there must be no
			// Drafting cell in the default branch.
			if draftingCellCount(html) != 0 {
				t.Errorf("default row must not carry a Drafting cell (superseded by D1/D2); html:\n%s", html)
			}
			labels := [4]string{"D", "P", "I", "V"}
			for i, label := range labels {
				want := `<span class="progress-cell progress-cell-` + tc.wantClasses[i] + `">` + label + `</span>`
				if !strings.Contains(html, want) {
					t.Errorf("cell %d (%s) shape mismatch; want %q in:\n%s", i, label, want, html)
				}
			}
		})
	}
}

// TestRenderProgressRowCellDOMPreservedAcrossBranches asserts
// p3 C-INV-1 (cell-anchor): both the default row (no
// progress_stages) and the declared row (progress_stages
// present) emit per-cell DOM elements — the F3b future per-cell
// attachment surface. Counting <span class="progress-cell ...">
// occurrences per branch is the structural check.
func TestRenderProgressRowCellDOMPreservedAcrossBranches(t *testing.T) {
	defaultRow := renderTree(t, buildTree([]parsedDoc{
		{Slug: "alpha", Status: "In progress"},
	}, nil))
	if got := strings.Count(defaultRow, `<span class="progress-cell `); got != 4 {
		t.Errorf("default row must emit 4 per-cell DOM elements (C-INV-1); got %d in:\n%s", got, defaultRow)
	}

	declared := renderTree(t, buildTree([]parsedDoc{
		{Slug: "alpha", Status: "In progress", ProgressStages: []string{"Spec", "Render"}},
	}, nil))
	// Declared branch: 1 Drafting cell (carries
	// progress-cell-drafting modifier) + 2 declared cells (bare
	// progress-cell). The C-INV-1 invariant: every cell is its
	// own DOM element regardless of which row drew it.
	if got := draftingCellCount(declared); got != 1 {
		t.Errorf("declared row must carry exactly one Drafting cell; got %d in:\n%s", got, declared)
	}
	if got := declaredCellCount(declared); got != 2 {
		t.Errorf("declared row must carry one cell per declared stage; got %d in:\n%s", got, declared)
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
// C4 carried forward through p3: the progress row renders
// alongside — not in place of — the preserved t2 surfaces
// (Status badge, actor markers). p3 F7 (C1) relocates the long
// description and the related-PR list inside a nested
// body-disclosure <details>; both still render, just behind the
// disclosure gesture. The long description now renders as
// markdown (paragraphs wrap in <p>); the related-PR anchors are
// unchanged because related-PRs are not the markdown body —
// they're authored URLs the existing isURL branch linkifies.
func TestRenderProgressRowCoexistsWithPreservedSurfaces(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "In progress",
			LongDescription: "The long body.",
			RelatedPRs:      []string{"https://github.com/o/r/pull/9"},
			ProgressStages:  []string{"Spec", "Render"}},
	}, map[string][]*ActiveWorkInstance{
		"alpha": {{Actor: "wst-1", Name: "agent-1", Slug: "alpha"}},
	})
	html := renderTree(t, roots)

	for _, want := range []string{
		`<span class="badge status-in-progress">In progress</span>`,
		`<span class="actor-marker">agent-1</span>`,
		// The long description now markdown-renders inside the
		// body-disclosure; goldmark wraps the single paragraph
		// in a <p> tag.
		`<div class="long-desc"><p>The long body.</p>`,
		`<a href="https://github.com/o/r/pull/9">https://github.com/o/r/pull/9</a>`,
		`<div class="progress-row">`,
		`<details class="body-disclosure">`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("preserved/added surface missing %q; html:\n%s", want, html)
		}
	}
}

// TestRenderBodyHeaderOnlyByDefault asserts p3 F7 (parent C1):
// every node box's first view is header-only — no .long-desc
// and no .related-prs render outside the body-disclosure
// wrapper. A node with no body content emits no body-disclosure
// either (the {{if or .LongDescription .RelatedPRs}} gate); a
// node with body content wraps both .long-desc and .related-prs
// inside the .body-disclosure <details>, so neither surfaces in
// the header-only first view.
func TestRenderBodyHeaderOnlyByDefault(t *testing.T) {
	// (1) field-less node: no body-disclosure emitted at all.
	bare := renderTree(t, buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed"},
	}, nil))
	if strings.Contains(bare, `<details class="body-disclosure">`) {
		t.Errorf("field-less node must not emit a body-disclosure; html:\n%s", bare)
	}
	if strings.Contains(bare, `<div class="long-desc">`) ||
		strings.Contains(bare, `<ul class="related-prs">`) {
		t.Errorf("field-less node must emit no body markup; html:\n%s", bare)
	}

	// (2) body-bearing node: body-disclosure wraps both
	// .long-desc and .related-prs; neither escapes outside the
	// disclosure to render in the header-only first view.
	bodied := renderTree(t, buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed",
			LongDescription: "Body.",
			RelatedPRs:      []string{"https://github.com/o/r/pull/1"}},
	}, nil))
	bodyDiscAt := strings.Index(bodied, `<details class="body-disclosure">`)
	if bodyDiscAt < 0 {
		t.Fatalf("body-bearing node must emit a body-disclosure; html:\n%s", bodied)
	}
	bodyDiscClose := strings.Index(bodied[bodyDiscAt:], `</details>`)
	if bodyDiscClose < 0 {
		t.Fatalf("body-disclosure has no closing tag; html:\n%s", bodied)
	}
	disclosed := bodied[bodyDiscAt : bodyDiscAt+bodyDiscClose+len(`</details>`)]
	if !strings.Contains(disclosed, `<div class="long-desc">`) {
		t.Errorf("long-desc must render inside the body-disclosure; disclosed=%q", disclosed)
	}
	if !strings.Contains(disclosed, `<ul class="related-prs">`) {
		t.Errorf("related-prs must render inside the body-disclosure; disclosed=%q", disclosed)
	}
	// And neither must render anywhere OUTSIDE the body-disclosure
	// in the surrounding node markup.
	outside := bodied[:bodyDiscAt] + bodied[bodyDiscAt+bodyDiscClose+len(`</details>`):]
	if strings.Contains(outside, `<div class="long-desc">`) {
		t.Errorf("long-desc must not render outside body-disclosure; outside=%q", outside)
	}
	if strings.Contains(outside, `<ul class="related-prs">`) {
		t.Errorf("related-prs must not render outside body-disclosure; outside=%q", outside)
	}
}

// TestRenderBodyDisclosureStripsAnchors asserts p3 F7 OD3 = I2b:
// when the body disclosure is open, the long description renders
// as markdown with <a> (anchor) tags stripped. Inline link text
// survives as plain text (the Walk continues into the Link node's
// children); autolinks emit their URL as escaped plain text. The
// markdown formatting around the links still renders (paragraph,
// emphasis, etc.) — only the <a> wrappers are suppressed.
func TestRenderBodyDisclosureStripsAnchors(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed",
			LongDescription: "See [the spec](./spec.md) and *details*.\n\nBare URL: <http://example.com/page>"},
	}, nil)
	html := renderTree(t, roots)

	// Find the body-disclosure region to scope assertions there.
	bodyDiscAt := strings.Index(html, `<details class="body-disclosure">`)
	if bodyDiscAt < 0 {
		t.Fatalf("body-disclosure missing; html:\n%s", html)
	}
	bodyDiscEnd := strings.Index(html[bodyDiscAt:], `</details>`)
	if bodyDiscEnd < 0 {
		t.Fatalf("body-disclosure close tag missing; html:\n%s", html)
	}
	disclosed := html[bodyDiscAt : bodyDiscAt+bodyDiscEnd]

	// No <a tags inside the disclosed body. The summary text
	// "Show description" does not contain anchors; the long-desc
	// content is the only candidate. (The related-PRs ul is the
	// separate authored-PR list — also inside body-disclosure but
	// it intentionally linkifies absolute URLs; this test
	// targets the markdown-body anchor strip, so we restrict to
	// the .long-desc element.)
	longDescAt := strings.Index(disclosed, `<div class="long-desc">`)
	if longDescAt < 0 {
		t.Fatalf("long-desc missing inside body-disclosure; disclosed=%q", disclosed)
	}
	longDescEnd := strings.Index(disclosed[longDescAt:], `</div>`)
	if longDescEnd < 0 {
		t.Fatalf("long-desc close tag missing; disclosed=%q", disclosed)
	}
	longDescHTML := disclosed[longDescAt : longDescAt+longDescEnd]
	if strings.Contains(longDescHTML, "<a ") || strings.Contains(longDescHTML, "<a>") {
		t.Errorf("F7 markdown-rendered body must carry no <a> tags; longDescHTML=%q", longDescHTML)
	}
	// Link text "the spec" survives as plain text.
	if !strings.Contains(longDescHTML, "the spec") {
		t.Errorf("link text must survive the anchor strip; longDescHTML=%q", longDescHTML)
	}
	// Surrounding markdown formatting (emphasis) renders.
	if !strings.Contains(longDescHTML, "<em>details</em>") {
		t.Errorf("non-link markdown must still render; longDescHTML=%q", longDescHTML)
	}
	// Autolink emits its URL as text, not <a>.
	if !strings.Contains(longDescHTML, "http://example.com/page") {
		t.Errorf("autolink URL must render as text after strip; longDescHTML=%q", longDescHTML)
	}
}

// TestRenderBodyDisclosureIndependentOfParentCollapse asserts
// p3 F7 (parent C1): the per-node outer <details> (the box
// collapse) and the nested body-disclosure <details> are
// structurally independent — they are two separate DOM elements
// at different nesting depths, each with its own `open` state.
// Opening one cannot toggle the other (the native <details>
// semantics; the test asserts the structural independence the
// observable independence rides on).
func TestRenderBodyDisclosureIndependentOfParentCollapse(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "In progress",
			LongDescription: "Body."},
	}, map[string][]*ActiveWorkInstance{
		"alpha": {{Actor: "wst-1", Name: "agent-1", Slug: "alpha"}},
	})
	html := renderTree(t, roots)

	// Both <details> exist as separate DOM elements.
	outerAt := strings.Index(html, `<details class="box box-root box-leaf"`)
	if outerAt < 0 {
		t.Fatalf("outer per-node <details> missing; html:\n%s", html)
	}
	bodyAt := strings.Index(html, `<details class="body-disclosure">`)
	if bodyAt < 0 {
		t.Fatalf("body-disclosure <details> missing; html:\n%s", html)
	}
	if !(outerAt < bodyAt) {
		t.Errorf("body-disclosure must nest inside the outer <details>; outerAt=%d bodyAt=%d",
			outerAt, bodyAt)
	}
	// The outer renders `open` (active subtree); the body
	// disclosure does NOT carry `open` (it defaults closed —
	// "header-only first view").
	outerOpenAt := strings.Index(html, `<details class="box box-root box-leaf" open>`)
	if outerOpenAt < 0 {
		t.Errorf("outer <details> must render open for active subtree; html:\n%s", html)
	}
	bodyOpenAt := strings.Index(html, `<details class="body-disclosure" open>`)
	if bodyOpenAt >= 0 {
		t.Errorf("body-disclosure must default closed (header-only first view); html:\n%s", html)
	}
}

// TestRenderModeAffordanceMap is the m2 t1 falsifier for the
// "Mode-affordance map drift" m2 Cross-Task Risk. Each row pins
// one production-reachable (NodeType × Status × has-children)
// triple D3 distinguishes: the rendered HTML either carries the
// exact D4 form shape for the target node OR carries no <form>
// element at all (the <form> presence is the unambiguous
// discriminator between positive and negative D3 rows). Renders
// via the standard renderTree helper so the assertion exercises
// the real template path the GET / handler takes.
//
// Coverage per the m2 t1 Validation Gate: each NodeType the
// renderer can produce — root, milestone, task, phase — at the
// Status values D3 enumerates for that NodeType; at least one
// unknown-Status row; at least one Deferred — <reason> row to
// pin the canonical-prefix branch (C4); the task-with-children
// case. The epic NodeType is structurally unreachable in
// production (slugs.parseSegment never assigns NodeTypeEpic;
// epic-rooted docs render as root), so the root row covers that
// case. The phase-with-children case is structurally unreachable
// in well-formed slugs (pN is the terminal segment per the
// grammar) so it carries no row.
func TestRenderModeAffordanceMap(t *testing.T) {
	const (
		planning = "Begin planning"
		impl     = "Begin implementation"
	)
	cases := []struct {
		name     string
		docs     []parsedDoc
		target   string
		wantMode string // "" means: no <form> anywhere in the render
	}{
		// root: parent-shape by D3 — no affordance regardless of
		// Status or has-children.
		{"root In draft", []parsedDoc{
			{Slug: "alpha", Status: "In draft"},
		}, "alpha", ""},
		{"root Proposed", []parsedDoc{
			{Slug: "alpha", Status: "Proposed"},
		}, "alpha", ""},

		// milestone: parent-shape by D3 — no affordance.
		{"milestone Proposed", []parsedDoc{
			{Slug: "alpha", Status: "Proposed"},
			{Slug: "alpha-m1", Status: "Proposed"},
		}, "alpha-m1", ""},

		// task leaf (root-level task, no phase children): the
		// affordance-carrying shape — every D3 Status branch.
		{"task leaf In draft → planning", []parsedDoc{
			{Slug: "alpha", Status: "Proposed"},
			{Slug: "alpha-t1", Status: "In draft"},
		}, "alpha-t1", planning},
		{"task leaf empty Status → planning (no-doc branch)", []parsedDoc{
			{Slug: "alpha", Status: "Proposed"},
			{Slug: "alpha-t1", Status: ""},
		}, "alpha-t1", planning},
		{"task leaf Proposed → implementation", []parsedDoc{
			{Slug: "alpha", Status: "Proposed"},
			{Slug: "alpha-t1", Status: "Proposed"},
		}, "alpha-t1", impl},
		{"task leaf In progress → no affordance", []parsedDoc{
			{Slug: "alpha", Status: "Proposed"},
			{Slug: "alpha-t1", Status: "In progress"},
		}, "alpha-t1", ""},
		{"task leaf Validating → no affordance", []parsedDoc{
			{Slug: "alpha", Status: "Proposed"},
			{Slug: "alpha-t1", Status: "Validating"},
		}, "alpha-t1", ""},
		{"task leaf Landed → no affordance", []parsedDoc{
			{Slug: "alpha", Status: "Proposed"},
			{Slug: "alpha-t1", Status: "Landed"},
		}, "alpha-t1", ""},
		{"task leaf Deferred → no affordance", []parsedDoc{
			{Slug: "alpha", Status: "Proposed"},
			{Slug: "alpha-t1", Status: "Deferred"},
		}, "alpha-t1", ""},
		{"task leaf Deferred — reason → canonical-prefix to no affordance", []parsedDoc{
			{Slug: "alpha", Status: "Proposed"},
			{Slug: "alpha-t1", Status: "Deferred — out of scope"},
		}, "alpha-t1", ""},
		{"task leaf unknown Status → no affordance (C6 graceful fall-through)", []parsedDoc{
			{Slug: "alpha", Status: "Proposed"},
			{Slug: "alpha-t1", Status: "Frobnicating"},
		}, "alpha-t1", ""},

		// task with phase children: parent-shape by D3 — the task
		// has no affordance. The phase children carry Landed so the
		// rendered HTML carries no <form> from any node, making the
		// negative assertion unambiguous.
		{"task with phase children → no affordance on the task", []parsedDoc{
			{Slug: "alpha", Status: "Proposed"},
			{Slug: "alpha-m1", Status: "Proposed"},
			{Slug: "alpha-m1-t1", Status: "Proposed"},
			{Slug: "alpha-m1-t1-p1", Status: "Landed"},
			{Slug: "alpha-m1-t1-p2", Status: "Landed"},
		}, "alpha-m1-t1", ""},

		// phase leaf: same D3 branches the task-leaf row exercises,
		// pinned at the phase NodeType to catch a switch that
		// silently miscategorizes one but not the other.
		{"phase In draft → planning", []parsedDoc{
			{Slug: "alpha", Status: "Proposed"},
			{Slug: "alpha-m1", Status: "Proposed"},
			{Slug: "alpha-m1-t1", Status: "Proposed"},
			{Slug: "alpha-m1-t1-p1", Status: "In draft"},
		}, "alpha-m1-t1-p1", planning},
		{"phase Proposed → implementation", []parsedDoc{
			{Slug: "alpha", Status: "Proposed"},
			{Slug: "alpha-m1", Status: "Proposed"},
			{Slug: "alpha-m1-t1", Status: "Proposed"},
			{Slug: "alpha-m1-t1-p1", Status: "Proposed"},
		}, "alpha-m1-t1-p1", impl},
		{"phase Landed → no affordance", []parsedDoc{
			{Slug: "alpha", Status: "Proposed"},
			{Slug: "alpha-m1", Status: "Proposed"},
			{Slug: "alpha-m1-t1", Status: "Proposed"},
			{Slug: "alpha-m1-t1-p1", Status: "Landed"},
		}, "alpha-m1-t1-p1", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			roots := buildTree(tc.docs, nil)
			html := renderTree(t, roots)

			if tc.wantMode == "" {
				if strings.Contains(html, "<form") {
					t.Errorf("expected no <form> element in render; html:\n%s", html)
				}
				return
			}

			// Positive case: the exact D4-conformant form for the
			// target node must appear in the render. Pinning the
			// full form shape exercises C2 (one value drives both
			// the hidden mode input and the button label) and C3
			// (the D4 form is emitted exactly).
			wantForm := fmt.Sprintf(
				`<form class="affordance-form" method="POST" action="/spawn"><input type="hidden" name="slug" value="%s"><input type="hidden" name="mode" value="%s"><button class="affordance-button" type="submit">%s</button></form>`,
				tc.target, tc.wantMode, tc.wantMode)
			if !strings.Contains(html, wantForm) {
				t.Errorf("missing expected affordance form\n  want: %s\n  got:\n%s", wantForm, html)
			}

			// C5: the form sits inside the node's <summary> region
			// (the always-visible region) and not inside the
			// .box-body div (which native <details> hides on a
			// closed box). Scoped to the form's own location: the
			// next </summary> after the form must precede the next
			// <div class="box-body"> after the form — otherwise
			// the form lives in the body, not the summary.
			formIdx := strings.Index(html, wantForm)
			rest := html[formIdx+len(wantForm):]
			endSummary := strings.Index(rest, `</summary>`)
			nextBody := strings.Index(rest, `<div class="box-body">`)
			if endSummary < 0 || (nextBody >= 0 && nextBody < endSummary) {
				t.Errorf("affordance form must render inside <summary>, before .box-body: next </summary>=%d next box-body=%d (offsets from form end)",
					endSummary, nextBody)
			}
		})
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
