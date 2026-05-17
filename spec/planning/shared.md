# Cross-Level Planning Rules

Cross-level rules that bind every plan-drafting session regardless of
doc-type (epic doc, milestone doc, task plan — see "Taxonomy: node-type
and doc-type" below). Per-level files ([`epic.md`](./epic.md),
[`milestone.md`](./milestone.md), [`task-plan.md`](./task-plan.md)) **reference**
these rules; they do not duplicate them. Rules that bind only at the
implementation layer (the task-plan doc-type, covering both task and
phase node-types) — Planning Depth, the `` `In draft` → `Proposed` ``
promotion gate, the Plan-to-PR Completion Gate — live in
[`task-plan.md`](./task-plan.md), not here, because epic and milestone docs do
not consume those gates.

If a rule below feels load-bearing only at one level, that's a signal
to move it to the per-level file rather than restate it here. The
discipline is: if a rule binds two or more planning levels, it lives
here once.

Four of the rules below — `Plan-doc review stance`, `Cross-Cutting
Invariants section`, `"Verified by:" annotations on load-bearing
claims`, `Falsifiability check on each load-bearing claim` — bind every
level but mean different things at each. Each carries a closing
"What this means at each level" line naming the rule's interpretation
at epic / milestone / plan level. Those interpretations describe
per-level reading of the rule, not additional per-level binding; they
do not re-introduce the per-level duplication trap the layered
authority structure was designed to avoid.

Rule prose throughout the planning spec uses compound forms — "task
plan," "phase plan," "task-level," "phase-level" — where bare "plan"
would create ambiguity between the doc-type and the broader plan-tree
category. See [`task-plan.md`](./task-plan.md) "Noun discipline" for
the full rule. Bare "plan" remains correct when the rule binds every
plan-tree doc (the cross-level rules below are the canonical case).

> **TODO (workflow portability).** Several rules in this spec
> ("Plan-doc review stance," "Section variance disclosure," the
> "Estimate Deviations" callout under "Plan content is a mix of
> rules and estimates," and most of the Plan-to-PR Completion Gate
> in [`task-plan.md`](./task-plan.md)) embed a GitHub-style PR
> review workflow — `## Review Stance` headings in the PR body,
> `## Estimate Deviations` after `## Documentation`, etc. The
> rules' protective intent is universal but the scaffolding
> assumes that workflow shape. Revisit when the spec's first
> non-GitHub consumer surfaces, or when v0.1 makes review-surface
> configurability worth the cost; for v0.0, the GitHub framing is
> kept as-is.

## Taxonomy: node-type and doc-type

The planning model carries two parallel concepts. Distinguishing them
keeps later rules unambiguous.

**Node-type** — hierarchical position in the plan tree. Four values:
**epic**, **milestone**, **task**, **phase**. Encoded by slug segment
prefixes (see "Plan-doc identity (slug)" below). Consumed by tools
that walk the tree and by the workstream-tracker API for work-instance
registration. The four-segment cap on slugs (epic → m → t → p) is a
forcing function: if a node legitimately needs a fifth segment, the
right fix is to revisit whether the parent at one of the higher levels
was decomposed correctly, not to extend the taxonomy.

**Doc-type** — which authoring-rule set applies. Three values: **epic
doc**, **milestone doc**, **task plan**. Bound by the per-level
authority files: [`epic.md`](./epic.md) carries epic-doc rules,
[`milestone.md`](./milestone.md) carries milestone-doc rules,
[`task-plan.md`](./task-plan.md) carries task-plan rules and binds both
task-level and phase-level docs because their authoring rules are
near-identical.

The axes are not the same thing. A single doc-type can apply to two
node-types ([`task-plan.md`](./task-plan.md) covers task and phase). Most of the
time the axes line up: a task-level node carries a task plan, a
phase-level node carries a phase plan under task-plan rules. The rules
permit deviations but don't require them.

One useful consequence: the **N=1 collapse**. When a task has exactly
one phase, the author writes one task-plan doc covering the single
phase's content (not two separate docs), and the agent registers one
work-instance at the task slug (not two). Both axes collapse together.

## Plan-doc identity (slug)

Every plan-tree doc — epic, milestone, task plan, phase plan — carries
a **slug** in YAML frontmatter as its stable identifier:

```yaml
---
slug: <slug>
Status: <status>
---
```

**Slug format.** Hierarchical, derived from the root slug. The root
slug is short and human-readable (e.g., `madrona-feedback` for an
epic, `docs-canonical-corrections` for a standalone task plan).
Descendants append their position to the parent:

- `madrona-feedback-m1` — milestone 1 of the epic
- `madrona-feedback-m1-t2` — task 2 of that milestone
- `madrona-feedback-m1-t2-p3` — phase 3 of that task

Segment prefixes are `m` for milestone, `t` for task, `p` for phase.
Root slugs use kebab-case (lowercase letters, digits, hyphens), and no
kebab-delimited token of a root slug may be a bare position segment
(`mN`/`tN`/`pN`) — roots are descriptive names, so `m1` is never a
valid root. After the root, position segments must appear in order and
at most once each: optional `mN`, then optional `tN`, then optional
`pN`. Maximum four segments end-to-end (root + m + t + p); the
four-segment cap is the forcing function described under "Taxonomy"
above.

