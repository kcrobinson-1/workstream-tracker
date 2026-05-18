---
slug: workstream-tracker-1-0-m2-t2
Status: In draft
short_description: Nested collapsible per-node Status boxes inside the forest region
---

# Task 2 (stub) — Expanded in-root nested-box render

> **Parent-promotion stub.** Seeded from
> [`m2-expanded-view-and-roster.md`](m2-expanded-view-and-roster.md)
> "Task Contracts" in the PR that flipped m2 to `Proposed`. Per
> [`shared.md`](../../../spec/planning/shared.md) "Parent-doc
> child contracts," this stub is exempt from the
> [`task-plan.md`](../../../spec/planning/task-plan.md) "Required
> and optional sections" rule and from "Plans describe contracts,
> not implementation" until this task's own `In draft` →
> `Proposed` drafting session runs. It carries only the locked
> inherited WHAT contract below; HOW is scoped at drafting.

## Inherited contract (locked at m2 milestone planning)

**End result + preserves.** Within t1's forest region, a root's
descendants render as nested epic → milestone → task → phase
boxes, each carrying its own Status badge, each independently
collapsible — replacing the flat nested-`<ul>` bullet render.
Owns the stacked-row narrow-window degrade (it introduces the
right-aligned Status/progress that can collide with left-aligned
label text). Preserves: v0.1's actor tags on nodes still render
on each box (the deferred actor-icons-on-progress-boxes work
assumes node actor tags remain — they must not regress);
walk-on-every-request unchanged; a stub and any node without
rich data still render.

**Feeds siblings.** Builds into t1's forest region. Provides the
per-node expanded box surface t3 renders progress boxes
*inside*. Independent of t4.

**Open HOW (decide at this task's drafting, per m2 Cross-Task
Decisions):** the collapsible mechanism vs. the v0.2 no-JS
posture must be decided explicitly, not regressed silently.
