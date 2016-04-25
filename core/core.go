package core

import (
	"errors"
	"time"

	"github.com/mindfork/mindfork/core/message"
	"github.com/mindfork/mindfork/core/scheduler"
	"github.com/mindfork/mindfork/core/scheduler/kernel"
	mfm "github.com/mindfork/mindfork/message"
)

// Timer is an interface for Time that can be mocked.
type Timer interface {
	Now() time.Time
}

// Core is a base implementation of Server.
type Core struct {
	Timer
	scheduler.Scheduler
}

// Default returns a Core using a new scheduler.Kernel.
func Default() *Core {
	return &Core{Scheduler: &kernel.Kernel{
		Intentions: make(map[int64]message.Intention),
		Roots:      make(map[int64]message.Intention),
		Free:       make(map[int64]message.Intention),
	}}
}

// Now implements Timer.Now for Core.
func (c *Core) Now() time.Time {
	if c.Timer == nil {
		return time.Now()
	}
	return c.Timer.Now()
}

// Serve implements server.Server by handling echo and health checks, then hands
// off control to the Message implementation if it's a Transformer.
func (c *Core) Serve(m mfm.Message) mfm.Message {
	if m == mfm.Message(nil) || m == nil {
		return mfm.MakeError(errors.New("nil Message"))
	}

	switch tM := m.(type) {
	case message.Echo:
		return message.Echo{When: c.Now()}
	case message.Source:
		return Source()
	case message.Intention:
		return c.Add(tM)
	default:
		return message.Error{Err: errors.New("unknown Message type")}
	}
}

// Source simply returns a path to the Mindfork source code.
// TODO: Figure out a way to make this integrate with current running version.
func Source() mfm.Message {
	return struct {
		Source  string
		License string
	}{"github.com/mindfork/mindfork", "Affero GPL"}
}
