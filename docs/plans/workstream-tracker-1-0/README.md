---
slug: workstream-tracker-1-0
Status: In draft
---

# workstream-tracker — 1.0 Launch

## Purpose

Scope the path from v0.1 to a 1.0 launch defined as "the version
neighborly-events adopts as its first real consumer." The doc
names the milestones the tool must clear before that adoption is
supportable and the open questions the epic must resolve along
the way. Working surface across the epic's lifecycle — sessions
promote open questions into committed milestones (or to Out of
Scope) as deliberation locks them, and refine per-milestone scope
when a milestone moves toward drafting.

## Why This Epic

v0.1 shipped as a thin slice: the plan tree renders, agents
register, the dogfood loop closes. None of v0.1's gaps block the
dogfood — but together they keep the tool from being something an
external project could adopt. Alphabetical order buries live
work, slug labels read poorly, per-node detail doesn't show where
each agent actually is, and registration is manual enough that
the attached work-instance markers aren't trustworthy.

**neighborly-events** is the natural first consumer: an active
project with parallel agent work, the contributor's other live
workstream, and a real test of whether the spec, agent rules, and
visualization hold up outside this repo. Validating against it is
what changes "tool the author uses on itself" into "tool that has
a consumer."

The single customer for 1.0 is the same person writing this plan,
so the exact line stays a little blurry — calls about which
features matter resolve under real use, not in advance. The
criterion is qualitative: opening the page has to make the
parallel-agent picture legible at a glance.

## Goal

Ship workstream-tracker at a quality bar where neighborly-events
can vendor the spec and shared agent rules, populate its plan
tree, and run real sessions against the tracker without local
patches or workarounds. The final integration milestone surfaces
whatever gaps remain so they close before 1.0 ships.

## Cross-Cutting Invariants

Rules that bind every milestone in this epic. Reviewer-flag
candidates for any milestone deliberation that brushes against
them.

- **Tool stays one-way.** Server reads agent events; it never
  pushes back. Drift prevention runs through intent in the
  project's files, not through tool→agent injection. (See
  [vision §5](../../../design/vision.md).)
- **Single contributor, single local environment.** No auth, no
  multi-tenant, no remote hosting. neighborly-events adopts a
  local install of the same tool. (See
  [vision §5, §8](../../../design/vision.md).)
- **Spec contract is additive across this epic.** Mid-epic spec
  changes (new optional frontmatter fields, etc.) are additive;
  no breaking changes between milestones. A vendor consumer can
  vendor any milestone's spec snapshot without local patches.
- **Plan-doc spec owns plan-doc shape.** The tool reads what the
  spec defines; the tool does not invent its own plan-doc
  conventions. Tool needs that require new frontmatter or
  status values require spec changes first.

## Milestone Structure

One landed milestone, one in-draft middle milestone, one
early-estimate middle milestone, and a final integration
milestone whose slug is allocated at registration time.

The original `m2` ("activity-first ordering") and `m3`
("sub-stage cells + richer actor presence") were early
estimates made before m1 landed. They are superseded here: the
new `m2` is the in-root expanded view, doc-declared progress
boxes, and session roster (drafted in
[`m2-expanded-view-and-roster.md`](m2-expanded-view-and-roster.md));
the across-roots tier ordering and cell-level actor presence
the old estimates named are re-homed to a later milestone
rather than lost.

**Sequencing rationale** — sequential dependency, not parallel:

- **m1 first** because it unblocks accurate dogfood (multiple
  work-instances per slug, automatic registration) and lays the
  descriptive-label foundation that m2 and m3 both consume.
- **m2 before m3** because m2 makes a single root's internal
  structure legible — expanded per-node Status boxes,
  doc-driven progress boxes, and a session roster. m3's
  across-roots tier ordering and cell-level actor presence sit
  on top of m2's per-node box and roster; the forest-level
  framing is more useful once each root reads correctly
  internally.
- **Final integration last** by definition: it tests the whole
  stack against a real external consumer.

**Milestones:**

- `workstream-tracker-1-0-m1` (Landed — all four tasks
  complete). **v0.2 — Read-experience improvements.** Multi-work-
  instance support per slug, automatic agent registration,
  descriptive tree labels (with frontmatter `short_description`
  parsing as a phase), expanded per-node display. Full task
  list and contracts in [`m1-v0-2.md`](m1-v0-2.md).

