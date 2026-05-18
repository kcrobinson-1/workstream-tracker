---
slug: spec-updates-contracts-and-gates
Status: Landed
---

# Spec Updates: Contracts and Gates

## Context

This plan corrects a cluster of internal inconsistencies in the
planning and backlog spec — the rule set consumer projects vendor
and the workstream-tracker repo dogfoods on itself. The gaps are
not new features: each is a place where the spec contradicts
itself or under-specifies a shape its own live plan docs already
need, forcing a disclosed section-variance work-around instead of
a rule the doc can follow.

It is being done now because the inconsistencies were hit
directly while drafting the `workstream-tracker-1-0` epic and its
m1 milestone, and the work-arounds those docs carry (variance-
disclosed "Milestone Contracts" / "Task Contracts" sections, an
epic doc stuck at `In draft` with no rule telling it to advance)
will keep accreting one disclosed variance at a time until the
spec is corrected. This is a **chosen** proactive cross-level
consistency fix, not a fix the section-variance rule yet mandates:
that rule fires only when the same variance recurs across two or
more docs of the *same* doc-type, and here the variance appears
once at epic level and once at milestone level — one instance per
type, below the per-type threshold. The same conceptual gap
manifesting at every parent level is the reason to codify the
cross-level rule now rather than wait for each per-type variance
to recur twice and trigger four separate codification PRs.

The surfaces this touches are entirely documentation: the
cross-level planning rule file, the three per-level planning rule
files (epic, milestone, task-plan), and the backlog format spec.
No product code, no schema, no API. The scoping deliberation
(six decisions D1–D6 with rejected alternatives and reality-check
inputs) was transient and was deleted at this task's terminal PR
per [`task-plan.md`](../../../spec/planning/task-plan.md) "Scoping
owns / plan owns"; it survives in git history. Everything durable
from it is carried in this plan's Contracts and Out of Scope.

## Goal

Correct the planning and backlog spec so that:

1. Parent-level plan docs (epic, milestone) and N≥2 task plans
   carry a required per-child **WHAT-contract** section, and the
   anti-scope rules permit that WHAT contract while still
   forbidding HOW-scoping.
2. Epic and milestone docs are bound by an `In draft` →
   `Proposed` promotion gate analogous to the one task/phase
   plans already have, so a PR that locks a parent doc's scope
   has a rule telling it to advance Status.
3. `backlog.md` codifies a navigable affordance to a graduated
   entry's plan beyond the slug-only `Graduated — <plan-slug>`
   form, **and** the motivating local entry adopts it in the
   same PR so the dogfood is closed rather than left slug-only.
4. The milestone spec names the milestone's direct child unit
   consistently with the slug grammar and level picker (task,
   not phase).

Verifiable when the five spec files are internally consistent on
these four points, the live `workstream-tracker-1-0` docs no
longer need section-variance disclosures for their child-contract
and Out-of-Scope sections, and the `plan-doc-child-contracts`
backlog entry carries the new graduated-plan affordance.

## Contracts

The final shape each spec file must satisfy. These are WHAT the
spec must establish, not the prose that establishes it — exact
wording is the implementing PR's choice.

### C1 — Cross-level child-contract rule (`shared.md`)

`shared.md` carries one cross-level rule establishing that a
parent doc states, per direct child, the child's **WHAT**
contract (end result, sibling interfaces, what it preserves) and
that the child's **HOW** (file inventory, signatures, commands,
validation gate) stays at the child's own doc. The rule names its
per-level section realization (Milestone Contracts at epic level,
Task Contracts at milestone level, Phase Contracts at task level
when N≥2) and is placed in `shared.md` because it binds three
doc-types — consistent with `shared.md`'s rule that a rule
binding two or more planning levels lives there once. Verified by:
[`../../../spec/planning/shared.md:14-17`](../../../spec/planning/shared.md).

### C2 — Per-level required-section additions

