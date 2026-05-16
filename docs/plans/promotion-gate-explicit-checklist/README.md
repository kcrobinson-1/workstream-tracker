---
slug: promotion-gate-explicit-checklist
Status: Landed
---

# Promotion-Gate Explicit Checklist

## Context

The planning spec's `` `In draft` → `Proposed` `` promotion gate
is the comprehensive self-review an agent runs before claiming a
plan doc is ready for code review. It enumerates four named
steps, but three obligations that bind the same doc through
separate always-on rules — that the doc carries its required
sections (with any divergence disclosed), that it prescribes
contracts rather than implementation, and that it conforms to the
broader cross-level guiding spec — are not named in the gate's
own step list. Because those rules are always-on, an agent
running exactly the four listed steps can flip Status to
`Proposed` without ever explicitly confirming them; the gate is
correct as a *rule* but not self-contained as a *checklist*, so
the runner has to remember adjacent rules from outside it.

This is being done now because the gap was hit live: it surfaced
during the `workstream-tracker-1-0-m1-t1` promotion walk, which
required ad-hoc augmentation to cover the unlisted checks, and
that experience was captured as this backlog entry rather than
papered over. Left unfixed, every future gate run re-derives the
same augmentation from memory or skips it silently.

The surface is entirely the planning spec — the rule set this
repo dogfoods and consumer projects vendor. Two files: the
task/phase gate in
[`task-plan.md`](../../../spec/planning/task-plan.md) and its
symmetric parent-doc twin in
[`shared.md`](../../../spec/planning/shared.md). No product code,
schema, API, or lifecycle token changes.

## Goal

Make both `` `In draft` → `Proposed` `` promotion gates
self-contained as checklists: an agent following only a gate's
listed steps also explicitly confirms required-sections presence
(with disclosed variance), no-implementation-prescription, and
conformance to the broader cross-level spec — without the gate
restating any of those always-on rules' bodies, and without
implying the rules bind only at the flip.

Verifiable when: reading either gate's step list end to end
covers all three obligations by name; the two gates remain
symmetric; no always-on rule's substantive prose is duplicated
into either gate; and the added prose frames the rules as
receiving their final explicit confirmation at the flip, not as
new mid-draft obligations.

## Contracts

The final shape each spec file must satisfy — WHAT the spec must
establish, not the prose that establishes it (exact wording is
the implementing PR's choice).

### C1 — Task/phase gate self-containment (`task-plan.md`)

The `` `In draft` → `Proposed` promotion gate `` step list in
[`task-plan.md`](../../../spec/planning/task-plan.md) gains a
final step that, citing by exact section title and restating no
rule body, directs the promoting agent to confirm at the flip:

- **Required-sections presence + disclosed variance** — against
  [`task-plan.md`](../../../spec/planning/task-plan.md)
  "Required and optional sections" and its companion
  [`shared.md`](../../../spec/planning/shared.md) "Section
  variance disclosure" (a missing required section, or an
  undisclosed divergence, is a gate failure).
- **No implementation prescription** — against
  [`shared.md`](../../../spec/planning/shared.md) "Plans
  describe contracts, not implementation."
- **Broader cross-level conformance** — a catch-all directing
  confirmation that the doc conforms to the always-on
  cross-level rules in
  [`shared.md`](../../../spec/planning/shared.md) (and the
  applicable per-level file) as a class, pointing at
  `shared.md`'s own rule index rather than re-listing it.

