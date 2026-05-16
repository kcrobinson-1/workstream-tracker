package site

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// parsedDoc carries the per-file outputs of the walker: the slug
// (frontmatter authoritative, per spec/planning/shared.md
// "Plan-doc identity (slug)"), the Status value if present, the
// optional short/long descriptions, and the path the doc was found
// at (for diagnostics).
type parsedDoc struct {
	Slug             string
	Status           string
	ShortDescription string
	LongDescription  string
	Path             string
}

// walkPlans walks plansPath as a per-root-folder layout per
// spec/planning-doc-location.md and returns all plan-tree docs
// found. Files with missing or invalid frontmatter are logged and
// skipped, not fatal — a single bad file shouldn't take the whole
// page down.
//
// If plansPath does not exist, walkPlans returns an empty slice
// and no error; the index page handles "no plans yet" gracefully.
func walkPlans(plansPath string) ([]parsedDoc, error) {
	info, err := os.Stat(plansPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("stat %q: %w", plansPath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", plansPath)
	}

	var docs []parsedDoc

	rootEntries, err := os.ReadDir(plansPath)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", plansPath, err)
	}

	for _, rootEntry := range rootEntries {
		if !rootEntry.IsDir() {
			// Per the layout convention, every plan-tree root
			// gets a folder. Top-level files in plansPath are
			// not part of the plan tree.
			continue
		}
		rootDir := filepath.Join(plansPath, rootEntry.Name())

		// Walk the root directory recursively to pick up the
		// root README.md, sibling descendant files, and the
		// scoping/ subfolder (which we currently skip — scoping
		// docs are transient and not rendered in the tree).
		err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if d.Name() == "scoping" {
					return fs.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(d.Name(), ".md") {
				return nil
			}
			doc, perr := parsePlanDoc(path)
			if perr != nil {
				slog.Warn("skipping plan doc", "path", path, "err", perr)
				return nil
			}
			docs = append(docs, doc)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walk %q: %w", rootDir, err)
		}
	}

	return docs, nil
}

// parsePlanDoc reads the file at path and extracts slug + Status
// from its YAML frontmatter. Returns an error if the file can't
// be read, the frontmatter can't be parsed, or no slug is
// present.
func parsePlanDoc(path string) (parsedDoc, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return parsedDoc{}, fmt.Errorf("read: %w", err)
	}

	md := goldmark.New(goldmark.WithExtensions(meta.Meta))
	ctx := parser.NewContext()
	md.Parser().Parse(text.NewReader(source), parser.WithContext(ctx))
	metaData := meta.Get(ctx)

	slug, _ := metaData["slug"].(string)
	if slug == "" {
		return parsedDoc{}, errors.New("missing or empty `slug` in frontmatter")
	}
	status, _ := metaData["Status"].(string)
	shortDescription, _ := metaData["short_description"].(string)

	return parsedDoc{
		Slug:             slug,
		Status:           status,
		ShortDescription: shortDescription,
		LongDescription:  markdownBody(source),
		Path:             path,
	}, nil
}

// markdownBody returns the document body — the content following
// the leading YAML frontmatter block — with surrounding whitespace
// trimmed. An absent or empty body yields an empty string. This
// only runs after the slug check, so a file with no leading
// frontmatter block has already been rejected upstream and never
// reaches here.
func markdownBody(source []byte) string {
	lines := strings.Split(string(source), "\n")
	if len(lines) == 0 || strings.TrimRight(lines[0], "\r") != "---" {
		return ""
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimRight(lines[i], "\r") == "---" {
			return strings.TrimSpace(strings.Join(lines[i+1:], "\n"))
		}
	}
	return ""
}
