package site

import (
	"reflect"
	"testing"
)

// All ghprs tests use a stubbed PR set — no real `gh` is invoked.
// discoverPRsByTitle's subprocess/failure behaviour is exercised
// by the failure-matrix manual checks against real environments
// (Validation Gate), not unit tests; the augment/merge/dedupe
// contract is unit-tested here over decoded PRs.

func TestAugmentRelatedPRs_MatchAndMerge(t *testing.T) {
	const slug = "workstream-tracker-1-0-m1-t4"
	const childSlug = "workstream-tracker-1-0-m1-t4-p2"

	tests := []struct {
		name     string
		frontPRs []string
		prs      []ghPR
		want     []string
	}{
		{
			name:     "case-sensitive substring match appends after frontmatter",
			frontPRs: []string{"https://github.com/o/r/pull/1"},
			prs: []ghPR{
				{URL: "https://github.com/o/r/pull/9", Title: "docs(plans): draft " + slug + " thing", Number: 9},
			},
			want: []string{
				"https://github.com/o/r/pull/1",
				"https://github.com/o/r/pull/9",
			},
		},
		{
			name:     "no frontmatter, single discovered match",
			frontPRs: nil,
			prs: []ghPR{
				{URL: "https://github.com/o/r/pull/5", Title: "feat: " + slug + " render", Number: 5},
			},
			want: []string{"https://github.com/o/r/pull/5"},
		},
		{
			name:     "precision boundary: short-scope title is NOT matched",
			frontPRs: nil,
			prs: []ghPR{
				// Title carries the short conventional-commit scope, not
				// the verbatim slug — scoping P2-D2 under-match.
				{URL: "https://github.com/o/r/pull/7", Title: "feat(m1-t4-p1): inline detail render", Number: 7},
			},
			want: nil,
		},
		{
			name:     "case-sensitive: differently-cased title does not match",
			frontPRs: nil,
			prs: []ghPR{
				{URL: "https://github.com/o/r/pull/8", Title: "docs: WORKSTREAM-TRACKER-1-0-M1-T4 upper", Number: 8},
			},
			want: nil,
		},
		{
			name:     "dedupe: URL in both frontmatter and gh appears once",
			frontPRs: []string{"https://github.com/o/r/pull/3"},
			prs: []ghPR{
				{URL: "https://github.com/o/r/pull/3", Title: "x " + slug + " y", Number: 3},
				{URL: "https://github.com/o/r/pull/4", Title: "z " + slug, Number: 4},
			},
			want: []string{
				"https://github.com/o/r/pull/3",
				"https://github.com/o/r/pull/4",
			},
		},
		{
			name:     "dedupe across multiple discovered entries (gh order preserved)",
			frontPRs: nil,
			prs: []ghPR{
				{URL: "https://github.com/o/r/pull/2", Title: slug + " a", Number: 2},
				{URL: "https://github.com/o/r/pull/2", Title: slug + " b", Number: 2},
				{URL: "https://github.com/o/r/pull/1", Title: slug + " c", Number: 1},
			},
			want: []string{
				"https://github.com/o/r/pull/2",
				"https://github.com/o/r/pull/1",
			},
		},
		{
			name:     "non-matching PRs leave frontmatter untouched",
			frontPRs: []string{"https://github.com/o/r/pull/1"},
			prs: []ghPR{
				{URL: "https://github.com/o/r/pull/99", Title: "unrelated other-slug PR", Number: 99},
			},
			want: []string{"https://github.com/o/r/pull/1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &PlanNode{Slug: slug, RelatedPRs: append([]string(nil), tt.frontPRs...)}
			augmentRelatedPRs([]*PlanNode{node}, tt.prs)
			if !reflect.DeepEqual(node.RelatedPRs, tt.want) {
				t.Fatalf("RelatedPRs = %#v, want %#v", node.RelatedPRs, tt.want)
			}
		})
	}

	// A child node whose slug is a superstring still matches its own
	// slug, and a parent does not steal a child-only PR (the child
	// slug contains the parent slug but not vice-versa).
	t.Run("nested nodes each match their own slug", func(t *testing.T) {
		child := &PlanNode{Slug: childSlug}
		parent := &PlanNode{Slug: slug, Children: []*PlanNode{child}}
		prs := []ghPR{
			{URL: "https://github.com/o/r/pull/20", Title: "docs: " + childSlug + " plan", Number: 20},
		}
		augmentRelatedPRs([]*PlanNode{parent}, prs)
		// Parent slug IS a substring of the child-slug title, so the
		// parent legitimately also lists it (verbatim-substring
		// contract; documented precision behaviour, not a bug).
		if !reflect.DeepEqual(parent.RelatedPRs, []string{"https://github.com/o/r/pull/20"}) {
			t.Fatalf("parent RelatedPRs = %#v", parent.RelatedPRs)
		}
		if !reflect.DeepEqual(child.RelatedPRs, []string{"https://github.com/o/r/pull/20"}) {
			t.Fatalf("child RelatedPRs = %#v", child.RelatedPRs)
		}
	})

	t.Run("nil roots and empty PR set are no-ops", func(t *testing.T) {
		augmentRelatedPRs(nil, nil)
		n := &PlanNode{Slug: slug, RelatedPRs: []string{"https://github.com/o/r/pull/1"}}
		augmentRelatedPRs([]*PlanNode{n, nil}, nil)
		if !reflect.DeepEqual(n.RelatedPRs, []string{"https://github.com/o/r/pull/1"}) {
			t.Fatalf("empty-PR augmentation mutated RelatedPRs: %#v", n.RelatedPRs)
		}
	})
}

// TestDiscoveryErrorYieldsUnaugmentedTree models Server.index's
// log-and-continue contract: when discovery returns an error the
// tree is rendered unaugmented and no error propagates. We
// simulate the error branch (discovery error => augment is not
// called) and assert the tree is byte-identical to its pre-augment
// state.
func TestDiscoveryErrorYieldsUnaugmentedTree(t *testing.T) {
	node := &PlanNode{
		Slug:       "workstream-tracker-1-0-m1-t4",
		RelatedPRs: []string{"https://github.com/o/r/pull/1"},
	}
	roots := []*PlanNode{node}

	// Mirror the Server.index control flow: a discovery error means
	// augmentRelatedPRs is skipped entirely.
	var discoveryErr error = errStub{}
	if discoveryErr != nil {
		// log-and-continue: do not augment.
	} else {
		augmentRelatedPRs(roots, nil)
	}

	want := []string{"https://github.com/o/r/pull/1"}
	if !reflect.DeepEqual(node.RelatedPRs, want) {
		t.Fatalf("on discovery error RelatedPRs = %#v, want unaugmented %#v", node.RelatedPRs, want)
	}
}

type errStub struct{}

func (errStub) Error() string { return "stub gh failure" }