- `epic.md`'s "Required and optional sections" list includes a
  **Milestone Contracts** entry (required when the epic locks any
  milestone's scope).
- `milestone.md`'s list includes a **Task Contracts** entry
  (required when the milestone locks any task's scope).
- `task-plan.md`'s list includes a **Phase Contracts** entry
  (required when the task plan is N≥2).

Each addition cross-references C1 as its authority rather than
restating the WHAT/HOW split. This codification is a chosen
cross-level consistency fix, not one the recurrence-to-codify
rule mandates: that rule's threshold is two or more docs of the
*same* doc-type, and the variance currently appears once per type
(one epic doc, one milestone doc). Codifying proactively avoids
waiting for each per-type variance to recur. Verified by:
[`../../../spec/planning/shared.md:443-446`](../../../spec/planning/shared.md);
[`../workstream-tracker-1-0/README.md:120-128`](../workstream-tracker-1-0/README.md);
[`../workstream-tracker-1-0/m1-v0-2.md:77-85`](../workstream-tracker-1-0/m1/README.md).

### C3 — Anti-scope rules permit WHAT, forbid HOW

`milestone.md`'s "Anti-goal: do not scope any phase" and
`epic.md`'s "Scope: what an epic does and does not say" each read
such that the per-child WHAT contract (C1/C2) is explicitly
permitted and the prohibition is narrowed to HOW-scoping (child-
level file inventory, contracts-as-implementation, risks,
execution steps). No reader can come away believing the anti-
scope rule bars the child-contract section. Verified by:
[`../../../spec/planning/milestone.md:35-43`](../../../spec/planning/milestone.md);
[`../../../spec/planning/epic.md:8-31`](../../../spec/planning/epic.md).

### C4 — Parent-doc promotion gate (`shared.md`)

`shared.md` carries an `In draft` → `Proposed` promotion gate
that binds epic and milestone docs, mirroring the task/phase gate
in `task-plan.md`. The gate:

- fires on the PR that locks the parent doc's scope (its child
  set and the per-child WHAT-contracts are decision-complete),
  flipping frontmatter Status `In draft` → `Proposed`;
- requires, before the flip, an end-to-end coherence read, a
  decision-completeness check on the child WHAT-contracts, the
  `Verified by:` walk at the level's existing "What this means at
  each level" interpretation, and re-confirmation of any reality-
  check inputs the parent doc rests on;
- lives in `shared.md` (not `task-plan.md`, not duplicated into
  `epic.md`/`milestone.md`) because it binds two doc-types and
  task-plan-layer files are not loaded by epic/milestone docs.

The Status tokens it references (`In draft`, `Proposed`) match
the canonical lifecycle by exact string. Verified by:
[`../../../spec/planning/task-plan.md:8-13`](../../../spec/planning/task-plan.md);
[`../../../spec/planning/task-plan.md:436-482`](../../../spec/planning/task-plan.md);
[`../../../spec/planning/shared.md:135-165`](../../../spec/planning/shared.md).

### C5 — Backlog graduated-link affordance (`backlog.md`)

`backlog.md`'s entry-lifecycle section codifies an optional,
bold-labeled line on a graduated entry carrying a repo-relative
path to the graduated plan doc, distinct from and additional to
the exact-match `Graduated — <plan-slug>` Status line. The slug
on the Status line is not wrapped in link syntax (it stays an
exact-match token). The affordance uses only link/path syntax
that exists today (no dependency on the deferred
`repo-rooted-doc-links` work). The same PR applies the codified
affordance to the motivating
[`plan-doc-child-contracts`](../../backlog.md#plan-doc-child-contracts)
entry in `docs/backlog.md` — adding the path line pointing at
this plan doc — so the entry that drove this work is the first
to carry the new form rather than being left slug-only. The
original scaffold's reason for keeping the entry slug-only ("a
pre-supplied ad-hoc link would bias the design decision") no
longer holds: the design is locked in scoping decision D5, so
applying its codified form is no longer ad-hoc. Verified by:
[`../../../spec/backlog.md:39-42`](../../../spec/backlog.md);
[`../../../spec/backlog.md:44-53`](../../../spec/backlog.md);
[`../../../spec/planning/shared.md:546-561`](../../../spec/planning/shared.md);
[`../../backlog.md:7-21`](../../backlog.md).

### C6 — Milestone child named "task" not "phase" (`milestone.md`)

`milestone.md`'s required-section list and surrounding prose name
the milestone's *direct child* unit consistently with the slug
grammar and level picker: the milestone's child is a **task**
(Task Status, Cross-Task Invariants/Decisions/Risks). "Phase"
remains correct only where the doc refers to the task's child one
level lower. The live milestone doc's "Task …" sections no longer
require a variance disclosure after this lands. Verified by:
[`../../../spec/planning/shared.md:97-106`](../../../spec/planning/shared.md);
[`../../../spec/planning/task-plan.md:22-37`](../../../spec/planning/task-plan.md);
[`../../../spec/planning/milestone.md:99-105`](../../../spec/planning/milestone.md);
[`../workstream-tracker-1-0/m1-v0-2.md:25-31`](../workstream-tracker-1-0/m1/README.md).

## Cross-Cutting Invariants

- **Format/data layer stays strictly additive; authoring
  obligations intentionally tighten.** Two distinct compatibility
  claims that must not be conflated. (1) *Format/data additivity
  — must hold.* No edit changes a frontmatter field, a Status or
  slug token, the slug grammar, or the file-layout convention a
  vendored snapshot's tooling parses; a consumer's existing docs
  still parse and render unchanged. (2) *Authoring-obligation
  tightening — intentional and disclosed.* Adding
  required-when-applicable child-contract sections (C2) and a
  parent-doc promotion gate (C4) makes a parent-doc author do
  more than before. That is a deliberate **contract tightening**,
  not a purely-additive change, and it is explicitly **not**
  covered by — nor is this task governed by — the
  `workstream-tracker-1-0` epic's additive-spec invariant (that
  invariant binds that epic's milestones; this is a separate
  standalone task with its own root slug). The tightening is
  acceptable pre-1.0 because the only consumer today is this repo
  dogfooding, no external project has vendored a snapshot yet
  (neighborly-events integration is a future, not-yet-started
  milestone), and the obligation already existed de facto as the
  disclosed variance the live docs carry — codification removes a
  work-around, it does not invent a new burden. The implementing
  PR body states the tightening explicitly under Estimate
  Deviations / Documentation rather than letting the additivity
  check pass silently. Verified by:
  [`../workstream-tracker-1-0/README.md:64-71`](../workstream-tracker-1-0/README.md)
  (the epic-scoped additive invariant this task is distinct
  from);
  [`../workstream-tracker-1-0/README.md:90-118`](../workstream-tracker-1-0/README.md)
  (neighborly-events integration is a later milestone, so no
  external snapshot is vendored yet).
- **Layered authority, no per-level duplication.** Cross-level
  rules (C1, C4) live in `shared.md` once; per-level files
  reference them, never restate them. A reviewer finding the same
  rule body in `shared.md` and a per-level file should treat that
  as a defect, not redundancy-for-safety. Verified by:
  [`../../../spec/planning/shared.md:1-18`](../../../spec/planning/shared.md).
- **Exact-match Status/slug tokens stay verbatim.** Any Status or
  backlog-slug token the edits reference is copied from its
  canonical definition, never paraphrased; C4 and C5 both touch
  exact-match-tracked tokens. Verified by:
  [`../../../spec/planning/shared.md:546-561`](../../../spec/planning/shared.md).
- **The live `workstream-tracker-1-0` docs are the acceptance
  oracle.** A correct edit is one after which those docs' child-
  contract and Out-of-Scope sections, and the milestone doc's
  Task-named sections, no longer need a variance disclosure. The
  same PR removes the now-unnecessary disclosures from those
  docs.

## Files to touch

> Estimate of the expected change shape, not a binding contract.
> Implementation may revise the inventory when a structural call
> requires it; deviations are handled per the Estimate Deviations
> callout in the implementing PR body.

**Modify:**

- [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md)
  — add C1 (cross-level child-contract rule) and C4 (parent-doc
  promotion gate).
- [`../../../spec/planning/epic.md`](../../../spec/planning/epic.md)
  — add Milestone Contracts to required/optional list (C2);
  refine the scope-anti-goal prose (C3).
- [`../../../spec/planning/milestone.md`](../../../spec/planning/milestone.md)
  — add Task Contracts to required/optional list (C2); refine
  the anti-scope prose (C3); rename child unit "phase" → "task"
  in required sections and prose (C6).
- [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md)
  — add Phase Contracts (N≥2) to required/optional list (C2);
  add a back-reference to the new `shared.md` parent-doc gate so
  the symmetry with the task/phase gate is navigable (C4).
- [`../../../spec/backlog.md`](../../../spec/backlog.md)
  — codify the optional graduated-plan path line (C5).
- [`../../backlog.md`](../../backlog.md)
  — apply the new C5 affordance to the
  `plan-doc-child-contracts` entry: add the path line pointing
  at this plan doc (the dogfood close-out for C5).
- [`../workstream-tracker-1-0/README.md`](../workstream-tracker-1-0/README.md)
  — Milestone Contracts variance disclosure removed (now cites
  the codified `epic.md` + `shared.md` rule). Status left
  `In draft`: the new C4 gate was judged **not** satisfied at
  PR time — the epic's child set is not locked (m2/m3 are
  proposed-not-locked, the final-integration slug is
  unallocated, and Open Questions remain open), so the gate's
  "child set locked + each child WHAT decision-complete"
  precondition fails.
