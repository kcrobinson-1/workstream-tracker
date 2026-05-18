---
slug: workstream-tracker-1-0-m2-t4
Status: In draft
short_description: Session roster (bound + unbound) with deliberately-unstructured JSON detail
---

# Task 4 (stub) — Session roster + work-item enrichment

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

**End result + preserves.** Replacing t1's roster-region
placeholder, a roster lists every registered active session —
bound (slug matches a plan-tree node) and unbound (slug absent
from the tree, the accepted orphan/typoed-slug residual) — that
v0.1's render currently drops. The register client/CLI sends
richer per-session data; roster entries are expandable, showing
session name, reported PRs, and other reported fields rendered
as a **deliberately unstructured raw-JSON view** (schema-loose
on purpose; "what is required" is intentionally deferred).
Preserves: registration stays observable best-effort and opt-in
— the roster surfaces sessions that chose to register and never
claims to see all of them (the existing best-effort/observable
registration tenet); walk-on-every-request unchanged; richer
data is additive (a session reporting only the v0.2 minimum
still lists).

**Feeds siblings.** Builds into t1's roster region (replaces its
placeholder); does not touch the forest region. Independent of
t2/t3; ships in parallel with the t2 → t3 chain. Delivers the
"every session can be accounted for" observability surface the
[`deterministic-interactive-registration`](../../backlog.md#deterministic-interactive-registration)
backlog entry's mitigation direction names.

**Open HOW (decide at this task's drafting, per m2 Cross-Task
Decisions):** whether enriched data is read by joining the event
log or via an additive `work_instances` column (bounded by the
additive-schema invariant); the precise backlog effect on
`deterministic-interactive-registration` (shift vs.
reference-only). The schema-loose JSON posture must NOT be
homogenized toward t3's declared-stages schema.
