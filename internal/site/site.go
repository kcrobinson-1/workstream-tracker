// Package site renders the workstream-tracker website. See
// design/v0.1-design.md Section 7 for the visualization scope.
//
// v0.0 serves a placeholder index page. v0.1 introduces the
// folder-walk + frontmatter parse + plan-tree forest render via
// templ components.
package site

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Router returns a chi.Router serving the workstream-tracker website.
func Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/", index)
	return r
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, indexHTML)
}

const indexHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>workstream-tracker</title>
</head>
<body>
  <h1>workstream-tracker</h1>
  <p>v0.0 bones — visualization not yet implemented.</p>
</body>
</html>
`
