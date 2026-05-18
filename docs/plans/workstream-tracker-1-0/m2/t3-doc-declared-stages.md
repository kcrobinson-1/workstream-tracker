---
slug: workstream-tracker-1-0-m2-t3
Status: In draft
short_description: Additive spec field + parser + per-node progress-box render driven by the doc
---

# Task 3 — Doc-declared progress stages (spec-first)

## Status

`In draft` — **pending human resolution of the open questions
below; not promoted, no promotion gate run.** This plan was
drafted just-in-time by a spawned planning session whose endpoint
is `In draft` with the genuine decisions decomposed and surfaced,
not resolved. The
[`task-plan.md`](../../../../spec/planning/task-plan.md)
`` `In draft` → `Proposed` `` promotion-gate self-review has
**not** been run and this plan makes no claim that it has. The
plan stays `In draft` until the human resolves the open questions,
at which point a resolving drafting session re-confirms the
reality-check inputs and runs the promotion gate.

Blocking open questions (decomposed with trade-offs and a
recommendation in
[`scoping/t3-doc-declared-stages.md`](scoping/t3-doc-declared-stages.md)
"Open questions for the human"):

- **OQ1 (central)** — field shape: ordered stage-label list vs.
  per-stage counts (shapes A1/A2/A3/B1/B2/C). Recommendation: A2.
- **OQ2** — field name and reserved-token spelling (exact-match
  discipline). Recommendation: `progress_stages`; Drafting box
  render-side with no doc-visible token.
- **OQ3** — Drafting box render-side-reserved vs. doc-declared
  (resolves jointly with OQ1).
- **OQ4** — per-level applicability: independent (no inheritance)
  vs. inherited. Recommendation: no inheritance.
- **OQ7** — box row gated by field-presence (α) vs. Status (β).
  Recommendation: α.
- **OQ5** — phase split N = 1 vs. N ≥ 2; downstream of OQ1.
  Assessment: most likely N = 1; left open, branch test runs at
  the resolving drafting.

These are WHAT-level behavior decisions; per the
[`task-plan.md`](../../../../spec/planning/task-plan.md)
just-in-time rule they cannot be deferred to phase drafting, so
the plan is held at `In draft` rather than carried as pending
inputs. Contracts below that depend on them are written
**conditionally / flagged**, not as locked decisions.

## Context

This plan covers a small, additive extension to the plan-doc
spec: a way for a plan doc to declare, in its own frontmatter, the
progress stages it expects to move through — and the
visualization rendering those stages as a row of progress boxes on
every node in the forest. Today the forest shows each node's
Status badge and, since t2 landed, a nested collapsible box per
node; there is no per-node sense of "how far along the stages this
work is." This task adds that sense, driven entirely by the plan
doc rather than inferred.

It is being done now because t2 (the expanded nested-box render)
has Landed and established the per-node box this work renders
inside, and because the m2 milestone deliberately deferred the one
hard decision — the exact frontmatter shape — to this task's
drafting. The shape is the spec-change surface itself: it is the
contract every consumer project that vendors the spec will author
against, which is why it is decomposed carefully and left for
human decision rather than invented.

The surfaces touched are the plan-tree frontmatter parser, the
forest-region render template it feeds, the optional/additive
section of the plan-doc spec where the new field is documented,
and their tests — plus the design doc's render-scope section. No
API, schema, route, or dependency change; the page stays
server-rendered and refresh-to-update, and a doc that omits the
field is unaffected.

## Goal

A plan doc may optionally declare its own ordered progress stages
in frontmatter. The plan-tree walker reads the field with the same
absence-tolerance the existing optional fields have. Every node in
the forest renders a row of progress boxes whose count and order
come from that node's doc. A doc that omits the field — including
the already-supported `slug` + `Status: In draft` stub — renders
the Drafting box only, with no error, skip, or broken layout. The
field is optional and additive: pre-existing docs and vendored
spec consumers are unaffected by its absence, and t4's roster
deliberately does not consume or schematize it.

Verifiable when (conditional on the open questions resolving):
rendering the dogfood `docs/plans/` tree shows, on every node box,
a progress-box row whose count and order match that node's
declared stages; a stub and any field-omitting doc render exactly
the Drafting box as an intentional observed state (not an empty or
errored row); the existing actor markers, Status badge, long
description, and related-PR list still render unchanged; and the
walk-on-every-request render path is unchanged.

## Contracts

