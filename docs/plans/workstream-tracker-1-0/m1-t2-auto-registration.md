---
slug: workstream-tracker-1-0-m1-t2
Status: Landed
short_description: Automatic agent registration
---

# Task — Automatic agent registration (workstream-tracker-1-0-m1-t2)

Task plan for the second task of
[`m1`](m1-v0-2.md) under the
[workstream-tracker-1-0 epic](README.md). Paired scoping doc:
[`scoping/m1-t2-auto-registration.md`](scoping/m1-t2-auto-registration.md)
(authoritative for the deliberation, rejected alternatives, and
the reality-check inputs; this plan owns the durable contract,
file inventory, validation surface, and risks). N = 1 — one
phase, content absorbed inline; no separate phase plan files,
per the scoping doc's plan-structure handoff and
[`task-plan.md`](../../../spec/planning/task-plan.md) "N = 1 task
plan: phase content absorbed inline."

## Context

Today nothing registers work-instances automatically. The HTTP
API to register one exists — t1 landed the exact-slug
create-or-attach flow — but an agent or a human would have to
hand-assemble and fire that HTTP request, and no repo
instruction even tells anyone to. The practical consequence:
the cross-agent visualization's actor markers are not
trustworthy, because whether a session shows up depends entirely
on someone remembering to register it by hand. This task makes
registration happen as a *consequence of starting work* rather
than as a thing a contributor must remember.

This is being done now because t1 (the multi-work-instance
create-or-attach API) just landed and unblocked it, and because
the milestone's credibility goal — "make the agent dogfood
credible" — rests on markers actually appearing.

