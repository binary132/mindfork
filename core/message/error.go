package message

import (
	"encoding/json"

	"github.com/mindfork/mindfork"
)

// Error is a mf.Error implementing encoding/json.Marshaler.
type Error mindfork.Error

// ErrStr is a simple marshalable Error struct.
type ErrStr struct {
	Err string
}

// MarshalJSON implements encoding/json.Marshaler on Error.
func (e Error) MarshalJSON() ([]byte, error) {
	if err := e.Err; err != nil {
		// If we can marshal the error using its own Marshaler, do so.
		if _, ok := err.(json.Marshaler); ok {
			return json.Marshal(mindfork.Error(e))
		}

		// Otherwise, marshal it using its Error() string method.
		return json.Marshal(ErrStr{err.Error()})
	}
	// If the inner error was nil, just use Error's Error() string method.
	return json.Marshal(ErrStr{mindfork.Error(e).Error()})
}
