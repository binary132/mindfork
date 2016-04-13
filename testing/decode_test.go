package testing_test

import (
	"github.com/mindfork/mindfork/message"
	"github.com/mindfork/mindfork/testing"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

var _ = message.Decoder(testing.Decode)

func (ts *TestingSuite) TestDecode(c *C) {
	for i, test := range []struct {
		should    string
		given     string
		expect    message.Message
		expectErr string
	}{{
		should: "error on broken JSON",
		given:  `{X":5}`,
		expectErr: `invalid character 'X' looking for beginning of ` +
			`object key string`,
	}, {
		should: "work for simple case",
		given:  `{"Type":"test","RawMessage":{"X":5}}`,
		expect: testing.Message{X: 5},
	}} {
		c.Logf("test %d: should %s", i, test.should)
		c.Logf("  given: %s", test.given)

		m, err := testing.Decode([]byte(test.given))

		if test.expectErr != "" {
			c.Check(err, ErrorMatches, test.expectErr)
			continue
		}

		c.Assert(err, jc.ErrorIsNil)

		c.Check(m, jc.DeepEquals, test.expect)
	}
}
