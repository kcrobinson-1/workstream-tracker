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

### nested-milestone-doc-layout

**Status:** Graduated — nested-milestone-doc-layout

**Plan:** [`docs/plans/nested-milestone-doc-layout/README.md`](plans/nested-milestone-doc-layout/README.md)

Nest milestone docs under per-milestone folders in the plan
layout.

A flat epic-rooted tree puts every milestone, task, and phase
doc as siblings in one folder; a multi-milestone epic
accumulates dozens of flat files (`workstream-tracker-1-0/` is
already ~13). The layout convention nests each milestone's doc
and descendants under `docs/plans/<root>/m<N>/`, with scoping at
`m<N>/scoping/`; standalone task plans are unchanged. The
visualization walker is already recursive and slug-identity
driven, so this is a spec-prose change plus a dogfood test, with
no walker/tree behavior change.

### plan-doc-child-contracts

**Status:** Graduated — spec-updates-contracts-and-gates

**Plan:** [`docs/plans/spec-updates-contracts-and-gates/README.md`](plans/spec-updates-contracts-and-gates/README.md)

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

**Status:** Graduated — stub-children-on-parent-promotion

**Plan:** [`docs/plans/stub-children-on-parent-promotion/README.md`](plans/stub-children-on-parent-promotion/README.md)

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

### stale-skeleton-on-parent-reopen

**Status:** Open

What happens to a still-pristine seeded child skeleton when its
parent is re-opened, its child contracts revised, and it is
re-promoted.

The stub-children-on-parent-promotion rule binds only that
re-seeding never clobbers a child already seeded, drafted, or
advanced. It deliberately leaves undecided what to do about a
child that is *still a pristine skeleton* whose parent contract
changed during a `Deferred → In draft` re-opening: such a
skeleton can carry a contract that lags the re-locked parent, so
the child's later drafting starts from stale requirements. This
surfaced in review of the implementing PR and was deferred rather
than decided reactively under bot pressure. Options among
several: leave the stale skeleton as an accepted observable
residual; auto-refresh inherited-contract text and
`short_description` (never the immutable `slug`); or require an
explicit per-child reconciliation step at re-promotion. It
intersects the broader `Deferred`-resumption semantics and
should be deliberated with them, not in isolation.

### promotion-gate-explicit-checklist

**Status:** Graduated — promotion-gate-explicit-checklist

**Plan:** [`docs/plans/promotion-gate-explicit-checklist/README.md`](plans/promotion-gate-explicit-checklist/README.md)

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

### deterministic-interactive-registration

**Status:** Open

Deterministic work-instance registration for interactive,
natural-language sessions.

m1-t2 ships interactive registration as an observable
best-effort grounded narration handshake. It is *not*
deterministic, and provably cannot be: resolving natural-language
intent to a canonical slug requires agent cognition, which
postdates session start, while a deterministic trigger must run
before it (the "registration circularity"). Observable
best-effort is the accepted, vision-faithful posture, not a
stopgap: the m1-t2 decision (locked at t2 drafting) treats the
grounded narration handshake as a faithful operationalization of
the long-term vision's prescribed mitigation, and "deterministic"
was never an upstream promise for registration. Achieving
determinism is therefore **not a 1.0 requirement** — its only
resolution home is the post-v0.2
[`tool-originated-task-sessions`](#tool-originated-task-sessions)
capability (the tool itself launches the agent, so the slug is
known by construction), which is well beyond this epic's scope;
the registration circularity is dissolved only by changing who
the launcher is, and that is where it gets dissolved. What *does*
bind: best-effort's acceptability rests on sole-consumer
compensation (the lone producer can notice and hand-fix a missed
marker), and that compensation evaporates when an external
project adopts the tool. **Tripwire: this must be
re-deliberated — not necessarily resolved — before the
neighborly-events integration milestone**, where sole-consumer
compensation no longer holds; that re-deliberation decides
whether best-effort is still acceptable at that point or whether
a backstop must be pulled forward, and is not a commitment that
registration becomes deterministic for 1.0.

This entry also tracks the **observability residual**: a missed
registration is unobservable from the rendered tree by
construction (an unregistered session emits no signal the tool
ever sees, so the tree cannot distinguish unregistered work from
no work). t2's only backstop is the in-session narration
handshake, which works solely while a contributor is present to
notice it — the same sole-consumer compensation the tripwire is
about. A tree-side heuristic ("a node with an active/Proposed
plan doc but no work-instance") is the candidate future
affordance, deferred under the same tripwire. This entry tracks
both the determinism gap and the observability gap until
resolved.

### tool-originated-task-sessions

**Status:** Open

The tool's own UX originates a planning/implementation session
from a plan-tree node.

Today the contributor opens an agent and states intent in
natural language; the tool only ever *observes* work. An
opportunity: a "plan this task" / "work this task" affordance on
a plan-tree node in the tool's UX that spawns the agent
session itself. Because the tool is the launcher and the click
carries the node's identity, the spawned session's canonical
slug is known *by construction* before any agent cognition —
which dissolves the registration circularity and makes
work-instance registration deterministic for free (no slug
resolution, no narration handshake needed for these sessions).
Beyond registration it closes a larger loop: the visualization
stops merely observing work and starts originating
correctly-attributed work. Well beyond v0.2 (the tool acting /
spawning agents is far past the epic's scope). One option among
several, opportunity-framed; this is the home where the
determinism deferred by
[`deterministic-interactive-registration`](#deterministic-interactive-registration)
is eventually achieved.

### templ-render-adoption

**Status:** Graduated — templ-render-adoption

**Plan:** [`docs/plans/templ-render-adoption/README.md`](plans/templ-render-adoption/README.md)

Revisit `templ` for HTML rendering, or formally accept stdlib
`html/template` as the v0.1+ choice.

[`design/v0.1-design.md`](../design/v0.1-design.md) §10 locked
`templ` for HTML rendering, but the v0.1 implementation
deliberately used the stdlib `html/template` instead, deferring
templ "until there are real reusable components" (the decision
and its rationale are recorded in git, commit `60040be`); §10
was reconciled to point here rather than left silently
contradicting the codebase. The m2-t1 region split (a shell
composing `forest` / `roster` / `node` sub-templates in their
own files) is the first plausible "real reusable components"
trigger, so the open question is whether to migrate the render
layer to templ for its type-safe component model or formally
accept `html/template` as the standing choice. One option among
several: a throwaway spike converting the forest/roster/node
templates to templ components to weigh the
ergonomics-versus-extra-dependency tradeoff §10 originally
cited, before committing either way.
