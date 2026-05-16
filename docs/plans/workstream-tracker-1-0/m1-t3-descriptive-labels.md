---
slug: workstream-tracker-1-0-m1-t3
Status: Landed
---

# t3 — Descriptive tree labels

## Context

The workstream-tracker plan-tree page today labels every node
with its raw slug chain (`workstream-tracker-1-0-m1-t3`). That is
unambiguous identity but reads poorly when scanning a forest —
the eye can't tell a milestone from a task without parsing the
slug. This task makes each node read as
`<Type> <ordinal>: <short description>` (e.g.,
`Task 3: Descriptive tree labels`) while keeping the full slug
one hover away as the stable identity for URLs and git refs.

It's being done now because it is the foundation of m1's
read-experience track: t4 (expanded per-node display) consumes
the parsed `short_description` and long-description this task
introduces, and m3's cell rendering consumes the same
frontmatter field. Doing it first unblocks that track without
touching the independent t1/t2 foundation track.

Surfaces touched, conceptually: the plan-doc frontmatter
contract (a new optional field documented in the vendored spec),
the plan-tree walker that reads frontmatter, the tree-label
render layer, and the plan-tree spec doc. No API, DB, or schema
surface is involved — this is a read-path-only change.

This is an N = 1 task plan; phase content is absorbed inline.
Deliberation, rejected alternatives, and reality-check inputs
live in the sibling scoping doc
([`scoping/m1-t3-descriptive-labels.md`](scoping/m1-t3-descriptive-labels.md)),
which this plan does not restate. The
`In draft → Proposed` promotion-gate self-review per
[`task-plan.md`](../../../spec/planning/task-plan.md) has been
run (end-to-end coherence, contract decision-completeness,
universal `Verified by:` walk, reality-check re-confirmation);
no open inputs remained (see the scoping doc's "Open decisions
carried to plan-drafting"), so Status is `Proposed`.

## Goal

Plan-tree nodes render as `<Type> <ordinal>: <short description>`
with the full slug preserved as a tooltip, an optional
`short_description` frontmatter field is documented in the spec
and parsed, and the markdown body is parsed into a
long-description value carried on each node for t4 to render.
Nodes without `short_description` fall back to the slug suffix;
existing plan-tree docs without the new field render correctly
with no other behavior change.

## Naming

- `short_description` — new optional YAML frontmatter key.
  kebab-free snake_case to match the existing `slug` / `Status`
  field style in frontmatter.
- `parsedDoc.ShortDescription`, `parsedDoc.LongDescription` —
  new string fields on the walker's per-file struct.
- `PlanNode.Label`, `PlanNode.ShortDescription`,
  `PlanNode.LongDescription` — new string fields on the render
  node. `Label` is the computed display string; the other two
  are the raw parsed values carried for t4.
- `nodeTypeDisplay` — unexported function from slug node type
  to title-case display word (`milestone`→`Milestone`, etc.).
- `slugs.Slug.Position() (int, bool)` — new exported accessor
  returning the terminal segment's position; `ok == false` for
  a root slug (no position segment). Mirrors the existing
  `NodeType()` accessor shape so `internal/site` can read the
  ordinal without re-implementing the slug grammar.

## Contracts

Final shapes. Estimate-shaped sections (Files to touch) are
labeled as estimates per
[`shared.md`](../../../spec/planning/shared.md) "Plan content is
a mix of rules and estimates."

These bullets state the observable end-state each surface must
reach. Implementation technique (parsing approach, Go idioms,
template internals) is non-binding guidance under Execution
Steps — the contract is the behavior, not the mechanism.

### Frontmatter / spec contract

- `short_description` is an **optional** string frontmatter
  field. A doc carrying it is rendered with it; a doc omitting
  it is rendered without warning, error, or skip — the field's
  absence is a valid, common state.
- The spec documents `short_description` as an optional,
  additive field that pre-existing docs and vendored consumers
  remain valid without. `Verified by:`
  [`spec/planning/shared.md` "Plan-doc identity (slug)"](../../../spec/planning/shared.md)
  is the section carrying the frontmatter field block the new
  field documents alongside.
- The plan-tree carries a long-description value derived from
  the document's markdown body (the content after the
  frontmatter block). t3 parses and carries it only; where and
  how it renders is t4's call. `Verified by:`
  [`m1-v0-2.md`](m1-v0-2.md) Cross-Task Decisions
  ("Long-description rendering location (t4)").

