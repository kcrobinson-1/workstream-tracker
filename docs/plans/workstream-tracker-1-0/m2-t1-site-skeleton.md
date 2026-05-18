---
slug: workstream-tracker-1-0-m2-t1
Status: Proposed
short_description: Two-region page shell — existing forest in the forest region, placeholder in the roster region
---

# t1 — Site skeleton (two-region shell)

## Context

The workstream-tracker page today renders one thing: the
plan-tree forest, stacked roots in a single centered
`max-width: 60rem` column. The approved m2 target
([`design/workstreams-view-m2.svg`](../../../design/workstreams-view-m2.svg))
is a two-region page — a tall plan-tree **forest** on the left
(~2/3) and a **session roster** on the right (~1/3), one page
scroll, the roster sitting at the top so it shows without
scrolling when the window is tall enough, and a too-narrow
window stacking the roster below the forest.

This task ships that page shell and nothing else. The forest
region gets the **existing** forest render moved into it
unchanged (same nodes, same Status badges, same actor markers,
same long-description and related-PR rendering); the roster
region gets a deliberate, intentional placeholder — an observed
"the roster lands later" state, not a blank div that reads as a
bug. It is being done first because it converts the
shared-page-shell contract that otherwise sits between t2
(forest enrichment) and t4 (roster) into a shipped artifact:
once the two named regions exist and are independently owned, t2
builds into the forest region and t4 into the roster region with
no shell-ownership or merge-order coordination between them. That
coordination-elimination is t1's independent value.

Surfaces touched, conceptually: only the website's single page
template and its CSS. No API, DB, schema, parser, walker, or
tree-builder surface is involved, and the per-request data path
is untouched — t1 reshapes layout markup around render output
that is byte-for-intent the same as today.

This is an **N = 1** task plan; phase content is absorbed
inline. It is a **narrow-surface** plan and deliberately skips
the separate scoping doc per
[`task-plan.md`](../../../spec/planning/task-plan.md)
"Narrow-surface plans may skip the scoping doc": it touches a
single subsystem (one UI template + its CSS), an estimated 1
code file, introduces no new public-API contract, introduces no
new cross-cutting invariant (the shell/forest-verbatim/no-JS
rules it honors are **inherited** from the parent milestone, not
introduced here), and uses no novel mechanism (CSS-only
two-column layout with `html/template`, both already in the
codebase). All five narrow-surface conditions hold, so the
planner explicitly invokes the carve-out; the reality-check
inputs the gate protects are absorbed inline in the
**Reality-check inputs** section below (plus per-contract
`Verified by:` citations), not dropped.