- [`../workstream-tracker-1-0/m1-v0-2.md`](../workstream-tracker-1-0/m1/README.md)
  — Task Contracts variance disclosure removed (now cites the
  codified `milestone.md` + `shared.md` rule); the Task Status /
  Cross-Task section naming now matches `milestone.md` after C6,
  so it is no longer a variance. **Deviation from estimate:**
  the Out-of-Scope variance disclosure is **retained** — this
  task's contracts (C1–C6) never added an `Out Of Scope` section
  to the milestone doc-type, and adding one would be the
  retroactive scope creep this plan's own Out of Scope forbids.
  It stays a legitimate, separately-disclosed variance outside
  this task's surface.

**Not touched (estimate):**

- Product code, schema, API — this is a spec-only correction.
- `planning-doc-location.md` — layout convention is unaffected
  by contract/gate/backlog rule changes.
- The `repo-rooted-doc-links` backlog entry — C5 is deliberately
  built on today's link syntax and does not graduate, shift, or
  depend on it.

## Validation Gate

This is a prose-spec change; validation is read-based, not
build-based. Before the implementing PR opens:

1. **Internal-consistency read.** Read the five edited spec files
   end to end; confirm no per-level file restates a cross-level
   rule body (layered-authority invariant), and that C3's anti-
   scope prose and C1/C2's child-contract requirement do not
   contradict each other.
