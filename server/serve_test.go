package server_test

import (
	"errors"

	mf "github.com/mindfork/mindfork"
	"github.com/mindfork/mindfork/server"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

type mockServer struct{}

func (m mockServer) Serve(msg mf.Message) mf.Message {
	return msg
}

func makeMaker(err error) mf.MessageMaker {
	return func(bs []byte) (mf.Message, error) {
		if err != nil {
			return nil, err
		}
		return struct{ M string }{string(bs)}, nil
	}
}

func (s *ServerSuite) TestServe(c *C) {
	for i, t := range []struct {
		should     string
		givenMaker mf.MessageMaker
		givenMsg   string
		expect     string
		expectErr  string
	}{{
		should:     "make some message",
		givenMaker: makeMaker(nil),
		givenMsg:   "hello",
		expect:     `{"M":"hello"}`,
	}, {
		should:     "have some error",
		givenMaker: makeMaker(errors.New("oops")),
		givenMsg:   "hello",
		expect:     `{"Err":"oops"}`,
	}} {
		c.Logf("test %d: should %s", i, t.should)
		bs, err := server.Serve(
			mockServer{}, t.givenMaker, []byte(t.givenMsg),
		)
		if t.expectErr != "" {
			c.Check(err, ErrorMatches, t.expectErr)
			continue
		}

		c.Assert(err, jc.ErrorIsNil)

		c.Check(string(bs), Matches, t.expect)
	}
}
