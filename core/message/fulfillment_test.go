package message_test

import (
	"github.com/mindfork/mindfork/core/message"
	mfm "github.com/mindfork/mindfork/message"

	. "gopkg.in/check.v1"
)

var _ = mfm.Message(&message.Fulfillment{})

func (m *MessageSuite) TestFulfillment(c *C) {}
