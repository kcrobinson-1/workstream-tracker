---
slug: workstream-tracker-1-0-m2-t4
Status: In draft
short_description: Session roster (bound + unbound) with named sessions and a deliberately-unstructured raw-JSON detail view
---

# t4 — Session roster + work-item enrichment

## Status

`In draft`. **History:** two review findings each regressed a
premature `Proposed`.

1. *Decision-completeness.* Two Contracts bullets (the roster's
   metadata read policy and the name/no-name fallback) deferred
   task-level *behavior* decisions to phase-plan drafting, which
   an N ≥ 2 orchestrating plan must own. Resolved concretely in
   Contracts (read policy: register-event metadata baseline with
   the latest later event's metadata overlaid key-by-key,
   per-request; label: reported `name` else the work-instance
   slug, never the `wst-<uuid>` actor, only the *literal
   slug-fallback formatting* render-time-deferred under "Bans on
   surface require rendering the consequence"); scoping
   Open-decisions and the p1/p2 stub Open-HOW reconciled to
   mechanism-only.
2. *Cross-doc currency.* The plan permitted t4 to edit
   `render.go`'s shared `indexData` / `renderIndex` data-path
   plumbing, but the authoritative milestone Cross-Task Invariant
   and the (`Landed`) t1 Region-ownership contract still made any
   t4 `render.go` edit a reviewer-flag violation, with the plan
   only *promising* to reconcile them in the implementing PR. A
   `Proposed` plan resting on an unreconciled cross-doc
   contradiction is drift; per the parent-doc-currency rule the
   reconciliation belongs in this drafting change. The milestone
   Cross-Task Invariant and t1's Region-ownership contract (and
   its Cross-Cutting restatement) are now amended **here** with
   the narrow data-path carve-out (shared `indexData` /
   `renderIndex` plumbing is expected and not reviewer-flag;
   changing a region body other than the task's own, or the shell
   layout / region boundary, stays reviewer-flag); the plan and
   scoping no longer defer that reconciliation.

The load-bearing HOW calls the parent milestone deferred to t4 are
resolved in the scoping doc
([`scoping/t4-session-roster.md`](scoping/t4-session-roster.md),
decisions D1–D5); no "input from prior task" is pending (t1 is
`Landed`; t4 is independent of t2/t3 per the milestone Sequencing
graph). What remains before `Proposed` is re-walking the
[`task-plan.md`](../../../../spec/planning/task-plan.md)
`In draft → Proposed` promotion gate with both findings resolved:
end-to-end + cross-doc coherence (the milestone/t1 carve-out now
makes the permitted `render.go` data-path edit consistent with the
authoritative invariants); Contracts decision-complete (the sole
residual deferral, the literal slug-fallback formatting, is
authorized by "Bans on surface require rendering the
consequence"); the broadened `Verified by:` rule re-confirmed
against the branch; required sections present and estimate-shaped
sections labeled; no implementation prescription (query mechanism
explicitly carved as p2 HOW). The two phase skeleton docs the
`Phase Contracts` section names —
[`t4-p1-bare-roster.md`](t4-p1-bare-roster.md) and
[`t4-p2-enrichment.md`](t4-p2-enrichment.md) — were seeded at an
earlier flip and left in place (no-clobber); only their inherited
Open-HOW framing was corrected to mechanism-only.

This is an **N ≥ 2 task plan** (orchestrating doc). Per-phase HOW
lives in the phase plans seeded at the `Proposed` flip; this doc
owns the task-level contract, the cross-phase invariants, and the
`Phase Contracts` WHAT split.

## Context

Today the workstream-tracker page renders a forest of plan-tree
nodes with raw per-node actor markers, and t1 landed a two-region
shell whose roster region is a deliberate placeholder. A registered
session whose slug does not match a plan-tree node is **silently
dropped** from the page — there is no surface on which an
orphan/typoed-slug session, or any session's reported detail, can
be seen. This task replaces the placeholder with a **session
roster**: a list of every registered active session, bound and
unbound, each entry showing a human-useful session name and an
expandable raw-JSON view of whatever else the session reported.

It is being done now because t1 made the roster region an
independently-owned surface and the milestone committed the
"every session can be accounted for" observability goal to this
milestone; t4 is the surface that delivers it. The conceptual
surfaces touched: the website's roster region and its request-time
data path, the registration client/CLI (it begins reporting richer
data), and contributor-facing docs (the registration handshake and
the design doc). No new database structure is introduced — reported
data rides the existing event-log `metadata`.

## Goal

Opening the page renders, in t1's roster region, a list of every
**active** registered session. Each session appears whether its
slug is **bound** (matches a plan-tree node) or **unbound** (no
matching node — the accepted orphan/typoed-slug residual that the
forest currently drops). Each entry's primary label is a
**human-useful name**, never the internal `wst-<uuid>` actor; each
entry is **expandable** to a deliberately unstructured raw-JSON
view of the session's other reported fields (reported PRs and
anything else). A session that reports only the v0.2 minimum still
lists (enrichment is additive). Registration stays observable
best-effort and opt-in — the roster surfaces sessions that chose to
register and never claims to see all work. The
walk-on-every-request render path is unchanged (no cache, no
file-watch).

