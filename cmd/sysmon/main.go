package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	samplingPeriodFlag = durationFlag(time.Second)
	// SamplingPeriod is how often do we sample /proc/stats.
	SamplingPeriod time.Duration

	defaultNSamples = 20
	// NSamples is how many samples to remember.
	NSamples int
)

func init() {
	flag.Var(&samplingPeriodFlag, "p", "sampling period")
	flag.IntVar(&NSamples, "n", defaultNSamples, "number of samples to remember")
}

func main() {
	flag.Usage = usage
	flag.Parse()
	SamplingPeriod = samplingPeriodFlag.duration()

	if len(flag.Args()) != 0 {
		usageAndExit()
	}

	fmt.Println("SamplingPeriod", SamplingPeriod)
	fmt.Println("NSamples", NSamples)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of sysmon:\n")
	fmt.Fprintf(os.Stderr, "\tsysmon [flags]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func usageAndExit() {
	usage()
	os.Exit(2)
}
