---
slug: tool-originated-task-sessions-m1-t1
short_description: Scoping — deterministic-path contract expression
---

# Scoping — m1 t1 Deterministic-Path Contract Expression

> **Spawned scoping session — decisions decomposed, then resolved
> by the reviewer in-session.** This is the *scoping /
> investigation* artifact for `tool-originated-task-sessions-m1-t1`,
> not the plan doc and not a plan-drafting session. It runs no
> promotion gate and opens no PR. Each HOW decision was decomposed
> into shapes against merged code and left open; the reviewer then
> resolved them during the scoping walk-through. They are recorded
> below as scoping decisions for the plan-drafting session to
> consume, not re-litigate. The durable plan doc remains the
> skeleton
> [`../t1-deterministic-path-contract.md`](../t1-deterministic-path-contract.md).
> Frontmatter note: per
> [`task-plan.md`](../../../../../spec/planning/task-plan.md)
> "Goal: scoping doc + plan doc", scoping docs do **not** carry a
> `Status` field (a scoping-doc Status can only become wrong), so
> this doc carries none (see SD5).

## Context preamble

t1 must express the deterministic, handshake-free,
slug-by-construction registration path as a **first-class additive
contract** in two durable carriers: the plan-doc spec
([`spec/planning/shared.md`](../../../../../spec/planning/shared.md))
and the session-registration agent rule
([`docs/agents/local/session-registration.md`](../../../../../docs/agents/local/session-registration.md)).
The contract: when a session's canonical slug is known *by
construction* (carried in, not resolved from a natural-language
prompt), registration is deterministic and the best-effort
grounded narration handshake does not apply. The expression is
**purely additive** — the interactive handshake for
natural-language sessions, the spec's exact-slug create-or-attach
posture, and the API/request/schema are all unchanged; the
deterministic path is a pure consumer of the already-merged
exact-slug flow. It carries **no** tool-origination / spawn-UX
content (that is the m2 milestone); it stays generic to
construction-known slugs, with the manual / CLI `--slug` argument
as the stand-in producer. The locked WHAT is inherited from the
milestone and is binding input, not to be loosened (see
[`../t1-deterministic-path-contract.md`](../t1-deterministic-path-contract.md)
"Inherited Contract" and [`../README.md`](../README.md) Cross-Task
Decisions D1/D2/D3, Cross-Task Invariants).

## Reality-check inputs (the plan must re-verify these at plan-drafting)

Code-grounded findings the t1 plan's contract claims rest on. Each
carries a `Verified by:` citation; the plan-drafting session
re-confirms them against then-merged code per
[`task-plan.md`](../../../../../spec/planning/task-plan.md)
"Reality-check pass before plan-drafting".

