package db

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestInitCreatesSchema(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()
	if err := Init(ctx, conn); err != nil {
		t.Fatalf("Init: %v", err)
	}

	for _, table := range []string{"events", "work_instances"} {
		var name string
		err := conn.QueryRowContext(ctx,
			`SELECT name FROM sqlite_master WHERE type='table' AND name=?`,
			table,
		).Scan(&name)
		if err != nil {
			t.Errorf("expected table %q to exist: %v", table, err)
		}
	}
}

func TestInitIsIdempotent(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()
	for i := 0; i < 3; i++ {
		if err := Init(ctx, conn); err != nil {
			t.Fatalf("Init run %d: %v", i+1, err)
		}
	}
}

// TestWorkInstanceSlugIsNotUnique verifies the relaxed schema: a
// slug carries N>=1 work-instances. Two rows at the same slug both
// insert on a fresh DB.
func TestWorkInstanceSlugIsNotUnique(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()
	if err := Init(ctx, conn); err != nil {
		t.Fatalf("Init: %v", err)
	}

	for i := 0; i < 2; i++ {
		_, err := conn.ExecContext(ctx,
			`INSERT INTO work_instances (id, slug, actor, state, created_at, last_updated_at, terminal_at)
			 VALUES (?, ?, ?, ?, ?, ?, NULL)`,
			uuid.NewString(), "same-slug", fmt.Sprintf("agent-%d", i), "active", i, i,
		)
		if err != nil {
			t.Fatalf("insert %d at duplicate slug should succeed under relaxed schema: %v", i, err)
		}
	}

	var count int
	if err := conn.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM work_instances WHERE slug = ?`, "same-slug",
	).Scan(&count); err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 2 {
		t.Errorf("work_instances at slug = %d, want 2", count)
	}
}

// TestRelaxationDropsExistingUniqueIndex exercises the Risk
// Register's "relaxation does not reach an existing dogfood
// database" hazard: a DB that already carries the legacy
// ux_work_instances_slug unique index must have it dropped by
// Init so a second insert at an existing slug succeeds. A bare
// create-if-not-exists edit would leave the unique index in place.
func TestRelaxationDropsExistingUniqueIndex(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()

	// Simulate the existing dogfood DB: tables plus the legacy
	// unique index on work_instances(slug).
	if _, err := conn.ExecContext(ctx, `
CREATE TABLE work_instances (
    id              TEXT    PRIMARY KEY,
    slug            TEXT    NOT NULL,
    actor           TEXT    NOT NULL,
    state           TEXT    NOT NULL,
    created_at      INTEGER NOT NULL,
    last_updated_at INTEGER NOT NULL,
    terminal_at     INTEGER
);
CREATE UNIQUE INDEX ux_work_instances_slug ON work_instances(slug);
`); err != nil {
		t.Fatalf("seed legacy schema: %v", err)
	}

	// Sanity: the legacy unique index rejects a duplicate slug.
	if _, err := conn.ExecContext(ctx,
		`INSERT INTO work_instances (id, slug, actor, state, created_at, last_updated_at, terminal_at)
		 VALUES (?, ?, ?, ?, ?, ?, NULL)`,
		uuid.NewString(), "dupe", "a0", "active", 0, 0,
	); err != nil {
		t.Fatalf("first legacy insert: %v", err)
	}
	if _, err := conn.ExecContext(ctx,
		`INSERT INTO work_instances (id, slug, actor, state, created_at, last_updated_at, terminal_at)
		 VALUES (?, ?, ?, ?, ?, ?, NULL)`,
		uuid.NewString(), "dupe", "a1", "active", 1, 1,
	); err == nil {
		t.Fatal("legacy unique index should reject a duplicate slug before Init runs")
	}

	// Init must drop the legacy unique index and create the
	// non-unique replacement.
	if err := Init(ctx, conn); err != nil {
		t.Fatalf("Init against existing-unique-index DB: %v", err)
	}

	var idxCount int
	if err := conn.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name=?`,
		"ux_work_instances_slug",
	).Scan(&idxCount); err != nil {
		t.Fatalf("query legacy index: %v", err)
	}
	if idxCount != 0 {
		t.Errorf("legacy ux_work_instances_slug still present after Init")
	}

	// A second insert at an existing slug now succeeds.
	if _, err := conn.ExecContext(ctx,
		`INSERT INTO work_instances (id, slug, actor, state, created_at, last_updated_at, terminal_at)
		 VALUES (?, ?, ?, ?, ?, ?, NULL)`,
		uuid.NewString(), "dupe", "a1", "active", 2, 2,
	); err != nil {
		t.Fatalf("second insert at existing slug after relaxation should succeed: %v", err)
	}
}