## Contracts

Task-level behavior contracts (the durable WHAT this task
guarantees). Per-phase HOW (file inventory, function shapes,
validation specifics) lives in the phase plans; the `Phase
Contracts` section below states each phase's WHAT. Estimate-shaped
sections are labeled per
[`shared.md`](../../../../spec/planning/shared.md) "Plan content is
a mix of rules and estimates."

### Roster membership

- The roster lists **every active work-instance**, independent of
  plan-tree membership: bound (slug matches a walked plan doc) and
  unbound (slug absent from the walked tree) both appear. Terminal
  (completed/abandoned) work-instances do not appear.
  `Verified by:`
  [`loadActiveWorkInstances` in site.go](../../../../internal/site/site.go)
  selects all rows with `state = active` (it does not filter by
  tree membership);
  [`buildTree` in tree.go](../../../../internal/site/tree.go)
  attaches `active[d.Slug]` only for parsed docs — that join, not
  the query, is the current drop site the roster bypasses.
- Each entry is classified **bound** or **unbound** by testing its
  slug against the same walked plan-tree slug set the forest is
  built from, evaluated **per request** (no cached membership).
  `Verified by:`
  [`Server.index` in site.go](../../../../internal/site/site.go)
  walks plans (`walkPlans`) and loads work-instances per HTTP
  request; the roster consumes the same per-request walk.
- The roster does not add a promote-into-tree / dismiss / attach
  affordance for unbound sessions — it lists only. (Milestone Out
  of Scope; reviewer-flag per the milestone Cross-Task Risk.)

### Session identity and naming

- An entry's **primary display label is a reported human name**,
  carried in the session's reported metadata. The internal
  `wst-<uuid>` actor is the identity key and is **never the roster
  label**. This task does not change the actor generator;
  idempotency stays keyed on `(slug, actor)`.
  `Verified by:`
  [`sessionActor` in cmd/workstream-tracker/register.go](../../../../cmd/workstream-tracker/register.go)
  generates `wst-<uuid>`;
  [`activeWorkInstanceID` in handlers.go](../../../../internal/api/handlers.go)
  keys idempotency on `(slug, actor, state=active)` — untouched by
  this task.
- An entry's display label is the reported metadata `name` value
  when present; when a session reports **no name**, the label
  falls back to the work-instance's **slug** — **never** the
  `wst-<uuid>` actor. Every work-instance has a slug (it is the
  registration key); only the optional reported name may be
  absent, so the slug fallback is always available. This
  name-then-slug rule is the task-level contract, fixed here and
  not deferred; the only render-time-deferred detail is the
  *literal formatting* of the slug-fallback label, authorized by
  "Bans on surface require rendering the consequence" and observed
  in p2's Validation Gate.
  `Verified by:`
  [`work_instances` in schema.go](../../../../internal/db/schema.go)
  declares `slug TEXT NOT NULL`, so every active row the roster
  lists has a slug to fall back to.
