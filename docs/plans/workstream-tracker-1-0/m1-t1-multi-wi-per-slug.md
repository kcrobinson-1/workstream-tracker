---
slug: workstream-tracker-1-0-m1-t1
Status: Landed
---

# Task — Multi-work-instance per slug

## Context

Today the tracker permits exactly one work-instance per slug,
enforced by a unique database index. The dogfood loop is lossy
as a result: a session that pauses and later resumes, or two
agents co-working the same plan-tree node, either collide on
the unique constraint or overwrite each other. This task lifts
that limit so a slug can carry N≥1 work-instances across its
lifetime — serially (resume after pause) and concurrently
(co-working actors) — and adds a registration flow that lets a
session register at a caller-supplied exact slug: creating that
slug's first work-instance, or attaching an additional one,
rather than only creating a new root or generating a new
descendant.

This is the first task of [m1](m1-v0-2.md) and the foundation
for m1's t2 (automatic agent registration). t2's model is: an
agent session starts, derives its slug from agent context
(e.g., the plan-file frontmatter it was told to implement — the
exact signal is chosen at t2 drafting per the milestone), and
registers automatically. The *first* agent on a freshly-drafted node is
the one that creates that node's first work-instance, so the
flow t2 consumes must support **first registration against an
exact, already-allocated slug** — not only attaching to a slug
that already has a work-instance. Auto-registration retries on
session start are then safe because a repeated register for the
same active `(slug, actor)` pair is idempotent and does not fan
out duplicate rows. The change
touches the SQLite schema and the registration handler; the
website read path and slug-allocation logic are deliberately
untouched (verified multi-WI-safe — see Files to touch). The
deliberation, rejected alternatives, and reality-check pass
behind these contracts live in the paired scoping doc,
[`scoping/m1-t1-multi-wi-per-slug.md`](scoping/m1-t1-multi-wi-per-slug.md).

This is an N=1 task plan: the work is one coherent outcome with
no intermediate point that ships independent value, so phase
content is absorbed inline per
[`task-plan.md`](../../../spec/planning/task-plan.md) "N = 1
task plan."

## Goal

A slug can carry N≥1 work-instance rows. The registration API
exposes a create-or-attach flow at a caller-supplied exact
slug. A repeated registration for an already-active `(slug,
actor)` pair is idempotent. v0.1 callers and the website are
unaffected. Verifiable when: registering at an exact slug that
has no prior work-instance creates the first one (t2's
first-registration case); two actors registering against one
slug both get active work-instances; the same actor registering
twice against one active slug gets one work-instance; a bare
root-create of an already-registered root still returns HTTP
409; and the website renders multiple actor markers on a single
node.

## Contracts

### Schema — slug is no longer unique

`work_instances.slug` no longer carries a unique constraint.
The schema that `db.Init` applies drops the existing
`ux_work_instances_slug` unique index if present and creates a
non-unique `idx_work_instances_slug` index on
`work_instances(slug)`. `db.Init` stays idempotent on both a
fresh database (drop is a no-op; non-unique index created) and
the existing dogfood database (unique index dropped; non-unique
index created), executed within the single schema-application
step `Init` already performs at startup before the server
serves traffic. The replacement index is retained, not merely
dropped: the slug column is still queried on the registration
paths — the descendant-allocation prefix scan and the
root-existence equality lookup — so an unindexed slug column is
a banned surface whose consequence (per-register slug scans on
every registration call) is not acceptable. The render path is
unaffected either way: it filters work-instances by state and
buckets them by slug in memory, so it neither benefits from nor
needs the slug index.

Verified by: `internal/db/schema.go:39-49` (single
create-if-not-exists block; `Init` idempotent via one
`ExecContext`); `internal/api/slugs.go:33` (prefix scan on
`slug` for descendant allocation) and `internal/api/slugs.go:61`
(`slug` equality lookup in `rootExists`); `internal/site/site.go:71-90`
(`loadActiveWorkInstances` filters by state, not slug).

### Registration request — additive `exact_slug` field

`RegisterRequest` gains one optional field, `exact_slug`
(JSON `exact_slug`, omitted by v0.1 callers). Its presence
selects a third registration flow, parallel to the existing
nil/non-nil `ParentPath` dispatch:

- `exact_slug` absent or empty: behavior is exactly v0.1 —
  nil `ParentPath` is root-create, non-nil is
  descendant-create. No observable change for existing callers,
  including the HTTP 201 response shape.
- `exact_slug` non-empty: the exact-slug flow (below).
  `ParentPath` and `RootSlug`-derived generation are not
  consulted.

