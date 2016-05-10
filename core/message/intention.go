package message

import (
	"errors"
	"time"

	mfm "github.com/mindfork/mindfork/message"
)

// Intention is the necessary information to make a Mindfork intention.
type Intention struct {
	ID int64 `json:",omitempty"`

	Who  string     `json:",omitempty"`
	What string     `json:",omitempty"`
	When *time.Time `json:",omitempty"`

	Bounty int `json:",omitempty"`
	// Urgency    int `json:",omitempty"`
	// Importance int `json:",omitempty"`
	// Cost       int `json:",omitempty"`

	Deps []int64 `json:",omitempty"`
}

// Validate validates an Intention.
func (i Intention) Validate() error {
	if i.Who == "" {
		return errors.New("Intention needs a Who")
	}
	return nil
}

// Fulfill fulfills the given Intention with the given Fulfillment.  If it
// succeeds, it will generate a new slice of Messages to be Served.
func (i *Intention) Fulfill(f Fulfillment) ([]mfm.Message, error) {
	return nil, errors.New("not implemented")
}
