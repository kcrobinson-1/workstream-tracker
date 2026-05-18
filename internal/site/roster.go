package site

import "html/template"

// rosterTemplates owns the roster region: the "roster" template
// (currently a deliberate, intentional placeholder — m2 t1 ships
// the shell, the bound/unbound session roster is m2 t4) and
// "roster-style" (the roster CSS the shell injects into <head>).
// This is the surface m2 t4 grows; it must stay a deliberate
// observed state (not a blank gap) and must not render
// data-shaped content that pre-empts t4.
const rosterTemplates = `
{{define "roster"}}<section class="roster-panel">
    <h2 class="roster-title">Sessions</h2>
    <p class="roster-placeholder">The session roster lands in a later task. Once built, this region will list every registered active session — bound and unbound — alongside the forest.</p>
  </section>{{end}}

{{define "roster-style"}}    .roster-panel { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 0.5rem; padding: 0.75rem 1rem; }
    .roster-title { margin: 0 0 0.25rem 0; font-size: 0.95rem; font-weight: 600; color: #0f172a; }
    .roster-placeholder { margin: 0; color: #6b7280; font-style: italic; font-size: 0.9em; }{{end}}`

func init() {
	template.Must(indexTmpl.Parse(rosterTemplates))
}
