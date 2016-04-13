package testing

import (
	"encoding/json"

	"github.com/mindfork/mindfork/message"
)

// Encode is a JSON mindfork.Encoder for testing.
func Encode(m message.Message) ([]byte, error) {
	return json.Marshal(m)
}
