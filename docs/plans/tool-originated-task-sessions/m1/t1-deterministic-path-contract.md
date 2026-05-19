---
slug: tool-originated-task-sessions-m1-t1
Status: Proposed
short_description: Deterministic-path contract expression
---

# m1 t1 — Deterministic-Path Contract Expression

## Context preamble

Today the only *registration* procedure written down anywhere is
the **interactive best-effort grounded narration handshake** in
[`docs/agents/local/session-registration.md`](../../../../docs/agents/local/session-registration.md)
(since PR #45 a *lifecycle* rule that also documents a symmetric
session-end completion handshake — orthogonal to t1): an agent
resolves a canonical slug from a natural-language prompt, confirms
it, invokes the register subcommand, echoes the real receipt, and
narrates any failure. That handshake exists because turning prose
intent into a slug is interpretation only the agent can do
mid-session.

But the registration *mechanism* in the merged binary is already
deterministic when the slug does not need interpreting: the
exact-slug create-or-attach flow honors a caller-supplied slug
verbatim with no natural-language resolution, and the
`workstream-tracker register --slug <slug>` subcommand drives
exactly that. When a session's canonical slug is known **by
construction** (carried in, not derived from a prompt), none of
the handshake's interpretation steps apply — yet no durable
artifact says so. A reader of the spec or the agent rule would
reasonably conclude the narration handshake is the *only*
registration path.

t1 closes that gap: it expresses the deterministic,
handshake-free, slug-by-construction path as a **first-class
additive contract** in the plan-doc spec and the
session-registration agent rule, so the m2 tool-originated-session
UX milestone has a sanctioned, durable contract to consume when
its spawn becomes the real construction-time slug producer. t1
ships **no code** and **no new registration surface** — the
mechanism already exists (milestone Cross-Task Decision D1); the
deliverable is the contract expression itself.

The HOW decisions were resolved at scoping (SD1–SD7 in the
paired, transient scoping doc
[`scoping/t1-deterministic-path-contract.md`](./scoping/t1-deterministic-path-contract.md));
this plan carries the durable contract those decisions produced,
with its own code-grounded verification.

## Inherited Contract (locked by the milestone)

Binding input from [`../README.md`](../README.md) (Task
Contracts, Cross-Task Invariants/Decisions); not loosened here.

**End result.** The deterministic, handshake-free,
slug-by-construction registration path is expressed as a
first-class additive contract in
[`spec/planning/shared.md`](../../../../spec/planning/shared.md)
and
[`docs/agents/local/session-registration.md`](../../../../docs/agents/local/session-registration.md):
when a session's canonical slug is known by construction (carried
in, not resolved from a natural-language prompt), registration is
deterministic and the narration handshake does not apply.

**Preserves, unchanged:** the interactive best-effort grounded
narration handshake; the spec's exact-slug create-or-attach
posture; the API endpoint / request / schema (the deterministic
path is a pure consumer). Adds **no** tool-origination or
spawn-UX content (m2); stays generic to construction-known slugs,
with m1's manual / CLI `--slug` argument as the stand-in producer.

**Sibling interface.** Produces the durable contract m2 consumes
when its spawn becomes the real construction-time slug producer;
sanctions the existing `workstream-tracker register --slug <slug>`
over the exact-slug create-or-attach flow — no new registration
surface (D1).

## Goal

Land an additive expression of the construction-known
deterministic registration path in two durable carriers
(`spec/planning/shared.md` and
`docs/agents/local/session-registration.md`) such that: the path
is sanctioned as a first-class sibling of the interactive
best-effort path; the interactive handshake, the exact-slug
create-or-attach posture, and the API/request/schema are provably
unchanged; the expression carries no tool-origination / spawn-UX
language; and the `rule-additions.md` discipline is satisfied for
the agent-rule edit. No code, no test, no new registration
surface.

## Cross-Cutting Invariants

These thread both edited files and the rule discipline; a
single-file self-review misses a violation that only shows when
two sites disagree.

