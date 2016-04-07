package message_test

import (
	"github.com/mindfork/mindfork/message"
	mft "github.com/mindfork/mindfork/testing"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

func (m *MessageSuite) TestReflectSet(c *C) {
	msg := new(message.Message)
	x := &mft.Message{X: 5}

	c.Assert(message.ReflectSet(msg, x), jc.ErrorIsNil)

	c.Check(*msg, jc.DeepEquals, message.Message(x))
}
