package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/alcortesm/sysmon/server"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	nSamples := 20
	samplingPeriod := time.Second
	s := server.New(nSamples, samplingPeriod)
	if err := s.Connect(); err != nil {
		log.Fatal(err)
	}

	<-c // block until os.Interrupt is receieved
	if err := s.Disconnect(); err != nil {
		log.Fatal(err)
	}
}
