---
slug: post-m2-ux-correction
Status: Proposed
short_description: Corrective UX work for the six product-facing gaps m2 shipped without a per-leaf product-validation gate
---

# post-m2 UX correction

## Status

`Proposed`. The recently-shipped milestone
[`workstream-tracker-1-0-m2`](../workstream-tracker-1-0/m2/README.md)
(Landed) delivered the two-region shell, the expanded nested-box
forest, the doc-declared progress-cell row, and the session
roster — and six product-facing UX gaps shipped alongside it
because m2 had no per-leaf product-validation gate. This task is
the corrective work, framed as bug fixes, all six findings
addressed under one plan. The broader per-leaf validation
discipline that catches future gaps already landed forward-only
in [`milestone.md`](../../../spec/planning/milestone.md) "Product
acceptance and per-leaf validation"; m2 is Landed and not
retrofitted.

This is an **N = 3 task plan** (an orchestrating doc; per-phase
HOW lives in three phase plans seeded at the
`` `In draft` → `Proposed` `` promotion-gate walk per
[`shared.md`](../../../spec/planning/shared.md) "Parent-promotion
stub seeding" — seeded in the same commit that flipped this plan
to `Proposed`: [`p1-cosmetic-defects.md`](p1-cosmetic-defects.md),
[`p2-humanize-forest-actor.md`](p2-humanize-forest-actor.md),
[`p3-contract-revisiting.md`](p3-contract-revisiting.md)). The
phase split, sibling sequencing, and per-phase WHAT contract live
in the `## Phase Contracts` section below.

