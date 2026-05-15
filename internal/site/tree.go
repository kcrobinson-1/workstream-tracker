package site

import (
	"regexp"
	"sort"
	"strings"
)

// PlanNode is one node in the rendered plan tree. The tree is
// built by joining frontmatter-derived facts (slug, Status) with
// runtime work-instance state queried from the DB.
type PlanNode struct {
	Slug          string
	Status        string
	NodeType      string // root | milestone | task | phase
	Children      []*PlanNode
	WorkInstances []*ActiveWorkInstance
}

// ActiveWorkInstance is one currently-active work-instance
// attached to a plan-tree node. Marker-shaped data only —
// terminal-state work-instances are filtered out before render.
type ActiveWorkInstance struct {
	Actor string
}

// segmentPattern matches a position-prefix segment (m1, t2, p3,
// etc.) used in descendant slugs. See spec/planning/shared.md
// "Plan-doc identity (slug)".
var segmentPattern = regexp.MustCompile(`^[mtp]\d+$`)

// nodeTypeFromSlug returns the node type for slug given the
// known root slug. The last hyphen-separated segment determines
// the type when it matches the m/t/p position-prefix pattern;
// otherwise the slug is the root.
func nodeTypeFromSlug(slug, rootSlug string) string {
	if slug == rootSlug {
		return "root"
	}
	suffix := strings.TrimPrefix(slug, rootSlug+"-")
	parts := strings.Split(suffix, "-")
	last := parts[len(parts)-1]
	if !segmentPattern.MatchString(last) {
		return "root" // can't classify; treat as root for safety
	}
	switch last[0] {
	case 'm':
		return "milestone"
	case 't':
		return "task"
	case 'p':
		return "phase"
	}
	return "root"
}

// parentSlug returns the parent slug for a descendant slug, or
// "" if the slug is a root or the parent can't be determined.
func parentSlug(slug, rootSlug string) string {
	if slug == rootSlug {
		return ""
	}
	suffix := strings.TrimPrefix(slug, rootSlug+"-")
	parts := strings.Split(suffix, "-")
	if !segmentPattern.MatchString(parts[len(parts)-1]) {
		return ""
	}
	if len(parts) == 1 {
		return rootSlug
	}
	return rootSlug + "-" + strings.Join(parts[:len(parts)-1], "-")
}

// buildTree groups parsed docs by their root and assembles each
// root into a PlanNode tree. Returns the roots in stable
// (alphabetical-by-slug) order.
//
// Active work-instances (active map: slug → []actor) are attached
// to their corresponding nodes by exact slug match; orphaned
// work-instances (slug not present in any parsed doc) are
// dropped.
func buildTree(docs []parsedDoc, active map[string][]*ActiveWorkInstance) []*PlanNode {
	// Group docs by their root slug. The root slug is the slug
	// of any doc that is *not* a descendant of another doc — its
	// last segment doesn't match the m/t/p pattern, OR no other
	// doc's slug is a prefix.
	//
	// Simpler heuristic: a doc is a root iff its slug doesn't
	// contain a `-m\d+`, `-t\d+`, or `-p\d+` suffix segment AT
	// THE END of any prefix-truncation of itself. In practice:
	// a doc is a root if its slug, walked segment-by-segment,
	// has no segment matching the position pattern.
	//
	// Even simpler: roots are the slugs that no other slug is a
	// strict-prefix-with-hyphen of, AND whose own last segment
	// isn't a position prefix. But the layout convention puts
	// each root in its own folder, and the README.md inside is
	// the root. We use the structural property: a slug is a root
	// if *its own structure* shows no descent — no segment in
	// the slug matches m/t/p.
	rootSlugs := map[string]bool{}
	for _, d := range docs {
		if isRootSlug(d.Slug) {
			rootSlugs[d.Slug] = true
		}
	}

	// Build nodes, indexed by slug.
	nodesBySlug := map[string]*PlanNode{}
	for _, d := range docs {
		root := findRootSlugFor(d.Slug, rootSlugs)
		nodesBySlug[d.Slug] = &PlanNode{
			Slug:          d.Slug,
			Status:        d.Status,
			NodeType:      nodeTypeFromSlug(d.Slug, root),
			WorkInstances: active[d.Slug],
		}
	}

	// Wire children to parents.
	var roots []*PlanNode
	for _, d := range docs {
		node := nodesBySlug[d.Slug]
		root := findRootSlugFor(d.Slug, rootSlugs)
		if d.Slug == root {
			roots = append(roots, node)
			continue
		}
		parent := parentSlug(d.Slug, root)
		if parentNode, ok := nodesBySlug[parent]; ok {
			parentNode.Children = append(parentNode.Children, node)
		}
		// If the parent isn't present (gap in the tree), the
		// node is silently dropped from the visualization. The
		// authoring side is responsible for keeping the tree
		// consistent.
	}

	// Stable order: roots alphabetical by slug, children too.
	sort.Slice(roots, func(i, j int) bool { return roots[i].Slug < roots[j].Slug })
	for _, n := range nodesBySlug {
		sort.Slice(n.Children, func(i, j int) bool { return n.Children[i].Slug < n.Children[j].Slug })
	}

	return roots
}

// isRootSlug reports whether a slug looks like a root (no
// position-prefix segments anywhere in its structure).
//
// This works because root slugs are kebab-case and don't contain
// segments matching the m\d+ / t\d+ / p\d+ pattern at any
// boundary. False positive risk: a project with a root slug
// containing a literal "m1" segment (e.g., "model-1-plan") would
// be incorrectly classified as a descendant. The four-segment
// slug cap and noun discipline reduce this risk in practice.
func isRootSlug(slug string) bool {
	for _, seg := range strings.Split(slug, "-") {
		if segmentPattern.MatchString(seg) {
			return false
		}
	}
	return true
}

// findRootSlugFor returns the root slug whose prefix matches the
// given slug, or the slug itself if it IS a root. Falls back to
// the slug itself if no root is found in the known set.
func findRootSlugFor(slug string, rootSlugs map[string]bool) string {
	if rootSlugs[slug] {
		return slug
	}
	parts := strings.Split(slug, "-")
	for i := len(parts) - 1; i > 0; i-- {
		candidate := strings.Join(parts[:i], "-")
		if rootSlugs[candidate] {
			return candidate
		}
	}
	return slug
}
