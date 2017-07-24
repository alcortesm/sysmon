package stat

// CPUer represents the time units spent since last boot and the
// ones spent in CPU tasks.  This two values allows to calculate CPU
// usage percentages.
type CPUer interface {
	// The total amount of time units spent in CPU since last boot.
	TotalCPU() int
	// The total amount of time units spent since last boot.
	Total() int
}

// CPUUsage calculates the CPU usage percentage during the period
// between the two given CPUusagers.
func CPUUsage(cur, prev CPUer) float64 {
	cpu := cur.TotalCPU() - prev.TotalCPU()
	total := cur.Total() - prev.Total()
	percentage := 100 * float64(cpu) / float64(total)
	if percentage > 100 {
		return 100
	} else if percentage < 0 {
		return 0
	}
	return percentage
}
