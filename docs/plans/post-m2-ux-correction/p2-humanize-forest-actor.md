---
slug: post-m2-ux-correction-p2
Status: Proposed
short_description: Humanize the forest actor — reported name (slug fallback) in the per-node actor-marker (F4)
---

# Phase 2 — Humanize forest actor (F4)

## Status

`Proposed`. The OD walk resolved all five open decisions
(**OD1 = OD1.a** — additive `Name` on `ActiveWorkInstance`
populated in `Server.index` between `loadSessionMetadata` and
`buildTree`; **OD2 = OD2.a** — additive `Slug` on
`ActiveWorkInstance` alongside `Name`; **OD3 = OD3.a** —
unbound stays roster-only; **OD4 = OD4.a** — semantic
presence-and-absence test assertions; **OD5 dissolved** —
subsumed by OD1.a) and the
[`task-plan.md`](../../../spec/planning/task-plan.md)
`` `In draft` → `Proposed` `` promotion gate was walked before
the flip. Drafting, the OD walk, and the gate walk all ran in
a single delegated drafting session; the gate-walk extended
the session past the standing spawned-drafting "stop at `In
draft`" bound at the contributor's explicit direction
(matching the t3 precedent recorded in
[`m2/README.md` Task Status](../workstream-tracker-1-0/m2/README.md)
— "extending past the spawn's original 'stop at `In draft`'
bound at the contributor's explicit direction"). Status flip
co-locates with the gate-walk record below in this commit.

### Gate-walk history

- **End-to-end coherence** — plan re-read in order; no
  contradiction between Goal, Contract **C1**, the inherited
  Cross-Cutting Invariants (C-INV-2 / C-INV-3 / C-INV-4 /
  C-INV-5; C-INV-1 explicitly out of this phase's surface),
  Files-to-touch, the per-phase Validation Gate, Self-Review
  Audits, Out of Scope, the Risk Register, and the resolved
  OD entries. C1 is the only contract clause; it locks the
  WHAT the parent task plan's Contract **C6** specifies for
  the forest-region surface.
- **Decision-completeness on Contracts** — every OD resolved
  with rationale (OD1, OD2, OD3, OD4 locked with vision- or
  code-grounded reasoning; OD5 dissolved by OD1.a). No
  "decided at plan-drafting," "shape later," or "spelling at
  plan time" phrasing in C1. The implementing PR's narrow
  render-altitude calls (exact span text, dark/light-mode
  observation) are explicitly authorized by
  [`shared.md`](../../../spec/planning/shared.md) "Plans
  describe contracts, not implementation."
- **Universal `Verified by:` walk** — every load-bearing claim
  in C1, the resolved ODs, the Cross-Cutting Invariants,
  Files-to-touch, and the Risk Register carries a code
  citation or a parent-doc / Landed-sibling-doc reference.
  Symbolic anchors (`node-header`, `actor-marker`,
  `ActiveWorkInstance`, `Server.index`, `loadSessionMetadata`,
  `buildTree`, `buildRoster`, `RosterEntry`) dominate per the
  anchor-preference rule.
- **Reality-check inputs re-confirmed** — every cited line
  range spot-checked against current branch code:
  [`forest.go:51`](../../../internal/site/forest.go) (the
  `node-header` `actor-marker` span renders `{{.Actor}}`
  today, the surface this phase edits);
  [`tree.go:39`](../../../internal/site/tree.go)
  (`ActiveWorkInstance` carries `ID` + `Actor`);
  [`tree.go:104`](../../../internal/site/tree.go)
  (`buildTree` signature `docs, active`);
  [`tree.go:148`](../../../internal/site/tree.go)
  (`WorkInstances: active[d.Slug]` exact-slug join — the
  attached/unbound discriminator);
  [`site.go`](../../../internal/site/site.go)
  (`Server.index` orchestration; `loadActiveWorkInstances` /
  `loadSessionMetadata` / `buildTree` / `buildRoster` ordering;
  `RosterEntry` shape; `buildRoster`'s
  pure-function-over-already-loaded-values precedent);
  [`forest_test.go`](../../../internal/site/forest_test.go)
  (file exists; semantic-test idiom matches the OD4.a
  posture). No drift.
- **Always-on rules** — required sections present (Status,
  Context, Goal, Contracts, Files to touch, Validation Gate);
  optional sections present where applicable (Cross-Cutting
  Invariants inherited from parent; Self-Review Audits; Out
  of Scope; Risk Register; Related Docs); no `Phase Contracts`
  because p2 is a phase plan with no sub-children. No descent
  to implementation prescription — OD1.a's data-flow shape
  (additive struct field + `Server.index` threading point) is
  contract-altitude, explicitly authorized by the milestone's
  t4-amended data-path carve-out
  ([`m2/README.md` Cross-Task Invariants](../workstream-tracker-1-0/m2/README.md))
  and by parent C-INV-2. No soft-commitment language. The
  `### Open decisions` sub-block stays in `## Status` per the
  parent README's OD-walk-outcomes precedent and is not a
  [`task-plan.md`](../../../spec/planning/task-plan.md)
  "Required and optional sections" variance.
