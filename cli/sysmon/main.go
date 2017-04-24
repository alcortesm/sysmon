package main

import (
	"fmt"
	"log"

	"github.com/alcortesm/sysmon/cpu"
)

func main() {
	load, err := cpu.New()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(load)
}
