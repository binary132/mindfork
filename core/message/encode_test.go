package message_test

import (
	coremsg "github.com/mindfork/mindfork/core/message"
	"github.com/mindfork/mindfork/message"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

func (ts *MessageSuite) TestEncode(c *C) {
	for i, test := range []struct {
		should string
		given  message.Message
		expect string
	}{{
		should: "marshal a simple Intention using JSON",
		given:  coremsg.Intention{},
		expect: `{}`,
	}} {
		c.Logf("test %d: should %s", i, test.should)
		c.Logf("  given: %v", test.given)

		bs, err := coremsg.Encode(test.given)
		c.Assert(err, jc.ErrorIsNil)

		c.Check(string(bs), Matches, test.expect)
	}
}
