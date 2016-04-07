package testing

import (
	"encoding/json"
	"io"

	"github.com/mindfork/mindfork/message"
)

const (
	Test message.Type = "test"
)

// MessageMaker is a testing implementation of mindfork.MessageMaker for JSON.
type MessageMaker struct {
}

// NewEncoder implements mindfork.MessageMaker NewEncoder.
func (t *MessageMaker) NewEncoder(w io.Writer) message.Encoder {
	return &Encoder{*json.NewEncoder(w)}
}

// NewDecoder implements mindfork.MessageMaker NewDecoder.
func (t *MessageMaker) NewDecoder(r io.Reader) message.Decoder {
	return &Decoder{*json.NewDecoder(r)}
}

// Message is a mindfork.Message that testing.MessageMaker exercises.
type Message struct {
	X int
	S string `json:",omitempty"`
}
