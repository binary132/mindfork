package testing

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mindfork/mindfork/message"
)

type jsonMessage struct {
	Type       message.Type
	RawMessage json.RawMessage
}

// Decode is a testing JSON Decoder.
func Decode(m []byte) (message.Message, error) {
	jm := new(jsonMessage)

	if err := json.Unmarshal(m, jm); err != nil {
		return nil, err
	}

	switch jm.Type {
	case Test:
		msg := new(Message)
		if err := json.Unmarshal(jm.RawMessage, msg); err != nil {
			return nil, err
		}
		return *msg, nil
	case "":
		return nil, errors.New("no message Type received")
	default:
		return nil, fmt.Errorf("unknown Type %q", jm.Type)
	}
}
