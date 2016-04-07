package message

import (
	"encoding/json"
	"errors"
	"fmt"

	mfm "github.com/mindfork/mindfork/message"
)

type jsonMessage struct {
	Type       mfm.Type
	RawMessage json.RawMessage
}

// Decoder implements mindfork.Decoder using encoding/json.Decoder.
type Decoder struct {
	json.Decoder
}

// Decode implements mindfork.Decoder for Decoder.
func (d *Decoder) Decode(m mfm.Message) error {
	jm := new(jsonMessage)

	if err := d.Decoder.Decode(jm); err != nil {
		return err
	}

	switch jm.Type {
	case TIntention:
		msg := new(Intention)
		if err := json.Unmarshal(jm.RawMessage, msg); err != nil {
			return err
		}
		if err := Validate(*msg); err != nil {
			return err
		}
		return mfm.ReflectSet(m, *msg)
	case "":
		return errors.New("no message Type received")
	default:
		return fmt.Errorf("unknown Type %q", jm.Type)
	}

}
