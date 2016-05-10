package message

import (
	"encoding/json"

	"github.com/mindfork/mindfork/message"
)

// Error is an alias for mf.Error implementing encoding/json.Marshaler.
type Error message.Error

// ErrStr is a simple marshalable Error struct.
type ErrStr struct {
	Err string
}

// MarshalJSON implements encoding/json.Marshaler on Error.
func (e Error) MarshalJSON() ([]byte, error) {
	if err := e.Err; err != nil {
		// If we can marshal the error using its own Marshaler, do so.
		if _, ok := err.(json.Marshaler); ok {
			return json.Marshal(message.Error(e))
		}

		// Otherwise, marshal it using its Error() string method.
		return json.Marshal(ErrStr{err.Error()})
	}

	// If the inner error was nil, just use Error's Error() string method.
	return json.Marshal(ErrStr{message.Error(e).Error()})
}
