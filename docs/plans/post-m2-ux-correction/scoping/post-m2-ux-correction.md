---
slug: post-m2-ux-correction
short_description: Scoping — corrective UX work for m2-shipped gaps
---

# Scoping — post-m2 UX correction

> **Spawned scoping session — decisions left open for the gate-walker.**
> This is the *scoping / investigation* artifact for the standalone
> graduated task `post-m2-ux-correction`, not the plan doc and not a
> plan-drafting session. It runs no promotion gate and opens no PR.
> Each open HOW decision is decomposed into shapes against merged
> code and left **open** for the user to resolve when the
> `` `In draft` → `Proposed` `` promotion gate is walked. The durable
> plan doc is the paired
> [`../README.md`](../README.md). Per
> [`task-plan.md`](../../../../spec/planning/task-plan.md) "Goal: scoping
> doc + plan doc", this scoping doc carries **no** `Status` field.

## Context preamble

The recently-shipped milestone `workstream-tracker-1-0-m2` (Landed)
delivered the two-region shell, the expanded nested-box forest, the
doc-declared progress-cell row, and the session roster. Six
product-facing UX gaps shipped alongside it because m2 had no
per-leaf product-validation gate to catch them (the gate now lives
in [`milestone.md`](../../../../spec/planning/milestone.md) "Product
acceptance and per-leaf validation" — landed forward-only, not
retrofitted into m2). This task is the *corrective UX work* —
framed as bug fixes, top-to-bottom product-facing — and is therefore
itself bound by that gate.

The surfaces touched: the page shell's body-level CSS rule
([`render.go`](../../../../internal/site/render.go)), the forest
region's per-node template
([`forest.go`](../../../../internal/site/forest.go)), and the
roster region's per-entry template
([`roster.go`](../../../../internal/site/roster.go)). No API,
schema, route, dependency, or rendering-runtime change; the page
stays server-rendered, walk-on-every-request, no-JS.

This is a **standalone graduated task**, not a milestone child.
Its only tracking surfaces are this plan tree, the `humanize-forest-actor`
backlog entry it graduates from, and a new sibling backlog entry
for the F3b carve-out.

## Reality-check inputs (the plan must re-verify these at plan-drafting)

Code-grounded findings the plan's contract claims rest on. Each
carries a `Verified by:` citation; the plan-drafting session
re-confirms them against then-merged code per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Reality-check pass before plan-drafting" and refreshes any drifted
line numbers per the
[`shared.md`](../../../../spec/planning/shared.md) anchor-preference
rule.

- **F1 (page chrome unreadable in OS dark mode).** The page-shell
  body CSS sets `color: #111827` but declares no
  `background-color` and no `color-scheme`. Under OS / browser dark
  mode the user-agent canvas paints dark while text stays dark; the
  `<h1>` and inter-card chrome become invisible. Cards remain
  legible only because each card-shaped element sets an explicit
  light background. *Verified by:*
  [`internal/site/render.go`](../../../../internal/site/render.go)
  body rule (`indexTmpl` `:27` — the shell-level CSS this fix
  touches); the per-region card backgrounds in
  [`forest.go`](../../../../internal/site/forest.go) `.box*` /
  [`roster.go`](../../../../internal/site/roster.go) `.roster-*` are
  the reason cards still appear.
- **F7 (hide long description by default; separate disclosure
  independent of tree collapse).** Today every node box renders its
  long description (the entire post-frontmatter markdown body
  trimmed at 40 lines) by default and as **raw / unrendered text**
  — `white-space: pre-wrap`, no markdown rendering. 25 of 36
  dogfood plan-tree docs (the walker reaches 36 per-doc render
  targets — verified by `find docs/plans -name "*.md" -not -path
  "*scoping*" | wc -l`) hit the 40-line truncation marker, so the
  long descriptions are clearly not the boxes' first-glance
  content. *Verified by:*
  [`markdownBody` in
  walker.go](../../../../internal/site/walker.go) (returns the
  trimmed full post-frontmatter body, `:160-177`);
  [`maxLongDescLines` / `truncateLongDesc` in
  render.go](../../../../internal/site/render.go) (`:93` cap +
  `:99-106` function); [`node-detail` and `.long-desc` in
  forest.go](../../../../internal/site/forest.go) (`:55-65` body
  template + `:99` style — the open-by-default render this
  task replaces).
- **F8 (collapse arrow misplaced).** The native `<details>`
  triangle renders at the summary's content edge and visibly
  overlaps the box border; it does not align with the label text
  baseline. The summary uses no explicit marker styling, so the
  user-agent default disclosure triangle is used. *Verified by:*
  [`node` template + `summary` style in
  forest.go](../../../../internal/site/forest.go) (`:67-75`
  `<details>`/`<summary>` markup; `:83-85` `summary { cursor:
  pointer; }` + `summary .box-header { display: flex; }` — no
  `list-style` or marker styling).
- **F9 (every roster entry clickable, including no-metadata
  ones).** The roster template conditionally renders the entry's
  expandable `<details>` only when the resolved metadata is
  non-empty (`{{if .Detail}}...{{else}}<plain row>{{end}}`); an
  entry whose session reported no metadata renders as a plain
  non-interactive row. The data the user wants disclosed for the
  no-metadata case — slug, actor id, bound/unbound, registered-at,
  last event — is already on the roster row's source rows or
  trivially derivable from the roster loader's input set, not new
  data. *Verified by:* [`roster` template in
  roster.go](../../../../internal/site/roster.go) (`:25-35` the
  `{{if .Detail}}` conditional — F9 supersedes this branch);
  [`t4-session-roster.md` "Reported data and the raw-JSON detail
  view"](../../workstream-tracker-1-0/m2/t4-session-roster.md) (the
  t4 plan whose "no metadata ⇒ plain row" clause this task
  supersedes).
