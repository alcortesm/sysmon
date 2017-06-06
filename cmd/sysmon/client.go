package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/alcortesm/sysmon"

	"github.com/godbus/dbus"
	"github.com/joliv/spark"
)

func client() error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	if err := makeSureServerIsRunning(conn); err != nil {
		return err
	}
	obj := conn.Object(sysmon.WellKnownBusName, sysmon.Path)
	ff, err := cpuUsageHistory(obj)
	if err != nil {
		return err
	}
	fmt.Println(report(ff))
	if err := conn.Close(); err != nil {
		return err
	}

	return nil
}

const listNames = "org.freedesktop.DBus.ListNames"

// TODO: search for a better way to know if the server is running
func isServerRunning(c *dbus.Conn) (bool, error) {
	var s []string
	err := c.BusObject().Call(listNames, 0).Store(&s)
	if err != nil {
		return false,
			fmt.Errorf("failed to get list of session dbus owned names:", err)
	}

	for _, v := range s {
		if v == sysmon.WellKnownBusName {
			return true, nil
		}
	}

	return false, nil
}

func cpuUsageHistory(o dbus.BusObject) ([]float64, error) {
	resp := o.Call(sysmon.CPUsUsageHistoryMethod, 0)
	if resp.Err != nil {
		return nil, resp.Err
	}
	if len(resp.Body) != 1 {
		return nil, fmt.Errorf("length of resp.Body should be 1 but is %d",
			len(resp.Body))
	}
	ff, ok := resp.Body[0].([]float64)
	if !ok {
		return nil, fmt.Errorf("response body is not a []float64")
	}
	return ff, nil
}

var tail = len(spark.Line([]float64{0.0, 1.0}))

func report(ff []float64) string {
	withMinAndMax := append(ff, 0, 100)
	plot := spark.Line(withMinAndMax)
	return plot[:len(plot)-tail]
}

func makeSureServerIsRunning(c *dbus.Conn) error {
	ok, err := isServerRunning(c)
	if err != nil {
		return err
	}
	if !ok {
		if err := runServer(); err != nil {
			return err
		}
		return waitForServer(c)
	}
	return nil
}

func runServer() error {
	cmd := exec.Command("sysmon-server")
	err := cmd.Start()
	if err != nil {
		return err
	}
	cmd.Process.Release()
	return nil
}

var (
	timeoutDuration    = time.Second
	poolServerDuration = 100 * time.Millisecond // TODO: implemente a backoff here
)

func waitForServer(c *dbus.Conn) error {
	timeout := time.After(timeoutDuration)
	for {
		pool := time.After(poolServerDuration)
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for server")
		case <-pool:
			ok, err := isServerRunning(c)
			if err != nil {
				return err
			}
			if ok {
				return nil
			}
		}
	}
}
