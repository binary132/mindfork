package core

import (
	"encoding/json"

	mf "github.com/mindfork/mindfork"
)

// Result is a Message end result.
type Result struct {
	mf.Message
	Err mf.Error `json:"omitempty"`
}

// MarshalJSON implements json.Marshaler for Result.
func (r *Result) MarshalJSON() ([]byte, error) {
	if err := r.Err; err.Err != nil {
		return json.Marshal(err)
	}

	return r.MarshalJSON()
}
