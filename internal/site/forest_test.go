package site

import (
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
	// leaf box, not an <li>.
	if !strings.Contains(html, `<details class="box box-root" open>`) {
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
	if !strings.Contains(html, "The session roster lands in a later task.") {
		t.Errorf("roster placeholder must still render when the forest is empty; html:\n%s", html)
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
