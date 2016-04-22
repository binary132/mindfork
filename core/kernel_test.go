package core_test

import (
	"reflect"

	"github.com/mindfork/mindfork/core"
	"github.com/mindfork/mindfork/core/message"
	mfm "github.com/mindfork/mindfork/message"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

func (cs *CoreSuite) TestKernelAdd(c *C) {
	for i, test := range []struct {
		given         []message.Intention
		expectResults []mfm.Message
		expectKern    []message.Intention
	}{{
		given: []message.Intention{{}, {}, {}},
		expectResults: []mfm.Message{
			message.Intention{},
			message.Intention{ID: 1},
			message.Intention{ID: 2},
		},
		expectKern: []message.Intention{{}, {ID: 1}, {ID: 2}},
	}} {
		c.Logf("test %d: should", i)

		k := core.NewKernel()

		results := make([]mfm.Message, len(test.given))

		for i, in := range test.given {
			results[i] = k.Add(in)
			c.Check(
				reflect.TypeOf(results[i]),
				Equals,
				reflect.TypeOf(mfm.Message(message.Intention{})),
			)
		}

		c.Check(results, jc.DeepEquals, test.expectResults)
		c.Check(k.Export(), jc.DeepEquals, test.expectKern)
	}
}
