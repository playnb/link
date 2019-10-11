package link

import (
	"github.com/playnb/link/codec"
	"github.com/playnb/link/connect"
	"github.com/playnb/util"
)

type Agent struct {
	UserData
	OnClose func()

	msgChan chan util.BuffData
	conn    connect.Conn
	cc      codec.Codec
}

func (agent *Agent) init(conn connect.Conn, pendingNum int) {
	agent.msgChan = make(chan util.BuffData, pendingNum)
	agent.conn = conn
}

func (agent *Agent) putChan(data util.BuffData) {
	if agent.cc != nil {
		agent.msgChan <- agent.cc.Decode(data)
	} else {
		agent.msgChan <- data
	}
}
func (agent *Agent) closeChan() {
	close(agent.msgChan)
}

func (agent *Agent) GetUniqueID() uint64 {
	return agent.conn.GetUniqueID()
}

func (agent *Agent) ReadChan() chan util.BuffData {
	return agent.msgChan
}

func (agent *Agent) WriteMsg(data util.BuffData) error {
	if agent.cc != nil {
		return agent.conn.WriteMsg(agent.cc.Encode(data))
	} else {
		return agent.conn.WriteMsg(data)
	}
}

func (agent *Agent) Close() {
	agent.conn.Close()
}

func (agent *Agent) CloseChan() chan bool {
	return agent.conn.CloseChan()
}
