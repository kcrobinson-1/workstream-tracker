---
slug: tool-originated-task-sessions-m2-t3
Status: In draft
short_description: SessionStart hook — .claude/settings.json entry running the deterministic register CLI
---

# m2 t3 — SessionStart Hook: `.claude/settings.json` Entry

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
Cross-Task Invariants, Cross-Task Decisions D1/D2,
Cross-Task Risks "Hook misconfiguration silently misses …" and
"The committed hook fires for every Claude Code session in
this repo …"); not loosened here.

**End result.** A **Claude Code `SessionStart` hook** (matcher
`startup`, `type: "command"`) is committed to a
**project-scoped** `.claude/settings.json` at the repo root
with a command that invokes the deterministic
`workstream-tracker register` subcommand (**no `--slug` flag**
— env-only slug carry via the CLI's existing `WST_SLUG`
reading; see D2 for why the no-flag form is load-bearing
rather than stylistic) via the agent rule's full module-path
form (so the hook resolves regardless of the contributor's
working directory). The hook fires at every Claude Code
session start in this repo.

**When `WST_SLUG` is set in the spawned session's
environment** (the case t2's spawn produces), the register CLI
attaches a work-instance to the carried slug by construction
before the model reasons — m1's deterministic path runs by
hook, not by agent narration.

**When `WST_SLUG` is unset** (the case for a contributor
opening Claude Code in this repo through any path other than
t2's spawn), the register CLI **short-circuits to a no-op**
per `runRegister`'s existing "no slug supplied … skipping
registration, session proceeds" branch — reached because
`flag.Parse` succeeds (no `--slug` flag to argument-error
against) and the resolved slug is empty. So existing
interactive natural-language sessions, the m1-landed
best-effort grounded narration handshake, and any
contributor's per-`.claude/settings.local.json` overlay all
see **byte-unchanged session-start behavior**.

The hook is the determinism-relevant integration moment t2's
spawn relies on.

**Preserves.** The m1-landed
[`docs/agents/local/session-registration.md`](../../../../docs/agents/local/session-registration.md)
deterministic-path rule is *consumed*, not edited; the
existing register CLI is consumed verbatim; the existing
interactive handshake is unaffected; no new registration code
path, endpoint, schema change, or CLI flag.

Whether the agent rule needs a small additive clarification —
that for tool-originated sessions the *hook* runs the
deterministic-path "Invoke" step and the agent's first
interaction is the "Echo real output" step — is HOW for t3's
own planning session against then-merged rule wording.

**Sibling interface.** Produces the deterministic-register
integration moment t2's spawn relies on (without this hook in
place, t2's spawn still launches `claude` but no work-instance
attaches — the slug-by-construction contract is half-wired).
Consumes m1's deterministic register CLI verbatim.

**Product acceptance.** Closes on technical gate; the
milestone's product validation lives on t4. (Interior node —
not a Mermaid-graph leaf.)

## What this skeleton does NOT cover

HOW. The exact JSON shape of the `.claude/settings.json` hook
entry, whether any small additive clarification to
[`session-registration.md`](../../../../docs/agents/local/session-registration.md)
is needed (the milestone names this as a HOW judgment t3
makes), the falsifiability shape of the no-op-branch /
register-branch tests, and the risk register — produced by
t3's own `In draft` → `Proposed` drafting session against
then-merged code, per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns" and
[`milestone.md`](../../../../spec/planning/milestone.md)
"Anti-goal: WHAT-contract each task, do not HOW-scope any task."

## Related Docs

- [`./README.md`](./README.md) — m2 milestone doc; the locked
  WHAT, Cross-Task Invariants/Decisions (especially D2),
  Cross-Task Risks (the two hook-relevant risks), Documentation
  Currency.
- [`./t2-spawn-endpoint.md`](./t2-spawn-endpoint.md) — sibling
  task; the spawn that produces the `WST_SLUG`-set environment
  the hook reacts to.
- [`../m1/t1-deterministic-path-contract.md`](../m1/t1-deterministic-path-contract.md) —
  the landed deterministic-path contract this hook consumes.
- [`../../../../docs/agents/local/session-registration.md`](../../../../docs/agents/local/session-registration.md) —
  the lifecycle rule; "The deterministic path
  (construction-known slug)" is the runtime contract the hook
  satisfies.
- [Claude Code `SessionStart` hook
  documentation](https://code.claude.com/docs/en/hooks.md) —
  D2 vendor grounding.
