package server

import (
	"fmt"
	"log"

	"github.com/alcortesm/sysmon"
	"github.com/alcortesm/sysmon/loadavg"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

// Server implements sysmon.Server.
type Server struct {
	conn *dbus.Conn
}

// New creates a Server.
func New() sysmon.Server {
	return &Server{}
}

func (s *Server) Connect() error {
	if s.conn != nil {
		return fmt.Errorf("already connected")
	}
	var err error
	if s.conn, err = dbus.SessionBus(); err != nil {
		return err
	}
	if err = claimBusName(s.conn, sysmon.WellKnownBusName); err != nil {
		return err
	}
	s.conn.Export(s, sysmon.Path, sysmon.InterfaceName)
	s.conn.Export(introspect.Introspectable(sysmon.IntrospectDataString),
		sysmon.Path, "org.freedesktop.DBus.Introspectable")
	fmt.Printf("Listening on %s...\n", sysmon.WellKnownBusName)
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

func (s *Server) Disconnect() error {
	if s.conn == nil {
		return fmt.Errorf("not connected")
	}

	return s.conn.Close()
}

func (s *Server) LoadAvgs() (string, *dbus.Error) {
	log.Println("Foo called")
	l, err := loadavg.New()
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}
	return l.String(), nil
}

func (s *Server) Dev() ([]float64, *dbus.Error) {
	return []float64{1.0, 1.1, 1.2}, nil
}
