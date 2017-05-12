package main

import (
	"bytes"
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
	ff, err := loadavgs(obj)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(report(ff))
	if err := conn.Close(); err != nil {
		log.Fatal(err)
	}
}

func loadavgs(o dbus.BusObject) ([]float64, error) {
	method := "com.github.alcortesm.sysmon1.LoadAvgs"
	resp := o.Call(method, 0)
	if resp.Err != nil {
		return nil, resp.Err
	}
	if len(resp.Body) != 1 {
		return nil, fmt.Errorf("length of resp.Body should be 1 but is %d",
			len(resp.Body))
	}
	ff, ok := resp.Body[0].([]float64)
	if !ok {
		return nil, fmt.Errorf("response body is not a []float64")
	}
	return ff, nil
}

func report(ff []float64) string {
	var buf bytes.Buffer
	sep := ""
	for _, f := range ff {
		buf.WriteString(sep)
		fmt.Fprintf(&buf, "%3.2f", f)
		sep = " "
	}
	return buf.String()
}
