package scheduler

import (
	"github.com/mindfork/mindfork/core/message"
	mfm "github.com/mindfork/mindfork/message"
)

// Ordering is a package constant which specifies Intention orderings.
type Ordering int

// Scheduler is a Core member which handles the scheduling of Intentions.  It
// should be implemented in a concurrency-safe way.
type Scheduler interface {
	// Add adds the given Intention to the job queue.  It returns the
	// Intention with an ID populated, or any Error.
	Add(message.Intention) mfm.Message

	// Available returns a slice of Intentions which have no dependencies,
	// ordered using the given Ordering.  This slice is a copy.
	Available(Ordering) []message.Intention

	// Fulfill applies a Fulfillment which resolves an Intention.  This is
	// only to be permitted if the Intention has no unresolved deps.
	Fulfill(message.Fulfillment) mfm.Message

	// // Export returns all Intentions known to the Scheduler.
	// Export() []message.Intention
}
