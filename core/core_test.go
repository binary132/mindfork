package core_test

import (
	"errors"
	"testing"

	mf "github.com/mindfork/mindfork"
	"github.com/mindfork/mindfork/core"

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
		given  mf.Message
		expect mf.Message
	}{{
		should: "return an error for a nil",
		given:  nil,
		expect: mf.Error{Err: errors.New("nil Message")},
	}, {
		should: "return an error for a nil Message",
		given:  mf.Message(nil),
		expect: mf.Error{Err: errors.New("nil Message")},
	}} {
		c.Logf("test %d: %s", i, t.should)
		c.Check(new(core.Core).Serve(t.given), jc.DeepEquals, t.expect)
	}
}
