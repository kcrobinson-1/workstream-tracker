# Scoping — m2 t2: Expanded in-root nested-box render

Transient deliberation doc for the
[`workstream-tracker-1-0-m2-t2`](../t2-expanded-render.md) task
plan. Deletes in batch with sibling scoping docs at the
milestone-terminal PR per
[`task-plan.md`](../../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns." No Status field — scoping docs are
not part of the plan-doc lifecycle.

## Context summary

t2 replaces the forest region's flat nested-bullet render with
nested, independently-collapsible per-node boxes — each box
carrying its own Status badge and its actor markers — so a root's
internal epic → milestone → task → phase structure reads as a
scannable hierarchy rather than an indented `<ul>`. It is drafted
now because t1 (the two-region site skeleton) has Landed and
established the forest region t2 grows into, and the m2 milestone
explicitly defers t2's one hard decision — the collapse mechanism
versus the v0.2 no-collapse posture — to this drafting session.
The only surfaces touched are the forest-region render template,
the tree builder that feeds it, and their tests.

## Decisions made at scoping time

Each decision is durable rationale the plan references by name
rather than restating; rejected alternatives live here and do not
migrate into the plan doc.

### S1 — Collapse mechanism: HTML-native `<details>`/`<summary>`, zero JavaScript

The nested boxes collapse via the browser's native
`<details>`/`<summary>` elements. No JavaScript, no new route, no
client state.

This supersedes a **narrower** decision, not the architecture.
The durable architectural tenet is *server-rendered,
refresh-to-update, no read-side JSON API, no server push* —
`Verified by:`
[`design/v0.1-design.md` §2](../../../../../design/v0.1-design.md)
("no client-side JavaScript polls for updates, no JSON API is
exposed for reads"), §9 Out of Scope ("Server-to-agent push,"
"Real-time updates without browser refresh"), and §10's `templ`
re-trigger bullet, which states the site is "server-rendered by
design" and that the revisit tripwire is render *complexity, not
interactivity*. `<details>` is static server-rendered HTML and
preserves every clause of that tenet.

What is superseded is the **task-scoped** "no expand/collapse, no
JavaScript" decision made at m1 t4 drafting — `Verified by:`
[m1 README "Cross-Task Decisions" → "Long-description rendering
location: inline static HTML, no panel (RESOLVED at t4
drafting)"](../../m1/README.md) and
[`m1-t4-p1-inline-detail-render.md` Goal](../../m1/t4-p1-inline-detail-render.md)
("No truncation, no JavaScript, no panel, no expand/collapse
(resolved at t4 drafting)"). That decision's *stated rationale*
was that a detail-on-click panel "would be the first JavaScript
and first non-`/` route in the codebase, disproportionate to the
bare-bones v0.2 render." Both costs are moot under `<details>`:
zero JavaScript, zero new routes. The m2 milestone explicitly
routes this call here — `Verified by:`
[m2 README "Cross-Task Decisions" → "Collapsible mechanism
(decide when t2 drafts)"](../README.md) — and requires the
decision be recorded against the v0.2 posture rather than
regressed silently; this section is that record.

**Rejected alternatives:**
- *JavaScript toggle (click handler + class flip).* Would
  introduce the codebase's first client-side script for a purely
  presentational affordance the platform offers natively.
  Rejected as disproportionate and as a genuine regression of
  the surviving architectural tenet, not just the t4 decision.
- *Detail-on-click panel / separate route.* The exact shape m1 t4
  rejected; same rejection holds, and `<details>` makes it
  unnecessary.
- *Stay fully flat (decline to collapse).* Contradicts the locked
  m2 t2 inherited contract and the approved mockup
  [`design/workstreams-view-m2.svg`](../../../../../design/workstreams-view-m2.svg),
  which shows caret affordances on nested boxes.

### S2 — Default expand state derived from active work, recomputed every request

`<details>` has no persistence and there is no client state, so
the server decides each node's initial open/closed state per
render. The rule: a node renders expanded if it, or any
descendant, has an active work-instance; otherwise collapsed.
This directly serves the milestone's qualitative bar ("the
parallel-agent picture is legible at a glance") and is a feature
of the walk-on-every-request model, not a limitation — the view
auto-focuses wherever agents are currently working, with no
stored UI state to go stale.

`Verified by:`
[`buildTree` in tree.go](../../../../../internal/site/tree.go)
wires children to parents and attaches `active[d.Slug]` to each
node before returning roots, so a post-order pass over the
assembled tree can compute "self-or-descendant has an active
work-instance" with no second data source;
[`loadActiveWorkInstances` in site.go](../../../../../internal/site/site.go)
already supplies the `slug → []ActiveWorkInstance` map
`buildTree` consumes.

**Rejected alternatives:**
- *Everything expanded by default.* Defeats the scannability goal
  on a large dogfood tree; the mockup deliberately shows idle
  subtrees collapsed.
- *Everything collapsed by default.* Hides exactly the active
  work the tool exists to surface.
- *Expand by Status (e.g. In progress).* Status is a plan-doc
  authoring fact, not a liveness signal; an idle `In progress`
  node with no session is noise, an active node mid-Status is the
  signal. Active work-instance presence is the correct axis.

### S3 — Expand state carried as an additive `PlanNode` field set in `buildTree`

"Self-or-descendant has an active work-instance" requires a tree
walk the `html/template` layer cannot express cleanly across the
recursive `node` template. The value is computed once in
`buildTree` (post-order, after children are wired) and carried as
a new additive boolean field on `PlanNode`; the template reads it
to emit the `<details open>` attribute conditionally.

`Verified by:`
[`PlanNode` in tree.go](../../../../../internal/site/tree.go) is
the render input struct already threaded to the template;
[`buildTree` in tree.go](../../../../../internal/site/tree.go)
already runs a post-wiring pass (the stable-order `sort.Slice`
loop) so an additional post-order computation is an established
shape, not a novel mechanism. Field name is plan-owned (the plan
"Naming" section decides the exact identifier).

**Rejected alternative:**
- *Template `FuncMap` helper that recurses.* Recursing the child
  slice from inside `html/template` to test for an active
  descendant is the stringly-typed indirection §10 names as the
  eventual `templ` tripwire; computing in Go and passing a bool
  is the established codebase pattern.

### S4 — Node-shape tests assert semantic/structural output, not byte-identity

t2 deliberately rewrites the `node` template, which invalidates
the existing byte-identity pin. The replacement assertions are
semantic/structural (the nested `<details>` is present, the
Status badge and actor markers are present, a leaf/stub still
renders, a field-less node still emits no `.long-desc` /
`.related-prs` markup) rather than a new byte-exact literal. The
explicit goal is **not byte-exact output**; re-ratcheting a fresh
byte literal would only break again at t3.

`Verified by:`
[`TestRenderNoFieldNodeUnchanged` in forest_test.go](../../../../../internal/site/forest_test.go)
is the current byte-identity pin on the `<ul>`/`<li>` node markup
t2 replaces;
[`TestRenderTwoRegionShell` in render_test.go](../../../../../internal/site/render_test.go)
sub-asserts the old exact node line inside t1's shell test;
[`design/v0.1-design.md` §10](../../../../../design/v0.1-design.md)
"`templ`" Risk Register direction is "asserting semantic (not
byte-exact) output on net-new partials," which this decision
applies. The cross-file edit consequence (a shell-test file
owned by t1 gets one node-line sub-assertion updated) is a
real boundary nuance — the plan's Contracts/Files section
sanctions it explicitly so review does not read it as a
shell-region violation under the m2 file-boundary invariant.

### S5 — Scoping doc is warranted (narrow-surface carve-out declined)

t2's file surface is small, but S1 is a genuine
options-considered decision that supersedes a prior locked
decision and carries a mandatory multi-source `Verified by:`
citation — exactly the content the
[`task-plan.md`](../../../../../spec/planning/task-plan.md)
"Narrow-surface plans may skip the scoping doc" carve-out
excludes (it applies only when options analysis "would not
produce material content"). This doc is the deliberate
invocation of the default "scoping first" path, not the
carve-out.

## Open decisions to make at plan-drafting

- **Phase split (N = 1 vs N ≥ 2).** The m2 milestone estimates t2
  plausibly N ≥ 2 (nested-box render, then the narrow-window
  per-node-row degrade). This is an estimate, not a contract;
  resolve it with the
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "PR-count predictions need a branch test." Scoping's read
  (non-binding): both the box render and the per-node-row degrade
  live in `forest.go` (template + CSS) and are sequence-steps
  toward one outcome with no independent value — shipping the
  boxes without the degrade ships a known-broken narrow surface —
  so N = 1 is likely unless the branch test shows > 300 LOC of
  substantive logic or > 5 subsystems. The plan's Status block
  records the resolved N with the branch-test result.

## Plan structure handoff

- Replace the parent-promotion stub at
  [`../t2-expanded-render.md`](../t2-expanded-render.md) with a
  full task plan. Required sections per
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "Required and optional sections": Status, context preamble,
  Goal, Contracts, Files to touch, Validation Gate.
- Optional sections that apply: Cross-Cutting Invariants (the
  inherited m2 invariants t2 must not regress — actor tags,
  walk-on-every-request, stub render, file boundary), Naming
  (the new `PlanNode` field from S3), Self-Review Audits
  (`validation-honesty` + the general checklist), Risk Register
  (the test-rewrite and narrow-window-degrade residuals), Out Of
  Scope (t3's progress boxes; t4's roster), Backlog Impact (none
  — state explicitly), Related Docs.
- Same change updates the m2 README parent doc: flip the t2 row
  in the Task Status table off "In draft (stub)", and resolve the
  "Collapsible mechanism (decide when t2 drafts)" entry in m2
  "Cross-Task Decisions" by pointing it at S1 (per the
  parent-doc-currency / plan-drafting-updates-parent-table rule).
- Status stays `In draft` until the phase-split open decision is
  resolved and the
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  `In draft → Proposed` promotion-gate self-review has run.

## Reality-check inputs (the plan must re-verify before promotion)

Load-bearing claims this scoping rests on; the plan's promotion
gate re-confirms each against current code.

- **Forest render surface.**
  [`forest.go`](../../../../../internal/site/forest.go) — the
  `forestTemplates` const holds the `forest` template, the
  recursive `node` template (badge, `{{range .WorkInstances}}`
  actor markers, long-desc, related-PRs, then `{{if .Children}}`
  `<ul>` recursion), and `forest-style`; an `init()` parses them
  into the shared `indexTmpl` tree. This is t2's entire owned
  render surface.
- **File boundary.** The shell
  [`render.go`](../../../../../internal/site/render.go) and the
  roster [`roster.go`](../../../../../internal/site/roster.go)
  are off-limits per the m2 "the shell is t1's" cross-task
  invariant; t2 edits `forest.go` and `tree.go` only (plus
  tests). Confirm the `{{template "forest" .}}` composition seam
  is unchanged.
- **Tree builder.**
  [`buildTree` / `PlanNode` in tree.go](../../../../../internal/site/tree.go)
  — children wired to parents, work-instances attached by exact
  slug, roots returned in stable order; confirm the post-order
  insertion point for S3's field still exists as described.
- **Walk-on-every-request.**
  [`Server.index` in site.go](../../../../../internal/site/site.go)
  walks plans + loads active work-instances per request; t2
  touches only the template/builder, not this handler — confirm
  no caching/goroutine is introduced.
- **Node levels.**
  [`slugs.go`](../../../../../internal/slugs/slugs.go) — slug
  parsing yields `root` / `milestone` / `task` / `phase` only
  (never `epic`; the root is the top box); `Position()` returns
  false for a root. Per-level box styling keys off
  `PlanNode.NodeType`.
- **Test pins.**
  [`forest_test.go`](../../../../../internal/site/forest_test.go)
  `TestRenderNoFieldNodeUnchanged` byte-pin and
  [`render_test.go`](../../../../../internal/site/render_test.go)
  `TestRenderTwoRegionShell` node-line sub-assertion both
  reference the pre-t2 markup; confirm they are still the only
  byte-level references to node shape.
- **Validation commands + UI capture.**
  [`docs/dev.md`](../../../../../docs/dev.md) — `go build ./...`,
  `go vet ./...`, `go test ./...`; UI capture is manual
  screenshots at ~1280px against a real `go run` with a
  populated `DB_PATH` and the `docs/plans/` tree (no screenshot
  tooling). The narrow-window per-node-row degrade has no
  automated coverage, so "Bans on surface require rendering the
  consequence" makes a manually-observed resized render a
  required gate step.
