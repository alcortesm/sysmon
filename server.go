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

func Server(quit <-chan bool) (err error) {
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

	err = claimBusName(conn, WellKnownBusName)
	if err != nil {
		return err
	}
	fmt.Println(conn.Names())

	l, err := loadavg.New()
	if err != nil {
		return err
	}
	fmt.Println(l)

	conn.Export(l, Path, InterfaceName)
	conn.Export(introspect.Introspectable(IntrospectDataString),
		Path, "org.freedesktop.DBus.Introspectable")
	fmt.Printf("Listening on %s...\n", WellKnownBusName)

	_ = <-quit
	return nil
}

func claimBusName(conn *dbus.Conn, name string) error {
	reply, err := conn.RequestName(name, dbus.NameFlagDoNotQueue)
	if err != nil {
		return err
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("bus name already taken: %s", name)
	}

	return nil
}
