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
read-experience track (t3 → t4): t3 laid the parsing foundation
specifically so t4 could consume it, and t4's "expanded
per-node display" is the last gap on that track between v0.1's
bare-bones render and v0.2's "scannable per node" goal. t4 is
**not** m1's last task overall — the foundation track's
`…-m1-t2` (Automatic agent registration) is still undrafted, so
m1 is not terminal when t4 completes (see "Terminal state"). Related PRs come from two sources the milestone t4 contract
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
`In progress`: its orchestration contracts (Phase Contracts,
Cross-Phase Decisions, Cross-Cutting Invariants, sequencing)
locked at drafting (promotion-gate self-review per
[`task-plan.md`](../../../spec/planning/task-plan.md) run), it
flipped `Proposed → In progress` when P1's implementing PR
(#15) merged, and it reaches `Landed` with **P2's** (the last
phase's) implementing PR. Both phase plans are now drafted: P1
[`m1-t4-p1-inline-detail-render.md`](m1-t4-p1-inline-detail-render.md)
is `Landed`, P2
[`m1-t4-p2-gh-discovery.md`](m1-t4-p2-gh-discovery.md) is
`Proposed` (drafted just-in-time after P1 merged, per scoping
D2/D5). Note: P2's PR is **t4's task-terminal** PR, **not** the
m1-milestone-terminal PR — m1 has tasks beyond t4 (see
"Terminal state" below).

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
  [`m1-t4-p2-gh-discovery.md`](m1-t4-p2-gh-discovery.md)
  (Status `Proposed`; drafted just-in-time now that P1's PR has
  merged, per scoping D2/D5 — spike run, promotion-gate
  self-review complete). Adds the codebase's first
  subprocess shell-out: one `gh pr list` per request, PRs whose
  title contains a node's slug merged/deduped into P1's
  related-PR set, best-effort with graceful degradation. The
  novel-mechanism spike (scoping D5) has been run — findings and
  the resolved P2 decisions (P2-D1…P2-D4) are in the scoping
  doc; the plan carries a Validation Gate that exercises the
  `gh`-unavailable failure matrix against real environments.
  P2 is t4's last phase: its implementing PR is t4's
  **task-terminal** PR (flips P2, this task plan, and the
  `m1-v0-2.md` t4 row to `Landed`). It is **not** the
  m1-terminal PR — m1 still has `…-m1-t2` undrafted, so the
  sibling scoping-doc batch deletion and milestone
  reconciliation defer to the later m1-terminal PR.

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
  Dedupes by plain absolute-URL string equality (D5 spike
  resolved scoping D3's deferred half — both sources are
  absolute URLs, no canonicalization), frontmatter entries
  first. No new frontmatter or spec field.
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
  the `gh` source and the merge/dedupe into P1's set. The
  dedupe key is **plain absolute-URL string equality**
  (resolved by the D5 spike — scoping P2-D3; both sources are
  absolute URLs, no canonicalization), frontmatter entries
  first.
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
once its orchestration contracts lock (done at drafting). It
flips `Proposed → In progress` when P1's implementing PR merges
(P1's phase plan flips `Landed` in that same PR), and
`In progress → Landed` with **P2's** implementing PR (the last
phase). P2's PR — the **t4 task-terminal** PR — performs only
the t4 close-out: flip P2 `→ Landed`, this task plan
`→ Landed`, and the [`m1-v0-2.md`](m1-v0-2.md) t4 **row**
`→ Landed`.

**P2's PR is NOT the m1-milestone-terminal PR.** m1 has tasks
beyond t4 — its Task Status table still carries
`workstream-tracker-1-0-m1-t2` (Automatic agent registration)
at `—` (undrafted), and t4 is the last task of the
*read-experience track*, not of m1. Per
[`task-plan.md`](../../../spec/planning/task-plan.md) path
conventions, the `scoping/` subfolder's contents delete **in
batch at the milestone-terminal PR** (sibling scoping docs
`m1-t1-*`, `m1-t3-*`, `m1-t4-*` together), and any milestone
Status / Backlog / Documentation-Currency reconciliation
happens there. Those m1-terminal actions are explicitly
**out of scope for P2's PR** and defer to whichever PR lands
m1's last remaining task.

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
