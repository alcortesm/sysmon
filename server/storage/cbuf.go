package storage

import (
	"fmt"
	"sync"

	"github.com/alcortesm/sysmon/stat"
)

// CBuf implements Storage with a thread safe, in-memory, circular buffer.
type CBuf struct {
	cbuf  []stat.CPUer
	first int
	len   int
	mutex sync.Mutex
}

func NewCBuf(capacity int) (*CBuf, error) {
	if capacity < 2 {
		return nil, fmt.Errorf("NewCBuf: capacity must be > 1, was %d",
			capacity)
	}
	return &CBuf{
		cbuf: make([]stat.CPUer, capacity),
	}, nil
}

// Insert implements Storage.
func (c *CBuf) Insert(s stat.CPUer) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.len == cap(c.cbuf) {
		c.cbuf[c.first] = s
		c.first = c.add(c.first, 1)
		return
	}
	c.cbuf[c.add(c.first, c.len)] = s
	c.len++
	return
}

func (c *CBuf) add(a, b int) int {
	return (a + b) % cap(c.cbuf)
}

// CPUUsage implements Storage.
func (c *CBuf) CPUUsage() []float64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.len < 2 {
		return []float64{}
	}
	n := c.len - 1
	ret := make([]float64, n)
	for i := 0; i < n; i++ {
		cur := c.add(c.first, i)
		next := c.add(cur, 1)
		ret[n-i-1] = stat.CPUUsage(c.cbuf[cur], c.cbuf[next])
	}
	return ret
}
