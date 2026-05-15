# Milestone Planning Sessions

Per-level planning playbook for **milestone-planning** sessions. Loads
[`shared.md`](./shared.md) for cross-level planning rules. This file
covers what is unique to the milestone level.

A milestone planning session establishes durable cross-phase
coordination for a multi-phase milestone. Run this session once at
the start of a milestone, before any per-phase planning.

- **Goal.** Produce the milestone doc: restated milestone goal, phase
  sequencing with dependency rationale, cross-phase invariants that
  thread multiple phases, cross-phase decisions that lock contracts
  between phases, milestone-level risks, doc-currency map across the
  milestone set. Path: `docs/plans/<epic-slug>/m<N>-<short-descriptor>.md`,
  where `<short-descriptor>` is an optional human-readable suffix. The
  slug in frontmatter is the identity; the filename is for browsing.
  See [`planning-doc-location.md`](../planning-doc-location.md) for
  the full layout convention.
- **Phase dependency graph.** The milestone doc's "Sequencing"
  section opens with a Mermaid `flowchart LR` block: each phase
  is a node, "blocks" relationships are arrows (an `A --> B`
  edge means A blocks B / B depends on A), the upstream
  milestone (e.g. M1 for an M2 doc) appears as a
  dependency-only node so prerequisites are explicit. Phase
  numbering reflects intended ship order, **not** strict
  dependency — readers default-assume `N.k` depends on
  `N.(k-1)`, which silently wastes time when phases are
  independent and could draft or ship in parallel. The graph
  makes parallelism visible at a glance instead of buried in
  prose; the prose still carries rationale (which phase ships
  first and why, terminal-PR conventions, cross-phase coupling
  beyond hard dependencies).
- **Anti-goal: do not scope any phase in this session.** Phase
  scoping and plan-drafting (the per-phase deliberation,
  contracts, file inventory, risks, and execution steps — split
  between the scoping doc and the plan doc per
  [`task-plan.md`](./task-plan.md) "Scoping owns / plan owns") belong to
  the phase planning session for each phase, run per the timing
  and pending-input rules in [`task-plan.md`](./task-plan.md) (which
  permit drafting in parallel with the prior phase's
  implementation or review under explicit citation requirements,
  but still bar scoping during the milestone session itself).
  Scoping any phase in the milestone session — even the first —
  risks recording assumptions that won't survive contact with
  merged code, and produces confident-feeling artifacts that may
  or may not be grounded. When phase A's scoping cites phase B's
  "Inputs From Siblings" section, both docs feel verified;
  neither is. Grounding lives in the per-cross-phase-decision
  verification rule below ("read the actual code that would be
  affected by each option").
- **Output set.** Milestone doc — durable; survives all phase
  work. Single output of this session. Phase scoping docs are
  produced by their respective phase planning sessions, not
  here; they delete in batch when the milestone's full set of
  phase plans exists (not as each phase plan lands), as part of
  the milestone's terminal PR or a focused cleanup PR. The
  reason: sibling scoping docs reference each other, so deleting
  one early creates link rot elsewhere. The milestone doc may
  override the batch-deletion rule for an unusual lifecycle, but
  should record the override explicitly. Cross-phase decision
  record lives inside the milestone doc, not as a separate file.
- **Cap.** Stop when iteration without ending hits — repeated
  rewrites of the same section, cross-phase decisions that
  re-open after being marked resolved, or new docs spawning
  without resolving existing ones. That iteration signal is the
  diminishing-returns indicator; remaining value comes from
  doing the work, not from more planning content. Planning time
  should be a fraction of implementation time, not a parallel
  effort.
- **Verify before recording any cross-phase decision.** For each
  cross-phase decision, read the actual code that would be
  affected by each option, not summaries from a research
  subagent. A decision recorded with options/pros/cons but
  without code-grounded option generation is a guess dressed as
  rigor — the option set itself can be wrong if the underlying
  mental model is wrong.
- **Defer rather than over-resolve.** If a cross-phase decision
  can be made later by the affected phase's planner without
  blocking earlier phases, mark it deferred with a clear "decide
  when phase N drafts" note. Premature resolution of deferrable
  decisions is a major source of wrong premises that propagate
  through the doc set.
- **PR-count predictions are not contracts.** Per-phase PR
  counts named in the milestone doc are estimates. The phase
  planning session re-derives the actual PR count using the
  rule in [`task-plan.md`](./task-plan.md) "PR-count predictions need a
  branch test"; splitting a phase into sub-phases at plan time
  is normal, not a process failure.

## Required and optional sections

A milestone doc carries the following sections.

**Required:**

- Status
- Goal
- Phase Status
- Sequencing (Mermaid `flowchart LR` per the rule above)
- Cross-Phase Invariants
- Cross-Phase Decisions
- Cross-Phase Risks
- Documentation Currency
- Backlog Impact (per [`backlog.md`](../backlog.md))
- Related Docs

**Optional, when content applies:**

- Pending Inputs From `<prior milestone>` — when drafting
  depends on outputs of a prior milestone not yet merged

Variance from this list — an unlisted section is appropriate,
or a required section genuinely doesn't apply — follows the
[`shared.md`](./shared.md) "Section variance disclosure" rule.
