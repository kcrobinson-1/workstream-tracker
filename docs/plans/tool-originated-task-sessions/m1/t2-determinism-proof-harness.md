---
slug: tool-originated-task-sessions-m1-t2
Status: In progress
short_description: Determinism proof harness
---

# m1 t2 — Determinism Proof Harness

## Context preamble

The deterministic, handshake-free, slug-by-construction
registration path is now a sanctioned contract (t1 landed it in
the plan-doc spec and the session-registration agent rule). But a
contract is only de-risked when something *proves* the merged code
actually behaves the way the contract claims — otherwise m2 would
have to co-develop the tool-originated-session UX *and* discover
whether the registration path it depends on really is
deterministic, at the same time.

t2 closes that gap with a single end-to-end test: drive the real
exact-slug create-or-attach path the way a construction-time slug
producer will (an explicit slug, no natural-language handshake),
and assert each determinism property t1's contract names — the
slug is honored verbatim, create-or-attach is idempotent on the
`(slug, actor)` identity while active, no natural-language
resolution is engaged, and the server stays repo-blind. It is a
**test-only** task: it adds no product code and no new
registration surface (the mechanism already exists); it touches
only the CLI end-to-end test surface under
`cmd/workstream-tracker/`. The deliverable is the proof artifact
that lets m2 consume a *proven* path rather than an asserted one.

The HOW decisions were resolved at scoping (SD1–SD6 in the
paired, transient scoping doc
[`scoping/t2-determinism-proof-harness.md`](./scoping/t2-determinism-proof-harness.md));
this plan carries the durable contract those decisions produced,
with its own code-grounded verification.

## Inherited Contract (locked by the milestone)

Binding input from [`../README.md`](../README.md) (Task
Contracts, Cross-Task Invariants, Cross-Task Decision D1); not
loosened here.

**End result.** An end-to-end test proves that, given a
construction-known slug supplied as an explicit argument with **no
narration handshake**, a work-instance registers deterministically
against the real exact-slug create-or-attach path:
identity-by-construction (slug honored verbatim), idempotent
create-or-attach on `(slug, actor, active)`, no
natural-language-resolution path engaged, server repo-blind. The
slug producer is the **minimal explicit `--slug` argument / test
harness** — the stand-in for m2's future construction-time
producer — scoped as the minimum to validate the contract, **not
a durable UX**.

**Preserves:** no product code surface added; the existing
interactive-path / CLI tests
([`cmd/workstream-tracker/register_test.go`](../../../../cmd/workstream-tracker/register_test.go)
and siblings) unchanged. t2 adds the *determinism-specific*
assertions rather than re-testing idempotency from scratch.

**Sibling interface.** Consumes t1's named determinism properties
as the assertions it must demonstrate; produces the **proof
artifact** that de-risks the registration contract before the m2
posture shift.

**Constraints carried from the milestone.** No new registration
surface — no endpoint, request/response field, schema, code path,
or CLI flag (Cross-Task Decision D1). "Deterministic" stays
epic-scoped: the proof asserts only the deterministic identity
step for a session that *did* launch, never "the launch cannot
fail" and never enrichment as load-bearing for attachment.

## Goal

Land a single end-to-end test that exercises the merged exact-slug
create-or-attach path through the CLI with a construction-known
slug and **no handshake**, asserting each determinism property
t1's contract names with a discriminator strong enough that the
assertion fails if the property does not hold (no
status-code-only assertion for the properties multiple causes
could satisfy). No product code, no new registration surface, the
existing CLI tests byte-unchanged.

## Cross-Cutting Invariants

These thread the Contracts, the test file, and the validation
gate; a single-contract self-review misses a violation that only
shows when two hold each other in tension.

- **No product code surface (D1).** Test files only. No
  `internal/**` or `cmd/**` *product* change; no schema,
  endpoint, or CLI flag; no new exported test helper (the named
  template's unexported helpers are reused as-is from the same
  `package main`).
- **Existing CLI tests byte-unchanged.**
  [`register_test.go`](../../../../cmd/workstream-tracker/register_test.go)
  is not modified; the determinism proof is a purely additive new
  file (scoping SD1).
