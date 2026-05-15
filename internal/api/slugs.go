package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kcrobinson-1/workstream-tracker/internal/slugs"
)

type dbReader interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

// generateDescendantSlug returns the next available slug under
// parent (rootSlug + parentPath). It queries existing
// work-instance slugs that match the prefix + type letter and
// returns one greater than the largest existing direct-child
// position.
//
// rootSlug is the root identifier; parentPath is the slug suffix
// of the parent (empty when the parent is the root); nodeType is
// the descendant's node type (must yield a non-empty letter).
func generateDescendantSlug(ctx context.Context, db dbReader, rootSlug, parentPath, nodeType string) (string, error) {
	prefix, err := slugs.DescendantPrefix(rootSlug, parentPath, nodeType)
	if err != nil {
		return "", err
	}

	rows, err := db.QueryContext(ctx,
		`SELECT slug FROM work_instances WHERE slug LIKE ?`,
		prefix+"%",
	)
	if err != nil {
		return "", fmt.Errorf("query existing slugs: %w", err)
	}
	defer rows.Close()

	var existing []string
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return "", fmt.Errorf("scan slug: %w", err)
		}
		existing = append(existing, slug)
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("iterate slugs: %w", err)
	}

	return slugs.NextDescendant(rootSlug, parentPath, nodeType, existing)
}

// rootExists reports whether a work-instance has been registered
// for the given root slug.
func rootExists(ctx context.Context, db dbReader, rootSlug string) (bool, error) {
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
