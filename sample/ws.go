package main

import (
	"github.com/playnb/link"
	"github.com/playnb/util"
	"github.com/playnb/util/log"
)

var serverOpt = &link.ServerOption{
	Addr:            ":1234",
	MaxMsgLen:       0,
	MaxConnNum:      0,
	PendingWriteNum: 0,
	HTTPTimeout:     0,
	CertFile:        "",
	KeyFile:         "",
	RelativePath:    "/echo",
	GinLogger:       log.GinLogger(),
}

func echo(agent *link.Agent) {
	log.Trace("echo创建连接: %d", agent.GetUniqueID())
	agent.OnClose = func() {
		log.Trace("echo断开连接: %d", agent.GetUniqueID())
	}
	for {
		data, ok := <-agent.ReadChan()
		if !ok {
			break
		}
		str := string(data.GetPayload())
		log.Trace("%d Recv: %s", agent.GetUniqueID(), str)

		agent.WriteMsg(util.MakeBuffDataBySlice([]byte("ECHO:"+str), 0))
	}
}

var serverOptPing = &link.ServerOption{
	Addr:            ":1234",
	MaxMsgLen:       0,
	MaxConnNum:      0,
	PendingWriteNum: 0,
	HTTPTimeout:     0,
	CertFile:        "",
	KeyFile:         "",
	RelativePath:    "/ping",
	GinLogger:       log.GinLogger(),
}

func ping(agent *link.Agent) {
	log.Trace("ping创建连接: %d", agent.GetUniqueID())
	agent.OnClose = func() {
		log.Trace("ping断开连接: %d", agent.GetUniqueID())
	}
	for {
		data, ok := <-agent.ReadChan()
		if !ok {
			break
		}
		str := string(data.GetPayload())
		log.Trace("%d Recv: %s", agent.GetUniqueID(), str)

		agent.WriteMsg(util.MakeBuffDataBySlice([]byte("Pong:"+str), 0))
	}
}

func ws() {
	{
		serv := &link.WSServer{}
		serv.OnAccept = func(agent *link.Agent) {
			go echo(agent)
		}
		serv.Start(serverOpt)
		log.Trace("WSServer启动 echo")
	}
	{
		serv := &link.WSServer{}
		serv.OnAccept = func(agent *link.Agent) {
			go ping(agent)
		}
		serv.Start(serverOptPing)
		log.Trace("WSServer启动 ping")
	}
}