A naive "every session registers automatically and silently on
start" is not achievable, and the reason is a logical
constraint, not an implementation gap: a trigger that fires
before the agent thinks has nothing to read except its launch
environment, while turning a natural-language request ("plan
that task") into a canonical slug is interpretation that happens
*after* the session starts (the "registration circularity," in
the scoping doc). For the way the sole contributor actually
works — open an agent, state intent in natural language — the
honest mechanism is therefore an **observable best-effort
grounded narration handshake**: the agent resolves the slug
mid-session and registers as a loud, fact-grounded handshake
rather than a silent best-effort. This is a faithful
operationalization of the long-term vision's prescribed
mitigation (make the manual invocation cheap, surface the gap),
not a weaker form of a determinism that was never promised
upstream.

A deterministic path — where something that already knows the
slug invokes registration before the agent runs — is
**explicitly out of scope** for this task: it has no current
consumer, and the circularity is only dissolved by changing who
the launcher is (the tool's own UX spawning the agent from a
clicked plan-tree node), which is a future capability captured
in the backlog (`tool-originated-task-sessions`), where
deterministic registration is free-by-construction. Building a
launcher seam now for a consumer that does not exist is scope
this task deliberately cuts.

The surfaces this task touches: a new command-line entry point
on the existing tool binary plus a thin internal HTTP client;
the milestone-doc contract reconciliation; and the backlog
reconciliation. The agent-rule, design, and contributor-doc
edits are named here and deferred to the implementing PR.

## Goal

Make work-instance registration a consequence of starting work
in an interactive natural-language session, via an observable
best-effort grounded narration handshake, without ever blocking
or slowing the session, and without coupling the tool to any one
agent harness. Registration consumes t1's exact-slug
create-or-attach flow unchanged; this task adds no API, schema,
or frontmatter change.

## Contracts

The contract is the registration command-line entry point and
the interactive grounded narration handshake. Each is stated at
contract altitude — the condition that must hold, not the
implementer's technique for making it hold.

### C1 — The registration command

- A subcommand on the existing tool binary performs one
  work-instance registration against a running server's
  `POST /work-instances`, consuming t1's exact-slug
  create-or-attach flow. Verified by:
  [`internal/api/handlers.go:58-62`](../../../internal/api/handlers.go)
  (non-empty `exact_slug` selects create-or-attach and bypasses
  the root/descendant requirements) and
  [`design/v0.1-design.md`](../../../design/v0.1-design.md) §4
  exact-slug case.
- **Inputs.** The canonical slug, the actor label, and the
  server base location, each supplied to the command as an
  explicit argument or environment variable. The command
  performs no plan-tree reading and no intent resolution — it
  is a thin, mechanical HTTP-POST of inputs it is handed (the
  agent resolves the slug and hands it in; see C2). The exact
  flag/environment-variable spellings are named in the Naming
  section.
- **Actor.** The command accepts an actor label and defaults to
  a generated per-session identifier; it must never default to
  the git user, because collapsing parallel agents onto one
  actor defeats the cross-agent-visibility goal. Verified by:
  [`internal/api/api.go`](../../../internal/api/api.go)
  `RegisterRequest.Actor` (free string) and
  [`internal/api/handlers.go:51-54`](../../../internal/api/handlers.go)
  (actor required, non-empty). The per-session id format and
  its reuse mechanism are named in the Naming section.
- **Single-attempt best-effort, non-blocking.** The command
  makes one short-timeout attempt. On any failure — server
  unreachable, non-success HTTP status, malformed response — it
  logs an explicit line and exits success; it never blocks,
  retries in a loop, or fails the session. Verified by:
  [`m1-v0-2.md`](m1-v0-2.md) t2 "Preserves" (server
  unreachability tolerated; session proceeds) and the planning
  spike, which exercised server-up (HTTP 201), server-down, and
  no-input paths and confirmed exit-success in every case.
- **Idempotent consumption.** Re-running the command for the
  same slug and actor while a prior registration is still active
  collapses to the existing work-instance (no duplicate row, no
  duplicate event), so restart/resume is safe. Verified by:
  [`internal/api/handlers.go:226-234`](../../../internal/api/handlers.go)
  (the `(slug, actor, state = active)` idempotency check) and
  [`design/v0.1-design.md`](../../../design/v0.1-design.md) §4
  idempotency; the spike's repeat scenario returned the same
  work-instance id.
- **First-registration safe.** When the slug has no prior
  work-instance, the command creates the slug's first one (the
  freshly-drafted-node case). Verified by:
  [`design/v0.1-design.md`](../../../design/v0.1-design.md) §4
  ("no precondition that a work-instance already exists for the
  slug") and the locked t1↔t2 cross-task decision in
  [`m1-v0-2.md`](m1-v0-2.md) "Cross-Task Decisions."

### C2 — The grounded narration handshake

For an interactive natural-language session, the agent performs
registration as an explicit, human-visible handshake **before
doing task work**, and every line of the handshake reports
observed fact, not a prose success claim:

- **Resolve.** The agent resolves the canonical slug from the
  prompt and the plan tree (the unavoidable cognition step,
  given the registration circularity). Verified by:
  [`spec/planning/shared.md`](../../../spec/planning/shared.md)
  "Plan-doc identity (slug)" (the slug is the authoritative
  identity, carried in plan-doc frontmatter) and
  [`design/v0.1-design.md`](../../../design/v0.1-design.md) §3
  (agents read plan files including the slug frontmatter).
- **Confirm.** The agent states the resolved plan-doc path and
  the resolved slug for the human to confirm before proceeding.
- **Echo real output.** The agent invokes the registration
  command and reports the *actual* command/server response — the
  real work-instance identifier and HTTP status it observed —
  not a paraphrased "registered successfully." A prose success
  claim with no observed output is a contract violation, because
  a loosely-followed instruction produces confident fake
  receipts; the handshake's value is that a present human can
  catch a missing or incongruent *fact*. This ties directly to
  the `validation-honesty` audit. Verified by:
  [`docs/agents/local/self-review-catalog.md`](../../../docs/agents/local/self-review-catalog.md)
  (`validation-honesty`: a claim that a check ran is only valid
  if it ran end-to-end on the current state).
- **Narrate failure explicitly.** If registration fails or the
  slug is unresolvable, the agent states that explicitly and
  actionably (the session will not appear in the tree; the
  human can run the registration command by hand), then
  proceeds with task work. A silent skip is a contract
  violation. This ties to the `error-surfacing-user-mutations`
  audit. Verified by:
  [`docs/agents/local/self-review-catalog.md`](../../../docs/agents/local/self-review-catalog.md)
  (`error-surfacing-user-mutations`: silent failure is the
  trap).
- **Proceed.** Registration never blocks task work; whether it
  succeeded, failed, or was skipped, the session continues after
  the handshake.

This lane is **observable best-effort**, not deterministic.
That is faithful to the long-term vision, which explicitly
accepts manual-for-some-sessions and prescribes "make the
manual invocation cheap and habit-forming … and surface
unregistered work somehow." Verified by:
[`design/vision.md:171`](../../../design/vision.md) (the
"registration depends on user discipline for some kinds of
sessions" feature-risk paragraph and its stated mitigation
direction). The exact narration line wording is render-time UX
copy; this contract fixes the *facts* each line must convey and
the echo-real-output / narrate-failure rule, and authorizes the
phrasing deferral per the structural-surface and copy-deferral
conventions.

## Cross-Cutting Invariants

Rules that must hold simultaneously across the command, the
narration rule, the agent-rule docs, and the validation
surface — each breaks silently if only one site honors it.

- **Never block or slow the session.** Every registration path
  (success, failure, unresolvable, server-down) ends with the
  session proceeding; registration is always a short-timeout
  single attempt with no blocking retry.
- **Consumer-agnostic; no committed harness coupling.** The
  tool's committed contract is the command plus the narration
  rule. No harness-specific trigger, hook, or settings file is a
  committed tool artifact.
- **Spec stays additive — t2 only consumes t1's flow.** This
  task adds no endpoint, no request/response field, no schema
  change, and no plan-doc frontmatter change. It is a pure
  consumer of t1's exact-slug create-or-attach flow. Verified
  by: [`m1-v0-2.md`](m1-v0-2.md) "Cross-Task Invariants"
  (backward-compatible API extension; spec changes stay
  additive).
- **Missing/failed registration must be noticeable to the
  contributor, via the narration handshake.** The in-session
  handshake (C2) is the observability backstop for the residual
  silent-skip: a present contributor sees the failure
  narration, or sees the expected handshake missing. The
  rendered tree cannot serve this — an unregistered session
  emits no signal the tool ever sees, so the tree cannot
  distinguish unregistered work from genuinely no work
  (unobservable by construction, a cousin of the registration
  circularity). The tree-side "this node likely has unregistered
  work" affordance is therefore deferred (Out of Scope below;
  scoping Decision 7) and tracked by the
  `deterministic-interactive-registration` backlog tripwire — it
  is the same residual.

