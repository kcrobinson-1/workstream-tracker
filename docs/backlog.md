# Backlog

Captured-but-unplanned work for workstream-tracker. Entries
graduate into plan-tree nodes when picked up; format and
lifecycle live in [`../spec/backlog.md`](../spec/backlog.md).

### repo-rooted-doc-links

**Status:** Open

Repo-rooted markdown link syntax in plan and spec docs.

Cross-doc links currently use file-relative paths like
`../../../design/vision.md`, which read poorly from deeply
nested plan docs and break when files move. The blocker is
that standard markdown resolves unprefixed paths relative to
the source file, leading `/` resolves to host root on GitHub,
and full URLs are bloated. One option among several: extend
the `{spec_root}` substitution mechanism in
`scripts/assemble.sh` with a sibling `{repo_root}` placeholder
applied to in-repo docs at write time.

### plan-doc-child-contracts

**Status:** Open

Codify child contracts as a required section at every parent
plan-tree level.

The current spec has Cross-Phase Invariants and Cross-Phase
Decisions (rules binding multiple children, contracts between
children) at the milestone level, but no section for per-child
**WHAT** contracts — what each child delivers, sibling
interfaces, preserves. The "Anti-goal: do not scope any phase"
rule in `milestone.md` addresses HOW-scoping but inadvertently
bars WHAT-contracting, leaving locked-scope milestone docs
with only child names and titles. The same gap repeats at the
epic level (Milestone Structure carries only brief titles, no
contracts) and the task-plan level (no Phase Contracts when
N≥2). One option among several: add a cross-level rule to
[`spec/planning/shared.md`](../spec/planning/shared.md)
describing the WHAT/HOW split, with matching
required-section additions in
[`spec/planning/epic.md`](../spec/planning/epic.md)
(Milestone Contracts),
[`spec/planning/milestone.md`](../spec/planning/milestone.md)
(Task Contracts), and
[`spec/planning/task-plan.md`](../spec/planning/task-plan.md)
(Phase Contracts when N≥2). The "Anti-goal: do not scope"
rules get refined to forbid HOW-scoping while requiring
WHAT-contracting.

Also surfaced: milestone.md's required sections use "Phase"
naming for the unit below the milestone, but the slug grammar
calls that unit a task; see also the naming-vs-grammar
question separately if it merits its own entry.
