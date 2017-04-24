package main

import (
	"fmt"
	"log"

	"github.com/alcortesm/sysmon/loadavg"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

const (
	name  = "com.github.alcortesm.sysmon"
	path  = "/com/github/alcortesm/sysmon"
	intro = `<node name="` + path + `">
	<interface name="` + name + `">
		<method name="OneMinLoadAvg">
			<arg direction="out" type="d"/>
		</method>
		<method name="FiveMinLoadAvg">
			<arg direction="out" type="d"/>
		</method>
		<method name="FifteenMinLoadAvg">
			<arg direction="out" type="d"/>
		</method>
		<method name="RunnableCount">
			<arg direction="out" type="i"/>
		</method>
		<method name="ExistsCount">
			<arg direction="out" type="i"/>
		</method>
		<method name="LastCreatedPID">
			<arg direction="out" type="i"/>
		</method>
	</interface>` + introspect.IntrospectDataString + `</node>`
)

func main() {
	l, err := loadavg.New()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(l)

	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		errClose := conn.Close()
		if err == nil {
			err = errClose
		}
	}()

	fmt.Println(conn.Names())

	reply, err := conn.RequestName(name, dbus.NameFlagDoNotQueue)
	if err != nil {
		log.Fatal(err)
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		log.Fatal("name already taken")
	}

	conn.Export(l, path, name)
	conn.Export(introspect.Introspectable(intro), path,
		"org.freedesktop.DBus.Introspectable")
	fmt.Printf("Listening on %s / %s ...\n", name, path)
	select {}
}
