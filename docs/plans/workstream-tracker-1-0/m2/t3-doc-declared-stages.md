---
slug: workstream-tracker-1-0-m2-t3
Status: In draft
short_description: Additive spec field + parser + per-node progress-cell render driven by the doc
---

# Task 3 — Doc-declared progress stages (spec-first)

## Status

`In draft` — **design decisions resolved by human input
(2026-05-18); promotion gate not yet run; not promoted.** A
spawned just-in-time planning session decomposed the genuine
decisions into shapes and surfaced them; the human resolved them
in-loop, aligned with the recommendations. The Contracts below are
now **locked** to those choices (no longer conditional). The plan
stays `Status: In draft` because the spawned session's scope bound
stops at `In draft` and does **not** run the
[`task-plan.md`](../../../../spec/planning/task-plan.md)
`` `In draft` → `Proposed` `` promotion gate; this plan makes no
claim that the gate has run. A resolving drafting session re-reads
the plan + scoping end-to-end, walks the universal `Verified by:`
rule, re-confirms the reality-check inputs against current code,
sketches the D6 branch test, and only then flips to `Proposed`.

Resolved decisions (full decomposition, rejected shapes, and
`Verified by:` grounding in
[`scoping/t3-doc-declared-stages.md`](scoping/t3-doc-declared-stages.md)
"Decisions resolved by human input"):

- **D1 (central, was OQ1)** — field shape: **A2**, an ordered
  list of stage-label strings, Drafting cell reserved render-side.
  One rendered cell ≈ one PR in the typical case. Per-stage counts
  (B2) deferred as an additive-linear future migration, not
  designed out.
- **D2 (was OQ2)** — frontmatter key **`progress_stages`**;
  rendered element named a **"progress cell"** (reconciles m2
  "progress box" toward the established `design/vision.md` §7
  "cell" vocabulary); no doc-visible Drafting token. The
  element-name spelling is the one judgment call made under the
  "make the reasonable call and continue" instruction — overridable
  by the human without reopening D1.
- **D3 (was OQ3)** — Drafting cell is render-side-reserved (one
  synthesis path serves the absent/stub case and a drafted doc's
  leading cell).
- **D4 (was OQ4)** — per-level applicability: **no inheritance**;
  each node's own doc governs its row.
- **D5 (was OQ7)** — row gated by **field-presence (α)**,
  Status-independent.
- **D6 (was OQ5)** — phase split assessed **N = 1** under D1 = A2;
  the branch-test sketch and the promotion gate are the resolving
  drafting session's work, not run here.

## Context

This plan covers a small, additive extension to the plan-doc
spec: a way for a plan doc to declare, in its own frontmatter, the
progress stages it expects to move through — and the
visualization rendering those stages as a row of progress cells on
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
the forest renders a row of progress cells whose count and order
come from that node's doc. A doc that omits the field — including
the already-supported `slug` + `Status: In draft` stub — renders
the Drafting cell only, with no error, skip, or broken layout. The
field is optional and additive: pre-existing docs and vendored
spec consumers are unaffected by its absence, and t4's roster
deliberately does not consume or schematize it.

Verifiable when: rendering the dogfood `docs/plans/` tree shows,
on every node box, a progress-cell row whose count and order match
that node's declared stages; a stub and any field-omitting doc
render exactly the Drafting cell as an intentional observed state
(not an empty or errored row); the existing actor markers, Status
badge, long description, and related-PR list still render
unchanged; and the walk-on-every-request render path is unchanged.

## Contracts

Final WHAT shape — **locked to the human-resolved decisions
D1–D5** (2026-05-18); no clause remains conditional. HOW
grounding and the decomposition with rejected shapes live in
[`scoping/t3-doc-declared-stages.md`](scoping/t3-doc-declared-stages.md)
"Decisions resolved by human input."

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
- **C2 — Field shape and name (locked: D1 = A2, D2).** The
  optional frontmatter field is named `progress_stages` and is an
  ordered YAML block sequence of stage-label strings, read by the
  existing tolerant block-sequence decoder. **Cell count = number
  of entries; order = sequence order; the label of each cell is
  its entry string** — so box count and order are derived from the
  doc, as the locked m2 WHAT requires. One rendered cell
  corresponds to one PR in the typical case. The shape stays a
  flat string sequence; growing a per-entry count later (the B2
  migration) is additive and out of scope here (Risk Register /
  Out of Scope). `Verified by:`
  [m2 README "Task Contracts" t3 row](README.md) (count + order
  come from the doc — the locked WHAT this realizes);
  [`stringList` in
  walker.go](../../../../internal/site/walker.go) (the exact
  tolerant block-sequence-of-strings decoder this reuses, as
  `related_prs` does);
  [`shared.md` "Decompose options into shapes before analyzing"
  + "Quote labels whose enforcement depends on exact-match
  matching"](../../../../spec/planning/shared.md) (why the shape
  was decomposed not invented, and the exact-match discipline the
  `progress_stages` token is copied verbatim under).
