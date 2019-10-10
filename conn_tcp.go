package link

import (
	"github.com/playnb/util"
	"net"
)

type TCPConn struct {
	*ImpConn
	msgParser MsgParser
}

//PendingWriteNum: 发送缓冲区大小
//msgParser: 黏包的解析器
func newTCPConn(conn *net.TCPConn, pendingWriteNum int, msgParser MsgParser) Conn {
	tcpConn := &TCPConn{}
	tcpConn.msgParser = msgParser
	tcpConn.ImpConn = newImpConn(conn, pendingWriteNum)
	tcpConn.doDestroy = func() {
		tcpConn.conn.(*net.TCPConn).SetLinger(0)
		tcpConn.conn.Close()
	}
	return tcpConn
}

func (tcpConn *TCPConn) ReadMsg() (util.BuffData, error) {
	return tcpConn.msgParser.Read(tcpConn.conn)
}

func (tcpConn *TCPConn) WriteMsg(data util.BuffData) error {
	data, err := tcpConn.msgParser.Write(data)
	if err != nil {
		return err
	}

	tcpConn.Lock()
	defer tcpConn.Unlock()
	tcpConn.doWrite(data)
	return nil
}
