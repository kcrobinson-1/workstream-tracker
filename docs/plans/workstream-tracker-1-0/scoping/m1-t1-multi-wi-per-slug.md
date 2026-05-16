# Scoping — workstream-tracker-1-0-m1-t1 (Multi-work-instance per slug)

Scoping doc for the first task of
[`m1`](../m1-v0-2.md) under the
[workstream-tracker-1-0 epic](../README.md). Pairs with the task
plan [`m1-t1-multi-wi-per-slug.md`](../m1-t1-multi-wi-per-slug.md).
Per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns," this doc owns the deliberation,
the decisions with their rejected alternatives, the
reality-check inputs, and the plan-structure handoff; it does
**not** restate the plan-owned file inventory, contracts, or
validation surface. No Status field, per the same spec
("Scoping docs do NOT carry a Status field").

## Context

Today a slug can carry exactly one work-instance row, enforced
at the database. The dogfood loop is therefore lossy: a paused-
then-resumed session, or two agents co-working the same node,
either collide or silently overwrite. t1 relaxes that so a slug
can carry N≥1 work-instances over its lifetime — serially and
concurrently — and exposes a create-or-attach registration flow
at a caller-supplied exact slug: a session can register the
slug's *first* work-instance or attach an additional one. This
is the foundation the milestone's t2 (automatic agent
registration) builds on. t2 derives a slug from agent context
(e.g., the plan-file frontmatter; the exact signal is chosen at
t2 drafting) and auto-registers; the first agent on a
freshly-drafted node is the one that creates that node's first
work-instance, so the flow must support first-registration at
an exact slug, not only attaching to a slug that already has a
work-instance. Without that, t2 cannot register the first
agent on any node; with idempotency, auto-registration retries
do not spawn duplicate active rows.

## Reality-check pass (load-bearing claims, verified)

