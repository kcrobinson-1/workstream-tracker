---
slug: workstream-tracker-1-0-m2-t2
Status: Landed
short_description: Nested collapsible per-node Status boxes inside the forest region
---

# Task 2 — Expanded in-root nested-box render

## Status

`Landed`. Implemented across the two estimated commits with no
deviation from the Commit Boundaries (commit 1: forest.go
nested-box render + CSS + semantic test rewrite; commit 2:
tree.go `ActiveInSubtree` post-order field + conditional
`<details open>` + C4 narrow-window degrade). The Validation
Gate is satisfied: `go build ./...`, `go vet ./...`,
`go test ./...` all pass, and the render + narrow-window degrade
were observed against a real `go run` reading the `docs/plans/`
dogfood tree and a populated `DB_PATH` (active subtrees open,
idle collapsed, stubs render as leaf boxes, the `@media
(max-width: 48rem)` degrade rule served). No estimate deviations.
The transient `scoping/` doc is
intentionally retained — it deletes in batch with sibling
scoping docs at the milestone-terminal PR (t3/t4 are still
undrafted stubs), not at this task-terminal PR.

The drafting record below is preserved for the decision history.

`Proposed`. The one named input — the phase split — was
**resolved: N = 1**, by the
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"PR-count predictions need a branch test." Evidence: the change
touches a single subsystem (`internal/site` render) and ≈ 15–20
LOC of substantive logic (the `buildTree` post-order pass — the
template/CSS rewrite is markup, not algorithmic logic), far below
the >5-subsystem / >300-LOC split thresholds. Candidate phase
boundaries were enumerated (render vs. degrade; model vs. render;
always-open vs. collapse; collapsible vs. active-work
default-open); the strongest (collapsible boxes | active-work
default-open) was rejected for N ≥ 2 because, unlike m1-t4's
subprocess P2, t2's default-open logic shares the same validation
apparatus (`go test` + the same manual UI capture) as the rest —
it is a commit boundary, not a distinct validation surface. The
real seam is preserved as a Commit Boundary (below), not a phase.

