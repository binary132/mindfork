package message

import (
	"encoding/json"

	mfm "github.com/mindfork/mindfork/message"
)

// Encoder implements mindfork.Encoder using encoding/json.Encoder.
type Encoder struct {
	json.Encoder
}

// Encode implements mindfork.Encoder for Encoder.
func (d *Encoder) Encode(m mfm.Message) error {
	return d.Encoder.Encode(m)
}
