package sysmon

import (
	"fmt"

	"github.com/alcortesm/sysmon/loadavg"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

const (
	WellKnownBusName = "com.github.alcortesm.sysmon1"
	InterfaceName    = WellKnownBusName
	Path             = "/com/github/alcortesm/sysmon1"
)

func Server() error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	defer func() {
		errClose := conn.Close()
		if err == nil {
			err = errClose
		}
	}()

	fmt.Println(conn.Names())

	reply, err := conn.RequestName(WellKnownBusName, dbus.NameFlagDoNotQueue)
	if err != nil {
		return err
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("name already taken: %s", WellKnownBusName)
	}

	l, err := loadavg.New()
	if err != nil {
		return err
	}
	fmt.Println(l)

	conn.Export(l, Path, InterfaceName)
	conn.Export(introspect.Introspectable(IntrospectDataString),
		Path, "org.freedesktop.DBus.Introspectable")
	fmt.Printf("Listening on %s, %s ...\n", InterfaceName, Path)
	select {}
}
