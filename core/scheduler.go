package core

import (
	"sync"

	"github.com/mindfork/mindfork/core/message"
	mfm "github.com/mindfork/mindfork/message"
)

// Scheduler is a Core member which handles the scheduling of tasks.
type Scheduler interface {
	Add(message.Intention) mfm.Message
}

// Kernel is the core Scheduler implementation.
type Kernel struct {
	sync.RWMutex
	Intentions []message.Intention
}

// Add implements Scheduler on Kernel.
func (k *Kernel) Add(i message.Intention) mfm.Message {
	k.Lock()
	defer k.Unlock()

	k.Intentions = append(k.Intentions, i)
	i.ID = int64(len(k.Intentions) - 1)

	return i
}
