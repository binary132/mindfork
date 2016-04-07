package message

import (
	"encoding/json"
	"io"

	mfm "github.com/mindfork/mindfork/message"
)

const (
	TIntention mfm.Type = "intention"
)

// MessageMaker is a testing implementation of mindfork.MessageMaker for JSON.
type MessageMaker struct {
}

// NewEncoder implements mindfork.MessageMaker NewEncoder.
func (t *MessageMaker) NewEncoder(w io.Writer) mfm.Encoder {
	return &Encoder{*json.NewEncoder(w)}
}

// NewDecoder implements mindfork.MessageMaker NewDecoder.
func (t *MessageMaker) NewDecoder(r io.Reader) mfm.Decoder {
	return &Decoder{*json.NewDecoder(r)}
}
