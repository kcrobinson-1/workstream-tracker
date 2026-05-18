---
slug: tool-originated-task-sessions-m1-t2
Status: In draft
short_description: Determinism proof harness
---

# m1 t2 — Determinism Proof Harness

> **Parent-promotion skeleton.** Seeded when
> [`tool-originated-task-sessions-m1`](./README.md) passed its
> `In draft` → `Proposed` promotion gate. Until this task's own
> planning session runs, it carries only the inherited locked
> contract below and is exempt from the task-plan
> required-sections rule and from
> [`shared.md`](../../../../spec/planning/shared.md) "Plans
> describe contracts, not implementation," per "Parent-doc child
> contracts → Parent-promotion stub seeding." The task planning
> session re-derives HOW (test inventory, fixtures, validation,
> risks) against merged code; the WHAT below is locked input, not
> to be loosened without reopening the milestone.

## Inherited Contract (locked by the milestone)

**End result.** An end-to-end test proves that, given a
construction-known slug supplied as an explicit argument with
**no narration handshake**, a work-instance registers
deterministically against the real exact-slug create-or-attach
path: identity-by-construction (the slug is honored verbatim),
idempotent create-or-attach on `(slug, actor, active)`, no
natural-language-resolution path engaged, server repo-blind. The
slug producer is the **minimal explicit `--slug` argument / test
harness** — the stand-in for m2's future construction-time
producer — scoped as the minimum to validate the registration
contract, **not a durable UX**.

**Preserves**: no product code surface added; the existing
interactive-path tests
([`cmd/workstream-tracker/register_test.go`](../../../../cmd/workstream-tracker/register_test.go)
and siblings) unchanged. The existing
`TestRegisterCommandSuccessAndIdempotentRepeat` already exercises
the CLI end-to-end against a real API server through the
exact-slug flow; t2 adds the *determinism-specific* assertions
(handshake-free, slug-by-construction, no NL resolution engaged)
rather than re-testing idempotency from scratch.

**Sibling interface.** Consumes t1's named determinism properties
as the assertions it must demonstrate; produces the **proof
artifact** that de-risks the registration contract before the m2
posture shift, so m2 consumes a proven path rather than
co-developing registration and the posture change at once.

**Constraints carried from the milestone.** No new registration
surface — no endpoint, request/response field, schema, code path,
or CLI flag (milestone Cross-Task Decision D1). "Deterministic"
stays epic-scoped: the proof asserts only the deterministic
identity step for a session that *did* launch, never "the launch
cannot fail" and never enrichment as load-bearing for
attachment.

## Open Questions carried into this task's planning

- Test inventory and fixture shape (what the
  [`register_test.go`](../../../../cmd/workstream-tracker/register_test.go)
  httptest-API pattern must add to assert the
  determinism-specific properties t1 names) — a HOW call for this
  task's planning against then-merged code, sequenced after t1 so
  the assertions cite the locked contract wording (milestone
  Sequencing).
