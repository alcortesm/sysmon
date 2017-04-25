package main

import (
	"log"

	"github.com/alcortesm/sysmon"
)

func main() {
	err := sysmon.Server()
	if err != nil {
		log.Fatal(err)
	}
}
