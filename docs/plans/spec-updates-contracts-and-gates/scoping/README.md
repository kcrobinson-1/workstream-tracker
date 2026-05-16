---
slug: spec-updates-contracts-and-gates
---

# Scoping — Spec Updates: Contracts and Gates

Transient scoping artifact for the
[`spec-updates-contracts-and-gates`](../README.md) task plan.
Deletes at the task-terminal PR per
[`task-plan.md`](../../../../spec/planning/task-plan.md) "Scoping
owns / plan owns." Carries no `Status` field by rule: the Status
lifecycle is for the durable plan, not the transient scoping doc.

## Context

Dogfooding workstream-tracker's own planning surfaced a cluster
of internal inconsistencies in the planning + backlog spec — the
artifacts consumer projects vendor. None block 1.0, but each was
hit while drafting the `workstream-tracker-1-0` epic and its m1
milestone, and each forced a disclosed section-variance
work-around rather than a spec the doc could follow. This task
batches the corrections into one spec-correctness pass. This
scoping doc owns the deliberation: the decisions made, the
alternatives rejected with code citations, the plan-structure
handoff, and the reality-check inputs the plan verifies before
promotion. It does not restate the plan's Contracts, Files-to-
touch inventory, or Validation Gate — those live in the plan and
are referenced here by name.

## Why this plan needs a scoping doc

