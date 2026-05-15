// Package db provides access to the workstream-tracker SQLite
// database. See design/v0.1-design.md Section 5 for the data model
// and Section 10 for the driver choice (modernc.org/sqlite —
// pure-Go, no CGO).
//
// v0.0 exposes a minimal Open function. Schema initialization and
// migrations land in v0.1 alongside the API implementation.
package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

// Open opens (or creates) the workstream-tracker SQLite database
// at path and verifies the connection.
func Open(path string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite", sqliteDSN(path))
	if err != nil {
		return nil, fmt.Errorf("open sqlite at %q: %w", path, err)
	}
	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ping sqlite at %q: %w", path, err)
	}
	return conn, nil
}

func sqliteDSN(path string) string {
	separator := "?"
	if strings.Contains(path, "?") {
		separator = "&"
	}
	return path + separator + "_pragma=busy_timeout(5000)&_txlock=immediate"
}
