package testing

import (
	"sync"

	"github.com/mindfork/mindfork/core/message"
	mfm "github.com/mindfork/mindfork/message"
)

// MockScheduler implements core.Scheduler for testing.
type MockScheduler struct {
	sync.RWMutex
	Intentions []message.Intention
}

// Add implements core.Scheduler on MockScheduler.
func (m *MockScheduler) Add(i message.Intention) mfm.Message {
	m.Lock()
	defer m.Unlock()

	m.Intentions = append(m.Intentions, i)
	i.ID = int64(len(m.Intentions) - 1)

	return i
}

// Peek returns a slice of Intentions which have no dependencies.
func (m *MockScheduler) Peek() []message.Intention {
	return m.Intentions
}

// Export returns all Intentions known to the Scheduler.
func (m *MockScheduler) Export() []message.Intention {
	return m.Intentions
}
