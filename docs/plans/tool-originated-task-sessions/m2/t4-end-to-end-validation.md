---
slug: tool-originated-task-sessions-m2-t4
Status: In draft
short_description: End-to-end product validation + milestone-terminal close-out
---

# m2 t4 — End-to-End Product Validation + Milestone-Terminal Close-Out

> **Parent-promotion skeleton — scope locked, HOW not yet
> drafted.** Seeded by the m2 milestone-planning session's
> `In draft` → `Proposed` promotion gate per
> [`shared.md`](../../../../spec/planning/shared.md) "Parent-doc
> child contracts → Parent-promotion stub seeding." The locked
> WHAT contract this skeleton carries comes from the [m2
> milestone doc](./README.md) Task Contracts row. Per the
> seeding rule the skeleton is **exempt from the
> [`task-plan.md`](../../../../spec/planning/task-plan.md)
> "Required and optional sections" rule and from
> [`shared.md`](../../../../spec/planning/shared.md) "Plans
> describe contracts, not implementation"** until its own
> drafting session opens — at which point that exemption ends
> and the doc grows into a full task plan.

## Inherited Contract (locked by the m2 milestone doc)

Binding input from [`./README.md`](./README.md) (Task Contracts,
Documentation Currency); not loosened here.

**End result.** A product reviewer performs the **full
end-to-end walkthrough** against the local server with a real
Claude Code launcher available: open the page, locate a
`Proposed` task/phase node, click **Begin implementation**,
observe a real Claude Code session start, observe the **real
register receipt** echo with the clicked node's canonical slug
(per the deterministic-path handshake rule), observe the
work-instance appear on the clicked node in the rendered
tree, observe the session run its initial prompt; symmetrically
for **Begin planning** against an `In draft`/no-doc node;
symmetrically observe that nodes offering no mode (epic /
milestone parents, in-flight task/phase nodes) show no
affordance. The walkthrough records approval in this plan;
product-validation findings are routed per
[`milestone.md`](../../../../spec/planning/milestone.md)
"Product-validation findings: fix now or defer."

The terminal PR then performs the **milestone-terminal
close-out**: batch-deletes the m2 `scoping/` subfolder (every
transient scoping doc t1's, t2's, and t3's planning sessions
produced), de-links any inbound references to those scoping
docs in any durable plan doc that survives the batch (the
non-link inline-code form), flips this milestone doc's
`Status` `Proposed` → `Landed`, advances the parent epic's m2
milestone row to `Landed` with its terminal PR link.

**Preserves.** Nothing about the production code path is
touched at this stage; t4 is a validation + close-out node,
not a code-producing task. The milestone retrospective (per
[`milestone.md`](../../../../spec/planning/milestone.md)
"Milestone retrospective") runs at the same boundary and routes
any accumulated findings forward through the backlog — it does
**not** gate t4's `Landed` flip.

**Sibling interface.** Sole Mermaid-graph leaf — converges t1,
t2, and t3 and is the milestone-terminal node.

**Product acceptance.** The full walkthrough above runs to
completion against a live local server + Claude Code launcher,
with the real register receipt observed for both modes and
parent-shape nodes verified to offer no affordance; approval
is recorded in this plan before its `Landed` flip.

## What this skeleton does NOT cover

HOW. The exact walkthrough script, the form of the recorded
approval, the precise close-out commit ordering and per-file
diff inventory, and the risk register — produced by t4's own
`In draft` → `Proposed` drafting session against then-merged
code, per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns" and
[`milestone.md`](../../../../spec/planning/milestone.md)
"Anti-goal: WHAT-contract each task, do not HOW-scope any task."

## Related Docs

- [`./README.md`](./README.md) — m2 milestone doc; the locked
  WHAT, Cross-Task Invariants/Decisions, Documentation
  Currency, Backlog Impact.
- [`./t1-mode-affordance-render.md`](./t1-mode-affordance-render.md),
  [`./t2-spawn-endpoint.md`](./t2-spawn-endpoint.md),
  [`./t3-session-start-hook.md`](./t3-session-start-hook.md) —
  the three interior tasks t4 converges and validates.
- [`../README.md`](../README.md) — parent epic; receives the
  milestone-row advance at t4's terminal PR.
