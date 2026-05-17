---
slug: stub-children-on-parent-promotion
Status: In draft
short_description: Seed skeleton child docs when a parent doc promotes
---

# Stub Children on Parent Promotion

## Context

When a parent plan doc (epic or milestone) passes its
`` `In draft` → `Proposed` `` promotion gate, its children exist
only as names in the parent's contract section — each child slug
has no doc until someone drafts it just-in-time. Two costs
follow: planned-but-unstarted work is invisible in the rendered
tree until its drafting session runs, and work-instance
auto-registration has no frontmatter slug to derive from before
any drafting happens.

This is being done now because the gap was captured during the
`workstream-tracker-1-0-m1-t1` review as a backlog entry that
explicitly depends on m1-t1's now-landed exact-slug
create-or-attach registration path; with that dependency
satisfied, the entry is graduatable. The surface is entirely the
planning spec this repo dogfoods and consumer projects vendor —
the parent-doc promotion gate and the parent-doc child-contracts
rule in [`shared.md`](../../../spec/planning/shared.md), plus the
slug-generation rule the reconciliation narrows. No product
code, schema, API, or lifecycle token changes (a slug+Status
stub is an already-supported render case).

The one load-bearing decision — D1 in
[`scoping/README.md`](scoping/README.md), loosening the
"descendant slugs are server-generated" rule — was ratified
2026-05-17 (input I1 resolved): a seeded slug is *declared* in
frontmatter with no server-side creation, *asserted* by the
child session at registration, and server-generation stays as
the no-pre-declared-slug fallback.

## Goal

When a PR flips a parent doc to `Proposed`, the same PR also
contains a skeleton doc for every child the parent's contract
section names, each initialized with a canonical slug, a
level-appropriate short and long description, and the parent's
WHAT-contract block plus any illustrative examples for that
child; and the parent doc presents its children as a table that
also states, at a high level, how each child's deliverable feeds
its siblings' inputs.

The spec change distinguishes two registers (D6): **artifact
requirements** — checkable facts about the seeded files and the
parent doc — and **expected agent behavior** — the prompted
seeding act, whose miss is an accepted, observable residual, not
a guarantee.

Verifiable when: the parent-doc gate's prose directs stub
seeding in the promoting PR; the slug-generation rule
unambiguously says a pre-declared frontmatter slug is asserted
(not server-generated), file-write performs no server-side
creation, and server-generation is the no-pre-declared-slug
fallback; the per-level files reference the new obligation
without restating it; and the added prose frames the seeding act
as a prompted, observable best-effort obligation rather than an
enforced guarantee.

## Contracts

The final shape the spec must establish — WHAT, not the prose
that establishes it (exact wording is the implementing PR's
choice). Split per D6.

### C1 — Artifact requirements (checkable; `must`)

The spec must establish these as checkable facts about a
promoting PR's diff:

- **A1 — A stub exists per named child.** For every child named
  in the parent's `Milestone Contracts` / `Task Contracts`
  section, a doc exists at that child's layout path in the same
  PR that flips the parent to `Proposed`.
- **A2 — Stub frontmatter is canonical and complete.** Each stub
  carries `slug` (author-supplied at promotion per C2),
  `Status: In draft` (D2), and `short_description`.
- **A3 — Stub body carries the inherited contract.** The
  parent's WHAT-contract block for that child, and any
  illustrative examples the parent states for it, are present in
  the stub so the child's drafting session starts from the
  locked contract.
- **A4 — Parent presents children as a wiring table.** The
  parent's child-contracts section presents each child in a
  table giving its short and long description and a high-level
  statement of how that child's deliverable feeds sibling
  children's inputs (the inter-child interface story, not a bare
  name list).