- **"Deterministic" stays epic-scoped.** Every assertion covers
  only the deterministic identity step for a session that *did*
  launch — not "the launch/producer cannot be wrong," not
  enrichment made load-bearing. Decisive for C2's exclusion of
  the `active`/serial-resume leg.
- **Falsifiability discipline.** No status-code-only assertion
  for "no NL-resolution engaged" (C3) or "server repo-blind"
  (C4); each carries a positive/negative discriminator a single
  `201` cannot satisfy (per
  [`shared.md`](../../../../spec/planning/shared.md)
  "Falsifiability check," `:944-977`).
- **Assertions mirror t1's landed named properties.** No property
  is asserted that t1's two carriers
  ([`shared.md`](../../../../spec/planning/shared.md) `:181-200`;
  [`session-registration.md`](../../../../docs/agents/local/session-registration.md)
  `:71-123`) do not express; t2 demonstrates the contract, it
  does not extend it.

## Contracts

The conditions that must hold for the implementation to satisfy
this plan. Exact test names, slug constants, and helper
extraction are the implementer's choice (contract altitude); only
the conditions below bind. Each maps to a determinism property
t1's contract names (scoping SD3's property→assertion map).

### C1 — Identity-by-construction: verbatim round-trip

A **position-segmented** construction-known slug (a root followed
by `mN`/`tN` segments, e.g. the named template's
`demo-root-m1-t2` shape) supplied through the bare `--slug`
argument registers, and the slug echoed in the success receipt is
**byte-identical** to the supplied slug. The position-segmented
shape is load-bearing: a slug carrying its own `mN`/`tN`
surviving verbatim distinguishes the exact-slug branch (honors
the caller's slug) from any derivation or server-side descendant
generation (which would mint its own position number). *Verified
by:*
[`internal/api/handlers.go`](../../../../internal/api/handlers.go)
`insertRegister` exact-slug branch (`slug = req.ExactSlug`,
"honor the caller's slug verbatim. No derivation…", `:197-201`);
[`cmd/workstream-tracker/register.go`](../../../../cmd/workstream-tracker/register.go)
`runRegister` success receipt prints `slug=` (`:89-91`).

### C2 — Idempotent create-or-attach on the `(slug, actor)` identity, while active

Two assertions: (a) a repeat invocation with the **same
`(slug, actor)`** while the prior work-instance is still active
collapses to the **same** `work_instance_id`; (b) a **different
actor on the same slug** yields a **distinct**
`work_instance_id`. (b) is the discriminator that proves the
attach identity is `(slug, actor)`, not slug alone — if attach
were keyed on slug, the second actor would wrongly collapse onto
the first. The `active`/serial-resume leg (terminal-then-resume)
is **out of scope** (Out Of Scope; Risk Register accepted
residual): driving a terminal transition exercises the completion
lifecycle, not the deterministic identity step the epic scopes
t2 to. *Verified by:*
[`internal/api/handlers.go`](../../../../internal/api/handlers.go)
`insertRegister` active-pair idempotency (`activeWorkInstanceID`
→ return existing id/slug, `:226-234`).

### C3 — No natural-language-resolution path engaged

Demonstrated by a **conjunction**, not a bare status code: C1's
verbatim round-trip of a position-segmented slug (rules out
interpretation/derivation) **and** C2(a)'s idempotent *attach* on
repeat (`201` + same id), which positively fingerprints the
exact-slug create-or-attach branch — the bare-root-create branch
`409`s on a repeat and the descendant-generation branch never
carries a caller exact slug, so neither can produce C1∧C2(a). A
status-`201`-only assertion is explicitly insufficient (the
falsifiability trap: many causes yield `201`). *Verified by:*
[`internal/api/handlers.go`](../../../../internal/api/handlers.go)
bare-root `errSlugConflict` on a repeat (`:216-224`) vs.
exact-slug attach (`:197-201`, `:226-234`).

### C4 — Server repo-blind: paired grammar-gate fingerprint

A **paired** positive/negative assertion: (a) a well-formed slug
that **cannot name any real plan-tree doc** (a clearly synthetic
root) registers `201` and round-trips verbatim — existence is
never checked; (b) a **grammatically malformed** slug (a token
with an uppercase letter or underscore, or a bare `root-p1`) is
rejected `400`. The pair fingerprints "the only acceptance gate
is grammar, and it is repo-blind," which neither leg alone
proves. *Verified by:*
[`internal/api/handlers.go`](../../../../internal/api/handlers.go)
exact-slug grammar gate (`IsWellFormed` else `400`, `:58-62`);
[`internal/slugs/slugs.go`](../../../../internal/slugs/slugs.go)
`IsWellFormed` ("does not verify the slug names a real plan-tree
doc — the server is repo-blind", `:73-117`).

### C5 — Handshake-free (structural)

The proof invokes `runRegister` with only `--slug` / `--actor` /
`--server` and an empty environment (`noEnv`) — no prompt, no
resolve step, no confirm interaction — and still obtains a real
`201` receipt. "No narration handshake" is shown by the
invocation shape (there is no prompt parameter and no confirm
gate to engage) plus the positive receipt; it is **not** a
separate status assertion. *Verified by:*
[`cmd/workstream-tracker/register.go`](../../../../cmd/workstream-tracker/register.go)
`runRegister` contains no resolve/confirm step (`:39-93`);
[`docs/agents/local/session-registration.md`](../../../../docs/agents/local/session-registration.md)
deterministic path drops Resolve/Confirm (`:71-123`).

### C6 — Additive, no new surface

The only change under `cmd/workstream-tracker/` is one new test
file in the existing `package main`; product code and
[`register_test.go`](../../../../cmd/workstream-tracker/register_test.go)
are byte-unchanged; the named template's helpers
(`startAPIServer`, `isolateCaches`, `widPattern`, `noEnv`,
`runRegister`) are reused without adding any exported symbol.
*Verified by:*
[`cmd/workstream-tracker/register_test.go`](../../../../cmd/workstream-tracker/register_test.go)
unexported package-`main` helpers (`:23-38`, `:40`, `:42`,
`:140-145`); milestone Cross-Task Decision D1
([`../README.md`](../README.md)).

## Files to touch — new / modify / intentionally not touched

> *Estimate of expected shape. Implementation may revise a
> structural call; the Contracts above bind, this inventory does
> not.*

- **New:** one determinism test file under
  [`cmd/workstream-tracker/`](../../../../cmd/workstream-tracker/)
  (e.g. `register_determinism_test.go`) in `package main`
  (scoping SD1).
- **Modify at closeout (milestone-terminal — see Documentation
  Currency):** this plan's own `Status` frontmatter; the
  [`../README.md`](../README.md) t2 Task Status row to its
  terminal state; the batch deletion of the transient
  [`scoping/`](./scoping/) subfolder (**both**
  [`scoping/t2-determinism-proof-harness.md`](./scoping/t2-determinism-proof-harness.md)
  **and**
  [`scoping/t1-deterministic-path-contract.md`](./scoping/t1-deterministic-path-contract.md)),
  because t2's terminal PR is the m1-terminal PR (t1's plan
  Documentation Currency explicitly defers scoping deletion to
  here).
- **Intentionally not modified (expected):**
  [`register_test.go`](../../../../cmd/workstream-tracker/register_test.go)
  (byte-unchanged); `internal/**`, `cmd/**` product code (no
  code — D1);
  [`design/v0.1-design.md`](../../../../design/v0.1-design.md) (a
  frozen v0.1 record, not an authority or reconciliation target —
  scoping SD5); the backlog (no entry — D1 / milestone Backlog
  Impact).

## Validation Gate

Test-only; the gate is verification, not a build.

1. **Test green.** `go test ./cmd/workstream-tracker/...` passes
   — the new determinism test and every existing test in the
   package. (Falsifier: any failure or compile break fails the
   gate.)
2. **Discriminators actually exercised.** Each of C1–C5 is run
   and asserted in the test, not announced — specifically C1's
   byte-identical receipt slug, C2(a) same-id collapse, C2(b)
   distinct-id for a different actor, C3's C1∧C2(a) conjunction,
   C4's `201`-verbatim / `400`-malformed pair, C5's
   `noEnv`-no-prompt invocation yielding a real `201`. (Falsifier:
   a "covered" claim for any discriminator the test does not
   actually assert fails the gate — `validation-honesty`.)
3. **No-new-surface diff check.** `git diff` shows
   `register_test.go` and all `internal/**` / `cmd/**` product
   code byte-unchanged; the only code change is the one new test
   file. This is a *surface* check, not a literal file allowlist
   — the milestone-terminal closeout edits named in Documentation
   Currency are expected and are not a violation. (Falsifier: any
   product / schema / endpoint / CLI-flag change, or a
   modification inside `register_test.go`, fails the gate.)
4. **Reality-check re-confirm.** *Every* input in the scoping
   doc's "Reality-check inputs" list — the full enumerated set,
   not only the subset surfaced in Contracts — still holds
   against then-merged code. This explicitly includes the
   `(slug, actor)` co-working independence branch (C2(b)), the
   bare-root-`409`-on-repeat discriminator (C3), and the
   `IsWellFormed`-`400`-on-malformed grammar gate (C4), not only
   the verbatim-honor citation. That full named list is the
   falsifier set; a stale entry anywhere in it fails the gate.

## Self-Review Audits

Diff surface maps to these seeded audits
([catalog](../../../../docs/agents/local/self-review-catalog.md)):

- **validation-honesty** — the falsifiability discriminators
  (C3, C4) and the full Validation Gate must be run end-to-end on
  the actual test before any "checks passed" claim; a prose
  success claim without `go test` output and the diff inspection
  is the trap.
- **readiness-gate-truthfulness** — the `In draft → Proposed`
  promotion claim is valid only if the gate's named conditions
  were actually verified, not announced best-effort.
- **trigger-map-currency** — adding a test file under
  `cmd/workstream-tracker/` adds no top-level dir and does not
  restructure the tree; confirm the
  [`AGENTS.md`](../../../../AGENTS.md) "Mandatory pre-edit reads"
  trigger map still matches reality (expected: unchanged — the
  map keys on `spec/**`, `docs/plans/**`, not `cmd/**`).

## Documentation Currency

Status-bearing docs the implementing PR must keep current (per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Plan-to-PR Completion Gate" and
[`shared.md`](../../../../spec/planning/shared.md) "Plan-doc
Status"). These edits are expected closeout, not a
"no-new-surface" violation (Validation Gate step 3).

- **This plan's own `Status`.** The implementing PR transitions
  it `Proposed → In progress` on start and `→ Landed` at close —
  no `Validating` (test-only; no post-merge gate). Frontmatter is
  authoritative.
- **The m1 `README.md` t2 Task Status row** to its terminal
  state, narrowing the adjacent task-status prose if needed.
- **Milestone-terminal scoping deletion.** t2 is the **last task
  in m1**, so t2's terminal PR is the **milestone-terminal PR**:
  it deletes the entire transient
  [`scoping/`](./scoping/) subfolder (both t1's and t2's scoping
  docs) per
  [`planning-doc-location.md`](../../../../spec/planning-doc-location.md)
  "scoping/ subfolder is transient" and
  [`task-plan.md`](../../../../spec/planning/task-plan.md)
  "Scoping owns / plan owns". The implementing session **must
  verify t1's plan `Status` is terminal** (`Landed`) before this
  milestone-terminal action — finishing t2's Contracts is
  task-terminal, but the scoping-batch deletion + milestone-row
  reconciliation is milestone-terminal and gated on the sibling
  (t1) state. (t1 is currently `Landed`; re-verify at
  implementation time.)

## Out Of Scope

- **The `active`/serial-resume idempotency leg** (register →
  terminal transition → re-register expecting a fresh id). That
  exercises the completion lifecycle, not the deterministic
  *identity* step the epic scopes t2 to; the `active` qualifier
  in the `(slug, actor, active)` key is satisfied by asserting
  collapse *while active* (the only state the proof creates). A
  documented boundary, not an omission (scoping SD4).
- **The m2 tool-originated-session UX** — node affordance, real
  construction-time slug producer, session-start hook. t2's
  producer is the bare `--slug` stand-in.
- **Any change to product registration code or the interactive
  path** — t2 only exercises already-merged code (D1).
- **Re-deriving `TestRegisterCommandSuccessAndIdempotentRepeat`**
  — the inherited contract says add the determinism-specific
  assertions, not re-test idempotency from scratch (scoping SD4).
- **Editing `design/v0.1-design.md`** — frozen v0.1 record; not
  load-bearing or reconciled against (scoping SD5).

## Risk Register

- **A status-code-only assertion silently passes for the wrong
  reason.** The named falsifiability trap for C3/C4. Mitigation:
  C3's conjunction and C4's paired fingerprint are contract
  conditions, and Validation Gate step 2 requires each
  discriminator be actually exercised.
- **Accepted residual: the `active`/serial-resume leg is
  unproven by t2.** Conscious scoping boundary (C2, Out Of
  Scope, scoping SD4) — recorded here so a reviewer flagging it
  gets the epic-scoped-"deterministic" hold-the-line rationale
  rather than a re-litigation; the leg is a completion-lifecycle
  property exercised elsewhere.
- **The new test couples to the named template's unexported
  helpers.** If a future change renames/moves `startAPIServer` /
  `isolateCaches` / `widPattern` / `runRegister`, the proof
  breaks at compile time (loud, not silent) — acceptable: it is
  same-package reuse, and a compile break is a correct signal,
  not a hidden regression.

## Backlog Impact

None. Per [`../README.md`](../README.md) "Backlog Impact", m1's
implementing PRs touch no backlog entry; the
graduate/split/carve-outs landed in the epic promotion PR. The
determinism best-effort gap remains tracked by the existing
`deterministic-interactive-registration` backlog entry; t2
proves the contract, it does not close that entry.

## Related Docs

- [`./scoping/t2-determinism-proof-harness.md`](./scoping/t2-determinism-proof-harness.md)
  — the paired transient scoping doc (SD1–SD6 deliberation +
  rejected alternatives; deletes at the milestone-terminal PR,
  which is t2's).
- [`./scoping/t1-deterministic-path-contract.md`](./scoping/t1-deterministic-path-contract.md)
  — t1's scoping doc; its "Reality-check inputs" enumerate the
  determinism facts t2's assertions mirror.
- [`./t1-deterministic-path-contract.md`](./t1-deterministic-path-contract.md)
  — t1's landed plan; the durable contract (C1–C4) t2 proves and
  the milestone-terminal scoping-deletion deferral.
- [`../README.md`](../README.md) — m1 milestone doc; locks t2's
  WHAT, Cross-Task Invariants, Cross-Task Decision D1,
  Documentation Currency.
- [`../../README.md`](../../README.md) — parent epic;
  Cross-Cutting Invariants and the epic-scoped meaning of
  "deterministic".
- [`../../../../spec/planning/shared.md`](../../../../spec/planning/shared.md)
  — "Deterministic assertion when the slug is construction-known"
  (the contract t2 mirrors); "Falsifiability check".
- [`../../../../spec/planning/task-plan.md`](../../../../spec/planning/task-plan.md)
  — scoping-vs-plan ownership, reality-check pass, the
  promotion gate.
- [`../../../../cmd/workstream-tracker/register_test.go`](../../../../cmd/workstream-tracker/register_test.go),
  [`../../../../cmd/workstream-tracker/register.go`](../../../../cmd/workstream-tracker/register.go),
  [`../../../../internal/api/handlers.go`](../../../../internal/api/handlers.go),
  [`../../../../internal/slugs/slugs.go`](../../../../internal/slugs/slugs.go),
  [`../../../../internal/registerclient/client.go`](../../../../internal/registerclient/client.go)
  — the merged code grounding C1–C6.
