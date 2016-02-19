package server

import (
	"encoding/json"

	mf "github.com/mindfork/mindfork"
)

// Error is a mf.Error implementing encoding/json.Marshaler.
type Error mf.Error

// ErrStr is a simple marshalable Error struct.
type ErrStr struct {
	Err string
}

// MarshalJSON implements encoding/json.Marshaler on Error.
func (e Error) MarshalJSON() ([]byte, error) {
	if err := e.Err; err != nil {
		// If we can marshal the error using its own Marshaler, do so.
		if _, ok := err.(json.Marshaler); ok {
			return json.Marshal(mf.Error(e))
		}

		// Otherwise, marshal it using its Error() string method.
		return json.Marshal(ErrStr{err.Error()})
	}
	// If the inner error was nil, just use Error's Error() string method.
	return json.Marshal(ErrStr{mf.Error(e).Error()})
}

// Server must be implemented in order to wire up a Mindfork service.
type Server interface {
	// Serve specifies the routing and responses of Messages.
	Serve(m mf.Message) mf.Message
}

// Wrap returns a function of []byte returning ()[]byte, error) which calls
// Serve on the Server with the given byte slice and MessageMaker.
func Wrap(s Server, m mf.MessageMaker) func([]byte) ([]byte, error) {
	return func(bs []byte) ([]byte, error) {
		return Serve(s, m, bs)
	}
}

// Serve handles the generation, routing, and responses of Messages using JSON.
func Serve(s Server, maker mf.MessageMaker, msg []byte) ([]byte, error) {
	m, err := maker(msg)
	if err != nil {
		return json.Marshal(Error{Err: err})
	}

	return json.Marshal(s.Serve(m))
}