- **F3a (status-driven coarse cell model).** No dogfood plan-tree
  doc declares `progress_stages` in frontmatter (`grep -l
  "progress_stages:" docs/plans/**/*.md` returns nothing), so every
  one of the 36 walked plan-tree docs renders an identical greyed
  Drafting pill, with adoption zero. The render path inserts a
  reserved Drafting cell unconditionally and ranges over the
  parsed list for additional cells. *Verified by:*
  [`node-progress` template in
  forest.go](../../../../internal/site/forest.go) (`:53` — the
  unconditional Drafting cell + `range .ProgressStages` per-cell
  render); [`parsedDoc.ProgressStages` in
  walker.go](../../../../internal/site/walker.go) (`:24-32` the
  parsed-list field this task does not touch); [`design/vision.md`
  §3 "Sub-stages within a node" + §7 "Smaller open
  questions"](../../../../design/vision.md) (the four-stage D/P/I/V
  vocabulary this task adopts as the default-cells label set; §7
  surfaces the "placeholder cells vs. only-what-the-plan-committed"
  open question — this task picks placeholder for the default-cells
  case, while the existing declared-stages render keeps the
  only-what's-committed posture).
- **F4 (humanize forest actor).** The forest's per-node
  `actor-marker` span renders the raw `wst-<uuid>` actor verbatim;
  the m2 t4 roster surfaces the reported `name` (slug fallback).
  Forest and roster therefore identify the same work-instance
  differently — the deliberate, surfaced inconsistency the
  `humanize-forest-actor` backlog entry was opened to track.
  *Verified by:* [`node-header` template in
  forest.go](../../../../internal/site/forest.go) (`:51` renders
  raw `{{.Actor}}` in the `actor-marker` span);
  [`buildTree` in tree.go](../../../../internal/site/tree.go)
  attaches the work-instance to the node (the existing data path
  already carries the work-instance row this task surfaces a
  different field from); [`docs/backlog.md` `### humanize-forest-actor`](../../../backlog.md)
  (the Open entry this task graduates).
- **F3b carve-out grounding (out of scope, captured as a new
  backlog entry).** The work-instance register payload carries no
  sub-stage attribution and the render path never fetches PR state;
  the active-cell actor marker and per-cell PR-state coloring
  therefore require data the current model does not carry.
  *Verified by:* [`RegisterRequest` in
  api.go](../../../../internal/api/api.go) (`:32-52` the registered
  fields — `RootSlug` / `ParentPath` / `ExactSlug` / `NodeType` /
  `Actor` / `Metadata`; no `current_stage` field);
  [`node-detail` template in
  forest.go](../../../../internal/site/forest.go) (`:55-65`
  `RelatedPRs` rendered as a flat list, no GitHub fetch);
  [`design/vision.md` §3 / §9](../../../../design/vision.md) (§3
  names *active* / *in-review* / *complete* as the eventual cell
  states + the actor-icon-on-cell expectation; §9 names the
  manual-invocation reliability risk — neither commits to the
  fetch path or the registration extension; both are future work).
- **F2 subsumption grounding.** The earlier "raw-markdown body"
  finding is **not** an independent fix — F7 hides the body by
  default, and how it is rendered when expanded becomes a
  sub-decision *of* F7's expanded view. The `truncateLongDesc`
  function and the `maxLongDescLines = 40` cap exist in
  [`render.go`](../../../../internal/site/render.go) `:93-106`; F7's
  expanded-view rendering decision is recorded in Open Decision
  OD3 below.

**Note on frozen-doc citations.** This task does **not** cite
[`design/v0.1-design.md`](../../../../design/v0.1-design.md) as an
authority. Per the user's standing guidance,
`design/v0.1-design.md` is a frozen v0.1 end-state record, not a
reconciliation target. All load-bearing claims about render
behavior are grounded in merged code under
[`internal/site/`](../../../../internal/site/) or in
[`design/vision.md`](../../../../design/vision.md).

## Settled by scoping (not open)

The six findings F1, F3a, F4, F7, F8, F9 are scope-locked. Each is
recorded here as a settled scoping decision; the plan doc carries
the durable WHAT contract that realizes it (the plan owns the
contract, this doc references the plan's section by name per the
[`task-plan.md`](../../../../spec/planning/task-plan.md) "Scoping
owns / plan owns" split). Rejected alternatives that were
pre-litigated by the user before the scoping session opened are
recorded here so the gate-walker can see the deliberation rather
than re-discover it.

### SD1 — F1: shell-level CSS makes the page declare a light context

The fix touches only the page-shell body rule — not the per-region
card bodies. The shell declares a light page context so the
user-agent canvas matches the dark text the shell already declares.
The exact CSS shape (`color-scheme: light`, an explicit
`background-color`, or both) is render-altitude — settled in the
implementing PR alongside the dark-mode walkthrough per
[`task-plan.md`](../../../../spec/planning/task-plan.md) "Bans on
surface require rendering the consequence." Locked: shell-level
edit only; the m2 region-body file-enforcement invariant
([m2 README "Cross-Task Invariants" → "the shell is
t1's"](../../workstream-tracker-1-0/m2/README.md)) is preserved —
the touched rule is in `render.go`'s shell, not in `forest.go` /
`roster.go`.

**Rejected.** *A1 — per-region card backgrounds carry the dark-mode
contract (each card sets its own light + dark backgrounds):*
multi-site agreement on a rule that has only one site of failure
(the page-shell canvas is the dark surface; cards are already
light); fragmenting the rule across files brushes the additive /
small-diff posture and the supersededer-of-no-decisions tenet. *A2
— honor OS dark mode in the page chrome by giving the shell + every
card a paired dark-mode palette:* the larger redesign this task
explicitly is not (the task is bug-fix-shaped, not a dark-mode
adoption; the cards being legible-in-light-context is the user's
locked acceptable outcome).

### SD2 — F7: collapse the body content behind a separate disclosure

Today every node box renders its long description open by default.
The corrected default first view shows only headers (label, status
badge, progress cells, actor markers). The long description and the
related-PR list move behind a **separate disclosure** opened by
explicit user gesture, **independent of** the parent
`<details>` tree-collapse that controls child boxes (the parent
collapse is a different affordance with a different scope and must
not couple to body disclosure). The disclosure mechanism (a button,
a tooltip, a nested `<details>`) is render-altitude — OD2 below.
Locked: header-only first view; separate disclosure; body
disclosure independent of child collapse; the
`maxLongDescLines = 40` truncation cap is preserved as a
defense-in-depth tail (the forest is not a document viewer; even
when disclosed, very long bodies still truncate with the existing
marker). This finding **subsumes** the earlier "F2 — raw markdown
body" framing: the rendering question becomes "when expanded, raw
vs. rendered markdown," which is OD3 below.

**Rejected.** *B1 — keep the body open by default but tighten the
truncation cap:* doesn't address the user's locked observation that
header-only is the right first view (25 of 36 boxes truncate
today, so even with a tighter cap the box is still body-dominant).
*B2 — couple the body disclosure to the parent `<details>` tree
collapse:* coupling the two affordances loads the parent collapse
with a second meaning the user does not want it to carry; the
disclosures must be independent.

### SD3 — F8: explicit summary marker

The native UA disclosure triangle is replaced by an explicit
summary-marker treatment that aligns with the label text baseline
and does not overlap the box border. The specific CSS shape — a
`list-style: none` on the summary plus an inline triangle inside
the header flex row is the conventional fix — is render-altitude
and is OD4 below. Locked: explicit summary marker; baseline-aligned
with the label; no overlap with the box border.

**Rejected.** *C1 — leave the triangle at the UA default and adjust
the box border to accommodate:* trades a render fix for a layout
contortion and does nothing for the alignment problem with the
label baseline. *C2 — replace `<details>`/`<summary>` with a
JavaScript-driven collapse:* contradicts the milestone-locked
no-JavaScript decision ([m2 README "Cross-Task Decisions" →
"Collapsible mechanism"](../../workstream-tracker-1-0/m2/README.md))
and inflates the diff far beyond a bug fix.

### SD4 — F9: every roster entry opens

Every roster entry — including one whose session reported no
metadata — opens to a disclosure of what is known about that
session: slug, actor id, bound/unbound, registered-at, last event,
etc. This **supersedes** the locked t4 contract clause "an entry
whose session reported nothing renders as a plain row with no
disclosure control"
([`t4-session-roster.md` "Session identity and
naming"](../../workstream-tracker-1-0/m2/t4-session-roster.md)). t4
is a task plan, not a milestone-level cross-task decision; per
[`task-plan.md`](../../../../spec/planning/task-plan.md) Decisions
in this task plan may supersede a sibling task plan's contract
without retro-editing the sibling. The supersession lives in this
plan's Decisions section; t4's Landed doc is **not** retro-edited.

**Rejected.** *D1 — keep the no-metadata row plain and add a
sidecar table of known-facts elsewhere:* multiplies surfaces for
one read affordance and fragments the "every session is
accountable" goal; the every-entry-opens framing keeps the read
affordance singular.

### SD5 — F3a: default D/P/I/V cells, Status-driven, when no `progress_stages` declared

When a doc declares no `progress_stages` in frontmatter (the only
state any dogfood doc is in today), the forest renders a default
**D / P / I / V** cell row colored by the node's `Status`:

- `Landed` ⇒ all four cells filled (e.g. green; exact value
  render-altitude).
- `In draft` ⇒ D filled, P / I / V dashed empty placeholders.
- `In progress` / `Proposed` / `Validating` / `Deferred` ⇒ one
  neutral shade across the non-D cells. The current data model
  carries no per-stage progression for in-progress nodes (see F3b
  carve-out below); no per-stage inference is invented here.

The label set D / P / I / V is the four-stage vocabulary
[`design/vision.md` §3](../../../../design/vision.md) names. This
**supersedes two locked contracts** that the plan's Decisions
section names explicitly:

- D1 — m2 t3 D5 (the `progress_stages` row is gated by **field
  presence**, Status-independent;
  [`t3-doc-declared-stages.md`
  Contracts](../../workstream-tracker-1-0/m2/t3-doc-declared-stages.md)
  C5; m2 README "Cross-Task Decisions" "Doc-declared-stages
  frontmatter shape" entry). The supersession: a doc that declares
  no `progress_stages` is no longer rendered as "exactly the one
  reserved Drafting cell" (the prior stub-case behavior); it now
  renders the default-cells row above. A doc that *does* declare
  `progress_stages` continues to render the field-presence-gated
  declared row unchanged.
- D2 — `stub-children-on-parent-promotion` stub-render contract
  ("`slug` + `Status: In draft` ⇒ exactly the Drafting box";
  [`stub-children-on-parent-promotion/README.md`
  Contracts](../../stub-children-on-parent-promotion/README.md) A2
  / A6). The supersession: a pristine seeded stub renders the
  default D filled + P / I / V dashed placeholders row, not the
  single Drafting cell. The stub render case is still preserved
  (the doc still renders, never errors); only the rendered shape
  changes.

Landed docs are **not** retro-edited; the supersession lives only
in this plan's Decisions. Both prior contracts remain accurate as
historical records of what shipped in their PRs.

**Rejected.** *E1 — backfill `progress_stages` onto every existing
plan-tree doc so the field-presence path covers every node:* the
user has explicitly declined `short_description` backfill on
pre-t3 docs (see standing feedback), and a more invasive
`progress_stages` backfill brushes the same posture; default cells
let the absent-field render carry useful information without doc
churn. *E2 — render no cell row at all when `progress_stages` is
absent:* loses the at-a-glance progress signal across the whole
dogfood tree and effectively makes the t3 affordance opt-in for
*any* visible progress at all. *E3 — extend the data model to
attribute work-instances to a current stage so the active cell can
be highlighted now:* that is the F3b carve-out and is out of scope
(see below).

### SD6 — F4: render the reported name in the forest

The forest's per-node `actor-marker` span renders the reported
session `name` (slug fallback, same name-then-slug rule t4 locked
for the roster), not the raw `wst-<uuid>` actor. Forest and roster
display the same identity for the same work-instance, closing the
m2-shipped inconsistency the `humanize-forest-actor` backlog entry
opened. The `wst-<uuid>` stays the internal identity key (no change
to the actor generator, no schema change). This **graduates** the
Open backlog entry; the entry's Status flips to `Graduated —
post-m2-ux-correction` in the implementing PR with the
`**Plan:**` line per [`spec/backlog.md` "Entry
lifecycle"](../../../../spec/backlog.md).

**Rejected.** *F1 — revisit the `wst-<uuid>` actor generator
itself:* the user has not asked for this and t4's `wst-<uuid>` is
serving as the idempotency key; replacing it is well beyond a UX
correction. *F2 — keep the forest showing `wst-<uuid>` and add a
hover-tooltip with the reported name:* preserves the inconsistency
the backlog entry was opened against; the name-in-the-forest fix
is the goal.

## F3b carve-out — captured as a new backlog entry, deferral-decision argued

The active-cell actor marker (an icon for the currently-active
agent sitting on a *specific* cell) and per-cell PR-state coloring
(in-review vs. complete, distinct from node `Status`) are
**deferred** and captured as a new
[`docs/backlog.md`](../../../backlog.md) entry
(`progress-cell-active-state-and-actor`, Open) in the implementing
PR. The deferral decision is argued per the
[`task-plan.md`](../../../../spec/planning/task-plan.md)
deferral-decision rule (re-trigger axis, migration cost
classification, named mitigation):

- **Re-trigger pain axis.** Until shipped, the forest tells you a
  session is *somewhere in* a node but not which stage, and an
  in-progress node shows no per-stage progression. The pain
  surfaces in real use the moment a single node carries multiple
  in-flight work-instances against different stages — a regime no
  current dogfood node is in, but one the multi-agent vision
  ([`vision.md`](../../../../design/vision.md) §3) commits the tool
  to.
- **Migration cost.** **Linear, not cliff.** This task's SD5
  default-cells row designs the cell DOM as the per-stage anchor
  any future per-cell actor marker or PR-state overlay attaches to
  additively — the cells exist as discrete elements, in declared
  order, and carry no inference about per-stage progression that a
  future per-stage attribution would have to undo. Re-introducing
  per-cell state later is additive on the cell DOM, not a cell-row
  rewrite.
- **Mitigation that keeps the deferral cheap.** This plan's
  Cross-Cutting Invariant **C-INV-1 cell-anchor** (named in the
  durable plan doc) commits the cell DOM as the per-stage anchor a
  future actor marker / state overlay attaches to. The invariant
  makes the deferral cheap by binding *this* task to leave the
  attachment point in place; F3b graduation later does not have to
  excavate one out of a render that closed off the affordance.
- **Required data not yet captured.** F3b also needs new data the
  current model does not carry: a sub-stage attribution on the
  work-instance (today `RegisterRequest` registers against a slug
  only;
  [api.go](../../../../internal/api/api.go) `:32-52`) and a
  PR-state source (today the render path renders `RelatedPRs` as a
  flat list and never fetches PR state;
  [forest.go](../../../../internal/site/forest.go) `:60-64`;
  [`design/vision.md` §3 / §9](../../../../design/vision.md) name
  the GitHub API as the eventual input but the codebase contains
  no proof). One option among several: schematize a small
  `current_stage` field on the work-instance register payload and
  source PR state from the GitHub API.

The backlog entry's body carries this analysis directly; the plan
doc's Out of Scope section references the entry by slug.

## F5 — accepted as-is

The narrow-roster word-break (a long slug mid-word-wrapping inside
the ~1/3 roster column under `word-break: break-all`) is
**accepted as the shipped tradeoff**, not corrected. The locked
roster CSS is in
[`roster.go`](../../../../internal/site/roster.go) `.roster-label`
(`:43`); mid-word wrapping is the deliberate cost of preventing
long-slug overflow in the narrow column. The plan doc's **Out of
Scope** records this with the rationale; no code change.

## Decisions resolved at OD walk

The six open decisions below were decomposed into shapes against
merged code per
[`shared.md`](../../../../spec/planning/shared.md) "Decompose
options into shapes before analyzing"; the contributor resolved
them in an OD-walk conversation that preceded the
`` `In draft` → `Proposed` `` promotion-gate walk. Outcomes are
folded into the durable plan doc's contracts and
`## Phase Contracts` section; the full decomposition with
rejected shapes is retained here as the durable scoping record.
Constraints that bound every decision: stay at contract altitude
in the plan doc (implementation prose lives in the implementing
PR); preserve the m2 file-enforced region invariants (shell
edits in `render.go` are allowed for F1 only — body rule, not
layout; forest edits in `forest.go`; roster edits in
`roster.go`); preserve the walk-on-every-request render path;
spec changes (none expected here) stay additive.

### OD1 — F3a cell-DOM shape: **dropped (implementation-altitude)**

**Resolution.** Dropped from the OD list. The user-facing
contract Cross-Cutting Invariant **C-INV-1 (cell-anchor)**
already names: the cell DOM is the per-stage anchor a future
F3b actor marker / state overlay attaches to additively. All
three decomposed shapes (G1 / G2 / G3) satisfy that contract;
which class shape ships is the implementer's call at the p3
implementing PR per
[`shared.md`](../../../../spec/planning/shared.md) "Plans
describe contracts, not implementation."

### OD1 (original framing, retained as the scoping record) — F3a cell-DOM shape: shared `.progress-cell` class with per-Status modifiers vs. a separate class family

The default-cells row (SD5) needs a DOM contract the F3b carve-out
can later hang a per-cell actor marker and per-cell PR-state overlay
on (Cross-Cutting Invariant C-INV-1 cell-anchor). Two shapes
worth distinguishing:

- **G1 — Shared `.progress-cell` element across default and
  declared rows, with modifier classes per Status-driven state
  (filled / placeholder / neutral).** One DOM contract for both
  row shapes; F3b's future actor-marker attaches to
  `.progress-cell` regardless of which row drew it. Smallest
  attachment surface; preserves the t3 cell-vocabulary (the
  rendered element stays a "progress cell"). The declared-stages
  per-cell label and the default-cells D/P/I/V label coexist as
  the cell's text content, distinguished by modifier class.
- **G2 — Separate `.progress-cell-default` / `.progress-cell`
  classes per row shape.** Two attachment surfaces; F3b's future
  overlay would have to know which row drew which cell. Larger
  CSS surface; cleaner separation but unhelpful for the
  cross-cutting invariant.
- **G3 — A nested `.progress-cell` inside a per-row container
  with the row's shape encoded only on the container.** Best
  separation of concerns, but the deepest DOM; the F3b overlay
  reads the same per-cell anchor regardless.

The decision shapes Cross-Cutting Invariant C-INV-1 in the plan
doc. *Verified by:* [`node-progress` template + `.progress-cell` /
`.progress-cell-drafting` styles in
forest.go](../../../../internal/site/forest.go) (`:53` template +
`:106-107` styles — the existing two-class precedent G1 extends).

### OD2 — F7 body disclosure mechanism: **resolved = H1 (nested `<details>`/`<summary>`)**

**Resolution.** H1: a nested native `<details>`/`<summary>`
inside the parent `<details>` per-node box. Same idiom the
forest already uses; no JavaScript; the two `<open>` states are
independent per the native `<details>` semantics. Folded into
the plan doc's contract **C2**. H2 (`:has()` button) and H3
(tooltip) rejected for the reasons below; H3 specifically loses
keyboard reachability and is not a `<details>`-equivalent
affordance.

### OD2 (original framing, retained as the scoping record) — F7 body disclosure mechanism: nested `<details>` vs. a styled button toggling a sibling region vs. a tooltip

The body disclosure (SD2) must be independent of the parent
`<details>` tree collapse. Three shapes:

- **H1 — Nested `<details>`/`<summary>` for the body, inside the
  parent `<details>` per-node box.** No JavaScript; same idiom the
  forest already uses; the body summary line carries the
  affordance ("Show description / Show details" — exact wording
  render-altitude per "Bans on surface require rendering the
  consequence"). Native `<details>` nesting is supported by every
  target browser and does not couple the two `<open>` states.
  Smallest diff; preserves the no-JS milestone tenet.
- **H2 — Styled button toggling a sibling region via the `:has()`
  selector + CSS-only state.** No JavaScript but adds a CSS state
  contract (`:has()` browser-support is post-Safari-15.4 / Chrome
  105); a more bespoke pattern than the project already uses.
- **H3 — Tooltip on a body-summary affordance.** Trades
  disclosure-on-click for hover-disclosure; loses keyboard-only
  reachability; not a `<details>`-equivalent affordance.

Cross-references OD3 (the rendered-when-expanded question is the
sub-decision *of* whichever shape OD2 lands on). *Verified by:*
[`node` template's existing native `<details>`/`<summary>` in
forest.go](../../../../internal/site/forest.go) (`:67-75` — the
no-JS idiom H1 reuses); [m2 README "Cross-Task Decisions" →
"Collapsible mechanism"](../../workstream-tracker-1-0/m2/README.md)
(the no-JS decision OD2 must stay consistent with).

### OD3 — F7 when expanded: **resolved = I2b (render markdown, strip `<a>` tags)**

**Resolution.** I2b: when the body disclosure opens, the body
content is rendered as markdown via the existing `goldmark`
dependency, with anchor (`<a>`) tags stripped at render time so
link text reads as plain text. Rationale: the bodies are
plan-doc markdown — code spans, inline emphasis, lists, headers
— which read poorly as raw `<pre>`-style text; rendering
delivers a large legibility win at low cost (one
`goldmark.Convert` call plus an HTML-tree pass to strip
anchors). The strip is the difference between this and I2a
(full render): plan-doc bodies use relative-path links like
`[design/vision.md](../../../../design/vision.md)` which the
forest's page URL would resolve to broken paths, so clickable
links would 404. Stripping anchors keeps the rendered body
readable without that footgun. I2c (rewrite links to source
plan-doc paths at render time) is more complex (loader
passthrough of each doc's source path through to render-time)
and not required by the user-facing contract. The
`maxLongDescLines = 40` truncation cap is preserved as the
defense-in-depth tail (truncate first, then render). Folded
into plan doc contract **C2**.

### OD3 (original framing, retained as the scoping record) — F7 when expanded: raw `<pre>`-style markdown vs. rendered markdown

When the body disclosure opens, the long-description content is
either:

- **I1 — Rendered raw / pre-wrapped (the current shape, minus
  open-by-default).** Smallest diff; the existing
  `truncateLongDesc` + `white-space: pre-wrap` rendering path is
  reused unchanged. The disclosed body still reads as
  source-markdown; this is not a regression because the body is no
  longer the first-glance content.
- **I2 — Rendered markdown (goldmark inline / block).** A
  pleasanter reading experience inside the disclosed body but
  requires wiring the existing `goldmark` parser (already a
  dependency for the frontmatter read,
  [`walker.go`](../../../../internal/site/walker.go) `:12-15`)
  through to render-time; the rendered HTML must then pass
  through `html/template`'s contextual-escape boundary safely
  (and a Markdown body in the dogfood tree can carry inline
  markup that loses information when escaped).

The user's posture is "default to whatever keeps the diff smallest
unless rendering is cheap." If the rendered path is essentially a
`goldmark.Convert` call into a `template.HTML` value the body
template trusts (because the markdown originates from the project's
own plan-tree docs the walker already read), the rendered path may
be small enough to prefer. The gate-walk decides; the plan doc's
WHAT contract is "when expanded the body is human-readable" — the
rendered-vs-raw call is HOW. *Verified by:*
[`walker.go` goldmark import](../../../../internal/site/walker.go)
(`:12-15` — goldmark is already a dep);
[`truncateLongDesc` in
render.go](../../../../internal/site/render.go) (`:99-106` — the
existing escape-boundary path).

### OD4 — F8 marker shape: **dropped (implementation-altitude)**

**Resolution.** Dropped from the OD list. The user-facing
contract **C3** already names: the marker aligns with the
label baseline and does not overlap the box border. Both J1
(`list-style: none` + inline triangle) and J2 (`summary::marker`
content) are decomposed shapes of the same contract; which
satisfies it ships is the implementer's call at the p1
implementing PR per
[`shared.md`](../../../../spec/planning/shared.md) "Plans
describe contracts, not implementation."

### OD4 (original framing, retained as the scoping record) — F8 marker shape: inline custom triangle inside the header flex vs. `summary::marker` styling

Two shapes that keep the no-JS native `<details>`/`<summary>`:

- **J1 — `summary { list-style: none; }` + an inline
  triangle element inside the header flex row, toggled
  open/closed via `details[open] summary > .triangle`
  attribute-selector CSS.** Most flexible; aligns naturally with
  the existing flex `box-header`. The user named this in the spec
  as the conventional fix.
- **J2 — `summary::marker` content rule (`summary::marker {
  content: ...; }`).** Lets the user-agent place the marker via
  its native code path; `::marker` content support is universal
  in target browsers, but `::marker` cannot be flex-aligned with
  arbitrary children, so the overlap-with-border problem the user
  named is not fully addressed by `::marker` alone.

The decision is render-altitude; the plan doc's WHAT contract is
"the marker aligns with the label baseline and does not overlap
the box border." *Verified by:* [the `summary` + `.box-header` flex
shape in forest.go](../../../../internal/site/forest.go) (`:83-85`).

### OD5 — F9 every-entry-opens disclosure shape: **resolved = K3 (known-facts header + raw-JSON, same shape every entry)**

**Resolution.** K3: every roster entry — bound or unbound, with
or without reported metadata — opens to the same structured
disclosure: a known-facts header (slug, actor id, bound /
unbound, registered-at, last event) above an optional raw-JSON
block carrying the deliberately-unstructured reported metadata
when present. Entries without reported metadata show "(no
reported metadata)" in place of the raw-JSON block, preserving
the same disclosure structure across every entry. The two
registers stay separated: always-known facts are structured
because they are schematized facts the loader knows about every
session; reported metadata stays schema-loose because it is the
deliberately-unstructured field t4 locked
([`t4-session-roster.md` "Cross-Cutting Invariants" →
"Schema-loose, not homogenized"](../../workstream-tracker-1-0/m2/t4-session-roster.md)).
The known-facts header is **not** a schema on the
reported-metadata side; it is a sibling surface for the
loader-known facts. Folded into plan doc contract **C4**. K1
(JSON for both) rejected because it conflates the schematized
always-known facts with the deliberately-unstructured reported
metadata in a single read register; K2 (table-vs-JSON per
case) rejected because two disclosure shapes give the reader
inconsistent presentation row-to-row. If `RosterEntry` doesn't
already carry registered-at / last-event, the p3 phase-plan
drafting identifies the additive loader extension and the p3
implementing PR captures the delta via the `## Estimate
Deviations` PR-body callout.

### OD5 (original framing, retained as the scoping record) — F9 every-entry-opens disclosure shape for a no-metadata entry

A no-metadata roster entry must open to a disclosure of what is
known about the session — slug, actor id, bound/unbound,
registered-at, last event. Two shapes:

- **K1 — Reuse the existing raw-JSON `<pre>` block (same shape
  the metadata-bearing entry uses today;
  [roster.go](../../../../internal/site/roster.go) `:29`), filling
  the JSON with the always-known fields the loader can produce.**
  Smallest diff; one disclosure shape for both metadata-bearing
  and no-metadata entries; preserves t4's deliberately-unstructured
  raw-JSON posture (the t4 plan's "schema-loose, not homogenized"
  invariant). Some always-known fields (registered-at, last event)
  may need a small loader extension if not already in the roster
  row.
- **K2 — A small known-facts table specific to the no-metadata
  case, distinct from the metadata-bearing entry's raw-JSON
  block.** Two disclosure shapes; cleaner read for the
  no-metadata case but adds a parallel render path.
- **K3 — A hybrid: a small known-facts header followed by the
  raw-JSON block (which is empty / "no reported metadata" for
  the no-metadata case, populated for the metadata-bearing
  case).** One row template, more uniform; slightly more
  complex than K1.

The user's posture is "open to show what's known about the
session — even when no metadata was reported." Whichever shape the
gate-walker picks, the supersession D3 in the plan's Decisions
section is the same. *Verified by:*
[`roster` template's existing `<details>` + `<pre class="roster-detail">` shape in
roster.go](../../../../internal/site/roster.go) (`:25-35` and
`:48`); [`RosterEntry` in
site.go](../../../../internal/site/site.go) (the loader's per-entry
fields — the gate-walker re-verifies which fields are already on
the entry and which need an additive loader extension).

### OD6 — Phase split: **resolved = L2 (N = 3 phases)**

**Resolution.** L2: N = 3 phases, ordered **p1 → p2 → p3**.

- **p1 — `post-m2-ux-correction-p1` (cosmetic defects).** F1
  + F8. Pure CSS in `render.go` body rule + `forest.go` summary
  marker. No plan-doc supersessions. Lowest risk; ships first;
  gives the rest of the demo a legible foundation.
- **p2 — `post-m2-ux-correction-p2` (humanize forest actor).**
  F4 only. Single-line template change in `forest.go`
  `actor-marker`, plus any small `WorkInstanceView` passthrough
  the p2 phase-plan drafting verifies. No plan-doc supersession;
  graduates the
  [`humanize-forest-actor`](../../../backlog.md#humanize-forest-actor)
  backlog entry (the entry's Status flip lives in this
  drafting change; the work that closes the entry's intent
  lands at p2).
- **p3 — `post-m2-ux-correction-p3` (contract-revisiting +
  task-terminal).** F3a + F7 + F9. Carries all three plan-doc
  supersessions (D1 / D2 / D3 in the plan's Decisions section)
  and the full product-acceptance walkthrough. Task-terminal —
  flips both p3's phase plan and the post-m2-ux-correction task
  plan `Validating → Landed` per
  [`task-plan.md`](../../../../spec/planning/task-plan.md) "Task
  plan terminal state when N ≥ 2."

**Rationale for placing the contract-revisiting work last
(p3).** The supersession review concentrates with the rendered
consequence — reviewers evaluate D1 / D2 / D3 against the
rendered result rather than approving them speculatively in a
middle phase that then drifts under p3's later observation.
p3's task-terminal status (heaviest contract surface + heaviest
review gate) is consistent with how t4-p2 carried t4's terminal
walkthrough
([`t4-session-roster.md` "Phase Contracts"](../../workstream-tracker-1-0/m2/t4-session-roster.md)).

L1 (N = 1) rejected because: while the branch-test thresholds
(>5 subsystems or >300 LOC) are satisfied for L1, the
contract-supersession review concentrates badly when bundled
with the cosmetic defects + the actor-name change — three
distinct review-coherence shapes in one PR. L2's split gives
each PR a single review-coherence shape (CSS defects; backlog
graduation; contract-revisiting) without re-introducing the
file-per-future-owner work
([`feedback_file_per_future_owner.md`](../../../../../../../.claude/projects/-Users-kyle-workspace-workstream-tracker/memory/feedback_file_per_future_owner.md)) that the gate-walk's stub-seeding step
covers by seeding one phase-skeleton file per phase.

Per
[`shared.md`](../../../../spec/planning/shared.md) "Parent-
promotion stub seeding," the gate-walk PR that flips this task
plan `In draft → Proposed` seeds three phase-skeleton files
(`post-m2-ux-correction-p1.md`, `-p2.md`, `-p3.md`) — not now;
this scoping doc is paired with an `In draft` plan and seeds
nothing.

### OD6 (original framing, retained as the scoping record) — Phase split: branch test outcome

Per [`task-plan.md`](../../../../spec/planning/task-plan.md)
"PR-count predictions need a branch test," the phase split is not
pre-baked. Two plausible shapes:

- **L1 — N = 1.** All six findings (F1, F3a, F4, F7, F8, F9)
  ship in a single coherent phase. The diff spans
  `render.go` (F1 body rule), `forest.go` (F3a default cells, F4
  actor name, F7 body disclosure, F8 marker), and `roster.go` (F9
  every-entry-opens) — three files, all `internal/site`. Modest
  LOC; review-coherence is plausible because every finding
  participates in the same "make m2 product-acceptable" narrative.
  The plan absorbs phase content inline.
- **L2 — N ≥ 2.** A natural split is by contract-shape: defects
  F1 / F8 (no contract supersessions; pure CSS) as one phase;
  contract-revisiting F3a / F7 / F9 (supersedes m2 t3 D5,
  stub-children D2, t4's "no metadata ⇒ plain row" D3) as a
  second phase; F4 as a third (smallest scope, graduates a
  backlog entry, no supersession). Splits review-coherence by
  whether each PR carries a contract supersession in the plan's
  Decisions section, which may make each smaller PR's review
  cleaner.

Per the user's standing **file-per-future-owner** rule: if the
gate-walk lands on L2, the gate-walk PR seeds one phase-skeleton
file per phase per
[`shared.md`](../../../../spec/planning/shared.md) "Parent-doc
child contracts" (Parent-promotion stub seeding) and the
[`task-plan.md`](../../../../spec/planning/task-plan.md)
task/phase gate's seeding step. If L1, the task plan absorbs phase
content inline and seeds nothing.

The decision is not pre-baked. *Verified by:*
[`task-plan.md` "PR-count predictions need a branch
test"](../../../../spec/planning/task-plan.md);
[`shared.md` "Parent-promotion stub
seeding"](../../../../spec/planning/shared.md).

## Plan-structure handoff (for the promotion-gate walk)

The plan doc is durable and carries the contracts the
gate-walk re-confirms. With OD1–OD6 resolved at the OD walk
(see "Decisions resolved at OD walk" above), the gate-walker's
remaining work is the standard promotion-gate self-review per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
`` `In draft` → `Proposed` `` promotion gate; this handoff
names the surfaces that walk applies to and the gate-specific
seeding step the flip carries.

- This is a **standalone graduated task plan**. The plan's only
  tracking surfaces are this plan tree, the graduating
  `humanize-forest-actor` backlog entry, and the new
  `progress-cell-active-state-and-actor` entry (both already
  mutated in this drafting change).
- **N = 3** (OD6 = L2). The plan doc carries a `## Phase
  Contracts` section with the per-phase WHAT and the
  **p1 → p2 → p3** sequence. The gate-walk PR seeds three
  phase-skeleton files (`post-m2-ux-correction-p1.md`, `-p2.md`,
  `-p3.md`) per
  [`shared.md`](../../../../spec/planning/shared.md)
  "Parent-promotion stub seeding" at the same time it flips the
  task plan `In draft → Proposed`.
- The task is **product-facing top-to-bottom**; the
  [`milestone.md`](../../../../spec/planning/milestone.md)
  "Product acceptance and per-leaf validation" rule binds. The
  full product-acceptance walkthrough is the task-terminal
  Validation Gate that runs at p3's implementing PR; per-phase
  gates live in each phase plan (seeded by the gate-walk PR).
  p3's implementing PR merges at Status `Validating`; the
  follow-up doc-only commit flips both p3's phase plan and the
  task plan `Validating → Landed` in a single commit per
  [`task-plan.md`](../../../../spec/planning/task-plan.md)
  "Task plan terminal state when N ≥ 2" + "Plan-to-PR
  Completion Gate" Post-merge-validation exception.
- **Files the plan contracts (per-phase, as estimate):**
  - **p1:** [`render.go`](../../../../internal/site/render.go)
    body rule (F1); [`forest.go`](../../../../internal/site/forest.go)
    summary marker styling (F8).
  - **p2:** [`forest.go`](../../../../internal/site/forest.go)
    `actor-marker` text (F4); possibly an additive
    `WorkInstanceView` passthrough in
    [`tree.go`](../../../../internal/site/tree.go) / the loader
    in [`site.go`](../../../../internal/site/site.go) — captured
    via the implementing PR's `## Estimate Deviations` if so.
  - **p3:** [`forest.go`](../../../../internal/site/forest.go)
    default-cells render (F3a), nested-`<details>` body
    disclosure rendering markdown with anchors stripped (F7);
    [`roster.go`](../../../../internal/site/roster.go) K3
    every-entry-opens render (F9); possibly an additive
    always-known-fields extension on `RosterEntry` /
    [`site.go`](../../../../internal/site/site.go) — captured
    via `## Estimate Deviations` if so.
  - **No edit** to the walker, the API, the schema, or any
    dependency in any phase.
- **Decisions the plan locks:** three supersessions
  (D1 = m2 t3 D5; D2 = stub-children stub-render contract;
  D3 = t4 "no metadata ⇒ plain row"), plus the OD-walk
  outcomes (folded into Contracts C2 / C4 and the OD-walk
  outcomes block in the plan's `## Status` section). The SD1–SD6
  settled scoping decisions are absorbed as plan-level locked
  WHAT.
- **Cross-Cutting Invariants the plan names:**
  C-INV-1 cell-anchor (the cell DOM is the per-stage attachment
  point a future F3b actor marker / state overlay hangs off
  additively; underwrites the
  `progress-cell-active-state-and-actor` backlog entry's
  linear-not-cliff migration claim); C-INV-2 region-boundary
  preservation (the m2 file-enforced shell / forest / roster
  region invariants hold — shell body-rule edits in `render.go`
  are the only `render.go` edit, no layout/composition change);
  C-INV-3 actor-tag-preserved (F4 changes the rendered text
  inside the `actor-marker` span; the span itself, its
  placement, and its attachment to every work-instance remain —
  the m2 milestone no-regress invariant); C-INV-4 walk-on-
  every-request preserved; C-INV-5 additive (no spec / API /
  schema / dependency change).
- **Validation Gate the plan owns (at p3):** the
  product-acceptance walkthrough — run the site against
  `docs/plans/` per
  [`docs/dev.md`](../../../../docs/dev.md); seed bound +
  unbound sessions, with and without metadata, via the
  `register` subcommand
  ([`cmd/workstream-tracker/register.go`](../../../../cmd/workstream-tracker/register.go))
  per [`docs/dev.md`](../../../../docs/dev.md) "Registering and
  completing a session"; open in **both** OS dark mode and OS
  light mode; observe each F-finding's corrected behavior in
  each mode. Approval is recorded against the walkthrough
  revision per
  [`task-plan.md`](../../../../spec/planning/task-plan.md)
  "Product-facing leaf tasks output a demo walkthrough." Per-
  phase Validation Gates (p1: dark/light chrome + marker; p2:
  actor-marker name; p3: full walkthrough subsuming all three
  phases) live in each phase plan when seeded.
- **Backlog mutations:** **already executed in this drafting
  change** (the three entries are in
  [`docs/backlog.md`](../../../backlog.md) at the close of the
  OD walk). The graduation-flip is on `humanize-forest-actor`;
  the new entries are `post-m2-ux-correction` (Graduated) and
  `progress-cell-active-state-and-actor` (Open). The work that
  closes the `humanize-forest-actor` intent lands at **p2**.
- **The gate-walk's universal `Verified by:` step** re-confirms
  every cited line range against then-merged code; line numbers
  in this scoping doc may drift between scoping and gate-walk
  and the symbolic anchors (function names, template names,
  selectors) are the load-bearing references per
  [`shared.md`](../../../../spec/planning/shared.md)
  "Anchor preference."
- **Seeding step the gate-walk PR carries.** Per
  [`shared.md`](../../../../spec/planning/shared.md)
  "Parent-promotion stub seeding" and
  [`task-plan.md`](../../../../spec/planning/task-plan.md) the
  task/phase gate's seeding step: three phase-skeleton files
  with pre-declared slugs `post-m2-ux-correction-p1` / `-p2` /
  `-p3`, `Status: In draft`, the inherited contract block from
  the parent's `## Phase Contracts` row, and a level-appropriate
  `short_description`.

## Related Docs

- [`../README.md`](../README.md) — the durable plan doc this
  scoping doc pairs with (Status: `In draft` until the gate-walk).
- [`../../workstream-tracker-1-0/m2/README.md`](../../workstream-tracker-1-0/m2/README.md)
  — the Landed milestone whose product gaps this task corrects.
- [`../../workstream-tracker-1-0/m2/t3-doc-declared-stages.md`](../../workstream-tracker-1-0/m2/t3-doc-declared-stages.md),
  [`../../workstream-tracker-1-0/m2/t4-session-roster.md`](../../workstream-tracker-1-0/m2/t4-session-roster.md),
  [`../../stub-children-on-parent-promotion/README.md`](../../stub-children-on-parent-promotion/README.md)
  — the three Landed contracts this task's Decisions section
  supersedes (D1, D3, D2 respectively); Landed docs are not
  retro-edited.
- [`../../../backlog.md`](../../../backlog.md) — the
  `humanize-forest-actor` entry this task graduates and the new
  `progress-cell-active-state-and-actor` entry it adds.
- [`../../../../design/vision.md`](../../../../design/vision.md)
  §3 (sub-stages D/P/I/V; cell states) and §7 (placeholder vs.
  strict — this task picks placeholder for the default-cells
  case).
- [`../../../../spec/planning/task-plan.md`](../../../../spec/planning/task-plan.md),
  [`../../../../spec/planning/shared.md`](../../../../spec/planning/shared.md),
  [`../../../../spec/planning/milestone.md`](../../../../spec/planning/milestone.md)
  — the rules this scoping doc and the plan are structured
  against.
- [`../../../../internal/site/render.go`](../../../../internal/site/render.go),
  [`../../../../internal/site/forest.go`](../../../../internal/site/forest.go),
  [`../../../../internal/site/roster.go`](../../../../internal/site/roster.go),
  [`../../../../internal/site/walker.go`](../../../../internal/site/walker.go),
  [`../../../../internal/site/site.go`](../../../../internal/site/site.go),
  [`../../../../internal/api/api.go`](../../../../internal/api/api.go),
  [`../../../../cmd/workstream-tracker/register.go`](../../../../cmd/workstream-tracker/register.go)
  — the merged code grounding the reality-check inputs.
  ([`design/v0.1-design.md`](../../../../design/v0.1-design.md)
  is intentionally **not** listed — frozen v0.1 record, not an
  authority.)
