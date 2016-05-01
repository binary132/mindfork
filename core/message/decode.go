package message

import (
	"encoding/json"
	"errors"
	"fmt"

	mfm "github.com/mindfork/mindfork/message"
)

// Decode is a JSON Decoder for core.message.
func Decode(m []byte) (mfm.Message, error) {
	jm := new(jsonMessage)

	if err := json.Unmarshal(m, jm); err != nil {
		return nil, err
	}

	switch jm.Type {
	case TIntention:
		msg := new(Intention)
		if err := json.Unmarshal(jm.RawMessage, msg); err != nil {
			return nil, err
		}
		if err := msg.Validate(); err != nil {
			return nil, err
		}
		return *msg, nil
	case TSource:
		return Source{}, nil
	case TEcho:
		return Echo{}, nil
	case "":
		return nil, errors.New("no message Type received")
	default:
		return nil, fmt.Errorf("unknown Type %q", jm.Type)
	}
}
