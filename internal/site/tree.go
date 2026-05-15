package site

import (
	"sort"

	"github.com/kcrobinson-1/workstream-tracker/internal/slugs"
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
		if slugs.IsStructuralRoot(d.Slug) {
			rootSlugs[d.Slug] = true
		}
	}

	// Build nodes, indexed by slug.
	nodesBySlug := map[string]*PlanNode{}
	for _, d := range docs {
		root := slugs.FindRoot(d.Slug, rootSlugs)
		parsed, err := slugs.Parse(d.Slug, root)
		nodeType := string(slugs.NodeTypeRoot)
		if err == nil {
			nodeType = string(parsed.NodeType())
		}
		nodesBySlug[d.Slug] = &PlanNode{
			Slug:          d.Slug,
			Status:        d.Status,
			NodeType:      nodeType,
			WorkInstances: active[d.Slug],
		}
	}

	// Wire children to parents.
	var roots []*PlanNode
	for _, d := range docs {
		node := nodesBySlug[d.Slug]
		root := slugs.FindRoot(d.Slug, rootSlugs)
		if d.Slug == root {
			roots = append(roots, node)
			continue
		}
		parsed, err := slugs.Parse(d.Slug, root)
		if err != nil {
			continue
		}
		parent := parsed.Parent()
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
