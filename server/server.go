package server

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/alcortesm/sysmon"
	"github.com/alcortesm/sysmon/stat"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

// Server implements sysmon.Server.
type Server struct {
	// amount of samples to remember
	nSamples int
	// sampling period
	period   time.Duration
	conn     *dbus.Conn
	mutex    *sync.Mutex
	stats    []*stat.S
	cpuUsage []float64
}

// New creates a new Server that samples /proc/loadavg every "period"
// seconds and remembers "nSamples" samples.
func New(nSamples int, period time.Duration) sysmon.Server {
	return &Server{
		nSamples: nSamples,
		period:   period,
		mutex:    &sync.Mutex{},
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
	fmt.Println("LoadAvgs cpu usage = ", s.cpuUsage)
	ret := make([]float64, len(s.cpuUsage))
	copy(ret, s.cpuUsage)
	return ret, nil
}

func (s *Server) run() {
	for {
		time.Sleep(s.period)
		st, err := stat.New()
		if err != nil {
			log.Fatal(err)
		}
		s.add(st)
	}
}

func (s *Server) add(st *stat.S) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	defer s.updateCPUUsage()
	if len(s.stats) < s.nSamples {
		s.stats = append([]*stat.S{st}, s.stats...)
		return
	}
	copy(s.stats[1:], s.stats[:len(s.stats)-1])
	s.stats[0] = st
}

func (s *Server) updateCPUUsage() {
	if len(s.stats) < 2 {
		return
	}

	current, _ := s.stats[0], s.stats[1]
	percentage := float64(current.TotalCPU()) / float64(current.Total())
	percentage *= 100
	fmt.Println(percentage, "%")

	if len(s.cpuUsage) < s.nSamples-1 {
		s.cpuUsage = append([]float64{percentage}, s.cpuUsage...)
		return
	}
	copy(s.cpuUsage[1:], s.cpuUsage[:len(s.cpuUsage)-1])
	s.cpuUsage[0] = percentage
}