The [`task-plan.md`](../../../../spec/planning/task-plan.md)
`In draft → Proposed` promotion-gate self-review ran (contract
decision-completeness, the universal `Verified by:` walk,
reality-check re-confirmation against current code, no descent
into implementation prescription; N = 1 so no phase skeletons),
and the flip past the locked-decision supersession (C2/S1) was
explicitly authorized by the contributor. Two findings the gate's
spec-conformance and end-to-end-coherence steps should have
caught were raised in review and applied as in-place corrections
(plan remains `Proposed`): the estimate-preface requirement
([`shared.md`](../../../../spec/planning/shared.md) "Plan content
is a mix of rules and estimates") now labels Files to touch and
Commit Boundaries; and the C3/C5 leaf-box contradiction (a leaf
has no disclosure control, so the default-open rule is scoped to
boxes with children) is resolved.

Drafting deliberation, rejected alternatives, and `Verified by:`
grounding for every decision below live in
`scoping/t2-expanded-render.md`
(S1–S5) and are referenced by name here, not restated, per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns."

## Context

This plan covers the plan-tree forest's visual structure. Today a
root's descendants render as an indented bullet list: every node
is one line of badge + label, nested under `<ul>`/`<li>`. This
task turns each node into its own box — a header carrying the
node's Status badge and the actor markers for any agent active on
it, with the node's children nested as boxes inside it, each box
independently collapsible.

It is being done now because the prior task (the two-region site
skeleton) has Landed and established the forest region this work
grows into, and because the milestone deliberately left one
decision — whether reintroducing collapse is compatible with the
read-experience tooling's "no expand/collapse" posture — to be
resolved at this task's drafting. Resolving it unblocks the
remaining forest work: the next task renders doc-driven progress
boxes *inside* the per-node box this task establishes.

The surfaces touched are the forest-region render template, the
in-memory tree the template renders, and their unit tests. No
API, schema, route, or dependency changes; the page stays
server-rendered and refresh-to-update.

## Goal

Within the forest region, a root and its descendants render as
nested epic/root → milestone → task → phase boxes. Each box shows
its node's Status badge and its active-work actor markers, and is
independently collapsible via native HTML, with no JavaScript and
no new route. Among boxes that have children, one with active
work in its subtree renders expanded and an idle subtree renders
collapsed; a leaf box has no children, no disclosure control, and
no expand state. The expand state is recomputed from live
work-instance data on every request, with no stored UI state. The flat nested-bullet render is fully replaced.
Nodes carrying no rich data — including a bare `slug` +
`Status: In draft` stub — still render as a valid box, and the
existing per-node actor markers, long description, and related-PR
list still render unchanged on every box.

Verifiable when: rendering the dogfood `docs/plans/` tree shows
each root's internal epic → milestone → task → phase structure as
nested collapsible boxes; a box with an active work-instance in
its subtree is open on load and an idle subtree is closed; every
box still shows its actor markers and detail; a stub renders as a
leaf box (no disclosure control) with just its badge and label;
and the
narrow-forest-column degrade is an observed, intentional layout,
not a collision.

## Contracts

Final WHAT shape. HOW grounding is in
`scoping/t2-expanded-render.md`.

- **C1 — Nested per-node boxes replace the bullet list.** The
  recursive node render emits each node as a box whose header
  carries the Status badge and the node's actor markers, and
  whose children are rendered as nested boxes within it. The
  `<ul>`/`<li>` descendant nesting is removed. Per-level visual
  treatment keys off the node type (`root` / `milestone` /
  `task` / `phase`; there is no `epic` node — the root is the
  top box).
- **C2 — Collapse is native `<details>`/`<summary>`, zero JS.**
  Each box with children is collapsible through the browser's
  native disclosure element. No client-side script, no new
  route, no read-side API. This is the resolution of the
  milestone's deferred collapse-mechanism decision; rationale
  and the supersession of the m1 t4 "no expand/collapse"
  decision are recorded in
  `scoping/t2-expanded-render.md`
  S1, and the m2 README "Cross-Task Decisions" entry is
  reconciled to point at it in the same change as this plan.
- **C3 — Default expand state from active work.** Expand/collapse
  applies only to a box that has a disclosure control — i.e., a
  box with children; a leaf box has none (C5) and therefore has
  no expand state at all (its badge, label, actor markers, and
  detail are always visible — there is nothing to disclose). A
  box *with children* renders expanded iff it, or any descendant,
  has an active work-instance; otherwise collapsed. The
  active-in-subtree signal is computed in the tree builder as a
  post-order pass after children are wired, carried on the
  render-input node struct as one new additive boolean field, and
  read by the template to emit the open state on collapsible
  boxes only. The "or any descendant" clause is what forces an
  active leaf's ancestor boxes open so the leaf is visible.
  Recomputed every request from live data — no persistence,
  consistent with the walk-on-every-request invariant.
- **C4 — Right-aligned Status region; t2 owns the per-node-row
  narrow-window degrade.** The box header places the label group
  left and the Status group right. t2 introduces this
  right-alignment and therefore owns the degrade when the forest
  column is narrow: the right group reflows below the label
  rather than colliding with or truncating it. The degrade is an
  intentional observed layout (verified per the Validation
  Gate), not an unhandled overflow. This is distinct from the
  region-stacking media query the shell already owns. The right
  group is the surface the next task renders doc-driven progress
  boxes into; t2 establishes the region but renders no progress
  boxes itself (out of scope).
- **C5 — Additive, degrade-safe render.** Every preserved
  surface still renders on each box: actor markers for active
  work-instances, long description, related-PR list. A node with
  none of these — including a `slug` + `Status: In draft` stub,
  and any node with no children — renders as a valid box (a leaf
  box has no disclosure control rather than an empty one). No
  node is dropped or errored by the new render.
- **C6 — Node-shape tests are semantic/structural, not
  byte-exact.** The rewritten node render is asserted by
  structure and presence (nested `<details>` present; badge and
  actor markers present; leaf/stub renders; field-less node
  emits no detail markup), not a byte-identity literal. The
  existing byte pin is replaced, not re-ratcheted to a new
  literal. One node-line sub-assertion in the t1-owned shell
  test file is updated in the same change; this single cross-file
  test edit is sanctioned here (node shape is t2's surface) so it
  is not read as a shell-region boundary violation under the m2
  file-boundary invariant. Grounding: scoping S4.

## Cross-Cutting Invariants

Inherited from the m2 milestone; t2 must not regress them. Each
is a reviewer-flag candidate.

- **Actor tags on nodes do not regress.** Every box renders the
  active-work actor markers the current node line renders. The
  deferred on-box actor-icon work assumes node-level markers
  remain.
- **Render path stays walk-on-every-request.** No caching, no
  file-watch, no goroutine, no in-memory build-up. C3's
  expand-state is recomputed per request inside the existing
  build path.
- **Shell is t1's; t2 stays in the forest region.** t2 edits the
  forest render and the tree builder only. The shell and roster
  files are not touched. The one sanctioned exception is the
  C6 shell-test sub-assertion, called out explicitly.
- **Stub render is preserved, not pre-empted.** A `slug` +
  `Status: In draft` stub remains a valid render (a leaf box);
  t2 does not add progress-box or declared-stage behavior — that
  is the next task.

## Naming

- One new exported boolean field on the render-input node struct
  in the tree builder, expressing "this node or a descendant has
  an active work-instance" (the C3 default-open signal). Exact
  identifier chosen at implementation; it is additive and
  internal to the `site` package.

## Files to touch

_Estimate of the expected file shape, not a binding rule.
Implementation may revise this list when a structural call
requires it — including touching a file listed under
"Intentionally not touched" (that line means "we don't expect to
need these," not "implementation must not touch them"). Any
deviation is normal and is handled via the PR-body Estimate
Deviations callout, with the plan reconciled to what shipped._

- **Modify:** the forest-region render template + styles (the
  node render, its CSS, the `<ul>` removal, the `<details>`
  structure, the header left/right groups and the narrow-window
  degrade); the tree builder (the additive C3 field + its
  post-order computation).
- **Modify (tests):** the forest-region test file (replace the
  byte-identity node pin with semantic/structural assertions per
  C6; add coverage for default-open-by-active-work, leaf/stub
  box, nested structure, narrow-window degrade is not asserted by
  unit test — it is a manual gate step); one node-line
  sub-assertion in the shell test file (C6, sanctioned).
- **Modify (docs):** the m2 README parent doc — Task Status table
  t2 row, and the "Cross-Task Decisions" collapsible-mechanism
  entry reconciled to C2/S1, in the same change as this plan
  (parent-doc currency); `design/v0.1-design.md` §7 updated to
  describe the nested collapsible render on the PR that lands the
  implementation.
- **Intentionally not touched:** the shell render file and the
  roster render file (m2 file boundary); the plan-tree walker /
  frontmatter parser (no new field is read from docs — that is
  the next task); the API, DB schema, register client, and
  routing; any dependency (`<details>` is native HTML).

## Validation Gate

- `go build ./...`, `go vet ./...`, `go test ./...` all pass
  (per [`docs/dev.md`](../../../../docs/dev.md)).
- Unit tests cover: nested-box structure replaces the bullet
  list; a box with an active work-instance in its subtree is
  open and an idle one is closed; actor markers / long
  description / related PRs still render on a box; a stub and a
  no-children node render as valid boxes; field-less node emits
  no detail markup (the preserved additive guard).
- **Manual UI capture** (no screenshot tooling exists; capture is
  manual per [`docs/dev.md`](../../../../docs/dev.md) UI Review):
  screenshots at ~1280px against a real `go run` reading a
  populated `DB_PATH` and the `docs/plans/` dogfood tree —
  confirming nested boxes, default expand/collapse keyed to live
  work, preserved actor markers/detail, and the stub render.
- **Narrow-window degrade observed, not assumed.** Per
  [`task-plan.md`](../../../../spec/planning/task-plan.md) "Bans
  on surface require rendering the consequence," the C4 degrade
  is verified by capturing a narrowed forest column and
  confirming the right group reflows below the label as an
  intentional layout — this is a required gate step, not optional.
- Self-review audits below are run at the implementing commit
  boundary.

## Commit Boundaries

_The N = 1 phase count is a resolved rule (the branch-test
outcome — see Status), not re-estimated here. The specific
two-commit split below is an **estimate** of cohesive review
chunks; the implementer may refine or reshuffle it, with any
deviation handled via the PR-body Estimate Deviations callout and
the plan reconciled to what shipped._

N = 1 (one phase, one implementing PR), two commits. This is the
real phase-D seam (collapsible render | active-work default-open)
preserved as a review boundary without N ≥ 2 orchestration.

- **Commit 1 — nested collapsible render.** The `forest.go` node
  template rewrite (`<ul>/<li>` → `<details>/<summary>` boxes,
  per-type styling, the flex header with the reserved right-side
  Status region) and CSS, plus the semantic/structural test
  rewrite replacing the byte-identity pin (C1, C2, C5, C6). A
  simple static default-open at this commit (all boxes open).
  Reviewable as "the forest renders nested collapsible boxes,
  no surface regressed."
- **Commit 2 — active-work default-open + degrade.** The
  `tree.go` additive field and its post-order computation, the
  template reading it for `<details open>`, the
  default-open-by-active-work tests, and the C4 narrow-window
  per-node-row degrade with its required observed-render gate
  step (C3, C4). Reviewable as "the view opens where work is
  active and degrades cleanly when narrow."

The plan is the only doc (no phase plan files). If commit 2's
diff unexpectedly balloons past the branch-test thresholds at
implementation, the
[`task-plan.md`](../../../../spec/planning/task-plan.md) N = 1 →
N ≥ 2 transition rule governs the split decision then.

## Self-Review Audits

From [`docs/agents/local/self-review-catalog.md`](../../../../docs/agents/local/self-review-catalog.md)
(no repo-specific audits exist):

- **validation-honesty** — a "the render works / the degrade is
  acceptable" claim is valid only if the render and the narrowed
  layout were actually observed, not asserted from the diff. This
  is the audit lens on the manual gate steps above.
- The general self-review checklist
  ([`how-to-use.md`](../../../../docs/agents/shared/self-review/how-to-use.md))
  layered on top.

No data / CI / runbook audit surfaces apply (no schema, no
pipeline, no operational doc touched).

## Risk Register

- **Test rewrite masks a real regression.** Replacing the
  byte-identity pin with semantic assertions loosens the pin by
  design (C6). Mitigation: the semantic assertions explicitly
  cover the preserved surfaces (actor markers, detail blocks,
  stub/leaf render, no-detail-on-field-less), so a dropped
  surface still fails a test; the manual UI capture is the
  backstop for visual regressions a structural assertion would
  not catch.
- **Narrow-window degrade ships as a collision.** t2 introduces
  the right-aligned region that can collide with the label
  (C4). Mitigation: the degrade is a required *observed* gate
  step, not a unit test — it cannot pass on diff inspection
  alone.
- **Default-open rule surprises on a large idle tree.** If most
  of the dogfood tree is idle, most boxes load collapsed.
  Mitigation: this is the intended behavior (S2) and the manual
  capture confirms active work is the thing surfaced; not a
  defect.

## Out of Scope

- **Doc-driven progress boxes.** The right-aligned region C4
  establishes is the host for the next task's progress boxes;
  t2 renders none and reads no new frontmatter field.
- **The session roster.** A separate region owned by a sibling
  task; t2 does not touch it.
- **On-box actor icons.** The deferred milestone work that
  assumes node-level actor markers remain; t2 preserves the
  markers and adds nothing here.

## Backlog Impact

None. No [`docs/backlog.md`](../../../../docs/backlog.md) entry
graduates, is deleted, split, or shifted by this task.

## Related Docs

- [`README.md`](README.md) — parent milestone; its Task Status
  table and collapsible-mechanism decision entry are reconciled
  in the same change as this plan.
- `scoping/t2-expanded-render.md`
  — drafting deliberation, rejected alternatives, and `Verified
  by:` grounding (S1–S5).
- [`../m1/README.md`](../m1/README.md) and
  [`../m1/t4-p1-inline-detail-render.md`](../m1/t4-p1-inline-detail-render.md)
  — the v0.2 "no expand/collapse" task decision C2/S1 supersede.
- [`../../../../design/v0.1-design.md`](../../../../design/v0.1-design.md)
  — §2/§9/§10 architecture (preserved); §7 updated by the
  implementing PR.
- [`../../../../design/workstreams-view-m2.svg`](../../../../design/workstreams-view-m2.svg)
  — the approved target the nested-box render realizes.
