package site

import (
	"testing"
)

func TestIsRootSlug(t *testing.T) {
	cases := []struct {
		slug string
		want bool
	}{
		{"madrona-feedback", true},
		{"docs-canonical-corrections", true},
		{"madrona-feedback-m1", false},
		{"madrona-feedback-m1-t2", false},
		{"madrona-feedback-m1-t2-p3", false},
		{"single", true},
	}
	for _, c := range cases {
		if got := isRootSlug(c.slug); got != c.want {
			t.Errorf("isRootSlug(%q) = %v, want %v", c.slug, got, c.want)
		}
	}
}

func TestNodeTypeFromSlug(t *testing.T) {
	root := "madrona-feedback"
	cases := []struct {
		slug string
		want string
	}{
		{"madrona-feedback", "root"},
		{"madrona-feedback-m1", "milestone"},
		{"madrona-feedback-m1-t2", "task"},
		{"madrona-feedback-m1-t2-p3", "phase"},
	}
	for _, c := range cases {
		if got := nodeTypeFromSlug(c.slug, root); got != c.want {
			t.Errorf("nodeTypeFromSlug(%q, %q) = %q, want %q", c.slug, root, got, c.want)
		}
	}
}

func TestParentSlug(t *testing.T) {
	root := "madrona-feedback"
	cases := []struct {
		slug string
		want string
	}{
		{"madrona-feedback", ""},
		{"madrona-feedback-m1", "madrona-feedback"},
		{"madrona-feedback-m1-t2", "madrona-feedback-m1"},
		{"madrona-feedback-m1-t2-p3", "madrona-feedback-m1-t2"},
	}
	for _, c := range cases {
		if got := parentSlug(c.slug, root); got != c.want {
			t.Errorf("parentSlug(%q, %q) = %q, want %q", c.slug, root, got, c.want)
		}
	}
}

func TestBuildTreeBasic(t *testing.T) {
	docs := []parsedDoc{
		{Slug: "alpha", Status: "In progress"},
		{Slug: "alpha-m1", Status: "Proposed"},
		{Slug: "alpha-m1-t1", Status: "In draft"},
		{Slug: "alpha-m1-t1-p1", Status: "In draft"},
		{Slug: "alpha-m1-t1-p2", Status: "In draft"},
		{Slug: "beta", Status: "Landed"},
	}
	roots := buildTree(docs, nil)

	if len(roots) != 2 {
		t.Fatalf("got %d roots, want 2", len(roots))
	}
	// Alphabetical order.
	if roots[0].Slug != "alpha" || roots[1].Slug != "beta" {
		t.Errorf("roots = %q, %q; want alpha, beta", roots[0].Slug, roots[1].Slug)
	}

	alpha := roots[0]
	if len(alpha.Children) != 1 || alpha.Children[0].Slug != "alpha-m1" {
		t.Fatalf("alpha children: %+v", alpha.Children)
	}
	m1 := alpha.Children[0]
	if len(m1.Children) != 1 || m1.Children[0].Slug != "alpha-m1-t1" {
		t.Fatalf("m1 children: %+v", m1.Children)
	}
	t1 := m1.Children[0]
	if len(t1.Children) != 2 {
		t.Fatalf("t1 children: %+v", t1.Children)
	}
	if t1.Children[0].Slug != "alpha-m1-t1-p1" || t1.Children[1].Slug != "alpha-m1-t1-p2" {
		t.Errorf("phase order = %q, %q; want alphabetical", t1.Children[0].Slug, t1.Children[1].Slug)
	}

	if roots[1].Slug != "beta" || len(roots[1].Children) != 0 {
		t.Errorf("beta should be a leaf root: %+v", roots[1])
	}
}

func TestBuildTreeAttachesActiveWorkInstances(t *testing.T) {
	docs := []parsedDoc{
		{Slug: "alpha", Status: "In progress"},
		{Slug: "alpha-m1", Status: "Proposed"},
	}
	active := map[string][]*ActiveWorkInstance{
		"alpha":         {{Actor: "agent-1"}},
		"alpha-m1":      {{Actor: "agent-2"}, {Actor: "agent-3"}},
		"unknown-slug":  {{Actor: "orphan"}}, // dropped: no matching node
	}
	roots := buildTree(docs, active)

	if len(roots[0].WorkInstances) != 1 || roots[0].WorkInstances[0].Actor != "agent-1" {
		t.Errorf("alpha workinstances: %+v", roots[0].WorkInstances)
	}
	m1 := roots[0].Children[0]
	if len(m1.WorkInstances) != 2 {
		t.Errorf("m1 workinstances: got %d, want 2", len(m1.WorkInstances))
	}
}

func TestStatusClass(t *testing.T) {
	cases := []struct {
		status, want string
	}{
		{"In draft", "in-draft"},
		{"Proposed", "proposed"},
		{"In progress", "in-progress"},
		{"Validating", "validating"},
		{"Validating — prod smoke", "validating"},
		{"Landed", "landed"},
		{"Deferred", "deferred"},
		{"Deferred — moved to next quarter", "deferred"},
		{"", "unknown"},
		{"Almost there", "unknown"},
	}
	for _, c := range cases {
		if got := statusClass(c.status); got != c.want {
			t.Errorf("statusClass(%q) = %q, want %q", c.status, got, c.want)
		}
	}
}
