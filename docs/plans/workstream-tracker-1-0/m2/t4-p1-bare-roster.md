---
slug: workstream-tracker-1-0-m2-t4-p1
Status: Proposed
short_description: Bare bound/unbound session roster replacing the t1 placeholder (no event join, no client change)
---

# t4 p1 — Bare bound/unbound roster

## Context

This is **phase 1** of the N ≥ 2 task plan
[`t4-session-roster.md`](t4-session-roster.md)
(`workstream-tracker-1-0-m2-t4`, "Session roster + work-item
enrichment"), drafted just-in-time against merged code per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Just-in-time scoping and plan drafting." Its scoping doc is
[`scoping/t4-p1-bare-roster.md`](scoping/t4-p1-bare-roster.md).

**What this covers.** The workstream-tracker page today renders a
plan-tree forest plus t1's deliberate roster-region *placeholder*.
A registered session whose slug does not match a plan-tree node is
silently dropped from the page — there is no surface where an
orphan or typoed-slug session can be seen. This phase replaces the
placeholder with a **bare session roster**: a list of every active
work-instance, each marked **bound** (its slug matches a walked
plan doc) or **unbound** (its slug is absent from the walked
tree), using only the slug/actor/state data the server already
has.

**Why now.** t1 made the roster region an independently-owned
surface and the milestone committed the "every session can be
accounted for" observability goal; the bare roster is the
independently-shippable first half of that — it makes
registered-but-unbound work visible immediately, before the p2
enrichment (named sessions, the event-log join, the expandable
raw-JSON detail) builds on it. Shipping p1 alone has observable
accountability value.

**What surfaces this touches.** Only the website's roster region
and its request-time data path: the per-request site loader, the
shared page-data struct it feeds, and the roster region body. **No
event log, no register client/CLI, no schema, no API, no forest
region, no shell layout.** Those are explicitly out of p1 (the
event-log join and the client change are p2).

The inherited WHAT is locked at task-plan level (the t4 Contracts
"Roster membership" and "Session identity and naming"); the only
p1-level decision was the pure classification-data-source
mechanism, resolved in scoping SD1 (reuse the in-request walked
parsed-doc set; no fresh read, no second query) and SD2 (slug is
the p1 label; the `wst-<uuid>` actor is carried for
ordering/identity but never rendered).

## Goal

Opening the page renders, in t1's roster region, a list of
**every active work-instance** the server knows about. Each entry
shows the work-instance **slug** as its label and is visibly
classified **bound** or **unbound**. A bound entry is one whose
slug matches a walked plan doc; an unbound entry is one whose slug
is absent from the walked tree — the orphan/typoed-slug sessions
the forest currently drops. A session that reported nothing beyond
the registration minimum still lists. When there are no active
work-instances at all, the region renders a deliberate, observed
"no active sessions" state — not a blank panel. The forest region,
the shell layout, and v0.1's per-node forest actor markers are
unchanged, and the render path stays walk-on-every-request (no
cache, no second read).

**Sibling-interface handoff to p2** (recorded here, verified at p2
drafting per
[`task-plan.md`](../../../../spec/planning/task-plan.md) "Cross-PR
coordination"): p1 establishes (a) the per-request roster loader
in `site.go` built as a pure function over the already-walked
parsed-doc set and the already-loaded active-work-instance map;
(b) the roster field on the shared `indexData`; (c) the
`roster.go` region body with the bound/unbound classification. p2
enriches the loader with the event-log join and adds the reported
name + expandable raw-JSON detail to the same region; it does not
need to restructure p1's loader signature or the region.

## Contracts

Task-level behavior contracts are owned by the parent task plan
[`t4-session-roster.md`](t4-session-roster.md) Contracts and are
cited, not restated, here per
[`shared.md`](../../../../spec/planning/shared.md)
scoping-vs-duplication discipline (the recursion across
task/phase boundaries in
[`task-plan.md`](../../../../spec/planning/task-plan.md) "How a
phase plan cites its parent task plan"). The bullets below are
p1's own realization shape — the observable end-state each surface
must reach. Estimate-shaped sections (Files to touch, Execution
Steps) are labeled as estimates per
[`shared.md`](../../../../spec/planning/shared.md) "Plan content
is a mix of rules and estimates."

### Roster membership and classification

- The roster lists **every active work-instance**, bound and
  unbound, with **no plan-tree-membership filter**. Terminal
  (completed/abandoned) work-instances do not appear. The list is
  built from the same per-request active-work-instance map
  `Server.index` already loads — the roster does not re-query.
  `Verified by:`
  [`loadActiveWorkInstances` in site.go](../../../../internal/site/site.go)
  selects all `state = active` rows (no tree filter) into a
  slug-keyed map;
  [`Server.index` in site.go](../../../../internal/site/site.go)
  loads that map once per request before `buildTree`.
- Each entry is classified **bound** iff its slug is a member of
  the **walked parsed-doc slug set** (the walker's output, the
  same value `buildTree` joins against), else **unbound**. The
  classification is evaluated **per request** against that
  in-request set — never a cached or post-`buildTree` set, so a
  parsed doc the tree later drops for a parent gap is still
  classified bound. `Verified by:`
  [`buildTree` in tree.go](../../../../internal/site/tree.go)
  consumes the parsed-doc slice and is the current
  unbound/parent-gap drop site; the roster classifies against the
  same slice the handler passes `buildTree`, bypassing that drop
  (scoping SD1).
- The bound vs. unbound state is **visibly distinguishable** in
  the rendered region (a per-entry marker), not merely a struct
  field. The unbound case is rendered, not implied. `Verified
  by:` the Validation Gate's manual render observation of a
  deliberately-unbound session against a real render (Bans on
  surface require rendering the consequence — parent task plan
  Self-Review Audits, `validation-honesty`).
- The roster **lists only** — it adds no promote-into-tree,
  dismiss, or attach affordance for unbound sessions. `Verified
  by:` parent milestone [`README.md`](README.md) Cross-Task Risk
  "Reviewer-flag any t4 drafting that adds a promote/dismiss
  affordance"; the roster tests assert no button/form/input
  element is emitted.

### Entry label and identity

- An entry's display label is the work-instance **slug**. The
  internal `wst-<uuid>` actor is **never rendered** in the roster
  region in p1 (it is carried in the loader output only as the
  identity key and the deterministic secondary sort key). p1 has
  no reported-name source — that arrives with p2's client change
  and event-log join — so every p1 label is the slug-fallback arm
  of the task-level name-then-slug rule. `Verified by:` the t4
  task plan Contracts "Session identity and naming" (label is
  reported name else slug, never the `wst-<uuid>` actor);
  [`work_instances` in schema.go](../../../../internal/db/schema.go)
  declares `slug TEXT NOT NULL`, so the slug fallback is always
  available; scoping SD2.
- Roster order is **deterministic across requests**: entries are
  ordered by `(slug, actor)` so the walk-on-every-request page
  does not reshuffle between loads. `Verified by:`
  [`activeWorkInstanceID` in handlers.go](../../../../internal/api/handlers.go)
  keys idempotency on `(slug, actor, state=active)`, so `(slug,
  actor)` is a stable per-instance identity available without the
  event log.

### Empty state

- When there are zero active work-instances, the roster region
  renders a **deliberate, observed "no active sessions" state** —
  a framed panel with the heading and an explanatory line, in the
  existing muted vocabulary — **not** a blank or zero-height
  region. `Verified by:` the parent task plan's Validation Gate
  inherits "Bans on surface require rendering the consequence"
  ([`task-plan.md`](../../../../spec/planning/task-plan.md)); the
  no-session consequence is observed in this plan's Validation
  Gate, not asserted from template source.

### Render path and region ownership

- The roster **region body** (its template + scoped CSS) lives in
  [`roster.go`](../../../../internal/site/roster.go), replacing
  the t1 placeholder define. The **forest region body and the
  shell layout are not changed.** p1 **does** edit the shared
  request-time data path — a roster loader in `site.go` and a
  roster field on `render.go`'s `indexData` / `renderIndex` —
  which the parent task plan's Cross-Cutting Invariants
  ("Region/shell boundary") and the milestone Cross-Task
  Invariant data-path carve-out explicitly permit and which is
  **not** reviewer-flag. `Verified by:`
  [`roster.go`](../../../../internal/site/roster.go) owns
  `{{define "roster"}}` / `{{define "roster-style"}}` parsed into
  `indexTmpl` via `init()`;
  [`indexData` in render.go](../../../../internal/site/render.go)
  carries only `Roots` + `PlansPath` today (a roster field is
  net-additive); the t4 task plan Contracts "Render path and
  ownership" + Cross-Cutting Invariants "Region/shell boundary".
- The render path stays **walk-on-every-request**: the roster is
  built from the per-request `walkPlans` + `loadActiveWorkInstances`
  values; no cache, file-watch, in-memory accumulation, or second
  read is introduced. `Verified by:`
  [`Server.index` in site.go](../../../../internal/site/site.go)
  performs exactly one `walkPlans` and one
  `loadActiveWorkInstances` per request; the roster loader is a
  pure function over those (scoping SD1).

## Cross-Cutting Invariants

p1 introduces no new cross-cutting rule. It inherits the parent
task plan's Cross-Cutting Invariants
([`t4-session-roster.md`](t4-session-roster.md) "Cross-Cutting
Invariants") — walk-on-every-request, schema-loose-not-homogenized,
v0.1 forest actor markers do not regress, region/shell boundary,
additive registration data — cited here per the parent-citation
discipline rather than duplicated. The two p1's diff most directly
brushes (reviewer-flag candidates here):

- **Walk-on-every-request** — satisfied by SD1 (pure function over
  the existing per-request values; no second read).
- **Region/shell boundary** — satisfied by keeping the region body
  in `roster.go` and limiting shared-path edits to the
  carved-out `indexData` / `renderIndex` / loader plumbing; the
  forest region body and shell layout are untouched.

## Files to touch

*Estimate of expected shape per
[`shared.md`](../../../../spec/planning/shared.md) "Plan content
is a mix of rules and estimates"; implementation may revise with
the deviation reported per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Plan-to-PR Completion Gate."*

**Modify (estimate):**

- `internal/site/site.go` — add the per-request roster loader (a
  pure function over the already-walked parsed-doc set + the
  already-loaded active map producing the classified entry list)
  and the `RosterEntry` type; wire it into `Server.index` to pass
  the roster into `indexData`. Shared request-time data path
  (carved out by the parent task plan / milestone invariant).
- `internal/site/render.go` — add a roster field to `indexData`
  and pass it through `renderIndex`. Shared composition (the
  data-path carve-out; the shell layout/composition is not
  otherwise touched).
- `internal/site/roster.go` — replace the placeholder
  `{{define "roster"}}` with the bound/unbound list render + the
  deliberate empty state; extend `{{define "roster-style"}}`. The
  roster region body — t4's owned surface.
- `internal/site/roster_test.go` — replace the t1 placeholder
  test with p1 coverage (membership, bound/unbound classification
  incl. a parent-gapped-but-walked doc still bound, slug-as-label
  with no uuid rendered, deterministic order, the empty state, no
  affordance element). t4's owned test surface.
- `internal/site/forest_test.go` — **one assertion only**:
  `TestRenderEmptyStateInForestRegion` currently asserts the
  literal t1 placeholder string "The session roster lands in a
  later task.", which p1 removes. Reconcile that single assertion
  to the roster's new state so the test's intent (an empty forest
  does not blank the roster) survives. This is a **plan-flagged,
  acknowledged cross-region test reconciliation**, not an
  unflagged region-boundary crossing: naming it here at plan time
  is the spec-required handling
  ([`task-plan.md`](../../../../spec/planning/task-plan.md) "if a
  reviewer flags a gap that should have been named at plan time,
  fix the plan first"). It is a t1-era estimate deviation
  (forest_test.go was "intentionally not touched" by t4) and is
  reported under `## Estimate Deviations` in the implementing PR.
  `Verified by:`
  [`TestRenderEmptyStateInForestRegion` in forest_test.go](../../../../internal/site/forest_test.go)
  asserts that literal placeholder substring.

**Intentionally not touched** *(estimate — where we don't expect
changes, not a hard prohibition):*

- `internal/site/forest.go`, `internal/site/tree.go`,
  `internal/site/walker.go` — forest region body, forest join,
  and walker unchanged (the v0.1 forest-actor no-regress
  invariant; the `RosterEntry` type lives in `site.go`, not
  `tree.go`, to avoid brushing the forest file).
- `internal/site/render_test.go` — the shared `renderTree` helper
  and the shell test stay as-is; p1 adds its own roster-render
  helper in `roster_test.go` rather than changing the shared
  helper's signature (which forest tests depend on).
- `internal/db/*`, `internal/api/*`,
  `internal/registerclient/*`, `cmd/workstream-tracker/*` — no
  schema, API, event-log, or register-client change (those are
  p2).
- Any caching/file-watch layer — none exists and none is
  introduced (walk-on-every-request).

## Validation Gate

The Go toolchain is the gate (no build/test wrapper exists;
`Verified by:` repo root has no Makefile/justfile, `scripts/` is
only `assemble.sh` — matches t1's confirmed finding and the
parent task plan Validation Gate). p1's gate:

- `gofmt -l internal cmd` reports no files; `go build ./...`,
  `go vet ./...`, `go test ./...` all clean, including the new
  roster cases and the reconciled `forest_test.go` assertion.
- New automated roster coverage asserts: every active
  work-instance (bound + unbound) appears; a bound and an unbound
  slug are distinguishably marked; a parsed doc the tree drops
  for a parent gap is still classified **bound**; the label is
  the slug and no `wst-<uuid>` substring is rendered as a label;
  order is deterministic; a session with only slug/actor still
  lists; the zero-session empty state renders the framed
  observed panel; no button/form/input affordance is emitted.
- **Manual render observation** (Bans on surface require
  rendering the consequence). Run the server
  (`go run ./cmd/workstream-tracker`, port 8080) against the
  dogfood tree, with at least one **bound** active session (the
  session-start handshake registered slug
  `workstream-tracker-1-0-m2-t4-p1`, which is a walked plan doc)
  and one **deliberately unbound** session (register a typoed
  slug). Visually confirm in a real browser render: every active
  session appears in the roster region; bound vs. unbound is
  distinguishable; each entry's label is its slug and no
  `wst-<uuid>` shows as a label; the forest region still renders
  its existing per-node actor markers unchanged; one page scroll
  preserved (shell layout unchanged). Then stop both sessions (or
  observe a tree with none) and confirm the deliberate
  "no active sessions" state renders as an intentional panel, not
  a blank gap. Observation is on the real render, not asserted
  from template source (`validation-honesty`).

## Execution Steps

*Estimate of ordering; deviating from a step is an estimate
deviation, not a contract breach.*

1. **Baseline validation.** On a clean tree, run the Validation
   Gate toolchain commands and confirm green *before* editing, so
   a pre-existing failure isn't misattributed.
2. **Branch hygiene.** Implement on a dedicated branch off
   current `main`; no unrelated changes ride along. Keep the
   plan-drafting commits (scoping + plan + promotion) distinct
   from the implementing commits.
3. **Loader + data path.** Add `RosterEntry` and the per-request
   roster loader in `site.go` (pure function over the
   already-walked parsed-doc set + the already-loaded active
   map), the `indexData` roster field in `render.go`, and wire
   `Server.index`.
4. **Region body.** Replace the `roster.go` placeholder define
   with the bound/unbound list + empty-state render and extend
   the roster CSS; keep the forest region body and shell layout
   untouched.
5. **Tests.** Rewrite `roster_test.go` for the p1 coverage above
   (add a roster-render helper there); reconcile the one stale
   `forest_test.go` placeholder assertion.
6. **Self-review.** Run the Self-Review Audits below against the
   diff.
7. **Final validation.** Re-run the full Validation Gate
   including the manual render observation (observe the
   bound/unbound render and the empty state, don't assume them).
8. **PR preparation.** PR body carries `## Review Stance` (the
   drafting commits' doc-diff) and `## Estimate Deviations` (the
   `forest_test.go` reconciliation); flip this plan's Status
   `Proposed → Landed` and advance the milestone Task Status
   `t4-p1` row in the same PR per the Plan-to-PR Completion Gate
   (Validation Gate fully satisfiable pre-merge — no post-merge
   gate — and p1 is a single implementing PR, so the default
   same-PR `Landed` flip applies; the N ≥ 2 parent task plan
   stays `Proposed` until p2 lands per
   [`task-plan.md`](../../../../spec/planning/task-plan.md) "Task
   plan terminal state when N ≥ 2").

## Self-Review Audits

From the parent task plan's Self-Review Audits
([`t4-session-roster.md`](t4-session-roster.md)), the subset
mapped to p1's diff surfaces:

- **validation-honesty** — the bound/unbound render and the empty
  state are "passed" only when observed in a real render, not
  asserted from template source.
- **readiness-gate-truthfulness** — the `Proposed → Landed` flip
  and the milestone `t4-p1` row advance only when every Goal /
  Contract / Validation item is satisfied or explicitly deferred
  in this plan.
- **rename-aware-diff-classification** — replacing the roster
  placeholder define and the placeholder test, and reconciling
  the cross-region `forest_test.go` assertion, must be
  hand-classified so the "forest region unchanged" and "shell
  unchanged" claims are real, not tooling-fooled by add+delete.

`error-surfacing-user-mutations` has no p1 surface (p1 adds no
user-initiated mutation — the register client/CLI change is p2);
it re-enters scope at p2.

## Out Of Scope

- The event-log join, reported-name display, and the expandable
  raw-JSON detail view — **p2** (the t4 task plan Phase Contracts
  p2 row). p1 renders slug-labelled entries only.
- Any register client/CLI change (sending metadata / a name) —
  **p2**.
- Any schema change — explicitly rejected task-wide in favor of
  the event-log join (t4 scoping D1); p1 introduces none.
- Triage **action** (promote/dismiss/attach for unbound
  sessions) — listing only; deferred past 1.0 (milestone Out of
  Scope; reviewer-flag).
- Humanizing the **forest** actor display / changing the
  `wst-<uuid>` generator — a t4-wide new backlog entry added in
  p2's PR, not done here; p1 does not touch the forest region.

## Risk Register

- **The placeholder-removal silently breaks a cross-region
  test.** `forest_test.go`'s empty-state test asserts the t1
  placeholder string. Mitigation: named in Files to touch as a
  plan-flagged reconciliation, hand-classified under the
  rename-aware audit, and the Validation Gate runs the full
  `go test ./...` so the reconciled assertion is exercised.
- **Roster reshuffles between requests.** Without a deterministic
  order the walk-on-every-request page would reorder entries each
  load. Mitigation: the entry-label contract fixes
  `(slug, actor)` ordering; a roster test asserts deterministic
  order.
- **The uuid leaks into the roster UX.** A naive entry render
  could show the `wst-<uuid>` actor. Mitigation: SD2 keeps the
  actor loader-internal; a roster test asserts no `wst-`-shaped
  label is rendered and the manual observation confirms it.
- **A second per-request read creeps in.** Re-walking or
  re-querying for the roster would violate walk-on-every-request.
  Mitigation: SD1 fixes the loader as a pure function over the
  existing per-request values; the loader takes those values as
  parameters rather than acquiring its own.

## Documentation Currency

- Milestone [`README.md`](README.md) — the Task Status `t4-p1`
  row and the stub-prose are reconciled **in this drafting
  change** (row → `Proposed`; prose: p1 is now a drafted phase
  plan, p2/t3 remain stubs) per
  [`task-plan.md`](../../../../spec/planning/task-plan.md)
  "Just-in-time scoping and plan drafting," mirroring the t1
  precedent. The implementing PR advances the same row and this
  plan's frontmatter Status onward to `Landed`.
- This plan's `Status` lifecycle: `In draft` while drafted
  (committed first), `Proposed` after the promotion-gate
  self-review walk, then `Proposed → Landed` in the implementing
  PR per the Plan-to-PR Completion Gate (Validation Gate fully
  satisfiable pre-merge — no post-merge gate, single implementing
  PR — so the default same-PR `Landed` flip applies; the N ≥ 2
  parent task plan stays `Proposed` until p2 lands).
- No `spec/`, `design/`, or `session-registration.md` currency is
  owed by p1 — those are p2's (the client/CLI + handshake change).

## Related Docs

- [`t4-session-roster.md`](t4-session-roster.md) — parent N ≥ 2
  task plan; owns the task-level Contracts, Cross-Cutting
  Invariants, and Phase Contracts p1 inherits and cites.
- [`scoping/t4-p1-bare-roster.md`](scoping/t4-p1-bare-roster.md)
  — p1 scoping doc (SD1 classification source, SD2 label/identity,
  reality-check inputs).
- [`t1-site-skeleton.md`](t1-site-skeleton.md) — the landed
  two-region shell + Region-ownership contract (with the
  data-path carve-out) p1 builds the roster region into.
- [`README.md`](README.md) — parent milestone; Cross-Task
  Invariants/Risks and the Task Status `t4-p1` row p1 advances.
- [`../../../../spec/planning/task-plan.md`](../../../../spec/planning/task-plan.md),
  [`../../../../spec/planning/shared.md`](../../../../spec/planning/shared.md)
  — the rules this phase plan is structured against.
