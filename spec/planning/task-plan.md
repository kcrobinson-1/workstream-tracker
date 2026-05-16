# Plan Drafting Sessions

Per-level planning playbook for **plan-drafting** sessions — both
task plans (the doc-type that lands when a task has one phase, or
that orchestrates N ≥ 2 phase plans) and phase plans (the structural
sub-doc of an N ≥ 2 task plan). Loads [`shared.md`](./shared.md) for
cross-level planning rules. This file covers what is unique to the
implementation layer — task plans and phase plans are the doc-types
under which code actually lands, which is why the Planning Depth
rule, the `` `In draft` → `Proposed` `` promotion gate, and the
Plan-to-PR Completion Gate all live here rather than in `shared.md`
(epic and milestone docs do not consume those rules, so they cannot
sit at the cross-level layer without misfiring).

The plan tree carries four conceptual levels (epic → milestone →
task → phase) and three doc-types in practice (epic doc, milestone
doc, task plan). Phase plans exist only as structural sub-docs of a
task plan when the task has N ≥ 2 phases (one file per phase); when
the task has N = 1 phase, the task plan IS the only doc, with phase
content absorbed inline.

**The level picker** is *independent value per unit vs. sequence
steps toward one outcome.* Applied between adjacent levels: each
milestone of an epic has independent stakeholder-facing value
relative to siblings; each task of a milestone has independent
stakeholder-facing value within the milestone's scope; phases of a
task lack independent value — they are sequence steps toward the
task's one outcome, and stopping after phase 1 ships nothing
coherent. The operational question for a new effort: "if I shipped
half of this, would what shipped have independent value, or is
value realized only when the sequence completes?" Independent-value
shape implies the level above; sequence-step shape implies the
level below. The recurring trap is mis-classifying a sequence-step
effort as the level above and importing per-PR cross-phase
contracts that produce review-time corrections one PR at a time;
the level picker plus the Planning Depth rule below catches both
failure modes at plan-drafting time.

## Rules for task plans and phase plans alike

The rules in this section bind every plan-drafting session at the
implementation layer regardless of whether the doc is a task plan or
a phase plan. Three of them — Planning Depth, the `` `In draft` →
`Proposed` `` promotion gate, the Plan-to-PR Completion Gate — bind
at this layer but not at epic or milestone layer, which is why they
live here rather than in [`shared.md`](./shared.md).

### Planning Depth

When asked to make a plan, do not compress the workflow to an
arbitrary short step count.

- include every execution gate that materially affects quality,
  even if that makes the plan longer than five steps
- keep baseline validation, branch hygiene, implementation,
  automated code-review feedback, documentation cleanup, final
  validation, and PR preparation as separate steps when they are
  all relevant
- for implementation plans, include an explicit
  documentation-current-state gate: identify which docs with
  status-oriented sections (for example `Current State`, `Current
  status`, rollout status, or phase status) are affected, and
  include updating them as a required step before handoff
- do not merge steps just to keep the plan visually compact
- if a new required step is added, insert it at the correct point
  in the sequence without weakening or collapsing the surrounding
  steps
- for implementation plans, make the plan decision-complete enough
  that another engineer or agent can execute it without inventing
  missing gates, validation, or handoff work
- for PR-sized work, name the intended commit boundaries before
  editing when practical, and keep review-fix commits distinct
  when they clarify the history
- name the self-review audits that apply to this PR's diff
  surfaces, drawn from the project's self-review catalog (if
  maintained). The plan should list audit names by surface (data /
  frontend / CI / runbook) so the implementer runs them at commit
  boundaries rather than rediscovering review feedback at PR-review
  time

### Just-in-time scoping and plan drafting

A plan-drafting session produces the plan an implementing PR
consumes. Run this session just-in-time before the plan's
implementation starts. Drafting the scoping and plan docs may
begin while the prior phase (for phase plans) or prior task (for
task plans) is still in implementation or review, provided every
still-pending decision in the prior effort is enumerated as a
named "input from prior phase" (or "input from prior task") in
the scoping doc and carried through the plan doc as an open
input the plan must verify before promotion. A pending input is
only valid if it cites the concrete surface where the decision is
being made — a PR number, a scoping-doc section heading, a
review-comment thread, an issue number, or a named cross-phase
decision in the milestone doc. Bare "TBD," "pending," or
unattributed-prose entries do not count and must be resolved
(either by citing the surface or by deciding the question now)
before they can be carried as inputs. The plan's `Status:` stays
`In draft` (not `Proposed`) until every named input has settled —
this is the canonical pre-`Proposed` label per the
`` `In draft` → `Proposed` `` promotion gate below, and that
gate's full self-review walk runs before the flip. Updating the
next plan's draft if the prior effort shifts during implementation
or review is an accepted cost — preferable to forcing serial
execution. Do **not** open plan-drafting in batch alongside the
prior phase's or prior task's planning session: that risks
recording assumptions before any code exists to ground them; the
relaxation here is about parallelism with **implementation or
review**, not with **planning**.

