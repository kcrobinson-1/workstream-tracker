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

**Status:** Graduated — spec-updates-contracts-and-gates

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

### stub-children-on-parent-promotion

**Status:** Open

Seed skeleton child docs when a parent doc passes its
promotion gate.

When a parent (epic/milestone) locks scope and promotes, its
children exist only as names in a table; the child slugs have
no doc until each is drafted just-in-time. Seeding skeleton
child docs (slug + Status frontmatter) in the parent's
promoting PR would make planned-but-unstarted nodes render in
the tree earlier and give auto-registration a frontmatter slug
to derive from before any drafting session. The opportunity
intersects two unsettled areas and should be deliberated with
them, not in isolation: the deferred "Triage zone in 1.0?"
question (a stub doc is not a work-instance row — orphan/
unattached WIs still need a home) and the spec rule that
descendant slugs are *server-generated* (convention-assigned
stub slugs could diverge from the server's allocation counter;
the milestone Task Status tables already name `-tN` slugs by
convention, hardening that latent inconsistency). One option
among several: a parent-promotion-gate rule in
[`spec/planning/`](../spec/planning/) that emits child
skeletons with a reconciliation story for server slug
allocation. Note: this does not remove the need for an
exact-slug create-or-attach registration path
(workstream-tracker-1-0-m1-t1) — stub-seeding depends on that
path rather than replacing it.

### promotion-gate-explicit-checklist

**Status:** Open

The `In draft` → `Proposed` promotion gate under-specifies
its minimum checks.

The gate in
[`spec/planning/task-plan.md`](../spec/planning/task-plan.md)
enumerates four steps (end-to-end coherence, Contracts
decision-completeness, universal `Verified by:` walk,
reality-check reconfirmation) but does not name three checks a
promoting agent is expected to perform at minimum:
required-sections presence, no-implementation-prescription, and
conformance to the broader guiding specs. Those are enforced by
separate always-on rules ("Required and optional sections",
"Plans describe contracts, not implementation", "Section
variance disclosure", etc.), so a gate run can pass its four
named steps without explicitly covering them — the gap surfaced
when running the gate on
`workstream-tracker-1-0-m1-t1` required ad-hoc augmentation.
One option among several: extend the gate's step list to
reference those rules by name (a superset checklist), so the
gate is self-contained rather than relying on the runner to
remember the adjacent always-on rules.
