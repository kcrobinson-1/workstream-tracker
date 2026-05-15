package db

import (
	"context"
	"database/sql"
	"fmt"
)

// schema is the workstream-tracker SQLite schema. v0.1 ships with
// one schema version applied as a single CREATE-IF-NOT-EXISTS
// block; a versioned migration mechanism lands when a backward-
// incompatible schema change becomes necessary.
//
// See design/v0.1-design.md Section 5 for the data model these
// tables encode.
const schema = `
CREATE TABLE IF NOT EXISTS events (
    id               TEXT    PRIMARY KEY,
    work_instance_id TEXT    NOT NULL,
    type             TEXT    NOT NULL,
    payload          TEXT,
    metadata         TEXT,
    received_at      INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_events_work_instance_id
    ON events(work_instance_id);

CREATE TABLE IF NOT EXISTS work_instances (
    id              TEXT    PRIMARY KEY,
    slug            TEXT    NOT NULL,
    actor           TEXT    NOT NULL,
    state           TEXT    NOT NULL,
    created_at      INTEGER NOT NULL,
    last_updated_at INTEGER NOT NULL,
    terminal_at     INTEGER
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_work_instances_slug
    ON work_instances(slug);
`

// Init applies the workstream-tracker schema to db. Idempotent.
func Init(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("init schema: %w", err)
	}
	return nil
}
