package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

var commands = map[string]func(*State, ...string){
	"client": (*State).client,
	"server": (*State).server,
}

type State struct{}

func main() {
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of sysmon:\n")
	fmt.Fprintf(os.Stderr, "\tsysmon [command] [flags]\n")
	printCommands()
	flag.PrintDefaults()
}

func usageAndExit(fs *flag.FlagSet) {
	usage()
	os.Exit(2)
}

func printCommands() {
	fmt.Fprintf(os.Stderr, "Sysmon commands:\n")
	var cmdStrs []string
	for cmd := range commands {
		cmdStrs = append(cmdStrs, cmd)
	}
	sort.Strings(cmdStrs)
	for _, cmd := range cmdStrs {
		fmt.Fprintf(os.Stderr, "\t%s\n", cmd)
	}
}
