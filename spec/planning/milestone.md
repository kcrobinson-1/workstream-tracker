# Milestone Planning Sessions

Per-level planning playbook for **milestone-planning** sessions. Loads
[`shared.md`](./shared.md) for cross-level planning rules. This file
covers what is unique to the milestone level.

A milestone's direct child is a **task**, not a phase — the level
picker is epic → milestone → task → phase, and the slug grammar
encodes the milestone's child with a `t` segment (see
[`shared.md`](./shared.md) "Plan-doc identity (slug)" and
[`task-plan.md`](./task-plan.md) "The level picker"). "Phase" is
the task's child, one level lower; it appears in this file only
where a rule genuinely reaches down to the task's sub-unit.

A milestone planning session establishes durable cross-task
coordination for a multi-task milestone. Run this session once at
the start of a milestone, before any per-task planning.

- **Goal.** Produce the milestone doc: restated milestone goal, task
  sequencing with dependency rationale, the per-task WHAT contract
  (see "Task Contracts" below), cross-task invariants that thread
  multiple tasks, cross-task decisions that lock contracts between
  tasks, milestone-level risks, doc-currency map across the
  milestone set. Path: `docs/plans/<epic-slug>/m<N>-<short-descriptor>.md`,
  where `<short-descriptor>` is an optional human-readable suffix. The
  slug in frontmatter is the identity; the filename is for browsing.
  See [`planning-doc-location.md`](../planning-doc-location.md) for
  the full layout convention.
- **Task dependency graph.** The milestone doc's "Sequencing"
  section opens with a Mermaid `flowchart LR` block: each task
  is a node, "blocks" relationships are arrows (an `A --> B`
  edge means A blocks B / B depends on A), the upstream
  milestone (e.g. M1 for an M2 doc) appears as a
  dependency-only node so prerequisites are explicit. Task
  numbering reflects intended ship order, **not** strict
  dependency — readers default-assume `N.k` depends on
  `N.(k-1)`, which silently wastes time when tasks are
  independent and could draft or ship in parallel. The graph
  makes parallelism visible at a glance instead of buried in
  prose; the prose still carries rationale (which task ships
  first and why, terminal-PR conventions, cross-task coupling
  beyond hard dependencies).
- **Anti-goal: WHAT-contract each task, do not HOW-scope any
  task in this session.** The milestone doc states each task's
  WHAT contract (end result, sibling interfaces, what it
  preserves) in its "Task Contracts" section per
  [`shared.md`](./shared.md) "Parent-doc child contracts" — that
  is required once task scope is locked, not barred. What is
  barred is the task's **HOW**: per-task file inventory,
  contracts-as-implementation, signature shapes, validation-gate
  specifics, risk register, and execution-step ordering. Task
  HOW-scoping and plan-drafting (the per-task deliberation —
  split between the scoping doc and the plan doc per
  [`task-plan.md`](./task-plan.md) "Scoping owns / plan owns")
  belong to the task planning session for each task, run per the
  timing and pending-input rules in
  [`task-plan.md`](./task-plan.md) (which permit drafting in
  parallel with the prior task's implementation or review under
  explicit citation requirements, but still bar HOW-scoping
  during the milestone session itself). HOW-scoping a task in the
  milestone session — even the first — risks recording
  assumptions that won't survive contact with merged code, and
  produces confident-feeling artifacts that may or may not be
  grounded. When task A's scoping cites task B's "Inputs From
  Siblings" section, both docs feel verified; neither is.
  Grounding lives in the per-cross-task-decision verification
  rule below ("read the actual code that would be affected by
  each option").
- **Output set.** Milestone doc — durable; survives all task
  work. Single output of this session. Task and phase scoping
  docs are produced by their respective planning sessions, not
  here; they delete in batch when the milestone's full set of
  task plans (and any phase plans under them) exists (not as
  each plan lands), as part of the milestone's terminal PR or a
  focused cleanup PR. The reason: sibling scoping docs reference
  each other, so deleting one early creates link rot elsewhere.
  The milestone doc may override the batch-deletion rule for an
  unusual lifecycle, but should record the override explicitly.
  Cross-task decision record lives inside the milestone doc, not
  as a separate file.
- **Cap.** Stop when iteration without ending hits — repeated
  rewrites of the same section, cross-task decisions that
  re-open after being marked resolved, or new docs spawning
  without resolving existing ones. That iteration signal is the
  diminishing-returns indicator; remaining value comes from
  doing the work, not from more planning content. Planning time
  should be a fraction of implementation time, not a parallel
  effort.
- **Verify before recording any cross-task decision.** For each
  cross-task decision, read the actual code that would be
  affected by each option, not summaries from a research
  subagent. A decision recorded with options/pros/cons but
  without code-grounded option generation is a guess dressed as
  rigor — the option set itself can be wrong if the underlying
  mental model is wrong.
- **Defer rather than over-resolve.** If a cross-task decision
  can be made later by the affected task's planner without
  blocking earlier tasks, mark it deferred with a clear "decide
  when task N drafts" note. Premature resolution of deferrable
  decisions is a major source of wrong premises that propagate
  through the doc set.
- **PR-count predictions are not contracts.** Per-task PR
  counts named in the milestone doc are estimates. The task
  planning session re-derives the actual PR count using the
  rule in [`task-plan.md`](./task-plan.md) "PR-count predictions
  need a branch test"; splitting a task into phases at plan time
  is normal, not a process failure.

## Required and optional sections

A milestone doc carries the following sections.

**Required:**

- Status
- Goal
- Task Status
- Sequencing (Mermaid `flowchart LR` per the rule above)
- Task Contracts — required once the milestone locks any task's
  scope; per [`shared.md`](./shared.md) "Parent-doc child
  contracts" (the per-task WHAT contract; task HOW stays in each
  task plan), which also governs the per-child table
  presentation and the parent-promotion skeleton seeding driven
  from this section when the milestone passes its promotion gate
- Cross-Task Invariants
- Cross-Task Decisions
- Cross-Task Risks
- Documentation Currency
- Backlog Impact (per [`backlog.md`](../backlog.md))
- Related Docs

**Optional, when content applies:**

- Pending Inputs From `<prior milestone>` — when drafting
  depends on outputs of a prior milestone not yet merged
- Out of Scope — when the milestone records boundary calls as
  final answers (e.g. deferred pieces re-homed to a later
  milestone). Promoted from a recurring disclosed variance: it
  appeared in `m1-v0-2.md` and again in
  `m2-expanded-view-and-roster.md` under the
  [`shared.md`](./shared.md) "Section variance disclosure"
  recurrence rule, so it is listed here rather than
  re-disclosed per doc.

Variance from this list — an unlisted section is appropriate,
or a required section genuinely doesn't apply — follows the
[`shared.md`](./shared.md) "Section variance disclosure" rule.
