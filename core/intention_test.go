package core_test

import (
	"time"

	mf "github.com/mindfork/mindfork"
	"github.com/mindfork/mindfork/core"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

func (cs *CoreSuite) TestMakeIntention(c *C) {
	for i, t := range []struct {
		should      string
		given       []byte
		expect      mf.Message
		expectError string
	}{{
		should:      "give an error for an untyped message",
		given:       sampleMessages("emptyMessage"),
		expectError: "unexpected end of JSON input",
	}, {
		should: "give an error for a wrongly-typed message",
		given:  sampleMessages("emptyString"),
		expectError: "json: cannot unmarshal string into Go value of " +
			"type core.Intention",
	}, {
		should:      "return an error for an empty Who",
		given:       sampleMessages("emptyObject"),
		expectError: "Intention needs a Who",
	}, {
		should: "return a valid Intention for valid JSON",
		given:  sampleMessages("validIntention"),
		expect: mf.Message(core.Intention{
			Who:  "Bodie",
			What: "To seek the Holy Grail",
		}),
	}, {
		should: "return a valid Intention for valid JSON with time",
		given:  sampleMessages("timedIntention"),
		expect: mf.Message(core.Intention{
			Who:  "Bodie",
			What: "Something neat",
			When: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local),
		}),
	}} {
		c.Logf("test %d: should %s", i, t.should)
		m, err := core.MakeIntention(t.given)
		if t.expectError != "" {
			c.Check(err, ErrorMatches, t.expectError)
			continue
		}
		c.Assert(err, jc.ErrorIsNil)
		c.Check(m, jc.DeepEquals, t.expect)
	}
}
