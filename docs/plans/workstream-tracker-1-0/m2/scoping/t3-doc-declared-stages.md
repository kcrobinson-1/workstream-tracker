# Scoping — m2 t3: Doc-declared progress stages (spec-first)

Transient deliberation doc for the
[`workstream-tracker-1-0-m2-t3`](../t3-doc-declared-stages.md) task
plan. Deletes in batch with sibling scoping docs at the
milestone-terminal PR per
[`task-plan.md`](../../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns." No Status field — scoping docs are
not part of the plan-doc lifecycle.

> **History of this doc.** A spawned just-in-time drafting session
> decomposed the genuine decisions into shapes and surfaced them
> OPEN (the original endpoint was a task plan at `Status: In
> draft`). The contributor then resolved D1–D6 in-loop
> (2026-05-18) and directed the session to walk the
> `` `In draft` → `Proposed` `` promotion gate in-session,
> consciously extending past the spawn's original "stop at
> `In draft`" bound at the contributor's explicit direction. The
> "Decisions made at scoping time" section records the
> scoping-method calls; "Decisions resolved by human input"
> records D1–D6 with rejected shapes; the gate was walked and the
> paired plan is now `Proposed`. No PR was opened by the session
> (separately out of scope).

## Context summary

t3 adds an additive, optional `spec/` frontmatter affordance that
lets a plan doc declare its own progress stages, teaches the
plan-tree walker to read it, and renders a row of progress cells
(D2's rendered-element name; the m2 milestone's umbrella term
stays "progress boxes") — one row per node at every level (root,
milestone, task, phase) — whose declared cell count and order come
from the doc. A doc that omits the field (the already-supported
`slug` + `Status: In draft` stub) renders only the reserved
Drafting cell and is never errored or skipped. The cell row
renders inside the per-node box t2 (Landed) established in the
forest region; t3 does not touch the shell or the roster, and
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

### S2 — The field-shape decision is decomposed here; the human resolved it 2026-05-18

The m2 milestone routed the field name/structure/per-stage-counts-
vs-ordered-list call to "when t3 drafts" under the
[`shared.md`](../../../../../spec/planning/shared.md) "Decompose
options into shapes" and exact-match-label discipline. This
spawned session decomposed it (D1–D6 below) with trade-off
analysis; the human resolved it in-loop on 2026-05-18 (aligned
with the recommendations) and then directed the session to walk
the `` `In draft` → `Proposed` `` promotion gate, consciously
extending past the spawn's original "stop at `In draft`" bound at
the contributor's explicit direction. Recording the decomposition
is the scoping-method act; the choices are the human's; the gate
was walked in-session and the plan is now `Proposed`. `Verified
by:` [m2 README "Cross-Task Decisions" →
"Doc-declared-stages frontmatter shape"](../README.md);
[`shared.md` "Decompose options into shapes before
analyzing"](../../../../../spec/planning/shared.md).

## Decisions resolved by human input (2026-05-18)

Each decision was decomposed into shapes with trade-offs grounded
in cited code/spec, then resolved by the human in-loop (aligned
with the recommendations). The verdicts below are the durable
rationale the paired plan's Contracts realize; rejected shapes are
retained as the decision record. The plan's Contracts are now
locked to these choices. After resolving them the contributor
directed the session to walk the
`` `In draft` → `Proposed` `` promotion gate in-session
(consciously extending past the spawn's original "stop at
`In draft`" bound, at the contributor's explicit direction); the
gate was walked and the plan is now `Proposed`.

**Conceptual model the human set (load-bearing for D1).** A
rendered progress unit corresponds to **one PR in the typical
case**. A node's declared list is therefore the ordered sequence
of shippable steps (≈ PRs) it expects, each with a short label;
the typical single-PR phase declares a length-1 list. The human
explicitly noted the one-box-per-PR assumption "may break"
someday — i.e. a step that needs more than one PR. That case is
exactly D1 shape B2 and is **deferred, not designed-out**; the
A2 → B2 migration is **additive and linear, not a cliff**: a
sequence of label strings can later gain an optional per-entry
count by becoming a sequence of {label, count} records without
breaking docs that authored the string form (the tolerant decoder
already drops non-conforming entries rather than erroring). This
keeps the deferral cheap and is recorded as a Risk Register entry
and an Out-of-Scope note in the plan. `Verified by:`
[`stringList` in walker.go](../../../../../internal/site/walker.go)
(the tolerant decoder whose drop-don't-error posture makes the
later record form additive);
[`design/vision.md` §7](../../../../../design/vision.md) (the
per-phase-PR-count axis B2 would later serve).

### D1 (was OQ1, central) — Field shape — RESOLVED: A2

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
cannot answer: **does one progress unit represent a stage, or a PR
within a stage?** That is exactly the vision §7 open question
("structured fields for the load-bearing numbers" vs. softer
signal) and is the per-stage-counts-vs-ordered-list axis the
milestone deferred here. **RESOLVED by the human (2026-05-18):
A2** — an ordered list of stage-label strings, Drafting reserved
render-side, reusing the `stringList` precedent; one rendered
unit ≈ one PR in the typical case (see the conceptual-model note
above). **Rejected:** A1 (no reserved Drafting box — forces the
absent/stub case into a second code path); A3 (doc-visible
`Drafting` magic token + two code paths); B1 (YAML map loses
order); **B2 (per-stage counts) deferred — not designed out**, it
is the additive future migration if the one-PR-per-unit
assumption breaks; C (record richness no current clause
justifies).

### D2 (was OQ2) — Field name and element name — RESOLVED

The field name, and any reserved stage token, become exact-match-
checked identifiers the parser and consumers key off, governed by
[`shared.md` "Quote labels whose enforcement depends on exact-match
matching"](../../../../../spec/planning/shared.md). Candidate
field names: `progress_stages`, `stages`, `declared_stages`,
`progress`. Candidate Drafting treatment: a render-side reserved
box with **no doc-visible token** (preferred — keeps the magic
string out of authored docs and out of the exact-match surface)
vs. a literal `Drafting` token (couples to D1 shape A3).
`Verified by:`
[`shared.md` "Quote labels whose enforcement depends on exact-match
matching"](../../../../../spec/planning/shared.md) (the discipline
this naming is bound by);
[`parsePlanDoc` in walker.go](../../../../../internal/site/walker.go)
(the frontmatter key the chosen name is read at, alongside the
existing `slug` / `Status` / `short_description` reads). The
milestone says the token is "decided here … not invented at
milestone level." **RESOLVED (2026-05-18):** frontmatter key
**`progress_stages`**; the rendered element is named a **"progress
cell"**, not "progress box" — the human flagged "progress box" as
needing a better name, and "cell" is the established vision
vocabulary ([`design/vision.md`
§7](../../../../../design/vision.md) uses "sub-stage cells" /
"cell count"), so "progress cell" reconciles the m2 wording toward
an existing term rather than inventing one and reads correctly at
PR granularity. Drafting: render-side reserved, **no doc-visible
token**. The element-name call is the one judgment call the
resolving turn made under the "make the reasonable call and
continue" instruction; it is overridable without reopening D1.
**Rejected:** field keys `stages` (overloaded), `declared_stages`,
`progress` (vague); retaining "progress box" (the human rejected
it); a literal `Drafting` token (couples to the rejected D1 shape
A3).

### D3 (was OQ3) — Drafting box render-side-reserved — RESOLVED

Tied to D1 A2 vs. A3. The contract requires a field-omitting stub
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
only). **RESOLVED (2026-05-18): render-side reserved Drafting
box** — declared list = post-drafting stages, one synthesis path
serves both the absent/stub case and a drafted doc's leading cell;
no doc-visible token (resolved jointly with D1/D2).

### D4 (was OQ4) — Per-level applicability: no inheritance — RESOLVED

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
**RESOLVED (2026-05-18): no inheritance** — each node's own doc
governs its row; absence ⇒ Drafting cell only, even under a
declaring ancestor. **Rejected:** inheritance (adds a render-time
parent/child walk; breaks the per-doc absence-benign precedent).

### D5 (was OQ7) — Row gated by field-presence (α) — RESOLVED

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
questions"](../../../../../design/vision.md). **RESOLVED
(2026-05-18): α** — field-presence gated, Status-independent, no
coupling to the exact-match Status lifecycle token. **Rejected:**
β (Status-gated; adds a Status-token coupling into the progress
render that must also satisfy "Unknown Status values render
gracefully").

### D6 (was OQ5) — Phase split — RESOLVED N = 1 (branch-test sketch at the in-session gate)

The m2 "Per-task phase splits" estimate guesses t3 plausibly N ≥ 2
(spec/parser, then progress-box render). Candidate boundaries
enumerated: (a) spec-field + parser read | progress-box render;
(b) single phase. The resolver is the
[`task-plan.md`](../../../../../spec/planning/task-plan.md)
"PR-count predictions need a branch test," run at the resolving
drafting after human input. **RESOLVED at the in-session gate
walk: N = 1.** The split hinges on D1: under A2 the substantive
logic is tiny (one tolerant field read reusing `stringList`, plus
a render loop in the forest node template), pointing strongly to
N = 1, single subsystem (`internal/site`) plus additive spec-doc
prose; under the deferred B2 the
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
branch test). With D1 = A2 the substantive logic is tiny (one
tolerant block-sequence read reusing `stringList`, plus a render
loop in the forest node template), so **RESOLVED: N = 1**, one
subsystem (`internal/site`) plus additive spec-doc prose, well
under the >5-subsystem / >300-LOC thresholds. The branch-test
sketch was run at the in-session promotion-gate walk; the plan's
Status section carries the sketched file list.

## Plan structure handoff

- Replace the parent-promotion stub at
  [`../t3-doc-declared-stages.md`](../t3-doc-declared-stages.md)
  with a full task plan (now `Status: Proposed` after the
  in-session gate walk) per
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "Required and optional sections": Status, context preamble,
  Goal, Contracts (locked to D1–D5; D6 phase split resolved
  N = 1 by the branch-test sketch), Files to touch
  (estimate-prefaced), Validation Gate.
- Optional sections that apply: Cross-Cutting Invariants
  (inherited m2 invariants — additive-spec, posture-tension,
  stub-render-preserved, walk-on-every-request, file boundary),
  Naming (locked per D2 — frontmatter key `progress_stages`,
  rendered element "progress cell" + any new `PlanNode` field),
  Self-Review
  Audits (`validation-honesty` + general checklist),
  Documentation Currency (status-bearing docs the implementing PR
  and the drafting/gate change touch — design §7, m2 README),
  Risk Register
  (additive-breaking-a-vendored-consumer; homogenization pressure;
  the box≈PR assumption breaking → B2 as the additive-linear
  migration), Out of Scope (t4's roster; on-box
  actor icons; AI prose extraction), Backlog Impact (none — state
  explicitly; reference the vision §7 prose-to-data and sub-stage-
  cell open questions as deliberated intersections, not resolved),
  Related Docs.
- The m2 README parent doc was updated for currency across the
  drafting and the gate flip: the t3 Task Status row
  `In draft (stub)` → `In draft` → `Proposed`, the t3 prose, and
  the "Doc-declared-stages frontmatter shape" Cross-Task Decisions
  entry (decomposed → RESOLVED by human input D1–D6 → gate
  walked).
- The plan is now `Proposed`: the contributor resolved the
  decisions in-loop and directed the in-session promotion-gate
  walk (end-to-end coherence, decision-completeness, universal
  `Verified by:`, reality-check re-confirmation, always-on rules,
  D6 branch-test sketch, then the flip), consciously extending
  past the spawn's original "stop at `In draft`" bound at the
  contributor's explicit direction. No PR is opened by this
  session (separately out of scope); the flip is a doc-only
  commit on the worktree branch.

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
  decoder D1 (A2) reuses.
- **Render host.**
  [`forest.go`](../../../../../internal/site/forest.go) — the
  recursive `node` template (`node-header`, `node-detail`, the
  collapsible/leaf box) is t2's Landed surface and the host the
  progress-cell row renders inside; `forest-style` is where any
  cell-row CSS lives. t3 stays in the forest region —
  [`render.go`](../../../../../internal/site/render.go) shell and
  [`roster.go`](../../../../../internal/site/roster.go) are
  off-limits per the m2 file-boundary invariant.
- **Render-input struct.**
  [`PlanNode` in tree.go](../../../../../internal/site/tree.go) —
  the struct threaded to the template; the `ActiveInSubtree`
  additive-field + post-order pass is the precedent for any
  additive field the D1 = A2 shape needs (the parsed
  stage list carried to the template).
- **Walk-on-every-request.**
  [`Server.index` in site.go](../../../../../internal/site/site.go)
  — walks plans + loads work-instances per request; t3 touches
  only the parser and the forest template/builder, introducing no
  cache/goroutine.
- **Node levels.**
  [`internal/slugs/slugs.go`](../../../../../internal/slugs/slugs.go)
  — `root` / `milestone` / `task` / `phase` (no `epic` node; the
  root is the top box); the cell row renders at every level.
- **Additive precedent + exact-match discipline.**
  [`shared.md` "Optional `short_description`
  field"](../../../../../spec/planning/shared.md), "Optional
  `related_prs` field", and "Quote labels whose enforcement
  depends on exact-match matching" — the optional/additive home
  `progress_stages` documents adjacent to, and the naming
  discipline D2 is bound by.
- **Validation commands + UI capture.**
  [`docs/dev.md`](../../../../../docs/dev.md) — `go build ./...`,
  `go vet ./...`, `go test ./...`; UI capture is manual
  screenshots at ~1280px against a real `go run` with a populated
  `DB_PATH` and the `docs/plans/` tree (no screenshot tooling).
