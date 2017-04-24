package main

import (
	"fmt"
	"log"

	"github.com/alcortesm/sysmon/loadavg"
)

func main() {
	l, err := loadavg.New()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(l)

}