- **A5 — Level-appropriate description register.** A stub's
  descriptions are written at its level: milestone stubs in
  product terms ("enable the user to X", "lay the technical
  groundwork for m2"); task/phase stubs in technical terms still
  grounded in the product outcome served.
- **A6 — Stub exemption.** A doc that is a parent-promotion stub
  is exempt from "Required and optional sections" and "Plans
  describe contracts, not implementation" until its own
  `` `In draft` → `Proposed` `` drafting session (D3).

Verified by:
[`../../../spec/planning/shared.md:551-592`](../../../spec/planning/shared.md)
(parent-doc child contracts — the WHAT block A3 copies and the
section A4 restructures);
[`../../../spec/planning/epic.md:43-67`](../../../spec/planning/epic.md)
and
[`../../../spec/planning/milestone.md:109-125`](../../../spec/planning/milestone.md)
(per-level Milestone/Task Contracts + Status sections A1/A2/A4
attach to);
[`../../../spec/planning/shared.md:265-286`](../../../spec/planning/shared.md)
(Status lifecycle backing A2/D2);
[`../../../spec/planning/task-plan.md:169-214`](../../../spec/planning/task-plan.md)
(the required-sections rule A6 grants a stub exemption from).

### C2 — Slug-generation reconciliation (`shared.md` "Slug generation")

The "Descendant slugs are server-generated" rule is **loosened,
not replaced**. A parent-promotion-seeded slug is **declared in
the stub's frontmatter** — author-supplied, format-validated,
exactly as a root slug is. The file-write is declaration only:
it performs **no server-side creation** (no node, work-instance,
or allocation call). The child session **asserts** the declared
slug to the server when it registers, via the exact-slug
create-or-attach flow (m1-t1, landed) — "use this slug," not
"give me a slug" — depended on, not replaced. Server
slug-generation remains the path for any registration with **no
pre-declared frontmatter slug**. A pre-declared slug is asserted,
never server-generated; the two paths are mutually exclusive, so
there is no allocator to diverge from.

Verified by:
[`../../../spec/planning/shared.md:127-131`](../../../spec/planning/shared.md)
(the rule being narrowed);
[`../../../spec/planning/shared.md:113-119`](../../../spec/planning/shared.md)
(slug immutability/identity the narrowing preserves);
[`../workstream-tracker-1-0/m1-t1-multi-wi-per-slug.md:113-142`](../workstream-tracker-1-0/m1-t1-multi-wi-per-slug.md)
(exact-slug create-or-attach flow);
[`../workstream-tracker-1-0/m1-t1-multi-wi-per-slug.md:377-389`](../workstream-tracker-1-0/m1-t1-multi-wi-per-slug.md)
(m1-t1 leaves stub-seeding to this entry, not replacing the
flow).

### C3 — Expected agent behavior (best-effort; observable; not guaranteed)

The parent-doc `` `In draft` → `Proposed` `` promotion gate
gains a step directing the promoting agent to seed the C1 stubs
from the locked parent contracts in the promoting PR, and to not
clobber a child already drafted or advanced when the gate
re-runs on a re-opened parent. The step is framed — in the
register the gate already uses — as a prompted obligation whose
missed or imperfect execution is an accepted, observable
residual addressed through prompt engineering and tree
visibility, **not** a determinism guarantee: the spec is
conformant by specifying the prompted behavior and keeping
misses observable, and does not promise they cannot be missed.

Verified by:
[`../../../spec/planning/shared.md:288-356`](../../../spec/planning/shared.md)
(the parent-doc gate — itself a prompted self-review, the
register C3 joins; re-run-from-scratch behavior the no-clobber
clause guards);
[`../../backlog.md`](../../backlog.md)
`deterministic-interactive-registration` (the "observable
best-effort … not deterministic" project vocabulary reused).

### C4 — No rule-body duplication; layered authority preserved

The new gate step cites "Parent-doc child contracts" and the
slug-generation rule by exact section title and restates neither
body; the per-level files
([`epic.md`](../../../spec/planning/epic.md),
[`milestone.md`](../../../spec/planning/milestone.md)) reference
the obligation rather than carrying it. A reviewer finding a
rule body restated inside the gate or a per-level file treats
that as a defect, not safety-by-copying.

Verified by:
[`../../../spec/planning/shared.md:551-558`](../../../spec/planning/shared.md)
(the layering discipline — a multi-doc-type rule lives in
`shared.md` once, per-level files reference it);
[`../promotion-gate-explicit-checklist/README.md:136-149`](../promotion-gate-explicit-checklist/README.md)
(the cite-by-name, no-duplication precedent for gate edits).

## Cross-Cutting Invariants

- **Declare vs. assert; no allocation for pre-declared slugs.**
  A pre-declared child slug is never server-allocated: file-write
  declares identity (no server-side creation), the child session
  asserts that slug at registration via exact-slug
  create-or-attach, and server slug-generation runs only when no
  frontmatter slug was pre-declared. No path may make file-write
  trigger server-side creation, or make the server generate over
  an already-declared slug — either re-opens the divergence C2
  dissolves.
- **Stub ≠ work-instance.** Seeding a stub never creates a
  work-instance and never pre-empts the deferred triage-zone
  open question; any prose implying a stub is runtime state is a
  defect (D5).
- **Two-register discipline.** Artifact requirements stay
  `must`; the seeding act stays best-effort/observable. Prose
  that promotes the seeding act to a guarantee, or that softens
  an artifact requirement to "should", breaks the register the
  project tenet requires (D6).
- **Exact-match tokens stay verbatim.** `In draft`, `Proposed`,
  `Milestone Contracts`, `Task Contracts`, and every cited
  section title are copied from their canonical definitions, not
  paraphrased.

## Files to touch

> Estimate of the expected change shape, not a binding contract.
> Deviations are handled per the Estimate Deviations callout in
> the implementing PR body.

**Modify:**

- [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md)
  — add the C3 seeding step to the parent-doc promotion gate;
  add the C1 stub-content + A4 table requirements and the A6
  exemption to "Parent-doc child contracts"; narrow "Slug
  generation" per C2.
- [`../../../spec/planning/epic.md`](../../../spec/planning/epic.md)
  and
  [`../../../spec/planning/milestone.md`](../../../spec/planning/milestone.md)
  — reference the new obligation from their Milestone/Task
  Contracts entries (cite-by-name, no restatement; C4).
- [`../../backlog.md`](../../backlog.md) — flip the entry
  to `Graduated` with the `**Plan:**` line (Backlog Impact).

**Not touched (estimate):**

- Product code, `internal/site/*`, schema, API, lifecycle
  tokens — a slug+Status stub is an already-supported render
  case (scoping D4); spec-prose-only.
- [`../../../spec/backlog.md`](../../../spec/backlog.md) and
  `planning-doc-location.md` — seeding reuses the existing
  layout convention and backlog lifecycle unchanged.

## Validation Gate

Prose-spec change; validation is read-based. Before the
implementing PR opens:

1. **Artifact-checkability read.** Each A1–A6 fact is stated in
   the spec as something checkable against a promoting PR's
   diff. Falsifier: an artifact requirement phrased as a
   best-effort aspiration.
2. **Register-separation check.** The seeding act (C3) is
   framed best-effort/observable; A1–A6 stay `must`. Falsifier:
   the seeding act stated as a guarantee, or an artifact
   requirement softened to "should/recommended".
3. **Single-allocator check.** After the C2 edit, reading "Slug
   generation" end to end yields exactly one allocation event
   for a seeded child (author-supplied at promotion) with no
   residual path implying a divergent server counter. Falsifier:
   prose a reader can take to mean a seeded slug is later
   server-reallocated.
4. **No-duplication / layering check.** Grep the diff: no rule
   body restated in the gate or per-level files; references by
   section title only. Falsifier: a hunk copying rule-body text.
5. **Triage-zone boundary check.** The spec states a stub
   creates no work-instance and does not pre-empt the deferred
   open question. Falsifier: the boundary absent or implying
   resolution of triage zone.
6. **Exact-match token check.** Grep the diff for every cited
   token/section title; confirm verbatim against canonical
   definitions.
7. **Link-resolution check.** Every relative link added or moved
   resolves from its editing file's location.

The implementing PR body carries a `## Review Stance` section
(spec-doc PRs default to the canonical stance) and an
`## Estimate Deviations` section per the Plan-to-PR Completion
Gate.

## Self-Review Audits

Drawn from
[`../../agents/local/self-review-catalog.md`](../../agents/local/self-review-catalog.md);
diff surface is documentation/spec:

- **readiness-gate-truthfulness** — this edits a readiness gate;
  audit that the added step actually causes seeding to be
  prompted at the flip and is framed as best-effort, not
  announcing an enforcement it doesn't have.
- **validation-honesty** — the Validation Gate is read-based;
  audit each step is performed end-to-end against the final
  diff, not asserted from the plan.
- **trigger-map-currency** — the change adds cross-references
  from the gate and per-level files to "Parent-doc child
  contracts" and "Slug generation"; audit every added pointer
  resolves to the current section title and none drifts.

## Out of Scope

- **Resolving the "Triage zone in 1.0?" open question.** A stub
  creates no work-instance; orphan/unattached work-instances
  stay the epic's deferred question. Deliberated as an
  intersection, not resolved (scoping D5).
- **Tree-rendering changes.** A slug+Status stub already renders
  (scoping D4); no `internal/site` change.
- **Server-side creation at file-write.** A seeded stub
  declares a slug and nothing else; building a runtime
  allocation/registration call into the gate is rejected
  (scoping D1 alt-a). Server slug-generation is untouched as the
  no-pre-declared-slug fallback.
- **New lifecycle states or a bespoke "stub" Status.** Stubs
  reuse `In draft` (scoping D2).
- **Retrofitting existing parent docs with seeded children or
  the new table shape.** The obligation binds promotions from
  this rule forward; no historical-doc backfill.

## Backlog Impact

Graduated from the
[`stub-children-on-parent-promotion`](../../backlog.md#stub-children-on-parent-promotion)
entry. That entry's Status is set to
`Graduated — stub-children-on-parent-promotion` with the
optional `**Plan:**` line pointing at this doc, in the planning
change that creates this plan — this standalone task's only
tracking surface (no parent epic/milestone row). No other
backlog entry graduates, deletes, splits, or shifts; the
`deterministic-interactive-registration` and triage-zone
concerns are *referenced* as deliberated intersections, not
shifted. Verified by:
[`../../backlog.md`](../../backlog.md);
[`../../../spec/backlog.md:43-72`](../../../spec/backlog.md)
(graduation + `**Plan:**` line lifecycle).

## Related Docs

- [`scoping/README.md`](scoping/README.md) — decisions D1–D6,
  rejected alternatives, and open input I1; transient, deletes
  at this task's terminal PR.
- [`../promotion-gate-explicit-checklist/README.md`](../promotion-gate-explicit-checklist/README.md)
  — the spec-prose-only gate-edit precedent this plan mirrors in
  shape and cite-by-name discipline.
- [`../workstream-tracker-1-0/m1-t1-multi-wi-per-slug.md`](../workstream-tracker-1-0/m1-t1-multi-wi-per-slug.md)
  — the landed exact-slug create-or-attach path C2 depends on,
  and whose review captured this entry.
- [`../workstream-tracker-1-0/README.md`](../workstream-tracker-1-0/README.md)
  — the epic carrying the deferred "Triage zone in 1.0?" open
  question this task intersects but does not resolve.
- [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md),
  [`../../../spec/planning/epic.md`](../../../spec/planning/epic.md),
  [`../../../spec/planning/milestone.md`](../../../spec/planning/milestone.md)
  — the files this task contracts.
