# Scoping — workstream-tracker-1-0-m1-t2 (Automatic agent registration)

Scoping doc for the second task of
[`m1`](../m1-v0-2.md) under the
[workstream-tracker-1-0 epic](../README.md). Pairs with the task
plan `docs/plans/workstream-tracker-1-0/m1-t2-auto-registration.md`.
Per [`task-plan.md`](../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns," this doc owns the deliberation, the
decisions with rejected alternatives, the reality-check inputs,
and the plan-structure handoff; it does not restate plan-owned
file inventory, contracts, or validation surface. No Status
field, per the same spec.

## Context

Nothing registers work-instances automatically today. The HTTP
API exists (`POST /work-instances`, with t1's create-or-attach
exact-slug flow, Landed), but an agent or human would have to
hand-assemble and fire that request, and nothing in `AGENTS.md`
or `docs/dev.md` even tells anyone to — so the cross-agent
visualization's actor markers are not trustworthy. t2 makes
registration happen as a consequence of starting work.

t2 ships **one lane**: the **interactive grounded narration
handshake**, the realization for the way the sole contributor
actually works (open an agent, state intent in natural
language). A deterministic "launcher" path — where something
that already knows the canonical slug invokes registration
before the agent runs — is **explicitly out of scope for t2**.
There is no current consumer for it, and the only compelling
future consumer (the tool's own UX spawning the planning agent
from a clicked plan-tree node) is a capability well beyond
v0.2. That future capability is captured in the backlog
(`tool-originated-task-sessions`), and deterministic
registration falls out of it for free — so there is no value in
building a launcher seam now for a consumer that does not exist.

## The load-bearing constraint: the registration circularity

This is why t2's lane is best-effort rather than a defect, and
why the deterministic path is a *future* capability, not a
v0.2 omission:

- Deterministic registration requires the canonical slug
  **before the agent runs** (a trigger that fires independent of
  agent cognition has nothing else to read).
- Resolving a natural-language request ("plan that task," "t2
  of m1 of my epic") into a canonical slug
  (`workstream-tracker-1-0-m1-t2`) requires interpreting intent
  against the plan tree — **agent cognition**, which happens
  *after* session start.
- Therefore, for a natural-language prompt, deterministic
  registration and "the human just says what they want" are
  **mutually exclusive**.

The circularity is dissolved only by **changing who the
launcher is**: if the tool's own UX is the launcher (the human
clicks "plan this task" on a specific plan-tree node and the
tool spawns the agent), the slug is known by construction
before any agent exists — the click carries the identity that
natural language cannot. That is the
`tool-originated-task-sessions` backlog opportunity, not v0.2
plumbing. Within t2's scope, the circularity holds, so t2's
honest answer is the observable best-effort handshake — which
is exactly the long-term vision's prescribed mitigation, not a
weaker form of a determinism that was never promised.

## Reality-check pass (load-bearing claims, verified)

- **No registration client or CLI exists.**
  `cmd/workstream-tracker/main.go` is server-only (router +
  serve; no subcommand, no client package). The cheap
  invocation the agent needs is net-new — t2 is not a
  narrow-surface change.
- **The vision frames this risk and prescribes exactly this
  mitigation.**
  [`design/vision.md`](../../../../design/vision.md) (the
  "registration depends on user discipline" feature-risk
  paragraph) explicitly accepts manual-for-some-sessions and
  prescribes: make the manual invocation cheap and
  habit-forming, and surface unregistered work so the
  contributor notices. The narration handshake operationalizes
  this faithfully. "Deterministic" appears nowhere upstream for
  registration; it was an analytical lens, never an upstream
  promise.
- **t1's flow + idempotency make repeated registration safe.**
  [`design/v0.1-design.md`](../../../../design/v0.1-design.md)
  §4 documents exact-slug create-or-attach and the
  `(slug, actor, state = active)` idempotency no-op; t1's
  `insertRegister` in `internal/api/handlers.go` implements it.
- **The slug is the authoritative identity, in plan-doc
  frontmatter.**
  [`spec/planning/shared.md`](../../../../spec/planning/shared.md)
  "Plan-doc identity (slug)". Resolving intent to it requires
  reading the plan tree — confirming the circularity's second
  premise and grounding the agent-resolves-the-slug step.
- **`docs/agents/shared/**` is vendored read-only.** `AGENTS.md`
  states the shared modules are vendored output, not edited
  directly; the repo-owned rule surfaces are `AGENTS.md` and
  `docs/agents/local/**`. This contradicts the milestone doc's
  Documentation-Currency line ("wire into
  `docs/agents/shared/`"); Decision 6 + the milestone-amendment
  record below correct it.
- **`actor` is a required free-form string.**
  `internal/api/api.go` `RegisterRequest.Actor` + the non-empty
  check in `internal/api/handlers.go`. Shapes Decision 3.

## Decisions made at scoping time

### Decision 1 — One lane (interactive grounded narration handshake); the deterministic launcher path is deferred, not built

Shapes considered:

- **1a — Single universal automatic mechanism (an instruction
  the agent follows at start).** Rejected: cannot cover the
  natural-language case deterministically (the circularity) and
  overclaims "reliable"; its silent failure is correlated and
  trust-collapsing.
- **1b — Ship a harness-specific session-start hook.** Rejected:
  harness-coupled and git-ignored, not vendorable.
- **1c — Two lanes (a deterministic launcher lane plus the
  interactive lane), both shipped in t2.** Rejected. The
  launcher lane has **no current consumer** (the sole
  contributor's workflow is interactive natural-language; the
  orchestrated sub-agent case is marginal and not worth a
  committed seam). Building a launcher-binding seam now is
  scope for a consumer that does not exist; the epic's "cut
  anything not on the path" applies. The compelling future
  consumer (tool-originated sessions) is a separate capability
  where determinism is free-by-construction — so the launcher
  path is *deferred to that backlog opportunity*, not built
  here.
- **1d — One lane: the interactive grounded narration
  handshake (chosen).** The agent resolves the slug
  mid-session and registers as an explicit, fact-grounded
  handshake before doing task work. Best-effort,
  agent-mediated, loud and checkable rather than silent.

Verified by: the circularity section;
[`design/vision.md`](../../../../design/vision.md) (mitigation
prescription); `cmd/workstream-tracker/main.go` (CLI net-new).

### Decision 2 — Slug source: the agent resolves it from prompt + plan tree and narrates it for human confirmation

The agent resolves the canonical slug from the prompt and the
plan tree (the unavoidable cognition step in this lane) and
**must narrate the resolved slug for human confirmation before
proceeding**. Unresolvable → narrate that explicitly and skip;
never a silent drop. Branch-name inference is rejected (observed
branches `impl/...`, `plans/...` carry the slug behind
non-uniform prefixes with no one-to-one mapping). Verified by:
[`spec/planning/shared.md`](../../../../spec/planning/shared.md)
"Plan-doc identity (slug)";
[`design/v0.1-design.md`](../../../../design/v0.1-design.md) §3.

### Decision 3 — `actor` is a supplied label, default a generated per-session id; never the git user

git user collapses all parallel agents into one marker,
defeating the cross-agent-visibility goal. The command accepts
an actor argument, defaulting to a generated per-session id
stable across a session's repeat invocations (so idempotency
collapses restarts); exact format is plan-owned. Verified by:
`internal/api/api.go` `RegisterRequest.Actor`;
`internal/api/handlers.go` non-empty check;
[`design/vision.md`](../../../../design/vision.md) (actor
distinguishes parallel work).

### Decision 4 — Best-effort, non-blocking, idempotent

One short-timeout attempt; on any failure the command logs and
exits success — never blocks or fails the session. The
narration *is* the surfacing in this lane. Repetition
(restart/resume) is safe because t1's `(slug, actor, active)`
idempotency collapses repeats. Verified by:
[`m1-v0-2.md`](../m1-v0-2.md) t2 "Preserves";
[`design/v0.1-design.md`](../../../../design/v0.1-design.md) §4
idempotency.

### Decision 5 — Narration must be grounded in real tool output, not a prose success claim

A prose "I have registered myself" is hallucinatable. Each
handshake line **echoes observed output**: the resolved doc
path, the resolved slug, and the *actual* command/server
response (real work-instance id and HTTP status). Failure is an explicit, actionable line, then proceed.
Ties to the `validation-honesty` audit. Verified by:
[`../../../agents/local/self-review-catalog.md`](../../../agents/local/self-review-catalog.md)
(`validation-honesty`); the silent-and-correlated failure
analysis (a present human catches a missing or incongruent
*fact*, not a prose claim).

### Decision 6 — The auto-registration rule lives in `AGENTS.md` + `docs/agents/local/`, not vendored shared

`docs/agents/shared/**` is vendored read-only, so the milestone
doc's "wire into `docs/agents/shared/`" is impossible. The
narration rule lives in `AGENTS.md` (a universal session rule)
and `docs/agents/local/**`. Verified by: `AGENTS.md` (shared
modules vendored; `docs/agents/local/**` is the repo-owned
rule-additions path).

### Decision 7 — Missing/failed registration is noticeable via the narration handshake, not the rendered tree

The residual silent-skip needs a backstop (the vision's second
mitigation: surface unregistered work so the contributor
notices). The achievable backstop in t2 is the **in-session
narration handshake** (Decision 5 / plan C2): a present
contributor sees the failure narration or the missing
handshake. The **rendered tree cannot** serve this — an
unregistered session emits no signal the tool ever sees, so the
tree cannot distinguish unregistered work from genuinely no
work. This is unobservable by construction, the same family of
constraint as the registration circularity, not a UI gap to be
closed in t2. A tree-side heuristic affordance (flag a node
with an active/Proposed plan doc but no work-instance) is a
genuine future capability, deferred and tracked by the
`deterministic-interactive-registration` backlog tripwire — it
is the same residual: the handshake backstop holds only while
the sole consumer is present to notice it, and that evaporates
at external adoption. Verified by:
[`design/vision.md`](../../../../design/vision.md) (surface
unregistered work so the contributor notices) and the
unobservable-by-construction argument above.

## On the spike (already run; no longer a required pre-plan gate)

The original two-lane scoping invoked "spike before plan for
novel mechanisms" for the launcher lane's session-start
trigger. With the launcher lane cut (Decision 1c), **there is no
novel mechanism left** — the lane that ships is instruction +
a thin HTTP-POST command, both established patterns. A throwaway
spike was nonetheless run during planning; its findings (the
registration command against a t1 server: server-up → HTTP 201
with a real id; server-down → log + exit success; repeat → same
id via idempotency; no-signal → register nothing) stand as
supporting evidence that the command's best-effort and
idempotent behavior is sound. They are recorded as evidence,
not as a gating spike the plan depended on.

## Milestone contract amendment required (executed in the t2 plan PR)

The milestone [`m1-v0-2.md`](../m1-v0-2.md) t2 "End result"
currently overclaims ("register on start without manual API
invocation … trustworthy because registration is reliable").
The circularity makes the unconditional form false. The t2 plan
PR amends it to the single-lane reality: registration happens
early in the session via a grounded narration handshake;
**without manual API invocation** is kept (true — the agent
invokes the command; no human hand-crafts an API call);
"reliable/trustworthy" → observable best-effort, with
deterministic registration deferred to the
`tool-originated-task-sessions` capability. Wording
reconciliation, not a contract gutting, mirroring t1's in-PR
milestone reconciliation. Because t2 drafting also resolves
questions the milestone still recorded as pending or pointed at
stale surfaces, the plan PR reconciles those in the same pass
(currency, not contract): the Task Status row, the t2
"Interfaces" bullet (slug signal was pending with a
now-rejected branch-name candidate), the Cross-Task Decisions
entry (slug derivation: deferred → RESOLVED), the Cross-Task
Risk "Slug derivation in t2 harder than estimated"
(→ RESOLVED, no escalation), and the Documentation-Currency
line that pointed t2 at vendored read-only
`docs/agents/shared/` (→ `AGENTS.md` + `docs/agents/local/`).
The `v0.1-schema` and `short_description` Cross-Task Risk
bullets are also stale relative to t1/t3 having Landed, but
that is those landed tasks' debt, not t2's, and is out of
scope here.

## Decisions handed off to plan-drafting

- The registration subcommand name and flag/env spellings; the
  per-session actor-id format and reuse mechanism (Contract +
  Naming).
- The exact narration line wording (UX copy — the *facts* each
  line conveys and the "echo real output / narrate failure"
  contract are fixed by Decision 5; phrasing is render/plan
  time, an authorized copy deferral).
- Decision-7 observability realization scope for t2 vs. a
  deferred richer affordance.
- N = 1 vs N ≥ 2 (recommend **N = 1**, inline): bounded surface
  (a `cmd` subcommand + thin client + the milestone amendment;
  agent-rule/design/dev edits deferred to the implementation
  PR).

## Reality-check inputs the plan must verify before promotion

- `cmd/workstream-tracker/main.go` still server-only (CLI
  net-new).
- t1's exact-slug create-or-attach + `(slug, actor, active)`
  idempotency intact in `internal/api`; confirm no regression.
- `docs/agents/shared/**` still vendored-read-only;
  `docs/agents/local/**` still the editable repo-owned path.
- [`design/vision.md`](../../../../design/vision.md)
  feature-risk paragraph wording (the milestone amendment
  references it).

## Plan-structure handoff

The plan doc (N = 1 recommended) will own, per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Required and optional sections":

- **Goal + context preamble** — this doc's Context is the seed.
- **Contracts** — the registration command (inputs;
  single-attempt best-effort; idempotent consumption of t1's
  exact-slug flow); the grounded narration handshake (resolve →
  confirm slug → echo real command/server output → narrate
  failure → proceed); each at contract altitude.
- **Files to touch** (estimate-labeled) — `cmd/workstream-tracker/`
  (new subcommand) + a thin internal HTTP client; the
  [`m1-v0-2.md`](../m1-v0-2.md) t2-contract amendment and the
  backlog reconciliation (executed in the plan PR);
  `AGENTS.md`/`docs/agents/local/`/`docs/dev.md`/
  `design/v0.1-design.md` §3/§4 named but DEFERRED to the t2
  implementation PR; `.claude/settings*` and `design/vision.md`
  intentionally not touched.
- **Cross-Cutting Invariants** — never block the session;
  consumer-agnostic (ship a command, no committed harness
  coupling); spec stays additive (t2 only consumes t1's flow);
  missing/failed registration is noticeable via the narration
  handshake, not the rendered tree (Decision 7).
- **Validation Gate** — `go build ./...`, `go vet ./...`,
  `go test ./...` (per [`docs/dev.md`](../../../dev.md) lines
  84-86), plus the interactive-lane exercise (narration states
  resolved doc + slug + real command/server response;
  induced failure → explicit failure narration + session
  proceeds; unresolvable → explicit skip narration), and the
  stale-neighbor-server validation hazard with its mitigation.
- **Self-Review Audits** — `validation-honesty`,
  `error-surfacing-user-mutations`.
- **Documentation Currency** — `AGENTS.md`, `docs/agents/local`,
  `docs/dev.md`, `design/v0.1-design.md` §3/§4, the milestone
  t2-contract amendment.
- **Backlog Impact** — this plan PR **shifts** the framing of
  `deterministic-interactive-registration` (it stays Open; t2
  ships the observable best-effort answer; the future
  determinism home becomes the new entry) and **adds**
  `tool-originated-task-sessions` (the tool-spawns-the-agent
  opportunity where deterministic registration is
  free-by-construction). Neither is graduated here.
- **Risk Register, Out of Scope, Related Docs** as content
  applies.

This scoping doc deletes in batch with sibling scoping docs at
the milestone-terminal PR per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns."
