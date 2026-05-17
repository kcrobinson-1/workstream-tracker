# Scoping — stub-children-on-parent-promotion

Transient deliberation paired with
[`../README.md`](../README.md). Deletes at this task's terminal
PR per [`task-plan.md`](../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns"; survives in git history. No Status
field (scoping docs don't carry one).

## Context summary

Graduating the
[`stub-children-on-parent-promotion`](../../../backlog.md#stub-children-on-parent-promotion)
backlog entry. The work is a planning-spec rule: when a parent
doc (epic/milestone) passes its `` `In draft` → `Proposed` ``
promotion gate, the promoting PR also seeds skeleton child docs
so planned-but-unstarted nodes render and auto-registration has
a frontmatter slug to derive from. This is a doc-only decision
plan — the deliverable is the recorded decisions plus the spec
prose; no product code ships (rendering of a slug+Status-only
doc already works, see D4).

## Why scoping (not narrow-surface skip)

Fails narrow-surface criterion 4: the change carries a real
cross-cutting concern — convention-named stub slugs vs. the
spec's server-generated descendant-slug rule must be reconciled
across `shared.md`, not just added. That reconciliation is the
load-bearing decision (D1) and is why a scoping deliberation
exists. Criteria 1–3, 5 hold (single subsystem: the planning
spec; ≤8 files; no public-API/code/schema change; no novel
mechanism).

## Decisions made at scoping time

### D1 — Slug-allocation reconciliation: seeded stub slug is author-supplied-at-promotion, not server-generated

**Decision.** A child slug seeded by the parent-promotion gate
is **declared in the stub's frontmatter — validated for format,
author-supplied, exactly as root slugs already are**. The
file-write is a *declaration of identity only*: it performs no
server-side creation — no node, no work-instance, no allocation
call. The child session later **asserts** the declared slug to
the server when it registers a work-instance, via the
exact-slug create-or-attach path (m1-t1, landed) — "use this
slug," not "give me a slug." The spec's "Descendant slugs are
server-generated" rule is **loosened, not replaced**: server
slug-generation remains the path for any registration with no
pre-declared frontmatter slug; a pre-declared slug is asserted,
never server-generated.

**Why this dissolves the divergence.** The backlog entry's
stated risk is two allocators (convention-assigned stub slugs
vs. the server's allocation counter) drifting apart. The
resolution is that for a pre-declared slug **there is no
allocation at all** — not at file-write (declaration only) and
not at registration (assertion of the already-declared string).
The server counter only runs when nothing was pre-declared, so
the two paths are mutually exclusive and cannot diverge. The
milestone Task Status tables already name `-tN` slugs by
convention; this makes that convention the *sanctioned declared
identity* rather than hardening an inconsistency. The exact-slug
create-or-attach path is the assertion mechanism — depended on,
not replaced.

**Rejected alternatives.**
- *(a) Server allocates at seed time* — parent-promotion calls
  the allocation API for each child. Rejected: promotion is a
  spec/PR-authoring moment, not a tool-runtime moment; coupling
  the gate to a live server call makes a prose-spec obligation
  depend on tool availability, and the project's tenet is
  prompt-driven agent behavior over runtime enforcement.
- *(c) Provisional placeholder slugs the server rewrites on
  first registration* — rejected: violates slug immutability
  ([`shared.md`](../../../../spec/planning/shared.md) "Slug is
  identity") and breaks every convention-named `-tN` reference
  already written into milestone Task Status tables.

Verified by:
[`../../../../spec/planning/shared.md:127-131`](../../../../spec/planning/shared.md)
(current "Slug generation" — root author-supplied, descendant
server-generated; the rule D1 narrows);
[`../../../../spec/planning/shared.md:113-119`](../../../../spec/planning/shared.md)
(slug encodes creation-time position, is immutable identity);
[`../workstream-tracker-1-0/m1-t1-multi-wi-per-slug.md:113-142`](../../workstream-tracker-1-0/m1-t1-multi-wi-per-slug.md)
(exact-slug create-or-attach flow the stub depends on);
[`../workstream-tracker-1-0/m1-t1-multi-wi-per-slug.md:377-389`](../../workstream-tracker-1-0/m1-t1-multi-wi-per-slug.md)
(m1-t1 Risk Register explicitly leaving stub-seeding to this
entry, not replacing the exact-slug flow).

### D2 — A seeded stub carries `Status: In draft`

**Decision.** Seeded stubs carry `Status: In draft` — the
canonical pre-`Proposed` token. An unstarted seeded child is a
plan doc whose drafting has not happened; `In draft` is exactly
that state and renders as an in-draft tree node with no extra
lifecycle value invented.

Verified by:
[`../../../../spec/planning/shared.md:265-286`](../../../../spec/planning/shared.md)
("Plan-doc Status" lifecycle + the do-not-invent-adjacent-states
rule that bars a bespoke "stub" status);
[`../../../../spec/planning/task-plan.md:98-102`](../../../../spec/planning/task-plan.md)
(`In draft` is the canonical pre-`Proposed` label).

### D3 — A seeded stub is exempt from required-sections / no-implementation-prescription until its own drafting session

**Decision.** A doc that is a parent-promotion stub is not held
to "Required and optional sections" or "Plans describe
contracts, not implementation" until its own `In draft` →
`Proposed` drafting session runs. Without this carve-out a stub
fails its own promotion gate trivially the instant it exists,
which would make the seeding obligation self-contradictory. The
stub carries only: frontmatter (slug, Status, short_description)
plus the parent's WHAT-contract block and any illustrative
examples for that child copied in (D6).

Verified by:
[`../../../../spec/planning/task-plan.md:169-214`](../../../../spec/planning/task-plan.md)
("Required and optional sections" — what a non-stub plan owes);
[`../../../../spec/planning/shared.md:288-356`](../../../../spec/planning/shared.md)
(the gate the stub is exempt from until its own drafting).

### D4 — Tree rendering is out of scope; it already works

**Decision.** No tool/render change. The site renders a node
from parsed frontmatter; a doc with only slug+Status (no long
description, no body) is an already-supported render case, so
"planned-but-unstarted nodes render earlier" is achieved by the
stub existing, not by new rendering code.

Verified by:
[`../../../../internal/site/render_test.go:73`](../../../../internal/site/render_test.go)
(`{Slug: "alpha", Status: "Proposed"}` with no LongDescription
is an existing valid render case).

### D5 — A stub creates no work-instance; the triage-zone open question is not pre-empted

**Decision.** Seeding a stub creates a plan-tree *doc* only —
never a work-instance. Orphan/unattached work-instances remain
the epic's deferred "Triage zone in 1.0?" open question; this
task neither resolves nor forecloses it, and the plan states the
boundary explicitly so a reviewer reads it as a deliberated
intersection (per the backlog entry's "deliberate with, not in
isolation" instruction), not an omission.

Verified by:
[`../workstream-tracker-1-0/README.md:199-204`](../../workstream-tracker-1-0/README.md)
(the deferred "Triage zone in 1.0?" open question);
[`../../../backlog.md`](../../../backlog.md) `stub-children-on-parent-promotion`
(the "intersects two unsettled areas, deliberate with them"
instruction);
[`../../../../spec/backlog.md:83-106`](../../../../spec/backlog.md)
(backlog entries get no work-instances; the API has no backlog
node-type).

### D6 — Requirements split into artifact requirements vs. expected agent behavior

**Decision.** The plan's contract is written in two registers:
**artifact requirements** (checkable facts about the seeded
files and the parent doc — legitimately `must`) and **expected
agent behavior** (what the promotion-gate prompt instructs the
agent to do — best-effort, observable, not guaranteed). The spec
prose added to the gate is framed as a prompted obligation whose
miss is an accepted, observable residual, not a determinism
guarantee — per the project tenet that agent-run process
failures are acceptable and the lever is prompt engineering +
visibility, not enforcement tooling.

Verified by:
[`../../../../spec/planning/shared.md:288-356`](../../../../spec/planning/shared.md)
(the parent-doc gate is itself a prompted self-review, not a
tool-enforced check — the precedent register the new step joins);
[`../../../backlog.md`](../../../backlog.md)
`deterministic-interactive-registration` ("observable
best-effort … not deterministic" — the existing project
vocabulary this reuses).

### D7 — Stub-seeding binds all three parent→child relationships, via both promotion gates symmetrically

**Decision.** The seeding obligation binds every relationship
the "Parent-doc child contracts" rule binds: epic→milestone
(`Milestone Contracts`), milestone→task (`Task Contracts`), and
**task-plan with N≥2 phases→phase (`Phase Contracts`)**. It is
therefore wired into *both* promotion gates symmetrically — the
parent-doc gate in
[`shared.md`](../../../../spec/planning/shared.md) (epic/milestone)
and the task/phase gate in
[`task-plan.md`](../../../../spec/planning/task-plan.md) (the
N≥2 task plan promoting). A task plan with N=1 absorbs phase
content inline and has no separate phase children to seed —
nothing to do there, stated as a boundary so implementation
doesn't try to stub inline phases.

**Why uniform extension, not a phase carve-out.** Caught in
review: the original plan named only `Milestone Contracts` /
`Task Contracts` and wired only the parent-doc gate, while
editing the three-doc-type "Parent-doc child contracts" rule.
That leaves two incoherent outcomes — the shared-rule edit
silently imposes stub-seeding on phase children with no gate to
trigger it, or the edit contradicts the rule's explicit
three-doc-type binding. Scoping stub-seeding to two of three
relationships would require carving a phase exception *into* a
rule whose whole point is binding all three at once. The
motivating value (render planned-but-unstarted nodes; give
registration a declared frontmatter slug) applies to phase
children identically — a phase plan is a rendered tree node and
its implementing PR carries a work-instance. Uniform extension
also matches the project's demonstrated gate-symmetry value (the
`promotion-gate-explicit-checklist` task explicitly preserved
task/phase ↔ parent-doc gate symmetry).

**Rejected alternative.** *Scope to epic/milestone only,
explicitly exempting `Phase Contracts`.* Rejected: it weakens
the feature against its own motivation for no stated reason, and
the exemption prose would have to live inside the three-doc-type
shared rule, which is the exact incoherence this decision
removes.

Verified by:
[`../../../../spec/planning/shared.md:551-592`](../../../../spec/planning/shared.md)
("Parent-doc child contracts" — explicitly binds three
doc-types; `Phase Contracts` is the task-plan-N≥2 realization);
[`../../../../spec/planning/shared.md:288-299`](../../../../spec/planning/shared.md)
(the parent-doc gate binds epic/milestone only — task/phase
plans do not load it, so a second gate is required);
[`../../../../spec/planning/task-plan.md:441-447`](../../../../spec/planning/task-plan.md)
(the symmetric task/phase gate the seeding step must also join);
[`../../../../spec/planning/task-plan.md:190-194`](../../../../spec/planning/task-plan.md)
(`Phase Contracts` required when N≥2; absorbed inline so absent
at N=1);
[`../promotion-gate-explicit-checklist/README.md:151-169`](../../promotion-gate-explicit-checklist/README.md)
(the two-gate-symmetry precedent this mirrors).

## Inputs resolved before promotion

- **I1 — D1 ratification — RESOLVED 2026-05-17.** D1 loosens a
  load-bearing cross-level rule ("Slug generation") and touches
  the deferred triage-zone boundary, so the flip was held until
  the user ratified the model. Ratified by the user with the
  declare-vs-assert refinement now folded into D1/C2: file-write
  declares identity (no server-side creation); the child session
  asserts the declared slug via exact-slug create-or-attach;
  server-generation stays as the no-pre-declared-slug fallback.
  No open inputs remain; promotion gate may run.

## Plan-structure handoff

Standalone task plan, N = 1 (no phases) — mirrors the
[`promotion-gate-explicit-checklist`](../../promotion-gate-explicit-checklist/README.md)
precedent (same subsystem, same spec-prose-only shape). Plan
doc at [`../README.md`](../README.md); no phase plan files.
