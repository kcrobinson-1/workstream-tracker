# Agent Guidance

Map of `docs/agents/`:

- [`shared/`](shared/) — vendored from `shared-agent-rules`.
  **Do not edit directly.** Edit the upstream module or the
  matching local overlay; re-run `scripts/assemble.sh` to
  regenerate.
- [`local/`](local/) — repo-specific content.
  - [`local/reference/`](local/reference/) — repo-specific
    constraint sets (architecture guardrails, styling tokens,
    framework conventions, anything declared as a mandatory
    pre-edit read).
  - [`local/workflows/`](local/workflows/) — repo-specific
    workflow additions or extensions to a shared workflow.
  - [`local/self-review-catalog.md`](local/self-review-catalog.md)
    — repo-specific audits. The seed audits imported from the
    shared library live under
    [`shared/self-review/seed-audits/`](shared/self-review/seed-audits/).
- [`shared.manifest.yaml`](shared.manifest.yaml) — the manifest
  that names which shared modules this repo vendors and at which
  version. Lives at `docs/agents/shared.manifest.yaml` (alongside
  this README); `scripts/assemble.sh` reads it by default.

Universal session rules live in [`/AGENTS.md`](/AGENTS.md) at
the repo root and the linked shared modules.

## Editing flow

- To change a shared module: open a proposal against the
  `shared-agent-rules` repo (see its `proposals/` directory).
- To change a local rule: edit the relevant `local/` file
  directly. The change description must follow
  [`shared/meta/rule-additions.md`](shared/meta/rule-additions.md).
- To add a local overlay that extends a shared module: create
  the corresponding file under `local/overlays/<module-path>.md`
  and re-run `scripts/assemble.sh`. The overlay is appended to
  the generated file under a `## Local additions` header.
