package server

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/alcortesm/sysmon"
	"github.com/alcortesm/sysmon/loadavg"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

const (
	// amount of samples to remember
	nSamples = 8
	// sampling period
	period = 3 * time.Second
)

// Server implements sysmon.Server.
type Server struct {
	conn    *dbus.Conn
	mutex   *sync.Mutex
	samples []float64
}

// New creates a new Server.
func New() sysmon.Server {
	return &Server{
		mutex: &sync.Mutex{},
	}
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

// LoadAvgs implements sysmon.Server.
func (s *Server) LoadAvgs() ([]float64, *dbus.Error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	ret := make([]float64, len(s.samples))
	copy(ret, s.samples)
	return ret, nil
}

func (s *Server) run() {
	for {
		time.Sleep(period)
		l, err := loadavg.New()
		if err != nil {
			log.Fatal(err)
		}
		s.add(l.OneMinLoadAvg)
	}
}

func (s *Server) add(f float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.samples) < nSamples {
		s.samples = append([]float64{f}, s.samples...)
		return
	}
	copy(s.samples[1:], s.samples[:len(s.samples)-1])
	s.samples[0] = f
}