### Goal: scoping doc + plan doc

Produce two artifacts that split ownership cleanly rather than
co-cover the same content: a plan scoping doc (transient —
deletes in batch with sibling scoping docs at the plan's
terminal PR, or at the milestone-terminal PR for phase plans
under a milestone) and a plan doc (durable — survives the
feature). Paths follow the in-repo plan layout (see the "Path
conventions" section below). The next bullet specifies what each
owns; both docs may carry a one-paragraph context summary for
orientation, and that is the only intentional overlap.
**Scoping docs do NOT carry a Status field** — the Status
lifecycle is for plans (which the implementing PR consumes);
scoping docs are transient deliberation that delete at
terminal-state, so the lifecycle doesn't apply to them and a
Status field on a scoping doc can only become wrong (a `Proposed`
scoping doc has no consumer for the promotion claim, and an
`In draft` scoping doc becomes stale the moment the sibling plan
flips to `Proposed`). Sibling plans alone carry Status.

### Scoping owns / plan owns

Because scoping deletes at the plan's (or milestone's) terminal
PR, the durable plan must end up with everything worth persisting
in record after the feature launches; restating the same content
in scoping during the scoping doc's lifetime burns drafting time
and creates drift risk every time one side updates without the
other.

- **Scoping owns** the deliberation prose with rejected
  alternatives (the "Decisions made at scoping time" section, each
  decision carrying `Verified by:` code citations), the open
  decisions to make at plan-drafting (handoff), the plan-structure
  handoff, and the reality-check inputs the plan must verify
  (handoff). This content has no audience after the plan lands —
  exactly why it lives somewhere transient.
- **Plan owns** the durable record that survives after the
  feature lands. The specific sections a task plan or phase plan
  carries are enumerated in "Required and optional sections"
  below.
- **Scoping does not restate plan-owned content.** Where a
  scoping decision touches the file inventory, a contract, an
  invariant, a validation procedure, or a risk, scoping
  references the plan's section by name ("the `EventContent` type
  defined in the plan…"); it does not duplicate the artifact. The
  reality-check gate (named below) operates on scoping's
  decisions — load-bearing claims about the codebase or supporting
  services that reality-check verifies — not on duplicated
  contract text.
- Recurring trap: scoping + plan together can balloon to
  thousands of lines because both docs end up carrying the full
  file inventory, full Contracts, full Cross-Cutting Invariants,
  full Risk Register, and full Self-Review Audits — heavy
  duplication for the same coverage, with drift risk every time
  one side updates without the other. The fix is to follow the
  ownership split above strictly: scoping references the plan's
  section by name rather than duplicating its content.

### Required and optional sections

A task plan or phase plan carries the following sections.

**Required:**

- Status (per the canonical lifecycle in [`shared.md`](./shared.md)
  "Plan-doc Status," with the `Proposed → Landed` transitions
  bound by "Plan-to-PR Completion Gate" below)