The "Narrow-surface plans may skip the scoping doc" carve-out in
[`task-plan.md`](../../../../spec/planning/task-plan.md) requires
all five conditions to hold. Condition 4 ("No new cross-cutting
invariant") fails: the cross-level WHAT/HOW-split rule this task
adds is itself a rule multiple spec files must agree on, which is
the load-bearing reason the carve-out names for *not* skipping.
Scoping is therefore mandatory, not optional. The
"Doc-only decision plans" carve-out also does not apply: this
plan ships concrete spec edits, not a recorded-decision artifact,
so the decision-resolved-at-scoping-time bar applies in full.

## Decisions made at scoping time

Each decision carries a `Verified by:` citation to the spec file
that grounds it. Spec files are the codebase for a spec-
correctness task; citing them is the code-grounding the
reality-check gate requires.

### D1 — The cross-level WHAT/HOW-split rule lives in `shared.md`

The rule that parent docs carry per-child WHAT contracts while
HOW stays at the child's own level binds epic, milestone, and
task-plan (N≥2) doc-types simultaneously. `shared.md`'s own
layering discipline mandates the home: a rule binding two or more
planning levels lives in `shared.md` once, not duplicated
per-level.

- **Rejected:** duplicate the rule into `epic.md`,
  `milestone.md`, and `task-plan.md`. Rejected because it
  re-introduces the per-level duplication trap `shared.md`
  explicitly exists to prevent.
- `Verified by:`
  [`spec/planning/shared.md:14-17`](../../../../spec/planning/shared.md)
  ("if a rule binds two or more planning levels, it lives here
  once").

### D2 — Codifying the child-contract sections is a chosen cross-level fix, not a recurrence-mandated one

The child-contract section already exists as a *disclosed
variance* in two live parent docs — but of *different doc-types*
(one epic doc, one milestone doc). The
section-variance-disclosure rule fires only when the same
variance recurs across two or more docs of the **same** type, so
with one instance per type the per-type threshold is **not** met
and the rule does **not** yet mandate codification. The decision
here is therefore a deliberate choice: the same conceptual gap
shows up at every parent level, and codifying the cross-level
rule now is preferable to waiting for each per-type variance to
recur twice and then running separate per-type codification PRs.
Codify proactively, framed as a chosen consistency fix — not as a
mandate the spec does not yet impose.

- Add "Milestone Contracts" to the `epic.md` required/optional
  list; "Task Contracts" to `milestone.md`; "Phase Contracts
  (when N≥2)" to `task-plan.md`.
- `Verified by:`
  [`spec/planning/shared.md:443-446`](../../../../spec/planning/shared.md)
  (recurrence → codify rule);
  [`docs/plans/workstream-tracker-1-0/README.md:120-128`](../../workstream-tracker-1-0/README.md)
  ("Milestone Contracts" disclosed variance, epic doc);
  [`docs/plans/workstream-tracker-1-0/m1-v0-2.md:77-85`](../../workstream-tracker-1-0/m1-v0-2.md)
  ("Task Contracts" disclosed variance, milestone doc).

### D3 — The anti-scope rules forbid HOW-scoping but must require WHAT-contracting

`milestone.md`'s "Anti-goal: do not scope any phase" and
`epic.md`'s "Scope: what an epic does and does not say" bar
HOW-scoping (file inventory, contracts, risks, execution steps at
the child level). Read literally they also bar the WHAT contract
the child-contract section requires, which is the contradiction
the dogfooding hit. The refinement narrows the prohibition to HOW
and explicitly requires the WHAT contract.

- `Verified by:`
  [`spec/planning/milestone.md:35-43`](../../../../spec/planning/milestone.md)
  ("Anti-goal: do not scope any phase");
  [`spec/planning/epic.md:8-31`](../../../../spec/planning/epic.md)
  ("Scope: what an epic does and does not say").

### D4 — The parent-doc promotion gate lives in `shared.md`, scoped to epic + milestone docs

The existing `In draft` → `Proposed` promotion gate binds task
and phase plans only — `task-plan.md` states epic and milestone
docs do not consume it, which is why it lives in `task-plan.md`
rather than `shared.md`. The symmetric parent-doc gate binds epic
+ milestone (two levels) and is the mirror case: it belongs in
`shared.md` per the same two-or-more-levels discipline, scoped to
parent docs. The gate's content adapts to parent docs — decision-
completeness runs against the child WHAT-contracts (D2), internal
coherence and the `Verified by:` walk apply at each level's
existing "What this means at each level" interpretation.

- **Rejected:** place the gate in `task-plan.md` alongside the
  task/phase gate. Rejected because `task-plan.md` is the
  implementation-layer authority file and epic/milestone docs do
  not load it; the gate would never be read by the docs it binds.
- **Rejected:** duplicate the gate into `epic.md` and
  `milestone.md`. Rejected per the same per-level duplication
  trap as D1.
- **Trigger:** the gate binds the PR that locks a parent doc's
  scope (its child set and child WHAT-contracts are decision-
  complete). That PR flips Status `In draft` → `Proposed`.
- `Verified by:`
  [`spec/planning/task-plan.md:8-13`](../../../../spec/planning/task-plan.md)
  (impl-layer gate lives in `task-plan.md` because epic/milestone
  do not consume it — establishes the symmetric placement
  argument);
  [`spec/planning/task-plan.md:436-482`](../../../../spec/planning/task-plan.md)
  (the task/phase gate whose parent-doc mirror this is);
  [`docs/plans/workstream-tracker-1-0/README.md:1-4`](../../workstream-tracker-1-0/README.md)
  (epic doc still `Status: In draft` after its scope-locking PR
  merged — the live evidence the parent-doc gate is missing).

### D5 — Backlog link affordance: an optional bold-labeled path line, not a linked slug

The backlog `Graduated — <plan-slug>` form carries the slug only;
the slug is the navigable identity but resolving it to a file
requires manual path construction. Options decomposed:

- **(a) Slug only (status quo).** Rejected: the navigation gap is
  the entry's reason for graduating.
- **(b) Add an optional bold-labeled line carrying a repo-
  relative path to the graduated plan doc**, distinct from the
  slug line. Chosen: it is additive, fits the existing optional-
  fields affordance, keeps the slug as the untouched exact-match
  token, and needs no link mechanism that does not yet exist.
- **(c) Wrap the `Graduated — <plan-slug>` slug itself in
  markdown-link syntax.** Rejected: the slug is matched by
  exact string for status tracking; wrapping it breaks the
  queryable token, the exact failure mode the exact-match label
  discipline names.
- **(d) Wait for `repo-rooted-doc-links`.** Rejected: that entry
  is Open / post-1.0; a spec-correctness pass must not block on
  deferred post-1.0 work, and the affordance must be expressible
  with today's link syntax.

Once (b) is the locked form, the originating scaffold's reason
for keeping the motivating `plan-doc-child-contracts` entry
slug-only — "pre-supplying an ad-hoc link would bias the design
decision" — no longer holds: applying the *codified* form is not
ad-hoc. So the same PR applies (b) to that entry, closing the
dogfood rather than shipping a backlog rule the motivating entry
itself does not follow. Carried into the plan as C5's
backlog-application clause and a Files-to-touch / Validation-Gate
entry for `docs/backlog.md`.

- `Verified by:`
  [`spec/backlog.md:39-42`](../../../../spec/backlog.md)
  (optional bold-labeled fields are an existing affordance
  projects may add);
  [`spec/backlog.md:44-53`](../../../../spec/backlog.md)
  (`Graduated — <plan-slug>` is the exact-match-tracked form);
  [`spec/planning/shared.md:546-561`](../../../../spec/planning/shared.md)
  (exact-match label discipline — grounds rejecting (c));
  [`docs/backlog.md:7-21`](../../../backlog.md)
  (`repo-rooted-doc-links` is Open / post-1.0 — grounds
  rejecting (d)).

### D6 — Phase-vs-Task naming: rename `milestone.md`'s child unit to "Task", fold into this task

`milestone.md`'s required sections name the milestone's direct
child "Phase" (Phase Status, Cross-Phase Invariants/Decisions/
Risks). The slug grammar and the level picker both name the
milestone's child a *task*; "phase" is the task's child, one
level lower. The live milestone doc already renamed these to
"Task …" as a disclosed variance. The originating backlog entry
offers fold-in-or-split; folding in is chosen because the rename
edits the same `milestone.md` required-section surface D2 already
touches, and splitting would fragment one cohesive correction.

- `Verified by:`
  [`spec/planning/shared.md:97-106`](../../../../spec/planning/shared.md)
  (slug grammar: milestone's child segment is `t` = task);
  [`spec/planning/task-plan.md:22-37`](../../../../spec/planning/task-plan.md)
  (level picker: milestone → task → phase);
  [`spec/planning/milestone.md:99-105`](../../../../spec/planning/milestone.md)
  (current required sections use "Phase" for the milestone child);
  [`docs/plans/workstream-tracker-1-0/m1-v0-2.md:25-31`](../../workstream-tracker-1-0/m1-v0-2.md)
  (live "Task Status" / "Cross-Task …" variance).

## Open decisions to make at plan-drafting

None. All six decisions above resolve at scoping time; this is a
spec-edit (code-shipping) plan, so the
decision-resolved-at-scoping-time bar applies in full and no
decision is carried open into drafting.

## Plan structure handoff

- **Doc-type / level.** Standalone task plan, root slug
  `spec-updates-contracts-and-gates`, no epic/milestone parent.
- **Phase count.** N = 1. The four correction clusters (child
  contracts, parent-doc gate, backlog affordance, Phase-vs-Task
  rename) are sequence-mates toward one outcome — an internally
  consistent spec snapshot — with no independent ship value, and
  the additive-snapshot invariant requires they land together so
  no vendored snapshot is internally inconsistent. Phase content
  is absorbed inline in the task plan; no phase plan files.
- **PR count.** 1 implementing PR. Branch-test estimate: ~5
  markdown spec files, prose-only additions, well under the
  >5-subsystem / >300-LOC split threshold.
- **Sections the plan carries.** Status, Context, Goal,
  Contracts, Files to touch (estimate-labeled), Validation Gate,
  Cross-Cutting Invariants, Self-Review Audits, Documentation
  currency, Out of Scope, Backlog Impact, Related Docs.

## Reality-check inputs

Load-bearing inputs the plan must re-confirm against current spec
files before the `In draft` → `Proposed` flip:

- The two-or-more-levels layering discipline still reads as in D1
  / D4. `Verified by:`
  [`spec/planning/shared.md:14-17`](../../../../spec/planning/shared.md).
- The section-variance recurrence-to-codify rule still reads as
  in D2. `Verified by:`
  [`spec/planning/shared.md:443-446`](../../../../spec/planning/shared.md).
- The child-contract variance has in fact recurred in ≥2 live
  parent docs of distinct doc-type. `Verified by:`
  [`docs/plans/workstream-tracker-1-0/README.md:120-128`](../../workstream-tracker-1-0/README.md);
  [`docs/plans/workstream-tracker-1-0/m1-v0-2.md:77-85`](../../workstream-tracker-1-0/m1-v0-2.md).
- The epic doc remains `Status: In draft` post-scope-lock (live
  evidence for D4). `Verified by:`
  [`docs/plans/workstream-tracker-1-0/README.md:1-4`](../../workstream-tracker-1-0/README.md).
- The exact-match label discipline still grounds D5(c)'s
  rejection. `Verified by:`
  [`spec/planning/shared.md:546-561`](../../../../spec/planning/shared.md).
- `repo-rooted-doc-links` is still Open / post-1.0 (grounds
  D5(d)'s rejection). `Verified by:`
  [`docs/backlog.md:7-21`](../../../backlog.md).
- **Unverifiable from this repo:** neighborly-events' backlog
  `Detail` format, named in the scaffold as "prior art to
  consult, not match." It lives in an external project not
  vendored here. D5 is decided from in-repo constraints alone and
  deliberately does not attempt to match an unseen format; this
  is recorded as a known gap, not a blocker.

## Related Docs

- [`../README.md`](../README.md) — the task plan this scopes.
- [`../../../backlog.md`](../../../backlog.md) — origin
  `plan-doc-child-contracts` entry.
- [`../../../../spec/planning/task-plan.md`](../../../../spec/planning/task-plan.md)
  — the rules this scoping + plan pair follows.
