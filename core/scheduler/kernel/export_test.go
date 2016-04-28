package kernel

import "github.com/mindfork/mindfork/core/message"

type Node struct {
	ParentBounties int
	// Parents        []int64
	// childCosts       int
	// parentUrgency    int
	// parentImportance int
	// visited          bool
}

// Nodes exports k.intentions for testing.
func Nodes(k *Kernel) map[int64]Node {
	nodes := make(map[int64]Node)
	for id, n := range k.intentions {
		newNode := Node{
			ParentBounties: n.parentBounties,
		}
		nodes[id] = newNode
	}
	return nodes
}

// Roots exports k.roots for testing.
func Roots(k *Kernel) map[int64]message.Intention {
	roots := make(map[int64]message.Intention)
	for id, i := range k.roots {
		roots[id] = i.Intention
	}

	return roots
}

// CheckCycle exports checkCycle for testing.
func CheckCycle(graph map[int64]message.Intention, from, to int64) []int64 {
	nodes := make(map[int64]node)
	for id, i := range graph {
		nodes[id] = node{Intention: i}
	}

	return checkCycle(nodes, from, to)
}