The step is framed as the promotion flip being the moment these
always-on rules get their explicit final confirmation —
mirroring the framing the gate's existing `Verified by:` step
already carries ("the promotion gate is when it gets applied
universally rather than to whichever claims happened to feel
technical during drafting"). It adds no new mid-draft
obligation. Verified by:
[`../../../spec/planning/task-plan.md:441-493`](../../../spec/planning/task-plan.md)
(current four-step gate);
[`../../../spec/planning/task-plan.md:169-214`](../../../spec/planning/task-plan.md)
("Required and optional sections");
[`../../../spec/planning/shared.md:312-426`](../../../spec/planning/shared.md)
("Plans describe contracts, not implementation");
[`../../../spec/planning/shared.md:557-575`](../../../spec/planning/shared.md)
("Section variance disclosure");
[`../../../spec/planning/shared.md:1-34`](../../../spec/planning/shared.md)
(the cross-level rule index the catch-all points at).

### C2 — Symmetric parent-doc gate self-containment (`shared.md`)

The "Parent-doc `` `In draft` → `Proposed` `` promotion gate"
step list in [`shared.md`](../../../spec/planning/shared.md)
gains the equivalent final step, in the same cite-by-name shape,
same always-on framing, restating no rule body. Its
required-sections reference resolves to the parent-doc per-level
realization — "Required and optional sections" in
[`epic.md`](../../../spec/planning/epic.md) /
[`milestone.md`](../../../spec/planning/milestone.md) — while the
"Section variance disclosure" and "Plans describe contracts, not
implementation" references and the broader-conformance catch-all
are the same cross-level rules C1 cites. After the edit the two
gates carry equivalent self-containment steps (allowing only for
the per-level realization difference), preserving the symmetry
the prior spec task established. Verified by:
[`../../../spec/planning/shared.md:257-310`](../../../spec/planning/shared.md)
(current parent-doc gate);
[`../spec-updates-contracts-and-gates/README.md`](../spec-updates-contracts-and-gates/README.md)
(C4 — the parent-doc gate built symmetric to the task/phase
gate);
[`../../../spec/planning/epic.md`](../../../spec/planning/epic.md)
and
[`../../../spec/planning/milestone.md`](../../../spec/planning/milestone.md)
(the per-level "Required and optional sections" lists the
parent-doc step's required-sections reference resolves to).

### C3 — No rule-body duplication; layered authority preserved

Neither added step copies an always-on rule's substantive prose
into the gate; each references the rule by its exact section
title only. A reviewer finding an always-on rule's body restated
inside either gate treats that as a defect, not
self-containment-by-copying — the gate is made self-contained by
*pointing*, per the layered-authority / no-duplication discipline
`shared.md` is built on. Verified by:
[`../../../spec/planning/shared.md:14-17`](../../../spec/planning/shared.md)
(layered-authority / no-duplication discipline);
[`../../../spec/planning/task-plan.md:477-481`](../../../spec/planning/task-plan.md)
(the existing `Verified by:` step — the cite-by-name precedent
shape both added steps follow).

## Cross-Cutting Invariants

- **Two-gate symmetry.** The task/phase gate
  ([`task-plan.md`](../../../spec/planning/task-plan.md)) and the
  parent-doc gate
  ([`shared.md`](../../../spec/planning/shared.md)) must carry
  equivalent self-containment steps; editing one without the
  other re-opens the asymmetry the prior spec task's C4 closed.
  The only permitted divergence is the required-sections
  reference resolving to a different per-level home (task-plan vs.
  epic/milestone).
- **Layered authority / no per-gate duplication.** The added
  steps cite always-on rules by exact section title and never
  restate their bodies; a rule body appearing inside a gate is a
  defect, not redundancy-for-safety.
- **Exact-match tokens stay verbatim.** `In draft`, `Proposed`,
  and every cited section title are copied from their canonical
  definitions, not paraphrased — the gate is matched and
  navigated by those exact strings.

## Files to touch

> Estimate of the expected change shape, not a binding contract.
> Deviations are handled per the Estimate Deviations callout in
> the implementing PR body.

**Modify:**

- [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md)
  — add the C1 self-containment step to the
  `` `In draft` → `Proposed` `` promotion gate's step list.
- [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md)
  — add the symmetric C2 step to the "Parent-doc
  `` `In draft` → `Proposed` `` promotion gate" step list.

**Not touched (estimate):**

- [`../../../spec/planning/epic.md`](../../../spec/planning/epic.md)
  and
  [`../../../spec/planning/milestone.md`](../../../spec/planning/milestone.md)
  — their "Required and optional sections" lists are *referenced*
  by the C2 step, not edited; the lists already exist.
- [`../../../spec/backlog.md`](../../../spec/backlog.md) and
  `planning-doc-location.md` — gate self-containment does not
  touch backlog format or layout convention.
- Product code, schema, API, lifecycle tokens — spec-prose-only.

## Validation Gate

Prose-spec change; validation is read-based, not build-based.
Before the implementing PR opens:

1. **Self-containment read.** Read each edited gate's step list
   end to end. Confirm a runner following only the listed steps
   now also covers required-sections presence + disclosed
   variance, no-implementation-prescription, and broader
   cross-level conformance. Falsifier: any of the three is still
   absent from a gate's step list after the edit.
2. **Symmetry check.** Confirm the task/phase gate and the
   parent-doc gate carry equivalent self-containment steps,
   diverging only in the per-level required-sections reference.
   Falsifier: one gate gains the step, the other does not.
3. **No-duplication check.** Grep the diff: no always-on rule's
   body prose is restated inside either gate; references are by
   section title only. Falsifier: a diff hunk copying rule-body
   text into a gate.
4. **Always-on framing check.** Confirm the added prose frames
   the rules as receiving their final explicit confirmation at
   the flip, not as obligations that begin at the flip.
   Falsifier: prose a reader can take to mean the rules don't
   bind before `Proposed`.
5. **Exact-match token check.** Grep the diff for `In draft`,
   `Proposed`, and every cited section title; confirm each is
   verbatim against its canonical definition, not paraphrased.
6. **Link-resolution check.** Confirm every relative link added
   or moved resolves to an existing target from the editing
   file's location.

The implementing PR body carries a `## Review Stance` section
(plan/spec-doc PRs default to the canonical stance) and an
`## Estimate Deviations` section per the Plan-to-PR Completion
Gate.

## Self-Review Audits

Drawn from
[`../../agents/local/self-review-catalog.md`](../../agents/local/self-review-catalog.md);
diff surface is documentation/spec:

- **readiness-gate-truthfulness** — this task edits a readiness
  gate; audit that the gate's amended step list actually causes
  the named confirmations to happen rather than announcing
  coverage it doesn't enforce.
- **validation-honesty** — the Validation Gate above is
  read-based; audit that each numbered step's "confirm X" is
  performed end-to-end against the final diff, not asserted from
  the plan.
- **trigger-map-currency** — the change adds cross-references
  from the gates to always-on rules; audit that every added "see
  X" pointer resolves to the rule's current section title and no
  pointer drifts.

## Out of Scope

- **Enumerating every cross-level rule into the gate.** The
  broader-conformance clause is a catch-all pointing at
  `shared.md`'s rule index; exhaustively listing each always-on
  rule in the gate is a maintenance/duplication trap and is not
  done (scoping D3).
- **Restating rule bodies for stronger self-containment.**
  Self-containment is achieved by cite-by-name, not by copying
  (scoping D2; C3).
- **Retrofitting historical or live plan docs.** No live doc
  carries a variance disclosure this task removes; the
  motivating `workstream-tracker-1-0-m1-t1` capture stays
  accurate post-graduation and is not edited (scoping D5).
- **Editing the epic/milestone "Required and optional sections"
  lists.** They are referenced by the C2 step, not changed.
- **New lifecycle states or tooling.** The gates reuse the
  existing `In draft` / `Proposed` tokens; no Status value or
  visualization change.

## Backlog Impact

Graduated from the
[`promotion-gate-explicit-checklist`](../../backlog.md#promotion-gate-explicit-checklist)
entry. That entry's Status is set to
`Graduated — promotion-gate-explicit-checklist` with the
optional `**Plan:**` path line pointing at this doc, in the
planning change that creates this plan (per the parent-doc
currency obligation — a graduated backlog entry is this
standalone task's only tracking surface; there is no parent
epic/milestone row). No other backlog entry graduates, deletes,
splits, or shifts. Verified by:
[`../../backlog.md`](../../backlog.md);
[`../../../spec/backlog.md:43-72`](../../../spec/backlog.md)
(graduation + `**Plan:**` line lifecycle).

## Related Docs

- The transient scoping deliberation (decisions D1–D5 with
  rejected alternatives) was deleted at this task's terminal PR
  per [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md)
  "Scoping owns / plan owns"; it survives in git history.
- [`../spec-updates-contracts-and-gates/README.md`](../spec-updates-contracts-and-gates/README.md)
  — the prior spec task whose C4 built the parent-doc gate
  symmetric to the task/phase gate; this task preserves that
  symmetry.
- [`../workstream-tracker-1-0/m1-t1-multi-wi-per-slug.md`](../workstream-tracker-1-0/m1-t1-multi-wi-per-slug.md)
  — the doc whose promotion walk surfaced and captured this gap.
- [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md)
  and [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md)
  — the two files this task contracts.
