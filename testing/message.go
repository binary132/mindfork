package testing

import (
	"encoding/json"
	"io"

	mf "github.com/mindfork/mindfork"
)

const (
	Test mf.Type = "test"
)

// MessageMaker is a testing implementation of mindfork.MessageMaker for JSON.
type MessageMaker struct {
}

// NewEncoder implements mindfork.MessageMaker NewEncoder.
func (t *MessageMaker) NewEncoder(w io.Writer) mf.Encoder {
	return &Encoder{*json.NewEncoder(w)}
}

// NewDecoder implements mindfork.MessageMaker NewDecoder.
func (t *MessageMaker) NewDecoder(r io.Reader) mf.Decoder {
	return &Decoder{*json.NewDecoder(r)}
}

type jsonMessage struct {
	Type       mf.Type
	RawMessage json.RawMessage
}

// Message is a mindfork.Message that testing.MessageMaker exercises.
type Message struct {
	X int
	S string `json:",omitempty"`
}

// tmpMessage aliases Message to avoid json.Unmarshal recursion.
type tmpMessage Message
