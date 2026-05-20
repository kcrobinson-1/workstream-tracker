---
slug: post-m2-ux-correction-p1
Status: Landed
short_description: Cosmetic defects — dark-mode chrome legibility (F1) + baseline-aligned collapse marker (F8)
---

# Phase 1 — Cosmetic defects (F1 + F8)

## Status

`Landed`. p1 ships pure-CSS fixes for two m2-shipped cosmetic
defects under the [`post-m2-ux-correction`](README.md) task plan:
dark-mode page chrome (F1) and the misplaced native `<details>`
disclosure marker (F8). p1 ships first in the p1 → p2 → p3
sequence so p2 and p3 observe their findings against legible
chrome and an aligned disclosure marker.

This phase qualifies as **narrow-surface** under
[`task-plan.md`](../../../spec/planning/task-plan.md)
"Narrow-surface plans may skip the scoping doc": single subsystem
(`internal/site`); two files touched; no new public-API contract
(CSS only); no new cross-cutting invariant (C-INV-2 / C-INV-4 /
C-INV-5 are inherited from the parent task plan unchanged); no
novel mechanism (`color-scheme` and `list-style` are established
CSS patterns). The carve-out compresses the *form* — a
**Reality-check inputs** section sits inline below in place of a
separate scoping doc — not the *function* (the falsifier protocol
against load-bearing claims still applies, including the
navigational `Verified by:` walk the parent plan's offline review
named as a gate-discipline gap).

This is an N = 1 phase plan; no further skeleton seeding fires
from p1's promotion.

### `In draft` → `Proposed` promotion gate walked

`Proposed` was reached at this drafting commit; the
[`task-plan.md`](../../../spec/planning/task-plan.md)
`` `In draft` → `Proposed` `` promotion gate was walked before
the flip:

- **End-to-end coherence** — plan re-read in order; no
  contradiction between Context, Goal, Reality-check inputs,
  Contracts C1 / C2, Cross-Cutting Invariants (inherited),
  Files to touch, Validation Gate, Self-Review Audits, Out of
  Scope, Risk Register.
- **Decision-completeness on Contracts** — both C-clauses are
  locked at WHAT altitude; the implementer's HOW (specific CSS
  shape for F1, specific marker technique for F8) is
  render-altitude and explicitly authorized by
  [`task-plan.md`](../../../spec/planning/task-plan.md)
  "Bans on surface require rendering the consequence" + the
  parent task plan's OD-walk record (the option sets are
  already decomposed in the parent's scoping doc and need no
  re-decomposition here).
- **Universal `Verified by:` walk** — every load-bearing claim
  carries a citation; the navigational pass (read the cited
  file/line range and confirm the target exists) was run, not
  just the conceptual check that a citation is present. The two
  primary code citations ([render.go:27](../../../internal/site/render.go)
  body rule lacking `background-color` / `color-scheme`;
  [forest.go:83-85](../../../internal/site/forest.go) summary
  rules lacking `list-style` / marker styling) were re-verified
  against `origin/main` at this branch's base.
- **Reality-check inputs re-confirmed** — the body-rule and
  summary-rule citations above are stable on current `origin/main`
  (`4a37462`); no drift from the task-plan-time citations.
- **Always-on rules** — required sections present (Status,
  Context, Goal, Contracts, Files to touch, Validation Gate);
  no implementation prescription (no fenced code, no executable
  predicates, no specific CSS spelling locked); no
  soft-commitment language.
- **Phase skeletons** — N = 1; none to seed.

### Implementation history

What shipped at the implementing PR, recorded here so the
durable plan doc describes the picked render-altitude shapes
(both deferrals on Contracts C1 / C2 are now closed).

