package kernel

import (
	"sort"

	"github.com/mindfork/mindfork/core/message"
	"github.com/mindfork/mindfork/core/scheduler"
)

// Available implements Scheduler.Available on Kernel.
func (k *Kernel) Available(ord scheduler.Ordering) []message.Intention {
	k.RLock()
	toSort := make([]node, len(k.free))
	result := make([]message.Intention, len(k.free))
	i := 0
	for _, in := range k.free {
		toSort[i] = in
		i++
	}
	k.RUnlock()

	sort.Sort(struct {
		nodes
		lesser
	}{toSort, whichOrd(toSort, ord)})

	// Copy sorted Intentions into Intention slice.
	for i, n := range toSort {
		result[i] = n.Intention
	}

	return result
}
