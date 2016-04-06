package testing_test

import (
	"bytes"

	"github.com/mindfork/mindfork"
	"github.com/mindfork/mindfork/testing"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

var _ = mindfork.Decoder(&testing.Decoder{})

func (ts *TestingSuite) TestDecode(c *C) {
	for i, test := range []struct {
		should    string
		given     string
		expect    mindfork.Message
		expectErr string
	}{{
		should: "error on broken JSON",
		given:  `{X":5}`,
		expectErr: `invalid character 'X' looking for beginning of ` +
			`object key string`,
	}, {
		should: "work for simple case",
		given:  `{"Type":"test","RawMessage":{"X":5}}`,
		expect: mindfork.Message(&testing.Message{X: 5}),
	}} {
		c.Logf("test %d: should %s", i, test.should)
		c.Logf("  given: %s", test.given)

		r := &bytes.Buffer{}
		mm := &testing.MessageMaker{}
		dec := mm.NewDecoder(r)

		_, err := r.WriteString(test.given)
		c.Assert(err, jc.ErrorIsNil)

		m := new(testing.Message)
		err = dec.Decode(m)

		if test.expectErr != "" {
			c.Check(err, ErrorMatches, test.expectErr)
			continue
		}

		c.Assert(err, jc.ErrorIsNil)

		c.Check(m, jc.DeepEquals, test.expect)
	}
}
