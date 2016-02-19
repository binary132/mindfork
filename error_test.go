package mindfork_test

import (
	"errors"

	mf "github.com/mindfork/mindfork"

	. "gopkg.in/check.v1"
)

type testError struct{}

func (t *testError) Error() string {
	return "TestError Result"
}

func (m *MindforkSuite) TestErrorError(c *C) {
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
		e := mf.Error{Err: t.given}
		c.Check(e.Error(), Equals, t.expect)
	}
}

func (m *MindforkSuite) TestMakeError(c *C) {
	for i, t := range []struct {
		should string
		given  error
		expect mf.Error
	}{{
		should: "make an Error containing a nil error",
		given:  nil,
		expect: mf.Error{Err: nil},
	}, {
		should: "make an Error wrapping a non-nil error",
		given:  errors.New("an error"),
		expect: mf.Error{Err: errors.New("an error")},
	}} {
		c.Logf("test %d: should %s", i, t.should)
		c.Check(mf.MakeError(t.given), DeepEquals, t.expect)
	}
}
