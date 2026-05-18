package site

import (
	"html/template"
	"io"
	"net/url"
	"strings"
)

// indexTmpl renders the v0.1 forest visualization: roots
// stacked, descendants nested as bullet lists, each node carrying
// a colored Status badge plus markers for any active
// work-instances. Bare-bones by design (see
// design/v0.1-design.md Section 7).
var indexTmpl = template.Must(template.New("index").Funcs(template.FuncMap{
	"statusClass": statusClass,
	"isURL":       isAbsoluteURL,
}).Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>workstream-tracker</title>
  <style>
    body { font-family: system-ui, -apple-system, sans-serif; padding: 1rem; max-width: 80rem; margin: 0 auto; color: #111827; }
    h1 { margin-top: 0; }
    /* Two-region shell (m2 t1): forest ~2/3 left, roster ~1/3
       right, top-aligned so the roster shows without scrolling on
       a tall window; the page itself is the only scroll. Below a
       desktop-narrow width the regions stack (roster under
       forest) — mobile is out of scope. */
    .layout { display: flex; gap: 1.5rem; align-items: flex-start; }
    .forest { flex: 2 1 0; min-width: 0; }
    .roster { flex: 1 1 0; min-width: 0; }
    @media (max-width: 60rem) { .layout { flex-direction: column; } }
    .roster-panel { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 0.5rem; padding: 0.75rem 1rem; }
    .roster-title { margin: 0 0 0.25rem 0; font-size: 0.95rem; font-weight: 600; color: #0f172a; }
    .roster-placeholder { margin: 0; color: #6b7280; font-style: italic; font-size: 0.9em; }
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
    .label { font-weight: 500; cursor: help; }
    .empty { color: #6b7280; font-style: italic; }
    .long-desc { white-space: pre-wrap; margin: 0.25rem 0 0.25rem 0; color: #374151; font-size: 0.9em; }
    .related-prs { margin: 0.25rem 0 0.25rem 0; padding-left: 1.25rem; font-size: 0.85em; }
    .related-prs li { padding: 0.1rem 0; }
  </style>
</head>
<body>
  <h1>workstream-tracker</h1>
  <div class="layout">
  <main class="forest">
  {{if .Roots}}
  {{range .Roots}}
  <div class="root">
    {{template "node" .}}
  </div>
  {{end}}
  {{else}}
  <p class="empty">No plan-tree roots found at <code>{{.PlansPath}}</code>.</p>
  {{end}}
  </main>
  <aside class="roster">
  <section class="roster-panel">
    <h2 class="roster-title">Sessions</h2>
    <p class="roster-placeholder">The session roster lands in a later task. Once built, this region will list every registered active session — bound and unbound — alongside the forest.</p>
  </section>
  </aside>
  </div>
</body>
</html>

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

// isAbsoluteURL reports whether s is a well-formed absolute URL
// (an http/https URL with a host), so the related-PR renderer can
// decide whether to emit an anchor. A non-URL entry (e.g. a
// `#123` shorthand) renders as escaped plain text instead — never
// a broken in-page anchor, an error, or a skip. This is a
// render-safety classification, not PR-reference parsing: P1
// derives no URL from shorthand, it only declines to linkify a
// non-URL.
func isAbsoluteURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return u.Host != ""
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
