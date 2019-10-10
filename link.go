package link

import "github.com/playnb/util"

type Agent struct {
	OnClose func()

	msgChan chan util.BuffData
	conn    Conn
}

func (agent *Agent) init(conn Conn, pendingNum int) {
	agent.msgChan = make(chan util.BuffData, pendingNum)
	agent.conn = conn
}

func (agent *Agent) GetUniqueID() uint64 {
	return agent.conn.GetUniqueID()
}

func (agent *Agent) ReadChan() chan util.BuffData {
	return agent.msgChan
}

func (agent *Agent) WriteMsg(data util.BuffData) error {
	return agent.conn.WriteMsg(data)
}

func (agent *Agent) Close() {
	agent.conn.Close()
}

func (agent *Agent) CloseChan() chan bool {
	return agent.conn.CloseChan()
}
