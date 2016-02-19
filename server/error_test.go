package server_test

import (
	"encoding/json"
	"errors"

	"github.com/mindfork/mindfork/server"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

var _ = json.Marshaler(&server.Error{})

type marshalableError struct{}

func (m *marshalableError) MarshalJSON() ([]byte, error) {
	return []byte(`{"some": "error"}`), nil
}
func (m *marshalableError) Error() string { return "some error" }

func (s *ServerSuite) TestErrorMarshalJSON(c *C) {
	for i, t := range []struct {
		should    string
		given     error
		expect    string
		expectErr string
	}{{
		should: "marshal an Error wrapping a nil",
		given:  nil,
		expect: `{"Err":""}`,
	}, {
		should: "marshal a simple error",
		given:  errors.New("an error"),
		expect: `{"Err":"an error"}`,
	}, {
		should: "marshal a json.Marshaler error",
		given:  &marshalableError{},
		expect: `{"Err":{"some":"error"}}`,
	}} {
		c.Logf("test %d: should %s", i, t.should)
		e := &server.Error{Err: t.given}
		bs, err := json.Marshal(e)
		if t.expectErr != "" {
			c.Check(err, ErrorMatches, t.expectErr)
			continue
		}
		c.Assert(err, jc.ErrorIsNil)
		c.Check(string(bs), Matches, t.expect)
	}
}