### Slugs accessor contract (`internal/slugs/slugs.go`)

- A descendant slug exposes its terminal segment's ordinal
  position; a root slug reports "no position." The exported
  surface gains exactly one method,
  `func (s Slug) Position() (int, bool)` (position + `true`,
  or `0, false` for a root), mirroring the existing
  `NodeType()` accessor shape. `Verified by:`
  [`slugs.go` exported `Slug` methods](../../../internal/slugs/slugs.go)
  are `String/Root/NodeType/Parent` only and `segments` is
  unexported — no terminal-position accessor exists today.
- The addition is purely additive: no existing exported
  signature or behavior changes, so `internal/api` (the other
  consumer of `slugs`) is unaffected. `Verified by:`
  [`slugs.go`](../../../internal/slugs/slugs.go) — the change
  adds a method and touches no existing one.

### Walker contract (`internal/site/walker.go`)

- A parsed doc carries its short description and long
  description alongside the existing slug/Status/path
  (`parsedDoc` gains `ShortDescription`, `LongDescription`).
- When `short_description` is present and a string, it is
  carried; when absent or non-string, the carried value is
  empty and the doc is neither warned nor skipped. `Verified
  by:` [`parsePlanDoc` in walker.go](../../../internal/site/walker.go)
  already applies this tolerant-read shape to `Status`
  (`status, _ := metaData["Status"].(string)`).
- The long description is the document's markdown body — the
  content following the leading frontmatter block — with
  surrounding whitespace trimmed; an empty or absent body
  carries an empty string and is not an error. A file with no
  leading frontmatter block has no `slug` and is skipped
  upstream exactly as today. `Verified by:`
  [`parsePlanDoc` in walker.go](../../../internal/site/walker.go)
  returns the "missing or empty `slug`" error before any body
  handling, so the no-frontmatter path is unchanged.

### Tree contract (`internal/site/tree.go`)

- Every rendered node carries a display `Label` and the raw
  short/long description (`PlanNode` gains `Label`,
  `ShortDescription`, `LongDescription`, copied from the
  matching `parsedDoc`).
- `Label` resolves by this grammar:
  - Descendant with `short_description`:
    `"<Type> <ordinal>: <short_description>"`, where `<Type>`
    is the display word for the node's terminal type and
    `<ordinal>` is its terminal-segment position.
  - Root with `short_description`: the `short_description`
    alone (a root has no position segment).
  - Any node without `short_description`: the **slug suffix**
    (slug text after the root; the full slug for a root).
  - `Verified by:`
    [`buildTree` in tree.go](../../../internal/site/tree.go)
    already resolves `root` and calls `slugs.Parse` per doc, so
    type and position are available where the node is built;
    `<ordinal>` reads `slugs.Slug.Position()` (new accessor)
    and `<Type>` reads `slugs.Slug.NodeType()`.

### Render contract (`internal/site/render.go`)

- A node's visible text is its `Label`. The full slug remains
  reachable as a hover tooltip on that text, exposed via the
  HTML `title` attribute (scoping D4), emitted in an
  auto-escaped context so the slug cannot break out of the
  attribute. `Verified by:`
  [`indexTmpl` in render.go](../../../internal/site/render.go)
  is built with `html/template` (imported at the top of the
  file), which auto-escapes attribute context; no manual
  escaping is added.
- The Status badge, work-instance markers, child nesting, and
  the empty-state path are unchanged.

## Cross-Cutting Invariants

- **`short_description` is optional at every read site.** The
  parser tolerates its absence (empty string, no warning), the
  label builder falls back to the slug suffix, and the body
  parser tolerates an empty body. A plan-tree doc lacking the
  field renders exactly as before except the displayed text
  becomes the slug suffix (full slug for roots) with the full
  slug in the tooltip. Three sites must agree on this fallback:
  `parsePlanDoc`, `buildTree`'s label builder, and the `node`
  template. The plan names the rule once so self-review walks
  all three.
- **Render path stays walk-on-every-request.** No caching, file
  watching, or in-memory build-up is introduced; the label is
  computed within the existing per-request build. `Verified
  by:` [`Server.index` in site.go](../../../internal/site/site.go)
  calls `walkPlans(s.plansPath)` then `buildTree` on every HTTP
  handler invocation.
