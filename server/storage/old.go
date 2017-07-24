package storage

import (
	"fmt"
	"sync"

	"github.com/alcortesm/sysmon/stat"
)

// Old implements Storage by inefficiently copyin slices all around.
type Old struct {
	capacity int
	stats    []stat.CPUer
	mutex    sync.Mutex
	cpuUsage []float64
}

// Returns a new Old server storage of the given capacity.
func NewOld(capacity int) (*Old, error) {
	if capacity < 1 {
		return nil,
			fmt.Errorf("cannot create Old: capacity must be > 0, was %d",
				capacity)
	}
	return &Old{capacity: capacity}, nil
}

// Insert implements Storage.
func (o *Old) Insert(s stat.CPUer) {
	o.mutex.Lock()
	defer o.updateCPUUsage()
	defer o.mutex.Unlock()
	if len(o.stats) < o.capacity {
		o.stats = append([]stat.CPUer{s}, o.stats...)
		return
	}
	copy(o.stats[1:], o.stats[:len(o.stats)-1])
	o.stats[0] = s
}

func (o *Old) updateCPUUsage() {
	if len(o.stats) < 2 {
		return
	}
	percentage := stat.CPUUsage(o.stats[0], o.stats[1])
	if len(o.cpuUsage) < o.capacity-1 {
		o.cpuUsage = append([]float64{percentage}, o.cpuUsage...)
		return
	}
	copy(o.cpuUsage[1:], o.cpuUsage[:len(o.cpuUsage)-1])
	o.cpuUsage[0] = percentage
}

// CPUUsage implements Storage.
func (o *Old) CPUUsage() []float64 {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	ret := make([]float64, len(o.cpuUsage))
	copy(ret, o.cpuUsage)
	return ret
}
