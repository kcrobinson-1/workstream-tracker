---
slug: templ-render-adoption
Status: Landed
short_description: Resolve templ vs html/template — accept html/template
---

# templ Render Adoption — Decision

## Context

The workstream-tracker website renders a plan-tree "forest" as
server-rendered HTML. The original v0.1 design locked a community
templating library (`templ`) for that rendering, but the
implementation deliberately shipped on the Go standard library's
`html/template` instead, recording that `templ` would be revisited
"once there are real reusable components." A later change split
the renderer into composable sub-templates (a page shell that
composes a forest region, a recursive node partial, and a roster
region) — arguably the "real reusable components" moment the
deferral named. This task resolves that deferred decision: adopt
`templ` for its type-safe component model, or formally accept
`html/template` as the standing choice.

It is being done now because the deferred decision has a named
trigger that has plausibly fired, and a standing design doc that
records the choice as "deferred, see backlog." Leaving it open
means every future render change re-litigates the same question
and the design doc stays in a provisional state. The surfaces
this touches are conceptual: the website's HTML rendering layer,
the project's design record, and the backlog — no API, schema,
or lifecycle behavior.

The deliverable is the recorded decision plus rejected
alternatives, driven by a throwaway spike. This is a **doc-only
decision plan**: it legitimately skips the `Contracts` and
`Files to touch — new` code-contract sections (a decision-plan
variance disclosed in the implementing PR body per the
[`shared.md`](../../../spec/planning/shared.md) "Section variance
disclosure" rule); the only file edits are doc reconciliation.

## Goal

Resolve the `templ` vs. `html/template` question with a
spike-grounded recommendation, surface it for maintainer review
through this plan and its implementing PR, and reconcile the
design record to the outcome.

Verifiable when: the decision and its rejected alternatives are
recorded here with spike evidence; `design/v0.1-design.md` §10
records the resolved standing choice (no longer "deferred, see
backlog") and its re-trigger condition; the backlog entry is
graduated; and the existing render falsifier tests pass
unchanged (no code migration ships).

## Decision

**Formally accept the Go stdlib `html/template` as the standing
HTML-rendering choice (option B2: accept with a documented
re-trigger condition). Do not migrate to `templ` now.**

The spike (see Spike Evidence) shows migration is not
output-neutral: `templ` minifies insignificant whitespace by
design, so the existing byte-identity falsifier pins fail under
any `templ` port and would have to be deliberately rewritten.
Migration also introduces a code-generation CLI plus a runtime
dependency, and a `templ generate` step into a validation gate
that is deliberately the Go toolchain alone. Against that cost,
the value (compile-time-checked components) is marginal at the
current surface — three thin presentational templates, one
recursive, with the roster still a placeholder. The original
deferral's bar ("real reusable components") is not yet met by a
file split that produced thin, stable partials.

**Documented re-trigger condition.** The tripwire is **render-
surface complexity**, not interactivity. The website is
server-rendered by design — real-time updates without refresh
and server-to-agent push are explicitly out of scope (see
[`design/v0.1-design.md`](../../../design/v0.1-design.md) §9),
so the v1.0+ roadmap (intent bar, tier sorting, sub-stage
D/P/I/V cells, triage zone, actor lineage) grows *more
server-rendered structure*, not client interactivity. Those
features are all expressible in `html/template`; the cost of
staying is developer-experience drag that accretes roughly
linearly with the surface. Revisit `templ` when that drag
becomes load-bearing: the partial set, recursion/conditional
depth, and amount of typed data threaded through the templates
grow such that the stringly-typed `{{template}}`/`FuncMap`
indirection and the absence of compile-time field checks are a
real readability and refactor-safety hazard (a rename of a
`PlanNode` field silently mis-rendering rather than failing to
compile is the canonical symptom). A growing file count alone,
absent that indirection/typed-data pressure, is **not** a
re-trigger; "interactive components" is explicitly **not** the
tripwire, because this UI is not headed there.

This decision is recorded for maintainer review via this plan's
implementing PR; it is not self-applied silently — the PR is the
review surface (there is no parent epic/milestone, so the
graduated backlog entry is this task's only tracking surface).

## Options Considered

Per [`shared.md`](../../../spec/planning/shared.md) "Decompose
options into shapes before analyzing," each top-level option was
split into sub-shapes before scoring.

- **A1 — full migration to `templ`.** Rejected. The spike's
  byte-identity miss means the falsifier pins must be
  deliberately rewritten; the type-safety benefit is marginal at
  three thin templates; adds a codegen CLI + runtime dependency
  and a generate step to a Go-toolchain-only gate.
- **A2 — hybrid (`templ` bodies + string shell, or vice
  versa).** Rejected. Worst of both: two render systems to
  maintain and reason about, the new dependency cost without a
  clean single idiom, and the same byte-identity break.
- **A3 — migrate but preserve exact bytes.** Rejected as not
  viable. `templ`'s whitespace normalization is intentional and
  cannot be configured to reproduce an `html/template`
  `{{range}}` whitespace artifact (`</li>\n  \n</ul>`); the
  spike confirms the collapse. There is no shape of A that keeps
  the existing pins byte-for-byte.
- **B1 — accept `html/template` silently.** Rejected. Recording
  the choice without a re-trigger condition reproduces exactly
  the ambiguity this task closes: a future render change would
  re-litigate the question with no recorded bar. Naming the
  re-trigger is cheap and is the whole point of "formally
  accept."
- **B2 — accept with a documented re-trigger condition.**
  **Chosen.** Lowest cost, no code churn, no falsifier rewrite,
  no new dependency; converts a provisional "deferred" design
  state into a standing decision with a concrete, falsifiable
  condition for reopening.

## Spike Evidence

Throwaway spike per [`task-plan.md`](../../../spec/planning/task-plan.md)
"Spike before plan for novel mechanisms" (`templ` is a new
framework idiom for this codebase). Branch `spike/templ-render`,
commit `d19445b`, **not merged**; spike code never promoted into
the implementing PR.

- `templ` v0.3.1020; CLI installed via `go install
  github.com/a-h/templ/cmd/templ`; runtime module
  `github.com/a-h/templ`.
- A faithful port of the `forest` and recursive `node` templates
  to `templ` components was written and `templ generate` run to
  emit the generated Go.
- Rendering the exact input the byte-identity falsifier pins
  (a `Proposed` root with one `Landed` child):
  - `templ` output (234 bytes): a single line with whitespace
    collapsed —
    `<div class="root"><span class="badge status-proposed">Proposed</span>…<ul><li>…</li></ul></div>`.
  - `html/template` pinned output (218 bytes): newlines, `  `
    indentation, and the `</li>\n  \n</ul>` text-template range
    artifact.
  - Result: **byte-identity miss.** `templ` minifies
    insignificant whitespace by design; the text-template range
    artifact cannot be reproduced.
- Dealbreaker found: the byte-identity falsifier pins
  (`TestRenderNoFieldNodeUnchanged`, `TestRenderTwoRegionShell`)
  cannot survive a `templ` port unchanged, and the project's
  validation gate has no codegen step to host `templ generate`.

Verified by: spike branch `spike/templ-render` @ `d19445b`
(reproduce with `go test ./internal/site/ -run
TestSpikeByteIdentity -v`);
[`internal/site/forest_test.go:71-97`](../../../internal/site/forest_test.go)
and
[`internal/site/render_test.go:24-47`](../../../internal/site/render_test.go)
(the pins the spike output was compared against);
[`go.mod`](../../../go.mod) and repo root (no `Makefile` /
`justfile` / `.github/workflows/` — the Go-toolchain-only gate).

## Validation Gate

This is a decision + doc-reconciliation task; no render code
changes. Before the implementing PR opens:

1. **No code migration shipped.** `git diff` touches only
   `docs/**` and `design/v0.1-design.md`; `internal/` and `cmd/`
   carry no functional change. Falsifier: any `internal/site`
   diff hunk altering render behavior.
2. **Render falsifiers pass unchanged.** `go test ./...` is
   clean and `TestRenderNoFieldNodeUnchanged` /
   `TestRenderTwoRegionShell` pass with their literals
   unmodified. Falsifier: either test edited, or a failure.
3. **Toolchain clean.** `gofmt -l internal cmd` empty,
   `go build ./...`, `go vet ./...`, `go test ./...` all pass —
   confirming the spike's `templ` dependency did not leak into
   the planning branch (`go.mod` carries no `a-h/templ`).
4. **§10 records the outcome, not "deferred."**
   `design/v0.1-design.md` §10 states `html/template` is the
   standing choice with the re-trigger condition, and no longer
   frames the decision as open/backlogged. The project-structure
   diagram's `cmd/ tool/ main.go` is corrected to
   `cmd/workstream-tracker/` (the real binary path). Falsifier:
   §10 still says the decision is deferred/tracked in the
   backlog, or the diagram still says `tool/`. Verified by:
   [`design/v0.1-design.md:211-222`](../../../design/v0.1-design.md)
   (the diagram, currently `cmd/ tool/ main.go`);
   [`cmd/workstream-tracker/main.go`](../../../cmd/workstream-tracker/main.go)
   (the real path).
5. **Backlog graduated.** The `templ-render-adoption` entry
   carries `Status: Graduated — templ-render-adoption` and a
   `**Plan:**` line; the slug token is bare (exact-match), the
   `**Plan:**` line is a separate markdown link. Falsifier: the
   entry still `Open`, or the Status slug wrapped in link syntax.
6. **Manual browser render observation.** `go run
   ./cmd/workstream-tracker` and load the index; confirm the
   forest renders unchanged (no visible regression) — the
   decision ships no render change, so the rendered page is
   identical to pre-task `main`.

The implementing PR body carries a `## Review Stance` section
(plan-doc PR — the canonical stance) and an `## Estimate
Deviations` section per the Plan-to-PR Completion Gate.

## Self-Review Audits

Drawn from
[`docs/agents/local/self-review-catalog.md`](../../agents/local/self-review-catalog.md);
diff surface is documentation + a decision record:

- **validation-honesty** — the Validation Gate above is
  read/observation-based; audit that each numbered step was
  performed end-to-end against the final diff, not asserted from
  the plan.
- **readiness-gate-truthfulness** — this plan reaches a terminal
  Status; audit that the named gate conditions actually held
  before the `Landed` flip rather than being announced.
- **trigger-map-currency** — the change rewrites cross-doc
  pointers (backlog → plan, §10 → outcome); audit that every
  added/changed link resolves and no pointer drifts.

## Risk Register

- **Deferral makes a future `templ` migration progressively more
  expensive — bounded, linear, mitigable.** Accepting
  `html/template` now means a later migration (if the re-trigger
  fires) ports more partials than exist today. The structural
  cost is *not* a cliff: the one-time costs (add the `templ`
  CLI, wire `templ generate` into the gate, rewrite the
  byte-identity falsifiers to semantic assertions) are fixed
  regardless of timing, the render layer is cleanly separated
  (`PlanNode` is a plain struct `templ` consumes unchanged, each
  template renders through an isolated Go entrypoint), and
  migration is leaf-up and component-by-component. The only cost
  that genuinely *accretes* is this repo's own byte-identity
  test ratchet: each new partial added under `html/template`
  ships with a byte-exact pin (the
  `TestRenderNoFieldNodeUnchanged` discipline), and every such
  pin is one more test to rewrite at migration time because
  `templ` minifies whitespace.
  **Mitigation (keeps the deferral cheap):** new partials added
  while `html/template` is the standing choice assert on
  *semantic* output (presence/structure of the rendered markup),
  not byte-exact snapshots. The existing pins
  (`TestRenderNoFieldNodeUnchanged`, `TestRenderTwoRegionShell`)
  stay as-is — they are load-bearing region-split falsifiers
  today and rewriting them is only warranted if the re-trigger
  fires — but the ratchet stops tightening for net-new render
  surface. Verified by:
  [`internal/site/forest_test.go:71-97`](../../../internal/site/forest_test.go)
  (the byte-exact pin pattern this mitigation declines to
  propagate); spike branch `spike/templ-render` @ `d19445b`
  (the `templ` whitespace-minification behavior that makes
  byte-exact pins migration-hostile).

## Out of Scope

- **Migrating the renderer to `templ`.** The decision is to not
  migrate; A1/A2/A3 are rejected with rationale above.
- **Revising the byte-identity falsifier tests.** They pass
  unchanged because no migration ships; deliberately rewriting
  them is only on the table if a future re-trigger reopens this.
- **Adding a codegen step / build wrapper to the project.** No
  `Makefile`/`justfile`/CI is introduced; the Go-toolchain-only
  gate is unchanged.
- **Re-reconciling §10's narrative history.** PR #29 already
  recorded `html/template` + the original deferral rationale;
  this task updates §10 to the *outcome* (standing choice +
  re-trigger), it does not re-explain the original divergence.
- **Removing the dangling `spike/templ-render` branch.** Left
  for reference per the spike worktree rule; not merged, not
  promoted.

## Backlog Impact

Graduated from the
[`templ-render-adoption`](../../backlog.md#templ-render-adoption)
entry. That entry's Status is set to
`Graduated — templ-render-adoption` with the optional
`**Plan:**` line pointing at this doc, in the same change that
creates this plan (per the backlog graduation lifecycle — this
standalone decision task's only tracking surface; no parent
epic/milestone row). No other backlog entry graduates, deletes,
splits, or shifts. Verified by:
[`../../backlog.md`](../../backlog.md);
[`../../../spec/backlog.md:43-72`](../../../spec/backlog.md)
(graduation + `**Plan:**` line lifecycle).

## Related Docs

- The transient scoping deliberation (decision space + bounding
  constraints with `Verified by:` citations) deletes at this
  task's terminal PR per
  [`task-plan.md`](../../../spec/planning/task-plan.md)
  "Scoping owns / plan owns"; it survives in git history.
- [`design/v0.1-design.md`](../../../design/v0.1-design.md) §10
  — the design record this task reconciles to the outcome.
- Spike branch `spike/templ-render` @ `d19445b` — the throwaway
  evidence; not merged.
