// Package slugs defines the plan-doc slug grammar shared by the
// API and site packages.
//
// A root slug is kebab-case descriptive text such as
// "madrona-feedback"; no kebab-delimited token of a root may be a
// bare position segment ("m1" is never a root). A descendant
// appends position segments to the root in strict, contiguous
// order — at most one "-mN" (milestone), then at most one "-tN"
// (task), then at most one "-pN" (phase). For example,
// "madrona-feedback-m1-t2-p1" is phase 1 under task 2 under
// milestone 1 of the "madrona-feedback" root.
package slugs

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// NodeType is a plan-doc node kind encoded in or associated with
// a slug.
type NodeType string

const (
	NodeTypeRoot      NodeType = "root"
	NodeTypeEpic      NodeType = "epic"
	NodeTypeMilestone NodeType = "milestone"
	NodeTypeTask      NodeType = "task"
	NodeTypePhase     NodeType = "phase"
)

var (
	rootPattern    = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	segmentPattern = regexp.MustCompile(`^[mtp]\d+$`)
)

// Slug is a parsed plan-doc slug relative to a known root slug.
type Slug struct {
	raw      string
	root     string
	segments []Segment
}

// Segment is one descendant position segment, such as "m1" or
// "t2".
type Segment struct {
	Raw      string
	NodeType NodeType
	Position int
}

// IsValidRoot reports whether slug is valid kebab-case root text.
// Root slugs are descriptive names: kebab-case (lowercase letters,
// digits, hyphens) and no kebab-delimited token may be a bare
// position segment (`mN`/`tN`/`pN`), so "m1" and "m1-foo" are not
// valid roots. This matches what IsStructuralRoot already assumes
// and the IsWellFormed root-portion grammar.
func IsValidRoot(slug string) bool {
	if !rootPattern.MatchString(slug) {
		return false
	}
	for _, token := range strings.Split(slug, "-") {
		if segmentPattern.MatchString(token) {
			return false
		}
	}
	return true
}

var wordPattern = regexp.MustCompile(`^[a-z0-9]+$`)

// IsWellFormed reports whether slug is a grammatically valid
// plan-doc slug standalone (without a known root): a kebab-case root
// (no token of which is a bare position segment), optionally
// followed by ordered position segments — at most one "mN", then at
// most one "tN", then at most one "pN", in that order and nothing
// else. It does not verify the slug names a real plan-tree doc — the
// server is repo-blind — only that it parses under the grammar. Used
// by the create-or-attach exact-slug flow, which honors a
// caller-supplied slug verbatim.
func IsWellFormed(slug string) bool {
	if slug == "" {
		return false
	}
	parts := strings.Split(slug, "-")

	// Consume the root: one or more kebab words, none of which may
	// be a bare position segment.
	i := 0
	for i < len(parts) && wordPattern.MatchString(parts[i]) && !segmentPattern.MatchString(parts[i]) {
		i++
	}
	if i == 0 {
		// No valid root word before the first position segment.
		return false
	}

	// Consume position segments as a contiguous prefix of the
	// sequence [m, t, p], at most one of each and in that exact
	// order. A slug encodes a full path from the root, so "tN"
	// requires a preceding "mN" and "pN" requires a preceding "tN"
	// — "root-p1" and "root-t1" are not well-formed.
	order := []byte{'m', 't', 'p'}
	oi := 0
	for ; i < len(parts); i++ {
		if !segmentPattern.MatchString(parts[i]) {
			return false
		}
		if oi >= len(order) || parts[i][0] != order[oi] {
			// Out of order, repeated, or non-contiguous.
			return false
		}
		oi++
	}
	return true
}

// ValidNodeType reports whether nodeType is a recognized API node
// type.
func ValidNodeType(nodeType string) bool {
	switch NodeType(nodeType) {
	case NodeTypeEpic, NodeTypeMilestone, NodeTypeTask, NodeTypePhase:
		return true
	}
	return false
}

// DescendantLetter returns the slug segment letter for a
// descendant node type. Epics and roots do not have descendant
// letters.
func DescendantLetter(nodeType string) (string, bool) {
	switch NodeType(nodeType) {
	case NodeTypeMilestone:
		return "m", true
	case NodeTypeTask:
		return "t", true
	case NodeTypePhase:
		return "p", true
	}
	return "", false
}

// DescendantPrefix returns the slug prefix used to allocate a
// direct child of parentPath under rootSlug. parentPath is empty
// for a child of the root, or a suffix such as "m1-t2" for a
// deeper child.
func DescendantPrefix(rootSlug, parentPath, nodeType string) (string, error) {
	letter, ok := DescendantLetter(nodeType)
	if !ok {
		return "", fmt.Errorf("node_type %q is not a descendant type", nodeType)
	}

	prefix := rootSlug + "-"
	if parentPath != "" {
		prefix += parentPath + "-"
	}
	return prefix + letter, nil
}

