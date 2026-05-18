---
slug: tool-originated-task-sessions-m1
Status: In draft
short_description: Unwired deterministic slug-carried registration path
---

# m1 — Unwired Deterministic Slug-Carried Registration Path

> **Parent-promotion skeleton.** Seeded when
> [`tool-originated-task-sessions`](../README.md) passed its
> `In draft` → `Proposed` promotion gate. Until this
> milestone's own drafting session runs, it carries only the
> inherited locked contract below and is exempt from the
> milestone-doc required-sections rule per
> [`shared.md`](../../../../spec/planning/shared.md) "Parent-doc
> child contracts → Parent-promotion stub seeding." The
> drafting session re-derives HOW (task breakdown, validation,
> risks) against merged code; the WHAT below is locked input,
> not to be loosened without reopening the epic.

## Inherited Contract (locked by the epic)

**End result.** Given a construction-known canonical slug, a
work-instance registers deterministically — no natural-language
resolution, no narration handshake. The deterministic path is
proved end-to-end before any tool-acting / agent-spawning UX
exists: the slug is supplied via an explicit slug argument (a
manual / CLI invocation) acting as the stand-in slug producer.

**Preserves.** The existing interactive (contributor-opened)
best-effort grounded narration handshake is unchanged — this is
an *additional* registration path, not a replacement. The
spec's existing exact-slug create-or-attach registration posture
is preserved.

**Sibling interface.** Produces the slug-carried registration
entrypoint that the tool-originated-session UX milestone's spawn
later wires into as the real construction-time slug producer.
m1's manual slug argument is the stand-in for that producer
until the UX milestone replaces it.

**Groundwork this lays.** This milestone is the determinism
proof in isolation: it de-risks the registration contract
(identity-by-construction, no resolution, no handshake) before
the larger observe-only → origination posture shift, so the UX
milestone consumes a proven path rather than co-developing
registration and the posture change at once.

## Open Question carried into this milestone's drafting

- Does this path reuse the spec's existing exact-slug
  create-or-attach registration path, or introduce new
  registration surface? A technical-direction call for this
  drafting session against merged code — not pre-decided by the
  epic.
