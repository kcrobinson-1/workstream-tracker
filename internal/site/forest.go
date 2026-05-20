package site

import "html/template"

// forestTemplates owns the forest region: the "forest" template
// (roots stacked, the no-roots empty-state), the recursive "node"
// template (each node a nested, independently-collapsible box —
// header with the Status badge and actor markers, body with the
// progress-cell row, long description, related PRs, and child
// boxes), and "forest-style" (the node-level CSS the shell injects
// into <head>). This is m2 t2's owned surface (expanded
// nested-box render); the shell (render.go) and roster
// (roster.go) are not t2's.
//
// m2 t3 adds the progress-cell row: every node renders a
// render-side-reserved Drafting cell followed by one cell per
// entry in its own doc's optional `progress_stages` frontmatter,
// in document order (N declared ⇒ N + 1 cells; a field-omitting
// doc or a `slug` + `Status: In draft` stub ⇒ exactly the one
// Drafting cell). Gating is by field presence only and is
// Status-independent; there is no inheritance — each node's own
// doc governs its own row.
//
// Every box — leaf or not — is a native <details>/<summary> (no
// JavaScript, no route — m2 t2 C2): the summary is the clickable
// header (always visible: label, Status badge, actor markers) and
// the body (progress row, long description, related PRs, child
// boxes) is hidden until the box is expanded. A box renders
// expanded (`open`) iff it, or any descendant, has an active
// work-instance (C3, via PlanNode.ActiveInSubtree); otherwise it
// is collapsed so the forest is a scannable header-only overview
// rather than a wall of plan-file text. A childless box keeps the
// box-leaf class for styling but is collapsible like any other.
// node-header / node-detail are factored out so the structure is
// identical regardless of whether the node has children.
//
// The long description is line-capped (truncateLongDesc) even when
// expanded: the forest gives context and the goal, not the whole
// plan file — read the doc itself for the full text.
const forestTemplates = `
{{define "forest"}}{{if .Roots}}
  {{range .Roots}}
  <div class="root">
    {{template "node" .}}
  </div>
  {{end}}
  {{else}}
  <p class="empty">No plan-tree roots found at <code>{{.PlansPath}}</code>.</p>
  {{end}}{{end}}

{{define "node-header"}}<span class="label-group"><span class="label" title="{{.Slug}}">{{.Label}}</span>{{range .WorkInstances}}<span class="actor-marker">{{if .Name}}{{.Name}}{{else}}{{.Slug}}{{end}}</span>{{end}}</span><span class="status-group"><span class="badge status-{{statusClass .Status}}">{{if .Status}}{{.Status}}{{else}}(no Status){{end}}</span></span>{{$aff := .Affordance}}{{if $aff}}<button class="affordance-button" type="submit" form="spawn-{{.Slug}}">{{$aff}}</button>{{end}}{{end}}

{{define "node-progress"}}<div class="progress-row">{{if .ProgressStages}}<span class="progress-cell progress-cell-drafting">Drafting</span>{{range .ProgressStages}}<span class="progress-cell">{{.}}</span>{{end}}{{else}}<span class="progress-cell progress-cell-{{progressCellClass .Status "d"}}">D</span><span class="progress-cell progress-cell-{{progressCellClass .Status "p"}}">P</span><span class="progress-cell progress-cell-{{progressCellClass .Status "i"}}">I</span><span class="progress-cell progress-cell-{{progressCellClass .Status "v"}}">V</span>{{end}}</div>{{end}}

{{define "node-detail"}}
{{- template "node-progress" .}}
{{- if or .LongDescription .RelatedPRs}}
<details class="body-disclosure"><summary class="body-disclosure-summary"><span class="triangle" aria-hidden="true">&#9656;</span>Show description</summary>
{{- if .LongDescription}}
<div class="long-desc">{{renderLongDescBody .LongDescription}}</div>
{{- end}}
{{- if .RelatedPRs}}
<ul class="related-prs">
{{range .RelatedPRs}}<li>{{if isURL .}}<a href="{{.}}">{{.}}</a>{{else}}{{.}}{{end}}</li>
{{end}}</ul>
{{- end}}
</details>
{{- end}}
{{- end}}

{{define "node"}}
<details class="box box-{{.NodeType}}{{if not .Children}} box-leaf{{end}}"{{if .ActiveInSubtree}} open{{end}}>
<summary><span class="box-header"><span class="triangle" aria-hidden="true">&#9656;</span>{{template "node-header" .}}</span></summary>
<div class="box-body">
{{- $aff := .Affordance}}{{if $aff}}<form class="affordance-form" id="spawn-{{.Slug}}" method="POST" action="/spawn"><input type="hidden" name="slug" value="{{.Slug}}"><input type="hidden" name="mode" value="{{$aff}}"></form>
{{end -}}
{{- template "node-detail" .}}
{{range .Children}}{{template "node" .}}
{{end}}</div>
</details>
{{end}}

{{define "forest-style"}}    .root { margin-bottom: 1.5rem; }
    .box { border: 1px solid #e5e7eb; border-radius: 0.5rem; margin: 0.4rem 0; background: #fff; }
    .box-root { background: #f9fafb; border-color: #d1d5db; }
    .box-milestone { background: #fcfcfd; }
    .box-phase { border-style: dashed; }
    .box-body { padding: 0 0.75rem 0.5rem 1rem; }
    /* m2 post-ux-correction p1 F8: suppress the UA-default
       <details> marker and render an inline triangle inside the
       header flex row, baseline-aligned with the label and rotated
       open via details[open]. list-style: none is the modern
       spec; ::-webkit-details-marker covers older WebKit. */
    summary { cursor: pointer; list-style: none; }
    summary::-webkit-details-marker { display: none; }
    summary .box-header { display: flex; }
    .box-header { display: flex; align-items: baseline; gap: 0.75rem; flex-wrap: wrap; padding: 0.4rem 0.75rem; }
    .triangle { display: inline-block; flex: 0 0 auto; font-size: 0.75em; color: #6b7280; transition: transform 0.1s ease; }
    details[open] > summary .triangle { transform: rotate(90deg); }
    .label-group { display: flex; align-items: baseline; gap: 0.4rem; flex-wrap: wrap; min-width: 0; }
    .status-group { flex: 0 0 auto; margin-left: auto; }
    .badge { display: inline-block; padding: 0.1rem 0.5rem; border-radius: 0.25rem; font-size: 0.85em; font-weight: 500; }
    .status-in-draft    { background: #fef3c7; color: #78350f; }
    .status-proposed    { background: #dbeafe; color: #1e3a8a; }
    .status-in-progress { background: #fed7aa; color: #7c2d12; }
    .status-validating  { background: #e9d5ff; color: #581c87; }
    .status-landed      { background: #d1fae5; color: #065f46; }
    .status-deferred    { background: #e5e7eb; color: #374151; }
    .status-unknown     { background: #f3f4f6; color: #6b7280; }
    .actor-marker { display: inline-block; background: #fef9c3; color: #713f12; padding: 0.05rem 0.4rem; border-radius: 0.25rem; font-size: 0.75em; }
    /* m2 t1: the per-node mode-affordance. The visible button
       lives in <summary>'s header row (trailing flex child after
       .status-group, per SD7); the associated <form> lives inside
       .box-body and is referenced by the button's form="<id>"
       attribute. The split keeps HTML5's phrasing-vs-flow content-
       model rule satisfied (<summary> permits phrasing content
       only; <form> is flow content) while still rendering the
       action in the always-visible region. HTML5 form-association
       by ID resolves via the DOM, so submission works whether the
       <details> is open or closed. The form's POST goes to t2's
       /spawn endpoint; a POST before t2 lands receives a 404 as
       the expected sibling-not-yet-shipped degrade. */
    .affordance-form { margin: 0; padding: 0; }
    .affordance-button { display: inline-block; padding: 0.15rem 0.6rem; border-radius: 0.25rem; background: #4f46e5; color: #fff; border: 1px solid #4338ca; font-size: 0.8em; font-weight: 500; cursor: pointer; font-family: inherit; flex: 0 0 auto; }
    .affordance-button:hover { background: #4338ca; }
    .label { font-weight: 500; cursor: help; }
    .empty { color: #6b7280; font-style: italic; }
    .long-desc { white-space: pre-wrap; margin: 0.25rem 0 0.25rem 0; color: #374151; font-size: 0.9em; }
    .related-prs { list-style: none; margin: 0.25rem 0 0.25rem 0; padding-left: 1.25rem; font-size: 0.85em; }
    .related-prs li { padding: 0.1rem 0; }
    /* m2 t3: the doc-declared progress-cell row. A
       render-side-reserved Drafting cell followed by one cell per
       declared progress_stages entry, in document order.
       p3 F3a (C2): when a doc declares no progress_stages, the
       default D / P / I / V row renders instead — every cell is
       its own DOM element across both branches per C-INV-1, so
       the F3b future per-cell attachment is preserved. */
    .progress-row { display: flex; flex-wrap: wrap; gap: 0.25rem; margin: 0.25rem 0 0.25rem 0; }
    .progress-cell { display: inline-block; padding: 0.1rem 0.5rem; border-radius: 0.25rem; font-size: 0.8em; background: #eef2ff; color: #3730a3; border: 1px solid #c7d2fe; }
    .progress-cell-drafting { background: #f3f4f6; color: #4b5563; border-color: #d1d5db; }
    /* p3 F3a default-row variants (OD1's conservative starting
       point: reuse the .status-landed and .status-unknown badge
       palette values so a cell's fill matches its Status badge). */
    .progress-cell-landed   { background: #d1fae5; color: #065f46; border-color: #6ee7b7; }
    .progress-cell-in-draft { background: #fef3c7; color: #78350f; border-color: #fcd34d; }
    .progress-cell-neutral  { background: #f3f4f6; color: #6b7280; border-color: #d1d5db; }
    .progress-cell-empty    { background: #ffffff; color: #9ca3af; border-style: dashed; border-color: #d1d5db; }
    /* p3 F7 (C1): the per-node body disclosure — a nested
       <details>/<summary> inside the per-node <details> box,
       independent of the parent tree-collapse. The summary uses
       the same inline-triangle treatment p1 introduced for the
       outer collapse (the global summary list-style:none rule
       above suppresses the UA marker for nested summaries too,
       so an inline triangle is needed for visual parity). */
    .body-disclosure { margin: 0.25rem 0 0.25rem 0; }
    .body-disclosure-summary { display: inline-flex; align-items: baseline; gap: 0.4rem; cursor: pointer; color: #4b5563; font-size: 0.85em; padding: 0.15rem 0; }
    .body-disclosure-summary .triangle { color: #9ca3af; }
    /* triangle rotation reuses the existing details[open] >
       summary .triangle rule above, which matches both the
       outer per-node summary and the nested body-disclosure
       summary independently — each tied to its own immediate
       parent <details>[open] state. The two disclosures stay
       structurally independent per parent C2. */
    /* m2 t2 C4: t2 owns the per-node-row narrow-window degrade.
       When the forest column is narrow the right-aligned Status
       group reflows below the label group as an intentional
       stacked layout, rather than colliding with or truncating
       the label. Distinct from the shell's region-stacking query
       (render.go, max-width 60rem). */
    @media (max-width: 48rem) {
      .box-header { flex-direction: column; align-items: flex-start; gap: 0.2rem; }
      .status-group { align-self: flex-start; margin-left: 0; }
      .affordance-button { align-self: flex-start; }
    }{{end}}`

func init() {
	template.Must(indexTmpl.Parse(forestTemplates))
}
