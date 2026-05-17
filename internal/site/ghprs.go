package site

// gh pr list auto-discovery (t4 P2). This is the codebase's first
// subprocess shell-out. One `gh` invocation per page request
// discovers merged/open PRs; a discovered PR's URL is attached to
// a node when the node's verbatim slug is a case-sensitive
// substring of the PR title (match done in Go over the decoded
// JSON, never via `gh --search`, so no plan-derived data ever
// reaches the command line — there is no command-injection
// surface). Every failure mode degrades to "no discoveries": the
// node still renders its frontmatter PRs (or none) and the page
// never fails. See
// docs/plans/workstream-tracker-1-0/m1-t4-p2-gh-discovery.md.

import (
	"context"
	"encoding/json"
	"os/exec"
	"strings"
	"time"
)

// ghTimeout bounds the single `gh` subprocess so a hung `gh`
// cannot wedge a page render beyond this deadline (scoping P2-D4:
// generous for a normal local round-trip, bounded against a
// hang). exec.CommandContext propagates the deadline as a kill.
const ghTimeout = 5 * time.Second

// ghWaitDelay bounds how long Output() blocks for I/O after the
// context deadline kills `gh`, then force-closes the inherited
// pipes (Cmd.WaitDelay). Without it a killed `gh`'s grandchild
// holding the stdout pipe wedges the request past ghTimeout. Kept
// short: by the time it fires the deadline has already elapsed
// and the result is discarded anyway.
const ghWaitDelay = 1 * time.Second

// ghPRListLimit caps the `gh pr list -L` result set. PRs beyond
// the cap are missed; acceptable under the best-effort posture
// (scoping D4 / P2-D1).
const ghPRListLimit = "200"

// ghPR decodes one element of `gh pr list --json url,title,number`.
// The field tags match gh's JSON keys; a future gh field rename
// yields empty/garbled discovery that degrades to frontmatter-only
// rather than crashing.
type ghPR struct {
	URL    string `json:"url"`
	Title  string `json:"title"`
	Number int    `json:"number"`
}

// discoverPRsByTitle runs exactly one `gh` subprocess and returns
// the decoded PRs. The command is a fixed argument vector — no
// shell, no interpolation of any plan/slug data into the command
// line; slug matching happens in Go (see augmentRelatedPRs). Every
// failure mode (missing `gh` binary, non-zero exit — not a git
// repo / not authenticated / API error, context-deadline timeout,
// empty/malformed JSON) returns a non-nil error and a nil slice;
// it never panics and never propagates an error to the page.
func discoverPRsByTitle(ctx context.Context) ([]ghPR, error) {
	ctx, cancel := context.WithTimeout(ctx, ghTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "gh", "pr", "list",
		"--state", "all",
		"--json", "url,title,number",
		"-L", ghPRListLimit,
	)
	// On the context deadline exec.CommandContext kills the `gh`
	// process, but Output() still blocks until the stdout pipe
	// closes — and a killed `gh` can leave a grandchild holding
	// that pipe open, which would wedge the request well past
	// ghTimeout (observed in the hung-`gh` failure-matrix check).
	// WaitDelay bounds the post-kill I/O wait and force-closes the
	// inherited pipes so a hung `gh` cannot wedge a page render
	// beyond ~ghTimeout + ghWaitDelay.
	cmd.WaitDelay = ghWaitDelay

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var prs []ghPR
	if err := json.Unmarshal(out, &prs); err != nil {
		return nil, err
	}
	return prs, nil
}

// augmentRelatedPRs walks the node forest and, for each node,
// appends discovered PR URLs whose title contains the node's
// verbatim slug as a case-sensitive substring (the in-process
// exact match, immune to GitHub's hyphen tokenisation — scoping
// P2-D1). The resulting RelatedPRs is the node's frontmatter
// entries in source order (authoritative — scoping D4/P2-D3),
// followed by discovered URLs not already present, in gh result
// order. Dedupe is plain absolute-URL string equality; no
// canonicalization. Frontmatter entries are never reordered or
// rewritten. Augmentation mutates the in-request PlanNode slices
// only; nothing is persisted or cached.
func augmentRelatedPRs(roots []*PlanNode, prs []ghPR) {
	for _, n := range roots {
		if n == nil {
			continue
		}

		seen := make(map[string]bool, len(n.RelatedPRs))
		for _, u := range n.RelatedPRs {
			seen[u] = true
		}

		for _, pr := range prs {
			if !strings.Contains(pr.Title, n.Slug) {
				continue
			}
			if seen[pr.URL] {
				continue
			}
			seen[pr.URL] = true
			n.RelatedPRs = append(n.RelatedPRs, pr.URL)
		}

		augmentRelatedPRs(n.Children, prs)
	}
}