- Context preamble (per "Plan opens with a plain-language
  context preamble" below)
- Goal
- Contracts (full final shape)
- Files to touch — new / modify / intentionally not touched
- Validation Gate

**Optional, when content applies:**

- Cross-Cutting Invariants — when ≥ 2 sites must agree on a
  rule (per [`shared.md`](./shared.md) "Cross-Cutting
  Invariants section")
- Phase Contracts — required when the task plan is N ≥ 2
  (the per-phase WHAT contract; phase HOW stays in each phase
  plan), per [`shared.md`](./shared.md) "Parent-doc child
  contracts". Absorbed inline; an N = 1 task plan has no phase
  children and omits it.
- Naming — when the plan introduces new identifiers
- Execution Steps — when implementer ordering beyond Commit
  Boundaries is needed
- Commit Boundaries — when the implementing PR is multi-commit
- Self-Review Audits — when diff surfaces map to entries in the
  project's self-review catalog
- Documentation currency in implementing PRs — when
  status-bearing docs are touched
- Out Of Scope — when the plan records boundary calls as final
  answers (deliberation prose lives in the scoping doc)
- Risk Register — when residual risks exist
- Backlog Impact — when backlog entries graduate, get
  deleted, split, or shift (per [`backlog.md`](../backlog.md))
- Related Docs

Variance from this list — an unlisted section is appropriate, or
a required section genuinely doesn't apply (decision plans whose
deliverable IS the recorded decisions may legitimately skip
Contracts or Files to touch, for example) — follows the
[`shared.md`](./shared.md) "Section variance disclosure" rule.

### Scoping precedes plan drafting; check before starting plan draft

Before opening the plan doc to write, verify the scoping doc
exists at this plan's canonical scoping path per the in-repo
plan layout (see "Path conventions" below) with substantive
scoping-owned content per the rule above — at minimum, a
"Decisions made at scoping time" section with at least one
decision carrying a `Verified by:` code citation, plus whichever
of "Open decisions to make at plan-drafting," "Plan structure
handoff," and "Reality-check inputs" the plan needs. Not empty,
not a stub, not a placeholder paragraph saying "scoping
pending." If the scoping doc does not exist or is a stub, do
scoping first as its own artifact; plan-drafting cannot start
without it. Without scoping content, the reality-check gate
below has nothing to operate on, and plan-drafting silently
collapses into scoping-during-drafting — exactly what scoping
exists to separate from drafting. The substantive-content list
named here is the falsifier for that gate; it intentionally
tracks scoping's owned content per "Scoping owns / plan owns"
above, not the plan-owned content (file inventory, contracts,
validation surface, risks) that earlier drafts of this rule
asked scoping to also carry. The check is a simple
file-existence + first-paragraph read, takes seconds, and
protects against the most common procedural skip when plan
drafting starts in a fresh agent session that did not produce
the scoping doc.

### Doc-only decision plans satisfy the substantive-content gate via cited open-question constraints, not resolved decisions

When the whole plan's output is the decision artifact — no code
ships, no contracts get implemented; the durable plan doc IS the
recorded decisions plus rejected alternatives — the scoping doc
may surface the decision space without resolving any decisions
at scoping-doc-open time. Decisions resolve through the
collaborative deliberation that constitutes the plan, then
absorb back into the durable artifact. The protective intent of
the rule above (prevent stub-scoping → drafting collapse) holds,
but is satisfied differently: the scoping doc must carry
code-grounded `Verified by:` citations on the constraints that
bound each open question, and on the reality-check inputs the
resolution will rest on. Citation-free constraints, generic
"TBD" placeholders, or open-question sections that don't
decompose the decision space against actual code still fail the
gate — the carve-out is for decision-deferral, not for skipping
the code grounding. Code-shipping plans do not get this
carve-out: the decision-resolved-at-scoping-time bar still
applies because the plan-drafting step that scoping precedes is
real for them.

### Narrow-surface plans may skip the scoping doc

A plan whose estimated implementation surface is bounded enough
that options-considered analysis and reality-check input
gathering would not produce material content may go straight to
drafting the plan doc. A plan qualifies as **narrow-surface**
when ALL of these hold:

1. **Single subsystem.** Touches one of: a UI surface (route,
   section, component), a data-model layer (table, migration),
   or backend logic (RPC, edge function, service handler). Plans
   spanning more than one of these classes are multi-subsystem
   and do not qualify.
2. **Bounded file count.** Plan estimates ≤ 8 files touched
   (excluding generated types and test files).
3. **No new public-API contract.** No new RPC, no new auth /
   authz boundary, no new route family. Schema-touch is allowed
   only when **purely additive** — net-new tables, net-new
   columns on existing tables, net-new indexes — without
   altering the meaning of existing rows or breaking existing
   read paths. Any change that modifies or removes an existing
   structure (column drop, type / default / nullability change,
   rename, FK refactor, access-control rewrite) requires the
   full path. The discriminator is "modifies or removes existing
   structures," not the verb.
4. **No new cross-cutting invariant.** The change does not
   introduce a rule that multiple files must agree on.
   Cross-cutting invariants are the load-bearing reason scoping
   docs exist; their absence is the load-bearing reason to skip.
5. **No novel mechanism.** Uses patterns already established in
   the codebase. Novel mechanisms (a new auth shape, a new
   framework idiom, a new privilege-elevation pattern) trigger
   the "Spike before plan for novel mechanisms" rule below and
   warrant scoping.

All five must hold. If any fails, draft the scoping doc as the
"Scoping precedes plan drafting" rule above requires. This
carve-out is an exception the planner explicitly invokes; the
default direction stays "scoping first."

**Verification protocols are not optional under this carve-out.**
Narrow-surface plans must still carry the reality-check inputs
the implementing pass walks (per the "Reality-check pass before
plan-drafting" rule below). The plan doc absorbs them inline
when the scoping doc is skipped — the *form* compresses (a
Reality-check inputs section in the plan rather than a separate
scoping artifact), the *function* (falsifier protocol against
load-bearing claims) does not. A narrow-surface plan whose doc
has no Reality-check inputs section has skipped the carve-out's
protective intent, not just its prose.

### Plan opens with a plain-language context preamble

Before any implementation specifics (file paths, framework
names, function signatures, phase-numbering shorthand), the plan
must contain 1–3 paragraphs that name three things in plain
English: **what this plan covers** (the surface or capability
under change, not the file paths), **why it's being done now**
in human terms (closes a loose end from a prior phase, lays
foundation for upcoming work, addresses a specific user-visible
gap), and **what surfaces this touches** at the conceptual level
(admin pages, routing layer, e2e fixtures, docs — not file
paths). Phase-numbering or task-numbering prose ("depends on 2.2
+ 2.3," "prerequisite for 2.5") describes the dependency graph,
not the motivation, and does not satisfy "why now." The preamble
can live at the top of `## Goal` or in a separate `## Context`
section before Goal; structure is implementer choice. The
protective check this rule enforces: a reader who hasn't read
the epic, milestone, or scoping doc can understand what problem
this plan solves and why anyone should care after reading the
plan's first ~250 words. Implementation-detail-first openings
are the recurring trap that motivates this rule — opening with a
mechanical migration or refactor description ("Migrate X from A
to B as the new platform surface") is true, complete,
mechanical, and silent on why anyone other than the plan author
should care.

### Reality-check pass before plan-drafting

Before promoting the scoping doc to plan-drafting, do a forced
reality-check pass on every load-bearing claim about the
codebase or supporting services.

- **Data-layer contracts.** Read the actual migration files for
  named tables, policies, RPCs; confirm the predicates, grants,
  and constraints exist as scoped.
- **RPC / service-handler behavior.** Read the function body —
  if the scoping says "widen the X gate," confirm the gate
  exists.
- **Database engine semantics.** Write one sentence that would
  falsify the claim and check whether it falsifies. Engine
  behavior under concurrent writes, transaction isolation, and
  permission evaluation is full of footguns; a one-sentence
  falsifier exposes most of them.
- **Type-system / signature contracts.** Read the function
  signature and at least one real call site.
- **URL contracts and route topology.** Confirm any URL the
  plan names is a distinct route the platform actually serves
  (not an inline-conditional render at the same URL); when an
  affordance points at "where the user does X today," verify
  that X has a URL distinct from the affordance's render site.
- **Validation-procedure claims.** For "step N validates Z" /
  "the existing fixture covers Y" / "passing `<command>`
  confirms Q," trace whether the procedure actually exercises
  the surface it claims.
- **Dev-tool semantics.** Read the project's actual config
  files before claiming runtime behavior — dev tools are
  config-dependent and general knowledge will not catch
  project-specific overrides.
- **External-service / framework / vendor behavior.** Read the
  upstream / vendor documentation for any claim about hosted
  dependency behavior the codebase doesn't contain proof of.
  "I think X does Y" is not a reality-check; "the docs at <URL>
  say Y" is.

If the reality-check finds a discrepancy, fix the scoping
before drafting the plan; do not carry wrong premises into plan
time. The recurring shape that motivates the dev-tool bullet
specifically: a config file's destinations or rules can produce
runtime behavior that contradicts a scoping claim ("`<dev-tool>`
will validate the new local routes" turns out wrong because the
config proxies to a deployed environment instead of the local
build) — checking the config catches this at scoping time
rather than during implementation.

### Prefer existing wrapper scripts over lower-level CLI invocations in plan validation steps

Before naming a validation command in a plan, search the
project's script registry (`package.json` `scripts`,
`Makefile`, `scripts/` directory, `justfile`, etc.) for an
existing wrapper. If a wrapper exists, name it — the wrapper is
what future contributors will run, and naming the lower-level
invocation silently skips meaningful orchestration the wrapper
does (local service stack, DB reset, function runtime, env-var
sourcing, fixture seeding). The lower-level command usually
still works, but it forces the implementer to reinvent setup
the wrapper already handles, which is a parallel-track
procedure rather than the project's canonical local path. This
rule is distinct from the reality-check gate above: that one
asks "does the named procedure exercise the right surface;"
this one asks "is the named procedure the canonical entry
point, or am I reinventing orchestration the project already
wrapped."

### Spike before plan for novel mechanisms

When the plan introduces a new mechanism (a new authorization
shape, a new cross-app boundary, a new privilege-elevation
pattern, a new framework idiom), build a 30-minute throwaway
spike that exercises the mechanism end-to-end before writing
the plan. The spike's job is to find dealbreakers — wrong
assumptions about runtime semantics, missing constraints,
hidden coupling.

**Worktree handling for spikes:** create a throwaway branch
named `spike/<phase-or-mechanism>` off the planning branch;
commit freely on the spike branch; do not merge it. When the
spike concludes, either delete the branch or leave it dangling
for reference and continue plan-drafting on the original branch
with a clean worktree. Spike code is never promoted into the
implementation PR — the plan describes the contract, the
implementation PR builds it from scratch. If a scratch script
or non-code artifact would help, write it under a git-ignored
scratch path the project conventions designate.

### Cap

Stop when continued drafting stops changing the contract —
repeated rewrites of the same section, decisions that re-open
after being marked resolved, or new content spawning without
resolving existing questions. That iteration signal is the
diminishing-returns indicator. If drafting feels like it's not
converging, reality-check the actual scope size and either
ship what's clearly right or escalate. Planning time should be
a fraction of implementation time, not a parallel effort.

### `In draft` → `Proposed` promotion gate

This gate binds task and phase plans only. Epic and milestone
docs have a symmetric `` `In draft` → `Proposed` `` gate that
lives in [`shared.md`](./shared.md) "Parent-doc `In draft` →
`Proposed` promotion gate" — it is placed there, not duplicated
here, because epic and milestone docs do not load this file.

Plan-drafting is not required to be a single pass. A multi-pass
session can lay out the plan's structure with explicit
deferrals, then resolve them, then flip Status. While the plan
is in the resolution-pending phase, it carries Status `In
draft`; flipping to `Proposed` claims the plan is ready for code
review — contracts are decision-complete, the doc is internally
coherent, and load-bearing claims are verified. The flip is the
natural moment for a comprehensive self-review pass:

- **Read end-to-end as a coherent whole.** Re-read the full plan
  + scoping doc in order. Look for contradictions between
  sections (a contract claim that conflicts with a Risk Register
  mitigation, a decision in the scoping doc whose framing no
  longer matches a later decision's resolution, a Files To Touch
  entry that contradicts a contract, a multi-step contract
  whose ordering doesn't match its own test contract).
- **Decision-completeness on Contracts.** Walk Contracts for
  deferral phrases ("plan-drafting picks," "final spelling at
  plan time," "shape decided later," or any deferral that names
  plan-drafting itself as the resolver). Each must be resolved
  concretely, explicitly authorized by another rule (e.g.,
  "Bans on surface require rendering the consequence" authorizes
  UX-copy deferral to render-time), or moved to the scoping
  doc's "Open decisions to make at plan-drafting" handoff.
  Plan-drafting is the moment that produced the plan; a contract
  that defers to "plan-drafting" defers to a moment that does
  not exist.
- **Walk the broadened `Verified by:` rule against every
  load-bearing claim.** The rule is named in
  [`shared.md`](./shared.md) "`Verified by:` annotations on
  load-bearing claims"; the promotion gate is when it gets
  applied universally rather than to whichever claims happened
  to feel "technical" during drafting.
- **Re-confirm the scoping doc's reality-check inputs against
  current code.** Line numbers, file existence, configuration
  shapes, and other inputs that drift from scoping →
  plan-drafting are reflected in the plan; stale references are
  updated.

Failures surface either as resolutions (apply edits before
flipping) or as plan-blockers that the user must triage before
the flip happens. A plan flipped to `Proposed` without this
walk is the same shape of drift as a plan flipped to `Landed`
without satisfying its Validation Gate — the Status claim is
wrong.

### Plan-to-PR Completion Gate

The gate that an implementing PR walks before merge, and the
Status lifecycle that gate produces. The canonical lifecycle
values are defined in [`shared.md`](./shared.md) "Plan-doc
Status"; this gate binds the transitions out of `Proposed`.

A PR that implements a plan must leave the plan in a terminal
state. "Most of the plan" is not "the plan." A plan doc that
still says `Proposed` or `In progress` after its implementation
merges is drift, and drift compounds into follow-up PRs that
re-review the same decisions.

- before opening the PR, walk every Goal, Test, Validation
  step, and Self-Review audit named in the plan; for each one
  confirm it is either satisfied in the PR or explicitly
  deferred **in the plan itself** with written rationale.
  Deferrals live in the plan, not in the PR body, not in an
  issue, not as an unwritten promise.
- flip the plan's Status line from `Proposed` / `In progress`
  to `Landed` (or to `Validating`, when post-merge validation
  is outstanding — see exception below) in the same PR that
  implements it. Plans in active multi-pass drafting may carry
  an interim `In draft` Status before `Proposed`; the
  `In draft` → `Proposed` flip is gated by the promotion-gate
  rule above. Do not record commit SHAs in the Status block —
  `git log` and `git blame` are authoritative for navigating
  from plan to history, and recording SHAs creates a chicken-
  and-egg problem (the SHA isn't known until after merge,
  which forces a follow-up commit whose only purpose is to
  record the previous commit's SHA). Same-PR flip to `Landed`
  is the default whenever the plan's Validation Gate can be
  fully satisfied pre-merge AND a single implementing PR (or
  a clearly-last-to-merge implementing PR) exists. Two named
  exceptions allow a follow-up close-out PR:

  - **Post-merge validation.** Plans whose Validation Gate
    names a check that can only run post-merge (production
    smoke, manual walkthrough, async user-facing
    verification) land in two phases: the implementing PR
    merges with Status `Validating` (optionally
    `Validating — <gate-descriptor>` for context, e.g.
    `Validating — prod smoke`); a follow-up doc-only commit
    flips Status to `Landed` and records the validation run
    URL or evidence once the post-merge check passes. The
    run URL is durable external evidence, unlike a commit
    SHA which is already in git.
  - **Parallel implementing PRs.** When a plan ships via N≥2
    implementing PRs that can land in any order (neither PR
    is clearly last), the close-out (Status flip plus any
    scoping-doc deletion or backlog-cleanup the plan binds
    to close-out) lands in a single-commit follow-up PR
    opened immediately after the last implementing PR
    merges. Plan Status moves `Proposed → In progress` when
    the first implementing PR merges and `In progress →
    Landed` when the close-out PR merges (or `In progress →
    Validating → Landed` if post-merge validation also
    applies). The drift window between the last implementing
    PR's merge and the close-out PR's opening is
    bookkeeping-bounded — on the order of minutes, not days
    — because the close-out PR carries no decision content.
    The plan's Contracts section invokes this exception by
    name; an N≥2-PR plan that does not invoke it defaults to
    "last-to-merge carries the close-out" (workable only
    when the order is knowable in advance).

  These two exceptions are the complete list; do not invent
  additional states or leave the flip to an informal
  post-merge promise. A plan-specific carve-out that adds a
  third exception belongs in this canonical rule, not in the
  plan's own Contracts section.
- ban soft-commitment words in plans: "optional but
  recommended," "consider adding," "nice to have," "probably
  should." A requirement is either in-scope or deferred —
  there is no third option. Soft commitments silently relax
  under review pressure and reappear as reviewer findings
  after merge.
- if a reviewer flags a gap that should have been named at
  plan time, fix the plan first (tighten the requirement or
  defer with rationale), then address the gap. Do not carry
  the gap as a post-merge follow-up without updating the
  plan.
- if a plan requirement cannot be fully satisfied in the
  intended PR, split the plan along a phase boundary before
  merging partial work so each phase's Status can flip
  independently, rather than merging a partially-satisfied
  plan and tracking the remainder informally.
- **Call out estimate deviations in the PR body, and update
  the plan to match what shipped.** When implementation
  diverges from an estimate-shaped section of the plan
  ("Files intentionally not touched" ended up touched, "Files
  to touch — new" missed a file, contract bullets gained a
  requirement, intended commit boundaries reshuffled, an
  execution step was unnecessary or had to be split), two
  things must happen in the same PR:
  - The **PR body** names the deviation explicitly under a
    `## Estimate Deviations` heading inserted immediately
    after `## Documentation` (or `N/A` if no deviations).
    Each entry is one or two sentences naming the
    estimative section, the actual outcome, and why the call
    was the right one — enough that a reviewer can audit the
    deviation without reading the diff cold. This is the
    rationale and audit trail.
  - The **plan doc** is updated so its estimate-shaped
    sections describe what actually shipped, not the
    pre-implementation guess. Walk every estimate-shaped
    section ("Files to touch — new / modify / intentionally
    not touched," per-module Contracts, Execution steps,
    Commit boundaries) and reconcile each against the merged
    diff. A plan that says "Files intentionally not touched:
    X" after we shipped an X edit is the same shape of drift
    the Status-flip rule already forbids — the plan must
    describe the implemented system, not the
    pre-implementation guess. The PR body says *why we
    deviated*; the plan says *what shipped*.

  Distinct from the rule-deviation path above: rule
  deviations (a Cross-Cutting Invariant turning out to be
  wrong, a contract that can't be satisfied, a Validation
  Gate command that doesn't exercise what it claims) require
  the plan rule itself to be rewritten in the same PR;
  estimate deviations require the estimate-shaped section to
  be updated to match reality, plus the PR-body callout.
  Pre-existing PR templates do not need the Estimate
  Deviations heading until they are next edited; PRs opened
  from this point forward must include the section, and
  plan-implementing PRs must reconcile the plan with what
  shipped per the bullet above.

## Task-plan-specific rules

These rules bind task plans (a doc with N = 1 phase absorbed
inline, or the orchestrating doc when N ≥ 2 phase plan files
exist). They do not bind phase plans, where the per-phase
contracts and per-PR coordination are scoped narrower.

### Cross-PR coordination

A task plan with N ≥ 2 phase plan files coordinates between
phases via the task plan itself (cross-phase decisions, shared
invariants, sequencing rationale), not through pre-locked
cross-phase contracts in each phase plan. When a phase plan
needs something from a sibling phase that hasn't shipped yet,
write down the assumption in the consuming phase plan and tag
it for verification at sibling-merge time. Don't wait for the
sibling phase to catch up; don't pre-coordinate every detail.
Recording assumptions to verify-on-merge is more honest than
committing to a contract neither side has built yet.

### Task plan terminal state when N ≥ 2

A task plan with N ≥ 2 phase plan files reaches `Landed`
Status when the last phase's implementing PR merges. Earlier
phase plans flip `Landed` individually as their implementing
PRs merge per the Plan-to-PR Completion Gate above; the task
plan's own Status flips with the last phase. The terminal-state
flip of the task plan happens in the same PR as the last
phase's Status flip — typically that PR's diff also closes out
the task's tracking surfaces (parent epic's milestone row when
the task lives under one, related backlog entries).

### N = 1 task plan: phase content absorbed inline

When a task has exactly one phase, the task plan IS the only
doc; no separate phase plan file is created. The task plan
absorbs phase content inline (Contracts, Files to touch,
Validation Gate, Self-Review Audits, Execution Steps). Treat
the entire "task and phase plans alike" rule set as binding on
this single doc; the Phase-plan-specific rules section below
applies to the inlined phase content (PR-count branch test,
bans-on-surface, etc.). The doc carries Status as a task plan,
not as a phase plan.

## Phase-plan-specific rules

These rules bind phase plans (the per-phase file inside an
N ≥ 2 task plan) and the inlined phase content of an N = 1
task plan. They do not bind task plans at the orchestration
layer (an N ≥ 2 task plan's Status flip is governed by the
task-plan-specific rule above, not by the PR-count branch test
below).

### PR-count predictions need a branch test

Before declaring "1 PR" in the plan's Status block, create the
branch and sketch the file list. If the diff would touch >5
distinct subsystems or >300 LOC of substantive logic, split.
Either ship as sub-phases (within the same phase doc, not as a
new plan-tree node — the four-segment slug cap precludes a
fifth node-level) or justify the size with concrete
review-coherence reasoning. The milestone doc's PR-count
estimate does not bind the phase plan.

### Bans on surface require rendering the consequence

When a plan writes "no X" / "minimum surface" / "intentionally
not done" for any observable surface (UX, log output,
operational dashboard, generated artifact), state in concrete
terms what the absence looks like and verify it. The
discipline is not "always add the thing" but "before banning
the surface, prove the no-X outcome is acceptable by observing
it." Optimizing for diff size produces plans that ship
regressed surfaces.

Verify the consequence in the environment that will actually
run it, not the local-dev environment that may diverge.
Local-dev tooling routinely papers over production-only
behavior (dev servers self-serve asset paths a production
proxy would miss, in-process runtimes fake out boundary
crossings the deployed environment enforces, dev databases
skip migrations a production database applies). The recurring
shape: a config-driven change works in local dev because dev
synthesizes the missing piece; production exposes the gap
because no synthesis happens. When the consequence depends on
deployed-environment behavior, run the check against a
production-shaped build, not the dev shortcut.

### When a behavior change shifts which code path a test exercises, re-audit every assertion — not just the changed input

A test's expectations on the runtime path it exercises can
break in non-obvious ways when an "input-only" change actually
shifts which code runs. A test that asserts on a specific
output (a button label, a log line, an emitted event) against
the old path may reach a different code path that produces
different output even though the test's surface inputs are
unchanged. Stability invariants on the *new* path cover what
that path must preserve; the plan must additionally walk the
test against the *new* code reached at every step the test
takes, not just the entrypoint. Cite the target code surface
for each step and verify the assertions resolve against the
actual behavior of that surface. The recurring shape: a plan
retargets a navigation/routing/dispatch affordance, the
destination invokes a different code path than the source,
and assertions the test "obviously preserves" actually break
because the new path's outputs don't carry them.

### Cross-boundary destinations need the boundary-crossing primitive, not the in-context shortcut

When a destination requires the host's routing or dispatch
layer to fire (a CDN rewrite to evaluate, a service router to
re-invoke, a cross-app proxy to rebuild context, a cache
freshness check to run), use the primitive that exits the
in-context layer and re-enters the host. In-context
equivalents (SPA-internal navigation APIs like `replace` /
`push`, framework link components, History-API methods like
`pushState` / `replaceState`; in-process function calls when
the contract requires going through a service dispatcher;
cache reads when the contract requires freshness validation)
produce the right surface result but the host's layer never
re-evaluates. Cross-boundary destinations need the primitive
that leaves the in-context layer: for SPAs, full document
navigation (`window.location.replace` / `assign`); for
in-process calls, the cross-process invocation; for cached
reads, an explicit cache-bypass.

The same trap shows up in reverse when a destination migrates
*out* of the in-context layer: existing in-context tooling
that still produces the migrated request (button handlers,
in-process callers, cached-read sites) must be audited and
converted to the boundary-crossing primitive; the request is
generated correctly but the in-context layer never escalates,
so the host never re-evaluates.

When a plan specifies any in-context API as a contract, walk
every destination and classify in-context vs. cross-boundary.

## Task plan / phase plan relationship

### N = 1 → N ≥ 2 transition

When a task started as N = 1 and a second PR enters scope
mid-flight, the task plan gains a `Phases` section describing
per-phase scope. Whether the original PR's contracts get edited
in-place vs. split into separate phase plan files is a call
made at the moment N grows, not pre-locked here. The N ≥ 2
split is warranted when the additional phase needs the full
apparatus (separate Contracts surface, separate Validation
Gate, separate Self-Review Audits); if the second phase fits
cleanly under the task plan's existing apparatus, keep one doc
and add the `Phases` section. The supersession-vs-edit-in-place
call follows from whether contracts in the original PR's plan
are still accurate after the second phase enters: if yes, edit
in place and add the new phase's content; if no, supersede the
original content with an explicit "see PR #N for the original
phase 1 shape" pointer and let the diff carry the rest.

### Path conventions

Per the in-repo plan layout (see
[`planning-doc-location.md`](../planning-doc-location.md)):

Every plan-tree root lives in its own folder at
`docs/plans/<root-slug>/`. The root doc is at
`<root-slug>/README.md`; descendants are sibling files inside
the same folder, named by the slug-suffix (the part of the
slug after the root) plus an optional `-<descriptor>` for
human readability. The slug in frontmatter is the identity;
the filename is for browsing.

- **Standalone task plan, N = 1 phase:**
  `docs/plans/<task-slug>/README.md` only. Phase content lives
  inline in the task plan.
- **Standalone task plan, N ≥ 2 phases:**
  `docs/plans/<task-slug>/README.md` (the task plan) plus
  `docs/plans/<task-slug>/p<N>-<descriptor>.md` for each
  phase plan.
- **Epic with milestones, tasks, phases:**
  `docs/plans/<epic-slug>/README.md` (the epic doc),
  `docs/plans/<epic-slug>/m<N>-<descriptor>.md` (milestone
  docs), `docs/plans/<epic-slug>/m<N>-t<T>-<descriptor>.md`
  (task plans), `docs/plans/<epic-slug>/m<N>-t<T>-p<P>-<descriptor>.md`
  (phase plans).
- **Scoping docs:** `docs/plans/<root-slug>/scoping/<slug-suffix>.md`
  mirroring the slug position of the doc they pair with. The
  `scoping/` subfolder is transient — its contents delete in
  batch at the milestone-terminal PR (or the task-terminal PR
  for standalone task plans) per "Scoping owns / plan owns"
  above.

The four-segment slug cap (root + m + t + p) holds across all
shapes; see [`shared.md`](./shared.md) "Plan-doc identity
(slug)" for the rule.

### How a phase plan cites its parent task plan

When a task plan has N ≥ 2 phase plan files, each phase plan
opens with a `## Context` section that names the parent task
plan and links to it. Rule references that the phase plan
inherits from the task plan (Cross-Cutting Invariants, Naming,
cross-phase decisions) cite the parent task plan's section by
name rather than duplicating the content; the scoping-vs-plan
duplication rule applies recursively across task/phase
boundaries — duplicating the task plan's Cross-Cutting
Invariants block into each phase plan creates the same drift
trap that duplicating shared invariants between scoping and
plan creates.

### Noun discipline

In rule prose under the planning spec and adjacent rule-bearing
files, prefer compound forms ("task plan," "phase plan,"
"task-level," "phase-level") over bare "plan" or bare "task"
wherever ambiguity could arise. Bare "task" is overloaded
against the colloquial / to-do-list / issue-tracker term; bare
"plan" is overloaded against the broader category of "any doc
under the plan tree." Both compounds disambiguate the doc-type
in prose. The directory name `docs/plans/` is a category label
and stays unchanged; the noun discipline binds prose, not
paths. Reviewers may flag bare-"plan" usage in rule prose
under the planning spec as a structural issue when the
surrounding sentence does not clearly anchor the term.

Bare "plan" remains correct when unambiguously generic
(cross-level rules in `shared.md`, generic plan-implementation
workflow, "plan tree" as a category reference). Where the bare
term creates ambiguity — typically when a rule prescribes
behavior at one doc-type level but phrases it as "the plan" —
use the compound form. The
recurring trap that motivates this discipline is the
mis-classification path: a sequence-step effort labeled "an
epic" carries per-PR contracts that should have been
task-plan content, surfacing as one-finding-at-a-time review
corrections. Noun discipline plus the level picker (above)
catches both failure modes at plan-drafting time.