## Naming

New identifiers this task introduces. Exact spellings are fixed
here so the command and the docs agree on one set of strings;
the implementing PR carries the literal definitions.

- **Subcommand name.** `register` — a subcommand of the existing
  `workstream-tracker` binary (the binary currently runs the
  server with no subcommand; the registration entry point is the
  first subcommand). The literal argument-parsing wiring is
  implementing-PR content.
- **Slug input.** A command flag (or equivalent environment
  variable `WST_SLUG`); the canonical slug verbatim, supplied
  by the agent after resolution.
- **Actor input.** A command flag (or equivalent environment
  variable `WST_ACTOR`). Default when unset: a generated
  per-session identifier (unique per session, never the git
  user, stable across a single session's repeat invocations so
  idempotency collapses restarts); the generation/reuse
  mechanism is implementing-PR content.
- **Server location.** Environment variable `WST_SERVER` (or an
  equivalent flag), defaulting to the local server's default
  address; aligns with the existing `PORT` default in
  [`cmd/workstream-tracker/main.go`](../../../cmd/workstream-tracker/main.go).

These spellings are this plan's estimate of the command's
surface; the implementing PR may adjust a spelling if a
structural call requires it, with the deviation called out per
the "Plan-to-PR Completion Gate."

## Files to touch

*Estimate.* The lists below are the planner's best guess at the
scope shape, not a binding rule. Implementation may revise them
when a structural call requires deviating; deviations are called
out in the implementing PR body per
[`task-plan.md`](../../../spec/planning/task-plan.md)
"Plan-to-PR Completion Gate."

### New

- `cmd/workstream-tracker/` — a new `register` subcommand entry
  point (the binary is currently server-only; Verified by:
  [`cmd/workstream-tracker/main.go`](../../../cmd/workstream-tracker/main.go),
  no subcommand dispatch).
- A thin internal HTTP client package (under `internal/`)
  that posts the exact-slug registration request and surfaces
  the observed status and body. Kept minimal — one request, no
  client framework.
- Test files covering the command's input handling,
  best-effort/non-blocking behavior, and idempotent-repeat
  behavior (excluded from the estimate count per the
  estimate-labeling convention).

### Modify

- [`docs/plans/workstream-tracker-1-0/m1-v0-2.md`](m1-v0-2.md)
  — the recorded milestone t2-contract amendment (executed in
  this plan PR; see "Milestone amendment" below).
- [`docs/backlog.md`](../../backlog.md) — add two net-new
  entries (executed in this plan PR; see "Backlog Impact"
  below): `deterministic-interactive-registration` and
  `tool-originated-task-sessions`.

### Documentation-Currency edits DEFERRED to the t2 implementation PR

