package message_test

import (
	"fmt"

	coremsg "github.com/mindfork/mindfork/core/message"
	mfm "github.com/mindfork/mindfork/message"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

var _ = mfm.Decoder(coremsg.Decode)

func (m *MessageSuite) TestDecode(c *C) {
	for i, t := range []struct {
		should    string
		given     string
		givenType string
		expect    mfm.Message
		expectErr string
	}{{
		should:    "fail to unmarshal broken JSON",
		givenType: `intention`,
		given:     `{`,
		expectErr: `unexpected end of JSON input`,
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
	}, {
		should:    "make a valid Source",
		givenType: string(coremsg.TSource),
		expect:    mfm.Message(coremsg.Source{}),
	}, {
		should:    "make a valid Echo",
		givenType: string(coremsg.TEcho),
		expect:    mfm.Message(coremsg.Echo{}),
	}} {
		c.Logf("test %d: should %s", i, t.should)
		var bs []byte
		if t.given == "" {
			bs = []byte(fmt.Sprintf(`{"Type":%q}`, t.givenType))
		} else {
			bs = []byte(fmt.Sprintf(
				`{"Type":%q,"RawMessage":%s}`,
				t.givenType, t.given,
			))
		}

		m, err := coremsg.Decode(bs)
		if t.expectErr != "" {
			c.Check(err, ErrorMatches, t.expectErr)
			continue
		}

		c.Assert(err, jc.ErrorIsNil)
		c.Check(m, jc.DeepEquals, t.expect)
	}
}
