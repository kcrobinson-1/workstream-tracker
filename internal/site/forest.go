package site

import "html/template"

// forestTemplates owns the forest region: the "forest" template
// (roots stacked, the no-roots empty-state), the recursive "node"
// template (each node a nested, independently-collapsible box —
// header with the Status badge and actor markers, body with long
// description, related PRs, and child boxes), and "forest-style"
// (the node-level CSS the shell injects into <head>). This is m2
// t2's owned surface (expanded nested-box render); the shell
// (render.go) and roster (roster.go) are not t2's.
//
// A box with children is a native <details>/<summary> (no
// JavaScript, no route — m2 t2 C2): its summary is the clickable
// header and its children render as nested boxes inside the body.
// A leaf box (no children) has no disclosure control and no
// expand state (C5) — its header and detail are always visible.
// node-header / node-detail are factored out so both the
// collapsible and leaf branches render the identical preserved
// surfaces (C5: actor markers, long description, related PRs).
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

{{define "node-header"}}<span class="label-group"><span class="label" title="{{.Slug}}">{{.Label}}</span>{{range .WorkInstances}}<span class="actor-marker">{{.Actor}}</span>{{end}}</span><span class="status-group"><span class="badge status-{{statusClass .Status}}">{{if .Status}}{{.Status}}{{else}}(no Status){{end}}</span></span>{{end}}

{{define "node-detail"}}
{{- if .LongDescription}}
<div class="long-desc">{{.LongDescription}}</div>
{{- end}}
{{- if .RelatedPRs}}
<ul class="related-prs">
{{range .RelatedPRs}}<li>{{if isURL .}}<a href="{{.}}">{{.}}</a>{{else}}{{.}}{{end}}</li>
{{end}}</ul>
{{- end}}
{{- end}}

{{define "node"}}
{{- if .Children}}
<details class="box box-{{.NodeType}}" open>
<summary><span class="box-header">{{template "node-header" .}}</span></summary>
<div class="box-body">
{{- template "node-detail" .}}
{{range .Children}}{{template "node" .}}
{{end}}</div>
</details>
{{- else}}
<div class="box box-{{.NodeType}} box-leaf">
<div class="box-header">{{template "node-header" .}}</div>
{{- template "node-detail" .}}
</div>
{{- end}}
{{end}}

{{define "forest-style"}}    .root { margin-bottom: 1.5rem; }
    .box { border: 1px solid #e5e7eb; border-radius: 0.5rem; margin: 0.4rem 0; background: #fff; }
    .box-root { background: #f9fafb; border-color: #d1d5db; }
    .box-milestone { background: #fcfcfd; }
    .box-phase { border-style: dashed; }
    .box-body { padding: 0 0.75rem 0.5rem 1rem; }
    summary { cursor: pointer; }
    summary .box-header { display: flex; }
    .box-header { display: flex; align-items: baseline; justify-content: space-between; gap: 0.75rem; flex-wrap: wrap; padding: 0.4rem 0.75rem; }
    .label-group { display: flex; align-items: baseline; gap: 0.4rem; flex-wrap: wrap; min-width: 0; }
    .status-group { flex: 0 0 auto; }
    .badge { display: inline-block; padding: 0.1rem 0.5rem; border-radius: 0.25rem; font-size: 0.85em; font-weight: 500; }
    .status-in-draft    { background: #fef3c7; color: #78350f; }
    .status-proposed    { background: #dbeafe; color: #1e3a8a; }
    .status-in-progress { background: #fed7aa; color: #7c2d12; }
    .status-validating  { background: #e9d5ff; color: #581c87; }
    .status-landed      { background: #d1fae5; color: #065f46; }
    .status-deferred    { background: #e5e7eb; color: #374151; }
    .status-unknown     { background: #f3f4f6; color: #6b7280; }
    .actor-marker { display: inline-block; background: #fef9c3; color: #713f12; padding: 0.05rem 0.4rem; border-radius: 0.25rem; font-size: 0.75em; }
    .label { font-weight: 500; cursor: help; }
    .empty { color: #6b7280; font-style: italic; }
    .long-desc { white-space: pre-wrap; margin: 0.25rem 0 0.25rem 0; color: #374151; font-size: 0.9em; }
    .related-prs { list-style: none; margin: 0.25rem 0 0.25rem 0; padding-left: 1.25rem; font-size: 0.85em; }
    .related-prs li { padding: 0.1rem 0; }{{end}}`

func init() {
	template.Must(indexTmpl.Parse(forestTemplates))
}
