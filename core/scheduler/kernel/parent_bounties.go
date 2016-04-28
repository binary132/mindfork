package kernel

func recalculateParentBounties(graph map[int64]node, ids ...int64) {
	var (
		seen = make(map[int64]bool)

		// We will explore downward from the given node to find
		// everything that needs to be recalculated.
		curr []int64
		next []int64 = ids

		// Bases are the nodes that need to be recalculated upward from.
		bases []int64
	)

	// Look for the base nodes in the given ID's network.
	for len(next) != 0 {
		curr, next = next, nil

		for _, id := range curr {
			node := graph[id]

			// If this node has no deps, add it to bases.
			if len(node.Deps) == 0 {
				bases = append(bases, id)
			}

			// Explore downward.  If we find something we haven't
			// seen yet, explore that next.
			for _, dep := range node.Deps {
				if !seen[dep] {
					next = append(next, dep)
				}
			}

			seen[id] = true
		}
	}

	updateAncestors(graph, bases...)
}

// updateAncestors calculates blockedBounty of each node.  It returns a slice of
// the ancestor map for each given ID.
func updateAncestors(graph map[int64]node, ids ...int64) []map[int64]bool {
	toReturn := make([]map[int64]bool, len(ids))

	//   - Find the union of the sets of all ancestors of its parents.
	//     > for each base:
	//     > if I have parents, find and update their ancestors first.
	//       + my ancestors = union of their ancestors and them.
	//       + my blocked bounty = the sum of Bounties of this set.
	for i, id := range ids {
		node := graph[id]

		node.parentBounties = 0

		if len(node.parents) == 0 {
			toReturn[i] = nil
			graph[id] = node
			continue
		}

		ancestors := make(map[int64]bool)

		for _, p := range node.parents {
			ancestors[p] = true
		}

		ancestors = union(
			ancestors, updateAncestors(graph, node.parents...)...,
		)

		for id := range ancestors {
			node.parentBounties += graph[id].Bounty
		}

		graph[id] = node
		toReturn[i] = ancestors
	}

	return toReturn
}

// assume any overlap has identical members.
func union(g map[int64]bool, grs ...map[int64]bool) map[int64]bool {
	if len(grs) == 0 || len(grs) == 1 && grs[0] == nil {
		return g
	}

	result := make(map[int64]bool)
	for id := range g {
		result[id] = true
	}
	for _, gr := range grs {
		for id := range gr {
			result[id] = true
		}
	}

	return result
}
