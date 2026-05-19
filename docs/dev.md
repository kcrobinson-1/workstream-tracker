# Development Guide

This document is the contributor-workflow source of truth for this
repo. [`../AGENTS.md`](../AGENTS.md) handles agent decision
discipline (pre-edit gates, scope guardrails, stop-and-report
conditions); this doc handles human contributor procedure (local
setup, daily workflow, validation commands, release flow). The
two are jointly authoritative — when they conflict, stop and
report rather than picking a side.

## Purpose

How a contributor works in this repository today — local setup,
validation, release flow. System vision lives in
[`../design/vision.md`](../design/vision.md); the first shipping
version's spec in [`../design/v0.1-design.md`](../design/v0.1-design.md);
the canonical plan-doc spec consumer projects adopt in
[`../spec/`](../spec/); agent decision discipline in
[`../AGENTS.md`](../AGENTS.md).

## Current Tooling

- **Go (1.21+; tested on 1.26)** — primary language for the
  server and CLI binary.
- **SQLite** — embedded persistence; schema lives under
  `internal/db/` and is applied on startup via idempotent
  CREATE-IF-NOT-EXISTS.
- **`slog`** — structured logging with per-request IDs.

Single-maintainer repo. No CI or build matrix exists yet — this
describes the current scale, not a constraint. Adding CI checks,
or build/lint/doc tooling (including external tools and new
dependencies), is fine when it earns its keep; it just hasn't been
needed yet.

## Repository Shape

- [`cmd/workstream-tracker/`](../cmd/workstream-tracker/) —
  entry point for the local server.
- [`internal/`](../internal/) — server packages.
  - `api/` — HTTP API for agent registration and event recording.
  - `site/` — HTML rendering of the plan-tree forest.
  - `db/` — SQLite access and schema.
  - `models/` — shared types.
- [`spec/`](../spec/) — the canonical plan-doc spec consumer
  projects adopt (also read by this repo's server when parsing
  `docs/plans/`).
- [`design/`](../design/) — vision and per-version design docs.
- [`docs/plans/`](plans/) — actual plan-tree docs.
- [`docs/agents/`](agents/) — vendored and local agent rules.

## Local Workflow

1. Clone the repo. Today the only dependencies are Go and an
   `sqlite3` driver pulled in via `go.mod` — that is the current
   state, not a no-new-dependencies rule.
2. Run the server:

   ```sh
   go run ./cmd/workstream-tracker
   ```

   Defaults: `PORT=8080`, `DB_PATH=./workstream-tracker.db`.
   Both are overridable via env vars:

   ```sh
   PORT=9000 DB_PATH=/tmp/wst.db go run ./cmd/workstream-tracker
   ```