- **F1 — picked CSS shape: both `color-scheme: light` and an
  explicit `background-color: #fff` on the page-shell body
  rule.** `background-color: #fff` is the load-bearing
  declaration: per the CSS body-background-canvas propagation
  rule, body's background paints the user-agent canvas when
  the `html` element has no background, which makes the
  canvas white in OS dark mode regardless of UA color-scheme
  policy. `color-scheme: light` is the belt-and-suspenders
  partner: it tells the UA the page is designed for the light
  scheme, so any UA-rendered widgets (scrollbars, form
  controls) follow suit. The "color-scheme only" shape
  (option A in the parent's scoping doc) was rejected at
  the implementing PR because `color-scheme: light` applied
  to `body` does not by itself paint the html canvas light —
  the canvas inherits from `:root`. The "both" shape closes
  this falsifier directly and is what shipped.
- **F8 — picked marker technique: J1 (inline triangle
  inside the `.box-header` flex row + `list-style: none` on
  summary + `::-webkit-details-marker { display: none }`
  fallback).** J1 was picked over J2 (`summary::marker`
  content styling) because `summary::marker` is positioned
  by the UA and cannot be baseline-aligned with the
  summary's flex children, which is the contract's load-
  bearing visual. The minimal structural element added
  inside `<summary>` is a `<span class="triangle"
  aria-hidden="true">` placed as the first child of
  `.box-header` (not inside `node-header` — preserving the
  C-INV-2 carve-out's narrow scope). The triangle rotates
  90° in the open state via
  `details[open] > summary .triangle { transform:
  rotate(90deg); }`. The pre-existing `justify-content:
  space-between` on `.box-header` was removed in favor of
  `margin-left: auto` on `.status-group` so the three-item
  flex row (triangle, label-group, status-group) lays out
  with the status badge still pinned to the right edge.

## Context

This plan covers two cosmetic defects m2 shipped without a
per-leaf product-validation gate to catch them. The page chrome
(`<h1>`, inter-card spacing) is unreadable under OS dark mode
because the page-shell body rule declares dark text but not a
light page context, so the user-agent canvas paints dark and
text disappears against it. The native `<details>` disclosure
triangle on each forest node box overlaps the box border and
does not align with the label text baseline, because the
summary uses no explicit marker styling. Both are visible at
every page view today; neither has a workaround.

The fixes are render-surface only, in two files: the page-shell
body rule in [render.go](../../../internal/site/render.go) for
F1, and the summary marker styling in
[forest.go](../../../internal/site/forest.go) for F8 — plus, if
the inline-triangle marker technique is picked for F8, a
minimal structural element inside `<summary>` as part of the
marker treatment. No API, schema, route, dependency, broader
template reshape, or render-runtime change. The page stays
server-rendered, walk-on-every-request, no-JavaScript.

It is being done now as the first phase of the
[`post-m2-ux-correction`](README.md) task because the
dark-mode-legible chrome and the baseline-aligned marker are
the foundation p2's actor-marker change and p3's
contract-revisiting render are observed against; running p1
first removes a confounder from p2's and p3's per-phase
walkthroughs.

## Goal

Open the page against `docs/plans/` per
[`docs/dev.md`](../../dev.md) and observe:

- **In OS dark mode.** The page chrome (`<h1>`, inter-card
  spacing) is legible against the dark canvas — text
  contrasts the background, no element is "ghosted." Cards
  remain legible (each card sets an explicit light
  background, unchanged).
- **In OS light mode.** No regression — the page chrome
  contrasts the light canvas the same as before p1.
- **In both modes.** Every node box's collapse marker (the
  disclosure triangle) sits aligned with the label text
  baseline and does not overlap the box border. The marker
  is visible in both open and closed `<details>` states.

All while m2's file-enforced region invariants hold (the only
`render.go` edit is the body rule — no shell-layout,
region-boundary, or sibling-region body change; the `forest.go`
edits are the summary marker styling for F8 and, if the
inline-triangle marker technique is picked, a minimal structural
element inside `<summary>` — no `actor-marker` change, no
`node-header` / `node-detail` / `node-progress` change, no
broader template reshape); the no-JS native
`<details>`/`<summary>` idiom is preserved
([m2 README "Cross-Task Decisions" → "Collapsible
mechanism"](../workstream-tracker-1-0/m2/README.md));
walk-on-every-request render path is unchanged.

## Reality-check inputs

Inline per the narrow-surface carve-out (no separate scoping
doc). Each carries a `Verified by:` citation re-confirmed
against the branch's base `origin/main` (`4a37462`) before this
flip per the gate's navigational pass.

- **F1 surface.** The page-shell body rule declares `color:
  #111827` but no `background-color` and no `color-scheme`.
  Under OS / browser dark mode the user-agent canvas paints
  dark while the declared dark text stays dark; `<h1>` and
  inter-card chrome become invisible. Cards remain legible
  because each card-shaped element sets an explicit light
  background. *Verified by:*
  [`internal/site/render.go`](../../../internal/site/render.go)
  body rule (line 27 — `body { font-family: ...; padding: 1rem;
  max-width: 80rem; margin: 0 auto; color: #111827; }`); the
  per-region card backgrounds at
  [`forest.go`](../../../internal/site/forest.go) `.box*` styles
  (lines 78–81) and
  [`roster.go`](../../../internal/site/roster.go) `.roster-*`
  styles are the reason cards still appear in dark mode.
- **F8 surface.** The forest `<details>` / `<summary>` markup
  uses no explicit marker styling; the user-agent default
  disclosure triangle is rendered. Today the summary rules are
  `summary { cursor: pointer; }` and `summary .box-header {
  display: flex; }` with no `list-style` or `summary::marker`
  rule, so the UA's default marker renders at the summary's
  content edge and overlaps the box border. *Verified by:*
  [`internal/site/forest.go`](../../../internal/site/forest.go)
  summary rules (lines 83–85); the `node` template's
  `<details>` / `<summary>` markup (lines 67–75) — uses native
  HTML, no JavaScript, no marker styling.
- **Card backgrounds set explicit light values today.** The F1
  fix is correctly scoped to the page-shell body rule only; per-
  region card bodies need no edit because they already set
  explicit light backgrounds:
  [`forest.go`](../../../internal/site/forest.go) `.box { ...
  background: #fff; }` and variants (`box-root: #f9fafb`,
  `box-milestone: #fcfcfd`);
  [`roster.go`](../../../internal/site/roster.go) `.roster-panel
  { background: #f8fafc; ... }`. These are unaffected by p1.
- **Native `<details>` is the locked m2 collapse mechanism.**
  The F8 fix is correctly scoped to summary marker styling
  inside the existing native `<details>` / `<summary>` shape;
  no JavaScript or alternative collapse mechanism is in scope.
  *Verified by:* [m2 README "Cross-Task Decisions" →
  "Collapsible mechanism"](../workstream-tracker-1-0/m2/README.md)
  (the locked m2 t2 resolution: native HTML `<details>` /
  `<summary>`, zero JavaScript, no new route).

The option sets the implementing PR will choose between
(F1's CSS shape: `color-scheme: light` only, explicit
`background-color` only, or both; F8's marker technique:
`list-style: none` + inline triangle in the header flex, or
`summary::marker` content styling) are already decomposed in
the parent's scoping doc
(`scoping/post-m2-ux-correction.md`
SD1 / SD3 and the OD-walk record for OD4) and are not
re-decomposed here. p1's contracts stay at WHAT altitude per
[`shared.md`](../../../spec/planning/shared.md) "Plans
describe contracts, not implementation."

## Contracts

Final WHAT shape. Both clauses inherit the parent task plan's
locked WHAT contracts unchanged; p1 carries them here so the
implementing PR's gate has a single-doc reference.

### C1 — Page chrome legible under OS dark mode (F1)

The page-shell body rule declares a light page context so the
user-agent canvas matches the dark text the shell already
declares. The fix is in the page-shell CSS rule only — not in
the per-region card bodies (each card already sets its own
light background; cards stay legible regardless). No
shell-layout, region-boundary, or sibling-region body edit.
The specific CSS shape (`color-scheme: light`, explicit
`background-color`, or both) is render-altitude per
[`task-plan.md`](../../../spec/planning/task-plan.md) "Bans
on surface require rendering the consequence" — settled in the
implementing PR alongside the dark-mode walkthrough. `Verified
by:` parent task plan
[Contract C1](README.md#c1--page-chrome-legible-under-os-dark-mode-f1)
(the WHAT this clause inherits);
[`internal/site/render.go`](../../../internal/site/render.go)
body rule (the surface this clause's implementing PR edits;
line 27 lacks `background-color` and `color-scheme` on
`origin/main` `4a37462`).

### C2 — Native disclosure marker baseline-aligned, no border overlap (F8)

The native `<details>` disclosure triangle on every node box
aligns with the label text baseline and does not overlap the
box border. The marker stays visible in both open and closed
`<details>` states. The no-JS native `<details>` /
`<summary>` idiom is preserved (the m2 t2 Collapsible-
mechanism decision holds); the only change is the marker's
visual treatment. The specific technique (`list-style: none`
+ inline triangle inside the header flex row, or
`summary::marker` content styling) is render-altitude per
[`task-plan.md`](../../../spec/planning/task-plan.md) "Bans
on surface require rendering the consequence" — settled in
the implementing PR alongside the marker-alignment
walkthrough. `Verified by:` parent task plan
[Contract C3](README.md#c3--native-disclosure-marker-aligned-with-the-label-baseline-f8)
(the WHAT this clause inherits);
[`internal/site/forest.go`](../../../internal/site/forest.go)
summary rules (the surface this clause's implementing PR
edits; lines 83–85 carry no `list-style` or
`summary::marker` rule on `origin/main` `4a37462`).

## Cross-Cutting Invariants

p1 binds and preserves the parent task plan's Cross-Cutting
Invariants unchanged; p1 introduces no new invariants. Cited
by name per [`shared.md`](../../../spec/planning/shared.md)
"Cross-Cutting Invariants section" and the layering
discipline that bars duplicating cross-cutting rules across
parent and child docs.

- **[C-INV-2 (region-boundary preservation)](README.md#cross-cutting-invariants).**
  The only `render.go` edit is the body rule (the page-shell
  CSS for F1); no shell-layout, `.layout` container, region-
  composition, or `indexData` / `renderIndex` plumbing change.
  The `forest.go` edit is scoped to the F8 summary marker
  styling. The `node-header` / `node-detail` / `node-progress`
  templates are not changed; the `actor-marker` span text is
  not changed; the `range .WorkInstances` shape is not changed;
  the roster region is not touched. If the inline-triangle
  marker technique (parent OD-walk's J1 shape) is picked, a
  minimal structural element may sit inside `<summary>` as
  part of the marker treatment — that is structural to the
  technique, in scope, and not a forbidden template reshape.
  Broader reshapes of the `node` / `node-header` / `node-detail`
  / `node-progress` templates are forbidden.
- **[C-INV-3 (actor-tag preserved)](README.md#cross-cutting-invariants).**
  p1 does not touch the `actor-marker` span or the
  `range .WorkInstances` shape; the m2 v0.1 actor-tag no-
  regress invariant holds trivially.
- **[C-INV-4 (walk-on-every-request preserved)](README.md#cross-cutting-invariants).**
  p1's edits are render-surface only (CSS, plus at most a
  minimal structural element inside `<summary>` if the inline-
  triangle marker technique is picked) — no caching, file-watch,
  or in-memory build-up is added; every render reads the current
  `docs/plans/` walk and the current work-instance state per
  HTTP request unchanged.
- **[C-INV-5 (additive, no spec / API / schema / dependency
  change)](README.md#cross-cutting-invariants).** p1's edits
  are render-surface only (CSS, plus at most a minimal
  structural element inside `<summary>` if the inline-triangle
  marker technique is picked); no frontmatter field, API
  endpoint, schema column, or dependency is added. `go.mod` is
  not touched.

(C-INV-1 cell-anchor is realized at p3, not p1; p1 is unaffected.)

## Files to touch

*Estimate of the expected file shape, not a binding rule.
Implementation may revise this list when a structural call
requires it; any deviation is handled via the PR-body
`## Estimate Deviations` callout with the plan reconciled to
what shipped per
[`task-plan.md`](../../../spec/planning/task-plan.md)
"Plan-to-PR Completion Gate."*

- **Modify:**
  - [`internal/site/render.go`](../../../internal/site/render.go)
    — the page-shell body rule (line 27) gains a light-context
    declaration. The specific declaration set (`color-scheme:
    light`, `background-color: <light>`, or both) is the
    implementer's call observed against the dark-mode
    walkthrough. No other `render.go` edit.
  - [`internal/site/forest.go`](../../../internal/site/forest.go)
    — the summary marker styling (lines 83–85 region) gains an
    explicit marker treatment. The specific technique
    (`list-style: none` + inline triangle in the header flex,
    or `summary::marker` content) is the implementer's call
    observed against the marker-alignment walkthrough. If the
    inline-triangle technique is picked, the `node` template
    (lines 67–75) may gain a small structural element inside
    the summary; that is structural to the technique, not a
    region-boundary or `node-progress` / `actor-marker` change.
- **Modify (tests):** the existing forest-region semantic
  test gains coverage for the F8 marker render (a `summary`
  carries an explicit marker rule); no F1 test is added — F1's
  acceptance is observed in the manual dark-mode walkthrough
  per "Bans on surface require rendering the consequence,"
  not asserted from a unit test (the rendered CSS string can
  carry the declarations without exercising the visual result;
  validation-honesty applies). The semantic-not-byte-exact
  posture m2 t2 / t3 established
  ([`t3-doc-declared-stages.md` C7](../workstream-tracker-1-0/m2/t3-doc-declared-stages.md))
  is preserved.
- **Intentionally not touched** *(estimate — where we don't
  expect changes, not a hard prohibition):*
  - The `node` / `node-header` / `node-detail` / `node-progress`
    templates in [`forest.go`](../../../internal/site/forest.go)
    (lines 51–65 region) — except for the small structural
    `summary` change F8 may require under the inline-triangle
    technique noted above.
  - [`internal/site/roster.go`](../../../internal/site/roster.go) —
    F9 (every-entry-opens) lives at p3, not p1.
  - The shell composition / `.layout` / `indexData` /
    `renderIndex` plumbing in
    [`render.go`](../../../internal/site/render.go) — only
    the body rule (line 27) is in scope for F1.
  - [`internal/site/walker.go`](../../../internal/site/walker.go),
    [`internal/site/tree.go`](../../../internal/site/tree.go),
    [`internal/site/site.go`](../../../internal/site/site.go) —
    parser, tree builder, loader are CSS-irrelevant.
  - [`internal/api/`](../../../internal/api/),
    [`internal/db/`](../../../internal/db/),
    [`internal/registerclient/`](../../../internal/registerclient/),
    [`cmd/workstream-tracker/`](../../../cmd/workstream-tracker/) —
    no API, schema, register-client, or CLI change.
  - [`design/vision.md`](../../../design/vision.md) and
    [`spec/planning/`](../../../spec/planning/) — no design
    or spec edit (p1 is bug-fix-shaped, not a contract
    extension).
  - The parent task plan
    [`README.md`](README.md) — no edit at p1; the parent's
    Status Task-tracking surface for p1 (no such surface
    exists today since the parent doesn't carry a per-phase
    Status table; the task plan's terminal-state flip
    happens at p3 per "Task plan terminal state when N ≥ 2").

## Validation Gate

p1's per-phase Validation Gate. The full product-acceptance
walkthrough across all six F-findings runs at p3 (task-
terminal); this gate covers F1 and F8 only. Approval at this
gate is observation-only — the post-merge approval-recording
step that closes `Validating → Landed` for a product-facing
leaf binds the task-terminal phase (p3), not earlier phases
per
[`shared.md`](../../../spec/planning/shared.md) "Plan-doc
Status" (the mandatory-`Validating` rule binds the leaf, and
p1 is not the leaf — p3 is). p1 closes on its technical gate
+ this observation.

### Reviewer-facing demo (scoped to F1 + F8)

**Setup.** From a clean working tree at the implementing PR's
head: start the server against the dogfood plan-tree per
[`docs/dev.md`](../../dev.md) "Local Workflow" —
`go run ./cmd/workstream-tracker` with defaults (`PORT=8080`,
`DB_PATH=./workstream-tracker.db`). No session seeding is
needed — p1 does not touch the progress-cell row, the body
disclosure, the roster, or the actor-marker; the existing
forest content is sufficient for F1 and F8 observation.

**Open the page in OS dark mode** at
`http://localhost:8080/`. Observe:

- **F1 acceptance.** The page chrome — the `<h1>`
  "workstream-tracker" header, the inter-card spacing — is
  **legible** against the dark canvas (text contrasts the
  background; no element is "ghosted" or invisible). Cards
  stay legible (each card sets an explicit light background,
  unchanged by p1).
- **F8 acceptance.** Every forest node box's collapse marker
  sits **aligned with the label text baseline** and does
  **not overlap** the box border. The marker is visible in
  both open and closed `<details>` states (toggle at least one
  node to confirm both states render the marker correctly).

**Open the page in OS light mode** at the same URL. Observe:

- **F1 no-regression.** The page chrome is legible in light
  mode the same as before p1 — the light-context declaration
  does not break the existing rendering in OS light mode.
- **F8 acceptance** (same as dark mode — marker baseline-
  aligned, no border overlap).

**Tear down.** Stop the server. No teardown-specific steps
needed (no seeded sessions, no DB state to clean).

### Toolchain gate

`gofmt -l internal cmd` reports no files; `go build ./...`,
`go vet ./...`, `go test ./...` all pass per
[`docs/dev.md`](../../dev.md) "Local Workflow." The
forest-region semantic test gains coverage for the F8 marker
render (the summary rule set carries an explicit marker
declaration); no F1 test is added per the validation-honesty
posture above.

## Self-Review Audits

From
[`docs/agents/local/self-review-catalog.md`](../../agents/local/self-review-catalog.md);
diff surface is CSS in two files + a forest-region semantic
test:

- **validation-honesty** — F1 / F8 acceptance is observed
  against a real `go run` rendering both OS modes, not
  asserted from the diff. The CSS strings in the rendered
  template are not the load-bearing check; the visual result
  is.
- The general self-review checklist
  ([`how-to-use.md`](../../agents/shared/self-review/how-to-use.md))
  layered on top.

No data / CI / runbook audit surfaces apply (render-surface
diff, no schema, pipeline, or operational doc touched).

## Out of Scope

- **F4 (humanize forest actor).** Lives at p2.
- **F3a / F7 / F9 (default cells / body disclosure / every-
  entry-opens roster).** Live at p3.
- **F3b (active-cell actor icon + per-cell PR-state coloring).**
  Captured as the
  [`progress-cell-active-state-and-actor`](../../backlog.md#progress-cell-active-state-and-actor)
  Open backlog entry (parent task plan Out of Scope); not p1.
- **F5 (narrow-roster word-break).** Accepted as the shipped
  tradeoff (parent task plan Out of Scope); no code change.
- **Full OS dark-mode adoption.** F1 makes the page chrome
  legible in OS dark mode by declaring a light context; a
  paired dark-mode palette across the shell + every card is
  the larger redesign p1 is explicitly not.
- **Card background changes.** The cards already set explicit
  light backgrounds; p1's F1 fix does not edit per-region
  card CSS.
- **Template structure beyond the F8 marker treatment.** If
  the inline-triangle marker technique requires a small
  structural element inside `<summary>`, that is in scope as
  part of the marker treatment itself. Any larger template
  reshape (e.g., changing the `node` template's outer
  `<details>` structure) is out of scope and belongs at p3 or
  a separate phase.

## Risk Register

- **The light-context declaration creates a visible change in
  OS light mode.** A `background-color: white` body rule could
  visually paint the canvas the same color it already appears
  in OS light mode (so no observable change), or the
  user-agent's existing light-mode canvas may differ from the
  explicit value chosen. Mitigation: the Validation Gate's
  OS-light-mode pass observes the no-regression result; if a
  light-mode regression appears, the implementing PR adjusts
  the declaration (e.g., picks the `color-scheme: light` only
  shape) and re-observes. The render-altitude deferral on
  C1's CSS shape exists precisely for this falsifier.
- **The inline-triangle marker requires a `<summary>`
  structural element that breaks the existing flex layout.**
  The summary already uses a flex `box-header` layout (line
  84: `summary .box-header { display: flex; }`). Adding a
  triangle element inside that flex row must respect the
  existing label / status-group alignment. Mitigation: the
  Validation Gate's marker-alignment observation catches a
  broken layout; if the inline-triangle technique fails the
  layout, the implementing PR falls back to the
  `summary::marker` content shape (the J2 option from the
  parent's OD-walk decomposition) and re-observes.
- **F8 marker styling crosses the m2 file-enforcement
  boundary.** Negligible: the `forest.go` summary rules are
  in the forest-region body, the file p1's F8 edit is
  scoped to. The C-INV-2 region-boundary preservation
  invariant holds for both files.

## Related Docs

- [`README.md`](README.md) — parent task plan; locked
  Contracts C1 / C3 this phase inherits as C1 / C2; Cross-
  Cutting Invariants C-INV-2 / C-INV-3 / C-INV-4 / C-INV-5
  this phase preserves; Phase Contracts table's p1 row this
  plan realizes; the task-terminal Validation Gate this
  phase's smaller gate is a subset of.
- `scoping/post-m2-ux-correction.md`
  — paired scoping doc for the parent task plan; SD1 (F1)
  and SD3 (F8) record the scoping-time deliberation this
  phase realizes; the OD-walk record for OD4 (F8 marker
  technique, dropped as implementation-altitude) names the
  J1 / J2 decomposition the implementing PR picks from.
- [`p2-humanize-forest-actor.md`](p2-humanize-forest-actor.md),
  [`p3-contract-revisiting.md`](p3-contract-revisiting.md) —
  sibling phase skeletons; ordered after p1.
- [`../workstream-tracker-1-0/m2/README.md`](../workstream-tracker-1-0/m2/README.md)
  "Cross-Task Invariants" / "Cross-Task Decisions →
  Collapsible mechanism" — the m2 region invariants and the
  no-JS native `<details>` decision p1's contracts preserve.
- [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md),
  [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md)
  — the rules this plan is structured against (in particular
  "Narrow-surface plans may skip the scoping doc", "Bans on
  surface require rendering the consequence", "Plans
  describe contracts, not implementation").
- [`../../dev.md`](../../dev.md) — local workflow the
  Validation Gate's reviewer-facing demo invokes.
