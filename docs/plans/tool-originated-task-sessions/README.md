---
slug: tool-originated-task-sessions
Status: In draft
short_description: Tool originates plan-tree sessions; deterministic registration falls out
---

# Tool-Originated Task Sessions

## Purpose

Scope a post-1.0 arc in which the workstream-tracker's own UX
*originates* a planning or implementation session from a
plan-tree node, instead of only ever observing sessions a
contributor opened elsewhere. Because the tool is the launcher
and the originating action carries the node's identity, the
spawned session's canonical slug is known *by construction*
before any agent cognition — which dissolves the registration
circularity and makes work-instance registration deterministic
for these sessions for free.

**The product picture.** The contributor is looking at the plan
tree in the tool. They see a node — a task whose phases aren't
worked out yet, or a phase that's ready to build — and instead
of leaving the tool, opening an agent somewhere else, and
describing the work in prose, they act on the node directly.
The node offers two choices: **Begin planning** (the node's
scope isn't settled yet) or **Begin implementation** (the node
is ready to build). Choosing one starts an agent session.
Claude Code is the first and only launcher the epic wires
concretely; the origination path is built so a second agent is
an added adapter, not a rewrite. The session comes up *already
knowing which node it is* — its canonical slug rode in with the
action — so it registers itself back on the tool with no
narration handshake and no hand-fixed marker, and the
work-instance appears on the very node the contributor clicked,
the loop closed by construction. Registration is **two-phase**:
a *deterministic identity step* (the injected slug attaches the
work-instance to the right node, guaranteed, before any agent
cognition), then a *best-effort context enrichment* — once the
agent has gathered what it is about to do, it sends a follow-up,
non-deterministic update carrying things that inherently require
its cognition (branch, worktree, a human short name, intended
work). The enrichment feeds surfaces the tool already renders
(the session roster's unstructured-detail view, the node's short
label); a missed enrichment degrades *richness*, not
*correctness* — the node still shows the agent is there. It also
comes up *already working*: a planning session opens on an
investigation/scoping prompt for that node; an implementation
session opens on an implementation prompt for it. The contributor's experience
collapses from "leave the tool, find an agent, describe the
work, hope it attributes itself correctly" to "click the node,
pick the mode, watch the right session start and attach
itself."

This doc's first job is to **settle the product vision and
technical direction before any milestone scope locks**. The
tool moving from observe-only to launching agents is the
largest posture change in the product's history; the milestone
breakdown is an estimate until the vision questions below are
resolved. Working surface across the epic's lifecycle: sessions
promote the Open Questions into committed milestone scope (or to
Out of Scope) as deliberation locks them.

## Why This Epic

Two backlog entries motivate this epic, and they are not
co-equal:

