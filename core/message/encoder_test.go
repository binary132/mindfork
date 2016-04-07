package message_test

import (
	"bytes"

	coremsg "github.com/mindfork/mindfork/core/message"
	"github.com/mindfork/mindfork/message"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

var _ = message.Encoder(&coremsg.Encoder{})

func (ts *MessageSuite) TestEncode(c *C) {
	w := &bytes.Buffer{}
	mm := &coremsg.MessageMaker{}
	enc := mm.NewEncoder(w)

	for i, test := range []struct {
		should string
		given  message.Message
		expect string
	}{{
		should: "marshal a simple Intention using JSON",
		given:  coremsg.Intention{},
		expect: `{"Who":"","What":""}`,
	}} {
		c.Logf("test %d: should %s", i, test.should)
		c.Logf("  given: %v", test.given)

		w.Reset()

		c.Assert(enc.Encode(test.given), jc.ErrorIsNil)

		c.Check(w.String(), Matches, test.expect+"\n")
	}
}
