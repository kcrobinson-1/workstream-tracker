package site

import (
	"html/template"
	"io"
	"strings"
)

// indexTmpl renders the v0.1 forest visualization: roots
// stacked, descendants nested as bullet lists, each node carrying
// a colored Status badge plus markers for any active
// work-instances. Bare-bones by design (see
// design/v0.1-design.md Section 7).
var indexTmpl = template.Must(template.New("index").Funcs(template.FuncMap{
	"statusClass": statusClass,
}).Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>workstream-tracker</title>
  <style>
    body { font-family: system-ui, -apple-system, sans-serif; padding: 1rem; max-width: 60rem; margin: 0 auto; color: #111827; }
    h1 { margin-top: 0; }
    .root { margin-bottom: 1.5rem; padding: 0.75rem 1rem; background: #f9fafb; border-radius: 0.5rem; }
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
    .slug { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 0.95em; }
    .empty { color: #6b7280; font-style: italic; }
  </style>
</head>
<body>
  <h1>workstream-tracker</h1>
  {{if .Roots}}
  {{range .Roots}}
  <div class="root">
    {{template "node" .}}
  </div>
  {{end}}
  {{else}}
  <p class="empty">No plan-tree roots found at <code>{{.PlansPath}}</code>.</p>
  {{end}}
</body>
</html>

{{define "node"}}
<span class="badge status-{{statusClass .Status}}">{{if .Status}}{{.Status}}{{else}}(no Status){{end}}</span><span class="slug">{{.Slug}}</span>
{{- range .WorkInstances }} <span class="actor-marker">{{.Actor}}</span>{{end}}
{{- if .Children}}
<ul>
  {{range .Children}}<li>{{template "node" .}}</li>
  {{end}}
</ul>
{{end}}
{{end}}
`))

// indexData is the template input passed to indexTmpl.
type indexData struct {
	Roots     []*PlanNode
	PlansPath string
}

// renderIndex writes the rendered index page to w.
func renderIndex(w io.Writer, data indexData) error {
	return indexTmpl.ExecuteTemplate(w, "index", data)
}

// statusClass converts a Status string to a CSS class fragment.
// Recognized values from spec/planning/shared.md "Plan-doc
// Status" map to specific classes; unrecognized values fall back
// to "unknown" per the spec's graceful-fallback rule.
func statusClass(status string) string {
	// Strip the em-dash freeform suffix from Deferred — <reason>.
	canonical := status
	if i := strings.Index(canonical, " — "); i >= 0 {
		canonical = canonical[:i]
	}
	switch canonical {
	case "In draft":
		return "in-draft"
	case "Proposed":
		return "proposed"
	case "In progress":
		return "in-progress"
	case "Validating":
		return "validating"
	case "Landed":
		return "landed"
	case "Deferred":
		return "deferred"
	}
	return "unknown"
}
