package main

import (
	"fmt"
	"log"

	"github.com/alcortesm/sysmon/cpu"
)

const loadAvgPath = "/proc/loadavg"

func main() {
	cpu, err := cpu.New(loadAvgPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cpu)
}
