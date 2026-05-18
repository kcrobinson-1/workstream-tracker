---
slug: tool-originated-task-sessions-m1-t1
Status: In draft
short_description: Deterministic-path contract expression
---

# m1 t1 — Deterministic-Path Contract Expression

> **Parent-promotion skeleton.** Seeded when
> [`tool-originated-task-sessions-m1`](./README.md) passed its
> `In draft` → `Proposed` promotion gate. Until this task's own
> planning session runs, it carries only the inherited locked
> contract below and is exempt from the task-plan
> required-sections rule and from
> [`shared.md`](../../../../spec/planning/shared.md) "Plans
> describe contracts, not implementation," per "Parent-doc child
> contracts → Parent-promotion stub seeding." The task planning
> session re-derives HOW (file inventory, validation, risks,
> rule-additions retire-or-merge target) against merged code; the
> WHAT below is locked input, not to be loosened without
> reopening the milestone.

## Inherited Contract (locked by the milestone)

**End result.** The deterministic, handshake-free,
slug-by-construction registration path is expressed as a
**first-class additive contract** in the plan-doc spec
([`spec/planning/shared.md`](../../../../spec/planning/shared.md))
and the session-registration agent rule
([`docs/agents/local/session-registration.md`](../../../../docs/agents/local/session-registration.md)):
when a session's canonical slug is known *by construction*
(carried in, not resolved from a natural-language prompt),
registration is deterministic and the narration handshake does
not apply.

**Preserves**, unchanged: the interactive best-effort grounded
narration handshake for natural-language sessions; the spec's
exact-slug create-or-attach posture; the API endpoint / request /
schema (no change — the deterministic path is a pure consumer of
the existing exact-slug create-or-attach flow). Adds **no**
tool-origination or spawn-UX content (that is m2); the expression
stays generic to construction-known slugs, with m1's manual / CLI
`--slug` argument as the stand-in producer.

**Sibling interface.** Produces the **durable contract** the m2
tool-originated-session UX milestone consumes when its spawn
becomes the real construction-time slug producer. The reused
entrypoint it sanctions is the existing
`workstream-tracker register --slug <slug>` over the exact-slug
create-or-attach flow — no new registration surface (milestone
Cross-Task Decision D1).

**Constraints carried from the milestone.** Edits to
[`docs/agents/local/**`](../../../../docs/agents/local/) follow
[`rule-additions.md`](../../../../docs/agents/shared/meta/rule-additions.md)
(name the rule retired/merged, or why none); edits to
[`spec/**`](../../../../spec/) follow the
[`AGENTS.md`](../../../../AGENTS.md) `spec-authoring` pre-edit
read. "Deterministic" stays epic-scoped: identity-by-construction
for a session that *did* launch, not "the launch cannot fail,"
not enrichment made load-bearing.

## Open Questions carried into this task's planning

- Exact section placement and wording in
  [`shared.md`](../../../../spec/planning/shared.md) and
  [`session-registration.md`](../../../../docs/agents/local/session-registration.md),
  and the `rule-additions.md` retire-or-merge target for the
  agent-rule addition (milestone Cross-Task Decision D3 — a HOW
  call for this task's planning, not pre-decided).
