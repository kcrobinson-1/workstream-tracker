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

Parent-level plan-doc governance: missing child contracts and
missing promotion gates.

Parent docs (epic, milestone) are underspecified versus
task/phase plans on two fronts. (1) No per-child WHAT-contract
section — what each child delivers, sibling interfaces,
preserves; the "Anti-goal: do not scope any phase" rule bars
HOW-scoping but inadvertently bars WHAT-contracting too,
leaving locked-scope docs with only child names. (2) No
`In draft` → `Proposed` promotion gate — the gate in
`task-plan.md` binds task/phase plans only, so a PR that locks
a milestone's or epic's scope has no rule telling it to flip
Status; this is why `m1-v0-2.md` sat at `In draft` after its
scope-locking PR merged. Both gaps repeat at epic level
(Milestone Structure carries only titles; no gate) and
task-plan level (no Phase Contracts when N≥2). One option
among several: a cross-level rule in
[`spec/planning/shared.md`](../spec/planning/shared.md) for
the WHAT/HOW split plus a parent-doc promotion gate, with
matching section/gate additions in
[`spec/planning/epic.md`](../spec/planning/epic.md),
[`spec/planning/milestone.md`](../spec/planning/milestone.md),
and
[`spec/planning/task-plan.md`](../spec/planning/task-plan.md)
(Phase Contracts when N≥2), and the "Anti-goal: do not scope"
rules refined to forbid HOW-scoping while requiring
WHAT-contracting.

Also surfaced: milestone.md's required sections use "Phase"
naming for the milestone's child, but the slug grammar calls
that unit a task; fold into the same spec edit or split into
its own entry when this graduates.
