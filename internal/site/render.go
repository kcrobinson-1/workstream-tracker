package site

import (
	"html/template"
	"io"
	"net/url"
	"strings"
)

// indexTmpl is the page shell: chrome, the two-region .layout
// container, and the composition of the forest and roster region
// templates. The region bodies + their CSS live in their own
// files (forest.go is t2's surface, roster.go is t4's) and are
// parsed into this same template tree by their package init()s,
// so no later task edits this shell to grow a region. Bare-bones
// by design (see design/v0.1-design.md Section 7).
var indexTmpl = template.Must(template.New("index").Funcs(template.FuncMap{
	"statusClass":      statusClass,
	"isURL":            isAbsoluteURL,
	"truncateLongDesc": truncateLongDesc,
}).Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>workstream-tracker</title>
  <style>
    body { font-family: system-ui, -apple-system, sans-serif; padding: 1rem; max-width: 80rem; margin: 0 auto; color: #111827; background-color: #fff; color-scheme: light; }
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
{{template "forest-style"}}
{{template "roster-style"}}
  </style>
</head>
<body>
  <h1>workstream-tracker</h1>
  <div class="layout">
  <main class="forest">
  {{template "forest" .}}
  </main>
  <aside class="roster">
  {{template "roster" .}}
  </aside>
  </div>
</body>
</html>
`))

// indexData is the template input passed to indexTmpl.
type indexData struct {
	Roots     []*PlanNode
	PlansPath string
	Roster    []RosterEntry
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

// maxLongDescLines is the hard cap on how many lines of a node's
// long description the forest renders, even when the box is
// expanded. Plan docs are whole markdown files; the forest is a
// context-and-goal overview, not a document viewer. Read the plan
// file itself for the full text. Tunable — bump it if the cap
// hides too much of the goal.
const maxLongDescLines = 40

// truncateLongDesc caps s at maxLongDescLines lines. A capped body
// gets a trailing marker line so the truncation is visible rather
// than silently dropping the tail. Returns a plain string, so
// html/template still contextually escapes it (no raw HTML).
func truncateLongDesc(s string) string {
	lines := strings.Split(s, "\n")
	if len(lines) <= maxLongDescLines {
		return s
	}
	kept := strings.Join(lines[:maxLongDescLines], "\n")
	return kept + "\n\n… (truncated — see the plan doc for the full text)"
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
