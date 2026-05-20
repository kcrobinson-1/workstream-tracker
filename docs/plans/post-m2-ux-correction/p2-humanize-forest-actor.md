---
slug: post-m2-ux-correction-p2
Status: In draft
short_description: Humanize the forest actor — reported name (slug fallback) in the per-node actor-marker (F4)
---

# Phase 2 — Humanize forest actor (F4)

## Status

`In draft`. This drafting session captures the phase's Goal,
Contract clause **C1** (the WHAT this phase realizes for parent
Contract **C6**), the Cross-Cutting Invariants the phase
inherits from the parent task plan
([C-INV-2 / C-INV-3 / C-INV-4 / C-INV-5](README.md#cross-cutting-invariants)),
the estimate-shaped Files-to-touch (pre-flagging the `tree.go` /
`site.go` deviation path the parent task plan named), the
per-phase Validation Gate (observation-only — the task-terminal
full product-acceptance walkthrough lives at p3), Out of Scope,
and the Risk Register. Five **open decisions** are surfaced
below for the `` `In draft` → `Proposed` `` promotion-gate walk
to drive to resolution; none are locked here.

### Open decisions

Each is decomposed against merged code; resolutions fold into
Contract **C1**, Files-to-touch, the Validation Gate, and the
implementing PR. None are locked in this drafting session.

- **OD1 — Data-flow shape for the reported `name`.** The
  forest's per-node `range .WorkInstances` iterates
  `*ActiveWorkInstance` values that carry `ID` and `Actor` only
  today ([`ActiveWorkInstance` in
  tree.go](../../../internal/site/tree.go); the `ID` field
  exists for the roster's event-log join, never rendered by the
  forest). The reported `name` is already resolved per request
  by [`loadSessionMetadata`](../../../internal/site/site.go) into
  a `map[string]sessionMeta` keyed by `wi.ID`, and consumed by
  `buildRoster` ([`buildRoster` in
  site.go](../../../internal/site/site.go)). Threading the
  resolved value to the forest's per-node template has three
  plausible shapes:
  - **OD1.a — Add an additive `Name` field (and a `Slug` field
    if OD2 lands on OD2.a) to `ActiveWorkInstance`, populated in
    `Server.index` between `loadSessionMetadata` and
    `buildTree`.** No `buildTree` signature change. The forest
    template reads `.Name` directly (no template function). The
    roster continues to read from `RosterEntry` and ignores the
    new field. *Touches:*
    [`tree.go`](../../../internal/site/tree.go) (struct
    additive),
    [`site.go`](../../../internal/site/site.go)
    (`Server.index` threading + `loadActiveWorkInstances` is
    untouched — the field is populated post-load, not in the SQL
    read), and
    [`forest.go`](../../../internal/site/forest.go) (template
    text).
  - **OD1.b — Pass a parallel `map[string]string` (id → name)
    into `buildTree`, look up in the template via a registered
    template function.** `ActiveWorkInstance` stays loader-key
    shape (no rendered fields). `buildTree`'s signature grows by
    one map argument; the template grows a function call. *New
    surface:* a render-time template function in the forest's
    template set. *Touches:*
    [`tree.go`](../../../internal/site/tree.go) (`buildTree`
    signature),
    [`site.go`](../../../internal/site/site.go) (call-site map
    derivation from `meta`), and
    [`forest.go`](../../../internal/site/forest.go) (template
    function registration + template text).
  - **OD1.c — Pass the full `meta map[string]sessionMeta`
    through `buildTree` (or onto a render context), look up in
    the template via a template function.** Same shape as OD1.b
    but threads the full resolved-metadata blob rather than a
    derived name-only map. Larger surface than C1 needs (the
    forest reads only `Name`), but leaves room for a future
    forest field to read further reported facts without a second
    threading pass.

  Tradeoff lens: OD1.a is the smallest reader-side change (the
  template stays a plain field read), at the cost of an additive
  field on `ActiveWorkInstance` that only the forest renders.
  OD1.b / OD1.c keep `ActiveWorkInstance` as loader-key shape
  but introduce a render-time template function the forest
  template set does not use today. All three preserve C-INV-4
  (single per-request resolution via `loadSessionMetadata`).
  *Verified by:*
  [`Server.index` order in
  site.go](../../../internal/site/site.go) (docs → active → meta
  → buildTree(docs, active) → buildRoster(docs, active, meta) —
  `meta` is already loaded before `buildTree` is called, so
  OD1.a can populate the new field between those two steps with
  no reordering);
  [`buildRoster` in
  site.go](../../../internal/site/site.go) (the existing
  "pure-function-over-already-loaded-values" precedent OD1.a/b/c
  each match).

- **OD2 — Slug fallback's source on the rendered value.** Parent
  Contract **C6** locks the rule as **name-then-slug**, the same
  rule t4 locked for the roster
  ([`t4-session-roster.md` "Session identity and naming"](../workstream-tracker-1-0/m2/t4-session-roster.md));
  the `wst-<uuid>` actor is the loader-internal identity key and
  is never a label (banned by parent C6 + C-INV-3). Confirm the
  *source* of the slug to fall back to:
  - **OD2.a — The work-instance's registered slug.** The slug
    the session registered against — the same value the roster
    falls back to in `RosterEntry.Slug` today. Requires a `Slug`
    field on `ActiveWorkInstance` (additive — only OD1.a's
    shape; OD1.b/c carry the slug through the threaded map).
  - **OD2.b — The enclosing node's slug** (`.Slug` on the
    parent `PlanNode`, available via outer-template scope —
    `$.Slug` or a pre-`range` template variable inside
    `node-header`). No additive field needed.

  Tradeoff lens: For every *attached* work-instance, OD2.a and
  OD2.b produce **identical** rendered output — `buildTree`
  attaches via exact slug match
  ([`tree.go:148`](../../../internal/site/tree.go) `active[d.Slug]`),
  so an attached work-instance's registered slug equals the
  enclosing node's slug. The difference is semantic: OD2.a keeps
  the rule keyed to the same identity the roster reads (the
  registered slug — the durable fact about the session); OD2.b
  ties the forest's fallback semantically to the node identity.
  Name-then-actor-id is **not** a candidate — banned by parent
  C6 + C-INV-3. OD2 interacts with OD1: OD1.a + OD2.a needs both
  `Name` and `Slug` on `ActiveWorkInstance`; OD1.a + OD2.b needs
  only `Name`.

- **OD3 — Unbound work-instance handling.** A work-instance
  registered against a slug not present in the walked tree is
  *unbound*. Today the forest drops it at `buildTree`'s join
  site ([`tree.go:104`](../../../internal/site/tree.go);
  `active[d.Slug]` reads only when the slug matches a parsed
  doc), and the roster lists it
  ([`buildRoster` in
  site.go](../../../internal/site/site.go) classifies unbound
  entries with `Bound: false`). Confirm:
  - **OD3.a — No change; unbound work-instances stay
    roster-only.** The forest renders only attached
  work-instances; the drop predates p2 and falls naturally out
    of the exact-slug join. The parent task plan's contract
    surface (C6) places the F4 humanization on the forest's
    *existing* per-node `actor-marker` attachments, not on a new
    surface for unbound sessions. Recorded in `## Out of Scope`.
  - **OD3.b — Render unbound work-instances somewhere in the
    forest.** Out of scope per the parent task plan's split
    between roster (which surfaces unbound) and forest (which
    renders attached); would require new render plumbing past
    F4.

  Tradeoff lens: OD3.a is the natural reading of parent C6 and
  the t4 / forest split — surfacing unbound in the forest would
  expand p2's scope past F4 with no contract support. Confirm
  OD3.a and record the drop site in `## Out of Scope` so
  reviewers don't read the pre-existing drop as a p2 regression.

- **OD4 — Test-coverage posture for the F4 contract.** The
  forest-region test surface
  ([`internal/site/forest_test.go`](../../../internal/site/forest_test.go))
  gains coverage for C1. Two shapes:
  - **OD4.a — Assert presence-and-absence semantically: the
    rendered `actor-marker` text contains the reported `name`
    (or the slug fallback when no `name` was reported) **and**
    does not contain the `wst-` actor prefix.** Same posture as
    m2 t3 C7's semantic-not-byte-exact assertions
    ([`t3-doc-declared-stages.md` Contracts](../workstream-tracker-1-0/m2/t3-doc-declared-stages.md)).
  - **OD4.b — Assert presence only** — the rendered
    `actor-marker` text contains the reported `name` (or the
    slug fallback). Smaller assertion surface; a regression that
    re-introduced the uuid alongside the name in the same span
    would not be caught.

  Tradeoff lens: OD4.a is the t3 C7 precedent and catches the
  regression OD4.b misses. Recommend OD4.a; the gate-walk
  confirms.

- **OD5 — Pre-flag the `tree.go` / `site.go` Estimate
  Deviation.** Parent task plan's `## Files to touch` names
  `forest.go` as the certain p2 edit and `tree.go` / `site.go`
  as estimated deviations "if the reported-`name` field needs an
  additive passthrough" ([`README.md` Files to touch — p2](README.md#files-to-touch)).
  All three OD1 shapes touch `tree.go` and `site.go` (additively
  in OD1.a; signature-and-call-site in OD1.b/c), so the
  deviation path is the **expected** shape, not the unexpected
  one. The phase plan's Files-to-touch should pre-flag this so
  the implementing PR can name it under `## Estimate Deviations`
  as a structural-call confirmation rather than a re-litigation
  of scope. This is a record-keeping decision (pre-flag yes/no),
  not a behavior decision.

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
the span. The data-flow shape (OD1), the slug fallback's source
on the rendered value (OD2), and any additive passthrough on
`ActiveWorkInstance` or through `buildTree` are render-altitude
/ data-plumbing decisions resolved at the gate-walk and the
implementing PR per
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

*Estimate of the expected file shape, not a binding rule.
Implementation may revise this list when a structural call
requires it. Any deviation is handled via the implementing PR's
`## Estimate Deviations` callout with the plan reconciled to
what shipped per
[`task-plan.md`](../../../spec/planning/task-plan.md)
"Plan-to-PR Completion Gate." The final file inventory is
OD1's resolution at the gate-walk; the parent task plan
pre-flagged the `tree.go` / `site.go` deviation path (OD5 above
asks whether the phase plan pre-flags it explicitly).*

- **Modify (certain):**
  - [`internal/site/forest.go`](../../../internal/site/forest.go)
    — the `node-header` template's `actor-marker` span's
    rendered text (**C1**). The `range .WorkInstances` and the
    span attachment stay; only the text inside the span changes.

- **Modify (OD1-dependent — Estimate Deviation pre-flagged):**
  - [`internal/site/tree.go`](../../../internal/site/tree.go) —
    additive `Name` (and `Slug` per OD2.a) field on
    `ActiveWorkInstance` under OD1.a; or `buildTree` signature
    growth under OD1.b / OD1.c. The parent task plan's `## Files
    to touch` p2 row already names this deviation path; this
    plan inherits the pre-flag (OD5 = pre-flag yes).
  - [`internal/site/site.go`](../../../internal/site/site.go) —
    `Server.index` threading to populate the new field from the
    already-loaded `meta` map between `loadSessionMetadata` and
    `buildTree` (OD1.a); or `buildTree` call-site map derivation
    (OD1.b / OD1.c). Same Estimate Deviation pre-flag.

- **Modify (tests):**
  - [`internal/site/forest_test.go`](../../../internal/site/forest_test.go)
    — semantic coverage for C1 per OD4: a node carrying a
    work-instance with a reported `name` renders the name inside
    the `actor-marker` span; a node carrying a work-instance
    whose session reported no `name` renders the slug fallback;
    in neither case does the rendered `actor-marker` text
    contain the `wst-` prefix (OD4.a — recommended). The
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
includes the C1 semantic assertions per OD4.

## Self-Review Audits

From
[`docs/agents/local/self-review-catalog.md`](../../agents/local/self-review-catalog.md);
diff surfaces are the forest-region template + the data path
between `loadSessionMetadata` and `buildTree`:

- **validation-honesty** — the C1 acceptance is observed
  against a real `go run` rendering with both name-bearing and
  no-name seeded sessions, not asserted from the diff or from
  the template source alone.
- **rename-aware-diff-classification** — additive struct field
  on `ActiveWorkInstance` (OD1.a) or signature growth on
  `buildTree` (OD1.b / OD1.c) must be hand-classified so the
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
  roster's.** If OD1's resolution introduces a second name-read
  path (a second call to `resolveSessionMeta`, a second event-
  log query, or a second key read on the metadata blob), the
  forest could show a different name than the roster for the
  same work-instance under heartbeat / state-transition
  overlay races. Mitigation: OD1.a / OD1.b / OD1.c all consume
  the **same** `loadSessionMetadata` result — the resolution
  runs once per request in `Server.index`; both surfaces read
  from the same `map[string]sessionMeta` keyed by `wi.ID`. The
  OD1 resolution must preserve this single-resolution
  invariant; reviewer-flag a second resolution path. Carried,
  not blocking.
- **Additive struct field is silently rendered elsewhere.** If
  OD1.a adds `Name` to `ActiveWorkInstance`, an unrelated
  consumer (e.g. a debug log or a `%+v` format) could pick it
  up. Mitigation: the forest's `node-header` template is the
  only intended reader; the roster does not iterate
  `ActiveWorkInstance` values (it consumes `RosterEntry` from
  `buildRoster`). Spot-check at implementing-PR review.

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
