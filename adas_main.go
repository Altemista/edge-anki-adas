// Copyright 2018 NTT Group

// Permission is hereby granted, free of charge, to any person obtaining a copy of this
// software and associated documentation files (the "Software"), to deal in the Software
// without restriction, including without limitation the rights to use, copy, modify,
// merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to the following
// conditions:

// The above copyright notice and this permission notice shall be included in all copies
// or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
// PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE
// FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package main

import (
	"runtime"
	"flag"
	"log"
	_ "net/http/pprof"
	"os"

	anki "github.com/okoeth/edge-anki-base"
	"net/http"
	"goji.io/pat"
	"github.com/rs/cors"
	"io/ioutil"
	"fmt"
	"goji.io"
)

// Logging
var mlog = log.New(os.Stdout, "EDGE-ANKI-OVTK: ", log.Lshortfile|log.LstdFlags|log.Lmicroseconds)

func init() {
	flag.Parse()
}

func main() {
	//runtime.GOMAXPROCS(1)
	mlog.Println("Processors: ", runtime.GOMAXPROCS(0))	

	// Set-up routes
	mux := goji.NewMux()

	track := anki.CreateTrack()

	// Set-up channels for status and commands
	cmdCh, statusCh, err := anki.CreateChannels("edge.adas", mux, &track)
	if err != nil {
		mlog.Fatalln("FATAL: Could not establish channels: %s", err)
	}

	// Go and drive cars on track
	go driveCars(track, cmdCh, statusCh)

	//statusCh <- anki.Status{}

	tc := NewAdasController(track, cmdCh)
	tc.AddHandlers(mux)
	mux.Handle(pat.Get("/html/*"), http.FileServer(http.Dir("html/dist/")))
	corsHandler := cors.Default().Handler(mux)


	indexFile, err := os.Open("index.html")
	if err != nil {
		mlog.Println(err)
	}
	index, err := ioutil.ReadAll(indexFile)
	if err != nil {
		mlog.Println(err)
	}

	mux.HandleFunc(pat.Get("/"), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(index))
	})

	mlog.Println("INFO: System is ready.")
	http.ListenAndServe("0.0.0.0:8003", corsHandler)
}
