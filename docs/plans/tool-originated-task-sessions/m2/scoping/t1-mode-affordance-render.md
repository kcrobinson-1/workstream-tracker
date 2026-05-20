---
slug: tool-originated-task-sessions-m2-t1
short_description: Scoping — mode-affordance render on plan-tree nodes
---

# Scoping — m2 t1 Mode-Affordance Render

Scoping doc paired with the
[`t1` task plan](../t1-mode-affordance-render.md). Transient per
[`task-plan.md`](../../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns" — deletes in batch at the
milestone-terminal PR (t4's PR, not t1's). Carries no `Status`
field per the same rule (the lifecycle is for plans, not for
transient scoping).

## Context

t1 implements the m2 milestone's **mode-affordance render**: per
the m2 milestone doc's locked **D3** static map and locked **D4**
form-shape contract, every plan-tree node in the rendered forest
either emits a small HTML form that POSTs to t2's future `/spawn`
endpoint, or emits no form, deterministically driven by the
node's already-rendered facts. The render side is the only
surface t1 touches; t2 (spawn endpoint) and t3 (SessionStart
hook) are independent surfaces.

This scoping doc lays out the implementation decisions the plan
locks. The locked WHAT comes from the
[m2 milestone doc](../README.md) and is not re-litigated here.

## Reality-check inputs (the plan's load-bearing premises)

Each item below is a code or vendor-doc fact the plan rests on.
Each is verified by reading the cited surface before drafting,
not by transitive citation.

- **`PlanNode` carries the three fields D3 keys on.** `NodeType`
  is set by `buildTree` from `slugs.Parse` (defaulting to
  `slugs.NodeTypeRoot` when the slug is unparseable); `Status` is
  the raw frontmatter value (empty string when the doc carries no
  `Status` field); `Children` is wired in the same pass before
  any post-order walk. *Verified by:*
  [`internal/site/tree.go`](../../../../../internal/site/tree.go)
  `PlanNode` struct definition and `buildTree` (field
  population);
  [`internal/site/forest_test.go`](../../../../../internal/site/forest_test.go)
  `TestRenderLeafAndStubBox` (a `Status: In draft` stub renders
  as a leaf box with its raw Status string visible).
- **The slug grammar's node types are exactly the five D3 keys
  on.** `slugs.NodeType` constants are `NodeTypeRoot`,
  `NodeTypeEpic`, `NodeTypeMilestone`, `NodeTypeTask`,
  `NodeTypePhase`. The map matches by string equality on
  `PlanNode.NodeType` (which carries the stringified node-type),
  so no additional grammar work is needed. *Verified by:*
  [`internal/slugs/slugs.go`](../../../../../internal/slugs/slugs.go)
  `NodeType` constants and `Slug.NodeType()`.
- **`statusClass` already does the `Deferred — <reason>`
  canonical-prefix strip.** The strip looks for ` — ` (space,
  em-dash, space) and truncates to the prefix. Any new
  Status-keyed code path that needs the canonical prefix must
  reuse this logic — the m2 risk register names "map drift" as
  the failure mode of re-implementing it. *Verified by:*
  [`internal/site/render.go`](../../../../../internal/site/render.go)
  `statusClass` (em-dash split before the lifecycle switch).
- **The forest template structure is `node` → (`summary` with
  `node-header`) + (`box-body` with `node-detail` and child
  `node` recursion).** `node-detail` is the body-rendered region;
  `node-header` is the summary-rendered region (always visible
  even when the `<details>` box is closed). Only content inside
  `node-header` satisfies m2's "the affordance is visible
  without expanding the node" constraint. *Verified by:*
  [`internal/site/forest.go`](../../../../../internal/site/forest.go)
  `forestTemplates` (the `define "node"` block, the `define
  "node-header"` block, the `define "node-detail"` block, the
  `<details>` element rendered with the header inside `<summary>`
  and the rest inside `<div class="box-body">`).
- **The forest is rendered server-side once per GET request with
  no JavaScript and one DB read.** The page walks the plan tree
  + reads active work-instances + reads session metadata, then
  renders. No JS file is referenced, the chi router exposes only
  `GET /`. *Verified by:*
  [`internal/site/render.go`](../../../../../internal/site/render.go)
  `indexTmpl` (the `<head>` carries no `<script>`);
  [`internal/site/site.go`](../../../../../internal/site/site.go)
  `Router()` (only `r.Get("/", s.index)`);
  [`internal/site/site.go`](../../../../../internal/site/site.go)
  `index` handler (the per-request walk + DB reads — the
  walk-on-every-request invariant the m2 Cross-Task Invariants
  bind).
- **`<details>` elements honor in-`<summary>` form submission
  without script.** A submit button inside `<summary>` activates
  the form per HTML5 semantics; the form's POST navigates the
  page. Browser-default click bubbling from the button to the
  `<summary>` is unobservable to the user because the navigation
  takes precedence over any local toggle effect. *Verified by:*
  [HTML Living Standard, "The summary
  element"](https://html.spec.whatwg.org/multipage/interactive-elements.html#the-summary-element)
  ("activation behaviour" of `summary` toggles `details` only
  when no inner interactive element handled the activation; form
  submit is an activation behaviour of the submit button).
  **OQ1** below tracks the open question of whether a stricter
  guarantee is needed here.
- **`t2`'s `POST /spawn` endpoint does not yet exist on `main`.**
  The form's `action="/spawn"` points at a route that t2 will
  add; a form that POSTs before t2 lands receives a 404 — the
  expected sibling-not-yet-shipped degrade per the m2 Sequencing
  ("three independent surfaces"). *Verified by:*
  [`internal/site/site.go`](../../../../../internal/site/site.go)
  `Router()` (no `r.Post` line; only `r.Get("/", s.index)`).
- **m2 D3 is the locked static map; m2 D4 is the locked form
  shape.** Both are decision-complete in the m2 milestone doc and
  consumed verbatim. *Verified by:*
  [m2 README, Cross-Task Decisions D3 and D4](../README.md).

## Decisions made at scoping time

Each decision below decomposes its option space into shapes per
[`shared.md`](../../../../../spec/planning/shared.md) "Decompose
options into shapes before analyzing," names the chosen shape,
and names the rejected alternatives with rationale.

### SD1 — Form-emission site: inside the node's `<summary>` flex row

Where in the forest template the affordance renders.

**Shapes considered:**

- **(a) Inside `<summary>`, in the `node-header` flex row** —
  always visible (the `<summary>` renders regardless of
  `<details>` open state), no expansion required.
- **(b) Inside `node-detail`** (under `<div class="box-body">`)
  — hidden when the `<details>` is closed, since native
  `<details>` hides the body until the open attribute is set.
- **(c) Outside `<details>` entirely, as a sibling element in
  the box wrapper** — would orphan the affordance from its node
  visually and require a wrapper-element refactor.
- **(d) Inside `<details>` but outside both `<summary>` and
  `<div class="box-body">`** — HTML5 hides every child of a
  closed `<details>` except `<summary>`, so this collapses to
  the same hidden behavior as (b).

**Chosen: (a).** m2's product acceptance contract
("the affordance is visible without expanding the node") forces
the affordance into the summary-rendered region. (b)/(d) hide
the affordance behind expansion; (c) requires structural surgery
on the node template's outer shape that the WHAT does not ask
for. Inside `<summary>`, the affordance sits adjacent to the
existing label-group / status-group flex children, reusing the
header's existing flex layout. *Verified by:*
[`internal/site/forest.go`](../../../../../internal/site/forest.go)
`forestTemplates` (`define "node-header"` is the
summary-rendered fragment; `define "node-detail"` is the
body-rendered fragment).

### SD2 — Static-map encoding: typed method on `PlanNode`

How D3's lookup is encoded in Go.

**Shapes considered:**

- **(a) Method on `PlanNode`** returning a typed enum value
  carrying "no affordance / planning / implementation," accessed
  from the template as an ordinary field-style template lookup
  on the receiver. Pure function of three fields already on
  `PlanNode`.
- **(b) Template `FuncMap` predicate** taking a `*PlanNode` and
  returning the same enum value, registered as a template
  function and called by name.
- **(c) Field on `PlanNode`** computed once in `buildTree`,
  matching the precedent of `ActiveInSubtree`.
- **(d) Free function in a new file**, called from the template
  via a `FuncMap` registration.

**Chosen: (a).** The lookup is a pure deterministic function of
three already-populated `PlanNode` fields with no tree-walk
dependency, so the precedent that justifies `ActiveInSubtree`'s
buildTree-computed field (a post-order tree fold) does not apply
here — there is nothing to fold. A method on the receiver
expresses the per-node lookup directly, tests via direct method
calls (no template render needed for the unit-level enumeration),
and reads naturally in the template as a receiver field-style
lookup. (b) and
(d) work but route through `FuncMap`, which is heavier
template-side machinery for a per-node lookup. (c) adds field
state for no reason since the value is cheap to recompute and
not shared across nodes. *Verified by:*
[`internal/site/tree.go`](../../../../../internal/site/tree.go)
`PlanNode` (existing receiver-method-friendly value type;
`ActiveInSubtree` is the post-order fold that justifies a field,
which the affordance lookup is not).

### SD3 — Affordance type: a typed enum carried as a string-valued constant set

The shape of the value the method returns.

**Shapes considered:**

- **(a) Named string type** with constants for "no affordance,"
  "Begin planning," "Begin implementation." Template renders the
  value directly when emitting the form button's label and the
  hidden `mode` input's value.
- **(b) Struct** carrying separate "mode value" and "button
  label" fields.
- **(c) Bare bool + label** — a boolean "has affordance" and a
  separate label string.

**Chosen: (a) — a named string type whose constant set is the
three D3 outcomes.** The constant string IS both the hidden
`mode` input's value and the button's visible label per D4 and
D3 (D3 names the mode strings exactly as "Begin planning" /
"Begin implementation"). Folding them into one value removes a
parallel mapping that drift could open between mode and label.
(b) creates a second mapping with no source-of-truth basis. (c)
loses the no-affordance/planning/implementation discrimination
at the type level. *Verified by:*
[m2 D3](../README.md) (the resolved mode strings are exactly the
button labels and exactly the mode values D4's form carries).

### SD4 — Form HTML shape: minimal D4-conformant form, no extra attributes

The form's HTML structure.

**Shapes considered:**

- **(a) Minimal D4 form** — `<form method="POST" action="/spawn">`,
  two hidden inputs (`name="slug"`, `name="mode"`), one
  `<button type="submit">` whose visible text is the affordance
  label.
- **(b) D4 form + visible mode dropdown** — instead of a fixed
  per-node mode, a `<select>` lets the user pick between the
  modes. Rejected: D3 maps each node to one specific mode (or
  none); a dropdown contradicts the static map.
- **(c) D4 form + visible slug input** — exposes the slug as an
  editable text field. Rejected: the slug is identity, carried
  by construction; a contributor-editable slug would re-introduce
  the slug-resolution problem the epic dissolves.
- **(d) D4 form + CSRF / origin tokens** — adds a hidden CSRF
  token. Rejected at t1's level: the trust boundary is t2's
  loopback-only binding (m2's Cross-Task Risk on
  `/spawn`-write-surface); t1 emits no token because t2's
  defense-in-depth contract owns the verification surface, and
  adding a token at t1 without t2 to consume it is dead surface.

**Chosen: (a).** D4 specifies the form shape; t1 emits exactly
that, no more. The form is keyboard-operable and screen-reader
friendly by virtue of being a standard HTML `<button
type="submit">` with its visible text as the accessible name —
no `aria-label` is needed. *Verified by:* [m2 D4](../README.md)
(the explicit form structure: hidden `slug` + `mode`, submit
button labeled per the mode).

### SD5 — Canonical-prefix reuse: extract a shared helper from `statusClass`

How the `Deferred — <reason>` prefix-strip is shared between
`statusClass` (existing) and the affordance lookup (new).

**Shapes considered:**

- **(a) Extract a shared helper** (e.g., `canonicalStatus(string)
  string`) in `render.go` that returns the prefix portion of a
  Status value, and rewrite `statusClass` to call it. The
  affordance lookup calls the same helper.
- **(b) Inline a parallel em-dash split in the affordance
  lookup**, mirroring `statusClass`'s existing inline split at a
  second site. Rejected: this is exactly the "map drift" failure
  mode the m2 Risk Register names (two sites for one rule,
  silently diverging).
- **(c) Key the affordance lookup on the result of `statusClass`
  (the CSS class string)**, not on the raw Status. Rejected:
  couples a render-side concern (CSS class naming) to a contract
  decision (the affordance map), inverting the dependency. A
  future class rename would silently change the affordance
  behavior.

**Chosen: (a).** The shared helper is the single source of truth
for the prefix-strip; both callers consume it. Implementation is
a small refactor inside `render.go` — `statusClass`'s body
becomes a one-liner over the helper's result. *Verified by:*
[`internal/site/render.go`](../../../../../internal/site/render.go)
`statusClass` (the existing inline strip is the body that
extracts cleanly into a helper without behavior change).

### SD6 — Has-children discrimination: render-time check, not a precomputed field

Whether "has children" is a render-time check or a precomputed
field.

**Shapes considered:**

- **(a) Render-time** — the affordance lookup reads the
  presence-or-absence of items in `PlanNode.Children` directly
  on every invocation. `Children` is wired in `buildTree` before
  any rendering, so the read is always safe.
- **(b) Precomputed field** on `PlanNode` set in `buildTree`,
  alongside `ActiveInSubtree`. Rejected: redundant with the
  slice's own state.

**Chosen: (a).** The slice's state is already-available; adding
a field duplicates information already representable. *Verified
by:*
[`internal/site/tree.go`](../../../../../internal/site/tree.go)
`buildTree` (the "Wire children to parents" pass populates
`Children` before any post-order walk and before any render).

### SD7 — CSS placement: trailing position in the header flex row, narrow-window stacked below

Where the rendered button visually sits in the always-visible
header.

**Shapes considered:**

- **(a) After the status-group**, at the far right of the
  default row-direction header. CTA placement convention
  ("rightmost is the action") and the status-group's existing
  `margin-left: auto` pin already creates a separator gap.
- **(b) Between label-group and status-group** — would push the
  status badge away from its current right-edge position.
- **(c) Before the label-group**, at the far left — visually
  unusual and conflates the action with the node's identity
  markers.

**Chosen: (a).** Trailing position preserves the existing label
→ markers → badge reading order; the button reads as a per-row
action. In the narrow-window degrade (the existing
`max-width: 48rem` block that flips `.box-header` to
`flex-direction: column`), the button stacks below the status
badge per the same column rules. The styling specifics — the
button's padding, border, color, hover state — are HOW the
implementing PR settles; t1's plan constrains the *placement
contract* (visible in summary, trailing position, does not
regress the post-m2-ux-correction-p1 stacked-left narrow-window
intent — the PR
[#57](https://github.com/kcrobinson-1/workstream-tracker/pull/57)
behavior the
`d39e2ba` commit pins), not the specific border-radius. *Verified
by:*
[`internal/site/forest.go`](../../../../../internal/site/forest.go)
`forestTemplates` (`.box-header` flex layout; the
`@media (max-width: 48rem)` block carrying the C4 stacked-left
intent that
[`d39e2ba`](https://github.com/kcrobinson-1/workstream-tracker/commit/d39e2ba)
restored after a regression on PR #57).

### SD8 — Validation gate: enumerated-map test as the falsifier

The shape of the test that catches "map drift" per the m2 Risk
Register.

**Shapes considered:**

- **(a) Enumerated table-driven test** with one row per D3 case
  (the five categories D3 distinguishes), each row asserting
  the rendered HTML for a fixture node contains exactly the
  expected affordance — or no form at all. The test renders via
  the existing `renderTree` helper, so the assertion exercises
  the real template path the production code takes.
- **(b) Unit-test the affordance method directly** with one
  case per D3 row, without touching the template. Rejected
  alone: tests the method in isolation but does not catch a
  template wiring bug (e.g., the template forgets to consult the
  method).
- **(c) Both (a) and (b)** — duplicates coverage; the table-
  driven render assertion subsumes the unit assertion because
  the rendered output reflects both the method's return and the
  template's faithful emission.

**Chosen: (a).** A single table-driven test, one row per
production-reachable (NodeType × Status × has-children) triple
D3 distinguishes, rendering via `renderTree` (the existing
helper used across forest_test.go) and asserting on rendered
HTML — same shape as the existing
[`TestRenderProgressRowDeclaredStages`](../../../../../internal/site/forest_test.go)
and
[`TestRenderProgressRowFieldlessOnlyDrafting`](../../../../../internal/site/forest_test.go)
table tests. Per-row fixture shape is implementation choice: a
case whose test node is a structural root needs only its own
parsedDoc; a deeper case needs the chain of intermediate
parents (the milestone's root, the task's milestone-and-root,
etc.) because
[`internal/site/tree.go`](../../../../../internal/site/tree.go)
`buildTree`'s "If the parent isn't present (gap in the tree),
the node is silently dropped from the visualization" branch
otherwise drops the test node before render; a with-children
case additionally needs at least one child phase doc so
`buildTree`'s "Wire children to parents" pass appends it to the
test node's `Children`. Coverage: each NodeType the renderer
can produce (`root`, `milestone`, `task`, `phase`) at the
applicable Status values per D3, plus at least one unknown-Status
row and at least one `Deferred — <reason>` row to pin canonical-
prefix behavior. The `phase`-with-children case is not in the
enumeration — per
[`internal/slugs/slugs.go`](../../../../../internal/slugs/slugs.go)
`pN` is the terminal segment, so `buildTree` cannot construct a
phase with children from well-formed slugs and that branch of
D3 is structurally unreachable in production. *Verified by:*
[`internal/site/forest_test.go`](../../../../../internal/site/forest_test.go)
`renderTree` invocation pattern (the existing precedent for
table-driven render assertions);
[`internal/site/tree.go`](../../../../../internal/site/tree.go)
`buildTree` (the parent-chain and child-wiring requirements the
fixture shape depends on);
[`internal/slugs/slugs.go`](../../../../../internal/slugs/slugs.go)
`IsWellFormed` and the `pN`-is-terminal grammar (the constraint
that rules out the phase-with-children case); m2's Risk Register
entry "Mode-affordance map drift" naming the enumerated-map
test as the falsifier.

### SD9 — Click bubbling inside `<summary>`: accept the navigation-supersedes posture, no contract added

When the contributor clicks the submit button inside the node's
`<summary>`, the browser may *also* toggle the parent `<details>`
open/closed state on the same click before the form's POST
navigates the page. The HTML spec leaves the bubbling behavior
implementation-defined for nested interactive elements; in
practice the visible toggle either does not happen
(Chromium/Firefox stop the bubbling for activated inner
interactive elements) or happens but is immediately superseded
by the page navigation. Neither outcome violates D4 or the no-JS
posture.

**Shapes considered:**

- **(a) Leave the side-effect unstated.** The plan adds no
  contract on bubbling. If the visible toggle ever happens, the
  POST navigation supersedes any local visible change before the
  user perceives it.
- **(b) Explicitly contract "no observable side-effect on the
  box's open state."** Would require either a vendor-doc citation
  pinning the no-bubble guarantee for the relevant browsers, or
  a browser-behavior test harness outside the Go test surface.

**Chosen: (a).** Consistent with the project's standing
accepted-failure-with-visibility posture: a harmless visible
flicker that resolves itself before the user perceives it is not
a contract surface worth a test harness. (b) buys a contract
whose observable difference the navigation-supersedes behavior
already erases. *Verified by:*
[HTML Living Standard, "The summary
element"](https://html.spec.whatwg.org/multipage/interactive-elements.html#the-summary-element)
(activation behaviour of `summary` toggles `details` only when
no inner interactive element handled the activation; form-submit
is an activation behaviour of the submit button); the
reality-check input above carrying the same citation as the
upstream framing.

## Open decisions to make at plan-drafting

None remaining. SD9 above resolves the only open question the
scoping surfaced (the click-bubbling discriminator carried as
OQ1 into the plan's first In-draft pass; the user resolved (A)
at the promotion gate, folded into SD9 here and reflected in the
plan's Risk Register entry).

## Plan structure handoff

The paired plan doc was drafted and promoted to **Status:
Proposed** (the seeded skeleton in place at
[`../t1-mode-affordance-render.md`](../t1-mode-affordance-render.md)
was replaced in the drafting commit and flipped to `Proposed` at
the promotion-gate commit on this branch). Section structure as
shipped:
- A Context preamble naming what the plan covers, why it is
  being done now, and what surfaces it touches at the conceptual
  level (per
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "Plan opens with a plain-language context preamble").
- **Goal.**
- **Cross-Cutting Invariants** consumed by reference from the m2
  README; one t1-local invariant added (the canonical-prefix
  single-source rule, SD5).
- **Contracts** C1…C7 covering (i) D3 implemented exactly;
  (ii) D4 emitted exactly; (iii) summary-region placement;
  (iv) canonical-prefix shared with `statusClass`; (v) no-JS /
  no-new-walk; (vi) unknown-Status falls through to no
  affordance; (vii) preserves existing rendered surfaces.
- **Files to touch (estimate)** —
  [`internal/site/tree.go`](../../../../../internal/site/tree.go)
  (the affordance method and its named type) +
  [`internal/site/render.go`](../../../../../internal/site/render.go)
  (the shared canonical-status helper) +
  [`internal/site/forest.go`](../../../../../internal/site/forest.go)
  (the summary-region form emission and supporting CSS) +
  [`internal/site/forest_test.go`](../../../../../internal/site/forest_test.go)
  (the enumerated-map test). Labeled estimate per the
  "Plan content is a mix of rules and estimates" rule.
- **Validation Gate** — SD8's enumerated-map test, plus
  `go test ./internal/site/…` clean and a manual single-node
  render check against the rendered page (the manual check is
  not the product walkthrough; that lives on t4).
- **Self-Review Audits** — render-side reuse audit
  (canonical-prefix sharing); no-JS audit; existing-surface
  preservation audit.
- **Out of Scope** — anything t2 or t3 owns; CSRF / token
  surfaces; in-session steering; per-mode prompt-file content;
  any backlog edit.
- **Risk Register** — the t1-side mitigations of the m2-level
  "Mode-affordance map drift" risk (the falsifier enumerated
  test), plus any t1-specific risks scoping surfaces.
- **Backlog Impact** — none (m2 README locks this at milestone
  level).
- **Related Docs.**

## Related Docs

- [`../t1-mode-affordance-render.md`](../t1-mode-affordance-render.md) —
  the paired durable plan doc this scoping doc fuels.
- [`../README.md`](../README.md) — m2 milestone doc; the locked
  WHAT, D3, D4, Cross-Task Invariants, Cross-Task Risks
  (including the map-drift risk SD8 mitigates).
- [`../../README.md`](../../README.md) — parent epic; the
  Cross-Cutting Invariants and resolved vision inputs the m2
  contract sits on top of.
- [`../../../../../spec/planning/task-plan.md`](../../../../../spec/planning/task-plan.md) —
  this scoping doc's authority; "Scoping owns / plan owns,"
  "Reality-check pass before plan-drafting."
- [`../../../../../spec/planning/shared.md`](../../../../../spec/planning/shared.md) —
  cross-level rules; "Decompose options into shapes,"
  "Verified by: annotations," "Plans describe contracts, not
  implementation."
