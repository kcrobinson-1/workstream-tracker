package site

import "html/template"

// forestTemplates owns the forest region: the "forest" template
// (roots stacked, descendants nested as bullet lists, the
// no-roots empty-state), the recursive "node" template (Status
// badge, actor markers, long description, related PRs, children),
// and "forest-style" (the node-level CSS the shell injects into
// <head>). This is the surface m2 t2 (expanded nested-box render)
// grows; it is relocated byte-for-intent from the pre-split
// render.go and must not change node shape here (t2's call).
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

{{define "node"}}
<span class="badge status-{{statusClass .Status}}">{{if .Status}}{{.Status}}{{else}}(no Status){{end}}</span><span class="label" title="{{.Slug}}">{{.Label}}</span>
{{- range .WorkInstances }} <span class="actor-marker">{{.Actor}}</span>{{end}}
{{- if .LongDescription}}
<div class="long-desc">{{.LongDescription}}</div>
{{- end}}
{{- if .RelatedPRs}}
<ul class="related-prs">
  {{range .RelatedPRs}}<li>{{if isURL .}}<a href="{{.}}">{{.}}</a>{{else}}{{.}}{{end}}</li>
  {{end}}
</ul>
{{- end}}
{{- if .Children}}
<ul>
  {{range .Children}}<li>{{template "node" .}}</li>
  {{end}}
</ul>
{{end}}
{{end}}

{{define "forest-style"}}    .root { margin-bottom: 1.5rem; padding: 0.75rem 1rem; background: #f9fafb; border-radius: 0.5rem; }
    ul { padding-left: 1.5rem; list-style: none; margin: 0.25rem 0; }
    li { padding: 0.25rem 0; }
    .badge { display: inline-block; padding: 0.1rem 0.5rem; border-radius: 0.25rem; font-size: 0.85em; margin-right: 0.5rem; font-weight: 500; }
    .status-in-draft    { background: #fef3c7; color: #78350f; }
    .status-proposed    { background: #dbeafe; color: #1e3a8a; }
    .status-in-progress { background: #fed7aa; color: #7c2d12; }
    .status-validating  { background: #e9d5ff; color: #581c87; }
    .status-landed      { background: #d1fae5; color: #065f46; }
    .status-deferred    { background: #e5e7eb; color: #374151; }
    .status-unknown     { background: #f3f4f6; color: #6b7280; }
    .actor-marker { display: inline-block; background: #fef9c3; color: #713f12; padding: 0.05rem 0.4rem; border-radius: 0.25rem; font-size: 0.75em; margin-left: 0.4rem; }
    .label { font-weight: 500; cursor: help; }
    .empty { color: #6b7280; font-style: italic; }
    .long-desc { white-space: pre-wrap; margin: 0.25rem 0 0.25rem 0; color: #374151; font-size: 0.9em; }
    .related-prs { margin: 0.25rem 0 0.25rem 0; padding-left: 1.25rem; font-size: 0.85em; }
    .related-prs li { padding: 0.1rem 0; }{{end}}`

func init() {
	template.Must(indexTmpl.Parse(forestTemplates))
}
