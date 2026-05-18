# Session-start work-instance registration

Repo-owned universal rule. The root [`AGENTS.md`](../../../AGENTS.md)
"Universal session rules" section points here; this file owns the
detail. Lives in `AGENTS.md` + `docs/agents/local/` (not the
vendored read-only `docs/agents/shared/**`).

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
   workstream-tracker register --slug <canonical-slug>
   ```

   (Slug also via `WST_SLUG`; actor via `--actor`/`WST_ACTOR`,
   defaulting to a generated per-session id — never the git user;
   server via `--server`/`WST_SERVER`, default
   `http://localhost:8080`. Optionally report a human session
   name via `--name`/`WST_NAME` — it becomes the session's roster
   label; omitting it lists the session under its slug. The name
   is the sole conventionally-read reported key; it is sent as
   request metadata and never blocks the best-effort attempt.)
4. **Narrate failure explicitly.** If the command reports a
   failure, or the slug is unresolvable, say so explicitly and
   actionably — the session will not appear in the tree; the
   contributor can run the command by hand against a running
   server. A silent skip is an `error-surfacing-user-mutations`
   violation.
5. **Proceed.** Registration never gates task work. Continue
   whether it succeeded, failed, or was skipped.

## Scope and residual

This handshake is the only observability backstop for a missed
registration: an unregistered session emits no signal the tool
ever sees, so the rendered tree cannot distinguish it from
genuinely no work. A present contributor noticing a missing or
incongruent handshake line is the backstop. The determinism and
tree-side-affordance gaps are consciously accepted best-effort,
tracked by the `deterministic-interactive-registration` backlog
entry.
