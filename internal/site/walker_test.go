package site

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// writeDoc writes a markdown file with the given frontmatter.
func writeDoc(t *testing.T, path, slug, status string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	body := "---\nslug: " + slug + "\n"
	if status != "" {
		body += "Status: " + status + "\n"
	}
	body += "---\n# placeholder\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestWalkPlansEmpty(t *testing.T) {
	dir := t.TempDir()
	docs, err := walkPlans(dir)
	if err != nil {
		t.Fatalf("walkPlans: %v", err)
	}
	if len(docs) != 0 {
		t.Errorf("docs = %v, want empty", docs)
	}
}

func TestWalkPlansNonExistent(t *testing.T) {
	docs, err := walkPlans(filepath.Join(t.TempDir(), "does-not-exist"))
	if err != nil {
		t.Fatalf("walkPlans: %v", err)
	}
	if len(docs) != 0 {
		t.Errorf("docs = %v, want empty", docs)
	}
}

func TestWalkPlansFindsRootAndDescendants(t *testing.T) {
	dir := t.TempDir()
	writeDoc(t, filepath.Join(dir, "epic-a", "README.md"), "epic-a", "In progress")
	writeDoc(t, filepath.Join(dir, "epic-a", "m1.md"), "epic-a-m1", "Proposed")
	writeDoc(t, filepath.Join(dir, "epic-a", "m1-t1.md"), "epic-a-m1-t1", "In draft")
	writeDoc(t, filepath.Join(dir, "task-b", "README.md"), "task-b", "Landed")

	docs, err := walkPlans(dir)
	if err != nil {
		t.Fatalf("walkPlans: %v", err)
	}

	got := []string{}
	for _, d := range docs {
		got = append(got, d.Slug)
	}
	sort.Strings(got)

	want := []string{"epic-a", "epic-a-m1", "epic-a-m1-t1", "task-b"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("docs[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestWalkPlansSkipsScopingFolder(t *testing.T) {
	dir := t.TempDir()
	writeDoc(t, filepath.Join(dir, "epic-a", "README.md"), "epic-a", "")
	writeDoc(t, filepath.Join(dir, "epic-a", "scoping", "m1-t1.md"), "scoping-doc-should-not-render", "")

	docs, err := walkPlans(dir)
	if err != nil {
		t.Fatalf("walkPlans: %v", err)
	}
	if len(docs) != 1 {
		t.Fatalf("docs = %v, want only the README", docs)
	}
	if docs[0].Slug != "epic-a" {
		t.Errorf("docs[0].Slug = %q, want epic-a", docs[0].Slug)
	}
}

// TestWalkPlansNestedMilestoneLayout dogfoods the per-milestone
// nesting convention (spec/planning-doc-location.md): a milestone
// doc at m<N>/README.md, its task/phase docs as siblings inside
// m<N>/, and scoping at m<N>/scoping/. The walker is slug-driven
// and recursive, so the discovered doc set must be identical to
// the flat equivalent and the per-milestone scoping/ must be
// skipped just like a root-level one.
func TestWalkPlansNestedMilestoneLayout(t *testing.T) {
	dir := t.TempDir()
	writeDoc(t, filepath.Join(dir, "epic-a", "README.md"), "epic-a", "In progress")
	writeDoc(t, filepath.Join(dir, "epic-a", "m1", "README.md"), "epic-a-m1", "Proposed")
	writeDoc(t, filepath.Join(dir, "epic-a", "m1", "t1-foo.md"), "epic-a-m1-t1", "In draft")
	writeDoc(t, filepath.Join(dir, "epic-a", "m1", "t1-p1-bar.md"), "epic-a-m1-t1-p1", "In draft")
	// Per-milestone scoping must be skipped just like root scoping.
	writeDoc(t, filepath.Join(dir, "epic-a", "m1", "scoping", "t1-p1.md"), "scoping-should-not-render", "")

	docs, err := walkPlans(dir)
	if err != nil {
		t.Fatalf("walkPlans: %v", err)
	}

	got := []string{}
	for _, d := range docs {
		got = append(got, d.Slug)
	}
	sort.Strings(got)

	want := []string{"epic-a", "epic-a-m1", "epic-a-m1-t1", "epic-a-m1-t1-p1"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v (per-milestone scoping must be skipped)", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("docs[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestWalkPlansSkipsTopLevelFiles(t *testing.T) {
	dir := t.TempDir()
	writeDoc(t, filepath.Join(dir, "stray.md"), "stray", "")
	writeDoc(t, filepath.Join(dir, "epic-a", "README.md"), "epic-a", "")

	docs, err := walkPlans(dir)
	if err != nil {
		t.Fatalf("walkPlans: %v", err)
	}
	if len(docs) != 1 {
		t.Errorf("docs = %v, want only epic-a", docs)
	}
}

func TestParsePlanDocShortDescriptionAndBody(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "doc.md")
	content := "---\nslug: epic-a-m1-t3\nStatus: Proposed\nshort_description: Descriptive tree labels\n---\n\n# t3 — Descriptive tree labels\n\nThe long description body.\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	doc, err := parsePlanDoc(path)
	if err != nil {
		t.Fatalf("parsePlanDoc: %v", err)
	}
	if doc.ShortDescription != "Descriptive tree labels" {
		t.Errorf("ShortDescription = %q, want %q", doc.ShortDescription, "Descriptive tree labels")
	}
	want := "# t3 — Descriptive tree labels\n\nThe long description body."
	if doc.LongDescription != want {
		t.Errorf("LongDescription = %q, want %q", doc.LongDescription, want)
	}
}

func TestParsePlanDocNoShortDescriptionNoBody(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "doc.md")
	if err := os.WriteFile(path, []byte("---\nslug: epic-a\nStatus: Landed\n---\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	doc, err := parsePlanDoc(path)
	if err != nil {
		t.Fatalf("parsePlanDoc: %v", err)
	}
	if doc.ShortDescription != "" {
		t.Errorf("ShortDescription = %q, want empty", doc.ShortDescription)
	}
	if doc.LongDescription != "" {
		t.Errorf("LongDescription = %q, want empty", doc.LongDescription)
	}
}

func TestParsePlanDocRelatedPRsPopulated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "doc.md")
	content := "---\nslug: epic-a-m1-t4\nStatus: Proposed\nrelated_prs:\n  - https://github.com/o/r/pull/1\n  - https://github.com/o/r/pull/2\n---\n# body\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	doc, err := parsePlanDoc(path)
	if err != nil {
		t.Fatalf("parsePlanDoc: %v", err)
	}
	want := []string{"https://github.com/o/r/pull/1", "https://github.com/o/r/pull/2"}
	if len(doc.RelatedPRs) != len(want) {
		t.Fatalf("RelatedPRs = %v, want %v", doc.RelatedPRs, want)
	}
	for i := range want {
		if doc.RelatedPRs[i] != want[i] {
			t.Errorf("RelatedPRs[%d] = %q, want %q", i, doc.RelatedPRs[i], want[i])
		}
	}
}

func TestParsePlanDocRelatedPRsAbsent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "doc.md")
	if err := os.WriteFile(path, []byte("---\nslug: epic-a\nStatus: Landed\n---\n# x\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	doc, err := parsePlanDoc(path)
	if err != nil {
		t.Fatalf("parsePlanDoc: %v", err)
	}
	if len(doc.RelatedPRs) != 0 {
		t.Errorf("RelatedPRs = %v, want empty", doc.RelatedPRs)
	}
}

func TestParsePlanDocRelatedPRsDropsNonStringElement(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "doc.md")
	// The middle entry is a YAML integer, not a string: it must be
	// dropped without erroring or skipping the doc.
	content := "---\nslug: epic-a-m1-t4\nrelated_prs:\n  - https://github.com/o/r/pull/1\n  - 42\n  - https://github.com/o/r/pull/3\n---\n# x\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	doc, err := parsePlanDoc(path)
	if err != nil {
		t.Fatalf("parsePlanDoc: %v", err)
	}
	want := []string{"https://github.com/o/r/pull/1", "https://github.com/o/r/pull/3"}
	if len(doc.RelatedPRs) != len(want) {
		t.Fatalf("RelatedPRs = %v, want %v (non-string dropped)", doc.RelatedPRs, want)
	}
	for i := range want {
		if doc.RelatedPRs[i] != want[i] {
			t.Errorf("RelatedPRs[%d] = %q, want %q", i, doc.RelatedPRs[i], want[i])
		}
	}
}

func TestParsePlanDocRelatedPRsNonSequence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "doc.md")
	// A scalar value where a sequence is expected: tolerated as an
	// empty list, no error.
	if err := os.WriteFile(path, []byte("---\nslug: epic-a\nrelated_prs: not-a-list\n---\n# x\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	doc, err := parsePlanDoc(path)
	if err != nil {
		t.Fatalf("parsePlanDoc: %v", err)
	}
	if len(doc.RelatedPRs) != 0 {
		t.Errorf("RelatedPRs = %v, want empty for non-sequence value", doc.RelatedPRs)
	}
}

func TestParsePlanDocProgressStagesPopulated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "doc.md")
	content := "---\nslug: epic-a-m1-t3\nStatus: Proposed\nprogress_stages:\n  - Spec\n  - Parser\n  - Render\n---\n# body\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	doc, err := parsePlanDoc(path)
	if err != nil {
		t.Fatalf("parsePlanDoc: %v", err)
	}
	want := []string{"Spec", "Parser", "Render"}
	if len(doc.ProgressStages) != len(want) {
		t.Fatalf("ProgressStages = %v, want %v", doc.ProgressStages, want)
	}
	for i := range want {
		if doc.ProgressStages[i] != want[i] {
			t.Errorf("ProgressStages[%d] = %q, want %q (document order preserved)", i, doc.ProgressStages[i], want[i])
		}
	}
}

