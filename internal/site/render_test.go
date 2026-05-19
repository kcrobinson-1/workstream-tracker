package site

import (
	"bytes"
	"strings"
	"testing"
)

// renderTree renders a forest of roots and returns the HTML.
// Shared by the shell, forest, and roster region tests.
func renderTree(t *testing.T, roots []*PlanNode) string {
	t.Helper()
	var buf bytes.Buffer
	if err := renderIndex(&buf, indexData{Roots: roots, PlansPath: "docs/plans"}); err != nil {
		t.Fatalf("renderIndex: %v", err)
	}
	return buf.String()
}

// TestRenderTwoRegionShell pins the m2 t1 two-region shell: the
// shell composes a forest region and a roster region (in that
// order) from their own templates. Region-internal behavior is
// covered by forest_test.go / roster_test.go.
func TestRenderTwoRegionShell(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed"},
	}, nil)
	html := renderTree(t, roots)

	for _, want := range []string{
		`<div class="layout">`,
		`<main class="forest">`,
		`<aside class="roster">`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("two-region shell missing %q; html:\n%s", want, html)
		}
	}
	// Forest region still renders the node (m2 t2: now a nested
	// per-node box, not the flat bullet line). One sanctioned
	// cross-file sub-assertion update — node shape is t2's surface
	// (m2 t2 C6), not a shell-region boundary violation.
	if !strings.Contains(html, `<details class="box box-root box-leaf">`) ||
		!strings.Contains(html, `<span class="label" title="alpha">alpha</span>`) {
		t.Errorf("forest region lost the node render; html:\n%s", html)
	}
	// The forest region opens before the roster region.
	if i, j := strings.Index(html, `class="forest"`), strings.Index(html, `class="roster"`); i < 0 || j < 0 || i > j {
		t.Errorf("forest region must precede roster region (forest=%d roster=%d)", i, j)
	}
}
