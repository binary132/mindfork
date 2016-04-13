package message

import (
	"encoding/json"

	mfm "github.com/mindfork/mindfork/message"
)

type jsonMessage struct {
	Type       mfm.Type
	RawMessage json.RawMessage
}

// Maker implements MessageMaker for mindfork Core Messages.
type Maker struct{}

// Encoder implements mindfork.MessageMaker Encoder.
func (t *Maker) Encoder() mfm.Encoder {
	return Encode
}

// Decoder implements mindfork.MessageMaker Decoder.
func (t *Maker) Decoder() mfm.Decoder {
	return Decode
}
