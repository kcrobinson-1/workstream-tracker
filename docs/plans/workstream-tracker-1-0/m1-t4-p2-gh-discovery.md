---
slug: workstream-tracker-1-0-m1-t4-p2
Status: Landed
short_description: gh pr list auto-discovery
---

# t4 P2 — `gh pr list` auto-discovery

## Context

This is phase 2 (the last phase) of the t4 task plan
([`m1-t4-expanded-per-node-display.md`](m1-t4-expanded-per-node-display.md)) —
the parent task plan owns the Cross-Phase Decisions,
Cross-Cutting Invariants, and sequencing this phase plan
inherits by reference rather than restating.

P1 (Landed, #15) renders each node's long description and its
author-curated `related_prs` (absolute-URL strings) inline. This
phase makes related PRs **also** auto-discoverable: a single
`gh pr list` call per request finds PRs whose title contains a
node's slug and merges them into that node's rendered PR list,
so a contributor does not have to hand-maintain `related_prs`
for every node. It is drafted now, just-in-time, because P1's
implementing PR has merged (scoping D2/D5) — the just-in-time
point for the phase that introduces the codebase's **first
subprocess shell-out**.

It is being done now because t4's milestone contract names two
related-PR sources ("an optional `related_prs` frontmatter
field, **and** auto-discovers via `gh pr list`"); P1 shipped the
first, P2 ships the second and completes **t4** (the
read-experience track's terminal task — not m1's last task;
`…-m1-t2` remains undrafted). Surfaces touched conceptually: a
new subprocess
boundary in the read path (the `gh` CLI), the site request
handler that augments the tree, and the plan-tree render the
discovered PRs flow into (P1's block, unchanged). No API, DB,
schema, frontmatter, or spec-field surface — P2 adds no new
field.

Deliberation, the D5 spike findings, rejected alternatives, and
the resolved P2 decisions (P2-D1…P2-D4) live in the sibling
scoping doc
([`scoping/m1-t4-expanded-per-node-display.md`](scoping/m1-t4-expanded-per-node-display.md)
"D5 spike — findings and resolved P2 decisions"), which this
plan does not restate. P2's previously-deferred inputs (gh
query/field semantics, dedupe key, timeout) are all resolved
there by the spike; no input remains open, so after the
`In draft → Proposed` promotion-gate self-review this plan is
`Proposed`.

## Goal

On every page render, one `gh pr list` invocation discovers
merged/open PRs and each plan-tree node's rendered related-PR
list additionally includes any discovered PR whose title
contains that node's verbatim slug, merged and deduped with the
node's frontmatter `related_prs` (frontmatter first and
authoritative). When `gh` is unavailable for any reason —
missing binary, not authenticated, no network, not a git repo,
non-zero exit, malformed JSON, or exceeding a bounded timeout —
the page renders exactly P1's output (frontmatter PRs or none)
and never errors. No caching, file-watch, or in-memory build-up
is introduced; the render path stays walk-on-every-request.

## Naming

- `ghPR` — unexported struct decoding the `gh pr list --json`
  array elements: `URL`, `Title`, `Number`.
- `discoverPRsByTitle(ctx) ([]ghPR, error)` — runs the single
  `gh` subprocess and returns the decoded PRs; all failure
  modes return a non-nil error and a nil slice.
- `augmentRelatedPRs(roots, prs)` — walks the node tree and, for
  each node, appends discovered PR URLs whose title contains the
  node's slug, deduped against existing `RelatedPRs`.
- `ghTimeout` — the bounded subprocess timeout constant
  (`5 * time.Second`, scoping P2-D4).
- `ghWaitDelay` — the post-kill I/O-wait bound
  (`1 * time.Second`, set as `Cmd.WaitDelay`). After the
  `ghTimeout` deadline kills `gh`, `Output()` would still block
  until the stdout pipe closes; a killed `gh` whose grandchild
  (e.g. a wrapping shell's `sleep`) still holds that pipe would
  wedge the request well past `ghTimeout`. `WaitDelay`
  force-closes the inherited pipes so the worst case is
  bounded at ~`ghTimeout + ghWaitDelay` (observed in the
  hung-`gh` failure-matrix check; see Estimate-shaped contract
  refinement note in the Subprocess contract).
- `ghPRListLimit` — the `-L` result cap constant (`200`).

## Contracts

Final shapes. Estimate-shaped sections (Files to touch, Commit
Boundaries) are labelled as estimates per
[`shared.md`](../../../spec/planning/shared.md) "Plan content is
a mix of rules and estimates." These bullets state the
observable end-state; technique is non-binding guidance under
Execution Steps.

### Subprocess contract

- Exactly **one** `gh` subprocess runs per page request,
  independent of node count: `gh pr list --state all --json
  url,title,number -L 200`, executed via `exec.CommandContext`
  with a 5-second timeout (`ghTimeout`, scoping P2-D4). The
  command is invoked with a fixed argument vector (no shell, no
  string interpolation of any plan-derived data into the
  command line) — slugs never reach the command line, so there
  is no command-injection surface. `Verified by:` the slug
  match is done in Go against the decoded JSON (Match contract
  below), not by passing the slug to `gh`.
- Every failure path returns "no discoveries," never an error
  to the caller and never a non-200 page: missing `gh` binary
  (exec error), non-zero exit (not a git repo / not
  authenticated / API error), context-deadline timeout,
  empty/malformed JSON. Each is logged once via `slog` at
  warn and the render proceeds with frontmatter PRs only.
  `Verified by:` D5 spike Finding 5 — every mode exits
  detectably (non-repo → `gh` non-zero "fatal: not a git
  repository"; missing → exec error); standard
  `exec`/exit-code/JSON-unmarshal error checks cover all of
  them.
- The timeout kills the subprocess (context cancellation
  propagated by `exec.CommandContext`) **and** `Cmd.WaitDelay`
  (`ghWaitDelay`) force-closes the inherited stdout/stderr pipes
  after the kill, so a hung `gh` cannot wedge a request beyond
  ~`ghTimeout + ghWaitDelay`. **Estimate-shaped contract
  refinement (PR-body Estimate Deviation):** the drafted
  contract attributed the whole bound to `exec.CommandContext`
  alone; the hung-`gh` failure-matrix check observed that
  `exec.CommandContext`'s kill does *not* unblock `Output()`
  while a killed `gh`'s grandchild still holds the stdout pipe
  (request wedged ~30s in the spike repro), so the implementation
  adds `Cmd.WaitDelay` to actually satisfy the contract's
  "cannot wedge beyond the timeout" intent — the durable
  guarantee is unchanged, the mechanism is corrected to what the
  observed behaviour requires. `Verified by:` the hung-`gh`
  manual check (request returned in ~6s, not ~30s, no leaked
  process); the 5-second deadline is scoping P2-D4 (generous for
  a normal local `gh` round-trip, bounded against a hang).

### Match contract

- A discovered PR is attached to a node iff the node's verbatim
  slug is a **case-sensitive substring of the PR title**.
  Matching is done in Go over the decoded `--json` array, never
  via GitHub `--search` (D5 Findings 1–2: `--search` tokenises
  on hyphens and over-matches; in-process substring is exact).
- The discovered value appended is the PR's `url` field — an
  absolute URL, the identical entry shape P1 established
  (scoping D3/D4, D5 Finding 4). No PR-reference
  canonicalization is performed in either phase.
- **Precision boundary (accepted consequence, scoping P2-D2).**
  Discovery finds only PRs whose title verbatim contains the
  slug. PRs whose title uses a short conventional-commit scope
  (e.g. `feat(m1-t4-p1): …`) are **not** discovered. This is a
  known, accepted loss: frontmatter `related_prs` is the
  complete and authoritative source (scoping D4); `gh`
  discovery is an additive convenience. The Validation Gate
  observes this consequence rather than asserting it.

### Merge contract

- For each node the rendered `RelatedPRs` is: the node's
  frontmatter entries in source order (authoritative, scoping
  D4), followed by discovered PR URLs not already present, in
  `gh` result order. Dedupe is plain absolute-URL **string
  equality** (scoping P2-D3, D5 Finding 4); a discovered URL
  already in the frontmatter list is dropped. No reordering or
  rewriting of frontmatter entries.
- Augmentation mutates the in-request `PlanNode.RelatedPRs`
  slices only; nothing is persisted or cached across requests.
  `Verified by:` augmentation runs inside `Server.index` per
  invocation (Integration contract below); m1 Cross-Task
  Invariant "Render path stays walk-on-every-request."

### Integration contract (`internal/site/site.go`)

- `Server.index` calls the discovery + augmentation between
  `buildTree` and `renderIndex`: build the tree (P1 shape),
  attempt discovery, augment the roots in place, render. A
  discovery failure is non-fatal — the same `roots` render
  unaugmented. `Verified by:`
  [`Server.index` in site.go](../../../internal/site/site.go)
  already sequences `walkPlans → loadActiveWorkInstances →
  buildTree → renderIndex` per request; P2 inserts one
  best-effort step before `renderIndex` and changes no existing
  step.
- The render template (P1's `node` block) is **unchanged**: it
  already renders `.RelatedPRs` as the absolute-URL-or-text
  list; discovered URLs flow through the existing block with no
  template edit. `Verified by:`
  [`indexTmpl` in render.go](../../../internal/site/render.go)
  ranges `.RelatedPRs` with the `isURL` link/text split shipped
  in P1 (#15).

## Cross-Cutting Invariants

Inherited from the parent task plan by reference (per
[`task-plan.md`](../../../spec/planning/task-plan.md) "How a
phase plan cites its parent task plan" — cite, don't duplicate):
see
[`m1-t4-expanded-per-node-display.md`](m1-t4-expanded-per-node-display.md)
"Cross-Cutting Invariants" (optional-frontmatter-degrades-
silently, walk-on-every-request, additive-spec-change,
bare-bones-no-JS). The one P2 most directly carries:
**walk-on-every-request** — P2's `gh` call runs inside the
per-request handler with no caching, file-watch, goroutine, or
in-memory build-up (scoping D4 rejected all three). P2 adds no
new spec field, so additive-spec-change is vacuous here; no
JavaScript; the P1 render block is unchanged.

## Files to touch

*Estimate of expected shape — implementation may revise if a
structural call requires it; deviations are reported per
[`task-plan.md`](../../../spec/planning/task-plan.md)
"Plan-to-PR Completion Gate" with the `## Estimate Deviations`
PR-body callout.*

**New:**

- `internal/site/ghprs.go` — `ghPR`, `discoverPRsByTitle`,
  `augmentRelatedPRs`, the `ghTimeout` / `ghWaitDelay` /
  `ghPRListLimit` constants.
- `internal/site/ghprs_test.go` — match/merge/dedupe unit
  tests with a stubbed PR set (table tests; no real `gh`); a
  test that the failure path yields no augmentation and no
  error.

**Modify:**

- `internal/site/site.go` — `Server.index` calls
  `discoverPRsByTitle` + `augmentRelatedPRs` between
  `buildTree` and `renderIndex`, log-and-continue on error.
- `design/v0.1-design.md` §7 — note the `gh pr list`
  auto-discovery source for related PRs (P1 added the inline
  detail; P2 records the second, best-effort source).
- `docs/plans/workstream-tracker-1-0/m1-t4-p2-gh-discovery.md`
  — this plan: `Proposed → Landed` in the implementing PR.
- `docs/plans/workstream-tracker-1-0/m1-t4-expanded-per-node-display.md`
  — the task plan: `In progress → Landed` (P2 is the last
  phase) in the implementing PR.
- `docs/plans/workstream-tracker-1-0/m1-v0-2.md` — the parent
  milestone **t4 row only** → `Landed` (mirrors the task plan).
  No milestone Status / Backlog / Documentation-Currency
  reconciliation here — m1 is not terminal (see Documentation
  currency).
- `docs/plans/workstream-tracker-1-0/scoping/m1-t4-expanded-per-node-display.md`
  — **NOT deleted in this PR.** Per
  [`task-plan.md`](../../../spec/planning/task-plan.md) path
  conventions the `scoping/` contents delete in batch at the
  **m1-terminal PR** (with sibling `m1-t1-*` / `m1-t3-*`
  scoping docs); m1 still has `…-m1-t2` undrafted, so deletion
  defers to that later PR. P2's PR leaves the scoping doc in
  place.

**Intentionally not touched** *(estimate — where we don't
expect changes, not a hard prohibition)*:

- `internal/site/render.go`, `tree.go`, `walker.go` — P1's
  render block, the `RelatedPRs` carry, and frontmatter parsing
  are unchanged; P2 only augments the slice before render.
- `internal/api/*`, `internal/db/*` — no API/DB/schema.
- `spec/planning/shared.md` — P2 adds no frontmatter field.
- Any cache/file-watch/goroutine layer — none introduced
  (walk-on-every-request).

## Validation Gate

No build wrapper exists. `Verified by:` repo root has no
`Makefile`/`justfile`; `scripts/` is only `assemble.sh`
(unrelated). The Go toolchain is the gate (per
[`docs/dev.md`](../../dev.md)):

- `gofmt -l internal` reports no files.
- `go build ./...` succeeds.
- `go vet ./...` clean.
- `go test ./...` passes, including new `ghprs_test.go` cases:
  match is case-sensitive substring on title; merge is
  frontmatter-first; dedupe is exact URL-string equality; a
  discovery error yields the unaugmented tree and no error.
- **Failure-matrix manual checks against real environments**
  (the distinct, environment-dependent gate that warranted P2
  being its own phase — scoping D2). Each is *observed*, not
  asserted from source:
  1. **`gh` authenticated, repo with PRs:** a node whose slug
     verbatim appears in a PR title shows that PR URL appended
     after its frontmatter entries; render still single page,
     no JS.
  2. **Precision boundary (scoping P2-D2):** a PR whose title
     uses a short scope (`feat(m1-t4-p1): …`) is **not**
     auto-listed on the `…-m1-t4` node — observe the
     under-match, do not assume it.
  3. **`gh` missing:** run the server with `gh` not on `PATH`;
     the page renders frontmatter PRs only, one warn log, HTTP
     200.
  4. **Not a git repo / no remote:** run from a dir `gh`
     rejects; same graceful outcome (D5 Finding 5).
  5. **Hung `gh`:** stub a `gh` that sleeps past `ghTimeout`;
     the request returns within ~5s with frontmatter-only and a
     timeout warn — confirm the page is not wedged.
  6. **Dedupe:** a PR URL listed in a node's frontmatter AND
     discovered by `gh` appears **once**.

## Execution Steps

Deviating from a step is an estimate deviation (PR-body
callout), not a contract breach.

1. **Baseline validation.** On a clean tree, run the full Go
   gate green *before* editing.
2. **Branch hygiene.** Dedicated branch off current `main`; no
   unrelated changes.
3. **Discovery core.** Add `ghprs.go`: `ghPR`,
   `discoverPRsByTitle` (one `exec.CommandContext`, 5s, fixed
   arg vector, all-failures-return-error), the constants. Unit
   tests with stubbed data (no real `gh` in tests).
4. **Match + merge.** Add `augmentRelatedPRs` (case-sensitive
   title-substring match; frontmatter-first; URL-equality
   dedupe). Table tests.
5. **Integration.** Wire `Server.index` to discover + augment
   between `buildTree` and `renderIndex`, log-and-continue on
   error. No other handler change.
6. **Docs.** `design/v0.1-design.md` §7 note the gh source.
7. **Self-review.** Run the Self-Review Audits below.
8. **Final validation.** Full Go gate + the six failure-matrix
   manual checks observed in real environments (steps 3–5 of
   the matrix need a no-`gh` / non-repo / hung-`gh` setup —
   actually run them).
9. **t4 task-terminal close-out (this is t4's task-terminal
   PR, NOT the m1-milestone-terminal PR).** Per
   [`task-plan.md`](../../../spec/planning/task-plan.md)
   "Plan-to-PR Completion Gate" / "Task plan terminal state
   when N ≥ 2": flip this plan `Proposed → Landed`; flip the
   task plan `In progress → Landed`; flip the `m1-v0-2.md` t4
   **row** → `Landed`. **Do NOT** delete the scoping doc and
   **do NOT** touch milestone Status / Backlog / Documentation
   Currency: m1 is not terminal — its Task Status table still
   carries `…-m1-t2` at `—` (undrafted), and per
   [`task-plan.md`](../../../spec/planning/task-plan.md) path
   conventions the `scoping/` batch deletion + milestone
   reconciliation happen at the **m1-terminal PR** (whichever
   PR lands m1's last remaining task), not here. Leaving the
   scoping doc in place is correct, not an omission.
10. **PR preparation.** PR body carries `## Estimate
    Deviations` (or `N/A`) and reconciles estimate-shaped plan
    sections with what shipped.

## Commit Boundaries

*Estimate of cohesive review chunks — implementer may refine.*
Single implementing PR (this phase = 1 PR; P2 is t4's last
phase). Expected commits: (a) `ghprs.go` discovery core +
tests; (b) match/merge/dedupe + `site.go` integration + tests;
(c) docs §7 + **t4 task-terminal** close-out — the three
Status flips only (P2 → Landed, task plan → Landed,
`m1-v0-2.md` t4 row → Landed). Commit (c) does **not** delete
the scoping doc or reconcile milestone Status/Backlog/Doc-
Currency: those are m1-terminal actions deferred to the later
m1-terminal PR (see Step 9 / Out Of Scope — `…-m1-t2` is
undrafted, m1 is not terminal). The order keeps each commit
building and test-green.

## Self-Review Audits

From
[`docs/agents/local/self-review-catalog.md`](../../agents/local/self-review-catalog.md),
run at step 7:

- **effect-cleanup** (subprocess/timeout surface) — P2 opens a
  subprocess with a context deadline; verify the
  `exec.CommandContext` + timeout actually kills a hung child
  and no goroutine/process leaks per request (the failure-
  matrix hung-`gh` check observes this, not asserts it).
- **error-surfacing** (subprocess failure surface) — every
  `gh` failure mode degrades to frontmatter-only with exactly
  one `slog` warn, never a 500 or a silent swallow with no
  log; the page contract holds on all five failure modes.
- **validation-honesty** (validation surface) — the
  failure-matrix checks count as passed only if actually run
  in the no-`gh` / non-repo / hung-`gh` environments, not
  reasoned from source.
- **readiness-gate-truthfulness** (Status / parent-doc /
  terminal-close-out surface) — the **three** t4 task-terminal
  flips (P2 → Landed, task plan → Landed, `m1-v0-2.md` t4 row →
  Landed) happen only after every Goal/Contract/Validation item
  is satisfied or deferred *in a plan*. The audit guards drift
  in **both** directions here: skipping a required flip, **and**
  performing m1-terminal actions early — deleting the scoping
  doc or flipping milestone Status/Backlog/Doc-Currency at this
  PR is wrong (m1 is not terminal; `…-m1-t2` is undrafted).
  Leaving the scoping doc in place is the correct outcome, not
  an omission to "fix."

rename-aware-diff-classification and trigger-map-currency have
no matching surface (no renames, no directory restructure).

## Out Of Scope

Final boundary calls (deliberation in the scoping doc):

- **Discovering short-scope PRs** (titles like
  `feat(m1-t4-p1): …` lacking the verbatim slug). Inherent
  under-match (D5 Finding 3 / scoping P2-D2); frontmatter
  `related_prs` is the authoritative complete source.
- **Caching / background refresh / per-node fan-out.** One
  per-request subprocess only; caching is the m1 Cross-Task
  Invariant ban (scoping D4).
- **PR-state styling** (open vs merged vs closed visual
  treatment). P2 lists discovered PR URLs through P1's existing
  block unchanged; richer PR presentation is later-milestone.
- **Mandating PR-title or branch conventions** to improve
  discovery precision — not P2's remit (scoping P2-D2).
- **m1-milestone-terminal close-out.** P2's PR is t4's
  *task-terminal* PR, not m1's milestone-terminal PR (m1's Task
  Status table still carries `…-m1-t2` at `—`, undrafted). The
  sibling `scoping/` batch deletion (`m1-t1-*` / `m1-t3-*` /
  `m1-t4-*` together) and milestone Status / Backlog /
  Documentation-Currency reconciliation are out of scope for
  P2's PR and defer to whichever PR lands m1's last remaining
  task, per
  [`task-plan.md`](../../../spec/planning/task-plan.md) path
  conventions. Leaving the t4 scoping doc in place after P2 is
  correct, not an omission.

## Risk Register

- **`gh` latency on the request path.** A normal `gh pr list`
  round-trip adds latency to every render. Accepted: single
  local user, per-request (scoping D4); bounded by the 5s
  `ghTimeout` so worst case is bounded, and a failure (incl.
  timeout) degrades to P1's instant frontmatter-only render.
- **Precision boundary surprises a viewer** (expected PR not
  auto-listed). Intended, not a regression (scoping P2-D2);
  frontmatter `related_prs` is the authoritative path. Called
  out so review doesn't read the under-match as a bug; the
  Validation Gate observes it.
- **Subprocess/timeout leak.** Mitigated by
  `exec.CommandContext` kill-on-deadline; the effect-cleanup
  audit + hung-`gh` manual check verify no wedged request or
  leaked process.
- **`gh --json` field rename on a future `gh` upgrade.** The
  decode targets `url,title,number`; a field rename yields
  empty/garbled discovery, which degrades to frontmatter-only
  (no crash) and is caught by the authenticated manual check.

## Documentation currency

- `design/v0.1-design.md` §7 — the `gh` auto-discovery source
  note lands in this implementing PR.
- This plan's `Status` is `Proposed` (promotion-gate
  self-review complete); it flips `Proposed → Landed` in the
  implementing PR.
- `m1-t4-expanded-per-node-display.md` (task plan) flips
  `In progress → Landed` in this PR — P2 is the last phase.
- `m1-v0-2.md` — **t4 row only** → `Landed` (mirrors the task
  plan). This PR is **not** the m1-terminal PR: m1's Task
  Status table still carries `…-m1-t2` at `—` (undrafted), so
  milestone Status / Backlog / Documentation-Currency
  reconciliation is **not** performed here — it happens at the
  later m1-terminal PR.
- `scoping/m1-t4-expanded-per-node-display.md` — **NOT
  deleted** in this PR. Per
  [`task-plan.md`](../../../spec/planning/task-plan.md) path
  conventions the `scoping/` contents delete in batch at the
  **m1-terminal PR** (with sibling `m1-t1-*` / `m1-t3-*`
  scoping docs); m1 is not terminal, so the doc stays in place
  and its references remain live. No link neutralization is
  needed in this PR.

## Backlog Impact

None. No backlog entry graduates, is deleted, split, or shifts
(see the parent task plan's Backlog Impact). The discovery
precision boundary (scoping P2-D2) is an accepted consequence,
not a captured backlog item; if richer slug→PR association is
wanted post-1.0 it would be raised then.

## Related Docs

- [`m1-t4-expanded-per-node-display.md`](m1-t4-expanded-per-node-display.md)
  — parent task plan (Cross-Phase Decisions, Cross-Cutting
  Invariants, sequencing inherited by reference).
- [`scoping/m1-t4-expanded-per-node-display.md`](scoping/m1-t4-expanded-per-node-display.md)
  — sibling scoping doc; D5 spike findings + P2-D1…P2-D4
  (transient; deleted at the later m1-terminal PR, not by
  P2's PR).
- [`m1-t4-p1-inline-detail-render.md`](m1-t4-p1-inline-detail-render.md)
  — P1 (Landed, #15); supplies the `RelatedPRs` carry and the
  render block P2 augments.
- [`m1-v0-2.md`](m1-v0-2.md) — parent milestone.
- [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md)
  — the rules this plan and its terminal close-out follow.
