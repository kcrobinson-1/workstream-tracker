package db

import (
	"context"
	"path/filepath"
	"testing"
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
