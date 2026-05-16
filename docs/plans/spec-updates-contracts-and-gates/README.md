---
slug: spec-updates-contracts-and-gates
Status: In draft
---

# Spec Updates: Contracts and Gates

> **Starting scaffold, not a gated plan.** Status is `In draft`.
> Full Contracts, Files-to-touch detail, and Validation Gate
> draft just-in-time before implementation, after a scoping doc
> per [task-plan.md](../../../spec/planning/task-plan.md)
> "Scoping precedes plan drafting." This doc records the work's
> shape so it can be tracked now; it deliberately skips the
> task-plan required-sections gate until it graduates to a
> proper plan.

## Context

Dogfooding workstream-tracker's own planning surfaced a cluster
of small, cohesive gaps in the planning and backlog spec — the
artifacts consumer projects vendor. None block 1.0, but each
makes the spec internally inconsistent and was hit directly
while drafting the `workstream-tracker-1-0` epic and its m1
milestone. This task batches the corrections into one
spec-correctness pass rather than spinning up a root per nit.

What it touches conceptually: the planning-spec rule set
(epic / milestone / task-plan / shared) and the backlog format
spec. No product code.

## Goal

Correct the planning and backlog spec so that parent-level plan
docs (epic, milestone) carry required per-child contract
sections and pass an `In draft` → `Proposed` promotion gate,
and backlog entries carry a codified link to their graduated
plan.

## Estimated Phases

Estimates, not commitments — phase count and structure draft at
this task's planning session per
[task-plan.md](../../../spec/planning/task-plan.md)
"Just-in-time scoping and plan drafting."

- **Child contracts.** Cross-level WHAT/HOW-split rule in
  `shared.md`; required-section additions — Milestone Contracts
  (`epic.md`), Task Contracts (`milestone.md`), Phase Contracts
  for N≥2 (`task-plan.md`). Refine the "Anti-goal: do not scope"
  rules to forbid HOW-scoping while requiring WHAT-contracting.
- **Parent-doc promotion gate.** An `In draft` → `Proposed`
  gate for epic and milestone docs (the existing gate in
  `task-plan.md` binds task/phase plans only), bound to the PR
  that locks the doc's scope.
- **Backlog link affordance.** A codified spot in `backlog.md`
  for a navigable link/path to the graduated plan, beyond the
  slug-only `Graduated — <plan-slug>` form. Prior art to
  consult, not match: neighborly-events' backlog `Detail`
  format.
- **Phase-vs-Task naming.** `milestone.md` names the
  milestone's child unit "Phase," but the slug grammar calls it
  a task. Reconcile, or fold into the child-contracts phase.

## Spec Surfaces (estimate)

- [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md)
- [`../../../spec/planning/epic.md`](../../../spec/planning/epic.md)
- [`../../../spec/planning/milestone.md`](../../../spec/planning/milestone.md)
- [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md)
- [`../../../spec/backlog.md`](../../../spec/backlog.md)

## Backlog Impact

Graduated from the `plan-doc-child-contracts` backlog entry,
which flips to `Graduated — spec-updates-contracts-and-gates`
in the same change per
[`../../../spec/backlog.md`](../../../spec/backlog.md) entry
lifecycle. The backlog entry is left carrying only the slug
(no navigable link) on purpose — codifying that link is one of
this task's own phases, and pre-supplying an ad-hoc link would
bias the design decision.

## Deferred to plan-drafting

Full Contracts, Files-to-touch inventory, Validation Gate, and
the scoping doc — drafted just-in-time before implementation
per [task-plan.md](../../../spec/planning/task-plan.md). This
scaffold intentionally does not pre-commit them.

## Related Docs

- [`../../backlog.md`](../../backlog.md) — origin backlog entry.
- [`../workstream-tracker-1-0/README.md`](../workstream-tracker-1-0/README.md)
  — the epic whose drafting surfaced these gaps (the
  variance-disclosed Task/Milestone Contracts sections there
  are the interim workaround this task codifies).
- [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md)
  — the rules a proper plan for this task will follow once it
  graduates from scaffold to gated plan.
