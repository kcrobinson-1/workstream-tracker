---
slug: post-m2-ux-correction-p1
Status: In draft
short_description: Cosmetic defects — dark-mode chrome legibility (F1) + baseline-aligned collapse marker (F8)
---

# Phase 1 — Cosmetic defects (F1 + F8)

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

This is phase 1 of three for the
[`post-m2-ux-correction`](README.md) task. p1 ships pure-CSS
fixes for two m2-shipped cosmetic defects — dark-mode page
chrome (F1) and the misplaced native `<details>` disclosure
marker (F8) — first in the sequence so p2 and p3 observe their
findings against a legible page. No plan-doc supersessions in
this phase.

## Inherited Contract

Copied from the parent task plan's
[`## Phase Contracts`](README.md#phase-contracts) p1 row at
parent promotion. Authoritative until the phase's drafting
session restates it inline.

**Short.** Cosmetic defects (F1 + F8).

**End result + preserves.** The page chrome is legible in OS
dark mode — `<h1>` and inter-card chrome contrast the dark
canvas (parent contract **C1**); the native disclosure marker
on each node box aligns with the label baseline and does not
overlap the box border (parent contract **C3**). No plan-doc
supersessions. Preserves: m2 file-enforced region invariants
(the only `render.go` edit is the body rule; no shell-layout /
region-boundary / sibling-region body edit); the no-JS native
`<details>`/`<summary>` idiom (m2 t2 Collapsible-mechanism
decision); walk-on-every-request render path.

**Feeds siblings.** Establishes the dark-mode-legible page
chrome and the baseline-aligned marker that p2 and p3 observe
their findings against (the demo walkthrough's dark + light
passes assume legible chrome). Does not change forest content
semantics; does not alter the roster region.

## Related Docs

- [`README.md`](README.md) — parent task plan; the locked
  Goal, Contracts (C1, C3), Cross-Cutting Invariants (C-INV-2,
  C-INV-4, C-INV-5), Files to touch (p1 estimate), and the
  task-terminal Validation Gate this phase contributes to.
- [`scoping/post-m2-ux-correction.md`](scoping/post-m2-ux-correction.md)
  — paired scoping doc; SD1 (F1) and SD3 (F8) record the
  scoping-time deliberation this phase realizes.
