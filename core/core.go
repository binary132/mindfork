package core

import (
	"errors"

	"github.com/mindfork/mindfork/message"

	coremsg "github.com/mindfork/mindfork/core/message"
)

// Core is a base implementation of the mindfork engine.
type Core struct {
}

// Serve implements server.Server
func (c *Core) Serve(m message.Message) message.Message {
	if m == message.Message(nil) || m == nil {
		return message.MakeError(errors.New("nil Message"))
	}

	switch tM := m.(type) {
	case coremsg.Intention:
		return c.Intend(tM)
	}

	return m
}

// Intend applies an Intention to a Core.
func (c *Core) Intend(i coremsg.Intention) Result {
	return i
}
