package kernel

import "github.com/mindfork/mindfork/core/scheduler"

// Orderings
const (
	ByID scheduler.Ordering = iota
	ByBounty
	ByDate
	ByScore
)

type sorter struct {
	nodes
	lesser
}

type nodes []node
type lesser interface {
	Less(int, int) bool
}

func whichOrd(ins []node, ord scheduler.Ordering) lesser {
	switch ord {
	case ByID:
		return byID(ins)
	case ByBounty:
		return byBounty(ins)
	case ByDate:
		return byDate(ins)
	case ByScore:
		return byScore(ins)
	default:
		return byDate(ins)
	}
}

type byID []node
type byBounty []node
type byDate []node
type byScore []node

// Len implements Interface.Len on defaultSort.
func (i nodes) Len() int { return len(i) }

// Swap implements Interface.Swap on defaultSort.
func (i nodes) Swap(j, k int) { i[j], i[k] = i[k], i[j] }

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
		return byBounty(b).Less(i, j)
	case wI == nil:
		return false
	case wJ == nil:
		return true
	case wI.Equal(*wJ):
		return byBounty(b).Less(i, j)
	default:
		return wI.Before(*wJ)
	}
}

func score(n node) int {
	return n.Bounty + n.parentBounties
}

// Less implements Interface.Less on byScore.
func (b byScore) Less(i, j int) bool {
	iScore, jScore := score(b[i]), score(b[j])

	if iScore == jScore {
		return byDate(b).Less(i, j)
	}

	return iScore > jScore
}
