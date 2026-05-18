---
slug: workstream-tracker-1-0-m2-t3
Status: In draft
short_description: Additive spec field + parser + per-node progress-box render driven by the doc
---

# Task 3 (stub) — Doc-declared progress stages (spec-first)

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

**End result + preserves.** An additive `spec/` frontmatter
affordance lets a plan doc declare its own progress stages; the
parser reads it; every node level (root, milestone, task, phase)
renders a row of progress boxes whose count and order come from
the doc. A stub (`slug` + `Status: In draft`, no declared
stages) renders only the Drafting box; the rest appear once the
doc reaches `Proposed`. Preserves: the spec change is optional
and additive (a doc omitting the field renders with no
error/skip, falling back to the Drafting-box-only / bare shape —
the already-supported stub render case
`stub-children-on-parent-promotion` relies on); existing
vendored consumers are unaffected; walk-on-every-request
unchanged.

**Feeds siblings.** Consumes t2's expanded per-node box as the
render host for the box row. Does not feed t4. The
declared-stages spec field is the spec-first deliverable; t4
deliberately does **not** consume or schematize it
(posture-tension invariant — do not homogenize the two postures).

**Open HOW (decide at this task's drafting, per m2 Cross-Task
Decisions):** the frontmatter field name/structure (per-stage
counts vs. ordered stage list) is the spec-change surface
itself; spelling is decided here under the "decompose options
into shapes" and exact-match-label discipline, not invented at
milestone level.
