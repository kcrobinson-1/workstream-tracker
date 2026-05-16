package slugs

import "testing"

func TestIsValidRoot(t *testing.T) {
	cases := []struct {
		slug string
		want bool
	}{
		{"madrona-feedback", true},
		{"single", true},
		{"m1", true},
		{"", false},
		{"Foo", false},
		{"foo_bar", false},
		{"foo bar", false},
		{"-leading", false},
		{"trailing-", false},
		{"foo--double", false},
	}

	for _, c := range cases {
		if got := IsValidRoot(c.slug); got != c.want {
			t.Errorf("IsValidRoot(%q) = %v, want %v", c.slug, got, c.want)
		}
	}
}

func TestIsWellFormed(t *testing.T) {
	cases := []struct {
		slug string
		want bool
	}{
		{"madrona-feedback", true},
		{"single", true},
		{"workstream-tracker-1-0-m1-t1", true},
		{"epic-m1", true},
		{"epic-m1-t2-p3", true},
		{"", false},
		{"Foo", false},
		{"foo_bar", false},
		{"foo bar", false},
		{"-leading", false},
		{"trailing-", false},
		{"foo--double", false},
		{"m1", false},             // bare position segment, no root word
		{"epic-m1-x9", false},     // non-position segment after position segment
		{"epic-m1-foo", false},    // non-position word after position segment
		{"epic-m1-t2-foo", false}, // trailing non-position word
	}

	for _, c := range cases {
		if got := IsWellFormed(c.slug); got != c.want {
			t.Errorf("IsWellFormed(%q) = %v, want %v", c.slug, got, c.want)
		}
	}
}

func TestValidNodeType(t *testing.T) {
	for _, nodeType := range []string{"epic", "milestone", "task", "phase"} {
		if !ValidNodeType(nodeType) {
			t.Errorf("ValidNodeType(%q) = false, want true", nodeType)
		}
	}
	for _, nodeType := range []string{"", "root", "stage"} {
		if ValidNodeType(nodeType) {
			t.Errorf("ValidNodeType(%q) = true, want false", nodeType)
		}
	}
}

func TestNextDescendant(t *testing.T) {
	existing := []string{
		"epic-x-m1",
		"epic-x-m1-t1",
		"epic-x-m2",
		"epic-x-m10",
		"epic-x-m10-t1",
		"epic-x-mx",
		"other-m11",
	}

	got, err := NextDescendant("epic-x", "", "milestone", existing)
	if err != nil {
		t.Fatalf("NextDescendant: %v", err)
	}
	if got != "epic-x-m11" {
		t.Errorf("NextDescendant milestone = %q, want epic-x-m11", got)
	}

	got, err = NextDescendant("epic-x", "m1", "task", existing)
	if err != nil {
		t.Fatalf("NextDescendant nested: %v", err)
	}
	if got != "epic-x-m1-t2" {
		t.Errorf("NextDescendant task = %q, want epic-x-m1-t2", got)
	}
}

func TestNextDescendantRejectsNonDescendantType(t *testing.T) {
	if _, err := NextDescendant("epic-x", "", "epic", nil); err == nil {
		t.Fatal("NextDescendant epic returned nil error, want error")
	}
}

func TestParse(t *testing.T) {
	cases := []struct {
		raw        string
		wantType   NodeType
		wantParent string
	}{
		{"madrona-feedback", NodeTypeRoot, ""},
		{"madrona-feedback-m1", NodeTypeMilestone, "madrona-feedback"},
		{"madrona-feedback-m1-t2", NodeTypeTask, "madrona-feedback-m1"},
		{"madrona-feedback-m1-t2-p3", NodeTypePhase, "madrona-feedback-m1-t2"},
	}

	for _, c := range cases {
		got, err := Parse(c.raw, "madrona-feedback")
		if err != nil {
			t.Fatalf("Parse(%q): %v", c.raw, err)
		}
		if got.String() != c.raw {
			t.Errorf("String() = %q, want %q", got.String(), c.raw)
		}
		if got.Root() != "madrona-feedback" {
			t.Errorf("Root() = %q, want madrona-feedback", got.Root())
		}
		if got.NodeType() != c.wantType {
			t.Errorf("NodeType() = %q, want %q", got.NodeType(), c.wantType)
		}
		if got.Parent() != c.wantParent {
			t.Errorf("Parent() = %q, want %q", got.Parent(), c.wantParent)
		}
	}
}

func TestParseRejectsMalformedDescendant(t *testing.T) {
	cases := []string{
		"madrona-feedback-mx",
		"madrona-feedback-m1-task",
		"other-m1",
	}

	for _, raw := range cases {
		if _, err := Parse(raw, "madrona-feedback"); err == nil {
			t.Errorf("Parse(%q) returned nil error, want error", raw)
		}
	}
}

func TestIsStructuralRoot(t *testing.T) {
	cases := []struct {
		slug string
		want bool
	}{
		{"madrona-feedback", true},
		{"docs-canonical-corrections", true},
		{"m1", false},
		{"madrona-feedback-m1", false},
		{"madrona-feedback-m1-t2", false},
		{"madrona-feedback-m1-t2-p3", false},
		{"single", true},
	}
	for _, c := range cases {
		if got := IsStructuralRoot(c.slug); got != c.want {
			t.Errorf("IsStructuralRoot(%q) = %v, want %v", c.slug, got, c.want)
		}
	}
}

func TestFindRoot(t *testing.T) {
	roots := map[string]bool{
		"alpha":       true,
		"alpha-extra": true,
		"beta":        true,
	}

	cases := []struct {
		slug string
		want string
	}{
		{"alpha", "alpha"},
		{"alpha-m1-t1", "alpha"},
		{"alpha-extra-m1", "alpha-extra"},
		{"beta-m2", "beta"},
		{"unknown-m1", "unknown-m1"},
	}
	for _, c := range cases {
		if got := FindRoot(c.slug, roots); got != c.want {
			t.Errorf("FindRoot(%q) = %q, want %q", c.slug, got, c.want)
		}
	}
}
