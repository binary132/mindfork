package message_test

import (
	"github.com/mindfork/mindfork/core/message"
	mfm "github.com/mindfork/mindfork/message"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

var _ = mfm.Message(message.Intention{})

func (m *MessageSuite) TestIntentionFulfill(c *C) {
	for i, t := range []struct {
		should      string
		intention   message.Intention
		fulfillment message.Fulfillment
		expectMsgs  []mfm.Message
		expectError string
	}{{}} {
		c.Logf("test %d: should %s", i, t.should)
		result, err := t.intention.Fulfill(t.fulfillment)
		if t.expectError != "" {
			c.Check(err, ErrorMatches, t.expectError)
			continue
		}
		c.Check(result, jc.DeepEquals, t.expectMsgs)
	}
}