- `workstream-tracker-1-0-m2` (Proposed — task scope locked at
  4 tasks; parent-promotion stubs seeded). **v0.3 — In-root
  expanded view, doc-declared progress boxes, and session
  roster.** A root's contents
  render as nested epic → milestone → task → phase boxes with
  per-node Status, collapsible, inside the active-work surface;
  every node level shows progress boxes whose count and order
  come from the doc itself (a stub showing only the Drafting
  box, which requires an additive spec change); and a new
  roster lists bound + unbound active sessions with an
  expandable deliberately-unstructured-JSON detail view. Full
  task list and contracts in
  [`m2-expanded-view-and-roster.md`](m2-expanded-view-and-roster.md).

- `workstream-tracker-1-0-m3` (early estimate — scope not yet
  locked). **Activity-first forest ordering and richer actor
  presence.** The re-homed deferred pieces: across-roots
  tier-based sorting (active surfaces top, in-flight middle,
  landed/abandoned compressed) with collapse of inactive
  branches; actor icons positioned on individual progress boxes
  (assuming m2's node-level actor tags remain); and the richer
  work-instance state vocabulary (`awaiting-user`,
  `awaiting-external`, `backgrounded`) in the API and DB. This
  is an estimate made before m2 lands; its task breakdown is
  the m3 milestone-drafting session's output, not fixed here.

- **Final integration milestone** (slug allocated at
  registration). **Neighborly-events integration.** Vendor
  `spec/` and shared agent rules into neighborly-events,
  populate its `docs/plans/<root>/`, run the first real sessions
  against the tracker, close gaps surfaced by real use. Gates
  the claim that 1.0 is real rather than aspirational.

## Milestone Contracts

Per-milestone **WHAT** contracts — end result, sibling
interfaces, preserves. The **HOW** for each milestone lives in
the milestone doc (and its constituent task plans) when it
drafts. Required section per
[`epic.md`](../../../spec/planning/epic.md) "Required and
optional sections" and [`shared.md`](../../../spec/planning/shared.md)
"Parent-doc child contracts."

### m1 — v0.2 Read-experience improvements

- **End result.** v0.1's bare-bones render lifts to a tree
  that's scannable per node (descriptive labels, frontmatter
  `short_description`, expanded per-node detail) and an agent
  dogfood loop that's trustworthy (multi-work-instance support
  per slug, automatic registration). Task-level breakdown in
  [`m1-v0-2.md`](m1-v0-2.md).
- **Interfaces.** Establishes the `short_description`
  frontmatter field and parsed body that m3's cell rendering
  consumes. Produces a render path that supports multi-WI per
  slug (consumed by m2's tier classification and m3's per-cell
  actor positioning).
- **Preserves.** v0.1's bare-bones render path remains as the
  fallback shape; spec changes are additive (no breaking
  changes to vendored consumers).

### m2 — In-root expanded view, doc-declared progress boxes, and session roster

- **End result.** A root's descendants render as nested
  epic → milestone → task → phase boxes with per-node Status,
  collapsible, inside the active-work surface; every node level
  renders progress boxes whose count and order come from the
  plan doc itself (a stub renders only the Drafting box); and a
  session roster lists bound + unbound active sessions with an
  expandable, deliberately-unstructured-JSON detail view. An
  additive `spec/` change introduces the doc-declared-stages
  affordance (spec-first).
- **Interfaces.** Consumes m1's multi-work-instance schema,
  descriptive labels, and per-node detail. Provides the
  per-node expanded box and the doc-declared-stages spec field
  that m3's across-roots ordering and on-box actor presence
  build upon. The roster is the observability surface the
  `deterministic-interactive-registration` backlog entry's
  mitigation direction names.
- **Preserves.** m1's labels, per-node detail, and
  multi-work-instance schema; v0.1's actor tags on nodes (the
  no-regress invariant m3's on-box actor work depends on); the
  walk-on-every-request render path; the additive-spec posture
  (a doc without the new field still renders).

### m3 — Activity-first forest ordering and richer actor presence

*Scope not yet locked — early estimate, superseding the
original `m2`/`m3` estimates. The m3 milestone-drafting session
locks the contract and re-derives the task breakdown against
the code m2 actually lands; per
[`shared.md`](../../../spec/planning/shared.md) "Parent-doc
child contracts," this names the milestone without sealing its
contract until that session runs.*

