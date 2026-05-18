# Scoping — m2 t3: Doc-declared progress stages (spec-first)

Transient deliberation doc for the
[`workstream-tracker-1-0-m2-t3`](../t3-doc-declared-stages.md) task
plan. Deletes in batch with sibling scoping docs at the
milestone-terminal PR per
[`task-plan.md`](../../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns." No Status field — scoping docs are
not part of the plan-doc lifecycle.

> **This scoping doc surfaces a decision space; it does not
> resolve it.** This is a spawned, just-in-time drafting session
> whose endpoint is a task plan at `Status: In draft` with the
> genuine decisions decomposed into shapes and left OPEN for the
> human. The "Decisions made at scoping time" section below
> records only the procedural calls a rule lets this session make
> (scoping-method calls); every product/spec decision is in "Open
> questions for the human," not resolved here. No promotion gate
> was run; no plan was promoted.

## Context summary

t3 adds an additive, optional `spec/` frontmatter affordance that
lets a plan doc declare its own progress stages, teaches the
plan-tree walker to read it, and renders a row of progress boxes —
one row per node at every level (root, milestone, task, phase) —
whose box count and order come from the doc. A doc that omits the
field (the already-supported `slug` + `Status: In draft` stub)
renders only the Drafting box and is never errored or skipped. The
box row renders inside the per-node box t2 (Landed) established in
the forest region; t3 does not touch the shell or the roster, and
t4 deliberately does not consume or schematize this field
(posture-tension invariant). The exact field shape is the
spec-change surface itself and is the central deferred decision
the m2 milestone routed to this drafting session — this doc
decomposes it into candidate shapes and surfaces it OPEN rather
than locking it.

## Decisions made at scoping time

Only scoping-method calls a rule explicitly authorizes this
session to make. No product or spec-field decision is recorded
here — those are in "Open questions for the human" below.

### S1 — Narrow-surface carve-out is DECLINED; scoping doc is warranted

[`task-plan.md`](../../../../../spec/planning/task-plan.md)
"Narrow-surface plans may skip the scoping doc" requires ALL five
conditions to hold. At least two fail, so the carve-out does not
apply and the default "scoping first" path is taken:

- **Criterion 4 (no new cross-cutting invariant) fails.** t3
  carries an additive/absent-tolerant guarantee threaded across
  the parser read, the render, and the stub/absent case
  simultaneously, plus the m2 posture-tension invariant (the
  field must stay a checkable tightening, never homogenized toward
  t4's schema-loose JSON). Cross-cutting invariants are the
  load-bearing reason scoping docs exist. `Verified by:`
  [m2 README "Cross-Task Invariants" → "Opposite spec postures
  are intentional"](../README.md) and "Spec changes stay
  additive"; [`task-plan.md` "Narrow-surface" criterion
  4](../../../../../spec/planning/task-plan.md).
- **The options-considered content is material.** The field shape
  is a genuine Choose-One decision with sub-shapes and mandatory
  multi-source `Verified by:` citations — exactly the content the
  carve-out excludes (it applies only when options analysis "would
  not produce material content"). `Verified by:`
  [`task-plan.md` "Narrow-surface plans may skip the scoping
  doc"](../../../../../spec/planning/task-plan.md);
  [`shared.md` "Decompose options into shapes before
  analyzing"](../../../../../spec/planning/shared.md).

This S1 is the deliberate invocation of the default path, mirroring
t2's S5; it is not a deferred decision.

### S2 — The field-shape decision is decomposed but NOT locked here

The m2 milestone routed the field name/structure/per-stage-counts-
vs-ordered-list call to "when t3 drafts" under the
[`shared.md`](../../../../../spec/planning/shared.md) "Decompose
options into shapes" and exact-match-label discipline. This
spawned session decomposes it (OQ1–OQ4, OQ7 below) with trade-off
analysis and a recommendation, then leaves it OPEN for the human
per the session's scope bound. Recording the decomposition is the
scoping-method act; choosing the shape is the human's. `Verified
by:` [m2 README "Cross-Task Decisions" → "Doc-declared-stages
frontmatter shape (decide when t3 drafts)"](../README.md);
[`shared.md` "Decompose options into shapes before
analyzing"](../../../../../spec/planning/shared.md).

## Open questions for the human

Each question is decomposed into shapes with trade-offs grounded
in cited code/spec. A recommendation is offered where the analysis
points clearly, but the choice is the human's; the paired task
plan stays `Status: In draft` until they are resolved, and its
Contracts that depend on them are written conditionally.

### OQ1 (central) — Field shape: ordered stage list vs. per-stage counts

The locked WHAT (m2 README t3 contract) requires only: a doc
declares its own progress stages; **box count and order come from
the doc**; all node levels; a field-omitting stub renders the
Drafting box only; optional and additive. `Verified by:`
[m2 README "Task Contracts" t3 row](../README.md). Within that
contract the milestone explicitly names the open axis as
"per-stage counts vs. an ordered stage list." Decomposed into
shapes (not category labels), per
[`shared.md`](../../../../../spec/planning/shared.md):

- **Shape A1 — bare ordered list of stage-label strings.** The
  field is a YAML block sequence of plain strings; box count is
  the number of entries, order is sequence order, each box's label
  is its string. This is byte-for-byte the decode shape the
  optional `related_prs` field already uses, readable by the
  existing tolerant `stringList` helper with no new parsing
  mechanism. `Verified by:`
  [`stringList` in walker.go](../../../../../internal/site/walker.go)
  (block-sequence-of-strings tolerant decoder, absent/non-sequence/
  non-string-element all yield empty-or-partial, never error);
  [`shared.md` "Optional `related_prs`
  field"](../../../../../spec/planning/shared.md) (the
  optional/additive block-sequence precedent).
- **Shape A2 — A1, with the Drafting box reserved/implicit (not a
  declared entry).** The field lists only post-drafting stages;
  the Drafting box is synthesized render-side and always precedes
  the declared stages. A field-omitting stub renders exactly the
  Drafting box because the synthesized box is the absent-case
  default; a drafted doc renders Drafting + its declared stages.
  One Drafting-box code path serves both stub and drafted docs.
  `Verified by:`
  [m2 README "Task Contracts" t3 row](../README.md) ("a stub …
  renders only the Drafting box; the rest appear once the doc
  reaches `Proposed`");
  [`stub-children-on-parent-promotion` README "Stub ≠
  work-instance"](../../../stub-children-on-parent-promotion/README.md)
  (a `slug` + `Status: In draft` stub is an already-supported
  render case that must stay valid).
- **Shape A3 — ordered list where the doc must include `Drafting`
  as its literal first entry.** A field-omitting stub still needs
  a synthesized Drafting box (a stub has no field by definition),
  so A3 forces every drafted doc to repeat a magic `Drafting`
  token while the absent/stub case still needs the synthesized
  default — two Drafting-box code paths and a doc-visible magic
  string. `Verified by:`
  [`parsePlanDoc` in walker.go](../../../../../internal/site/walker.go)
  (a stub omits the field entirely; the absent read returns the
  zero value, so the Drafting box cannot come from the field for a
  stub regardless of A2/A3).
- **Shape B1 — map of stage-name to count.** Encodes a count per
  named stage but a YAML mapping has no reliable order under the
  pinned goldmark-meta / `yaml.v2` decode, so "order comes from
  the doc" is not satisfiable without a second ordered-key
  convention. `Verified by:`
  [`stringList` in walker.go](../../../../../internal/site/walker.go)
  comment (the goldmark-meta / `gopkg.in/yaml.v2` v2.3.0 pin whose
  block-sequence ordering A1/A2 rely on; a mapping decode has no
  equivalent order guarantee).
- **Shape B2 — ordered list of {stage, count} records.** Order is
  preserved (a sequence); each record carries a stage label and a
  count (e.g. how many PRs that stage calls for). This is the
  structured-fields answer to the vision's named example "how many
  pull requests the plan calls for at each phase," and it makes a
  "progress box" mean a PR within a stage rather than a stage.
  Requires a new typed decoder (a sequence of maps), which adds
  its own absence-tolerance surface distinct from `stringList`.
  `Verified by:`
  [`design/vision.md` §7 "How the tool turns plan-document prose
  into the data it renders"](../../../../../design/vision.md)
  (names per-phase PR counts as the load-bearing example and
  leaves structured-fields-vs-AI-extraction open);
  [`stringList` in walker.go](../../../../../internal/site/walker.go)
  (the sequence-of-strings helper B2 cannot reuse — a
  sequence-of-maps decode is net-new).
- **Shape C — list of single-key {label} (or richer) records.**
  Functionally A1 with a heavier per-entry shape; extensible later
  but no current contract clause justifies the extra structure,
  and it shares B2's net-new sequence-of-maps decode cost without
  B2's count payoff. `Verified by:`
  [m2 README "Task Contracts" t3 row](../README.md) (the contract
  asks for count + order only; per-entry richness is uncontracted).

**Analysis.** Every shape except B1 satisfies "count and order
come from the doc." A1/A2 reuse the exact `related_prs` tolerant
decoder, giving the strongest alignment with the cited
optional/additive precedent and the lowest absent-case risk; B2/C
add a net-new sequence-of-maps decoder with its own
absence-tolerance surface. The crux is a product question the code
cannot answer: **does one progress box represent a stage, or a PR
within a stage?** That is exactly the vision §7 open question
("structured fields for the load-bearing numbers" vs. softer
signal) and is the per-stage-counts-vs-ordered-list axis the
milestone deferred here. **Recommendation: A2** (ordered list of
stage-label strings, Drafting reserved render-side) as the
lowest-cost shape satisfying every locked clause and reusing the
`stringList` precedent, with **B2 explicitly deferred to a later
milestone unless the human's product intent is that progress boxes
count PRs**. This is a recommendation, not a resolution — **OPEN**.

### OQ2 — Field name and reserved-token spelling (exact-match discipline)

The field name, and any reserved stage token, become exact-match-
checked identifiers the parser and consumers key off, governed by
[`shared.md` "Quote labels whose enforcement depends on exact-match
matching"](../../../../../spec/planning/shared.md). Candidate
field names: `progress_stages`, `stages`, `declared_stages`,
`progress`. Candidate Drafting treatment: a render-side reserved
box with **no doc-visible token** (preferred — keeps the magic
string out of authored docs and out of the exact-match surface)
vs. a literal `Drafting` token (couples to OQ1 shape A3).
`Verified by:`
[`shared.md` "Quote labels whose enforcement depends on exact-match
matching"](../../../../../spec/planning/shared.md) (the discipline
this naming is bound by);
[`parsePlanDoc` in walker.go](../../../../../internal/site/walker.go)
(the frontmatter key the chosen name is read at, alongside the
existing `slug` / `Status` / `short_description` reads). The
milestone says the token is "decided here … not invented at
milestone level," but per this session's scope bound it is
surfaced with a recommendation (name: `progress_stages`; Drafting:
render-side, no token), not locked — **OPEN**.

### OQ3 — Is the Drafting box render-side-reserved or doc-declared?

Tied to OQ1 A2 vs. A3. The contract requires a field-omitting stub
to render exactly the Drafting box; a stub has no field, so the
Drafting box must be synthesizable from absence regardless. Open:
does the render always prepend a reserved Drafting box (declared
list = post-drafting stages, one code path), or does the declared
list fully enumerate including a literal first stage (two code
paths + magic token)? `Verified by:`
[`parsePlanDoc` in walker.go](../../../../../internal/site/walker.go)
(absent field ⇒ zero value, so a stub's Drafting box cannot derive
from the field);
[m2 README "Task Contracts" t3 row](../README.md) (stub ⇒ Drafting
only). **Recommendation: render-side reserved Drafting box** —
**OPEN** (resolves jointly with OQ1).

### OQ4 — Per-level applicability: independent or inherited?

The contract says every node level renders a row whose count/order
come from "the doc." Open: must each level's own doc carry the
field (a node whose own doc omits it ⇒ Drafting-box-only, even if
its parent declares stages — no inheritance), or does a node
inherit a parent's declared stages? No-inheritance matches the
absent-tolerant posture of `short_description` / `related_prs`
(each doc's own frontmatter governs that node, absence is benign)
and avoids coupling the progress render to a tree walk. `Verified
by:`
[`shared.md` "Optional `short_description`
field"](../../../../../spec/planning/shared.md) and "Optional
`related_prs` field" (per-doc, absence-benign, no inheritance
precedent);
[`buildTree` / `PlanNode` in
tree.go](../../../../../internal/site/tree.go) (per-doc fields are
attached per node; no field today inherits across parent/child).
**Recommendation: no inheritance** — **OPEN**.

### OQ7 — Box row gated by field-presence or by Status?

The milestone sentence "a stub … renders only the Drafting box;
the rest appear once the doc reaches `Proposed`" admits two
readings:

- **Shape α — field-presence gated (Status-independent).** A doc
  renders Drafting-only iff the field is absent; a doc that
  declares stages renders them regardless of its Status. Matches
  the `short_description` / `related_prs` absent-tolerance
  precedent exactly and keeps the progress render free of any
  exact-match Status-token dependency. `Verified by:`
  [`parsePlanDoc` in walker.go](../../../../../internal/site/walker.go)
  (Status and the new field are independent frontmatter reads;
  nothing couples them today);
  [`shared.md` "Optional `related_prs`
  field"](../../../../../spec/planning/shared.md) (absence is the
  only render discriminator in the additive precedent).
- **Shape β — Status-gated.** Non-Drafting boxes render only once
  Status has reached `Proposed`, coupling the progress render to
  the exact-match Status lifecycle token and interacting with the
  "Unknown Status values render gracefully" rule. Matches a
  literal reading of the milestone sentence. `Verified by:`
  [`shared.md` "Plan-doc Status" + "Unknown Status values render
  gracefully"](../../../../../spec/planning/shared.md) (the
  exact-match token surface β would couple the progress render
  to);
  [`statusClass` use in the node template,
  forest.go](../../../../../internal/site/forest.go) (Status is
  consumed today only for the badge class, not for gating which
  sub-elements render — β would be a new coupling).

The vision's "smaller open questions" frames the adjacent axis
directly — render placeholder cells for not-yet-committed nodes
(friendlier) vs. render only what the plan committed to (stricter,
never misleading). `Verified by:`
[`design/vision.md` §7 "Smaller open
questions"](../../../../../design/vision.md). **Recommendation: α**
(field-presence gated; simpler, matches the cited additive
precedent, no Status coupling) — **OPEN**.

### OQ5 — Phase split (N = 1 vs N ≥ 2)

The m2 "Per-task phase splits" estimate guesses t3 plausibly N ≥ 2
(spec/parser, then progress-box render). Candidate boundaries
enumerated: (a) spec-field + parser read | progress-box render;
(b) single phase. The resolver is the
[`task-plan.md`](../../../../../spec/planning/task-plan.md)
"PR-count predictions need a branch test," run at the resolving
drafting after human input — not here. The split **hinges on
OQ1**: under A2 the substantive logic is tiny (one tolerant field
read reusing `stringList`, plus a render loop in the forest node
template), pointing strongly to N = 1, single subsystem
(`internal/site`) plus additive spec-doc prose; under B2 the
net-new sequence-of-maps decoder enlarges it but still likely N =
1. `Verified by:`
[`stringList` in walker.go](../../../../../internal/site/walker.go)
(the A2-path read is an existing helper, ≈ no new algorithmic
logic);
[`forest.go` node template](../../../../../internal/site/forest.go)
(the render host is one recursive template; a box-row partial is
markup, not subsystem count);
[m2 README "Cross-Task Decisions" → "Per-task phase
splits"](../README.md) (estimate, not contract; re-derived at the
branch test). **Assessment: most likely N = 1 regardless of OQ1**,
but left **OPEN** because it is downstream of OQ1 and the branch
test runs at the post-input resolving drafting, per the session
scope bound.

## Plan structure handoff

- Replace the parent-promotion stub at
  [`../t3-doc-declared-stages.md`](../t3-doc-declared-stages.md)
  with a full task plan at `Status: In draft` per
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "Required and optional sections": Status, context preamble,
  Goal, Contracts (written **conditionally** where they depend on
  OQ1–OQ5/OQ7), Files to touch (estimate-prefaced), Validation
  Gate.
- Optional sections that apply: Cross-Cutting Invariants
  (inherited m2 invariants — additive-spec, posture-tension,
  stub-render-preserved, walk-on-every-request, file boundary),
  Naming (conditional — the new frontmatter field + any new
  `PlanNode` field, spelled only once OQ1/OQ2 resolve), Self-Review
  Audits (`validation-honesty` + general checklist), Risk Register
  (additive-breaking-a-vendored-consumer; homogenization pressure;
  shape churn if OQ1 reopens), Out of Scope (t4's roster; on-box
  actor icons; AI prose extraction), Backlog Impact (none — state
  explicitly; reference the vision §7 prose-to-data and sub-stage-
  cell open questions as deliberated intersections, not resolved),
  Related Docs.
- **Same change** updates the m2 README parent doc for currency:
  the t3 Task Status row off "In draft (stub)" and the t3 prose to
  reflect a drafted `In draft` plan (NOT `Proposed`); the
  "Doc-declared-stages frontmatter shape (decide when t3 drafts)"
  Cross-Task Decisions entry points at this scoping section and is
  marked **OPEN pending human input** (NOT resolved).
- Status stays `In draft`; no promotion gate is run and the plan
  is not promoted to `Proposed` — the open questions above are the
  blockers, surfaced for the human.

## Reality-check inputs (the plan must re-verify before any future promotion)

Load-bearing claims this scoping rests on; a future resolving
drafting re-confirms each against current code before any
promotion (this session does not promote).

- **Frontmatter read surface.**
  [`parsedDoc` / `parsePlanDoc` / `stringList` in
  walker.go](../../../../../internal/site/walker.go) — the
  goldmark-meta read where the new optional field is added,
  alongside `slug` / `Status` / `short_description` /
  `related_prs`; `stringList` is the tolerant block-sequence
  decoder OQ1 A1/A2 reuse.
- **Render host.**
  [`forest.go`](../../../../../internal/site/forest.go) — the
  recursive `node` template (`node-header`, `node-detail`, the
  collapsible/leaf box) is t2's Landed surface and the host the
  progress-box row renders inside; `forest-style` is where any
  box-row CSS lives. t3 stays in the forest region —
  [`render.go`](../../../../../internal/site/render.go) shell and
  [`roster.go`](../../../../../internal/site/roster.go) are
  off-limits per the m2 file-boundary invariant.
- **Render-input struct.**
  [`PlanNode` in tree.go](../../../../../internal/site/tree.go) —
  the struct threaded to the template; the `ActiveInSubtree`
  additive-field + post-order pass is the precedent for any
  additive field a chosen OQ1 shape needs (e.g. a parsed
  stage list carried to the template).
- **Walk-on-every-request.**
  [`Server.index` in site.go](../../../../../internal/site/site.go)
  — walks plans + loads work-instances per request; t3 touches
  only the parser and the forest template/builder, introducing no
  cache/goroutine.
- **Node levels.**
  [`internal/slugs/slugs.go`](../../../../../internal/slugs/slugs.go)
  — `root` / `milestone` / `task` / `phase` (no `epic` node; the
  root is the top box); the box row renders at every level.
- **Additive precedent + exact-match discipline.**
  [`shared.md` "Optional `short_description`
  field"](../../../../../spec/planning/shared.md), "Optional
  `related_prs` field", and "Quote labels whose enforcement
  depends on exact-match matching" — the optional/additive home
  the new field documents adjacent to, and the naming discipline
  OQ2 is bound by.
- **Validation commands + UI capture.**
  [`docs/dev.md`](../../../../../docs/dev.md) — `go build ./...`,
  `go vet ./...`, `go test ./...`; UI capture is manual
  screenshots at ~1280px against a real `go run` with a populated
  `DB_PATH` and the `docs/plans/` tree (no screenshot tooling).
