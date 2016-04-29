package kernel

import (
	"fmt"
	"sort"
	"sync"

	"github.com/mindfork/mindfork/core/message"
	"github.com/mindfork/mindfork/core/scheduler"
	mfm "github.com/mindfork/mindfork/message"
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

type intentions []node
type byID []node
type byBounty []node
type byDate []node

// Len implements Interface.Len on defaultSort.
func (i intentions) Len() int { return len(i) }

// Swap implements Interface.Swap on defaultSort.
func (i intentions) Swap(j, k int) { i[j], i[k] = i[k], i[j] }

// Less implements Interface.Less on byID.
func (b byID) Less(i, j int) bool { return b[i].ID < b[j].ID }

// Less implements Interface.Less on byBounty.  The one with the greater Bounty
// comes first.
func (b byBounty) Less(i, j int) bool { return b[i].Bounty > b[j].Bounty }

// Less implements Interface.Less on byDate.  If both dates are nil, they will
// be compared by Value.  If only one date is nil, the other will come first.
// If both are non-nil, the sooner one will come first in the list.
func (b byDate) Less(i, j int) bool {
	wI, wJ := b[i].When, b[j].When
	switch {
	case wI == nil && wJ == nil:
		return b[i].Bounty > b[j].Bounty
	case wI == nil:
		return false
	case wI == nil:
		return true
	default:
		return wI.Before(*wJ)
	}
}

// New sets up a new Kernel.
func New() *Kernel {
	return &Kernel{
		roots:      make(map[int64]node),
		intentions: make(map[int64]node),
		free:       make(map[int64]node),
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

// Add implements Scheduler.Add on Kernel.
func (k *Kernel) Add(i message.Intention) mfm.Message {
	switch {
	case i.ID != 0:
		return k.addExisting(i)
	}

	return k.addNew(i)
}

// Available implements Scheduler.Available on Kernel.
func (k *Kernel) Available(o scheduler.Ordering) []message.Intention {
	k.RLock()
	toSort := make([]node, len(k.free))
	result := make([]message.Intention, len(k.free))
	i := 0
	for _, in := range k.free {
		toSort[i] = in
		i++
	}
	k.RUnlock()

	ord := o
	if ord == nil {
		ord = byID(toSort)
	}

	sort.Sort(struct {
		intentions
		scheduler.Ordering
	}{toSort, ord})

	// Copy sorted Intentions into Intention slice.
	for i, n := range toSort {
		result[i] = n.Intention
	}

	return result
}

func (k *Kernel) addNew(i message.Intention) mfm.Message {
	k.RLock()
	if err := k.checkNew(i); err != nil {
		k.RUnlock()
		return err
	}
	k.RUnlock()

	k.Lock()
	defer k.Unlock()

	// Make sure that we haven't mutated the graph again in the meantime...
	if err := k.checkNew(i); err != nil {
		return err
	}

	k.nextID++
	i.ID = k.nextID

	newNode, newID := node{Intention: i}, i.ID

	for _, dep := range newNode.Deps {
		// Remove the new deps from roots.
		delete(k.roots, dep)
		// Add the new node to the child nodes' parents.
		this := k.intentions[dep]
		this.parents = append(this.parents, newID)
		k.intentions[dep] = this
	}

	// Recalculate its children's:
	//  - parentUrgency
	//  - parentImportance

	k.intentions[newID] = newNode

	recalculateParentBounties(k.intentions, newNode.Deps...)

	newNode = k.intentions[newID]
	k.roots[newID] = newNode

	if len(newNode.Deps) == 0 {
		k.free[newID] = newNode
	}

	return newNode.Intention
}

// checkNew does NOT lock k.  Use either RLock or Lock first.
func (k *Kernel) checkNew(i message.Intention) mfm.Message {
	for _, id := range i.Deps {
		if _, ok := k.intentions[id]; !ok {
			return message.Error{
				Err: fmt.Errorf("no such Intention %d", id),
			}
		}
	}

	return nil
}

func (k *Kernel) addExisting(i message.Intention) mfm.Message {
	k.RLock()
	if err := k.checkExisting(i); err != nil {
		k.RUnlock()
		return err
	}

	k.RUnlock()

	k.Lock()
	defer k.Unlock()

	// Make sure that we haven't mutated the graph again in the meantime...
	if err := k.checkExisting(i); err != nil {
		return err
	}

	old, oldID := k.intentions[i.ID], i.ID
	affectedIDs := append([]int64(nil), i.Deps...)
	newNode := old.copy()
	newNode.Intention = i

	k.intentions[oldID] = newNode

	seen := make(map[int64]bool)
	oldSeen := make(map[int64]bool)

	// Remove the new deps from roots.
	for _, dep := range newNode.Deps {
		delete(k.roots, dep)
		seen[dep] = true
	}

	// Any old deps that weren't seen in the new deps might be new roots.
	for _, dep := range old.Deps {
		oldSeen[dep] = true
		if seen[dep] {
			continue
		}
		child := k.intentions[dep]
		affectedIDs = append(affectedIDs, dep)
		// If the child had only one parent, this was that parent.  If
		// this does not have that child as a dep any more, the child is
		// therefore orphaned and is a new root.
		if lp := len(child.parents); lp == 1 {
			child.parents = nil
			k.intentions[dep] = child
			k.roots[dep] = child
		} else if lp > 1 {
			// Find the parent and remove it.
		removeParent:
			for iter, p := range child.parents {
				if p == oldID {
					child.parents = append(
						child.parents[:iter],
						child.parents[iter+1:]...,
					)
					k.intentions[dep] = child
					break removeParent
				}
			}
		}
	}

	// Any new deps that weren't in the old deps need this node's ID added
	// to their parents.
	for dep := range seen {
		if !oldSeen[dep] {
			newChild := k.intentions[dep]
			newChild.parents = append(newChild.parents, oldID)
			k.intentions[dep] = newChild
		}
	}

	recalculateParentBounties(k.intentions, affectedIDs...)
	newNode = k.intentions[oldID]

	// Add or remove this node from the free set.
	if len(newNode.Deps) == 0 {
		k.free[oldID] = newNode
	} else {
		delete(k.free, oldID)
	}

	if len(newNode.parents) == 0 {
		k.roots[oldID] = newNode
	}

	return newNode.Intention
}

// checkExisting does NOT lock k.  Use either RLock or Lock first.
func (k *Kernel) checkExisting(i message.Intention) mfm.Message {
	if _, ok := k.intentions[i.ID]; !ok {
		return message.Error{Err: fmt.Errorf("no such Intention %d", i.ID)}
	}

	for _, id := range i.Deps {
		if _, ok := k.intentions[id]; !ok {
			return message.Error{Err: fmt.Errorf("no such Intention %d", id)}
		}

		// Make sure we're not introducing a cycle.
		if cyc := checkCycle(k.intentions, i.ID, id); cyc != nil {
			return message.Error{Err: fmt.Errorf(
				"cycle requested: %v", cyc,
			)}
		}
	}

	return nil
}

// checkCycle will panic bad edges exist or are given.  Check that 'from' and
// 'to' exist first!
func checkCycle(graph map[int64]node, from, to int64) []int64 {
	var (
		seen = make(map[int64]bool)
		curr []node
		next = []node{graph[to]}
	)

	seen[from] = true

	for {
		curr, next = next, nil
		for _, in := range curr {
			seen[in.ID] = true
			for _, child := range in.Deps {
				if seen[child] {
					return []int64{in.ID, child}
				}

				next = append(next, graph[child])
			}
		}
		if len(next) == 0 {
			return nil
		}
	}
}
