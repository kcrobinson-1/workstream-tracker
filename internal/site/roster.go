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
// rule; the actor stays loader-internal).
//
// p3 F9 K3 (C3): every roster entry — bound or unbound, with
// or without reported metadata — opens to the same K3-shape
// disclosure: a known-facts header (slug, actor id, registered,
// last event) above either the raw-JSON block (when metadata
// was reported) or the no-metadata sentinel (when none was).
// The same outer disclosure structure across every entry is
// the K3 visual goal. The actor id surfaces only inside the K3
// header as a deliberately-labeled facts-block field; the
// .roster-label element continues to render the name-then-slug
// fallback unchanged, so the wst-<uuid> actor never re-leaks
// into a roster label (the t4 + p2 identity-rendering rule).
// The disclosure is a read affordance only — no promote /
// dismiss / attach action.
const rosterTemplates = `
{{define "roster"}}<section class="roster-panel">
    <h2 class="roster-title">Sessions</h2>
    {{if .Roster}}
    <ul class="roster-list">
      {{range .Roster}}<li class="roster-entry roster-{{if .Bound}}bound{{else}}unbound{{end}}"><details class="roster-disclosure"><summary class="roster-summary"><span class="roster-label{{if not .Name}} roster-slug{{end}}">{{if .Name}}{{.Name}}{{else}}{{.Slug}}{{end}}</span><span class="roster-tag">{{if .Bound}}bound{{else}}unbound{{end}}</span></summary><div class="roster-body"><dl class="roster-facts"><dt>Slug:</dt><dd class="roster-fact-slug">{{.Slug}}</dd><dt>Actor id:</dt><dd class="roster-fact-actor">{{.Actor}}</dd><dt>Registered:</dt><dd>{{formatEventTime .RegisteredAt}}</dd><dt>Last event:</dt><dd>{{if .LastEventAt}}{{formatEventTime .LastEventAt}}{{else}}{{formatEventTime .RegisteredAt}}{{end}}</dd></dl>{{if .Detail}}<pre class="roster-detail">{{.Detail}}</pre>{{else}}<p class="roster-empty-meta">(no reported metadata)</p>{{end}}</div></details></li>
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
    .roster-empty { margin: 0; color: #6b7280; font-style: italic; font-size: 0.9em; }
    /* p3 F9 K3: the disclosed body. The known-facts header is
       a two-column dl (label / value); the raw-JSON block or
       no-metadata sentinel sits below. Every entry across all
       four observable states (a)/(b)/(c)/(d) renders the same
       outer structure. */
    .roster-body { margin: 0.3rem 0 0.2rem 0; }
    .roster-facts { display: grid; grid-template-columns: max-content 1fr; gap: 0.15rem 0.5rem; margin: 0 0 0.4rem 0; font-size: 0.85em; }
    .roster-facts dt { color: #6b7280; font-weight: 500; }
    .roster-facts dd { margin: 0; color: #111827; word-break: break-all; }
    .roster-fact-slug { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }
    .roster-fact-actor { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; color: #6b7280; }
    .roster-empty-meta { margin: 0.4rem 0 0.2rem 0; color: #6b7280; font-style: italic; font-size: 0.85em; }{{end}}`

func init() {
	template.Must(indexTmpl.Parse(rosterTemplates))
}
