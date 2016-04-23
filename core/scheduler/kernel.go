package scheduler

import (
	"fmt"
	"sync"

	"github.com/mindfork/mindfork/core/message"
	mfm "github.com/mindfork/mindfork/message"

	"github.com/gonum/graph"
	"github.com/gonum/graph/simple"
	"github.com/gonum/graph/traverse"
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
	depNodes := make([]Intention, len(i.Deps))

	for index, dep := range i.Deps {
		if dep == i.ID {
			return message.Error{Err: fmt.Errorf("cannot depend on self")}
		} else if to := ints.Node(int(dep)); to != nil {
			depNodes[index] = to.(Intention)
		} else {
			return message.Error{Err: fmt.Errorf("no such intention %d", dep)}
		}
	}

	i.ID = int64(ints.NewNodeID())
	node := Intention{i}

	ints.AddNode(node)
	for _, dep := range depNodes {
		ints.SetEdge(simple.Edge{F: node, T: dep, W: 0})
	}

	return i
}

// Available implements Scheduler.Available on Kernel.
func (k *Kernel) Available() []message.Intention {
	k.RLock()
	defer k.RUnlock()

	ints := k.intentions

	var result []message.Intention

	w := traverse.BreadthFirst{EdgeFilter: func(e graph.Edge) bool {
		return len(e.To().(Intention).Deps) > 0 &&
			len(e.From().(Intention).Deps) > 0
	}}

	w.WalkAll(graph.Undirect{G: ints}, nil, nil, func(n graph.Node) {
		i := n.(Intention)
		result = append(result, i.Intention)
	})

	return result
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
