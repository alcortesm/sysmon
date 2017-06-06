package main

import "fmt"

func (s *State) server(args ...string) {
	fmt.Printf("sysmon server: %s\n", args)
}
