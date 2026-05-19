---
slug: tool-originated-task-sessions-m1-t2
short_description: Scoping — determinism proof harness
---

# Scoping — m1 t2 Determinism Proof Harness

> **Spawned scoping session — decisions decomposed, then resolved
> in-session.** This is the *scoping / investigation* artifact for
> `tool-originated-task-sessions-m1-t2`, not the plan doc and not a
> plan-drafting session. It runs no promotion gate and opens no PR.
> Each HOW decision was decomposed into shapes against merged code
> per [`shared.md`](../../../../../spec/planning/shared.md)
> "Decompose options into shapes before analyzing", then resolved
> in-session. They are recorded below as scoping decisions for the
> plan-drafting session to consume, not re-litigate. The durable
> plan doc remains the skeleton
> [`../t2-determinism-proof-harness.md`](../t2-determinism-proof-harness.md).
> Frontmatter note: per
> [`task-plan.md`](../../../../../spec/planning/task-plan.md)
> "Goal: scoping doc + plan doc", scoping docs do **not** carry a
> `Status` field (a scoping-doc Status can only become wrong), so
> this doc carries none (SD6; mirrors t1's scoping SD5).

## Context preamble

t2 produces a **determinism-specific end-to-end test** proving the
contract t1 expressed. Given a construction-known slug supplied as
an explicit argument with **no narration handshake**, a
work-instance must register deterministically against the real
exact-slug create-or-attach path: identity-by-construction (slug
honored verbatim), idempotent create-or-attach on
`(slug, actor, active)`, no natural-language-resolution path
engaged, server repo-blind. The slug producer is the **minimal
explicit `--slug` argument / test harness** — the stand-in for
m2's future construction-time producer, scoped as the minimum to
validate the contract, **not a durable UX**. t2 adds **no product
code surface** and leaves the existing CLI tests unchanged; the
existing
[`register_test.go`](../../../../../cmd/workstream-tracker/register_test.go)
`TestRegisterCommandSuccessAndIdempotentRepeat` httptest-API
pattern is the named template. The locked WHAT is inherited from
the milestone and is binding input, not to be loosened (see
[`../t2-determinism-proof-harness.md`](../t2-determinism-proof-harness.md)
"Inherited Contract" and [`../README.md`](../README.md) Task
Contracts, Cross-Task Invariants, Cross-Task Decision D1).

t2's assertions must **mirror the determinism properties t1's
now-landed contract names** — not invent properties the sanctioned
contract does not express (milestone Sequencing: "writing the
proof first risks asserting properties the sanctioned contract
does not actually express"). The contract carriers are merged on
`main`: the
[`spec/planning/shared.md`](../../../../../spec/planning/shared.md)
"Deterministic assertion when the slug is construction-known"
subsection (`:181-200`) and the
[`session-registration.md`](../../../../../docs/agents/local/session-registration.md)
"## The deterministic path (construction-known slug)" section
(`:71-123`). The enumerated determinism facts t1's scoping doc
verified ([`./t1-deterministic-path-contract.md`](./t1-deterministic-path-contract.md)
"Reality-check inputs") are the assertion target set.

## Reality-check inputs (the plan must re-verify these at plan-drafting)

Code-grounded findings t2's proof rests on. Each carries a
`Verified by:` citation; the plan-drafting session re-confirms
them against then-merged code per
[`task-plan.md`](../../../../../spec/planning/task-plan.md)
"Reality-check pass before plan-drafting". Grounded on merged code
on `main` at commit `9e3893e` (t1 PR #47 landed); not on
`design/v0.1-design.md` (SD7).

- **The exact-slug branch honors the caller's slug verbatim — no
  derivation, no descendant generation, no root-conflict check.**
  *Verified by:*
  [`internal/api/handlers.go`](../../../../../internal/api/handlers.go)
  `insertRegister` exact-slug branch (`slug = req.ExactSlug`,
  comment "honor the caller's slug verbatim. No derivation, no
  descendant generation, no root-conflict check", `:197-201`).
- **The only acceptance gate on an exact slug is grammar, and it
  is repo-blind.** The handler validates `slugs.IsWellFormed`
  only (400 on malformed) before the verbatim branch; well-formed
  is grammar-only and does not verify the slug names a real
  plan-tree doc. *Verified by:*
  [`internal/api/handlers.go`](../../../../../internal/api/handlers.go)
  `registerWorkInstance` exact-slug validation (`req.ExactSlug != ""`
  → `slugs.IsWellFormed` else 400, `:58-62`);
  [`internal/slugs/slugs.go`](../../../../../internal/slugs/slugs.go)
  `IsWellFormed` ("It does not verify the slug names a real
  plan-tree doc — the server is repo-blind", `:73-117`).
- **Idempotency is keyed on `(slug, actor, state = active)`; a
  different actor on the same slug is independent (co-working);
  the bare-root-create branch 409s on a repeat (a discriminator
  the exact-slug attach does not share).** *Verified by:*
  [`internal/api/handlers.go`](../../../../../internal/api/handlers.go)
  `insertRegister` active-pair idempotency (`activeWorkInstanceID`
  → return existing, `:226-234`) and the bare-root `errSlugConflict`
  branch (`:216-224`); INSERTs target only the existing `events` /
  `work_instances` tables, no new column (`:247-266`).
- **The single-request client posts `exact_slug` with no retry,
  no derivation, no plan-tree reading; `--name` enrichment rides
  the optional `metadata` blob and is not load-bearing for
  attachment.** *Verified by:*
  [`internal/registerclient/client.go`](../../../../../internal/registerclient/client.go)
  package doc ("no retry, no derivation, no plan-tree reading",
  `:1-7`) and `Register` (posts `exact_slug` / `actor` / optional
  `metadata`, single request, `:42-93`).
- **The CLI is deterministic and handshake-free given a slug:
  `runRegister` performs no resolve/confirm step, always exits 0,
  prints the real receipt (including `slug=`) on success, narrates
  failure explicitly.** *Verified by:*
  [`cmd/workstream-tracker/register.go`](../../../../../cmd/workstream-tracker/register.go)
  `runRegister` (`:39-93`); the "narration handshake" exists only
  as the interactive agent procedure in
  [`session-registration.md`](../../../../../docs/agents/local/session-registration.md)
  "## The handshake (interactive natural-language session)"
  (`:25-69`), not in the binary.
- **The named template already drives the CLI end-to-end against
  a real API server through the exact-slug flow, with a
  position-segmented slug round-tripping verbatim and an
  idempotent same-actor repeat collapsing to the same
  work-instance.** *Verified by:*
  [`cmd/workstream-tracker/register_test.go`](../../../../../cmd/workstream-tracker/register_test.go)
  `TestRegisterCommandSuccessAndIdempotentRepeat` (slug
  `demo-root-m1-t2`, asserts `http_status=201` and same
  `work_instance_id` on repeat, `:44-80`).
- **The named template's package-`main` test helpers are reusable
  from a sibling file in the same package without exporting
  anything new.** *Verified by:*
  [`cmd/workstream-tracker/register_test.go`](../../../../../cmd/workstream-tracker/register_test.go)
  `startAPIServer` (`:23-38`), `noEnv` (`:40`), `widPattern`
  (`:42`), `isolateCaches` (`:140-145`) — all unexported in
  `package main`, reachable from any `*_test.go` in
  `cmd/workstream-tracker/`.
- **The deterministic-path contract t2 mirrors is landed and
  names exactly these properties.** *Verified by:*
  [`spec/planning/shared.md`](../../../../../spec/planning/shared.md)
  "Deterministic assertion when the slug is construction-known"
  (`:181-200`);
  [`docs/agents/local/session-registration.md`](../../../../../docs/agents/local/session-registration.md)
  "## The deterministic path (construction-known slug)"
  (`:71-123`).

## Settled by inherited contract / merged code (not open)

Locked upstream or proven in code; recorded so the scoping
decisions below are legible by contrast. Not reopenable by the
plan-drafting session.

- **D1 — no new registration surface.** No endpoint,
  request/response field, schema, registration code path, or CLI
  flag ([`../README.md`](../README.md) Cross-Task Decision D1).
  t2 adds **test files only**; the proof exercises already-merged
  code unchanged.
- **No product code surface; existing CLI tests unchanged.**
  Inherited Preserves clause
  ([`../t2-determinism-proof-harness.md`](../t2-determinism-proof-harness.md)).
- **"Deterministic" stays epic-scoped.** The proof asserts only
  the deterministic identity step for a session that *did* launch
  — not "the launch cannot fail," not enrichment as load-bearing
  for attachment ([`../README.md`](../README.md) Cross-Task
  Invariants; epic Cross-Cutting Invariants). A binding **scope
  guard** on every decision below, decisive for SD3-P2.
- **The contract t2 mirrors is t1's landed wording.** t2 demonstrates
  the properties t1's two carriers name; it does not extend or
  reinterpret them (milestone Sequencing; t2 Sibling interface).

## Decisions made at scoping time

Each decision was decomposed into shapes per
[`shared.md`](../../../../../spec/planning/shared.md) "Decompose
options into shapes before analyzing" (`:979-1009`), then
resolved. Constraints binding every decision: D1 no-new-surface;
test-files-only; existing CLI tests byte-unchanged; epic-scoped
"deterministic"; assertions mirror t1's named properties; the
falsifiability check ([`shared.md`](../../../../../spec/planning/shared.md)
`:944-977`) — a status-code-only assertion multiple causes
satisfy is the trap.

### SD1 — Test artifact placement: a new sibling test file in the same `package main`

**Decision.** The determinism proof is a **new sibling test file**
under `cmd/workstream-tracker/` (e.g.
`register_determinism_test.go`) in the existing `package main`,
reusing the named template's unexported helpers (`startAPIServer`,
`isolateCaches`, `widPattern`, `noEnv`, `runRegister`) verbatim.
[`register_test.go`](../../../../../cmd/workstream-tracker/register_test.go)
is left **byte-unchanged**. The exact filename is implementation,
fixed by the plan/PR; scoping fixes the *shape* (new file, same
package).

**Rejected.** *A1 — add the determinism assertions as new
functions inside `register_test.go`:* smallest file count, and
the helpers are in-scope, but the milestone frames t2's
deliverable as a distinct **proof artifact** ("adds the
determinism-specific assertions … rather than re-testing
idempotency from scratch") — a standalone file is more legible
as that artifact and keeps the named template trivially unchanged.
*A3 — a new external test package / integration directory:*
cannot reach the unexported `package main` helpers without
exporting them, which manufactures new test surface for no
benefit and adds structure the minimal proof does not need.

**Verified by:** the reusable helpers,
[`cmd/workstream-tracker/register_test.go`](../../../../../cmd/workstream-tracker/register_test.go)
(`:23-38`, `:40`, `:42`, `:140-145`); the milestone "proof
artifact" framing, [`../README.md`](../README.md) Task Contracts
t2 row.

### SD2 — Producer shape: the bare explicit `--slug` argument

**Decision.** The slug producer is the **bare explicit `--slug`
argument** driven through `runRegister` exactly as the named
template does
(`runRegister([]string{"--slug", <slug>, "--actor", …, "--server", ts.URL}, noEnv, …)`),
with the construction-known slug a **literal string constant** in
the test (the stand-in for m2's future producer). A thin
*test-local* receipt-parsing helper (extracting
`work_instance_id` / `http_status` / echoed `slug=` from stdout)
is permitted **only** to DRY the assertions across cases; it is
never a "producer" abstraction.

**Rejected.** *B2 — a producer helper/abstraction wrapping
invoke+parse as a named producer:* reads as durable scaffolding,
risks the milestone's "minimum to validate, not a durable UX"
bound, and a wrapper that swallows the raw receipt can mask the
very signals SD3 asserts on. *B3 — a non-test producer in
product/cmd code:* directly locked out by D1 (no new code
surface) and the inherited "not a durable UX" clause.

**Verified by:** the named-template invocation shape,
[`cmd/workstream-tracker/register_test.go`](../../../../../cmd/workstream-tracker/register_test.go)
`TestRegisterCommandSuccessAndIdempotentRepeat` (`:48-55`);
inherited "minimal explicit `--slug` argument / test harness …
not a durable UX"
([`../t2-determinism-proof-harness.md`](../t2-determinism-proof-harness.md)).

### SD3 — Property → assertion map (with falsifiability discriminators)

**Decision.** Each determinism property t1's contract names maps
to a concrete assertion with an explicit discriminator. A
status-`201`-only assertion is rejected for every property whose
falsifier multiple causes satisfy (the falsifiability trap).

- **P1 — identity-by-construction (slug honored verbatim).**
  Input a **position-segmented** construction-known slug (e.g.
  `demo-root-m1-t2`-shaped); assert the receipt's echoed `slug=`
  is **byte-identical** to the input. *Discriminator:* a
  position-segmented slug surviving verbatim distinguishes the
  exact-slug branch from server-side descendant generation (which
  would mint its own position number) and from any
  interpretation/derivation. *Verified by:*
  [`internal/api/handlers.go`](../../../../../internal/api/handlers.go)
  `:197-201`;
  [`cmd/workstream-tracker/register.go`](../../../../../cmd/workstream-tracker/register.go)
  receipt prints `slug=` (`:89-91`).

- **P2 — idempotent create-or-attach on `(slug, actor, active)`.**
  Assert (a) a same-`(slug, actor)` repeat **while active**
  collapses to the **same** `work_instance_id`; (b) a
  **different actor, same slug** yields a **distinct**
  `work_instance_id`. (b) is the key discriminator: if attach
  were keyed on slug alone the second actor would collapse onto
  the first. The **`active`/serial-resume leg is out of t2's
  charter** (see Rationale below). *Verified by:*
  [`internal/api/handlers.go`](../../../../../internal/api/handlers.go)
  `:226-234`.

- **P3 — no natural-language-resolution path engaged.** Not a
  bare `201`. The signal is the **conjunction**: (i) P1's
  verbatim round-trip of a position-segmented slug (rules out
  interpretation/derivation) **and** (ii) P2(a)'s idempotent
  *attach* on repeat (201 + same id), which positively
  fingerprints the **exact-slug create-or-attach branch** —
  ruling out the bare-root-create branch (which 409s on a repeat)
  and the descendant-generation branch (no caller exact slug at
  all). *Verified by:*
  [`internal/api/handlers.go`](../../../../../internal/api/handlers.go)
  bare-root `errSlugConflict` on repeat (`:216-224`) vs.
  exact-slug attach (`:197-201`, `:226-234`).

- **P4 — server repo-blind.** Assert a **paired** positive /
  negative fingerprint, not a single status: (a) a well-formed
  slug that **cannot name any real plan-tree doc** (a clearly
  synthetic root) → `201` + verbatim (existence is never checked);
  (b) a **grammatically malformed** slug (e.g. an uppercase /
  underscore token, or a bare `root-p1`) → `400` (the only gate
  is grammar). The pair fingerprints "grammar-only, repo-blind,"
  which a single `201` cannot. *Verified by:*
  [`internal/api/handlers.go`](../../../../../internal/api/handlers.go)
  `:58-62`;
  [`internal/slugs/slugs.go`](../../../../../internal/slugs/slugs.go)
  `IsWellFormed` repo-blind grammar (`:73-117`).

- **P5 — handshake-free / no narration handshake.** *Structural,
  not a separate status assertion.* The harness invokes
  `runRegister` with only `--slug`/`--actor`/`--server` and
  `noEnv` — there is no prompt, no interpretation input, no
  confirm step — and still yields a real `201` receipt. The
  absence of a handshake is shown by the invocation shape (no
  prompt parameter, no confirm interaction; ties to SD2); the
  positive receipt is the proof it needs none. *Verified by:*
  [`cmd/workstream-tracker/register.go`](../../../../../cmd/workstream-tracker/register.go)
  `runRegister` has no resolve/confirm (`:39-93`);
  [`docs/agents/local/session-registration.md`](../../../../../docs/agents/local/session-registration.md)
  deterministic path drops Resolve/Confirm (`:82-88`).

**Rejected (per property).** *Status-code-only assertions for P3
and P4* — the named falsifiability trap: `201` has many causes
(correct exact-slug attach, an unrelated accept) and cannot, alone,
prove NL-resolution was not engaged or that acceptance was
grammar-gated rather than repo-gated. Replaced by the
conjunction / paired-fingerprint discriminators above.

### SD4 — Idempotency-assertion shape: the determinism delta, not a re-derivation

**Decision.** The determinism file asserts the **`(slug, actor)`
key discriminator** (different-actor-same-slug → distinct id,
SD3-P2(b)) anchored by a **single** same-`(slug, actor)` collapse
(SD3-P2(a)). It does **not** re-derive
`TestRegisterCommandSuccessAndIdempotentRepeat` from scratch and
does **not** drive a terminal transition. It cross-references the
named template rather than duplicating its coverage.

**Rejected.** *C1 — re-prove same-actor collapse from scratch:*
the inherited contract explicitly says "adds the
determinism-specific assertions … rather than re-testing
idempotency from scratch." *C3 — full key including the
`active`/serial-resume leg* (register → `complete` → re-register →
expect a fresh id): driving a terminal transition exercises the
**completion lifecycle**, not the **deterministic identity** the
epic scopes t2 to.

**Rationale / accepted residual.** The `active` qualifier in the
`(slug, actor, active)` key is satisfied by asserting collapse
*while active* — the only state the proof creates — not by also
proving terminal-then-resume. The serial-resume idempotency leg
is a lifecycle property, exercised by the completion-handshake
surface, **not** t2's determinism charter ("'Deterministic' stays
epic-scoped: the deterministic identity step for a session that
*did* launch"). Recorded as a deliberate, documented scoping
boundary so a reviewer flagging the omission gets the
hold-the-line rationale rather than a re-litigation.

**Verified by:**
[`cmd/workstream-tracker/register_test.go`](../../../../../cmd/workstream-tracker/register_test.go)
`TestRegisterCommandSuccessAndIdempotentRepeat` (the
not-re-derived baseline, `:44-80`);
[`internal/api/handlers.go`](../../../../../internal/api/handlers.go)
active-pair idempotency (`:226-234`); epic-scoped "deterministic"
([`../README.md`](../README.md) Cross-Task Invariants).

### SD5 — `design/v0.1-design.md` is a frozen v0.1 record, not load-bearing

**Decision.** `design/v0.1-design.md` describes only the v0.1 end
state and is **not an authority, currency constraint, or
reconciliation target** for t2. Every reality-check input above is
grounded directly in merged code. The milestone's "Documentation
Currency — currency-check §4, no edit expected" framing is
satisfied trivially: t2 edits only test files under
`cmd/workstream-tracker/`, and the design doc is historical.

**Verified by:** the merged-code Reality-check inputs above; no
`design/v0.1-design.md` citation by design (mirrors t1 scoping
SD6 and the standing project convention).

### SD6 — this scoping doc carries no `Status` field

**Decision.** No `Status` frontmatter, conforming to the spec
over any spawn-instruction `Status` wording.

**Verified by:**
[`task-plan.md`](../../../../../spec/planning/task-plan.md) "Goal:
scoping doc + plan doc" (scoping docs do not carry a `Status`
field; `:122-129`).

## Open decisions to make at plan-drafting (handoff)

**None remain open.** Every decomposed HOW decision (SD1–SD6) was
resolved in this scoping session. What remains for the
plan-drafting / implementing PR is **implementation-altitude**,
not an open scoping decision:

- The literal test function names, the exact synthetic slug
  constants for P4, and the precise malformed-slug input for
  P4(b) — implementation choices the plan states as contracts and
  the PR fixes in code.
- Whether the thin SD2 receipt-parsing helper is worth extracting
  vs. inlined per case — a code-readability call, not a scoping
  decision, bounded by SD2 (never a producer abstraction).

## Plan-structure handoff

For the plan-drafting session (scoping references the plan's
sections by name per
[`task-plan.md`](../../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns" `:131-167`; it does not duplicate
plan-owned content here):

- t2 is a **test-only task** (no product code; D1 no-new-surface
  locked). The plan's Contracts state the assertion conditions
  (the SD3 property→discriminator map); exact prose/names are the
  implementer's.
- **Files the plan will contract:** a new test file under
  [`cmd/workstream-tracker/`](../../../../../cmd/workstream-tracker/)
  (SD1); [`register_test.go`](../../../../../cmd/workstream-tracker/register_test.go)
  byte-unchanged;
  [`design/v0.1-design.md`](../../../../../design/v0.1-design.md)
  not touched (SD5);
  no `internal/**` / `cmd/**` product change (D1).
- **Validation gate the plan will carry:** `go test
  ./cmd/workstream-tracker/...` green; the SD3 falsifiability
  discriminators are each *run*, not announced (the
  `validation-honesty` trap); a diff check that the named
  template and product code are unchanged (no-new-surface is a
  surface check, not a literal allowlist — status-closeout edits
  below are expected).
- **Reality-check re-confirm:** the full enumerated
  "Reality-check inputs" list above is the falsifier set; a stale
  entry anywhere fails the gate. This explicitly includes the
  `(slug, actor)` co-working independence branch and the
  bare-root-409-on-repeat discriminator P3 leans on, not only the
  verbatim-honor citation.
- **Milestone-terminal closeout the plan must own.** t2 is the
  **last task in m1**, so t2's terminal PR is the
  **milestone-terminal PR**: it owns its own plan `Status`
  transition, the [`../README.md`](../README.md) t2 Task Status
  row → terminal, **and** the batch deletion of the transient
  `scoping/` subfolder (both this doc and
  [`./t1-deterministic-path-contract.md`](./t1-deterministic-path-contract.md),
  per t1's plan Documentation Currency, which explicitly defers
  scoping deletion to the milestone-terminal PR — i.e. here). The
  plan-drafting session must verify t1's status is terminal
  before scheduling that milestone-terminal action; finishing
  t2's Contracts is task-terminal, but the scoping-deletion +
  milestone-row reconciliation is **milestone-terminal** and
  gated on the sibling (t1) state.
- **PR-count:** a single test PR is the estimate (one new test
  file), subject to the
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "PR-count predictions need a branch test" rule — an estimate,
  not a commitment.

## Related docs

- [`../t2-determinism-proof-harness.md`](../t2-determinism-proof-harness.md)
  — the durable t2 skeleton carrying the locked inherited
  contract (this scoping doc pairs with it).
- [`./t1-deterministic-path-contract.md`](./t1-deterministic-path-contract.md)
  — t1's paired scoping doc; its "Reality-check inputs" enumerate
  the determinism facts t2's assertions mirror. Both delete in
  the milestone-terminal batch (t2's PR).
- [`../t1-deterministic-path-contract.md`](../t1-deterministic-path-contract.md)
  — t1's landed plan; carries the durable contract (C1–C4) t2
  proves and the milestone-terminal scoping-deletion deferral.
- [`../README.md`](../README.md) — m1 milestone doc; locks t2's
  WHAT, Cross-Task Invariants/Decision D1, Documentation Currency.
- [`../../README.md`](../../README.md) — parent epic;
  Cross-Cutting Invariants and the epic-scoped meaning of
  "deterministic".
- [`../../../../../spec/planning/shared.md`](../../../../../spec/planning/shared.md)
  — "Deterministic assertion when the slug is construction-known"
  (the contract t2 mirrors), "Decompose options into shapes",
  "Falsifiability check".
- [`../../../../../spec/planning/task-plan.md`](../../../../../spec/planning/task-plan.md)
  — scoping-vs-plan ownership, reality-check pass, the
  scoping-doc-no-Status rule.
- [`../../../../../cmd/workstream-tracker/register_test.go`](../../../../../cmd/workstream-tracker/register_test.go),
  [`../../../../../cmd/workstream-tracker/register.go`](../../../../../cmd/workstream-tracker/register.go),
  [`../../../../../internal/api/handlers.go`](../../../../../internal/api/handlers.go),
  [`../../../../../internal/slugs/slugs.go`](../../../../../internal/slugs/slugs.go),
  [`../../../../../internal/registerclient/client.go`](../../../../../internal/registerclient/client.go)
  — the merged code grounding the reality-check inputs and SD1–SD4.
  (`design/v0.1-design.md` is intentionally **not** listed — SD5.)
