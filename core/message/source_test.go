package message_test

import (
	"github.com/mindfork/mindfork/core/message"
	mfm "github.com/mindfork/mindfork/message"
)

var _ = mfm.Message(new(message.Source))
