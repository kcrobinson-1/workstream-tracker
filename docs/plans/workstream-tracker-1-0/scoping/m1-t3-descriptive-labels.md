# Scoping — t3 Descriptive tree labels

Transient deliberation doc for `workstream-tracker-1-0-m1-t3`.
Deletes in batch at the m1-terminal PR per
[`task-plan.md`](../../../../spec/planning/task-plan.md) "Scoping
owns / plan owns." Carries no Status field (scoping docs are
transient; the sibling plan alone carries Status).

Pairs with the task plan at
[`../m1-t3-descriptive-labels.md`](../m1-t3-descriptive-labels.md).

## Context summary

t3 is the first task on the read-experience track of m1
([`m1-v0-2.md`](../m1-v0-2.md)). It forks directly from the v0.1
baseline and is independent of the t1/t2 foundation track, so it
can draft and implement without waiting on t1 or t2 — there are
no pending "input from prior task" entries to carry. The task
lifts plan-tree node labels from full slug chains
(`workstream-tracker-1-0-m1-t3`) to human-readable
`<Type> <ordinal>: <short description>` strings, introduces an
optional `short_description` frontmatter field, and parses the
markdown body into a long-description value that t4 will later
render.

## Decisions made at scoping time

Each decision carries a `Verified by:` code citation. Rejected
alternatives are recorded here (deliberation prose has no
audience after the plan lands); the durable contract lives in
the plan doc and is not restated here.

### D1 — N = 1, phase content absorbed inline

t3 ships as a single-phase task plan; no separate phase plan
files. Parsing the frontmatter/body and rendering the label are
sequence-steps toward one outcome — shipping "parsing only"
renders nothing user-visible, and "rendering only" has no parsed
fields to read — so per the level picker they lack independent
value and do not warrant separate plan-tree nodes. The surface
is small and converges on one integration point.

- **Rejected:** N = 2 (a parse phase, then a render phase), the
  shape the milestone doc anticipated as a non-binding estimate.
  Rejected because the two halves share one PR's worth of
  review surface and splitting would create a phase boundary
  with no independently shippable artifact.
- `Verified by:` [`buildTree` in tree.go](../../../../internal/site/tree.go)
  already calls `slugs.Parse(d.Slug, root)` per doc, so the
  label can be computed at the existing single integration point
  without a new tree walk; the render side is one template
  (`{{define "node"}}` in
  [render.go](../../../../internal/site/render.go)).

### D2 — Spec edit ships in the same PR as the parser/render change

The `spec/planning/shared.md` documentation of the new optional
`short_description` field lands in the same implementing PR as
the parser and render changes, not as a separate spec-only PR.
This resolves the milestone's deferred "Spec change PR shape
(t3)" decision ([`m1-v0-2.md`](../m1-v0-2.md) Cross-Task
Decisions).

- **Rejected:** separate spec-only PR landing first. Rejected
  because it opens a window where the vendored spec documents a
  field that no code reads, and the field + its sole reader are
  one contract.