The [`task-plan.md`](../../../spec/planning/task-plan.md)
`In draft → Proposed` promotion-gate self-review has been run:
read end-to-end for cross-section coherence; Contracts walked
for deferral phrases (the only deferrals — exact CSS technique
and the placeholder's final wording — are mechanism / render-time
UX copy authorized by "the contract is the behavior, not the
mechanism" and "Bans on surface require rendering the
consequence," not deferrals to plan-drafting itself); the
broadened `Verified by:` rule applied to every load-bearing
claim; reality-check inputs re-confirmed against the current
branch; required sections present with the narrow-surface
Reality-check inputs section disclosed here; no content
descended to implementation prescription. N = 1 (PR-count
branch test: one code file — `render.go` — plus tests and doc
currency, well under the split threshold), so no phase skeletons
are seeded. No open inputs remained, so Status is `Proposed`.

## Goal

Opening the page renders the approved two-region layout: a
forest region (~2/3, left) carrying the **existing** plan-tree
forest render verbatim, and a roster region (~1/3, right)
carrying a deliberate, intentional placeholder. One page scroll
(the tall forest is what scrolls); the roster region is
top-aligned so it is visible without scrolling when the window
height allows. Below a desktop-narrow width threshold the roster
region stacks **below** the forest region (mobile out of scope).
v0.1's per-node actor tags, long descriptions, related-PR lists,
the empty-state ("no roots found"), and the
walk-on-every-request render path are all unchanged — only the
layout frame around the forest is new, plus the roster
placeholder.

## Reality-check inputs

The narrow-surface carve-out compresses the scoping doc's
reality-check pass into this section per
[`task-plan.md`](../../../spec/planning/task-plan.md)
"Verification protocols are not optional under this carve-out."
Each load-bearing codebase claim below was checked against the
current branch with a one-sentence falsifier; the per-contract
`Verified by:` citations point at the same surfaces.

- **`indexTmpl` is the single page template and a single
  centered column.** Falsifier: "there is more than one
  page template, or the body is already multi-column." Checked:
  [`render.go`](../../../internal/site/render.go) defines exactly
  one `template.Must(... "index" ...)`; `<body>` is one
  `{{range .Roots}}`/`{{else}}` block inside
  `body { … max-width: 60rem; margin: 0 auto }`. False —
  single template, single column. This is the only file t1's
  code change touches.
- **The forest render to preserve is the Roots/empty block +
  `node` template.** Falsifier: "node rendering lives outside
  `indexTmpl` (a partial/file include)." Checked:
  [`render.go`](../../../internal/site/render.go) — the
  `{{if .Roots}} … {{range .Roots}}<div class="root">` block,
  the `{{else}}<p class="empty">No plan-tree roots found …`
  branch, and `{{define "node"}}` (actor markers ranging
  `.WorkInstances`, long-desc, related-prs) are all inline in
  the one template. False — nothing to chase across files; the
  move is within one template body.
- **The data path is untouched by a layout-only change.**
  Falsifier: "the roster needs request-time data, so a
  handler/query/`indexData` change is unavoidable." Checked:
  [`site.go`](../../../internal/site/site.go) `Server.index`
  walks plans + loads work-instances and calls
  `renderIndex(w, indexData{Roots, PlansPath})`;
  [`render.go`](../../../internal/site/render.go) `indexData`
  has only `Roots` + `PlansPath`. The roster placeholder is
  static markup needing no data. False — t1 is template/CSS
  only; `site.go`/`tree.go`/`walker.go` stay untouched.
- **The approved target is the m2 SVG and it specifies the
  scroll/stack behavior.** Falsifier: "the SVG doesn't actually
  say one-scrollbar / roster-top-right / stacked-narrow."
  Checked:
  [`workstreams-view-m2.svg`](../../../design/workstreams-view-m2.svg)
  line 29 states it verbatim; the narrow-window note states
  region stacking, mobile out of scope. False — the shell
  behavior is anchored, not invented.
- **No build/test wrapper exists; the Go toolchain is the
  gate.** Falsifier: "a Makefile/justfile/script wraps
  build+test." Checked: repo root has no `Makefile`/`justfile`;
  `scripts/` is only `assemble.sh` (agent-rules vendoring).
  False — `gofmt`/`go build`/`go vet`/`go test` is the gate.
- **A render-test harness already exists to extend.** Falsifier:
  "there is no render test file / no `renderIndex` test entry
  point." Checked:
  [`render_test.go`](../../../internal/site/render_test.go) has
  a `renderTree` helper calling
  `renderIndex(&buf, indexData{Roots, PlansPath})` and existing
  cases; the new region/placeholder cases extend it. False — the
  harness exists.
- **A second tree exists for the manual render observation.**
  Falsifier: "only the dogfood tree exists, so the gate can't
  observe a second forest." Checked: `docs/plans/demo-workstream/`
  is a populated second plan-tree root. False — the manual
  observation can render both.

## Contracts

Final shapes. Estimate-shaped sections (Files to touch, Commit
Boundaries) are labeled as estimates per
[`shared.md`](../../../spec/planning/shared.md) "Plan content is
a mix of rules and estimates."

These bullets state the observable end-state each surface must
reach. Implementation technique (exact CSS property choice, flex
vs. grid, breakpoint pixel value) is non-binding guidance under
Execution Steps — the contract is the rendered behavior, not the
mechanism.

### Page-shell contract (`internal/site/render.go`)

- The page renders **two regions side by side**: a forest region
  on the left at roughly two-thirds width and a roster region on
  the right at roughly one-third width. The proportion is
  approximate (the design is `~2/3` / `~1/3`, not pixel-exact);
  the contract is "forest is the dominant left region, roster is
  the secondary right region," anchored to
  [`design/workstreams-view-m2.svg`](../../../design/workstreams-view-m2.svg).
  `Verified by:`
  [`indexTmpl` in render.go](../../../internal/site/render.go)
  is today a single centered `body { max-width: 60rem }` column
  with one content block (`{{range .Roots}}`); it is the single
  template t1 restructures into the two regions, and the
  parent-milestone Cross-Task Invariant names it as such.
- **One page scroll.** The page scrolls as a whole (the tall
  forest is the scroll content); neither region gets its own
  independent scrollbar or `overflow` scroll container. The
  roster region is **top-aligned** within the layout so it is
  visible without scrolling when the window is tall enough — it
  is not vertically centered or pushed down by the forest's
  height. `Verified by:`
  [`design/workstreams-view-m2.svg`](../../../design/workstreams-view-m2.svg)
  line 29 ("One page, one scrollbar. … Roster sits top-right so
  it shows without scrolling; the tall forest is what you scroll
  through.").
- **Narrow-window degrade is region stacking only.** Below a
  desktop-narrow width threshold the roster region moves to
  **below** the forest region (full-width stacked, no
  truncation, no horizontal scroll). Mobile/phone layout is out
  of scope, so the threshold is a reasonable desktop-narrow
  value, not a phone-tuned breakpoint, and below it the page is
  the two regions stacked, not a redesigned small-screen view.
  The **per-node-row** narrow degrade (right-aligned
  Status/progress colliding with left-aligned label text) is
  **not** in t1 and is not introduced here — t1 ships the
  existing flat forest render, which has no right-aligned
  per-row element to collide. `Verified by:` the parent
  milestone Task Contract for t1 ("The narrow-window
  *per-node-row* degrade is **not** here … it is t2's") and
  [`render.go`](../../../internal/site/render.go) — the existing
  `node` template renders the badge/label/markers left-to-right
  with no right-aligned column.
- The single visual vocabulary (card/box/spacing/type language)
  is anchored to
  [`design/workstreams-view-m2.svg`](../../../design/workstreams-view-m2.svg);
  the shell does not invent a divergent panel/card language, and
  the roster placeholder uses the same panel framing the design
  shows for the roster column. `Verified by:` the
  parent-milestone Cross-Task Invariant "The shell is t1's …
  Both regions share one visual vocabulary anchored to
  workstreams-view-m2.svg."

### Forest-region contract (`internal/site/render.go`)

- The forest region carries the **existing forest render
  verbatim**: the `{{if .Roots}} … {{range .Roots}}<div
  class="root">{{template "node" .}}</div> … {{else}}<p
  class="empty">No plan-tree roots found …</p>{{end}}` block and
  the entire `{{define "node"}}` template are moved **inside**
  the forest region wrapper with **no change to node shape,
  nesting, badges, actor markers, long-description, or
  related-PR rendering**. The `node` template body is not edited
  at all. `Verified by:`
  [`indexTmpl` in render.go](../../../internal/site/render.go)
  lines 47–77 — the Roots/empty block and `node` template that
  must move unchanged; the milestone names "no node-shape change
  — that is t2."
- The empty-state path is preserved and renders **inside the
  forest region**: when there are no roots, the forest region
  shows the existing `No plan-tree roots found at <code>…</code>`
  message (not a blank region), while the roster region still
  shows its placeholder. `Verified by:`
  [`render.go`](../../../internal/site/render.go) `{{else}}`
  branch (line 54) — the empty-state markup that must survive
  the move.
- All existing CSS classes the forest render depends on
  (`.root`, `.badge`, `.status-*`, `.actor-marker`, `.label`,
  `.empty`, `.long-desc`, `.related-prs`, `ul`/`li`) keep their
  current visual behavior; the layout change adds the
  region/column rules around them and may widen the page
  container, but does not restyle a node's internals. `Verified
  by:` [`render.go`](../../../internal/site/render.go) lines
  24–43 — the existing `<style>` block; the diff adds layout
  rules and does not alter these node-level rules.
- v0.1's per-node actor tags still render on every node exactly
  as today (the deferred "actor icons on progress boxes" work
  assumes node-level actor tags remain). `Verified by:`
  [`indexTmpl` "node" template in render.go](../../../internal/site/render.go)
  line 61 ranges `.WorkInstances` into `actor-marker` spans;
  t1 does not touch that line.

### Roster-region contract (`internal/site/render.go`)

- The roster region renders a **deliberate, intentional
  placeholder**: a panel using the design's roster-column
  framing with a heading (the roster's name, e.g. "Sessions")
  and an explanatory line, styled with the existing muted
  empty-state vocabulary, that reads as a designed "this lands
  in a later task" state — **not** an empty `<div>`, a broken
  gap, or a zero-height region. The placeholder is static
  template markup; it consumes no new data. `Verified by:` the
  parent-milestone Cross-Task Risk "The skeleton placeholder
  ships as a blank/broken gap … t1's roster region must render a
  deliberate, observed placeholder" and the "Bans on surface
  require rendering the consequence" rule in
  [`task-plan.md`](../../../spec/planning/task-plan.md) — the
  Validation Gate observes the rendered placeholder, it is not
  asserted from template source.
- The placeholder does **not** list, query, or hint at session
  data, and adds no promote/dismiss or any roster affordance —
  it states only that the roster surface arrives later. The
  roster (bound + unbound listing) and any session data are
  entirely t4's. `Verified by:` parent-milestone Cross-Task Risk
  "Reviewer-flag any t4 drafting that adds a promote/dismiss
  affordance" — t1 must not pre-empt t4 by rendering a
  data-shaped placeholder.

### Data-path contract (no change)

- `indexData`, `renderIndex`, `Server.index`,
  `loadActiveWorkInstances`, `walkPlans`, `buildTree`, the
  `walker.go` parser, and `tree.go` are **not touched**. The
  roster placeholder is static markup in the template; it
  introduces no new template field, no new query, and no new
  request-time work. `Verified by:`
  [`Server.index` in site.go](../../../internal/site/site.go)
  walks plans + loads work-instances per request and calls
  `renderIndex(w, indexData{Roots: roots, PlansPath: …})`; t1
  changes only the template body `renderIndex` executes, not its
  inputs or the handler.

## Cross-Cutting Invariants

t1 introduces no new cross-cutting rule; this section records
the **inherited** parent-milestone invariants whose enforcement
surface this task's diff brushes against, so self-review and
reviewers know they are reviewer-flag candidates here (per the
milestone's "Reviewer-flag candidates when any per-task drafting
brushes against these"). They are cited, not re-derived, per
[`shared.md`](../../../spec/planning/shared.md)
scoping-vs-duplication discipline.

- **The shell is t1's; siblings build into regions, never the
  shell.** t1 establishes the two named regions and the
  page-shell behavior; no later task alters the shell or a
  region boundary. (Parent milestone "Cross-Task Invariants" —
  first bullet.) t1's diff IS the shell, so every shell decision
  here is the locked surface t2/t4 build against.
- **Render path stays walk-on-every-request.** No caching,
  file-watch, or in-memory build-up. t1 trivially preserves this
  because it changes no data path. `Verified by:`
  [`Server.index` in site.go](../../../internal/site/site.go)
  calls `walkPlans` + `loadActiveWorkInstances` per request; t1
  does not touch `site.go`.
- **v0.1 actor tags on nodes must not regress.** Preserved by
  the Forest-region contract (the `node` template is moved
  unedited). (Parent milestone "Cross-Task Invariants" — actor
  bullet.)
- **Single visual vocabulary anchored to the m2 mockup.**
  Enforced by the Page-shell contract's last bullet.

## Files to touch

*Estimate of expected shape — implementation may revise if a
structural call requires it; deviations are reported per
[`task-plan.md`](../../../spec/planning/task-plan.md)
"Plan-to-PR Completion Gate."*

**Modify:**

- `internal/site/render.go` — restructure `indexTmpl`'s
  `<body>` into the two-region layout (forest region wrapping
  the existing Roots/empty block + `node` template unchanged;
  roster region with the deliberate placeholder); add the
  layout/region CSS rules and widen the page container to fit
  two columns. The `node` template body and the existing
  node-level CSS rules are not edited.
- `internal/site/render_test.go` — add coverage that a rendered
  page contains both region containers, that the forest region
  still contains the existing forest output (a known node's
  badge/label) and the empty-state when there are no roots, and
  that the roster region contains the deliberate placeholder
  text (so the no-roster consequence is observed, not assumed).
- `design/v0.1-design.md` — §7 "What the Website Renders" is
  updated to describe the two-region shell (forest region +
  deliberate roster placeholder) per the parent-milestone
  Documentation Currency entry; this lands in the implementing
  PR.
- `docs/plans/workstream-tracker-1-0/m2-expanded-view-and-roster.md`
  — the parent milestone's Task Status t1 row and the
  "each task is a … stub" prose were updated by *this drafting
  change* (row no longer `In draft (stub)`; see Documentation
  currency). The implementing PR re-touches this file only to
  advance the t1 row to `Landed`.

**New:**

- None beyond this plan doc.

**Intentionally not touched** *(estimate — where we don't expect
changes, not a hard prohibition)*:

- `internal/site/site.go`, `internal/site/tree.go`,
  `internal/site/walker.go`, `internal/site/ghprs.go` — no
  data-path, parser, or handler change; t1 is template/CSS only.
- The `{{define "node"}}` template body and node-level CSS
  (`.root`, `.badge`, `.status-*`, `.actor-marker`, `.label`,
  `.long-desc`, `.related-prs`) — moved/wrapped, not restyled
  (no node-shape change — that is t2).
- `internal/api/*`, `internal/db/*` — no API/DB/schema surface.
- Any caching/file-watch layer — none exists and none is
  introduced (walk-on-every-request invariant).

## Validation Gate

No project build/test wrapper exists. `Verified by:` repo root
has no `Makefile` or `justfile`, and `scripts/` contains only
`assemble.sh` (an agent-rules vendoring tool, unrelated to
build/test). The canonical Go toolchain is the gate:

- `gofmt -l internal` reports no files (formatting clean).
- `go build ./...` succeeds.
- `go vet ./...` clean.
- `go test ./...` passes, including the new render cases:
  - both region containers are present in the rendered HTML;
  - with roots, the forest region still contains the existing
    forest output (a seeded node's Status badge + label) and the
    `node`-template output is byte-for-intent unchanged;
  - with no roots, the forest region contains the existing
    "No plan-tree roots found" empty-state and the roster region
    still contains its placeholder;
  - the roster region contains the deliberate placeholder text.
- **Manual render observation (Bans on surface require rendering
  the consequence).** Run the site against the existing dogfood
  tree and the `demo-workstream` tree and visually confirm: (a)
  forest left ~2/3 with the current tree rendered exactly as
  before; (b) roster right ~1/3 showing the deliberate
  placeholder as an intentional state, not a blank/broken gap —
  the no-roster consequence is *observed*, not assumed; (c) one
  page scrollbar, roster top-aligned and visible without
  scrolling at a tall window; (d) narrowing the window below the
  threshold stacks the roster below the forest with no
  horizontal scroll or truncation.

## Execution Steps

Ordering and gates for the implementing pass. The "how" the
Contracts deliberately omit (exact CSS, flex vs. grid,
breakpoint value) is non-binding guidance here — deviating from
an Execution Step is an estimate deviation, not a contract
breach.

1. **Baseline validation.** On a clean tree, run the Validation
   Gate commands and confirm they pass *before* editing, so a
   pre-existing failure isn't misattributed.
2. **Branch hygiene.** Implement on a dedicated branch off
   current main; no unrelated changes ride along.
3. **Shell restructure.** In `indexTmpl`, wrap the existing
   Roots/empty block + `node` template in a forest-region
   container and add a roster-region container with the
   deliberate placeholder; add the two-column layout CSS
   (suggested, non-binding: a flex row, forest `flex: 2`, roster
   `flex: 1`, container `align-items: flex-start` for
   top-alignment, page-level scroll only, and a `flex-wrap` /
   width media query that stacks the roster below the forest
   below a desktop-narrow threshold) and widen the page
   container. Do not edit the `node` template body or node-level
   CSS.
4. **Tests.** Add the render-test cases from the Validation
   Gate (region presence, forest output unchanged with and
   without roots, roster placeholder text).
5. **Docs currency.** Update `design/v0.1-design.md` §7 to
   describe the two-region shell; the parent milestone row/prose
   were updated at drafting time (Documentation currency).
6. **Self-review.** Run the Self-Review Audits below against the
   diff.
7. **Final validation.** Re-run the full Validation Gate
   including the manual render observation (observe the roster
   placeholder and the stacked degrade, don't assume them).
8. **PR preparation.** PR body carries the
   `## Estimate Deviations` section (or `N/A`); flip this plan's
   Status and the parent t1 row to `Landed` in the same PR per
   the Plan-to-PR Completion Gate.

## Commit Boundaries

*Estimate of cohesive review chunks — the implementer may
refine.* Single implementing PR (N = 1). Expected commits:
(a) the `render.go` shell restructure + CSS; (b) render tests;
(c) `design/v0.1-design.md` §7 update + this plan's
`Proposed → Landed` flip and the parent t1 row advance. Order
lets each commit build and test green.

## Self-Review Audits

Audits from
[`docs/agents/local/self-review-catalog.md`](../../agents/local/self-review-catalog.md)
that map to this PR's diff surfaces, run at step 6:

- **rename-aware-diff-classification** (render surface) — moving
  the Roots/empty block and `node` template inside a new
  forest-region wrapper is a relocation tooling may show as
  add+delete. Hand-classify the moved lines and confirm the
  forest render is byte-for-intent unchanged (no node-shape
  edit) before any "forest render preserved" claim.
- **validation-honesty** (validation surface) — the render
  tests and the manual observation are only "passed" if run
  end-to-end on the final state, including the roster
  placeholder and the stacked narrow-window degrade *observed*
  in a real render, not asserted from template source.
- **readiness-gate-truthfulness** (Status / parent-doc surface)
  — the `Proposed → Landed` flip and the parent milestone t1
  row advance only after every Goal, Contract, and Validation
  step is satisfied or explicitly deferred *in this plan*; no
  best-effort "looks done" flip.

Other seeded audits have no matching surface: t1 adds no effects
/ listeners (effect-cleanup), no user-initiated mutations
(error-surfacing-user-mutations; read-path only), and no
directory restructure or new top-level dirs / Intent-matching
files (trigger-map-currency; one file's template body is
re-wrapped, the package layout is unchanged).

## Out Of Scope

- The session roster itself — bound/unbound listing, session
  data, expandable JSON, any roster affordance — is **t4**. t1
  ships only the placeholder; the placeholder must not render
  data-shaped content that pre-empts t4.
- Nested per-node Status boxes, collapse/expand, and the
  per-node-row narrow-window degrade are **t2** (forest region
  internals). t1 ships the existing flat forest render
  unchanged.
- Doc-declared progress boxes and the spec affordance are
  **t3**.
- Forest activity-tier sections (across-roots Active / In-flight
  / Landed grouping), actor icons on progress boxes, richer
  work-instance states, the triage *action*, and the intent
  strip — parent-milestone Out of Scope (re-homed to a later
  milestone or deferred past 1.0); not introduced here.
- Mobile/phone layout — out of scope; the narrow-window degrade
  is region stacking at a desktop-narrow threshold, not a
  small-screen redesign.
- The collapsible-mechanism / no-JS-posture decision is **t2's**
  call (parent milestone Cross-Task Decisions); t1 stays no-JS
  by shipping the existing static render and a static
  placeholder, and does not pre-decide t2's posture.

## Risk Register

- **The shell restructure silently regresses the forest
  render.** Moving the Roots/empty block and `node` template
  into the forest-region wrapper could drop the empty-state,
  reorder markup, or alter a node-level CSS rule. Mitigation:
  the Forest-region contract requires the `node` body and
  node-level CSS be untouched; the rename-aware audit
  hand-classifies the moved lines; the Validation Gate renders
  the existing tree *and* the no-roots empty-state and confirms
  both survive inside the forest region.
- **The roster placeholder ships as a blank/broken gap.**
  Mitigation: the Roster-region contract requires a deliberate,
  framed, intentional state and the Validation Gate *observes*
  it rendered (Bans on surface require rendering the
  consequence); a zero-height or empty-div placeholder fails
  the gate.
- **The placeholder pre-empts t4.** A data-shaped placeholder
  (fake session rows, a promote/dismiss affordance) would
  rework t4's surface. Mitigation: the Roster-region contract
  bans data-shaped/affordance content; reviewer-flag per the
  inherited posture/triage invariants.
- **The widened page container or column rules visually
  regress node internals.** Mitigation: layout rules are added
  around the regions; node-level CSS is in "intentionally not
  touched"; the manual render observation compares the forest
  region against the prior render.

## Documentation currency

- `design/v0.1-design.md` — §7 "What the Website Renders" is
  updated to describe the two-region shell (forest region with
  the existing render, deliberate roster placeholder) in the
  implementing PR, per the parent-milestone Documentation
  Currency entry.
- `docs/plans/workstream-tracker-1-0/m2-expanded-view-and-roster.md`
  — the parent milestone's Task Status t1 row is updated to
  mirror this plan's frontmatter Status as part of *this*
  drafting change (no longer `In draft (stub)`; the stub note
  no longer holds now that the plan is drafted), and the
  "Each task is a seeded parent-promotion **stub**" prose is
  reconciled to note t1 is now a drafted plan while t2–t4
  remain stubs. Subsequent row values track this plan's Status
  (`In draft → Proposed → Landed`) and land with the PR that
  performs each flip.
- This plan's `Status` lifecycle: `In draft` while drafted, then
  `Proposed` after the promotion-gate self-review (this drafting
  PR's second commit), then `Proposed → Landed` in the
  implementing PR per the Plan-to-PR Completion Gate.

## Backlog Impact

None. No backlog entry graduates, is deleted, split, or shifts.
The `deterministic-interactive-registration` intersection the
parent milestone names is **t4's** (the roster delivers the
observability mitigation); t1 ships only the placeholder and
touches no backlog entry.

## Related Docs

- [`m2-expanded-view-and-roster.md`](m2-expanded-view-and-roster.md)
  — parent milestone; t1 Task Contract, Cross-Task Invariants,
  Cross-Task Risks, and Documentation Currency this plan honors.
- [`README.md`](README.md) — parent epic.
- [`../../../design/workstreams-view-m2.svg`](../../../design/workstreams-view-m2.svg)
  — the approved two-region target this shell renders.
- [`../../../design/v0.1-design.md`](../../../design/v0.1-design.md)
  — §7 (render scope) is updated by the implementing PR.
- [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md)
  — the rules this plan is structured against (narrow-surface
  carve-out, promotion gate, completion gate).
- [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md)
  — cross-level planning rules (additive posture, `Verified by:`
  discipline, parent-doc child contracts).
