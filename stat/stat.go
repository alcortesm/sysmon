package stat

import (
	"fmt"
	"os"
)

// S values contain the amount of time spent by the CPU in different
// states, measured in USER_HZ (1/100 secs in most systems).
//
// This assumes Linux > 2.6.33.
type S struct {
	// Time spent in user mode.
	User int
	// Time spent in user mode with low priority.
	Nice int
	// Time spent in system mode.
	System int
	// Time spent in the idle task.
	Idle int
	// Time waiting for I/O to complete.
	IOWait int
	// Time servicing interrupts.
	IRQ int
	// Time servicing soft interrupts.
	SoftIRQ int
	// Stolen time (time spent in other OS when running in a virtualized
	// environment.
	Steal int
	// Time spent running virtual CPU for guest OSs under the control of
	// the kernel.
	Guest int
	// Time spent running a niced guest.
	NiceGuest int
}

const (
	path      = "/proc/stat"
	cpuFormat = "cpu  %d %d %d %d %d %d %d %d %d %d\n"
)

func New() (_ *S, err error) {
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

	var s S
	_, err = fmt.Fscanf(f, cpuFormat,
		&s.User,
		&s.Nice,
		&s.System,
		&s.Idle,
		&s.IOWait,
		&s.IRQ,
		&s.SoftIRQ,
		&s.Steal,
		&s.Guest,
		&s.NiceGuest,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *S) String() string {
	return fmt.Sprintf(cpuFormat,
		s.User,
		s.Nice,
		s.System,
		s.Idle,
		s.IOWait,
		s.IRQ,
		s.SoftIRQ,
		s.Steal,
		s.Guest,
		s.NiceGuest,
	)
}

func (s *S) Total() int {
	return s.User +
		s.Nice +
		s.System +
		s.Idle +
		s.IOWait +
		s.IRQ +
		s.SoftIRQ +
		s.Steal
}

func (s *S) TotalIdle() int {
	return s.Idle + s.IOWait
}

func (s *S) TotalCPU() int {
	return s.Total() - s.TotalIdle()
}
