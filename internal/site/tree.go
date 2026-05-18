package site

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kcrobinson-1/workstream-tracker/internal/slugs"
)

// PlanNode is one node in the rendered plan tree. The tree is
// built by joining frontmatter-derived facts (slug, Status,
// descriptions) with runtime work-instance state queried from the
// DB.
type PlanNode struct {
	Slug             string
	Status           string
	NodeType         string // root | milestone | task | phase
	Label            string // computed display string; see buildLabel
	ShortDescription string
	LongDescription  string
	RelatedPRs       []string
	Children         []*PlanNode
	WorkInstances    []*ActiveWorkInstance

	// ActiveInSubtree is true iff this node, or any descendant,
	// has an active work-instance. Computed post-order in
	// buildTree after children are wired; the forest template
	// reads it to decide the default expand state of a
	// collapsible box (m2 t2 C3). Recomputed every request from
	// live data — no persistence.
	ActiveInSubtree bool
}

// ActiveWorkInstance is one currently-active work-instance
// attached to a plan-tree node. Marker-shaped data only —
// terminal-state work-instances are filtered out before render.
type ActiveWorkInstance struct {
	// ID is work_instances.id (the generated PRIMARY KEY) — the
	// key events.work_instance_id references. Carried for the
	// roster's p2 event-log join (loadActiveWorkInstances selects
	// it; buildRoster joins reported metadata on it). The forest
	// node-header renders only Actor, so this additive, unrendered
	// field cannot regress the v0.1 forest actor markers.
	ID    string
	Actor string
}

// buildLabel computes a node's display label per the t3 grammar:
//
//   - Descendant with short_description:
//     "<Type> <ordinal>: <short_description>".
//   - Root with short_description: the short_description alone (a
//     root has no position segment).
//   - Any node without short_description: the slug suffix — the
//     slug text after the root, or the full slug for a root.
//
// parsed/parseErr are buildTree's slugs.Parse result for d.Slug;
// an unparseable slug is treated like a root (no ordinal).
func buildLabel(d parsedDoc, root string, parsed slugs.Slug, parseErr error) string {
	if d.ShortDescription == "" {
		if d.Slug == root {
			return d.Slug
		}
		return strings.TrimPrefix(d.Slug, root+"-")
	}
	if parseErr != nil {
		return d.ShortDescription
	}
	pos, ok := parsed.Position()
	if !ok {
		return d.ShortDescription
	}
	return fmt.Sprintf("%s %d: %s", nodeTypeDisplay(parsed.NodeType()), pos, d.ShortDescription)
}

// nodeTypeDisplay maps a slug node type to its title-case display
// word for labels (milestone -> "Milestone", task -> "Task",
// phase -> "Phase", epic -> "Epic", root -> "Root").
func nodeTypeDisplay(nt slugs.NodeType) string {
	switch nt {
	case slugs.NodeTypeEpic:
		return "Epic"
	case slugs.NodeTypeMilestone:
		return "Milestone"
	case slugs.NodeTypeTask:
		return "Task"
	case slugs.NodeTypePhase:
		return "Phase"
	default:
		return "Root"
	}
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
			Slug:             d.Slug,
			Status:           d.Status,
			NodeType:         nodeType,
			Label:            buildLabel(d, root, parsed, err),
			ShortDescription: d.ShortDescription,
			LongDescription:  d.LongDescription,
			RelatedPRs:       d.RelatedPRs,
			WorkInstances:    active[d.Slug],
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

	// Post-order pass (after children are wired): compute the C3
	// default-open signal — a node is active-in-subtree iff it,
	// or any descendant, has an active work-instance. The
	// "or any descendant" clause is what forces an active leaf's
	// ancestor boxes open so the leaf is visible.
	for _, r := range roots {
		markActiveInSubtree(r)
	}

	return roots
}

// markActiveInSubtree sets n.ActiveInSubtree post-order and
// returns it, so a parent's value reflects already-computed
// children.
func markActiveInSubtree(n *PlanNode) bool {
	active := len(n.WorkInstances) > 0
	for _, c := range n.Children {
		if markActiveInSubtree(c) {
			active = true
		}
	}
	n.ActiveInSubtree = active
	return active
}