2. **Acceptance-oracle check.** Re-read
   [`../workstream-tracker-1-0/README.md`](../workstream-tracker-1-0/README.md)
   and
   [`../workstream-tracker-1-0/m1-v0-2.md`](../workstream-tracker-1-0/m1/README.md);
   confirm every section those docs disclosed as a variance is
   now spec-covered, and that the disclosures are removed in the
   same PR. The falsifier: if any of those docs still needs a
   variance disclosure after the edits, C2/C3/C6 are
   under-satisfied.
3. **Exact-match token check.** Grep the diff for every Status
   and backlog-slug token introduced; confirm each is verbatim
   against its canonical definition, not paraphrased.
4. **Format/data additivity check.** Confirm no frontmatter
   field, Status or slug token, slug grammar, or layout
   convention is changed — a vendored consumer's existing docs
   must still parse and render unchanged. The falsifier: a diff
   hunk that alters a token, field name, or grammar a snapshot's
   tooling reads.
5. **Authoring-obligation-tightening acknowledgment.** Confirm
   the implementing PR body explicitly states that C2/C4 tighten
   parent-doc authoring obligations (not a silent
   "only-adds-prose" pass). The falsifier: the diff adds a
   required-when-applicable section or gate but no PR-body line
   names the tightening — exactly the silent-strictness failure
   this gate exists to catch.
