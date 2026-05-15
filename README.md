# workstream-tracker

A local tool for cross-agent visibility into parallel planning and
implementation work. The tool reads a structured plan tree from a
contributor's repo, accepts agent registration calls for
work-instances active against plan-tree nodes, and renders a
visualization that joins the two so the contributor can see what's
happening across their parallel work.

## Layout

- **[`design/`](design/)** — vision, design docs, and diagrams for
  the tool itself. Start with [design/vision.md](design/vision.md)
  for the long-form framing or
  [design/v0.1-design.md](design/v0.1-design.md) for the first
  shipping version's spec.
- **[`spec/`](spec/)** — the canonical plan-doc spec consumer
  projects adopt. The structure planning agents follow when
  authoring plan docs the tool visualizes. See
  [spec/planning/shared.md](spec/planning/shared.md) for the
  cross-level rules and [spec/backlog.md](spec/backlog.md) for the
  backlog format.
- **[`cmd/tool/`](cmd/tool/)** — entry point for the local server.
- **[`internal/`](internal/)** — server packages: `api/` (HTTP API
  for agents), `site/` (HTML rendering), `db/` (SQLite access),
  `models/` (shared types).

## Running locally

Requires Go 1.21 or later (tested on 1.26).

```sh
go run ./cmd/tool
```

The server listens on `:8080` by default and creates a SQLite
database at `./workstream-tracker.db` on first run. Visit
[http://localhost:8080/](http://localhost:8080/) to see the
placeholder index page; check `/health` for liveness.

Both are configurable via env vars:

```sh
PORT=9000 DB_PATH=/tmp/wst.db go run ./cmd/tool
```

To build a standalone binary:

```sh
go build -o wst ./cmd/tool
./wst
```

To run the test suite:

```sh
go test ./...
```

## Status

v0.0 — spec is drafted; tool skeleton boots and exposes the API
endpoint shapes (currently 501 stubs) plus a placeholder index
page. SQLite schema (events + work_instances) is applied on
startup. The actual API logic, plan-tree walk, and visualization
land in subsequent milestones; see
[design/v0.1-design.md](design/v0.1-design.md) for the full v0.1
target.