// NextDescendant returns the next direct-child slug by scanning
// existing slugs that may share the same prefix. Deeper
// descendants are ignored, so "root-m1-t1" does not affect the
// next milestone under "root".
func NextDescendant(rootSlug, parentPath, nodeType string, existing []string) (string, error) {
	prefix, err := DescendantPrefix(rootSlug, parentPath, nodeType)
	if err != nil {
		return "", err
	}

	tail := regexp.MustCompile(`^` + regexp.QuoteMeta(prefix) + `(\d+)$`)
	maxN := 0
	for _, slug := range existing {
		m := tail.FindStringSubmatch(slug)
		if m == nil {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		if n > maxN {
			maxN = n
		}
	}

	return prefix + strconv.Itoa(maxN+1), nil
}

// Parse parses raw as either root or descendant slug relative to
// rootSlug.
func Parse(raw, rootSlug string) (Slug, error) {
	if rootSlug == "" {
		return Slug{}, fmt.Errorf("root slug is required")
	}
	if raw == rootSlug {
		return Slug{raw: raw, root: rootSlug}, nil
	}

	prefix := rootSlug + "-"
	if !strings.HasPrefix(raw, prefix) {
		return Slug{}, fmt.Errorf("slug %q is not under root %q", raw, rootSlug)
	}

	parts := strings.Split(strings.TrimPrefix(raw, prefix), "-")
	segments := make([]Segment, 0, len(parts))
	for _, part := range parts {
		segment, err := parseSegment(part)
		if err != nil {
			return Slug{}, err
		}
		segments = append(segments, segment)
	}

	return Slug{raw: raw, root: rootSlug, segments: segments}, nil
}

// IsStructuralRoot reports whether slug contains no descendant
// position segments. This is used when inferring roots from a
// plan-doc forest; it is intentionally distinct from IsValidRoot,
// which validates the API's root slug text.
func IsStructuralRoot(slug string) bool {
	for _, segment := range strings.Split(slug, "-") {
		if segmentPattern.MatchString(segment) {
			return false
		}
	}
	return true
}

// FindRoot returns the known root slug that prefixes raw. If raw
// itself is a known root, raw is returned. If no known root
// matches, raw is returned as a conservative fallback.
func FindRoot(raw string, rootSlugs map[string]bool) string {
	if rootSlugs[raw] {
		return raw
	}
	parts := strings.Split(raw, "-")
	for i := len(parts) - 1; i > 0; i-- {
		candidate := strings.Join(parts[:i], "-")
		if rootSlugs[candidate] {
			return candidate
		}
	}
	return raw
}

// String returns the original slug text.
func (s Slug) String() string {
	return s.raw
}

// Root returns the root slug this slug was parsed against.
func (s Slug) Root() string {
	return s.root
}

// NodeType returns root for a root slug, otherwise the type
// encoded by the final descendant segment.
func (s Slug) NodeType() NodeType {
	if len(s.segments) == 0 {
		return NodeTypeRoot
	}
	return s.segments[len(s.segments)-1].NodeType
}

// Position returns the terminal descendant segment's position and
// true; for a root slug (no descendant segments) it returns 0 and
// false. Mirrors the NodeType accessor shape so callers can read
// the ordinal without re-implementing the slug grammar.
func (s Slug) Position() (int, bool) {
	if len(s.segments) == 0 {
		return 0, false
	}
	return s.segments[len(s.segments)-1].Position, true
}

// Parent returns the parent slug for a descendant, or an empty
// string for a root slug.
func (s Slug) Parent() string {
	if len(s.segments) == 0 {
		return ""
	}
	if len(s.segments) == 1 {
		return s.root
	}

	parts := make([]string, 0, len(s.segments)-1)
	for _, segment := range s.segments[:len(s.segments)-1] {
		parts = append(parts, segment.Raw)
	}
	return s.root + "-" + strings.Join(parts, "-")
}

func parseSegment(raw string) (Segment, error) {
	if !segmentPattern.MatchString(raw) {
		return Segment{}, fmt.Errorf("invalid slug segment %q", raw)
	}
	position, err := strconv.Atoi(raw[1:])
	if err != nil {
		return Segment{}, fmt.Errorf("parse slug segment %q: %w", raw, err)
	}

	nodeType := NodeTypeRoot
	switch raw[0] {
	case 'm':
		nodeType = NodeTypeMilestone
	case 't':
		nodeType = NodeTypeTask
	case 'p':
		nodeType = NodeTypePhase
	}

	return Segment{Raw: raw, NodeType: nodeType, Position: position}, nil
}
