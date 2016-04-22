package core

import (
	"fmt"
	"sync"

	"github.com/mindfork/mindfork/core/message"
	mfm "github.com/mindfork/mindfork/message"

	"github.com/gonum/graph/simple"
)

// Kernel is the core Scheduler implementation.  It uses github.com/gonum/graph
// to implement a directed acyclic graph of Intentions by dependency.
type Kernel struct {
	sync.RWMutex
	intentions *simple.DirectedGraph
}

// NewKernel makes a new Kernel with a populated Intention graph.
func NewKernel() *Kernel {
	return &Kernel{intentions: simple.NewDirectedGraph(0, 0)}
}

// Intention implements gonum/graph.Node for message.Intention.
type Intention struct {
	message.Intention
}

// ID implements Node.ID for Intention.
func (i Intention) ID() int {
	return int(i.Intention.ID)
}

// Add implements Scheduler.Add on Kernel.
func (k *Kernel) Add(i message.Intention) mfm.Message {
	k.Lock()
	defer k.Unlock()

	ints := k.intentions

	i.ID = int64(ints.NewNodeID())
	node := Intention{i}

	ints.AddNode(node)
	for _, dep := range i.Deps {
		if to := ints.Node(int(dep)); to != nil {
			ints.SetEdge(simple.Edge{F: node, T: to, W: 1})
		} else {
			return message.Error{Err: fmt.Errorf("no such intention %d", dep)}
		}
	}

	return i
}

// Peek implements Scheduler.Peek on Kernel.
func (k *Kernel) Peek() []message.Intention {
	return nil
}

// Export implements Scheduler.Export on Kernel.
func (k *Kernel) Export() []message.Intention {
	k.RLock()
	defer k.RUnlock()

	result := make([]message.Intention, len(k.intentions.Nodes()))
	for i, in := range k.intentions.Nodes() {
		result[i] = in.(Intention).Intention
	}

	return result
}
