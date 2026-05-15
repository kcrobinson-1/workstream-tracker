package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// rootSlugPattern matches the kebab-case format the spec requires
// for root slugs: lowercase letters and digits, with hyphens
// between segments. See spec/planning/shared.md "Plan-doc identity
// (slug)".
var rootSlugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// nodeTypeLetter returns the slug segment prefix for a node type
// (m, t, p), or "" if the type isn't a descendant type.
//
// Epics and standalone task plans are roots, not descendants —
// they don't carry a position prefix.
func nodeTypeLetter(nodeType string) string {
	switch nodeType {
	case "milestone":
		return "m"
	case "task":
		return "t"
	case "phase":
		return "p"
	}
	return ""
}

// validNodeType reports whether nodeType is one of the recognized
// values. Hierarchy validity (e.g., a phase under a task, not under
// an epic) is the author's concern, not the server's.
func validNodeType(nodeType string) bool {
	switch nodeType {
	case "epic", "milestone", "task", "phase":
		return true
	}
	return false
}

// generateDescendantSlug returns the next available slug under
// parent (rootSlug + parentPath). It queries existing
// work-instance slugs that match the prefix + type letter and
// picks the smallest unused position number starting from 1.
//
// rootSlug is the root identifier; parentPath is the slug suffix
// of the parent (empty when the parent is the root); nodeType is
// the descendant's node type (must yield a non-empty letter).
func generateDescendantSlug(ctx context.Context, db *sql.DB, rootSlug, parentPath, nodeType string) (string, error) {
	letter := nodeTypeLetter(nodeType)
	if letter == "" {
		return "", fmt.Errorf("node_type %q is not a descendant type", nodeType)
	}

	prefix := rootSlug + "-"
	if parentPath != "" {
		prefix += parentPath + "-"
	}
	prefix += letter

	rows, err := db.QueryContext(ctx,
		`SELECT slug FROM work_instances WHERE slug LIKE ?`,
		prefix+"%",
	)
	if err != nil {
		return "", fmt.Errorf("query existing slugs: %w", err)
	}
	defer rows.Close()

	// Match prefix followed by digits and nothing else — excludes
	// deeper descendants (e.g., m1 vs m1-t2) the LIKE would also
	// surface.
	tail := regexp.MustCompile(`^` + regexp.QuoteMeta(prefix) + `(\d+)$`)
	maxN := 0
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return "", fmt.Errorf("scan slug: %w", err)
		}
		m := tail.FindStringSubmatch(slug)
		if m == nil {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		if n > maxN {
			maxN = n
		}
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("iterate slugs: %w", err)
	}

	return prefix + strconv.Itoa(maxN+1), nil
}

// rootExists reports whether a work-instance has been registered
// for the given root slug.
func rootExists(ctx context.Context, db *sql.DB, rootSlug string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM work_instances WHERE slug = ?)`,
		rootSlug,
	).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf("check root existence: %w", err)
	}
	return exists, nil
}
