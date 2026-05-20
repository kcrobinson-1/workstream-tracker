---
slug: post-m2-ux-correction-p3
Status: In draft
short_description: Contract-revisiting + task-terminal — default D/P/I/V cells (F3a) + nested-details body disclosure (F7) + every-entry-opens roster (F9); carries the full product-acceptance walkthrough
---

# Phase 3 — Contract-revisiting + task-terminal (F3a + F7 + F9)

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

This is phase 3 of three for the
[`post-m2-ux-correction`](README.md) task, **task-terminal**.
p3 ships the three contract-revisiting findings — the default
D / P / I / V Status-driven cell row when no `progress_stages`
is declared (F3a), the header-only default node-box view with
body content behind a nested `<details>` (F7), and every roster
entry opening to the K3 known-facts + raw-JSON disclosure (F9).
It carries all three plan-doc supersessions (D1 / D2 / D3 in
the parent task plan's `## Status` →
`### Supersessions of sibling contracts` sub-block), the full
product-acceptance walkthrough across all six findings, and the
post-m2-ux-correction task plan's
`Validating → Landed` flip per
[`task-plan.md`](../../../spec/planning/task-plan.md) "Task
plan terminal state when N ≥ 2." Ships last in the sequence so
the supersession review concentrates with the rendered
consequence.

## Inherited Contract

Copied from the parent task plan's
[`## Phase Contracts`](README.md#phase-contracts) p3 row at
parent promotion. Authoritative until the phase's drafting
session restates it inline.

**Short.** Contract-revisiting + task-terminal (F3a + F7 + F9).

**End result + preserves.** A node whose doc declares no
`progress_stages` renders the default D / P / I / V
Status-driven cell row (parent contract **C5**; supersedes m2
t3 D5 under **D1** and the stub-children stub-render contract
under **D2**); every node box's first view is header-only, with
the body content (long description + related PRs) behind a
nested `<details>`/`<summary>` disclosure (OD2 = H1) rendering
markdown with `<a>` tags stripped (OD3 = I2b) when expanded,
independent of the parent tree-collapse (parent contract
**C2**); every roster entry opens to the K3 known-facts +
raw-JSON disclosure (parent contract **C4**; supersedes t4's
"no metadata ⇒ plain row" under **D3**). **Task-terminal** —
carries the full product-acceptance walkthrough across all
six findings + the post-m2-ux-correction task plan's
`Validating → Landed` flip. Preserves: m2 file-enforced region
invariants (forest edits in `forest.go`, roster edits in
`roster.go`, no shell edit); the existing field-presence-gated
declared-stages render unchanged for any doc that declares
`progress_stages`; t4's schema-loose-not-homogenized invariant
(the K3 known-facts header is a sibling surface for
loader-known facts, not a schema imposed on the
reported-metadata blob); the m2 v0.1 actor-tag no-regress
invariant (p3 does not touch the `actor-marker` span p2 ships).

**Feeds siblings.** Realizes all three plan-doc supersessions
in one phase so the supersession review concentrates with the
rendered consequence. Carries Cross-Cutting Invariant
**C-INV-1 (cell-anchor)** — the default-cells row's per-cell
DOM is the F3b future attachment surface. The task plan's
terminal close-out (Status flip; the Backlog Impact mutations
already executed in the drafting change need no further surgery
here).

## Related Docs

- [`README.md`](README.md) — parent task plan; the locked
  Goal, Contracts C2 / C4 / C5, all Cross-Cutting Invariants
  C-INV-1 through C-INV-5, Files to touch (p3 estimate;
  possible additive `RosterEntry` extension for registered-at /
  last-event), the task-terminal Validation Gate (the full
  product-acceptance walkthrough), and the
  `## Status` → `### Supersessions of sibling contracts`
  sub-block recording D1 / D2 / D3.
- [`scoping/post-m2-ux-correction.md`](scoping/post-m2-ux-correction.md)
  — paired scoping doc; SD2 (F7), SD4 (F9), SD5 (F3a) record
  the scoping-time deliberation this phase realizes; the
  OD-walk outcomes block locks OD2 = H1, OD3 = I2b, OD5 = K3.
- [`../workstream-tracker-1-0/m2/t3-doc-declared-stages.md`](../workstream-tracker-1-0/m2/t3-doc-declared-stages.md),
  [`../workstream-tracker-1-0/m2/t4-session-roster.md`](../workstream-tracker-1-0/m2/t4-session-roster.md),
  [`../stub-children-on-parent-promotion/README.md`](../stub-children-on-parent-promotion/README.md)
  — the three Landed contracts this phase supersedes (D1, D3,
  D2 respectively). Not retro-edited.
- [`../../backlog.md`](../../backlog.md)
  `progress-cell-active-state-and-actor` — the F3b carve-out
  whose linear-not-cliff migration cost depends on this
  phase's C-INV-1 (cell-anchor) realization.