- **Additive-only.** Neither edit changes the interactive
  handshake text, the exact-slug create-or-attach posture, or the
  API/request/schema. The interactive handshake section **and**
  the PR #45 session-end completion section of
  `session-registration.md` are both **byte-unchanged**; the
  `shared.md` exact-slug assertion posture text is unchanged. The
  new content is purely additional.
- **No new registration surface.** No endpoint, request/response
  field, schema, registration code path, or CLI flag — the
  expression describes the *existing* flow (D1).
- **Generic to construction-known slugs.** No
  tool-origination, spawn-UX, node-affordance, or Claude-Code
  language in either carrier; the only named producer is the
  manual / CLI `--slug` argument as a stand-in.
- **"Deterministic" stays epic-scoped.** The wording asserts
  identity-by-construction for a session that *did* launch — not
  "the launch/producer cannot be wrong," not enrichment made
  load-bearing for attachment.
- **Rule/spec-change discipline.** The `spec/**` edit follows the
  [`AGENTS.md`](../../../../AGENTS.md) `spec-authoring` pre-edit
  read; the `docs/agents/local/**` edit satisfies
  [`rule-additions.md`](../../../../docs/agents/shared/meta/rule-additions.md)
  in the PR/change description.

## Contracts

The conditions that must hold for the implementation to satisfy
this plan. Exact prose/headings are the implementer's choice
(contract altitude); only the conditions below bind.

### C1 — `spec/planning/shared.md`: additive first-class subsection

A new subsection is added to
[`shared.md`](../../../../spec/planning/shared.md) adjacent to the
"Slug generation" subsection (which carries the exact-slug
create-or-attach assertion posture). It must state that when a
session's canonical slug is known by construction — carried in,
not resolved from a natural-language prompt — assertion via the
existing exact-slug create-or-attach path is deterministic and no
natural-language-resolution step precedes it, and that this is an
additional path that does not replace the interactive one. It
must be generic to construction-known slugs (the manual / CLI
`--slug` argument named as the stand-in producer; no
tool-origination/spawn-UX language) and must not alter the
existing "Slug generation" assertion-posture text. *Verified by:*
the assertion posture it builds on,
[`spec/planning/shared.md`](../../../../spec/planning/shared.md)
"Slug generation" (`:126-186`); the mechanism it describes is a
pure consumer,
[`internal/api/handlers.go`](../../../../internal/api/handlers.go)
`insertRegister` exact-slug branch (`:197-201`).

### C2 — `docs/agents/local/session-registration.md`: additive sibling section, no confirm-equivalent