Final WHAT shape. Clauses that depend on an unresolved open
question are flagged **[conditional — OQN]** and state the
contract under the recommended shape with the alternative noted;
they lock when the human resolves the question. HOW grounding and
the decomposition with rejected shapes live in
[`scoping/t3-doc-declared-stages.md`](scoping/t3-doc-declared-stages.md).

- **C1 — Additive optional frontmatter field.** A plan doc may
  carry an optional frontmatter field declaring its progress
  stages. The field is optional and additive: a doc omitting it
  remains valid and renders with no warning, error, or skip;
  pre-existing docs and vendored spec consumers are unaffected.
  The plan-tree walker reads it with the same absence-tolerance
  posture as `short_description` / `related_prs` (absent /
  wrong-typed / partially-malformed never errors or skips the
  doc). The field is documented in the plan-doc spec adjacent to
  the existing optional fields. `Verified by:`
  [`stringList` / `parsePlanDoc` in
  walker.go](../../../../internal/site/walker.go) (the tolerant
  optional-field read pattern this extends);
  [`shared.md` "Optional `related_prs`
  field"](../../../../spec/planning/shared.md) (the
  optional/additive precedent and its absence-benign guarantee).
- **C2 — [conditional — OQ1/OQ2] Field shape and name.** The
  field encodes the node's progress stages such that **box count
  and order are derived from the doc**. Under the recommended
  shape (scoping OQ1 **A2**, OQ2 name `progress_stages`): an
  ordered YAML block sequence of stage-label strings, read by the
  existing tolerant block-sequence decoder, box count = number of
  entries, order = sequence order, label per box = the entry
  string. If the human selects a per-stage-count shape (OQ1 **B2**)
  instead, this clause becomes an ordered sequence of
  {stage, count} records and box semantics shift from
  one-box-per-stage to one-box-per-PR-within-stage; the parser
  acquires a net-new typed decoder with its own absence-tolerance.
  This clause locks on OQ1/OQ2 resolution. `Verified by:`
  [m2 README "Task Contracts" t3 row](README.md) (count + order
  come from the doc — the locked WHAT this conditional realizes);
  [`stringList` in
  walker.go](../../../../internal/site/walker.go) (the A2-shape
  decoder reused vs. the B2-shape decoder that is net-new);
  [`shared.md` "Decompose options into shapes before
  analyzing"](../../../../spec/planning/shared.md) (why the shape
  is decomposed, not invented at milestone level).
