package testing

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mindfork/mindfork"
	"github.com/mindfork/mindfork/core/message"
)

// Decoder implements mindfork.Decoder using encoding/json.Decoder.
type Decoder struct {
	json.Decoder
}

// Decode implements mindfork.Decoder for Decoder.
func (d *Decoder) Decode(m mindfork.Message) error {
	jm := new(jsonMessage)

	if err := d.Decoder.Decode(jm); err != nil {
		return err
	}

	switch jm.Type {
	case Test:
		msg := new(Message)
		if err := json.Unmarshal(jm.RawMessage, msg); err != nil {
			return err
		}
		return message.ReflectSet(m, *msg)
	case "":
		return errors.New("no message Type received")
	default:
		return fmt.Errorf("unknown Type %q", jm.Type)
	}

}
