package db

import (
	"context"
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

func TestWorkInstanceSlugIsUnique(t *testing.T) {
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
			uuid.NewString(), "same-slug", "agent", "active", i, i,
		)
		if i == 0 && err != nil {
			t.Fatalf("first insert: %v", err)
		}
		if i == 1 && err == nil {
			t.Fatal("second insert with duplicate slug succeeded; want unique constraint error")
		}
	}
}
