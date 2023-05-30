package main

import (
	lg "github.com/rereandy/log"
)

var log lg.Logger

func init() {
	log = lg.NewLogger("test")
}

func main() {
	log.Infof("test")
}