- The forest's existing raw per-node actor markers are
  **unchanged** by this task (the v0.1 actor-tag no-regress
  invariant). The resulting forest-shows-uuid / roster-shows-name
  inconsistency is intentional and tracked by a new backlog entry
  (see Backlog Impact), not resolved here.
  `Verified by:`
  [`node` template in forest.go](../../../../internal/site/forest.go)
  ranges `.WorkInstances` into `actor-marker` spans rendering raw
  `{{.Actor}}` — this task adds no edit there.

### Reported data and the raw-JSON detail view

- Richer per-session data is **read by joining the event log**, not
  from a new `work_instances` column. The register client/CLI
  begins sending request `metadata`; the roster loader joins active
  work-instances to their events to surface it. No schema change.
  `Verified by:`
  [`events` / `work_instances` in schema.go](../../../../internal/db/schema.go)
  (`events.metadata TEXT` exists; `work_instances` has no detail
  column; schema is CREATE-IF-NOT-EXISTS with no migration runner);
  [`insertRegister` in handlers.go](../../../../internal/api/handlers.go)
  already writes `nullableJSON(req.Metadata)` to `events.metadata`;
  [`registerclient.Register` in client.go](../../../../internal/registerclient/client.go)
  currently sends only `{exact_slug, actor}`.
- The reported data is rendered as a **deliberately unstructured
  raw-JSON view** — schema-loose on purpose; "what is required" is
  intentionally deferred. No task-side schema is imposed on the
  reported shape, and this data is **not** schematized toward t3's
  declared-stages field (the milestone posture-tension invariant).
  The sole conventionally-read key is the optional `name` (used
  only for the entry label per the naming contract above) — a
  display convention, not a required field or a schema; every
  other reported key stays arbitrary and unschematized, and a
  session omitting `name` still lists.
  `Verified by:` the milestone
  [`README.md`](README.md) Cross-Task Invariant "Opposite spec
  postures are intentional — do not homogenize."
- **Metadata read policy (task-level contract, fixed here).** A
  session is represented by its `register` event's metadata as the
  **identity baseline**, with the **latest** later event's
  (heartbeat / state-transition) metadata **overlaid key-by-key**
  when present: a later event's keys win; keys absent from later
  events keep the register value. It is a **per-request** read
  (the walk-on-every-request invariant). The query mechanism that
  realizes this (join shape, event-type filter, latest-event
  selection) is p2 implementation HOW, not a plan deferral.
  `Verified by:`
  [`insertRegister` / `insertHeartbeat` / `insertStateTransition`
  in handlers.go](../../../../internal/api/handlers.go) each write
  `nullableJSON(...metadata)` into `events.metadata` tagged by
  `type` (`register` / `heartbeat` / `state_transition`), so the
  baseline-plus-overlay rule is derivable from the event rows;
  [`idx_events_work_instance_id` in schema.go](../../../../internal/db/schema.go)
  indexes the per-request read.
- Reported PRs shown in the roster are what the **session
  reported**, distinct from the forest's `gh`-title PR discovery —
  the two are not conflated.
  `Verified by:`
  [`discoverPRsByTitle` in site.go](../../../../internal/site/site.go)
  is title-matched forest discovery; the roster surfaces
  session-reported fields from event metadata, a separate source.

### Render path and ownership

- The roster **region body** (its template + scoped CSS) lives in
  [`roster.go`](../../../../internal/site/roster.go), replacing the
  t1 placeholder; the **forest region body and the shell layout
  are not changed**. t4 **does** edit the shared request-time data
  path — `site.go` (a roster loader) and `render.go`'s shared
  `indexData` / `renderIndex` (a roster field) — which the
  milestone explicitly named as t4's surface. This relaxes t1's
  "t4 edits only `roster.go`" file-enforcement to "shared
  data-path plumbing through `indexData` is expected; changing a
  *region body* other than the roster's, or the shell layout, is
  still reviewer-flag". The milestone Cross-Task Invariant and
  t1's Region-ownership contract (and its Cross-Cutting
  restatement) are **amended with this carve-out in this drafting
  change** — not deferred to the implementing PR — so no `Proposed`
  cross-doc contradiction stands (see Documentation Currency).
  `Verified by:`
  [`roster.go`](../../../../internal/site/roster.go) owns
  `{{define "roster"}}`;
  [`render.go`](../../../../internal/site/render.go) `indexData`
  carries only `Roots` + `PlansPath` today;
  the milestone [`README.md`](README.md) Cross-Task Decisions names
  "Where richer session data is stored/read (decide when t4
  drafts)" as t4's call.
