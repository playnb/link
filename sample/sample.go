package main

import (
	"github.com/playnb/util/log"
	"time"
)

func init() {
}

func main() {
	log.InitPanic("../tmp")
	log.Init(log.DefaultLogger("../tmp", "run"))
	defer log.Flush()

	ws()
	for {
		time.Sleep(time.Second)
	}
}
