package message

import (
	"errors"
	"time"
)

// Intention is the necessary information to make a Mindfork intention.
type Intention struct {
	ID   int64
	Who  string
	What string
	When *time.Time `json:",omitempty"`
	Deps []int64    `json:",omitempty"`
}

// Validate validates an Intention.
func Validate(i Intention) error {
	if i.Who == "" {
		return errors.New("Intention needs a Who")
	}
	return nil
}
