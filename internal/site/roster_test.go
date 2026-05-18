package site

import (
	"strings"
	"testing"
)

// Roster-region render tests. The roster is a deliberate,
// intentional placeholder in m2 t1 (the bound/unbound session
// roster is t4); this pins that it renders as an observed state
// — a framed panel with a heading and an explanatory line — not
// a blank gap, and that it does not pre-empt t4 with data-shaped
// content. This is t4's owned test surface.

func TestRenderRosterPlaceholder(t *testing.T) {
	roots := buildTree([]parsedDoc{
		{Slug: "alpha", Status: "Proposed"},
	}, nil)
	html := renderTree(t, roots)

	for _, want := range []string{
		`<section class="roster-panel">`,
		`<h2 class="roster-title">Sessions</h2>`,
		"The session roster lands in a later task.",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("roster placeholder missing %q; html:\n%s", want, html)
		}
	}

	// The placeholder is inside the roster region.
	rosterOpen := strings.Index(html, `<aside class="roster">`)
	panelAt := strings.Index(html, `<section class="roster-panel">`)
	if rosterOpen < 0 || panelAt < 0 || rosterOpen > panelAt {
		t.Errorf("placeholder must render inside the roster region (roster=%d panel=%d)", rosterOpen, panelAt)
	}

	// It must not pre-empt t4 with a roster affordance.
	for _, banned := range []string{"promote", "dismiss", "<button", "<form", "<input"} {
		if strings.Contains(strings.ToLower(html), banned) {
			t.Errorf("roster placeholder must not render a data/affordance element (%q); html:\n%s", banned, html)
		}
	}
}
