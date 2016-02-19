package core

import (
	"errors"

	mf "github.com/mindfork/mindfork"
)

// Core is a base implementation of the mindfork engine.
type Core struct {
}

// Serve implements server.Server
func (c *Core) Serve(m mf.Message) mf.Message {
	if m == mf.Message(nil) || m == nil {
		return mf.MakeError(errors.New("nil Message"))
	}

	switch tM := m.(type) {
	case Intention:
		return c.Intend(tM)
	}

	return m
}

// Intend applies an Intention to a Core.
func (c *Core) Intend(i Intention) Result {
	return Result{Message: "foo"}
}