6. **Backlog dogfood check.** Confirm
   [`../../backlog.md`](../../backlog.md)'s
   `plan-doc-child-contracts` entry carries the new C5 path line
   in the same PR. The falsifier: `backlog.md` codifies the
   affordance but the motivating entry is still slug-only after
   merge.
7. **Link-resolution check.** Confirm every relative link added
   or moved resolves to an existing target from the editing
   file's location.

The implementing PR body carries a `## Review Stance` section
(plan/spec-doc PRs default to the canonical stance) and an
`## Estimate Deviations` section per the Plan-to-PR Completion
Gate.

## Self-Review Audits

Drawn from
[`../../agents/local/self-review-catalog.md`](../../agents/local/self-review-catalog.md).
Surface for this diff is documentation/spec:

- **readiness-gate-truthfulness** — C4 introduces a readiness
  gate; audit that the gate's prose actually verifies the named
  conditions (coherence read, decision-completeness, `Verified
  by:` walk, reality-check re-confirm) rather than announcing
  readiness on best-effort.
- **validation-honesty** — the Validation Gate above is read-
  based; audit that each numbered step's claim ("confirm X") is
  one the reader actually performs end-to-end on the final diff,
  not asserted from the plan.
- **trigger-map-currency** — C1 and C4 relocate where governance
  rules live (cross-level vs. per-level); audit that every per-
  level file's cross-reference points at the new `shared.md`
  home and no stale "see X" pointer drifts.

## Documentation currency

The spec files *are* the documentation under change; no separate
status-bearing doc tracks them. The one currency obligation is
the acceptance-oracle cleanup already in Files to touch: the
`workstream-tracker-1-0` epic and m1 docs lose their now-stale
variance disclosures in the same PR, and the epic doc's Status is
re-evaluated against the new C4 gate at PR time.

## Out of Scope

- **External prior art.** neighborly-events' backlog `Detail`
  format (named in the originating scaffold as "prior art to
  consult, not match") is in an external project not vendored
  here. C5 is decided from in-repo constraints; matching an
  unseen format is explicitly not attempted. This is a known,
  accepted gap, not a blocker.
- **Repo-rooted link syntax.** The `repo-rooted-doc-links`
  backlog entry stays Open / post-1.0. C5's affordance uses
  today's relative-path syntax and does not pre-empt that
  decision.
- **Retroactive sweep of all historical plan docs.** Only the
  two live `workstream-tracker-1-0` docs that carry the recurring
  variance are reconciled; demo-workstream and other docs are not
  retrofitted by this pass.
- **New plan-doc lifecycle states or tooling.** C4 reuses the
  existing `In draft` / `Proposed` tokens; no new Status value,
  no visualization change.

## Backlog Impact

Graduated from the
[`plan-doc-child-contracts`](../../backlog.md#plan-doc-child-contracts)
entry. That entry's Status is already
`Graduated — spec-updates-contracts-and-gates`; the Phase-vs-Task
naming note appended to it (its "Also surfaced" paragraph) is
folded into this task as C6/D6 rather than split into a separate
entry, per the entry's own fold-or-split offer. This PR also
edits that entry's body to add the new C5 graduated-plan path
line (the dogfood close-out): a content edit to an
already-graduated entry, not a new graduate/delete/split/shift
of its lifecycle state. No other backlog entry graduates,
deletes, splits, or shifts. The
[`repo-rooted-doc-links`](../../backlog.md#repo-rooted-doc-links)
entry is deliberately untouched (see Out of Scope). Verified by:
[`../../backlog.md:23-59`](../../backlog.md).

## Related Docs

- [`../../backlog.md`](../../backlog.md) — origin
  `plan-doc-child-contracts` entry.
- [`../workstream-tracker-1-0/README.md`](../workstream-tracker-1-0/README.md)
  and
  [`../workstream-tracker-1-0/m1-v0-2.md`](../workstream-tracker-1-0/m1/README.md)
  — the live docs whose disclosed variances are this task's
  acceptance oracle.
- [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md)
  — the rules this task plan is structured against.
