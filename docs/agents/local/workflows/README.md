# Repo-specific workflow additions

Placeholder directory for repo-specific workflow files. The
shared workflows (`implementation.md`, `debugging.md`,
`review-fixes.md`, `ui-review.md`) live under
[`../../shared/workflows/`](../../shared/workflows/) and are vendored;
do not edit them here.

This directory holds:

- **Net-new repo-specific workflows.** If the repo has a session
  type the shared workflows don't cover (e.g., a deployment-
  walkthrough workflow, a data-migration workflow), add a new
  file here and route to it from the root `AGENTS.md` session-
  type table.
- **Extensions to a shared workflow.** If the repo needs to
  layer additional rules on top of a shared workflow without
  changing the shared content, create the addition here as a
  separate file and load it after the shared workflow. (For
  small in-line additions to a shared workflow, prefer the
  overlay mechanism — see
  [`../../README.md`](../../README.md) "Editing flow.")

When in doubt about whether the addition is repo-specific or
universal, file a proposal upstream rather than carrying it
here long-term.
