package core_test

import (
	"errors"
	"testing"
	"time"

	"github.com/mindfork/mindfork/core"
	"github.com/mindfork/mindfork/core/message"
	coretest "github.com/mindfork/mindfork/core/testing"
	mfm "github.com/mindfork/mindfork/message"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

// Test hooks up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type CoreSuite struct{}

var _ = Suite(&CoreSuite{})

func (cs *CoreSuite) TestServe(c *C) {
	tNow := time.Now()
	mfCore := &core.Core{
		Timer:     coretest.TestTimer(tNow),
		Scheduler: &coretest.MockScheduler{},
	}

	for i, t := range []struct {
		should string
		given  mfm.Message
		expect mfm.Message
	}{{
		should: "return an error for a nil",
		given:  nil,
		expect: message.Error{Err: errors.New("nil Message")},
	}, {
		should: "return an error for a nil Message",
		given:  mfm.Message(nil),
		expect: message.Error{Err: errors.New("nil Message")},
	}, {
		should: "echo for an Echo",
		given:  message.Echo{},
		expect: message.Echo{When: tNow},
	}, {
		should: "return source for a Source",
		given:  message.Source(struct{}{}),
		expect: struct {
			Source  string
			License string
		}{"github.com/mindfork/mindfork", "Affero GPL"},
	}, {
		should: "Intend for an Intention",
		given:  message.Intention{},
		expect: message.Intention{ID: 0},
	}, {
		should: "return error for unknown type",
		given:  mfm.Message(5),
		expect: message.Error{Err: errors.New("unknown Message type")},
	}} {
		c.Logf("test %d: %s", i, t.should)

		c.Check(mfCore.Serve(t.given), jc.DeepEquals, t.expect)
	}
}
