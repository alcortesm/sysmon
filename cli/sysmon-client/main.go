package main

import (
	"fmt"
	"log"

	"github.com/alcortesm/sysmon"
	"github.com/godbus/dbus"
)

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}
	obj := conn.Object(sysmon.WellKnownBusName, sysmon.Path)
	s, err := loadavgs(obj)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)

	_, err = dev(obj)
	if err != nil {
		log.Fatal(err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal(err)
	}
}

func loadavgs(o dbus.BusObject) (string, error) {
	method := "com.github.alcortesm.sysmon1.LoadAvgs"
	resp := o.Call(method, 0)
	if resp.Err != nil {
		return "", resp.Err
	}
	s, ok := resp.Body[0].(string)
	if !ok {
		return "", fmt.Errorf("response body is not a string")
	}
	return s, nil
}

func dev(o dbus.BusObject) ([]float64, error) {
	method := "com.github.alcortesm.sysmon1.Dev"
	resp := o.Call(method, 0)
	if resp.Err != nil {
		return nil, resp.Err
	}
	//fmt.Println("len of body:", len(resp.Body))
	s, ok := resp.Body[0].([]float64)
	if !ok {
		return nil, fmt.Errorf("response body is not a []float64")
	}
	return s, nil
}
