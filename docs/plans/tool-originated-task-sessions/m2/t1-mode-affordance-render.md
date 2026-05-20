---
slug: tool-originated-task-sessions-m2-t1
Status: Proposed
short_description: Mode-affordance render on plan-tree nodes
---

# m2 t1 — Mode-Affordance Render on Plan-Tree Nodes

## Context

t1 is the page-side surface of the
[m2 milestone](./README.md): the rendered plan-tree forest gains
a per-node **mode-affordance** — either a small HTML form that
POSTs to t2's future `/spawn` endpoint with the node's slug and
the chosen mode, or no form at all — driven by the locked **D3**
static map over facts the node already renders. The render side
is the only surface this task touches; t2 (spawn endpoint) and
t3 (SessionStart hook) are three independent surfaces that
converge on t4 per the milestone's Sequencing.

**Why now.** m2 promoted to `Proposed` at PR
[#58](https://github.com/kcrobinson-1/workstream-tracker/pull/58)
with the WHAT contract locked, the three interior tasks seeded
as parallel-draftable skeletons. The locked D3 mode-offer map
and D4 form-shape contract make t1 implementable today against
merged code without waiting on t2's `/spawn` handler or t3's
SessionStart hook (the form's `action="/spawn"` points at a
not-yet-shipped route; a POST before t2 lands receives a 404,
the expected sibling-not-yet-shipped degrade the m2 Sequencing
anticipates). Drafting t1 now keeps the three surfaces moving in
parallel toward t4's convergence.

**Surfaces.** The plan-tree forest's per-node render in
[`internal/site/forest.go`](../../../../internal/site/forest.go);
its supporting computed-value layer in
[`internal/site/tree.go`](../../../../internal/site/tree.go);
the canonical-prefix Status helper in
[`internal/site/render.go`](../../../../internal/site/render.go);
the render-side test surface in
[`internal/site/forest_test.go`](../../../../internal/site/forest_test.go).
No spec, agent rule, or backlog edit.

## Goal

Implement the m2 milestone's mode-affordance contract for the
rendered forest: every plan-tree node either carries the
deterministically-derived **Begin planning** / **Begin
implementation** affordance — emitted as the locked D4 form
inside the always-visible summary region of the node's
`<details>` box — or carries no affordance, per the locked D3
static map. Land the supporting changes (the canonical-status
helper, the typed affordance value, the enumerated-map test) so
the static map has a single source of truth in code that future
edits cannot silently break.

## Inherited contract (consumed by reference)

The locked WHAT for this task is the [m2 milestone doc's Task
Contracts row for t1](./README.md). The Cross-Task Invariants,
Cross-Task Decisions **D3** (the static map) and **D4** (the
form shape), and the Cross-Task Risk "Mode-affordance map drift"
named in that doc are this task's load-bearing premises and are
not restated here. The Contracts below add the HOW the milestone
deferred to this drafting session — they do not loosen, redefine,
or duplicate any premise locked at milestone level.

## Cross-Cutting Invariants

The m2 README binds invariants that thread the four-task set;
t1 reads and respects them — they are consumed by reference, not
re-asserted here. See the [m2 README "Cross-Task
Invariants"](./README.md) for the full list (no new registration
surface, fencing properties hold, single-local environment,
agent-adapter seam additive, "deterministic" stays epic-scoped,
no JavaScript / no walk-on-every-request regression).

One t1-local invariant the milestone does not lock:

- **The `Deferred — <reason>` canonical-prefix logic has exactly
  one source of truth in render-side code.** Any code path that
  needs to compare a Status string against the canonical
  lifecycle tokens (the existing `statusClass` and the new
  affordance lookup) consumes the same helper rather than
  inlining a parallel strip. The m2-level "Mode-affordance map
  drift" risk lives at milestone level; this invariant is the
  t1-specific mechanism that closes it at the render-side reuse
  site. *Verified by:*
  [`internal/site/render.go`](../../../../internal/site/render.go)
  `statusClass` (the existing inline em-dash split is the body
  that extracts cleanly into a shared helper).

## Contracts

### C1 — D3 implemented exactly, as a render-side lookup

The mode-affordance map locked in
[m2 Cross-Task Decision D3](./README.md) is implemented as a
deterministic per-node lookup whose inputs are exactly the three
already-populated `PlanNode` fields D3 keys on: `NodeType`,
`Status`, and the presence-or-absence of items in `Children`.
The output is one of three values: no affordance, the "Begin
planning" affordance, or the "Begin implementation" affordance.
The lookup runs at render time over fields the per-request walk
has already populated, so it introduces no new walk and no new
DB read — the m2 Cross-Task Invariant "No JavaScript / no walk-
on-every-request regression" binds. The Status comparison runs
against the canonical-prefix form, not the raw Status — a
`Deferred — <reason>` node maps to the canonical `Deferred` row
of D3, not to the unknown-Status fall-through (this is what C4
binds at the render-side reuse site). *Verified by:*
[m2 Cross-Task Decision D3](./README.md);
[`internal/site/tree.go`](../../../../internal/site/tree.go)
`PlanNode` struct (the three fields the lookup keys on are
already populated by `buildTree` before any render call).

### C2 — One value drives both the form's `mode` input and the button label

The render-side representation of the affordance is shaped so
the same value the template emits as the form's hidden `mode`
input is the value it emits as the button's visible label —
there is no second mapping from "mode value" to "button label,"
so a future edit cannot let the two drift. The D3 strings
("Begin planning" / "Begin implementation") are exact-match-
checked downstream — t2's `/spawn` handler reads the `mode`
form field and dispatches per its value, so paraphrasing the
constant in one site without the other would break the
contract. *Verified by:*
[m2 Cross-Task Decision D3](./README.md) (the resolved mode
strings); [m2 Cross-Task Decision D4](./README.md) (the form
emits exactly `name="slug"` and `name="mode"` hidden inputs);
[`spec/planning/shared.md`](../../../../spec/planning/shared.md)
"Quote labels whose enforcement depends on exact-match matching"
(the discipline this contract realizes).

### C3 — D4 form emitted exactly, with the affordance value as both label and mode

When the affordance method returns a non-sentinel mode, the
template emits the form shape locked at D4: an HTML `<form>`
element with `method="POST"` and `action="/spawn"`, exactly two
hidden `<input>` elements (`name="slug"` carrying the node's
slug, `name="mode"` carrying the affordance constant), and one
`<button type="submit">` whose visible text is the same
affordance constant. When the affordance method returns the
sentinel, the template emits **no form element at all** for that
node — not an empty form, not a disabled form, not a hidden
form. *Verified by:*
[m2 Cross-Task Decision D4](./README.md) (the exact form
elements t1's render is the page-side surface of).

### C4 — Canonical-prefix logic is shared, not duplicated

`render.go` carries one helper that returns the canonical-prefix
portion of a Status string (the part before ` — ` for a
`Deferred — <reason>` value, the whole string otherwise);
`statusClass` consumes that helper, and so does the affordance
lookup. There is no second em-dash split in t1's new code. This
realizes the t1-local invariant above and closes the m2-level
"Mode-affordance map drift" risk at the render-side reuse site.
*Verified by:*
[`internal/site/render.go`](../../../../internal/site/render.go)
`statusClass` (the existing em-dash split is the body the
helper extracts; rewriting `statusClass` over the helper is
behavior-preserving).

### C5 — Affordance renders in the always-visible summary region

The form, when emitted, sits inside the node's `<summary>` flex
row — alongside the existing label-group and status-group flex
children — so the affordance is visible whether the
`<details>` box is open or closed. Nothing about the form's
emission renders only inside `<div class="box-body">` (which
HTML5 hides on a closed `<details>`). This realizes the m2
product-acceptance constraint "the affordance is visible
without expanding the node." *Verified by:*
[`internal/site/forest.go`](../../../../internal/site/forest.go)
`forestTemplates` (the `define "node"` block places
`node-header` inside `<summary>` and `node-detail` inside
`<div class="box-body">`; the `<details>` element hides its
body until the `open` attribute is set);
[m2 Task Contracts row for t1](./README.md) (product acceptance:
"the affordance is visible without expanding the node").

### C6 — Unknown Status falls through to no affordance

A node whose Status value, after canonical-prefix stripping, is
not one of the six canonical tokens D3 enumerates
(`In draft`, `Proposed`, `In progress`, `Validating`, `Landed`,
`Deferred`) emits no affordance. The empty-Status case is
included in this fall-through unless D3's "In draft / no doc"
branch maps it explicitly (per D3 it does, when the node is a
leaf-shape task or phase). The graceful-fall-through behavior
matches `statusClass`'s `unknown` fallback per the spec's
graceful-render rule for unrecognized Status values.
*Verified by:*
[`spec/planning/shared.md`](../../../../spec/planning/shared.md)
"Plan-doc Status" (unknown values render gracefully — t1
extends this from the badge surface to the affordance surface);
[m2 Cross-Task Decision D3](./README.md) (the explicit "or any
unknown Status → no affordance" clause).

### C7 — Existing rendered surfaces preserved

Every node still renders its existing label, Status badge,
active-work actor markers, progress-cell row, long description
(line-capped per the existing rule), and related-PR list — none
of those are modified. The forest's expand/collapse default
(`ActiveInSubtree`-driven), the no-JavaScript posture, and the
walk-on-every-request invariant are unchanged. The form's
insertion may necessarily reflow the header flex row (a new flex
child between or after existing ones); reflow is not
modification. *Verified by:*
[`internal/site/forest.go`](../../../../internal/site/forest.go)
`forestTemplates` (the preserved-surface definitions in
`node-header`, `node-progress`, `node-detail`, the
`ActiveInSubtree`-driven `open` attribute);
[`internal/site/forest_test.go`](../../../../internal/site/forest_test.go)
(the standing tests that pin label, badge, actor-marker,
progress-row, long-description, and related-PR rendering — t1's
test additions extend the suite without retiring any).

## Files to touch (estimate)

Estimate of the expected change shape per
[`spec/planning/shared.md`](../../../../spec/planning/shared.md)
"Plan content is a mix of rules and estimates — label which is
which." Implementation may revise this when a structural call
requires deviating; deviations are surfaced per the "Estimate
Deviations" callout rule.

**New files:** none expected — every change rides existing
files.

**Modify:**

- [`internal/site/tree.go`](../../../../internal/site/tree.go) —
  add the named-string affordance type, its constant set, and
  the per-node method (`PlanNode.Affordance()` returning the
  affordance value) implementing D3.
- [`internal/site/render.go`](../../../../internal/site/render.go) —
  extract the shared canonical-status helper; rewrite
  `statusClass` over it.
- [`internal/site/forest.go`](../../../../internal/site/forest.go) —
  template surface: emit the D4 form inside the summary's
  `node-header` (or a small sibling sub-template invoked from
  `node-header`; the exact sub-template factoring is HOW the
  implementing PR picks per the
  [`shared.md`](../../../../spec/planning/shared.md) "Plans
  describe contracts, not implementation" rule); add the
  affordance's CSS rules to `forest-style`.
- [`internal/site/forest_test.go`](../../../../internal/site/forest_test.go) —
  add the enumerated-map table-driven test (the Validation
  Gate's primary falsifier per the m2 Cross-Task Risk
  "Mode-affordance map drift").

**Intentionally not touched:**

- [`cmd/workstream-tracker/main.go`](../../../../cmd/workstream-tracker/main.go)
  `runServer` — the loopback-only binding tightening is t2's
  surface per the m2 Task Contracts row for t2, not t1's.
- [`internal/site/site.go`](../../../../internal/site/site.go)
  `Router()` — the `POST /spawn` handler mount is t2's surface.
- [`docs/agents/local/session-registration.md`](../../../../docs/agents/local/session-registration.md) —
  any additive clarification is t3's surface per m2 Documentation
  Currency.
- [`.claude/settings.json`](../../../../.claude/settings.json) —
  the SessionStart hook entry is t3's surface.
- [`design/v0.1-design.md`](../../../../design/v0.1-design.md) —
  frozen v0.1 end-state record; not reconciled against (per
  m2 Documentation Currency "Currency check, no edit expected").

## Validation Gate

The technical gate t1's implementing PR walks before merge. t1
is an **interior node** in the m2 Mermaid graph (it has an
outgoing edge to t4), so per
[`spec/planning/task-plan.md`](../../../../spec/planning/task-plan.md)
"Product-facing leaf tasks output a demo walkthrough" t1 closes
on this technical gate alone — the milestone's product
walkthrough lives on the t4 leaf and is not duplicated here.

1. **Enumerated-map test (the primary falsifier).** A new
   table-driven test in
   [`forest_test.go`](../../../../internal/site/forest_test.go)
   carries one row per production-reachable (NodeType × Status ×
   has-children) triple D3 distinguishes — covering each
   NodeType the renderer can produce (`root`, `milestone`,
   `task`, `phase`) at the Status values D3 enumerates for that
   NodeType, plus at least one unknown-Status row and at least
   one `Deferred — <reason>` row to pin canonical-prefix
   behavior, plus the task-with-children case (no affordance).
   D3's product-level **epic** category maps to `NodeType:
   "root"` at the renderer level — per
   [`internal/slugs/slugs.go`](../../../../internal/slugs/slugs.go)
   `parseSegment` and `Slug.NodeType()`, the slug parser never
   assigns `NodeTypeEpic`; epic-rooted docs and standalone-task
   roots both render as `NodeType: "root"`, so the `root` test
   row covers the epic case and no separate `epic` row is
   required (the unreachable case is documented in the
   coverage list above for the same structural reason as
   phase-with-children).
   Per-row fixture shape is implementation choice: cases whose
   test node is a structural root need only their own
   `parsedDoc`; deeper cases need the parsedDoc chain of
   intermediate parents because
   [`tree.go`](../../../../internal/site/tree.go) `buildTree`
   silently drops a node whose parent slug is absent from the
   slice; the with-children case additionally needs at least one
   child phase doc so `buildTree`'s wire-children pass populates
   the test node's `Children`. The phase-with-children case is
   **not** in the enumeration — per
   [`slugs.go`](../../../../internal/slugs/slugs.go) `pN` is the
   terminal segment in the slug grammar, so a well-formed
   parsedDoc set cannot construct a phase with children and that
   branch of D3 is structurally unreachable in production. Each
   row renders via `renderTree` and asserts the rendered HTML
   either contains exactly the expected affordance form shape
   (the `<form>` element with the hidden inputs and the button
   labeled per the mode) **or** contains no `<form>` element at
   all, matching D3's row outcome. The assertion distinguishes
   positive-affordance from negative-affordance cases on the
   same render output (the `<form>` element's presence is the
   unambiguous discriminator), so the falsifier-check rule's
   "multiple causes produce the same observation" failure mode
   does not apply.
2. **`go test ./internal/site/...` clean.** The standing suite
   (every preserved-surface test enumerated in C7) keeps passing
   alongside the new enumerated-map test.
3. **Manual single-node render check.** Run `go run
   github.com/kcrobinson-1/workstream-tracker/cmd/workstream-tracker`
   against the local checkout's plan tree; open the rendered
   page in a browser; visually confirm at least one
   leaf-shape task in `Proposed` Status renders the "Begin
   implementation" button visible without expanding the box, at
   least one in `In draft` renders the "Begin planning" button,
   and at least one parent-shape task (carrying phase children)
   renders no button. The manual check is the t1-local sanity
   pass against the rendered page — it is **not** the m2
   product walkthrough (which lives on t4 against a live Claude
   Code launcher and the full spawn/hook chain).

## Self-Review Audits

Audits the implementer walks before opening the PR, drawn from
the diff surfaces this plan touches.

- **Canonical-prefix single-source audit.** Search
  [`internal/site/`](../../../../internal/site/) for any
  remaining inline `" — "` em-dash split after the helper
  extraction in C4. Any second site is the t1-local Cross-
  Cutting Invariant violated.
- **No-JS / no-new-walk audit.** Search the diff for any new
  `<script>` element, any new client-side event handler attribute
  (`onclick`, `onsubmit`, etc.), any new DB read or
  per-request file walk introduced in the render path. None
  should appear.
- **Preserved-surfaces audit.** Run the standing
  forest_test.go tests for label, badge, actor-marker, progress-
  row, long-description, related-PR rendering; confirm each
  still passes against the touched template. Specifically
  confirm `TestRenderForestActorMarkerNameBearing`,
  `TestRenderProgressRowCoexistsWithPreservedSurfaces`, and
  `TestRenderSummaryCarriesExplicitMarkerTreatment` still pass
  — these are the cross-cutting cohabitation pins for the
  summary region the new affordance shares.
- **Sibling-not-yet-shipped degrade audit.** Confirm the form's
  `action="/spawn"` produces a 404 (not a 500, not a silent
  failure) when t2's handler is not yet mounted, by submitting
  the form against the current `Router()` (which exposes only
  `GET /`). The 404 is the expected sibling-not-yet-shipped
  degrade the m2 Sequencing anticipates; a 500 or silent
  failure would mask the missing sibling.

## Documentation Currency

Per the [m2 Documentation Currency](./README.md) clause for t1,
this task edits the forest template surface
([`forest.go`](../../../../internal/site/forest.go)) and the
render-side test surface
([`forest_test.go`](../../../../internal/site/forest_test.go)),
plus the supporting computed-value file
([`tree.go`](../../../../internal/site/tree.go)) and the
render helper
([`render.go`](../../../../internal/site/render.go)). No
status-bearing doc under
[`spec/**`](../../../../spec/) or
[`docs/agents/local/**`](../../../../docs/agents/local/) is
touched (the mode-offer map is a product surface, not a contract
under those rule trees). The frozen
[`design/v0.1-design.md`](../../../../design/v0.1-design.md) is
not reconciled against, per the project's standing posture.

## Out of Scope

- **The `POST /spawn` handler, its argv shape, and the loopback-
  only binding tightening.** All of these are t2's surface per
  the m2 Task Contracts row for t2.
- **The committed `.claude/settings.json` SessionStart hook
  entry and any additive clarification to
  [`session-registration.md`](../../../../docs/agents/local/session-registration.md).**
  Both are t3's surface per the m2 Task Contracts row for t3.
- **The end-to-end product walkthrough.** Lives on t4 against
  a live local server with a real Claude Code launcher; t1's
  manual check in the Validation Gate is the t1-local sanity
  pass, not the walkthrough.
- **CSRF / origin / token surfaces in the form.** t1 emits a
  minimal D4-conformant form per SD4; the trust boundary is
  t2's loopback-only binding per the m2 Cross-Task Risk
  "New `/spawn` endpoint expands the tool's write surface."
- **Per-mode prompt-file content.** Per
  [m2 Out of Scope](./README.md) "The exact wording of the
  planning / implementation prompts," the prompt wording is
  settled in t2's implementing PR; t1's render carries no
  prompt content.
- **Any backlog edit.** Per the [m2 Backlog
  Impact](./README.md) lock, m2's implementing PRs touch no
  backlog entries.
- **Any tree-side observability heuristic** (a "node has a plan
  doc but no work-instance" indicator). Per the [m2 Out of
  Scope](./README.md) clause, that lives in the
  [`unregistered-work-unobservable`](../../../../docs/backlog.md#unregistered-work-unobservable)
  backlog entry.

## Risk Register

t1-local risk mitigations carrying t1's share of the m2-level
risks plus any t1-specific residuals.

- **Map drift (the t1-side of the m2 Cross-Task Risk
  "Mode-affordance map drift").** A missed combinator (an
  unknown Status whose canonical-prefix the affordance lookup
  doesn't share with `statusClass`, a task-with-phases case the
  has-children check misclassifies, a node-type the affordance
  method's switch silently miscategorizes) could surface an
  affordance on the wrong node. Mitigation: (a) Contract C4
  binds the canonical-prefix single source; (b) the enumerated-
  map test in the Validation Gate is the falsifier the m2 Risk
  Register names — every D3 row is a fixture; (c) the canonical-
  prefix single-source audit in Self-Review explicitly searches
  the diff for a parallel inline strip.
- **Browser-default click bubbling from the submit button to
  the parent `<summary>`.** When the contributor clicks the
  button, the browser may also toggle the `<details>` open
  state before the form's POST navigates. Under merged-code
  behavior the toggle effect is harmless because the navigation
  supersedes any local visible change; no D4 or no-JS posture
  constraint is violated either way. Mitigation: **accepted
  under the project's standing accepted-failure-with-visibility
  posture** — the plan adds no contract on bubbling, per the
  paired scoping doc's
  [SD9](./scoping/t1-mode-affordance-render.md). The HTML spec's
  activation-behaviour clause for `<summary>` already covers the
  no-bubble case for modern Chromium/Firefox; when the rare
  bubble does happen, the form's POST navigation supersedes any
  visible toggle before the user perceives it.
- **CSS regression on the narrow-window degrade.** The added
  flex child in `.box-header` could regress the post-m2-ux-
  correction-p1 stacked-left intent
  ([`d39e2ba`](https://github.com/kcrobinson-1/workstream-tracker/commit/d39e2ba)).
  Mitigation: the placement contract in C5 binds visible-in-
  summary; SD7 scopes trailing-position + stacks-below in the
  narrow degrade; the manual single-node render check in
  the Validation Gate exercises the rendered header at both
  default and narrow widths.

## Related Docs

- [`./README.md`](./README.md) — m2 milestone doc; the locked
  WHAT, Cross-Task Invariants/Decisions (especially D3 and D4),
  Cross-Task Risks (especially "Mode-affordance map drift"),
  Documentation Currency.
- [`./scoping/t1-mode-affordance-render.md`](./scoping/t1-mode-affordance-render.md) —
  paired scoping doc carrying SD1…SD8 with rejected alternatives
  and the reality-check inputs; transient, deletes at t4's
  milestone-terminal PR.
- [`./t2-spawn-endpoint.md`](./t2-spawn-endpoint.md) — sibling
  task; consumes t1's affordance form submission as the launch
  trigger.
- [`./t4-end-to-end-validation.md`](./t4-end-to-end-validation.md) —
  the milestone-graph leaf t1 converges on; carries the product
  walkthrough that observes t1's affordance from a reviewer's
  vantage.
- [`../README.md`](../README.md) — parent epic; Cross-Cutting
  Invariants and the resolved vision inputs the milestone
  contract sits on top of.
- [`../../../../spec/planning/task-plan.md`](../../../../spec/planning/task-plan.md) —
  the rules this task plan is structured against; "Required and
  optional sections," "Plan-to-PR Completion Gate," the
  `In draft` → `Proposed` promotion gate.
- [`../../../../spec/planning/shared.md`](../../../../spec/planning/shared.md) —
  cross-level rules; "Plans describe contracts, not
  implementation," "Verified by: annotations," "Plan-doc
  Status."
- [`../../../../internal/site/tree.go`](../../../../internal/site/tree.go),
  [`../../../../internal/site/forest.go`](../../../../internal/site/forest.go),
  [`../../../../internal/site/render.go`](../../../../internal/site/render.go),
  [`../../../../internal/site/forest_test.go`](../../../../internal/site/forest_test.go) —
  the four files this task's estimate touches.
