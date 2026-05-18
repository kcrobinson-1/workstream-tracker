---
slug: tool-originated-task-sessions-m2
Status: In draft
short_description: Tool-originated session UX — spawn, slug-carry, hook registration
---

# m2 — Tool-Originated Session UX

> **Parent-promotion stub — scope NOT yet locked.** Seeded for
> tree visibility and slug identity while
> [`tool-originated-task-sessions`](../README.md) is `Proposed`.
> Unlike a locked-contract skeleton, **the content below is
> directional, not a sealed WHAT contract.** This milestone's
> WHAT is sealed at its own milestone-planning session, run
> *against merged m1 code* — the drafting session re-derives the
> contract, the task/phase breakdown, validation, and risks; it
> must not treat the directional text here as locked. The stub
> is exempt from the milestone-doc required-sections rule until
> that session runs, per
> [`shared.md`](../../../../spec/planning/shared.md) "Parent-doc
> child contracts → Parent-promotion stub seeding."

## Directional intent (NOT locked)

The plan-tree-node affordance that spawns the agent session and
becomes the *real* construction-time slug producer feeding m1's
path — the actual tool-originated-task-sessions capability and
the determinism-resolution home.

Directionally: an explicit contributor action on a node spawns a
real local interactive session, handed the mode-appropriate
prompt, carrying the node's slug out-of-band; a session-start
hook registers it deterministically via m1's path; the agent
self-provisions its worktree and reports it as best-effort
enrichment. Preserves observe-only for every session the tool
did not originate.

**Sibling interface.** Consumes m1's slug-carried registration
entrypoint as the real construction-time slug producer (m1's
manual slug argument is the stand-in this milestone replaces).

## Resolved vision inputs this milestone inherits

These are settled by the epic; the drafting session starts from
them rather than re-opening them (see the epic's "Open Questions
Resolved By This Epic"):

- **Spawn shape.** Real local interactive session, prompt handed
  over, slug carried out-of-band, deterministic registration via
  a session-start hook (Claude Code `SessionStart` verified).
- **One-way reconciliation.** The rule stands; this is the one
  bounded, fenced exception (human-initiated, one-shot at birth,
  identity-not-correction, running session stays observe-only).
- **Workspace origin.** Dissolved — the tool has no workspace
  role; the prompt instructs a fresh worktree off `origin/main`;
  the worktree name returns as display-only enrichment.
- **Level/mode.** The mode is the prompt; which mode a node
  offers is a shallow static map over node-type + plan-doc
  Status, no deep logic.

## Open questions carried into this milestone's drafting

Not scope-locking for the epic; this milestone's to resolve at
its own planning, against merged m1 code:

- How the plan-doc spec and the session-registration agent rule
  express the deterministic, handshake-free path additively.
- Whether to consume m1's exact-slug create-or-attach path or
  new surface — follows m1's own resolution of that question
  once m1 lands.

Task and phase breakdown are this milestone-planning session's
output, **not** enumerated here.
