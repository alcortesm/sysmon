package cpu

import (
	"fmt"
	"os"
)

// Cpu values contain data about the CPU usage.
type Cpu struct {
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
	scanFormat   = "%f %f %f %d/%d %d"
	stringFormat = "%.2f %.2f %.2f %d/%d %d"
)

// New returns a Cpu value taken by reading the file at path, interpreted in
// /proc/loadavg format.
func New(path string) (_ *Cpu, err error) {
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

	var cpu Cpu
	_, err = fmt.Fscanf(f, scanFormat,
		&cpu.OneMinLoadAvg,
		&cpu.FiveMinLoadAvg,
		&cpu.FifteenMinLoadAvg,
		&cpu.RunnableCount,
		&cpu.ExistCount,
		&cpu.LastCreatedPID,
	)
	if err != nil {
		return nil, err
	}

	return &cpu, nil
}

// String returns a human readable representation of a Cpu value as a string.
func (c *Cpu) String() string {
	return fmt.Sprintf(stringFormat,
		c.OneMinLoadAvg,
		c.FiveMinLoadAvg,
		c.FifteenMinLoadAvg,
		c.RunnableCount,
		c.ExistCount,
		c.LastCreatedPID,
	)
}