- **C3 — Drafting cell for the absent/stub case (locked: D3).** A
  doc that omits the field — including a `slug` + `Status: In
  draft` stub — renders exactly the Drafting cell, as an
  intentional observed state, never an empty or errored row. The
  Drafting cell is a single render-side synthesis serving both the
  absent/stub case and the leading cell of a drafted doc; there is
  **no doc-visible `Drafting` token** and the declared list holds
  only post-drafting stages. `Verified by:`
  [`parsePlanDoc` in
  walker.go](../../../../internal/site/walker.go) (a stub omits
  the field; the absent read returns the zero value, so the
  Drafting cell cannot derive from the field for a stub);
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
- **C5 — Render-gating discriminator (locked: D5 = α).** A node
  renders its declared-stage cells iff the `progress_stages` field
  is present; absent ⇒ Drafting cell only. The discriminator is
  field presence, **Status-independent**, with no coupling to the
  exact-match Status lifecycle token (so it never interacts with
  the "Unknown Status values render gracefully" rule). `Verified
  by:`
  [`parsePlanDoc` in
  walker.go](../../../../internal/site/walker.go) (Status and the
  new field are independent reads today; α adds no coupling, β
  adds one);
  [`shared.md` "Plan-doc Status" / "Unknown Status values render
  gracefully"](../../../../spec/planning/shared.md) (the
  exact-match Status surface β would couple to).