- The render path stays **walk-on-every-request**: roster data is
  loaded per HTTP request; no caching, file-watch, or in-memory
  build-up is introduced.
  `Verified by:`
  [`Server.index` in site.go](../../../../internal/site/site.go)
  calls `walkPlans` + `loadActiveWorkInstances` per request — the
  roster load joins this same per-request path.

## Phase Contracts

Per-phase **WHAT** (end result, sibling-interface handoff,
preserved behavior); per-phase **HOW** is scoped just-in-time in
each phase plan against merged code. Seeded as skeleton docs by the
PR that flips this task plan to `Proposed`, per
[`shared.md`](../../../../spec/planning/shared.md) "Parent-doc
child contracts." The split is an estimate re-derived via the
PR-count branch test, not a contract (scoping D5).

| Slug | Short | Long (end result + preserves) | Feeds siblings |
|------|-------|-------------------------------|----------------|
| `workstream-tracker-1-0-m2-t4-p1` | Bare bound/unbound roster | The roster region lists every active session, each classified bound or unbound, replacing t1's placeholder, using only existing slug/actor/state data (no event join, no client change). Preserves: forest region + shell unchanged; walk-on-every-request; the v0.1 forest actor markers; a session reporting nothing still lists. | Establishes the roster loader, the shared `indexData` roster field, and the bound/unbound classification p2 enriches. Does not touch the event log or the register client. |
| `workstream-tracker-1-0-m2-t4-p2` | Enrichment + named sessions | The register client/CLI reports metadata (incl. a human name); the roster loader gains the event-log join; roster entries show the reported name (defined non-uuid fallback) and an expandable deliberately-unstructured raw-JSON detail view. Preserves: no schema change (event-log join only); schema-loose posture (not homogenized toward t3); registration stays observable best-effort/opt-in; walk-on-every-request. | Consumes p1's roster loader + region. Closes the milestone's observability goal; triggers the backlog split + new entries (Backlog Impact). |

Cross-phase coordination is via this task plan (the cross-phase
decisions and the invariants below), not pre-locked cross-phase
contracts in each phase plan, per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Cross-PR coordination."

## Cross-Cutting Invariants

Rules ≥ 2 sites must hold simultaneously across this task's
phases/files. Reviewer-flag candidates when any phase drafting
brushes these.

- **Walk-on-every-request.** Every roster data read happens per
  HTTP request; no cache/file-watch/in-memory accumulation is
  added in any phase. (Inherited milestone/epic invariant.)
- **Schema-loose, not homogenized.** The reported-data view stays
  arbitrary JSON; no phase imposes a task-side schema or aligns it
  to t3's declared-stages field. (Milestone posture-tension
  invariant.)
- **v0.1 forest actor markers do not regress.** No phase edits the
  forest `node` actor-marker render; the roster's naming is a new
  surface, not a forest change. (Milestone actor invariant.)
- **Region/shell boundary.** Only the roster region body changes;
  shared data-path plumbing through `indexData` is permitted (D3);
  the forest region body and shell layout are not changed by any
  phase.
- **Additive registration data.** A session that reports no
  metadata/name still lists; enrichment never makes the v0.2
  minimum a non-lister. (Milestone preserve.)

## Files to touch

*Estimate of expected shape per
[`shared.md`](../../../../spec/planning/shared.md) "Plan content is
a mix of rules and estimates"; the per-phase plans carry the
authoritative inventory, refined against merged code.*

**Modify (estimate):**

- `internal/site/site.go` — add a roster loader (active
  work-instances + bound/unbound classification; p2 adds the event
  join). Shared data path.
- `internal/site/render.go` — add a roster field to `indexData`
  and pass it through `renderIndex`. Shared composition (D3).
- `internal/site/roster.go` — replace the placeholder with the
  roster render (p1) + named-session + expandable raw-JSON detail
  (p2). The roster region body.
- `internal/site/roster_test.go` — roster render coverage
  (membership, bound/unbound, name fallback, raw-JSON detail).
