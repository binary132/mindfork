package kernel

import (
	"fmt"

	"github.com/mindfork/mindfork/core/message"
	mfm "github.com/mindfork/mindfork/message"
)

// Add implements Scheduler.Add on Kernel.
func (k *Kernel) Add(i message.Intention) mfm.Message {
	switch {
	case i.ID != 0:
		return k.addExisting(i)
	}

	return k.addNew(i)
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