**The slug encodes position at creation time, NOT current position.**
If the plan tree is later reorganized — a phase becomes a different
number, a task moves under a different milestone — the slug **stays
the same**. The slug is the identifier; the descriptive
position-encoding is a creation-time convenience. Tools that need
*current* position read it from the tree (folder hierarchy or DB),
never by parsing the slug.

This is a footgun: a reader seeing `madrona-feedback-m1-t2-p3` will
assume the doc is currently at milestone 1, task 2, phase 3, which may
or may not be true after reorganization. Authors and reviewers treat
the slug as opaque identity; only initial-creation slug generation
parses position from it.

**Slug generation.** Root slugs are author-supplied and validated for
format on registration with the workstream-tracker API. Descendant
slugs are server-generated by default: the agent supplies the parent
context and node-type, and the server picks the next available
position number and returns the constructed slug.

A descendant slug may instead be **pre-declared** in a plan-tree
doc's frontmatter — the case a parent-promotion stub creates (see
"Parent-doc child contracts" below). A pre-declared slug is
**author-supplied** rather than server-generated (the same
author-supplied/validated posture root slugs have, as opposed to
server construction), but it is validated against the **full
descendant slug-grammar** in "Slug format" above — the root slug
followed by its ordered `mN`/`tN`/`pN` position segments — which
is the same well-formed-descendant validation the exact-slug
create-or-attach registration path applies. Root-slug validation
is **not** applied to it: a root slug forbids the very
`mN`/`tN`/`pN` position segments a descendant slug requires, so
validating a pre-declared child slug as a root would reject every
well-formed descendant.

