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

## The deterministic path (construction-known slug)

When a session's canonical slug is **known by construction** —
carried into the session rather than resolved from a
natural-language prompt — the interpretation the handshake above
exists for has already been done. This is a sibling of the
registration handshake, not a replacement for it, and it changes no
registration endpoint, request field, schema, or code path: it is
the same exact-slug create-or-attach path, asserted with a slug
that needs no interpreting.

Steps 1–2 (**Resolve**, **Confirm**) **do not apply**: there is no
prompt to interpret into a slug, and there is **no
confirm-equivalent** — no pre-invoke human-catch step. Nothing is
being interpreted, so there is nothing for a present human to catch
before invocation; the real-receipt echo below is the sole
observability surface. The remaining steps are unchanged and still
best-effort — it never blocks or fails the session:

1. **Invoke.** Run the registration command with the
   construction-known slug. The manual / CLI `--slug` argument is
   the stand-in producer of such a slug.

   ```sh
   go run github.com/kcrobinson-1/workstream-tracker/cmd/workstream-tracker register --slug <construction-known-slug>
   ```

   Invoke it through the full module path, with the same
   environment-variable and `--actor`/`--server`/`--name`
   equivalents and the same rationale as the interactive handshake
   above.
2. **Echo real output.** Report the *actual* command/server
   response — the real work-instance id and HTTP status it printed.
   The same `validation-honesty` obligation applies: report
   observed fact, not a prose success claim. This real-receipt echo
   is the sole observability surface, catching a
   wrong-by-construction slug *post-hoc*, not pre-invoke.
3. **Narrate failure explicitly.** If the command reports a
   failure, say so explicitly and actionably — the session will not
   appear in the tree; the contributor can run the command by hand
   against a running server. A silent skip is an
   `error-surfacing-user-mutations` violation.
4. **Proceed.** Registration never gates task work. Continue
   whether it succeeded, failed, or was skipped.

Because identity is fixed by construction, this is a deterministic,
handshake-free registration path. "Deterministic" scopes only to
identity being fixed by construction for a session that *did*
launch — not that the producer of the slug cannot itself be wrong.
A slug that is wrong *by construction* (a typo'd `--slug` today)
attaches an orphan work-instance: the already-accepted, deferred
orphan/unattached-work-instance residual, kept observable post-hoc
by the real-receipt echo.

## The session-end handshake (completion)

When the session's work is finished, close the bracket. Same
best-effort contract: it never blocks or fails the session.

1. **Invoke.** Run the completion subcommand, passing the
   `work_instance_id` the session-start register receipt printed
   (carry that value forward from step 3 of the start handshake).

   ```sh
   go run github.com/kcrobinson-1/workstream-tracker/cmd/workstream-tracker complete --id <work_instance_id>
   ```

   (Use `abandon` instead of `complete` if the work is being
   dropped rather than finished. Id also via `WST_WI_ID`; server
   via `--server`/`WST_SERVER`, default `http://localhost:8080`.)
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

The completion id is supplied explicitly: `complete`/`abandon`
take the `work_instance_id` from the register receipt via
`--id`/`WST_WI_ID`. There is no hidden cross-invocation state — the
session (or harness) carries the id from the start handshake to the
end one, the same value the start handshake already echoed.
