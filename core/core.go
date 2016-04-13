package core

import (
	"errors"
	"time"

	"github.com/mindfork/mindfork/message"

	coremsg "github.com/mindfork/mindfork/core/message"
)

// Core is a base implementation of Server and Mutator.
type Core struct {
}

// Serve implements server.Server by handling echo and health checks, then hands
// off control to the Message implementation if it's a Transformer.
func (c *Core) Serve(m message.Message) message.Message {
	if m == message.Message(nil) || m == nil {
		return message.MakeError(errors.New("nil Message"))
	}

	switch tM := m.(type) {
	case coremsg.Echo:
		return m
	case coremsg.Source:
		return Source()
	case coremsg.Intention:
		return c.Intend(tM)
	default:
		return coremsg.Error{Err: errors.New("unknown Message type")}
	}
}

// Intend applies an Intention to a Core.
func (c *Core) Intend(i coremsg.Intention) coremsg.Result {
	return coremsg.Result{coremsg.Intention{Who: "", What: "", When: (*time.Time)(nil)}, coremsg.Error{nil}}
}

// Source simply returns a path to the Mindfork source code.
// TODO: Figure out a way to make this integrate with current running version.
func Source() coremsg.Result {
	return coremsg.Result{struct {
		Source  string
		License string
	}{"github.com/mindfork/mindfork", "Affero GPL"}, coremsg.Error{nil}}
}
