package message_test

import (
	"errors"

	"github.com/mindfork/mindfork/message"

	. "gopkg.in/check.v1"
)

type testError struct{}

func (t *testError) Error() string {
	return "TestError Result"
}

func (m *MessageSuite) TestErrorError(c *C) {
	for i, t := range []struct {
		should string
		given  error
		expect string
	}{{
		should: "produce empty result for nil Error",
		given:  nil,
		expect: "",
	}, {
		should: "marshal a normal Error",
		given:  errors.New("an error"),
		expect: "an error",
	}, {
		should: "marshal a custom Error",
		given:  &testError{},
		expect: "TestError Result",
	}} {
		c.Logf("test %d: should %s", i, t.should)
		e := message.Error{Err: t.given}
		c.Check(e.Error(), Equals, t.expect)
	}
}

func (m *MessageSuite) TestMakeError(c *C) {
	for i, t := range []struct {
		should string
		given  error
		expect message.Error
	}{{
		should: "make an Error containing a nil error",
		given:  nil,
		expect: message.Error{Err: nil},
	}, {
		should: "make an Error wrapping a non-nil error",
		given:  errors.New("an error"),
		expect: message.Error{Err: errors.New("an error")},
	}} {
		c.Logf("test %d: should %s", i, t.should)
		c.Check(message.MakeError(t.given), DeepEquals, t.expect)
	}
}
