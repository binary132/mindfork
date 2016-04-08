package core_test

import (
	"errors"
	"testing"

	"github.com/mindfork/mindfork/core"
	coremsg "github.com/mindfork/mindfork/core/message"
	"github.com/mindfork/mindfork/message"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

// Test hooks up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type CoreSuite struct{}

var _ = Suite(&CoreSuite{})

func (cs *CoreSuite) TestServe(c *C) {
	mfCore := new(core.Core)

	for i, t := range []struct {
		should string
		given  message.Message
		expect message.Message
	}{{
		should: "return an error for a nil",
		given:  nil,
		expect: message.Error{Err: errors.New("nil Message")},
	}, {
		should: "return an error for a nil Message",
		given:  message.Message(nil),
		expect: message.Error{Err: errors.New("nil Message")},
	}, {
		should: "echo a non-nil Intention",
		given:  message.Message(coremsg.Intention{}),
		expect: message.Message(coremsg.Intention{}),
	}} {
		c.Logf("test %d: %s", i, t.should)
		c.Check(mfCore.Serve(t.given), jc.DeepEquals, t.expect)
	}
}
