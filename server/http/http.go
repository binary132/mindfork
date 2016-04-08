package http

import (
	"bytes"
	"fmt"
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
	s server.Server, m message.MessageMaker,
) func(*htr.Router, string) *htr.Router {
	return func(r *htr.Router, path string) *htr.Router {
		r.POST(path, RawBody(s, m))

		return r
	}
}

// RawURL is an httprouter.Handle which handles messages on path, with escaped
// Message contents in the URL.
func RawURL(s server.Server, m message.MessageMaker, path string) htr.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps htr.Params) {
		var (
			de  message.Decoder
			en  = m.NewEncoder(w)
			msg = new(message.Message)
		)

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
		de = m.NewDecoder(bytes.NewReader([]byte(qu)))

		if err := de.Decode(msg); err != nil {
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

		if err := en.Encode(s.Serve(*msg)); err != nil {
			problem := fmt.Sprintf(
				"failed to encode message: %s",
				err,
			)
			http.Error(w, problem, http.StatusInternalServerError)
			return
		}
	}
}

// RawBody is an httprouter.Handle which handles messages on path, with escaped
// Message contents in the URL.
func RawBody(s server.Server, m message.MessageMaker) htr.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps htr.Params) {
		var (
			de  = m.NewDecoder(r.Body)
			en  = m.NewEncoder(w)
			msg = new(message.Message)
		)

		if err := de.Decode(msg); err != nil {
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

		if err := en.Encode(s.Serve(*msg)); err != nil {
			problem := fmt.Sprintf(
				"failed to encode message: %s",
				err,
			)
			http.Error(w, problem, http.StatusInternalServerError)
			return
		}
	}
}
