# workstream-tracker

A local tool for cross-agent visibility into parallel planning and
implementation work. The tool reads a structured plan tree from a
contributor's repo, accepts agent registration calls for
work-instances active against plan-tree nodes, and renders a
visualization that joins the two so the contributor can see what's
happening across their parallel work.

## Layout

- **[`design/`](design/)** — vision, design docs, and diagrams for
  the tool itself. Start with [design/vision.md](design/vision.md)
  for the long-form framing or
  [design/v0.1-design.md](design/v0.1-design.md) for the first
  shipping version's spec.
- **[`spec/`](spec/)** — the canonical plan-doc spec consumer
  projects adopt. The structure planning agents follow when
  authoring plan docs the tool visualizes. See
  [spec/planning/shared.md](spec/planning/shared.md) for the
  cross-level rules and [spec/backlog.md](spec/backlog.md) for the
  backlog format.

## Status

v0.0 — spec is drafted, tool implementation is pending.
