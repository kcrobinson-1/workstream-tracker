package site

import (
	"testing"
)

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
		"alpha":        {{Actor: "agent-1"}},
		"alpha-m1":     {{Actor: "agent-2"}, {Actor: "agent-3"}},
		"unknown-slug": {{Actor: "orphan"}}, // dropped: no matching node
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

func TestBuildTreeLabel(t *testing.T) {
	docs := []parsedDoc{
		{Slug: "alpha", ShortDescription: "Workstream tracker"},
		{Slug: "alpha-m1", ShortDescription: "Read experience"},
		{Slug: "alpha-m1-t3", ShortDescription: "Descriptive tree labels"},
		{Slug: "alpha-m1-t3-p2", ShortDescription: "Render"},
		{Slug: "beta"},
		{Slug: "beta-m2"},
		{Slug: "beta-m2-t4"},
	}
	roots := buildTree(docs, nil)

	want := map[string]string{
		"alpha":          "Workstream tracker",
		"alpha-m1":       "Milestone 1: Read experience",
		"alpha-m1-t3":    "Task 3: Descriptive tree labels",
		"alpha-m1-t3-p2": "Phase 2: Render",
		"beta":           "beta",
		"beta-m2":        "m2",
		"beta-m2-t4":     "m2-t4",
	}

	var walk func(n *PlanNode)
	seen := map[string]string{}
	walk = func(n *PlanNode) {
		seen[n.Slug] = n.Label
		for _, c := range n.Children {
			walk(c)
		}
	}
	for _, r := range roots {
		walk(r)
	}

	for slug, wantLabel := range want {
		if seen[slug] != wantLabel {
			t.Errorf("Label[%q] = %q, want %q", slug, seen[slug], wantLabel)
		}
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