func TestParsePlanDocProgressStagesAbsent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "doc.md")
	// A `slug` + `Status: In draft` stub: the field is absent and
	// the doc must parse with no error and an empty stage list.
	if err := os.WriteFile(path, []byte("---\nslug: epic-a-m1-t3\nStatus: In draft\n---\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	doc, err := parsePlanDoc(path)
	if err != nil {
		t.Fatalf("parsePlanDoc: %v", err)
	}
	if len(doc.ProgressStages) != 0 {
		t.Errorf("ProgressStages = %v, want empty for absent field", doc.ProgressStages)
	}
}

func TestParsePlanDocProgressStagesDropsNonStringElement(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "doc.md")
	// A YAML integer in the middle of the sequence: dropped, never
	// erroring or skipping the doc (the additive-tolerance posture).
	content := "---\nslug: epic-a-m1-t3\nprogress_stages:\n  - Spec\n  - 7\n  - Render\n---\n# x\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	doc, err := parsePlanDoc(path)
	if err != nil {
		t.Fatalf("parsePlanDoc: %v", err)
	}
	want := []string{"Spec", "Render"}
	if len(doc.ProgressStages) != len(want) {
		t.Fatalf("ProgressStages = %v, want %v (non-string dropped)", doc.ProgressStages, want)
	}
	for i := range want {
		if doc.ProgressStages[i] != want[i] {
			t.Errorf("ProgressStages[%d] = %q, want %q", i, doc.ProgressStages[i], want[i])
		}
	}
}

func TestParsePlanDocProgressStagesNonSequence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "doc.md")
	// A scalar where a sequence is expected: tolerated as an empty
	// list, no error, the doc is not skipped.
	if err := os.WriteFile(path, []byte("---\nslug: epic-a\nprogress_stages: not-a-list\n---\n# x\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	doc, err := parsePlanDoc(path)
	if err != nil {
		t.Fatalf("parsePlanDoc: %v", err)
	}
	if len(doc.ProgressStages) != 0 {
		t.Errorf("ProgressStages = %v, want empty for non-sequence value", doc.ProgressStages)
	}
}

func TestParsePlanDocMissingSlug(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "no-slug.md")
	if err := os.WriteFile(path, []byte("---\nStatus: Proposed\n---\n# x\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	_, err := parsePlanDoc(path)
	if err == nil {
		t.Error("expected error for missing slug, got nil")
	}
}