3. Visit [http://localhost:8080/](http://localhost:8080/) for
   the plan-tree forest view; `/health` for liveness.
4. Build a standalone binary:

   ```sh
   go build -o workstream-tracker ./cmd/workstream-tracker
   ```

The SQLite file at `DB_PATH` is created on first run. Delete it
to start from a clean state.

## Registering and completing a session

Work-instance registration is automatic for an interactive
natural-language session: the agent runs the registration
subcommand as part of the session-start narration handshake, and
the symmetric completion subcommand at session end (see
[`../AGENTS.md`](../AGENTS.md) "Session work-instance lifecycle
handshake"). The manual command invocations are the documented
fallback when the handshake did not run or failed:

```sh
go run github.com/kcrobinson-1/workstream-tracker/cmd/workstream-tracker register --slug <canonical-slug>
go run github.com/kcrobinson-1/workstream-tracker/cmd/workstream-tracker complete
```

Invoke through the full module path — the repo ships no installed
`workstream-tracker` binary on `PATH`, and a cwd-relative
`go run ./cmd/...` only resolves from the repo root (it breaks
once a session has `cd`'d into a package directory). The module
path resolves against the current module from any directory in
the checkout. `--slug` (or
`WST_SLUG`) is the canonical plan-doc slug; `--actor` (or
`WST_ACTOR`) defaults to a generated per-session id; `--server`
(or `WST_SERVER`) defaults to `http://localhost:8080`. The command
makes one short, best-effort attempt: on any failure it prints an
explicit line and exits success — it never blocks or fails the
session. Re-running `register` for the same slug and actor
collapses to the existing work-instance, so restart/resume is
safe.

`complete` (or `abandon`, when the work is being dropped rather
than finished) records the terminal state. It resolves the
work-instance id from the receipt `register` cached in this
worktree; pass `--id` (or `WST_WI_ID`) when the receipt was not
cached — e.g. a separate process or worktree. Without a terminal
transition the tree keeps showing the session active.

### Maximizing reliable registration (interactive sessions)

The handshake is best-effort: an agent can get pulled into the
substance of a rich first prompt and skip registration. It never
fails the session, so a missed registration is silent — the
session just never appears in the tree. To make a miss unlikely
and immediately visible, split session start into two turns:

1. **Make the first prompt minimal and non-analytical.** Name
   the target node in plain language and ask only for
   registration plus a plan read — explicitly bounding the
   output. Example:

   ```text
   Register this session for the Tree Rendering task in the
   demo-workstream epic, then read its plan doc. Don't analyze,
   plan, or raise issues yet — stop after the registration
   receipt and a one-line summary of what the plan covers.
   ```

   Name the node, not the slug: resolving the canonical slug
   from the plan doc's `slug:` frontmatter is the agent's job.
2. **Verify the receipt before continuing.** Confirm the output
   echoes a real work-instance id and HTTP status. A prose
   "registered successfully" with no id/status is not
   confirmation — it is the failure this check exists to catch
   (see [`../AGENTS.md`](../AGENTS.md) "Session-start
   work-instance registration").
3. **Send the real instructions on turn 2,** once the receipt
   is confirmed.

This is an interactive-session technique only. Headless or
scheduled sessions have no human between turns to gate on the
receipt; for those the manual command above is the backstop.

## Validation Commands

The cadence and discipline behind validation live in
[`agents/shared/validation/philosophy.md`](agents/shared/validation/philosophy.md):
validation honesty, continuous validation, baseline failure
handling, test boundary discipline. This section names the
specific commands this repo uses to satisfy that discipline.

- **Build:** `go build ./...` — every package compiles.
- **Vet:** `go vet ./...` — static checks pass.
- **Tests:** `go test ./...` — full unit-test suite.

Run `go test ./...` before any push. Run `go build ./...` after
any cross-package refactor. The pre-edit-gate's baseline-
validation step in
[`agents/shared/core/pre-edit-gate.md`](agents/shared/core/pre-edit-gate.md)
fires `go test ./...` against the unmodified tree before edits
begin.

## Self-Review Before Push

Before pushing, walk the audits matching the diff's surfaces:

- **General checklist** (correctness / drift / downstream impact /
  scope discipline) —
  [`agents/shared/self-review/how-to-use.md`](agents/shared/self-review/how-to-use.md).
- **Seeded universal audits** —
  [`agents/shared/self-review/seed-audits/`](agents/shared/self-review/seed-audits/).
- **Project-specific audits** —
  [`agents/local/self-review-catalog.md`](agents/local/self-review-catalog.md);
  add audits as patterns surface (see
  [`agents/shared/self-review/mechanism.md`](agents/shared/self-review/mechanism.md)
  for the trigger-twice add / automated-coverage retire lifecycle).

## Release Flow

PR body shape and commit conventions live in
[`agents/shared/pr-conventions/`](agents/shared/pr-conventions/).

Today the repo is single-maintainer and `main` is not protected.
There is no CI, no merge queue, no tag-based deploy. The release
flow is:

1. Branch off `main` for the change (or commit directly to `main`
   for small contributor-facing edits).
2. Run `go build ./...` and `go test ./...`.
3. Self-review per the section above.
4. Open a PR or commit to `main` directly.

When the repo gains real consumers, this section gets stricter
(CI gates, branch protection, tag-based artifacts). Update this
section in the same change that introduces those gates.

## UI Review Workflow

The universal pattern lives in
[`agents/shared/workflows/ui-review.md`](agents/shared/workflows/ui-review.md).
The UI today is server-rendered HTML at `/` (plan-tree forest)
with no JS framework. Capture flow: take screenshots in a
desktop-width viewport (~1280px wide) against a real `go run`
of the server reading from a populated `DB_PATH` and a populated
`docs/plans/` tree. There is no Playwright or screenshot tooling
yet; capture is manual.

---

_This file is the canonical contributor-workflow source of truth.
[`../AGENTS.md`](../AGENTS.md) references it as the "Development
workflow source of truth." When agent rules conflict with this
doc, stop and report the conflict — the two are jointly
authoritative._
