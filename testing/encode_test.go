package testing_test

import (
	"github.com/mindfork/mindfork/message"
	"github.com/mindfork/mindfork/testing"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

var _ = message.Encoder(testing.Encode)

func (ts *TestingSuite) TestEncode(c *C) {
	for i, test := range []struct {
		should string
		given  message.Message
		expect string
	}{{"encode a simple Message", testing.Message{X: 5}, `{"X":5}`}} {
		c.Logf("test %d: should %s", i, test.should)
		c.Logf("  given: %v", test.given)

		bs, err := testing.Encode(test.given)
		c.Assert(err, jc.ErrorIsNil)

		c.Check(string(bs), Matches, test.expect)
	}
}
