package message

import (
	"errors"
	"time"
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
func Validate(i Intention) error {
	if i.Who == "" {
		return errors.New("Intention needs a Who")
	}
	return nil
}
