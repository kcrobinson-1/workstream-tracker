---
slug: workstream-tracker-1-0-m2-t4-p1
Status: In draft
short_description: Bare bound/unbound session roster replacing the t1 placeholder (no event join, no client change)
---

# t4 p1 (stub) — Bare bound/unbound roster

> **Parent-promotion stub.** Seeded from
> [`t4-session-roster.md`](t4-session-roster.md) "Phase Contracts"
> in the PR that flipped the t4 task plan to `Proposed`. Per
> [`shared.md`](../../../../spec/planning/shared.md) "Parent-doc
> child contracts," this stub is exempt from the
> [`task-plan.md`](../../../../spec/planning/task-plan.md) "Required
> and optional sections" rule and from "Plans describe contracts,
> not implementation" until this phase's own `In draft` →
> `Proposed` drafting session runs. It carries only the locked
> inherited WHAT below; HOW is scoped just-in-time against merged
> code at this phase's drafting.

## Inherited contract (locked at t4 task-plan drafting)

**End result + preserves.** The roster region lists every active
work-instance, each classified **bound** (slug matches a walked
plan doc) or **unbound** (slug absent from the walked tree),
replacing t1's deliberate placeholder, using only existing
slug/actor/state data — no event-log join, no register-client
change. Preserves: the forest region body and the shell layout
unchanged; walk-on-every-request (per-request roster load, no
cache); the v0.1 forest per-node actor markers; a session that
reports nothing still lists.

**Feeds siblings.** Establishes the roster loader, the shared
`indexData` roster field, and the bound/unbound classification
that p2 enriches. Does not touch the event log or the register
client/CLI.

**Open HOW (decide at this phase's drafting, per the t4 scoping
doc's Open decisions):** the bound/unbound classification data
source (the walked plan-tree slug set the handler already has vs.
a fresh slug read), bounded by the walk-on-every-request
invariant. See
[`scoping/t4-session-roster.md`](scoping/t4-session-roster.md)
Open decision 4 and the t4 task plan's Cross-Cutting Invariants.