- **End result (estimated).** The forest renders with
  across-roots tier-based sorting (active top, in-flight
  middle, landed/abandoned compressed) with inactive branches
  collapsed; actor icons sit on specific progress boxes rather
  than only as node-level tags; the work-instance state
  vocabulary expands to `awaiting-user`, `awaiting-external`,
  `backgrounded`.
- **Interfaces (estimated).** Consumes m1's multi-WI schema,
  m2's per-node expanded box, doc-declared progress boxes, and
  session roster. Adds richer work-instance states to the API
  and DB.
- **Preserves (estimated).** m2's per-node boxes and roster;
  v0.1's node-level actor tags remain as the fallback when
  on-box actor placement isn't populated; existing render path
  for nodes without the richer state vocabulary.

### Final integration milestone — Neighborly-events integration

- **End result.** neighborly-events successfully vendors
  `spec/` and shared agent rules, populates its plan tree,
  and runs real sessions against the workstream-tracker
  without local patches or workarounds. Gaps surfaced by real
  use close before 1.0 ships.
- **Interfaces.** Consumes the entire stack assembled across
  m1, m2, and m3. Exposes the contract surface (spec, agent
  rules, API, render layer) to a real external consumer.
- **Preserves.** All earlier milestone deliverables. The
  integration is additive — a new consumer adopting the
  existing system, not a rework.

## Open Questions Newly Opened

Calls deferred until m1-m3 are sized — the answers may shift
once the surface area is clearer.

- **Triage zone in 1.0? — RESOLVED (split).** Resolved when
  m2's scope was locked. The **session roster** — a surface
  listing bound + unbound active sessions so every session
  *can* be accounted for — **is a committed 1.0 goal**,
  delivered by `m2`
  ([`m2-expanded-view-and-roster.md`](m2-expanded-view-and-roster.md)).
  The **triage *action*** (acting on an unbound session to
  promote it into the tree or dismiss it) is **deferred past
  1.0**: the roster makes unbound work visible; deciding its
  home is a later capability. Reflected in
  [`design/vision.md`](../../../design/vision.md) §4 (the
  narrow "triage zone" framing replaced by the broader
  best-effort "session presence" concept).
- **Actor lineage in 1.0?** Probably no — single contributor
  carries the lineage in their head — but flag so the call is
  conscious.
- **Multi-repo or hosted instance in 1.0?** Almost certainly no
  (single local user) — flag in case the neighborly-events
  integration surfaces a real need.
- **How does "1.0 done" get decided?** The "enough website to
  be usable" criterion is qualitative and the customer is the
  same as the producer. Without an external pressure point,
  drift past 1.0 into a perpetual 0.x is a real risk. Candidate
  circuit-breakers: a self-imposed neighborly-events adoption
  deadline; a feature-count cap; declaring 1.0 the moment the
  integration milestone produces no new gap PRs.

## Out of Scope

- **Post-1.0 evolution:** capabilities that improve the tool
  after neighborly-events is on it. Additional consumer
  onboarding, hosted multi-tenant deployment, advanced
  analytics, full editorial fields (owner / dates /
  dependencies / tags / acceptance criteria) beyond what v0.2
  introduces.
- **Authentication, access control, multi-contributor
  support.** Per the single-contributor cross-cutting
  invariant.
- **Time-series history of state.** The view shows current
  state; reconstructing past state from the event log is not
  exposed to the visualization. Per
  [vision §8](../../../design/vision.md).
- **Adaptive hierarchies for other contributors' workflows.**
  The tool stays wired to this contributor's conventions
  through 1.0; generalization is a post-1.0 path that runs
  through a private beta phase. Per
  [vision §5](../../../design/vision.md).
- **Intent layer (intent strip + authoring nudges).** Resolved
  at the m1 retrospective as a deliberate deferral past 1.0:
  there is no concrete pain point today — the "plan docs lose
  their intent" problem the intent layer was meant to solve is
  currently handled by capturing intent directly inside scope
  docs during scoping. It stays core product vision and a
  possible post-1.0 direction, not cancelled. This is a
  roadmap/sequencing call, not backlog-sized work; recorded in
  [vision §4](../../../design/vision.md).
- **Anything not on the path to "neighborly-events can adopt
  this."** If a candidate area can be cut without blocking
  that adoption, it gets cut.

## Backlog Impact

