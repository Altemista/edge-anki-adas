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
	"encoding/json"
	"net/http"

	anki "github.com/okoeth/edge-anki-base"
	"github.com/okoeth/muxlogger"
	"goji.io"
	"goji.io/pat"
)

type (
	// AdasController represents the controller for working with this app
	AdasController struct {
		track []anki.Status
		cmdCh chan anki.Command
	}
)

// NewAdasController provides a reference to an IncomingController
func NewAdasController(t []anki.Status, ch chan anki.Command) *AdasController {
	return &AdasController{track: t, cmdCh: ch}
}

// AddHandlers inserts new greeting
func (oc *AdasController) AddHandlers(mux *goji.Mux) {
	mux.HandleFunc(pat.Get("/v1/adas/status"), oc.GetStatus) // Omitting logger for GetStatus
	mux.HandleFunc(pat.Post("/v1/adas/command"), muxlogger.Logger(mlog, oc.PostCommand))
}

// GetStatus retrieves latest status
func (oc *AdasController) GetStatus(w http.ResponseWriter, r *http.Request) {
	// TODO: There is a race condition around concurrent access to track
	sj, err := json.Marshal(oc.track)
	if err != nil {
		mlog.Println("ERROR: Error de-marshaling TheStatus")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(sj)
}

// PostCommand sends a command via Kafka to the controller
func (oc *AdasController) PostCommand(w http.ResponseWriter, r *http.Request) {
	// Read command from request
	cmd := anki.Command{}
	err := json.NewDecoder(r.Body).Decode(&cmd)
	if err != nil {
		mlog.Printf("ERROR: Error decoding request body: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mlog.Printf("INFO: Sending command to channel")
	oc.cmdCh <- cmd
	mlog.Printf("INFO: Command processed by channel")
	w.WriteHeader(http.StatusOK)
}
