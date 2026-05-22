---
slug: tool-originated-task-sessions-m2-t2
short_description: Scoping — spawn integration: /spawn endpoint + claude --bg exec
---

# Scoping — m2 t2 Spawn Endpoint

Scoping doc paired with the
[`t2` task plan](../t2-spawn-endpoint.md). Transient per
[`task-plan.md`](../../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns" — deletes in batch at the
milestone-terminal PR (t4's PR, not t2's). Carries no `Status`
field per the same rule (the lifecycle is for plans, not for
transient scoping).

## Context

t2 implements the m2 milestone's **spawn integration**: a new
`POST /spawn` endpoint on the workstream-tracker local server,
the `os/exec` invocation of `claude --bg --worktree <name>
--append-system-prompt-file <prompt-file> "<initial-task>"`
with `WST_SLUG=<slug>` in the child environment, and the
loopback-only server-binding tightening that ships in the same
change. The spawn endpoint is the *launcher* surface t1's
affordance form POSTs to; the determinism-resolution wiring
completes when t3's SessionStart hook runs the deterministic
register CLI in the spawned session.

This scoping doc lays out the HOW the m2 milestone deferred to
this drafting session. The locked WHAT — D1 (slug-via-`WST_SLUG`),
D4 (spawn-mechanism shape), D5 (worktree provisioning via
`--worktree`), and the Cross-Task Invariants — comes from the
[m2 milestone doc](../README.md) and is consumed verbatim, not
re-derived.

## Reality-check inputs (the plan's load-bearing premises)

Each item below is a code or vendor-doc fact the plan rests on.
Each is verified by reading the cited surface before drafting,
not by transitive citation.

- **The chi router at the top level mounts `apiServer.MountRoutes`
  under `/work-instances` and `siteServer.Router()` at `/`; the
  site router today exposes only `r.Get("/", s.index)`.** The new
  `POST /spawn` handler will mount inside `siteServer.Router()`
  beside the existing `r.Get("/", s.index)` line — that is the
  natural sibling location per the m2 README Documentation
  Currency ("the spawn endpoint and its tests under
  `internal/site/`"). *Verified by:*
  [`internal/site/site.go`](../../../../../internal/site/site.go)
  `Router()` returning `chi.NewRouter()` with `r.Get("/", s.index)`;
  [`cmd/workstream-tracker/main.go`](../../../../../cmd/workstream-tracker/main.go)
  `runServer` mounting `r.Mount("/", siteServer.Router())` beside
  `r.Route("/work-instances", apiServer.MountRoutes)` and
  `r.Get("/health", health)`.
- **The current server binding is `addr := ":" + port` —
  all-interfaces.** The m2 Cross-Task Risk on `/spawn` expanding
  the write surface leans on t2 changing this to loopback-only
  in the same change that mounts `/spawn`. The constraint is
  load-bearing: anyone on the LAN could otherwise fire
  `POST /spawn`, and the trust boundary the design assumes
  ("anyone on loopback is by construction the contributor on
  the same machine") would not hold. *Verified by:*
  [`cmd/workstream-tracker/main.go`](../../../../../cmd/workstream-tracker/main.go)
  `runServer` `addr := ":" + port`.
- **The deterministic register CLI reads its slug from the
  `--slug` flag or the `WST_SLUG` env var, the actor from
  `--actor`/`WST_ACTOR`, the server from `--server`/`WST_SERVER`
  (default `http://localhost:8080`), and the optional human
  display name from `--name`/`WST_NAME`.** D1's "env var
  carrier" mechanism rides the existing `WST_SLUG` fallback in
  `runRegister`; no new CLI flag or schema. D5's enrichment
  leg rides the existing `--name`/`WST_NAME` reading. *Verified
  by:*
  [`cmd/workstream-tracker/register.go`](../../../../../cmd/workstream-tracker/register.go)
  `runRegister` (lines 39–93: the `--slug`/`WST_SLUG`,
  `--actor`/`WST_ACTOR`, `--server`/`WST_SERVER`,
  `--name`/`WST_NAME` reads; the empty-slug short-circuit at
  lines 51–55).
- **The exact-slug create-or-attach path's grammar guard is
  `slugs.IsWellFormed`.** Defense-in-depth slug validation in
  the `/spawn` handler uses this same predicate so a hand-
  crafted POST cannot push a malformed slug into the
  registration leg or into the os/exec argv. The predicate
  enforces the kebab-case root + ordered position-segment
  grammar without consulting any plan-tree state, so the
  handler stays plan-tree-blind (no walk on POST). *Verified
  by:*
  [`internal/slugs/slugs.go`](../../../../../internal/slugs/slugs.go)
  `IsWellFormed` (the standalone grammar check used by the
  `ExactSlug` flow in
  [`internal/api/handlers.go`](../../../../../internal/api/handlers.go)
  `registerWorkInstance`).
- **`claude --bg` prints `backgrounded · <short-id> · <name>`
  (or `backgrounded · <short-id>` when `--name` is absent) on
  the first stdout line, followed by indented management
  commands.** The handler captures stdout and extracts the
  short id from the first-line pattern. *Verified by:*
  [Claude Code agent-view docs, "From your shell"
  example](https://code.claude.com/docs/en/agent-view#from-your-shell)
  ("After backgrounding, Claude prints the session's short ID
  and the commands for managing it" — quoting the literal
  output block `backgrounded · 7c5dcf5d · flaky-test-fix /
  claude agents             list sessions / claude attach
  7c5dcf5d    open in this terminal / …`).
- **`claude --worktree <name>` starts the session in an
  isolated git worktree at `<repo>/.claude/worktrees/<name>`;
  if no name is given, one is auto-generated.** D5's "worktree
  provisioning is the `--worktree` flag" rests on this. The
  spawn invocation passes the slug-derived name explicitly.
  *Verified by:* [Claude Code CLI reference,
  `--worktree`/`-w`](https://code.claude.com/docs/en/cli-reference#cli-flags)
  ("Start Claude in an isolated git worktree at
  `<repo>/.claude/worktrees/<name>`. If no name is given, one
  is auto-generated").
- **`claude --append-system-prompt-file <path>` appends the
  file's contents to the default system prompt; this flag works
  in both interactive and non-interactive (including `--bg`)
  modes.** D4's "mode-prompt carrier" rests on this. The four
  system-prompt flags (replace vs append × inline vs file) are
  the documented surface; `--append-system-prompt-file` is the
  one whose semantics fit the per-mode bootstrap content while
  preserving Claude Code's default tool guidance and safety
  instructions. *Verified by:* [Claude Code CLI reference,
  `--append-system-prompt-file`](https://code.claude.com/docs/en/cli-reference#cli-flags)
  ("Load additional system prompt text from a file and append
  to the default prompt"); [System prompt flags
  table](https://code.claude.com/docs/en/cli-reference#system-prompt-flags)
  ("All four work in both interactive and non-interactive
  modes").
- **`claude --bg` returns immediately after registering the
  session with the supervisor; the supervisor is a per-user
  long-lived process that owns the agent's lifetime.** D4's
  "fire-and-forget" posture rests on this — the workstream-
  tracker server does not become a process manager. *Verified
  by:* [Claude Code CLI reference,
  `--bg`](https://code.claude.com/docs/en/cli-reference#cli-flags)
  ("Start the session as a background agent and return
  immediately. Prints the session ID and management commands");
  [agent-view docs, "The supervisor
  process"](https://code.claude.com/docs/en/agent-view#the-supervisor-process)
  ("Background sessions are hosted by a per-user supervisor
  process, separate from your terminal and from agent view.
  The supervisor starts automatically the first time you
  background a session or open agent view, and you don't
  manage it directly. … Each background session is its own
  Claude Code process, managed by the supervisor rather than
  tied to your terminal").
- **Claude Code's SessionStart hook with matcher `startup`
  fires for new sessions. The hook is configured in
  `.claude/settings.json` with `type: "command"`. Non-zero
  hook exit codes do **not** block the session.** t2's spawn
  *assumes* t3's hook is in place; this scoping doc consumes
  the hook contract verbatim (t3 owns the entry). *Verified
  by:* [Claude Code SessionStart hook
  documentation](https://code.claude.com/docs/en/hooks)
  ("Runs when Claude Code starts a new session or resumes an
  existing session"; "startup: New session"; "SessionStart
  hooks cannot block").
- **Whether the supervisor passes the invoker's environment
  variables (set on the `claude --bg` process) through to the
  spawned background session — and therefore to the
  SessionStart hook process — is **not explicitly stated** in
  the vendor docs.** The hooks doc says "Handlers run in the
  current directory with Claude Code's environment," which is
  ambiguous as to *which* Claude Code process owns that
  environment when the session is spawned by the supervisor
  rather than directly by `claude --bg`. The agent-view doc's
  "Configuration flags from the original launch carry through
  to the backgrounded session" clause names *flags*, not env
  vars, and refers to `/bg` backgrounding, not `claude --bg`
  spawning. This is the load-bearing assumption D1 (env-var
  slug carrier) rests on, and the implementing PR must verify
  it before writing the rest of the handler — see
  [**OQ1**](#open-decisions-to-make-at-plan-drafting) below.
  *Verified by:*
  [Claude Code SessionStart hook documentation,
  "Handlers run in the current directory with Claude Code's
  environment"](https://code.claude.com/docs/en/hooks);
  [agent-view "From inside a session" — note this is `/bg`,
  not
  `--bg`](https://code.claude.com/docs/en/agent-view#from-inside-a-session)
  ("Configuration flags from the original launch carry through
  to the backgrounded session, so its MCP servers, settings,
  and fallback model remain in effect").
- **Background sessions start in the invoker's working
  directory; `--worktree <name>` moves them into the named
  isolated worktree at session start.** Without `--worktree`,
  Claude auto-creates a worktree under
  `<repo>/.claude/worktrees/` *before editing files* (later in
  the session). t2 passes `--worktree` explicitly so the
  worktree name is known at spawn time (D5), supporting the
  enrichment-leg correlation. *Verified by:* [agent-view docs,
  "How file edits are
  isolated"](https://code.claude.com/docs/en/agent-view#how-file-edits-are-isolated)
  ("Every background session, whether started from agent view,
  `/bg`, or `claude --bg`, starts in your working directory.
  Before editing files, Claude moves the session into an
  isolated git worktree under `.claude/worktrees/`").
- **Claude Code's `claude --name <name>` flag sets the
  session's display name in agent view; the same name appears
  in `--bg` output after the short id.** Independent of
  workstream-tracker's `WST_NAME`/`--name` reading (which sets
  the roster's display label). The two are independent
  mechanisms; SD9 below addresses whether the spawn should set
  one, both, or neither and to what value. *Verified by:*
  [Claude Code CLI reference,
  `--name`/`-n`](https://code.claude.com/docs/en/cli-reference#cli-flags)
  ("Set a display name for the session, shown in `/resume` and
  the terminal title"); [agent-view docs, "From your shell"
  example](https://code.claude.com/docs/en/agent-view#from-your-shell)
  (the `--bg --name` output example).
- **m2 D1, D4, and D5 are locked at milestone level; the
  Cross-Task Invariants thread t2 by reference.** This drafting
  derives HOW from the locked WHAT and does not re-litigate it.
  *Verified by:* [m2 README, Cross-Task Decisions D1, D4,
  D5](../README.md); [m2 README, Cross-Task
  Invariants](../README.md).

## Decisions made at scoping time

Each decision below decomposes its option space into shapes per
[`shared.md`](../../../../../spec/planning/shared.md) "Decompose
options into shapes before analyzing," names the chosen shape,
and names the rejected alternatives with rationale.

### SD1 — `/spawn` mounts inside `site.Server.Router()`, beside `GET /`

Where the `POST /spawn` handler lives in the chi-router topology.

**Shapes considered:**

- **(a) Inside `site.Server.Router()`** in
  [`internal/site/site.go`](../../../../../internal/site/site.go),
  as `r.Post("/spawn", s.spawn)` beside the existing `r.Get("/",
  s.index)`. The `Server` struct's `db` and `plansPath`
  dependencies are available; spawn-specific dependencies (an
  injectable exec function, the prompt-files config) extend the
  same struct.
- **(b) At the top-level router** in
  [`cmd/workstream-tracker/main.go`](../../../../../cmd/workstream-tracker/main.go)
  `runServer`, as a third `r.Route("/spawn", …)` sibling to
  `/work-instances` and `/`. Spawn becomes its own top-level
  surface, structurally parallel to the API surface.
- **(c) New `internal/spawn` package** with its own `Server`
  and `Router()` mounted from `main.go`. Cleanly separates the
  page-write surface from the page-read surface; adds one more
  internal package.

**Chosen: (a).** The `/spawn` endpoint is conceptually paired
with the page-side `GET /` — t1's affordance form's
`action="/spawn"` POSTs to a route the contributor experiences
as the same page's submit target. Co-locating the read and
write halves of the page-write surface in `internal/site/` keeps
the page-write surface contained. (b) and (c) introduce
structural separation that the m2 surface scope does not ask
for; (c) adds a third internal package for a single handler.
The site package's read-only history is incidental, not a
structural rule. *Verified by:*
[`internal/site/site.go`](../../../../../internal/site/site.go)
`Router()` (the natural sibling location); [m2 Documentation
Currency for t2](../README.md) ("the spawn endpoint and its
tests under `internal/site/`").

### SD2 — Loopback binding spelled `127.0.0.1:<port>`

How `runServer`'s `addr` is changed from all-interfaces to
loopback-only.

**Shapes considered:**

- **(a) `addr := "127.0.0.1:" + port`** — IPv4-only loopback.
  Single canonical spelling for the loopback interface; the OS
  binds exactly to 127.0.0.1.
- **(b) `addr := "localhost:" + port`** — relies on the
  resolver's `localhost` entry. Most systems resolve to
  127.0.0.1 (and also ::1, depending on /etc/hosts and the
  resolver order), but the binding is no longer
  spelling-deterministic.
- **(c) `addr := "[::1]:" + port`** — IPv6-only loopback.
  Refuses IPv4 connections; a contributor whose browser tries
  127.0.0.1 fails.
- **(d) Bind to both `127.0.0.1` and `[::1]`** via two
  listeners. Most defensive against dual-stack systems; doubles
  the listener wiring for a posture that 127.0.0.1 alone covers
  for the contributor's browser.
- **(e) Env override** (`WST_BIND_ADDR` defaulting to
  `127.0.0.1`) so future scenarios can relax. Trades a
  load-bearing constraint for a configurable knob whose default
  satisfies the constraint.

**Chosen: (a).** `127.0.0.1` is the canonical IPv4 loopback
spelling; combined with the existing `PORT` env var, the
binding is the deterministic minimum spelling for the trust
boundary D4 leans on. The standing single-contributor /
single-local-environment invariant rules out the env-override
shape (e) — relaxing the binding is the kind of feature
pressure the fencing properties Cross-Task Invariant guards
against; future scenarios re-evaluate the constraint, they
don't pre-bake a knob. (b)/(c)/(d) are correctness or
spelling-determinism regressions. *Verified by:*
[`cmd/workstream-tracker/main.go`](../../../../../cmd/workstream-tracker/main.go)
`runServer` (the current `addr := ":" + port` form replaces
1:1); [m2 Cross-Task Risk, "New `/spawn` endpoint expands the
tool's write surface"](../README.md) (the load-bearing premise
the loopback-only binding leans into); [m2 Cross-Task
Invariants, "Origination stays in the contributor's single
local environment"](../README.md) (the standing constraint
ruling out the env-override knob).

### SD3 — exec invocation: `os/exec.Command` with positional argv, no shell, explicitly constructed `Env`

How the `os/exec` call is shaped to avoid shell injection and to
control the child's environment.

**Shapes considered:**

- **(a) `exec.Command("claude", argv...)` with a fixed
  positional argv** (`--bg`, `--worktree`, `<worktree-name>`,
  `--append-system-prompt-file`, `<prompt-file>`,
  `<initial-task>`) and an explicitly constructed `cmd.Env`
  that copies the parent env plus the spawn-specific
  additions (`WST_SLUG`, `WST_NAME`). No shell, no string
  interpolation of caller data into a shell command line.
- **(b) `exec.Command("sh", "-c", "claude --bg …")`** — shell
  interpolation of caller-supplied slug/mode into a command
  string. Rejected: shell-injection vector; the slug or mode
  string would need escaping every call site.
- **(c) Inherit the parent env wholesale (`cmd.Env` unset)**
  and rely on `os.Setenv("WST_SLUG", …)` before the exec.
  Rejected: leaks process-wide state, races against concurrent
  spawns (the workstream-tracker server runs the chi handlers
  concurrently), and pollutes the workstream-tracker server's
  own env.

**Chosen: (a).** Positional argv with no shell removes the
injection vector; explicitly constructed `cmd.Env` (parent env
copy + spawn additions, scoped to the cmd) keeps per-request
state contained. The exact wrapper-vs-direct construction of
`cmd.Env` is implementation choice; the contract is *no shell
interpolation* and *no process-wide env mutation*. *Verified
by:* Go standard-library `os/exec` semantics ([package
docs](https://pkg.go.dev/os/exec#Cmd) — `Cmd.Env`: "If Env is
nil, the new process uses the current process's environment");
[m2 Cross-Task Risk, "New `/spawn` endpoint expands the tool's
write surface"](../README.md) (the canonical-argv discipline
this realizes).

### SD4 — Per-mode prompt files at `docs/spawn-prompts/{planning,implementation}.md`, filesystem-resolved at request time

Where the per-mode prompt bodies (the content
`--append-system-prompt-file` reads) live, and how the handler
finds them.

**Shapes considered:**

- **(a) `docs/spawn-prompts/planning.md` and
  `docs/spawn-prompts/implementation.md`**, two markdown files
  at fixed paths under a new top-level `docs/spawn-prompts/`
  subfolder. The handler resolves the path at request time
  (filesystem read) relative to a known config root, mirroring
  the existing `plansPath` posture.
- **(b) `internal/site/spawn-prompts/{planning,implementation}.md`
  embedded into the Go binary via `embed.FS`**, extracted to a
  per-server-process tempdir at startup. Prompts travel with
  the binary; no runtime filesystem dependency.
- **(c) `docs/agents/local/spawn-prompts/{planning,implementation}.md`**
  under the existing agent-rule scope. Semantic mismatch — the
  prompts are spawn-bootstrap content, not authoring-rule
  content.
- **(d) Hard-coded inline string constants** in
  `internal/site/spawn.go`. No version-control affordance for
  the prompt body separate from code; edits force a Go file
  change. Rejected: the m2 milestone explicitly names the
  prompt bodies as version-controlled *files*.

**Chosen: (a).** A new top-level `docs/spawn-prompts/` mirrors
`docs/plans/` and `docs/agents/`: a discoverable, filesystem-
resolved category of human-authored content the server reads at
request time. The path posture matches `plansPath` (an env-
overridable config root resolving to `docs/plans` by default);
analogously `SPAWN_PROMPTS_DIR` defaults to `docs/spawn-prompts`
and is dependency-injected on `Server`. (b)'s embed-then-extract
mechanism adds startup-time complexity and a per-process
tempdir for no net gain; the prompts edit cycle is a PR + merge
either way. (c) misplaces the content semantically. **The
exact wording of the prompt bodies is settled in t2's
implementing PR, not in this scoping doc**, per
[m2 Out of Scope](../README.md) "The exact wording of the
planning / implementation prompts." This SD locks the
*existence*, *shape* (two markdown files, one per mode), and
*storage path* per [m2 Documentation Currency for
t2](../README.md). *Verified by:*
[`cmd/workstream-tracker/main.go`](../../../../../cmd/workstream-tracker/main.go)
`runServer` (the precedent for env-overridable filesystem-
resolved config: `PLANS_PATH` defaulting to `docs/plans`); [m2
Documentation Currency for t2](../README.md) (the prompt files'
existence and shape live at milestone level; the paths are
HOW for this session).

### SD5 — Worktree name = `<slug>-<mode-suffix>` where mode-suffix is the kebab-case lowering of the mode string's terminal word

How the spawn derives the `--worktree <name>` argument from the
slug + mode.

**Shapes considered:**

- **(a) `<slug>-<mode-suffix>`** — concat the slug and a short
  mode token (`planning` from "Begin planning"; `implementation`
  from "Begin implementation"). Stable, readable, deterministic.
  Re-spawning the same (slug, mode) collides on the worktree
  name deterministically — the collision is *observable* (the
  contributor sees the error), which is the right diagnostic
  shape.
- **(b) Slug verbatim as `--worktree <slug>`.** Loses the mode
  disambiguation — planning and implementation spawns for the
  same node collide on one worktree, but they're conceptually
  separate work.
- **(c) Generated random suffix** (e.g.,
  `<slug>-<mode-suffix>-<uuid8>`). Every spawn gets a fresh
  worktree; orphan worktrees accumulate under
  `.claude/worktrees/` over time; loses the readability and the
  collision-as-diagnostic property.
- **(d) Don't pass `--worktree`** — let Claude Code auto-
  generate. Breaks D5's "name known at spawn time" property the
  enrichment-leg correlation rests on; agent-view's auto-name
  is opaque from the workstream-tracker's vantage.

**Chosen: (a).** Stable + readable + collision-as-diagnostic.
The mode-suffix derivation is a small in-handler lookup keyed
on the `ModeAffordance` constants t1 defined in
[`internal/site/tree.go`](../../../../../internal/site/tree.go)
(`ModeAffordancePlanning` → `planning`,
`ModeAffordanceImplementation` → `implementation`); the
implementing PR settles whether the lookup is a switch, a map
literal, or a method on `ModeAffordance` (HOW). Collision
behavior with Claude Code's `--worktree <name>` when the named
worktree directory already exists is **OQ2** below — the
contract is "the collision is observable to the contributor"
regardless of whether Claude Code errors or reuses; SD5's
choice does not change as a function of OQ2's resolution.
*Verified by:*
[`internal/site/tree.go`](../../../../../internal/site/tree.go)
`ModeAffordance` constants and `Affordance()` (the source of
truth for the mode strings t1 locked); [m2 D5](../README.md)
("the worktree name is therefore known at spawn time by the
launcher").

### SD6 — Session-id capture: parse the first stdout line for the `backgrounded · <short-id>` prefix

How the handler extracts the session id from `claude --bg`
output.

**Shapes considered:**

- **(a) Capture stdout via `cmd.Output()`; parse the first
  line for the literal prefix `backgrounded · ` followed by
  the short id (whitespace-delimited, the second word).** The
  vendor doc's example output is `backgrounded · 7c5dcf5d ·
  flaky-test-fix`. The handler reads up to a fixed byte budget
  (the vendor output is small — a header line + four
  indented management-command lines), splits on newline,
  takes the first line, and strips the `backgrounded · `
  prefix. Defense-in-depth: if the prefix doesn't match, the
  handler still returns a 200 acknowledgement page reporting
  the raw stdout for the contributor to read manually (the
  spawn ran; the parse failed loudly, observable).
- **(b) Use `claude agents --json` after `claude --bg` to
  query the supervisor for the latest session.** Adds a second
  exec call per spawn; the result's ordering and freshness
  semantics aren't specified for "newly-spawned-this-request"
  scoping.
- **(c) Don't capture the session id; surface the raw stdout
  to the contributor.** Wastes the structured-output
  affordance the vendor provides.

**Chosen: (a).** Parse the first-line prefix; fall back to
"display the raw stdout" if the parse misses. The parse pattern
is a string-prefix match (not a fragile regex over the whole
output). *Verified by:* [Claude Code agent-view, "From your
shell" example](https://code.claude.com/docs/en/agent-view#from-your-shell)
(the literal `backgrounded · 7c5dcf5d` output prefix).

### SD7 — POST response: server-rendered HTML acknowledgement page with the session id and the `claude attach <id>` invocation

What the contributor sees after their browser POSTs the
affordance form.

**Shapes considered:**

- **(a) Server-rendered HTML response** showing "Session
  `<short-id>` started. Run `claude attach <short-id>` in your
  terminal to attach. ↩ Back to the plan tree." with the raw
  command in a `<code>` block ready to copy. No-JS, no
  client-side machinery.
- **(b) HTTP 303 redirect to `/?spawned=<short-id>`** with a
  query-param flash that the index renders as a banner.
  Requires extending the GET handler to read and render the
  flash; mixes spawn-state into the read surface.
- **(c) Plain-text body** (e.g., `Session 7c5dcf5d started.
  Run: claude attach 7c5dcf5d`). Discoverable but visually
  thin; the no-back-link makes the contributor use the browser
  back button.
- **(d) JSON response** (no JS to consume it). Wrong shape for
  a browser-form POST.

**Chosen: (a).** A small HTML acknowledgement page is the
no-JS minimum that's still readable and copy-able. Uses the
existing `indexTmpl` template registration pattern — a sibling
template definition (`spawnTmpl` or a `define "spawn-ack"`
inside the existing template set) keeps the rendering
machinery contained. *Verified by:*
[`internal/site/render.go`](../../../../../internal/site/render.go)
`indexTmpl` (the template registration pattern the spawn page
follows); the affordance form's `action="/spawn"` (a regular
HTML form submit, so the response is rendered in the browser
window).

### SD8 — Defense-in-depth input validation: slug-grammar check + mode-set membership; no plan-tree walk on POST

What the `/spawn` handler validates before invoking `exec`.

**Shapes considered:**

- **(a) Validate the slug against `slugs.IsWellFormed`;
  validate the mode against the set of `ModeAffordance`
  constants minus `ModeAffordanceNone`; reject either failure
  with HTTP 400.** No plan-tree walk on POST; the handler
  trusts the slug's grammar but does not verify the slug names
  a real plan-tree node (the existing exact-slug create-or-
  attach path is similarly plan-tree-blind per m1's
  trust-the-caller posture).
- **(b) Walk the plan tree on POST and verify the (node,
  mode) pair matches D3's affordance map.** Adds a per-POST
  filesystem walk regression of the walk-on-every-request
  invariant's spirit (the invariant is named for GET, but a
  per-POST walk doubles the read posture for a write surface).
  The gating site for "this node may offer this mode" is t1's
  render — POST is defense-in-depth, not the gate.
- **(c) Skip validation entirely; rely on the loopback
  boundary** (anyone on loopback is the contributor; the
  contributor wouldn't hand-craft a malformed POST). Removes
  the defense-in-depth layer the m2 Cross-Task Risk explicitly
  calls for ("the endpoint's inputs (slug, mode) are
  **defense-in-depth-validated** against the slug grammar and
  the D3 mode-offer map").

**Chosen: (a).** Grammar + mode-set membership is the
load-bearing defense-in-depth the m2 Risk names. Skipping the
plan-tree walk preserves the read-side single-walk discipline.
Status-shape validation ("this slug's node Status currently
matches the mode's D3 row") is **not** in scope — that's t1's
gating site and a hand-crafted POST for a stale (slug, mode)
pair still satisfies the m2 trust boundary on loopback. The
orphan-attachment residual (a wrong-by-construction slug
attaches a tree-invisible work-instance) is the same residual
m1 already accepted per
[m1 README, "wrong-by-construction slug attaches an orphan
work-instance"](../../m1/README.md) — m2 does not re-litigate
it. *Verified by:* [m2 Cross-Task Risk, "New `/spawn` endpoint
expands the tool's write surface"](../README.md) (the
defense-in-depth language); [m2 Cross-Task Risk,
"Wrong-by-construction slug …"](../README.md) (the orphan
residual the m1 trust-the-caller posture covers).

### SD9 — Child env: set `WST_SLUG` + `WST_NAME` + pass `claude --name` to align workstream-tracker roster and Claude agent-view labels

How D1 and D5 surface in the child's environment and `claude`
argv.

**Shapes considered:**

- **(a) Set `WST_SLUG=<slug>` and `WST_NAME=<worktree-name>` in
  `cmd.Env`; pass `--name <worktree-name>` in the `claude` argv
  alongside `--bg --worktree <worktree-name>`.** The
  workstream-tracker register CLI reads `WST_NAME` as the
  optional `--name` enrichment metadata (PR
  [#38](https://github.com/kcrobinson-1/workstream-tracker/pull/38)),
  surfacing it as the roster's display label. Claude Code's
  `--name` sets the agent-view session display name. Both
  point at the same human-readable string (the worktree name
  per SD5), keeping the cross-tool correlation tight.
- **(b) Set `WST_SLUG` only; leave `WST_NAME` unset and omit
  `claude --name`.** Roster falls back to slug-as-label;
  agent-view auto-generates from the prompt. Loses the
  enrichment leg D5 directs the spawn to ride.
- **(c) Set `WST_NAME` to a richer string** (e.g.,
  `<slug> — planning`) different from the worktree name.
  Diverges from D5's "pass that same name through the existing
  `--name`/`WST_NAME` enrichment leg" wording.

**Chosen: (a).** D5 directs WST_NAME = worktree name; the
parallel `claude --name <same>` aligns the agent-view label.
Both fields decorate an already-correctly-attached work-
instance per D5's best-effort framing; a missing enrichment
degrades richness, not correctness. *Verified by:*
[`cmd/workstream-tracker/register.go`](../../../../../cmd/workstream-tracker/register.go)
`runRegister` (the `WST_NAME`/`--name` read path: lines
44–73, the `WST_NAME` env fallback and JSON-metadata
construction); [m2 D5](../README.md) (the "pass that same name
through the existing `--name`/`WST_NAME` enrichment leg"
clause); [Claude Code CLI reference,
`--name`/`-n`](https://code.claude.com/docs/en/cli-reference#cli-flags)
(the agent-view display name).

### SD10 — Initial task argv: short imperative carrying slug + mode (e.g., `"Plan tool-originated-task-sessions-m2-t2."`)

The positional argv string the `claude` invocation accepts as
its first-user-message anchor.

**Shapes considered:**

- **(a) Short imperative with slug + mode verb:** `"Plan
  <slug>."` for planning spawns; `"Implement <slug>."` for
  implementation spawns. Anchors the session's first user
  message; the bulk of the per-mode instructions lives in
  `--append-system-prompt-file` (the system-prompt persona).
- **(b) Empty positional argument** (or omit it entirely).
  Risks a `claude --bg` requirement for a non-empty initial
  prompt; the vendor example always passes one.
- **(c) Long-form prompt body in the positional argument.**
  Trips argv-length limits on long enough bodies; mixes
  system-prompt content with user-message content.

**Chosen: (a).** The split between system prompt (mode
persona, via `--append-system-prompt-file`) and user message
(the specific task, via the positional argv) tracks Claude
Code's own model: system prompt = persona, user message = task.
The exact wording of the initial task is HOW for the
implementing PR per [m2 Out of Scope, "The exact
wording"](../README.md). *Verified by:* [Claude Code CLI
reference, `claude "<query>"`
example](https://code.claude.com/docs/en/cli-reference#cli-commands)
(`claude --bg "investigate the flaky test"` — a positional
imperative initial prompt).

### SD11 — Server dependency wiring: inject the exec function on `Server` for testability

How the handler's `os/exec` dependency is structured so unit
tests can drive it without invoking a real `claude` binary.

**Shapes considered:**

- **(a) Add an `execCommand` field on `*site.Server`** (type
  `func(name string, args ...string) *exec.Cmd`) that defaults
  to `exec.Command` and that tests override with a fake. Tests
  verify the argv and env construction; an integration test
  with a stubbed `claude` script on `PATH` exercises the
  stdout-parse path.
- **(b) No injection; tests invoke a real `claude` binary or
  skip the test surface entirely.** The implementing PR's
  validation gate loses unit-level coverage of the argv shape;
  the orphan-attachment risk loses a falsifier.
- **(c) Refactor the exec wrapper into a new package** with
  its own interface. Over-engineered for a single handler.

**Chosen: (a).** Adding one function field on `Server` is the
minimal injection — same pattern as Go's standard `http.Server`
listener override or `slog`'s handler wrapping. Tests construct
a `*site.Server` with an `execCommand` that returns a deterministic
`*exec.Cmd` configured via Go's standard `helperCommand`
pattern (a no-op os/exec of the test binary itself with a
trigger flag), or via the documented `exec.Cmd.Args` /
`Cmd.Env` introspection without actually running the cmd.
*Verified by:* Go standard-library `os/exec.Cmd` field
semantics; the existing `internal/site` Server constructor
shape ([`internal/site/site.go`](../../../../../internal/site/site.go)
`New(db, plansPath)` — extending with `execCommand` is
additive).

## Open decisions to make at plan-drafting

Open questions the human resolves at the
`In draft` → `Proposed` promotion gate. Both are vendor-
behavior facts the implementing PR will need to verify
empirically before writing the dependent code; they are
**not** decided by this scoping pass.

### OQ1 — Does the supervisor pass invoker-set environment variables (`WST_SLUG`, `WST_NAME`) through to the spawned background session and its SessionStart hook?

**Background.** D1's slug-carry mechanism rests on the
spawned session — and therefore its SessionStart hook — seeing
`WST_SLUG` set in the `/spawn` handler's `cmd.Env`. The flow
the design assumes:

```
/spawn handler                    Claude Code supervisor                 spawned session
  cmd.Env = parent + WST_SLUG ── claude --bg ──>   ?  ──>   session process w/ WST_SLUG
                                                            │
                                                            └── SessionStart hook process
                                                                   inherits WST_SLUG
```

The middle `?` is the unknown. The supervisor is a per-user
long-lived process started by the *first* `claude --bg`
invocation, not by *this* `claude --bg` invocation; subsequent
invocations communicate with the existing supervisor via a
socket. Whether the supervisor passes the new invoker's env
vars through to the new session is **not specified** in vendor
docs. The relevant clauses:

- Hooks doc: *"Handlers run in the current directory with
  Claude Code's environment."* — ambiguous as to *which* Claude
  Code's environment when the session is supervisor-spawned.
- Agent-view doc: *"Configuration flags from the original
  launch carry through to the backgrounded session"* — names
  *flags*, refers to `/bg` backgrounding (an interactive
  session moving to background), not `claude --bg` spawning.

**Resolution options:**

- **(A) Assume yes; implementing PR verifies first.** The
  implementing PR's first step is a `bash`-level reality
  check: `WST_SLUG=test claude --bg "echo \$WST_SLUG"`
  followed by `claude logs <id>` to confirm. If yes, proceed
  with the rest of the implementation; if no, escalate to (B).
- **(B) Use `--settings <per-spawn-json>` as the carrier
  instead of an env var.** The `/spawn` handler synthesizes a
  per-spawn temp `.json` file whose `hooks.SessionStart`
  command is `go run … register --slug <literal-slug>` (slug
  baked into the command, not read from env); the handler
  passes `--settings <path>` to `claude --bg`. This bypasses
  the env-var question entirely. The cost: t3's design shifts
  — the static committed `.claude/settings.json` becomes the
  no-`WST_SLUG`-set no-op shell, and the per-spawn `--settings`
  overlay is the workhorse. **This shift is t3's surface, not
  t2's** — but t2's `/spawn` would be the producer of the
  per-spawn JSON.
- **(C) Run register from `/spawn` directly, before invoking
  `claude --bg`, using the actor identity of the
  workstream-tracker server process.** Rejected at this
  scoping: it changes the work-instance's actor from the
  spawned session's actor to the workstream-tracker process's
  actor, breaking m1's "actor identifies which agent is doing
  work" contract.

**Why this is OPEN, not chosen here.** The (A) → (B) fallback
chain is a real design change (t3's surface shifts under (B)).
Picking (A) without empirical verification would commit t2 to
an assumption that, if wrong, requires a milestone-level
reconciliation. The implementing PR is the right place to
run the verification and trigger the fallback if needed; the
plan-drafting session should not commit to (A) without that
verification. The promotion gate may resolve this by either
(i) running the verification as a spike now (which would
spawn a real Claude Code session — risky from a drafting
session), or (ii) explicitly authorizing the
"implementing PR verifies first, escalates to (B) on failure"
path with the t3 contract reconciliation pre-negotiated.

### OQ2 — When `claude --bg --worktree <name>` is invoked and `<repo>/.claude/worktrees/<name>` already exists, what happens?

**Background.** SD5 chose worktree name = `<slug>-<mode-suffix>`,
which collides deterministically on re-spawn of the same (slug,
mode). The collision is *observable* either way per the
"deterministic ≠ cannot fail" Cross-Task Invariant — but the
specific user-facing failure shape matters for SD7 (the
acknowledgement page) and the m2 Risk Register entry on
hook-misconfiguration.

**Vendor doc gap.** The CLI reference says only "If no name is
given, one is auto-generated." Behavior for a given-name
collision is not documented. Three plausible vendor behaviors:

- **(a) Error out before spawning** with a message like
  `worktree "<name>" already exists`. SD7's acknowledgement
  page renders the error; the contributor `git worktree
  remove`s the stale worktree and re-spawns.
- **(b) Reuse the existing worktree** (no isolation regression
  if the prior session is gone). SD7 still works; the spawn
  succeeds; the worktree's current state from the prior
  session is the new session's starting point.
- **(c) Auto-append a numeric suffix** to disambiguate. SD7
  works but D5's "name known at spawn time" property bends —
  the `--worktree <name>` argv isn't the *resulting* name.

**Resolution options:**

- **(A) Let Claude Code's behavior speak**, whichever of
  (a)/(b)/(c) it is. The handler returns the raw stdout/stderr
  to the contributor; SD7's acknowledgement page renders both
  on success and on failure paths. The implementing PR
  verifies which of (a)/(b)/(c) holds and writes the
  acknowledgement page accordingly.
- **(B) Pre-flight detection in the `/spawn` handler.** Before
  invoking `claude`, the handler checks for the directory's
  existence and either errors with a workstream-tracker-side
  message or appends a per-request suffix. Adds filesystem
  introspection logic that the vendor's own behavior may
  duplicate or contradict.

**Why this is OPEN, not chosen here.** (A) is the right
default but the user-facing failure shape (what the
acknowledgement page says when the spawn fails) depends on
which vendor behavior applies. The implementing PR verifies and
writes the acknowledgement page; the contract is "the
collision is observable to the contributor", which SD5 + SD7
already guarantee. The promotion gate may collapse OQ2 by
explicitly authorizing the "implementing PR verifies; SD7
renders the raw vendor failure" path without pre-committing to
(a)/(b)/(c).

## Plan structure handoff

The paired plan doc replaces the seeded skeleton at
[`../t2-spawn-endpoint.md`](../t2-spawn-endpoint.md) with
**Status: `In draft`**. Per the project rule on spawned
drafting, the plan stays at `In draft` (does **not** flip to
`Proposed`) — the human resolves OQ1 + OQ2 at the promotion
gate before flipping. Section structure as drafted:

- **Context** preamble naming what the plan covers, why now,
  and what surfaces it touches at the conceptual level (per
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "Plan opens with a plain-language context preamble").
- **Goal.**
- **Cross-Cutting Invariants** consumed by reference from the
  m2 README; one t2-local invariant added (the canonical-argv
  + no-shell-interpolation rule the m2 Cross-Task Risk
  realizes at the handler-internal level).
- **Contracts C1…C8** covering: chi-router mount + handler
  signature (C1); loopback-only binding ships in the same
  change (C2); exec invocation shape + env (C3); slug-grammar
  + mode-set defense-in-depth validation (C4); per-mode prompt
  files exist at `docs/spawn-prompts/` (C5); session-id capture
  + acknowledgement page (C6); fire-and-forget posture / the
  workstream-tracker server is not a process manager (C7);
  preserved invariants (no-JS, no GET-side walk regression,
  agent-adapter seam additive, observe-only for non-originated
  sessions) (C8).
- **Files to touch (estimate)** —
  [`internal/site/site.go`](../../../../../internal/site/site.go)
  (router mount + Server dependency extension);
  [`internal/site/spawn.go`](../../../../../internal/site/) (new —
  the handler + supporting types) +
  [`internal/site/spawn_test.go`](../../../../../internal/site/)
  (new — the unit + integration tests);
  [`cmd/workstream-tracker/main.go`](../../../../../cmd/workstream-tracker/main.go)
  (loopback-only `addr` + `Server` wiring with `SPAWN_PROMPTS_DIR`);
  `docs/spawn-prompts/planning.md` and
  `docs/spawn-prompts/implementation.md` (new — the per-mode
  prompt bodies whose **wording is the implementing PR's**, per
  m2 Out of Scope). Estimate labeled per the "Plan content is a
  mix of rules and estimates" rule.
- **Validation Gate** — the technical gate (unit tests against
  the injected exec function; an integration test against a
  stub `claude` script that exercises stdout parsing;
  loopback-binding `Addr` check; defense-in-depth 400 cases);
  the manual single-spawn sanity check against a real `claude`
  if available locally; the OQ1 verification step the
  implementing PR runs first. The product walkthrough lives on
  t4 (interior-node node-not-leaf shape per
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "Product-facing leaf tasks").
- **Self-Review Audits** — no-shell-interpolation audit
  (search the diff for any `sh -c` or string-formatted shell
  command); no-new-walk-on-GET audit (the new POST handler
  does not call `walkPlans` on the GET path); preserved-
  surface audit (the existing index render path is unchanged).
- **Documentation Currency** — internal-package surfaces
  inside [`internal/site/`](../../../../../internal/site/) and
  [`cmd/workstream-tracker/main.go`](../../../../../cmd/workstream-tracker/main.go);
  no spec, agent-rule, or backlog edit (per m2 Documentation
  Currency for t2 + Backlog Impact).
- **Out of Scope** — t1's affordance form (Landed); t3's
  `.claude/settings.json` entry; t4's product walkthrough; the
  exact wording of the per-mode prompt bodies; the session-end
  completion bracket (PR #45); any backlog edit; any in-session
  tool→agent flow.
- **Risk Register** — the t2-side mitigations of the m2-level
  "`/spawn` expands the write surface" risk; OQ1's potential
  cascade to t3 if env-var pass-through fails; OQ2's
  user-facing acknowledgement-shape implications.
- **Open Questions** — OQ1 and OQ2 surfaced verbatim from this
  scoping doc (the human resolves them at the gate, not the
  drafting session).
- **Backlog Impact** — none (m2 README locks this at
  milestone level).
- **Related Docs.**

## Related Docs

- [`../t2-spawn-endpoint.md`](../t2-spawn-endpoint.md) — the
  paired durable plan doc this scoping doc fuels (replaces the
  seeded skeleton).
- [`../README.md`](../README.md) — m2 milestone doc; the locked
  WHAT, D1/D4/D5, Cross-Task Invariants, Cross-Task Risks.
- [`../t1-mode-affordance-render.md`](../t1-mode-affordance-render.md) —
  Landed sibling task; produces the affordance form `/spawn`
  consumes.
- [`../t3-session-start-hook.md`](../t3-session-start-hook.md) —
  sibling task skeleton; OQ1's resolution may shift t3's
  static-hook design.
- [`../../README.md`](../../README.md) — parent epic; the
  Cross-Cutting Invariants and resolved vision inputs the m2
  contract sits on top of.
- [`../../../../../spec/planning/task-plan.md`](../../../../../spec/planning/task-plan.md) —
  this scoping doc's authority; "Scoping owns / plan owns,"
  "Reality-check pass before plan-drafting."
- [`../../../../../spec/planning/shared.md`](../../../../../spec/planning/shared.md) —
  cross-level rules; "Decompose options into shapes,"
  "Verified by: annotations," "Plans describe contracts, not
  implementation."
- [Claude Code CLI
  reference](https://code.claude.com/docs/en/cli-reference) —
  the vendor surface D4 / D5 verify against.
- [Claude Code agent-view
  documentation](https://code.claude.com/docs/en/agent-view) —
  the supervisor + worktree + `--bg` output-format surface.
- [Claude Code hook
  documentation](https://code.claude.com/docs/en/hooks) —
  the SessionStart hook contract t3's surface satisfies.
