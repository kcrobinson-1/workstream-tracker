---
slug: tool-originated-task-sessions-m2-t2
Status: In draft
short_description: Spawn integration — /spawn endpoint + claude --bg exec
---

# m2 t2 — Spawn Integration: `/spawn` Endpoint + `claude --bg` Exec

> **Parent-promotion skeleton — scope locked, HOW not yet
> drafted.** Seeded by the m2 milestone-planning session's
> `In draft` → `Proposed` promotion gate per
> [`shared.md`](../../../../spec/planning/shared.md) "Parent-doc
> child contracts → Parent-promotion stub seeding." The locked
> WHAT contract this skeleton carries comes from the [m2
> milestone doc](./README.md) Task Contracts row. Per the
> seeding rule the skeleton is **exempt from the
> [`task-plan.md`](../../../../spec/planning/task-plan.md)
> "Required and optional sections" rule and from
> [`shared.md`](../../../../spec/planning/shared.md) "Plans
> describe contracts, not implementation"** until its own
> drafting session opens — at which point that exemption ends
> and the doc grows into a full task plan.

## Inherited Contract (locked by the m2 milestone doc)

Binding input from [`./README.md`](./README.md) (Task Contracts,
Cross-Task Invariants, Cross-Task Decisions D1/D2/D4/D5,
Cross-Task Risks); not loosened here.

**End result.** A new `POST /spawn` endpoint on the
workstream-tracker local server accepts the slug + mode
submitted by t1's affordance form and **fire-and-forget execs
the Claude Code launcher in background-session mode** with the
slug carried out-of-band via `WST_SLUG` (D1) and a
mode-appropriate prompt body handed in via
`--append-system-prompt-file`. The launcher invocation uses
Claude Code's `--bg` (background session — Claude Code's
supervisor process owns the agent's process lifetime; the
workstream-tracker does **not** become a process manager) and
`--worktree <name>` (Claude Code provisions a fresh worktree
at `<repo>/.claude/worktrees/<name>` automatically — D5). The
endpoint captures the printed session id from stdout and
acknowledges it back to the page; the contributor's takeable
session is `claude attach <id>` in their own terminal.

**The same change that mounts `/spawn` also constrains the
server binding to a loopback address.** The current
`addr := ":" + port` in
[`cmd/workstream-tracker/main.go`](../../../../cmd/workstream-tracker/main.go)
`runServer` binds all interfaces; t2 changes it to a loopback
`Addr` (typically `127.0.0.1:<port>`; exact `Addr` spelling
HOW for this task's planning). The loopback binding is the
load-bearing security premise that D4 and the
`/spawn`-write-surface Cross-Task Risk both rest on; without
it at `/spawn`-mount, anyone reaching the host on the LAN
could fire `POST /spawn`.

The spawn **assumes t3's SessionStart hook is in place**: with
the hook present, the spawned session auto-registers via m1's
deterministic path before the model reasons; without it, the
spawn still launches `claude` but no work-instance attaches —
observable by the tree's emptiness for the launched slug. t2
does not commit the hook entry itself (that is t3's surface).

**Preserves.** The deterministic register CLI is unmodified
(the existing `--slug` argument and `WST_SLUG` env var are
honored verbatim); the interactive best-effort handshake for
natural-language sessions is unmodified; observe-only is
preserved for every session the tool did not originate; no
new endpoint is added to the **registration** surface (the
new `/spawn` endpoint is the *launcher* surface — orthogonal
to registration; the registration path is the unchanged
exact-slug create-or-attach flow). The fencing properties
(human-initiated, one-shot at birth, identity-not-correction,
running session stays observe-only) hold at every site t2
touches: the form submit is the human-initiated trigger; the
spawn is one-shot at session birth; the slug carried is
identity-not-correction; the *running* session stays
file-driven and observe-only — `/spawn` is a one-shot launch
surface, not an in-session steering channel. The
agent-adapter seam stays additive: an alternative agent
launcher is a new exec shape inside the handler, not a
rewrite.

**Sibling interface.** Consumes t1's mode-affordance form
submission as the launch trigger; consumes t3's SessionStart
hook as the deterministic-register integration moment;
produces the **real construction-time slug producer** that
drives m1's slug-carried registration path end-to-end.

**Product acceptance.** Closes on technical gate; the
milestone's product validation lives on t4. (Interior node —
not a Mermaid-graph leaf.)

## What this skeleton does NOT cover

HOW. The chi-router mount point, the exact `os/exec` argv
shape, the per-mode prompt-file paths, the response shape the
page reads to surface the session id, the defense-in-depth
input-validation specifics, and the risk register — all
produced by t2's own `In draft` → `Proposed` drafting session
against then-merged code, per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns" and
[`milestone.md`](../../../../spec/planning/milestone.md)
"Anti-goal: WHAT-contract each task, do not HOW-scope any task."

## Related Docs

- [`./README.md`](./README.md) — m2 milestone doc; the locked
  WHAT, Cross-Task Invariants/Decisions, Documentation
  Currency.
- [`./t3-session-start-hook.md`](./t3-session-start-hook.md) —
  sibling task; the SessionStart hook that t2's spawn depends
  on for the deterministic register to fire.
- [`../README.md`](../README.md) — parent epic; Cross-Cutting
  Invariants and the resolved vision inputs.
- [Claude Code CLI
  reference](https://code.claude.com/docs/en/cli-reference.md) —
  D4 / D5 vendor grounding.
