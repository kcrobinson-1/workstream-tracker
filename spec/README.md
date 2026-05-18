# spec/

The canonical plan-doc spec consumer projects adopt to participate
in workstream-tracker visualization. This directory is the
authority for plan-tree structure, frontmatter conventions, the
Status lifecycle, the doc-type taxonomy, the per-root-folder
layout, and the backlog format. Other repositories consume this
spec — most lightly via vendored copies; eventually via a
gitignored symlink + bootstrap script (deferred to v0.1).

## Files

- **[planning/shared.md](planning/shared.md)** — cross-level
  rules that bind every plan-drafting session regardless of
  doc-type. The foundational vocabulary lives here:
  - **Taxonomy**: node-type vs. doc-type axes, the four-segment
    cap, the N=1 collapse.
  - **Plan-doc identity (slug)**: kebab-case roots, hierarchical
    descendant slugs (`<root>-mN-tN-pN`), encoded-at-creation
    not current-position.
  - **Plan-doc Status**: canonical lifecycle (`In draft` →
    `Proposed` → `In progress` → `Validating` → `Landed`, plus
    `Deferred — <reason>`), two-channel state model (Status vs.
    work-instance state), graceful fallback for unknown values.
  - Cross-cutting authoring rules: contract altitude, "Verified
    by:" annotations, falsifiability checks, exact-match label
    discipline, options-into-shapes decomposition, and others.
- **[planning/epic.md](planning/epic.md)** — epic-doc rules:
  scope (the *what* and *why* of a multi-milestone arc),
  required and optional sections, path conventions.
- **[planning/milestone.md](planning/milestone.md)** —
  milestone-planning rules: cross-phase coordination, sequencing
  (Mermaid `flowchart LR`), cross-phase decisions, the
  anti-goal "do not scope any phase in this session."
- **[planning/task-plan.md](planning/task-plan.md)** — the
  implementation-layer authority. Binds both task plans and
  phase plans (the merged-doc-type approach); covers the
  Planning Depth rule, the `In draft` → `Proposed` promotion
  gate, the Plan-to-PR Completion Gate, scoping-doc semantics,
  reality-check pass, narrow-surface plan carve-out, and more.
- **[planning-doc-location.md](planning-doc-location.md)** —
  the in-repo layout convention: every plan-tree root in its
  own folder under `docs/plans/<root-slug>/`, the root doc at
  `<root-slug>/README.md`, each epic milestone nested in its own
  `m<N>/` folder, scoping docs in a transient `scoping/`
  subfolder.
- **[backlog.md](backlog.md)** — the backlog format: flat list
  of pre-plan work items at `docs/backlog.md`, stable
  kebab-case slugs in a separate namespace from plan-tree
  slugs, lifecycle (`Open` → `Graduated — <plan-slug>` →
  deletion).

## Where to start

If you're authoring a plan doc, jump to the per-doc-type file
that matches what you're writing
([epic.md](planning/epic.md), [milestone.md](planning/milestone.md),
or [task-plan.md](planning/task-plan.md)) — each loads
[shared.md](planning/shared.md) for the cross-level rules.

If you're new to the spec, read in this order:
1. [planning/shared.md](planning/shared.md) Taxonomy + Plan-doc
   identity + Plan-doc Status sections (the foundational
   vocabulary).
2. [planning-doc-location.md](planning-doc-location.md) (where
   docs live on disk).
3. The per-doc-type file matching your authoring task.

## TODO

The TODO note at the top of [planning/shared.md](planning/shared.md)
flags one known portability gap: several rules embed a GitHub-style
PR review workflow (`## Review Stance` headings in PR bodies, the
Estimate Deviations callout, etc.). The protective intent of those
rules is universal but the scaffolding assumes that workflow shape.
Revisit when the spec's first non-GitHub consumer surfaces.
