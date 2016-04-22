package main

import (
	"log"
	"net/http"

	"github.com/mindfork/mindfork/core"
	coremsg "github.com/mindfork/mindfork/core/message"
	mfh "github.com/mindfork/mindfork/server/http"

	htr "github.com/julienschmidt/httprouter"
)

func main() {
	httpMux := htr.New()

	server := &core.Core{Scheduler: core.NewKernel()}

	path := "localhost:25000"

	log.Printf("Mindfork listening on %v", path)

	// httpMux.GET("/dump", func(w http.ResponseWriter, r *http.Request, ps htr.Params) {
	// 	bs, err := json.MarshalIndent(server.Export(), "", "\t")
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	}
	// 	w.Write(bs)
	// })

	// Run the Core as a mindfork.Server using the httprouter.Router.
	err := http.ListenAndServe(
		path,
		mfh.Serve(server, &coremsg.Maker{})(httpMux, "/"),
	)
	if err != nil {
		log.Panic(err)
	}
}
