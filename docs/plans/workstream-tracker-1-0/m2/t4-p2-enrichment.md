---
slug: workstream-tracker-1-0-m2-t4-p2
Status: In draft
short_description: Client/CLI metadata + named sessions + event-log join + expandable raw-JSON detail
---

# t4 p2 (stub) — Enrichment + named sessions

> **Parent-promotion stub.** Seeded from
> [`t4-session-roster.md`](t4-session-roster.md) "Phase Contracts"
> in the PR that flipped the t4 task plan to `Proposed`. Per
> [`shared.md`](../../../../spec/planning/shared.md) "Parent-doc
> child contracts," this stub is exempt from the
> [`task-plan.md`](../../../../spec/planning/task-plan.md) "Required
> and optional sections" rule and from "Plans describe contracts,
> not implementation" until this phase's own `In draft` →
> `Proposed` drafting session runs. It carries only the locked
> inherited WHAT below; HOW is scoped just-in-time against merged
> code at this phase's drafting.

## Inherited contract (locked at t4 task-plan drafting)

**End result + preserves.** The register client/CLI reports
request `metadata` (including a human session name); the roster
loader gains the **event-log join** to surface it; roster entries
show the reported name (with a defined, human-readable fallback
that is **never** the `wst-<uuid>` actor) and an **expandable
deliberately-unstructured raw-JSON detail view** of the session's
other reported fields. Preserves: no schema change (event-log join
only, no `work_instances` column); the schema-loose posture (the
reported shape is not schematized toward t3's declared-stages
field); registration stays observable best-effort and opt-in (a
send failure is surfaced, not silently dropped; a v0.2-minimum
session still lists); walk-on-every-request; the v0.1 forest actor
markers unchanged.

**Feeds siblings.** Consumes p1's roster loader, shared
`indexData` roster field, and region. Closes the milestone's
"every session can be accounted for" observability goal. Its
implementing PR executes the
[`deterministic-interactive-registration`](../../../backlog.md#deterministic-interactive-registration)
**split** and adds the new observability-residual and
forest-actor-humanization backlog entries, and lands the
`session-registration.md` / `design/v0.1-design.md` currency.

**Open HOW (mechanism, scoped at this phase's drafting).** The
behavior decisions are already task-level contract and inherited,
not open here: the metadata read policy (register baseline +
latest-event key-by-key overlay, per-request) and the
name-then-slug label rule (never the `wst-<uuid>`; only the
literal slug-fallback formatting is render-time UX copy). p2's
remaining HOW is pure mechanism: the metadata-join query shape
that realizes the decided baseline+overlay rule; the client/CLI
plumbing to send the conventionally-read `name` key (the
handshake-metadata path and a CLI affordance such as
`--name` / `WST_NAME`); and the expandable raw-JSON detail
rendering. See
[`scoping/t4-session-roster.md`](scoping/t4-session-roster.md)
"Decisions resolved at plan-drafting" and the t4 task plan
Contracts + Cross-Cutting Invariants.
