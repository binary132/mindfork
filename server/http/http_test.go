package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/mindfork/mindfork/core"
	coremsg "github.com/mindfork/mindfork/core/message"
	coretest "github.com/mindfork/mindfork/core/testing"
	"github.com/mindfork/mindfork/message"
	"github.com/mindfork/mindfork/server"
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

var cores = map[string]server.Server{
	"testing": &st.Server{},
}

func (h *HTTPSuite) TestServe(c *C) {
	tNow := time.Now()
	tNowJson, _ := json.Marshal(tNow)

	for i, test := range []struct {
		should         string
		server         server.Server
		maker          message.Maker
		path           string
		arg            string
		expectBody     string
		expectMessages []message.Message
		expectCode     int
	}{{
		should: "fail on a broken message",
		server: cores["testing"],
		maker:  &mft.MessageMaker{},
		path:   "",
		arg:    `{"Type":foo"}`,
		expectBody: `failed to decode message: invalid character 'o'` +
			` in literal false (expecting 'a')` + "\n",
		expectMessages: nil,
		expectCode:     500,
	}, {
		should:         "echo a message",
		server:         cores["testing"],
		maker:          &mft.MessageMaker{},
		path:           "",
		arg:            `{"Type":"test","RawMessage":{"X":5}}`,
		expectBody:     `{"X":5}`,
		expectMessages: []message.Message{mft.Message{X: 5}},
		expectCode:     200,
	}, {
		should:     "echo another message",
		server:     cores["testing"],
		maker:      &mft.MessageMaker{},
		path:       "",
		arg:        `{"Type":"test","RawMessage":{"X":6}}`,
		expectBody: `{"X":6}`,
		expectMessages: []message.Message{
			mft.Message{X: 5}, mft.Message{X: 6},
		},
		expectCode: 200,
	}, {
		should:     "echo a message via core.Core",
		server:     &core.Core{Timer: coretest.TestTimer(tNow)},
		maker:      &coremsg.Maker{},
		path:       "",
		arg:        `{"Type":"echo"}`,
		expectBody: fmt.Sprintf(`{"When":%s}`, tNowJson),
		expectCode: 200,
	}} {
		c.Logf("test %d: should %s", i, test.should)
		c.Logf("  arg: %s", test.arg)

		htr := httprouter.New()
		mfh.Serve(test.server, test.maker)(htr, "/")

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
		c.Check(w.Body.String(), Equals, test.expectBody)
		if testServer, ok := test.server.(*st.Server); ok {
			c.Check(testServer.Messages, jc.DeepEquals, test.expectMessages)
		}

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
			` in literal false (expecting 'a')` + "\n",
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
		c.Check(w.Body.String(), Equals, test.expectBody)
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
			` in literal false (expecting 'a')` + "\n",
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
		c.Check(w.Body.String(), Equals, test.expectBody)
		c.Check(w.Code, Equals, test.expectCode)
	}
}
