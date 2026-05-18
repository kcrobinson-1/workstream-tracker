# Scoping — templ-render-adoption

Transient deliberation for the standalone decision task
`templ-render-adoption`. Deletes in batch at this task's
terminal PR per [`task-plan.md`](../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns"; survives in git history.

This is a **doc-only decision task**: the deliverable is the
recorded decision (migrate to `templ` vs. formally accept
`html/template`) plus rejected alternatives, driven by a
throwaway spike. Per
[`task-plan.md`](../../../../spec/planning/task-plan.md)
"Doc-only decision plans satisfy the substantive-content gate
via cited open-question constraints, not resolved decisions,"
this scoping doc surfaces the decision space with code-grounded
`Verified by:` citations on the constraints that bound the open
question; it does not resolve the decision at scoping-open time.
Resolution happens through the deliberation that constitutes the
plan and absorbs back into the durable plan doc.

## Context

[`design/v0.1-design.md`](../../../../design/v0.1-design.md) §10
originally locked `templ` for HTML rendering; v0.1 implementation
deliberately used the stdlib `html/template` instead, deferring
templ "until there are real reusable components" (commit
`60040be`). m2-t1 split the renderer into a shell composing
`forest` / `roster` / `node` sub-templates — the first plausible
"real reusable components" trigger. §10 was reconciled (PR #29)
to record `html/template` and point at the
`templ-render-adoption` backlog entry; this task resolves the
deferred decision and updates §10 to the **outcome**.

## The decision space (decomposed into shapes)

Per [`shared.md`](../../../../spec/planning/shared.md)
"Decompose options into shapes before analyzing." The backlog
names two top-level options; each hides sub-shapes:

- **A. Migrate to `templ`.**
  - **A1 — full migration.** All of `forest` / `node` / `roster`
    + the shell become `templ` components; `html/template`
    removed.
  - **A2 — hybrid.** `templ` for component bodies, a string
    shell (or vice versa). Two render systems coexist.
  - **A3 — migrate, preserve exact bytes.** Migrate but reproduce
    the existing `html/template` byte output so the falsifier
    pins survive unchanged.
- **B. Formally accept `html/template`.**
  - **B1 — accept silently.** Update §10, no re-trigger named.
  - **B2 — accept with a documented re-trigger condition.**
    Record `html/template` as the standing choice and state the
    concrete condition under which `templ` is revisited.

## Constraints that bound the decision (code-grounded)

- **Byte-identity falsifier pins.** `TestRenderNoFieldNodeUnchanged`
  pins exact node markup including a text-template range artifact
  (`</li>\n  \n</ul>` — two spaces then newline, a meaningless
  whitespace artifact of `html/template`'s `{{range}}`).
  `TestRenderTwoRegionShell` pins the shell composition. These
  must pass unchanged unless deliberately revised with rationale.
  Verified by:
  [`internal/site/forest_test.go:71-97`](../../../../internal/site/forest_test.go)
  (the `wantNode` literal);
  [`internal/site/render_test.go:24-47`](../../../../internal/site/render_test.go)
  (two-region shell pin).
- **`templ` minifies insignificant whitespace by design.** The
  spike (branch `spike/templ-render`, commit `d19445b`, not
  merged) ported `forest`/`node` faithfully and rendered the
  exact `TestRenderNoFieldNodeUnchanged` input: `templ` emitted
  `<div class="root"><span ...>Proposed</span>...<ul><li>...</li></ul></div>`
  (234 B, single line, whitespace collapsed); the
  `html/template`-pinned bytes are 218 B with newlines, `  `
  indents, and the `</li>\n  \n</ul>` artifact — a **byte-identity
  miss**. Reproducing a text-template range whitespace artifact
  in `templ` is not achievable (its whitespace normalization is
  intentional and not configurable to that end), so A3 is not a
  viable shape and A1/A2 force a deliberate falsifier rewrite.
  Verified by: spike branch `spike/templ-render` @ `d19445b`
  (`go test ./internal/site/ -run TestSpikeByteIdentity -v`
  output recorded in that commit).
- **Validation gate is the Go toolchain only — no codegen step.**
  No `Makefile`, no `justfile`, no `.github/workflows/`. The gate
  is `gofmt -l internal cmd` / `go build ./...` / `go vet ./...`
  / `go test ./...` plus a manual browser render observation.
  `templ` requires the `templ` CLI (`go install
  github.com/a-h/templ/cmd/templ`) and a `templ generate` step
  emitting committed `*_templ.go` — a new build-tool dependency
  and a generated-vs-source drift surface the gate would have to
  grow. Verified by: repo root (`ls Makefile justfile` empty,
  no `.github/workflows/`);
  [`go.mod`](../../../../go.mod) (no `a-h/templ`).
- **Render surface size.** The component surface is `forest`,
  recursive `node`, `roster` placeholder, and the shell —
  ~3 thin templates, the `roster` still a placeholder until m2
  t4. Verified by:
  [`internal/site/forest.go`](../../../../internal/site/forest.go),
  [`internal/site/roster.go`](../../../../internal/site/roster.go),
  [`internal/site/render.go`](../../../../internal/site/render.go).
- **Design-doc minimal-dependency posture.**
  [`design/v0.1-design.md`](../../../../design/v0.1-design.md)
  §10 frames Go for "a single static binary" and treats each
  community dependency as a deliberated cost; the original templ
  deferral rationale ("until there are real reusable components")
  is the standing bar. Verified by:
  [`design/v0.1-design.md:195-226`](../../../../design/v0.1-design.md).

## Open decision to resolve at plan-drafting

- **D1 — migrate (`templ`) vs. formally accept (`html/template`),
  and which sub-shape.** Bounded by the constraints above; the
  spike is the dealbreaker evidence. Resolved in the plan doc's
  Decision section, surfaced for maintainer review via the
  implementing PR (no parent epic/milestone — the graduated
  backlog entry is the only tracking surface).

## Reality-check inputs the plan must verify before promotion

- The two falsifier tests still pin the bytes cited above
  (`forest_test.go` `TestRenderNoFieldNodeUnchanged`,
  `render_test.go` `TestRenderTwoRegionShell`) at promotion time.
- §10 current text and the `cmd/ tool/ main.go` diagram error
  (real path `cmd/workstream-tracker/main.go`) still present.
- Backlog entry `templ-render-adoption` graduated (Status +
  `**Plan:**` line) in this same change.

## Plan-structure handoff

Decision task → the durable plan doc carries: Context preamble,
Goal, Decision (the resolved D1 + the documented re-trigger
condition), Options Considered (rejected shapes with rationale),
Spike Evidence, Validation Gate, Out Of Scope, Backlog Impact,
Related Docs. It legitimately **skips** `Contracts` and `Files
to touch — new` in their code-contract sense (the deliverable is
the decision; the only file edits are doc reconciliation) — a
decision-plan variance disclosed in the PR body's
`## Documentation` per
[`shared.md`](../../../../spec/planning/shared.md) "Section
variance disclosure."
