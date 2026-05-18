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

// TestRenderNoFieldNodeUnchanged pins the exact node-line markup
// for a node carrying neither a long description nor related PRs.
// The literal below is today's output for that line and its
// children; the conditional detail blocks must add nothing, so a
// regression here means the additive contract was broken. This is
// also the byte-identity falsifier for the region-ownership
// split: relocating the node template into forest.go must not
// change this output.
func TestRenderNoFieldNodeUnchanged(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed"},
		{Slug: "alpha-m1", Status: "Landed"},
	}, nil)
	html := renderTree(t, roots)

	// Exact pre-change node region, captured by rendering the
	// original template (the `  ` after `</li>` is the range
	// trailing indent — preserved verbatim to pin byte-identity).
	const wantNode = "<span class=\"badge status-proposed\">Proposed</span><span class=\"label\" title=\"alpha\">alpha</span>\n" +
		"<ul>\n" +
		"  <li>\n" +
		"<span class=\"badge status-landed\">Landed</span><span class=\"label\" title=\"alpha-m1\">m1</span>\n" +
		"</li>\n" +
		"  \n" +
		"</ul>"
	if !strings.Contains(html, wantNode) {
		t.Errorf("field-less node line changed (additive contract broken).\nwant substring:\n%q\ngot:\n%s", wantNode, html)
	}
	// Guard the emitted detail markup, not the CSS rule names (the
	// stylesheet legitimately defines `.long-desc`/`.related-prs`).
	if strings.Contains(html, `<div class="long-desc">`) ||
		strings.Contains(html, `<ul class="related-prs">`) {
		t.Errorf("field-less render emitted detail markup; html:\n%s", html)
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
