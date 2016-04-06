package testing_test

import (
	"bytes"

	"github.com/mindfork/mindfork"
	"github.com/mindfork/mindfork/testing"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

var _ = mindfork.Encoder(&testing.Encoder{})

func (ts *TestingSuite) TestEncode(c *C) {
	w := &bytes.Buffer{}
	mm := &testing.MessageMaker{}
	enc := mm.NewEncoder(w)

	for i, test := range []struct {
		should string
		given  mindfork.Message
		expect string
	}{{"say hello", testing.Message{X: 5}, `{"X":5}`}} {
		c.Logf("test %d: should %s", i, test.should)
		c.Logf("  given: %v", test.given)

		w.Reset()

		c.Assert(enc.Encode(test.given), jc.ErrorIsNil)

		c.Check(w.String(), Matches, test.expect+"\n")
	}
}