- **Phase-skeleton seeding (N ≥ 2 only)** — does not apply.
  p2 is a phase plan with no further phase children; the
  parent task plan is N = 3 (p1 / p2 / p3) and seeded its
  three phase skeletons when *it* was promoted.

### Open decisions

Each was decomposed against merged code; resolutions fold into
Contract **C1**, Files-to-touch, the Validation Gate, and the
implementing PR. Resolved entries lead with **Resolved =** and
fold the chosen shape into the durable plan; the original
framing is retained underneath as the scoping record (parent
README precedent).

- **OD1 — Data-flow shape for the reported `name`. Resolved =
  OD1.a** (additive `Name` field on `ActiveWorkInstance`,
  populated in `Server.index` between `loadSessionMetadata` and
  `buildTree`; no `buildTree` signature change; the forest
  template reads `.Name` directly). *Rationale (vision-grounded,
  not precedent-grounded).* The vision-named forest extensions
  ahead — F3b's per-cell `current_stage`
  ([`progress-cell-active-state-and-actor`](../../backlog.md#progress-cell-active-state-and-actor)),
  the work-instance state vocabulary
  ([`vision.md` §3](../../../design/vision.md) `awaiting-user` /
  `awaiting-external` / `backgrounded`), actor kind for
  color/hover ([`vision.md` §3](../../../design/vision.md)),
  actor lineage ([`vision.md` §3 / §4](../../../design/vision.md))
  — are sourced from a *mix* of schema growth (additional
  `work_instances` columns) and reported metadata. OD1.a is the
  only shape that unifies both sources at one boundary:
  `Server.index` populates `ActiveWorkInstance` from
  `loadActiveWorkInstances` (schema-sourced fields) and from
  `loadSessionMetadata` (metadata-sourced fields), both onto the
  same value the forest renders. OD1.c's `sessionMeta`-only
  routing reaches only the metadata-sourced half; OD1.b ages
  worst as each new field needs another map on `indexData`.
  Beyond the future-feature pipeline, **F4 itself is the cost of
  having shipped with divergent forest/roster value-shapes** —
  `RosterEntry` carried `Name`, `ActiveWorkInstance` did not;
  that structural asymmetry is what allowed the rendering
  divergence to ship. OD1.a closes the asymmetry structurally
  (single mental model of "what the tool displays about a
  work-instance"; a future field added to one surface naturally
  pulls toward the other); OD1.b/c close the F4 rendering
  symptom while preserving the conditions that allowed it.
  Folded into Contract **C1** and the Files-to-touch certain set
  below. *Verified by:*
  [`Server.index` order in
  site.go](../../../internal/site/site.go) (docs → active → meta
  → buildTree → buildRoster — `meta` is already loaded before
  `buildTree` is called, so OD1.a populates the new field
  between those two steps with no reordering); [`ActiveWorkInstance`
  in tree.go](../../../internal/site/tree.go) (the existing `ID`
  field, added in t4-p2 for the roster's event-log join even
  though the forest doesn't render `ID`, sets the
  additive-field-for-cross-region-needs precedent OD1.a
  extends); [`buildRoster` in
  site.go](../../../internal/site/site.go) (the existing
  "pure-function-over-already-loaded-values" precedent OD1.a
  matches symmetrically).

  - **OD1 (original framing, retained as scoping record).** The
    forest's per-node `range .WorkInstances` iterates
    `*ActiveWorkInstance` values that carry `ID` and `Actor`
    only today ([`ActiveWorkInstance` in
    tree.go](../../../internal/site/tree.go); the `ID` field
    exists for the roster's event-log join, never rendered by
    the forest). The reported `name` is already resolved per
    request by
    [`loadSessionMetadata`](../../../internal/site/site.go) into
    a `map[string]sessionMeta` keyed by `wi.ID`, and consumed by
    `buildRoster` ([`buildRoster` in
    site.go](../../../internal/site/site.go)). Threading the
    resolved value to the forest's per-node template was
    decomposed against three plausible shapes:
    - **OD1.a — Additive `Name` field (and `Slug` if OD2 lands
      on OD2.a) on `ActiveWorkInstance`, populated in
      `Server.index` between `loadSessionMetadata` and
      `buildTree`.** No `buildTree` signature change. Template
      reads `{{.Name}}` directly.
    - **OD1.b — Parallel `map[string]string` (id → name) on
      `indexData`, looked up in the template via the `index`
      builtin (`{{index $.Names .ID}}`).** `ActiveWorkInstance`
      stays loader-key shape. New template surface (outer-scope
      `$` + `index` builtin, neither used today).
    - **OD1.c — Full `map[string]sessionMeta` through
      `indexData`, lookup via `{{(index $.Meta .ID).Name}}`.**
      Same template surface as OD1.b plus a `.Name` step. Routes
      more data than C1 needs; future-extensibility for further
      metadata-sourced reads.

    Tradeoff lens at decision time: future forest fields will be
    sourced from a *mix* of schema growth (state vocab, lineage,
    F3b `current_stage`, possibly kind) and reported metadata;
    (a)/(d-style) value-shape carries both at one boundary,
    (c)'s metadata-only routing reaches only one half. Symmetry
    with the roster's `RosterEntry.Name` is load-bearing because
    F4 itself is the F4-shaped bug — the structural asymmetry
    between `RosterEntry` (value-carries-display) and
    `ActiveWorkInstance` (loader-key-shape) is what allowed the
    forest/roster identity disagreement m2 shipped. All three
    shapes preserve C-INV-4 single-resolution.

- **OD2 — Slug fallback's source on the rendered value.
  Resolved = OD2.a** (the work-instance's registered slug; add
  a `Slug string` field to `ActiveWorkInstance` alongside the
  `Name` field OD1.a lands, populated in the same
  `Server.index` pass from the `active`-map key in scope).
  *Rationale.* Under OD1.a the forest's per-node template
  reads display fields off the work-instance value; OD2.a
  continues that convention for the slug fallback rather than
  splitting it across two scopes (`Name` on the value, `Slug`
  via outer-template `$node.Slug`). The rendered output is
  identical to OD2.b for every attached work-instance —
  `buildTree`'s exact-slug join
  ([`tree.go:148`](../../../internal/site/tree.go)
  `active[d.Slug]`) guarantees the work-instance's registered
  slug equals the enclosing node's slug — so OD2 is not
  load-bearing for behavior. The choice is posture-only, and
  OD2.a's posture matches `RosterEntry.Name` / `RosterEntry.Slug`
  ([site.go:286](../../../internal/site/site.go)): both
  surfaces read identity off the work-instance / roster-entry
  value, both fall back name-then-slug, both populate from the
  same per-request resolution. The future-feature pattern
  ahead (F3b `current_stage`, work-instance state vocab, actor
  kind, actor lineage) will all land on `ActiveWorkInstance`,
  not via outer-template scope, so OD2.a establishes the
  forest-template convention every future field follows
  without further negotiation. Name-then-actor-id was **not**
  a candidate — banned by parent **C6** + **C-INV-3**.

  - **OD2 (original framing, retained as scoping record).**
    Parent Contract **C6** locks the rule as **name-then-slug**,
    the same rule t4 locked for the roster
    ([`t4-session-roster.md` "Session identity and
    naming"](../workstream-tracker-1-0/m2/t4-session-roster.md)).
    Two shapes were decomposed for the *source* of the slug:
    - **OD2.a — The work-instance's registered slug.**
      Requires a `Slug` field on `ActiveWorkInstance`
      (additive — only OD1.a's shape; OD1.b/c carry the slug
      through the threaded map).
    - **OD2.b — The enclosing node's slug** (`.Slug` on the
      parent `PlanNode`, available via outer-template scope).
      No additive field needed; introduces a template-scope
      reference the forest's `node-header` doesn't use today.

    Tradeoff lens at decision time: identical rendered output
    for every attached work-instance (`buildTree`'s exact-slug
    join equates the two sources). OD2 is posture-only — the
    same axis as OD1's value-shape symmetry, applied to the
    fallback half.

- **OD3 — Unbound work-instance handling. Resolved = OD3.a**
  (no change; unbound stays roster-only). *Rationale
  (vision-grounded).*
  [`vision.md` §4](../../../design/vision.md) explicitly names
  the roster as the surface for unbound sessions: "set aside
  from the plan-tree forest is a roster of the sessions the
  tool knows about ... a session that is unbound ... still
  appears, in the roster, rather than vanishing."
  [§3](../../../design/vision.md) puts unattached work in a
  triage zone, not in the forest. The forest/roster split is a
  vision-level commitment, not a t4 implementation choice. The
  pre-existing `buildTree` exact-slug join
  ([`tree.go:104`](../../../internal/site/tree.go);
  `active[d.Slug]` reads only when the slug matches a parsed
  doc) already implements that split correctly. Parent Contract
  **C6** targets the forest's *per-node* `actor-marker` —
  there is no `actor-marker` to humanize on a node that doesn't
  exist — so OD3.a is also the scope-grounded reading. The
  triage *action* (promote-into-tree / dismiss) is deferred
  past 1.0 per [`vision.md` §4](../../../design/vision.md) and
  the parent epic's resolved "Triage zone in 1.0?" open
  question; OD3.a doesn't close that door — a future bound
  transition flows through `buildTree`'s existing join
  automatically. Recorded in `## Out of Scope` and observed in
  the Validation Gate's "Observe" step (seeded session 3, an
  unbound entry that appears in the roster but not in the
  forest).

  - **OD3 (original framing, retained as scoping record).** A
    work-instance registered against a slug not present in the
    walked tree is *unbound*. Today the forest drops it at
    `buildTree`'s join site
    ([`tree.go:104`](../../../internal/site/tree.go);
    `active[d.Slug]` reads only when the slug matches a parsed
    doc), and the roster lists it
    ([`buildRoster` in
    site.go](../../../internal/site/site.go) classifies unbound
    entries with `Bound: false`). Two shapes were decomposed:
    - **OD3.a — No change; unbound stays roster-only.** The
      drop predates p2 and falls naturally out of the
      exact-slug join.
    - **OD3.b — Render unbound work-instances somewhere in the
      forest.** Out of scope per the parent task plan's split
      between roster (which surfaces unbound) and forest (which
      renders attached); would require new render plumbing
      past F4 and contradicts the vision-level split.

- **OD4 — Test-coverage posture for the F4 contract. Resolved
  = OD4.a** (assert presence-and-absence semantically; the
  rendered `actor-marker` text contains the reported `name`
  (or the slug fallback when no `name` was reported) *and*
  does not contain the `wst-` actor prefix). *Rationale.* The
  `wst-<uuid>` actor is **not** being deprecated — it remains
  live, present on every `ActiveWorkInstance.Actor` field
  ([`loadActiveWorkInstances` in
  site.go](../../../internal/site/site.go)) and on every
  `work_instances.actor` schema column, used as the
  registration idempotency key on `(slug, actor)` (parent task
  plan `## Out of Scope` "no actor-generator change"). The
  absence assertion is therefore testing against **current
  live data**, not a deprecated relic pattern. It catches the
  specific silent-regression mode the
  [`vision.md` §9](../../../design/vision.md) feature risk
  names ("sessions go untracked ... without an obvious
  symptom") in this surface: a future template "tidy" like
  `{{.Name}} ({{.Actor}})` would pass OD4.b (the name is
  present) while re-introducing the uuid leak F4 corrects.
  OD4.a fails that case via the absence side. Posture matches
  m2 t3 C7's semantic-not-byte-exact assertions
  ([`t3-doc-declared-stages.md` Contracts](../workstream-tracker-1-0/m2/t3-doc-declared-stages.md));
  cost is one additional `NotContains(rendered, "wst-")` per
  case.

  - **OD4 (original framing, retained as scoping record).**
    The forest-region test surface
    ([`internal/site/forest_test.go`](../../../internal/site/forest_test.go))
    gains coverage for C1. Two shapes were decomposed:
    - **OD4.a — Presence-and-absence semantically.** Assert
      `name`/slug present and `wst-` absent.
    - **OD4.b — Presence only.** Assert `name`/slug present;
      uuid-alongside-name regression not caught.

- **OD5 — Pre-flag the `tree.go` / `site.go` Estimate
  Deviation. Resolved = dissolved by OD1.a.** Under OD1.a,
  `tree.go` (additive `Name` field) and `site.go`
  (`Server.index` threading) are **certain** edits, not
  estimate deviations. The `## Files to touch` section below
  promotes them from "OD1-dependent — pre-flagged" to "Modify
  (certain)." The implementing PR records the edits as the
  planned scope, not as an `## Estimate Deviations` callout;
  the parent task plan's `## Files to touch` p2 row already
  named the path as expected so this isn't a re-litigation of
  scope either way.

## Context

This is **phase 2 of three** for the
[`post-m2-ux-correction`](README.md) task. p2 closes the
m2-shipped forest/roster identity disagreement: the roster
identifies an active work-instance by its reported `name` (slug
fallback per the rule t4 locked for the roster), while the
forest identifies the same work-instance by the raw
`wst-<uuid>` actor it generated at registration time. The
disagreement is the deliberate scope of the
[`humanize-forest-actor`](../../backlog.md#humanize-forest-actor)
backlog entry — graduated `Open → Graduated —
post-m2-ux-correction` in the parent task-plan drafting change;
this phase realizes the intent the entry was opened to track.

p2 ships **after** p1 (cosmetic defects) so the
dark-mode-legible chrome p1 establishes is the baseline against
which the F4 actor-marker change is visually observed in both
OS modes; p2 ships **before** p3 (contract-revisiting) so the
actor-marker change is not entangled with p3's larger render
reshuffle (the default-cells row, the nested-`<details>` body
disclosure, the K3 every-entry-opens disclosure). No plan-doc
supersessions are realized in this phase. The graduating
backlog entry's Status flip lives in the parent task-plan
drafting change (already executed); the work that **closes**
the entry's intent lands here.

## Goal

Open the page against `docs/plans/` per
[`docs/dev.md`](../../dev.md) and observe:

- A bound active session whose reported metadata carries a
  `name` renders that name inside the attached node's
  `actor-marker` span — never the raw `wst-<uuid>` actor.
- A bound active session that reported no `name` (or no
  metadata at all) renders the **slug fallback** inside the
  attached node's `actor-marker` span — the same fallback the
  roster shows for the same work-instance.
- The forest and the roster display the **same identity** for
  the same active work-instance: a name-bearing session reads
  identically in both surfaces; a no-name session reads as its
  slug in both surfaces.
- An unbound work-instance (a session registered against a slug
  not in the walked tree) is **not** rendered in the forest —
  the roster remains its only surface. Pre-existing behavior of
  `buildTree`'s exact-slug join, recorded here for reviewer
  completeness (OD3.a).

All while: the parent task plan's Cross-Cutting Invariants hold
— **C-INV-2 (region-boundary preservation)** is honored under
the milestone's t4-amended data-path carve-out
([`m2/README.md` Cross-Task Invariants → data-path carve-out](../workstream-tracker-1-0/m2/README.md));
edits are scoped to the forest-region body
([`forest.go`](../../../internal/site/forest.go)) plus any
additive data-flow extension in
[`tree.go`](../../../internal/site/tree.go) and
[`site.go`](../../../internal/site/site.go) that OD1's
resolution requires. **C-INV-3 (actor-tag preserved):** the
`actor-marker` span itself, its placement inside `label-group`,
and the per-node `range .WorkInstances` attachment remain —
only the rendered text inside the span changes.
**C-INV-4 (walk-on-every-request preserved):** no caching, no
file-watch, no in-memory build-up; the reported `name` rides
the already-per-request
[`loadSessionMetadata`](../../../internal/site/site.go) read.
**C-INV-5 (additive, no spec / API / schema / dependency
change):** no frontmatter field, no API endpoint or
request/response field, no schema column, no dependency added.
The `wst-<uuid>` actor remains the internal identity key — no
actor-generator change, no schema change.

## Contracts

### C1 — Forest's per-node `actor-marker` renders the reported `name` (slug fallback) (F4)

This phase **realizes** parent task plan Contract **C6**
([`README.md` Contracts](README.md#contracts)); C1 below
restates the WHAT for this phase's surface and is the
authoritative contract clause for the implementing PR's review.

The forest's per-node `actor-marker` span renders the reported
session `name` when the work-instance's session reported one,
and the **slug** fallback when it did not — never the raw
`wst-<uuid>` actor. The fallback rule is **name-then-slug**, the
same rule t4 locked for the roster
([`t4-session-roster.md` "Session identity and naming"](../workstream-tracker-1-0/m2/t4-session-roster.md));
the `wst-<uuid>` actor remains loader-internal identity (the
idempotency key keyed on `(slug, actor)`), never a label. The
`actor-marker` span itself, its placement inside `label-group`,
and the per-node `range .WorkInstances` attachment remain
unchanged — this contract changes only the text rendered inside
the span. The data-flow shape locks at **OD1.a** (additive `Name` field
on `ActiveWorkInstance`, populated in `Server.index` between
`loadSessionMetadata` and `buildTree`); the slug fallback's
source on the rendered value locks at **OD2.a** (additive
`Slug` field on `ActiveWorkInstance` alongside `Name`,
populated in the same pass from the `active`-map key); see
Status → Open decisions. The implementing PR's narrow
render-altitude calls (exact span text, dark/light mode
observation) are settled per
[`shared.md`](../../../spec/planning/shared.md) "Plans describe
contracts, not implementation." The per-request resolution of
the reported `name` continues to live in the already-existing
[`loadSessionMetadata`](../../../internal/site/site.go) loader
— the forest consumes the same resolved-metadata map the roster
consumes; no second resolution path is introduced.

`Verified by:` [`node-header` template's `actor-marker` span at
forest.go:51](../../../internal/site/forest.go) (the span this
phase edits the rendered text of); [`ActiveWorkInstance` at
tree.go:39](../../../internal/site/tree.go) (the per-node
work-instance shape this phase additively extends per OD1's
resolution); [`loadSessionMetadata` and `resolveSessionMeta` in
site.go](../../../internal/site/site.go) (the
already-per-request name resolution this phase consumes — no
new loader, no new query); [`buildRoster` in
site.go](../../../internal/site/site.go) (the existing
"pure-function-over-already-loaded-values" precedent the forest
extends); [`t4-session-roster.md` "Session identity and
naming"](../workstream-tracker-1-0/m2/t4-session-roster.md) (the
name-then-slug rule reused).

## Cross-Cutting Invariants

Inherited from the parent task plan's
[`## Cross-Cutting Invariants`](README.md#cross-cutting-invariants);
the phase's implementing change must hold every one
simultaneously. Reviewer-flag candidates when the implementing
PR brushes any of these.

- **C-INV-2 (region-boundary preservation).** The forest-region
  body in [`forest.go`](../../../internal/site/forest.go) is the
  primary edit. The milestone's t4-amended data-path carve-out
  ([`m2/README.md` Cross-Task Invariants](../workstream-tracker-1-0/m2/README.md))
  authorizes the additive struct field on `ActiveWorkInstance`
  ([`tree.go`](../../../internal/site/tree.go)) and the
  `Server.index` threading in
  [`site.go`](../../../internal/site/site.go) that OD1's
  resolution may require. The shell
  ([`render.go`](../../../internal/site/render.go)) and the
  roster region body
  ([`roster.go`](../../../internal/site/roster.go)) are not
  touched.
- **C-INV-3 (actor-tag preserved).** The `actor-marker` span,
  its placement inside `label-group`, and the
  `range .WorkInstances` attachment remain. A diff that removes
  the span, removes the per-work-instance range, removes the
  work-instance-to-node attachment, or collapses multiple
  work-instances into a single marker is the defect, not the
  fix.
- **C-INV-4 (walk-on-every-request preserved).** No caching,
  no file-watch, no in-memory build-up. The reported `name`
  rides the existing per-request `loadSessionMetadata` read;
  no second resolution path or second event-log query is
  introduced.
- **C-INV-5 (additive, no spec / API / schema / dependency
  change).** No frontmatter field, no API endpoint or
  request/response field, no schema column, no dependency added.
  Render + data-path threading only. C-INV-1 (cell-anchor) is
  out of this phase's surface — p3 carries the cell-DOM
  contract.

## Files to touch

*Per OD1.a + OD2.a, the file inventory is `forest.go` +
`tree.go` + `site.go` + `forest_test.go` — all certain. The
parent task plan's `## Files to touch` p2 row pre-flagged the
`tree.go` / `site.go` deviation path as expected (OD5 dissolved
by OD1.a), so the implementing PR records these as planned
scope, not as `## Estimate Deviations` callouts. If the
implementation discovers a genuine structural deviation past
this set, the implementing PR handles it via `## Estimate
Deviations` per
[`task-plan.md`](../../../spec/planning/task-plan.md)
"Plan-to-PR Completion Gate," with the plan reconciled to what
shipped.*

- **Modify (certain):**
  - [`internal/site/forest.go`](../../../internal/site/forest.go)
    — the `node-header` template's `actor-marker` span's
    rendered text (**C1**). The `range .WorkInstances` and the
    span attachment stay; only the text inside the span changes.
  - [`internal/site/tree.go`](../../../internal/site/tree.go) —
    additive `Name string` and `Slug string` fields on
    `ActiveWorkInstance` (per OD1.a + OD2.a). The struct's
    existing `ID` field — added in t4-p2 for the roster's
    event-log join even though the forest doesn't render `ID`
    — is the additive-field-for-cross-region-needs precedent
    these extend.
  - [`internal/site/site.go`](../../../internal/site/site.go) —
    `Server.index` populates the two new fields on each
    `ActiveWorkInstance` between `loadSessionMetadata` and
    `buildTree` in a single pass over `active`'s entries:
    `Name` from `meta[wi.ID].Name`, `Slug` from the
    `active`-map's slug key (in scope from
    `for slug, wis := range active`). No
    `loadActiveWorkInstances` change; no `buildTree` signature
    change; no second resolution path.

- **Modify (tests):**
  - [`internal/site/forest_test.go`](../../../internal/site/forest_test.go)
    — semantic coverage for C1 per **OD4.a**: a node carrying
    a work-instance with a reported `name` renders the name
    inside the `actor-marker` span; a node carrying a
    work-instance whose session reported no `name` renders the
    slug fallback; in neither case does the rendered
    `actor-marker` text contain the `wst-` prefix (the absence
    side catches the silent-regression mode where a future
    template change surfaces both `Name` and `Actor`). The
    semantic-not-byte-exact posture m2 t3 established
    ([`t3-doc-declared-stages.md` Contracts](../workstream-tracker-1-0/m2/t3-doc-declared-stages.md))
    is preserved.

- **Intentionally not touched** *(estimate — where we don't
  expect changes, not a hard prohibition):*
  - [`internal/site/render.go`](../../../internal/site/render.go)
    — no shell edit (p1 owned the F1 body rule; p2 has no
    shell-level work).
  - [`internal/site/roster.go`](../../../internal/site/roster.go)
    — the roster's identity rendering is unchanged; the forest
    reconciles *to* the roster, not the other way around.
  - [`internal/site/walker.go`](../../../internal/site/walker.go)
    — no frontmatter field; `parsedDoc` shape unchanged.
  - [`internal/api/`](../../../internal/api/),
    [`internal/db/`](../../../internal/db/),
    [`internal/registerclient/`](../../../internal/registerclient/),
    [`cmd/workstream-tracker/`](../../../cmd/workstream-tracker/)
    — no API, schema, register-client, or CLI change.
  - [`spec/planning/`](../../../spec/planning/) — no spec
    change.
  - The Landed plan docs the parent task plan references —
    [`m2/README.md`](../workstream-tracker-1-0/m2/README.md),
    [`m2/t4-session-roster.md`](../workstream-tracker-1-0/m2/t4-session-roster.md)
    — are not retro-edited.

## Validation Gate

This is the **per-phase** Validation Gate (observation-only).
The task-terminal full product-acceptance walkthrough across all
six findings lives at p3's implementing PR per the parent task
plan's [`## Phase Contracts`](README.md#phase-contracts). p2's
gate observes the F4 corrected behavior in isolation; the phase
plan's Status flips `In progress → Validating → Landed` per
[`task-plan.md`](../../../spec/planning/task-plan.md); the
parent task plan's own Status flip waits for p3.

### Setup

From a clean working tree at the implementing PR's head: start
the server against the dogfood plan-tree per
[`docs/dev.md`](../../dev.md) "Local Workflow" —
`go run ./cmd/workstream-tracker` with the defaults
(`PORT=8080`, `DB_PATH=./workstream-tracker.db`; the SQLite file
seeds itself on first run).

### Seed sessions

From a separate shell, using the module-path form per
[`docs/dev.md`](../../dev.md) "Registering and completing a
session":

1. A **bound session with a `--name`**: register against any
   walked plan-tree slug present in `docs/plans/` (e.g.
   `post-m2-ux-correction`) with `--name "Demo bound"` so the
   reported metadata carries a name.
2. A **bound session with no `--name`**: register against a
   *different* walked plan-tree slug without `--name`. The
   session reports no name; the work-instance attaches to the
   node by exact slug match but renders the slug fallback in
   the forest.
3. An **unbound session** (any name): register against a typoed
   slug (e.g. `post-m2-ux-correction-typo`). This session
   appears in the roster but not in the forest; included to
   observe the OD3.a no-change behavior.

Each session stays active (the `register` subcommand exits 0
without holding state).

### Observe

Open `http://localhost:8080/` and observe:

- **C1 acceptance — name-bearing session.** The forest node
  attached to seeded session 1 renders `Demo bound` inside its
  `actor-marker` span. The roster shows the same session with
  the same label `Demo bound`. The forest and roster identify
  the same work-instance identically.
- **C1 acceptance — no-name session.** The forest node attached
  to seeded session 2 renders the slug fallback inside its
  `actor-marker` span — the same value the roster renders as
  that entry's label. No `wst-<uuid>` text appears in the
  forest's `actor-marker` span for either bound session.
- **C-INV-3 acceptance — no-regress on attachment.** Each
  node's `actor-marker` span is still present and still placed
  inside `label-group`; a node with multiple active
  work-instances renders one `actor-marker` per work-instance
  (the per-node `range .WorkInstances` attachment is intact).
- **OD3.a acceptance — unbound stays roster-only.** Seeded
  session 3 appears in the roster (as an unbound entry) but
  **not** in the forest. The pre-existing `buildTree`
  join-side drop is preserved; p2 does not regress it and does
  not extend the forest to surface unbound sessions.

### Tear down

Close the seeded sessions per
[`docs/dev.md`](../../dev.md) "Registering and completing a
session" for each `work_instance_id` receipt.

### Toolchain gate

`gofmt -l internal cmd` reports no files; `go build ./...`,
`go vet ./...`, `go test ./...` all pass per
[`docs/dev.md`](../../dev.md) "Local Workflow." New / extended
test coverage in
[`internal/site/forest_test.go`](../../../internal/site/forest_test.go)
includes the C1 semantic assertions per OD4.a (presence of
`Name`/slug fallback **and** absence of the `wst-` prefix).

## Self-Review Audits

From
[`docs/agents/local/self-review-catalog.md`](../../agents/local/self-review-catalog.md);
diff surfaces are the forest-region template + the data path
between `loadSessionMetadata` and `buildTree`:

- **validation-honesty** — the C1 acceptance is observed
  against a real `go run` rendering with both name-bearing and
  no-name seeded sessions, not asserted from the diff or from
  the template source alone.
- **rename-aware-diff-classification** — the additive `Name`
  field on `ActiveWorkInstance` (OD1.a) and the new
  `Server.index` population pass must be hand-classified so the
  "actor-marker text changes only" and "C-INV-3 attachment
  preserved" claims are real, not tooling-fooled.

## Out of Scope

- **Changing the `wst-<uuid>` actor generator.** C1 changes the
  rendered text inside the `actor-marker` span; the underlying
  actor remains the internal identity / idempotency key. Out of
  scope per parent task plan's `## Out Of Scope` and C-INV-5;
  scoping doc SD6 Rejected F1.
- **Hover-tooltip preserving the uuid.** Rejected at scoping —
  scoping doc SD6 Rejected F2 (preserves the very inconsistency
  the backlog entry was opened against).
- **Rendering unbound work-instances in the forest.** OD3.a:
  the forest renders attached work-instances only (the
  `buildTree` exact-slug join). The roster surfaces unbound
  work-instances per t4; the forest does not. Pre-existing
  behavior preserved.
- **Touching the roster's identity rendering.** The roster's
  name-then-slug-fallback render lives in
  [`roster.go`](../../../internal/site/roster.go); p2
  reconciles the forest *to* the roster, not the other way
  around.
- **Adding fields to the reported-metadata schema.** No new
  conventionally-read keys beyond the existing `name` t4
  contract names; per C-INV-5.
- **Backlog file mutations.** The `humanize-forest-actor`
  entry's `Open → Graduated — post-m2-ux-correction` Status
  flip lives in the parent task-plan drafting change (already
  executed). This phase's implementing PR realizes the entry's
  *intent* but performs no backlog-file mutation.

## Risk Register

- **Slug fallback ladder accidentally falls through to the
  actor.** If OD1's resolution wires a name-then-fallback ladder
  in the template, a mis-wired ladder could read the
  `wst-<uuid>` actor when both `name` and the fallback-source
  value are empty. Mitigation: the slug is **always available**
  for an attached work-instance — `buildTree` joins on exact
  slug match
  ([`tree.go:148`](../../../internal/site/tree.go)
  `active[d.Slug]`), and `work_instances.slug` is `NOT NULL`
  per t4's verified schema citation
  ([`t4-session-roster.md` "Session identity and naming"](../workstream-tracker-1-0/m2/t4-session-roster.md)),
  so the slug fallback's value is never absent. The OD4.a
  semantic test asserts the **absence** of the `wst-` prefix to
  catch a mis-wired ladder. Carried, not blocking.
- **Per-work-instance `Name` resolution diverges from the
  roster's.** A second name-read path (a second call to
  `resolveSessionMeta`, a second event-log query, or a second
  key read on the metadata blob) would let the forest show a
  different name than the roster for the same work-instance
  under heartbeat / state-transition overlay races.
  Mitigation: under OD1.a the threading point is structurally
  shared — `Server.index` populates `ActiveWorkInstance.Name`
  from the same `map[string]sessionMeta` `buildRoster` reads to
  populate `RosterEntry.Name`. The resolution runs once per
  request in `loadSessionMetadata`; both surfaces consume the
  same map keyed by `wi.ID`. Reviewer-flag any commit that
  introduces a second `resolveSessionMeta` call, a second
  event-log query, or a separate metadata-blob key read in the
  forest's loader. Carried, not blocking.
- **Additive struct field is silently rendered elsewhere.**
  Under OD1.a, the additive `Name` field on
  `ActiveWorkInstance` could be picked up by an unrelated
  consumer (e.g. a debug log or a `%+v` format). Mitigation:
  the forest's `node-header` template is the only intended
  reader; the roster does not iterate `ActiveWorkInstance`
  values (it consumes `RosterEntry` from `buildRoster`).
  Spot-check at implementing-PR review.

## Related Docs

- [`README.md`](README.md) — parent task plan; the locked
  Contract **C6**, Cross-Cutting Invariants C-INV-2 / C-INV-3 /
  C-INV-4 / C-INV-5 inherited here (C-INV-1 cell-anchor is out
  of this phase's surface), the Files-to-touch p2 estimate
  (with the OD1-driven `tree.go` / `site.go` deviation path
  pre-flagged), and the task-terminal Validation Gate this
  phase's per-phase gate feeds into at p3.
- [`scoping/post-m2-ux-correction.md`](scoping/post-m2-ux-correction.md)
  — paired scoping doc; **SD6** records the scoping-time
  deliberation for F4 (decomposed shape: render reported name
  with slug fallback; rejected: actor-generator change at SD6
  Rejected F1, hover-tooltip at SD6 Rejected F2).
- [`../../backlog.md`](../../backlog.md)
  `humanize-forest-actor` — the graduated backlog entry whose
  intent this phase closes; its `Open → Graduated —
  post-m2-ux-correction` Status flip lives in the parent
  task-plan drafting change.
- [`../workstream-tracker-1-0/m2/t4-session-roster.md`](../workstream-tracker-1-0/m2/t4-session-roster.md)
  — the name-then-slug rule t4 locked for the roster (reused
  by C1) and the milestone-amended data-path carve-out (which
  authorizes p2's `tree.go` / `site.go` deviation path under
  C-INV-2). Not retro-edited.
- [`../workstream-tracker-1-0/m2/t3-doc-declared-stages.md`](../workstream-tracker-1-0/m2/t3-doc-declared-stages.md)
  — the semantic-not-byte-exact test posture (C7) OD4's
  recommended shape reads against.
- [`../workstream-tracker-1-0/m2/README.md`](../workstream-tracker-1-0/m2/README.md)
  — milestone Cross-Task Invariants (region-boundary file
  enforcement + the t4-amended data-path carve-out that
  authorizes p2's deviation path; the v0.1 actor-tag no-regress
  invariant C-INV-3 inherits from). Not retro-edited.
- [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md),
  [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md)
  — the rules this phase plan is structured against.
- [`../../dev.md`](../../dev.md) — the local workflow and
  session-registration commands the Validation Gate invokes.
