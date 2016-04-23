package scheduler

import (
	"github.com/mindfork/mindfork/core/message"
	mfm "github.com/mindfork/mindfork/message"
)

// Cayley implements core.Scheduler using a github.com/google/cayley HTTP query.
type Cayley struct{}

// Add implements Scheduler.Add on Cayley
func (c *Cayley) Add(message.Intention) mfm.Message { return nil }

// Available implements Scheduler.Available on Cayley.
func (c *Cayley) Available() []message.Intention { return nil }

// Export implements Scheduler.Export on Cayley.
func (c *Cayley) Export() []message.Intention { return nil }
