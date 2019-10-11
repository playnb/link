package connect

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/playnb/util"
	"net"
	"time"
)

type wsconn struct {
	*websocket.Conn
}

func (w wsconn) Read(b []byte) (n int, err error) {
	//_, buf, err := w.Conn.ReadMessage()
	//return len(buf), err
	return 0, errors.New("not implement")
}

func (w wsconn) Write(b []byte) (n int, err error) {
	return len(b), w.Conn.WriteMessage(websocket.BinaryMessage, b)
}

func (w wsconn) SetDeadline(t time.Time) error {
	e := w.Conn.SetReadDeadline(t)
	if e != nil {
		return e
	}
	e = w.Conn.SetWriteDeadline(t)
	return e
}

type WSConn struct {
	*ImpConn
	maxMsgLen int
}

func NewWSConn(conn *websocket.Conn, pendingWriteNum int, maxMsgLen int) Conn {
	wsConn := &WSConn{}
	wsConn.maxMsgLen = maxMsgLen
	wsConn.ImpConn = newImpConn(&wsconn{Conn: conn}, pendingWriteNum)
	wsConn.doDestroy = func() {
		wsConn.conn.(*wsconn).UnderlyingConn().(*net.TCPConn).SetLinger(0)
		wsConn.conn.Close()
	}
	return wsConn
}

func (wsConn *WSConn) ReadMsg() (util.BuffData, error) {
	_, b, err := wsConn.conn.(*wsconn).ReadMessage()
	return util.MakeBuffDataBySlice(b, 0), err
}

func (wsConn *WSConn) WriteMsg(data util.BuffData) error {
	msgLen := len(data.GetPayload())
	if msgLen > wsConn.maxMsgLen {
		return errors.New("message too long")
	} else if msgLen < 1 {
		return errors.New("message too short")
	}

	wsConn.Lock()
	defer wsConn.Unlock()
	wsConn.doWrite(data)
	return nil
}
