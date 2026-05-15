# Backlog

The backlog is a flat list of work items captured but not yet
planned. Each entry is a candidate for graduation into a plan-tree
node when picked up. Cross-referenced from
[`planning/epic.md`](./planning/epic.md) and
[`planning/task-plan.md`](./planning/task-plan.md) "Backlog Impact"
sections.

## Location

`docs/backlog.md` at the consumer repo root, sibling to
`docs/plans/`. Single file; not a folder.

## Entry format

Each entry is a sub-section (level 3 heading) carrying a stable
kebab-case slug as the heading anchor and a short body:

```markdown
### <entry-slug>

**Status:** Open

One-line title summarizing the work.

One to three sentences describing the problem, the constraint,
or the opportunity. Frame entries by goal/problem, not solution.
One illustrative option allowed; mark it as one option among
several, not the prescription.
```

**Slug rules.** Kebab-case (lowercase letters, digits, hyphens),
unique within the backlog file. Backlog slugs are flat — no
hierarchical segments, no `m`/`t`/`p` prefixes. They occupy a
separate slug namespace from plan-tree slugs.

**Optional fields.** Projects may add tags, priority/tier markers,
or links to related artifacts as additional bold-labeled lines or
a list under the body. The spec doesn't enumerate these — project
conventions own them.

## Entry lifecycle

The `Status:` line carries one of:

- `Open` — captured, not yet picked up.
- `Graduated — <plan-slug>` — picked up; a plan-tree node now
  carries the work, identified by the named plan slug. The
  `<plan-slug>` is the root slug of the plan-tree node (epic root
  or task-plan root). Em-dash freeform context permitted after
  the slug, mirroring the `Deferred — <reason>` pattern in
  [`planning/shared.md`](./planning/shared.md) "Plan-doc Status."

When work is no longer relevant — cancelled outright, absorbed
into another entry, or otherwise resolved without a plan-tree
node ever being created — the entry is **deleted** in a PR that
records the rationale in the commit message. Deleted entries
survive in git history; the spec does not carry a `Closed` or
`Cancelled` Status because the visualization treats absent
entries as closed (no node to render, no work to track).

## Plan-tree relationship

Backlog entries are NOT plan-tree nodes. They occupy a separate
concept the tool understands:

- **Backlog entries don't get work-instances.** The
  workstream-tracker API has no `backlog` node-type. Until an
  entry graduates, it's a captured intent without runtime state.
- **Graduation creates a plan-tree node with a new root slug
  by default.** When the plan author picks an entry up, they
  pick the plan's root slug (a new epic slug or task-plan slug)
  and set the backlog entry's Status to `Graduated — <plan-slug>`
  pointing at it. Reusing the backlog slug as the plan slug is
  allowed when natural but not required; the two slugs live in
  separate namespaces and a name collision between them is
  permitted.
- **Plan docs reference backlog entries via the `Backlog Impact`
  section.** Each entry the plan touches is named by its slug
  with a one-line statement of effect: **graduate** (entry's
  Status flips to `Graduated — <plan-slug>`), **delete** (entry
  was absorbed or cancelled by the plan; removed in the same
  PR), **split** (entry decomposed into multiple new entries),
  or **shift** (entry's framing changed but it stays Open). The
  PR that lands the plan also updates the backlog file
  accordingly.

## Visualization (v0.0)

The workstream-tracker visualization renders the backlog as a
separate panel or list, distinct from the plan tree. Open
entries are shown; graduated entries are hidden by default but
discoverable via the plan-tree node they graduated into. The
backlog has no Status-driven node colors and no work-instance
markers — it's a captured-intent surface, not a runtime surface.

## Why backlog entries aren't plan-tree nodes

Two reasons:

1. **Different lifecycle shape.** Plan-tree nodes have a
   contract-and-validation lifecycle (`In draft` → `Proposed` →
   `In progress` → `Validating` → `Landed`). Backlog entries
   don't — they're captured-but-undefined intent. Forcing them
   into the same lifecycle either dilutes the lifecycle's
   semantics or pollutes the backlog with state it doesn't need.
2. **No descendants.** Plan-tree nodes have hierarchical
   structure (epic → milestone → task → phase). Backlog entries
   don't — if work is structured enough to have descendants,
   it's structured enough to be a plan. Keeping the backlog flat
   keeps it cheap to capture into.

The graduation handoff (`Graduated — <plan-slug>`) is the seam
between the two surfaces. Once an entry graduates, it's
shoulder-tapped to its plan-tree node and the plan-tree
lifecycle takes over.
