package main

import (
	"log"
	"time"

	"github.com/alcortesm/sysmon"
)

func main() {
	quit := make(chan bool)
	go func() {
		err := sysmon.Server(quit)
		if err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(time.Second * 5)
	quit <- true
}
