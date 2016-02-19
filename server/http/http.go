package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	mf "github.com/mindfork/mindfork"
	"github.com/mindfork/mindfork/server"

	htr "github.com/julienschmidt/httprouter"
)

// Serve returns a function on httprouter.Router that binds the given Mindfork
// Server on the given path handling the Server's Handle.  Note that if Handle
// mutates the Server, this function will then cause mutation.  To avoid races,
// Handle must therefore be synchronized or pure.
func Serve(
	s server.Server, m mf.MessageMaker,
) func(*htr.Router, string) *htr.Router {
	return func(r *htr.Router, path string) *htr.Router {
		if path == "/" {
			path = ""
		}

		// r.POST(path+"/:type/:m", RESTfully(s, m, "type"))
		r.POST(path+"/:m", Raw(s, m))

		return r
	}
}

func Raw(s server.Server, m mf.MessageMaker) htr.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps htr.Params) {
		q := ps.ByName("m")

		qu, err := url.QueryUnescape(q)
		if err != nil {
			WriteError(w, err)
		}

		bs, err := server.Serve(s, m, []byte(qu))
		if err != nil {
			WriteError(w, err)
		}

		if _, err = w.Write(bs); err != nil {
			WriteError(w, err)
		}
		if q != "" {
			WriteError(w, errors.New(
				"no message received",
			))

			return
		}
	}
}

// // RESTfully is an httprouter.Handle handling the given Server's Handle with the
// // given param name.
// func RESTfully(s server.Server, m mf.MessageMaker, name string) htr.Handle {
// 	return func(w http.ResponseWriter, r *http.Request, ps htr.Params) {
// 		// First find out whether the type was passed.
// 		if t := ps.ByName(name); t != "" {
// 			mbs, err := json.Marshal(mf.MessageBytes{
// 				Type: t,
// 				RawMessage: ps.ByName("name string")
// 			}
// 			bs, err := server.Serve(s, m, []byte(msg))
// 			if err != nil {
// 				WriteError(w, err)
// 			}
//
// 			if _, err = w.Write(bs); err != nil {
// 				WriteError(w, err)
// 			}
//
// 			return
// 		}
//
// 		WriteError(w, errors.New(
// 			"message must be passed as /<type>/msg",
// 		))
// 	}
// }

// WriteError writes an error to the ResponseWriter as a server.Error Message.
func WriteError(w http.ResponseWriter, e error) {
	bs, err := json.Marshal(server.Error{Err: e})
	if err != nil {
		http.Error(w, fmt.Sprintf(
			"failed to marshal error message: %s", err.Error(),
		), http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(bs); err != nil {
		http.Error(w, fmt.Sprintf(
			"failed to write error message: %s", err.Error(),
		), http.StatusInternalServerError)
	}
}
