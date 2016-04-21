package message_test

import (
	"fmt"
	"time"

	message "github.com/mindfork/mindfork/core/message"
	"github.com/mindfork/mindfork/core/testing"
	mfm "github.com/mindfork/mindfork/message"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

var _ = mfm.Decoder(message.Decode)

func (m *MessageSuite) TestDecode(c *C) {
	tExpect, err := time.Parse(time.RFC3339, "2009-11-10T23:00:00Z")
	c.Assert(err, jc.ErrorIsNil)

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
		given:     testing.SampleMessages("emptyObject"),
		expectErr: `unknown Type "x"`,
	}, {
		should:    "fail to make invalid Intention",
		givenType: string(message.TIntention),
		given:     testing.SampleMessages("emptyObject"),
		expectErr: `Intention needs a Who`,
	}, {
		should:    "make a valid Intention",
		givenType: string(message.TIntention),
		given:     testing.SampleMessages("validIntention"),
		expect: message.Intention{
			Who:  "Bodie",
			What: "To seek the Holy Grail",
		},
	}, {
		should:    "make a valid Intention",
		givenType: string(message.TIntention),
		given:     testing.SampleMessages("timedIntention"),
		expect: message.Intention{
			Who:  "Bodie",
			What: "Something neat",
			When: &tExpect,
		},
	}, {
		should:    "make a valid Source",
		givenType: string(message.TSource),
		given:     testing.SampleMessages("validIntention"),
		expect:    message.Source{},
	}, {
		should:    "make a valid Echo",
		givenType: string(message.TEcho),
		expect:    message.Echo{},
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

		m, err := message.Decode(bs)
		if t.expectErr != "" {
			c.Check(err, ErrorMatches, t.expectErr)
			continue
		}

		c.Assert(err, jc.ErrorIsNil)
		c.Check(m, jc.DeepEquals, t.expect)
	}
}
