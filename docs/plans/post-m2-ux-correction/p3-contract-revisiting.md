---
slug: post-m2-ux-correction-p3
Status: Validating
short_description: Contract-revisiting + task-terminal — default D/P/I/V cells (F3a) + nested-details body disclosure (F7) + every-entry-opens roster (F9); carries the full product-acceptance walkthrough
---

# Phase 3 — Contract-revisiting + task-terminal (F3a + F7 + F9)

## Status

`Validating`. The implementing PR shipped F3a + F7 + F9 with the
gate-walk-locked OD resolutions (goldmark-native Renderer
customization for OD3; additive `RegisteredAt int64` +
`LastEventAt int64` on `sessionMeta` and `RosterEntry` per OD6,
populated inside `loadSessionMetadata`'s existing baseline /
latest distinction loop; render-altitude deferrals on OD1 / OD2
/ OD4 / OD5 picked at this PR per the conservative starting
points named at the gate). The full six-finding task-terminal
walkthrough was observed against a `go run` rendering on the
implementing branch — F1 / F8 carried forward from p1 (Landed
PR #57), F4 carried forward from p2 (Landed PR #60), F3a / F7 /
F9 shipped here. Per
[`shared.md`](../../../spec/planning/shared.md) "Plan-doc
Status," the mandatory-`Validating` rule binds the product-
facing leaf, so the implementing PR merges at this state; the
post-merge doc-only commit that records product approval flips
both this phase plan and the parent task plan
`Validating → Landed` in a single commit per
[`task-plan.md`](../../../spec/planning/task-plan.md) "Task
plan terminal state when N ≥ 2."

p3 is the **task-terminal** phase of the
[`post-m2-ux-correction`](README.md) task: it ships the three
contract-revisiting findings (F7 nested-`<details>` body
disclosure; F3a default D / P / I / V Status-driven cell row;
F9 K3 every-entry-opens roster), realizes all three plan-doc
supersessions D1 / D2 / D3 (recorded at parent drafting in the
parent task plan's `## Status` →
`### Supersessions of sibling contracts` sub-block; this phase's
implementing PR carries the rendered consequence), carries the
**full product-acceptance walkthrough across all six post-m2
findings** (F1 + F3a + F4 + F7 + F8 + F9), and flips both the
phase plan's and the parent task plan's `Validating → Landed`
in a single post-merge doc-only commit per
[`task-plan.md`](../../../spec/planning/task-plan.md) "Task plan
terminal state when N ≥ 2."

The parent task plan's
[`## Phase Contracts`](README.md#phase-contracts) p3 row locked
the WHAT this phase realizes; the open decisions below — walked
at the `` `In draft` → `Proposed` `` promotion gate this commit
flips — were the HOW questions. The parent already resolved the
inputs that bind this phase at WHAT altitude: OD2 = H1 (nested
`<details>` for the body disclosure), OD3 = I2b (render
markdown with `<a>` tags stripped), OD5 = K3 (known-facts
header + raw-JSON disclosure for every roster entry), and
OD6 = L2 (the N = 3 phase shape this phase realizes). Those
resolutions are not re-litigated here. The OD walk at this
commit took a sweep posture — the contributor reviewed the
open decisions and expressed no preferences on any, which the
walk folded as "lock the WHAT-altitude shape (OD3 / OD6) at the
conservative-and-verified choice; lock every render-altitude
decision (OD1 / OD2 / OD4 / OD5) as deferred to the implementing
PR per [`task-plan.md`](../../../spec/planning/task-plan.md)
'Bans on surface require rendering the consequence'; absorb the
test-coverage OD (OD8) as planned coverage; affirm the
locked-by-structure entries (OD7 / OD9)." Resolutions fold into
each OD entry below, into the Contracts section, and into the
Files-to-touch certain set.

p1 (cosmetic defects F1 + F8) and p2 (humanize forest actor F4)
both Landed before this drafting session; p3 is the last phase
in the p1 → p2 → p3 sequence per parent
[`## Phase Contracts`](README.md#phase-contracts). The
mandatory-`Validating` rule
([`shared.md`](../../../spec/planning/shared.md) "Plan-doc
Status") binds only at the task-terminal leaf — and p3 *is* the
leaf — so this phase's lifecycle is `In draft → Proposed →
Validating → Landed`, **not** the `In draft → Proposed →
Landed` shape p1 and p2 used. The implementing PR merges at
Status `Validating`; a follow-up doc-only commit records
product approval against this plan's Validation Gate revision
and flips both p3 and the parent task plan to `Landed` in a
single commit.

### `In draft` → `Proposed` promotion gate walked

`Proposed` was reached at this drafting commit; the
[`task-plan.md`](../../../spec/planning/task-plan.md)
`` `In draft` → `Proposed` `` promotion gate was walked before
the flip.

- **Sibling-divergent-lifecycle check (gate-walk step 1).** p3
  is the task-terminal product-leaf; the
  mandatory-`Validating` rule
  ([`shared.md`](../../../spec/planning/shared.md) "Plan-doc
  Status") binds the leaf and p3 *is* the leaf, so the
  lifecycle is `In draft → Proposed → Validating → Landed` —
  **divergent from** sibling
  [`p1-cosmetic-defects.md`](p1-cosmetic-defects.md) "Status"
  and
  [`p2-humanize-forest-actor.md`](p2-humanize-forest-actor.md)
  "Status," both of which used the interior-phase `In draft →
  Proposed → Landed` shape that skips `Validating`. This is
  the **correct** divergence: p1 and p2 are interior phases
  and the leaf-keyed rule routes them past `Validating`; p3
  is the leaf and binds it. The wording in the Status block
  above + the Validation Gate section's "Approval recording"
  sub-section + parent task plan's
  [`## Phase Contracts`](README.md#phase-contracts) p3 row all
  reflect this divergence consistently. (This step is
  pre-empted up front because the analogous sibling-comparison
  gap surfaced at p2's PR review per
  [`feedback_wait_before_push_or_reply.md`](../../../) —
  catching it at step 1 rather than at review.)
- **End-to-end coherence.** Plan re-read in order; no
  contradiction between Status, Open decisions (OD1–OD9 with
  resolutions folded), Context, Goal, Contracts C1 / C2 / C3,
  Cross-Cutting Invariants (inherited C-INV-1 through
  C-INV-5; C-INV-1 explicitly realized at this phase),
  Files-to-touch (with the parent's pre-flagged `site.go`
  deviation path now closed at "certain" per OD6
  verification), the task-terminal Validation Gate (the full
  six-finding walkthrough; the divergent `Validating`
  lifecycle), Self-Review Audits, Out of Scope, Risk Register,
  and Related Docs.
- **Decision-completeness on Contracts.** Every C-clause
  locked at WHAT altitude:
  - **C1 (F7).** Header-only first view + nested `<details>`
    body disclosure + markdown-render-with-anchors-stripped
    are all parent-locked at OD2 = H1 and OD3 = I2b; OD3's
    technique choice was locked at this gate (goldmark-native
    Renderer customization — the build-vs-adopt fit since
    [`walker.go`](../../../internal/site/walker.go) lines
    12–15 already import the dep). The exact summary text
    (OD2) and exact goldmark renderer interface (OD3) are
    render-altitude / implementation-altitude deferrals
    explicitly authorized by
    [`task-plan.md`](../../../spec/planning/task-plan.md)
    "Bans on surface require rendering the consequence" +
    [`shared.md`](../../../spec/planning/shared.md) "Plans
    describe contracts, not implementation."
  - **C2 (F3a).** Default D / P / I / V row shape locked at
    parent C5 (three buckets across the seven `statusClass`
    values); OD1's palette + dashed-placeholder treatment are
    render-altitude deferrals explicitly authorized by the
    same Bans-on-surface + Plans-describe-contracts rules.
    The cell DOM shape is **not** a render-altitude deferral
    — it is load-bearing under C-INV-1 (cell-anchor) and is
    locked at this contract.
  - **C3 (F9).** K3 every-entry-opens shape locked at parent
    C4 + OD5; OD6's data-flow extension locked at this gate
    (additive `RegisteredAt int64` + `LastEventAt int64` on
    `sessionMeta` and `RosterEntry`, populated inside
    `loadSessionMetadata`'s existing baseline / latest
    distinction loop, threaded through `buildRoster`'s
    existing per-entry read; `LastEventType` **not** added —
    the type sits on the per-event payload, not the merged
    blob, and the minimal-additive C-INV-5 posture defers it
    until a concrete need arises). OD4's K3 field order +
    labels and OD5's no-metadata sentinel exact wording are
    render-altitude deferrals.

  No "decided at plan-drafting," "shape later," or "spelling
  at plan time" phrasing in any C-clause; all deferrals are
  explicitly named with the rule that authorizes them.
- **Universal `Verified by:` walk.** Every load-bearing claim
  in Contracts, Open decisions, Cross-Cutting Invariants,
  Files-to-touch, and the Risk Register carries a citation;
  symbolic anchors (`node-progress`, `node-detail`,
  `node-header`, `actor-marker`, `roster` template,
  `RosterEntry`, `sessionMeta`, `loadSessionMetadata`,
  `buildRoster`, `statusClass`, `truncateLongDesc`,
  `maxLongDescLines`, `markdownBody`, `goldmark` imports,
  `IsWellFormed`) dominate per the anchor-preference rule.
  The navigational pass (read the cited file/symbol and
  confirm the target exists) was run during drafting and
  re-confirmed at this gate.
- **Reality-check inputs re-confirmed.** Every cited surface
  spot-checked against current branch code (base
  `origin/main` `2142789`, the p2-Landed commit):
  [`forest.go`](../../../internal/site/forest.go)
  `node-progress` (line 53, the field-presence-gated row this
  phase generalizes) and `node-detail` (line 55, the
  long-description + related-PR render this phase collapses
  behind a nested disclosure) and the per-node `<details>`
  shape inside `node` (lines 67–75, including p1's
  inline-triangle marker at line 69 the F7 nested-`<summary>`
  may reuse);
  [`roster.go`](../../../internal/site/roster.go) `roster`
  template (lines 25–35, the `.Detail`-gated disclosure this
  phase generalizes; line 43 `.roster-label` carrying the
  name-then-slug fallback that the K3 known-facts header's
  `actor id` field must not silently re-introduce uuids
  through);
  [`site.go`](../../../internal/site/site.go) `sessionMeta`
  (line 131, the struct OD6 extends with `RegisteredAt` +
  `LastEventAt`), `loadSessionMetadata` (line 148, the
  baseline-vs-latest loop already in scope; `latestAt[wid]`
  tracked but unexposed at line 191; baseline `received_at`
  in scope at the register-event scan at line 187),
  `RosterEntry` (line 297, the consumer struct extended via
  `buildRoster` line 323), `Server.index` (line 47, the
  per-request orchestration matching p2's data-path
  precedent);
  [`render.go`](../../../internal/site/render.go)
  `statusClass` (line 112, seven recognized Status values +
  unknown fallback; the OD1 mapping's grounding) and
  `truncateLongDesc` (line 99, the cap the F7 truncate-first
  pipeline preserves) and `maxLongDescLines` (line 93);
  [`walker.go`](../../../internal/site/walker.go) goldmark
  import set (lines 12–15, the dep OD3 builds on without a
  new addition) and `markdownBody` (line 166, the
  plain-text-returning source the F7 markdown-render pipeline
  reads). No drift; all citations stable.
- **Always-on rules.** Required sections present (Status,
  Context, Goal, Contracts, Files to touch, Validation Gate);
  optional sections present where applicable (Cross-Cutting
  Invariants inherited from parent; Self-Review Audits; Out
  of Scope; Risk Register; Related Docs). No descent to
  implementation prescription — the OD6 additive struct
  fields are contract-altitude per the milestone's t4-amended
  data-path carve-out
  ([`m2/README.md` Cross-Task
  Invariants](../workstream-tracker-1-0/m2/README.md)) and the
  p2 OD1.a precedent; the goldmark renderer choice is
  named at technique level (goldmark-native) not interface
  level (the exact extension point is implementation-altitude);
  no fenced code; no executable predicates; no inline
  template-syntax expressions in prose (the lesson from p2's
  PR review on the now-Landed plan). No soft-commitment
  language. The `### Open decisions` sub-block stays under
  `## Status` per the parent README and p2 precedent — not a
  [`task-plan.md`](../../../spec/planning/task-plan.md)
  "Required and optional sections" variance.
- **Phase-skeleton seeding (N ≥ 2 only).** Does not apply.
  p3 is a phase plan with no further phase children; the
  parent task plan is N = 3 (p1 / p2 / p3) and seeded its
  three phase skeletons when *it* was promoted.

### Implementation history

What shipped at the implementing PR. The OD walk had already
locked the WHAT-altitude shapes (OD3 / OD6); render-altitude
deferrals (OD1 / OD2 / OD4 / OD5) closed at this PR per the
conservative starting points named at the gate, so this block
is one-line confirmations rather than rationale records.

- **OD6 data-flow extension (C3 + OD6.a).** Added `RegisteredAt
  int64` and `LastEventAt int64` to `sessionMeta` and
  `RosterEntry` in
  [`site.go`](../../../internal/site/site.go); populated inside
  the existing `loadSessionMetadata` loop from values already
  in scope. The `WHERE metadata IS NOT NULL` filter was dropped
  from the loader query so register events with no metadata
  still surface their `received_at` as the K3 "registered-at"
  field (observable conditions (b) and (d) need this), and
  later events with no metadata still contribute their
  `received_at` to the K3 "Last event" field. The loop tracks
  three maps: `registeredAt[wid]` (register event), `lastEventAt[wid]`
  (absolute latest later event across ALL events — the K3
  "Last event" surface), and `latestMetadataAt[wid]` (gates
  the write to `latest[wid]`, the Detail-fold source). The
  third map's existence — added at the post-PR-#62 review fix
  — **preserves** the t4 task plan's "Metadata read policy"
  contract and the
  [`destructive-metadata-updates`](../../backlog.md#destructive-metadata-updates)
  backlog entry's deferred semantics: a no-metadata heartbeat
  contributes nothing to Detail and does not mask a previously
  surfaced metadata-bearing event. The "every active
  work-instance gets a sessionMeta entry now even when none
  was reported" semantic change is the K3 enabler: the
  loader's `out` map writes for every id, not just the
  metadata-bearing ones. Loader-level regression coverage in
  [`roster_test.go`](../../../internal/site/roster_test.go)
  `TestLoadSessionMetadataNoMetadataLaterDoesNotMaskMetadataBearing`
  pins the contract with the specific
  (register-with-metadata, heartbeat-with-metadata,
  heartbeat-with-NO-metadata) sequence the post-#62 bot review
  named.
- **F3a default-cells row (C2).** Generalized
  [`forest.go`](../../../internal/site/forest.go)
  `node-progress` template to the field-presence branch:
  `if .ProgressStages` keeps the m2 t3 declared-stages row
  unchanged (the supersession applies only to the no-field
  branch); `else` renders the default D / P / I / V row using
  the new `progressCellClass` template helper in
  [`render.go`](../../../internal/site/render.go). Added four
  CSS variants `.progress-cell-landed` / `.progress-cell-in-
  draft` / `.progress-cell-neutral` / `.progress-cell-empty`
  per OD1's conservative starting point (reuse the
  `.status-landed` and `.status-unknown` badge palette values
  so a cell fill matches its Status badge). Cell DOM stays
  per-cell across both branches per C-INV-1 (cell-anchor) —
  verified by `TestRenderProgressRowCellDOMPreservedAcrossBranches`.
- **F7 body-disclosure (C1).** Wrapped the long description +
  related-PR list in a nested `<details class="body-disclosure">`
  in
  [`forest.go`](../../../internal/site/forest.go)
  `node-detail` template, with the same inline-triangle marker
  treatment p1 introduced for the outer collapse. Added the
  `renderLongDescBody` template helper in
  [`render.go`](../../../internal/site/render.go): truncates
  first via the existing `truncateLongDesc` cap (the
  defense-in-depth tail), then markdown-renders via the
  package-level `bodyMarkdown` goldmark instance, then returns
  `template.HTML`. The goldmark instance carries a custom
  `stripLinksRenderer` at priority 100 (lower than the default
  html renderer at 1000); under goldmark's "lower priority
  registers last, overwrites" pattern this overrides `KindLink`
  + `KindAutoLink` to emit no `<a>` wrappers — the Walk
  continues into Link children so the link text survives as
  plain text; AutoLink emits its label as escaped plain text.
  Picked OD2's summary text as "Show description" — discoverable
  and matches the discloses-the-description intent. The
  `goldmark/renderer/html` package was added to scope at
  import — `walker.go`'s pre-existing goldmark imports already
  covered `goldmark`, `goldmark-meta`, `goldmark/parser`,
  `goldmark/text`; `render.go` added `goldmark`, `goldmark/ast`,
  `goldmark/renderer`, `goldmark/util` (all from the same
  already-listed `go.mod` `github.com/yuin/goldmark v1.8.2`
  module — no `go.mod` change).
- **F9 K3 every-entry-opens (C3).** Generalized
  [`roster.go`](../../../internal/site/roster.go) `roster`
  template: removed the `{{if .Detail}}...{{else}}...{{end}}`
  branch (the t4 plain-row case D3 supersedes); every entry
  now renders a `<details class="roster-disclosure">` with a
  summary identical to the prior shape, and the disclosed
  body carries (a) the K3 known-facts header — a
  `<dl class="roster-facts">` two-column grid of slug, actor
  id, registered timestamp, last-event timestamp — followed by
  either the raw-JSON block (when `.Detail` present) or the
  `<p class="roster-empty-meta">(no reported metadata)</p>`
  sentinel (when not). Picked OD4 field order: slug, actor id,
  registered, last-event (visual reading order); each label
  follows `Label:` / value convention. Picked OD5 sentinel
  wording: "(no reported metadata)" matching the parent C4
  prose. Added the `formatEventTime` helper in
  [`render.go`](../../../internal/site/render.go) — UTC
  RFC-3339-without-T format ("2006-01-02 15:04:05 UTC") for
  human readability; zero renders as em-dash; the template
  falls back from `.LastEventAt` to `.RegisteredAt` when the
  former is zero. Added five `.roster-*` CSS classes for the
  K3 body / facts grid / sentinel.
- **OD8 test coverage (semantic-not-byte-exact per m2 t3 C7;
  presence-and-absence posture per p2 OD4.a).** All three
  candidate-assertion blocks shipped as planned coverage. Six
  prior tests updated for the superseded contracts (the prior
  test names like `TestRenderLongDescriptionInline`,
  `TestRenderProgressRowFieldlessOnlyDrafting`,
  `TestRenderProgressRowStubOnlyDrafting`,
  `TestRenderRosterExpandableDetailVsPlainRow` reflected the
  m2 t3 / t4 / stub-children contracts now superseded by D1 /
  D2 / D3; they're replaced by tests reflecting the new
  contracts, with new names like
  `TestRenderLongDescriptionInsideBodyDisclosure`,
  `TestRenderDefaultProgressRowPerStatusBucket`,
  `TestRenderProgressRowCellDOMPreservedAcrossBranches`,
  `TestRenderRosterEveryEntryOpensToK3Disclosure`,
  `TestRenderBodyHeaderOnlyByDefault`,
  `TestRenderBodyDisclosureStripsAnchors`,
  `TestRenderBodyDisclosureIndependentOfParentCollapse`,
  `TestRenderRosterFourObservableStatesAllOpen`,
  `TestRenderRosterK3TimestampsRender`). The
  `TestRenderRosterNameLabelAndSlugFallback` assertion was
  narrowed to scan `.roster-label` elements only (a new
  `rosterLabelContents` helper) — the K3 known-facts header
  IS allowed to surface the actor id as a deliberately-labeled
  facts-block field, so the prior whole-roster uuid-absence
  assertion was too strict. Forest + roster tests pass; the
  full toolchain gate (gofmt, build, vet, test ./...) passes.
- **Reviewer-facing walkthrough observed.** Started
  `go run ./cmd/workstream-tracker` against the worktree's
  `docs/plans` tree; registered four sessions covering the
  four observable conditions ((a) name-bearing bound, (b)
  no-name bound, (c) name-bearing unbound, (d) no-name
  unbound); curl'd the rendered page and verified per-finding
  acceptances structurally: F1 body rule carries
  `background-color: #fff` + `color-scheme: light`; F3a
  default-row variants present across the seven Status values
  in their three buckets, In-draft leaves render D filled +
  P / I / V dashed empty placeholders (verified on
  `demo-workstream-m1-t2` / `demo-workstream-m1-t2-p1` /
  `tool-originated-task-sessions-m2-*` / `workstream-tracker-1-0`);
  F4 forest actor-marker renders the name (slug fallback) for
  every bound session, no wst-<uuid> in any actor-marker; F7
  body-disclosure renders nested inside the per-node
  `<details>`, default-closed (header-only first view), opens
  to markdown-rendered body with zero `<a>` tags inside any
  `.long-desc` element; F8 inline-triangle marker treatment
  preserved (the `summary { list-style: none; }` rule + the
  `<span class="triangle">` in every summary); F9 each of the
  four observable states renders a roster-disclosure with K3
  known-facts header, two with the raw-JSON block (the
  name-bearing entries) and two with the no-metadata sentinel
  (the no-name entries), no wst-<uuid> in any `.roster-label`
  across all four states. The OS dark / light mode observation
  is the human-side check the post-merge approval-recording
  commit captures; the structural acceptance above is the
  reproducible-without-the-diff portion the implementing PR
  carries into review.
- **No `## Estimate Deviations` callout.** The file inventory
  shipped exactly as the gate-locked Files-to-touch named —
  `forest.go` + `roster.go` + `site.go` + `render.go` +
  `forest_test.go` + `roster_test.go` (the `render.go` slot
  was already on the certain list for the three new helpers).
  The OD6 verification at the gate had already promoted
  `site.go` from "possibly" to "certain" so the additive
  struct fields are recorded as planned scope. No additional
  files touched.
- **Slug-grammar limitation hit (OD9 confirmed out-of-scope).**
  Session registration for `post-m2-ux-correction-p3` was
  rejected by
  [`internal/slugs/slugs.go`](../../../internal/slugs/slugs.go)
  `IsWellFormed` per the known constraint (same as p1 and p2
  drafting sessions); the implementing session proceeded
  without a registered `work_instance_id`. Non-blocking per
  the CLI's "session proceeds without registration" fallback.
  Reviewer-facing walkthrough used naturally-well-formed
  slugs (`demo-workstream-m1-t1`, `nonexistent-demo`) so the
  acceptance observation was unaffected.

### Open decisions

Decomposed against merged code; resolutions folded at the OD
walk into Contracts, Files to touch, the Validation Gate, and
the implementing PR. Each entry leads with **Resolved =** and
carries a `Verified by:` citation. Per the parent README
OD-walk-outcomes precedent and the
[`task-plan.md`](../../../spec/planning/task-plan.md) "Required
and optional sections" rule, this `### Open decisions`
sub-block stays under `## Status` and is not a separate
top-level section.

The walk took a sweep posture (the contributor expressed no
preferences on any OD; the walk folded WHAT-altitude items at
the conservative-and-verified choice and render-altitude items
as deferrals to the implementing PR per
[`task-plan.md`](../../../spec/planning/task-plan.md) "Bans on
surface require rendering the consequence" +
[`shared.md`](../../../spec/planning/shared.md) "Plans describe
contracts, not implementation"). The original framing of each
OD (the decomposed shapes the gate-walk read against) is
retained as a scoping record nested under each resolved entry,
matching the
[`p2-humanize-forest-actor.md` Open decisions](p2-humanize-forest-actor.md)
precedent.

- **OD1 — F3a default-cells per-Status shape and color
  treatment. Resolved = render-altitude deferral (palette
  values + dashed-placeholder CSS technique).** The three
  shape buckets are parent-C5-locked: `Landed` ⇒ all four
  cells filled (conventionally green to match the existing
  `status-landed` palette); `In draft` ⇒ D filled and P / I /
  V rendered as dashed empty placeholders; every other Status
  (`In progress`, `Proposed`, `Validating`, `Deferred`, and
  the `unknown` fallback) ⇒ one neutral shade across all four
  cells — covering every value `statusClass` returns into
  three buckets. The exact filled color, the exact neutral
  shade, and the dashed-placeholder CSS treatment (a
  border-style change, a fill-with-stripe pattern, or a
  low-opacity empty cell) are settled at the implementing PR
  against the F3a acceptance bullet, per
  [`task-plan.md`](../../../spec/planning/task-plan.md) "Bans
  on surface require rendering the consequence" +
  [`shared.md`](../../../spec/planning/shared.md) "Plans
  describe contracts, not implementation." The conservative
  starting point named at this gate is **reuse the
  `.status-landed` and `.status-unknown` palette values from
  [`forest.go`](../../../internal/site/forest.go)** so the
  cell fill matches the existing status badge for the same
  Status (a sibling-palette alternative remains available if
  the implementing PR finds a visual reason to diverge). The
  cell DOM shape itself is **not** render-altitude — it is
  load-bearing under C-INV-1 (cell-anchor) and stays as
  every-cell-is-its-own-element across both row branches.
  *Verified by:* [`render.go`](../../../internal/site/render.go)
  `statusClass` (the seven recognized Status values + the
  `unknown` fallback this OD's mapping reads against);
  [`forest.go`](../../../internal/site/forest.go) `.progress-row`
  / `.progress-cell` / `.progress-cell-drafting` /
  `.status-landed` / `.status-unknown` styles (the existing
  palette the conservative starting point references);
  [`design/vision.md` §3 / §7](../../../design/vision.md) (the
  D / P / I / V vocabulary and the placeholder-vs-strict open
  question this contract picks placeholder for).

- **OD2 — F7 body-disclosure summary text. Resolved =
  render-altitude deferral.** Parent OD2 = H1 is locked (the
  body disclosure is a nested `<details>` / `<summary>` inside
  the parent per-node `<details>` box); the summary's
  rendered text (e.g., "Show description," "Show details,"
  "Show more," "Body," or a chevron-only glyph paired with a
  screen-reader label) is settled at the implementing PR
  against the F7 acceptance bullet, per the same
  Bans-on-surface + Plans-describe-contracts rules. *Verified
  by:*
  [`task-plan.md`](../../../spec/planning/task-plan.md) "Bans
  on surface require rendering the consequence";
  [parent Contract **C2**](README.md#c2--header-only-default-view-body-content-behind-a-separate-disclosure-f7)
  (the locked WHAT the wording must serve).

- **OD3 — F7 markdown-render and anchor-stripping technique.
  Resolved = goldmark-native Renderer customization (the
  build-vs-adopt fit; no new dependency).** Parent OD3 = I2b
  locks the WHAT (render markdown, strip `<a>` tags); the
  technique decomposition at the gate-walk surfaced three
  candidates: (a) custom `goldmark.Renderer` that emits no
  `<a>` openers (anchor stripping at the renderer layer); (b)
  HTML tree walk over the rendered fragment to strip `<a>`
  tags (adds a transitive HTML-parser dep); (c) regex over
  the rendered HTML string (fragile against nested or
  attribute-bearing anchor tags). The walk picked **(a)**
  because
  [`walker.go`](../../../internal/site/walker.go) lines 12–15
  already import `goldmark`, `goldmark-meta`,
  `goldmark/parser`, and `goldmark/text` — the dep is in scope
  for the renderer-level customization, no new package added,
  and the strip step runs at the rendering layer rather than
  as a post-process — closing the contextual-escape boundary
  risk cleanly (the rendered HTML reaching the template never
  carries an `<a>` tag in the first place). The exact
  goldmark renderer interface — whether the customization
  lives in the renderer's HTML registration
  (`renderer.NewRenderer` / `html.WithExtension`) or in an
  AST transformer (`parser.WithASTTransformers`) — is
  **implementation-altitude** and is the implementing PR's
  call. *Verified by:*
  [`walker.go`](../../../internal/site/walker.go) lines 12–15
  (the existing goldmark import set the customization
  extends);
  [`walker.go`](../../../internal/site/walker.go)
  `markdownBody` (the plain-text source the F7 pipeline
  renders);
  [parent Contract **C2**](README.md#c2--header-only-default-view-body-content-behind-a-separate-disclosure-f7)
  (the locked WHAT);
  [parent Risk Register](README.md#risk-register)
  "rendered-vs-raw contextual-escape boundary" (the falsifier
  the renderer-layer choice closes structurally).

- **OD4 — F9 K3 known-facts header field ordering and labels.
  Resolved = render-altitude deferral.** Parent C4 locks the
  five fields (slug, actor id, bound / unbound, registered-at,
  last event); the rendered order, the human-readable labels
  (e.g., `Slug:` vs. unlabelled; `Actor:` vs. `Actor id:`;
  bound / unbound as a badge versus a free word;
  `Registered:` versus `Registered at:`; the timestamp format
  — full ISO 8601, "relative" / "30 minutes ago", or both),
  and whether the order matches the existing roster row's
  `.roster-label` + `.roster-tag` order or introduces its own
  are all settled at the implementing PR against the F9
  acceptance bullet, per the same Bans-on-surface +
  Plans-describe-contracts rules. *Verified by:*
  [parent Contract **C4**](README.md#c4--roster-entry-opens-for-every-entry-f9)
  (the five-field set);
  [`roster.go`](../../../internal/site/roster.go) `roster`
  template (the existing roster row composition the K3
  disclosure body sits beneath).

- **OD5 — F9 no-metadata case display string. Resolved =
  render-altitude deferral.** Parent C4 names "(no reported
  metadata)" in the prose; the exact rendered string
  (capitalization; en-dash versus em-dash; whether to include
  a one-line explanation pointing at the always-known fields
  above) is settled at the implementing PR. The load-bearing
  observation C4 locks — and OD5 does **not** defer — is that
  no-metadata and metadata-bearing entries render the **same
  outer disclosure structure**: the disclosed body is not
  absent for a no-metadata entry; it carries the known-facts
  header plus an explicit sentinel where the raw-JSON block
  would otherwise sit. *Verified by:*
  [parent Contract **C4**](README.md#c4--roster-entry-opens-for-every-entry-f9);
  parent
  [Self-Review Audits](README.md#self-review-audits)
  "error-surfacing-user-mutations" (the audit lens that
  catches an empty-`<details>`-reads-as-broken regression).

- **OD6 — F9 RosterEntry additive fields (registered-at +
  last-event). Resolved = OD6.a (additive `RegisteredAt int64`
  + `LastEventAt int64` on `sessionMeta` and `RosterEntry`,
  populated inside `loadSessionMetadata`'s existing baseline /
  latest distinction; `LastEventType` NOT added).** Verified
  concretely against current code at the drafting and
  re-confirmed at this gate: the K3 known-facts header needs
  registered-at and last-event timestamps;
  [`sessionMeta` in
  site.go](../../../internal/site/site.go) carries only
  `Name` and `Detail` today; but inside
  [`loadSessionMetadata`](../../../internal/site/site.go) the
  loop already distinguishes the `register` event (baseline)
  from later events by `events.type` and already tracks the
  latest-later-event timestamp internally to identify the
  latest reading by `events.received_at`. Both timestamps
  (the register event's `received_at` available at the
  register-event scan step; the latest-later-event's
  `received_at` already tracked) are discarded today; the K3
  header consumes them. The extension shape: extend
  `sessionMeta` with `RegisteredAt int64` and
  `LastEventAt int64` (Unix-epoch nanoseconds, the same storage
  shape
  [`schema.go`](../../../internal/db/schema.go)
  `events.received_at` carries — every event write in
  [`api/handlers.go`](../../../internal/api/handlers.go) uses
  `now.UnixNano()`), populate inside the existing
  loader loop from values already in scope, propagate to
  `RosterEntry` through `buildRoster`'s already-existing
  per-entry metadata read, and read from the roster template.
  Data-flow shape is identical in structure to p2's OD1.a
  (additive value-fields on the per-request resolved-metadata
  struct, populated inside the existing loader pass,
  threaded to consumers through the already-existing reader
  site). **This is a certain extension, not an Estimate
  Deviation** — closes the parent's pre-flagged "possibly"
  estimate to "certain" at this gate per the
  [`feedback_deferral_decision_must_analyze_growth_and_migration_cost.md`](../../../)
  growth-and-migration-cost posture (the deferral is cheaper
  named than left open). The **sub-question on
  `LastEventType`** (whether the known-facts header also
  surfaces the latest event's type — heartbeat /
  state-transition / register — alongside its timestamp)
  resolves to **not added** at this gate: the type sits on
  the per-event payload, not the merged metadata blob (the
  raw-JSON block already shows the merged content under K3),
  and the C-INV-5 additive-minimal posture defers any further
  field until a concrete need arises. The implementing PR may
  re-open this if the F9 acceptance walkthrough exposes a
  need; otherwise the two timestamp fields are the certain
  scope. *Verified by:*
  [`site.go`](../../../internal/site/site.go)
  `loadSessionMetadata` (the baseline vs. latest distinction
  already made; the latest-event timestamp already tracked
  internally; the baseline `received_at` available at the
  register-event scan step but currently unwritten);
  [`site.go`](../../../internal/site/site.go) `sessionMeta`
  (the struct K3 extends with the two timestamp fields);
  [`site.go`](../../../internal/site/site.go) `RosterEntry`
  (the consumer struct K3 propagates the timestamps onto via
  `buildRoster`'s per-entry metadata read);
  [`schema.go`](../../../internal/db/schema.go)
  `events.received_at` (the source column's storage shape);
  [`p2-humanize-forest-actor.md` "Open
  decisions" OD1.a](p2-humanize-forest-actor.md) (the
  precedent shape this extension mirrors symmetrically).

- **OD7 — Per-phase Validation Gate shape and `Validating`
  lifecycle confirmation. Resolved = task-terminal lifecycle
  + walkthrough as documented; approval-recording = prose in
  the post-merge commit message (no sub-template).** p3 is the
  task-terminal product-leaf; the mandatory-`Validating` rule
  ([`shared.md`](../../../spec/planning/shared.md) "Plan-doc
  Status") binds the leaf and p3 *is* the leaf, so the
  lifecycle is `In draft → Proposed → Validating → Landed`
  — divergent from p1 and p2's `In draft → Proposed →
  Landed` interior shape; the divergence is correct because
  p3 is the leaf and the siblings are not. The Validation
  Gate carries the **full product-acceptance walkthrough
  across all six findings**, not a phase-scoped subset, per
  the parent
  [Validation Gate](README.md#validation-gate); the four
  observable conditions ((a) name-bearing bound, (b) no-name
  bound, (c) name-bearing unbound, (d) no-name unbound) match
  the parent's locked set. The implementing PR merges at
  Status `Validating`; a post-merge doc-only commit records
  product approval (who approved, the walkthrough revision
  approved by commit SHA, the date of approval) and flips
  both this phase plan and the parent task plan
  `Validating → Landed` in a single commit per
  [`task-plan.md`](../../../spec/planning/task-plan.md) "Task
  plan terminal state when N ≥ 2." Approval-recording shape
  resolves to **prose in the post-merge commit message**,
  matching the parent README's existing "Approval recording"
  sub-section format — the three required fields are short
  enough to live as commit-message prose; a sub-template
  would add ceremony without closing a real falsifier.
  *Verified by:*
  [`shared.md`](../../../spec/planning/shared.md) "Plan-doc
  Status" (the mandatory-`Validating` leaf-keyed rule);
  [`task-plan.md`](../../../spec/planning/task-plan.md) "Task
  plan terminal state when N ≥ 2" (the joint-flip rule);
  [`milestone.md`](../../../spec/planning/milestone.md)
  "Product acceptance and per-leaf validation" (the
  approval-recording shape);
  [parent Validation Gate](README.md#validation-gate) (the
  six-finding walkthrough this phase's gate carries verbatim
  + the four observable conditions);
  [`p1-cosmetic-defects.md`](p1-cosmetic-defects.md)
  "Validation Gate" and
  [`p2-humanize-forest-actor.md`](p2-humanize-forest-actor.md)
  "Validation Gate" (the interior-phase shape p3
  deliberately diverges from).

- **OD8 — Test-coverage posture for F3a / F7 / F9. Resolved =
  all three sub-blocks below are the planned coverage (the
  candidate assertions ARE the test scope, not a choice
  among options).** The semantic-not-byte-exact posture m2 t3
  established
  ([`t3-doc-declared-stages.md`
  Contracts](../workstream-tracker-1-0/m2/t3-doc-declared-stages.md)
  C7) is preserved; p2's OD4.a presence-and-absence posture
  is the precedent for catching silent-regression modes via
  the absence side. The implementing PR's test fixtures
  realize each block below.

  - **F3a coverage candidates.** A `Landed` node renders all
    four cells filled; an `In draft` node renders the D cell
    filled and the P / I / V cells in the dashed-empty
    treatment; nodes whose Status is `In progress`,
    `Proposed`, `Validating`, `Deferred`, or any unrecognized
    value render all four cells in the neutral shade; a doc
    that **declares** `progress_stages` continues to render
    the m2 t3 declared-stages row exactly as before (the
    field-presence branch is preserved per parent **D1**
    scope); every cell of every row is its own DOM element
    regardless of which row drew it (the load-bearing C-INV-1
    invariant). The dashed-vs-filled distinction can be
    asserted by class name (e.g., a `progress-cell-empty` or
    similar) so the test reads the structural intent, not the
    pixel-level color.
  - **F7 coverage candidates.** A node's first view contains
    no `.long-desc` element and no `.related-prs` element
    (header-only-by-default); opening the body disclosure
    renders the long description with markdown formatting
    elements present (e.g., `<p>`, `<em>`, `<strong>`,
    `<code>`, `<ul>` / `<li>` depending on the input fixture)
    and contains no `<a>` open tags (the anchor-stripping
    invariant); the parent tree-collapse `<details>` and the
    body disclosure `<details>` are distinct DOM elements at
    different nesting depths and toggling one does not
    auto-toggle the other (the independence invariant the
    parent Goal locks).
  - **F9 coverage candidates.** Every roster entry across the
    four observable states (a) / (b) / (c) / (d) renders a
    `<details>` outer wrapper (no entry renders as a plain
    row); the K3 known-facts header is present on every
    entry's disclosed body; the raw-JSON block is present
    when reported metadata exists and is replaced by the
    no-metadata sentinel string when it does not; no
    `wst-<uuid>` text surfaces in any roster label (the t4 +
    p2-locked invariant on the roster identity rendering
    that the K3 known-facts header's "actor id" field must
    not silently re-introduce — the actor id is a
    deliberately-labeled facts-block field, not a free-text
    label collision with the name-then-slug rule).

  *Verified by:*
  [`t3-doc-declared-stages.md`
  Contracts](../workstream-tracker-1-0/m2/t3-doc-declared-stages.md)
  C7 (the semantic-not-byte-exact posture);
  [`forest_test.go`](../../../internal/site/forest_test.go)
  (the existing forest-region semantic test surface);
  [`roster_test.go`](../../../internal/site/roster_test.go)
  (the existing roster-region semantic test surface);
  [`p2-humanize-forest-actor.md` "Open
  decisions" OD4.a](p2-humanize-forest-actor.md) (the
  presence-and-absence sibling precedent).

- **OD9 — Backlog mutations and slug-grammar non-scope.
  Resolved = no backlog-file mutations at this phase's
  implementing PR; the `post-m2-ux-correction` entry's close
  rides the post-merge `Validating → Landed` commit; the
  slug-grammar bug stays out of scope per PR #59's gate
  redesign decline.** The parent task plan's
  [`## Backlog Impact`](README.md#backlog-impact) recorded
  three mutations executed in the parent drafting change:
  the `humanize-forest-actor` graduation, the
  `post-m2-ux-correction` entry add, and the
  `progress-cell-active-state-and-actor` entry add for the
  F3b carve-out. **This phase's implementing PR does not
  mutate the backlog file** — both graduations are already
  Graduated, and the F3b entry is already Open. The
  `post-m2-ux-correction` entry's eventual close
  (`Graduated → Landed` per
  [`spec/backlog.md` "Entry
  lifecycle"](../../../spec/backlog.md)) lands in this same
  post-merge doc-only commit that flips the plan, since the
  entry's `Plan:` pointer is to this task plan and the
  parent's terminal-state flip is the closing trigger. Open
  at this phase: the open-decisions walk explicitly declines
  to graduate the
  [`internal/slugs/slugs.go`](../../../internal/slugs/slugs.go)
  `IsWellFormed` slug-grammar bug as part of p3's scope —
  the bug rejects every slug under the `post-m2-ux-correction`
  root because `m2` is parsed as a forbidden bare-position
  token in the root portion (same constraint that affected
  p1 and p2 drafting sessions; non-blocking — the CLI's
  fallback is "session proceeds without registration"). The
  bug is a candidate for its own scope per the parent's gate
  redesign decline (PR #59); not for p3. *Verified by:*
  [parent `## Backlog Impact`](README.md#backlog-impact);
  [`spec/backlog.md` "Entry
  lifecycle"](../../../spec/backlog.md);
  [`internal/slugs/slugs.go`](../../../internal/slugs/slugs.go)
  `IsWellFormed` (the surface; not in scope here).

## Context

This is **phase 3 of three** for the
[`post-m2-ux-correction`](README.md) task; p3 is **task-terminal**.
p3 ships the three contract-revisiting findings — the F7
header-only-by-default node-box view with body content behind a
nested `<details>` disclosure (parent Contract **C2**), the F3a
default D / P / I / V Status-driven cell row when no
`progress_stages` is declared (parent Contract **C5**), and the
F9 K3 every-entry-opens roster disclosure (parent Contract
**C4**). The three findings ship together rather than as
separate phases because all three are the **plan-doc
supersessions** D1 / D2 / D3 the task is built around: D1
(m2 t3 D5's field-presence-gates-existence becomes
field-presence-gates-shape), D2 (stub-children's exactly-the-
Drafting-cell becomes the default D / P / I / V row for a
pristine stub), and D3 (t4's no-metadata-plain-row becomes the
K3 known-facts disclosure for every entry). Concentrating the
supersession review at one PR is the parent's locked phasing
rationale ([`## Phase Contracts`](README.md#phase-contracts) p3
row) — the supersession review and the rendered consequence
sit together, not split across phases.

p3 ships **after** p1 and p2 ([both Landed](README.md)) so the
dark-mode-legible page chrome (p1's F1) and the baseline-aligned
collapse marker (p1's F8) are the visual baseline against which
the F3a / F7 / F9 changes are observed in the task-terminal
walkthrough, and so the F4 reported-name actor-marker (p2's
contract C1) is the identity baseline the walkthrough's
forest / roster identity-matching observation relies on. With
p1 + p2 Landed at this drafting time, p3's implementing PR's
walkthrough exercises the full corrected product across all six
findings — not just F3a + F7 + F9 in isolation.

p3 carries the **task-terminal Validation Gate**: the full
product-acceptance walkthrough across all six findings, the
`Validating → Landed` lifecycle for this phase plan
(mandatory-`Validating` binds the product-facing leaf per
[`shared.md`](../../../spec/planning/shared.md) "Plan-doc
Status"), and the joint flip of both this phase plan and the
parent task plan in a single post-merge doc-only commit per
[`task-plan.md`](../../../spec/planning/task-plan.md) "Task plan
terminal state when N ≥ 2." No plan-doc supersessions are
realized at this phase's implementing PR — they are already
recorded in the parent's `## Status` →
`### Supersessions of sibling contracts` sub-block (executed
at the parent drafting change). What this phase realizes is the
rendered consequence the supersessions describe.

## Goal

Open the page against `docs/plans/` per
[`docs/dev.md`](../../dev.md) and observe, in both OS dark mode
and OS light mode:

- Every node box's first view is **header-only** — the label,
  the status badge, the default-cells row, the actor markers.
  The long description and the related-PR list sit behind a
  separate disclosure that opens only on explicit user gesture,
  realized as a nested `<details>` / `<summary>` inside the
  per-node `<details>` box. When the disclosure is open, the
  long description renders as markdown (code, emphasis,
  headers, lists, etc.) with link text rendered as plain text
  (anchor tags stripped). The parent tree-collapse and the
  body disclosure are independent — opening the body does not
  toggle child boxes, and opening child boxes does not toggle
  the body (F7; parent Contract **C2**).
- Every node box renders a four-cell row labeled **D / P / I /
  V**, colored by the node's `Status`: `Landed` ⇒ all four
  cells filled; `In draft` ⇒ D filled and P / I / V dashed
  empty placeholders; every other Status (`In progress`,
  `Proposed`, `Validating`, `Deferred`, and the `unknown`
  fallback) ⇒ a neutral shade across all four cells. A doc
  that declares `progress_stages` continues to render its
  declared row unchanged (F3a; parent Contract **C5**;
  supersedes m2 t3 D5 under D1 and the stub-children
  stub-render contract under D2 — both supersessions recorded
  in the parent task plan's `## Status` block, realized here).
- Every roster entry — bound or unbound, with or without
  reported metadata — opens to the same K3-shape disclosure:
  a structured known-facts header (slug, actor id, bound /
  unbound, registered-at, last event) above an optional
  raw-JSON block. Entries with reported metadata show the
  deliberately-unstructured reported blob in the raw-JSON
  block; entries without reported metadata show an explicit
  no-metadata sentinel in place of the block — the **same
  outer disclosure structure** across every entry (F9; parent
  Contract **C4**; supersedes t4's "no metadata ⇒ plain row"
  under D3 — recorded in the parent task plan's `## Status`
  block, realized here).
- The full product-acceptance walkthrough covers all six
  post-m2 findings — F1 (page chrome legible under OS dark
  mode), F3a, F4 (reported-name forest actor-marker), F7, F8
  (baseline-aligned summary marker), F9 — observed against the
  current branch's render, not asserted from the diff (parent
  Validation Gate; the task-terminal walkthrough this phase's
  Validation Gate carries verbatim).

All while: the parent task plan's Cross-Cutting Invariants hold
— **C-INV-1 (cell-anchor)** is realized at this phase (the
default-cells row's per-cell DOM is the F3b future
attachment surface); **C-INV-2 (region-boundary preservation)**
is honored under the milestone's t4-amended data-path carve-out
([`m2/README.md` Cross-Task Invariants](../workstream-tracker-1-0/m2/README.md))
— forest edits live in [`forest.go`](../../../internal/site/forest.go),
roster edits in [`roster.go`](../../../internal/site/roster.go),
and any additive data-flow extension (OD6) lives in
[`site.go`](../../../internal/site/site.go); no shell edit;
**C-INV-3 (actor-tag preserved)** holds trivially (this phase
does not touch the `actor-marker` span p2 ships); **C-INV-4
(walk-on-every-request preserved)** holds (no caching,
file-watch, or in-memory build-up; the OD6 timestamps ride the
already-per-request
[`loadSessionMetadata`](../../../internal/site/site.go) read);
**C-INV-5 (additive, no spec / API / schema / dependency
change)** holds (no frontmatter field, no API endpoint or
request/response field, no schema column, no dependency added —
the goldmark dep used by OD3's anchor-stripping is already
imported by [`walker.go`](../../../internal/site/walker.go)).

## Contracts

This phase **realizes** parent task plan Contracts **C2**, **C5**,
and **C4** ([`README.md` Contracts](README.md#contracts)); the
phase-local C1 / C2 / C3 below restate the WHAT for this phase's
surface and are the authoritative contract clauses for the
implementing PR's review. The plan-doc supersessions D1 / D2 / D3
that bind the F3a + F9 contracts are recorded in the parent's
`## Status` → `### Supersessions of sibling contracts` sub-block
and are referenced here, not re-stated, per the
[`feedback_single_source_of_truth_per_claim.md`](../../../)
single-source-of-truth-per-claim posture.

### C1 — Header-only default view; body content behind a nested-`<details>` disclosure (F7)

Every node box's first view shows only headers: the label, the
Status badge, the default-cells row (C2 below), the actor markers
(per parent **C6**, ship-Landed at p2). The long description and
the related-PR list move behind a **separate disclosure** opened
by explicit user gesture, realized as a **nested `<details>` /
`<summary>` inside the per-node `<details>` box** (parent OD2
locked at H1) — same idiom the forest already uses for its
per-node collapse, no JavaScript, and the two `<open>` states
are independent per the native `<details>` semantics. The
disclosure is **independent of** the parent `<details>`
tree-collapse — the parent collapse controls child boxes; the
body disclosure controls the long description and the
related-PR list; the two affordances do not couple.

When the body disclosure is open, the long description is
**rendered as markdown** via the existing `goldmark` dependency
([`walker.go`](../../../internal/site/walker.go) lines 12–15;
already imported for the frontmatter / AST pass), with **anchor
(`<a>`) tags stripped** at render time so link text reads as
plain text (parent OD3 locked at I2b — link navigation would
resolve against the forest page URL rather than the source plan
doc's path, producing 404s; stripping anchors keeps the
rendered body legible without that footgun). The
`maxLongDescLines = 40` truncation cap
([`render.go`](../../../internal/site/render.go)
`maxLongDescLines` / `truncateLongDesc`) is **preserved as the
defense-in-depth tail** — truncate first, then render — so a
very long body still terminates with the existing truncation
marker even when disclosed. The rendered HTML must satisfy the
contextual-escape boundary
[`html/template`](../../../internal/site/render.go) enforces:
the goldmark output reaches the template as
`template.HTML` only when the source is trusted (the
plan-tree's own walked files; the same trust boundary the
existing `markdownBody` plain-text path already crosses
implicitly), and the anchor-strip step runs before the
template ingest so no stray `<a href>` survives into the
rendered page.

The implementing PR's narrow render-altitude calls (the exact
summary text per OD2; the exact goldmark renderer interface per
OD3) are explicitly authorized by
[`task-plan.md`](../../../spec/planning/task-plan.md) "Bans on
surface require rendering the consequence" +
[`shared.md`](../../../spec/planning/shared.md) "Plans describe
contracts, not implementation."

`Verified by:`
[parent Contract **C2**](README.md#c2--header-only-default-view-body-content-behind-a-separate-disclosure-f7)
(the WHAT this clause inherits);
[`forest.go`](../../../internal/site/forest.go) `node-detail`
template (the surface this clause's implementing PR collapses
into the nested disclosure — today the long description and the
related-PR list render unconditionally beneath `node-progress`);
[`forest.go`](../../../internal/site/forest.go) `node` template
(the existing per-node `<details>` the nested disclosure sits
inside; the F7 nested-`<details>` reuses the no-JS idiom this
template establishes);
[`render.go`](../../../internal/site/render.go)
`maxLongDescLines` / `truncateLongDesc` (the truncation cap this
clause preserves);
[`walker.go`](../../../internal/site/walker.go) `markdownBody`
(the source of the disclosed content — unchanged in shape; the
clause adds a markdown-render step on the read site, not a
walker change);
[`walker.go`](../../../internal/site/walker.go) lines 12–15
(the existing `goldmark` import set the anchor-stripping
technique builds on, unchanged in scope).

### C2 — Default D / P / I / V cell row, Status-driven, when no `progress_stages` declared (F3a)

A node whose doc declares no `progress_stages` renders a default
four-cell row labeled D / P / I / V, in that order, colored by
the node's `Status`:

- `Landed` ⇒ all four cells **filled** (exact filled color is
  render-altitude per OD1; conventionally green to match the
  existing `.status-landed` palette).
- `In draft` ⇒ D **filled**; P / I / V rendered as **dashed
  empty placeholders** (exact dashed-treatment shape is
  render-altitude per OD1).
- Every other Status the
  [`render.go`](../../../internal/site/render.go) `statusClass`
  helper recognizes (`In progress`, `Proposed`, `Validating`,
  `Deferred`) and the `unknown` fallback ⇒ all four cells
  rendered in one **neutral shade** (exact neutral value is
  render-altitude per OD1; no per-stage progression is
  inferred — the F3b deferral named in parent **Out of Scope**
  and in the
  [`progress-cell-active-state-and-actor`](../../backlog.md#progress-cell-active-state-and-actor)
  backlog entry is what holds the per-stage signal back until
  the data model supports it).

A doc that **declares** `progress_stages` continues to render
the field-presence-gated declared row exactly as m2 t3
specified — a reserved Drafting cell followed by one cell per
declared `progress_stages` entry, in document order. No
data-model change; no parser change; no frontmatter spelling
change; no inheritance. The render gate becomes "declared row
if `progress_stages` present, else default row"; both branches
render a cell row, and **the cell DOM is the same shape** —
every cell is its own DOM element regardless of which row drew
it. This is the load-bearing C-INV-1 (cell-anchor) invariant the
F3b future work attaches to: the per-cell DOM is the per-stage
anchor a future per-cell actor marker or per-cell PR-state
overlay attaches to additively. A render that collapses the row
to a single combined indicator, or encodes the row only on the
row container with no per-cell element, is the defect — not the
fix.

The Drafting cell is not labeled "Drafting" in this default row
— the label is "D" per the
[`design/vision.md` §3](../../../design/vision.md) D / P / I / V
vocabulary the default-cells row adopts. The doc-visible
`Drafting` token banned by m2 t3 D3 remains banned in
frontmatter; this contract introduces no `Drafting` token into
any doc-visible surface.

This contract supersedes m2 t3 D5 under parent **D1** (a doc
that declares no `progress_stages` no longer renders "exactly
the one reserved Drafting cell"; the field-presence
discriminator now gates *shape*, not existence) and the
stub-children stub-render contract under parent **D2** (a
pristine seeded stub renders the default D filled + P / I / V
dashed placeholders row, not the single Drafting cell — and
no parser change is needed; an absent `progress_stages` field
returns the zero value, so the default-cells row derives from
`Status` alone). Both supersessions are **recorded** in the
parent task plan's `## Status` →
`### Supersessions of sibling contracts` sub-block and
referenced here; this phase's implementing PR realizes the
rendered consequence without retro-editing the superseded
Landed docs.

`Verified by:`
[parent Contract **C5**](README.md#c5--default-d--p--i--v-cell-row-status-driven-when-no-progress_stages-declared-f3a)
(the WHAT this clause inherits + the OD1 render-altitude
deferral the implementing PR closes);
[parent
`### Supersessions of sibling contracts`](README.md#supersessions-of-sibling-contracts-recorded-at-scoping)
**D1** + **D2** (the supersession statements this clause
realizes the rendered consequence of);
[`forest.go`](../../../internal/site/forest.go) `node-progress`
template (the surface this clause's implementing PR generalizes
— today renders the field-presence-gated reserved Drafting cell
plus declared cells; the contract adds the no-`progress_stages`
default-row branch);
[`render.go`](../../../internal/site/render.go) `statusClass`
(the seven Status values the default-row's three shape buckets
read against);
[parent C-INV-1 (cell-anchor)](README.md#cross-cutting-invariants)
(the per-cell DOM invariant the F3b future work attaches to —
load-bearing under this contract);
[`design/vision.md` §3 / §7](../../../design/vision.md) (the
D / P / I / V vocabulary the default row adopts).

### C3 — Roster entry opens for every entry (F9)

Every roster entry — bound or unbound, with or without reported
metadata — opens to the same **K3-shape disclosure** (parent
OD5): a structured **known-facts header** (slug, actor id, bound
/ unbound, registered-at, last event) above an optional
**raw-JSON block** carrying the deliberately-unstructured
reported metadata when present. Every entry's disclosure shows
the same known-facts header. Entries with reported metadata
also show the raw-JSON block carrying the merged
register-baseline + latest-event metadata
[`site.go`](../../../internal/site/site.go) `resolveSessionMeta`
already returns under the t4 schema-loose posture. Entries
without reported metadata show an explicit **no-metadata
sentinel string** (exact wording is render-altitude per OD5) in
place of the raw-JSON block, so the reader sees the **same
outer disclosure structure** regardless of whether the session
reported metadata.

The two registers stay separated: the **always-known facts** are
structured because they are schematized facts the loader knows
about every active work-instance (the row's
[`models.StateActive`](../../../internal/models) classification
sources the slug + actor id; the bound / unbound classification
is `buildRoster`'s `bound[slug]` read; the registered-at
timestamp is the `register`-event's `received_at` already in
scope inside `loadSessionMetadata`; the last-event timestamp is
the same loader's already-tracked `latestAt[wid]` value);
the **reported metadata** stays schema-loose because it is the
deliberately-unstructured field t4 locked
([`t4-session-roster.md` "Cross-Cutting
Invariants"](../workstream-tracker-1-0/m2/t4-session-roster.md)).
t4's "schema-loose, not homogenized" invariant **holds** — the
K3 known-facts header is **not** a schema on the
reported-metadata side; it is a sibling surface for the
loader-known facts. The forest's name-then-slug rendering rule
(parent **C6**, Landed at p2) is **not** weakened by this
clause: the K3 known-facts header carries the `actor id` as a
deliberately-labeled facts-block field (e.g., `Actor id:
wst-<uuid>`) and is not the entry's display label — the
roster's existing `.roster-label` element continues to render
the name-then-slug fallback as the entry's identity, unchanged
from p2.

The data-flow extension this contract requires is the OD6
additive: `sessionMeta`
([`site.go`](../../../internal/site/site.go)) gains
`RegisteredAt int64` and `LastEventAt int64`
(Unix-epoch seconds, the same shape
[`schema.go`](../../../internal/db/schema.go)
`events.received_at` stores), populated inside the existing
baseline / latest distinction loop in `loadSessionMetadata` —
the data is currently tracked and discarded; this clause writes
it through. `RosterEntry` gains the same two timestamp fields
through `buildRoster`'s already-existing `meta[wi.ID]` read
site. **No new query**, no new event-log read, no new
resolution path: the values ride the per-request read already
performed; the C-INV-4 single-resolution invariant is
structurally preserved (both surfaces consume the same
per-request `meta` map keyed by `wi.ID`).

This contract supersedes t4's "no metadata ⇒ plain row" under
parent **D3**: every roster entry opens to the same K3-shape
disclosure, including one whose session reported no metadata.
The supersession is **recorded** in the parent task plan's
`## Status` → `### Supersessions of sibling contracts`
sub-block and referenced here; this phase's implementing PR
realizes the rendered consequence without retro-editing t4.

The implementing PR's render-altitude calls (the exact field
order and labels per OD4; the no-metadata sentinel string per
OD5) are explicitly authorized by
[`task-plan.md`](../../../spec/planning/task-plan.md) "Bans on
surface require rendering the consequence" +
[`shared.md`](../../../spec/planning/shared.md) "Plans describe
contracts, not implementation."

`Verified by:`
[parent Contract **C4**](README.md#c4--roster-entry-opens-for-every-entry-f9)
(the WHAT this clause inherits + the K3 shape lock);
[parent
`### Supersessions of sibling contracts`](README.md#supersessions-of-sibling-contracts-recorded-at-scoping)
**D3** (the supersession statement this clause realizes the
rendered consequence of);
[`roster.go`](../../../internal/site/roster.go) `roster`
template (the existing metadata-conditional disclosure this
clause generalizes — today the `Detail`-gated branch becomes
unconditional and the no-`Detail` plain-row branch is replaced
by the K3 known-facts + sentinel shape);
[`site.go`](../../../internal/site/site.go) `sessionMeta`
(the struct extended with `RegisteredAt` + `LastEventAt`);
[`site.go`](../../../internal/site/site.go) `RosterEntry`
(the consumer struct extended via `buildRoster`'s
`meta[wi.ID]` read);
[`site.go`](../../../internal/site/site.go)
`loadSessionMetadata` (the loader where the timestamp values
are already in scope — the existing baseline / latest
distinction is the data the contract writes through);
[`t4-session-roster.md` "Cross-Cutting
Invariants"](../workstream-tracker-1-0/m2/t4-session-roster.md)
(the schema-loose-not-homogenized invariant this clause
preserves — the K3 known-facts header is a sibling surface,
not a schema imposition).

## Cross-Cutting Invariants

Inherited from the parent task plan's
[`## Cross-Cutting Invariants`](README.md#cross-cutting-invariants);
this phase's implementing change must hold every one
simultaneously. Reviewer-flag candidates when the implementing
PR brushes any of these. Cited by name per
[`shared.md`](../../../spec/planning/shared.md) "Cross-Cutting
Invariants section" — the layering discipline bars duplicating
cross-cutting rules across parent and child docs, so the
per-clause statement lives at the parent and is referenced here.

- **[C-INV-1 (cell-anchor)](README.md#cross-cutting-invariants).**
  **Realized at this phase.** The default-cells row (C2 above)
  and the declared-cells row (m2 t3) share the cell-level DOM
  shape: every cell is its own DOM element regardless of which
  row drew it. The F3b future work
  ([`progress-cell-active-state-and-actor`](../../backlog.md#progress-cell-active-state-and-actor)
  backlog entry; parent Out of Scope; deferred at the parent
  drafting with the linear-not-cliff migration claim grounded
  in this invariant) attaches a per-cell actor marker or
  per-cell PR-state overlay to that DOM additively. The
  load-bearing contract is "per-cell DOM anchor exists,
  unconditionally"; the per-cell rendered shape (the dashed-
  empty treatment OD1 picks, for instance) is render-altitude
  but does not collapse the anchor. A render that encodes the
  row only on the row container with no per-cell element, or
  collapses the row to a single combined indicator, is the
  defect — and is what the F3b carve-out's "linear, not cliff"
  migration cost depends on this phase **not** producing.
- **[C-INV-2 (region-boundary preservation)](README.md#cross-cutting-invariants).**
  Forest edits live in
  [`forest.go`](../../../internal/site/forest.go) (the F7
  `node-detail` collapse into a nested disclosure and the
  F3a `node-progress` generalization to the default-vs-declared
  branch); roster edits live in
  [`roster.go`](../../../internal/site/roster.go) (the F9
  every-entry-opens K3 render). The shell — the body rule,
  `.layout`, the `indexData` / `renderIndex` plumbing, the
  region composition — is not touched (p1 owned the F1 body
  rule; no further `render.go` edit is in this phase). The
  milestone's t4-amended data-path carve-out
  ([`m2/README.md` Cross-Task
  Invariants](../workstream-tracker-1-0/m2/README.md)) authorizes
  the OD6 additive struct fields on `sessionMeta` and
  `RosterEntry` in
  [`site.go`](../../../internal/site/site.go), with the
  populate-pass inside the existing `loadSessionMetadata` /
  `buildRoster` flow — the same data-path carve-out scope p2's
  `ActiveWorkInstance.Name` / `.Slug` extension shipped under.
- **[C-INV-3 (actor-tag preserved)](README.md#cross-cutting-invariants).**
  Holds trivially. This phase does **not** touch the
  `actor-marker` span p2 shipped — its placement inside
  `label-group`, the per-node `range .WorkInstances`
  attachment, and the name-then-slug rendering rule (parent
  Contract **C6**) all stand unchanged. The F9 K3 known-facts
  header surfaces the actor id as a deliberately-labeled
  facts-block field but does **not** re-introduce a `wst-<uuid>`
  text into any roster *label* — the entry's identity continues
  to be the name-then-slug fallback the roster's
  `.roster-label` element already renders (Landed at p2; p3
  preserves it).
- **[C-INV-4 (walk-on-every-request preserved)](README.md#cross-cutting-invariants).**
  No caching, file-watch, or in-memory build-up is introduced.
  The OD6 timestamps ride the **already-per-request**
  [`loadSessionMetadata`](../../../internal/site/site.go) read —
  the loader's existing baseline / latest distinction is the
  source; no second event-log query is added. The F7
  markdown-render step is per-render-pass over the already-
  per-request walked `LongDescription` value, with no caching.
- **[C-INV-5 (additive, no spec / API / schema / dependency
  change)](README.md#cross-cutting-invariants).** No frontmatter
  field, no API endpoint or request/response field, no schema
  column, no dependency added. The goldmark dep used by OD3's
  markdown render is already imported by
  [`walker.go`](../../../internal/site/walker.go) for the
  frontmatter / AST pass; the F7 render step uses the existing
  import. `go.mod` is not touched.

## Files to touch

*Estimate of the expected file shape, not a binding rule.
Implementation may revise this list when a structural call
requires it — including touching a file listed under
"Intentionally not touched." Any deviation is handled via the
PR-body `## Estimate Deviations` callout with the plan
reconciled to what shipped per
[`task-plan.md`](../../../spec/planning/task-plan.md) "Plan-to-PR
Completion Gate." The OD6 verification at this drafting time
**promotes** the `site.go` extension from "possibly" (parent
task plan estimate) to "certain" here — the parent's pre-flagged
deviation path is closed at planning time per
[`feedback_deferral_decision_must_analyze_growth_and_migration_cost.md`](../../../).*

- **Modify (certain):**
  - [`internal/site/forest.go`](../../../internal/site/forest.go)
    — **F3a (C2):** the `node-progress` template generalizes
    from "reserved Drafting cell + declared cells" to a
    field-presence branch ("default D / P / I / V row when
    `ProgressStages` is empty, declared row otherwise"). The
    cell DOM shape is preserved per **C-INV-1**: every cell
    stays its own element regardless of which row drew it.
    The `.progress-row` / `.progress-cell` / per-status cell
    styles in `forest-style` gain the new shape's class
    variants (e.g., the dashed-empty placeholder variant per
    OD1) but the row container's structure is unchanged.
    **F7 (C1):** the `node-detail` template collapses the
    long-description and the related-PR list into a nested
    `<details>` / `<summary>` disclosure inside the per-node
    `<details>`; the disclosed body renders the long
    description as markdown with anchors stripped per OD3,
    preserving the `truncateLongDesc` cap as the
    defense-in-depth tail. The nested `<summary>`'s marker
    treatment is render-altitude (the implementing PR may
    reuse the inline-triangle pattern p1 introduced for the
    parent collapse, or pick a deliberately-different marker
    so the two nested disclosures read distinctly — OD2's
    summary-text decision feeds this).
  - [`internal/site/roster.go`](../../../internal/site/roster.go)
    — **F9 (C3):** the `roster` template generalizes from
    "metadata-conditional disclosure" to "every-entry-opens K3
    disclosure." The current `Detail`-gated branch becomes the
    *metadata-bearing* sub-branch of the K3 disclosure; the
    current no-`Detail` plain-row branch is replaced by the K3
    known-facts header + the no-metadata sentinel string. The
    existing `.roster-summary` / `.roster-detail` / `.roster-tag`
    styles in `roster-style` gain the K3 known-facts header's
    rendered shape (a `.roster-facts` or similar block,
    render-altitude per OD4).
  - [`internal/site/site.go`](../../../internal/site/site.go)
    — **OD6 (C3 data-flow extension):** `sessionMeta` gains
    `RegisteredAt int64` and `LastEventAt int64`
    (Unix-epoch seconds, matching
    [`schema.go`](../../../internal/db/schema.go)
    `events.received_at`'s storage shape).
    `loadSessionMetadata`'s existing baseline / latest
    distinction loop populates both fields from already-in-scope
    values — the register-event branch already reads the
    register event's `received_at`; the latest-event branch
    already tracks the latest later event's `received_at` (the
    existing `latestAt[wid]` value). Both timestamps reach
    `sessionMeta` at the end of the loop, and `RosterEntry`
    gains the same two timestamp fields, populated through
    `buildRoster`'s already-existing per-entry read of the
    per-request metadata map. No new query; no second event-log
    read; no `loadActiveWorkInstances` change.

- **Modify (tests):**
  - [`internal/site/forest_test.go`](../../../internal/site/forest_test.go)
    — semantic coverage per **OD8**: the F3a default-row
    Status-driven shape across the seven `statusClass` values;
    the declared-stages row's preservation when
    `progress_stages` is declared; the cell DOM shape's
    preservation per C-INV-1 (every cell is its own element);
    the F7 header-only-by-default observation (no `.long-desc`
    and no `.related-prs` in the closed state); the F7
    body-disclosed render's markdown formatting elements
    present and `<a>` open tags absent (the anchor-stripping
    invariant); the F7 independence between the parent
    tree-collapse and the body disclosure.
  - [`internal/site/roster_test.go`](../../../internal/site/roster_test.go)
    — semantic coverage per **OD8**: every entry across the
    four observable states (a) / (b) / (c) / (d) renders a
    `<details>` outer wrapper; the K3 known-facts header is
    present on every entry; the raw-JSON block is present iff
    reported metadata exists; no `wst-<uuid>` surfaces in any
    `.roster-label` element (the t4 + p2 invariant on roster
    identity rendering — the K3 known-facts header's `actor
    id` field is the only deliberately-uuid-bearing surface);
    the `RegisteredAt` / `LastEventAt` timestamps render in
    the known-facts header.

- **Modify (docs):** none. The
  [`design/vision.md` §3 / §7](../../../design/vision.md)
  placeholder-vs-strict open question that the parent task
  plan references is unchanged at this phase — the parent
  task plan already established the cross-reference framing;
  this phase realizes the rendered consequence. The
  [`design/vision.md` §3](../../../design/vision.md) D / P / I /
  V vocabulary the default-cells row adopts is unchanged. The
  m2 README, the m2 t3 plan, the m2 t4 plan, and the
  stub-children plan are **not** retro-edited; their
  superseded clauses remain accurate as historical records of
  what shipped (Landed docs are immutable history per
  [`milestone.md`](../../../spec/planning/milestone.md)
  "already-`Landed` milestone docs are immutable history ...
  and are not retrofitted", generalized to Landed task plans
  here as the parent task plan's posture).

- **Intentionally not touched** *(estimate — where we don't
  expect changes, not a hard prohibition):*
  - [`internal/site/render.go`](../../../internal/site/render.go)
    — no shell edit. p1 owned the F1 body rule; this phase has
    no shell-level work. The `maxLongDescLines` /
    `truncateLongDesc` helpers stay unchanged (the F7
    truncate-first-then-render pipeline preserves the existing
    cap as the defense-in-depth tail; the helpers don't gain
    a rendered-markdown variant — the markdown render is a
    forest-side step on the truncated plain-text output, not
    a render.go change).
  - [`internal/site/walker.go`](../../../internal/site/walker.go)
    — no frontmatter field; no `parsedDoc` shape change; no
    `markdownBody` change. The existing goldmark import is
    consumed by the F7 render step at the forest-side surface,
    not by an extension here.
  - [`internal/site/tree.go`](../../../internal/site/tree.go)
    — no `ActiveWorkInstance` change; no `WorkInstanceView`
    change. p2's
    [`tree.go`](../../../internal/site/tree.go) extensions
    (`Name string`, `Slug string`) are the per-node
    work-instance display fields; this phase consumes them
    unchanged through the existing `node-header` template.
  - [`internal/api/`](../../../internal/api/),
    [`internal/db/`](../../../internal/db/),
    [`internal/registerclient/`](../../../internal/registerclient/),
    [`cmd/workstream-tracker/`](../../../cmd/workstream-tracker/)
    — no API, schema, register-client, or CLI change.
    [`internal/slugs/slugs.go`](../../../internal/slugs/slugs.go)
    is not touched (the slug-grammar bug that rejects this
    task's slugs at CLI registration is explicitly out of
    scope per OD9).
  - [`spec/planning/`](../../../spec/planning/) — no spec
    change.
  - The Landed plan docs the parent task plan supersedes —
    [`m2/README.md`](../workstream-tracker-1-0/m2/README.md),
    [`m2/t3-doc-declared-stages.md`](../workstream-tracker-1-0/m2/t3-doc-declared-stages.md),
    [`m2/t4-session-roster.md`](../workstream-tracker-1-0/m2/t4-session-roster.md),
    [`stub-children-on-parent-promotion/README.md`](../stub-children-on-parent-promotion/README.md)
    — are not retro-edited.
  - [`docs/backlog.md`](../../backlog.md) — no mutations at
    this phase's implementing PR (all three were executed at
    the parent drafting change per parent
    [`## Backlog Impact`](README.md#backlog-impact)). The
    `post-m2-ux-correction` entry closes at the post-merge
    doc-only commit that flips the parent task plan
    `Validating → Landed`, in the same commit, not at the
    implementing PR.

## Validation Gate

This is the **task-terminal** Validation Gate. p3 is the
product-facing leaf of the
[`post-m2-ux-correction`](README.md) task; the
mandatory-`Validating` rule
([`shared.md`](../../../spec/planning/shared.md) "Plan-doc
Status") binds, and this phase's lifecycle is `In draft →
Proposed → Validating → Landed` — divergent from p1 and p2's
`In draft → Proposed → Landed` interior shape (a deliberate
divergence per **OD7**: p1 / p2 are interior phases under the
leaf-keyed rule and so skip `Validating`; p3 is the leaf and so
binds it). The gate carries the **full product-acceptance
walkthrough across all six post-m2 findings** (F1 + F3a + F4 +
F7 + F8 + F9), not a phase-scoped subset, per the parent's
[Validation Gate](README.md#validation-gate) — the per-phase
walkthrough for the task-terminal phase **is** the task-terminal
walkthrough. p1's and p2's per-phase gates exercised their
findings in isolation at their implementing PRs; this gate
exercises all six together so the product approval evaluates
the complete corrected product.

Per
[`task-plan.md`](../../../spec/planning/task-plan.md)
"Product-facing leaf tasks output a demo walkthrough," this
section produces a concrete step-by-step product-reviewer-facing
demo walkthrough, reproducible without reading the diff. Per
[`shared.md`](../../../spec/planning/shared.md) "Plan-doc
Status," the implementing PR merges at Status `Validating`; a
follow-up doc-only commit records product approval against this
walkthrough revision and flips **both** this phase plan and the
parent task plan `Validating → Landed` in a **single commit**
per
[`task-plan.md`](../../../spec/planning/task-plan.md) "Task plan
terminal state when N ≥ 2." The same post-merge commit closes
the
[`post-m2-ux-correction`](../../backlog.md#post-m2-ux-correction)
backlog entry per
[`spec/backlog.md` "Entry
lifecycle"](../../../spec/backlog.md) — the entry's `Plan:`
pointer is to the parent task plan, and the parent's
`Validating → Landed` flip is the closing trigger.

### Reviewer-facing demo walkthrough

**Setup.** From a clean working tree at the implementing PR's
head: start the server against the dogfood plan-tree per
[`docs/dev.md`](../../dev.md) "Local Workflow" —
`go run ./cmd/workstream-tracker` with the defaults
(`PORT=8080`, `DB_PATH=./workstream-tracker.db`; the SQLite file
seeds itself on first run).

### Observable conditions

The demo exercises four work-instance states; the reviewer
arranges for each to hold at some point during the walkthrough.
**How those states are produced is the reviewer's choice** —
naturally-active sessions against the dogfood tree, the CLI's
`register` / `complete` subcommands per
[`docs/dev.md`](../../dev.md) "Registering and completing a
session," or direct database fixturing all qualify. Constraints
of any particular seeding path (e.g. the slug-grammar limitation
in [`internal/slugs/slugs.go`](../../../internal/slugs/slugs.go)
`IsWellFormed` that rejects certain root-substring patterns,
which makes `post-m2-ux-correction-*` slugs unregistrable via
the CLI today — out of scope per **OD9**) are **not**
constraints of this gate; the gate is about the rendered output
given the states, not the path that produced them. The four
states (`<DemoName>` and `<DemoUnboundName>` are reviewer-chosen
names):

- **(a) Name-bearing bound session.** An active work-instance
  attached to a plan-tree node; reported metadata carries
  `<DemoName>`.
- **(b) No-name bound session.** An active work-instance
  attached to a plan-tree node; session reported no metadata.
- **(c) Name-bearing unbound session.** An active work-instance
  whose slug is not in the walked plan-tree; reported metadata
  carries `<DemoUnboundName>`.
- **(d) No-name unbound session.** An active work-instance
  whose slug is not in the walked plan-tree; session reported
  no metadata.

The roster reads each as active in the walk-on-every-request
render.

### Observe (OS dark mode)

Open the page in OS dark mode at `http://localhost:8080/`.
Observe:

- **F1 acceptance (Landed at p1).** The page chrome — the
  `<h1>` "workstream-tracker" header, the inter-card spacing —
  is **legible** against the dark canvas (text contrasts the
  background; no element is "ghosted" or invisible). Cards
  stay legible (each card sets an explicit light background).
  Carried into the task-terminal walkthrough so the product
  approval evaluates the complete corrected product, not just
  the p3-scoped findings.
- **F8 acceptance (Landed at p1).** Every node box's collapse
  marker (the triangle p1 introduced inside `.box-header`) sits
  **aligned with the label text baseline** and does **not
  overlap** the box border. The triangle is visible in both
  open and closed `<details>` states (toggle at least one node
  to confirm both states render the marker correctly).
- **F7 acceptance (this phase).** Every node box's first view
  shows **headers only**: the label, the Status badge, the
  default-cells row (F3a below), the actor markers. **No long
  description and no related-PR list renders until the body
  disclosure is explicitly opened.** The body disclosure is a
  **nested `<details>`/`<summary>`** inside the node box's
  outer `<details>` (parent OD2 = H1). Opening the body
  disclosure shows the long description **rendered as markdown**
  (code spans, emphasis, headers, lists, etc., per the existing
  `goldmark` dep) with **link text rendered as plain text**
  (anchor tags stripped per parent OD3 = I2b; no clickable
  links to chase into 404s). Related PRs render below the body.
  **The parent `<details>` tree-collapse does not auto-toggle
  in response** (the two affordances are independent — opening
  the body does not open child boxes; opening child boxes does
  not open the body); toggle each separately to confirm.
- **F3a acceptance (this phase).** Every node box renders a
  four-cell row labeled **D P I V**. A `Landed` node renders
  all four cells filled; a stub `In draft` node (the
  `post-m2-ux-correction` plan tree contains its own `In draft`
  Status entries) renders the D cell filled and the P / I / V
  cells as dashed empty placeholders; a node whose Status is
  `In progress`, `Proposed`, `Validating`, `Deferred`, or
  unrecognized renders all four cells in the neutral shade. A
  doc whose frontmatter **declares** `progress_stages` (none
  currently exist in the dogfood tree; the walkthrough seeds
  one temporarily by adding a `progress_stages` entry to any
  leaf plan doc and reloading) renders the m2 t3 declared-
  stages row — Drafting + per-declared-entry cells — exactly
  as before. Revert the temporary frontmatter declaration
  before continuing. The cell DOM is per-cell on both branches
  (the C-INV-1 invariant; reviewer-spot-check that each cell
  is its own element so the F3b future overlay attachment is
  preserved).
- **F4 acceptance (Landed at p2).** Every node box that
  carries a registered session shows the **reported session
  `name`** inside its `actor-marker` span (or the slug fallback
  when no `name` was reported), never the raw `wst-<uuid>`
  actor. State (a) renders `<DemoName>` in the forest; state
  (b) renders its slug in the forest; **the forest and roster
  display the same identity** for the same session. Carried
  into the task-terminal walkthrough.
- **F9 acceptance (this phase).** Every roster entry — across
  states (a) through (d) — is **openable** to the same
  **K3-shape disclosure**: a known-facts header (slug, actor
  id, bound / unbound, registered-at, last event — exact order
  + labels per parent OD4) above a raw-JSON block. Entries
  with reported metadata ((a), (c)) show the deliberately-
  unstructured reported metadata t4 shipped in the raw-JSON
  block; entries without reported metadata ((b), (d)) show the
  no-metadata sentinel string in place of the block. **The
  same outer disclosure structure across every entry** is the
  K3 visual goal — a reviewer who clicks any entry sees the
  same shape. No `wst-<uuid>` text appears in any
  `.roster-label`; the actor id surfaces only inside the K3
  known-facts header as a deliberately-labeled facts-block
  field, preserving the t4 + p2 identity-rendering rule.

### Observe (OS light mode)

Open the page in OS light mode at the same URL. Observe each
acceptance above; F1 specifically must demonstrate **no
regression** in light mode (the shell's light-context
declaration p1 introduced does not break the page chrome in OS
light mode). F8's marker treatment is visually identical in
both modes. F7 / F3a / F9 acceptances are the same as in dark
mode — the contracts are color-independent (modulo the F3a
filled / neutral / dashed treatment OD1 picks; the visual
should be legible in both modes).

### Tear down

Restore the pre-demo state via whatever path produced the
observable conditions; the CLI's `complete --id <id>` /
`abandon` subcommands per
[`docs/dev.md`](../../dev.md) "Registering and completing a
session" are one such path. The roster returns to its prior
state on the next page reload.

### Approval recording

Per
[`milestone.md`](../../../spec/planning/milestone.md) "Product
acceptance and per-leaf validation" and
[`shared.md`](../../../spec/planning/shared.md) "Plan-doc
Status," the post-merge doc-only commit that flips
`Validating → Landed` for both p3 and the parent task plan
records: **who approved**, the **walkthrough revision approved**
(a commit SHA of *this* Validation Gate section), and the
**date of approval**. The same commit closes the
[`post-m2-ux-correction`](../../backlog.md#post-m2-ux-correction)
backlog entry. The implementing PR's `## Estimate Deviations`
section names any deviation from this plan's `## Files to
touch` estimates per the
[`task-plan.md`](../../../spec/planning/task-plan.md) Plan-to-PR
Completion Gate; the OD6 verification at this drafting time
**closed** the parent task plan's pre-flagged `site.go`
deviation path (the extension is named as "certain" above), so
unless the implementing PR discovers a structurally different
shape, no deviation is expected on the OD6 surface.

### Toolchain gate

`gofmt -l internal cmd` reports no files; `go build ./...`,
`go vet ./...`, `go test ./...` all pass per
[`docs/dev.md`](../../dev.md) "Local Workflow." Test coverage
includes the semantic assertions per **OD8** in
[`forest_test.go`](../../../internal/site/forest_test.go) and
[`roster_test.go`](../../../internal/site/roster_test.go): the
F3a default-row shape per Status; the F7 header-only-by-default
observation + body-disclosed render + anchor absence + parent
/ body disclosure independence; the F9 every-entry-opens
K3 disclosure + raw-JSON / no-metadata sentinel branch + no
`wst-<uuid>` in any roster label.

## Self-Review Audits

From
[`docs/agents/local/self-review-catalog.md`](../../agents/local/self-review-catalog.md);
diff surfaces are the forest-region template + the roster-region
template + the data-flow extension in `site.go` +
forest / roster semantic tests:

- **validation-honesty** — every F-finding acceptance is
  observed against a real `go run` rendering in both OS modes,
  not asserted from the diff or from the template source
  alone. The CSS strings, the template structure, and the
  rendered HTML are not the load-bearing checks; the visual /
  semantic result observed in the demo walkthrough is.
- **error-surfacing-user-mutations** — the F9 every-entry-opens
  contract must surface the no-metadata case observably: the
  K3 known-facts header is genuinely populated with the
  always-known fields and the no-metadata sentinel is rendered
  in place of the raw-JSON block, not an empty `<details>` body
  that reads as broken. The OD5 sentinel string is render-
  altitude; the audit lens catches the "disclosure opens to
  nothing" regression that would slip past a closed-state-only
  test.
- **rename-aware-diff-classification** — the F7 body-template
  reshuffle (the `node-detail` collapse into a nested
  disclosure), the F3a default-cells render relocation inside
  `node-progress`, the F9 K3 render's split of the existing
  `Detail`-gated branch into the K3 known-facts header +
  conditional raw-JSON branch, and the OD6 additive struct
  fields on `sessionMeta` / `RosterEntry` with the
  `Server.index` / `loadSessionMetadata` populate-pass — all
  relocate render and data-flow responsibilities across
  multiple sites. Hand-classify so the "header-only by
  default," "C-INV-1 cell-anchor preserved," "every roster
  entry opens to the same outer shape," and "single-resolution
  C-INV-4 preserved" claims are real, not tooling-fooled.

## Out of Scope

- **F1 (page chrome under OS dark mode) as a phase-realized
  contract.** Landed at p1; carried into the task-terminal
  walkthrough as an acceptance observation, not as a contract
  this phase realizes.
- **F4 (forest's per-node actor-marker reports the session
  name) as a phase-realized contract.** Landed at p2; same
  carry-forward posture.
- **F8 (native disclosure marker aligned with the label
  baseline) as a phase-realized contract.** Landed at p1; same
  carry-forward posture.
- **F2 (raw markdown body rendering) as an independent fix.**
  Subsumed by F7 (parent **Out of Scope**); the raw-vs-rendered
  question was resolved at parent OD3 = I2b (render markdown
  with anchors stripped).
- **F3b (active-cell actor marker; per-cell PR-state
  coloring).** Captured as the
  [`progress-cell-active-state-and-actor`](../../backlog.md#progress-cell-active-state-and-actor)
  Open backlog entry (parent **Out of Scope**); not this
  phase. The cell DOM C-INV-1 invariant this phase realizes is
  the F3b deferral's load-bearing "linear, not cliff" migration
  claim.
- **F5 (narrow-roster word-break).** Accepted as the shipped
  tradeoff (parent **Out of Scope**); no code change.
- **OS dark mode adoption beyond F1.** A paired dark-mode
  palette across the shell + every card is the larger redesign
  parent **Out of Scope** declines. This phase's F3a / F7 / F9
  color treatments are render-altitude per OD1 / OD4 / OD5 and
  must be legible in both OS modes (the OS-light-mode
  walkthrough pass observes the F3a row legibility); no full
  dark-mode palette is introduced.
- **`wst-<uuid>` actor generator change.** The actor stays the
  internal identity / idempotency key (parent **Out of Scope**
  and C-INV-5); this phase's F9 K3 known-facts header surfaces
  the actor id as a deliberately-labeled facts-block field but
  does not change the actor's generator or its storage role.
- **Spec / API / schema / dependency change.** Per C-INV-5.
  The OD6 additive struct fields are render-side data-flow,
  not a schema column; the goldmark dep is already imported by
  [`walker.go`](../../../internal/site/walker.go), not new.
- **Retro-editing the m2 / t3 / t4 / stub-children Landed plan
  docs.** Parent **Out of Scope** — the D1 / D2 / D3
  supersessions are recorded in the parent's `## Status`
  block, not by amending the superseded Landed docs.
- **Backlog file mutations at this phase's implementing PR.**
  Parent `## Backlog Impact` recorded all three mutations
  (humanize-forest-actor graduation, post-m2-ux-correction
  add, progress-cell-active-state-and-actor add) executed at
  the parent drafting change. This phase's implementing PR
  performs no backlog-file mutation. The post-merge doc-only
  commit that flips `Validating → Landed` closes the
  `post-m2-ux-correction` entry per
  [`spec/backlog.md` "Entry
  lifecycle"](../../../spec/backlog.md), in the same commit.
- **Slug-grammar fix in
  [`internal/slugs/slugs.go`](../../../internal/slugs/slugs.go).**
  The
  `IsWellFormed` bug that rejects every slug under the
  `post-m2-ux-correction` root because `m2` is a forbidden
  bare-position token in the root portion is not in scope here
  (OD9); it remains a candidate for its own scope per PR #59's
  gate-redesign decline. Non-blocking — the CLI's fallback is
  "session proceeds without registration."

## Risk Register

- **F3a default-cells row makes a node look "more done" than
  it is** (carried from parent Risk Register). A `Landed` node
  fills all four D / P / I / V cells (green under the
  conventional OD1 candidate), which a reader could misread as
  "every stage succeeded" rather than "the node reached
  terminal state." Mitigation: the cells in the default row
  are Status-driven *by design*, and the row's job is
  at-a-glance Status signaling, not per-stage truth; the
  declared-stages affordance remains the path for a doc that
  wants per-stage truth, and the cell labels D / P / I / V are
  the
  [`design/vision.md` §3](../../../design/vision.md)
  vocabulary the contributor population reads as "stage
  positions," not "stage outcomes." The Validation Gate
  walkthrough observes both the `Landed` and the `In draft`
  cases in dark and light mode so the reader's interpretation
  matches the rendered output. Carried, not blocking.
- **F7 body-disclosure open state survives child-collapse
  toggles unexpectedly** (carried from parent Risk Register).
  The disclosure is independent of the parent `<details>`
  tree-collapse — a reader could collapse the parent and find
  the body disclosure still reads as "open" on the next reload
  (a native `<details>` retains its `open` attribute on the
  current render but not across page reloads — `<details>` is
  not auto-persisted). Mitigation: the disclosure is a
  per-render-pass affordance, not user-state-bearing; the
  Validation Gate's F7 acceptance explicitly observes "the two
  affordances are independent" by toggling each separately.
  Carried, not blocking.
- **C-INV-1 cell-anchor silently weakened by an implementation
  that merges the default and declared rows into a single row
  with a per-shape branch on the row container** (carried from
  parent Risk Register). Mitigation: C-INV-1 names the
  per-cell attachment surface as the load-bearing contract;
  the reviewer-flag posture catches any render that encodes
  shape only on the row container. The OD8 test posture
  asserts every cell is its own element across both row
  branches, so a row-container-only encoding fails the
  semantic test.
- **F7 markdown render lands on the contextual-escape boundary
  unsafely** (carried from parent Risk Register; concrete
  shape: rendered HTML reaching the template as raw bytes
  rather than `template.HTML` allows the template's contextual
  escape to double-encode the markup, breaking the rendered
  result; conversely, marking untrusted bytes as
  `template.HTML` allows arbitrary inline markup through).
  Mitigation: the goldmark render runs on the trusted-source
  walked `LongDescription` value (the same trust boundary
  `markdownBody`'s plain-text path already crosses
  implicitly); the anchor-strip step runs before the template
  ingest so no stray `<a>` survives; the rendered output
  reaches the template as `template.HTML` only after the
  strip step verifies the output. The implementing PR picks
  the goldmark Renderer interface per OD3; the gate-walk
  observes the rendered body in the demo walkthrough.
- **OD6 additive timestamps drift from the loader's source of
  truth.** The K3 known-facts header's `registered-at` and
  `last event` timestamps are sourced from
  [`loadSessionMetadata`](../../../internal/site/site.go)'s
  existing baseline / latest distinction — the only correct
  source per
  [`schema.go`](../../../internal/db/schema.go)
  `events.received_at`. A future change that introduces a
  parallel timestamp resolution (e.g., a second event-log
  query, or a heartbeat-table read in a downstream caller)
  would let the K3 header show a different timestamp than the
  loader's source of truth. Mitigation: under the OD6 shape
  named above, the population pass runs **inside**
  `loadSessionMetadata`'s existing loop; both
  `sessionMeta.RegisteredAt` and `sessionMeta.LastEventAt` are
  populated in the same scan that already computes the
  baseline / latest metadata distinction. Reviewer-flag any
  commit that introduces a second event-log read, a parallel
  timestamp computation, or a separate
  `resolveSessionMeta`-style path on the roster side. The
  C-INV-4 single-resolution invariant catches this
  structurally. Carried, not blocking.
- **F9 K3 known-facts header re-introduces the `wst-<uuid>`
  identity into a roster *label* surface** (concrete
  regression-mode risk; not the same as the actor id field).
  The K3 header carries the actor id as a deliberately-labeled
  facts-block field — that is the intentional surface. A
  future template tidy that surfaces both the K3 header's
  `Actor id:` value and replaces the entry's `.roster-label`
  with a uuid would re-introduce the t4 + p2-corrected
  identity-disagreement bug. Mitigation: the OD8 test posture
  asserts the `.roster-label` element never carries the
  `wst-` prefix across every observable state. Carried, not
  blocking.

## Related Docs

- [`README.md`](README.md) — parent task plan; the locked
  Goal, Contracts C2 / C4 / C5 this phase realizes as C1 / C3 /
  C2 (the phase-local ordering differs from the parent for
  readability — F7 → F3a → F9 matches the demo walkthrough's
  observation order), all Cross-Cutting Invariants C-INV-1
  through C-INV-5 inherited here, the `## Files to touch` p3
  estimate (with the `site.go` deviation path now closed at
  "certain" per OD6 verification), the task-terminal
  Validation Gate the per-phase gate carries verbatim, the
  `## Status` → `### Supersessions of sibling contracts`
  sub-block recording D1 / D2 / D3, the `## Backlog Impact`
  block recording the mutations executed at parent drafting
  and the closing trigger this phase's terminal flip executes,
  and the OD-walk-outcomes record locking OD2 / OD3 / OD5 /
  OD6 at the WHAT-altitude inputs this phase plan reads
  against.
- [`scoping/post-m2-ux-correction.md`](scoping/post-m2-ux-correction.md)
  — paired scoping doc; **SD2** (F7), **SD4** (F9), **SD5**
  (F3a) record the scoping-time deliberation this phase
  realizes — the locked WHAT this phase carries, the rejected
  shapes (B1 / B2 for F7; D1 for F9; E1 / E2 / E3 for F3a),
  and the Reality-check inputs the parent's gate walk
  re-confirmed against the implementing branch's base.
- [`p1-cosmetic-defects.md`](p1-cosmetic-defects.md) — sibling
  p1, Landed. The dark-mode-legible chrome (F1) and the
  baseline-aligned collapse marker (F8) this phase's
  walkthrough observes as carry-forward acceptances. p1's
  inline-triangle marker treatment is the precedent the F7
  nested-`<summary>` marker call may reuse per OD2's
  render-altitude scope.
- [`p2-humanize-forest-actor.md`](p2-humanize-forest-actor.md)
  — sibling p2, Landed. The reported-name forest actor-marker
  (F4) this phase's walkthrough observes as a carry-forward
  acceptance. p2's OD1.a precedent (additive value-fields on
  the per-request resolved-metadata struct populated inside
  the existing loader pass, threaded to consumers through the
  already-existing read site) is the data-flow precedent OD6
  follows symmetrically for `sessionMeta.RegisteredAt` /
  `.LastEventAt` and `RosterEntry`'s timestamp fields. p2's
  OD4.a semantic-test posture (presence-and-absence) is the
  test-coverage precedent OD8 reads against — especially the
  absence side that catches silent-regression modes.
- [`../workstream-tracker-1-0/m2/README.md`](../workstream-tracker-1-0/m2/README.md)
  — milestone Cross-Task Invariants (region-boundary file
  enforcement + the t4-amended data-path carve-out the OD6
  extension ships under) and Cross-Task Decisions (the
  "Doc-declared-stages frontmatter shape" D5 entry that
  parent D1 supersedes; the no-JavaScript collapse-mechanism
  decision the F7 nested-`<details>` reuses unchanged). Not
  retro-edited.
- [`../workstream-tracker-1-0/m2/t3-doc-declared-stages.md`](../workstream-tracker-1-0/m2/t3-doc-declared-stages.md)
  — the Landed t3 plan whose D5 contract (field-presence
  gates existence) is superseded by parent D1 (field-presence
  gates shape). The C7 semantic-not-byte-exact test posture
  precedent OD8 reads against. Not retro-edited.
- [`../workstream-tracker-1-0/m2/t4-session-roster.md`](../workstream-tracker-1-0/m2/t4-session-roster.md)
  — the Landed t4 plan whose "no metadata ⇒ plain row"
  contract is superseded by parent D3, and whose "Session
  identity and naming" name-then-slug rule the F9 K3
  disclosure preserves (the K3 header surfaces actor id as a
  facts-block field, not a label). The schema-loose-not-
  homogenized invariant C3 preserves. Not retro-edited.
- [`../stub-children-on-parent-promotion/README.md`](../stub-children-on-parent-promotion/README.md)
  — the Landed stub-render contract whose A2 / A6 clauses
  (slug + `Status: In draft` ⇒ exactly the Drafting box)
  are superseded by parent D2 (a pristine stub renders the
  default D filled + P / I / V dashed placeholders row).
  Not retro-edited.
- [`../../backlog.md`](../../backlog.md) — the
  `post-m2-ux-correction` entry whose `Plan:` pointer is at
  the parent task plan and whose terminal close fires at the
  post-merge `Validating → Landed` flip this phase's gate
  records; the
  `progress-cell-active-state-and-actor` Open entry capturing
  the F3b carve-out whose "linear, not cliff" migration cost
  depends on this phase's C-INV-1 (cell-anchor) realization.
- [`../../../design/vision.md`](../../../design/vision.md) §3
  (sub-stages D / P / I / V vocabulary the F3a default-cells
  row adopts; eventual cell-state vocabulary *none* / *active*
  / *in-review* / *complete* this phase's per-cell DOM anchor
  preserves room for) and §7 (placeholder-vs-strict open
  question this phase's F3a default-cells case picks
  placeholder for, by design).
- [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md),
  [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md),
  [`../../../spec/planning/milestone.md`](../../../spec/planning/milestone.md),
  [`../../../spec/backlog.md`](../../../spec/backlog.md) —
  the rules this phase plan is structured against (in
  particular "Task plan terminal state when N ≥ 2," "Plans
  describe contracts, not implementation," "Bans on surface
  require rendering the consequence," the mandatory-`Validating`
  leaf-keyed rule, "Product acceptance and per-leaf
  validation," and the backlog entry-lifecycle close-on-plan-
  Landed trigger).
- [`../../dev.md`](../../dev.md) — local workflow and
  session-registration commands the Validation Gate
  walkthrough invokes; documents the same slug-grammar
  limitation OD9 declines to fix at this phase.