- **C6 — Per-level applicability (locked: D4 = no inheritance).**
  Each node's own doc governs its row; a node whose own doc omits
  the field renders the Drafting cell only even if an ancestor
  declares stages. No render-time parent/child stage-resolution
  walk is introduced. `Verified by:`
  [`shared.md` "Optional `short_description` field" / "Optional
  `related_prs` field"](../../../../spec/planning/shared.md)
  (the per-doc, absence-benign, no-inheritance precedent);
  [`buildTree` / `PlanNode` in
  tree.go](../../../../internal/site/tree.go) (per-doc fields are
  attached per node today; no field inherits across parent/child).
- **C7 — Render-side tests are semantic/structural.** The
  progress-cell row is asserted by structure and presence (the row
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
  falling back to the Drafting-cell-only shape, mirroring the
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
  render exactly the Drafting cell and creates no work-instance and
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

Locked per D2. The new optional frontmatter key is
**`progress_stages`** (an ordered block sequence of stage-label
strings, D1 = A2), and the parsed list is carried to the template
on one additive `PlanNode` field (precedent: t2's
`ActiveInSubtree` additive field; exact Go identifier chosen at
implementation, internal to the `site` package). The rendered
element is a **"progress cell"** (reconciles the m2 "progress
box" wording toward `design/vision.md` §7's established "cell"
vocabulary). `progress_stages` is an exact-match-checked token per
[`shared.md` "Quote labels whose enforcement depends on exact-match
matching"](../../../../spec/planning/shared.md) and is copied
verbatim into the spec, not paraphrased. The element-name spelling
("progress cell") is the one judgment call the resolving turn made
under the "make the reasonable call and continue" instruction and
is human-overridable without reopening D1. `Verified by:`
[`PlanNode` in tree.go](../../../../internal/site/tree.go)
(the additive-field precedent);
[`shared.md` exact-match-label
discipline](../../../../spec/planning/shared.md).

## Files to touch

_Estimate of the expected file shape, not a binding rule.
Implementation may revise this list when a structural call
requires it — including touching a file listed under
"Intentionally not touched." Any deviation is handled via the
PR-body Estimate Deviations callout with the plan reconciled to
what shipped._

- **Modify:** the plan-tree frontmatter parser (the new
  `progress_stages` read — reuses the existing tolerant
  block-sequence decoder, as `related_prs` does, per D1 = A2);
  the forest-region render template + styles (the progress-cell
  row partial inside the existing per-node box, plus its CSS);
  the render-input node struct (one additive field carrying the
  parsed stage list to the template — precedent: t2's
  `ActiveInSubtree`).
- **Modify (spec):** the plan-doc spec's optional/additive
  section — `progress_stages` documented adjacent to
  `short_description` / `related_prs`, optional + additive + the
  defined absent behavior (Drafting cell only). This is the
  spec-first deliverable.
- **Modify (tests):** the forest-region test file (semantic
  assertions per C7: row count/order for a declaring doc;
  Drafting-only for a field-omitting doc and a stub; absent /
  malformed field never errors); the parser test file (the new
  field's tolerant read, absence/wrong-type/partial-malformed).
- **Modify (docs):** the m2 README parent doc — Task Status t3
  row and t3 prose to a drafted `In draft` plan, and the
  "Doc-declared-stages frontmatter shape" Cross-Task Decisions
  entry pointed at the scoping section and recorded RESOLVED by
  human input (D1–D5), promotion gate still pending (this drafting
  change, per parent-doc currency); `design/v0.1-design.md`
  §7 updated to describe the doc-driven progress-cell row on the PR
  that lands the implementation.
- **Intentionally not touched:** the shell render file
  (`render.go`) and the roster render file (`roster.go`) — m2
  file boundary; the API, DB schema, register client, routing;
  any dependency (`progress_stages` reuses the existing
  goldmark-meta / `yaml.v2` block-sequence decode per D1 = A2 — no
  new dependency).

## Validation Gate

The design decisions are locked; the resolving drafting session
runs the promotion-gate self-review and finalizes the flip (this
session does neither). The gate below is the implementing PR's
Validation Gate.

- `go build ./...`, `go vet ./...`, `go test ./...` all pass (per
  [`docs/dev.md`](../../../../docs/dev.md)).
- Unit tests cover: a declaring doc renders a progress-cell row of
  the declared count and order; a field-omitting doc and a `slug`
  + `Status: In draft` stub render exactly the Drafting cell; an
  absent / wrong-typed / partially-malformed field never errors or
  drops the node (the additive guarantee); the preserved t2
  surfaces (Status badge, actor markers, long description,
  related-PR list, collapse) still render on every box.
- **Manual UI capture** (no screenshot tooling exists; capture is
  manual per [`docs/dev.md`](../../../../docs/dev.md) UI Review):
  screenshots at ~1280px against a real `go run` reading a
  populated `DB_PATH` and the `docs/plans/` dogfood tree —
  confirming the progress-cell row on declaring nodes and the
  Drafting-cell-only render on the stub.
- **Drafting-cell-only is observed, not assumed.** Per
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
  absent-behavior (Drafting cell only), mirroring the landed
  `short_description` / `related_prs` precedent. `Verified by:`
  [`shared.md` "Optional `related_prs`
  field"](../../../../spec/planning/shared.md).
- **The two spec postures get homogenized under review pressure.**
  A reviewer or later task may "tidy" this field toward t4's
  schema-loose JSON or vice versa. Mitigation: the posture-tension
  cross-cutting invariant names reconciliation as the defect.
- **The one-cell-per-PR assumption breaks (a step needs >1 PR).**
  The human noted this may happen someday; D1 = A2 does not model
  per-step counts. Re-trigger axis: a node whose work genuinely
  fans out into multiple PRs per declared stage. Migration cost:
  **linear/additive, not a cliff** — the flat string sequence can
  later become a sequence of {label, count} records (D1 shape B2)
  while docs that authored the string form keep working, because
  the tolerant decoder drops non-conforming entries rather than
  erroring. Mitigation: B2 is recorded as the named additive
  escape hatch (Out of Scope below + scoping D1); no pre-emptive
  structure is added now, keeping the deferral cheap. `Verified
  by:`
  [`stringList` in walker.go](../../../../internal/site/walker.go)
  (drop-don't-error tolerance that makes the later record form
  additive);
  [`design/vision.md` §7](../../../../design/vision.md) (the
  per-phase-PR-count axis B2 would serve).

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
- **Per-step PR counts (D1 shape B2).** Declaring a count per
  stage (so one cell = one PR within a multi-PR stage) is the
  named additive future migration if the one-cell-per-PR
  assumption breaks; t3 ships the flat string sequence only and
  adds no per-entry count now.
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
placeholder vs. only-what-committed sub-stage cells (adjacent to
D5; surfaced, not decided here).

## Related Docs

- [`README.md`](README.md) — parent milestone; its Task Status t3
  row, t3 prose, and "Doc-declared-stages frontmatter shape"
  Cross-Task Decisions entry are reconciled to this drafting in
  the same change (pointed at scoping, recorded RESOLVED by human
  input D1–D5, promotion gate still pending).
- [`scoping/t3-doc-declared-stages.md`](scoping/t3-doc-declared-stages.md)
  — the decision-space decomposition (D1–D6 with shapes,
  trade-offs, `Verified by:` grounding, rejected alternatives, and
  the human resolutions) this plan's Contracts realize.
- [`t2-expanded-render.md`](t2-expanded-render.md) and
  [`scoping/t2-expanded-render.md`](scoping/t2-expanded-render.md)
  — the Landed sibling whose per-node box hosts the progress-cell
  row, and the semantic-not-byte-exact test precedent (C7).
- [`../../stub-children-on-parent-promotion/README.md`](../../stub-children-on-parent-promotion/README.md)
  — the landed stub render case the Drafting-cell-only behavior
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