- **Spec change is additive.** The `short_description` addition
  introduces no breaking change for vendored spec consumers;
  existing frontmatter without the field stays valid.

## Files to touch

*Estimate of expected shape — implementation may revise if a
structural call requires it; deviations are reported per
[`task-plan.md`](../../../spec/planning/task-plan.md)
"Plan-to-PR Completion Gate."*

**Modify:**

- `internal/slugs/slugs.go` — add the exported
  `Slug.Position() (int, bool)` accessor (additive; required
  for the label ordinal — see Slugs accessor contract).
- `internal/slugs/slugs_test.go` — cover `Position()` for a
  task/milestone/phase slug (returns position, `true`) and a
  root slug (returns `0, false`).
- `internal/site/walker.go` — `parsedDoc` fields;
  `parsePlanDoc` short-description + body extraction.
- `internal/site/tree.go` — `PlanNode` fields; label build in
  `buildTree` (uses `slugs.Slug.Position()`); `nodeTypeDisplay`
  helper.
- `internal/site/render.go` — `node` template: display `.Label`,
  add `title="{{.Slug}}"`; the `.slug` monospace CSS rule is
  renamed to `.label` (prose, not slug text) so no dead/misnamed
  rule is left behind.
- `spec/planning/shared.md` — document optional
  `short_description` field (additive).
- `docs/plans/workstream-tracker-1-0/m1-v0-2.md` — the t3 row
  and phase-structure note were already updated by the drafting
  change (row now `In draft`, note resolved to N = 1); the
  implementing PR re-touches this file only to advance the row
  to `Landed` (see Documentation currency).

**New:**

- None beyond this plan doc and its sibling scoping doc.

