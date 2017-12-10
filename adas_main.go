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
	"flag"
	"log"
	"net/http"
	"os"

	anki "github.com/okoeth/edge-anki-base"
	"github.com/rs/cors"
	"goji.io"
	"goji.io/pat"
)

// Logging
var mlog = log.New(os.Stdout, "EDGE-ANKI-OVTK: ", log.Lshortfile|log.LstdFlags)

func init() {
	flag.Parse()
}

func main() {
	// Set-up channels for status and commands
	cmdCh, statusCh, err := anki.CreateChannels("edge.adas")
	if err != nil {
		mlog.Fatalln("FATAL: Could not establish channels: %s", err)
	}
	track := anki.CreateTrack()

	// Go and drive cars on track
	go driveCars(track, cmdCh, statusCh)

	// Set-up routes
	mux := goji.NewMux()
	tc := NewAdasController(track, cmdCh)
	tc.AddHandlers(mux)
	mux.Handle(pat.Get("/html/*"), http.FileServer(http.Dir("html/dist/")))
	corsHandler := cors.Default().Handler(mux)
	mlog.Println("INFO: System is ready.")
	http.ListenAndServe("0.0.0.0:8003", corsHandler)
}