A backlog file at [`../../backlog.md`](../../backlog.md) was
created during this epic's deliberation with one entry —
[`repo-rooted-doc-links`](../../backlog.md#repo-rooted-doc-links)
— surfaced for post-1.0 work. The epic itself doesn't graduate
from a backlog entry (framed directly as an epic) and doesn't
graduate, delete, split, or shift any backlog entry per the
effect taxonomy in
[`../../../spec/backlog.md`](../../../spec/backlog.md).

The Open Questions above remain intra-epic deliberations; if
any get split out as standalone work, they become backlog
entries at that moment.

## Risk Register

- **"1.0 done" can't be disproven.** Single customer plus
  qualitative criterion means the producer can shift the bar
  indefinitely. Mitigation: pick one circuit-breaker (see Open
  Questions above) before m3 sizes; commit to it.
- **neighborly-events integration surfaces a gap requiring a
  new milestone.** A pre-1.0 milestone may turn out
  insufficient against real use. Mitigation: the final
  milestone is explicitly the gap-closing step; expect at
  least some feedback loop rather than treating it as
  schedule slip.
- **Deferred open questions re-open after m1 lands.** The two
  highest-weight questions are now resolved at the m1
  retrospective: the intent layer is deferred past 1.0 (recorded
  in Out of Scope), and the triage-zone question is split (the
  session roster is committed to m2; the triage *action* is
  deferred past 1.0). The residual risk is the remaining
  lower-weight questions (actor lineage, multi-repo/hosted,
  "1.0 done" criterion). Mitigation: revisit the remaining open
  questions explicitly at the m1/m2 retrospective, not
  passively.
- **Mid-epic spec changes break vendor consumers.** Additive
  posture is the invariant, but a milestone that doesn't
  realize it's introducing a breaking change is a real
  surface. Mitigation: spec-changing milestones carry an
  explicit "is this additive?" check in their planning.
- **Agent auto-registration (m1, t2) harder than estimated —
  RESOLVED, no escalation.** The risk assumed fragile slug
  derivation from agent context. t2 drafting reframed it via the
  registration circularity: deterministic pre-session derivation
  is impossible, so the agent resolves the slug mid-session and
  registers via a grounded narration handshake (observable
  best-effort, not deterministic). t2 landed within m1 without
  escalating to its own milestone; the residual determinism gap
  is tracked by the
  [`deterministic-interactive-registration`](../../backlog.md#deterministic-interactive-registration)
  backlog entry with a re-deliberation tripwire at the
  neighborly-events integration milestone.

## Sizing Summary

Per-milestone task counts. m2's count is proposed by its
drafted milestone doc (pending review); m3 and the final
milestone remain estimates pending their planning sessions.

- **m1**: 4 tasks. Landed (all four tasks complete). See
  [`m1-v0-2.md`](m1-v0-2.md).
- **m2**: 4 tasks, scope locked (Proposed; stubs seeded): site skeleton
  (two-region shell); expanded nested-box render; doc-declared
  progress stages (spec-first); session roster + work-item
  enrichment. The skeleton (t1) ships the approved side-by-side
  shell so t2/t4 build into independently-owned regions in
  parallel; t2/t3/t4 each plausibly N ≥ 2, phase splits
  re-derived at task drafting. See
  [`m2-expanded-view-and-roster.md`](m2-expanded-view-and-roster.md).
- **m3**: estimate pending the m3 milestone-drafting session.
  Re-homed deferred pieces: across-roots tier classification
  and sort, collapse of inactive branches, on-box actor-icon
  positioning, richer work-instance state vocabulary in
  schema/API. Task count is the m3 session's output, not fixed
  here.
- **Final integration**: 1-3 tasks estimated. Vendoring,
  population, first-real-session validation, gap-fix cleanup.

## Related Docs

- [`m1-v0-2.md`](m1-v0-2.md) — m1 milestone doc (drafted,
  scope locked).
- [`../../../design/vision.md`](../../../design/vision.md) —
  the long-term vision; this epic's scope decisions cite it
  for what defers and why.
- [`../../../design/v0.1-design.md`](../../../design/v0.1-design.md)
  — the design v0.1 ships against; this epic builds on
  v0.1's data model and API shape.
- [`../../../spec/`](../../../spec/) — the plan-doc spec
  consumer projects (starting with neighborly-events) vendor.
- [`../../../spec/planning/epic.md`](../../../spec/planning/epic.md)
  — the rules this epic doc is structured against.
