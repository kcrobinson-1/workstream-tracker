# Agent Instructions

This file is the **router** for AI coding agents working in this repo.
It carries the universal rules every session needs (imported from the
shared library), the session-type routing table that names which files
an agent reads for the work at hand, and pointers to the
contributor-workflow source of truth.

Per-session-type playbooks and topic-organized constraint sets live
under [`docs/agents/`](docs/agents/); see
[`docs/agents/README.md`](docs/agents/README.md) for the directory
map.

## Repo orientation

`workstream-tracker` is a local tool that gives a contributor
cross-agent visibility into parallel planning and implementation
work. The local server reads a structured plan tree from the
contributor's repo, accepts agent registration calls for
work-instances active against plan-tree nodes, and renders an
HTML visualization that joins the two. The repo also owns the
canonical plan-doc spec (`spec/`) that consumer projects adopt.

Top-level layout:

- [`cmd/workstream-tracker/`](cmd/workstream-tracker/) — entry
  point for the local server.
- [`internal/`](internal/) — server packages: `api/` (HTTP API
  for agents), `site/` (HTML rendering), `db/` (SQLite access),
  `models/` (shared types).
- [`spec/`](spec/) — the canonical plan-doc spec consumer
  projects adopt. Start with
  [`spec/planning/shared.md`](spec/planning/shared.md) and the
  per-level files (`epic.md`, `milestone.md`, `task-plan.md`).
- [`design/`](design/) — vision, design docs, and diagrams for
  the tool itself.
- [`docs/plans/`](docs/plans/) — actual plan-tree docs the
  server visualizes (one root-slug subdir per plan tree).
- [`docs/agents/`](docs/agents/) — vendored and local agent
  rules (this file's router targets).

Canonical entry points:
- [`README.md`](README.md) — what the repo is, how to run it.
- [`design/v0.1-design.md`](design/v0.1-design.md) — first
  shipping version's spec.
- [`docs/dev.md`](docs/dev.md) — contributor workflow source of
  truth (see next section).

## Development workflow source of truth

[`docs/dev.md`](docs/dev.md) is the contributor-workflow source of
truth for this repo. It handles human contributor procedure (local
setup, validation commands, release flow); this AGENTS.md handles
agent decision discipline. The two are jointly authoritative — if
they conflict, stop and report rather than picking a side.

<!-- audit-coverage: R-02 (agent-rule vs contributor-doc conflict → stop and report) -->

## Session-type routing

Pick the row that best fits the work at hand and read the named
files. Universal rules below apply to every session and are not
enumerated in the table.

| If your session is… | Read these files |
|---|---|
| Implementation work without a plan doc to consume | [`docs/agents/shared/workflows/implementation.md`](docs/agents/shared/workflows/implementation.md) |
| Implementing a documented plan | [`docs/agents/shared/workflows/implementation.md`](docs/agents/shared/workflows/implementation.md) + plan-implementation rules from [`spec/planning/task-plan.md`](spec/planning/task-plan.md) + the plan's own `Cross-Cutting Invariants` and named self-review audits |
| Addressing review feedback | [`docs/agents/shared/workflows/review-fixes.md`](docs/agents/shared/workflows/review-fixes.md) |
| Debugging a failing validation | [`docs/agents/shared/workflows/debugging.md`](docs/agents/shared/workflows/debugging.md) |
| UI review / screenshot capture (if applicable) | [`docs/agents/shared/workflows/ui-review.md`](docs/agents/shared/workflows/ui-review.md) |

Reference files under [`docs/agents/local/reference/`](docs/agents/local/reference/)
are topic-organized constraint sets specific to this repo. They are
not optional lookups; the workflow files name when each fires.

## Mandatory pre-edit reads

Each entry is binding when its Triggers fire. The list is sparse
on purpose; add an entry when a class of edit has repeatedly
slipped without consulting a binding doc.

### plan-doc-authoring

Intent: edits to plan-tree docs must follow the canonical plan-doc
spec for slug hierarchy, status lifecycle, frontmatter shape, and
cross-cutting invariants.
Triggers:
  - `docs/plans/**`
Read: [`spec/planning/shared.md`](spec/planning/shared.md), then
the level-specific file in `spec/planning/` matching the doc's
plan-tree level (`epic.md`, `milestone.md`, or `task-plan.md`).

### spec-authoring

Intent: edits to the plan-doc spec itself change the contract every
consumer follows. Read the spec's own README before touching it.
Triggers:
  - `spec/**`
Read: [`spec/README.md`](spec/README.md) and
[`spec/planning/shared.md`](spec/planning/shared.md).

## Universal session rules

The universal rule set lives in vendored shared modules. Every
session loads them; they are not optional.

- [Pre-Edit Gate](docs/agents/shared/core/pre-edit-gate.md)
- [Scope Guardrails and Stop-And-Report](docs/agents/shared/core/scope-and-stop.md)
- [Change Boundaries](docs/agents/shared/core/change-boundaries.md)
- [Anti-Patterns](docs/agents/shared/core/anti-patterns.md)
- [Sub-Agent Delegation](docs/agents/shared/delegation/sub-agent-delegation.md)

Repo-owned universal rule (not vendored shared):

- [Session-start work-instance registration](docs/agents/local/session-registration.md)
  — before task work, perform the observable best-effort grounded
  narration handshake (resolve slug → confirm → echo the real
  `workstream-tracker register` receipt → narrate failure
  explicitly → proceed). Never blocks or fails the session.

## Self-review

Before finishing, run the audits from
[`docs/agents/local/self-review-catalog.md`](docs/agents/local/self-review-catalog.md)
that match the diff's surfaces. The catalog mechanism is documented
in [`docs/agents/shared/self-review/mechanism.md`](docs/agents/shared/self-review/mechanism.md);
the seed audits ship under
[`docs/agents/shared/self-review/seed-audits/`](docs/agents/shared/self-review/seed-audits/);
project-specific audits live in the local catalog.

Layer the general self-review checklist from
[`docs/agents/shared/self-review/how-to-use.md`](docs/agents/shared/self-review/how-to-use.md)
on top of the catalog walk.

## Adding to this rule set

Changes that add or modify rule content under
[`docs/agents/local/**`](docs/agents/local/) or this file must
follow
[`docs/agents/shared/meta/rule-additions.md`](docs/agents/shared/meta/rule-additions.md):
name a rule the new rule retires or merges into, or state why no
existing rule could be retired.

The shared modules under
[`docs/agents/shared/**`](docs/agents/shared/) are vendored
output; do not edit them directly. To change shared content,
file a proposal upstream (see
`shared-agent-rules/proposals/` in the shared repo).
