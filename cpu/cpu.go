package cpu

import (
	"fmt"
	"os"
)

const path = "/proc/loadavg"

// Cpu values contain data about the CPU usage.
type Cpu struct {
	// OneMinLoadAvg is the average of the number of jobs in the run
	// queue during the last minute.
	OneMinLoadAvg float64
}

func New() (_ *Cpu, err error) {
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

	var parsed float64
	_, err = fmt.Fscan(f, &parsed)
	if err != nil {
		return nil, err
	}

	return &Cpu{
		OneMinLoadAvg: parsed,
	}, nil
}

func (c *Cpu) String() string {
	format := "cpu average load (1 minute) = %.0f%%"
	return fmt.Sprintf(format, 100*c.OneMinLoadAvg)
}