A parent-promotion stub's pre-declared slug additionally must be
**parent-scoped**: it is the promoting parent's own slug extended
by exactly one further position segment of the child's level
(`mN` for an epic's milestone child, `tN` for a milestone's task
child, `pN` for a task plan's phase child) — not merely
globally-well-formed under "Slug format." A globally-well-formed
slug that names a different root, or the parent's own slug
unextended, is an invalid seed: the promoting PR's diff is
checkable against this (the seeded slug starts with the
promoting parent's slug plus one level-appropriate segment), and
a seed that fails it is a seeding defect to fix before the flip.
This is an authoring-time constraint on what a correct seed
writes; it is distinct from — and does not add — any server-side
parent-context check at assertion time. The exact-slug
create-or-attach path still trusts the caller's slug and can
attach to a typoed or non-existent node; that runtime
trust-the-caller property is the orphan/unattached-work-instance
residual already accepted and deferred to the epic's triage-zone
concern (it is not re-litigated or closed here).

Writing the slug into a file is a
*declaration of identity only*: it
performs no server-side creation — no node, no work-instance, no
allocation call, no consumption of the server's position counter. The
slug becomes known to the server only when a session for that
descendant later **asserts** it via the exact-slug create-or-attach
registration path ("register at exactly this slug"), not by the
server generating one. Server-side descendant generation remains the
path for any registration with no pre-declared frontmatter slug. The
two paths are mutually exclusive per descendant — a slug is either
pre-declared and asserted, or server-generated, never both — so no
second allocator exists to diverge from the position counter.

**Slug is identity; path is layout.** The path the file lives at is
governed by the layout convention (see
[`planning-doc-location.md`](../planning-doc-location.md)). At
creation time the slug determines the initial path; after that, the
slug is immutable and the path can change freely without touching the
slug.

**Optional `short_description` field.** A plan-tree doc may carry an
optional `short_description` string in frontmatter, alongside `slug`
and `Status`:

```yaml
---
slug: <slug>
Status: <status>
short_description: <one-line human-readable summary>
---
```

It is a short human-readable summary used to render the doc's
plan-tree node as `<Type> <ordinal>: <short_description>` (e.g.,
`Task 3: Descriptive tree labels`) instead of the raw slug chain.
The field is **optional and additive**: a doc that omits it remains
valid and renders with no warning, error, or skip (the node falls
back to its slug suffix, with the full slug always reachable via the
node's tooltip). Pre-existing docs and vendored spec consumers are
unaffected by its absence.

**Optional `related_prs` field.** A plan-tree doc may carry an
optional `related_prs` list in frontmatter — a YAML block sequence
of **absolute-URL strings**, each identifying a pull request
related to the node (e.g. a `.../pull/N` link):

```yaml
---
slug: <slug>
Status: <status>
related_prs:
  - https://github.com/<owner>/<repo>/pull/<n>
  - https://github.com/<owner>/<repo>/pull/<m>
---
```

It renders the doc's plan-tree node with an inline, author-curated
list of links to those PRs. Each entry is an absolute URL
(`http://`/`https://`); `#NNN` / `owner/repo#NNN` shorthand is not
expanded — a non-URL entry renders as plain text rather than a
link. Like `short_description`, the field is **optional and
additive**: a doc that omits it (or carries an empty list) remains
valid and renders with no warning, error, or skip, and pre-existing
docs and vendored spec consumers are unaffected by its absence.

## Plan-doc Status

Every plan-tree doc carries a `Status` field in frontmatter, alongside
`slug`. The canonical lifecycle values are:

- `In draft` — drafting underway; the doc is not yet ready for an
  implementing agent.
- `Proposed` — drafting complete; the doc has passed its promotion
  gate and is ready for implementation.
- `In progress` — an implementing agent has begun work; not yet
  landed.
- `Validating` — code has merged but a post-merge validation gate
  (manual walkthrough, prod smoke, async verification) is outstanding.
  **Optional**: projects without a meaningful post-merge gap skip
  this state and transition `In progress` → `Landed` directly.
- `Landed` — implementation merged; validation (if any) passed; the
  plan reached its goal.
- `Deferred — <reason>` — drafting is intentionally paused. The
  em-dash and freeform reason are part of the convention; the
  canonical prefix `Deferred` is what tooling matches on. See
  "`Deferred` Status semantics" below.

A `## Status` markdown heading inside the doc is permitted and useful
for human readers, but **the frontmatter field is authoritative**.
When the two disagree, the frontmatter wins.

**Unknown Status values render gracefully.** Tools that visualize the
plan tree fall back to a neutral color for any unrecognized value and
surface it as-is in the UI. Authors experimenting with project-local
extensions don't break the tool; reviewers see the unfamiliar value
and can decide whether to upstream it.

**Status is one of two state channels.** The other is **work-instance
state** (`active` / `completed` / `abandoned`), tracked in the
workstream-tracker DB by the agent registration API. Plan-doc Status
describes the durable lifecycle of the *planning artifact*;
work-instance state describes the runtime fact of whether an agent is
*currently doing work* against a node. The two evolve independently:

- A plan-doc node can have zero, one, or many work-instances attached
  over its lifetime.
- An agent registering doesn't flip Status; a work-instance completing
  doesn't auto-advance Status.
- They correlate at terminal moments (the PR that flips Status to
  `Landed` is typically the same PR the agent completes against), but
  no rule binds them in lockstep.

In the workstream-tracker visualization, plan-doc Status drives node
*color* and work-instance state drives the *marker* showing which
agent is currently attached.

### `Deferred` Status semantics

Plans, milestone docs, and scoping docs whose drafting is
intentionally paused — typically because the work has been resequenced
behind other priorities, an upstream dependency hasn't decided yet,
or the deliverable has been moved to a future epic — carry Status
`Deferred` (exact-match canonical prefix), followed by an em-dash and
freeform human-readable context: `Deferred — <reason>`.

The canonical `Deferred` prefix is what status-tracking queries match
against; the post-em-dash reason is freeform context for humans, NOT
exact-match-checked. Examples:

- `Deferred — moved behind higher-priority work in the sequence`
- `Deferred — pending decision on upstream dependency`
- `Deferred — moved to future epic`

While Deferred, the doc's content is **non-prescriptive**: future
planning sessions that resume the work re-derive every goal,
sequencing, decision, invariant, and risk against the actually-merged
code at resume time and are not bound by the choices recorded in the
deferred draft. The protective intent is the recurring trap a Deferred
state otherwise hides: a future planner reads the deferred doc, treats
its decisions as settled because they look complete, and silently
inherits assumptions that the original drafting session never expected
to bind. Authors of a Deferred doc must state the non-prescriptive
framing explicitly in the doc's leading prose so the future reader
can't miss it.

State transitions out of `Deferred`:

- **`Deferred` → `In draft` (resumption).** When the work becomes
  next-up, the resuming planner flips Status back to `In draft` and
  re-runs its `` `In draft` → `Proposed` `` promotion gate from
  scratch (for task and phase plans, the gate in
  [`task-plan.md`](./task-plan.md); for epic and milestone docs,
  the parent-doc gate below). The previous deliberation
  becomes input to consider, not contract to respect.
- **`Deferred` → (deletion).** If the work is cancelled outright
  (epic re-scoped to drop it, or absorbed by a different epic), the
  doc is deleted in the same PR that records the cancellation
  rationale; the deferred content survives in git history.

Do not invent additional states adjacent to `Deferred` (e.g.
`Deferred draft`, `Deferred — non-prescriptive`, `Paused`, `On hold`,
`Frozen`). The em-dash freeform reason is the supported affordance for
context. This rule is the exact-match label discipline (see "Quote
labels whose enforcement depends on exact-match matching" below)
applied to the Status lifecycle: `Deferred` is the canonical token,
and adjacent descriptive variants break queryability the same way
paraphrased Status strings do.

## Parent-doc `In draft` → `Proposed` promotion gate

Epic docs and milestone docs carry the same `Status` lifecycle as
task and phase plans (see "Plan-doc Status" above), but the
`` `In draft` → `Proposed` `` promotion gate in
[`task-plan.md`](./task-plan.md) binds task and phase plans only —
epic and milestone docs do not load that file. This is the
symmetric gate for parent docs. It lives here (not in
[`task-plan.md`](./task-plan.md), not duplicated into
[`epic.md`](./epic.md) and [`milestone.md`](./milestone.md))
because it binds two doc-types and the layering discipline keeps a
two-or-more-level rule here once.

A parent doc's drafting need not be a single pass. A multi-pass
session can lay out the milestone or task structure with explicit
deferrals, then resolve them, then flip Status. While resolution
is pending the doc carries Status `In draft`; flipping to
`Proposed` claims the parent doc is ready for its children's
planning sessions to consume — the child set is locked and each
child's WHAT contract (see "Parent-doc child contracts" below) is
decision-complete.

The flip fires on the PR that locks the parent doc's scope.
Before the flip:

- **Read end-to-end as a coherent whole.** Re-read the full
  parent doc. Look for contradictions between sections (a child
  WHAT contract that conflicts with a Risk Register mitigation, a
  sequencing-rationale claim a child contract contradicts, an
  Out-of-Scope deferral a child contract silently re-includes).
- **Decision-completeness on child contracts.** Walk the
  `Milestone Contracts` / `Task Contracts` section for deferral
  phrases that name the parent-doc drafting session itself as the
  resolver. Each child's WHAT is either resolved concretely or
  the child is explicitly marked scope-not-yet-locked (carrying
  its name without a contract) — not deferred to a moment that
  has already passed.
- **Walk the `Verified by:` rule against every load-bearing
  claim** at this level's interpretation (see "`Verified by:`
  annotations on load-bearing claims" → "What this means at each
  level" below).
- **Re-confirm reality-check inputs** the parent doc rests on
  against current code; stale references are updated.
- **Confirm the always-on authoring rules, explicitly.** Three
  further obligations bind the parent doc continuously through
  always-on rules and get their final explicit confirmation here,
  mirroring how the `Verified by:` step applies an always-on rule
  universally at the flip: (1) every required section is present
  and any divergence is disclosed, per the "Required and optional
  sections" rule in this doc's applicable per-level file and
  "Section variance disclosure" below; (2) no content has
  descended to implementation prescription, per "Plans describe
  contracts, not implementation" below; (3) the parent doc
  otherwise conforms to the broader always-on spec — the
  cross-level rules indexed above and its per-level file's rules.
  Naming them here makes the gate self-contained — not a new
  gate-only requirement; the task/phase gate in
  [`task-plan.md`](./task-plan.md) carries the symmetric step.
- **Seed child skeleton docs.** In the same PR that flips the
  parent to `Proposed`, seed a skeleton doc for every child the
  locked child-contracts section names, per "Parent-doc child
  contracts" below ("Parent-promotion stub seeding") — which
  governs what a skeleton carries, the slug-declaration
  mechanism, the skeleton exemption, and the best-effort,
  observable, no-clobber-on-re-run framing this step inherits
  rather than restates. The task/phase gate in
  [`task-plan.md`](./task-plan.md) carries the symmetric step for
  a task plan's N ≥ 2 phase children.

Failures surface either as resolutions (apply edits before
flipping) or as blockers the user triages before the flip. A
parent doc flipped to `Proposed` without this walk is the same
drift shape as a task plan flipped to `Proposed` without its
promotion-gate walk — the Status claim is wrong. The `In draft`
and `Proposed` tokens are matched by exact string per "Plan-doc
Status" above. Transitions out of `Proposed` for parent docs are
not gated here — a parent doc has no implementing PR of its own;
its terminal state follows from its children per the per-level
files.

## Plans describe contracts, not implementation

A plan describes what the implementation must achieve — the
conditions that must hold for the implementation to satisfy the
plan. The implementation is the specific choices the implementer
makes to satisfy the contract.

When plan content descends to the level of implementation choice,
one of two failure modes follows:

- **The plan no longer fully covers the contract.** When a plan
  names a subset of a category, only the named subset is
  constrained; the rest of the category is left uncovered.
- **The plan contradicts the implementer's reasonable choice.**
  When a plan prescribes a specific sequence or technique, it
  contradicts the implementation the moment the implementer's
  reasonable choice differs.

Either failure mode is a sign the plan was written below the
right altitude: above it, the contract is invariant under the
implementer's reasonable choices; below it, it isn't.

The protective intent is structural, not formatting preference —
implementation-shaped content in plans attracts implementation-shaped
review, which compounds the rule layer rather than the product.

### Structural surface

**No fenced code blocks of any kind in plan or scoping docs.**
Inline backticks for identifiers, file paths with optional
`:line` suffixes, and one-line type or function names embedded
in prose are fine — those are citations, not code under
contract. **Inline backticks for executable expressions or
predicate spellings are still code, even at one line** — a
literal predicate, a regex literal, or a function-body fragment
attracts the same code-shaped review fenced blocks do and
belongs in the implementing PR alongside them. Anything more
belongs in the PR that implements the plan (commit message,
code, comment), not in the plan doc itself.

### Reviewer-fix discipline

When a reviewer flags content that violates this rule, the
response is shape-specific:

- **Structural violation (code-shaped content in a plan):**
  remove or summarize the snippet. Do not fix the code in place.
  Code-correctness iteration belongs in the PR that implements
  the plan.
- **Altitude violation (plan descended to implementation
  choice):** loosen the prescription to contract altitude. Do not
  patch the technique in place. Patching makes the plan hostage
  to the next implementation choice; loosening makes the plan
  resilient.

Exception in either case: if the comment surfaces a genuine
design flaw whose phrasing happens to be code or technique (an
ordering race, an invariant violation, a coverage gap), fix the
prose contract in the plan and move the code or technique to
the implementation PR — don't fix both in the plan.

### Recurring traps

Illustrative examples of the rule's failure mode in practice.
Not an exhaustive list; new traps are appended as they're found.

- **Subset enumerated, or content prescribed, where the contract
  is about category or shape.** Two shapes of the same failure
  mode — the plan descends below contract altitude into specifics
  that the contract didn't ask for. (a) *Subset-enumerated.* A
  plan enumerates a specific subset of a category (privilege
  types, error classes, role names) when the contract calls for
  the full category. Only the enumerated subset is constrained;
  the rest is left uncovered. (b) *Content-prescribed.* A plan
  prescribes specific content (specific config text, specific
  prose, specific copy strings) when the contract is about shape.
  The plan-doc takes on the same factual-drift surface the
  implementing doc would carry — and gets no benefit, because the
  implementing PR can verify the shape against the authoritative
  carrier in the same change. The fix for either shape is to
  describe what the implementing doc must do (the category or the
  shape), not which specific members or content satisfy it.
- **Specific source named where the contract is about coverage.**
  A plan names a specific introspection source (catalog view,
  library function) when the contract is about behavior the
  source must produce. When the named source's coverage gap
  exposes itself, the prescribed technique no longer satisfies
  the contract. The fix is to describe what coverage the source
  must produce, not which source is used.
- **Trajectory prescribed where only the end state is in the
  contract.** A plan describes a specific sequence of
  intermediate steps — what order to do things, what commands to
  run, how to decompose the work — when the contract is only
  about what must hold at the end. Trajectory is the
  implementer's choice; the contract is the post-implementation
  state. Over-prescribing trajectory either contradicts the
  implementer's reasonable choice (creating internal-coherence
  bugs) or imports middle-state commitments the plan shouldn't
  care about. The fix is to describe the end state and let the
  implementer pick the trajectory. State transitions that ARE in
  the contract — a validation gate that must pass before a
  subsequent step, an output that a later step consumes, an
  artifact that must reach a specific state before another can
  act on it — are contract states, not trajectory; Planning
  Depth's requirement to "insert steps at the correct point in
  the sequence" applies to those contract states and does not
  extend to prescribing trajectory between them.

### Reviewer-side companion

Plans stay at contract altitude across all sections; reviewers may
flag descents to implementation choice as a structural issue rather
than reviewing the specific technique. The companion rule on the
reviewer side lives in the next section.

## Plan-doc review stance

Plan- and scoping-doc PRs carry the same review hazard the "Plans
describe contracts, not implementation" rule above addresses on the
supply side: code-trained reviewers default to line-level findings
(syntax-shaped, command-ordering-shaped) even on prose-shaped
content, and the cascade compounds when each finding produces either
a snippet patch or a new rule. The supply-side fix (no code blocks,
no technique prescription) is cheaper if the demand-side stance is
also explicit, so reviewers don't reach for the wrong toolkit at
review-open time.

A PR whose primary diff is in `docs/plans/**` (plan doc or scoping
doc) carries a `## Review Stance` section in the PR body. Canonical
wording, which the PR may copy verbatim or adapt to the doc's
content:

> **Review Stance.** Review for: factual accuracy of claims about
> the codebase and supporting services (with attention to
> `Verified by:` citations); internal coherence between contract /
> invariant / validation sections; decision-completeness on
> Contracts; that estimate-shaped sections are labeled per the
> "Plan content is a mix of rules and estimates" rule in
> `shared.md`; that load-bearing claims pass the "Falsifiability
> check" rule in the same file. **Do not review for line-level
> correctness** on identifiers, signatures, commands, or
> implementation sequencing embedded in prose — per "Plans
> describe contracts, not implementation" in `shared.md`, fenced
> code blocks and technique prescriptions belong in the
> implementing PR, and anything short enough to live inline in
> prose is a citation, not code or technique under contract.

The stance lives in the PR body (where reviewers see it at
review-open time), not in the plan doc itself. Adding it to the
durable doc would be its own form of bloat — every plan and scoping
doc would carry the same paragraph. The PR body is the natural
place because the stance shifts the *review interaction*, not the
durable artifact.

Plan- and scoping-doc PRs include the `## Review Stance` section.

**What this means at each level.** The same Review Stance binds every level;
what shifts is the load-bearing surface the stance protects. At
epic level the stance keeps reviewers from prescribing
per-milestone scope or per-task / per-phase technique that hasn't
been planned yet. At milestone level the stance keeps reviewers
from prescribing per-task trajectory or per-PR contracts that
belong to the task plan. At plan level the stance keeps reviewers from
prescribing per-file technique or implementation sequencing that
belongs to the implementing PR.

## Cross-Cutting Invariants section

List the cross-cutting invariants that thread through multiple files
in their own `## Cross-Cutting Invariants` subsection, distinct from
per-file contracts. Per-file contracts describe what one module
does; cross-cutting invariants describe relationships that must
hold simultaneously at every call site and break silently when one
site drifts (examples: "a shared reference clock advances on every
user action that changes filtered output," "every dialog exposes
an accessible name via `aria-label` or `aria-labelledby`," "derived
state for modal return-focus must survive the close transition,
not null out with the trigger state"). Aim for 2–4 one-line
invariants. Without naming these, implementer self-review checks
each file in isolation and misses bugs that only appear when two
sites disagree about the same rule; reviewer rounds then
rediscover the gap one call site at a time. The plan's job is to
name the rule once so self-review can walk every site against it.

**What this means at each level.** At epic level, the invariants thread
across milestones — a capability constraint multiple milestones
must respect, or a posture decision that binds the whole arc. At
milestone level, the invariants thread across tasks — a contract
every task must preserve, or a coordination rule that binds the
task set. At plan level, the invariants thread across files
within the plan's implementing PR(s) — the original framing of
this rule.

## Parent-doc child contracts

A parent doc (epic, milestone) and a task plan with N ≥ 2 phases
states, for each direct child, that child's **WHAT** contract;
the child's **HOW** stays in the child's own doc. This binds
three doc-types at once (epic → milestone, milestone → task,
task → phase), so per the layering discipline it lives here once
and the per-level files reference it rather than restating it.

- **WHAT the parent states per child:** the child's end result
  (what is true when the child is done), the sibling interfaces
  it produces or consumes (what later siblings build on or feed
  it), and what it preserves (existing behavior that must still
  hold after the child lands). One short block per child, on the
  order of 2–4 lines.
- **HOW the parent does NOT state:** the child's file inventory,
  function or signature shapes, specific commands, validation-
  gate specifics, execution-step ordering, or risk register.
  Those are scoped against actually-merged code at the child's
  own planning session and live in the child's doc.
- **Per-level section name.** This rule's realization is a
  required-when-applicable section in each parent doc-type:
  `Milestone Contracts` in an epic doc, `Task Contracts` in a
  milestone doc, `Phase Contracts` in a task plan when N ≥ 2.
  The per-level files ([`epic.md`](./epic.md),
  [`milestone.md`](./milestone.md),
  [`task-plan.md`](./task-plan.md)) carry the section in their
  "Required and optional sections" list and cite this rule as
  the authority rather than restating the WHAT/HOW split.
- **When required.** The section is required once the parent doc
  locks the child's scope (the child set is fixed and each
  child's WHAT is decision-complete). A parent doc whose child
  scope is still open carries the child names without contracts
  until the locking session fills them.

The recurring trap this closes: the parent-doc anti-scope rule
in the applicable per-level file bars the parent from scoping a
child's HOW and was read as also barring the child's WHAT —
leaving a scope-locked parent doc able to name its children but
not contract them. WHAT-contracting is
required; HOW-scoping stays barred. The companion phrasing lives
in those per-level anti-scope rules.

**What this means at each level.** At epic level the children are
milestones and the WHAT is milestone-granular (end result,
cross-milestone interfaces, what the milestone preserves). At
milestone level the children are tasks. At task level (N ≥ 2)
the children are phases; because phases are sequence-steps toward
one task outcome rather than independent-value units, the WHAT is
thinner — chiefly the sibling-interface handoff and preserved
behavior, since a phase has no independent end result of its own.

**Per-child table presentation.** The `Milestone Contracts` /
`Task Contracts` / `Phase Contracts` section presents its children
as a table, one row per child, carrying at minimum: the child's
slug; a short description (one line); a long description (the
child's end result and what it preserves); and a high-level
statement of how that child's deliverable feeds its siblings'
inputs (the inter-child interface, stated as the wiring between
children, not a bare per-child name). The table is the WHAT
contract's presentation form — it carries the same WHAT the
per-child bullets above define and adds no HOW. A section that
lists child names without the sibling-wiring statement has named
its children but not contracted how they compose.

**Parent-promotion stub seeding.** When a parent doc locks its
child scope and passes its `` `In draft` → `Proposed` `` promotion
gate (the parent-doc gate above for an epic or milestone doc; the
task/phase gate in [`task-plan.md`](./task-plan.md) for a task
plan with N ≥ 2 phases), the same PR that flips the parent to
`Proposed` also seeds a skeleton doc for every child the locked
child-contracts section names. A seeded skeleton carries:

- **Frontmatter:** a `slug` pre-declared per "Slug generation"
  above — the promoting parent's own slug extended by exactly one
  level-appropriate position segment for this child
  (parent-scoped, author-supplied, format-validated; the
  file-write declares identity and performs no server-side
  creation) — a `Status` of `In draft`, and a
  `short_description`.
- **Inherited contract:** the parent's WHAT block for that child
  and any illustrative examples the parent states for it, copied
  into the skeleton so the child's own drafting session starts
  from the locked contract rather than rediscovering it.
- **Level-appropriate register:** the descriptions are written at
  the child's level — milestone skeletons in product terms (what
  the milestone enables, or the technical groundwork it lays for a
  later milestone); task and phase skeletons in technical terms
  still grounded in the product outcome the work serves.

A task plan with N = 1 absorbs its phase content inline and has no
separate phase children — there is nothing to seed at N = 1.

A seeded skeleton is, until its own `` `In draft` → `Proposed` ``
drafting session runs, **exempt** from the "Required and optional
sections" rule in its per-level file and from "Plans describe
contracts, not implementation" above: a skeleton legitimately
carries only the frontmatter and inherited-contract content above.
The exemption ends when that child's drafting session begins; the
skeleton is not a standing variance and needs no "Section variance
disclosure" call-out while it remains a skeleton.

Seeding is a **prompted, best-effort, observable** obligation, not
a determinism guarantee. The promotion-gate step that prompts it
(below, and the symmetric task/phase step in
[`task-plan.md`](./task-plan.md)) directs the promoting agent to
do the seeding; a missed or imperfect seeding is an accepted
residual addressed through prompt engineering and the rendered
tree's visibility of an un-seeded child, not through tool
enforcement. The contract is satisfied by specifying the prompted
behavior and keeping a miss observable; it does not promise a miss
cannot happen. On a gate re-run against a re-opened parent (a
`Deferred → In draft` resumption per "Plan-doc Status" above),
re-seeding splits by what the child has become:

- A child **still a pristine skeleton** (not yet drafted, no
  content beyond the seeded frontmatter and inherited contract)
  is **re-synced** to the re-locked parent contract: if the
  parent's child-contract text, slug, or descriptions changed
  during the re-opening, the skeleton's inherited-contract and
  frontmatter content is refreshed to match, since the skeleton's
  only purpose is to carry the *currently locked* contract
  forward and a stale skeleton would start the child's drafting
  from outdated requirements.
- A child **already drafted or advanced past skeleton** is left
  as-is; re-seeding never clobbers content a drafting session
  has put there. A divergence between such a child and a changed
  parent contract is reconciled by that child's own drafting,
  not by re-seeding overwriting it.

## Section variance disclosure

Each plan-tree doc-type carries a "Required and optional sections"
rule in its per-level file ([`epic.md`](./epic.md),
[`milestone.md`](./milestone.md), [`task-plan.md`](./task-plan.md)) naming
the sections that doc-type typically uses. When a particular doc
legitimately diverges from that listed shape — either by adding a
section not on the list, or by skipping a required section because
it genuinely doesn't apply (a decision plan whose deliverable is
the recorded decisions skipping Contracts, for example) — the PR
introducing the divergence calls it out in the PR body's
`Documentation` section: which doc, which header, whether added
or skipped, and why the listed shape didn't fit.

Variances are expected to be rare and illustrative. When the same
variance recurs across two or more docs of the same type, the next
PR that touches it updates the per-level required-or-optional list
rather than letting the variance accrete as one-off prose.

## Plan content is a mix of rules and estimates — label which is which

A plan doc carries two kinds of content: **rules** that bind the
implementation (Cross-Cutting Invariants, Contracts, Validation
Gate, Self-Review Audits, Out Of Scope deferrals) and **estimates**
of what the implementation will look like (file inventory under
"Files to touch — new / modify / intentionally not touched," step
counts, commit boundaries, sometimes per-section LOC predictions).
Estimates are the planner's best guess at plan time about scope
shape; reality during implementation may surface that an estimate
was wrong without any rule being wrong. Plan authors **must**
structure the doc so the distinction is visible to both human
reviewers and implementing agents:

- Sections that bind (Cross-Cutting Invariants, Contracts,
  Validation Gate, Goal, Self-Review Audits, Risk Register
  mitigations, Out Of Scope) are rule-shaped by section name and
  don't need extra labeling.
- Sections that estimate ("Files to touch — new," "Files to
  touch — modify," "Files intentionally not touched," "Execution
  steps" sequencing, "Commit boundaries") **must** carry a
  one-line preface naming them as estimates of the expected shape,
  explicitly admitting that implementation may revise them when
  a structural call requires deviating. The list `Files
  intentionally not touched` is the recurring trap — its name
  reads as a hard prohibition but the underlying claim is
  "we don't expect to need to touch these," not "implementation
  must not touch these." Same for "intended commit boundaries":
  the planner's split is an estimate of cohesive review chunks;
  the implementer can refine.
- Implementers reading a plan: distinguish before deviating.
  Deviating from a rule means the rule is wrong and the plan
  needs to be revised in this PR before the deviation lands;
  deviating from an estimate is normal and is handled via the
  "Estimate Deviations" callout in the PR body (see
  [`task-plan.md`](./task-plan.md) "Plan-to-PR Completion Gate"). When the
  call is unclear, ask.

Plans must label their estimative sections per the bullet above.

## "Verified by:" annotations on load-bearing claims

Load-bearing claims in the plan about the codebase or supporting
services (including data-layer contracts, RPC behavior, type-system
or function-signature contracts, validation-procedure claims,
dev-tool semantics, URL contracts and route topology, copy that
names artifacts or destinations, framework / vendor behavior —
hosting-platform routing and CDN behavior, database engine
semantics, application-framework conventions, build- and test-tool
semantics — and any other claim asserting something specific about
how the codebase or an external service behaves) **must** carry an
inline "Verified by:" reference to the source that proves them.
Acceptable sources: a code citation (file path with optional symbol
or section anchor, or `:N-M` line range — see "Anchor preference"
below), generated test output, an already-merged sibling artifact,
or the upstream / vendor documentation URL for claims about
external-service behavior the codebase doesn't contain proof of.
"Per scoping doc" or "per epic" are not acceptable verification
sources. Claims that cannot carry a verification reference are
re-phrased as assumptions (clearly tagged as such) or removed.
This is not formatting preference — it is the protective check that
keeps the reality-check pass (named in [`task-plan.md`](./task-plan.md)
"Reality-check pass before plan-drafting") from being rolled back
during plan-drafting. The trigger enumeration is illustrative,
not exhaustive: the rule binds any load-bearing claim about the
codebase or supporting services, and agents do not get to argue
"my claim isn't on the list, so the rule doesn't apply."

**Anchor preference: prefer symbol- or section-anchored references
over `:N-M` line ranges when the cited target has a stable name.**
Function names, JSON key paths, exported symbols, markdown section
headings — any stable identifier the cited file already carries —
survive structural edits to the file (a strip of unrelated lines
shifts every `:N-M` citation but doesn't move the symbol). Line
numbers remain permitted (and remain required at write time per
the retrieve-before-citing rule) but should be treated as
**directional navigation aids, not exact contracts** — line
numbers drift as files get edited, and there is no docs-equivalent
of an IDE's "rename all references" tool to keep them in sync
across the plan tree. Reviewers do **not** flag stale or imprecise
`:N-M` citations as findings; the symbolic content of the
surrounding prose is the load-bearing anchor, and an off-by-N line
range is a navigation hint to refresh, not a review issue. This
preference is not a prohibition: paragraph-level precision inside a
long section, files without stable named targets (plain config,
flat data definitions), and citations the author judges clearer
with line numbers all stay valid uses.

**What this means at each level.** At epic level, citations target
capability framing and external constraints (vendor docs that
define the capability surface, prior-art product decisions, the
upstream policies the epic depends on). At milestone level,
citations target cross-task coordination decisions and the
upstream/downstream contracts the milestone locks. At plan level,
citations target code, generated test output, or vendor docs for
the specific contracts the plan binds — the original framing of
this rule.

## Quote labels whose enforcement depends on exact-match matching

When a plan references a label whose value is checked or queried
by exact-string match (Status strings used for plan-state tracking,
branch naming conventions automation watches for, exact phrases a
rule forbids paraphrasing of), copy-paste from the source with a
`path:line` citation rather than retyping. Paraphrasing silently
weakens rules whose enforcement value depends on the exact string.

The Status lifecycle defined under "Plan-doc Status" above is the
canonical case: tokens like `In draft`, `Proposed`, `In progress`,
`Validating`, `Landed`, and the `Deferred` prefix in
`Deferred — <reason>` are matched by exact string. Descriptive
variants (`Drafted`, `Validation pending`, `Almost there`) silently
break the queryability invariant; tools and reviewers cannot tell
the variant apart from a project-local extension.

Ordinary identifiers (env-var names, file paths, function names,
fixture names) do not need this treatment — code blocks and
adjacent file references already carry the spelling, and citing
every identifier adds noise without protective value. Citation is
required only when the plan claims something specific about an
identifier's wording, when the plan is the artifact introducing
it, or when downstream automation or status tracking depends on
its exact spelling. Apply the same exact-match discipline to
trigger clauses: when citing a rule's trigger ("the trigger that
catches this plan"), read every clause and quote the one that
catches your case, not the first one that looks relevant.

## Falsifiability check on each load-bearing claim

For every claim the plan presents as load-bearing pre-merge proof
("step N validates Z," "fixture X covers Y," "passing `build:web`
confirms Q"), walk through the falsifier in your head: what
observation would prove the claim wrong, and could the named
procedure surface that observation? If the exercise reveals the
validation is ambiguous (multiple causes produce the same
observation, the named procedure cannot distinguish them),
tighten the procedure and record the tightened version with its
discriminator. If the falsifier is obvious and the procedure
clearly catches it, the exercise is its own reward — no recording
needed; most validation bullets fall in this category.

The load-bearing case is exactly when the exercise *changes* the
procedure. Recurring shape: a validation step ("HTTP request
returns 4xx") has multiple failure causes that produce the same
observation — the new route fires correctly returning 404, the
route doesn't exist returning 404, an upstream rule blocks the
request returning 4xx. The named test cannot distinguish the
desired-positive signal from the failures it was meant to catch.
The fix is an identity-fingerprint procedure (capture positive and
negative response signatures and assert against both) or a more
discriminating check.

**What this means at each level.** At epic level, the falsifier targets
capability and constraint claims ("would shipping milestone 2
alone surface this constraint?" "would the named upstream
dependency remove this capability?"). At milestone level, the
falsifier targets claims about cross-task coupling and sequencing
rationale ("can task X actually consume what task Y produces?").
At plan level, the falsifier targets validation procedures and
per-contract claims as the rule body's examples illustrate — the
original framing of this rule.

## Decompose options into shapes before analyzing

When a Choose-One decision (scoping decision, framework / library
choice, alternative-evaluation in a Risk Register) lays out the
candidate set, the first step is **enumeration**, not analysis.
Each named option can hide multiple sub-shapes with materially
different cost/benefit profiles. Before accepting or rejecting any
option, ask: **are there sub-shapes — variants of how this option
could be implemented — that would change the analysis?** Decompose
first, then analyze each shape on its own merits.

The recurring failure mode is category-level analysis: the option's
category name is treated as a single thing, the rejection rationale
holds against one shape that happens to be in the category, and
sibling shapes slip through unevaluated. The decision looks
well-reasoned in retrospect because the rationale holds against the
strawman it cited — the failure is invisible until implementation
surfaces a constraint (or unlocks a benefit) that an unevaluated
sibling shape would have caught.

Categories that frequently hide multiple shapes: "server-mediated
write" hides "stored-procedure call from the client" and "service
wrapping the stored procedure"; "abstraction" hides
extracted-helper, shared-module, and full class hierarchy;
"framework integration" hides minimal adapter,
wrapper-with-escape-hatch, and full rewrite; "caching layer" hides
per-request, per-session, and shared CDN. If a candidate's name
covers more than one viable shape, split it before scoring.

Scoping sessions decompose options into shapes before locking the
candidate set.

## Anti-pattern: planning artifacts that only cite each other

If the plan, scoping doc, and milestone doc all cite each other
for the same load-bearing claim, the claim is unverified. Fluent
cross-doc citation is not verification. Each load-bearing claim
needs at least one citation to actual code, generated test output,
an already-merged sibling artifact, or upstream / vendor
documentation for external-service-behavior claims the codebase
doesn't contain proof of.
