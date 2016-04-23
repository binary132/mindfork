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

	i.ID = int64(len(m.Intentions))
	m.Intentions = append(m.Intentions, i)

	return i
}

// Available implements core.Scheduler on MockScheduler.
func (m *MockScheduler) Available() []message.Intention {
	return m.Intentions
}

// Export implements core.Scheduler on MockScheduler.
func (m *MockScheduler) Export() []message.Intention {
	return m.Intentions
}
