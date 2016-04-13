package testing

import "github.com/mindfork/mindfork/message"

const (
	Test message.Type = "test"
)

// MessageMaker is a testing implementation of mindfork.MessageMaker for JSON.
type MessageMaker struct {
}

// Encoder implements mindfork.MessageMaker Encoder.
func (t *MessageMaker) Encoder() message.Encoder {
	return Encode
}

// Decoder implements mindfork.MessageMaker Decoder.
func (t *MessageMaker) Decoder() message.Decoder {
	return Decode
}

// Message is a mindfork.Message that testing.MessageMaker exercises.
type Message struct {
	X int
	S string `json:",omitempty"`
}
