---
slug: workstream-tracker-1-0-m2-t1
Status: In draft
short_description: Two-region page shell — existing forest in the forest region, placeholder in the roster region
---

# Task 1 (stub) — Site skeleton (two-region shell)

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

**End result + preserves.** The page renders as the approved
two-region layout — a forest region (~2/3, left) and a roster
region (~1/3, right), single page scroll, roster placed to be
visible without scrolling when window height allows; a
too-narrow viewport stacks roster below forest (mobile out of
scope). The **existing** plan-tree forest render is placed in
the forest region unchanged; the roster region renders a
deliberate, intentional placeholder (an observed state, not a
blank/broken gap). Preserves: the current forest rendering
behavior verbatim (no node-shape change — that is t2); v0.1's
actor tags on nodes; the walk-on-every-request render path.

**Feeds siblings.** Establishes the two named regions every
later task builds into: t2 enriches the *forest region*'s
internals; t4 replaces the roster region placeholder. No later
task owns or alters the shell or a sibling's region. The
narrow-window *per-node-row* degrade is **not** here (t1 ships
the existing flat render, which has no right-aligned collision)
— it is t2's.

The approved layout reference is
[`../../../design/workstreams-view-m2.svg`](../../../design/workstreams-view-m2.svg).
