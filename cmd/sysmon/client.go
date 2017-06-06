package main

import "fmt"

func (s *State) client(args ...string) {
	fmt.Printf("sysmon client: %s\n", args)
}