These are named here but **not edited in this plan PR**. They
describe behavior that does not exist until the implementing PR
ships, so editing them now would document an unbuilt system:

- [`AGENTS.md`](../../../AGENTS.md) — the narration-handshake
  universal session rule. The auto-registration rule lives here
  and in the local agent-rules tree, **not** in the vendored
  shared modules, because `docs/agents/shared/**` is vendored
  read-only. Verified by: `AGENTS.md` ("The shared modules …
  are vendored output; do not edit them directly") — this
  corrects the milestone doc's stale "wire into
  `docs/agents/shared/`" Documentation-Currency line (scoping
  Decision 6).
- `docs/agents/local/` — the repo-owned rule-additions surface
  for the narration rule.
- [`docs/dev.md`](../../dev.md) — note that registration is
  automatic via the handshake and that a manual command
  invocation is the fallback path.
- [`design/v0.1-design.md`](../../../design/v0.1-design.md)
  §3/§4 — reflect that registration is consumed via the new
  command through the narration handshake (no API/schema
  change; the §4 API shape is unchanged).

### Intentionally not touched

*Estimate.* "We do not expect to need to touch these," not a
hard prohibition.

- `internal/api/` and `internal/db/` — this task is a pure
  consumer of t1's flow; no API, handler, or schema change is
  expected. Touching them would signal a contract violation of
  the "spec stays additive" invariant.
- `spec/planning/**` — no plan-doc spec change; the slug and
  Status contracts are consumed, not modified.
- `.claude/settings*` — harness configuration.
  `settings.local.json` is git-ignored machine-local state;
  `settings.json` is committable shared config. Verified by:
  `.gitignore`, which ignores only `.claude/settings.local.json`,
  `.claude/shell-snapshots/`, and `.claude/worktrees/` and notes
  that shared `.claude` config is still committable. This task
  touches neither — not because of gitignore, but because of the
  "consumer-agnostic; no committed harness coupling"
  cross-cutting invariant: no harness-specific configuration is
  a committed tool artifact.
- [`design/vision.md`](../../../design/vision.md) — an optional
  one-clause honesty clarification ("automatic only with a
  launch-supplied slug") was floated in scoping as the
  maintainer's call, not assumed. This plan does not edit the
  vision and does not require the clarification.

## Validation Gate

The implementing PR satisfies this gate before its Status flips
to `Landed`. The three canonical repo commands plus the
interactive-lane exercise.

- **Build.** `go build ./...` — every package compiles.
  Verified by: [`docs/dev.md:84`](../../dev.md).
- **Vet.** `go vet ./...` — static checks pass. Verified by:
  [`docs/dev.md:85`](../../dev.md).
- **Tests.** `go test ./...` — full unit-test suite. Verified
  by: [`docs/dev.md:86`](../../dev.md).
- **Command behavior.** With a slug and actor supplied and the
  server up, the command registers a work-instance that becomes
  visible in the tree; with the server down, the command logs
  and exits success; a repeat invocation for the same slug and
  actor returns the same work-instance (idempotent, no
  duplicate). The discriminator between "registered" and "looked
  like it registered" is the observed work-instance identifier
  and HTTP status plus the work-instance appearing in the
  rendered tree — not the command's exit code alone (exit code
  is success on failure by contract). The planning spike already
  demonstrated server-up (HTTP 201 with a real work-instance
  id), server-down (exit success, session proceeds), repeat
  (same id), and no-input (register nothing, exit success)
  against a t1 server.
- **Narration-handshake exercise.** An interactive prompt
  produces a narration that states the resolved plan-doc path,
  the resolved slug, and the real command/server response
  (work-instance id, HTTP status); an induced registration
  failure produces an explicit
  failure narration and the session still proceeds; an
  unresolvable prompt produces an explicit skip narration. The
  discriminator is that each narrated line corresponds to an
  observed fact (path, slug, real response) — a narration that
  asserts success with no echoed work-instance id or status is a
  failure of this gate, not a pass.

**Validation-environment hazard (must heed during the gate).**
This repository is checked out as many parallel git worktrees,
and a stale tool server from another worktree can already be
bound to the default server port. An exercise that posts to the
default address may silently hit the wrong server (a pre-t1
build that rejects the exact-slug flow), producing a false
negative. The gate must confirm it is exercising the build
under test — bind the server under test to a known free address
and point the command at that address explicitly, and confirm
the observed work-instance id was created by this server. This
hazard was hit live during the spike (a stale port-8080 server
returned `root_slug is required` for a well-formed exact-slug
request until the server under test was moved to a free port).

