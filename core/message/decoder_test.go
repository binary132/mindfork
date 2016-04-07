package message_test

import (
	"bytes"
	"fmt"

	coremsg "github.com/mindfork/mindfork/core/message"
	mfm "github.com/mindfork/mindfork/message"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

func (m *MessageSuite) TestDecode(c *C) {
	for i, t := range []struct {
		should    string
		given     string
		givenType string
		expect    mfm.Message
		expectErr string
	}{{
		should:    "fail to unmarshal broken JSON",
		givenType: `\"`,
		given:     sampleMessages("emptyMessage"),
		expectErr: `invalid character '}' looking for beginning of ` +
			`value`,
	}, {
		should:    "fail to make unknown Message type",
		givenType: "x",
		given:     sampleMessages("emptyObject"),
		expectErr: `unknown Type "x"`,
	}, {
		should:    "fail to make invalid Intention",
		givenType: string(coremsg.TIntention),
		given:     sampleMessages("emptyObject"),
		expectErr: `Intention needs a Who`,
	}, {
		should:    "make a valid Intention",
		givenType: string(coremsg.TIntention),
		given:     sampleMessages("validIntention"),
		expect: mfm.Message(coremsg.Intention{
			Who:  "Bodie",
			What: "To seek the Holy Grail",
		}),
	}} {
		c.Logf("test %d: should %s", i, t.should)
		bs := []byte(fmt.Sprintf(
			`{%q:%q,%q:%s}`,
			"Type", t.givenType, "RawMessage", t.given,
		))

		m := new(mfm.Message)
		de := (&coremsg.MessageMaker{}).NewDecoder(bytes.NewReader(bs))
		err := de.Decode(m)
		if t.expectErr != "" {
			c.Check(err, ErrorMatches, t.expectErr)
			continue
		}

		c.Assert(err, jc.ErrorIsNil)
		c.Check(*m, jc.DeepEquals, t.expect)
	}
}
