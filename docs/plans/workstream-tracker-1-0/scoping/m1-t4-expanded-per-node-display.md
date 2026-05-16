# Scoping — t4 Expanded per-node display

Transient deliberation doc for `workstream-tracker-1-0-m1-t4`.
Deletes in batch with sibling scoping docs at the m1-terminal PR
per [`task-plan.md`](../../../../spec/planning/task-plan.md)
"Scoping owns / plan owns." Carries no Status field (scoping
docs are transient; the sibling plans alone carry Status).

Pairs with the task plan at
[`../m1-t4-expanded-per-node-display.md`](../m1-t4-expanded-per-node-display.md)
(N ≥ 2 orchestrating doc) and its phase plans.

## Context summary

t4 is the last task on the read-experience track of m1
([`m1-v0-2.md`](../m1-v0-2.md)) and the terminal task of the
milestone. t3 (Landed, #7) already lifted node labels to
`<Type> <ordinal>: <short_description>` and — as an interface
for t4 — already parses the markdown body into
`parsedDoc.LongDescription` and carries it onto
`PlanNode.LongDescription`; that value is **parsed but never
rendered today**. t4 makes each node surface detail beyond the
Status badge: the long description, and the PRs related to the
node. Related PRs come from two sources the milestone t4
contract names: an optional `related_prs` frontmatter field, and
auto-discovery via `gh pr list` keyed on the slug.

t4 forks from t3 only (read-experience track). It is independent
of the t1/t2 foundation track (multi-WI / auto-registration),
which it does not touch — so there is no pending "input from
prior task" entry from t1 or t2 to carry. The single input from
a prior task is t3's `LongDescription` carry, which is **Landed**
(not pending): re-confirmed against current code in the
reality-check pass below.

## Reality-check pass (load-bearing claims, verified)

Each claim was checked against the code at scoping time; the
plans re-verify them at draft time (see "Reality-check inputs
the plans must verify").

- **t3's long-description carry is live and unrendered.**
  `parsePlanDoc` (`internal/site/walker.go:104-129`) sets
  `parsedDoc.LongDescription = markdownBody(source)`;
  `buildTree` (`internal/site/tree.go:122-131`) copies it to
  `PlanNode.LongDescription`. The `{{define "node"}}` block
  (`internal/site/render.go:54-63`) emits only the badge, the
  `.Label` span (with `title="{{.Slug}}"`), `.WorkInstances`,
  and `.Children`. `LongDescription` has **no emission site** —
  P1 is purely additive to the template, joining an existing
  field. **This is the load-bearing finding behind D2 and the
  P1 contract.**
- **`related_prs` reads via the same tolerant frontmatter
  pattern as `short_description`.** `parsePlanDoc`
  (`internal/site/walker.go:115-128`) already does
  `metaData["slug"].(string)` and
  `metaData["short_description"].(string)` with `, _ :=`
  absence-tolerance. A list-valued field surfaces from
  goldmark-meta as `[]interface{}`, not `[]string`, so the
  read needs an element-wise string assertion — a shape
  difference from the scalar fields, called out so the P1
  contract specifies it rather than assuming `.([]string)`
  works. **Assumption (tagged):** goldmark-meta decodes a YAML
  block sequence to `[]interface{}` of `string`; the P1 plan's
  reality-check re-confirms this against the `goldmark-meta`
  version pinned in `go.mod` before promotion.
- **The render path is a single static `html/template` with no
  JavaScript and no client state.** `render.go:14-64` is one
  `template.Must` with an inline `<style>` block and no
  `<script>`. `design/v0.1-design.md` §7 states "No
  collapse-and-expand affordances ... Everything is rendered on
  one page, refresh-to-update." A detail-on-click panel would
  be the first JavaScript and the first non-`/` route in the
  codebase. **Load-bearing for D1.**
- **No subprocess execution exists anywhere in the codebase.**
  `grep -rn "os/exec" --include="*.go"` is empty; `go.mod`
  declares no `gh`-related dependency. The `gh pr list`
  auto-discovery (P2) is the **first** subprocess shell-out in
  the codebase — a novel mechanism per
  [`task-plan.md`](../../../../spec/planning/task-plan.md)
  "Spike before plan for novel mechanisms." **Load-bearing for
  D4 and D5.**
- **The render path is walk-on-every-request and m1 forbids
  introducing caching.** `Server.index`
  (`internal/site/site.go:43-66`) calls `walkPlans` then
  `buildTree` per HTTP handler invocation.
  [`m1-v0-2.md`](../m1-v0-2.md) Cross-Task Invariant "Render
  path stays walk-on-every-request" states "No task in m1
  introduces caching, file watching, or in-memory build-up;
  that optimization is deferred until a real performance pain
  point surfaces." Any per-request `gh` call therefore cannot
  be memoised or background-refreshed within m1. **Load-bearing
  for D4.**
- **The tool is single-contributor, single local environment.**
  [`README.md`](../README.md) cross-cutting invariant "Single
  contributor, single local environment. No auth, no
  multi-tenant, no remote hosting." A per-request local `gh`
  invocation is evaluated against a single local user loading a
  page, not a hosted multi-user request rate. **Load-bearing
  for D4's "per-request is acceptable" call.**
- **`spec/planning/shared.md` has an optional-field block t3
  established.** `shared.md:135-154` documents
  `short_description` as "optional and additive." `related_prs`
  documents adjacent to it with the identical posture (D3),
  mirroring t3 scoping D2.
- **The canonical validation gate is the Go toolchain; no
  wrapper script exists.** [`docs/dev.md`](../../../dev.md)
  lines 84-86 name Build (`go build ./...`), Vet
  (`go vet ./...`), Tests (`go test ./...`). Repo root has no
  `Makefile`/`justfile`; `scripts/` holds only `assemble.sh`
  (agent-rules vendoring, unrelated to build/test). Each phase
  plan names these plus its own manual checks.

## Decisions made at scoping time

Each decision carries a `Verified by:` citation. Rejected
alternatives live here (no audience after the plans land); the
durable contracts live in the plan docs and are not restated.

### D1 — Surface shape: inline static rendering, no JavaScript panel

This resolves the milestone's deferred "Long-description
rendering location (t4)" decision
([`m1-v0-2.md`](../m1-v0-2.md) Cross-Task Decisions). The long
description and related PRs render **inline beneath the node**
as static HTML emitted by the existing recursive
`{{define "node"}}` template — no detail-on-click panel, no
expand/collapse interaction.

- **Rejected: detail-on-click panel.** A panel requires
  client-side JavaScript (or a new per-node route) to show/hide
  detail. That would be the first JavaScript and the first
  non-`/` route in the codebase — a novel mechanism well beyond
  the bare-bones v0.1/v0.2 render ethos, and disproportionate
  to surfacing two text fields and a short link list.
- **Rejected: collapse/expand "tree-inline expansion."** The
  milestone's phrasing offered "tree-inline expansion" as an
  option, but interactive expand/collapse also needs JavaScript
  (`<details>` is the only no-JS option and renders a native
  disclosure widget that diverges from the bare-bones inline
  style without buying a real interaction model for a
  single-user local page). Static always-rendered inline detail
  is the minimum surface that satisfies the contract; if detail
  density becomes a real scanning problem, an interaction model
  is a later-milestone call, not a t4 call.
- **Rejected: both.** Strictly more surface than either, with
  no consumer for the redundancy on a single-user local tool.
- `Verified by:`
  [`render.go`](../../../../internal/site/render.go) is one
  `html/template` with an inline `<style>` and no `<script>`;
  [`design/v0.1-design.md`](../../../../design/v0.1-design.md)
  §7 ("No collapse-and-expand affordances ... rendered on one
  page, refresh-to-update"); the `node` template already
  recurses over `.Children`, so an inline detail block composes
  into it without new routing.

**Bans-on-surface consequence (per
[`task-plan.md`](../../../../spec/planning/task-plan.md) "Bans
on surface require rendering the consequence"):** "no
interaction model" means a node with a long multi-paragraph
description renders that whole description inline every page
load, lengthening the page. This is acceptable for the current
corpus (plan-tree docs have short bodies; the milestone Out of
Scope defers density/ordering work) and is the observed
consequence the P1 plan's manual validation step must actually
render and look at, not assume.

### D2 — N ≥ 2: P1 inline detail render, P2 `gh` auto-discovery

t4 ships as an N ≥ 2 task plan: a P1 phase (render
`LongDescription` + parse/render the optional `related_prs`
frontmatter field — pure read-path, no external dependency) and
a P2 phase (`gh pr list` auto-discovery merged into the related
-PR set). The task plan is the orchestrating doc; each phase has
its own phase plan file.

- **Rejected: N = 1, all content inline (the t3 precedent).**
  t3 bundled parse+render as N = 1 because both halves were one
  pure-read-path PR's worth of review surface with no
  intermediate independently-shippable artifact. t4 differs on
  the load-bearing axis: P2 introduces the codebase's first
  subprocess shell-out (reality-check), which carries a
  *distinct* Validation Gate (the failure matrix — no `gh`
  binary, no auth, no network, non-repo directory, non-zero
  exit, malformed JSON, hang — must be exercised against
  environments, not just unit-tested) and *distinct*
  Self-Review surface (subprocess error-surfacing,
  process/timeout cleanup) that the pure-template P1 does not
  have. Per
  [`task-plan.md`](../../../../spec/planning/task-plan.md)
  "N = 1 → N ≥ 2 transition," the split is warranted precisely
  when the additional phase "needs the full apparatus (separate
  Contracts surface, separate Validation Gate, separate
  Self-Review Audits)." Bundling them into one gate would let
  the risky subprocess review hide behind a green pure-template
  test run.
- **Why this is phases, not a re-decomposition into two
  tasks.** The level picker
  ([`task-plan.md`](../../../../spec/planning/task-plan.md))
  asks whether a half-ship has independent value. P1 alone
  ships independent value (descriptions + author-curated PRs
  render). P2 alone ships **nothing** — it augments P1's
  already-rendered PR list with auto-discovered entries; with
  no P1 PR-render surface there is no artifact. P2 is a
  sequence-step refinement of P1's PR surface, which is the
  phase signal, not the task signal. The milestone also owns
  task decomposition and locked m1 at four tasks; re-cutting t4
  into two tasks is out of t4's authority. Phases under the
  locked task t4 is the correct grain.
- **Just-in-time consequence.** P2's novel-mechanism spike (D5)
  runs at P2 drafting, which happens just-in-time *after* P1
  lands, per
  [`task-plan.md`](../../../../spec/planning/task-plan.md)
  "Just-in-time scoping and plan drafting." Isolating P2 keeps
  the spike at its just-in-time moment instead of forcing it
  now for code that drafts later. **Only the task plan and the
  P1 phase plan are drafted in this session;** the P2 phase
  plan is drafted just-in-time after P1's implementing PR
  merges. The task plan enumerates P2 as a phase with a named
  handoff so the structure is committed without pre-locking P2
  cross-phase contracts (per
  [`task-plan.md`](../../../../spec/planning/task-plan.md)
  "Cross-PR coordination").
- **Disclosure:** this is a finely-balanced call; the t3
  precedent and the small total surface genuinely pull toward
  N = 1. The deciding factor recorded for review re-litigation
  is the distinct environment-dependent Validation Gate and
  Self-Review surface of the subprocess phase, plus the
  just-in-time spike interaction.
- `Verified by:` reality-check "No subprocess execution exists
  anywhere" (P2's distinct mechanism);
  [`m1-t3-descriptive-labels.md`](../m1-t3-descriptive-labels.md)
  (the N = 1 precedent being departed from, recorded for
  contrast); [`m1-v0-2.md`](../m1-v0-2.md) "Task Status" legend
  and Sizing (milestone owns the four-task lock).

### D3 — `related_prs` is a new optional frontmatter field; additive spec change ships in P1's PR

`related_prs` is an **optional** frontmatter field: a YAML
block sequence of PR references. A doc omitting it renders with
no warning, error, or skip — identical absence-tolerance to
`short_description`. The `spec/planning/shared.md` documentation
of the field lands in the same implementing PR as P1's parser/
render change, not as a separate spec-only PR.

- **Rejected: separate spec-only PR landing first.** Rejected
  for the same reason t3 rejected it (t3 scoping D2): it opens
  a window where the vendored spec documents a field no code
  reads, and the field plus its sole reader are one contract.
- **PR-reference value shape — RESOLVED for P1.** Each entry is
  a string. P1 accepts and renders entries verbatim as the
  link/label; it does **not** parse, validate, or canonicalize
  PR-reference syntax. Canonical-identity normalisation (needed
  only to dedupe frontmatter entries against `gh`-discovered
  ones) is **P2's** contract, resolved at P2 drafting after the
  spike (D5) establishes the `gh --json` shape. P1 has no `gh`
  source to dedupe against, so it needs no canonical key —
  recording this split here prevents P1 from over-building a
  normaliser with no consumer.
- `Verified by:`
  [`spec/planning/shared.md`](../../../../spec/planning/shared.md)
  lines 135-154 (the optional-field block `short_description`
  established; `related_prs` documents adjacent with identical
  additive posture);
  [`m1-t3-descriptive-labels.md` scoping D2](m1-t3-descriptive-labels.md)
  (same-PR spec-change precedent).

### D4 — `gh` auto-discovery is best-effort, per-request, gracefully degrading; frontmatter is authoritative

(P2 contract direction; recorded at scoping so P2 drafting
inherits a decided posture, not an open question.) The
`gh pr list` query runs inside the existing per-request walk —
**no caching, no background refresh, no in-memory build-up** —
because the m1 Cross-Task Invariant forbids introducing any of
those within the milestone. Frontmatter `related_prs` is
authoritative and **always** rendered; `gh`-discovered PRs are
merged in additively. Every `gh` failure mode — binary absent,
not authenticated, no network, not a git repo, non-zero exit,
malformed JSON, or exceeding a bounded timeout — is
logged-and-skipped and **never** fails the page: the node still
renders its frontmatter PRs (or none).

- **Rejected: cache/TTL the `gh` result.** Directly violates
  [`m1-v0-2.md`](../m1-v0-2.md) Cross-Task Invariant "Render
  path stays walk-on-every-request ... No task in m1 introduces
  caching, file watching, or in-memory build-up."
- **Rejected: background/async refresh goroutine.** Same
  invariant (in-memory build-up) plus it adds concurrency a
  single-user local page does not warrant.
- **Rejected: drop `gh` auto-discovery, frontmatter only.** The
  milestone t4 contract explicitly requires "auto-discovers via
  `gh pr list` grep against the slug." Dropping it is a
  contract breach, not a t4-level scope call; if P2's spike
  surfaces a true dealbreaker, that escalates to the milestone
  per the milestone's own risk-escalation posture, it is not
  silently descoped here.
- **Acceptability of per-request shell-out.** Justified by the
  single-contributor / single-local-environment cross-cutting
  invariant: one local user loading a page tolerates a bounded
  local `gh` call; this is not a hosted multi-request service.
  The bounded timeout (concrete value resolved at P2 drafting)
  caps the worst case.
- `Verified by:` [`m1-v0-2.md`](../m1-v0-2.md) Cross-Task
  Invariant (caching ban); [`README.md`](../README.md)
  cross-cutting invariant (single local environment);
  reality-check "render path is walk-on-every-request."

### D5 — P2's subprocess is a novel mechanism; a just-in-time spike is required at P2 drafting

The `gh pr list` shell-out is the codebase's first subprocess
(reality-check). Per
[`task-plan.md`](../../../../spec/planning/task-plan.md) "Spike
before plan for novel mechanisms," P2 drafting must run a
30-minute throwaway spike (branch `spike/m1-t4-gh-prlist`, not
merged) before writing the P2 plan. The spike's job is to find
dealbreakers in the *external-tool semantics* the milestone
phrase "grep against the slug" leaves open:

- which `gh pr list` field the slug actually matches —
  `--search` over title/body, branch-name match, or a label —
  and whether that returns the intended PRs for this repo's
  slug-to-branch convention;
- the exact `gh pr list --json <fields>` output shape P2 will
  decode (drives D3's deferred canonical-identity key);
- real behaviour of the failure matrix in D4 (no binary / no
  auth / no network / non-repo) so P2's graceful-degradation
  contract rests on observed exit/stderr behaviour, not
  assumption.

- **Rejected: skip the spike, design P2 from `gh` docs alone.**
  The slug→PR match semantics is exactly the
  "wrong-assumption-about-external-tool" risk the spike rule
  targets; `os/exec` plumbing is well-trodden but *what query
  returns the right PRs* is not, and getting it wrong ships a
  feature that silently returns nothing.
- `Verified by:` reality-check "No subprocess execution exists
  anywhere";
  [`task-plan.md`](../../../../spec/planning/task-plan.md)
  "Spike before plan for novel mechanisms" (worktree handling:
  throwaway branch, never promoted into the implementation PR).

## Reality-check inputs the plans must verify before promotion

Re-confirm at plan-draft time (line numbers are navigation
aids; the symbolic anchors are load-bearing):

- `internal/site/walker.go` `parsePlanDoc` / `parsedDoc` —
  still reads frontmatter via the tolerant
  `metaData[k].(T)`-with-`, _ :=` pattern, and
  `LongDescription` is still set from `markdownBody(source)`
  (P1's `related_prs` read grafts onto this exact shape; P1's
  long-description render depends on the carry being live).
- `internal/site/tree.go` `buildTree` / `PlanNode` — still
  copies `LongDescription` onto every node and constructs nodes
  in the one loop (P1 adds a `RelatedPRs` field alongside, same
  loop).
- `internal/site/render.go` `{{define "node"}}` — still the
  sole node-emission site and still a no-JS single template
  (D1's "inline, no panel" depends on this; if a script/route
  appeared, D1 re-opens).
- `goldmark-meta` version pinned in `go.mod` — a YAML block
  sequence still decodes to `[]interface{}` of `string` (the
  tagged assumption behind P1's list-field read).
- `spec/planning/shared.md` lines ~135-154 — the optional
  `short_description` block is still the adjacency point for
  the additive `related_prs` documentation (D3).
- (P2, at P2 drafting) `grep -rn "os/exec"` is still empty and
  `go.mod` still declares no `gh` dependency — i.e. P2 is still
  introducing the first subprocess (D5's novel-mechanism
  premise).

## Plan-structure handoff

- **Doc-type:** task plan, **N ≥ 2** (D2). Orchestrating doc at
  [`../m1-t4-expanded-per-node-display.md`](../m1-t4-expanded-per-node-display.md);
  phase plans at `../m1-t4-p1-inline-detail-render.md` and
  (drafted just-in-time, post-P1) `../m1-t4-p2-gh-discovery.md`.
- **This session drafts:** the task plan + the P1 phase plan
  only. P2's phase plan is drafted just-in-time after P1's
  implementing PR merges (D2, D5).
- **Task plan owns** (per
  [`task-plan.md`](../../../../spec/planning/task-plan.md)
  "Required and optional sections" + "Task-plan-specific
  rules"): the context preamble, the Phase Contracts (per-phase
  WHAT), Cross-Phase Decisions (D1/D3/D4 framing and the
  P1↔P2 related-PR merge boundary), Cross-Cutting Invariants
  threading both phases (additive-frontmatter, no-caching,
  bare-bones-no-JS), the sequencing rationale, and the terminal
  -state rule (task plan flips `Landed` with P2's PR; it also
  closes the parent milestone t4 row).
- **P1 phase plan owns:** Goal + context, Contracts (the
  `related_prs` frontmatter field, `parsedDoc`/`PlanNode`
  `RelatedPRs` carry, the inline `node`-template detail block,
  the additive `spec/planning/shared.md` edit), Files to touch,
  Validation Gate (Go toolchain + manual render incl. the D1
  bans-on-surface observation and a field-less node), Naming,
  Self-Review Audits, Execution Steps, Risk Register,
  Documentation currency.
- **P2 phase plan (future) owns:** the spike (D5), the `gh`
  invocation + bounded timeout, the graceful-degradation
  contract (D4), the frontmatter↔`gh` merge/dedupe by the
  canonical key (D3's deferred half), and a Validation Gate
  that exercises the failure matrix against real environments.
- **Section variance from
  [`task-plan.md`](../../../../spec/planning/task-plan.md)
  required+optional list:** none expected; an unlisted section,
  if added, is disclosed per
  [`shared.md`](../../../../spec/planning/shared.md) "Section
  variance disclosure."

## Open decisions carried to plan-drafting

- **P1 — none.** D1, D3 (P1 half), and the structure decision
  (D2) resolve every call P1 needs; t3's `LongDescription`
  carry is **Landed**, not a pending input. The P1 plan
  proceeds through the `In draft → Proposed` promotion gate in
  this session.
- **P2 — deferred to P2's just-in-time drafting**, each citing
  a concrete surface per
  [`task-plan.md`](../../../../spec/planning/task-plan.md)
  "Just-in-time scoping and plan drafting":
  - `gh pr list` query/field semantics and `--json` shape —
    resolved by the D5 spike at P2 drafting (cited surface:
    this doc's D5; spike branch `spike/m1-t4-gh-prlist`).
  - Canonical PR-identity dedupe key (D3 deferred half) —
    resolved at P2 drafting once the spike fixes the `gh`
    shape (cited surface: D3 + D5).
  - Bounded `gh` timeout value (D4) — resolved at P2 drafting
    (cited surface: D4).
  These are P2-plan open inputs, not P1 blockers; the task plan
  records them as P2's named handoff and they do not gate P1's
  promotion.