**Intentionally not touched** *(estimate — these are where we
don't expect to need changes, not a hard prohibition)*:

- `internal/api/*`, `internal/db/*` — no API/DB/schema surface;
  t3 is read-path only.
- Any caching/file-watch layer — none exists and none is
  introduced (walk-on-every-request invariant).

## Validation Gate

No project build/test wrapper exists. `Verified by:` repo root
has no `Makefile` or `justfile` and `scripts/` contains only
`assemble.sh` (an agent-rules vendoring tool — see its header
comment — unrelated to build/test). The canonical Go toolchain
is the gate:

- `gofmt -l internal` reports no files (formatting clean).
- `go build ./...` succeeds.
- `go vet ./...` clean.
- `go test ./...` passes, including new cases:
  - walker: a doc with `short_description` + body parses both;
    a doc without either yields empty strings and no error.
  - tree: label is `Task <n>: <desc>` for a task with the
    field; falls back to the slug suffix without it; root with
    the field uses the bare description; root without it uses
    the full slug.
- Manual: render a sample plan tree and confirm (a) labels read
  as `<Type> <ordinal>: <desc>`, (b) hovering shows the full
  slug, (c) a field-less node shows its slug suffix and still
  carries the slug tooltip — i.e., observe the no-field
  consequence rather than assume it.

## Execution Steps

Ordering and gates for the implementing pass. The "how" the
Contracts deliberately omit (parsing approach, idioms) is
non-binding guidance here — deviating from an Execution Step is
an estimate deviation, not a contract breach.

1. **Baseline validation.** On a clean tree, run the Validation
   Gate commands and confirm they pass *before* editing, so a
   pre-existing failure isn't misattributed.
2. **Branch hygiene.** Implement on a dedicated branch off the
   current main; no unrelated changes ride along.
3. **Slugs accessor.** Add `Slug.Position()` + its test;
   `go test ./internal/slugs/` green in isolation.
4. **Walker.** Add `parsedDoc` fields and the
   short-description/body reads. Suggested (non-binding)
   approach: read `short_description` through the existing
   `meta.Get` map with a tolerant string assertion; derive the
   body by taking the source after the leading frontmatter
   block and trimming. Add walker tests.
5. **Tree + render.** Add `PlanNode` fields, the label grammar,
   the `nodeTypeDisplay` helper, and the template change
   (`Label` as text, slug into `title`). Add tree tests.
6. **Spec + parent-doc.** Add the optional `short_description`
   field to `spec/planning/shared.md` (additive); the parent
   milestone row/note were already updated at drafting time
   (Documentation currency).
7. **Self-review.** Run the Self-Review Audits below against
   the diff.
8. **Final validation.** Re-run the full Validation Gate,
   including the manual render check (observe the no-field
   consequence, don't assume it).
9. **PR preparation.** PR body carries the
   `## Estimate Deviations` section (or `N/A`); flip this
   plan's Status and the parent row to `Landed` in the same PR
   per the Plan-to-PR Completion Gate.

## Commit Boundaries

*Estimate of cohesive review chunks — the implementer may
refine.* Single implementing PR (N = 1). Expected commits:
(a) slugs accessor + test; (b) walker + tree + render + tests;
(c) spec field addition + plan/parent Status flip. Splitting
(a) keeps the shared-package change reviewable on its own; the
order lets each commit build and test green.

## Self-Review Audits

Audits from
[`docs/agents/local/self-review-catalog.md`](../../agents/local/self-review-catalog.md)
that map to this PR's diff surfaces, run at step 7:

- **validation-honesty** (validation surface) — the Validation
  Gate and manual render check are only "passed" if run
  end-to-end on the final state, including the
  no-`short_description` node observed in a real render, not
  asserted from the template source.
- **readiness-gate-truthfulness** (Status / parent-doc surface)
  — the `Proposed → Landed` flip and the parent milestone row
  advance only after every Goal, Contract, and Validation step
  is actually satisfied or explicitly deferred *in this plan*;
  no best-effort "looks done" flip.

Other seeded audits (effect-cleanup,
error-surfacing-user-mutations, rename-aware-diff-classification,
trigger-map-currency) have no matching surface here: t3 adds no
effects/listeners, no user mutations (read-path only), no
renames, and no directory restructure.

## Out Of Scope

- Long-description **rendering** (location/shape) — t4's call
  per [`m1-v0-2.md`](m1-v0-2.md). t3 only parses and carries it.
- Activity-first ordering, sub-stage cells, richer
  work-instance states — later milestones.
- Editorial frontmatter fields beyond `short_description`
  (owner, dates, tags) — deferred per the milestone Out of
  Scope.
- A dedicated copy-to-clipboard affordance — the `title`
  tooltip satisfies the preserve-the-slug contract at minimum
  surface.

## Risk Register

- **Body extraction misclassifies a doc with no leading
  frontmatter block.** Mitigation: such a doc has no `slug` and
  is skipped upstream before any body handling (per the Walker
  contract's `Verified by:`), so the long-description path
  never runs on it. Covered by a walker test for the
  no-frontmatter case regardless of the extraction approach
  chosen.
- **Spec edit ripples to adjacent frontmatter rules.**
  Mitigation: the edit is purely additive and adjacent to
  "Plan-doc identity (slug)"; the implementing PR walks
  `shared.md` for adjacent frontmatter-shape references and
  confirms none change meaning (parent-milestone Cross-Task
  Risk).
- **Fallback changes display for field-less nodes** (full slug
  → slug suffix). This is intended per the milestone t3
  contract, not a regression; the slug stays reachable via the
  tooltip. Called out so review doesn't read it as accidental.

## Documentation currency

- `spec/planning/shared.md` — the `short_description` field
  addition lands in the implementing PR (same PR, per scoping
  D2).
- `docs/plans/workstream-tracker-1-0/m1-v0-2.md` — the parent
  milestone's Task Status table t3 row is updated to `In draft`
  as part of *this* drafting change (the `—` legend means
  "not drafted," which no longer holds), and the t3
  phase-structure note is resolved to N = 1. Subsequent row
  values track this plan's Status (`In draft → Proposed →
  Landed`) and land with the PR that performs each flip.
- This plan's `Status` is `Proposed` (promotion-gate
  self-review complete); it flips `Proposed → Landed` in the
  implementing PR per the Plan-to-PR Completion Gate.

## Backlog Impact

None. No backlog entry graduates, is deleted, split, or shifts.
The [`repo-rooted-doc-links`](../../backlog.md#repo-rooted-doc-links)
entry (parent-epic deliberation) is post-1.0 and untouched by t3.

## Related Docs

- [`m1-v0-2.md`](m1-v0-2.md) — parent milestone; t3 task
  contract and deferred decisions.
- [`README.md`](README.md) — parent epic.
- [`scoping/m1-t3-descriptive-labels.md`](scoping/m1-t3-descriptive-labels.md)
  — sibling scoping doc (deliberation, transient).
- [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md)
  — the rules this plan is structured against.
- [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md)
  — cross-level planning rules; the `short_description` spec
  edit lands here.
