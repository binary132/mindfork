package http_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/mindfork/mindfork/core"
	"github.com/mindfork/mindfork/core/message"
	mfh "github.com/mindfork/mindfork/server/http"
	st "github.com/mindfork/mindfork/server/testing"
	mft "github.com/mindfork/mindfork/testing"

	jc "github.com/juju/testing/checkers"
	"github.com/julienschmidt/httprouter"
	. "gopkg.in/check.v1"
)

func (h *HTTPSuite) BenchmarkTesting(c *C) {
	htr := httprouter.New()
	mfh.Serve(&st.Server{}, &mft.MessageMaker{})(htr, "/")

	c.Logf("http echo benchmark: ")
	br := bytes.NewReader([]byte(`{"Type":"test","RawMessage":{"X":5}}`))

	req, err := http.NewRequest(
		"POST",
		"http://example.com/",
		br,
	)
	c.Assert(err, jc.ErrorIsNil)

	for i := 0; i < c.N; i++ {
		w := httptest.NewRecorder()
		htr.ServeHTTP(w, req)
	}
}

func (h *HTTPSuite) BenchmarkCore(c *C) {
	htr := httprouter.New()
	mfh.Serve(&core.Core{}, &message.Maker{})(htr, "/")

	c.Logf("http echo benchmark: ")
	br := bytes.NewReader([]byte(`{"Type":"intention","RawMessage":{"Who":"User","What":"Run a test"}}`))

	req, err := http.NewRequest(
		"POST",
		"http://example.com/",
		br,
	)
	c.Assert(err, jc.ErrorIsNil)

	for i := 0; i < c.N; i++ {
		w := httptest.NewRecorder()
		htr.ServeHTTP(w, req)
	}
}