Every claim below was checked against the code at scoping time;
the plan must re-verify them at draft time (see "Reality-check
inputs the plan must verify").

- **The slug uniqueness constraint is a single unique index, and
  `Init` will not relax it on an existing DB.** `internal/db/schema.go:39-40`
  declares a unique index named `ux_work_instances_slug` on
  `work_instances(slug)`.
  `Init` (`internal/db/schema.go:44-49`) runs the whole schema
  string through one `ExecContext` and is idempotent only in the
  create-if-not-exists sense. A DB that already has the unique
  index keeps it: a create-if-not-exists index statement is a
  no-op against an index that already exists, so relaxing the
  constraint on
  the existing dogfood DB requires an explicit drop, not just an
  edit to the create statement. **This is the load-bearing
  finding that shapes Decision 1.**
- **Registration is a strict two-flow dispatch with no
  exact-slug path.** `internal/api/handlers.go:61-71` and
  `:167-191` dispatch on `RegisterRequest.ParentPath`
  (`internal/api/api.go:39-45`): nil → root-create (slug =
  `RootSlug`), non-nil → descendant-create (server-generated
  slug). Root-create rejects a duplicate root slug as
  `errSlugConflict` → HTTP 409 (`handlers.go:80-81,183-191`).
  There is no request shape that registers a WI at a
  caller-supplied exact slug: root-create only accepts a root
  slug (and rejects duplicates), and descendant-create
  *generates* the slug rather than honoring one the caller
  already wrote into a doc's frontmatter. The exact-slug
  create-or-attach flow is net-new, not an extension of an
  existing branch — this is why t1 does not qualify for the
  narrow-surface scoping carve-out.
- **The server never reads the repo.**
  [`design/v0.1-design.md`](../../../../design/v0.1-design.md)
  states the server is repo-blind (the website walks the repo;
  the server only sees API calls and the DB). Consequence: the
  exact-slug flow cannot validate that a caller-supplied slug
  names a real plan-tree doc. "Known to the tracker" can only
  mean "grammar-valid," never "a doc exists for it" — which is
  the load-bearing constraint behind Decision 6.
- **Idempotency on `(slug, actor, active)` is unimplemented.**
  `insertRegister` (`internal/api/handlers.go:160-229`) always
  mints a fresh `uuid` row with `models.StateActive`; there is no
  `(slug, actor, state)` lookup or upsert anywhere. The
  milestone's "registration is idempotent on `(slug, actor,
  active)`" interface is net-new behavior, not a tightening of
  existing logic.
- **The read side is already multi-WI-safe — no site change
  needed.** `loadActiveWorkInstances`
  (`internal/site/site.go:71-90`) selects `slug, actor` for
  active rows and appends into a `map[string][]*ActiveWorkInstance`
  keyed per slug; the template ranges `.WorkInstances`
  (`internal/site/render.go:55-56`). N≥1 per slug already renders.
- **Descendant slug allocation tolerates duplicate slugs.**
  `NextDescendant` (`internal/slugs/slugs.go:101-124`) takes the
  **max** matched position, not a count, over the rows
  `generateDescendantSlug` reads via a `slug` prefix scan
  (`internal/api/slugs.go:32-34`). Duplicate slugs in the result
  set cannot inflate the next position. Allocation is unaffected
  by relaxation.
- **`rootExists` is EXISTS-based and duplicate-insensitive.**
  `internal/api/slugs.go:58-68` is an EXISTS query keyed by
  `slug` equality. It is used both for root-case conflict detection
  (`handlers.go:184-190`) and the descendant-case root precondition
  (`handlers.go:171-177`). Relaxing uniqueness does not change its
  truth value for either caller.

## Decisions made at scoping time

### Decision 1 — Relax the constraint by an inline drop-and-recreate in the schema block, not a migration framework

Shapes considered:

- **1a — Inline drop-then-create-non-unique in the `Init` schema
  string.** Add an explicit drop of the existing unique index
  followed by a non-unique index create, inside the same schema
  string `Init` already runs every startup. Idempotent on both a
  fresh DB (drop is a no-op, non-unique index created) and the
  existing dogfood DB (unique index dropped, non-unique index
  created). **Chosen.**
- **1b — Versioned migration step.** Rejected. The milestone doc's
  Out of Scope explicitly defers this: "Schema migration system.
  v0.2's single schema relaxation inlines into the existing setup;
  a real migration framework lands when a non-trivial schema
  change requires it." Verified by:
  [`m1-v0-2.md`](../m1-v0-2.md) "Out of Scope" → Schema migration
  system bullet.
- **1c — Leave the unique index, dedupe in application code.**
  Rejected: contradicts the end result (the DB itself must permit
  N≥1 rows per slug) and pushes a uniqueness invariant into every
  write site, which is the opposite of the milestone's
  cross-task "backward-compatible API extension" posture.

Verified by: `internal/db/schema.go:39-49` (single
create-if-not-exists block; `Init` idempotent via one
`ExecContext`); [`m1-v0-2.md`](../m1-v0-2.md) Out of Scope.

The non-unique replacement index is retained (not just dropped):
the slug column is still queried on the registration paths —
the descendant-allocation prefix scan (`api/slugs.go:33`) and
the `rootExists` equality lookup (`api/slugs.go:61`) — so
dropping the index without replacement would regress those
lookups. The render path does not query by slug
(`loadActiveWorkInstances` filters by state and buckets in Go,
per the reality-check bullet above), so it is unaffected either
way. "Bans on surface require rendering the consequence" (per
[`task-plan.md`](../../../../spec/planning/task-plan.md))
applies: removing the index outright is a banned surface whose
consequence (slug scans on every registration call) is not
acceptable.

### Decision 2 — Create-or-attach is an additive third flow on `RegisterRequest`, keyed by an explicit exact-slug field

Shapes considered:

- **2a — New optional field on `RegisterRequest`** naming the
  exact slug to register at; when present the server skips slug
  derivation, descendant generation, and the root-conflict
  check, and registers a WI at that slug (creating the slug's
  first WI if none exists, else attaching — see Decision 6).
  Additive: v0.1 callers that never set the field keep the exact
  two-flow behavior. **Chosen (shape).**
- **2b — Overload `ParentPath`/`RootSlug` with a new node-type or
  sentinel.** Rejected: it mutates the meaning of existing fields,
  breaking the "agents using the v0.1 API shape keep functioning"
  preserve clause (`m1-v0-2.md` t1 contract).
- **2c — Separate endpoint.** Rejected for v0.2: a third route
  family is a heavier public-API surface than the milestone's
  additive posture wants, and the dispatch already lives in one
  handler.

Scoping settles that the flow is an additive third branch in
`registerWorkInstance` selected by an explicit exact-slug
signal, parallel to the existing nil/non-nil `ParentPath`
dispatch. The plan resolves the field spelling as `exact_slug`
(named for its semantic: "use this slug verbatim; do not derive
or generate") — this was a "plan-owned Contract" handoff item
and is now closed in the plan, not left open.

Verified by: `internal/api/api.go:39-45` (current
`RegisterRequest`); `internal/api/handlers.go:61-71,167-191`
(two-flow dispatch); [`m1-v0-2.md`](../m1-v0-2.md) t1 "Preserves."

### Decision 3 — A repeat register for an existing active `(slug, actor)` is a no-op that returns the existing work-instance id

Shapes of "idempotent on `(slug, actor, active)`":

- **3a — Return the existing WI id; insert no new row and no new
  `register` event.** Chosen as the contract intent. A second
  register for a `(slug, actor)` that already has an `active` row
  yields the same id, so t2's auto-registration retries (server
  restart, double session launch) cannot fan out duplicate active
  rows.
- **3b — Return the existing id but still append a `register` (or
  heartbeat) event.** Rejected. The plan's Idempotency contract
  resolves the idempotent path as a pure no-op (no new row, no
  new event); emitting an event for a non-action would pollute
  the event log.
- **3c — Always insert a new row.** Rejected: directly contradicts
  the milestone's "registration is idempotent on `(slug, actor,
  active)`" interface and reintroduces the duplicate-active-row
  failure t1 exists to prevent.

The idempotency key is specifically `(slug, actor, state=active)`:
a new register is permitted once the prior WI for that pair has
reached a terminal state (serial reuse), and is permitted
concurrently for a *different* actor (co-working). Only an
already-active same-actor pair collapses.

Verified by: `internal/api/handlers.go:160-229` (no existing
`(slug, actor, state)` lookup; always inserts fresh
uuid+`StateActive`); [`m1-v0-2.md`](../m1-v0-2.md) t1
Interfaces ("registration is idempotent on `(slug, actor,
active)`").

### Decision 4 — Bare root-create still rejects a duplicate root slug; only the explicit exact-slug flow bypasses the conflict

Shapes considered:

- **4a — Remove root-conflict rejection entirely.** Rejected: an
  accidental double root-create (or a v0.1 caller that re-runs)
  would silently attach instead of erroring, losing the safety net
  the milestone's "existing root-create … flows continue to work"
  preserve clause requires.
- **4b — Route every duplicate root-create through the
  exact-slug flow.** Rejected for the same reason: it changes the
  observable behavior of the existing flow for unchanged callers.
- **4c — Only the explicit exact-slug signal (Decision 2)
  bypasses root-conflict; bare root-create keeps its 409.**
  Chosen.

This keeps `rootExists` correct for both its callers without a
split (Decision 5).

Verified by: `internal/api/handlers.go:183-191` (root-case
`errSlugConflict`); [`m1-v0-2.md`](../m1-v0-2.md) t1 "Preserves."

### Decision 5 — `rootExists` is not split

Because Decision 4c keeps the root-conflict semantic, the
overloaded `rootExists` (`internal/api/slugs.go:58-68`) stays
correct for both call sites under relaxation (it is EXISTS-based,
so duplicate rows do not change its result). No rename or split
is in scope; the plan records this as a reality-check input to
re-confirm rather than a change.

Verified by: `internal/api/slugs.go:58-68`;
`internal/api/handlers.go:171-177,184-190`.

### Decision 6 — The exact-slug flow has no prior-WI precondition (create-or-attach), because t2 consumes it for first-registration

An earlier draft of this flow required that a work-instance
already exist for the slug ("attach to existing only"), to
avoid creating "orphan" work-instances. A review surfaced that
this is incompatible with the milestone's t2 interface.

Shapes considered:

- **6a — Require a prior WI row (attach-to-existing only).**
  Rejected. t2 derives a slug from plan-file frontmatter and
  auto-registers on session start; the *first* agent on a
  freshly-drafted node is the one that would create the first
  WI, so a prior-WI precondition makes first-registration
  impossible — a chicken-and-egg. Concrete falsifier in this
  repo: [`m1-v0-2.md`](../m1-v0-2.md) lists
  `workstream-tracker-1-0-m1-t1…-t4` in its Task Status table
  with no WI rows; under a live t2 the agent implementing t1
  itself would derive `workstream-tracker-1-0-m1-t1` and be
  rejected.
- **6b — t2 uses descendant-create for first-registration.**
  Rejected: descendant-create *generates* a slug rather than
  honoring the frontmatter slug, diverging the WI slug from the
  doc's immutable identity (`shared.md` "Plan-doc identity").
  It also pushes the hard problem into t2 instead of solving it
  once in t1.
- **6c — No prior-WI precondition; the exact-slug flow is
  create-or-attach.** Chosen. Validate slug grammar only; if no
  WI exists for the slug, create the first; else attach
  (subject to Decision 3 idempotency). Orphan-ness (a slug with
  no doc) is *not* policed by t1 — the server cannot read the
  repo to verify a doc exists, and unattached work-instances
  are the epic's deferred triage-zone concern.

Verified by: [`m1-v0-2.md`](../m1-v0-2.md) t2 Interfaces
("Consumes t1's create-or-attach flow") and t1 Interfaces;
[`design/v0.1-design.md`](../../../../design/v0.1-design.md)
(server is repo-blind, per the reality-check bullet above);
[`README.md`](../README.md) "Triage zone in 1.0?" open
question (where orphan/unattached WIs are tracked). This
decision drove the matching edit to the milestone's t1/t2
contract wording (attach-to-existing → create-or-attach), made
in the same change per the Plan-to-PR Completion Gate's
"fix the plan first" rule.

## Decisions handed off to plan-drafting (all now resolved)

These were Contract-/Naming-/structure-level calls deferred to
the plan doc. All are now resolved in the plan; recorded here
for traceability, none still open:

- **Exact `RegisterRequest` field name — RESOLVED** as
  `exact_slug` (see Decision 2; plan Naming + Contracts).
- **Idempotent-path event semantics — RESOLVED.** The plan's
  Idempotency contract chooses 3a: the idempotent no-op writes
  **no new row and no new `register` event**. The 3b
  heartbeat-on-idempotent sub-shape is rejected, not deferred.
- **Replacement index name — RESOLVED** as
  `idx_work_instances_slug` (plan Naming), matching the
  existing non-unique-index convention; the dropped index is
  `ux_work_instances_slug`.
- **N=1 vs N≥2 phase decision — RESOLVED: N=1, phase content
  inline.** The surface is bounded and the internal shape is
  sequence-step (schema relaxation → multi-WI insert →
  exact-slug flow → idempotency) with no intermediate point
  that ships independent stakeholder value — the level-picker
  signal for one inline task plan rather than separate phase
  plans.
- **Self-review audit selection — RESOLVED.** The plan locks
  `validation-honesty` and `error-surfacing-user-mutations`
  from
  [`../../../agents/local/self-review-catalog.md`](../../../agents/local/self-review-catalog.md);
  the latter targets the exact-slug flow's malformed-slug
  rejection (the only rejection path after Decision 6 removed
  the prior-WI precondition).

## Reality-check inputs the plan must verify before promotion

Re-confirm at plan-draft time (line numbers are navigation aids;
the symbolic anchors are load-bearing):

- `internal/db/schema.go` — the slug index is still a single
  unique index named `ux_work_instances_slug`, and `Init` still
  runs the schema as one `ExecContext` block (Decision 1's
  drop-then-recreate depends on both).
- `internal/api/handlers.go` `insertRegister` / `registerWorkInstance`
  — the nil/non-nil `ParentPath` dispatch and the `errSlugConflict`
  root path are unchanged (Decisions 2, 4 graft onto this exact
  structure).
- `internal/site/site.go` `loadActiveWorkInstances` — still
  aggregates a per-slug slice (Decision: no site-package change in
  t1; if this drifted, the no-change claim breaks).
- `internal/slugs/slugs.go` `NextDescendant` — still max-based
  over matched positions (descendant allocation stays
  duplicate-tolerant; no slug-package change in t1).
- `internal/api/slugs.go` `rootExists` — still EXISTS-based
  (Decision 5's no-split claim depends on this).

## Plan-structure handoff

The plan doc (N=1 recommended) will own, per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Required and optional sections":

- **Goal + context preamble** — plain-language framing (this
  doc's Context section is the seed; the plan restates it
  briefly, the only intentional overlap).
- **Contracts** — the relaxed schema invariant, the additive
  `RegisterRequest` `exact_slug` field, the create-or-attach
  third-flow behavior (no prior-WI precondition, per Decision
  6), and the `(slug, actor, active)` idempotency rule, each at
  contract altitude.
- **Files to touch** (estimate-labeled) — `internal/db/schema.go`,
  `internal/api/api.go`, `internal/api/handlers.go`, tests under
  `internal/api/` and `internal/db/`; `internal/api/slugs.go` and
  `internal/site/` reviewed-not-touched per Decisions 5 and the
  read-side finding.
- **Cross-Cutting Invariants** — backward-compatible API
  extension (v0.1 callers unaffected); idempotency on `(slug,
  actor, active)`; the walk-on-every-request render path is not
  touched (per `m1-v0-2.md` Cross-Task Invariants).
- **Validation Gate** — the repo's three canonical validation
  commands `go build ./...`, `go vet ./...`, and `go test ./...`
  (Build / Vet / Tests, defined at
  [`docs/dev.md`](../../../dev.md) lines 84-86 and anchored to
  `agents/shared/validation/philosophy.md`), plus a manual
  register matrix: exact-slug register with no prior WI → first
  WI created (t2's case); two different actors at one slug → two
  active WIs; same actor twice at one active slug → one WI
  (idempotent); bare root-create of an existing root slug →
  still HTTP 409; malformed `exact_slug` → 400. `go vet` is
  load-bearing here: the new `exact_slug` JSON struct tag and
  the exact-slug flow's new `writeError` format strings are vet
  surfaces `go build`/`go test` need not catch.
- **Documentation Currency** —
  [`design/v0.1-design.md`](../../../../design/v0.1-design.md)
  §4 (API: the new create-or-attach exact-slug flow) and §5
  (data model: slug no longer unique), updated in the
  implementing PR.
- **Out of Scope, Risk Register, Related Docs** as content
  applies.

This scoping doc deletes in batch with sibling scoping docs at
the milestone-terminal PR per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns."
