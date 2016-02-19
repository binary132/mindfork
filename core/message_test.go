package core_test

import (
	"fmt"

	mf "github.com/mindfork/mindfork"
	"github.com/mindfork/mindfork/core"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

func (cs *CoreSuite) TestMakeMessage(c *C) {
	for i, t := range []struct {
		should    string
		given     []byte
		givenType []byte
		expect    mf.Message
		expectErr string
	}{{
		should:    "fail to unmarshal broken JSON",
		givenType: []byte(`\"`),
		given:     sampleMessages("emptyMessage"),
		expectErr: `invalid character '}' looking for beginning of ` +
			`value`,
	}, {
		should:    "fail to make unknown Message type",
		given:     []byte(`"{}"`),
		givenType: []byte("x"),
		expectErr: `wrong`,
	}, {
		should:    "fail to make invalid Intention",
		given:     sampleMessages("emptyObject"),
		expectErr: `Intention needs a Who`,
	}, {
		should: "make a valid Intention",
		given:  sampleMessages("validIntention"),
		expect: mf.Message(core.Intention{
			Who:  "Bodie",
			What: "To seek the Holy Grail",
		}),
	}} {
		c.Logf("test %d: should %s", i, t.should)
		bs := []byte(fmt.Sprintf(
			`{%q:%q,%q:%s}`,
			"Type", t.givenType, "RawMessage", t.given,
		))

		m, err := core.MakeMessage(bs)
		if t.expectErr != "" {
			c.Check(err, ErrorMatches, t.expectErr)
			continue
		}

		c.Assert(err, jc.ErrorIsNil)
		c.Check(m, jc.DeepEquals, t.expect)
	}
}
