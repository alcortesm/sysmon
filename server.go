package sysmon

import (
	"fmt"

	"github.com/alcortesm/sysmon/loadavg"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

const (
	Name = "com.github.alcortesm.sysmon"
	Path = "/com/github/alcortesm/sysmon"
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

	reply, err := conn.RequestName(Name, dbus.NameFlagDoNotQueue)
	if err != nil {
		return err
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("name already taken: %s", Name)
	}

	l, err := loadavg.New()
	if err != nil {
		return err
	}
	fmt.Println(l)

	conn.Export(l, Path, Name)
	conn.Export(introspect.Introspectable(IntrospectDataString),
		Path, "org.freedesktop.DBus.Introspectable")
	fmt.Printf("Listening on %s, %s ...\n", Name, Path)
	select {}
}
