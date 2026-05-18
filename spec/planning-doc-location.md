# Planning Doc Location

The in-repo layout convention for plan-tree docs. Cross-referenced
from [`planning/shared.md`](./planning/shared.md) (slug rule),
[`planning/epic.md`](./planning/epic.md) (epic path),
[`planning/milestone.md`](./planning/milestone.md) (milestone
path), and [`planning/task-plan.md`](./planning/task-plan.md) (task and
phase paths).

## Layout convention

Every plan-tree root lives in its own folder under
`docs/plans/`. The root doc is at `<root-slug>/README.md`. Each
milestone of an epic lives in its own per-milestone folder
`<root-slug>/m<N>/`; that milestone's doc is the `README.md`
inside it, and the milestone's task and phase docs are sibling
files in the same `m<N>/` folder. A standalone task plan (no
milestones) keeps its descendants directly in the root folder.

The folder name is the root slug; the `m<N>/` segment carries
the milestone position. Inside `m<N>/`, a descendant's filename
encodes the slug-suffix *after the milestone segment* (`t<T>`,
`t<T>-p<P>`) plus an optional human-readable descriptor — the
path's folder segment, not the filename, carries the `m<N>`
disambiguation.

```
docs/plans/
  <root-slug>/                            # root folder
    README.md                             # root doc (epic / standalone task plan)
    m<N>/                                 # per-milestone folder
      README.md                           # milestone doc
      t<T>-<descriptor>.md                # task plan
      t<T>-p<P>-<descriptor>.md           # phase plan
      scoping/                            # transient (per-milestone)
        <suffix-after-mN>-<descriptor>.md # scoping doc
  <standalone-task-slug>/                 # standalone task plan root
    README.md                             # the task plan (N = 1)
    p<P>-<descriptor>.md                  # phase plan (standalone, N ≥ 2)
    scoping/                              # transient
      <slug-suffix>-<descriptor>.md       # scoping doc
```

Concrete examples:

```
docs/plans/
  madrona-feedback/                       # epic root
    README.md                             # epic doc
    m1/                                   # milestone 1 folder
      README.md                           # milestone 1 doc
      t1-ingest.md                        # task 1 of milestone 1
      t1-p1-schema.md                     # phase 1 sub-plan
      t1-p2-backfill.md                   # phase 2 sub-plan
      scoping/
        t1-p1-schema.md                   # scoping for m1 phase 1

  docs-canonical-corrections/             # standalone task plan root
    README.md                             # the task plan (N = 1)
```

## Why this shape

- **Every root gets a folder, no top-level files in `docs/plans/`.**
  Removes the special case that would otherwise distinguish
  standalone task plans (single file) from epic-rooted trees
  (folder).
- **Per-milestone folders bound folder size.** A multi-milestone
  epic with phased tasks accumulates many descendant docs; nesting
  each milestone's docs under `m<N>/` keeps any one folder to a
  single milestone's worth of files instead of a flat list that
  grows without bound as the epic adds milestones. The cost is a
  one-time file relocation when a plan changes shape — a standalone
  task plan that grows a milestone, or a milestone reorganization —
  which is accepted: the slug is immutable identity (see
  [`planning/shared.md`](./planning/shared.md) "Plan-doc identity
  (slug)"), so moving the file changes layout only, never identity,
  and the rendered tree is unaffected.
- **`README.md` for the root doc and each milestone doc.** GitHub
  renders it automatically when a reader navigates into the folder.
  Survives doc-type changes (a folder that starts as a standalone
  task plan and grows into an epic doesn't need a rename; an
  `m<N>/` folder's `README.md` stays the milestone doc regardless
  of how many task/phase siblings it gains).
- **Filename = slug-suffix + optional descriptor.** The slug in
  frontmatter is the identity of record (see
  [`planning/shared.md`](./planning/shared.md) "Plan-doc identity
  (slug)"). Inside an `m<N>/` folder the filename encodes the
  slug-suffix *after* the milestone segment; the `m<N>` folder
  carries that part of the position. The filename is purely for
  human browsing and may carry an optional descriptor; the tool
  reads slug from frontmatter, never from filename or folder path.
- **`scoping/` subfolder is transient.** Its contents delete in
  batch at the milestone-terminal PR (or task-terminal PR for
  standalone task plans) per
  [`planning/task-plan.md`](./planning/task-plan.md) "Scoping owns / plan
  owns." Co-locating scoping under each milestone's `m<N>/scoping/`
  makes that batch delete a clean subtree removal rather than a
  slug-prefix-filtered selection from a shared folder.

## Slug-to-path mapping

The full slug-to-path mapping under this layout:

| Slug                                    | Path                                                              |
|-----------------------------------------|-------------------------------------------------------------------|
| `<root>` (epic or standalone task)      | `docs/plans/<root>/README.md`                                     |
| `<root>-m<N>`                           | `docs/plans/<root>/m<N>/README.md`                                |
| `<root>-m<N>-t<T>`                      | `docs/plans/<root>/m<N>/t<T>-<descriptor>.md`                     |
| `<root>-m<N>-t<T>-p<P>`                 | `docs/plans/<root>/m<N>/t<T>-p<P>-<descriptor>.md`                |
| `<root>-m<N>-...` paired scoping        | `docs/plans/<root>/m<N>/scoping/<suffix-after-mN>-<descriptor>.md`|
| `<root>-p<P>` (standalone task, N ≥ 2)  | `docs/plans/<root>/p<P>-<descriptor>.md`                          |
| `<root>-...` standalone paired scoping  | `docs/plans/<root>/scoping/<slug-suffix>-<descriptor>.md`         |

The `-<descriptor>` portion is optional in every row; pure
slug-suffix filenames (e.g., `t1.md`, `t1-p2.md`) are valid.
Standalone task plans have no milestone segment, so their
descendants and `scoping/` stay directly under `<root>/` exactly
as before — the `m<N>/` nesting applies only to epic-rooted
trees.

## Per-epic milestone numbering

Each epic counts milestones from `m1` independently. Sibling
epics may reuse the same milestone numbers without collision
because the slug's root segment (`<epic-slug>`) disambiguates;
the path's `<root-slug>/m<N>/` folder segments carry the same
disambiguation.

## Migration of pre-existing flat trees

A consumer project that adopted an earlier flat layout (every
descendant a sibling file directly under `<root>/`) is **not
required** to relocate existing trees. The visualization tool
discovers docs by their frontmatter slug regardless of folder
depth (see "Tool behavior" below), so a flat tree and a nested
tree both render correctly and can coexist. The convention going
forward is: new epics use the per-milestone nesting, and an
existing flat epic restructures organically only when its next
milestone is drafted or otherwise touched. There is no mandated
backfill. (This repository dogfoods the spec and migrates its own
trees on its own schedule; that is a repo-local choice, not a
rule the spec imposes on consumers.)

## Tool behavior

The visualization tool walks each `docs/plans/<root-slug>/`
folder **recursively** and identifies every plan-tree doc by the
`slug` in its YAML frontmatter; the file's path is used only for
diagnostics, never to derive identity or tree position (tree
hierarchy is computed from the slug — see
[`planning/shared.md`](./planning/shared.md) "Plan-doc identity
(slug)"). The walk skips any directory named `scoping` at any
depth, so a per-milestone `m<N>/scoping/` is transient and
non-rendered just like a root-level `scoping/`. Because identity
and hierarchy are slug-derived, the tool is agnostic to whether
a tree is flat or per-milestone nested; the layout above is the
authoring convention, not a parsing constraint the tool
enforces. Project-customizable layout (a `workstream.toml` or
similar) is deferred.
