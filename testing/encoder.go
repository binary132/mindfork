package testing

import (
	"encoding/json"

	"github.com/mindfork/mindfork/message"
)

// Encoder implements mindfork.Encoder using encoding/json.Encoder.
type Encoder struct {
	json.Encoder
}

// Encode implements mindfork.Encoder for Encoder.
func (d *Encoder) Encode(m message.Message) error {
	return d.Encoder.Encode(m)
}