## Self-Review Audits

Run at the implementing PR's commit boundaries, drawn from
[`docs/agents/local/self-review-catalog.md`](../../../docs/agents/local/self-review-catalog.md):

- **`validation-honesty`** — the narration must report observed
  command/server output, not a prose success claim; a
  "registered successfully" line with no echoed work-instance id
  or HTTP status is exactly the trap this audit catches.
  Surface: the agent-rule docs and the narration contract.
- **`error-surfacing-user-mutations`** — registration is a
  user-affecting mutation; its failure must surface as a clear,
  actionable failure narration, never a silent drop or a crash.
  Surface: the command's failure path and the failure
  narration.

## Documentation Currency

Status-bearing or contract-bearing docs the implementing PR
must update **in that same PR** (named here, not edited in this
plan PR):

- [`AGENTS.md`](../../../AGENTS.md) — add the narration-handshake
  universal rule (in `AGENTS.md` + `docs/agents/local/`, not
  vendored shared, per scoping Decision 6).
- `docs/agents/local/` — the narration rule's repo-owned
  surface.
- [`docs/dev.md`](../../dev.md) — registration is automatic via
  the handshake; the manual command invocation is the documented
  fallback.
- [`design/v0.1-design.md`](../../../design/v0.1-design.md)
  §3/§4 — the registration-via-command narration reality (no
  API/schema change).
- [`docs/plans/workstream-tracker-1-0/m1-v0-2.md`](m1-v0-2.md)
  and [`docs/backlog.md`](../../backlog.md) — the t2-contract
  amendment and the backlog reconciliation are executed in
  **this** plan PR (see next sections), not deferred.

The implementing PR walks every Goal, Validation, and
Self-Review item above and either satisfies it or records an
explicit deferral with rationale in this plan doc, per the
"Plan-to-PR Completion Gate."

## Milestone amendment (executed in this plan PR)

The milestone [`m1-v0-2.md`](m1-v0-2.md) t2 "End result"
currently overclaims: agent sessions "register on start without
manual API invocation" and markers are "trustworthy because
registration is reliable." The unconditional form is false under
the registration circularity (interactive registration cannot be
deterministic). This plan PR amends the milestone's t2 contract
to the single-lane reality — registration happens early in the
session via a grounded narration handshake; observable
best-effort, not deterministic; deterministic registration
deferred to the `tool-originated-task-sessions` backlog
capability — keeping "without manual API invocation" (true: the
agent invokes the command; no human hand-crafts an API call).
This is a wording reconciliation, not a contract gutting,
mirroring how t1 reconciled the milestone's t1/t2 wording
in-PR.

Because t2 drafting also *resolves* questions the milestone
still recorded as pending or pointed at stale surfaces, this
plan PR reconciles those in the same pass (each a currency fix,
not a contract change):

- the Task Status row for `workstream-tracker-1-0-m1-t2`
  (`—` → drafted; flips to `Proposed` with the Status flip);
- the t2 "Interfaces" bullet (the slug signal was presented as
  pending with a now-rejected branch-name candidate → restated
  to the resolved reality);
- the Cross-Task Decisions entry (slug derivation moved from
  "still deferred to t2 drafting" to RESOLVED, with the section
  intro reconciled);
- the Cross-Task Risk "Slug derivation in t2 harder than
  estimated" (marked RESOLVED at t2 drafting; no escalation);
- the Documentation-Currency line that pointed t2 at the
  vendored read-only `docs/agents/shared/` (corrected to
  `AGENTS.md` + `docs/agents/local/` per scoping Decision 6).

(The `v0.1-schema` and `short_description` Cross-Task Risk
bullets are also stale relative to t1/t3 having Landed, but
that is those landed tasks' reconciliation debt, not t2's, and
is out of scope here.)

## Risk Register

Residual risks after the contract above.

- **Narration silent skip.** A loosely-followed narration
  instruction can still degrade to a missing handshake. This is
  the *irreducible* residual: the handshake itself **is** the
  observability backstop (Cross-Cutting Invariants / scoping
  Decision 7), so a skipped handshake has no in-tool backstop —
  the rendered tree cannot surface it (unobservable by
  construction). Mitigation is therefore not another in-tool
  mechanism but: the `validation-honesty` audit (catches
  fake/absent receipts at self-review) and the present
  sole-consumer noticing the *absence* of the expected
  handshake. Consciously accepted as best-effort for this
  milestone, not indefinitely — tracked by the
  `deterministic-interactive-registration` backlog entry with a
  tripwire at the neighborly-events integration milestone (where
  sole-consumer compensation evaporates).
