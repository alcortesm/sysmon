package storage

import "github.com/alcortesm/sysmon/stat"

// Storage represents a collection of recent samples of sysmon stats.
type Storage interface {
	// Insert adds a new stat sample to the collection and drops outdated ones.
	Insert(stat.CPUer)
	// CPUUsage calculates the recent values of CPU usage utilization
	// based on the stats samples in the collection and returns them
	// sorted from recent to older.
	CPUUsage() []float64
}
