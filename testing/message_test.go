package testing_test

import (
	"testing"

	mf "github.com/mindfork/mindfork"
	t "github.com/mindfork/mindfork/testing"

	. "gopkg.in/check.v1"
)

var (
	_ = mf.MessageMaker(&t.MessageMaker{})
	_ = mf.Message(&t.Message{})
)

// Test hooks up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type TestingSuite struct{}

var _ = Suite(&TestingSuite{})

// func (ts *TestingSuite) TestEncoder(c *C) {
// 	var bs bytes.Buffer
// 	en := (&t.MessageMaker{}).NewEncoder(&bs)
//
// 	for i, test := range []struct {
// 		should    string
// 		given     t.Message
// 		expect    string
// 		expectErr string
// 	}{{
// 		should: "make and exercise a new json Encoder OK",
// 		given:  t.Message{X: 5},
// 		expect: `{"X":5}`,
// 	}, {
// 		should: "reset and write a new Message",
// 		given:  t.Message{S: "hello"},
// 		expect: `{"X":0,"S":"hello"}`,
// 	}} {
// 		c.Logf("test %d: should %s", i, test.should)
// 		bs.Reset()
//
// 		err := en.Encode(test.given)
// 		if test.expectErr != "" {
// 			c.Check(err, ErrorMatches, test.expectErr)
// 			continue
// 		}
//
// 		c.Assert(err, jc.ErrorIsNil)
// 		c.Check(string(bs.Bytes()), Matches, test.expect+"\n")
// 	}
// }
//
// func (ts *TestingSuite) TestDecoder(c *C) {
// 	var bs bytes.Buffer
// 	de := (&t.MessageMaker{}).NewDecoder(&bs)
//
// 	for i, test := range []struct {
// 		should    string
// 		given     string
// 		expect    mf.Message
// 		expectErr string
// 	}{{
// 		should: "make and exercise a new json Decoder OK",
// 		given:  `{"Type":"test","RawMessage":{"S":"hello","X":5}}`,
// 		expect: t.Message{S: "hello", X: 5},
// 	}, {
// 		should: "reset and reuse the Decoder",
// 		given:  `{"Type":"test","RawMessage":{"S":"foo"}}`,
// 		expect: t.Message{S: "foo", X: 0},
// 	}, {
// 		should:    "reset and have an error",
// 		given:     `{"Type":"","RawMessage":{"S":"foo"}}`,
// 		expectErr: "no message Type received",
// 	}, {
// 		should: "reset and have no error",
// 		given:  `{"Type":"test","RawMessage":{"X":3}}`,
// 		expect: t.Message{X: 3},
// 	}} {
// 		c.Logf("test %d: should %s", i, test.should)
//
// 		bs.Reset()
// 		_, err := bs.WriteString(test.given)
// 		c.Assert(err, jc.ErrorIsNil)
//
// 		c.Logf("  using bytes: %s", bs.String())
//
// 		var m t.Message
// 		err = de.Decode(&m)
// 		if test.expectErr != "" {
// 			c.Check(err, ErrorMatches, test.expectErr)
// 			continue
// 		}
//
// 		c.Assert(err, jc.ErrorIsNil)
// 		c.Check(m, jc.DeepEquals, test.expect)
// 	}
// }