- **Validation false-negative from a stale neighbor server.**
  Recorded in the Validation Gate as a hazard with a concrete
  mitigation (bind to a free port, point the command
  explicitly, confirm the id originated from the server under
  test). Called out because it was hit live in the spike.

## Backlog Impact

This plan PR **adds two new backlog entries** (captures —
not one of the four canonical backlog effects; disclosed as a
section-variance in the PR body). Relative to this PR's base,
neither entry exists yet, so both are net-new:

- **Add** `deterministic-interactive-registration`: authored to
  track the determinism gap that t2's observable best-effort
  narration handshake consciously accepts (the future home
  where determinism is actually achieved is the second entry
  below, not a "launch-ceremony helper"). It also tracks the
  **observability residual** — a missed registration is
  unobservable from the rendered tree by construction, so the
  handshake is t2's only backstop and the deferred tree-side
  heuristic affordance sits under the same tripwire (revisit
  before the neighborly-events integration milestone, where
  sole-consumer compensation no longer holds).
- **Add** `tool-originated-task-sessions`: the opportunity where
  the tool's UX spawns the planning agent from a clicked
  plan-tree node. Because the tool is the launcher and the click
  carries the node's identity, the spawned session's slug is
  known by construction — dissolving the registration
  circularity and yielding deterministic registration for free.
  Flagged as the home where deterministic registration is
  eventually achieved; opportunity-framed, one option among
  several, well beyond v0.2.

Neither entry graduates here; both remain Open. No entry is
deleted.

## Out of Scope

Boundary calls recorded as final answers for this task; the
deliberation prose is in the scoping doc.

- **A deterministic launcher path / committed session-start
  hook / marker-writing helper.** Cut entirely (scoping Decision
  1c): no current consumer, and the circularity is only
  dissolved by the future `tool-originated-task-sessions`
  capability where determinism is free-by-construction.
  Concrete consequence of the absence: there is no way to
  register a session that does not run the agent-driven
  handshake; such sessions go unregistered and are surfaced by
  neither the tree (unobservable by construction) nor a
  handshake (there is none) — they are simply invisible. That
  is the consciously accepted best-effort residual tracked by
  the `deterministic-interactive-registration` backlog tripwire,
  not a gap t2 closes.
- **A tree-side "this node likely has unregistered work"
  affordance.** The observability backstop the invariant
  requires is delivered in t2 by the narration handshake (C2),
  not by the rendered tree. A tree-side heuristic signal — e.g.,
  flagging a node whose plan doc is active/Proposed but that has
  no work-instance — is a genuine future capability but is
  unobservable from per-session data (the tool never sees an
  unregistered session) and exceeds this task's surface.
  Concrete consequence, stated honestly: in the rendered tree a
  node whose session never registered looks the same as a node
  with genuinely no active work — the tree does not distinguish
  them, and t2 does not claim it does (the handshake is where
  the contributor notices, not the tree). Deferred and tracked
  by the `deterministic-interactive-registration` backlog
  tripwire (same residual: the handshake backstop holds only
  while the sole consumer is present to notice it).
- **The optional `design/vision.md` honesty clause.** Left to
  the maintainer; not part of this task's contract.

## Related Docs

- [`scoping/m1-t2-auto-registration.md`](scoping/m1-t2-auto-registration.md)
  — the paired scoping doc (deliberation, rejected
  alternatives, reality-check inputs).
- [`m1-v0-2.md`](m1-v0-2.md) — parent milestone (v0.2).
- [`README.md`](README.md) — parent epic
  (`workstream-tracker-1-0`).
- [`m1-t1-multi-wi-per-slug.md`](m1-t1-multi-wi-per-slug.md) —
  the task that landed the exact-slug create-or-attach flow this
  task consumes.
- [`design/vision.md`](../../../design/vision.md) — the
  long-term vision; the feature-risk paragraph this task's
  handshake operationalizes.
- [`design/v0.1-design.md`](../../../design/v0.1-design.md) —
  §3 (plan tree / slug frontmatter) and §4 (the API and
  exact-slug create-or-attach flow).
- [`spec/planning/task-plan.md`](../../../spec/planning/task-plan.md),
  [`spec/planning/shared.md`](../../../spec/planning/shared.md)
  — the plan-doc authoring rules this plan is structured
  against.
