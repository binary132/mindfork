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
	srv := testing.Server{}
	m, n := message.Message(5), message.Message("hello")
	c.Check(srv.Serve(m), jc.DeepEquals, m)
	c.Check(srv.Messages, jc.DeepEquals, []message.Message{m})
	c.Check(srv.Serve(n), jc.DeepEquals, n)
	c.Check(srv.Messages, jc.DeepEquals, []message.Message{m, n})
}
