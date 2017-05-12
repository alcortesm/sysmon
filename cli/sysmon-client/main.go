package main

import (
	"fmt"
	"log"

	"github.com/alcortesm/sysmon"
	"github.com/godbus/dbus"
)

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}
	obj := conn.Object(sysmon.WellKnownBusName, sysmon.Path)
	path := obj.Path()
	fmt.Println(path, path.IsValid())
	call := obj.Call(sysmon.WellKnownBusName+"/OneMinLoadAvg", 0)
	if call.Err != nil {
		log.Fatal(call.Err)
	}
	if err := conn.Close(); err != nil {
		log.Fatal(err)
	}
}