Verified by: `internal/api/api.go:39-45` (current
`RegisterRequest`); `internal/api/handlers.go:61-71,167-191`
(two-flow dispatch the new branch grafts onto).

### Exact-slug flow (create-or-attach)

When `exact_slug` is non-empty, the server registers a
work-instance at exactly that slug — no slug derivation, no
descendant generation, no root-conflict check. There is **no
precondition that a work-instance already exists for the
slug**: if none exists, this registration creates the slug's
first work-instance (the case t2 needs for the first agent on a
freshly-drafted node); if one or more already exist, this
attaches an additional work-instance (subject to the
idempotency rule below). Effects:

- The slug must be well-formed under the existing slug grammar
  (root text, or root plus valid `m`/`t`/`p` segments).
  Malformed input is a 400. The server does **not** verify the
  slug corresponds to a real plan-tree doc — it cannot, because
  the server never reads the repo (per
  [`design/v0.1-design.md`](../../../design/v0.1-design.md)).
  A registration at a slug with no corresponding doc (an
  "orphan" work-instance) is permitted by t1 and is the epic's
  deferred triage-zone concern, not t1's to police (see
  [`README.md`](README.md) "Triage zone in 1.0?").
- On success the server appends a `register` event and inserts
  a work-instance row at the given slug with initial state
  active, returning the work-instance id and the slug in the
  existing `RegisterResponse` shape with HTTP 201.

This is the flow m1's t2 consumes for automatic
first-registration; the milestone's t1/t2 contracts name it
create-or-attach (see [`m1-v0-2.md`](m1-v0-2.md) t1 and t2
Interfaces).

### Idempotency on `(slug, actor, active)`

Before inserting a work-instance row in any flow, if an active
work-instance already exists for the exact resolved `(slug,
actor)` pair, no new row and no new `register` event are
written; the existing work-instance id and slug are returned in
the existing `RegisterResponse` shape with HTTP 201. The
idempotency key is `(slug, actor, state = active)`:

- A new registration is permitted once the prior work-instance
  for that pair has reached a terminal state (serial reuse
  after pause/resume).
- A new registration is permitted concurrently for a different
  actor on the same slug (co-working).
- Only an already-active same-actor pair collapses to the
  existing row.

For bare root-create the existing duplicate-root rejection
fires before this check (see Preserves), so idempotency is
observable on the exact-slug flow and on serial resume. The
check-then-insert is serialized by the existing immediate-lock
transaction, so two concurrent registrations for the same pair
cannot both insert.

Verified by: `internal/api/handlers.go:160-229` (no existing
`(slug, actor, state)` lookup; insert path); `internal/db/db.go:32-38`
(`_txlock=immediate` serializes the read+insert in
`insertRegister`).

### Preserves

- v0.1 root-create (nil `ParentPath`, no `exact_slug`) still
  validates root slug format and rejects an already-registered
  root with HTTP 409. Only the explicit exact-slug flow
  bypasses root-conflict.
- v0.1 descendant-create (non-nil `ParentPath`) still generates
  the next position slug; allocation is unaffected because the
  allocator takes the max matched position, not a count, so
  duplicate slugs cannot perturb it.
- The website read path is unchanged: it already aggregates a
  per-slug slice of active work-instances and renders one
  marker per instance.
- The `RegisterRequest`/`RegisterResponse` JSON shapes and the
  HTTP 201 success status are unchanged for all flows.

Verified by: `internal/api/handlers.go:183-191` (root-conflict
path); `internal/slugs/slugs.go:101-124` (max-based
allocation); `internal/site/site.go:71-90` and
`internal/site/render.go:55-56` (per-slug slice + marker
range).

## Cross-Cutting Invariants

- **Backward-compatible API extension.** With `exact_slug`
  absent, every existing request produces the identical row,
  event, slug, and HTTP response it produced in v0.1. The new
  field is the only request-shape change and is optional.
- **At most one active work-instance per `(slug, actor)`.**
  Holds at every WI-insert site (root-create, descendant-create,
  exact-slug). Self-review must walk all three against this
  rule, not just the new exact-slug branch.
- **Render path stays walk-on-every-request.** No caching, file
  watching, or in-memory build-up is introduced; the website
  continues to walk the plan tree and query work-instance state
  per request, per [m1-v0-2.md](m1-v0-2.md) Cross-Task
  Invariants.

## Naming

