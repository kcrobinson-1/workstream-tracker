---
slug: tool-originated-task-sessions-m2-t2
Status: In draft
short_description: Spawn integration — /spawn endpoint + claude --bg exec
---

# m2 t2 — Spawn Integration: `/spawn` Endpoint + `claude --bg` Exec

## Context

t2 is the launcher-side surface of the
[m2 milestone](./README.md): the workstream-tracker local
server gains a new `POST /spawn` endpoint that consumes t1's
mode-affordance form, fire-and-forget execs `claude --bg
--worktree <name> --append-system-prompt-file <path>
"<initial-task>"` with `WST_SLUG=<slug>` in the child's
environment, and tightens the server's binding to loopback-only
in the same change. The endpoint is the *real construction-time
slug producer* that drives m1's deterministic register path
end-to-end (once t3's SessionStart hook is in place).

**Why now.** m2 promoted to `Proposed` at PR
[#58](https://github.com/kcrobinson-1/workstream-tracker/pull/58)
with t1, t2, and t3 seeded as three independent draftable
surfaces converging on t4. t1 Landed at PR
[#65](https://github.com/kcrobinson-1/workstream-tracker/pull/65)
with the affordance form's `action="/spawn"` already pointing
at this task's not-yet-shipped route. Drafting t2 now keeps
the three surfaces moving in parallel toward t4's convergence.

**Surfaces.** The chi router and request handlers in
[`internal/site/`](../../../../internal/site/) (the new
`POST /spawn` handler beside the existing `GET /`); the server
binding in
[`cmd/workstream-tracker/main.go`](../../../../cmd/workstream-tracker/main.go)
`runServer` (loopback-only `Addr`); two new version-controlled
per-mode prompt files at `docs/spawn-prompts/`. No spec,
agent-rule, or backlog edit.

## Goal

Implement the m2 milestone's spawn-mechanism contract (D4) and
its slug-carry mechanism (D1), grounded in the worktree
provisioning vendor surface (D5): every well-formed POST to
`/spawn` carrying a slug + mode the affordance form emits
deterministically launches a Claude Code background session in
an isolated worktree, with the slug carried out-of-band so
that t3's SessionStart hook attaches a work-instance to that
exact node by construction before the model reasons. Ship the
loopback-only binding tightening in the same change so the new
write surface lands with the trust boundary D4 leans on
already in place.

## Inherited contract (consumed by reference)

The locked WHAT for this task is the [m2 milestone doc's Task
Contracts row for t2](./README.md). The Cross-Task Invariants,
Cross-Task Decisions **D1** (slug-via-`WST_SLUG`), **D4**
(spawn-mechanism shape: server-side fire-and-forget exec of
`claude --bg --worktree …` from a new `POST /spawn`,
loopback-only binding tightening in the same change), and
**D5** (worktree provisioning via `--worktree`), and the
Cross-Task Risks "New `/spawn` endpoint expands the tool's
write surface," "Wrong-by-construction slug attaches an orphan
work-instance," and "The fencing properties erode under
feature pressure" named in that doc are this task's
load-bearing premises and are not restated here. The Contracts
below add the HOW the milestone deferred to this drafting
session — they do not loosen, redefine, or duplicate any
premise locked at milestone level.

## Cross-Cutting Invariants

The m2 README binds invariants that thread the four-task set;
t2 reads and respects them — they are consumed by reference,
not re-asserted here. See the [m2 README "Cross-Task
Invariants"](./README.md) for the full list (no new
registration surface, fencing properties hold, single-local
environment, agent-adapter seam additive, "deterministic"
stays epic-scoped, no JavaScript / no walk-on-every-request
regression).

One t2-local invariant the milestone does not lock:

- **No shell interpolation of caller data anywhere in the
  exec path.** The `os/exec.Command` invocation passes a fixed
  positional argv (`"claude"`, `"--bg"`, `"--worktree"`,
  `<worktree-name>`, `"--append-system-prompt-file"`,
  `<prompt-path>`, `<initial-task>`); the child's environment
  is set via `cmd.Env` (a constructed list, not
  `os.Setenv`); no `sh -c`, no `fmt.Sprintf` of caller data
  into a shell command line, no `bash -c`. The slug arrives
  only as `WST_SLUG` in `cmd.Env` and as the
  argv positional `<worktree-name>`/`<initial-task>` after
  passing the SD8 grammar check. The m2-level "`/spawn`
  expands the tool's write surface" risk lives at milestone
  level; this invariant is the t2-specific mechanism that
  closes it at the exec call site. *Verified by:* Go standard-
  library [`os/exec.Cmd.Env`
  semantics](https://pkg.go.dev/os/exec#Cmd) (the
  explicitly-constructed argv + Env shape the contract
  realizes); [m2 Cross-Task Risk, "New `/spawn`
  endpoint…"](./README.md) (the milestone-level surface this
  closes at the handler site).

## Contracts

### C1 — `POST /spawn` mounts on `site.Server.Router()` beside `GET /`

The new handler is registered on the same chi.Router instance
the existing index handler is registered on, as
`r.Post("/spawn", s.spawn)` alongside `r.Get("/", s.index)`.
The mount point is *not* the top-level router in `runServer`
(that one still mounts `/work-instances` and `/health` at the
top), and *not* a new internal package. The `Server` struct
gains a small set of additional fields for the new handler's
dependencies (an injectable exec function per C3 / SD11, the
prompt-files config root per C5 / SD4). *Verified by:*
[`internal/site/site.go`](../../../../internal/site/site.go)
`Server` struct and `Router()`; [m2 Documentation Currency for
t2](./README.md) ("the spawn endpoint and its tests under
`internal/site/`").

### C2 — `runServer`'s server `Addr` is loopback-only, in the same change that mounts `/spawn`

`runServer` in
[`cmd/workstream-tracker/main.go`](../../../../cmd/workstream-tracker/main.go)
constructs the `http.Server`'s `Addr` as `"127.0.0.1:" + port`,
not `":" + port`. The change ships in the same PR that mounts
`/spawn` (not a follow-up PR); the m2 Cross-Task Risk's
load-bearing security premise demands the binding be tight
*at* `/spawn` mount, not eventually. No new env var to relax
the binding is introduced — the single-contributor /
single-local-environment invariant rules out the configurable
knob. The `PORT` env var continues to govern the port; only
the interface tightens. *Verified by:*
[`cmd/workstream-tracker/main.go`](../../../../cmd/workstream-tracker/main.go)
`runServer` (the current `addr := ":" + port` form the change
replaces); [m2 Cross-Task Risk, "New `/spawn` endpoint expands
the tool's write surface"](./README.md) (the security premise
this contract realizes); [m2 Cross-Task Invariants, "Origination
stays in the contributor's single local
environment"](./README.md) (the standing constraint ruling
out the env-override knob).

### C3 — Exec invocation: fixed positional argv, no shell, explicitly-constructed `cmd.Env`

When the handler invokes `claude`, it does so via
`os/exec.Command(name, args...)` with `name = "claude"` and
a fixed positional argv that includes `--bg`,
`--worktree <name>`, `--name <name>`,
`--append-system-prompt-file <path>`, and the positional
`<initial-task>` argument — in that order or any order that
preserves each flag's adjacency to its value (the order is
HOW for the implementing PR). The child's environment is
constructed by copying the parent's environment (so PATH and
home-directory variables flow through to Claude Code's own
needs) and appending the spawn-specific additions
`WST_SLUG=<slug>` and `WST_NAME=<worktree-name>` per SD9; no
process-wide `os.Setenv` is used. No invocation routes through
a shell. The contract is "no shell interpolation; per-spawn
env scoped to the cmd," not the exact ordering of argv
elements. *Verified by:* Go [`os/exec.Cmd`
semantics](https://pkg.go.dev/os/exec#Cmd) (`Args`, `Env`
fields); [Claude Code CLI reference,
`--bg`/`--worktree`/`--name`/`--append-system-prompt-file`](https://code.claude.com/docs/en/cli-reference#cli-flags);
the t2-local Cross-Cutting Invariant on shell interpolation
above.

### C4 — Defense-in-depth input validation: slug grammar + mode-set membership

Before invoking exec, the handler validates the incoming form
fields:

- The `slug` form field passes `slugs.IsWellFormed`; on
  failure the handler returns HTTP 400 with a human-readable
  message naming the validation failure.
- The `mode` form field exact-string-matches one of the two
  non-sentinel `ModeAffordance` constants from t1
  (`ModeAffordancePlanning` = `"Begin planning"`,
  `ModeAffordanceImplementation` = `"Begin implementation"`);
  any other value, including the empty string and the
  no-affordance sentinel, returns HTTP 400.

No plan-tree walk runs on POST (the GET-side walk-on-every-
request invariant binds the GET surface only; symmetric
discipline here keeps POST plan-tree-blind). Status-shape
validation ("this slug's node currently offers this mode")
is **not** in scope at t2; t1's render is the gating site,
and the orphan-attachment residual a hand-crafted stale POST
could cause is the same residual m1 already accepted (a
typoed `--slug` today attaches an orphan; m2 produces but
does not change that residual — see the m2 Cross-Task Risk
"Wrong-by-construction slug attaches an orphan
work-instance"). *Verified by:*
[`internal/slugs/slugs.go`](../../../../internal/slugs/slugs.go)
`IsWellFormed` (the predicate the contract reuses); [m2
Cross-Task Risk, "New `/spawn` endpoint expands the tool's
write surface"](./README.md) (the defense-in-depth language);
[`internal/site/tree.go`](../../../../internal/site/tree.go)
`ModeAffordance` constants (the exact-match tokens C2's
discipline locks); [m2 Cross-Task Risk,
"Wrong-by-construction…"](./README.md) (the orphan residual
this scope explicitly does not re-litigate).

### C5 — Per-mode prompt files exist at `docs/spawn-prompts/{planning,implementation}.md`

Two version-controlled markdown files, one per
non-sentinel `ModeAffordance`, exist at the fixed paths
`docs/spawn-prompts/planning.md` and
`docs/spawn-prompts/implementation.md`. The `/spawn` handler
reads the appropriate file at request time via filesystem read
(no embed, no startup-time extraction) and passes the
absolute or working-directory-relative path to `claude` via
`--append-system-prompt-file <path>`. The config root for the
spawn-prompts directory is governed by a `SPAWN_PROMPTS_DIR`
env var defaulting to `docs/spawn-prompts` (the symmetric
posture to `PLANS_PATH`); the file basenames within the
directory are fixed by mode (`planning.md`,
`implementation.md`). The **wording** of the prompt bodies is
the implementing PR's surface per [m2 Out of Scope, "The exact
wording of the planning / implementation
prompts"](./README.md). The contract this section locks is
*existence at known paths*, *shape* (markdown files, one per
mode), and *the request-time read mechanism* — not the
wording. *Verified by:*
[`cmd/workstream-tracker/main.go`](../../../../cmd/workstream-tracker/main.go)
`runServer` (the precedent for env-overridable filesystem-
resolved config: `PLANS_PATH` / `DB_PATH` patterns); [m2
Documentation Currency for t2](./README.md) ("the per-mode
prompt files … existence and shape … lives in this
milestone"); [Claude Code CLI reference,
`--append-system-prompt-file`](https://code.claude.com/docs/en/cli-reference#cli-flags)
(the consumer of the path).

### C6 — Session-id capture + acknowledgement page

The handler captures `claude --bg`'s stdout, parses the first
line for the `backgrounded · <short-id>` prefix per SD6, and
renders a server-side HTML acknowledgement page (no JS, no
flash, no redirect) showing the short id and the `claude
attach <short-id>` invocation the contributor copies into
their own terminal. On a parse miss, the handler still
returns HTTP 200 with the raw captured stdout (and stderr,
if any) rendered for the contributor to read — a parse miss
is loud, not silent, satisfying the "deterministic ≠ cannot
fail" Cross-Task Invariant. The acknowledgement page is its
own template (template definition `spawn-ack` or sibling) and
shares the existing index page's shell CSS so it visually
belongs to the tool. *Verified by:* [Claude Code agent-view
docs, "From your shell"
example](https://code.claude.com/docs/en/agent-view#from-your-shell)
(the literal `backgrounded · 7c5dcf5d · …` output prefix);
[`internal/site/render.go`](../../../../internal/site/render.go)
`indexTmpl` (the template registration pattern the
acknowledgement page follows); [m2 Cross-Task Invariants,
"'Deterministic' stays epic-scoped"](./README.md) (the
accepted-failure-with-visibility framing the parse-miss
fallback realizes).

### C7 — Fire-and-forget posture: the workstream-tracker server is not a process manager

`claude --bg` returns immediately after the supervisor has
the session in hand (vendor-documented); the `/spawn` handler
reads stdout to completion, parses, renders the
acknowledgement page, and returns. There is no goroutine
managing the spawned `claude` process, no process-table
tracking, no parent-child relationship the workstream-tracker
server maintains. The Claude Code supervisor
(`claude daemon status`) owns the agent's process lifetime;
the workstream-tracker server returns to its observe-only
posture for that session immediately after the spawn
acknowledgement. The single non-GET surface the binary now
exposes carries no persistent state beyond the database
write the SessionStart hook produces via the deterministic
register CLI. *Verified by:* [Claude Code agent-view docs,
"The supervisor
process"](https://code.claude.com/docs/en/agent-view#the-supervisor-process)
(the supervisor's lifetime-ownership model); [m2
Cross-Task Decision D4](./README.md) (the fire-and-forget
posture clause).

### C8 — Preserved invariants: no-JS, no GET-side walk regression, agent-adapter seam additive, observe-only for non-originated sessions

The new POST surface introduces no client-side JavaScript
(neither the affordance form's submit nor the
acknowledgement page renders any `<script>`); the existing
GET-`/` index handler is unchanged (no new walk, no new DB
read on the GET path; the POST path's defense-in-depth
validation is plan-tree-blind per C4); the exec invocation's
shape is contained inside the handler such that an
alternative agent launcher (a hypothetical `claude-fork-x`,
`copilot`, etc.) is a new exec shape inside the same handler
(or a sibling handler keyed off mode), not a router or
package rewrite; sessions the tool did **not** originate
(every Claude Code session a contributor opens outside the
`/spawn` surface) continue to use m1's interactive
best-effort grounded narration handshake — the deterministic
path is the *added* mechanism, not a replacement. *Verified
by:* [m2 Cross-Task Invariants, "Agent-adapter seam stays
additive"](./README.md); [m2 Cross-Task Invariants, "No
JavaScript / no walk-on-every-request regression"](./README.md);
[m2 Cross-Task Invariants, "Fencing properties hold at every
site"](./README.md) (the running-session-stays-observe-only
clause).

## Files to touch (estimate)

Estimate of the expected change shape per
[`spec/planning/shared.md`](../../../../spec/planning/shared.md)
"Plan content is a mix of rules and estimates — label which is
which." Implementation may revise this when a structural call
requires deviating; deviations are surfaced per the "Estimate
Deviations" callout rule.

**New files:**

- `internal/site/spawn.go` — the `POST /spawn` chi handler,
  the slug + mode validation, the prompt-file resolution, the
  `os/exec.Command` invocation, the stdout parse, and the
  acknowledgement page rendering. May internally split into
  a small types file if the implementing PR judges the
  handler grows past a comfortable single-file size (HOW per
  the [`shared.md`](../../../../spec/planning/shared.md)
  "Plans describe contracts, not implementation" rule).
- `internal/site/spawn_test.go` — the table-driven defense-
  in-depth-validation tests (one row per (slug-grammar, mode)
  case the handler distinguishes), the integration test
  against the injected exec function exercising the stdout-
  parse path with a fake `claude` binary that emits the
  vendor-documented `backgrounded · …` prefix.
- `docs/spawn-prompts/planning.md` — the per-mode prompt body
  the affordance's `Begin planning` row passes to
  `--append-system-prompt-file`. The *wording* is the
  implementing PR's surface per [m2 Out of Scope](./README.md);
  the file exists at this path under this PR.
- `docs/spawn-prompts/implementation.md` — symmetric for
  `Begin implementation`.

**Modify:**

- [`internal/site/site.go`](../../../../internal/site/site.go) —
  add the `r.Post("/spawn", s.spawn)` line beside the existing
  `r.Get("/", s.index)`; extend `Server` and `New(...)` with
  the new dependencies the spawn handler needs (the injectable
  exec function per C3 / SD11, the spawn-prompts directory per
  C5).
- [`cmd/workstream-tracker/main.go`](../../../../cmd/workstream-tracker/main.go) —
  change `addr` to `"127.0.0.1:" + port` (C2); read the
  `SPAWN_PROMPTS_DIR` env var (default `docs/spawn-prompts`)
  and pass it to `site.New(...)` alongside the existing
  `plansPath`; the default `claude` exec function is wired in
  on this construction site.

**Intentionally not touched:**

- [`internal/site/forest.go`](../../../../internal/site/forest.go),
  [`internal/site/tree.go`](../../../../internal/site/tree.go),
  [`internal/site/render.go`](../../../../internal/site/render.go) —
  t1's surface, Landed. The GET-`/` index path stays unchanged
  by C8.
- [`.claude/settings.json`](../../../../.claude/settings.json) —
  the SessionStart hook entry is t3's surface per m2
  Documentation Currency.
- [`docs/agents/local/session-registration.md`](../../../../docs/agents/local/session-registration.md) —
  any additive clarification is t3's surface per m2
  Documentation Currency.
- [`cmd/workstream-tracker/register.go`](../../../../cmd/workstream-tracker/register.go) —
  the register CLI is consumed verbatim per m2 Cross-Task
  Invariant "No new registration surface."
- [`internal/api/handlers.go`](../../../../internal/api/handlers.go) —
  the exact-slug create-or-attach path is consumed verbatim;
  no new endpoint, schema, or request/response field.
- [`design/v0.1-design.md`](../../../../design/v0.1-design.md) —
  frozen v0.1 end-state record; not reconciled against per
  m2 Documentation Currency "Currency check, no edit
  expected."

## Validation Gate

The technical gate t2's implementing PR walks before merge. t2
is an **interior node** in the m2 Mermaid graph (it has an
outgoing edge to t4), so per
[`spec/planning/task-plan.md`](../../../../spec/planning/task-plan.md)
"Product-facing leaf tasks output a demo walkthrough" t2
closes on this technical gate alone — the milestone's product
walkthrough lives on t4 and is not duplicated here.

1. **OQ1 verification (the first step the implementing PR
   runs).** Before writing the rest of the handler, the
   implementing PR runs an empirical check: `WST_SLUG=oq1-test
   claude --bg --name oq1-verify "Print WST_SLUG=$WST_SLUG and
   exit."`, then `claude logs <id>` (or `claude attach <id>`),
   verifies whether the spawned session's environment
   reflects `WST_SLUG=oq1-test`. If **yes**, the spawn
   mechanism this plan locks works as designed; proceed. If
   **no**, escalate to OQ1 option (B) — synthesized per-spawn
   `.json` carrier via `--settings` — and reconcile the t3
   contract before continuing. The verification result is
   recorded in the implementing PR's body under `## Estimate
   Deviations` (or under `## Documentation` if the result
   confirms the assumption with no design change). This is the
   load-bearing gate step the rest of the plan rests on.
2. **OQ2 verification.** A follow-up empirical check: `claude
   --bg --worktree oq2-collision-test "Initial task."` twice
   in a row, observing the second invocation's behavior. The
   implementing PR records which of OQ2's (a)/(b)/(c) holds
   and writes C6's acknowledgement page to handle the observed
   failure shape gracefully.
3. **Slug-grammar + mode-set defense-in-depth test (the
   primary unit-level falsifier).** A new table-driven test
   in `internal/site/spawn_test.go` carries one row per
   defense-in-depth case: a well-formed slug + valid mode
   reaches the exec branch (the injected fake exec function
   records the invocation and the test asserts on argv shape
   + `cmd.Env` content); a malformed slug returns 400; an
   unknown mode returns 400; the empty mode returns 400;
   each non-200 case carries the expected human-readable
   error message. The test does not run a real `claude`
   binary (the injected exec function is the seam SD11
   locks).
4. **`go test ./...` clean.** The standing suite (every site
   package test t1 ships, the api package tests, the slugs
   tests) keeps passing alongside the new spawn tests.
5. **Loopback-only `Addr` check.** Either a small unit test
   of `runServer`'s `Addr` derivation, or a manual local
   verification that connecting to the workstream-tracker
   server from a non-loopback interface (e.g., the host's
   LAN address) fails connection. The minimum bar is the
   string-level check that `Addr` starts with `127.0.0.1:`;
   a deeper integration check is a self-review extension if
   the implementing PR judges the assertion worth the
   binding.
6. **Stdout-parse integration test.** An integration test
   exercises the C6 stdout parse with a fake `claude` script
   (placed on the test binary's `PATH` via t.Setenv or
   injected as the exec function's stdout) that emits the
   vendor-documented `backgrounded · <id>` prefix on the
   first line followed by indented management-command lines.
   The test asserts: (i) on the well-formed-prefix path, the
   acknowledgement page renders the short id correctly; (ii)
   on the malformed-stdout fallback path (no `backgrounded ·`
   prefix), the acknowledgement page still returns 200 and
   surfaces the raw stdout verbatim.
7. **Manual single-spawn check.** If a real `claude` binary
   is locally available, run `go run
   github.com/kcrobinson-1/workstream-tracker/cmd/workstream-tracker`
   against the local checkout's plan tree, click a `Proposed`
   leaf node's `Begin implementation` button (or curl-emulate
   the form POST), observe the acknowledgement page rendering
   a real short id, and verify `claude attach <id>` opens the
   spawned session. **t3's hook is not yet shipped** during
   t2's standalone verification, so the work-instance does
   *not* appear in the tree — observably empty under the
   spawned slug, matching the inherited contract's
   "without [the hook] in place, t2's spawn still launches
   `claude` but no work-instance attaches" clause. The manual
   check is the t2-local sanity pass — it is **not** the m2
   product walkthrough (which lives on t4 against a live
   Claude Code launcher and the full spawn/hook chain).

## Self-Review Audits

Audits the implementer walks before opening the PR, drawn from
the diff surfaces this plan touches.

- **No-shell-interpolation audit.** Search the diff for any
  `sh -c`, `bash -c`, `fmt.Sprintf(...)` whose result is later
  passed to an `exec` shell call, or any string-formatted
  shell command line containing caller-supplied slug or mode
  values. None should appear. This is the t2-local
  Cross-Cutting Invariant violated otherwise.
- **No-process-wide-`os.Setenv` audit.** Search the diff for
  any `os.Setenv` call introduced by the handler. None should
  appear — per-spawn env additions ride `cmd.Env`, not
  process-wide mutation. The workstream-tracker server runs
  chi handlers concurrently; a process-wide setenv races.
- **No-new-walk-on-GET audit.** Confirm the existing `Server.index`
  handler and the GET-path `walkPlans` + `loadActiveWorkInstances`
  + `loadSessionMetadata` chain are unchanged by this PR (the
  C8 invariant). The POST handler shares the `Server` struct
  but does not consume any of the GET-path's read functions.
- **Loopback-binding audit.** Confirm the `http.Server`'s
  `Addr` field is constructed from `"127.0.0.1:" + port`
  (or equivalent) and no path through `runServer` produces
  an all-interfaces or LAN-interface bind.
- **Affordance-form contract audit (cross-task).** Re-read
  t1's [C2](./t1-mode-affordance-render.md) and
  [C3](./t1-mode-affordance-render.md) contracts: t1 emits
  exactly `name="slug"` and `name="mode"` hidden inputs in
  the form. The `/spawn` handler reads these field names
  verbatim (no paraphrase). If t1's contract were ever
  paraphrased, the silent break would surface here; the audit
  guards the cross-task wiring contract.
- **Verified-by + load-bearing-claim audit.** Walk every load-
  bearing claim in the PR description and the plan against
  the
  [`shared.md`](../../../../spec/planning/shared.md)
  "Verified by: annotations" rule. Any new claim about the
  codebase or vendor surface introduced post-drafting carries
  a citation.

## Documentation Currency

Per the [m2 Documentation Currency clause for t2](./README.md),
this task edits the spawn endpoint and its tests under
[`internal/site/`](../../../../internal/site/), the server
binding in
[`cmd/workstream-tracker/main.go`](../../../../cmd/workstream-tracker/main.go),
and the new per-mode prompt files under `docs/spawn-prompts/`.
No status-bearing doc under
[`spec/**`](../../../../spec/) or
[`docs/agents/local/**`](../../../../docs/agents/local/) is
touched (the SessionStart hook entry and any additive
clarification to
[`session-registration.md`](../../../../docs/agents/local/session-registration.md)
are t3's surface). The frozen
[`design/v0.1-design.md`](../../../../design/v0.1-design.md)
is not reconciled against, per the project's standing posture.

## Out of Scope

- **The mode-affordance form (t1's surface).** Landed at PR
  [#65](https://github.com/kcrobinson-1/workstream-tracker/pull/65);
  the form's contract is consumed verbatim by C4.
- **The committed `.claude/settings.json` SessionStart hook
  entry and any additive clarification to
  [`session-registration.md`](../../../../docs/agents/local/session-registration.md).**
  Both are t3's surface per the m2 Task Contracts row for t3.
  If OQ1's resolution requires t3 to shift from a single
  static hook entry to a static-shell-plus-per-spawn-overlay
  shape, the redesign lands on t3, not t2 — t2 may *produce*
  the per-spawn JSON file (under OQ1 option B) but t3 owns
  the hook contract end-to-end.
- **The end-to-end product walkthrough.** Lives on t4 against
  a live local server with a real Claude Code launcher; t2's
  manual single-spawn check in the Validation Gate is the
  t2-local sanity pass, not the walkthrough.
- **The exact wording of the per-mode prompt bodies.** Per
  [m2 Out of Scope, "The exact wording of the planning /
  implementation prompts"](./README.md), the wording is
  settled in the implementing PR (the C5 contract locks the
  files' existence + paths + shape, not the wording).
- **The completion-side bracket of the session lifecycle.**
  PR [#45](https://github.com/kcrobinson-1/workstream-tracker/pull/45)
  added the symmetric `complete`/`abandon` handshake;
  per [m2 Out of Scope, "The session-end completion
  bracket"](./README.md), m2 wires only the start half. The
  spawned session is the agent process responsible for its
  own completion handshake at session close per the
  lifecycle rule's session-end section; the workstream-
  tracker server does not introduce a tool-driven completion
  mechanism.
- **Any backlog edit.** Per the [m2 Backlog
  Impact](./README.md) lock, m2's implementing PRs touch no
  backlog entries.
- **Any in-session tool→agent flow.** Per [m2 Out of Scope,
  "In-session tool→agent flow of any kind"](./README.md), m2
  crosses the one-way line only at session *origination*.
  The `/spawn` endpoint is a one-shot launch surface; the
  running session stays observe-only and file-driven.
- **A configurable bind-address knob.** SD2's loopback
  binding is not env-overridable per the single-contributor
  / single-local-environment invariant; future scenarios
  re-evaluate the constraint rather than pre-bake a knob.
- **Validation of (node, mode) Status-shape on POST.** Per
  C4, status-shape ("this slug's node currently offers this
  mode") is t1's gating site; the orphan-attachment residual
  a hand-crafted stale POST could cause is the same
  m1-accepted residual m2 does not re-litigate.

## Risk Register

t2-local risk mitigations carrying t2's share of the m2-level
risks plus t2-specific residuals.

- **Supervisor env-var pass-through (the t2-side of OQ1, the
  primary load-bearing risk).** D1's slug-carry rests on the
  spawned session seeing `WST_SLUG` from the `/spawn`
  handler's `cmd.Env`. Vendor docs do not explicitly
  guarantee this. Mitigation: (a) the OQ1-verification step
  is the *first* gate in the Validation Gate; (b) if the
  verification fails, the design escalates to the per-spawn
  `--settings` synthesized-JSON carrier (OQ1 option B) — t2
  produces the JSON, t3's static hook becomes the no-`WST_SLUG`
  shell; (c) the implementing PR records the verification
  result and any escalation in `## Estimate Deviations`. The
  m2-level reconciliation surface is the milestone doc; if
  this drafting's D1 framing turns out wrong, the milestone's
  D1 wording is updated in the implementing PR (rule
  deviation per the Plan-to-PR Completion Gate).
- **Worktree-name collision behavior (the t2-side of OQ2).** SD5
  derives a deterministic name (`<slug>-<mode-suffix>`) so
  re-spawn collides deterministically. The contract is "the
  collision is observable to the contributor" regardless of
  whether Claude Code errors or reuses the existing worktree.
  Mitigation: C6's parse-miss fallback already covers the
  surface-error case (the raw stdout/stderr renders in the
  acknowledgement page); the implementing PR verifies which
  of OQ2's (a)/(b)/(c) holds and tunes the acknowledgement
  page's prose accordingly. No design change as a function of
  OQ2's resolution.
- **`/spawn` write-surface expansion (the m2 Cross-Task
  Risk).** Mitigations live in C1–C4 + the no-shell-
  interpolation invariant: the new POST is loopback-only
  (C2), defense-in-depth-validates inputs (C4), uses a
  fixed-argv exec with no shell interpolation (the
  t2-local Cross-Cutting Invariant + C3), and the
  fire-and-forget posture (C7) means no persistent
  workstream-tracker-side state accumulates. The agent-
  adapter seam (C8) keeps the shape inside the handler
  additive.
- **Stdout-parse fragility under future vendor format
  changes.** SD6's first-line `backgrounded · ` prefix match
  is documented in the agent-view docs' "From your shell"
  example, but the docs do not pin the format as a stable
  API. Mitigation: C6's parse-miss fallback (200 with raw
  stdout) is the load-bearing acceptance posture — a future
  format change does not break the spawn correctness, only
  degrades the acknowledgement-page polish to "here's the
  raw output." The integration test in the Validation Gate
  asserts both the well-formed-prefix and the fallback
  paths.
- **`claude` binary not on `PATH`.** A contributor running
  the workstream-tracker server without `claude` installed
  hits an `exec: "claude": executable file not found in
  $PATH` error from `os/exec`. Mitigation: the `/spawn`
  handler returns the exec failure verbatim in the
  acknowledgement page (the same observable-failure posture
  as C6's parse-miss fallback) — the contributor sees what
  went wrong and installs Claude Code; the workstream-
  tracker server does not crash, does not retry, does not
  poll. Per the "deterministic stays epic-scoped"
  invariant, a launch that fails to start is observable,
  not silently swallowed.
- **Concurrent POSTs racing on the same slug+mode.** Two
  near-simultaneous POSTs for the same (slug, mode) would
  both attempt to spawn `claude --bg --worktree <name>`
  with the same name. Whether OQ2's resolution surfaces
  this as a hard error or a soft collision, the second
  POST's acknowledgement page renders the failure (per
  C6's parse-miss fallback). The workstream-tracker server
  does not serialize spawn requests; the underlying race
  is bounded by Claude Code's own worktree-name semantics.
- **Stale work-instance from a prior spawn.** A previous
  spawn for the same slug + actor that didn't symmetrically
  `complete` leaves an active work-instance row whose
  (slug, actor) pair the new spawn's register would
  idempotency-collapse onto, per
  [`internal/api/handlers.go`](../../../../internal/api/handlers.go)
  `activeWorkInstanceID`. The new spawn's session would
  appear to "attach" to the old work-instance. This is the
  same residual the
  [`deterministic-interactive-registration`](../../../../docs/backlog.md#deterministic-interactive-registration)
  backlog tripwire already names for the idle-window-shared
  actor case; m2 does not extend or re-litigate it. The
  generated per-session actor's idle window
  ([`cmd/workstream-tracker/register.go`](../../../../cmd/workstream-tracker/register.go)
  `sessionActorIdleWindow`) bounds the practical scope.

## Open Questions

Two open questions surfaced by this scoping pass; both are
vendor-behavior facts the implementing PR will verify
empirically as the first step of the Validation Gate. The
human resolves them at the promotion gate (either by
authorizing the "implementing PR verifies and escalates"
path or by deciding to run a spike now).

- **[OQ1 — Does the supervisor pass invoker-set environment
  variables through to the spawned background session and its
  SessionStart hook?](./scoping/t2-spawn-endpoint.md#oq1--does-the-supervisor-pass-invoker-set-environment-variables-wst_slug-wst_name-through-to-the-spawned-background-session-and-its-sessionstart-hook)**
  Load-bearing for D1; fallback option (B) shifts t3's surface.
- **[OQ2 — When `claude --bg --worktree <name>` is invoked and
  the worktree directory already exists, what
  happens?](./scoping/t2-spawn-endpoint.md#oq2--when-claude---bg---worktree-name-is-invoked-and-reposclaudeworktreesname-already-exists-what-happens)**
  Surfaces in C6's acknowledgement page; no design change as
  a function of resolution.

See [the paired scoping doc](./scoping/t2-spawn-endpoint.md)
for the full decision-space decomposition of each.

## Backlog Impact

None (m2 README locks this at milestone level per [m2 Backlog
Impact](./README.md)).

## Related Docs

- [`./README.md`](./README.md) — m2 milestone doc; the locked
  WHAT, Cross-Task Invariants, Cross-Task Decisions
  (especially D1, D4, D5), Cross-Task Risks, Documentation
  Currency.
- [`./scoping/t2-spawn-endpoint.md`](./scoping/t2-spawn-endpoint.md) —
  paired scoping doc carrying SD1…SD11 with rejected
  alternatives, the reality-check inputs, and the open-
  question deliberation; transient, deletes at t4's
  milestone-terminal PR.
- [`./t1-mode-affordance-render.md`](./t1-mode-affordance-render.md) —
  Landed sibling task; produces the affordance form `/spawn`
  consumes (the `name="slug"` and `name="mode"` hidden inputs
  C4 reads).
- [`./t3-session-start-hook.md`](./t3-session-start-hook.md) —
  sibling task skeleton; consumes t2's spawn (the `WST_SLUG`-
  set environment); OQ1's resolution may shift t3's surface
  from static-hook to static-shell-plus-per-spawn-overlay.
- [`./t4-end-to-end-validation.md`](./t4-end-to-end-validation.md) —
  the milestone-graph leaf t2 converges on; carries the
  product walkthrough that observes t2's spawn end-to-end.
- [`../m1/README.md`](../m1/README.md) — Landed m1 milestone;
  the deterministic register CLI t2 consumes verbatim.
- [`../m1/t1-deterministic-path-contract.md`](../m1/t1-deterministic-path-contract.md) —
  Landed; the durable deterministic-path contract t2's
  spawn produces the real construction-time slug for.
- [`../../README.md`](../../README.md) — parent epic;
  Cross-Cutting Invariants and the resolved vision inputs
  the milestone contract sits on top of.
- [`../../../../spec/planning/task-plan.md`](../../../../spec/planning/task-plan.md) —
  the rules this task plan is structured against; "Required
  and optional sections," "Plan-to-PR Completion Gate," the
  `In draft` → `Proposed` promotion gate (run by the human,
  not this drafting session).
- [`../../../../spec/planning/shared.md`](../../../../spec/planning/shared.md) —
  cross-level rules; "Plans describe contracts, not
  implementation," "Verified by: annotations," "Plan-doc
  Status," "Decompose options into shapes," "Quote labels
  whose enforcement depends on exact-match matching."
- [`../../../../internal/site/site.go`](../../../../internal/site/site.go),
  [`../../../../internal/site/tree.go`](../../../../internal/site/tree.go),
  [`../../../../internal/site/render.go`](../../../../internal/site/render.go),
  [`../../../../cmd/workstream-tracker/main.go`](../../../../cmd/workstream-tracker/main.go),
  [`../../../../cmd/workstream-tracker/register.go`](../../../../cmd/workstream-tracker/register.go),
  [`../../../../internal/slugs/slugs.go`](../../../../internal/slugs/slugs.go),
  [`../../../../internal/api/handlers.go`](../../../../internal/api/handlers.go) —
  the seven files this task's estimate touches or
  consumes-by-reference.
- [Claude Code CLI
  reference](https://code.claude.com/docs/en/cli-reference) —
  the vendor surface D4 / D5 verify against (`--bg`,
  `--worktree`, `--name`, `--append-system-prompt-file`,
  `claude attach`, `claude daemon status`).
- [Claude Code agent-view
  documentation](https://code.claude.com/docs/en/agent-view) —
  the supervisor + worktree + `--bg` output-format surface
  (the `backgrounded · <short-id> · <name>` stdout format C6
  parses; the "How file edits are isolated" worktree
  semantics D5 relies on).
- [Claude Code hook
  documentation](https://code.claude.com/docs/en/hooks) —
  the SessionStart hook contract t3 satisfies; the
  environment-inheritance ambiguity OQ1 surfaces.
