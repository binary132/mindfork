package testing

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	mf "github.com/mindfork/mindfork"
)

const (
	Test mf.Type = "test"
)

// MessageMaker is a testing implementation of mindfork.MessageMaker.
type MessageMaker struct {
}

// NewEncoder implements mindfork.MessageMaker NewEncoder.
func (t *MessageMaker) NewEncoder(w io.Writer) mf.Encoder {
	return json.NewEncoder(w)
}

// NewDecoder implements mindfork.MessageMaker NewDecoder.
func (t *MessageMaker) NewDecoder(r io.Reader) mf.Decoder {
	return json.NewDecoder(r)
}

type jsonMessage struct {
	Type       mf.Type
	RawMessage json.RawMessage
}

// UnmarshalJSON implements encoding/json.Unmarshaler on Message.
func (t *Message) UnmarshalJSON(bs []byte) error {
	jm := jsonMessage{}
	err := json.Unmarshal(bs, &jm)
	if err != nil {
		return mf.Error{Err: err}
	}

	switch jm.Type {
	case Test:
		tm := tmpMessage{}
		err := json.Unmarshal(jm.RawMessage, &tm)
		if err != nil {
			return mf.Error{Err: err}
		}
		*t = Message(tm)
		return nil
	case "":
		return mf.Error{Err: errors.New("no message Type received")}
	default:
		return mf.Error{Err: fmt.Errorf("unknown Type %q", jm.Type)}
	}
}

// Message is a mindfork.Message that testing.MessageMaker exercises.
type Message struct {
	X int
	S string `json:",omitempty"`
}

// tmpMessage aliases Message to avoid json.Unmarshal recursion.
type tmpMessage Message
