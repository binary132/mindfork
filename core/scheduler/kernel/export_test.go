package kernel

import "github.com/mindfork/mindfork/core/message"

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
