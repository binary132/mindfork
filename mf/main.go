package main

import (
	"log"
	"net/http"

	"github.com/mindfork/mindfork/core"
	mfh "github.com/mindfork/mindfork/server/http"

	htr "github.com/julienschmidt/httprouter"
)

func main() {
	httpMux := htr.New()

	path := "localhost:25000"

	log.Printf("Mindfork listening on %v", path)

	// Run the Core as a mindfork.Server using the httprouter.Router.
	err := http.ListenAndServe(
		path,
		mfh.Serve(&core.Core{}, core.MakeMessage)(httpMux, "/"),
	)
	if err != nil {
		log.Panic(err)
	}
}
