---
slug: tool-originated-task-sessions-m2-t1
Status: In draft
short_description: Mode-affordance render on plan-tree nodes
---

# m2 t1 — Mode-Affordance Render on Plan-Tree Nodes

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
Cross-Task Invariants, Cross-Task Decisions D1/D3/D4); not
loosened here.

**End result.** The plan-tree forest renders a
**mode-affordance** on each node — either **Begin planning**,
**Begin implementation**, or neither — driven by the locked
static map (D3) over the node's already-rendered facts
(node-type + Status + has-children). The affordance is a
**plain HTML `<form method="POST" action="/spawn">`** carrying
the node's slug and the chosen mode in hidden inputs and a
submit button labeled per the mode; clicking it POSTs to t2's
spawn endpoint (D4). A node that offers no mode renders no
form.

**Preserves.** Every node still renders its existing label,
Status badge, work-instance markers, progress-cell row, long
description, and related-PR list — none of those are modified.
The forest's expand/collapse default, the **no-JavaScript**
posture (the page stays pure server-rendered HTML+CSS with
native `<details>`/`<summary>`; a `<form>` submit needs no
script), and the **walk-on-every-request** invariant for the
GET render are unchanged. No new DB read, no new walk.

**Sibling interface.** Produces the **affordance form surface**
t2's `/spawn` endpoint consumes. The render side and the launch
side are coupled by D1 (slug-via-`WST_SLUG`) and D4 (the form's
`action`, method, and field names); both are locked at the
milestone level.

**Product acceptance.** A product reviewer opens the page and
observes that every leaf-shape task and phase node renders the
correct affordance per the locked mode map (Begin planning on
`In draft` / no-doc; Begin implementation on `Proposed`;
neither on `In progress`/`Validating`/`Landed`/`Deferred`);
parent-shape nodes (epic/milestone, and tasks with phase
children) render no affordance; the affordance is visible
without expanding the node.

## What this skeleton does NOT cover

HOW. File inventory, function/template signatures, the exact
shape of the form's HTML, validation-gate specifics, and the
risk register are out of scope for this skeleton — those are
produced by t1's own `In draft` → `Proposed` drafting session
against then-merged code, per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns" and
[`milestone.md`](../../../../spec/planning/milestone.md)
"Anti-goal: WHAT-contract each task, do not HOW-scope any task."

## Related Docs

- [`./README.md`](./README.md) — m2 milestone doc; the locked
  WHAT, Cross-Task Invariants/Decisions, Documentation
  Currency.
- [`../README.md`](../README.md) — parent epic; Cross-Cutting
  Invariants and the resolved vision inputs.
