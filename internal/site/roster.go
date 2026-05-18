package site

import "html/template"

// rosterTemplates owns the roster region: the "roster" template
// (the bound/unbound session roster — every active work-instance
// listed, each marked bound when its slug matches a walked plan
// doc or unbound otherwise, with a deliberate "no active
// sessions" state) and "roster-style" (the roster CSS the shell
// injects into <head>). This is t4's owned region body.
//
// m2 t4 p2 enrichment: the entry label is the reported `name`
// when the session reported one, else the work-instance slug —
// never the wst-<uuid> actor (the task-level name-then-slug
// rule; the actor stays loader-internal). An entry whose session
// reported metadata is expandable via native HTML
// <details>/<summary> (the repo's no-JS idiom) to a
// deliberately-unstructured raw-JSON view of that resolved
// metadata; an entry whose session reported nothing renders as a
// plain row with no disclosure control (additive — absence is
// not a drop), mirroring the forest leaf-box-without-children
// idiom. The roster still lists only — no promote/dismiss/attach
// affordance (the disclosure is a read affordance, not triage).
const rosterTemplates = `
{{define "roster"}}<section class="roster-panel">
    <h2 class="roster-title">Sessions</h2>
    {{if .Roster}}
    <ul class="roster-list">
      {{range .Roster}}<li class="roster-entry roster-{{if .Bound}}bound{{else}}unbound{{end}}">{{if .Detail}}<details class="roster-disclosure"><summary class="roster-summary"><span class="roster-label{{if not .Name}} roster-slug{{end}}">{{if .Name}}{{.Name}}{{else}}{{.Slug}}{{end}}</span><span class="roster-tag">{{if .Bound}}bound{{else}}unbound{{end}}</span></summary><pre class="roster-detail">{{.Detail}}</pre></details>{{else}}<span class="roster-label{{if not .Name}} roster-slug{{end}}">{{if .Name}}{{.Name}}{{else}}{{.Slug}}{{end}}</span><span class="roster-tag">{{if .Bound}}bound{{else}}unbound{{end}}</span>{{end}}</li>
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
    .roster-disclosure { flex: 1 1 0; min-width: 0; }
    .roster-summary { display: flex; align-items: baseline; gap: 0.5rem; cursor: pointer; list-style: revert; }
    .roster-label { color: #111827; word-break: break-all; }
    .roster-slug { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }
    .roster-tag { margin-left: auto; flex: none; padding: 0.05rem 0.4rem; border-radius: 0.25rem; font-size: 0.75em; font-weight: 500; }
    .roster-bound .roster-tag { background: #d1fae5; color: #065f46; }
    .roster-unbound .roster-tag { background: #fee2e2; color: #991b1b; }
    .roster-detail { margin: 0.4rem 0 0.2rem 0; padding: 0.5rem; background: #ffffff; border: 1px solid #e2e8f0; border-radius: 0.25rem; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 0.8em; white-space: pre-wrap; word-break: break-all; overflow-x: auto; }
    .roster-empty { margin: 0; color: #6b7280; font-style: italic; font-size: 0.9em; }{{end}}`

func init() {
	template.Must(indexTmpl.Parse(rosterTemplates))
}
