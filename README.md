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
- **[`cmd/workstream-tracker/`](cmd/workstream-tracker/)** — entry point for the local server.
- **[`internal/`](internal/)** — server packages: `api/` (HTTP API
  for agents), `site/` (HTML rendering), `db/` (SQLite access),
  `models/` (shared types).

## Running locally

Requires Go 1.21 or later (tested on 1.26).

```sh
go run ./cmd/workstream-tracker
```

The server listens on `:8080` by default and creates a SQLite
database at `./workstream-tracker.db` on first run. Visit
[http://localhost:8080/](http://localhost:8080/) to see the
placeholder index page; check `/health` for liveness.

Both are configurable via env vars:

```sh
PORT=9000 DB_PATH=/tmp/wst.db go run ./cmd/workstream-tracker
```

To build a standalone binary:

```sh
go build -o workstream-tracker ./cmd/workstream-tracker
./workstream-tracker
```

To run the test suite:

```sh
go test ./...
```

## Status

v0.1 core loop is functionally complete:

- The local server reads plan-tree docs from
  `docs/plans/<root-slug>/`, parses YAML frontmatter for `slug`
  and `Status`, builds the tree from the slug hierarchy, and
  renders the forest as HTML with Status-colored badges and
  active-actor markers.
- The API accepts agent registration calls (root or descendant,
  with server-generated descendant slugs) and event recording
  (heartbeat or terminal state transitions to `completed` /
  `abandoned`). Each handler runs the event-log append +
  current-state update inside a single transaction.
- SQLite schema (`events`, `work_instances`) is applied on
  startup via idempotent CREATE-IF-NOT-EXISTS.
- `slog`-structured logging with per-request IDs; `/health`
  endpoint; graceful shutdown on SIGINT/SIGTERM.

Out of scope for v0.1 (deferred to later versions): intent
layer UI, sub-stage cells, triage zone, tier-based sorting,
actor lineage, multi-repository or multi-contributor support,
server-to-agent push. See
[design/v0.1-design.md](design/v0.1-design.md) Section 9 for
the full deferred list.
