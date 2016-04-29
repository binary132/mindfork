package kernel

import (
	"sync"

	"github.com/mindfork/mindfork/core/message"
)

type node struct {
	message.Intention

	parents        []int64
	parentBounties int
	// childCosts       int
	// parentUrgency    int
	// parentImportance int
	// visited          bool
}

func (n node) copy() node {
	newParents := make([]int64, len(n.parents))
	copy(newParents, n.parents)

	return node{
		Intention:      n.Intention,
		parents:        newParents,
		parentBounties: n.parentBounties,
	}
}

// Kernel is the core Scheduler implementation.  It holds Intentions in volatile
// memory.
type Kernel struct {
	sync.RWMutex

	roots      map[int64]node
	intentions map[int64]node
	free       map[int64]node

	nextID int64
}

// New sets up a new Kernel.
func New() *Kernel {
	return &Kernel{
		roots:      make(map[int64]node),
		intentions: make(map[int64]node),
		free:       make(map[int64]node),
	}
}
