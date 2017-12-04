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

package muxlogger

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Logger is the actual implementation of the middleware
func Logger(log *log.Logger, handler http.HandlerFunc) http.HandlerFunc {
	logger := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("INFO: Request: %s/%d.%d %s %s", r.Proto, r.ProtoMajor, r.ProtoMinor, r.Method, r.URL.Path)
		log.Printf("INFO:   Header: %s", logHeader(log, r.Header))
		log.Printf("INFO:   Body: %s", logBody(log, r))
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(logger)
}

// logHeader prints the header
func logHeader(log *log.Logger, header http.Header) string {
	hs := fmt.Sprintf("{ ")
	for name, values := range header {
		for _, value := range values {
			hs = fmt.Sprintf("%s %s: %s, ", hs, name, value)
		}
	}
	return fmt.Sprintf("%s }", hs)
}

// logBody prints the header
func logBody(log *log.Logger, r *http.Request) string {
	body, err := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	if err != nil {
		log.Printf("ERROR: Error sending HTTP request: %v", err)
		return ""
	}
	return string(body)
}
