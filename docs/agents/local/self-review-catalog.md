# Self-Review Audit Catalog

Named, reusable audits to run on a diff before push. Each
audit is a reviewer's lens — a specific class of issue that
reviewers (human or AI) would otherwise flag.

The catalog mechanism — entry format, lifecycle, how the
catalog runs — lives in
[`../shared/self-review/mechanism.md`](../shared/self-review/mechanism.md).

The general self-review checklist that runs *alongside* the
named audits lives in
[`../shared/self-review/how-to-use.md`](../shared/self-review/how-to-use.md).

## Seeded from shared library

The following audits are vendored from
`shared-agent-rules/library/self-review/seed-audits/` and live
under
[`../shared/self-review/seed-audits/`](../shared/self-review/seed-audits/).
They apply to every consuming repo; trigger them against the
diff surface, the same way as repo-specific audits below.

- **effect-cleanup** — any effect that subscribes / opens /
  schedules / installs a listener has a matching cleanup on
  every exit path.
- **error-surfacing-user-mutations** — user-initiated
  mutations surface errors back to the user; silent failure
  is the trap.
- **validation-honesty** — a claim that a check ran is only
  valid if the check ran end-to-end on the current state.
- **rename-aware-diff-classification** — renames that tooling
  treats as add+delete need a hand-classification pass before
  "no behavior changed" claims.
- **readiness-gate-truthfulness** — a gate that announces
  readiness must actually verify the named conditions, not
  announce on best-effort.
- **trigger-map-currency** — when the codebase restructures
  (renamed/moved directories, new top-level dirs, new files
  matching a declared Intent), the Mandatory pre-edit reads
  trigger map is re-evaluated against the new state so
  deterministic triggers do not silently drift away from their
  Intent.

See [`../shared/self-review/seed-audits/`](../shared/self-review/seed-audits/)
for each audit's Trigger / Check / Example.

## Repo-specific audits

<!--
  Add audits here as the team encounters them. Follow the
  Trigger / Check / Example format from
  shared/self-review/mechanism.md.

  Add a new audit when EITHER:
  - The same class of issue is flagged on two or more distinct
    changes, OR
  - A single high-severity finding has a root cause that
    clearly generalizes beyond the specific file.

  Drop an audit when an automated check enforces it, or when
  the pattern hasn't recurred in a long while.
-->

_No repo-specific audits yet. Add them under this heading as
they emerge._
