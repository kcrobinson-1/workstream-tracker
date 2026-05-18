# Scoping — t4 Session roster + work-item enrichment

Scoping doc for [`t4-session-roster.md`](../t4-session-roster.md)
(`workstream-tracker-1-0-m2-t4`). Transient: deletes in batch with
its sibling scoping docs at the milestone-terminal PR per
[`task-plan.md`](../../../../../spec/planning/task-plan.md) "Scoping
owns / plan owns." Carries no `Status` field (scoping docs are
deliberation, not lifecycle artifacts).

## Context

The workstream-tracker page renders a forest of plan-tree nodes
with per-node actor markers; t1 landed the two-region shell with a
deliberate placeholder in the roster region. t4 replaces that
placeholder with a real **session roster** that lists every
registered active session — including the unbound (orphan /
typoed-slug) sessions the forest join silently drops — and enriches
each entry with reported session data (a human-useful name, the
PRs the session reports, and other reported fields) rendered as a
deliberately unstructured raw-JSON view. This scoping doc records
the decisions taken with the contributor on the load-bearing HOW
calls the parent milestone deferred to t4 drafting; the durable
contract lives in the plan doc, not here.

## Decisions made at scoping time

Each decision carries a `Verified by:` code citation per
[`shared.md`](../../../../../spec/planning/shared.md) "`Verified
by:` annotations on load-bearing claims."

### D1 — Enriched session data is read by joining the event log; no schema change

The roster reads richer per-session data from the existing
`events.metadata` column rather than adding a column to
`work_instances`. The register client/CLI starts populating the
already-plumbed request `metadata`; the roster loader joins active
work-instances to their events to surface it.

Rationale: the write path already exists end-to-end — the API
accepts `RegisterRequest.Metadata` and writes it to
`events.metadata` (and the same on heartbeat / state-transition);
no server storage change is needed, only the client beginning to
send it. The `work_instances`-column alternative requires a guarded
`ALTER TABLE` against the existing dogfood DB (the schema is
CREATE-IF-NOT-EXISTS only, with no migration runner and a logged
migration-deferral risk), introduces a second source of truth, and
buys only a join-free read. The event log stays the single
append-only source of truth, so a future materialized column
remains a linear-cost derive-forward, not a cliff.

`Verified by:`
[`events` table in schema.go](../../../../../internal/db/schema.go)
has a `metadata TEXT` column; `work_instances` does not.
[`db.Init` in schema.go](../../../../../internal/db/schema.go)
applies the schema as a single CREATE-IF-NOT-EXISTS block plus a
hand-rolled `DROP INDEX IF EXISTS` — there is no versioned
migration mechanism.
[`insertRegister` in handlers.go](../../../../../internal/api/handlers.go)
writes `nullableJSON(req.Metadata)` into `events.metadata`;
[`insertHeartbeat` / `insertStateTransition` in handlers.go](../../../../../internal/api/handlers.go)
do the same. [`RegisterRequest.Metadata` in api.go](../../../../../internal/api/api.go)
is `json.RawMessage`.
[`registerclient.Register` in client.go](../../../../../internal/registerclient/client.go)
currently marshals only `{exact_slug, actor}` — the additive
client change is the only producer-side work.

### D2 — Backlog effect on `deterministic-interactive-registration` is a split

The milestone deferred the precise effect (it named only "shift vs.
reference-only"). With the contributor's plan to take
`deterministic-interactive-registration` up as its own independent
workstream, the resolved effect is a **split**:

- `deterministic-interactive-registration` is reduced to the
  **determinism gap only** (registration circularity, sole-consumer
  compensation, the neighborly-events tripwire). It stays `Open`
  and is the entry the independent determinism workstream
  graduates.
- A **new entry** captures the **observability residual**, narrowed
  by t4: registered-but-unbound work is now surfaced by the roster,
  so the residual is *unregistered* work only (a session that never
  registered emits no signal the tool ever sees). Stays `Open`; the
  tree-side-heuristic candidate remains deferred under the same
  tripwire.

Rationale: the split makes the two plans structurally unable to
couple — the determinism workstream graduates a clean
determinism-only entry with no t4-authored prose, and t4 owns only
the new observability entry regardless of graduation order. This
exceeds the milestone's enumerated options (shift / reference-only);
it is authorized because the milestone explicitly deferred the
effect to t4's Backlog Impact and the contributor introduced the
independent-workstream premise the milestone did not have. The
backlog-file mutation executes in the t4 implementing PR(s) per
[`backlog.md`](../../../../../spec/backlog.md) "Plan-tree
relationship" (the PR that lands the plan updates the backlog); the
*effect decision* is recorded now in the plan's Backlog Impact and
the milestone's deferral-note resolution.

`Verified by:`
[`deterministic-interactive-registration` in docs/backlog.md](../../../../backlog.md#deterministic-interactive-registration)
states "This entry tracks both the determinism gap and the
observability gap until resolved" and names the neighborly-events
tripwire and the tree-side-heuristic candidate; the m2 milestone
[`README.md`](../README.md) Backlog Impact records the effect as a
"Named open question: the precise effect (shift vs.
reference-only)" deferred to t4.

### D3 — t4 may edit `render.go` for shared roster data plumbing

The roster needs request-time data (the active-session list with
reported metadata, plus the walked plan-tree slug set to classify
bound vs. unbound). That requires a roster field on the shared
`indexData` and a roster loader in `site.go`. t1's
Region-ownership contract made a `render.go` edit "reviewer-flag";
the contributor's intent is that file separation is a goal, not a
hard restriction. Decision: t4 edits `render.go`'s shared
`indexData` / `renderIndex` plumbing and `site.go`'s data path as
needed; the roster **region body** (template + CSS) stays in
`roster.go`, and the **forest region body** and **shell layout**
are untouched. The milestone Cross-Task Invariant and t1's
Region-ownership contract (and its Cross-Cutting restatement) are
**amended with this carve-out in the drafting change itself** (not
deferred to the implementing PR): shared data-path plumbing through
`render.go`'s `indexData` / `renderIndex` and `site.go`'s loader is
expected; changing a region body other than the task's own, or the
shell layout / region boundary, stays reviewer-flag. This keeps the
plan tree free of a `Proposed` cross-doc contradiction (a review
finding caught the earlier defer-to-implementing-PR framing).

`Verified by:`
[`Region-ownership contract in t1-site-skeleton.md`](../t1-site-skeleton.md)
("t4 edits only `roster.go` … A later-task PR that edits
`render.go` … is reviewer-flag");
[`indexData` / `renderIndex` in render.go](../../../../../internal/site/render.go)
carry only `Roots` + `PlansPath` today;
[`Server.index` in site.go](../../../../../internal/site/site.go)
is the single `/` handler that walks plans and loads
work-instances per request.

### D4 — Sessions get a human-useful name; the raw `wst-<uuid>` actor is never shown in the roster UX

The roster's primary per-session label is a reported human name,
carried in the additive metadata (D1's path). The internal
`wst-<uuid>` actor stays the identity key and is **not changed by
this task** (idempotency stays keyed on `(slug, actor)`; the
per-session actor generator is untouched). The forest's existing
per-node actor markers still render the raw actor this task (the
v0.1 actor-tag no-regress invariant) — a deliberate, surfaced
inconsistency; humanizing the forest actor display is captured as a
new backlog entry, not done here.

`Verified by:`
[`sessionActor` in cmd/workstream-tracker/register.go](../../../../../cmd/workstream-tracker/register.go)
generates the `wst-<uuid>` actor;
[`activeWorkInstanceID` in handlers.go](../../../../../internal/api/handlers.go)
keys idempotency on `(slug, actor, state=active)`;
[`node` template in forest.go](../../../../../internal/site/forest.go)
ranges `.WorkInstances` into `actor-marker` spans rendering the raw
`{{.Actor}}` — the surface the no-regress invariant protects and
the cleanup entry will later address.

### D5 — N ≥ 2: bare bound/unbound roster, then enrichment

The branch-test sketch (see "Plan structure handoff") puts the work
at two phases with independent observable end-states and distinct
validation gates. p1 (bare roster: every active session listed,
bound + unbound, replacing the placeholder, using only
slug/actor/state) is independently shippable accountability value;
p2 (client/CLI enrichment + named-session display + expandable
raw-JSON detail + the event-log join) builds on it. This matches
the milestone's N ≥ 2 estimate and is re-confirmed here via the
PR-count branch test, not inherited as a contract.

`Verified by:` the bare-roster read needs no event join —
[`loadActiveWorkInstances` in site.go](../../../../../internal/site/site.go)
already selects `slug, actor` for active rows; p1 extends only the
selected columns and the bound/unbound classification, while the
event-log join (D1) is isolated to p2's surface.

## Decisions resolved at plan-drafting

These were surfaced here as open and **resolved concretely in the
task plan's Contracts** during plan-drafting (a review finding
caught the first two having been wrongly deferred to phase-plan
drafting; they are now task-level contract, not deferred):

1. **Metadata read policy** — RESOLVED: `register`-event metadata
   as the identity baseline, the latest later event's metadata
   overlaid key-by-key, per-request. The query mechanism is p2
   HOW. See task plan Contracts → "Reported data" ("Metadata read
   policy").
2. **Name + no-name fallback** — RESOLVED: label is the reported
   `name` else the work-instance slug, never the `wst-<uuid>`
   actor; only the literal slug-fallback formatting is
   render-time-deferred under "Bans on surface require rendering
   the consequence." See task plan Contracts → "Session identity
   and naming."
3. **Reported-name wire field** — RESOLVED: the sole
   conventionally-read key is the optional `name`; all other
   reported keys stay arbitrary (schema-loose preserved). Producer
   surface (handshake metadata + a CLI affordance) is p2's
   client/CLI plumbing HOW; the `session-registration.md`
   handshake update is in p2's Documentation Currency.

Genuinely phase-level mechanism (not a plan-level open decision —
the WHAT is already in the task plan's "Roster membership"
contract: per-request classification against the same walked
plan-tree slug set):

- **Unbound classification source (p1 mechanism HOW).** Whether p1
  reuses the slug set the handler already walked vs. a fresh read
  — an implementation choice scoped at p1 drafting, bounded by the
  walk-on-every-request invariant. Not a deferred contract.

## Plan structure handoff

- **Doc shape:** N ≥ 2 task plan — orchestrating
  [`t4-session-roster.md`](../t4-session-roster.md) with a
  `Phase Contracts` section; phase skeleton docs
  (`t4-p1-…`, `t4-p2-…`) are seeded by the PR that flips the task
  plan to `Proposed` per
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  promotion-gate seeding, **not** at `In draft`.
- **Phase split (estimate, re-derived via the branch test, not a
  contract):**
  - **p1 — bare bound/unbound roster.** Roster loader + roster.go
    template replacing the placeholder + shared `indexData` roster
    field + bound/unbound classification + tests. No event join,
    no client change. Observable end-state: every active session
    (bound + unbound) is listed.
  - **p2 — enrichment + named sessions.** Register client/CLI
    sends metadata incl. name; roster loader gains the event-log
    join; roster.go gains named-session display + expandable
    raw-JSON detail; `session-registration.md` and
    `design/v0.1-design.md` currency. Observable end-state: roster
    entries show the reported name and an expandable raw-JSON
    detail view.
- **Cross-phase coordination** is via the task plan (cross-phase
  decisions, the D1–D5 invariants), not pre-locked cross-phase
  contracts in each phase plan, per
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "Cross-PR coordination."
- **Parent-doc updates in the drafting change** (per the
  plan-drafting-updates-the-parent-doc rule): the milestone
  [`README.md`](../README.md) Task Status t4 row, the stub-prose
  reconciliation, and resolution of the two parent-doc deferral
  notes (Cross-Task Decisions storage call → D1; Backlog Impact
  effect question → D2).

## Reality-check inputs

Load-bearing codebase claims the plan's promotion gate must
re-confirm against the branch at drafting time (each stated as a
one-sentence falsifier):

- *"The register write path for metadata is fully plumbed
  server-side; only the client doesn't send it."* Falsifier: the
  API drops `metadata`, or `events.metadata` doesn't exist.
  Check: [`api.go`](../../../../../internal/api/api.go)
  `RegisterRequest.Metadata`,
  [`handlers.go`](../../../../../internal/api/handlers.go)
  `nullableJSON(req.Metadata)` insert,
  [`schema.go`](../../../../../internal/db/schema.go)
  `events.metadata`.
- *"The forest join is where unbound is dropped, not the query."*
  Falsifier: `loadActiveWorkInstances` filters by tree membership.
  Check: [`site.go`](../../../../../internal/site/site.go)
  `loadActiveWorkInstances` selects all active rows;
  [`tree.go`](../../../../../internal/site/tree.go) `buildTree`
  attaches `active[d.Slug]` only for parsed docs (the drop site).
- *"The roster region body is isolated to `roster.go`; the shell
  composes it."* Falsifier: the roster define lives in
  `render.go`. Check:
  [`roster.go`](../../../../../internal/site/roster.go) owns
  `{{define "roster"}}` / `{{define "roster-style"}}`;
  [`render.go`](../../../../../internal/site/render.go) composes
  `{{template "roster" .}}`.
- *"`indexData` carries only `Roots` + `PlansPath` today, so a
  roster field is a net-additive struct change."* Falsifier:
  `indexData` already has a roster/session field. Check:
  [`render.go`](../../../../../internal/site/render.go)
  `indexData`.
- *"No build/test wrapper exists; the Go toolchain is the gate."*
  Falsifier: a Makefile/justfile/script wraps build+test. Check:
  repo root + `scripts/` (only `assemble.sh`). (Matches t1's
  confirmed finding.)
- *"`render_test.go` / `roster_test.go` give a render-test harness
  to extend."* Falsifier: no render test entry point. Check:
  [`render_test.go`](../../../../../internal/site/render_test.go)
  `renderTree` helper;
  [`roster_test.go`](../../../../../internal/site/roster_test.go)
  placeholder coverage from t1.
- *"A second populated plan tree exists for manual render
  observation."* Falsifier: only the dogfood tree exists. Check:
  `docs/plans/demo-workstream/` is a populated second root.
