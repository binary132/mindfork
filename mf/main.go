package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mindfork/mindfork/core"
	coremsg "github.com/mindfork/mindfork/core/message"
	"github.com/mindfork/mindfork/core/scheduler"
	mfh "github.com/mindfork/mindfork/server/http"

	htr "github.com/julienschmidt/httprouter"
)

func main() {
	httpMux := htr.New()

	server := &core.Core{Scheduler: scheduler.NewKernel()}

	path := "localhost:25000"

	log.Printf("Mindfork listening on %v", path)

	httpMux.GET("/debug/:f", func(w http.ResponseWriter, r *http.Request, ps htr.Params) {
		var f func() []coremsg.Intention

		switch ps.ByName("f") {
		case "":
			fallthrough
		case "export":
			f = server.Export
		case "avail":
			f = server.Available
		default:
			http.Error(w, "unknown command", http.StatusBadRequest)
			return
		}

		bs, err := json.MarshalIndent(f(), "", "\t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(bs)
	})

	// Run the Core as a mindfork.Server using the httprouter.Router.
	err := http.ListenAndServe(
		path,
		mfh.Serve(server, &coremsg.Maker{})(httpMux, "/"),
	)
	if err != nil {
		log.Panic(err)
	}
}
