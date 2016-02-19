package http_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	mfh "github.com/mindfork/mindfork/server/http"

	. "gopkg.in/check.v1"
)

// Test hooks up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type HTTPSuite struct{}

var _ = Suite(&HTTPSuite{})

// func (h *HTTPSuite) TestServe(c *C) {
// 	handler := func(w http.ResponseWriter, r *http.Request) {
// 		http.Error(w, "something failed", http.StatusInternalServerError)
// 	}
//
// 	req, err := http.NewRequest("POST", "http://example.com/foo", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	w := httptest.NewRecorder()
// 	handler(w, req)
//
// 	fmt.Printf("%d - %s", w.Code, w.Body.String())
// }

func (h *HTTPSuite) TestWriteError(c *C) {
	for i, t := range []struct {
		should string
		err    string
		expect string
	}{{
		should: "write an error response",
		err:    "oops",
		expect: `{"Err":"oops"}`,
	}} {
		c.Logf("test %d: should %s", i, t.should)
		w := httptest.NewRecorder()

		mfh.WriteError(w, errors.New(t.err))

		c.Check(w.Body.String(), Matches, t.expect)
	}
}
