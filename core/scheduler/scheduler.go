package scheduler

import (
	"github.com/mindfork/mindfork/core/message"
	mfm "github.com/mindfork/mindfork/message"
)

// Ordering is a method that implements sort.Interface.Less.
type Ordering interface {
	Less(int, int) bool
}

// Scheduler is a Core member which handles the scheduling of Intentions.  It
// should be implemented in a concurrency-safe way.
type Scheduler interface {
	// Add adds the given Intention to the job queue.  It returns the
	// Intention with an ID populated, or any Error.
	Add(message.Intention) mfm.Message

	// Available returns a slice of Intentions which have no dependencies,
	// ordered using the given Ordering.  This slice is a copy.
	Available(Ordering) []message.Intention

	// Next gets the next

	// // Export returns all Intentions known to the Scheduler.
	// Export() []message.Intention
}
