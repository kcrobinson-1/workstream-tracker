# Planning Doc Location

The in-repo layout convention for plan-tree docs. Cross-referenced
from [`planning/shared.md`](./planning/shared.md) (slug rule),
[`planning/epic.md`](./planning/epic.md) (epic path),
[`planning/milestone.md`](./planning/milestone.md) (milestone
path), and [`planning/task-plan.md`](./planning/task-plan.md) (task and
phase paths).

## Layout convention

Every plan-tree root lives in its own folder under
`docs/plans/`. The root doc is at `<root-slug>/README.md`;
descendants are sibling files inside the same folder. The
folder name is the root slug; the filename of each descendant
encodes the slug-suffix (the part of the slug after the root)
plus an optional human-readable descriptor.

```
docs/plans/
  <root-slug>/                            # root folder
    README.md                             # root doc (epic / task plan)
    m<N>-<descriptor>.md                  # milestone doc
    m<N>-t<T>-<descriptor>.md             # task plan
    m<N>-t<T>-p<P>-<descriptor>.md        # phase plan
    scoping/                              # transient
      <slug-suffix>-<descriptor>.md       # scoping doc
```

Concrete examples:

```
docs/plans/
  madrona-feedback/                       # epic root
    README.md                             # epic doc
    m1.md                                 # milestone 1 doc
    m1-t1.md                              # task 1 of milestone 1
    m1-t1-p1.md                           # phase 1 sub-plan
    m1-t1-p2.md                           # phase 2 sub-plan
    scoping/
      m1-t1-p1.md                         # scoping for phase 1

  docs-canonical-corrections/             # standalone task plan root
    README.md                             # the task plan (N = 1)
```

## Why this shape

- **Every root gets a folder, no top-level files in `docs/plans/`.**
  Removes the special case that would otherwise distinguish
  standalone task plans (single file) from epic-rooted trees
  (folder). A task plan that grows a second phase doesn't need to
  be moved or restructured; sibling phase docs just appear next
  to it.
- **`README.md` for the root doc.** GitHub renders it
  automatically when a reader navigates into the folder.
  Survives doc-type changes (a folder that starts as a task plan
  and grows into an epic doesn't need a rename).
- **Filename = slug-suffix + optional descriptor.** The slug in
  frontmatter is the identity of record (see
  [`planning/shared.md`](./planning/shared.md) "Plan-doc identity
  (slug)"). The filename is purely for human browsing and may
  carry an optional descriptor after the slug-suffix to aid
  readers scanning the folder; the tool reads slug from
  frontmatter, never from filename.
- **`scoping/` subfolder is transient.** Its contents delete in
  batch at the milestone-terminal PR (or task-terminal PR for
  standalone task plans) per
  [`planning/task-plan.md`](./planning/task-plan.md) "Scoping owns / plan
  owns."

## Slug-to-path mapping

The full slug-to-path mapping under this layout:

| Slug                                    | Path                                                              |
|-----------------------------------------|-------------------------------------------------------------------|
| `<root>`                                | `docs/plans/<root>/README.md`                                     |
| `<root>-m<N>`                           | `docs/plans/<root>/m<N>-<descriptor>.md`                          |
| `<root>-m<N>-t<T>`                      | `docs/plans/<root>/m<N>-t<T>-<descriptor>.md`                     |
| `<root>-m<N>-t<T>-p<P>`                 | `docs/plans/<root>/m<N>-t<T>-p<P>-<descriptor>.md`                |
| `<root>-...` paired scoping             | `docs/plans/<root>/scoping/<slug-suffix>-<descriptor>.md`         |

The `-<descriptor>` portion is optional in every row; pure
slug-suffix filenames (e.g., `m1.md`, `m1-t1-p2.md`) are valid.

## Per-epic milestone numbering

Each epic counts milestones from `m1` independently. Sibling
epics may reuse the same milestone numbers without collision
because the slug's root segment (`<epic-slug>`) disambiguates;
the path's folder segment carries the same disambiguation.

## v0.0 hardcoded

For v0.0 of the workstream-tracker spec, this layout is hardcoded
— the visualization tool walks `docs/plans/<root-slug>/` and
expects the structure above. Project-customizable layout (a
`workstream.toml` or similar) is deferred.