- `internal/registerclient/client.go`,
  `cmd/workstream-tracker/register.go` — send reported metadata
  incl. name (`--name` / `WST_NAME`). p2 only.
- `docs/agents/local/session-registration.md`,
  `design/v0.1-design.md` — currency (see Documentation Currency).
  p2 only.

**Intentionally not touched** *(estimate — where we don't expect
changes, not a hard prohibition):*

- `internal/db/schema.go` — no schema change (event-log join).
- `internal/site/forest.go`, `internal/site/tree.go` — forest
  region body / forest join unchanged (the no-regress invariant).
- `internal/api/*` — the register/event metadata path is already
  plumbed; no API change expected.

## Validation Gate

The Go toolchain is the gate (no build/test wrapper exists;
`Verified by:` repo root has no Makefile/justfile, `scripts/` is
only `assemble.sh` — matches t1's confirmed finding). Per-phase
gates are authoritative in the phase plans; the task-terminal gate
(satisfied when p2's implementing PR merges) is:

- `gofmt -l internal cmd` reports no files; `go build ./...`,
  `go vet ./...`, `go test ./...` clean, including new roster
  cases.
- **Manual render observation** (Bans on surface require rendering
  the consequence): against the dogfood tree **and** the
  `demo-workstream` tree plus a deliberately registered
  unbound-slug session, visually confirm: every active session
  appears; bound vs. unbound is distinguishable; named sessions
  show the reported name and no `wst-<uuid>` is shown as a label;
  the no-name fallback renders as a designed state; an entry
  expands to the raw-JSON detail; the forest still renders its
  existing actor markers unchanged; one page scroll preserved.

## Execution Steps

*Estimate of ordering; per-phase plans carry the binding
sequence.* p1 then p2 (p2 depends on p1's roster loader/region).
Each phase: baseline validation → dedicated branch → implement →
self-review audits → final validation incl. manual render
observation → PR prep with `## Estimate Deviations` and the plan
Status flip per the
[`task-plan.md`](../../../../spec/planning/task-plan.md)
Plan-to-PR Completion Gate. The task plan reaches `Landed` when
p2's implementing PR merges (the task-terminal flip co-locates the
milestone t4-row update and the backlog split).

## Self-Review Audits

From
[`self-review-catalog.md`](../../../agents/local/self-review-catalog.md),
mapped to this task's surfaces (per-phase plans bind the per-PR
set):

- **validation-honesty** — the roster render and the
  unbound/name/JSON consequences are "passed" only when observed
  in a real render, not asserted from template source.
- **error-surfacing-user-mutations** — the register client/CLI
  reporting metadata must surface a send failure observably
  (registration stays best-effort/opt-in; a silent drop is the
  defect).
- **readiness-gate-truthfulness** — the per-phase and task-plan
  Status flips and the milestone t4-row advance only when each
  Goal/Contract/Validation item is satisfied or explicitly
  deferred in the plan.
- **rename-aware-diff-classification** — relocating/extending the
  roster template must be hand-classified so "forest unchanged"
  and "shell unchanged" claims are real, not tooling-fooled.

## Out Of Scope

- Triage **action** (promote-into-tree / dismiss / attach for
  unbound sessions) — listing only; deferred past 1.0 (milestone
  Out of Scope; reviewer-flag any phase adding it).
- Changing the `wst-<uuid>` actor generator or humanizing the
  **forest** actor display — captured as a new backlog entry, not
  done here.
- A `work_instances` detail column / any schema change — explicitly
  rejected in favor of the event-log join (scoping D1).
- t2/t3 surfaces (forest internals, declared-stages) — independent;
  not touched.

## Risk Register

- **Per-request events join cost.** Reading metadata via an
  event-log join on every page load grows with heartbeat volume.
  Mitigation: the join is keyed by the indexed
  `idx_events_work_instance_id` and scoped to active
  work-instances; it stays consistent with the existing
  per-request load posture. `Verified by:`
  [`schema.go`](../../../../internal/db/schema.go)
  `idx_events_work_instance_id`. Carried, not blocking.
- **Schema-loose view homogenized under review pressure.** A
  reviewer/phase may "tidy" the raw JSON toward t3's schema.
  Mitigation: the posture-tension invariant names this as the
  defect, not the fix.
- **Backlog split / determinism-workstream ordering.** If the
  determinism workstream graduates the entry before t4 lands, the
  split must compose with that. Mitigation: t4 owns only the new
  observability entry; the determinism entry is reduced
  in-place and stays `Open` regardless of graduation order (see
  Backlog Impact).
- **Name fallback leaks the uuid.** A missing-name path could fall
  back to `wst-<uuid>`. Mitigation: the naming contract bans it and
  the phase Validation Gate observes the no-name render.

## Documentation Currency

Status/contract-bearing docs this task's PRs update:

- Milestone [`README.md`](README.md) — the Task Status t4 row, the
  stub-prose reconciliation, and the two parent-doc deferral-note
  resolutions (Cross-Task Decisions storage call → scoping D1;
  Backlog Impact effect question → scoping D2) are updated **in
  this drafting change** per the plan-drafting-updates-the-parent
  rule; subsequent row values track this plan's Status with the PR
  that performs each flip.
- t1 [`t1-site-skeleton.md`](t1-site-skeleton.md) Region-ownership
  contract + its Cross-Cutting restatement, and the milestone
  Cross-Task Invariant — **amended with the data-path carve-out in
  this drafting change** (not the implementing PR): shared
  `render.go` `indexData` / `renderIndex` plumbing and `site.go`'s
  loader are expected; region-body / shell-layout crossings stay
  reviewer-flag (scoping D3). Editing the `Landed` t1 doc is a
  forward-constraint currency fix against the authoritative
  milestone invariant, not a t1 behavior change.
- [`docs/agents/local/session-registration.md`](../../../agents/local/session-registration.md)
  — the handshake reports a session name; updated in p2's PR.
- [`design/v0.1-design.md`](../../../../design/v0.1-design.md) §5/§7
  — roster surface + the event-log-join read of reported data;
  updated in p2's PR.

## Backlog Impact

Per [`spec/backlog.md`](../../../../spec/backlog.md) effect
taxonomy. The effect decision is recorded here at drafting; the
backlog-file mutations execute in the implementing PR(s) per
[`backlog.md`](../../../../spec/backlog.md) "Plan-tree
relationship."

- [`deterministic-interactive-registration`](../../../backlog.md#deterministic-interactive-registration)
  — **split** (scoping D2). The entry is reduced to the
  **determinism gap only** (registration circularity, sole-consumer
  compensation, the neighborly-events tripwire) and stays `Open` —
  it is the entry the independent determinism workstream
  graduates. The split executes in p2's implementing PR.
- **New entry — observability residual** (e.g.
  `unregistered-work-unobservable`, `Open`). Captures the residual
  t4 narrows: registered-but-unbound work is now surfaced by the
  roster, so the residual is *unregistered* work only; the
  tree-side-heuristic candidate stays deferred under the same
  neighborly-events tripwire. Added in p2's implementing PR.
- **New entry — humanize forest actor display** (`Open`). The
  forest still renders the raw `wst-<uuid>`; the roster now shows
  names. Replace the raw uuid in the forest UX (and/or revisit the
  actor generator). Added in p2's implementing PR; scope-framed,
  not prescribed.

## Related Docs

- [`README.md`](README.md) — parent milestone; t4 Task Contract,
  Cross-Task Invariants/Decisions/Risks this plan honors and whose
  deferral notes it resolves.
- [`scoping/t4-session-roster.md`](scoping/t4-session-roster.md) —
  the scoping doc (decisions D1–D5, open decisions, reality-check
  inputs).
- [`t1-site-skeleton.md`](t1-site-skeleton.md) — the landed shell;
  Region-ownership contract this task reconciles.
- [`t3-doc-declared-stages.md`](t3-doc-declared-stages.md) — the
  opposite spec posture this task must not homogenize toward.
- [`../README.md`](../README.md) — parent epic.
- [`../../../../spec/planning/task-plan.md`](../../../../spec/planning/task-plan.md),
  [`../../../../spec/planning/shared.md`](../../../../spec/planning/shared.md)
  — the rules this plan is structured against.
