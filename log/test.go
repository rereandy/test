package main

import (
	lg "github.com/rereandy/zaplog"
	"time"
)

func main() {
	d3 := time.NewTicker(3 * time.Second)
	d5 := time.NewTicker(5 * time.Second)
	d7 := time.NewTicker(7 * time.Second)
	for {
		select {
		case <-d3.C:
			lg.Infof("test info log")
		case <-d5.C:
			lg.Warnf("test warn log")
		case <-d7.C:
			lg.Errorf("test error log")
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
