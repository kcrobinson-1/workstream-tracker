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

One drafted milestone, two proposed middle milestones, and a
final integration milestone whose slug is allocated at
registration time.

**Sequencing rationale** — sequential dependency, not parallel:

- **m1 first** because it unblocks accurate dogfood (multiple
  work-instances per slug, automatic registration) and lays the
  descriptive-label foundation that m2 and m3 both consume.
- **m2 before m3** because tier-based forest ordering frames
  where the eye looks first; sub-stage cells are leaf-level
  detail that's more useful once the forest reads correctly.
- **Final integration last** by definition: it tests the whole
  stack against a real external consumer.

**Milestones:**

- `workstream-tracker-1-0-m1` (drafted, scope locked at four
  tasks). **v0.2 — Read-experience improvements.** Multi-work-
  instance support per slug, automatic agent registration,
  descriptive tree labels (with frontmatter `short_description`
  parsing as a phase), expanded per-node display. Full task
  list and contracts in [`m1-v0-2.md`](m1-v0-2.md).

- `workstream-tracker-1-0-m2` (proposed). **Activity-first
  ordering.** Tier-based sorting (active surfaces top, in-flight
  middle, landed/abandoned compressed), expand active-work paths
  inside active roots, collapse inactive branches. The change
  that surfaces "what's live now."

- `workstream-tracker-1-0-m3` (proposed). **Sub-stage cells and
  richer actor presence.** Per-leaf `D`/`P`/`I`/`V` cells colored
  by state (none / active / in-review / complete). Actor icons
  sit on specific cells rather than as bare chips on the node.
  Requires richer work-instance state vocabulary
  (`awaiting-user`, `awaiting-external`, `backgrounded`) in the
  API and DB.

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
drafts. Section added as a variance from
[`epic.md`](../../../spec/planning/epic.md)'s required+optional
list per the cross-level "child contracts at parent levels"
pattern; see PR body for the shared-spec edit backlog entry.

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

### m2 — Activity-first ordering

- **End result.** The forest renders with tier-based sorting:
  active surfaces top, in-flight middle, landed/abandoned
  compressed. Active-work paths inside active roots expand by
  default; inactive branches collapse.
- **Interfaces.** Consumes work-instance state from m1's
  multi-WI schema. Provides the visual ordering context that
  m3's per-cell rendering sits within.
- **Preserves.** m1's labels and detail rendering. The
  alphabetical ordering m1 ships with remains the fallback
  when no work-instance state distinguishes roots.

### m3 — Sub-stage cells and richer actor presence

- **End result.** Each leaf renders `D`/`P`/`I`/`V` cells
  colored by state (none / active / in-review / complete).
  Actor icons sit on specific cells. The work-instance state
  vocabulary expands to include `awaiting-user`,
  `awaiting-external`, and `backgrounded`.
- **Interfaces.** Consumes m1's multi-WI schema and m2's tier
  ordering. Adds richer work-instance states to the API and
  DB; agent rules from m1's t2 carry the new state-transition
  surface.
- **Preserves.** Existing render path for nodes without cell
  data; v0.1's actor-marker chip shape remains as the
  fallback when cells aren't populated.

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

- **Intent layer in 1.0?** The vision treats the intent strip
  and authoring nudges as core. Without them, the
  drift-prevention pillar is absent; with them, 1.0 grows by at
  least one substantial milestone (backlog parsing, intent-doc
  rendering, graduation flow). Open.
- **Triage zone in 1.0?** Distinct from the intent layer — the
  triage zone is where uncategorized work-instances sit until
  attached to a plan-tree node. Could land in 1.0 even if the
  broader intent layer doesn't, since it solves a narrower
  registration problem (exploratory sessions that don't yet
  know their slug).
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
- **Deferred open questions re-open after m1 lands.** If intent
  layer or triage zone become 1.0-blocking only after m1
  merges, the milestone set grows retroactively and m2/m3
  priorities may flip. Mitigation: revisit open questions
  explicitly at m1 retrospective, not passively.
- **Mid-epic spec changes break vendor consumers.** Additive
  posture is the invariant, but a milestone that doesn't
  realize it's introducing a breaking change is a real
  surface. Mitigation: spec-changing milestones carry an
  explicit "is this additive?" check in their planning.
- **Agent auto-registration (m1, t2) harder than estimated.**
  Reliable slug derivation from agent context involves design
  calls that may resist m1's task budget. Mitigation: m1's
  task list is locked but unsized; if t2 blows up, escalate
  to its own milestone before m1 ships.

## Sizing Summary

Per-milestone task counts. Estimates pending milestone
planning sessions for m2, m3, and the final milestone.

- **m1**: 4 tasks. Locked. See [`m1-v0-2.md`](m1-v0-2.md).
- **m2**: 2-4 tasks estimated. Render-layer tier classification
  and sort, plus collapse/expand affordances.
- **m3**: 3-5 tasks estimated. Schema and API for richer
  states, cell-level data parsing, cell renderer, actor-icon
  positioning.
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