A new section is added to
[`session-registration.md`](../../../../docs/agents/local/session-registration.md)
stating that when the canonical slug is known by construction the
resolve/confirm interpretation steps do not apply, while the
invoke / echo-the-real-receipt / narrate-failure-explicitly /
proceed steps remain the same observable best-effort steps, and
that registration never gates the session. It must carry **no
confirm-equivalent** (no pre-invoke human-catch step): the
real-receipt echo is the sole observability surface, catching a
wrong-by-construction slug post-hoc — consistent with the
already-accepted, deferred orphan/unattached-work-instance
residual. **Placement (re-grounded against the PR #45 rebase):**
the file is now a *lifecycle* rule carrying, in order, Why
(`:15-23`), "## The handshake (interactive natural-language
session)" (`:25-69`), "## The session-end handshake (completion)"
(`:71-95`, added by #45), and "## Scope and residual" (`:97-114`).
The new section is placed **after the interactive handshake and
before the session-end completion section** — it is a sibling of
the *registration* handshake, not a completion concern. The
interactive handshake section **and** the #45 completion section
must both be **byte-unchanged** (additive-only). *Verified by:*
the rebased file structure the new section sits within,
[`docs/agents/local/session-registration.md`](../../../../docs/agents/local/session-registration.md)
("## The handshake (interactive natural-language session)",
`:25-69`; "## The session-end handshake (completion)", `:71-95`);
the binary already performs no resolve/confirm given a slug,
[`cmd/workstream-tracker/register.go`](../../../../cmd/workstream-tracker/register.go)
`runRegister` (`:39-93`); the server trusts the caller's slug and
is repo-blind (orphan residual),
[`internal/slugs/slugs.go`](../../../../internal/slugs/slugs.go)
`IsWellFormed` (`:73-117`).

### C3 — `rule-additions.md` discipline satisfied for the agent-rule edit

The change/PR description for the C2 edit states the
[`rule-additions.md`](../../../../docs/agents/shared/meta/rule-additions.md)
trade-off explicitly as answer **(b)**: the deterministic
construction-known path covers a class of registration no
existing local rule addressed, so no rule is retired or merged.
The now-interactive-only "Scope and residual" determinism clause
in `session-registration.md` is **left unchanged** (a documented,
accepted imprecision — keeping the edit strictly additive and not
touching interactive-adjacent prose the
`deterministic-interactive-registration` backlog tripwire
anchors). *Verified by:*
[`rule-additions.md`](../../../../docs/agents/shared/meta/rule-additions.md)
"The rule" (a)/(b);
[`AGENTS.md`](../../../../AGENTS.md) "Adding to this rule set"
(`:143`); the clause left unchanged (re-confirmed: #45 left it
verbatim),
[`docs/agents/local/session-registration.md`](../../../../docs/agents/local/session-registration.md)
"Scope and residual" (`:97-114`).

### C4 — `AGENTS.md` pointer coherence

[`AGENTS.md`](../../../../AGENTS.md) is **not** edited: under C2
the interactive handshake is byte-unchanged, so the "Universal
session rules" pointer summary stays accurate, and the router
pointer's "this file owns the detail" reaches the new section.
Re-grounded against the PR #45 rebase: #45 already rewrote that
pointer to "Session work-instance lifecycle handshake",
summarizing the full register **and** completion lifecycle (and
the `go run <module-path>` invocation). That makes it even more
plainly a non-exhaustive overview, not a sole-path assertion, so
the conclusion is unchanged. The implementer performs a one-line
coherence read of that summary against the final C1/C2 wording;
an additive coherence clause is added **only if** the final
wording makes the summary read as asserting the narration
handshake is the *sole* registration path (less likely
post-#45). *Verified by:* [`AGENTS.md`](../../../../AGENTS.md)
"Universal session rules" lifecycle-handshake pointer (`:116`,
rebased post-#45 summary); [`../README.md`](../README.md)
"Documentation Currency".

## Files to touch — new / modify / intentionally not touched

> *Estimate of expected shape. Implementation may revise a
> structural call; the Contracts above bind, this inventory does
> not.*

- **Modify:**
  [`spec/planning/shared.md`](../../../../spec/planning/shared.md)
  (C1 — additive subsection) and
  [`docs/agents/local/session-registration.md`](../../../../docs/agents/local/session-registration.md)
  (C2 — additive section).
- **Modify at closeout (status currency — not new surface; see
  Documentation Currency below):** this plan's own `Status`
  frontmatter (`Proposed` → … → `Landed`) and the
  [`../README.md`](../README.md) t1 Task Status row to its
  terminal state. The `In draft → Proposed` half was done at
  promotion; the promotion-gate parent-doc convention does **not**
  cover the terminal half, so the implementing PR owns it.
- **Intentionally not modified (expected):**
  [`AGENTS.md`](../../../../AGENTS.md) (C4 — only if the
  sole-path coherence check fires);
  [`design/v0.1-design.md`](../../../../design/v0.1-design.md)
  (a frozen v0.1 end-state record — not an authority or
  reconciliation target for this work; not touched);
  `internal/**`, `cmd/**` (no code — D1); the backlog (no
  backlog entry — D1 / milestone Backlog Impact).
- **New:** none. The paired scoping doc already exists, is
  transient, and is deleted at the **milestone-terminal** PR —
  not by t1's implementing PR.

## Validation Gate

Doc-only; the gate is verification, not a build.

1. **Additive-only diff check.** `git diff` shows the interactive
   handshake section **and** the PR #45 session-end completion
   section of `session-registration.md` byte-unchanged, and the
   `shared.md` "Slug generation" assertion-posture text unchanged;
   the only changes are net-additive blocks. (Falsifier: any
   deletion/modification inside the interactive section, the
   completion section, or the exact-slug posture text fails the
   gate.)
2. **No-new-surface check.** No new registration / code surface:
   no `internal/**`, `cmd/**`, schema, endpoint, or CLI-flag
   change (D1). The diff is the two carrier docs, plus
   conditionally `AGENTS.md` (C4), plus the status-closeout edits
   the Documentation Currency section names (this plan's own
   `Status`; the m1 `README.md` t1 Task Status row). This is a
   *surface* check, **not a literal file allowlist** — the
   status-closeout edits are expected and are not a violation.
   (Falsifier: any `internal/**` / `cmd/**` / schema / endpoint /
   CLI-flag change fails the gate; a `Status` or parent-row
   currency edit does not.)
3. **Generic-language check.** Neither added block contains
   tool-origination, spawn-UX, node-affordance, or Claude-Code
   terms; the only named producer is the manual / CLI `--slug`
   argument.
4. **`spec-authoring` pre-edit read** performed for the
   `shared.md` edit (`spec/README.md` + `spec/planning/shared.md`),
   and **`rule-additions.md` (b)-answer** stated in the change
   description for the `session-registration.md` edit.
5. **AGENTS.md coherence read** performed against the final
   C1/C2 wording; C4 clause added only if the sole-path reading
   fires.
6. **Reality-check re-confirm.** *Every* input named in the
   scoping doc's "Reality-check inputs" list — the full enumerated
   set, not only the subset surfaced in Contracts C1–C4 — still
   holds against then-merged code. This explicitly includes the
   `internal/registerclient/client.go` behavior and the
   `cmd/workstream-tracker/register_test.go`
   (`TestRegisterCommandSuccessAndIdempotentRepeat`) end-to-end
   determinism proof the Context preamble leans on, not just the
   carrier-contract citations. That full named list is the
   falsifier set; a stale entry anywhere in it fails the gate.

## Self-Review Audits

Diff surface maps to these seeded audits
([catalog](../../../../docs/agents/local/self-review-catalog.md)):

- **validation-honesty** — the Validation Gate above must be run
  end-to-end on the actual diff before any "checks passed" claim;
  a prose success claim without the diff inspection is the trap.
- **readiness-gate-truthfulness** — the `In draft → Proposed`
  promotion claim is only valid if the gate's named conditions
  were actually verified, not announced best-effort.
- **trigger-map-currency** — adding a section to
  `docs/agents/local/session-registration.md` does not restructure
  the tree, but confirm the `AGENTS.md` "Mandatory pre-edit reads"
  trigger map (`spec/**`, `docs/plans/**`) still matches reality
  after the edit.

## Documentation Currency

Status-bearing docs the implementing PR must keep current (per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Plan-to-PR Completion Gate" and
[`shared.md`](../../../../spec/planning/shared.md) "Plan-doc
Status"). These edits are expected closeout, not a
"no-new-surface" violation (Validation Gate step 2).

- **This plan's own `Status`.** The implementing PR transitions
  it `Proposed → In progress` on start and `→ Landed` at close —
  no `Validating` (doc-only; no post-merge gate). Frontmatter is
  authoritative.
- **The m1 `README.md` t1 Task Status row.** Currently
  `Proposed (task plan drafted; scoping SD1–SD7 resolved)`
  ([`../README.md`](../README.md) "Task Status"). The implementing
  PR's terminal commit moves that row to the matching terminal
  state and narrows the adjacent "t1's task planning is complete"
  prose if needed. This is the **task-terminal** parent-row
  reconciliation: the `In draft → Proposed` half landed at
  promotion, and the promotion-gate parent-doc convention has no
  symmetric step at `Landed`, so ownership falls here. Without
  this, the parent row drifts stale after t1 lands.
- **`AGENTS.md`** — only if C4's sole-path coherence check fires
  (then it is a carrier edit, not status currency).
- **Not t1's to delete: the paired scoping doc.**
  [`scoping/t1-deterministic-path-contract.md`](./scoping/t1-deterministic-path-contract.md)
  is transient but is deleted in the **milestone-terminal** batch
  (m1's close, after t2) per
  [`planning-doc-location.md`](../../../../spec/planning-doc-location.md)
  "scoping/ subfolder is transient" and
  [`task-plan.md`](../../../../spec/planning/task-plan.md)
  "Scoping owns / plan owns" — t1's task-terminal PR is not the
  milestone-terminal PR, so it must **not** delete scoping.

## Out Of Scope

- **The m2 tool-originated-session UX** — node affordance, real
  construction-time slug producer, session-start hook, spawn
  shape. m1's producer is the manual/CLI stand-in.
- **Any change to the interactive registration path** or the
  independent `interactive-registration-tripwire`. t1 adds a
  sibling; it does not touch, resolve, or pull forward the
  interactive posture.
- **New registration code surface** — endpoint, request field,
  schema, registration code path, CLI flag (D1).
- **Editing `design/v0.1-design.md`** — frozen v0.1 end-state
  record; not load-bearing for this work and not reconciled
  against.
- **The tree-side observability heuristic**
  (`unregistered-work-unobservable`).

## Risk Register

- **Wording drifts into tool-origination / spawn-UX language,
  pulling m2 vision forward.** Mitigation: the
  generic-to-construction-known invariant and Validation Gate
  step 3; the only named producer is the manual/CLI stand-in.
- **The additive edit inadvertently weakens the interactive
  handshake.** Mitigation: the byte-unchanged invariant on the
  interactive section and Validation Gate step 1; C2's
  no-confirm-equivalent change is isolated to the new section.
- **Accepted residual: the "Scope and residual" determinism
  clause stays interactive-only-true but generally worded.**
  Consciously not fixed (C3) to keep the edit strictly additive
  and clear of the backlog-tripwire anchor; recorded here so a
  reviewer flagging it gets the hold-the-line rationale rather
  than a re-litigation.
- **Accepted residual: a wrong-by-construction slug attaches an
  orphan work-instance** (no pre-invoke confirm, C2). This is the
  already-accepted, deferred orphan/unattached residual, out of
  t1's charter; the real-receipt echo keeps it observable
  post-hoc.

## Backlog Impact

None. Per [`../README.md`](../README.md) "Backlog Impact", m1's
implementing PRs touch no backlog entry; the
graduate/split/carve-outs landed in the epic promotion PR.

## Related Docs

- [`./scoping/t1-deterministic-path-contract.md`](./scoping/t1-deterministic-path-contract.md)
  — the paired transient scoping doc (SD1–SD7 deliberation +
  rejected alternatives; deletes at the milestone-terminal PR).
- [`../README.md`](../README.md) — m1 milestone doc; locks t1's
  WHAT, Cross-Task Invariants/Decisions, Documentation Currency.
- [`../../README.md`](../../README.md) — parent epic;
  Cross-Cutting Invariants and the epic-scoped meaning of
  "deterministic".
- [`../../../../spec/planning/shared.md`](../../../../spec/planning/shared.md)
  — "Slug generation" (the exact-slug assertion posture t1
  expresses additively).
- [`../../../../docs/agents/local/session-registration.md`](../../../../docs/agents/local/session-registration.md),
  [`../../../../AGENTS.md`](../../../../AGENTS.md),
  [`../../../../docs/agents/shared/meta/rule-additions.md`](../../../../docs/agents/shared/meta/rule-additions.md)
  — the agent-rule carriers and the rule-additions discipline.
- [`../../../../internal/api/handlers.go`](../../../../internal/api/handlers.go),
  [`../../../../internal/slugs/slugs.go`](../../../../internal/slugs/slugs.go),
  [`../../../../cmd/workstream-tracker/register.go`](../../../../cmd/workstream-tracker/register.go)
  — the merged code grounding C1–C3.