- [`tool-originated-task-sessions`](../../backlog.md#tool-originated-task-sessions)
  is the deliverable capability — a "plan this task" / "work
  this task" affordance on a plan-tree node that spawns the
  agent session itself.
- [`deterministic-interactive-registration`](../../backlog.md#deterministic-interactive-registration)
  is a *property that falls out of* that capability, not a
  parallel track. Determinism for interactive,
  natural-language sessions "provably cannot be" achieved: the
  registration circularity is structural — resolving
  natural-language intent to a canonical slug requires agent
  cognition, which postdates session start, while a
  deterministic trigger must run before it. The circularity is
  dissolved *only* by changing who the launcher is. (Verified
  by [`docs/backlog.md` →
  `deterministic-interactive-registration`](../../backlog.md#deterministic-interactive-registration).)

So this epic delivers tool-originated sessions; deterministic
registration is a resolved consequence of that delivery, not a
separable earlier deliverable. The backlog entries say this to
each other:
[`deterministic-interactive-registration`](../../backlog.md#deterministic-interactive-registration)
names tool-origination as its "only resolution home," and
[`tool-originated-task-sessions`](../../backlog.md#tool-originated-task-sessions)
names itself "the home where the determinism deferred by
`deterministic-interactive-registration` is eventually
achieved." (Verified by
[`docs/backlog.md`](../../backlog.md#tool-originated-task-sessions).)

Both entries are framed "well beyond v0.2 / beyond the 1.0
epic's scope." This epic is therefore post-1.0 and successor to
the
[`workstream-tracker-1-0`](../workstream-tracker-1-0/README.md)
epic for this capability — not a milestone inside it.

## Goal

Ship the affordance by which the contributor originates a
plan-tree session from a node in the tool — choosing **Begin
planning** or **Begin implementation** — such that an agent
session starts, registers itself back on the node it was
launched from (canonical slug carried by construction;
deterministic; no slug resolution, no narration handshake), and
opens already working from a scoping prompt or an implementation
prompt respectively. Claude Code is the concrete first launcher;
adding another agent is an adapter, not a rewrite. The vision
and technical direction for *how the tool spawns a session* —
and how that reconciles with the product's standing observe-only
posture — are settled in this doc before milestone scope locks.

## Inherited Context

External decisions that bound this epic's scope. They are not
re-opened here; the epic's vision work reconciles against them.

- **The 1.0 epic's observe-only and single-contributor
  invariants.** "Tool stays one-way" (server reads agent
  events, never pushes back, no tool→agent injection) and
  "single contributor, single local environment" (no auth, no
  multi-tenant, no remote hosting) bind the product through
  1.0. This epic is the deliberate, bounded, post-1.0 revisit
  of the *first* of those; it does **not** revisit the second.
  (Verified by [`workstream-tracker-1-0` Cross-Cutting
  Invariants](../workstream-tracker-1-0/README.md) and
  [`design/vision.md` §5](../../../design/vision.md).)
- **The locked m1-t2 best-effort registration decision.**
  Interactive (contributor-opened) registration ships as an
  observable best-effort grounded narration handshake, and
  "deterministic" was never an upstream promise for it; the
  registration circularity proves interactive registration
  cannot be deterministic. This epic does not contradict that
  decision — it changes who the launcher is for a *different*
  class of session. (Verified by [`docs/backlog.md` →
  `deterministic-interactive-registration`](../../backlog.md#deterministic-interactive-registration).)
- **The spec's exact-slug create-or-attach registration path.**
  The spec already defines registration "at exactly this slug"
  with server-side trust of the caller's slug. The
  construction-known slug a tool-originated session carries is
  the natural producer for that path. (Verified by
  [`shared.md` "Slug generation"](../../../spec/planning/shared.md).)
- **The session-presence vision frame.** `design/vision.md`
  §4's best-effort "session presence" concept and §5's
  observe-only posture are the vision surface this epic's
  origination action narrowly extends. (Verified by
  [`design/vision.md` §4, §5](../../../design/vision.md).)

## Cross-Cutting Invariants

Rules that bind every milestone in this epic. Reviewer-flag
candidates for any milestone deliberation that brushes against
them.

- **This epic deliberately crosses the "tool stays one-way"
  line, narrowly.** The
  [`workstream-tracker-1-0`](../workstream-tracker-1-0/README.md)
  epic's invariant — server reads agent events, never pushes
  back, no tool→agent injection — is revisited *here*,
  post-1.0, by design. The crossing is bounded to session
  *origination* from an explicit contributor action on a node;
  it is not general tool→agent injection, drift-correction, or
  in-session steering. (Verified by
  [`workstream-tracker-1-0` Cross-Cutting
  Invariants](../workstream-tracker-1-0/README.md) and
  [`design/vision.md` §5](../../../design/vision.md).)
- **Determinism is a property of tool-origination, not a
  separable earlier deliverable.** No milestone delivers
  "deterministic interactive registration" — that phrase
  contradicts the locked backlog decision that interactive
  registration provably cannot be deterministic. Determinism
  for real sessions exists only once the tool is the launcher.
- **Slug-by-construction is the whole mechanism.** The
  dissolution of the registration circularity rests entirely
  on the launcher carrying the node's canonical slug into the
  spawned session. Any milestone that weakens "the originating
  action carries the node's canonical slug" removes the epic's
  reason to exist.
- **What "deterministic" is scoped to.** Two boundaries, one
  word. (a) *Given a launched session, not the launch
  succeeding* — the guarantee is that a session which *did*
  start registers itself; the spawn itself is best-effort and a
  failed or cancelled launch must be *observable*, not silently
  swallowed (the project's accepted-failure-with-visibility
  posture). (b) *Identity, not enrichment* — registration is
  two-phase: the deterministic identity step binds the
  work-instance to the right node from the injected slug before
  any agent cognition, while the follow-up context enrichment
  (branch, worktree, short name, intended work) inherently
  requires cognition and is explicitly best-effort; a missed
  enrichment degrades richness, not correctness. No milestone
  may read "deterministic" as "the spawn cannot fail" or make
  enrichment load-bearing for attachment. The best-effort
  follow-up does **not** reintroduce the registration
  circularity — identity is already established before cognition,
  so enrichment only decorates an already-correctly-attached
  instance — which is also why tool-originated sessions carry no
  interactive-path observability residual: attachment is
  guaranteed whenever the session launched.
- **Origination stays in the contributor's single local
  environment.** The spawned session runs where the contributor
  already runs agents; the tool launching a session introduces
  no hosting, daemon-for-others, or multi-tenant surface. The
  1.0-epic single-contributor/single-local invariant is carried
  forward unchanged. (Verified by [`workstream-tracker-1-0`
  Cross-Cutting Invariants](../workstream-tracker-1-0/README.md).)
- **Agent-adapter seam: Claude Code is the first launcher, not
  the only conceivable one.** The epic wires Claude Code
  concretely and may hard-code it for v1, but the origination
  and slug-carry path must not bake in Claude-Code-only
  assumptions that would make a second agent a rewrite rather
  than an added adapter. This is a forward-compatibility
  posture, not a commitment to ship a second agent in this
  epic.
- **This epic neither resolves nor cancels the 1.0-epic
  interactive-registration tripwire.** The
  [`deterministic-interactive-registration`](../../backlog.md#deterministic-interactive-registration)
  tripwire — re-deliberate best-effort registration before the
  1.0 epic's neighborly-events integration milestone — binds
  independently and *earlier* than this post-1.0 epic ships.
  This epic is the eventual determinism home; it does not pull
  determinism forward to satisfy that tripwire, and the
  tripwire may pull a different backstop forward on its own
  schedule. (Verified by
  [`docs/backlog.md` →
  `deterministic-interactive-registration`](../../backlog.md#deterministic-interactive-registration)
  and [`workstream-tracker-1-0` Risk
  Register](../workstream-tracker-1-0/README.md).)

## Milestone Structure

*Estimate pending milestone planning. The milestone count and
split are not locked: this is a vision-first epic, and the
technical-direction Open Questions below must resolve before
the UX milestone's scope can lock. Per
[`epic.md`](../../../spec/planning/epic.md) "Scope," these
per-milestone details are estimates, not binding specs, until
each milestone's planning session re-derives them against
merged code.*

**Sequencing rationale** — sequential dependency, not parallel:

- **m1 first** because the deterministic registration path can
  be built and proven independently of the tool-acting posture
  shift, using a manual slug producer. De-risking the
  registration contract first keeps the larger vision question
  from being entangled with a registration-mechanics question.
- **The UX milestone second** because it is the larger
  vision/technical-direction question (how the tool spawns a
  session; how that reconciles with the observe-only posture)
  and it consumes m1's proven slug-carried registration path
  as the thing the spawn wires into.

**Milestones (estimated):**

- `tool-originated-task-sessions-m1` (estimate — scope not yet
  locked). **Unwired deterministic slug-carried registration
  path.** A registration entrypoint that, given a
  construction-known canonical slug, registers a work-instance
  deterministically — no natural-language resolution, no
  narration handshake. Exercised via an explicit slug argument
  (a manual / CLI invocation) as the stand-in slug producer; no
  agent-spawning UX. Proves the deterministic path end-to-end
  before the tool-acting posture shift.

- **Tool-originated session UX milestone(s)** (estimate — scope
  not yet locked, gated on the Open Questions below). The
  plan-tree-node affordance that spawns the agent session and
  becomes the *real* construction-time slug producer feeding
  m1's path. This is the actual
  [`tool-originated-task-sessions`](../../backlog.md#tool-originated-task-sessions)
  capability and the determinism-resolution home. Milestone
  count and structure are the milestone-planning sessions'
  output, not fixed here, and cannot lock until the one
  remaining scope-locking Open Question resolves — the
  spawn-shape choice. (One-way-invariant reconciliation,
  workspace origin, and the level/mode matrix are resolved; see
  Open Questions Resolved By This Epic.)

## Milestone Contracts

Per-milestone **WHAT** contracts. The **HOW** for each
milestone lives in the milestone doc when it drafts. Required
section per [`epic.md`](../../../spec/planning/epic.md)
"Required and optional sections" and
[`shared.md`](../../../spec/planning/shared.md) "Parent-doc
child contracts." Scope is **not yet locked** for any milestone
in this epic — it is vision-first by construction; per
[`shared.md`](../../../spec/planning/shared.md) "Parent-doc
child contracts," this names the milestones without sealing
their contracts until their planning sessions run and the
Open Questions resolve.

| Milestone | Short description | End result and what it preserves (WHAT, estimated) | Sibling interface (estimated) |
|---|---|---|---|
| `tool-originated-task-sessions-m1` | Unwired deterministic slug-carried registration path | Given a construction-known canonical slug, a work-instance registers deterministically with no resolution and no narration handshake; exercised via an explicit slug argument. Preserves the existing interactive best-effort registration path unchanged (this is an additional path, not a replacement) and the spec's existing exact-slug create-or-attach posture. | Produces the slug-carried registration entrypoint the UX milestone's spawn wires into as the construction-time slug producer. |
| Tool-originated session UX milestone(s) | The node affordance that spawns a session and carries its slug by construction | *Scope not yet locked.* When done, an explicit contributor action on a plan-tree node spawns an agent session whose canonical slug is carried by construction, registered deterministically via m1's path. Preserves the observe-only behavior for all sessions not originated by the tool. | Consumes m1's slug-carried registration entrypoint as the real construction-time slug producer. |

## Open Questions Resolved By This Epic

Long-standing deferrals this epic's existence settles.

- **Where does registration determinism get achieved?**
  Resolved: *here*. The backlog and the 1.0 epic deferred
  determinism with "its only resolution home is
  tool-origination"; this epic is that home. The determinism
  gap stops being an open deferral and becomes a scoped
  deliverable. (Verified by [`docs/backlog.md` →
  `deterministic-interactive-registration`](../../backlog.md#deterministic-interactive-registration).)
- **How does tool-origination reconcile with the "tool stays
  one-way" vision invariant?** Resolved: *the rule stands; this
  epic adds one bounded exception, and the exception does not
  erode the rule's value.*
  - **The rule's value.** The one-way invariant exists so the
    tool never becomes a control authority over agents — no
    tool→agent feedback loop — which keeps the project's files
    the source of intent and keeps agents autonomous observed
    actors rather than tool-driven ones.
  - **The exception.** The tool gains exactly one origination
    action: at an explicit contributor action on a node, it
    spawns a session and hands it the node's slug and an
    opening prompt.
  - **Why the value survives.** The exception is fenced by four
    properties — it is *human-initiated* (not the tool acting
    autonomously), *one-shot at session birth* (not
    continuous), carries *identity, not behavioral correction*,
    and leaves the *running session fully file-driven and
    observe-only*. The two-phase split reinforces this: even
    registration and enrichment are agent→tool, so nothing
    about a *running* session ever flows tool→agent. Strip any
    one fencing property and the value would erode — which is
    why those properties are load-bearing, not decorative.
  - **What stays forbidden.** Steering, correcting, or
    injecting into a *running* session (the tool noticing
    drift and pushing a fix, re-prioritizing mid-flight,
    feeding tool state back into agent context). That fails the
    fencing and remains out of scope; drift prevention still
    runs through intent in the files.
  - **Authoritative home.** The invariant's authoritative
    source is [`design/vision.md`
    §5](../../../design/vision.md); adopting this epic implies
    the vision's one-way framing eventually carries this named
    carve-out. The epic flags that rather than letting vision
    §5 read as an absolute it silently contradicts. (Verified
    by [`workstream-tracker-1-0` Cross-Cutting
    Invariants](../workstream-tracker-1-0/README.md) and
    [`design/vision.md` §5](../../../design/vision.md).)

- **Where does a tool-originated session's workspace come
  from?** Dissolved: *the tool has no workspace role.* The
  opening prompt tells the session to start its work in a fresh
  worktree off the latest `origin/main`; the agent does so as
  ordinary agent behavior, and the resulting worktree name comes
  back as display-only best-effort enrichment. This is not a
  tool decision among options — it is already fully covered by
  the launch-not-leash invariant (the tool's only action is
  hand over the prompt) and the two-phase split (worktree is
  agent-reported enrichment, not a tool concern). The tool never
  provisions, selects, or acts on the workspace.
- **Which node levels offer "Begin planning" vs "Begin
  implementation," and how does the mode map to the session
  that starts?** Resolved minimally. The mode *is* the prompt:
  "Begin planning" hands over a scoping/investigation prompt
  for the node, "Begin implementation" hands over an
  implementation prompt — the tool decides nothing about how
  the session then plans or implements (launch-not-leash). Which
  mode(s) a node offers is a **shallow static map over two
  values the tool already holds** — the node-type it already
  walks the tree to render, and the plan-doc Status it already
  parses to color the node: parents (epic/milestone) are never
  "implement" targets (their work is their children); a node
  with no doc or `In draft` offers planning; `Proposed` offers
  implementation; `In progress` / `Validating` / `Landed` offer
  neither. This is a lookup table, not inference. **Bound:** if
  a milestone finds it needs deep logic to pick the button,
  that is a signal to revisit this resolution, not to add
  inference power to the tool — the minimal-tool posture is the
  invariant, the table is the means.

## Open Questions Newly Opened

These are the vision and technical-direction calls the epic
must resolve before the UX milestone's scope can lock. They are
the substance of the "settle vision first" framing.

- **How does the tool spawn an agent session?** The core
  technical-direction question — the tool today only observes;
  "the tool acting / spawning agents" is a fundamentally new
  capability surface. A first investigation (local-only, per
  the single-local invariant) decomposed the candidate space
  into three shapes, to be analyzed at the UX milestone's
  planning session, not locked here:
  - **Shape A — spawn the interactive Claude Code CLI; inject
    the slug out-of-band; register via a session-start hook.**
    The tool launches a real interactive session the
    contributor can take over; the launcher passes the
    construction-known slug through an out-of-band channel
    (environment or file); a session-start hook runs the
    deterministic registration before the model reasons. This
    is the **leading candidate**: it is local, keeps the
    session interactive, makes registration deterministic by
    construction, and reuses the existing
    [`session-registration.md`](../../../docs/agents/local/session-registration.md)
    model with slug *resolution* replaced by an injected known
    slug and the narrated best-effort step replaced by the
    hook. (Verified by [Claude Code hooks
    documentation](https://code.claude.com/docs/en/hooks.md)
    for the session-start hook surface.)
  - **Shape B — headless / print-mode invocation with the slug
    in the prompt.** Fully scriptable but not a session the
    contributor can take over and work in; weaker
    register-before-model-work guarantee. Fails the product
    vision of a takeable working session, though plausibly
    useful as m1's unwired test-harness producer.
  - **Shape C — Agent SDK / cloud managed agents.** A
    different, cloud product surface with no session-start-hook
    equivalent; out of local-first focus and in tension with
    the single-local invariant. Recorded as the **future cloud
    adapter** the agent-adapter-seam invariant keeps the door
    open for — explicitly not committed by this epic. (Verified
    by [Claude Agent SDK / managed-agents
    documentation](https://platform.claude.com/docs/en/managed-agents/quickstart.md).)

  Exact CLI flags and hook-event spellings are HOW: deferred to
  the UX milestone's planning session against the then-current
  Claude Code, and require real-doc verification before they
  become load-bearing.
- **Does m1's slug-carried path reuse the existing exact-slug
  create-or-attach registration path, or introduce new
  surface?** The spec already defines an exact-slug
  create-or-attach registration path ("register at exactly
  this slug"). (Verified by
  [`shared.md` "Slug generation"](../../../spec/planning/shared.md).)
  Whether m1 is a thin slug-passing invocation over that path
  or new registration surface is a technical-direction call for
  m1's milestone-planning session against merged code.
- **How do the spec and agent rules express a deterministic,
  handshake-free registration path?** A tool-originated session
  skips the best-effort grounded narration handshake because
  its slug is known by construction. The plan-doc spec and the
  session-registration agent rule currently describe only the
  best-effort interactive path; they will need an additive
  expression of the deterministic path. Downstream of vision,
  but flagged so spec coherence is a conscious milestone input.
- **What is the relationship between this epic and the 1.0-epic
  interactive-registration tripwire?** Stated as a
  cross-cutting invariant above (independent; not resolved
  here). Flagged here so the call stays conscious through the
  epic's lifecycle rather than being silently absorbed.

## Out of Scope

- **General tool→agent communication.** This epic crosses the
  one-way line only for session *origination*. In-session
  steering, drift-correction, or any tool→agent push beyond the
  spawn is out of scope and remains governed by the standing
  observe-only posture.
- **Resolving or cancelling the 1.0-epic
  interactive-registration tripwire.** Per the cross-cutting
  invariant, that obligation binds independently and earlier;
  this epic does not satisfy it and does not pull determinism
  forward to meet it.
- **Changing the interactive (contributor-opened) registration
  path.** Sessions a contributor opens outside the tool keep
  the observable best-effort grounded narration handshake; this
  epic adds a deterministic path for tool-originated sessions,
  it does not replace the interactive one.
- **Shipping a second agent.** Per the agent-adapter-seam
  invariant, the epic keeps the seam open but commits only
  Claude Code; delivering an additional agent adapter is
  post-epic.
- **The exact scoping / implementation prompt content.** That
  a planning session opens on a scoping prompt and an
  implementation session on an implementation prompt is in
  scope; the specific prompt wording is implementation, settled
  in the implementing PR, not contracted here.
- **The tree-side observability heuristic.** A tree-side
  "node with an active plan doc but no work-instance" heuristic
  is part of the
  [`deterministic-interactive-registration`](../../backlog.md#deterministic-interactive-registration)
  residual that stays in the backlog under the split below; it
  is not a deliverable of this epic.

## Backlog Impact

Per the effect taxonomy in
[`spec/backlog.md`](../../../spec/backlog.md):

- [`tool-originated-task-sessions`](../../backlog.md#tool-originated-task-sessions)
  — **graduate.** This epic is the plan-tree node that carries
  it; the entry's Status flips to
  `Graduated — tool-originated-task-sessions` with a `**Plan:**`
  line pointing here.
- [`deterministic-interactive-registration`](../../backlog.md#deterministic-interactive-registration)
  — **split.** The entry tracks two separable threads. (1) The
  *determinism-resolution* thread — determinism's "only
  resolution home" — graduates into this epic (the epic node
  now carries it; no separate Status flip on the entry for this
  thread, because the entry is not wholly graduated). (2) The
  *interactive-registration best-effort tripwire and
  observability residual* — re-deliberate best-effort before
  the 1.0 epic's neighborly-events integration milestone; the
  tree cannot distinguish unregistered work from no work — binds
  independently of this post-1.0 epic and *earlier* than it
  ships. That thread **stays Open under the same
  `deterministic-interactive-registration` slug**, with its body
  narrowed to the surviving obligation and a note that the
  determinism-resolution thread is now homed here. The slug is
  preserved deliberately: multiple 1.0-epic plan docs reference
  it by anchor, and the independent tripwire must not be
  silently absorbed into this far-future epic.

## Risk Register

- **Tool-acting posture is a vision-level shift, not just a
  feature.** Moving the tool from observe-only to launching
  agents is the largest conceptual change in the product's
  history and tensions with a 1.0-epic cross-cutting invariant.
  Mitigation: the epic resolves the one-way-invariant
  reconciliation in-doc (rule stands; one bounded, fenced
  exception; value preserved; vision §5 named as the carve-out's
  authoritative home) before any milestone scope locks —
  vision-first by construction, not by convention.
- **Determinism mis-scoped as an early deliverable.**
  Re-introducing "deterministic interactive registration" as an
  early milestone would contradict the locked backlog decision.
  Mitigation: the cross-cutting invariant pins determinism as a
  property of tool-origination; m1 is explicitly the
  slug-carried path with a *manual* slug producer, never
  "deterministic interactive registration."
- **The independent tripwire gets absorbed and lost.** Marking
  the whole
  [`deterministic-interactive-registration`](../../backlog.md#deterministic-interactive-registration)
  entry Graduated would silently move the earlier-binding
  tripwire into this far-future epic, where it would not fire
  in time. Mitigation: Backlog Impact records a split — the
  tripwire/observability residual stays Open under the same
  slug with inbound anchors preserved.
- **m1 builds throwaway scaffolding.** The explicit-slug
  invocation that exercises m1's path may be discarded once the
  UX producer lands. Mitigation: scope m1's manual producer at
  m1 milestone planning as the minimum to validate the
  registration contract (a test harness), not a durable UX.

## Sizing Summary

*Estimates only; nothing locked. Per-milestone task counts are
each milestone-planning session's output.*

- **m1**: estimate pending the m1 milestone-planning session.
  A slug-carried deterministic registration entrypoint plus a
  manual slug producer to exercise it.
- **Tool-originated session UX milestone(s)**: count and split
  are an estimate that cannot firm up until the one remaining
  scope-locking Open Question resolves — the spawn-shape choice.
  (One-way reconciliation, workspace origin, and the level/mode
  matrix are resolved.)

## Related Docs

- [`../../backlog.md`](../../backlog.md) — the two motivating
  entries; this epic graduates one and splits the other.
- [`../workstream-tracker-1-0/README.md`](../workstream-tracker-1-0/README.md)
  — the 1.0 epic; this epic is its post-1.0 successor for this
  capability and inherits the "tool stays one-way" invariant it
  deliberately revisits, plus the interactive-registration
  tripwire it must not absorb.
- [`../../../design/vision.md`](../../../design/vision.md) —
  §4 (session presence) and §5 (the observe-only posture this
  epic narrowly crosses); the epic's vision work reconciles
  against it.
- [`../../../spec/planning/epic.md`](../../../spec/planning/epic.md)
  — the rules this epic doc is structured against.
- [`../../../spec/backlog.md`](../../../spec/backlog.md) — the
  graduate/split effect taxonomy this epic's Backlog Impact
  follows.
- [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md)
  — the exact-slug create-or-attach registration path m1's
  Open Question references, and the parent-doc child-contract
  rule this doc's Milestone Contracts follows.
- [`../../../docs/agents/local/session-registration.md`](../../../docs/agents/local/session-registration.md)
  — the best-effort grounded narration handshake the
  interactive path uses; tool-originated sessions diverge from
  it (slug known by construction), the divergence this epic's
  spec-coherence Open Question must express.
