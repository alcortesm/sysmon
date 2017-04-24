// Package loadavg allows to read the contents of the /proc/loadavg file.
package loadavg

import (
	"fmt"
	"os"
)

// L values contain data read from /proc/loadavg.
type L struct {
	// OneMinLoadAvg is the average of the number of jobs in the run
	// queue during the last minute.
	OneMinLoadAvg float64
	// FiveMinLoadAvg is the average of the number of jobs in the run
	// queue during the last 5 minutes.
	FiveMinLoadAvg float64
	// OneMinLoadAvg is the average of the number of jobs in the run
	// queue during the last 15 minutes.
	FifteenMinLoadAvg float64
	// The number of currently runnable kernel scheduling entities
	// (processes, threads).
	RunnableCount int
	// The number of kernel scheduling entities that currently exist on
	// the system.
	ExistCount int
	// The PID of the process that  was  most recently created on the
	// system.
	LastCreatedPID int
}

const (
	path         = "/proc/loadavg"
	scanFormat   = "%f %f %f %d/%d %d"
	stringFormat = "%.2f %.2f %.2f %d/%d %d"
)

// New access /proc/loadavg and returns an L value with its data.
func New() (_ *L, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		errClose := f.Close()
		if err == nil {
			err = errClose
		}
	}()

	var l L
	_, err = fmt.Fscanf(f, scanFormat,
		&l.OneMinLoadAvg,
		&l.FiveMinLoadAvg,
		&l.FifteenMinLoadAvg,
		&l.RunnableCount,
		&l.ExistCount,
		&l.LastCreatedPID,
	)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

// String returns the data in L as a string in the same format as in
// /proc/loadavg.
func (c *L) String() string {
	return fmt.Sprintf(stringFormat,
		c.OneMinLoadAvg,
		c.FiveMinLoadAvg,
		c.FifteenMinLoadAvg,
		c.RunnableCount,
		c.ExistCount,
		c.LastCreatedPID,
	)
}