- **The deterministic registration mechanism already exists
  end-to-end in merged code; t1 adds no registration code
  surface.** The server's exact-slug branch validates slug grammar
  only and then honors the caller's slug verbatim — no
  natural-language resolution, no descendant generation, no
  root-conflict check. *Verified by:*
  [`internal/api/handlers.go`](../../../../../internal/api/handlers.go)
  `registerWorkInstance` exact-slug validation (`req.ExactSlug != ""`
  → `slugs.IsWellFormed` only, `:58-62`) and `insertRegister`
  exact-slug branch (`slug = req.ExactSlug`, comment "honor the
  caller's slug verbatim. No derivation, no descendant generation,
  no root-conflict check", `:197-201`).
- **Slug well-formedness is grammar-only and repo-blind.**
  *Verified by:* [`internal/slugs/slugs.go`](../../../../../internal/slugs/slugs.go)
  `IsWellFormed` (package/function doc: "It does not verify the
  slug names a real plan-tree doc — the server is repo-blind",
  `:73-117`).
- **The single-request client posts exactly that, with no retry,
  no derivation, no plan-tree reading.** *Verified by:*
  [`internal/registerclient/client.go`](../../../../../internal/registerclient/client.go)
  package doc and `Register` (`:1-7`, `:42-93`).
- **The CLI is already deterministic and handshake-free given a
  slug: it carries `--slug`/`WST_SLUG`, `--actor`/`WST_ACTOR`,
  `--server`/`WST_SERVER` and an optional `--name`/`WST_NAME`,
  performs no natural-language resolution, always exits 0, prints
  the real observed receipt on success, and narrates failure
  explicitly.** *Verified by:*
  [`cmd/workstream-tracker/register.go`](../../../../../cmd/workstream-tracker/register.go)
  `runRegister` (`:38-96`).
- **An end-to-end test already proves determinism: a verbatim slug
  supplied as an explicit argument registers against the real
  exact-slug create-or-attach flow and an idempotent repeat
  collapses to the same work-instance.** *Verified by:*
  [`cmd/workstream-tracker/register_test.go`](../../../../../cmd/workstream-tracker/register_test.go)
  `TestRegisterCommandSuccessAndIdempotentRepeat` (verbatim
  `demo-root-m1-t2`, asserts `http_status=201` and same
  `work_instance_id` on repeat, `:46-81`).
- **The "narration handshake" is not in the binary** — it is a
  layered agent procedure described only in the interactive
  section of
  [`docs/agents/local/session-registration.md`](../../../../../docs/agents/local/session-registration.md)
  ("## The handshake (interactive natural-language session)",
  `:18-55`); `runRegister` contains no resolve/confirm step.
  *Verified by:* the two citations above
  (`register.go:38-96`; `session-registration.md:18-55`).
- **The exact-slug flow adds no endpoint, request/response field,
  or schema change, and PR #38's `--name` enrichment is the
  orthogonal best-effort *enrichment* leg, not the deterministic
  identity step.** The exact-slug branch posts to the existing
  `/work-instances` endpoint and writes only the existing `events`
  / `work_instances` columns (no new column for the exact-slug or
  metadata path); the deterministic identity step (exact-slug +
  actor) is unchanged by #38, which rides the pre-existing optional
  `metadata` param. *Verified by:*
  [`internal/api/handlers.go`](../../../../../internal/api/handlers.go)
  `insertRegister` INSERT statements into existing `events` /
  `work_instances` tables (`:247-266`);
  [`internal/registerclient/client.go`](../../../../../internal/registerclient/client.go)
  `Register` (posts `exact_slug` / `actor` / optional `metadata`
  to the existing endpoint, `:42-93`);
  [`cmd/workstream-tracker/register.go`](../../../../../cmd/workstream-tracker/register.go)
  `runRegister` name→metadata handling (`:68-78`). (Grounded on
  merged code, not `design/v0.1-design.md` — see SD6.)

## Settled by inherited contract / merged code (not open)

Locked upstream or proven in code; recorded so the scoping
decisions below are legible by contrast. Not reopenable by the
plan-drafting session.

- **D1 — reuse the existing exact-slug create-or-attach path; no
  new registration surface.** Resolved by the milestone session
  ([`../README.md`](../README.md) Cross-Task Decision D1) and
  re-verified above. t1 writes no endpoint, request/response
  field, schema change, registration code path, or CLI flag.
- **D2 — the spec + agent-rule expression is m1's, taken up as
  t1.** Resolved by the milestone session
  ([`../README.md`](../README.md) Cross-Task Decision D2). The
  expression stays generic to construction-known slugs and must
  not encode the m2 node-affordance producer.
- **The expression is additive only.** The interactive handshake,
  the exact-slug create-or-attach posture, and the
  API/request/schema remain unchanged
  ([`../README.md`](../README.md) Cross-Task Invariants "Additive
  only", "No new registration surface").
- **"Deterministic" stays epic-scoped.** Identity-by-construction
  for a session that *did* launch — not "the launch cannot fail",
  not enrichment made load-bearing for attachment
  ([`../README.md`](../README.md) Cross-Task Invariants
  "'Deterministic' stays epic-scoped"; epic
  [`../../README.md`](../../README.md) Cross-Cutting Invariants
  "What 'deterministic' is scoped to"). A binding **wording
  guard** on every decision below.

## Decisions made at scoping time

Each decision was decomposed into shapes per
[`shared.md`](../../../../../spec/planning/shared.md) "Decompose
options into shapes before analyzing", then resolved by the
reviewer in-session. Constraints binding every decision:
additive-only; no new registration surface; no tool-origination /
spawn-UX language; generic to construction-known slugs with the
CLI `--slug` as stand-in producer; the epic-scoped meaning of
"deterministic"; the [`AGENTS.md`](../../../../../AGENTS.md)
`spec-authoring` pre-edit read for any `spec/**` edit and the
[`rule-additions.md`](../../../../../docs/agents/shared/meta/rule-additions.md)
discipline for any `docs/agents/local/**` edit.

### SD1 — `shared.md` placement: a new first-class sibling subsection

**Decision.** Express the deterministic path as a **new sibling
subsection** in
[`shared.md`](../../../../../spec/planning/shared.md), placed
adjacent to "Slug generation" (the exact-slug create-or-attach
assertion posture), naming construction-known deterministic
registration as a first-class concept. The exact prose is
implementation, settled in the implementing PR — scoping fixes the
*shape* (new subsection), not the wording.

**Rejected.** *A1 — an additive paragraph inside the existing
"Slug generation" subsection:* smallest surface, but "Slug
generation" is scoped to how slugs come to exist and are
validated, not to registration-handshake applicability; folding
the concept there overloads the subsection and weakens the
"first-class" requirement. *A3 — a pointer-only cross-reference in
`shared.md` with the substance in the agent rule:* smallest spec
surface, but the locked WHAT names `shared.md` as a *carrier* of
the first-class contract; pointer-only under-delivers and risks
loosening a locked input.

**Verified by:**
[`shared.md`](../../../../../spec/planning/shared.md) "Slug
generation" (`:126-186`);
[`internal/api/handlers.go`](../../../../../internal/api/handlers.go)
`insertRegister` exact-slug branch (the described path is a pure
consumer of the existing flow, `:197-201`).

### SD2 — `session-registration.md` placement: a new sibling section, interactive untouched

**Decision.** Add a **new sibling section** after the interactive
handshake in
[`session-registration.md`](../../../../../docs/agents/local/session-registration.md)
(e.g. a "deterministic path (construction-known slug)" section):
because identity is given by construction, the resolve step does
not apply; invoke / echo-the-real-receipt / narrate-failure /
proceed remain the same observable best-effort steps. The
interactive handshake section is left **textually unchanged**.

**Rejected.** *B2 — refactor to a shared invoke/echo/narrate/
proceed core with thin "interactive" and "deterministic"
subsections:* removes duplication but restructures the interactive
section; even behavior-preserving, this reads as a *change* to a
path the additive-only invariant says stays unchanged and brushes
the milestone's "agent-rule edit inadvertently weakens the
interactive handshake" risk. *B3 — a minimal note under "Why this
rule exists":* under-expresses the locked first-class sibling
section.

**Verified by:**
[`session-registration.md`](../../../../../docs/agents/local/session-registration.md)
current structure (`:8-65`); [`../README.md`](../README.md)
Cross-Task Invariants "Additive only" and Cross-Task Risks
"Agent-rule edit inadvertently weakens the interactive handshake".

### SD3 — no confirm-equivalent in the deterministic section

**Decision.** The deterministic section carries **no Confirm
step**. Because identity is given by construction there is nothing
to interpret and no pre-invoke human gate; the real-receipt echo
(the same step the interactive path already performs) is the
**sole observability surface**, catching a wrong-by-construction
slug *post-hoc*, not pre-invoke.

**Rejected.** *S2 — a lightweight confirm-equivalent (echo the
slug at invoke, no pause):* partially duplicates the receipt echo
and reads as "still a handshake", muddying the deterministic /
handshake-free contract. *S3 — bare-drop plus an explicit
rationale sentence:* acceptable, but the extra sentence was judged
unnecessary; the omission is self-explanatory next to the
interactive section.

**Rationale / accepted residual.** A wrong-slug-*by-construction*
(a typo'd `--slug` now; a mis-targeted producer in m2's future) is
the **already-accepted, deferred orphan / unattached-work-instance
residual** — the server trusts the caller's slug and is repo-blind
— and is out of t1's charter. Dropping pre-invoke confirm is
consistent with the epic-scoped meaning of "deterministic"
(identity for a session that *did* launch, not "the producer
cannot be wrong").

**Verified by:**
[`internal/api/handlers.go`](../../../../../internal/api/handlers.go)
`insertRegister` exact-slug branch trusts the caller's slug
(`:197-201`);
[`internal/slugs/slugs.go`](../../../../../internal/slugs/slugs.go)
`IsWellFormed` repo-blind (`:73-117`);
[`shared.md`](../../../../../spec/planning/shared.md) "Slug
generation" orphan/trust-the-caller residual note (`:160-179`).

### SD4 — `rule-additions.md`: answer (b), no retire (new failure class)

**Decision.** Satisfy the
[`rule-additions.md`](../../../../../docs/agents/shared/meta/rule-additions.md)
discipline with answer **(b)**: the deterministic
construction-known path covers a class of registration no existing
local rule addressed, so **no rule is retired or merged**. The
plan-drafting / implementing PR states this trade-off explicitly
in its change description.

**Rejected.** *C1 — retire/narrow the "Scope and residual"
determinism clause:* that clause is interactive-path scoped and is
referenced by the `deterministic-interactive-registration` backlog
tripwire; editing it brushes the additive-only invariant and the
interactive-handshake-weakening risk. *C3 — hybrid (b) plus a
coherence-only narrowing of that clause:* same interactive-adjacent
edit risk.

**Accepted residual.** The "Scope and residual" determinism
sentence in `session-registration.md` remains slightly
over-claiming (it states a now-interactive-only accepted gap in
general terms). Accepted as a documented wart, deliberately **not**
fixed, to keep the edit strictly additive and avoid touching
interactive-adjacent prose the backlog tripwire anchors on.

**Verified by:**
[`rule-additions.md`](../../../../../docs/agents/shared/meta/rule-additions.md)
"The rule" (a)/(b); [`AGENTS.md`](../../../../../AGENTS.md)
"Adding to this rule set" (`:137-144`);
[`session-registration.md`](../../../../../docs/agents/local/session-registration.md)
"Scope and residual" (`:57-65`).

### SD5 — scoping docs carry no `Status` field

**Decision.** This scoping doc carries **no `Status` frontmatter
field**, conforming to the spec over the spawned-session
instruction's `Status: In draft` (which was struck).

**Verified by:**
[`task-plan.md`](../../../../../spec/planning/task-plan.md) "Goal:
scoping doc + plan doc" (scoping docs do not carry a `Status`
field; a scoping-doc Status can only become wrong).

### SD6 — `design/v0.1-design.md` is a frozen v0.1 record, not load-bearing

**Decision.** `design/v0.1-design.md` describes only the v0.1 end
state and is **not an authority, currency constraint, or
reconciliation target** for this future-feature work. t1 does not
cite §4/§5 as load-bearing, does not reconcile its wording against
them, and does not edit that doc. All "pure consumer / no new
surface" facts are grounded directly in merged code (see
Reality-check inputs). The milestone's "currency-check §4, no edit
expected" framing is satisfied trivially because the doc is
historical.

**Verified by:** the merged-code Reality-check inputs above;
no `design/v0.1-design.md` citation by design.

### SD7 — no `AGENTS.md` pointer edit

**Decision.** **No edit to `AGENTS.md`.** Under SD2 the interactive
handshake is textually unchanged, so the `AGENTS.md` "Universal
session rules" summary (which describes the interactive handshake)
stays accurate; a router pointer may be non-exhaustive and the
deterministic section is reachable via the existing "this file
owns the detail" pointer.

**Rejected.** *F2 — a minimal additive coherence clause in the
`AGENTS.md` summary:* re-trips the `rule-additions.md` discipline
for an `AGENTS.md` change and widens the additive-only blast
radius; warranted only if the final SD1/SD2 wording reads as
asserting the handshake is the *sole* registration path — a
one-line coherence read the plan-drafting session performs against
the concrete wording, not expected to fire under SD2.

**Verified by:** [`AGENTS.md`](../../../../../AGENTS.md)
"Universal session rules" session-registration pointer
(`:116-121`); [`../README.md`](../README.md) "Documentation
Currency" ("t1 must also keep `AGENTS.md`'s session-registration
pointer coherent if the agent-rule's shape changes").

## Open decisions to make at plan-drafting (handoff)

**None remain open.** Every decomposed HOW decision (SD1–SD7) was
resolved in this scoping session. What remains for the
plan-drafting / implementing PR is **implementation-altitude**,
not an open scoping decision:

- The literal prose of the new `shared.md` subsection (SD1) and
  the new `session-registration.md` section (SD2). Per
  [`shared.md`](../../../../../spec/planning/shared.md) "Plans
  describe contracts, not implementation", the plan states the
  additive contract the two carriers must express; the exact
  wording is settled in the implementing PR alongside the
  `spec-authoring` and `rule-additions` pre-edit reads.
- A one-line coherence read of the `AGENTS.md` summary against the
  final SD1/SD2 wording (the SD7 contingency check) — a validation
  step, not a decision.

## Plan-structure handoff

For the plan-drafting session:

- t1 is a **doc-only contract-expression task** (no code ships;
  D1's no-new-surface is locked). The
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  doc-only-decision-plan framing applies.
- **Files the plan will contract:**
  [`spec/planning/shared.md`](../../../../../spec/planning/shared.md)
  (new sibling subsection, SD1) and
  [`docs/agents/local/session-registration.md`](../../../../../docs/agents/local/session-registration.md)
  (new sibling section, SD2). `AGENTS.md` is **not** expected to
  change (SD7). `design/v0.1-design.md` is **not** touched (SD6).
- **Validation gate the plan will carry:** the `spec-authoring`
  pre-edit read for the `spec/**` edit; the `rule-additions.md`
  (b)-answer trade-off articulation for the `docs/agents/local/**`
  edit (SD4); the SD7 coherence read; confirmation the interactive
  handshake section is byte-unchanged (additive-only).
- **PR-count:** a single doc PR is the estimate (narrow doc
  surface), subject to the
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "PR-count predictions need a branch test" rule — an estimate,
  not a commitment.

## Related docs

- [`../t1-deterministic-path-contract.md`](../t1-deterministic-path-contract.md)
  — the durable t1 skeleton carrying the locked inherited
  contract (this scoping doc pairs with it).
- [`../README.md`](../README.md) — the m1 milestone doc; locks
  t1's WHAT and carries Cross-Task Decisions D1/D2/D3, Invariants,
  Risks, and Documentation Currency.
- [`../../README.md`](../../README.md) — the parent epic;
  Cross-Cutting Invariants and the epic-scoped meaning of
  "deterministic".
- [`../../../../../spec/planning/shared.md`](../../../../../spec/planning/shared.md)
  — "Slug generation" (the exact-slug assertion posture t1
  expresses additively) and the scoping-doc / decompose-options
  rules this doc follows.
- [`../../../../../spec/planning/task-plan.md`](../../../../../spec/planning/task-plan.md)
  — scoping-vs-plan ownership, reality-check pass, the
  doc-only-decision-plan carve-out this scoping doc relies on.
- [`../../../../../docs/agents/local/session-registration.md`](../../../../../docs/agents/local/session-registration.md),
  [`../../../../../AGENTS.md`](../../../../../AGENTS.md),
  [`../../../../../docs/agents/shared/meta/rule-additions.md`](../../../../../docs/agents/shared/meta/rule-additions.md)
  — the agent-rule carriers and the rule-additions discipline.
- [`../../../../../internal/api/handlers.go`](../../../../../internal/api/handlers.go),
  [`../../../../../internal/slugs/slugs.go`](../../../../../internal/slugs/slugs.go),
  [`../../../../../internal/registerclient/client.go`](../../../../../internal/registerclient/client.go),
  [`../../../../../cmd/workstream-tracker/register.go`](../../../../../cmd/workstream-tracker/register.go),
  [`../../../../../cmd/workstream-tracker/register_test.go`](../../../../../cmd/workstream-tracker/register_test.go)
  — the merged code grounding the reality-check inputs.
  (`design/v0.1-design.md` is intentionally **not** listed — SD6.)
