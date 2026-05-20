package site

import (
	"bytes"
	"html/template"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// indexTmpl is the page shell: chrome, the two-region .layout
// container, and the composition of the forest and roster region
// templates. The region bodies + their CSS live in their own
// files (forest.go is t2's surface, roster.go is t4's) and are
// parsed into this same template tree by their package init()s,
// so no later task edits this shell to grow a region. Bare-bones
// by design (see design/v0.1-design.md Section 7).
var indexTmpl = template.Must(template.New("index").Funcs(template.FuncMap{
	"statusClass":        statusClass,
	"isURL":              isAbsoluteURL,
	"truncateLongDesc":   truncateLongDesc,
	"progressCellClass":  progressCellClass,
	"renderLongDescBody": renderLongDescBody,
	"formatEventTime":    formatEventTime,
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

// canonicalStatus returns the canonical-prefix portion of status —
// the part before " — " for a freeform-suffixed value like
// "Deferred — <reason>", or the whole string otherwise. m2 t1 C4:
// the single source of truth for the Status canonical-prefix
// strip; statusClass, progressCellClass, and PlanNode.Affordance
// all consume it so a Deferred — <reason> node reads as Deferred
// at every site. No second em-dash split lives in internal/site/.
func canonicalStatus(status string) string {
	if i := strings.Index(status, " — "); i >= 0 {
		return status[:i]
	}
	return status
}

// statusClass converts a Status string to a CSS class fragment.
// Recognized values from spec/planning/shared.md "Plan-doc
// Status" map to specific classes; unrecognized values fall back
// to "unknown" per the spec's graceful-fallback rule.
func statusClass(status string) string {
	switch canonicalStatus(status) {
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

// progressCellClass returns the per-cell CSS class suffix the
// default D / P / I / V progress row uses for a node with the
// given Status at the given cell position ("d" / "p" / "i" /
// "v"). Three shape buckets per p3 F3a:
//   - Landed: all four cells "landed" (filled green; matches
//     .status-landed badge palette).
//   - In draft: D = "in-draft" (filled amber); P/I/V = "empty"
//     (dashed border, no fill — the placeholder treatment that
//     reads "stage not yet started" while preserving the
//     per-cell DOM C-INV-1 anchor for F3b's future attachment).
//   - Everything else (In progress, Proposed, Validating,
//     Deferred, any unrecognized Status the statusClass helper
//     falls back to "unknown" for): all four "neutral" (filled
//     grey; matches .status-unknown badge palette).
//
// position is "d" only-distinguished from the others (the
// In-draft branch flips D to filled and the rest to empty);
// passing positions other than those four returns the same
// non-"d" treatment per bucket.
func progressCellClass(status, position string) string {
	switch canonicalStatus(status) {
	case "Landed":
		return "landed"
	case "In draft":
		if position == "d" {
			return "in-draft"
		}
		return "empty"
	}
	return "neutral"
}

// bodyMarkdown is the goldmark instance that renders the F7
// disclosed long-description body. The custom stripLinksRenderer
// runs at priority 100 (lower than the default html.NewRenderer
// at 1000), which under goldmark's "higher-priority registers
// first, lower-priority registers last and overwrites" pattern
// (renderer.Render's init loop, html.go line 130+) overrides
// the default Link / AutoLink renderers so anchor tags are
// suppressed structurally — never reach the rendered HTML.
// OD3 = I2b: render markdown with <a> tags stripped (the link
// text still renders via the Walk-into-children pass on Link
// nodes; AutoLink nodes emit their label as plain escaped text).
var bodyMarkdown = goldmark.New(
	goldmark.WithRendererOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&stripLinksRenderer{}, 100),
		),
	),
)

// stripLinksRenderer overrides KindLink + KindAutoLink so the
// rendered HTML carries no <a> tags. Link text still appears as
// plain text (Walk continues into children for KindLink; the
// AutoLink's label is emitted as escaped plain text). All other
// markdown constructs (paragraphs, emphasis, code spans,
// headers, lists, etc.) render through the default html
// renderer unchanged.
type stripLinksRenderer struct{}

func (r *stripLinksRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindLink, renderLinkStripped)
	reg.Register(ast.KindAutoLink, renderAutoLinkStripped)
}

// renderLinkStripped emits nothing for the link wrapper; the
// Walk continues into the link's text children, which the
// default Text/String renderers handle, so the link's display
// text appears as plain text in the rendered output.
func renderLinkStripped(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

// renderAutoLinkStripped emits the autolink's label (the URL
// itself for a bare <http://...>) as escaped plain text on the
// entering pass, then walks children (typically none for an
// AutoLink). Without the strip, the default renderer would emit
// <a href="...">label</a>; here the label survives as plain
// escaped text.
func renderAutoLinkStripped(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	al, ok := n.(*ast.AutoLink)
	if !ok {
		return ast.WalkContinue, nil
	}
	_, _ = w.Write(util.EscapeHTML(al.Label(source)))
	return ast.WalkContinue, nil
}

// renderLongDescBody renders the per-node long description as
// markdown HTML with <a> tags stripped per F7 OD3 = I2b. The
// truncateLongDesc cap runs first (truncate, then render) so a
// very long body still terminates with the existing truncation
// marker even when disclosed — the cap is the defense-in-depth
// tail per parent C2. The empty case short-circuits to an empty
// HTML value so the template can use a non-empty check
// unchanged. Render failures fall back to escaped plain text
// rather than failing the page render.
func renderLongDescBody(s string) template.HTML {
	if s == "" {
		return ""
	}
	truncated := truncateLongDesc(s)
	var buf bytes.Buffer
	if err := bodyMarkdown.Convert([]byte(truncated), &buf); err != nil {
		return template.HTML(template.HTMLEscapeString(truncated))
	}
	return template.HTML(buf.String())
}

// formatEventTime formats a Unix-epoch nanosecond timestamp for
// the F9 K3 known-facts header — the storage shape every event
// write in internal/api/handlers.go writes (now.UnixNano()),
// matching events.received_at's column. Zero (no event observed)
// renders as the em-dash placeholder. Otherwise UTC,
// RFC3339-shaped without the T separator so it reads cleanly to
// a human at a glance: "2026-05-20 14:33:21 UTC".
func formatEventTime(t int64) string {
	if t == 0 {
		return "—"
	}
	return time.Unix(0, t).UTC().Format("2006-01-02 15:04:05 UTC")
}
