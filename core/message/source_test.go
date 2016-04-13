package message_test

import (
	coremsg "github.com/mindfork/mindfork/core/message"
	"github.com/mindfork/mindfork/message"
)

var _ = message.Message(new(coremsg.Source))
