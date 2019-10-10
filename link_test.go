package link

import (
	"github.com/playnb/util/log"
	"testing"
)

func init() {
	log.InitPanic("../tmp")
	log.Init(log.DefaultLogger("../tmp", "run"))
	defer log.Flush()
}

var serverOpt = &ServerOption{
	Addr:            "127.0.0.1:1234",
	MaxMsgLen:       0,
	MaxConnNum:      0,
	PendingWriteNum: 0,
	HTTPTimeout:     0,
	CertFile:        "",
	KeyFile:         "",
}

func Test_WSServer(t *testing.T) {

}
