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

	"github.com/wvanbergen/kafka/consumergroup"

	"github.com/Shopify/sarama"
	"github.com/okoeth/edge-anki-base/anki"
	"github.com/rs/cors"
	"goji.io"
	"goji.io/pat"
)

// Variable mlog is the logger for the app
var mlog = log.New(os.Stdout, "EDGE-ANKI-OVTK: ", log.Lshortfile|log.LstdFlags)

// TheStatus carries the latest status information
var TheStatus = [3]anki.Status{}

// TheProducer provides a reference to the Kafka producer
var TheProducer sarama.AsyncProducer

// TheConsumer provides a reference to the Kafka producer
var TheConsumer *consumergroup.ConsumerGroup

func init() {
	flag.Parse()
}

func main() {
	// Initialise Cars
	TheStatus[0].CarNo = "1"
	TheStatus[1].CarNo = "2"
	TheStatus[2].CarNo = "3"

	// Set-up Kafka
	kafkaServer := os.Getenv("KAFKA_SERVER")
	if kafkaServer == "" {
		mlog.Printf("INFO: Using 127.0.0.1 as default KAFKA_SERVER.")
		kafkaServer = "127.0.0.1"
	}
	p, err := anki.CreateKafkaProducer(kafkaServer + ":9092")
	if err != nil {
		mlog.Fatalf("ERROR: Cannot create Kafka producer: %s", err)
	}
	TheProducer = p
	c, err := anki.CreateKafkaConsumer(kafkaServer+":2181", "edge.overtake", TheStatus)
	if err != nil {
		mlog.Fatalf("ERROR: Cannot create Kafka consumer: %s", err)
	}
	TheConsumer = c

	// Go and watch the track
	go watchTrack()

	// Set-up routes
	mux := goji.NewMux()
	tc := NewOvertakeController()
	tc.AddHandlers(mux)
	mux.Handle(pat.Get("/html/*"), http.FileServer(http.Dir("html/dist/")))
	corsHandler := cors.Default().Handler(mux)
	mlog.Printf("INFO: System is ready.\n")
	http.ListenAndServe("0.0.0.0:8002", corsHandler)
}
