# Scoping — promotion-gate-explicit-checklist

Transient scoping deliberation for the
[`promotion-gate-explicit-checklist`](../../../backlog.md#promotion-gate-explicit-checklist)
backlog entry. Deletes in batch at this standalone task plan's
terminal PR per
[`task-plan.md`](../../../../spec/planning/task-plan.md) "Scoping
owns / plan owns"; survives in git history. Carries no Status
field (scoping docs are transient deliberation; the sibling plan
alone carries Status).

## Context

The `` `In draft` → `Proposed` `` promotion gate for task and
phase plans
([`task-plan.md`](../../../../spec/planning/task-plan.md)
"`In draft` → `Proposed` promotion gate") enumerates four named
steps: end-to-end coherence read, decision-completeness on
Contracts, the universal `Verified by:` walk, and reality-check
re-confirmation. A promoting agent that runs exactly those four
steps can still flip Status `Proposed` without ever explicitly
confirming three obligations that bind the doc anyway via
separate always-on rules: required-sections presence (+ its
variance-disclosure companion), the no-implementation-prescription
rule, and conformance to the broader cross-level guiding spec.
Those rules are always-on, so the gate's silence on them is not a
correctness hole in the rule set — it is a self-containment hole
in the gate's *checklist*: the runner has to remember the
adjacent rules from outside the gate. The gap was hit live during
the `workstream-tracker-1-0-m1-t1` promotion walk and required
ad-hoc augmentation (recorded at that doc's Backlog Impact
capture).

## Decisions made at scoping time

### D1 — Fix both gates, keep them symmetric (not task-plan.md alone)

The backlog entry names only
[`task-plan.md`](../../../../spec/planning/task-plan.md)'s gate.
But [`shared.md`](../../../../spec/planning/shared.md) carries a
symmetric "Parent-doc `In draft` → `Proposed` promotion gate" for
epic and milestone docs with the same four-bullet shape and the
identical self-containment gap. The prior spec task deliberately
built that parent-doc gate symmetric to the task/phase gate (its
C4), and `shared.md` explicitly frames a parent doc flipped
without the walk as "the same drift shape as a task plan flipped
to `Proposed` without its promotion-gate walk." Fixing only the
task/phase gate would re-open exactly the asymmetry C4 closed.

**Decision:** the task graduates as a fix to *both* gates,
keeping them symmetric. **Rejected alternative:** task-plan.md
only (matches the backlog entry's literal text but re-creates a
gate asymmetry the prior task closed, and would predictably spawn
a follow-up entry the first time the gap is hit at parent level).
Verified by:
[`../../../../spec/planning/shared.md:257-310`](../../../../spec/planning/shared.md);
[`../../../../spec/planning/task-plan.md:441-493`](../../../../spec/planning/task-plan.md);
[`../../spec-updates-contracts-and-gates/README.md`](../../spec-updates-contracts-and-gates/README.md)
(C4 established the parent-doc gate's symmetry with the
task/phase gate).

### D2 — Reference the always-on rules by exact section name; do not restate their bodies

The backlog entry's one illustrative option is "extend the gate's
step list to reference those rules by name (a superset
checklist)." Restating each always-on rule's content inside the
gate would violate the layered-authority / no-duplication
discipline `shared.md` is built on ("if a rule binds two or more
planning levels it lives here once; per-level files reference
it") and would create a drift surface every time the referenced
rule changes.

**Decision:** add a final-confirmation step to each gate's step
list that names the always-on rules by their exact section
titles and points the runner at them, restating none of their
content — the same cite-by-name shape the gate's existing
`Verified by:` step already uses. **Rejected alternative:** copy
the rules' substantive prose into the gate (self-contained in the
strongest sense, but duplicative, drift-prone, and a direct
violation of the no-duplication invariant). Verified by:
[`../../../../spec/planning/shared.md:14-17`](../../../../spec/planning/shared.md)
(layered-authority / no-duplication discipline);
[`../../../../spec/planning/task-plan.md:477-481`](../../../../spec/planning/task-plan.md)
(the existing `Verified by:` step is itself a cite-by-name
reference to an always-on rule applied at the flip — the
precedent shape).

### D3 — Name two concrete rules explicitly; cover "broader spec conformance" as a class, not an enumeration

The backlog entry names three checks. Two map to concrete,
nameable rules: required-sections presence →
"Required and optional sections" (per-level file) plus its
companion "Section variance disclosure"
([`shared.md`](../../../../spec/planning/shared.md)); and
no-implementation-prescription → "Plans describe contracts, not
implementation"
([`shared.md`](../../../../spec/planning/shared.md)). The third —
"conformance to the broader guiding specs" — has no single
section; enumerating every cross-level rule into the gate would
re-create the duplication trap D2 rejects and rot as `shared.md`
grows.

**Decision:** the added step names the two concrete rules by
exact title and then carries one catch-all clause directing the
runner to confirm conformance to the always-on cross-level rules
in [`shared.md`](../../../../spec/planning/shared.md) (and the
applicable per-level file) as a class — pointing at `shared.md`'s
own rule index rather than re-listing it. **Rejected
alternative:** enumerate each always-on `shared.md` rule by name
in the gate (exhaustive but a maintenance/duplication burden that
silently goes stale). Verified by:
[`../../../../spec/planning/task-plan.md:169-214`](../../../../spec/planning/task-plan.md)
("Required and optional sections");
[`../../../../spec/planning/shared.md:557-575`](../../../../spec/planning/shared.md)
("Section variance disclosure");
[`../../../../spec/planning/shared.md:312-426`](../../../../spec/planning/shared.md)
("Plans describe contracts, not implementation");
[`../../../../spec/planning/shared.md:1-34`](../../../../spec/planning/shared.md)
(the cross-level rule index the catch-all points at).

### D4 — Frame the added step as the final explicit confirmation of always-on rules, not a new mid-draft obligation

The three rules bind every drafting session continuously; the
gate is the *last* checkpoint, not their origin. Phrasing them as
brand-new "gate steps" would wrongly imply they don't bind until
the flip.

**Decision:** the added prose frames the step as the promotion
flip being the moment these always-on rules get their explicit
final confirmation — mirroring verbatim the framing the gate's
existing `Verified by:` step already uses ("the promotion gate is
when it gets applied universally rather than to whichever claims
happened to feel technical during drafting"). The new step adds
no new obligation; it makes an existing always-on obligation's
final checkpoint explicit inside the gate. **Rejected
alternative:** phrase as net-new gate-only steps (misleads a
reader into treating the rules as gate-scoped). Verified by:
[`../../../../spec/planning/task-plan.md:477-481`](../../../../spec/planning/task-plan.md)
(the precedent framing this decision mirrors).

### D5 — No live-doc acceptance-oracle edit; the edited gate prose is self-verifying

The prior spec task had live docs carrying variance disclosures
it removed in the same PR. This task is different: it makes the
gate's checklist self-contained. The motivating doc
(`workstream-tracker-1-0-m1-t1`, currently `Proposed`) recorded
the gap as a backlog *capture*, not a variance disclosure; its
capture sentence stays historically accurate after the entry
graduates (a graduated entry keeps its slug + anchor). Nothing in
a live doc needs editing for the fix to be verifiable —
acceptance is read off the edited gate prose itself.

**Decision:** no edit to `m1-t1-multi-wi-per-slug.md` or any
other live plan doc. **Rejected alternative:** retrofit the
motivating doc (no disclosure to remove there; editing it would
be scope creep). Verified by:
[`../../workstream-tracker-1-0/m1-t1-multi-wi-per-slug.md:352-365`](../../workstream-tracker-1-0/m1-t1-multi-wi-per-slug.md)
(the capture note that stays accurate post-graduation).

## Plan structure handoff

- **Doc-type / shape:** standalone task plan, N = 1 phase
  (single `README.md`, phase content inline). Spec-prose-only
  change; no code, schema, or API.
- **Files the plan will contract:** two —
  [`task-plan.md`](../../../../spec/planning/task-plan.md) (the
  task/phase gate) and
  [`shared.md`](../../../../spec/planning/shared.md) (the
  symmetric parent-doc gate), per D1. The plan owns the exact
  contract shape, file inventory, and validation gate.
- **PR shape:** single PR; both edits are one coherent
  self-contained-checklist change.

## Reality-check inputs (re-confirm at promotion)

Load-bearing claims this scoping rests on; the plan's promotion
walk re-confirms each against current code:

- [`task-plan.md`](../../../../spec/planning/task-plan.md) gate
  at lines 441-493 currently enumerates exactly four steps and
  its `Verified by:` step (≈477-481) is already a cite-by-name
  reference to an always-on rule applied at the flip.
- [`shared.md`](../../../../spec/planning/shared.md) parent-doc
  gate at lines 257-310 carries the symmetric four-bullet shape
  with the same self-containment gap.
- "Plans describe contracts, not implementation"
  ([`shared.md`](../../../../spec/planning/shared.md) ≈312-426),
  "Section variance disclosure" (≈557-575), and "Required and
  optional sections"
  ([`task-plan.md`](../../../../spec/planning/task-plan.md)
  ≈169-214 / the per-level epic & milestone files) exist with
  those exact section titles.
- The motivating capture at
  [`m1-t1-multi-wi-per-slug.md:352-365`](../../workstream-tracker-1-0/m1-t1-multi-wi-per-slug.md)
  is a backlog capture (not a variance disclosure) and stays
  accurate after this entry graduates.
