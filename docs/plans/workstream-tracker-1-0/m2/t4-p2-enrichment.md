---
slug: workstream-tracker-1-0-m2-t4-p2
Status: Proposed
short_description: Client/CLI metadata + named sessions + event-log join + expandable raw-JSON detail
---

# t4 p2 — Enrichment + named sessions

## Context

This is **phase 2** (the task-terminal phase) of the N ≥ 2 task
plan [`t4-session-roster.md`](t4-session-roster.md)
(`workstream-tracker-1-0-m2-t4`, "Session roster + work-item
enrichment"), drafted just-in-time against merged code per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Just-in-time scoping and plan drafting." Its scoping doc is
[`scoping/t4-p2-enrichment.md`](scoping/t4-p2-enrichment.md).

**What this covers.** p1 shipped a bare roster: every active
session listed and classified bound/unbound, labelled by slug.
p2 makes the roster *useful* — sessions report a human name and
other data, and the roster shows the name plus an expandable
view of everything else the session reported. Today a registered
session is an opaque slug; after p2 a contributor scanning the
roster sees a human-meaningful name and can expand any entry to
see the raw reported detail (PRs, branch, anything the session
chose to send).

**Why now.** p1 made registered-but-unbound work *visible*; it
left every entry an opaque slug because there was no reported-name
source and no client sending one. p2 closes the milestone's
"every session can be accounted for" observability goal: a
contributor can tell *which* session each roster row is, not just
that one exists. It is the second, dependent half of t4 — it
builds directly on p1's loader, shared `indexData` field, and
region.

**What surfaces this touches.** The website's roster region and
its request-time data path (the per-request loader gains an
event-log read; the region body gains the name + expandable
detail); the registration client and CLI (they begin reporting
`metadata` including a name); and contributor-facing docs (the
registration handshake reference and the design doc). **No schema
change** (the reported data rides the existing `events.metadata`
column, read by joining the event log on `work_instances.id`);
**no forest region or shell-layout change** (the v0.1 forest
actor markers do not regress).

The inherited WHAT is locked at task-plan level (t4 Contracts
"Session identity and naming," "Reported data and the raw-JSON
detail view," "Metadata read policy"); the only p2-level
decisions are the pure mechanism, resolved in scoping SD1–SD5
(loader/struct id extension + companion per-request event read;
additive `ActiveWorkInstance.ID`; Go-side baseline+overlay fold;
schema-agnostic client with the CLI owning the `name` convention;
native `<details>` detail with a plain row for no-metadata
entries).

## Goal

Opening the page renders, in t1's roster region, every active
session p1 already lists — but each entry now shows a
**human-useful name** when the session reported one, falling back
to the work-instance **slug** when it did not, and **never** the
internal `wst-<uuid>` actor. An entry whose session reported
metadata is **expandable** to a deliberately-unstructured raw-JSON
view of that reported metadata (the session's PRs and anything
else it sent); an entry whose session reported nothing renders as
a plain row with no disclosure control (enrichment is additive —
absence is not a drop, and the v0.2-minimum session still lists).
The reported data is read by **joining the event log** on
`work_instances.id` — the `register` event's metadata as the
identity baseline with the latest later event's metadata overlaid
key-by-key, per request — with **no schema change**. The register
client/CLI begins sending request `metadata` including the
conventionally-read `name` key, via a CLI affordance, and a
metadata send failure is **surfaced** by the client (registration
stays observable best-effort and opt-in). The forest region, the
shell layout, and v0.1's per-node forest actor markers are
unchanged, and the render path stays walk-on-every-request (no
cache; the per-request event read is a distinct source, not a
duplicate walk).

p2 has no later sibling — it is the task-terminal phase, so its
implementing PR also closes out t4 (Backlog Impact, Documentation
Currency, and the p2 + t4 + milestone-row terminal flips); there
is no sibling-interface handoff to record.

## Contracts

Task-level behavior contracts are owned by the parent task plan
[`t4-session-roster.md`](t4-session-roster.md) Contracts and are
cited, not restated, here per
[`shared.md`](../../../../spec/planning/shared.md)
scoping-vs-duplication discipline (the recursion across
task/phase boundaries in
[`task-plan.md`](../../../../spec/planning/task-plan.md) "How a
phase plan cites its parent task plan"). The bullets below are
p2's own realization shape — the observable end-state each
surface must reach. Estimate-shaped sections (Files to touch,
Execution Steps) are labeled as estimates per
[`shared.md`](../../../../spec/planning/shared.md) "Plan content
is a mix of rules and estimates."

### Reported metadata read (event-log join, no schema change)

- The roster loader surfaces each active work-instance's reported
  metadata by **reading the event log keyed on
  `work_instances.id`** — extending `loadActiveWorkInstances` and
  the active-work-instance struct to carry that id (the
  cross-phase assumption p1 recorded and this plan verifies at
  drafting per the parent task plan's Cross-PR coordination), then
  performing one further **per-request** read over exactly the
  active id set. No `work_instances` column is added; no schema
  migration is introduced. `Verified by:`
  [`loadActiveWorkInstances` in site.go](../../../../internal/site/site.go)
  selects only `slug, actor` today;
  [`events.work_instance_id` / `idx_events_work_instance_id` /
  `work_instances.id` in schema.go](../../../../internal/db/schema.go)
  is the FK + index + primary key the read is keyed on;
  [`insertRegister` in handlers.go](../../../../internal/api/handlers.go)
  writes the generated `wid` as both `work_instances.id` and the
  register event's `work_instance_id`.
- The represented metadata for a session is its `register`
  event's metadata as the **identity baseline**, with the
  **latest** later event's metadata **overlaid key-by-key** (later
  keys win; keys absent from the later event keep the register
  value), resolved **per request**. This realizes the parent task
  plan's "Metadata read policy" contract; the read is the
  mechanism, the policy is the inherited rule. `Verified by:`
  [`insertRegister`/`insertHeartbeat`/`insertStateTransition` in
  handlers.go](../../../../internal/api/handlers.go) each write
  `nullableJSON(...metadata)` into `events.metadata` tagged by
  `type` with a server `received_at`, so baseline (the `register`
  row) and latest-later (max `received_at` among non-`register`
  rows) are both derivable; a `NULL` blob contributes no keys
  (per [`nullableJSON`](../../../../internal/api/handlers.go)),
  not an error.
- The read is **walk-on-every-request-consistent**: it is a
  per-request read of a distinct source (the event log), scoped
  to the active id set and using the existing event index — not a
  cache, file-watch, in-memory accumulation, or a second walk of
  already-walked data. `Verified by:`
  [`Server.index` in site.go](../../../../internal/site/site.go)
  performs the per-request `walkPlans` + `loadActiveWorkInstances`
  the event read composes alongside (parent task plan
  Cross-Cutting Invariants "Walk-on-every-request," cited not
  restated).

### Session identity and naming (realized)

- An entry's display label is the reported metadata **`name`**
  value when present, else the work-instance **slug**, **never**
  the `wst-<uuid>` actor. p2 realizes the reported-name arm of the
  parent task plan's task-level name-then-slug rule (p1 realized
  the slug-fallback arm); the actor remains loader-internal (the
  deterministic secondary sort key only), never rendered in the
  roster. The literal slug-fallback label formatting is
  render-time UX deferred under "Bans on surface require rendering
  the consequence" and observed in the Validation Gate.
  `Verified by:` the parent task plan Contracts "Session identity
  and naming" (the inherited rule);
  [`sessionActor` in cmd/workstream-tracker/register.go](../../../../cmd/workstream-tracker/register.go)
  generates the `wst-<uuid>` string the roster must never show as
  a label.
- `name` is the **sole conventionally-read key**; every other
  reported key stays arbitrary and unschematized. The reported
  shape is **not** schematized toward t3's declared-stages field.
  `Verified by:` the parent task plan Cross-Cutting Invariants
  "Schema-loose, not homogenized" (the milestone posture-tension
  invariant, cited not restated).
- The forest's existing per-node actor markers are **unchanged**.
  p2 adds no edit to the forest region body; the only `tree.go`
  change is an additive, unrendered `ActiveWorkInstance.ID` field
  (scoping SD2) which cannot alter a byte of the forest
  actor-marker output. `Verified by:`
  [`node-header` in forest.go](../../../../internal/site/forest.go)
  ranges `.WorkInstances` emitting only `{{.Actor}}`; no other
  `ActiveWorkInstance` field is rendered there.

### Expandable raw-JSON detail

- An entry whose session has resolved reported metadata is
  **expandable** to a deliberately-unstructured raw-JSON view of
  that metadata, using native HTML `<details>`/`<summary>` (the
  repo's established no-JS idiom; no JavaScript, no new route, no
  read API). The view imposes **no task-side schema** on the
  reported shape. `Verified by:`
  [design/v0.1-design.md §7](../../../../design/v0.1-design.md)
  records the render as native `<details>`/`<summary>`, "no
  JavaScript and no new route";
  [the m2 milestone Cross-Task Decisions "Collapsible mechanism"](README.md)
  locks that idiom; [`roster.go`](../../../../internal/site/roster.go)
  owns the roster region body parsed into `indexTmpl`.
- An entry whose session reported **no** metadata renders as a
  **plain row with no disclosure control** — a deliberate,
  observed state mirroring the forest leaf-box-without-children
  idiom, not an empty `<details>`. The v0.2-minimum session
  (slug + actor only) still lists. `Verified by:` the parent task
  plan Cross-Cutting Invariants "Additive registration data"
  (cited not restated); the consequence is observed in this
  plan's Validation Gate per "Bans on surface require rendering
  the consequence," not asserted from template source.
- The roster still **lists only** — p2 adds no
  promote/dismiss/attach affordance for unbound sessions (the
  `<details>` disclosure is not a triage control). `Verified by:`
  parent milestone [`README.md`](README.md) Cross-Task Risk
  "Reviewer-flag any t4 drafting that adds a promote/dismiss
  affordance"; the roster tests assert no button/form/input
  element is emitted.

### Register client/CLI reports metadata

- The register client begins sending request `metadata`; it stays
  **schema-agnostic** (it forwards an arbitrary metadata value, it
  does not type the `name` key). The CLI owns the `name`
  convention: a `--name` flag with a `WST_NAME` environment
  fallback; when a name is supplied the CLI sends the conventional
  one-key metadata object, when absent it sends **no** metadata
  (not an empty name). `Verified by:`
  [`registerclient.Register` in client.go](../../../../internal/registerclient/client.go)
  marshals only `{exact_slug, actor}` today;
  [`RegisterRequest.Metadata` in api.go](../../../../internal/api/api.go)
  is `json.RawMessage` and
  [`insertRegister` in handlers.go](../../../../internal/api/handlers.go)
  already writes `nullableJSON(req.Metadata)` to `events.metadata`
  — the server side is fully plumbed;
  [`runRegister` in cmd/workstream-tracker/register.go](../../../../cmd/workstream-tracker/register.go)
  reads `--slug`/`--actor`/`--server` via a `flag.FlagSet` + env
  fallback the affordance extends.
- A metadata send failure is **surfaced, not silently dropped**.
  `name` rides the single existing register request (there is no
  separate metadata POST), so a send failure *is* the
  register-attempt failure the existing best-effort path already
  prints to stderr while still exiting success — registration
  stays observable best-effort and opt-in. `Verified by:`
  [`runRegister` in cmd/workstream-tracker/register.go](../../../../cmd/workstream-tracker/register.go)
  prints an explicit actionable stderr line on a failed attempt
  and returns 0; the `error-surfacing-user-mutations` audit is
  satisfied by that surfaced path and observed (not asserted) in
  the Validation Gate.

### Render path and region ownership

- The roster **region body** (template + scoped CSS) stays in
  [`roster.go`](../../../../internal/site/roster.go); p2 extends
  it with the named-session label and the expandable detail.
  The **forest region body and the shell layout are not changed.**
  p2 edits the shared request-time data path — the `site.go`
  loader/companion-read and `render.go`'s `indexData` field stay
  the carved-out, **not** reviewer-flag surface the parent task
  plan Cross-Cutting Invariants "Region/shell boundary" and the
  milestone data-path carve-out permit. `Verified by:`
  [`roster.go`](../../../../internal/site/roster.go) owns
  `{{define "roster"}}` / `{{define "roster-style"}}`;
  [`indexData` in render.go](../../../../internal/site/render.go)
  already carries the p1 `Roster` field (p2 enriches the entry
  struct, not the shell composition); parent task plan Contracts
  "Render path and ownership" + Cross-Cutting Invariants
  "Region/shell boundary" (cited not restated).

## Cross-Cutting Invariants

p2 introduces no new cross-cutting rule. It inherits the parent
task plan's Cross-Cutting Invariants
([`t4-session-roster.md`](t4-session-roster.md) "Cross-Cutting
Invariants") — walk-on-every-request, schema-loose-not-homogenized,
v0.1 forest actor markers do not regress, region/shell boundary,
additive registration data — cited here per the parent-citation
discipline rather than duplicated. The three p2's diff most
directly brushes (reviewer-flag candidates here):

- **Walk-on-every-request** — satisfied by SD1-B: the metadata
  read is one per-request read of a distinct source scoped to the
  active id set; `buildRoster` stays a pure function over the
  per-request values.
- **Schema-loose, not homogenized** — satisfied by SD4 (client
  stays schema-agnostic) and the raw-JSON detail imposing no
  task-side schema; `name` is a display convention, not a schema.
- **v0.1 forest actor markers do not regress** — satisfied by the
  only `tree.go` change being an additive unrendered
  `ActiveWorkInstance.ID` field; the forest `node-header` renders
  only `{{.Actor}}` and is observed unchanged.

## Files to touch

*Estimate of expected shape per
[`shared.md`](../../../../spec/planning/shared.md) "Plan content
is a mix of rules and estimates"; implementation may revise with
the deviation reported per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Plan-to-PR Completion Gate."*

**Modify (estimate):**

- `internal/site/site.go` — extend `loadActiveWorkInstances` to
  also select `work_instances.id`; add the companion per-request
  event-metadata read over the active id set (register baseline +
  latest-later key-by-key overlay folded in Go); extend
  `RosterEntry` with the resolved name and the raw-JSON detail;
  keep `buildRoster` a pure function over the per-request values;
  wire the event read into `Server.index`. Shared request-time
  data path (the carved-out surface).
- `internal/site/tree.go` — **one additive line**: an `ID` field
  on `ActiveWorkInstance` (scoping SD2). This is a deviation from
  the t4/p1 *estimate* that `tree.go` is "intentionally not
  touched"; it is additive and unrendered, the forest no-regress
  *invariant* holds (the `node-header` renders only `{{.Actor}}`),
  and it is the minimal honoring of p1's recorded
  carry-`work_instances.id` handoff. Plan-flagged here,
  hand-classified under `rename-aware-diff-classification`, and
  reported under `## Estimate Deviations` in the implementing PR.
  `Verified by:`
  [`ActiveWorkInstance` in tree.go](../../../../internal/site/tree.go)
  is the struct `loadActiveWorkInstances` returns;
  [`node-header` in forest.go](../../../../internal/site/forest.go)
  renders only `{{.Actor}}` from it.
- `internal/site/roster.go` — extend the region body: the entry
  label is the name else slug (never the uuid); an entry with
  reported metadata renders as a native `<details>`/`<summary>`
  with the deliberately-unstructured raw JSON in the body, an
  entry with none renders as a plain row with no disclosure;
  extend the roster CSS. The roster region body — t4's owned
  surface.
- `internal/site/roster_test.go` — extend coverage: name-as-label,
  slug fallback when no name, no `wst-<uuid>` rendered, the
  baseline+overlay resolution (later event keys win, register
  keys survive), the expandable detail present for a
  metadata-bearing entry and absent (plain row) for a
  no-metadata entry, no affordance element. t4's owned test
  surface.
- `internal/registerclient/client.go` — `Register` gains an
  arbitrary metadata parameter forwarded verbatim into the
  request `metadata` field (no `name` typing — schema-agnostic).
- `cmd/workstream-tracker/register.go` — a `--name` flag with a
  `WST_NAME` environment fallback; when a name is present
  `runRegister` builds the conventional one-key metadata object
  and passes it to `Register`, when absent it sends no metadata;
  the existing surfaced-failure / exit-0 best-effort path is
  unchanged.
- `cmd/workstream-tracker/*_test.go` / `internal/registerclient/*_test.go`
  — register-client/CLI coverage for the metadata-forwarding and
  the name-present / name-absent paths, if a test entry point
  exists for these packages (estimate; the implementing PR
  confirms the existing test surface).
- `docs/agents/local/session-registration.md` — the handshake
  reference notes the optional session name (`--name` / `WST_NAME`)
  in its flag list. Documentation Currency.
- `design/v0.1-design.md` — §7 (the roster surface: named
  sessions + expandable raw-JSON detail) and §5 (the roster reads
  reported data by joining the event log; no schema change).
  Documentation Currency.
- `docs/backlog.md` — the `deterministic-interactive-registration`
  split + the two new entries (Backlog Impact). t4 task-terminal.
- `docs/plans/workstream-tracker-1-0/m2/t4-p2-enrichment.md`,
  `docs/plans/workstream-tracker-1-0/m2/t4-session-roster.md`,
  `docs/plans/workstream-tracker-1-0/m2/README.md` — the
  task-terminal Status flips (p2 + t4 → `Landed`; milestone
  `t4-p2` + `t4` rows → `Landed`).

**Intentionally not touched** *(estimate — where we don't expect
changes, not a hard prohibition):*

- `internal/site/forest.go`, `internal/site/walker.go` — forest
  region body and walker unchanged (the v0.1 forest-actor
  no-regress invariant). `tree.go` is touched only for the
  additive `ID` field above — it is **not** in this
  not-touched list (the estimate deviation is named, not hidden).
- `internal/site/render.go` — the `indexData` `Roster` field and
  the shell composition already exist from p1; p2 enriches the
  entry struct in `site.go`, not the shell. No `render.go` edit
  expected.
- `internal/db/schema.go` — no schema change (event-log join);
  explicitly rejected task-wide (t4 scoping D1).
- `internal/api/*` — the register/event metadata write path is
  already plumbed; no API change expected (SD4 `Verified by:`).
- Any caching/file-watch layer — none exists and none is
  introduced (walk-on-every-request).

## Validation Gate

The Go toolchain is the gate (no build/test wrapper exists;
`Verified by:` repo root has no Makefile/justfile, `scripts/` is
only `assemble.sh` — matches t1's confirmed finding and the
parent task plan Validation Gate). p2's gate:

- `gofmt -l internal cmd` reports no files; `go build ./...`,
  `go vet ./...`, `go test ./...` all clean, including the new
  roster and register-client/CLI cases.
- New automated coverage asserts: a session reporting a `name`
  renders that name as the entry label and no `wst-<uuid>`
  substring is rendered as a label; a session reporting no name
  falls back to the slug label; the baseline+overlay resolution
  (a later event key overrides the register value; a register key
  absent from the later event survives); an entry with reported
  metadata emits the `<details>`/`<summary>` detail and an entry
  with none emits a plain row and no `<details>`; no
  button/form/input affordance is emitted; the register client
  forwards supplied metadata and sends none when no name is given.
- **Manual render observation** (Bans on surface require rendering
  the consequence; observed on the real render, not asserted from
  template source — `validation-honesty`). Run the server
  (`go run ./cmd/workstream-tracker`, port 8080) against the
  **dogfood tree** and the **`demo-workstream` tree**, with: a
  **named bound** session (the session-start handshake
  re-registered with `--name` for slug
  `workstream-tracker-1-0-m2-t4-p2`, a walked plan doc); a
  **named unbound** session (a deliberately typoed slug,
  `--name`); a **no-name** session (registered without `--name`);
  and a **no-metadata / v0.2-minimum** session. Visually confirm
  in a real browser/curl render: the named sessions show their
  reported name and **no** `wst-<uuid>` appears as a label; the
  no-name session shows the slug fallback as a designed state;
  an entry with metadata expands to the raw-JSON detail and the
  no-metadata entry is a plain row with no disclosure; the
  empty/bare-minimum session still lists; the forest region still
  renders its existing per-node actor markers unchanged; one page
  scroll preserved (shell layout unchanged). Then induce a client
  metadata send failure (register against a stopped/unreachable
  server) and confirm the client surfaces it on stderr and the
  session still proceeds (best-effort).

## Execution Steps

*Estimate of ordering; deviating from a step is an estimate
deviation, not a contract breach.*

1. **Baseline validation.** On a clean tree, run the Validation
   Gate toolchain commands and confirm green *before* editing, so
   a pre-existing failure isn't misattributed.
2. **Branch hygiene.** Implement on a dedicated branch off
   current `main`; no unrelated changes ride along. Keep the
   plan-drafting commits (scoping + plan + promotion) distinct
   from the implementing commits.
3. **Loader + event read.** Extend `loadActiveWorkInstances`
   (select `id`) and `ActiveWorkInstance` (additive `ID`); add
   the companion per-request event-metadata read + the Go-side
   baseline/overlay fold; extend `RosterEntry`; wire
   `Server.index`. Keep `buildRoster` pure.
4. **Region body.** Extend the `roster.go` template + CSS for the
   name-or-slug label and the native `<details>` detail (plain
   row when no metadata); keep the forest region body and shell
   layout untouched.
5. **Register client/CLI.** Add the schema-agnostic metadata
   parameter to `registerclient.Register`; add `--name` /
   `WST_NAME` to `runRegister` (no metadata when absent).
6. **Tests.** Extend `roster_test.go` and the register-client/CLI
   tests for the coverage above.
7. **Self-review.** Run the Self-Review Audits below against the
   diff at commit boundaries.
8. **Final validation.** Re-run the full Validation Gate
   including the manual render observation and the client
   send-failure observation (observe the consequences, don't
   assume them).
9. **t4 close-out + PR preparation.** Execute the Backlog Impact
   mutations and Documentation Currency updates; flip p2 and t4
   to `Landed` and advance the milestone `t4-p2` + `t4` rows; PR
   body carries `## Review Stance`, `## Documentation`,
   `## Estimate Deviations` (the `tree.go` additive field), and a
   `## Test plan` checklist, per the Plan-to-PR Completion Gate.

## Self-Review Audits

From the parent task plan's Self-Review Audits
([`t4-session-roster.md`](t4-session-roster.md)), the full set
maps to p2's diff surfaces:

- **validation-honesty** — the name/slug-fallback render, the
  expandable-vs-plain-row consequence, and the
  baseline+overlay resolution are "passed" only when observed in
  a real render / exercised by a real test, not asserted from
  template or query source.
- **error-surfacing-user-mutations** — the register client/CLI
  metadata reporting must surface a send failure observably; a
  silent drop is the defect. Re-enters scope at p2 (it had no p1
  surface).
- **readiness-gate-truthfulness** — the p2 + t4 `Proposed →
  Landed` flips and the milestone `t4-p2` + `t4` row advances
  only when every Goal/Contract/Validation item is satisfied or
  explicitly deferred in the plan.
- **rename-aware-diff-classification** — the additive
  `ActiveWorkInstance.ID` field and the roster template
  extension must be hand-classified so the "forest unchanged" and
  "shell unchanged" claims are real, not tooling-fooled by
  add+delete.

## Out Of Scope

- Triage **action** (promote/dismiss/attach for unbound
  sessions) — listing only; deferred past 1.0 (milestone Out of
  Scope; reviewer-flag any phase adding it). The `<details>`
  disclosure is a read affordance, not a triage control.
- Humanizing the **forest** actor display / changing the
  `wst-<uuid>` generator — captured as a new backlog entry in
  this PR (Backlog Impact), not done here; the forest region
  body is unchanged.
- Any `work_instances` schema column / migration — explicitly
  rejected task-wide in favor of the event-log join (t4 scoping
  D1).
- Schematizing the reported data toward t3's declared-stages
  field — the posture-tension invariant names this as the
  defect, not the fix.
- Sending reported fields beyond `name` from this CLI — the
  raw-JSON detail *displays* arbitrary reported fields, but the
  p2 producer surface is the `name` convention only; richer
  producer reporting is not scoped here.

## Risk Register

- **Per-request event read cost.** Reading metadata per page
  load grows with event volume. Mitigation: the read is scoped to
  the active id set and keyed by the indexed
  `idx_events_work_instance_id`, consistent with the existing
  per-request load posture (parent task plan Risk Register
  "Per-request events join cost"). `Verified by:`
  [`idx_events_work_instance_id` in schema.go](../../../../internal/db/schema.go).
  Carried, not blocking.
- **Non-object reported metadata.** The schema-loose posture
  permits any JSON; a non-object blob has no defined key-by-key
  merge. Mitigation: the raw-JSON detail renders it verbatim and
  the conventional `name` read simply finds no key — an
  impl-robustness edge handled by rendering, not by imposing a
  schema (that would breach the posture-tension invariant). The
  Validation Gate's manual observation includes the conventional
  object case; the non-object case degrades to "shown verbatim,
  slug label."
- **The uuid leaks into the roster UX via the name path.** A
  naive fallback could surface `wst-<uuid>`. Mitigation: the
  task-level naming contract bans it; a roster test asserts no
  `wst-`-shaped label and the manual observation confirms the
  no-name fallback is the slug.
- **The schema-loose view homogenized under review pressure.** A
  reviewer/phase may "tidy" the raw JSON toward t3's schema.
  Mitigation: the posture-tension invariant names this as the
  defect, not the fix; the detail view imposes no task-side
  schema.
- **Backlog split / determinism-workstream ordering.** If the
  determinism workstream graduates the entry before t4 lands, the
  split must compose with that. Mitigation: t4 owns only the new
  observability entry; the determinism entry is reduced in-place
  and stays `Open` regardless of graduation order (parent task
  plan Backlog Impact).

## Documentation Currency

- Milestone [`README.md`](README.md) — the Task Status `t4-p2`
  row and the stub-prose are reconciled **in this drafting
  change** (`t4-p2` row → `Proposed`; prose: p2 is now a drafted
  phase plan, t3 remains a stub) per
  [`task-plan.md`](../../../../spec/planning/task-plan.md)
  "Just-in-time scoping and plan drafting," mirroring the p1
  precedent. The implementing PR advances the `t4-p2` **and**
  `t4` rows and the p2 + t4 plan Status onward to `Landed` (the
  N ≥ 2 task-plan terminal flip co-locates with the last phase
  per
  [`task-plan.md`](../../../../spec/planning/task-plan.md) "Task
  plan terminal state when N ≥ 2").
- [`docs/agents/local/session-registration.md`](../../../agents/local/session-registration.md)
  — the handshake reference notes the optional session name
  (`--name` / `WST_NAME`); updated in p2's implementing PR (the
  t4 Documentation Currency names this doc).
- [`design/v0.1-design.md`](../../../../design/v0.1-design.md)
  §7/§5 — §7's roster sentence is updated from "deliberate
  placeholder / later task" to the shipped roster (bound/unbound,
  named sessions, expandable raw-JSON detail); §5 notes the
  roster reads reported data via the event-log join with no
  schema change. Updated in p2's implementing PR (the t4
  Documentation Currency names this doc).
- This plan's `Status` lifecycle: `In draft` while drafted
  (committed first), `Proposed` after the promotion-gate
  self-review walk, then `Proposed → Landed` in the implementing
  PR per the Plan-to-PR Completion Gate (Validation Gate fully
  satisfiable pre-merge — no post-merge gate, single implementing
  PR — so the default same-PR `Landed` flip applies; as the
  task-terminal phase, the same PR flips the parent t4 task plan
  to `Landed` too).

## Backlog Impact

Per [`spec/backlog.md`](../../../../spec/backlog.md) effect
taxonomy. The effect decisions were recorded at t4 task-plan /
scoping drafting (t4 Backlog Impact, scoping D2); p2's
implementing PR executes the **backlog-file mutations** (it is
the t4 task-terminal PR):

- [`deterministic-interactive-registration`](../../../backlog.md#deterministic-interactive-registration)
  — **split.** Reduce the entry to the **determinism gap only**
  (registration circularity, sole-consumer compensation, the
  neighborly-events tripwire); it stays `Open` and is the entry
  the independent determinism workstream graduates. Remove the
  observability-gap prose it currently also carries.
- **New entry — observability residual** (`Open`). Captures the
  residual t4 narrows: registered-but-unbound work is now
  surfaced by the roster, so the residual is *unregistered* work
  only; the tree-side-heuristic candidate stays deferred under
  the same neighborly-events tripwire.
- **New entry — humanize forest actor display** (`Open`). The
  forest still renders the raw `wst-<uuid>`; the roster now shows
  names. Replace the raw uuid in the forest UX (and/or revisit
  the actor generator). Scope-framed, not prescribed.

## Related Docs

- [`t4-session-roster.md`](t4-session-roster.md) — parent N ≥ 2
  task plan; owns the task-level Contracts, Cross-Cutting
  Invariants, Phase Contracts, Backlog Impact, and Documentation
  Currency p2 inherits and cites.
- [`scoping/t4-p2-enrichment.md`](scoping/t4-p2-enrichment.md) —
  p2 scoping doc (SD1–SD5 with rejected alternatives,
  reality-check inputs).
- [`t4-p1-bare-roster.md`](t4-p1-bare-roster.md) — the landed
  phase 1; its loader, `indexData` `Roster` field, region, and
  the recorded carry-`work_instances.id` handoff p2 builds on.
- [`README.md`](README.md) — parent milestone; Cross-Task
  Invariants/Decisions/Risks and the Task Status `t4-p2` + `t4`
  rows p2 advances.
- [`../../../../spec/planning/task-plan.md`](../../../../spec/planning/task-plan.md),
  [`../../../../spec/planning/shared.md`](../../../../spec/planning/shared.md)
  — the rules this phase plan is structured against.
