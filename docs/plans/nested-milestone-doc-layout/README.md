---
slug: nested-milestone-doc-layout
Status: Landed
short_description: Nest milestone docs under per-milestone m<N>/ folders
---

# Nested milestone doc layout

## Context

A flat epic-rooted plan tree puts every milestone, task, and
phase doc as sibling files in a single folder. Once an epic has
several milestones with phased tasks, that folder becomes dozens
of files a human has to scan flat — `workstream-tracker-1-0/`
is already ~13 files and still growing. This task changes the
in-repo layout convention so each milestone's doc and its
descendants nest into a per-milestone `m<N>/` folder, keeping any
one folder to a single milestone's worth of files.

It's being done now because the flat shape is actively getting
unwieldy in this repo's own dogfooded trees and will only get
worse as `workstream-tracker-1-0` adds milestones. The surfaces
touched are conceptual: the layout-convention spec prose this
repo dogfoods and consumer projects vendor, a dogfood test that
locks the visualization tool's tolerance of the new shape, and
this repo's own plan-tree folders.

A reality-check (recorded in the scoping doc) established the
load-bearing fact that bounds scope: the visualization walker is
already recursive and derives node identity from the frontmatter
slug, never from the path, so nesting milestone docs is a
path-only change with no tree-rendering behavior change.

## Goal

The in-repo layout convention places each milestone's doc, its
task and phase docs, and their scoping docs inside a per-milestone
folder; standalone task plans are untouched; the slug remains the
immutable identity and the path is pure layout; and the
visualization tool's path-shape independence is stated correctly
in the spec and locked by an executable test rather than only
asserted in prose.

Verifiable when: the layout convention and the task-plan "Path
conventions" both describe the per-milestone nesting and their
slug-to-path mapping reflects it; the stale "hardcoded flat"
paragraph is corrected to match the recursive slug-driven walker;
a walker test exercises the nested shape and asserts the rendered
doc set plus scoping-skip match the flat equivalent; and the
consumer-migration policy is stated explicitly.

## Phases

