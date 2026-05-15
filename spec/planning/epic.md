# Epic Drafting

Per-level planning playbook for **epic-drafting** sessions. Loads
[`shared.md`](./shared.md) for cross-level planning rules. This file
covers what is unique to the epic level.

## Scope: what an epic does and does not say

Epics scope the *what* and *why* of a multi-milestone arc:
capability targets, cross-cutting invariants, milestone
sequencing rationale, milestone-level risks, and the open
questions the epic resolves or opens. Epics should *not*
prescribe per-milestone phase counts, per-phase content,
per-phase PR counts, validation-gate specifics, documentation
lists, or self-review audit sets. Those details belong to the
milestone planning session for each milestone, against
actually-merged code at milestone-start.

When an epic does name per-milestone details (during initial
epic drafting, before the milestone planning sessions have run),
tag them explicitly as estimates pending milestone planning, not
as binding specs. Sizing summaries in epics carry the same
caveat: per-milestone phase and PR counts are early estimates,
not commitments.

The milestone planning session re-derives the actual phase
shape and the milestone doc supersedes the epic's estimates.
The milestone doc PR also reconciles the epic's prescriptive
paragraphs — either rewriting them to match the milestone-doc
shape, or marking them as pre-milestone-planning estimates and
pointing to the milestone doc as canonical.

## Required and optional sections

An epic doc carries the following sections.

**Required:**

- Status
- Purpose
- Why This Epic
- Goal
- Cross-Cutting Invariants
- Out Of Scope
- Milestone Structure
- Backlog Impact (per [`backlog.md`](../backlog.md))
- Risk Register
- Related Docs

**Optional, when content applies:**

- Milestone Status — when milestones are in flight and a
  status table aids the reader
- Open Questions Resolved By This Epic — when the epic
  resolves prior open questions worth recording
- Open Questions Newly Opened — when the epic surfaces new
  questions for downstream planning
- Sizing Summary — when per-milestone PR or time estimates
  exist
- Documentation currency in implementing PRs — when
  status-bearing docs need updates across the epic's lifecycle
- Inherited Context — when prior epics or external decisions
  bound the scope worth restating

Variance from this list — an unlisted section is appropriate,
or a required section genuinely doesn't apply — follows the
[`shared.md`](./shared.md) "Section variance disclosure" rule.

## Path conventions

An epic gets its own folder at `docs/plans/<epic-slug>/`, with
the epic doc at `docs/plans/<epic-slug>/README.md` and
per-milestone / per-phase docs as siblings inside the same
folder (filename patterns named in [`milestone.md`](./milestone.md)
and [`task-plan.md`](./task-plan.md)). The full layout convention is in
[`planning-doc-location.md`](../planning-doc-location.md).

Per-epic milestone numbering is canonical: each epic counts from
`m1` independently, and sibling epics may reuse the same
milestone numbers without collision because the slug's root
segment (`<epic-slug>`) disambiguates.
