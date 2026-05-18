# Scoping — Nested milestone doc layout

Transient deliberation for the `nested-milestone-doc-layout`
standalone task plan. Deletes in batch at this task's terminal
PR per [`task-plan.md`](../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns"; survives in git history. Carries no
`Status` field (scoping docs don't).

## Context

A flat epic-rooted plan tree puts every milestone, task, and
phase doc as sibling files in one folder. A multi-milestone epic
accumulates dozens of flat siblings — `workstream-tracker-1-0/`
is already at ~13 files and still growing. This task changes the
in-repo layout convention so milestone docs nest into a
per-milestone folder, bounding folder size to one milestone's
worth of files. The surface is the layout-convention spec prose
this repo dogfoods and consumers vendor, plus a dogfood walker
test; the reality-check below establishes that no tree/walker
behavior changes.

## Decisions made at scoping time

### D1 — Milestone-level nesting only

A milestone's doc lives in a per-milestone folder
`docs/plans/<root>/m<N>/`; that milestone's task and phase docs
are flat files inside `m<N>/`. Rejected alternatives: **(a)
status-quo flat** — a multi-milestone epic accumulates dozens of
flat siblings, which is the problem being solved. **(b) deeper
per-task folders `m<N>/t<T>/`** — a task usually has few phases,
so per-task folders add depth and more move-on-growth churn for
little browsing gain; milestone-level nesting captures almost all
the win. Verified by:
[`../../workstream-tracker-1-0`](../../workstream-tracker-1-0)
(13 flat files in one epic folder today);
[`planning-doc-location.md:19-45`](../../../../spec/planning-doc-location.md)
(current flat convention).

### D2 — Milestone doc filename is `README.md` inside `m<N>/`

The milestone doc is `docs/plans/<root>/m<N>/README.md`. Rejected
alternative: keep `m<N>.md` inside `m<N>/` — the root-folder
rationale already in the spec (GitHub renders `README.md` on
folder navigation; it survives doc-type changes) applies
identically to the milestone folder, so reusing `README.md` is
consistent, not a new rule. Verified by:
[`planning-doc-location.md:55-58`](../../../../spec/planning-doc-location.md)
(the `README.md`-for-the-root rationale this extends).

### D3 — Scoping co-locates per-milestone at `m<N>/scoping/`

Scoping docs for a milestone's tasks and phases live at
`docs/plans/<root>/m<N>/scoping/`; a standalone task plan's
scoping stays `docs/plans/<root>/scoping/`. Rejected
alternative: one root-level `scoping/` for the whole epic —
scoping deletes in batch at the milestone-terminal PR, so
co-locating under `m<N>/` makes that batch a clean subtree
removal instead of a slug-prefix-filtered subset of a shared
folder. Verified by:
[`planning-doc-location.md:67-70`](../../../../spec/planning-doc-location.md),
[`task-plan.md:843-848`](../../../../spec/planning/task-plan.md)
(scoping transient, deletes at the milestone-terminal PR).

### D4 — Standalone task plans unchanged

A standalone task plan has no milestones, so it stays
`docs/plans/<root>/README.md` with `<root>/scoping/` exactly as
today. The blast radius is epic-rooted trees only. Verified by:
[`task-plan.md:830-836`](../../../../spec/planning/task-plan.md)
(standalone task-plan path conventions).

### D5 — No walker or tree-builder code change

`walkPlans` recurses with `filepath.WalkDir` over each root
folder and derives identity from the frontmatter `slug`; the
tree hierarchy is built by parsing slugs, never paths; the
`scoping` skip matches the directory name at any depth. Nesting
milestone docs changes paths only, never slugs, so the rendered
tree is byte-identical and `m<N>/scoping/` is still skipped. The
change ships a dogfood walker test that exercises the nested
shape and asserts the rendered doc set plus the scoping-skip
match the flat equivalent — proving the tolerance by execution,
not asserting it in prose. Verified by:
[`walker.go:73`](../../../../internal/site/walker.go)
(recursive `WalkDir`),
[`walker.go:77-80`](../../../../internal/site/walker.go)
(`scoping` skipped by directory name at any depth),
[`walker.go:117`](../../../../internal/site/walker.go)
(slug from frontmatter),
[`walker.go:124-131`](../../../../internal/site/walker.go)
(`Path` is diagnostics-only);
[`tree.go:117`](../../../../internal/site/tree.go) and
[`tree.go:144-151`](../../../../internal/site/tree.go)
(hierarchy from slug parse + `parsed.Parent()`, no path read);
[`walker_test.go:76-91`](../../../../internal/site/walker_test.go)
(existing scoping-skip test, currently flat).

### D6 — Consumer-project migration policy: hybrid

Resolved 2026-05-18 (was the open decision below). Because the
walker tolerates both shapes (D5), the spec **mandates nothing
for consumer projects**: new epics nest, and an existing flat
epic restructures organically only when its next milestone is
drafted or otherwise touched — no required backfill, mirroring
the "no historical-doc backfill" precedent. Independently, **this
repo migrates its own two dogfood trees**
(`workstream-tracker-1-0/`, `demo-workstream/`) in this task's
implementing PR, so the repo does not ship a convention its own
trees violate. Rejected: **big-bang** (every adopting project
migrates all trees in the adopting PR) — highest one-time
consumer burden for consistency the tool does not need, since
both shapes render. Rejected: **forward-and-organic with no
self-migration** — leaves this repo shipping a convention its own
largest tree violates, which is poor dogfooding. Verified by:
[`../../stub-children-on-parent-promotion/README.md:371-373`](../../stub-children-on-parent-promotion/README.md)
(the "no historical-doc backfill" precedent the consumer policy
mirrors).

### D7 — Correct the stale "v0.0 hardcoded [flat]" paragraph

[`planning-doc-location.md:94-99`](../../../../spec/planning-doc-location.md)
says the tool "walks `docs/plans/<root-slug>/` and expects the
structure above" as hardcoded. That is factually wrong today:
the walker is recursive and slug-driven (D5), not flat-hardcoded.
The same edit that revises the convention corrects this paragraph
to state the walker is path-shape-agnostic and slug-identity
driven. This rides the same surface; it is not a separate task.
Verified by:
[`walker.go:41-100`](../../../../internal/site/walker.go)
(recursive, slug-driven) contradicting
[`planning-doc-location.md:94-99`](../../../../spec/planning-doc-location.md).

## Open decisions to make at plan-drafting

None outstanding. D6 (consumer-project migration policy) was the
only open decision; it resolved 2026-05-18 to the hybrid and is
recorded above under "Decisions made at scoping time."

## Plan structure handoff

Single standalone N=1 task plan; one implementing PR. Estimated
work: layout-convention spec prose (`planning-doc-location.md`,
`task-plan.md` "Path conventions"), one nested-shape dogfood
walker test, and — if D6 resolves to (c) — migrating this repo's
two epic trees. Suggested commit boundaries: (1) spec prose +
dogfood test; (2) this repo's own tree migration. The plan-doc's
estimate-shaped sections are labeled per
[`shared.md`](../../../../spec/planning/shared.md) "Plan content
is a mix of rules and estimates."

## Reality-check inputs the plan must verify

- **Walker is recursive and slug-identity driven** (D5
  citations). Re-confirm `walker.go` / `tree.go` line anchors at
  promotion; they drift as files are edited.
- **No other repo path-shape coupling.** Tree structure is
  slug-derived; the only path coupling is "each root is a
  top-level folder" ([`walker.go:60-66`](../../../../internal/site/walker.go)),
  which milestone nesting *inside* a root folder does not touch.
  Registration is slug/API-based; `main.go:51` only sets
  `plansPath`. Verified by reading `walker.go`, `tree.go`,
  [`main.go:51`](../../../../cmd/workstream-tracker/main.go).
- **`assemble.sh` substitutes link roots, not structure** — the
  spec is vendored via link-root substitution, unaffected by a
  layout-prose change. Confirm at promotion that no assembly
  step encodes the flat path shape.