- **C3 — [conditional — OQ1/OQ3] Drafting box for the
  absent/stub case.** A doc that omits the field — including a
  `slug` + `Status: In draft` stub — renders exactly the Drafting
  box, as an intentional observed state, never an empty or errored
  row. Under the recommendation (render-side reserved Drafting box,
  declared list = post-drafting stages), the Drafting box is a
  single render-side synthesis serving both the absent/stub case
  and the leading box of a drafted doc; there is no doc-visible
  `Drafting` token. If OQ3 resolves to a doc-declared first stage,
  this clause additionally requires a synthesized Drafting box for
  the absent case (the stub has no field), i.e. two code paths.
  Locks on OQ1/OQ3. `Verified by:`
  [`parsePlanDoc` in
  walker.go](../../../../internal/site/walker.go) (a stub omits
  the field; the absent read returns the zero value, so the
  Drafting box cannot derive from the field for a stub);
  [`stub-children-on-parent-promotion` README "Stub ≠
  work-instance"](../../stub-children-on-parent-promotion/README.md)
  (a `slug` + `Status: In draft` stub is an already-supported
  render case that stays valid; seeding it creates no
  work-instance and does not resolve the deferred triage
  question).
- **C4 — Progress-box row renders inside t2's per-node box, all
  levels.** The row renders at every node level (root, milestone,
  task, phase — there is no `epic` node; the root is the top box)
  inside the per-node box t2 established, in the forest region
  only. t3 adds the row; it does not alter the box shell, the
  collapse mechanism, the Status badge, the actor markers, the
  long description, or the related-PR list — those still render
  unchanged on every box. t3 edits the forest render and the
  parser only; the shell (`render.go`) and roster (`roster.go`)
  are not touched. `Verified by:`
  [`forest.go` `node` / `node-header` / `node-detail`
  templates](../../../../internal/site/forest.go) (the Landed t2
  box the row renders inside, and the preserved surfaces);
  [m2 README "Cross-Task Invariants" → "the shell is
  t1's"](README.md) (the file-boundary invariant t3 stays inside;
  forest.go is the forest region's owned surface);
  [`internal/slugs/slugs.go`](../../../../internal/slugs/slugs.go)
  (the node levels the row renders at).
- **C5 — [conditional — OQ7] Render-gating discriminator.** Under
  the recommendation (OQ7 **α**, field-presence gated): a node
  renders its declared-stage boxes iff the field is present;
  absent ⇒ Drafting box only; the discriminator is field presence,
  Status-independent, with no coupling to the exact-match Status
  lifecycle token. If the human selects **β** (Status-gated),
  non-Drafting boxes render only once Status reaches `Proposed`,
  adding a Status-token coupling into the progress render that
  must coexist with the "Unknown Status values render gracefully"
  rule. Locks on OQ7. `Verified by:`
  [`parsePlanDoc` in
  walker.go](../../../../internal/site/walker.go) (Status and the
  new field are independent reads today; α adds no coupling, β
  adds one);
  [`shared.md` "Plan-doc Status" / "Unknown Status values render
  gracefully"](../../../../spec/planning/shared.md) (the
  exact-match Status surface β would couple to).
- **C6 — [conditional — OQ4] Per-level applicability.** Under the
  recommendation (no inheritance): each node's own doc governs its
  row; a node whose own doc omits the field renders the Drafting
  box only even if an ancestor declares stages. If the human
  selects inheritance, the render acquires a parent/child walk for
  stage resolution. Locks on OQ4. `Verified by:`
  [`shared.md` "Optional `short_description` field" / "Optional
  `related_prs` field"](../../../../spec/planning/shared.md)
  (the per-doc, absence-benign, no-inheritance precedent);
  [`buildTree` / `PlanNode` in
  tree.go](../../../../internal/site/tree.go) (per-doc fields are
  attached per node today; no field inherits across parent/child).
- **C7 — Render-side tests are semantic/structural.** The
  progress-box row is asserted by structure and presence (the row
  is present with the expected box count/order for a declaring
  doc; a field-omitting doc and a stub render exactly the Drafting
  box; absent/malformed field never errors or drops the node;
  preserved t2 surfaces still present), not a byte-identity
  literal — consistent with t2's C6 semantic-assertion posture so
  a fresh byte pin does not just break at the next render task.
  `Verified by:`
  [m2-t2 plan C6 + `scoping/t2-expanded-render.md`
  S4](t2-expanded-render.md) (the semantic-not-byte-exact
  precedent for net-new render partials this task inherits).

## Cross-Cutting Invariants

Inherited from the m2 milestone; t3 must not regress them. Each is
a reviewer-flag candidate.

- **Spec changes stay additive.** The new field is optional and
  additive only — no breaking change to vendored spec consumers,
  consistent with the parent-epic "Spec contract is additive"
  invariant. This task carries an explicit "is this additive?"
  check: a doc omitting the field renders with no error/skip,
  falling back to the Drafting-box-only shape, mirroring the
  landed `short_description` / `related_prs` precedent. `Verified
  by:` [m2 README "Cross-Task Invariants" → "Spec changes stay
  additive"](README.md);
  [`shared.md` "Optional `related_prs`
  field"](../../../../spec/planning/shared.md).
- **Opposite spec postures are intentional — do not homogenize.**
  t3 tightens the spec with a checkable additive declared-stages
  field; it must NOT be reconciled toward t4's deliberately
  schema-loose reported JSON, nor t4's JSON loosened toward this
  field. Any drafting that reconciles the two postures is the
  defect, not the fix. `Verified by:`
  [m2 README "Cross-Task Invariants" → "Opposite spec
  postures"](README.md).
- **Stub render is preserved, not pre-empted.** A `slug` +
  `Status: In draft` stub remains a valid render; t3 makes it
  render exactly the Drafting box and creates no work-instance and
  does not resolve the deferred triage-action question. `Verified
  by:`
  [`stub-children-on-parent-promotion` README "Stub ≠
  work-instance"](../../stub-children-on-parent-promotion/README.md);
  [m2 README "Cross-Task Invariants" → "Stub render case is
  preserved"](README.md).
- **Render path stays walk-on-every-request.** No caching,
  file-watch, or in-memory build-up; the field is read and the row
  rendered inside the existing per-request build path. `Verified
  by:`
  [`Server.index` in site.go](../../../../internal/site/site.go).
- **Shell is t1's; t3 stays in the forest region.** t3 edits the
  parser and the forest render only; `render.go` and `roster.go`
  are not touched. `Verified by:`
  [m2 README "Cross-Task Invariants" → "the shell is
  t1's"](README.md).

## Naming

**[conditional — OQ1/OQ2; spelled once resolved.]** The plan
introduces one new optional frontmatter field (recommended name
`progress_stages`) and, depending on OQ1, possibly one additive
`PlanNode` field carrying the parsed stage list to the template
(precedent: t2's `ActiveInSubtree` additive field). Exact
identifiers are exact-match-checked tokens per
[`shared.md` "Quote labels whose enforcement depends on exact-match
matching"](../../../../spec/planning/shared.md); they are chosen
when OQ1/OQ2 resolve and copied verbatim into the spec, not
paraphrased. `Verified by:`
[`PlanNode` in tree.go](../../../../internal/site/tree.go)
(the additive-field precedent);
[`shared.md` exact-match-label
discipline](../../../../spec/planning/shared.md).

## Files to touch

_Estimate of the expected file shape, not a binding rule, and
itself partly conditional on OQ1/OQ5. Implementation may revise
this list when a structural call requires it — including touching
a file listed under "Intentionally not touched." Any deviation is
handled via the PR-body Estimate Deviations callout with the plan
reconciled to what shipped._

- **Modify:** the plan-tree frontmatter parser (the new optional
  field read; under OQ1 A2 this reuses the existing tolerant
  block-sequence decoder, under B2 it adds a net-new typed
  decoder); the forest-region render template + styles (the
  progress-box row partial inside the existing per-node box, plus
  its CSS); the render-input node struct if the chosen OQ1 shape
  needs a parsed value carried to the template (precedent: t2's
  additive field).
- **Modify (spec):** the plan-doc spec's optional/additive
  section — the new field documented adjacent to
  `short_description` / `related_prs`, optional + additive + the
  defined absent behavior (Drafting box only). This is the
  spec-first deliverable.
- **Modify (tests):** the forest-region test file (semantic
  assertions per C7: row count/order for a declaring doc;
  Drafting-only for a field-omitting doc and a stub; absent /
  malformed field never errors); the parser test file (the new
  field's tolerant read, absence/wrong-type/partial-malformed).
- **Modify (docs):** the m2 README parent doc — Task Status t3
  row and t3 prose to a drafted `In draft` plan, and the
  "Doc-declared-stages frontmatter shape" Cross-Task Decisions
  entry pointed at the scoping section and marked OPEN (this
  drafting change, per parent-doc currency); `design/v0.1-design.md`
  §7 updated to describe the doc-driven progress-box row on the PR
  that lands the implementation.
- **Intentionally not touched:** the shell render file
  (`render.go`) and the roster render file (`roster.go`) — m2
  file boundary; the API, DB schema, register client, routing;
  any dependency (the field reuses the existing goldmark-meta /
  `yaml.v2` decode under OQ1 A2; B2 adds no dependency, only Go
  decode code).

## Validation Gate

Conditional on the open questions resolving; the resolving
drafting finalizes the gate and runs the promotion-gate
self-review (this session does neither).

- `go build ./...`, `go vet ./...`, `go test ./...` all pass (per
  [`docs/dev.md`](../../../../docs/dev.md)).
- Unit tests cover: a declaring doc renders a progress-box row of
  the declared count and order; a field-omitting doc and a `slug`
  + `Status: In draft` stub render exactly the Drafting box; an
  absent / wrong-typed / partially-malformed field never errors or
  drops the node (the additive guarantee); the preserved t2
  surfaces (Status badge, actor markers, long description,
  related-PR list, collapse) still render on every box.
- **Manual UI capture** (no screenshot tooling exists; capture is
  manual per [`docs/dev.md`](../../../../docs/dev.md) UI Review):
  screenshots at ~1280px against a real `go run` reading a
  populated `DB_PATH` and the `docs/plans/` dogfood tree —
  confirming the progress-box row on declaring nodes and the
  Drafting-box-only render on the stub.
- **Drafting-box-only is observed, not assumed.** Per
  [`task-plan.md`](../../../../spec/planning/task-plan.md) "Bans on
  surface require rendering the consequence," the absent-field
  Drafting-box render is verified by observing a field-omitting
  doc / stub render as an intentional state, not inferred from the
  diff — a required gate step.
- Self-review audits below run at the implementing commit
  boundary.

## Self-Review Audits

From
[`docs/agents/local/self-review-catalog.md`](../../../../docs/agents/local/self-review-catalog.md)
(no repo-specific audits exist); diff surface is parser + frontend
+ spec prose:

- **validation-honesty** — a "the row renders / the stub shows
  only Drafting" claim is valid only if observed against a real
  `go run`, not asserted from the diff. The lens on the manual
  gate steps.
- The general self-review checklist
  ([`how-to-use.md`](../../../../docs/agents/shared/self-review/how-to-use.md))
  layered on top.

No data / CI / runbook audit surfaces apply (no schema, pipeline,
or operational doc touched).

## Risk Register

- **t3's spec change is read as breaking by a vendored consumer.**
  Mitigation: the additive-spec invariant + the explicit "is this
  additive?" check; the field is optional with a defined
  absent-behavior (Drafting box only), mirroring the landed
  `short_description` / `related_prs` precedent. `Verified by:`
  [`shared.md` "Optional `related_prs`
  field"](../../../../spec/planning/shared.md).
- **The two spec postures get homogenized under review pressure.**
  A reviewer or later task may "tidy" this field toward t4's
  schema-loose JSON or vice versa. Mitigation: the posture-tension
  cross-cutting invariant names reconciliation as the defect.
- **Field-shape churn if OQ1 reopens after implementation
  starts.** The shape is the consumer-facing spec contract;
  reopening it after docs author against it is costly. Mitigation:
  OQ1 is held as a pre-implementation human decision (plan stays
  `In draft`), not deferred into implementation; the scoping doc
  decomposes it with a recommendation so the decision is cheap to
  make now and expensive to defer.

## Out of Scope

- **The session roster and its reported JSON.** t4's surface and
  its deliberately schema-loose posture; t3 neither consumes nor
  schematizes it.
- **AI / probabilistic prose-to-data extraction.** The vision §7
  open question of structured-fields vs. AI extraction vs. hybrid
  is intersected, not resolved: t3 delivers the structured-fields
  path for declared stages; whether softer signals are
  AI-extracted later stays a vision open question. `Verified by:`
  [`design/vision.md` §7 "How the tool turns plan-document prose
  into the data it renders"](../../../../design/vision.md).
- **On-box actor icons.** Deferred milestone work assuming
  node-level actor markers remain; t3 preserves the markers and
  adds nothing here.
- **Triage action for stubs/unbound work.** A stub still creates
  no work-instance; the deferred triage-action question is
  untouched.

## Backlog Impact

None. No [`docs/backlog.md`](../../../../docs/backlog.md) entry
graduates, is deleted, split, or shifted by this task. Two
[`design/vision.md` §7](../../../../design/vision.md) open
questions are *referenced as deliberated intersections, not
resolved*: "How the tool turns plan-document prose into the data
it renders" (t3 is the structured-fields cut for declared stages;
AI extraction stays open) and the "smaller open question" of
placeholder vs. only-what-committed sub-stage cells (interacts
with OQ7; surfaced, not decided).

## Related Docs

- [`README.md`](README.md) — parent milestone; its Task Status t3
  row, t3 prose, and "Doc-declared-stages frontmatter shape"
  Cross-Task Decisions entry are reconciled to this drafting in
  the same change (pointed at scoping, marked OPEN — not
  resolved).
- [`scoping/t3-doc-declared-stages.md`](scoping/t3-doc-declared-stages.md)
  — the decision-space decomposition (OQ1–OQ5/OQ7 with shapes,
  trade-offs, `Verified by:` grounding, and recommendations) this
  plan's conditional Contracts realize.
- [`t2-expanded-render.md`](t2-expanded-render.md) and
  [`scoping/t2-expanded-render.md`](scoping/t2-expanded-render.md)
  — the Landed sibling whose per-node box hosts the progress-box
  row, and the semantic-not-byte-exact test precedent (C7).
- [`../../stub-children-on-parent-promotion/README.md`](../../stub-children-on-parent-promotion/README.md)
  — the landed stub render case the Drafting-box-only behavior
  anchors to.
- [`../../../design/vision.md`](../../../design/vision.md) §7 —
  the prose-to-data and sub-stage-cell open questions t3
  intersects but does not resolve.
- [`../../../design/v0.1-design.md`](../../../design/v0.1-design.md)
  §7 — the render-scope section updated by the implementing PR.
- [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md)
  and [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md)
  — the optional/additive + exact-match-label + decompose-shapes
  rules this task is structured against.
