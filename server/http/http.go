package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/mindfork/mindfork/message"
	"github.com/mindfork/mindfork/server"

	htr "github.com/julienschmidt/httprouter"
)

// Serve returns a function on httprouter.Router that binds the given Mindfork
// Server on the given path handling the Server's Handle.  Note that if Handle
// mutates the Server, this function will then cause mutation.  To avoid races,
// Handle must therefore be synchronized or pure.
func Serve(
	s server.Server, m message.Maker,
) func(*htr.Router, string) *htr.Router {
	return func(r *htr.Router, path string) *htr.Router {
		r.POST(path, RawBody(s, m))

		return r
	}
}

// RawURL is an httprouter.Handle which handles messages on path, with escaped
// Message contents in the URL.
func RawURL(s server.Server, m message.Maker, path string) htr.Handle {
	encode, decode := m.Encoder(), m.Decoder()
	return func(w http.ResponseWriter, r *http.Request, ps htr.Params) {
		q := ps.ByName(path)
		if q == "" {
			http.Error(
				w,
				"no message passed",
				http.StatusBadRequest,
			)
			return
		}

		qu, err := url.QueryUnescape(q)
		if err != nil {
			problem := fmt.Sprintf(
				"failed to unescape URL: %s",
				err,
			)
			http.Error(w, problem, http.StatusBadRequest)
			return
		}

		msg, err := decode([]byte(qu))
		if err != nil {
			problem := fmt.Sprintf(
				"failed to decode message: %s",
				err,
			)
			http.Error(
				w,
				problem,
				http.StatusInternalServerError,
			)
			return
		}

		bs, err := encode(s.Serve(msg))
		if err != nil {
			problem := fmt.Sprintf(
				"failed to encode message: %s",
				err,
			)
			http.Error(w, problem, http.StatusInternalServerError)
			return
		}

		if _, err = w.Write(bs); err != nil {
			problem := fmt.Sprintf(
				"failed to write response: %s",
				err,
			)
			http.Error(
				w,
				problem,
				http.StatusInternalServerError,
			)
		}
	}
}

// RawBody is an httprouter.Handle which handles messages on path, with escaped
// Message contents in the URL.
func RawBody(s server.Server, m message.Maker) htr.Handle {
	encode, decode := m.Encoder(), m.Decoder()
	return func(w http.ResponseWriter, r *http.Request, ps htr.Params) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			problem := fmt.Sprintf(
				"failed to read request body: %s",
				err,
			)
			http.Error(
				w,
				problem,
				http.StatusBadRequest,
			)
			return
		}

		msg, err := decode(body)
		if err != nil {
			problem := fmt.Sprintf(
				"failed to decode message: %s",
				err,
			)
			http.Error(
				w,
				problem,
				http.StatusInternalServerError,
			)
			return
		}

		bs, err := encode(s.Serve(msg))
		if err != nil {
			problem := fmt.Sprintf(
				"failed to encode message: %s",
				err,
			)
			http.Error(w, problem, http.StatusInternalServerError)
			return
		}

		if _, err = w.Write(bs); err != nil {
			problem := fmt.Sprintf(
				"failed to write response: %s",
				err,
			)
			http.Error(
				w,
				problem,
				http.StatusInternalServerError,
			)
		}
	}
}
