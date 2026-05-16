---
slug: workstream-tracker-1-0-m1-t4-p1
Status: Proposed
short_description: Inline per-node detail render
---

# t4 P1 — Inline per-node detail render

## Context

This is phase 1 of the t4 task plan
([`m1-t4-expanded-per-node-display.md`](m1-t4-expanded-per-node-display.md)) —
the parent task plan owns the Cross-Phase Decisions,
Cross-Cutting Invariants, and sequencing this phase plan inherits
by reference rather than restating.

The plan-tree page today renders, per node, a Status badge, the
descriptive label (with the slug in a hover tooltip), and
active-work-instance markers. t3 (Landed, #7) already parses each
doc's markdown body into a long-description value and carries it
onto every rendered node, but **no template emits it**. This
phase renders that long description inline beneath the node, and
adds an optional `related_prs` frontmatter field whose
author-curated PRs render alongside it. It is the first half of
m1's terminal "expanded per-node display" task; it is pure
read-path and ships independent value (descriptions and
manually-listed PRs become visible) without waiting on P2's
`gh` auto-discovery.

Surfaces touched conceptually: the plan-doc frontmatter contract
(one new optional field documented in the vendored spec), the
plan-tree walker, the tree-node render template, and the
plan-tree spec doc. No API/DB/schema, no subprocess, no
JavaScript. Deliberation, rejected alternatives, and
reality-check inputs live in the sibling scoping doc
([`scoping/m1-t4-expanded-per-node-display.md`](scoping/m1-t4-expanded-per-node-display.md)),
which this plan does not restate.

P1 has no pending input from a prior task: t3's
`LongDescription` carry is **Landed** (re-confirmed against
current code in the scoping reality-check), and t1/t2 are the
independent foundation track P1 does not touch. The
`In draft → Proposed` promotion-gate self-review per
[`task-plan.md`](../../../spec/planning/task-plan.md)
(end-to-end coherence, contract decision-completeness, universal
`Verified by:` walk, reality-check re-confirmation) has been run;
no open inputs remained, so Status is `Proposed`.

## Goal

Every plan-tree node renders, beneath its existing
badge/label/markers line and before its children, (a) its long
description as a static prose block when non-empty, and (b) its
related PRs as a static list when `related_prs` is present and
non-empty. An optional `related_prs` frontmatter field is
documented in the spec and parsed with the same absence-tolerance
as `short_description`. A node with neither a body nor
`related_prs` renders a byte-unchanged badge/label/markers line.
No truncation, no JavaScript, no panel, no expand/collapse
(scoping D1).

## Naming

- `related_prs` — new optional YAML frontmatter key; a block
  sequence of strings. snake_case to match `short_description`.
- `parsedDoc.RelatedPRs` `[]string` — new field on the walker's
  per-file struct.
- `PlanNode.RelatedPRs` `[]string` — new field on the render
  node, copied from the matching `parsedDoc`.
- `stringList` — unexported walker helper turning the tolerant
  `related_prs` frontmatter value (a goldmark-meta sequence
  decoded as `[]interface{}`, or absent) into `[]string`,
  dropping non-string elements.

## Contracts

Final shapes. Estimate-shaped sections (Files to touch, Commit
Boundaries) are labeled as estimates per
[`shared.md`](../../../spec/planning/shared.md) "Plan content is
a mix of rules and estimates." These bullets state the observable
end-state each surface must reach; implementation technique is
non-binding guidance under Execution Steps.

### Frontmatter / spec contract

- `related_prs` is an **optional** frontmatter field: a YAML
  block sequence of strings. A doc carrying it renders with it;
  a doc omitting it renders with no warning, error, or skip —
  absence is a valid, common state. A present-but-empty list
  behaves identically to absence.
- Each entry is an **absolute-URL string** identifying a PR
  (an `https://`/`http://` URL, e.g. a `…/pull/N` link). The
  spec documents `related_prs` entries as absolute URLs. P1
  does **not** accept or expand `#NNN` / `owner/repo#NNN`
  shorthand and does **not** canonicalize PR-reference syntax:
  shorthand support is out of scope (see Out Of Scope below),
  and canonical-identity normalisation (needed only to dedupe
  frontmatter entries against `gh`-discovered ones) is P2's
  contract (scoping D3). P1 has no second source to dedupe
  against, so it needs no canonical key. A `gh`-discovered PR
  is already an absolute URL (P2's `gh … --json url`), so the
  two phases share one entry shape and no P1 canonicalization
  is implied.
- The spec documents `related_prs` as an optional, additive
  field of absolute-URL strings that pre-existing docs and
  vendored consumers remain valid without. `Verified by:`
  [`spec/planning/shared.md` "Plan-doc identity (slug)"](../../../spec/planning/shared.md)
  lines ~135-154 carry the `short_description` optional-field
  block; `related_prs` documents adjacent with identical
  posture.

### Walker contract (`internal/site/walker.go`)

- `parsedDoc` gains `RelatedPRs []string`. `parsePlanDoc` reads
  the `related_prs` frontmatter key through the same tolerant
  pattern as the scalar fields: a missing key, a non-sequence
  value, or a sequence with non-string elements never errors or
  skips the doc — it yields an empty (or partial) `[]string`.
  `Verified by:`
  [`parsePlanDoc` in walker.go](../../../internal/site/walker.go)
  already applies a tolerant comma-ok string assertion when
  reading `Status` from frontmatter (absence yields the zero
  value, never an error); the list field carries the same
  absence-tolerance element-wise via `stringList`.
- A YAML block sequence decodes from goldmark-meta as
  `[]interface{}` whose elements are `string`. `Verified by:`
  `goldmark-meta v1.1.0` (pinned in
  [`go.mod`](../../../go.mod)) imports `gopkg.in/yaml.v2 v2.3.0`
  and unmarshals frontmatter into `map[string]interface{}`
  (`goldmark-meta@v1.1.0/meta.go:18,140-141`); under yaml.v2 a
  block sequence into `interface{}` is `[]interface{}` of
  `string` — confirmed at this plan's promotion gate, recorded
  in scoping reality-check. `stringList` asserts each element
  to a string and drops non-strings; asserting the value as a
  `[]string` directly would fail, so that shortcut is not used.

### Tree contract (`internal/site/tree.go`)

- `PlanNode` gains `RelatedPRs []string`, copied from the
  matching `parsedDoc` in the existing `buildTree` node-build
  loop alongside `LongDescription`. No new tree walk. `Verified
  by:` [`buildTree` in tree.go](../../../internal/site/tree.go)
  builds each `PlanNode` in one loop already copying
  `LongDescription`.

### Render contract (`internal/site/render.go`)

- The `{{define "node"}}` block, after the existing
  badge/label/markers line and before the `.Children` `<ul>`,
  emits: a long-description block when `.LongDescription` is
  non-empty, and a related-PR list when `.RelatedPRs` is
  non-empty. Both are omitted entirely when their source is
  empty — a node with neither produces output byte-identical to
  today's for that line and its children. `Verified by:`
  [`indexTmpl` in render.go](../../../internal/site/render.go)
  is `html/template` (auto-escaping in element/attribute
  context), so descriptions and PR strings are
  injection-safe without manual escaping; the existing block
  already conditionally emits `.WorkInstances` and `.Children`,
  establishing the conditional-emit pattern.
- Each related-PR entry renders as a hyperlink whose href and
  visible text are both the entry URL, in the `html/template`
  auto-escaped attribute and text contexts. An entry that is
  **not** a well-formed absolute URL renders as escaped plain
  text with no anchor — never an error, skip, or broken
  in-page (`#…`) link. This url-vs-text classification is a
  render-safety decision (do not emit a broken/anchor href),
  not PR-reference parsing: P1 derives no URL from shorthand,
  it only declines to linkify a non-URL. The list order is the
  frontmatter order, unchanged.
- The long description renders as plain text (whitespace/newline
  preserved via CSS, not parsed as HTML/markdown) — it is doc
  body text, emitted in `html/template` auto-escaped text
  context. New CSS classes (`.long-desc`, `.related-prs`) are
  added to the existing inline `<style>`; no existing rule is
  removed or renamed (nothing to rename — this is purely
  additive surface).
- The Status badge, label, slug tooltip, work-instance markers,
  child nesting, and empty-state path are unchanged.

## Cross-Cutting Invariants

This phase plan inherits the task plan's Cross-Cutting
Invariants by reference (per
[`task-plan.md`](../../../spec/planning/task-plan.md) "How a
phase plan cites its parent task plan" — cite, do not
duplicate): see
[`m1-t4-expanded-per-node-display.md`](m1-t4-expanded-per-node-display.md)
"Cross-Cutting Invariants" (optional-frontmatter-degrades-
silently, walk-on-every-request, additive-spec-change,
bare-bones-no-JS). The one with ≥ 2 sites that must agree
*within P1*: **optional `related_prs` degrades silently** —
three sites must agree (the `stringList` read in `parsePlanDoc`,
the `RelatedPRs` copy in `buildTree`, and the conditional emit
in the `node` template). The plan names it once so the
end-to-end self-review walks all three.

## Files to touch

*Estimate of expected shape — implementation may revise if a
structural call requires it; deviations are reported per
[`task-plan.md`](../../../spec/planning/task-plan.md)
"Plan-to-PR Completion Gate" with the `## Estimate Deviations`
PR-body callout.*

**Modify:**

- `internal/site/walker.go` — `parsedDoc.RelatedPRs` field; the
  `stringList` helper; the tolerant `related_prs` read in
  `parsePlanDoc`.
- `internal/site/walker_test.go` — cover `related_prs`: a doc
  with a populated list parses it; a doc without it yields an
  empty slice and no error; a doc whose `related_prs` has a
  non-string element drops that element and does not error.
- `internal/site/tree.go` — `PlanNode.RelatedPRs` field; copy
  it in the `buildTree` node-build loop.
- `internal/site/tree_test.go` — assert `PlanNode.RelatedPRs`
  is carried from the parsed doc (and is empty, not nil-panic,
  when absent).
- `internal/site/render.go` — the `node`-template detail block
  (long-description + related-PR list, each conditional) and
  the additive `.long-desc` / `.related-prs` CSS rules.
- `internal/site/render_test.go` — if present, assert the
  rendered HTML contains the long description and PR entries
  when set, and that a field-less/body-less node's node-line
  output is unchanged; if no render test file exists, add one
  covering these.
- `spec/planning/shared.md` — document the optional
  `related_prs` field (additive), adjacent to the
  `short_description` block.
- `design/v0.1-design.md` — §7 "What the Website Renders":
  reflect that a node now also renders inline long description
  and author-curated related PRs (P2 will note the `gh` source
  later).
- `docs/plans/workstream-tracker-1-0/m1-v0-2.md` — the t4 row,
  the resolved "Long-description rendering location (t4)"
  decision, and the "t4 phase structure" note were updated by
  the *drafting* change; this implementing PR re-touches the
  file only to advance the t4 row to mirror P1's landing (see
  Documentation currency).

**New:**

- None beyond the t4 plan docs and the sibling scoping doc
  (created at drafting, not by this PR).

**Intentionally not touched** *(estimate — where we don't
expect changes, not a hard prohibition)*:

- `internal/api/*`, `internal/db/*` — no API/DB/schema surface;
  P1 is read-path only.
- Any subprocess / `os/exec` / `gh` integration — that is P2's
  surface, deliberately absent from P1.
- Any caching/file-watch layer — none exists and none is
  introduced (walk-on-every-request invariant).

## Validation Gate

No project build/test wrapper exists. `Verified by:` repo root
has no `Makefile`/`justfile`; `scripts/` contains only
`assemble.sh` (agent-rules vendoring — unrelated to build/test).
The canonical Go toolchain is the gate (per
[`docs/dev.md`](../../dev.md) lines 84-86):

- `gofmt -l internal` reports no files (formatting clean).
- `go build ./...` succeeds.
- `go vet ./...` clean (the new struct field tags and template
  string are vet surfaces).
- `go test ./...` passes, including the new cases:
  - walker: a doc with `related_prs` parses the list; a doc
    without it yields an empty slice and no error; a
    non-string element is dropped without error.
  - tree: `PlanNode.RelatedPRs` is carried from the parsed
    doc and is empty (not panicking) when absent.
  - render: an absolute-URL `related_prs` entry emits an
    anchor whose href and text are that URL; a non-URL entry
    (e.g. `#123`) emits escaped plain text with no anchor and
    no `href="#123"`; the long description renders when set; a
    node with neither emits an unchanged badge/label line.
- Manual: render a sample plan tree and confirm (a) a node with
  a body shows the long description inline; (b) a node with
  absolute-URL `related_prs` shows a clickable PR list inline;
  (c) a node with **neither** shows exactly the prior single
  badge/label/markers line with nothing extra — observe the
  no-field consequence (scoping D1 bans-on-surface), do not
  assume it; (d) a node with a deliberately long multi-paragraph
  body renders the full body inline (the accepted D1
  consequence — look at the page length, confirm acceptable for
  the current corpus); (e) a node with a non-URL `related_prs`
  entry shows that entry as plain text, not a link and not a
  broken in-page anchor — observe the non-URL fallback
  consequence, do not assume it.

## Execution Steps

Ordering and gates for the implementing pass. Deviating from a
step is an estimate deviation (PR-body callout), not a contract
breach.

1. **Baseline validation.** On a clean tree, run the full
   Validation Gate and confirm green *before* editing, so a
   pre-existing failure isn't misattributed.
2. **Branch hygiene.** Implement on a dedicated branch off
   current main; no unrelated changes ride along.
3. **Walker.** Add `parsedDoc.RelatedPRs`, the `stringList`
   helper, and the tolerant `related_prs` read per the Walker
   contract (absence, a non-sequence value, or a sequence with
   non-string elements each yield an empty or partial slice,
   never an error). Add walker tests (populated / absent /
   non-string-element). Parsing technique is the implementer's
   choice within that contract.
4. **Tree.** Add `PlanNode.RelatedPRs`; copy it in the
   `buildTree` loop next to `LongDescription`. Add the tree
   carry test.
5. **Render.** Add the conditional long-description + related-PR
   block to the `node` template and the additive CSS. Add/extend
   the render test.
6. **Spec + design docs.** Add the optional `related_prs` field
   to `spec/planning/shared.md` (additive); update
   `design/v0.1-design.md` §7.
7. **Parent-doc currency.** Advance the
   `m1-v0-2.md` t4 row per Documentation currency below (the
   decision/phase-structure notes were already resolved at
   drafting time).
8. **Self-review.** Run the Self-Review Audits below against the
   diff.
9. **Final validation.** Re-run the full Validation Gate
   including all four manual render observations (observe the
   no-field and long-body consequences, don't assume them).
10. **PR preparation.** PR body carries the
    `## Estimate Deviations` section (or `N/A`). Flip this phase
    plan's Status `Proposed → Landed` and advance the parent
    milestone t4 row in the same PR per the Plan-to-PR
    Completion Gate. The **task plan** stays `In draft` (P2
    pending) — do not flip it in this PR.

## Commit Boundaries

*Estimate of cohesive review chunks — the implementer may
refine.* Single implementing PR (this phase = 1 PR). Expected
commits: (a) walker + tree carry + their tests; (b) render
block + CSS + render test; (c) spec field + `design` §7 +
parent-doc t4-row advance + this plan's Status flip. The order
lets each commit build and test green.

## Self-Review Audits

Audits from
[`docs/agents/local/self-review-catalog.md`](../../agents/local/self-review-catalog.md)
mapped to this PR's diff surfaces, run at step 8:

- **validation-honesty** (validation surface) — the Validation
  Gate and all four manual render observations count as "passed"
  only if run end-to-end on the final state, including the
  no-field node and the long-body node observed in a real
  render, not asserted from template source.
- **readiness-gate-truthfulness** (Status / parent-doc surface)
  — the P1 `Proposed → Landed` flip and the parent milestone
  t4-row advance happen only after every Goal, Contract, and
  Validation step is satisfied or explicitly deferred *in this
  plan*; the task plan is **not** flipped (P2 pending), which a
  too-eager "task looks done" flip would get wrong.

Other seeded audits have no matching surface in P1:
effect-cleanup (P1 opens no effect/subscription/process — that
is P2's `gh` surface), error-surfacing-user-mutations (P1 is
read-path, no user mutation), rename-aware-diff-classification
(no rename — purely additive fields/CSS),
trigger-map-currency (no directory restructure).

## Out Of Scope

Final boundary calls (deliberation prose is in the scoping
doc):

- **`#NNN` / `owner/repo#NNN` shorthand in `related_prs`.** P1
  accepts absolute URLs only; shorthand entries render as plain
  text (the non-URL fallback), not as expanded links. Expanding
  shorthand requires repo-context resolution that is
  PR-reference canonicalization — deferred with P2's
  canonical-identity work (scoping D3), and a candidate
  follow-up if real plan docs accumulate shorthand entries.
- **Dedupe / canonical PR identity.** P1 has a single
  (frontmatter) source; the canonical key and merge belong to
  P2 (scoping D3).
- **`gh` auto-discovery.** P2's surface, deliberately absent
  from P1.

## Risk Register

- **goldmark-meta list decoding (confirmed; regression
  guard).** The decode shape (`[]interface{}` of `string` via
  `gopkg.in/yaml.v2 v2.3.0`) was verified against the pinned
  `goldmark-meta v1.1.0` at this plan's promotion gate, so it
  is not an open risk. Residual risk is a future dependency
  bump changing the shape: the walker test with a populated
  `related_prs` list asserts the parsed slice, so a shape
  change is caught at `go test` (not in production) — the test
  is the durable regression guard.
- **Long body inflates the page (scoping D1 consequence).**
  Intended, not a regression: the current corpus has short
  bodies and density/ordering work is milestone Out of Scope.
  Called out so review does not read inline full-body rendering
  as an accidental miss; the manual gate step (d) makes the
  consequence observed rather than assumed.
- **Fallback path for field-less nodes.** A node without
  `related_prs`/body must render byte-identically to today on
  its node line. Covered by the render test and manual step (c);
  flagged so review confirms additive-only, no node-line
  regression.

## Documentation currency

- `spec/planning/shared.md` — the additive `related_prs` field
  doc lands in this implementing PR (same PR, scoping D3).
- `design/v0.1-design.md` §7 — updated in this PR to reflect
  inline per-node detail.
- `docs/plans/workstream-tracker-1-0/m1-v0-2.md` — the parent
  milestone's t4 row, its resolved "Long-description rendering
  location (t4)" decision, and the "t4 phase structure" note
  were updated as part of the *drafting* change (the `—` legend
  meant "not drafted," which no longer held). This implementing
  PR re-touches the file only to advance the t4 row to mirror
  P1's landing.
- This plan's `Status` is `Proposed` (promotion-gate
  self-review complete); it flips `Proposed → Landed` in this
  implementing PR per the Plan-to-PR Completion Gate. The
  parent task plan
  [`m1-t4-expanded-per-node-display.md`](m1-t4-expanded-per-node-display.md)
  stays `In draft` (P2 pending) and is **not** flipped here.

## Backlog Impact

None. No backlog entry graduates, is deleted, split, or shifts
(see the parent task plan's Backlog Impact).

## Related Docs

- [`m1-t4-expanded-per-node-display.md`](m1-t4-expanded-per-node-display.md)
  — parent task plan (Cross-Phase Decisions, Cross-Cutting
  Invariants, sequencing this plan inherits by reference).
- [`scoping/m1-t4-expanded-per-node-display.md`](scoping/m1-t4-expanded-per-node-display.md)
  — sibling scoping doc (deliberation, transient).
- [`m1-v0-2.md`](m1-v0-2.md) — parent milestone.
- [`m1-t3-descriptive-labels.md`](m1-t3-descriptive-labels.md)
  — t3 (Landed); supplies the `LongDescription` carry P1
  renders.
- [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md)
  — the rules this plan is structured against.
- [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md)
  — cross-level planning rules; the `related_prs` spec edit
  lands here in this PR.