- `Verified by:` [`spec/planning/shared.md` "Plan-doc identity
  (slug)"](../../../../spec/planning/shared.md) is the section
  documenting frontmatter fields; the additive `short_description`
  entry sits adjacent to it.

### D3 — t3 parses the long description but does not render it

The markdown body (everything after the frontmatter fence) is
parsed, trimmed, and carried on the node as a long-description
value. t3 does **not** render it — the milestone explicitly
defers "Long-description rendering location (t4)" to t4. t3's
visible scope ends at the label; the parsed long-description is
an interface t4 consumes.

- **Rejected:** rendering a truncated long description inline in
  t3. Rejected because it preempts t4's surface-shape call
  (tree-inline vs. detail panel) the milestone reserves for t4.
- `Verified by:` [`m1-v0-2.md`](../m1-v0-2.md) Cross-Task
  Decisions ("Long-description rendering location (t4)") and the
  t3 Task Contract ("the markdown body becomes a parsed
  long-description value").

### D4 — Full slug preserved via the `title` attribute

The full slug stays accessible as the HTML `title` attribute on
the label element (native browser tooltip; text is
selectable/copyable from the DOM). No separate copy widget.

- **Rejected:** a dedicated copy-to-clipboard button. Rejected
  as surface the v0.1-bare-bones render doesn't warrant yet;
  `title` satisfies the contract's "tooltip or copy affordance"
  with the minimum surface.
- `Verified by:` [`indexTmpl` in
  render.go](../../../../internal/site/render.go) is built with
  `html/template`, which auto-escapes in attribute context, so
  emitting the slug into `title="…"` is injection-safe.

### D5 — Label grammar and fallback

Label = `<Type> <ordinal>: <short_description>`. `<Type>` is a
display map over the slug's terminal node type
(`milestone`→`Milestone`, `task`→`Task`, `phase`→`Phase`);
`<ordinal>` is the terminal segment's position. Root nodes have
no position segment — their label is the `short_description`
alone, or the slug when absent. Any node missing
`short_description` falls back to the **slug suffix** (the slug
text after the root; the full slug for a root).

- **Rejected:** falling back to the full slug chain (v0.1's
  current display). Rejected because the milestone t3 contract
  specifies the slug-suffix fallback explicitly, and the full
  slug remains reachable via D4's tooltip.
- `Verified by:` [`Slug.NodeType()` in
  slugs.go](../../../../internal/slugs/slugs.go) exposes the
  terminal node type; the milestone t3 contract gives the
  worked example `Task 1: Wire auto-registration`. The ordinal
  is **not** reachable from the current exported surface —
  `Slug.segments` is unexported and there is no terminal-
  position accessor — so D6 adds one rather than parsing the
  slug grammar inside `internal/site`.

### D6 — Add an exported `Slug.Position()` accessor

`internal/site` needs the terminal segment's position for the
label ordinal. The exported `slugs` surface
(`String/Root/NodeType/Parent`) has no such accessor and
`Slug.segments` is unexported. t3 adds
`func (s Slug) Position() (int, bool)` (position + `true`, or
`0, false` for roots), mirroring the `NodeType()` shape.

- **Rejected:** re-splitting the slug string inside
  `internal/site` to recover the ordinal. Rejected because the
  `slugs` package exists specifically to own the slug grammar
  shared between `internal/api` and `internal/site`;
  duplicating the parse defeats that.
- **Disclosure:** `slugs` is a shared package consumed by
  `internal/api` too. The addition is purely additive (new
  method only), so existing API call sites are unaffected;
  flagged here because a shared-package API addition is plan
  content, not an implicit implementer call.
- `Verified by:` [`Slug` exported methods in
  slugs.go](../../../../internal/slugs/slugs.go) are
  `String/Root/NodeType/Parent` only; `segments []Segment` is
  unexported, confirming no current path from a parsed `Slug`
  to its terminal `Segment.Position`.

## Reality-check inputs

Load-bearing claims the plan rests on, falsifier-checked against
current code. The plan's own contract sections are referenced by
name, not duplicated here.

- **Frontmatter is a leading `---`-fenced block; the body is
  everything after the closing fence.** Falsifier: "a plan doc
  starts with body text before `---`." Checked — every doc under
  `docs/plans/` opens with `---\nslug: …\n---`; the test helper
  [`writeDoc` in
  walker_test.go](../../../../internal/site/walker_test.go)
  encodes the same shape. Body extraction by locating the second
  `---` line is sound for the corpus. **Assumption (tagged):**
  goldmark-meta does not surface the post-frontmatter body as a
  string, so `parsePlanDoc` must split the raw source itself
  rather than ask the parser for the body.
- **`short_description` reads via the existing
  `map[string]interface{}` + string-assertion pattern.**
  `Verified by:` [`parsePlanDoc` in
  walker.go](../../../../internal/site/walker.go) already does
  `metaData["slug"].(string)`; `short_description` follows the
  identical idiom with `, _ :=` tolerance for absence.
- **The label can be computed without a new tree walk.**
  `Verified by:` [`buildTree` in
  tree.go](../../../../internal/site/tree.go) already resolves
  `root`, calls `slugs.Parse`, and constructs each `PlanNode`;
  the label-build slots into that existing loop.
- **Render change is localized to one template block.**
  `Verified by:` the `{{define "node"}}` block in
  [render.go](../../../../internal/site/render.go) is the only
  site emitting `.Slug` for display.
- **No API/DB/schema involvement.** t3 is read-path only.
  `Verified by:` the milestone Cross-Task Invariant "Render path
  stays walk-on-every-request" and the absence of any
  `work_instances` or `RegisterRequest` touch in the file
  inventory the plan names.

## Plan-structure handoff

- Doc-type: task plan, N = 1 (D1). Path:
  [`../m1-t3-descriptive-labels.md`](../m1-t3-descriptive-labels.md).
- Required sections present in the plan: Status, Context
  preamble, Goal, Contracts, Files to touch, Validation Gate.
- Optional sections the plan carries: Cross-Cutting Invariants
  (the optional-field fallback threads parser + label-builder +
  template), Naming (new identifiers), Out Of Scope, Risk
  Register, Documentation currency, Related Docs.
- Section variance from
  [`task-plan.md`](../../../../spec/planning/task-plan.md)
  required+optional list: none expected; if the plan adds an
  unlisted section the implementing PR discloses it per
  [`shared.md`](../../../../spec/planning/shared.md) "Section
  variance disclosure."

## Open decisions carried to plan-drafting

None. t3 has no pending input from prior tasks (independent of
t1/t2), and D1–D5 resolve every cross-task decision the
milestone deferred to t3. The plan doc may proceed straight to
the promotion-gate self-review before flipping
`In draft → Proposed`.
