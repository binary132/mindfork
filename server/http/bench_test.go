package http_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/mindfork/mindfork/core"
	"github.com/mindfork/mindfork/core/message"
	coretest "github.com/mindfork/mindfork/core/testing"
	mfm "github.com/mindfork/mindfork/message"
	"github.com/mindfork/mindfork/server"
	mfh "github.com/mindfork/mindfork/server/http"
	"github.com/mindfork/mindfork/server/testing"
	mft "github.com/mindfork/mindfork/testing"

	jc "github.com/juju/testing/checkers"
	"github.com/julienschmidt/httprouter"
	. "gopkg.in/check.v1"
)

func benchHTTPServer(c *C, sv server.Server, m mfm.Maker, bs []byte) {
	htr := httprouter.New()

	mfh.Serve(sv, m)(htr, "/")

	req, err := http.NewRequest(
		"POST",
		"http://example.com/",
		bytes.NewReader(bs),
	)
	c.Assert(err, jc.ErrorIsNil)

	for i := 0; i < c.N; i++ {
		req.Body = ioutil.NopCloser(bytes.NewReader(bs))

		w := httptest.NewRecorder()
		htr.ServeHTTP(w, req)

		c.Assert(w.Code, Equals, 200)
	}
}

func (h *HTTPSuite) BenchmarkTesting(c *C) {
	benchHTTPServer(
		c,
		&testing.Server{},
		&mft.MessageMaker{},
		[]byte(`{"Type":"test","RawMessage":{"X":5}}`),
	)
}

func (h *HTTPSuite) BenchmarkCoreIntention(c *C) {
	benchHTTPServer(
		c,
		&core.Core{Scheduler: core.NewKernel()},
		&message.Maker{},
		[]byte(`{"Type":"intention","RawMessage":`+
			coretest.SampleMessages("validIntention")+`}`),
	)
}

func (h *HTTPSuite) BenchmarkCoreEcho(c *C) {
	benchHTTPServer(
		c,
		&core.Core{Scheduler: core.NewKernel()},
		&message.Maker{},
		[]byte(`{"Type":"echo"}`),
	)
}
