---
slug: post-m2-ux-correction-p2
Status: In draft
short_description: Humanize the forest actor — reported name (slug fallback) in the per-node actor-marker (F4)
---

# Phase 2 — Humanize forest actor (F4)

> **Parent-promotion skeleton.** This file was seeded by the
> [`post-m2-ux-correction`](README.md) task plan's
> `` `In draft` → `Proposed` `` promotion gate per
> [`shared.md`](../../../spec/planning/shared.md) "Parent-promotion
> stub seeding." It carries only the frontmatter and the inherited
> contract block below; per the same rule, it is **exempt** from
> the [`task-plan.md`](../../../spec/planning/task-plan.md)
> "Required and optional sections" rule and from
> [`shared.md`](../../../spec/planning/shared.md) "Plans describe
> contracts, not implementation" until this phase's own
> `` `In draft` → `Proposed` `` drafting session begins. When that
> drafting starts, the skeleton exemption ends and the standard
> task-plan rules apply.

## Context

This is phase 2 of three for the
[`post-m2-ux-correction`](README.md) task. p2 humanizes the
forest's per-node actor display — rendering the reported
session `name` (slug fallback per the name-then-slug rule t4
locked for the roster) inside the `actor-marker` span, never
the raw `wst-<uuid>` actor. Ships after p1 (cosmetic defects)
so the dark-mode-legible chrome is in place for visual
verification; ships before p3 so the actor-marker change isn't
entangled with p3's contract-revisiting render reshuffle. No
plan-doc supersessions in this phase; the
[`humanize-forest-actor`](../../backlog.md#humanize-forest-actor)
backlog entry's Status flip lives in the task-plan drafting
change — the work that closes the entry's intent lands here.

## Inherited Contract

Copied from the parent task plan's
[`## Phase Contracts`](README.md#phase-contracts) p2 row at
parent promotion. Authoritative until the phase's drafting
session restates it inline.

**Short.** Humanize forest actor (F4).

**End result + preserves.** The forest's per-node `actor-marker`
span renders the reported session `name` (slug fallback per the
name-then-slug rule t4 locked for the roster
[`t4-session-roster.md` "Session identity and naming"](../workstream-tracker-1-0/m2/t4-session-roster.md)),
not the raw `wst-<uuid>` actor (parent contract **C6**);
graduates the
[`humanize-forest-actor`](../../backlog.md#humanize-forest-actor)
backlog entry (the entry's `Graduated — post-m2-ux-correction`
Status flip lives in the task-plan drafting change; the work
that closes the entry's intent lands here). No plan-doc
supersession. Preserves: m2 v0.1 actor-tag no-regress invariant
(the `actor-marker` span itself, its placement inside
`label-group`, attachment to every work-instance via per-node
`range .WorkInstances`); the `wst-<uuid>` actor stays the
internal identity key (no actor-generator change, no schema
change); walk-on-every-request.

**Feeds siblings.** Closes the m2-shipped forest/roster
identity disagreement before p3's contract-revisiting work
reshapes the surrounding render. Independent of p1's scope
(different file regions).

## Related Docs

- [`README.md`](README.md) — parent task plan; the locked
  Goal, Contract C6, Cross-Cutting Invariants (C-INV-3
  actor-tag-preserved, C-INV-4, C-INV-5), Files to touch (p2
  estimate; possible additive `WorkInstanceView` passthrough),
  and the task-terminal Validation Gate this phase contributes
  to.
- [`scoping/post-m2-ux-correction.md`](scoping/post-m2-ux-correction.md)
  — paired scoping doc; SD6 records the scoping-time
  deliberation this phase realizes.
- [`../../backlog.md`](../../backlog.md) `humanize-forest-actor`
  — the graduated backlog entry whose intent this phase closes.
- [`../workstream-tracker-1-0/m2/t4-session-roster.md`](../workstream-tracker-1-0/m2/t4-session-roster.md)
  — the name-then-slug rule t4 locked for the roster, reused
  here for the forest.
