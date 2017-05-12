package main

import (
	"fmt"
	"log"

	"github.com/alcortesm/sysmon"
	"github.com/godbus/dbus"
	"github.com/joliv/spark"
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

var tail = len(spark.Line([]float64{0.0, 1.0}))

func report(ff []float64) string {
	withMinAndMax := append(ff, 0.0, 1.0)
	plot := spark.Line(withMinAndMax)
	return plot[:len(plot)-tail]
}
