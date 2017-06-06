package main

import "time"

type durationFlag time.Duration

func (d *durationFlag) String() string {
	return time.Duration(*d).String()
}

func (d *durationFlag) Set(value string) error {
	tmp, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	*d = durationFlag(tmp)
	return nil
}

func (d *durationFlag) duration() time.Duration {
	return time.Duration(*d)
}
