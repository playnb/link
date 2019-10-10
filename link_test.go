package link

import (
	"github.com/playnb/util"
	"github.com/playnb/util/log"
	"testing"
	"time"
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
	serv := &WSServer{}
	var agent *Agent
	serv.OnAccept = func(a *Agent) {
		agent = a
		log.Trace("OnAccept: %v", agent)
	}
	serv.Start(serverOpt)
	for {
		time.Sleep(time.Second)
		if agent != nil {
			data, ok := <-agent.ReadChan()
			if !ok {
				return
			}
			log.Trace("ReadChan: %s", string(data.GetPayload()))
			agent.WriteMsg(util.MakeBuffDataBySlice([]byte("Hello wss"), 0))
		}
	}
}
