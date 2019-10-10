package main

import (
	"github.com/playnb/util/log"
	"time"
)

func init() {
	log.InitPanic("../tmp")
	log.Init(log.DefaultLogger("../tmp", "run"))
	defer log.Flush()
}

func main() {
	ws()
	for {
		time.Sleep(time.Second)
	}
}