This is a **standalone graduated task**, not a milestone child.
Its only tracking surfaces are this plan tree, the graduating
[`humanize-forest-actor`](../../backlog.md#humanize-forest-actor)
backlog entry, and the new sibling
[`progress-cell-active-state-and-actor`](../../backlog.md#progress-cell-active-state-and-actor)
entry carrying the F3b carve-out.

### OD-walk outcomes (resolved at this drafting session)

The paired scoping doc
[`scoping/post-m2-ux-correction.md`](scoping/post-m2-ux-correction.md)
decomposed six open decisions (OD1–OD6) against merged code. The
contributor walked them in this drafting session; outcomes are
folded into the durable record below. The scoping doc carries the
full decomposition, rejected shapes, and `Verified by:`
grounding; the resolutions are summarized here so a future reader
needn't trace each lock through scoping.

- **OD1 — dropped (implementation-altitude).** The cell-DOM shape
  satisfying Cross-Cutting Invariant **C-INV-1 (cell-anchor)** is
  the implementer's call at the implementing-PR per
  [`shared.md`](../../../spec/planning/shared.md) "Plans describe
  contracts, not implementation"; the contract is "per-cell DOM
  anchor exists for F3b's future overlay" and all three
  decomposed shapes satisfy it.
- **OD2 — H1 (nested `<details>` for the body disclosure).** The
  F7 body disclosure is a nested `<details>`/`<summary>` inside
  the parent `<details>` per-node box; same idiom the forest
  already uses; no JavaScript; the two `<open>` states are
  independent per the native `<details>` semantics. Folded into
  contract **C2**.
- **OD3 — I2b (render markdown with `<a>` tags stripped).** The
  disclosed body content is rendered via the existing
  `goldmark` dependency, with anchor tags stripped at render
  time so link text reads as plain text (no broken navigation
  from the forest's page-URL context). The `maxLongDescLines = 40`
  truncation cap is preserved as the defense-in-depth tail.
  Folded into contract **C2**.
- **OD4 — dropped (implementation-altitude).** The marker
  technique satisfying **C3** (baseline alignment, no border
  overlap) is the implementer's call.
- **OD5 — K3 (known-facts header + raw-JSON disclosure for every
  entry).** Every roster entry opens to a structured known-facts
  block (slug, actor id, bound / unbound, registered-at, last
  event) above an optional raw-JSON block carrying the
  deliberately-unstructured reported metadata when present. The
  registers stay separated (always-known facts above the
  schema-loose reported blob) so t4's "schema-loose, not
  homogenized" invariant
  ([`t4-session-roster.md` "Cross-Cutting Invariants"](../workstream-tracker-1-0/m2/t4-session-roster.md))
  is preserved. Folded into contract **C4**. If `RosterEntry`
  doesn't already carry registered-at / last-event, the gate-walk
  identifies the additive loader extension and the implementing
  PR captures it via the `## Estimate Deviations` PR-body
  callout.
- **OD6 — L2 (N = 3 phases).** Phase shape and ordering are in
  `## Phase Contracts` below. The `## Validation Gate` runs at
  task-terminal (p3's implementing PR); per-phase gates live in
  the three phase plans seeded at the promotion-gate-walk.

### Supersessions of sibling contracts (recorded at scoping)

Three explicit supersessions of locked sibling contracts.
Landed docs are **not** retro-edited; the supersession lives in
this sub-block and is referenced by the implementing PR body and
by the relevant Contracts clauses below. Established at scoping
time per the paired scoping doc (SD4, SD5); housed here within
Status rather than as a top-level `## Decisions` section to match
the t3 / t4 / stub-children precedent and avoid a [task-plan.md](../../../spec/planning/task-plan.md)
"Required and optional sections" variance.

- **D1 — Supersedes m2 t3 D5 (`progress_stages` row gated by
  field presence, Status-independent).** A doc that declares no
  `progress_stages` no longer renders "exactly the one reserved
  Drafting cell"; it renders the default D / P / I / V row
  colored by node `Status` per the Goal. A doc that *does*
  declare `progress_stages` continues to render the
  field-presence-gated declared row unchanged. The
  field-presence discriminator still gates the **shape** the row
  takes (default vs. declared) — it no longer gates the row's
  *existence*. `Verified by:`
  [`t3-doc-declared-stages.md` Contracts C3/C5](../workstream-tracker-1-0/m2/t3-doc-declared-stages.md)
  (the superseded clauses);
  [m2 README "Cross-Task Decisions" → "Doc-declared-stages
  frontmatter shape"](../workstream-tracker-1-0/m2/README.md)
  (the milestone-level record of the same locked decision);
  [`node-progress` template + `.progress-cell*` styles in
  forest.go](../../../internal/site/forest.go) (the render
  surface this task extends). Realized at p3.
- **D2 — Supersedes the `stub-children-on-parent-promotion`
  stub-render contract ("`slug` + `Status: In draft` ⇒ exactly
  the Drafting box").** A pristine seeded stub renders the
  default D filled + P / I / V dashed placeholders row, not the
  single Drafting cell. The stub render case is still preserved
  (the doc still renders, never errors); only the rendered cell
  shape changes. `Verified by:`
  [`stub-children-on-parent-promotion/README.md` Contracts A2 /
  A6](../stub-children-on-parent-promotion/README.md) (the
  superseded clauses); [`parsePlanDoc` in
  walker.go](../../../internal/site/walker.go) (a stub omits
  `progress_stages`, the absent read returns the zero value, so
  the default-cells row derives from `Status` alone — no parser
  change). Realized at p3.
- **D3 — Supersedes the t4 task plan's "no metadata ⇒ plain
  row" contract.** Every roster entry — including one whose
  session reported no metadata — opens to a K3-shape disclosure
  showing what is known about the session (slug, actor id,
  bound / unbound, registered-at, last event) above an optional
  raw-JSON block for reported metadata (OD5 = K3, resolved at
  the OD walk). `Verified by:` [`t4-session-roster.md`
  "Session identity and naming"](../workstream-tracker-1-0/m2/t4-session-roster.md)
  (the superseded clause); [`roster` template in
  roster.go](../../../internal/site/roster.go) (the
  `{{if .Detail}}` branch that gates the existing
  conditional disclosure — the surface this task removes the
  no-metadata branch from). Realized at p3.

The supersessions apply only to the named clauses; the
unsuperseded portions of the three Landed contracts continue to
hold. Landed docs are not retro-edited — the supersession is
declared in this plan's `## Status` →
`### Supersessions of sibling contracts` sub-block rather than
by amending m2 t3 / m2 t4 / stub-children in place. The
structural route (a new plan corrects deferred / not-caught
findings on Landed work) follows
[`milestone.md`](../../../spec/planning/milestone.md)
"Product-validation findings: fix now or defer," whose
forward-only graduation seam is the closest spec authority; the
"supersession lives in the new plan" half of the move is an
unwritten convention this plan establishes (no spec rule
explicitly authorizes or bans it; a follow-on backlog entry
should capture the spec-gap).

### `In draft` → `Proposed` promotion gate walked

`Proposed` was reached when drafting completed (paired scoping
doc; OD1–OD6 resolved at the OD walk; D1 / D2 / D3 supersessions
established at scoping) and the
[`task-plan.md`](../../../spec/planning/task-plan.md)
`` `In draft` → `Proposed` `` promotion gate was walked before
the flip. History:

- **End-to-end coherence** — plan + scoping re-read in order; no
  contradiction between Contracts C1–C6, Phase Contracts
  (p1 / p2 / p3), Cross-Cutting Invariants C-INV-1 through
  C-INV-5, the OD-walk outcomes, the Supersessions block above,
  the per-phase Files-to-touch estimate, the task-terminal
  Validation Gate, and the per-phase Validation Gates that live
  in the seeded phase skeletons.
- **Decision-completeness on Contracts** — every C-clause locked;
  three render-altitude deferrals (C1 dark-mode CSS shape, C3
  marker technique, C5 exact filled color) are explicitly
  authorized by [`task-plan.md`](../../../spec/planning/task-plan.md)
  "Bans on surface require rendering the consequence" + [`shared.md`](../../../spec/planning/shared.md)
  "Plans describe contracts, not implementation"; the C4
  `RosterEntry` loader-extension uncertainty is named as an
  `## Estimate Deviations` path, not a deferral. Phase Contracts
  WHATs locked per the N≥2-task-plan-promotes-independently
  posture (cross-phase WHAT is decision-complete; HOW is the
  phase plan's own scoping).
- **Universal `Verified by:` walk** — every load-bearing claim
  in Contracts, Cross-Cutting Invariants, the Supersessions
  block, and the OD-walk outcomes block carries a citation;
  symbol/template/section anchors dominate per the anchor-
  preference rule.
- **Reality-check inputs re-confirmed** — every cited line
  range in the scoping doc's Reality-check inputs section
  spot-checked against current branch code: render.go (`:27`
  body rule, `:93-106` truncation), forest.go (`:51`
  actor-marker, `:53` node-progress, `:67-75` node template),
  roster.go (`:25-35` template, `:43` label CSS), walker.go
  (`:12-15` goldmark, `:160-177` markdownBody), api.go
  (`:32-52` RegisterRequest). No drift; all citations stable.
- **Always-on rules** — required sections present (Status,
  Context, Goal, Contracts, Files to touch, Validation Gate);
  Phase Contracts present (required when N ≥ 2); the
  prior-draft `## Decisions` section was refolded into this
  Status block as `### Supersessions` to avoid a
  [task-plan.md](../../../spec/planning/task-plan.md) "Required
  and optional sections" variance; no descent to implementation
  prescription (OD1 / OD4 dropped for being implementation-
  altitude; no fenced code; no executable predicates); no
  soft-commitment language.
- **Phase skeletons seeded** — three phase skeletons created in
  this gate-flip commit per
  [`shared.md`](../../../spec/planning/shared.md)
  "Parent-promotion stub seeding":
  [`p1-cosmetic-defects.md`](p1-cosmetic-defects.md) (slug
  `post-m2-ux-correction-p1`),
  [`p2-humanize-forest-actor.md`](p2-humanize-forest-actor.md)
  (slug `post-m2-ux-correction-p2`),
  [`p3-contract-revisiting.md`](p3-contract-revisiting.md)
  (slug `post-m2-ux-correction-p3`). Each carries `Status: In
  draft`, a level-appropriate `short_description`, and the
  inherited Phase Contracts row.

## Context

Today the workstream-tracker page renders an expanded nested-box
forest of plan-tree docs and a session roster side-by-side. The
m2-shipped result is functional but product-rough: the page is
unreadable under OS dark mode because the shell never declares a
light context, every node box opens its full plan-doc body as
unrendered text by default, the native disclosure triangle
overlaps the box border, the progress-cell row is information-free
across the dogfood tree because no doc has adopted the
field-presence-gated declared-stages affordance, a roster entry
whose session reported no metadata is non-interactive, and the
forest and roster identify the same work-instance differently
(raw `wst-<uuid>` vs. the reported `name`).

This task corrects all six product gaps as a single corrective
piece of work. The corrections themselves are bug-fix-shaped:
shell-level CSS for the dark-mode-unreadable chrome; a per-node
disclosure that hides the long description behind explicit user
gesture; an explicit summary-marker treatment; a default
Status-driven D / P / I / V cell row for nodes whose doc declares
no `progress_stages`; every roster entry made openable; the
reported session name rendered in the forest. Each fix is small
in its own surface; the bundling reflects the user's reframing
(treat the six gaps as one corrective task) and the modest LOC
across the three forest-region / shell-region / roster-region
files.

It is being done now because the m2-shipped result is in active
contributor use and the rough edges are friction at every page
view; deferring the corrections past the next milestone hardens
them into "how the tool looks." The surfaces touched: the
page-shell body rule, the forest region's per-node template, and
the roster region's per-entry template — three files, one
internal package. No API, schema, route, dependency, or
rendering-runtime change; the page stays server-rendered,
walk-on-every-request, no-JavaScript.

## Goal

Open the page against `docs/plans/` per
[`docs/dev.md`](../../dev.md) and observe, in both OS dark mode
and OS light mode:

- The page chrome (`<h1>`, inter-card spacing) is legible in
  both modes (F1).
- Every node box's first view is header-only — label, status
  badge, the default-cells row, actor markers. The long
  description and related-PR list sit behind an explicit
  disclosure independent of the parent tree-collapse (F7).
- The native disclosure triangle on each node box aligns with
  the label baseline and does not overlap the box border (F8).
- Every node box renders a D / P / I / V cell row colored by
  the node's `Status`: `Landed` fills all four cells; `In draft`
  fills D and renders P / I / V as dashed empty placeholders;
  every other `Status` renders the non-D cells in one neutral
  shade. A doc that declares `progress_stages` continues to
  render its declared row unchanged (F3a; the supersession of m2
  t3 D5 and the stub-children stub-render contract is recorded
  in `## Status` → `### Supersessions of sibling contracts` as
  D1 / D2).
- The forest's per-node actor marker shows the reported session
  `name` (slug fallback), matching the roster's identity for the
  same work-instance (F4; graduates the
  `humanize-forest-actor` backlog entry).
- Every roster entry opens to a disclosure showing what is
  known about the session — including an entry whose session
  reported no metadata (F9; supersedes the t4 "no metadata ⇒
  plain row" contract).

All while: m2's region invariants hold
([m2 README "Cross-Task Invariants"](../workstream-tracker-1-0/m2/README.md)
"the shell is t1's" — body-rule edits in `render.go` are the
only shell edit; no shell-layout, region-boundary, or
sibling-region body change); the walk-on-every-request render
path is unchanged; the `wst-<uuid>` actor stays the internal
identity key (no actor-generator change, no schema change); the
per-node `actor-marker` span itself, its placement, and its
attachment to every work-instance remain (the m2 v0.1
actor-tag no-regress invariant is honored — F4 changes only the
rendered text inside the span); the F3b future work (active-cell
actor marker and per-cell PR-state coloring) remains additively
attachable via the cell DOM contract C-INV-1 below.

## Contracts

Final WHAT shape. HOW grounding (specific CSS, specific
disclosure mechanism, specific marker treatment, the
rendered-vs-raw call for the disclosed body) lives in the
implementing PR; the scoping doc decomposes each into shapes
(OD1–OD6) for the gate-walk to resolve.

### C1 — Page chrome legible under OS dark mode (F1)

The page-shell body rule declares a light page context so the
user-agent canvas matches the dark text the shell already
declares. The fix is in the page-shell CSS rule only — not in
the per-region card bodies (each card already sets its own light
background and is unaffected by the shell change). No
shell-layout, region-boundary, or sibling-region body edit. The
m2 file-enforced region invariant
([m2 README "Cross-Task Invariants" → "the shell is
t1's"](../workstream-tracker-1-0/m2/README.md)) holds — the body
rule lives in `render.go`'s shell, the only `render.go` change
this task ships. `Verified by:`
[`render.go` body rule in
`indexTmpl`](../../../internal/site/render.go) (the shell-level
CSS this task edits — no `background-color` and no `color-scheme`
today).

### C2 — Header-only default view; body content behind a separate disclosure (F7)

Every node box's first view shows only headers: the label, the
status badge, the default-cells row (C5), the actor markers (C6).
The long description and the related-PR list sit behind a
**separate disclosure** opened by explicit user gesture, realized
as a **nested `<details>`/`<summary>` inside the parent
`<details>` per-node box** (OD2 = H1) — same idiom the forest
already uses, no JavaScript, and the two `<open>` states are
independent per the native `<details>` semantics. The disclosure
is **independent of** the parent `<details>` tree-collapse — the
parent collapse controls child boxes; the body disclosure controls
the long description + related PRs. The two affordances do not
couple. When the disclosure is open, the body content **is
rendered as markdown** via the existing `goldmark` dependency,
with anchor (`<a>`) tags stripped at render time so link text
reads as plain text (OD3 = I2b) — link navigation would resolve
against the forest page URL rather than the source plan-doc's
path, producing 404s; stripping anchors keeps the rendered body
legible without that footgun. The `maxLongDescLines = 40`
truncation cap is preserved as the defense-in-depth tail
(truncate first, then render). `Verified by:`
[`node-detail` template + `.long-desc` style in
forest.go](../../../internal/site/forest.go) (the open-by-default
body render this task collapses);
[`maxLongDescLines` /
`truncateLongDesc` in
render.go](../../../internal/site/render.go) (the truncation cap
this task preserves); [`markdownBody` in
walker.go](../../../internal/site/walker.go) (the source of the
disclosed content — unchanged); the existing `goldmark` import
([`walker.go` `:12-15`](../../../internal/site/walker.go)) keeps
the dependency surface unchanged.

### C3 — Native disclosure marker aligned with the label baseline (F8)

The native `<details>` triangle is replaced by an explicit
summary-marker treatment that aligns with the label text baseline
and does not overlap the box border. The implementation shape
(`list-style: none` plus an inline triangle inside the header
flex, or `summary::marker` content styling) is OD4 in the
scoping doc and is resolved at the gate-walk. The no-JS native
`<details>`/`<summary>` idiom is preserved; the milestone
no-JavaScript decision ([m2 README "Cross-Task Decisions" →
"Collapsible mechanism"](../workstream-tracker-1-0/m2/README.md))
holds. `Verified by:` [the `summary` + `.box-header` flex shape
in forest.go](../../../internal/site/forest.go).

### C4 — Roster entry opens for every entry (F9)

Every roster entry — bound or unbound, with or without reported
metadata — opens to a disclosure showing what is known about the
session in the **K3 shape** (OD5): a structured **known-facts
header** (slug, actor id, bound / unbound, registered-at, last
event) above an optional **raw-JSON block** carrying the
deliberately-unstructured reported metadata when present. Every
entry's disclosure shows the same known-facts header; entries
with reported metadata also show the raw-JSON block; entries
without reported metadata show "(no reported metadata)" in place
of the block, so the reader sees the same disclosure structure
regardless of whether the session reported metadata. The two
registers stay separated: the **always-known facts** are
structured because they are schematized facts the loader knows
about every session; the **reported metadata** stays
schema-loose because it is the deliberately-unstructured field
t4 locked. t4's "schema-loose, not homogenized" invariant
([`t4-session-roster.md` "Cross-Cutting Invariants"](../workstream-tracker-1-0/m2/t4-session-roster.md))
holds — the known-facts header is **not** a schema on the
reported-metadata side; it is a sibling surface for the
loader-known facts. `Verified by:`
[`roster` template in
roster.go](../../../internal/site/roster.go) (the existing
metadata-conditional disclosure this task generalizes);
[`RosterEntry` in
site.go](../../../internal/site/site.go) (the loader's per-entry
fields available for the known-facts header — the gate-walk
identifies any additive passthrough for registered-at /
last-event, captured via the implementing PR's
`## Estimate Deviations`).

### C5 — Default D / P / I / V cell row, Status-driven, when no `progress_stages` declared (F3a)

A node whose doc declares no `progress_stages` renders a default
four-cell row labeled D / P / I / V, in that order, colored by
the node's `Status`:

- `Landed` ⇒ all four cells **filled** (exact filled color is
  render-altitude; conventionally green to match the existing
  `status-landed` palette).
- `In draft` ⇒ D **filled**; P / I / V rendered as **dashed
  empty placeholders**.
- `In progress`, `Proposed`, `Validating`, `Deferred` (and any
  unrecognized Status, which the existing
  [`statusClass` fallback in
  render.go](../../../internal/site/render.go) maps to
  `status-unknown`) ⇒ D + the three non-D cells rendered in one
  **neutral shade** (no per-stage progression inferred; F3b
  carve-out below). The Drafting cell is not labeled "Drafting"
  in this default row — the label is "D" per the
  [`design/vision.md` §3](../../../design/vision.md) D / P /
  I / V vocabulary the default-cells row adopts. The doc-visible
  `Drafting` token banned by m2 t3 D3 remains banned. `Verified
  by:` [`statusClass` in
  render.go](../../../internal/site/render.go) (every recognized
  Status value + the `unknown` fallback the default-cells
  Status→shape mapping reads); [`design/vision.md` §3
  "Sub-stages within a node" + §7 "Smaller open
  questions"](../../../design/vision.md) (the D / P / I / V
  vocabulary + the placeholder-vs-strict open question, which
  this contract picks placeholder for the default-cells case);
  [m2 README "Cross-Task Decisions" → "Doc-declared-stages
  frontmatter shape" D5](../workstream-tracker-1-0/m2/README.md)
  (the field-presence discriminator preserved in scope under D1
  above).

A doc that **declares** `progress_stages` continues to render
the field-presence-gated declared row exactly as m2 t3 specified
— reserved Drafting cell followed by one cell per declared
entry, in document order. No data-model change; no parser
change; no frontmatter spelling change. The render gate becomes
"declared row if `progress_stages` present, else default row";
both branches render a cell row, and the cell DOM is the same
shape per C-INV-1 below.

### C6 — Reported session `name` rendered in the forest's `actor-marker` (F4)

The forest's per-node `actor-marker` span renders the reported
session `name` (slug fallback per the same name-then-slug rule
t4 locked for the roster
[`t4-session-roster.md` "Session identity and naming"](../workstream-tracker-1-0/m2/t4-session-roster.md)),
not the raw `wst-<uuid>` actor. The `actor-marker` span itself,
its placement inside the `label-group`, and its attachment to
every work-instance via the per-node `range .WorkInstances` are
**preserved** — this contract changes only the rendered text
inside the span. The m2 v0.1 actor-tag no-regress invariant
([m2 README "Cross-Task Invariants" → "v0.1 actor tags on nodes
must not regress"](../workstream-tracker-1-0/m2/README.md)) is
honored; the work-instance row carries the `name` field already
available to the roster (see OD5's verification of which
loader-side fields the forest will need to consume; the
`name`-or-slug fallback rule is t4-locked, so no new identity
contract is introduced). The `wst-<uuid>` actor remains the
internal identity key — no actor-generator change, no schema
change. `Verified by:` [`node-header` template in
forest.go](../../../internal/site/forest.go) (the
`actor-marker` span this task edits the rendered text of);
[`buildTree` + `WorkInstanceView` in
tree.go](../../../internal/site/tree.go) (the per-node
work-instance attachment the forest reads from; the gate-walk
re-verifies which session-`name`-bearing field is already on the
view and surfaces any additive forest-side passthrough);
[`t4-session-roster.md` "Session identity and naming"](../workstream-tracker-1-0/m2/t4-session-roster.md)
(the name-then-slug rule reused).

## Phase Contracts

Per-phase **WHAT** (end result, sibling-interface handoff,
preserved behavior); per-phase **HOW** is scoped just-in-time in
each phase plan against then-merged code per
[`task-plan.md`](../../../spec/planning/task-plan.md) "Cross-PR
coordination." Seeded as skeleton docs in the same commit that
flipped this task plan `In draft → Proposed`, per
[`shared.md`](../../../spec/planning/shared.md) "Parent-doc child
contracts" / "Parent-promotion stub seeding":
[`p1-cosmetic-defects.md`](p1-cosmetic-defects.md),
[`p2-humanize-forest-actor.md`](p2-humanize-forest-actor.md),
[`p3-contract-revisiting.md`](p3-contract-revisiting.md). The
split locked at the OD-walk (OD6 = L2) is recorded here as the
locked WHAT contract; the per-phase HOW (file inventory, function
shapes, validation-gate specifics, risks) is each phase plan's
own scoping when the phase drafts.

Sequence: **p1 → p2 → p3** (phases are sequence-steps toward
the task's one outcome, not parallel, per
[`task-plan.md`](../../../spec/planning/task-plan.md) "The level
picker"). p3 is task-terminal: its implementing PR carries the
full product-acceptance walkthrough (the `## Validation Gate`
below), records approval, and flips both the phase plan's and
this task plan's Status `Validating → Landed` per
[`task-plan.md`](../../../spec/planning/task-plan.md) "Task plan
terminal state when N ≥ 2."

| Slug | Short | Long (end result + preserves) | Feeds siblings |
|------|-------|-------------------------------|----------------|
| `post-m2-ux-correction-p1` | Cosmetic defects (F1 + F8) | The page chrome is legible in OS dark mode — `<h1>` and inter-card chrome contrast the dark canvas (C1); the native disclosure marker on each node box aligns with the label baseline and does not overlap the box border (C3). No plan-doc supersessions. Preserves: m2 file-enforced region invariants (the only `render.go` edit is the body rule; no shell-layout / region-boundary / sibling-region body edit); the no-JS native `<details>`/`<summary>` idiom (m2 t2 Collapsible-mechanism decision); walk-on-every-request render path. | Establishes the dark-mode-legible page chrome and the baseline-aligned marker that p2 and p3 observe their findings against (the demo walkthrough's dark + light passes assume legible chrome). Does not change forest content semantics; does not alter the roster region. |
| `post-m2-ux-correction-p2` | Humanize forest actor (F4) | The forest's per-node `actor-marker` span renders the reported session `name` (slug fallback per the name-then-slug rule t4 locked for the roster), not the raw `wst-<uuid>` actor (C6); graduates the [`humanize-forest-actor`](../../backlog.md#humanize-forest-actor) backlog entry (the entry's `Graduated — post-m2-ux-correction` Status flip lives in this drafting change; the work that closes the entry's intent lands here). No plan-doc supersession. Preserves: m2 v0.1 actor-tag no-regress invariant (the `actor-marker` span itself, its placement inside `label-group`, attachment to every work-instance via per-node `range .WorkInstances`); the `wst-<uuid>` actor stays the internal identity key (no actor-generator change, no schema change); walk-on-every-request. | Closes the m2-shipped forest/roster identity disagreement before p3's contract-revisiting work reshapes the surrounding render. Independent of p1's scope (different file regions). |
| `post-m2-ux-correction-p3` | Contract-revisiting + task-terminal (F3a + F7 + F9) | A node whose doc declares no `progress_stages` renders the default D / P / I / V Status-driven cell row (C5; supersedes m2 t3 D5 under **D1** and the stub-children stub-render contract under **D2**); every node box's first view is header-only, with the body content (long description + related PRs) behind a nested `<details>`/`<summary>` disclosure (OD2 = H1) rendering markdown with `<a>` tags stripped (OD3 = I2b) when expanded, independent of the parent tree-collapse (C2); every roster entry opens to the K3 known-facts + raw-JSON disclosure (C4; supersedes t4's "no metadata ⇒ plain row" under **D3**). **Task-terminal** — carries the full product-acceptance walkthrough across all six findings + the post-m2-ux-correction task plan's `Validating → Landed` flip. Preserves: m2 file-enforced region invariants (forest edits in `forest.go`, roster edits in `roster.go`, no shell edit); the existing field-presence-gated declared-stages render unchanged for any doc that declares `progress_stages`; t4's schema-loose-not-homogenized invariant (the K3 known-facts header is a sibling surface for loader-known facts, not a schema imposed on the reported-metadata blob); the m2 v0.1 actor-tag no-regress invariant (p3 does not touch the `actor-marker` span p2 ships). | Realizes all three plan-doc supersessions in one phase so the supersession review concentrates with the rendered consequence. Carries Cross-Cutting Invariant **C-INV-1 (cell-anchor)** — the default-cells row's per-cell DOM is the F3b future attachment surface. The task plan's terminal close-out (Status flip; the Backlog Impact mutations already executed in this drafting change need no further surgery here). |

Cross-phase coordination is via this task plan (the cross-phase
Decisions D1 / D2 / D3 below, the Cross-Cutting Invariants), not
through pre-locked cross-phase contracts in each phase plan, per
[`task-plan.md`](../../../spec/planning/task-plan.md) "Cross-PR
coordination."

## Cross-Cutting Invariants

Rules ≥ 2 sites must hold simultaneously across this task's
files. Reviewer-flag candidates when any implementing change
brushes these.

- **C-INV-1 (cell-anchor).** The cell DOM is the per-stage
  anchor a future per-cell actor marker or per-cell PR-state
  overlay attaches to additively. The default-cells row (C5)
  and the declared-cells row (m2 t3) share the cell-level DOM
  shape — the rendered cell is the attachment point regardless
  of which row drew it. No render decision in this task closes
  off the attachment surface (e.g., by collapsing the row to a
  single combined indicator, or by encoding the row only on the
  row container with no per-cell element). This invariant
  underwrites the F3b
  [`progress-cell-active-state-and-actor`](../../backlog.md#progress-cell-active-state-and-actor)
  backlog entry's "linear, not cliff" migration claim.
  `Verified by:` [`node-progress` template +
  `.progress-cell*` styles in
  forest.go](../../../internal/site/forest.go) (the cell-as-DOM-element
  precedent the default row extends).
- **C-INV-2 (region-boundary preservation).** The m2
  file-enforced region invariants hold. Shell body-rule edits
  in `render.go` are the only `render.go` change (C1); the
  shell layout, the `.layout` container, the
  `{{template "forest" .}}` / `{{template "roster" .}}`
  composition, and the `indexData` / `renderIndex` plumbing
  are not touched. Forest edits live in `forest.go`; roster
  edits live in `roster.go`. A diff that crosses these file
  boundaries the wrong way is reviewer-flag. `Verified by:`
  [m2 README "Cross-Task Invariants" → "the shell is
  t1's"](../workstream-tracker-1-0/m2/README.md) (the
  file-enforced contract); [the m2 t4 data-path carve-out in
  the same Cross-Task Invariants
  section](../workstream-tracker-1-0/m2/README.md) (the
  precedent — t4 was *explicitly* granted a shared data-path
  edit; this task is *not* and stays inside region bodies + the
  shell body rule).
- **C-INV-3 (actor-tag preserved).** Every work-instance still
  attaches to its node via the per-node `range .WorkInstances`
  in the forest; the `actor-marker` span is still rendered for
  every attached work-instance. C6 changes only the rendered
  text inside the span. A change that removes the
  `actor-marker` span, removes the per-work-instance range,
  removes the work-instance-to-node attachment, or collapses
  multiple work-instances into a single marker is the defect,
  not the fix. `Verified by:` [m2 README "Cross-Task Invariants"
  → "v0.1 actor tags on nodes must not
  regress"](../workstream-tracker-1-0/m2/README.md); [the
  `node-header` template's `range .WorkInstances` in
  forest.go](../../../internal/site/forest.go).
- **C-INV-4 (walk-on-every-request preserved).** No caching,
  file-watch, or in-memory build-up is introduced. Every render
  reads the current `docs/plans/` walk and the current
  work-instance state per HTTP request. `Verified by:`
  [`Server.index` in
  site.go](../../../internal/site/site.go) (the per-request
  walk this task does not alter).
- **C-INV-5 (additive, no spec / API / schema change).** The
  task adds no frontmatter field, no API endpoint or
  request/response field, no schema column, and no dependency.
  The fix is render-altitude across three files in
  `internal/site/`. `Verified by:`
  [`RegisterRequest` in
  api.go](../../../internal/api/api.go) (the registered fields
  this task does not extend); [`events` /
  `work_instances` schema in
  schema.go](../../../internal/db/schema.go) (the schema this
  task does not change);
  [`go.mod`](../../../go.mod) (the dependency set this task
  does not extend).

## Files to touch

*Estimate of the expected file shape per phase, not a binding
rule. Implementation may revise this list when a structural call
requires it — including touching a file listed under
"Intentionally not touched." Any deviation is handled via the
PR-body `## Estimate Deviations` callout with the plan
reconciled to what shipped per
[`task-plan.md`](../../../spec/planning/task-plan.md) "Plan-to-PR
Completion Gate." Per-phase plans (seeded at the promotion-gate
walk) carry the authoritative per-phase file inventory.*

- **p1 (F1 + F8) — Modify:**
  - [`internal/site/render.go`](../../../internal/site/render.go)
    — the body rule for F1 (C1). No layout or composition change.
  - [`internal/site/forest.go`](../../../internal/site/forest.go)
    — F8 explicit summary marker (C3).
- **p2 (F4) — Modify:**
  - [`internal/site/forest.go`](../../../internal/site/forest.go)
    — F4 reported-name render inside `actor-marker` (C6). The
    `range .WorkInstances` and the span attachment stay; only
    the rendered text changes.
  - Possibly [`internal/site/tree.go`](../../../internal/site/tree.go)
    and/or [`internal/site/site.go`](../../../internal/site/site.go)
    if the reported-`name` field needs an additive passthrough
    onto the per-node work-instance view; the p2 phase plan
    decides at drafting and the implementing PR captures any
    delta via `## Estimate Deviations`.
- **p3 (F3a + F7 + F9) — Modify:**
  - [`internal/site/forest.go`](../../../internal/site/forest.go)
    — F3a default-cells render (C5; supersedes the unconditional
    Drafting cell + `range .ProgressStages` shape under D1 / D2);
    F7 nested-`<details>` body disclosure (C2; the per-node body
    template moves the long description + related PRs behind the
    nested disclosure rendered as markdown with anchors stripped).
  - [`internal/site/roster.go`](../../../internal/site/roster.go)
    — F9 every-entry-opens K3 render (C4; supersedes the
    `{{if .Detail}}` conditional disclosure under D3).
  - Possibly [`internal/site/site.go`](../../../internal/site/site.go)
    for additive `RosterEntry` fields (registered-at, last-event)
    if not already attached — the p3 phase plan decides at
    drafting.
- **Modify (tests, per phase):** each phase carries semantic
  test coverage for its contracts — p1 covers F1 / F8 render
  outputs; p2 covers F4 actor-marker text; p3 covers
  default-cells row shape per Status, body-disclosed long
  description (rendered, no anchors), every-entry-opens K3
  disclosure. The semantic-not-byte-exact posture m2 t2 / t3
  established
  ([`t3-doc-declared-stages.md` C7](../workstream-tracker-1-0/m2/t3-doc-declared-stages.md))
  is preserved.
- **Modify (docs):** [`design/vision.md`](../../../design/vision.md)
  §3's "placeholder cells vs. only-what-the-plan-committed" open
  question is *referenced as picked for the default-cells case
  by this task*; the surrounding open-question framing remains
  because the strict reading still applies for the declared-row
  case. Whether the §3 / §7 prose is edited here (a one-line
  cross-reference) or left to the implementing PR's content
  call is a small currency question — the
  [`task-plan.md`](../../../spec/planning/task-plan.md)
  documentation-currency-gate convention covers this. The m2
  README and the t3 / t4 / stub-children plan docs are
  **not** retro-edited; their superseded clauses remain
  accurate as historical records of what shipped.
- **Intentionally not touched** *(estimate — where we don't
  expect changes, not a hard prohibition):*
  - [`internal/site/walker.go`](../../../internal/site/walker.go)
    — no frontmatter field added; `parsedDoc` shape unchanged.
  - [`internal/site/tree.go`](../../../internal/site/tree.go) —
    the per-node work-instance attachment and the
    `WorkInstanceView` field set are estimated unchanged for C6;
    the gate-walk verifies whether the reported-`name` field is
    already attached or needs a small additive passthrough (in
    which case `tree.go` and possibly `site.go`'s loader come in,
    handled via Estimate Deviations).
  - [`internal/site/site.go`](../../../internal/site/site.go) —
    the loader is estimated unchanged for C4 / C6; OD5 K1's
    no-metadata always-known-fields may require an additive
    extension, handled via Estimate Deviations.
  - [`internal/api/`](../../../internal/api/),
    [`internal/db/`](../../../internal/db/),
    [`internal/registerclient/`](../../../internal/registerclient/),
    [`cmd/workstream-tracker/`](../../../cmd/workstream-tracker/)
    — no API, schema, register-client, or CLI change.
  - [`spec/planning/`](../../../spec/planning/) — no spec
    change (additive or otherwise); the supersessions live in
    this plan's `## Status` → `### Supersessions of sibling
    contracts` sub-block, not in the spec.
  - The Landed plan docs this task supersedes —
    [`m2/README.md`](../workstream-tracker-1-0/m2/README.md),
    [`m2/t3-doc-declared-stages.md`](../workstream-tracker-1-0/m2/t3-doc-declared-stages.md),
    [`m2/t4-session-roster.md`](../workstream-tracker-1-0/m2/t4-session-roster.md),
    [`stub-children-on-parent-promotion/README.md`](../stub-children-on-parent-promotion/README.md)
    — are not retro-edited.

## Validation Gate

This is the **task-terminal** Validation Gate — the full
product-acceptance walkthrough that runs at p3's implementing PR
(the task-terminal phase per `## Phase Contracts` above).
Per-phase Validation Gates live in each phase plan (seeded at
the promotion-gate walk per
[`shared.md`](../../../spec/planning/shared.md)
"Parent-promotion stub seeding"): each per-phase gate carries
the technical observation specific to its findings (p1: dark/
light chrome legibility, marker alignment; p2: forest
actor-marker reads the reported name, forest/roster identity
matches; p3: default-cells row per Status, header-only first
view with the nested-`<details>` body disclosure, every roster
entry opens to the K3 known-facts + raw-JSON shape). The
task-terminal walkthrough below subsumes all per-phase
observations into one reviewer-facing demo so the product
approval evaluates the complete corrected product across all
six findings together.

Per
[`task-plan.md`](../../../spec/planning/task-plan.md)
"Product-facing leaf tasks output a demo walkthrough," this
section produces a concrete step-by-step product-reviewer-facing
demo walkthrough. The walkthrough is reproducible without
reading the diff. Per
[`shared.md`](../../../spec/planning/shared.md) "Plan-doc
Status," the leaf is product-facing and the mandatory
`Validating` state binds: p3's implementing PR merges at Status
`Validating`; the follow-up doc-only commit records approval
against this walkthrough revision and flips both p3's phase plan
and this task plan `Validating → Landed` per the
[`task-plan.md`](../../../spec/planning/task-plan.md) Plan-to-PR
Completion Gate Post-merge-validation exception, in a single
commit per "Task plan terminal state when N ≥ 2."

### Reviewer-facing demo walkthrough

**Setup.** From a clean working tree at the implementing PR's
head: start the server against the dogfood plan-tree per
[`docs/dev.md`](../../dev.md) "Local Workflow" —
`go run ./cmd/workstream-tracker` with the defaults
(`PORT=8080`, `DB_PATH=./workstream-tracker.db`; the SQLite file
seeds itself on first run).

**Observable conditions.** The demo exercises four
work-instance states; the reviewer arranges for each to hold at
some point during the walkthrough. **How those states are
produced is the reviewer's choice** — naturally-active sessions
against the dogfood tree, the CLI's `register` / `complete`
subcommands per
[`docs/dev.md`](../../dev.md) "Registering and completing a
session," or direct database fixturing all qualify. Constraints
of any particular seeding path (e.g. the slug-grammar
limitation in
[`internal/slugs/slugs.go`](../../../internal/slugs/slugs.go)
`IsWellFormed` that rejects certain root-substring patterns,
which makes `post-m2-ux-correction-*` slugs unregistrable via
the CLI today) are **not** constraints of this gate; the gate
is about the rendered output given the states, not the path
that produced them. The four states the F4 / F9 acceptance
bullets reference (`<DemoName>` is a reviewer-chosen name):

- **(a) Name-bearing bound session.** Active work-instance
  attached to a plan-tree node; reported metadata carries
  `<DemoName>`.
- **(b) No-name bound session.** Active work-instance attached
  to a plan-tree node; session reported no metadata.
- **(c) Name-bearing unbound session.** Active work-instance
  whose slug is not in the walked plan-tree; reported metadata
  carries `<DemoUnboundName>`.
- **(d) No-name unbound session.** Active work-instance whose
  slug is not in the walked plan-tree; session reported no
  metadata.

The roster reads each as active in the walk-on-every-request
render.

**Open the page in OS dark mode** at `http://localhost:8080/`.
Observe:

- **F1 acceptance.** The page chrome — the `<h1>`
  "workstream-tracker" header, the inter-card spacing — is
  **legible** against the dark canvas (text contrasts the
  background, never invisible). Cards stay legible (each card
  sets an explicit light background). No element is "ghosted"
  against the canvas.
- **F8 acceptance.** Every node box's collapse marker (the
  triangle) sits **aligned with the label text baseline** and
  does **not overlap** the box border. The triangle is visible
  in both open and closed `<details>` states.
- **F7 acceptance.** Every node box's first view shows
  **headers only**: the label, the status badge, the
  default-cells row (F3a below), the actor markers. **No long
  description and no related-PR list is rendered until the
  body disclosure is explicitly opened.** The body disclosure
  is a **nested `<details>`/`<summary>`** inside the node
  box's outer `<details>` (OD2 = H1). Opening the body
  disclosure shows the long description **rendered as
  markdown** (code spans, emphasis, headers, lists, etc., per
  the existing `goldmark` dep) with **link text rendered as
  plain text** (anchor tags stripped per OD3 = I2b; no
  clickable links to chase into 404s). Related PRs render
  below the body. **The parent `<details>` tree-collapse does
  not auto-toggle in response** (the two affordances are
  independent — opening the body does not open child boxes;
  opening child boxes does not open the body).
- **F3a acceptance.** Every node box renders a four-cell row
  labeled **D P I V**. A `Landed` node renders all four cells
  filled; a stub `In draft` node (the
  `post-m2-ux-correction` plan tree contains its own `In draft`
  Status) renders the D cell filled and the P / I / V cells
  as dashed empty placeholders; a node whose Status is
  `In progress`, `Proposed`, `Validating`, `Deferred`, or
  unrecognized renders the four cells in one neutral shade.
  A doc whose frontmatter **declares** `progress_stages`
  (none currently exist in the dogfood tree; the walkthrough
  seeds one temporarily by adding `progress_stages: [P1, P2]`
  to any leaf plan doc and reloading) renders the m2 t3
  declared-stages row — Drafting + per-declared-entry cells —
  exactly as before. Revert the temporary frontmatter
  declaration before continuing.
- **F4 acceptance.** Every node box that carries a registered
  session shows the **reported session `name`** inside its
  `actor-marker` span (or the slug fallback when no `name` was
  reported), never the raw `wst-<uuid>` actor. State (a)
  renders `<DemoName>` in the forest; state (b) renders its
  slug in the forest; **the forest and roster display the same
  identity** for the same session.
- **F9 acceptance.** Every roster entry — across states (a)
  through (d) — is **openable** to the same **K3-shape
  disclosure** (OD5 = K3): a known-facts header (slug, actor
  id, bound / unbound, registered-at, last event) above a
  raw-JSON block. Entries with reported metadata ((a), (c))
  show the deliberately-unstructured reported metadata t4
  shipped in the raw-JSON block; entries without reported
  metadata ((b), (d)) show "(no reported metadata)" in place
  of the block. The same disclosure structure across every
  entry is the K3 visual goal.

**Open the page in OS light mode** at the same URL. Observe each
acceptance above; F1 specifically must demonstrate **no
regression** in light mode (the shell's light-context declaration
does not break the page chrome in OS light mode).

**Tear down.** Restore the pre-demo state via whatever path
produced the observable conditions; the CLI's `complete --id
<id>` / `abandon` subcommands per
[`docs/dev.md`](../../dev.md) "Registering and completing a
session" are one such path. The roster returns to its prior
state on the next page reload.

### Approval recording

Per
[`milestone.md`](../../../spec/planning/milestone.md) "Product
acceptance and per-leaf validation" and
[`shared.md`](../../../spec/planning/shared.md) "Plan-doc
Status," the post-merge doc-only commit that flips `Validating
→ Landed` records: who approved, the walkthrough revision
approved (a commit SHA of *this* Validation Gate section), and
the date of approval. The PR-body `## Estimate Deviations`
section names any deviation from this plan's estimates per the
[`task-plan.md`](../../../spec/planning/task-plan.md) Plan-to-PR
Completion Gate.

### Toolchain gate

`gofmt -l internal cmd` reports no files; `go build ./...`,
`go vet ./...`, `go test ./...` all pass per
[`docs/dev.md`](../../dev.md) "Local Workflow." Test coverage
includes: the default-cells row shape per Status; the body
disclosure renders no long description and no related-PR list in
the closed state; an openable roster entry for every metadata
permutation; the forest's `actor-marker` renders the reported
`name` (slug fallback).

## Self-Review Audits

From
[`docs/agents/local/self-review-catalog.md`](../../agents/local/self-review-catalog.md);
diff surfaces are frontend render + frontend tests:

- **validation-honesty** — every F-finding acceptance is
  observed against a real `go run` rendering both OS modes, not
  asserted from the diff. The lens on the manual walkthrough
  above.
- **error-surfacing-user-mutations** — the F9 every-entry-opens
  contract must surface the no-metadata case observably (the
  disclosure is genuinely populated with the always-known fields,
  not an empty `<details>` body that reads as broken).
- **rename-aware-diff-classification** — the F7 body-template
  reshuffle and the F3a default-cells render relocate render
  responsibilities inside the forest template; hand-classify so
  "header-only by default" and "F3b cell-anchor preserved"
  claims are real, not tooling-fooled.

## Out Of Scope

- **F2 (raw markdown body rendering) as an independent fix.**
  Subsumed by F7. The "raw vs. rendered when expanded"
  question is OD3 in the scoping doc, resolved at the
  gate-walk as a sub-decision of F7.
- **F3b (active-cell actor marker; per-cell PR-state coloring).**
  Captured as the
  [`progress-cell-active-state-and-actor`](../../backlog.md#progress-cell-active-state-and-actor)
  backlog entry in the implementing PR. The cell DOM
  cross-cutting invariant C-INV-1 keeps the future deferral
  cheap (linear, not cliff). The data the future work needs
  — sub-stage attribution on the work-instance and a PR-state
  source — is not present in the current model (verified in
  scoping); this task does not introduce either.
- **F5 (narrow-roster word-break).** Accepted as the shipped
  tradeoff. The locked
  [`roster.go`](../../../internal/site/roster.go) `.roster-label`
  uses `word-break: break-all` deliberately to prevent long-slug
  overflow in the ~1/3 column; mid-word wrapping is the accepted
  cost. No code change.
- **OS dark mode adoption beyond F1.** F1 makes the page chrome
  legible in OS dark mode by declaring a light context; a full
  paired dark-mode palette across the shell + every card is a
  larger redesign and explicitly not in scope.
- **`wst-<uuid>` actor generator change.** F4 changes the
  rendered text inside the `actor-marker` span; the underlying
  `wst-<uuid>` actor stays the internal identity key
  (idempotency, ownership). No actor-generator change.
- **Spec / API / schema / dependency change.** None expected;
  the contracts above are all render-altitude.
- **Retro-editing the m2 / t3 / t4 / stub-children Landed plan
  docs.** The supersession lives in this plan's `## Status` →
  `### Supersessions of sibling contracts` sub-block, not in the
  superseded Landed docs themselves; Landed docs are not
  retro-edited per the
  [`milestone.md`](../../../spec/planning/milestone.md)
  "already-`Landed` milestone docs are immutable history ...
  and are not retrofitted" framing (which the spec scopes to
  milestone docs, with the principle generalized here to
  Landed task plans).

## Risk Register

- **The default-cells row makes a node look "more done" than
  it is.** A `Landed` node fills all four D / P / I / V cells
  (green), which a reader could misread as "every stage
  succeeded" rather than "the node reached terminal state."
  Mitigation: the cells in the default row are
  Status-driven *by design*, and the row's job is at-a-glance
  Status signaling, not per-stage truth; the
  declared-stages affordance remains the path for a doc that
  wants per-stage truth, and the cell labels D / P / I / V are
  the
  [`design/vision.md` §3](../../../design/vision.md)
  vocabulary the contributor population reads as "stage
  positions," not "stage outcomes." The Validation Gate
  walkthrough observes both the `Landed` and the `In draft`
  cases in dark and light mode so the reader's interpretation
  matches the rendered output. Carried, not blocking.
- **The body disclosure's open state survives child-collapse
  toggles unexpectedly.** F7's disclosure is independent of
  the parent `<details>` tree-collapse — a reader could
  collapse the parent and find the body disclosure still
  reads as "open" on the next reload (a native `<details>`
  retains its `open` attribute on the current render but not
  across page reloads — `<details>` is not auto-persisted).
  Mitigation: the disclosure is a per-render-pass affordance,
  not user-state-bearing; the Validation Gate explicitly
  observes "the two affordances are independent" by
  toggling each separately. Carried, not blocking.
- **The F3b cell-anchor invariant is silently weakened by an
  implementation that merges the default and declared rows
  into a single row with a per-shape branch on the row
  container.** Mitigation: C-INV-1 names the per-cell
  attachment surface as the load-bearing contract;
  reviewer-flag any render that encodes shape only on the
  row container.
- **The body content's rendered-vs-raw choice (OD3) lands on
  rendered and a markdown body in a dogfood doc carries
  inline markup that fails the contextual-escape
  boundary.** Mitigation: the WHAT contract C2 is
  "human-readable when expanded"; the gate-walk picks rendered
  only if the path is small and safe (a `goldmark.Convert`
  into a trusted-source `template.HTML` value).
  [`walker.go`](../../../internal/site/walker.go) goldmark
  import is already present, so the dependency surface does
  not grow.

## Backlog Impact

Per
[`spec/backlog.md`](../../../spec/backlog.md) effect taxonomy.
Three [`docs/backlog.md`](../../backlog.md) mutations land in
**this drafting change** (the same PR that creates this plan
doc + the paired scoping doc) per the
`stub-children-on-parent-promotion` precedent
([`../stub-children-on-parent-promotion/README.md` Backlog
Impact](../stub-children-on-parent-promotion/README.md)), not
deferred to a phase implementing PR — the graduation flip and
the new-sibling captures are scope-tracking surfaces of the
planning act itself, not implementation outcomes. The work that
*closes* the `humanize-forest-actor` entry's intent lands at
**p2** (F4); the Status flip on the entry happens earlier (at
this drafting change), the entry's intent is *realized* at p2.

- **graduate** the existing
  [`humanize-forest-actor`](../../backlog.md#humanize-forest-actor)
  entry. Its Status flips from `Open` to `Graduated —
  post-m2-ux-correction`, with the `**Plan:**` line pointing at
  [`docs/plans/post-m2-ux-correction/README.md`](README.md) per
  [`spec/backlog.md` "Entry
  lifecycle"](../../../spec/backlog.md). The entry's body prose
  is preserved. The realizing phase is **p2**.
- **add** a new entry for *this* task:
  [`post-m2-ux-correction`](../../backlog.md#post-m2-ux-correction),
  Status `Graduated — post-m2-ux-correction`, with the
  `**Plan:**` line pointing at this doc. The entry is created
  pre-Graduated rather than as an Open-then-graduated pair
  because this task did not pre-exist as a backlog item — it was
  framed and graduated in the same step. Closes when p3 flips
  `Validating → Landed` (task-terminal).
- **add** a new sibling entry capturing the F3b carve-out:
  [`progress-cell-active-state-and-actor`](../../backlog.md#progress-cell-active-state-and-actor),
  Status `Open`. The entry's body carries the deferral analysis
  (re-trigger axis, linear-not-cliff migration cost, named
  mitigation = this plan's C-INV-1, realized at **p3**).

## Related Docs

- [`scoping/post-m2-ux-correction.md`](scoping/post-m2-ux-correction.md)
  — the paired scoping doc (decomposed HOW with rejected
  alternatives; open decisions OD1–OD6 for the gate-walk;
  reality-check inputs; F3b deferral analysis).
- [`../workstream-tracker-1-0/m2/README.md`](../workstream-tracker-1-0/m2/README.md)
  — the Landed m2 milestone whose product gaps this task
  corrects.
- [`../workstream-tracker-1-0/m2/t3-doc-declared-stages.md`](../workstream-tracker-1-0/m2/t3-doc-declared-stages.md)
  — the Landed t3 plan whose D5 contract this task supersedes
  under D1. Not retro-edited.
- [`../workstream-tracker-1-0/m2/t4-session-roster.md`](../workstream-tracker-1-0/m2/t4-session-roster.md)
  — the Landed t4 plan whose "no metadata ⇒ plain row" contract
  this task supersedes under D3. Not retro-edited.
- [`../stub-children-on-parent-promotion/README.md`](../stub-children-on-parent-promotion/README.md)
  — the Landed stub-render contract this task supersedes under
  D2. Not retro-edited.
- [`../../backlog.md`](../../backlog.md) — the
  `humanize-forest-actor` entry this task graduates and the
  new `progress-cell-active-state-and-actor` entry capturing
  the F3b carve-out.
- [`../../../design/vision.md`](../../../design/vision.md) §3
  (sub-stages D/P/I/V; eventual cell states *none* / *active* /
  *in-review* / *complete*) and §7 (placeholder-vs-strict open
  question this task picks placeholder for the default-cells
  case).
- [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md),
  [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md),
  [`../../../spec/planning/milestone.md`](../../../spec/planning/milestone.md)
  — the rules this plan is structured against.
- [`../../dev.md`](../../dev.md) — the local workflow and
  session-registration commands the Validation Gate walkthrough
  invokes.
