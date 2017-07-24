package server

import (
	"fmt"
	"log"
	"time"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"

	"github.com/alcortesm/sysmon"
	"github.com/alcortesm/sysmon/server/storage"
	"github.com/alcortesm/sysmon/stat"
)

// Server implements sysmon.Server.
type Server struct {
	// sampling period
	period  time.Duration
	conn    *dbus.Conn
	storage storage.Storage
}

// New creates a new sysmon.Server that samples /proc/stats every "period"
// and remembers "nSamples" samples.
func New(nSamples int, period time.Duration) (sysmon.Server, error) {
	s, err := storage.NewOld(nSamples)
	if err != nil {
		return nil, err
	}

	return &Server{
		period:  period,
		storage: s,
	}, nil
}

// Connect implements sysmon.Server.
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
	go s.run()
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

// Disconnect implements sysmon.Server.
func (s *Server) Disconnect() error {
	if s.conn == nil {
		return fmt.Errorf("not connected")
	}

	return s.conn.Close()
}

// CPUsUsageHistory implements sysmon.Server.
func (s *Server) CPUsUsageHistory() ([]float64, *dbus.Error) {
	return s.storage.CPUUsage(), nil
}

func (s *Server) run() {
	for {
		time.Sleep(s.period)
		st, err := stat.New()
		if err != nil {
			log.Fatal(err)
		}
		s.storage.Insert(st)
	}
}
