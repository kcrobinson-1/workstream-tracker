package site

import "html/template"

// rosterTemplates owns the roster region: the "roster" template
// (m2 t4 p1 — the bare bound/unbound session roster: every active
// work-instance listed, each marked bound when its slug matches a
// walked plan doc or unbound otherwise, with a deliberate
// "no active sessions" state) and "roster-style" (the roster CSS
// the shell injects into <head>). This is t4's owned region body;
// p2 enriches the same region with reported names and an
// expandable raw-JSON detail. The entry label is the slug — p1
// has no reported-name source and the wst-<uuid> actor is never
// rendered (scoping SD2). It must stay a deliberate observed
// state (not a blank gap) and adds no promote/dismiss affordance.
const rosterTemplates = `
{{define "roster"}}<section class="roster-panel">
    <h2 class="roster-title">Sessions</h2>
    {{if .Roster}}
    <ul class="roster-list">
      {{range .Roster}}<li class="roster-entry roster-{{if .Bound}}bound{{else}}unbound{{end}}"><span class="roster-slug">{{.Slug}}</span><span class="roster-tag">{{if .Bound}}bound{{else}}unbound{{end}}</span></li>
      {{end}}
    </ul>
    {{else}}
    <p class="roster-empty">No active sessions. Registered sessions appear here while they are working — bound to a plan-tree node or not — and none are active right now.</p>
    {{end}}
  </section>{{end}}

{{define "roster-style"}}    .roster-panel { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 0.5rem; padding: 0.75rem 1rem; }
    .roster-title { margin: 0 0 0.5rem 0; font-size: 0.95rem; font-weight: 600; color: #0f172a; }
    .roster-list { list-style: none; margin: 0; padding: 0; }
    .roster-entry { display: flex; align-items: baseline; gap: 0.5rem; padding: 0.2rem 0; font-size: 0.9em; }
    .roster-slug { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; color: #111827; word-break: break-all; }
    .roster-tag { margin-left: auto; flex: none; padding: 0.05rem 0.4rem; border-radius: 0.25rem; font-size: 0.75em; font-weight: 500; }
    .roster-bound .roster-tag { background: #d1fae5; color: #065f46; }
    .roster-unbound .roster-tag { background: #fee2e2; color: #991b1b; }
    .roster-empty { margin: 0; color: #6b7280; font-style: italic; font-size: 0.9em; }{{end}}`

func init() {
	template.Must(indexTmpl.Parse(rosterTemplates))
}
