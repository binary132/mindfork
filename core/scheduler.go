package core

import (
	"github.com/mindfork/mindfork/core/message"
	mfm "github.com/mindfork/mindfork/message"
)

// Scheduler is a Core member which handles the scheduling of Intentions.
type Scheduler interface {
	// Add adds the given Intention to the job queue.  It returns the
	// Intention with an ID populated, or any Error.
	Add(message.Intention) mfm.Message

	// Available returns a slice of Intentions which have no dependencies.
	// This slice is a copy.
	Available() []message.Intention

	// Next gets the next

	// // Export returns all Intentions known to the Scheduler.
	// Export() []message.Intention
}
