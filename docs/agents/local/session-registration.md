# Session work-instance lifecycle handshake

Repo-owned universal rule. The root [`AGENTS.md`](../../../AGENTS.md)
"Universal session rules" section points here; this file owns the
detail. Lives in `AGENTS.md` + `docs/agents/local/` (not the
vendored read-only `docs/agents/shared/**`).

This rule owns the full session work-instance lifecycle: the
session-start registration handshake and the symmetric session-end
completion handshake. The completion half is not a separate rule —
it is the closing bracket of the same observable best-effort
narration contract; an open registration with no completion leaves
the tree showing work that has finished.

## Why this rule exists

The cross-agent visualization is only trustworthy if a session
that starts work actually appears in the tree. Nothing registers a
session automatically and silently: turning a natural-language
prompt into a canonical slug is interpretation that only the agent
can do, mid-session. So registration is an **observable
best-effort grounded narration handshake** — loud and fact-checked,
not silent.

## The handshake (interactive natural-language session)

Before doing task work, perform this handshake. It never blocks or
fails the session — whatever happens, you proceed.

1. **Resolve.** Resolve the canonical slug from the prompt and the
   plan tree (the plan doc's `slug` frontmatter is authoritative).
2. **Confirm.** State the resolved plan-doc path and the resolved
   slug so a present human can catch a wrong target before you
   proceed.
3. **Echo real output.** Invoke the registration command and
   report the *actual* command/server response — the real
   work-instance id and HTTP status it printed. A "registered
   successfully" claim with no echoed id/status is a
   `validation-honesty` violation: report observed fact, not a
   prose success claim.

   ```sh
   go run github.com/kcrobinson-1/workstream-tracker/cmd/workstream-tracker register --slug <canonical-slug>
   ```

   Invoke it through the full module path. The repo ships no
   installed `workstream-tracker` binary on `PATH`, so a bare
   `workstream-tracker …` will not resolve; and a cwd-relative
   `go run ./cmd/...` only resolves from the repo root, so it
   breaks once the session has `cd`'d into a package directory.
   The module path resolves against the current module from any
   directory in the checkout.

   (Slug also via `WST_SLUG`; actor via `--actor`/`WST_ACTOR`,
   defaulting to a generated per-session id — never the git user;
   server via `--server`/`WST_SERVER`, default
   `http://localhost:8080`.)
4. **Narrate failure explicitly.** If the command reports a
   failure, or the slug is unresolvable, say so explicitly and
   actionably — the session will not appear in the tree; the
   contributor can run the command by hand against a running
   server. A silent skip is an `error-surfacing-user-mutations`
   violation.
5. **Proceed.** Registration never gates task work. Continue
   whether it succeeded, failed, or was skipped.

## The session-end handshake (completion)

When the session's work is finished, close the bracket. Same
best-effort contract: it never blocks or fails the session.

1. **Invoke.** Run the completion subcommand. It resolves the
   work-instance id from the receipt the session-start register
   cached in this worktree — no id threading required.

   ```sh
   go run github.com/kcrobinson-1/workstream-tracker/cmd/workstream-tracker complete
   ```

   (Use `abandon` instead of `complete` if the work is being
   dropped rather than finished. Id also via `--id`/`WST_WI_ID`
   when the register receipt was not cached — e.g. a separate
   process or worktree; server via `--server`/`WST_SERVER`,
   default `http://localhost:8080`.)
2. **Echo real output.** Report the actual receipt — the real
   event id and HTTP status it printed. A prose "marked complete"
   with no echoed id/status is a `validation-honesty` violation.
3. **Narrate failure explicitly.** If the command reports a
   failure, or no work-instance id is resolvable, say so
   explicitly and actionably — the tree will keep showing this
   session active; the contributor can run the command by hand
   against a running server. A silent skip is an
   `error-surfacing-user-mutations` violation.

## Scope and residual

This handshake is the only observability backstop for a missed
registration: an unregistered session emits no signal the tool
ever sees, so the rendered tree cannot distinguish it from
genuinely no work. A present contributor noticing a missing or
incongruent handshake line is the backstop. A missed session-end
completion is symmetric: the tree keeps showing the session active
until the contributor notices the stale marker. The determinism and
tree-side-affordance gaps are consciously accepted best-effort,
tracked by the `deterministic-interactive-registration` backlog
entry.

### The cached-receipt contract

`complete`/`abandon` resolve the work-instance id from a receipt
`register` caches per worktree. That cache is **owner-scoped**: it
records the registering actor alongside the id, and a terminal
command uses it only when the actor it resolves matches. This is
what keeps a stale or another session's receipt from being marked
terminal — the failure class the implementation contract in
`cmd/workstream-tracker/register.go` ("The wi-cache contract")
states in full. Practical consequences for a session:

- **Pin the actor across both halves.** With no explicit actor a
  per-session id is generated and reused within the worktree; if
  the harness sets `WST_ACTOR` (or you pass `--actor`), set the
  *same* value for `complete` as for `register`, or pass `--id`.
  A mismatch makes `complete` safely skip and say so — it never
  marks the wrong work-instance terminal.
- **Accepted residuals** (same boundary as the backlog entry
  above, not engineered away): two sessions sharing one resolved
  actor in one worktree still collapse onto one receipt; an
  explicit `--id`/`WST_WI_ID` is honored without the owner check
  as a deliberate operator override.
