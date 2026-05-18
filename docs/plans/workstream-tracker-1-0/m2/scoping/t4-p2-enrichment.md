# Scoping — t4 p2 Enrichment + named sessions

Scoping doc for [`t4-p2-enrichment.md`](../t4-p2-enrichment.md)
(`workstream-tracker-1-0-m2-t4-p2`), phase 2 of the N ≥ 2 task
plan [`t4-session-roster.md`](../t4-session-roster.md). Transient:
this doc lives in the milestone's `m2/scoping/` subfolder (p2 is a
phase under the m2 milestone, **not** a standalone task plan), so
its contents delete **in batch at the milestone-terminal PR** per
[`task-plan.md`](../../../../../spec/planning/task-plan.md) "Path
conventions" — "The `scoping/` subfolder is transient — its
contents delete in batch at the milestone-terminal PR (or the
task-terminal PR for standalone task plans)." m2 still has open
siblings (t2 `Proposed`, t3 stub) at p2 plan time, so the
subfolder is **not** deleted by p2's task-terminal PR. Carries no
`Status` field (scoping docs are deliberation, not lifecycle
artifacts).

## Context

p1 landed the bare bound/unbound roster: every active
work-instance listed and classified, slug as the label, via the
pure `buildRoster` over the per-request walked parsed-doc set and
the active-work-instance map. **p2** enriches the same region: the
register client/CLI begins reporting request `metadata` (including
a human session name); the roster loader gains the event-log join
that surfaces it; each entry shows the reported name (slug
fallback, never the `wst-<uuid>` actor) and an expandable,
deliberately-unstructured raw-JSON detail view of the session's
reported metadata. The behavior WHAT is locked at task-plan level
(t4 Contracts → "Session identity and naming" / "Reported data and
the raw-JSON detail view" / "Metadata read policy"); the
genuinely-open p2 questions are pure mechanism the t4 scoping doc
handed to this phase: the loader/struct extension and metadata-join
shape, the client/CLI plumbing to send the conventionally-read
`name` key, and the expandable raw-JSON detail rendering. This
scoping doc records those decisions (with rejected alternatives)
and the reality-check inputs the plan's promotion gate must
re-confirm; the durable contract lives in the plan doc, not here.

## Decisions made at scoping time

Each decision carries a `Verified by:` code citation per
[`shared.md`](../../../../../spec/planning/shared.md) "`Verified
by:` annotations on load-bearing claims."

### SD1 — Carry `work_instances.id` on the loader output and resolve metadata in a companion per-request event read; `buildRoster` stays a pure classifier

p1 recorded (its plan Goal "Sibling-interface handoff to p2" and
scoping SD2) that p2 must extend `loadActiveWorkInstances` and the
entry struct to carry `work_instances.id` — the key
`events.work_instance_id` references — because p1 selects only
`slug, actor`. The open p2 mechanism is *where* the event-log read
that realizes the baseline+overlay rule lives. Decomposed into
shapes before analyzing per
[`shared.md`](../../../../../spec/planning/shared.md) "Decompose
options into shapes before analyzing":

- **A — combined join inside `loadActiveWorkInstances`.** One
  query joins `work_instances` to `events` and resolves the
  baseline-plus-overlay metadata in SQL (a per-work-instance
  latest-later-event selection). Rejected as the realizing
  mechanism: the baseline+key-by-key-overlay rule pushed into a
  single SQL statement is the exact code-shaped complexity the
  no-fenced-code / no-inline-predicate plan-doc rule keeps out of
  the plan and the key-by-key JSON merge is not expressible in
  the sqlite the driver exposes without per-row post-processing
  anyway, so the merge lands in Go regardless; folding the
  current-state read and the event read into one statement also
  couples a clean current-state loader to the event-log concern.
- **B — companion per-request event reader; `buildRoster` stays
  pure.** `loadActiveWorkInstances` is extended only to also
  select `work_instances.id` (the minimal honoring of the p1
  handoff). `Server.index` then performs one further per-request
  read over exactly the active work-instance id set, scoped to
  those ids and keyed by the indexed `events.work_instance_id`,
  and the register-baseline / latest-later-event key-by-key
  overlay is folded in Go. `buildRoster` remains a pure function
  over the per-request values (now also the resolved
  per-work-instance metadata), unchanged in shape from p1.
- **C — `buildRoster` queries the DB itself.** Rejected for the
  same reason p1 scoping SD1 rejected its B: it makes the roster
  builder impure and acquires its own read instead of consuming
  the per-request values `Server.index` already holds, eroding
  the walk-on-every-request "one per-request read posture, loader
  takes values as parameters" shape p1 established.

Decision: **B**. The event read is genuinely new data (the
session's reported metadata) not otherwise loaded in the request,
so it is not a duplicate read of an existing source the
walk-on-every-request invariant forbids — that invariant bars a
*cache / file-watch / in-memory accumulation* and a *second walk
of already-walked data*, not a per-request read of a distinct
source. The read is scoped to the active id set and uses the
existing `events` index, consistent with the t4 task plan Risk
Register "Per-request events join cost" mitigation. `buildRoster`
staying pure preserves p1's recorded "what p2 reuses without
restructuring" surface (classification, `indexData` plumbing,
region).

`Verified by:`
[`loadActiveWorkInstances` in site.go](../../../../../internal/site/site.go)
selects `slug, actor` only (no `id`) for `state = active` rows;
[`Server.index` in site.go](../../../../../internal/site/site.go)
calls `walkPlans` then `loadActiveWorkInstances` then `buildTree`
/ `buildRoster` once per HTTP request — the scope into which one
further per-request event read over the active id set composes;
[`idx_events_work_instance_id` in schema.go](../../../../../internal/db/schema.go)
indexes `events(work_instance_id)`, the read's keying column;
[`events` / `work_instances` in schema.go](../../../../../internal/db/schema.go)
shows `events.work_instance_id` is the FK onto `work_instances.id`
(the `TEXT PRIMARY KEY`), not `actor` or `(slug, actor)`.

### SD2 — The id is carried by an additive field on `ActiveWorkInstance`; an additive unrendered struct field does not regress the forest

`loadActiveWorkInstances` returns
`map[string][]*ActiveWorkInstance`; "extend the loader to carry
`work_instances.id`" (the p1 handoff) means the id rides the
struct that function already returns. `ActiveWorkInstance` is
defined in `tree.go`, which the t4 task plan and p1 plan listed
under their *estimate* "Intentionally not touched" (forest join
unchanged). Shapes:

- **A — add an `ID` field to `ActiveWorkInstance`.** One line,
  additive. The forest `node-header` template ranges
  `.WorkInstances` rendering only `{{.Actor}}`; a new unrendered
  field cannot change a single byte of the forest actor-marker
  output, so the v0.1 forest-actor no-regress *invariant* (a
  render invariant) holds. It is a deviation from the t4/p1
  *estimate* that `tree.go` is untouched — an estimate, not a
  rule, reconciled via `## Estimate Deviations` and the plan
  update per
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "Plan-to-PR Completion Gate."
- **B — a parallel id-carrying structure in `site.go` that
  avoids `tree.go`.** Rejected: it forks the active-work-instance
  list into two parallel shapes keyed the same way, diverges from
  p1's explicitly-recorded single-struct handoff ("extend
  `loadActiveWorkInstances` and `RosterEntry`"), and buys only
  the avoidance of a one-line additive field on a struct the
  loader already returns — a worse design to preserve an estimate
  that is not a prohibition.

Decision: **A**. The no-regress invariant protects the *rendered
forest actor markers*, not the literal "never edit `tree.go`"
estimate; an additive unrendered field preserves the invariant,
and the deviation is plan-flagged in Files to touch, hand-classified
under `rename-aware-diff-classification`, and the forest markers
are observed unchanged in the Validation Gate.

`Verified by:`
[`ActiveWorkInstance` in tree.go](../../../../../internal/site/tree.go)
is the struct `loadActiveWorkInstances` populates and `buildTree`
attaches to nodes;
[`node-header` template in forest.go](../../../../../internal/site/forest.go)
ranges `.WorkInstances` emitting only `{{.Actor}}` inside
`actor-marker` spans — no other `ActiveWorkInstance` field is
rendered in the forest;
[`loadActiveWorkInstances` in site.go](../../../../../internal/site/site.go)
constructs `&ActiveWorkInstance{Actor: actor}` today.

### SD3 — The register-baseline / latest-later-event key-by-key overlay is folded in Go, per work-instance

The t4 task-level "Metadata read policy" contract fixes the rule
(register-event metadata is the identity baseline; the **latest**
later event's metadata overlaid key-by-key; per-request). The open
mechanism is where the fold runs. Decision: the companion reader
(SD1-B) reads, per active work-instance id, the `register` event's
`metadata` and the single latest later event's `metadata` (by the
server-side arrival timestamp), and the key-by-key overlay (later
keys win; keys absent from the later event keep the register
value) is computed in Go. The SQL is only an indexed, active-id-
scoped row read; no overlay semantics live in the query (the
scoping rationale for SD1-A). Non-object reported metadata (the
schema-loose posture permits any JSON) has no key-by-key merge
defined; the raw-JSON detail still renders it verbatim and the
conventional `name` read simply finds no key — an impl-robustness
edge carried in the plan Risk Register, not a schema imposed to
make it regular (that would breach the posture-tension invariant).

`Verified by:`
[`insertRegister` / `insertHeartbeat` / `insertStateTransition`
in handlers.go](../../../../../internal/api/handlers.go) each
write `nullableJSON(...metadata)` into `events.metadata` tagged by
`type` (`register` / `heartbeat` / `state_transition`) with
`received_at` the server arrival time — so per-work-instance
baseline (the `register` row) and latest-later (max `received_at`
among non-`register` rows) are both derivable from the event
rows;
[`nullableJSON` in handlers.go](../../../../../internal/api/handlers.go)
stores `NULL` when no metadata was sent, so the baseline / overlay
read must treat an absent blob as "contributes no keys," not as an
error.

### SD4 — The client stays schema-agnostic; the CLI owns the `name` convention and omits metadata when no name is given

The sole conventionally-read key is the optional `name` (t4
Contracts "Reported data"); every other reported key stays
arbitrary (schema-loose). Shapes for the producer surface:

- **A — `registerclient.Register` gains a typed `name`
  parameter.** Rejected: it schematizes the client toward the one
  conventionally-read key, which is exactly the homogenization
  the milestone posture-tension invariant ("do not schematize t4's
  reported data") names as the defect; the client is a thin
  request marshaller, not the place the `name` convention lives.
- **B — `registerclient.Register` forwards arbitrary `metadata`;
  the CLI builds the object.** `Register` gains a raw metadata
  parameter it marshals into the request `metadata` field
  verbatim (no shape knowledge). `runRegister` reads
  `--name` / `WST_NAME`, and when a name is present builds the
  conventional one-key object; when absent it sends **no**
  metadata at all (not an empty-string `name`), so the v0.2
  minimum still lists and the roster falls back to the slug.

Decision: **B**. It keeps the schema-loose posture at the client
boundary and localizes the `name` convention to the CLI/handshake
producer. There is no separate metadata POST: `name` rides the
single existing register request, so a metadata send failure *is*
the register-attempt failure the existing best-effort path already
surfaces to stderr while still exiting success — the
`error-surfacing-user-mutations` audit is satisfied by the
existing surfaced-failure path, with the plan's Validation Gate
observing it rather than asserting it.

`Verified by:`
[`registerclient.Register` in client.go](../../../../../internal/registerclient/client.go)
marshals only `{exact_slug, actor}` today and already returns a
non-2xx / unreachable / malformed-body outcome as an error the
caller decides is non-fatal;
[`runRegister` in cmd/workstream-tracker/register.go](../../../../../cmd/workstream-tracker/register.go)
already prints an explicit actionable stderr line on a failed
attempt and still returns 0 (best-effort), and reads
`--slug`/`--actor`/`--server` via a `flag.FlagSet` + env fallback
the `--name`/`WST_NAME` affordance extends;
[`RegisterRequest.Metadata` in api.go](../../../../../internal/api/api.go)
is `json.RawMessage` and
[`insertRegister` in handlers.go](../../../../../internal/api/handlers.go)
already writes `nullableJSON(req.Metadata)` to `events.metadata` —
the server side is fully plumbed; the client is the only
producer-side change.

### SD5 — Expandable detail is native `<details>`/`<summary>`; an entry with no reported metadata renders as a plain row, no disclosure

The repo is server-rendered with no JavaScript and no read API;
t2 resolved the collapsible mechanism to native HTML
`<details>`/`<summary>` (milestone Cross-Task Decisions). The
roster region body is `roster.go`'s owned surface. Decision: an
entry with resolved reported metadata renders as a
`<details>`/`<summary>` element — summary carries the
name-or-slug label and the bound/unbound tag, the body carries
the deliberately-unstructured raw JSON; an entry with **no**
reported metadata renders as a plain row with no disclosure
control, mirroring the established forest idiom where a leaf box
with no children has no disclosure control. The literal label
formatting (slug-fallback copy) and the raw-JSON presentation
(pretty-print, escaping) are render-time UX deferred under "Bans
on surface require rendering the consequence" and observed in the
plan's Validation Gate, not prescribed here.

`Verified by:`
[`roster.go`](../../../../../internal/site/roster.go) owns
`{{define "roster"}}` / `{{define "roster-style"}}` parsed into
`indexTmpl` via `init()` — the region body p2 extends;
[design/v0.1-design.md §7](../../../../../design/v0.1-design.md)
records the render as native HTML `<details>`/`<summary>`, "no
JavaScript and no new route," and the leaf-box-no-disclosure
idiom (a box with children is a `<details>`, a leaf has no
disclosure control);
[the m2 milestone Cross-Task Decisions
"Collapsible mechanism — RESOLVED at t2 drafting"](../README.md)
locks native `<details>`/`<summary>`, zero JavaScript.

## Decisions resolved at the task-plan level (not re-opened here)

Resolved as task-level contract during t4 task-plan drafting,
inherited not re-litigated — listed so the plan does not
re-open them:

- Metadata read policy (register baseline + latest-later-event
  key-by-key overlay, per-request) — t4 Contracts "Metadata read
  policy." p2 mechanism is SD3; the *rule* is not p2's to change.
- Name-then-slug label rule, `wst-<uuid>` banned from the roster
  UX — t4 Contracts "Session identity and naming." p2 realizes
  the reported-name arm; only the literal slug-fallback formatting
  is render-time-deferred.
- Event-log join, no schema change; the sole conventionally-read
  key is `name`, all other keys arbitrary — t4 Contracts
  "Reported data" / scoping D1.
- N ≥ 2 split (bare roster, then enrichment) — t4 scoping D5.
- The `deterministic-interactive-registration` split + new
  observability-residual and forest-actor-humanization backlog
  entries — t4 Backlog Impact / scoping D2; executed in p2's
  implementing PR, decision not re-opened.

## Plan structure handoff

- **Doc shape:** phase plan (the `p2` file of the N ≥ 2 task
  plan [`t4-session-roster.md`](../t4-session-roster.md)). Opens
  with a `## Context` naming the parent task plan; inherited rule
  references (Cross-Cutting Invariants, naming, metadata read
  policy) cite the parent task plan's section by name rather than
  duplicating it, per
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "How a phase plan cites its parent task plan."
- **PR-count branch test:** estimated **N = 1 PR**. The diff is
  the loader/struct id extension + companion event reader in
  `site.go`, the metadata-resolution + name/detail fields on
  `RosterEntry`, the `roster.go` region body (named-session
  display + expandable raw-JSON detail) and CSS, the roster
  tests, the `registerclient` metadata parameter, the
  `register.go` `--name`/`WST_NAME` affordance, and the doc
  currency + backlog mutations — one subsystem (the site render
  path) plus the thin register client/CLI, well under the
  >5-subsystem / >300-LOC split threshold in
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "PR-count predictions need a branch test." p2 is the
  task-terminal phase, so its implementing PR also carries the
  t4 close-out (Backlog Impact, doc currency, the p2 + t4 +
  milestone-row terminal flips).
- **Cross-phase coordination** is via the parent task plan (the
  D1–D5 decisions and its Cross-Cutting Invariants), not
  pre-locked cross-phase contracts, per
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "Cross-PR coordination." p1's recorded handoff
  (`loadActiveWorkInstances` + entry struct extended to carry
  `work_instances.id`) is **verified against merged code at this
  drafting** (SD1/SD2 `Verified by:`), per
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "Cross-PR coordination" (verify the recorded assumption at the
  consuming phase's drafting).
- **Parent-doc update owed by this p2 drafting.** Per
  [`task-plan.md`](../../../../../spec/planning/task-plan.md)
  "Just-in-time scoping and plan drafting" (drafting updates the
  parent tracking table in the same change), this drafting change
  updates the milestone [`README.md`](../README.md) Task Status
  `t4-p2` row to `Proposed` and reconciles its stub-prose (p2 is
  now a drafted phase plan; t3 remains a stub) — mirroring the
  p1 precedent. The t4 *task-plan*-level obligations (milestone
  t4 row, the two deferral-note resolutions, the backlog-split
  *decision*) were discharged in earlier t4 PRs and are not
  re-touched; the backlog-file *mutations* are p2 implementation,
  not p2 drafting.

## Reality-check inputs

Load-bearing codebase claims the plan's promotion gate must
re-confirm against the branch at drafting time (each a
one-sentence falsifier):

- *"`loadActiveWorkInstances` selects only `slug, actor` (no
  `id`), so p2 must extend it and the entry struct to carry
  `work_instances.id`."* Falsifier: it already selects `id`, or
  `RosterEntry` already carries it. Check:
  [`loadActiveWorkInstances` / `RosterEntry` in site.go](../../../../../internal/site/site.go).
- *"`events.work_instance_id` references `work_instances.id` (the
  generated `TEXT PRIMARY KEY`), indexed; not keyed by actor or
  `(slug, actor)`."* Falsifier: the event log keys on actor.
  Check: [`schema.go`](../../../../../internal/db/schema.go)
  (`events.work_instance_id`, `idx_events_work_instance_id`,
  `work_instances.id PRIMARY KEY`);
  [`insertRegister` in handlers.go](../../../../../internal/api/handlers.go)
  writes the generated `wid` as both `work_instances.id` and the
  register event's `work_instance_id`.
- *"`register`/`heartbeat`/`state_transition` events each store
  `nullableJSON(metadata)` into `events.metadata` tagged by
  `type`, so register-baseline + latest-later overlay is
  derivable from event rows; an absent blob is `NULL`."*
  Falsifier: metadata is dropped, or untyped. Check:
  [`insertRegister`/`insertHeartbeat`/`insertStateTransition`/`nullableJSON`
  in handlers.go](../../../../../internal/api/handlers.go).
- *"The register metadata write path is fully plumbed
  server-side; only the client doesn't send it."* Falsifier: the
  API drops `metadata`, or `events.metadata` doesn't exist.
  Check: [`RegisterRequest.Metadata` in api.go](../../../../../internal/api/api.go),
  [`insertRegister` in handlers.go](../../../../../internal/api/handlers.go),
  [`events.metadata` in schema.go](../../../../../internal/db/schema.go).
- *"`registerclient.Register` marshals only `{exact_slug,
  actor}`; `runRegister` surfaces a failed attempt to stderr and
  still exits 0."* Falsifier: the client already sends metadata,
  or a failure is silently swallowed. Check:
  [`registerclient.Register` in client.go](../../../../../internal/registerclient/client.go),
  [`runRegister` in cmd/workstream-tracker/register.go](../../../../../cmd/workstream-tracker/register.go).
- *"`ActiveWorkInstance` is in `tree.go`; the forest `node-header`
  renders only `{{.Actor}}` from it, so an additive `ID` field
  doesn't regress the forest."* Falsifier: the forest renders
  another `ActiveWorkInstance` field. Check:
  [`ActiveWorkInstance` in tree.go](../../../../../internal/site/tree.go),
  [`node-header` in forest.go](../../../../../internal/site/forest.go).
- *"Native `<details>`/`<summary>` is the established no-JS idiom
  and `roster.go` owns the roster region body parsed into
  `indexTmpl`."* Falsifier: the roster define lives in
  `render.go`, or the project uses JS for collapse. Check:
  [`roster.go`](../../../../../internal/site/roster.go),
  [design/v0.1-design.md §7](../../../../../design/v0.1-design.md),
  [m2 README Cross-Task Decisions](../README.md).
- *"No build/test wrapper exists; the Go toolchain is the gate."*
  Falsifier: a Makefile/justfile/script wraps build+test. Check:
  repo root has no Makefile/justfile; `scripts/` is only
  `assemble.sh` (matches t1's confirmed finding and the t4 task
  plan Validation Gate).
- *"A render-test harness exists to extend."* Falsifier: no
  render test entry point. Check:
  [`renderRoster` in roster_test.go](../../../../../internal/site/roster_test.go),
  [`renderTree` in render_test.go](../../../../../internal/site/render_test.go).
- *"A populated second plan tree exists for manual render
  observation."* Falsifier: only the dogfood tree exists. Check:
  `docs/plans/demo-workstream/` is a populated second root.
- *"The registration handshake yields a real active session for
  manual observation."* Falsifier: `register` does not create an
  active work-instance reachable by the site. Check: the
  session-start handshake against the running server returned
  `http_status=201` with a real `work_instance_id` for slug
  `workstream-tracker-1-0-m2-t4-p2`; re-running with `--name`
  and a deliberately typoed slug yields the named + unbound
  observation cases.
