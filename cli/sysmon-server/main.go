package main

import (
	"log"
	"time"

	"github.com/alcortesm/sysmon/server"
)

func main() {
	s := server.New()
	if err := s.Connect(); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 5)
	if err := s.Disconnect(); err != nil {
		log.Fatal(err)
	}
}
