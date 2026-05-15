# Repo-specific reference constraint sets

Placeholder directory for repo-specific constraint sets. The
shared rule library does not ship reference files because the
reference layer (architecture guardrails, styling tokens,
framework conventions, persistence-layer rules) is necessarily
project-specific.

This directory holds:

- **Architecture guardrails.** What modules own what; what
  boundaries are load-bearing; what duplication is explicitly
  forbidden.
- **Styling / design-token discipline.** How themable vs.
  structural styles are organized; what counts as a valid
  token use.
- **Framework conventions.** Adapter patterns, environment
  reading, singleton lifecycles, framework-specific gates.
- **Persistence-layer rules.** Migration discipline, schema
  conventions, trust-boundary specifics for the data store
  the repo uses.
- **Any other binding constraint set** named in
  the root `AGENTS.md` "Mandatory pre-edit reads" section.

These files are loaded at the moment named in the relevant
workflow file (typically the pre-edit gate). They are not
optional lookups; the workflow points at them as binding.

## Naming convention

Use kebab-case filenames that describe the topic, not the
trigger: `architecture-guardrails.md`,
`styling-tokens.md`, `persistence-conventions.md`. The trigger
(which paths fire the read) lives in the root `AGENTS.md`
"Mandatory pre-edit reads" section, not in the filename.
