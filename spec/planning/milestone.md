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
  milestone set. Path: `docs/plans/<epic-slug>/m<N>/README.md` —
  the milestone doc is the `README.md` inside the milestone's own
  `m<N>/` folder; that milestone's task and phase docs are
  siblings in the same folder. The slug in frontmatter is the
  identity; the filename and folder path are for browsing.
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
- **A multi-leaf graph needs a terminal convergence node.** When
  the `flowchart LR` ends in more than one terminal leaf (≥ 2
  tasks with no outgoing edge), there is no graph-determined
  "last" task, so the milestone-terminal close-out — the batch
  scoping-doc deletion, the milestone doc's own
  `Proposed → Landed` flip, and the parent epic's milestone-row
  advance — has no unambiguous owner and silently rides an
  arbitrary leaf's PR. Such a graph must add one terminal
  **validate-and-land** node that depends on every leaf; that
  node owns the milestone-terminal close-out and lands as a
  dedicated close-out PR — the milestone-graph generalization of
  the [`task-plan.md`](./task-plan.md) "Parallel implementing
  PRs" exception, which binds only a single task plan's N ≥ 2 PR
  set and so does not reach across a milestone's task graph. A
  graph that ends in exactly one leaf — or whose prose names a
  clearly-last task — needs no convergence node: the close-out
  rides with that last PR per the terminal-PR conventions above.
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
  **The same terminal/cleanup PR that deletes the scoping docs
  must also de-link every inbound reference to them in the same
  change.** Durable milestone and task/phase plans routinely cite
  scoping by live markdown link (a decision's `Verified by:`, a
  Related Docs entry); those links dangle the instant the targets
  are deleted, so a close-out that deletes scoping without the
  de-link sweep ships broken links and is reviewer-flag. The
  canonical de-link form is a **non-link inline-code reference**
  (`` `scoping/<name>.md` ``): the decision identifiers already
  carried in the citing prose ("scoping D1", "S4", "SD2")
  preserve the provenance, and the deleted doc remains in git
  history — so no content is lost by dropping the link. This is a
  consequence of the scoping-is-transient / plan-owns-the-durable-
  record split (see [`task-plan.md`](./task-plan.md) "Scoping owns
  / plan owns"): a durable plan should not hold a live link to a
  doc the spec guarantees will be deleted. The milestone doc may
  override the batch-deletion rule for an
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

## Milestone retrospective

A milestone is a checkpoint: its purpose is not only to ship its
task set but to **stop, evaluate the shipped result against real
use and product feedback, and capture the adjustments that
evaluation surfaces** before the next vision is taken up. The
milestone-planning session at the start locks the task set; the
retrospective at the end is the symmetric bookend where that set
is judged against reality. The "decisions made at the m1
retrospective" that later milestone docs cite are this step —
previously unmodeled, now defined here.

The retrospective is **distinct from the milestone-terminal
close-out**. The terminal close-out (owned by the last or
convergence PR per the convergence-node rule above) is
mechanical: batch scoping-doc deletion, the milestone doc's
`Proposed → Landed` flip, the parent epic's milestone-row
advance. The retrospective is evaluative, runs at the same
boundary, and does **not** gate that flip — a milestone reaches
`Landed` when its tasks land and validate; gating it on
open-ended product feedback would mean milestones never close.

- **Output channel is the backlog, not this milestone.** Every
  adjustment the retrospective surfaces is captured as a
  [`backlog.md`](../backlog.md) entry (the existing
  "captured but not yet planned" surface), never as a new task
  retrofitted into the just-locked, now-`Landed` milestone and
  never by reopening its child set. The milestone-planning
  session ran once at the start and is over; the child set is
  immutable after it locks.
- **Acting on adjustments is the next cycle's scope.** An entry
  the team chooses to act on graduates per
  [`backlog.md`](../backlog.md) "Entry lifecycle" — into the
  **next milestone's** drafting session (which is thereby
  *informed by* the retrospective, the concrete "adjust before
  the next vision" seam) or into a standalone graduated task
  plan when the adjustment is independent of any milestone.
  Graduation is the seam; the retrospective itself creates no
  plan-tree node and reopens none.
- **No retrospective doc artifact.** The retrospective is a
  step, not a deliverable doc. Its durable trace is the backlog
  entries it produces — zero or more; a retrospective that
  surfaces nothing actionable is a valid outcome — plus any
  decision it makes that a later doc cites. It adds no section
  to the milestone doc.

## Product acceptance and per-leaf validation

A milestone-planning session decides not only what each task
must technically do (its WHAT contract in Task Contracts) but
**what a product reviewer must be able to do or see when the
task lands** — the task's *product acceptance*. This is a
demonstrable, product-facing statement ("a reviewer can open
the page and collapse any node," not "the `node` template emits
a `<details>`"), recorded alongside each task's contract in Task
Contracts. Milestone planning owns the *what to demo*; the task
implementation owns the *how to demo it* — the step-by-step
walkthrough, per [`task-plan.md`](./task-plan.md) "Product-facing
leaf tasks output a demo walkthrough."

Not every task has a product surface. A purely internal task —
a refactor, a tooling or build change, anything a product
reviewer cannot meaningfully *do or see* — gets **no**
product-acceptance contract and is validated by its technical
gate alone. This is the same `no-product-surface work` carve-out
[`shared.md`](./shared.md) "Plan-doc Status" names when scoping
mandatory `Validating`; the two specs are deliberately aligned
so a planner cannot reach opposite conclusions from them.

**Every product-facing Mermaid-graph leaf is product-validated
before the milestone closes.** Such a leaf (a task with no
outgoing edge that the milestone gave a product-acceptance
contract) must not flip `Landed` on merged code alone — its
product acceptance must be demonstrated and approval recorded. A
leaf with **no** product surface carries no product-acceptance
contract, closes on its technical gate alone, and is exempt from
this section. For each product-facing leaf, milestone planning
picks one of two shapes and records the choice in Sequencing:

- **In-task validation box.** The leaf task carries the
  validation in its own lifecycle: its implementing PR merges at
  Status `Validating`, the product-acceptance walkthrough is
  run, approval is recorded in the plan, then it flips `Landed`
  (per [`shared.md`](./shared.md) "Plan-doc Status" — `Validating`
  is mandatory, not skippable, for such a leaf).
- **Dedicated validation task node.** A separate validation task
  node is added that depends on the product-surface task.
  Adding it **changes the graph topology**: the validation node
  now has no outgoing edge, so **it** is the Mermaid-graph leaf,
  and the product-surface task it depends on becomes an
  **interior node** (it has an outgoing edge to the validation
  node). The validation node owns the product validation and, as
  the leaf, carries the mandatory `Validating` per
  [`shared.md`](./shared.md) "Plan-doc Status"; the interior task
  it depends on is not a leaf, so that leaf-keyed rule does not
  bind it and it closes on its own technical gate. The
  product-acceptance obligation always lives on whichever node is
  the leaf — never on two nodes, never in conflict. Use this
  shape when the demo spans siblings or needs a deliberate
  reviewer hand-off rather than riding the task's own PR.

This **composes with, and does not replace, the terminal
convergence node** (see "A multi-leaf graph needs a terminal
convergence node" above). The convergence node owns the
*mechanical* milestone-terminal close-out (batch scoping
deletion, the milestone `Proposed → Landed` flip, the parent
epic's milestone-row advance); per-leaf product validation is
the *product approval* and is a separate obligation. A
convergence node may also serve as the dedicated validation
node when planning chooses the dedicated-node shape for the
leaves it converges, but product validation is never silently
folded into the mechanical close-out — if the convergence node
carries it, Sequencing says so explicitly.

This binds milestone-planning sessions from this point forward;
already-`Landed` milestone docs are immutable history (per the
retrospective immutability framing above) and are not
retrofitted.

## Product-validation findings: fix now or defer

When a product-acceptance walkthrough surfaces an issue, the
milestone faces a triage with exactly two outcomes. There is no
third "note it and move on" — an open finding is either fixed
before the leaf closes or explicitly accepted and routed
forward.

- **Fix in the current milestone.** The gap blocks acceptable
  shipping. The leaf has not met its product acceptance, so it
  **stays `Validating`** (never flips `Landed` on the unmet
  contract); the fix lands — a follow-up PR on the leaf, or on
  its dedicated validation node — the walkthrough is re-run,
  approval is recorded, then `Landed`. The milestone does not
  close on that leaf until then.
- **Defer to a later milestone.** The gap is acceptable to ship
  with. In the **deferring leaf's own PR** (the same PR that
  flips it `Landed` — not deferred to milestone end): the leaf's
  product-acceptance contract is **explicitly narrowed to carve
  out the finding** (recorded in Task Contracts — the leaf then
  genuinely meets its *as-narrowed* acceptance and may flip
  `Landed`; it is never `Landed` against an unmet contract), and
  the carved-out issue is captured **immediately** as a
  [`backlog.md`](../backlog.md) entry. The **milestone
  retrospective** (see "Milestone retrospective" above) then
  *reviews and routes* the milestone's accumulated deferrals
  forward — graduating them into the next milestone or a
  standalone task per that mechanism. The retrospective is the
  forward-routing seam, **not** the capture gate: leaf-level
  `Landed` never waits on the milestone-end retrospective, which
  is non-gating by construction.

**The deciding factor — "is shipping with this gap acceptable?"
— is a subjective product call the spec does not adjudicate.**
There is no threshold, rubric, or severity scale here on
purpose: a false-precision rule would relax under pressure
exactly when it matters. What the spec requires is that the
call be **made explicitly and recorded** — which path was
taken, a one-line rationale, and for a deferral the contract
carve-out plus the backlog slug. This is the project's
visibility-over-determinism tenet applied to acceptance: an
accepted, visibly-recorded carve-out beats both an unrecorded
judgment and a rubric pretending the call is objective.

Deferral is not a way to skip the gate: the walkthrough is
still run and the finding still recorded — "defer" records an
accepted, visible gap, it does not mean validation didn't
happen. A leaf flipped `Landed` with an open product-validation
finding that has neither a recorded fix nor a recorded
deferral (carve-out + backlog entry) is Status drift and is
reviewer-flag.

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
