package scheduler

import (
	"fmt"
	"sort"
	"sync"

	"github.com/mindfork/mindfork/core/message"
	mfm "github.com/mindfork/mindfork/message"
)

type defaultSort []message.Intention

// Len implements Interface.Len on defaultSort.
func (d defaultSort) Len() int { return len(d) }

// Less implements Interface.Less on defaultSort.
func (d defaultSort) Less(i, j int) bool { return d[i].ID < d[j].ID }

// Swap implements Interface.Swap on defaultSort.
func (d defaultSort) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

// Kernel is the core Scheduler implementation.  It holds Intentions in volatile
// memory.
type Kernel struct {
	sync.RWMutex

	Roots      map[int64]message.Intention
	Intentions map[int64]message.Intention
	Free       map[int64]message.Intention

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
func (k *Kernel) Available() []message.Intention {
	k.RLock()
	result := make([]message.Intention, len(k.Free))
	i := 0
	for _, in := range k.Free {
		result[i] = in
		i++
	}
	k.RUnlock()

	sort.Sort(defaultSort(result))

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

	k.Intentions[i.ID] = i
	k.Roots[i.ID] = i

	// Remove the new deps from Roots.
	for _, child := range i.Deps {
		delete(k.Roots, child)
	}

	if len(i.Deps) == 0 {
		k.Free[i.ID] = i
	}

	return i
}

// checkNew does NOT lock k.  Use either RLock or Lock first.
func (k *Kernel) checkNew(i message.Intention) mfm.Message {
	for _, id := range i.Deps {
		if _, ok := k.Intentions[id]; !ok {
			return message.Error{Err: fmt.Errorf("no such Intention %d", id)}
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

	old := k.Intentions[i.ID]
	k.Intentions[i.ID] = i

	seen := make(map[int64]bool)

	// Remove the new deps from Roots.
	for _, dep := range i.Deps {
		delete(k.Roots, dep)
		seen[dep] = true
	}

	// Remove any old deps from Roots that weren't seen in the new deps.
	for _, dep := range old.Deps {
		if !seen[dep] {
			delete(k.Roots, dep)
		}
	}

	// Add or remove this node from the Free set.
	if len(i.Deps) == 0 {
		k.Free[i.ID] = i
	} else {
		delete(k.Free, i.ID)
	}

	return i
}

// checkExisting does NOT lock k.  Use either RLock or Lock first.
func (k *Kernel) checkExisting(i message.Intention) mfm.Message {
	if _, ok := k.Intentions[i.ID]; !ok {
		return message.Error{Err: fmt.Errorf("no such Intention %d", i.ID)}
	}

	for _, id := range i.Deps {
		if _, ok := k.Intentions[id]; !ok {
			return message.Error{Err: fmt.Errorf("no such Intention %d", id)}
		}

		// Make sure we're not introducing a cycle.
		if cyc := checkCycle(k.Intentions, i.ID, id); cyc != nil {
			return message.Error{Err: fmt.Errorf(
				"cycle requested: %v", cyc,
			)}
		}
	}

	return nil
}

// checkCycle will panic bad edges exist or are given.  Check that 'from' and
// 'to' exist first!
func checkCycle(graph map[int64]message.Intention, from, to int64) []int64 {
	var (
		seen = make(map[int64]bool)
		curr []message.Intention
		next = []message.Intention{graph[to]}
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
