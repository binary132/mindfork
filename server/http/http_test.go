package http_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	mfh "github.com/mindfork/mindfork/server/http"
	st "github.com/mindfork/mindfork/server/testing"
	mft "github.com/mindfork/mindfork/testing"

	jc "github.com/juju/testing/checkers"
	"github.com/julienschmidt/httprouter"
	. "gopkg.in/check.v1"
)

// Test hooks up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type HTTPSuite struct{}

var _ = Suite(&HTTPSuite{})

func (h *HTTPSuite) TestServe(c *C) {
	htr := httprouter.New()
	mfh.Serve(&st.Server{}, &mft.MessageMaker{})(htr, "/")

	for i, test := range []struct {
		should     string
		path       string
		arg        string
		expectBody string
		expectCode int
	}{{
		should: "fail on a broken message",
		path:   "",
		arg:    `{"Type":foo"}`,
		expectBody: `failed to decode message: invalid character 'o'` +
			` in literal false (expecting 'a')`,
		expectCode: 500,
	}, {
		should:     "echo a message",
		path:       "",
		arg:        `{"Type":"test","RawMessage":{"X":5}}`,
		expectBody: `{"X":5}`,
		expectCode: 200,
	}} {
		c.Logf("test %d: should %s", i, test.should)
		c.Logf("  arg: %s", test.arg)

		w := httptest.NewRecorder()
		query := fmt.Sprintf("http://example.com/%s",
			test.path,
		)
		c.Logf("  query: %s", query)
		c.Logf("  body: %s", test.arg)

		req, err := http.NewRequest(
			"POST",
			query,
			bytes.NewReader([]byte(test.arg)),
		)
		c.Assert(err, jc.ErrorIsNil)

		htr.ServeHTTP(w, req)
		c.Check(w.Body.String(), Equals, test.expectBody+"\n")
		c.Check(w.Code, Equals, test.expectCode)
	}
}

func (h *HTTPSuite) TestRawURL(c *C) {
	htr := httprouter.New()
	htr.POST("/:m", mfh.RawURL(&st.Server{}, &mft.MessageMaker{}, "m"))

	for i, test := range []struct {
		should     string
		path       string
		arg        string
		expectBody string
		expectCode int
	}{{
		should: "fail on a broken message",
		path:   "",
		arg:    `{"Type":foo"}`,
		expectBody: `failed to decode message: invalid character 'o'` +
			` in literal false (expecting 'a')`,
		expectCode: 500,
	}, {
		should:     "echo a message",
		path:       "",
		arg:        `{"Type":"test","RawMessage":{"X":5}}`,
		expectBody: `{"X":5}`,
		expectCode: 200,
	}} {
		c.Logf("test %d: should %s", i, test.should)
		c.Logf("  arg: %s", test.arg)

		w := httptest.NewRecorder()
		query := fmt.Sprintf("http://example.com/%s%s",
			test.path,
			url.QueryEscape(test.arg),
		)
		c.Logf("  query: %s", query)

		req, err := http.NewRequest(
			"POST",
			query,
			nil,
		)
		c.Assert(err, jc.ErrorIsNil)

		htr.ServeHTTP(w, req)
		c.Check(w.Body.String(), Equals, test.expectBody+"\n")
		c.Check(w.Code, Equals, test.expectCode)
	}
}

func (h *HTTPSuite) TestRawBody(c *C) {
	htr := httprouter.New()
	htr.POST("/", mfh.RawBody(&st.Server{}, &mft.MessageMaker{}))

	for i, test := range []struct {
		should     string
		path       string
		arg        string
		expectBody string
		expectCode int
	}{{
		should: "fail on a broken message",
		path:   "",
		arg:    `{"Type":foo"}`,
		expectBody: `failed to decode message: invalid character 'o'` +
			` in literal false (expecting 'a')`,
		expectCode: 500,
	}, {
		should:     "echo a message",
		path:       "",
		arg:        `{"Type":"test","RawMessage":{"X":5}}`,
		expectBody: `{"X":5}`,
		expectCode: 200,
	}} {
		c.Logf("test %d: should %s", i, test.should)

		w := httptest.NewRecorder()
		query := fmt.Sprintf("http://example.com/%s",
			test.path,
		)
		c.Logf("  query: %s", query)
		c.Logf("  body: %s", test.arg)

		req, err := http.NewRequest(
			"POST",
			query,
			bytes.NewReader([]byte(test.arg)),
		)
		c.Assert(err, jc.ErrorIsNil)

		htr.ServeHTTP(w, req)
		c.Check(w.Body.String(), Equals, test.expectBody+"\n")
		c.Check(w.Code, Equals, test.expectCode)
	}
}
