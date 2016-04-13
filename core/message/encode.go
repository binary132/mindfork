package message

import (
	"encoding/json"

	mfm "github.com/mindfork/mindfork/message"
)

// Encode is a JSON Encoder for core.message.
func Encode(m mfm.Message) ([]byte, error) {
	return json.Marshal(m)
}