This task drafted N = 1 (one PR: spec prose + dogfood test +
this repo's own tree migration). During implementation the repo
owner directed that this repo's own existing plan-tree docs are
**not** to be moved yet. That splits the work along a phase
boundary per
[`task-plan.md`](../../../spec/planning/task-plan.md) "N = 1 →
N ≥ 2 transition": Phase 2 fits cleanly under this doc's existing
Contracts and Validation apparatus (it is a mechanical,
slug-preserving relocation already covered by C5 and Validation
step 4), so no separate phase plan file is created — the work is
tracked here.

- **Phase 1 — convention + tool tolerance (landed, PR #28).**
  The layout-convention spec prose (C1–C4), the
  consumer-migration policy statement (C5, consumer half), and
  the dogfood walker test. No existing plan docs moved.
- **Phase 2 — this repo's own tree migration (this PR;
  task-terminal).** Relocated all 22 descendant docs of
  `workstream-tracker-1-0/` and `demo-workstream/` into `m<N>/`
  folders via `git mv` (slug-preserving — every frontmatter slug
  verified byte-identical to pre-move), satisfying the
  repo-self-migration half of C5 and Validation step 4. All
  ~290 relative links into/out of the moved files were
  re-relativized so every link that resolved pre-move resolves
  post-move. This is the task-terminal PR: Status flips to
  `Landed` and the scoping doc is deleted here.

## Contracts

The final shape the change must establish — WHAT, not the prose
that establishes it (exact wording is the implementing PR's
choice).

### C1 — Per-milestone nesting in the layout convention

The layout convention in `planning-doc-location.md` and the
"Path conventions" in `task-plan.md` must establish, for an
epic-rooted tree: the milestone doc is at
`docs/plans/<root>/m<N>/README.md`; that milestone's task and
phase docs are flat files inside `docs/plans/<root>/m<N>/`;
scoping for those docs is at `docs/plans/<root>/m<N>/scoping/`.
The slug-to-path mapping table reflects these nested paths.
(scoping D1–D3)

### C2 — Standalone task plans unchanged

The convention continues to place a standalone (no-milestone)
task plan at `docs/plans/<root>/README.md` with
`docs/plans/<root>/scoping/`, with no nesting. (scoping D4)

### C3 — Slug is identity; nesting is a path-only change

The convention reaffirms that the slug is the immutable identity
and the path is layout: nesting changes paths, never slugs, and a
later plan-type change that relocates a file does not change its
slug. No new rule may derive identity or hierarchy from the path
shape. Verified by:
[`../../../internal/site/tree.go:117`](../../../internal/site/tree.go)
and
[`../../../internal/site/tree.go:144-151`](../../../internal/site/tree.go)
(hierarchy built from `slugs.Parse` + `parsed.Parent()`, never
from the path);
[`../../../spec/planning/shared.md:181-186`](../../../spec/planning/shared.md)
("Slug is identity; path is layout").

### C4 — Walker path-shape independence: stated and tested

The spec's tool-behavior text states the walker discovers docs by
recursive walk plus frontmatter slug and is agnostic to
flat-versus-nested layout, and the stale paragraph that claims a
hardcoded flat structure is corrected accordingly. A walker test
exercises the nested shape (milestone `README.md` in `m<N>/`, a
task doc in `m<N>/`, a scoping doc in `m<N>/scoping/`) and
asserts the rendered doc set and the scoping-skip are identical
to the flat equivalent. (scoping D5, D7) Verified by:
[`../../../internal/site/walker.go:73`](../../../internal/site/walker.go)
(recursive `filepath.WalkDir`),
[`../../../internal/site/walker.go:77-80`](../../../internal/site/walker.go)
(`scoping` skipped by directory name at any depth),
[`../../../internal/site/walker.go:117`](../../../internal/site/walker.go)
(identity from frontmatter `slug`);
[`../../../internal/site/walker_test.go:76-91`](../../../internal/site/walker_test.go)
(existing flat scoping-skip test the nested test parallels);
[`../../../spec/planning-doc-location.md:94-99`](../../../spec/planning-doc-location.md)
(the stale hardcoded-flat paragraph C4 corrects).

### C5 — Consumer-migration policy stated explicitly (hybrid)

The convention states explicitly that an adopting project is
**not** required to migrate pre-existing flat trees: new epics
nest, and an existing flat epic restructures organically only
when its next milestone is drafted or otherwise touched (no
mandated backfill). Independently of the consumer rule, this
repo's own `workstream-tracker-1-0/` and `demo-workstream/` trees
are migrated to the nested shape in this task's implementing PR.
(scoping D6, resolved 2026-05-18)

## Cross-Cutting Invariants

- **Slug ≠ path.** Every doc the convention names is reachable by
  its frontmatter slug regardless of nesting depth; no spec rule
  or test may reintroduce path-derived identity or hierarchy.
- **Both shapes render until migrated.** Until a project migrates,
  flat and nested epic trees both render correctly; any spec
  sentence or test that implies flat-only is a defect.
- **`scoping/` skipped at any depth.** A `scoping/` folder is
  non-rendered wherever it appears — root or `m<N>/`; the
  per-milestone scoping folder must stay non-rendered.
- **Exact-match tokens verbatim.** `In draft`, `Proposed`, and
  every cited section title are copied from their canonical
  source, not paraphrased.

## Files to touch

> Estimate of the expected change shape, not a binding contract.
> Deviations are handled per the Estimate Deviations callout in
> the implementing PR body.

**Modify (Phase 1 — what shipped, reconciled to actual):**

- [`../../../spec/planning-doc-location.md`](../../../spec/planning-doc-location.md)
  — the layout convention, its layout tree + concrete examples,
  the "Why this shape" rationale, the slug-to-path mapping table,
  the per-epic-numbering note, and the stale "v0.0 hardcoded
  [flat]" paragraph (rewritten as a "Tool behavior" section);
  added a "Migration of pre-existing flat trees" section
  (C1, C3, C4, C5).
- [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md)
  — the "Path conventions" epic-with-milestones and scoping
  bullets (C1, C2).
- [`../../../spec/planning/epic.md`](../../../spec/planning/epic.md),
  [`../../../spec/planning/milestone.md`](../../../spec/planning/milestone.md),
  [`../../../spec/README.md`](../../../spec/README.md) — the
  per-level path-convention pointers and the spec index summary,
  which also stated the old flat shape and would otherwise drift
  from the convention. **Estimate deviation:** the original
  estimate named only `planning-doc-location.md` +
  `task-plan.md`; the layout shape is restated in three more
  rule-bearing places that had to move in lockstep (C1/C3
  coherence). Called out in the PR body.
- [`../../../internal/site/walker_test.go`](../../../internal/site/walker_test.go)
  — added `TestWalkPlansNestedMilestoneLayout` asserting the
  nested doc set and the `m<N>/scoping/` skip (C4).

**Phase 2 — this repo's tree migration (this PR; what shipped):**

- `docs/plans/workstream-tracker-1-0/**` and
  `docs/plans/demo-workstream/**` — 22 descendant docs `git mv`d
  into `m<N>/` folders (milestone docs → `m<N>/README.md`;
  task/phase docs → `m<N>/<suffix-after-mN>.md`). Frontmatter
  slugs byte-identical to pre-move (verified).
- Relative links across the repo re-relativized so every link
  that resolved pre-move still resolves post-move (inbound links
  from the epic `README.md`s, standalone plans, and backlog;
  outbound links from moved files to spec/sibling/cross-milestone
  targets).

**Not touched:**

- [`../../../internal/site/walker.go`](../../../internal/site/walker.go),
  [`../../../internal/site/tree.go`](../../../internal/site/tree.go)
  — no behavior change; the walker is already recursive and
  slug-driven (scoping D5). Only a test is added.
- [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md)
  — "Slug is identity; path is layout" already says what C3
  needs; no edit.
- `cmd/**`, `internal/api/**`, schema — registration and API are
  slug-based, not path-shape coupled.
- `scripts/assemble.sh` — substitutes link roots, not structure.

## Validation Gate

Prose-plus-test change; validation is read- and test-based.
Before the implementing PR opens:

1. **Convention-coherence read.** `planning-doc-location.md`'s
   layout tree and slug-to-path table and `task-plan.md`'s "Path
   conventions" all agree on per-milestone nesting with no
   residual flat-only mapping. Falsifier: a table row or bullet
   still maps a milestone to `docs/plans/<root>/m<N>.md`.
2. **Stale-paragraph check.** The "v0.0 hardcoded" paragraph no
   longer claims the tool expects a flat structure. Falsifier:
   text still asserts a hardcoded flat layout.
3. **Walker dogfood test.** `go test ./internal/site/` passes and
   the new nested test asserts the nested doc set equals the flat
   equivalent and `m<N>/scoping/` is skipped. Falsifier: the test
   is absent, or asserts tolerance in a comment rather than by
   execution.
4. **Repo-tree conformance (Phase 2 — satisfied this PR).**
   Every file under `workstream-tracker-1-0/` and
   `demo-workstream/` sits at its convention path; `go test
   ./...` green; rendered roots-and-children set unchanged
   pre/post move (slugs byte-identical → the slug-driven tree is
   identical by construction, and `TestWalkPlansNestedMilestone
   Layout` proves walker tolerance of the nested shape).
   Falsifier (a moved file's slug changed, or the tree diff is
   non-empty): checked — 22/22 slugs byte-identical to
   origin/main, tests green.
5. **Link-resolution check (satisfied this PR).** Every relative
   link in the repo was resolved against the post-move tree.
   Result: zero migration-caused broken links. The only
   unresolved links are pre-existing (present on origin/main,
   unrelated to this migration): five `docs/agents/**` links and
   one `stub-children-on-parent-promotion` link in the old m2
   milestone doc that was broken before the move and was
   faithfully re-relativized to the same intended (still-missing)
   target — neither introduced nor fixed here (out of scope; see
   Out of Scope).
6. **Exact-match token check.** Every cited section title and
   Status token is verbatim against its canonical definition.

The implementing PR body carries a `## Review Stance` section
(spec-doc PRs default to the canonical stance) and an
`## Estimate Deviations` section per the Plan-to-PR Completion
Gate.

## Out of Scope

- **Deeper per-task `m<N>/t<T>/` folders.** scoping D1 alt-b
  rejected; revisit only if per-task file counts become the pain.
- **Mandating consumer big-bang migration.** Per scoping D6
  (resolved hybrid) the spec mandates nothing for consumers; this
  repo migrates only its own trees.
- **Walker/tree behavior changes.** scoping D5 — code is
  unchanged; only a dogfood test is added.
- **Project-customizable layout (`workstream.toml`).** Already
  deferred by the layout convention; untouched.
- **Backlog and standalone-task scoping conventions.** Unchanged.
- **Pre-existing broken links.** Six links broken on
  origin/main before this migration (five under `docs/agents/**`,
  one `stub-children-on-parent-promotion` reference in the old
  m2 milestone doc) are not fixed here. They are unrelated to
  per-milestone nesting; the migration only re-relativizes them
  so a moved file's pre-existing broken link points at the same
  intended (still-missing) target rather than a newly-wrong one.
  Fixing them is tracked separately, not bundled into this
  layout change.

## Backlog Impact

A new entry `nested-milestone-doc-layout` is captured in
[`../../backlog.md`](../../backlog.md) and graduated
(`Graduated — nested-milestone-doc-layout` plus the optional
`**Plan:**` line pointing at this doc) in the planning change
that creates this plan — this standalone task's only tracking
surface (no parent epic or milestone row). No existing backlog
entry graduates, deletes, splits, or shifts. Verified by:
[`../../../spec/backlog.md:43-72`](../../../spec/backlog.md)
(graduation + `**Plan:**` line lifecycle).

## Self-Review Audits

Drawn from
[`../../agents/local/self-review-catalog.md`](../../agents/local/self-review-catalog.md);
the diff surface is documentation/spec plus a mass file move:

- **rename-aware-diff-classification** — the tree migration is
  mass file moves tooling sees as add+delete; hand-classify that
  no slug or rendered behavior changed before any "no behavior
  changed" claim.
- **trigger-map-currency** — moved/renamed directories; re-check
  any declared pre-edit trigger map or agent layout assumption
  against the new nested shape so deterministic triggers don't
  drift.
- **validation-honesty** — the gate is read- and test-based;
  audit each step ran end-to-end against the final diff, not
  asserted from the plan.

## Related Docs

- `scoping/README.md` — the transient scoping deliberation
  (D1–D7, rejected alternatives, the resolved D6); deleted in
  this task-terminal PR per "Scoping owns / plan owns", survives
  in git history.
- [`../../../spec/planning-doc-location.md`](../../../spec/planning-doc-location.md),
  [`../../../spec/planning/task-plan.md`](../../../spec/planning/task-plan.md),
  [`../../../spec/planning/shared.md`](../../../spec/planning/shared.md)
  — the convention surface this task contracts and the
  slug-is-identity rule it relies on unchanged.
- [`../../../internal/site/walker.go`](../../../internal/site/walker.go),
  [`../../../internal/site/tree.go`](../../../internal/site/tree.go)
  — the recursive, slug-driven implementation that makes nesting
  a path-only change.
- [`../stub-children-on-parent-promotion/README.md`](../stub-children-on-parent-promotion/README.md)
  — the spec-prose-only change precedent and the
  "no historical-doc backfill" pattern scoping D6 option (a)
  mirrors.