- `exact_slug` — the new optional `RegisterRequest` JSON
  field and its Go struct field. Named for its semantic ("use
  this slug verbatim; do not derive or generate"), which
  covers both the create and the attach case.
- `idx_work_instances_slug` — the non-unique replacement index,
  named to match the existing non-unique index convention
  (`idx_events_work_instance_id` at
  `internal/db/schema.go:26-27`). The dropped index is
  `ux_work_instances_slug`.

## Files to touch

*Estimate of the expected shape; implementation may revise a
specific entry when a structural call requires it, reported per
the Estimate Deviations rule in the implementing PR.*

**New:** None — the change is modifications plus test
additions; no new source files.

**Modify:**

- `internal/db/schema.go` — drop `ux_work_instances_slug`,
  create non-unique `idx_work_instances_slug`, within the
  schema string `Init` applies.
- `internal/api/api.go` — add the optional `exact_slug` field
  to `RegisterRequest`.
- `internal/api/handlers.go` — add the exact-slug
  (create-or-attach) branch to
  `registerWorkInstance`/`insertRegister`; add the
  `(slug, actor, active)` idempotency check before WI insert.
- `internal/slugs/slugs.go` — add `IsWellFormed`, a standalone
  full-slug grammar validator, and harmonize the existing
  `IsValidRoot`/`IsStructuralRoot` validators with it. Scope grew
  beyond the original `IsWellFormed` addition in response to two
  Codex P2 review comments: `IsWellFormed` now enforces *ordered,
  contiguous, at-most-one* `mN→tN→pN` segments after the root
  (previously each trailing part was only checked individually),
  and the root-portion grammar now forbids any kebab-delimited
  token that is a bare position segment (`mN`/`tN`/`pN`).
  `IsValidRoot` was deliberately narrowed to match — `m1` and
  `m1-foo` are no longer valid roots — bringing it into agreement
  with the pre-existing stricter `IsStructuralRoot`. The
  exact-slug flow's "well-formed under the existing slug grammar"
  contract needs root-independent validation; `slugs.Parse`
  requires a known root the exact-slug flow does not have, so a
  new pure validator was added rather than reusing `Parse`. This
  is an estimate deviation from the original "intentionally not
  touched: `internal/slugs/slugs.go`" line (recorded in the
  implementing PR's Estimate Deviations).
- `internal/api/api_test.go` — exact-slug create, exact-slug
  attach (multi-actor), idempotency (id + no-new-event
  discriminator), serial resume, malformed-slug,
  root-conflict-bypass, and v0.1-unaffected cases.
- `internal/db/db_test.go` — relaxation on a fresh DB and
  relaxation against a database that already has the legacy
  unique index (drop-and-recreate path).
- `internal/slugs/slugs_test.go` — `IsWellFormed` unit cases
  (added alongside the new validator above; same estimate
  deviation), plus expanded `IsValidRoot`/`IsWellFormed` cases
  asserting segment-order enforcement and rejection of
  segment-pattern root tokens (review-driven).

**Intentionally not touched** (verified at scoping and at
implementation):

- `internal/site/site.go`, `internal/site/render.go` — read
  path already multi-WI-safe (re-confirmed:
  `loadActiveWorkInstances` filters by state and buckets per
  slug in memory).
- `internal/api/slugs.go` — EXISTS-based `rootExists` stays
  correct under relaxation; `rootExists` is not split
  (Decision 5 in scoping). `generateDescendantSlug`'s max-based
  allocation in `internal/slugs/slugs.go` (`NextDescendant`) is
  unchanged; only the additive `IsWellFormed` was added to that
  file.

## Validation Gate

- `go build ./...`, `go vet ./...`, and `go test ./...` all
  pass. These are the three canonical validation commands this
  repo defines (Build / Vet / Tests) — verified by
  [`docs/dev.md`](../../dev.md) lines 84-86, anchored to
  `agents/shared/validation/philosophy.md`. `go vet` is
  load-bearing for this diff specifically: vet's struct-tag
  check covers the new `exact_slug` JSON tag and vet's printf
  check covers the exact-slug flow's new `writeError` format
  strings, neither of which `go build` or `go test`
  necessarily surfaces.
- Manual register matrix against a running server:
  - Two different actors register against one existing slug →
    two active work-instances; the website renders two markers
    on that node.
  - The same actor registers twice against one active slug →
    one work-instance; the second response returns the first
    work-instance id (idempotent, no new event).
  - Bare root-create of an already-registered root slug → HTTP
    409 (unchanged).
  - Exact-slug register against a slug with no prior
    work-instance → the slug's first work-instance is created
    (t2's first-registration case).
  - Malformed `exact_slug` (not valid under the slug grammar)
    → 400; no row created.
  - Relaxation applied to a database that already had the unique
    index: a second insert at an existing slug succeeds.

The idempotency case names a discriminator (the returned id
equals the first work-instance's id and no new event row
appears) so a false pass — a second row created but the test
only checking HTTP 201 — cannot slip through.

## Self-Review Audits

Run at commit boundary, drawn from
[`docs/agents/local/self-review-catalog.md`](../../agents/local/self-review-catalog.md):

- **validation-honesty** — the Validation Gate's idempotency
  and multi-actor checks actually exercise the relaxed schema
  and the new branch, not a path that passes regardless.
- **error-surfacing-user-mutations** — the exact-slug flow's
  malformed-slug rejection surfaces as an explicit API error
  matching the existing register error shapes, rather than
  failing silently or 500-ing.

## Documentation Currency

- [`design/v0.1-design.md`](../../../design/v0.1-design.md) §4
  (API) — document the create-or-attach exact-slug flow and the
  `exact_slug` field.
- [`design/v0.1-design.md`](../../../design/v0.1-design.md) §5
  (data model) — record that `work_instances.slug` is no longer
  unique. Both updated in the implementing PR.
- [`spec/planning/shared.md`](../../../spec/planning/shared.md)
  "Plan-doc identity (slug)" — slug-grammar clarification added
  in this PR (review-driven): root slugs may not contain a bare
  position-segment token (`m1` is never a root), and after the
  root, position segments must appear in order and at most once
  each (`mN`, then `tN`, then `pN`). This codifies what
  `IsStructuralRoot` already assumed and what the exact-slug
  "well-formed under the existing slug grammar" contract now
  upholds.

## Risk Register

- **Relaxation does not reach an existing dogfood database.**
  A bare edit to the create statement would leave an existing
  unique index in place (create-if-not-exists is a no-op
  against the existing object). Mitigation: the contract
  requires an explicit drop-then-create-non-unique inside the
  schema step `Init` already runs at startup, before serving;
  the `db_test.go` case exercises the already-has-unique-index
  path.
- **`isUniqueConstraint` slug branch becomes partially dead.**
  `insertRegister` maps a unique-constraint failure to
  `errSlugConflict`; after relaxation the slug can no longer
  trigger it, but the `id` primary key still can. Mitigation:
  the branch is retained for the `id` PK; no behavior change is
  claimed for it, and the plan does not remove it.
- **Concurrent same-pair registrations.** Two simultaneous
  registers for one `(slug, actor)` could both pass the
  idempotency read and both insert. Mitigation: `insertRegister`
  runs under the immediate-lock transaction configured at
  `internal/db/db.go:32-38`, serializing the read and insert;
  no additional locking is in scope.
- **Exact-slug flow trusts the caller; orphan work-instances
  are possible.** Because the server never reads the repo (per
  [`design/v0.1-design.md`](../../../design/v0.1-design.md)),
  the exact-slug flow cannot verify the slug names a real
  plan-tree doc; a caller can register at a grammar-valid slug
  with no corresponding doc. This is accepted by design:
  dropping the prior-WI precondition is required for t2's
  first-registration, and orphan/unattached work-instances are
  the epic's deferred triage-zone concern, not t1's. Mitigation:
  validate slug grammar only; do not add a doc-existence check
  (it would reintroduce the t1/t2 chicken-and-egg). Tracked via
  the [`README.md`](README.md) "Triage zone in 1.0?" open
  question.

## Backlog Impact

No existing backlog entry graduates, is deleted, split, or
shifted by this task. One new entry,
[`stub-children-on-parent-promotion`](../../backlog.md#stub-children-on-parent-promotion),
was captured during this task's review (a capture, not one of
the four effects): it records the parent-promotion stub-seeding
idea and explicitly notes it depends on this task's exact-slug
flow rather than replacing it. A second entry,
[`promotion-gate-explicit-checklist`](../../backlog.md#promotion-gate-explicit-checklist),
was captured during this plan's promotion walk (also a
capture): it records that the spec's promotion gate
under-specifies required-sections / no-implementation-prescription
/ spec-conformance as named steps. Orphan/unattached
work-instances (see Risk Register) remain the epic's deferred
"Triage zone in 1.0?" open question, not a backlog entry.

## Related Docs

- [`m1-v0-2.md`](m1-v0-2.md) — parent milestone; t1 contract
  and Cross-Task Invariants.
- [`scoping/m1-t1-multi-wi-per-slug.md`](scoping/m1-t1-multi-wi-per-slug.md)
  — paired scoping doc (decisions, rejected alternatives,
  reality-check pass).
- [`README.md`](README.md) — parent epic.
- [`design/v0.1-design.md`](../../../design/v0.1-design.md) §4,
  §5 — API and data model t1 builds on.
- [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md)
  — the rules this plan is structured against.
