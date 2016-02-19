package core

import (
	"encoding/json"
	"fmt"

	mf "github.com/mindfork/mindfork"
)

const (
	TIntention mf.Type = "intention"
)

// MakeMessage is a mindfork.MessageMaker that can make a core.Intention.
func MakeMessage(bs []byte) (mf.Message, error) {
	mb := mf.MessageBytes{}
	if err := json.Unmarshal(bs, &mb); err != nil {
		return nil, err
	}

	switch t := mb.Type; t {
	case TIntention:
		return MakeIntention(mb.RawMessage)
	default:
		return nil, fmt.Errorf("unknown message kind %q", t)
	}
}
