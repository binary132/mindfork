package core_test

import (
	"errors"
	"testing"

	"github.com/mindfork/mindfork/core"
	"github.com/mindfork/mindfork/message"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

// Test hooks up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type CoreSuite struct{}

var _ = Suite(&CoreSuite{})

func (cs *CoreSuite) TestServe(c *C) {
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
	}} {
		c.Logf("test %d: %s", i, t.should)
		c.Check(new(core.Core).Serve(t.given), jc.DeepEquals, t.expect)
	}
}
