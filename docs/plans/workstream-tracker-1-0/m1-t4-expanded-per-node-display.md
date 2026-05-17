---
slug: workstream-tracker-1-0-m1-t4
Status: In progress
short_description: Expanded per-node display
---

# t4 — Expanded per-node display

## Context

The workstream-tracker plan-tree page today shows, per node, a
Status badge, the `<Type> <ordinal>: <short description>` label
(with the full slug in a hover tooltip), and active-work-instance
markers. t3 (Landed, #7) already parses each doc's markdown body
into a long-description value and carries it onto every rendered
node — but **nothing renders it yet**. This task makes each node
surface that detail: the long description, and the pull requests
related to the node.

It is being done now because it is the terminal task of m1's
read-experience track and the last task of the milestone: t3 laid
the parsing foundation specifically so t4 could consume it, and
the milestone's "expanded per-node display" is the last gap
between v0.1's bare-bones render and v0.2's "scannable per node"
goal. Related PRs come from two sources the milestone t4 contract
names: an optional `related_prs` frontmatter field, and
auto-discovery via `gh pr list` keyed on the slug.

Surfaces touched, conceptually: the plan-doc frontmatter contract
(one new optional field, documented in the vendored spec), the
plan-tree walker that reads frontmatter, the tree-node render
layer, the plan-tree spec doc, and — for auto-discovery only — a
new subprocess shell-out to the `gh` CLI. No API, DB, or schema
surface is involved; this is a read-path task.

This is an **N ≥ 2 task plan**: this doc is the orchestrating
task plan; phase content lives in per-phase plan files. The
deliberation, rejected alternatives, and reality-check inputs
live in the sibling scoping doc
([`scoping/m1-t4-expanded-per-node-display.md`](scoping/m1-t4-expanded-per-node-display.md)),
which this plan does not restate. This task plan's `Status` is
`Proposed`: its orchestration contracts (Phase Contracts,
Cross-Phase Decisions, Cross-Cutting Invariants, sequencing) are
decision-complete and the `In draft → Proposed` promotion-gate
self-review per
[`task-plan.md`](../../../spec/planning/task-plan.md) has been
run. P2's *phase plan* is drafted just-in-time after P1's
implementing PR merges (scoping D2/D5) — that is a legitimate
future doc under the just-in-time rule, **not** an unsettled
decision in this task plan, so it does not hold the task plan at
`In draft`. The task plan flips `Proposed → In progress` when
P1's implementing PR merges and `→ Landed` with P2's (last)
implementing PR (see "Terminal state" below).

## Goal

Each plan-tree node renders, beneath its existing badge/label/
markers line, its long description and a list of related PRs.
Related PRs are sourced from an optional `related_prs` frontmatter
field (always authoritative) and, additively, from `gh pr list`
auto-discovery keyed on the slug. The detail renders inline as
static HTML (no JavaScript, no panel, no expand/collapse). Nodes
without a long description or related PRs render exactly as before
plus an unchanged badge/label line — the addition is purely
additive and degrades gracefully when `gh` is unavailable.

## Phases

t4 ships in two phases (scoping D2). Phase plans are separate
files per
[`task-plan.md`](../../../spec/planning/task-plan.md) path
conventions.

- **P1 — Inline per-node detail render.**
  [`m1-t4-p1-inline-detail-render.md`](m1-t4-p1-inline-detail-render.md)
  (Status `Proposed`). Renders `LongDescription` inline; adds
  the optional `related_prs` frontmatter field, parses it, and
  renders author-curated PRs inline. Pure read-path, no external
  dependency. Independently shippable: the page surfaces
  descriptions and manually-listed PRs.
- **P2 — `gh pr list` auto-discovery.**
  `m1-t4-p2-gh-discovery.md` (not yet drafted — drafted
  just-in-time after P1's PR merges, per scoping D2/D5). Adds
  the codebase's first subprocess shell-out: `gh pr list` keyed
  on the slug, merged/deduped into P1's related-PR set,
  best-effort with graceful degradation. Carries a
  novel-mechanism spike (scoping D5) and a Validation Gate that
  exercises the `gh`-unavailable failure matrix against real
  environments.

P1 → P2 is a sequence: P2 augments P1's already-rendered PR
surface and ships no artifact without it (scoping D2).

## Phase Contracts

Per-phase **WHAT** contracts. The **HOW** (file inventory,
signatures, commands, validation gate) lives in each phase plan.
Required for an N ≥ 2 task plan per
[`shared.md`](../../../spec/planning/shared.md) "Parent-doc child
contracts."

### P1 — Inline per-node detail render

- **End result.** Every node renders its long description and
  its author-curated related PRs inline beneath the existing
  badge/label/markers line, as static HTML. A node with neither
  renders an unchanged single line. The full long description is
  rendered (no truncation, no interaction — scoping D1's
  bans-on-surface consequence).
- **Interfaces.** Consumes t3's `PlanNode.LongDescription`
  (Landed). Introduces the optional `related_prs` frontmatter
  field — a sequence of **absolute-URL strings** — and a
  `parsedDoc`/`PlanNode` `RelatedPRs` carry that P2 reads and
  augments. The absolute-URL entry shape is the P1↔P2 contract:
  P2's `gh`-discovered PRs are already absolute URLs, so both
  phases share one shape and P1 carries no canonicalization.
  Documents `related_prs` in `spec/planning/shared.md` as
  optional/additive.
- **Preserves.** Existing badge, label, slug tooltip,
  work-instance markers, child nesting, and the empty-state path
  are unchanged. A doc without `related_prs` or a body renders
  with no warning, error, or skip.

### P2 — `gh pr list` auto-discovery

- **End result.** A node's related-PR list additionally includes
  PRs discovered by `gh pr list` keyed on the slug, merged and
  deduped with the frontmatter set. When `gh` is unavailable for
  any reason, the node still renders its frontmatter PRs (or
  none) and the page never fails.
- **Interfaces.** Reads and extends P1's `RelatedPRs` carry.
  Establishes the canonical PR-identity dedupe key (scoping D3's
  deferred half) once the D5 spike fixes the `gh --json` shape.
  No new frontmatter or spec field.
- **Preserves.** P1's frontmatter PRs stay authoritative and are
  always rendered. The render path stays walk-on-every-request
  with no caching, file-watch, or in-memory build-up introduced
  (m1 Cross-Task Invariant).

## Cross-Phase Decisions

Decisions that thread both phases, owned here per
[`task-plan.md`](../../../spec/planning/task-plan.md)
"Cross-PR coordination" (the task plan coordinates phases; phase
plans do not pre-lock cross-phase contracts). Deliberation and
rejected alternatives are in the scoping doc and not restated.

- **Surface shape is inline static HTML, no JavaScript (scoping
  D1).** Binds both phases: P2's discovered PRs render through
  the same inline `node`-template block P1 establishes. Resolves
  the milestone's deferred "Long-description rendering location
  (t4)."
- **`related_prs` is one optional, additive frontmatter field
  of absolute-URL strings; its spec change ships in P1's PR
  (scoping D3).** P2 adds no new field and inherits the
  absolute-URL entry shape (no PR-syntax canonicalization in
  either phase; `gh`'s `url` output is already this shape).
- **Frontmatter PRs are authoritative; `gh` discovery is
  best-effort and additive (scoping D4).** The P1↔P2 boundary:
  P1 owns the frontmatter source and the render block; P2 owns
  the `gh` source and the merge/dedupe into P1's set by a
  canonical key. P2's canonical-key spelling is a P2-plan open
  input (scoping D3 deferred half / D5 spike), recorded as P2's
  named handoff rather than pre-locked here.
- **P2's subprocess is a novel mechanism requiring a
  just-in-time spike at P2 drafting (scoping D5).** The spike
  branch is `spike/m1-t4-gh-prlist`, never merged into the
  implementation PR.

## Cross-Cutting Invariants

Rules ≥ 2 sites must agree on, threading the phases.

- **Optional frontmatter degrades silently.** `related_prs`
  absence is a valid common state: the parser carries an empty
  list (no warning/skip), the render block emits nothing for an
  empty list, and a field-less node's badge/label line is
  byte-unchanged. P1 names this once so its self-review walks
  parser + template together.
- **Render path stays walk-on-every-request.** No caching, file
  watching, or in-memory build-up is introduced in either phase;
  P2's `gh` call runs inside the existing per-request walk.
  `Verified by:`
  [`Server.index` in site.go](../../../internal/site/site.go)
  calls `walkPlans` then `buildTree` per HTTP handler
  invocation; [`m1-v0-2.md`](m1-v0-2.md) Cross-Task Invariant
  "Render path stays walk-on-every-request."
- **Spec change is additive.** The `related_prs` addition
  introduces no breaking change for vendored spec consumers;
  existing frontmatter without it stays valid (epic cross-cutting
  invariant "Spec contract is additive across this epic").
- **Bare-bones render, no JavaScript.** Neither phase introduces
  client-side script or a non-`/` route; detail is static
  server-rendered HTML. `Verified by:`
  [`render.go`](../../../internal/site/render.go) is one no-JS
  `html/template`; [`design/v0.1-design.md`](../../../design/v0.1-design.md)
  §7.

## Section variance

Per [`shared.md`](../../../spec/planning/shared.md) "Section
variance disclosure": as an N ≥ 2 orchestrating task plan this
doc skips the task-plan-required inline *Contracts (full final
shape)*, *Files to touch*, and a concrete *Validation Gate* —
those genuinely don't apply at the orchestration layer and are
delegated to the phase plans, with the N ≥ 2-required *Phase
Contracts* section substituting for inline Contracts per
[`task-plan.md`](../../../spec/planning/task-plan.md) "Required
and optional sections." The PR introducing this doc repeats
this disclosure in its body's Documentation section.

## Validation Gate

Per-phase Validation Gates live in each phase plan (the gates
differ materially — P1 is the Go toolchain plus a manual render
check; P2 adds the `gh`-unavailable failure matrix against real
environments). The task plan's terminal gate is satisfied when
both phase plans reach `Landed` per their own gates and the
parent milestone t4 row is closed (see Terminal state).

## Terminal state

Per [`task-plan.md`](../../../spec/planning/task-plan.md) "Task
plan terminal state when N ≥ 2": this task plan is `Proposed`
once its orchestration contracts lock (done in the drafting
PR). It flips `Proposed → In progress` when P1's implementing
PR merges (P1's phase plan flips `Landed` in that same PR), and
`In progress → Landed` with **P2's** implementing PR (the last
phase). An undrafted P2 *phase plan* is a just-in-time future
doc, not an unsettled input — it does not gate this task plan's
`Proposed`. P2's implementing PR also performs the parent
milestone [`m1-v0-2.md`](m1-v0-2.md) t4-row close-out, and since
t4 is m1's last task the milestone is then itself terminal,
which that m1-terminal PR additionally handles (sibling
scoping-doc batch deletion, milestone Status).

## Documentation currency

- [`m1-v0-2.md`](m1-v0-2.md) — the parent milestone's Task
  Status t4 row mirrors this task plan's Status: set to
  `Proposed` (linked to this plan; the `—` legend means "not
  drafted," which no longer holds). Its deferred "Long-
  description rendering location (t4)" decision is resolved to
  scoping D1, and a "t4 phase structure (resolved at t4
  drafting)" note is added (N ≥ 2 per scoping D2), mirroring the
  existing t3 note. Subsequent row values track the task plan's
  lifecycle (`Proposed → In progress → Landed`) and land with
  the PR that performs each flip.
- `spec/planning/shared.md` — the additive `related_prs` field
  doc lands in P1's implementing PR (scoping D3); owned by the
  P1 phase plan's Documentation currency.
- [`design/v0.1-design.md`](../../../design/v0.1-design.md) §7 —
  "What the Website Renders" currently says the node shows only
  badge + label + marker. P1's PR updates §7 to reflect inline
  per-node detail; P2's PR notes the `gh` auto-discovery source.
  Each phase plan owns its own §7 edit.

## Backlog Impact

None. No backlog entry graduates, is deleted, split, or shifts.
The [`repo-rooted-doc-links`](../../backlog.md#repo-rooted-doc-links)
entry is post-1.0 and untouched by t4; the governance entries
(`plan-doc-child-contracts`, `stub-children-on-parent-promotion`,
`promotion-gate-explicit-checklist`) are unrelated to t4's
read-path surface.

## Related Docs

- [`m1-v0-2.md`](m1-v0-2.md) — parent milestone; t4 task
  contract and the deferred decision this task resolves.
- [`README.md`](README.md) — parent epic.
- [`scoping/m1-t4-expanded-per-node-display.md`](scoping/m1-t4-expanded-per-node-display.md)
  — sibling scoping doc (deliberation, transient).
- [`m1-t4-p1-inline-detail-render.md`](m1-t4-p1-inline-detail-render.md)
  — P1 phase plan.
- [`m1-t3-descriptive-labels.md`](m1-t3-descriptive-labels.md)
  — t3 (Landed); supplies the `LongDescription` carry P1
  consumes.
- [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md)
  — the rules this plan is structured against.
- [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md)
  — cross-level planning rules; the `related_prs` spec edit
  lands in P1.
