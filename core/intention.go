package core

import (
	"encoding/json"
	"errors"
	"time"

	mf "github.com/mindfork/mindfork"
)

// MakeIntention is a mindfork.MessageMaker that only knows how to make an
// Intention using JSON.
func MakeIntention(bs []byte) (mf.Message, error) {
	i := Intention{}
	if err := json.Unmarshal(bs, &i); err != nil {
		return nil, err
	}
	if err := Validate(i); err != nil {
		return nil, err
	}
	return i, nil
}

// Intention is the necessary information to make a Mindfork intention.
type Intention struct {
	Who  string
	What string
	When time.Time `json:"omitempty"`
}

// Validate validates an Intention.
func Validate(i Intention) error {
	if i.Who == "" {
		return errors.New("Intention needs a Who")
	}
	return nil
}
