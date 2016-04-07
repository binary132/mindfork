package testing_test

import (
	"github.com/mindfork/mindfork/message"
	"github.com/mindfork/mindfork/server"
	"github.com/mindfork/mindfork/server/testing"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

var _ = server.Server(&testing.Server{})

func (s *TestingSuite) TestServe(c *C) {
	m := message.Message(5)
	srv := testing.Server{}
	c.Check(srv.Serve(m), jc.DeepEquals, m)
}
