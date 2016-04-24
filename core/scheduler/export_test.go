package scheduler

import "github.com/mindfork/mindfork/core/message"

// CheckCycle exports checkCycle for testing.
func CheckCycle(graph map[int64]message.Intention, from, to int64) []int64 {
	return checkCycle(graph, from, to)
}
