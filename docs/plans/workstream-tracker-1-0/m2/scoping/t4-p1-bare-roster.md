# Scoping — t4 p1 Bare bound/unbound roster

Scoping doc for [`t4-p1-bare-roster.md`](../t4-p1-bare-roster.md)
(`workstream-tracker-1-0-m2-t4-p1`), phase 1 of the N ≥ 2 task
plan [`t4-session-roster.md`](../t4-session-roster.md). Transient:
this doc lives in the milestone's `m2/scoping/` subfolder (p1 is a
phase under the m2 milestone, **not** a standalone task plan), so
its contents delete **in batch at the milestone-terminal PR** per
[`task-plan.md`](../../../../../spec/planning/task-plan.md) "Path
conventions" — "The `scoping/` subfolder is transient — its
contents delete in batch at the milestone-terminal PR (or the
task-terminal PR for standalone task plans)." Carries no `Status`
field (scoping docs are deliberation, not lifecycle artifacts).

## Context

t1 landed a two-region page shell whose roster region is a
deliberate placeholder. The t4 task plan replaces that placeholder
with a real session roster across two phases. **p1** ships the
bare roster: every active work-instance listed and classified
**bound** (slug matches a walked plan doc) or **unbound** (slug
absent from the walked tree) — the orphan/typoed-slug sessions the
forest join silently drops — using only the existing
slug/actor/state data, with **no event-log join and no
register-client change** (those are p2). The behavior WHAT is
already locked at task-plan level (t4 Contracts → "Roster
membership"); the only genuinely-open p1 question is the pure
mechanism the t4 scoping doc handed to this phase: the data source
for the bound/unbound classification. This scoping doc records
that decision (with rejected alternatives) and the reality-check
inputs the plan's promotion gate must re-confirm; the durable
contract lives in the plan doc, not here.

## Decisions made at scoping time

Each decision carries a `Verified by:` code citation per
[`shared.md`](../../../../../spec/planning/shared.md) "`Verified
by:` annotations on load-bearing claims."

### SD1 — Classify against the in-request walked parsed-doc set; no fresh read, no second query

The t4 scoping doc's only handed-down p1 mechanism is: does p1
reuse the slug set the handler already walked, or do a fresh read?
Decomposed into shapes before analyzing per
[`shared.md`](../../../../../spec/planning/shared.md) "Decompose
options into shapes before analyzing":

- **A1 — reuse the in-request parsed-doc set.** `Server.index`
  already calls `walkPlans` → the `[]parsedDoc` slice and
  `loadActiveWorkInstances` → the active map once per request,
  before `buildTree`. The roster is built as a pure function of
  those two already-in-hand values: the bound set is the set of
  parsed-doc slugs, and every active work-instance is classified
  by membership in it.
- **A2 — reuse the *built tree* (`roots` after `buildTree`).**
  Rejected. `buildTree` drops a parsed doc whose parent is absent
  (a tree gap) from the rendered `roots`, so deriving the bound
  set from post-`buildTree` nodes would classify a
  parent-gapped-but-walked doc as **unbound**, contradicting the
  task-level contract "bound = slug matches a walked plan doc."
  The correct membership set is the walker's parsed-doc output —
  the same input `buildTree` consumes — not its drop-filtered
  output.
- **B — fresh read (roster re-walks `plansPath` / re-queries the
  DB independently).** Rejected. A second per-request walk or
  query duplicates work already done in the same request and
  introduces a second source of truth that can diverge *within a
  single request* (the roster's walk racing the forest's walk
  against a mid-request file edit). The walk-on-every-request
  invariant is "one per-request read, not a cache," not "as many
  per-request reads as components want"; the t4 task plan's
  "Roster membership" contract explicitly says the roster
  "consumes the same per-request walk."

Decision: **A1**. The roster loader is a pure function over the
`[]parsedDoc` and the active-work-instance map already loaded by
`Server.index`; it performs no walk and no query of its own. The
bound set is the parsed-doc slug set (the walker's output, the
same value `buildTree` joins against), so the classification
matches the contract even for a doc the tree later drops for a
parent gap.

`Verified by:`
[`Server.index` in site.go](../../../../../internal/site/site.go)
calls `walkPlans` then `loadActiveWorkInstances` then `buildTree`,
all per HTTP request — the two pre-`buildTree` values the roster
loader consumes already exist in that scope;
[`loadActiveWorkInstances` in site.go](../../../../../internal/site/site.go)
selects `slug, actor` for `state = active` rows with no tree-membership
filter and returns a `slug → []*ActiveWorkInstance` map (the
unbound rows are present in this map; they are not dropped here);
[`buildTree` in tree.go](../../../../../internal/site/tree.go)
attaches `active[d.Slug]` only for parsed docs and silently drops
a node whose parent is absent — that join is the current
unbound/parent-gap drop site the roster bypasses by classifying
against the parsed-doc set directly.

### SD2 — p1 entry label is the slug; the `wst-<uuid>` actor is carried for ordering/identity but never rendered

p1 has no name source (the reported-`name` metadata arrives only
with p2's client change + event-log join). Under the t4
task-level naming contract the label is the reported name else the
work-instance slug, **never** the `wst-<uuid>` actor; with no name
source in p1, every p1 entry's label is therefore the slug
fallback. The actor is still needed in the loader output as the
identity key (p2 joins the event log on it) and as the
deterministic secondary sort key (multiple active work-instances
can share a slug under `(slug, actor)` idempotency, and the
walk-on-every-request page must render a stable order).

- **Rejected alternative — render the actor as secondary entry
  text in p1.** It would surface the `wst-<uuid>` the task-level
  naming contract bans from the roster UX and pre-empt p2's
  named-session display with a uuid-shaped placeholder (the same
  failure mode t1's placeholder was forbidden from: pre-empting a
  later task's surface). The actor stays loader-internal data in
  p1, not rendered.

`Verified by:`
[`sessionActor` in cmd/workstream-tracker/register.go](../../../../../cmd/workstream-tracker/register.go)
generates the `wst-<uuid>` actor (the string the roster must not
show as a label);
[`work_instances` in schema.go](../../../../../internal/db/schema.go)
declares `slug TEXT NOT NULL`, so every active row the roster
lists has a slug to fall back to;
[`activeWorkInstanceID` in handlers.go](../../../../../internal/api/handlers.go)
keys idempotency on `(slug, actor, state=active)`, so `(slug,
actor)` is the per-instance identity the loader sorts and p2
later joins on.

## Decisions resolved at the task-plan level (not re-opened here)

These were resolved as task-level contract during t4 task-plan
drafting and are inherited, not p1 decisions — listed so the plan
does not re-litigate them:

- Metadata read policy (register baseline + latest-event overlay,
  per-request) — t4 Contracts "Reported data"; p2 mechanism.
- Name-then-slug label rule, uuid banned — t4 Contracts "Session
  identity and naming." p1 realizes the slug-fallback arm (SD2);
  the reported-name arm is p2.
- Event-log join, no schema change; client/CLI metadata send — t4
  Contracts "Reported data" / scoping D1; entirely p2.
- N ≥ 2 split (bare roster, then enrichment) — t4 scoping D5.

## Plan structure handoff

- **Doc shape:** phase plan (the `p1` file of the N ≥ 2 task plan
  [`t4-session-roster.md`](../t4-session-roster.md)). Opens with a
  `## Context` naming the parent task plan; inherited rule
  references (Cross-Cutting Invariants, naming, cross-phase
  decisions) cite the parent task plan's section by name rather
  than duplicating it, per
  [`task-plan.md`](../../../../../spec/planning/task-plan.md) "How
  a phase plan cites its parent task plan."
- **PR-count branch test:** estimated **N = 1 PR**. The diff is
  the roster loader + classification in `site.go`, a roster field
  on `indexData` in `render.go`, the `roster.go` region body, and
  the roster tests, plus the one stale cross-region assertion
  reconciliation (see reality-check input below) — one subsystem
  (the site render path), well under the >5-subsystem / >300-LOC
  split threshold in
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "PR-count predictions need a branch test."
- **Cross-phase coordination** is via the parent task plan (the
  D1–D5 decisions and its Cross-Cutting Invariants), not
  pre-locked cross-phase contracts in the phase plan, per
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "Cross-PR coordination." p1 records the loader/`indexData`/region
  shape p2 consumes as the sibling-interface handoff in the plan's
  Goal/Contracts, tagged for p2 to verify at its drafting.
- **Parent-doc update owed by this p1 drafting.** Per
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "Just-in-time scoping and plan drafting" (drafting updates the
  parent tracking table in the same change), this drafting change
  updates the milestone [`README.md`](../README.md) Task Status
  `t4-p1` row to `Proposed` and reconciles its stub-prose (p1 is
  now a drafted phase plan; p2 and t3 remain stubs) — mirroring
  the t1 precedent. The t4 *task-plan*-level obligations
  (milestone t4 row, the two deferral-note resolutions) were
  already discharged in the t4 task-plan drafting PR (#34) and are
  not re-touched. The implementing PR advances only the onward
  Status values (`t4-p1` row and this plan → `Landed`), not a
  fresh parent reconciliation.

## Reality-check inputs

Load-bearing codebase claims the plan's promotion gate must
re-confirm against the branch at drafting time (each a
one-sentence falsifier):

- *"`loadActiveWorkInstances` returns every active work-instance
  including unbound ones, with no tree-membership filter."*
  Falsifier: the query filters by plan-tree membership, or omits
  unbound rows. Check:
  [`loadActiveWorkInstances` in site.go](../../../../../internal/site/site.go)
  selects `slug, actor WHERE state = active` and returns a
  slug-keyed map of all such rows.
- *"The unbound drop is in `buildTree`, not the query, so the
  roster bypasses it by classifying against the parsed-doc set."*
  Falsifier: the query already drops unbound, or `buildTree` does
  not drop parent-gapped docs. Check:
  [`buildTree` in tree.go](../../../../../internal/site/tree.go)
  attaches `active[d.Slug]` only for parsed docs and drops a node
  whose parent is absent.
- *"`indexData` carries only `Roots` + `PlansPath` today, so a
  roster field is a net-additive struct change."* Falsifier:
  `indexData` already has a roster/session field. Check:
  [`indexData` in render.go](../../../../../internal/site/render.go).
- *"The roster region body is isolated to `roster.go`; the shell
  composes it via `{{template "roster" .}}` and the funcs are set
  on the shared `indexTmpl`."* Falsifier: the roster define lives
  in `render.go`, or the roster needs a template func not on
  `indexTmpl`. Check:
  [`roster.go`](../../../../../internal/site/roster.go) owns
  `{{define "roster"}}` / `{{define "roster-style"}}` parsed into
  `indexTmpl` via `init()`;
  [`render.go`](../../../../../internal/site/render.go) shell calls
  `{{template "roster" .}}`.
- *"A forest-region test file asserts the soon-to-be-removed
  roster placeholder string, so p1 must reconcile one cross-region
  assertion."* Falsifier: the only roster-content assertion lives
  in `roster_test.go`. Check:
  [`TestRenderEmptyStateInForestRegion` in forest_test.go](../../../../../internal/site/forest_test.go)
  asserts the literal placeholder text "The session roster lands
  in a later task." — p1 removes that text, so that one assertion
  must be reconciled to the new roster state (named in the plan's
  Files to touch as an acknowledged, plan-flagged cross-region
  reconciliation, not an unflagged region-boundary crossing).
- *"No build/test wrapper exists; the Go toolchain is the gate."*
  Falsifier: a Makefile/justfile/script wraps build+test. Check:
  repo root has no Makefile/justfile; `scripts/` is only
  `assemble.sh` (matches t1's confirmed finding and the t4 task
  plan Validation Gate).
- *"A render-test harness exists to extend."* Falsifier: no render
  test entry point. Check:
  [`renderTree` in render_test.go](../../../../../internal/site/render_test.go)
  calls `renderIndex(&buf, indexData{...})`;
  [`roster_test.go`](../../../../../internal/site/roster_test.go)
  carries t1's placeholder coverage (the t4-owned test surface p1
  rewrites).
- *"A second populated plan tree exists for manual render
  observation."* Falsifier: only the dogfood tree exists. Check:
  `docs/plans/demo-workstream/` is a populated second root.
- *"The registration handshake yields a real bound active session
  for manual observation."* Falsifier: `register` does not create
  an active work-instance reachable by the site. Check: the
  session-start handshake against the running server returned
  `http_status=201` with a real `work_instance_id` for slug
  `workstream-tracker-1-0-m2-t4-p1` (a slug that IS a walked plan
  doc → a bound session); a deliberately typoed slug registered
  for the gate is the unbound case.
