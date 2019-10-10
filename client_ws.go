package link

import (
	"github.com/gorilla/websocket"
	"github.com/playnb/util/log"
	"sync"
	"time"
)

type ClientOption struct {
	Addr             string
	MaxMsgLen        int
	MaxConnNum       int
	PendingWriteNum  int
	HandshakeTimeout time.Duration
}

func (opt *ClientOption) Check() {
	if opt.MaxConnNum <= 0 {
		opt.MaxConnNum = 100
	}
	if opt.PendingWriteNum <= 0 {
		opt.PendingWriteNum = 200
	}
	if opt.MaxMsgLen <= 0 {
		opt.MaxMsgLen = 4096
	}
	if opt.HandshakeTimeout <= 0 {
		opt.HandshakeTimeout = 10 * time.Second
	}
}

type WSClient struct {
	sync.Mutex

	option *ClientOption
	dialer websocket.Dialer
	conn   Conn
}

func (client *WSClient) dial() error {
	client.dialer = websocket.Dialer{
		HandshakeTimeout: client.option.HandshakeTimeout,
	}
	conn, _, err := client.dialer.Dial(client.option.Addr, nil)
	if err != nil {
		return err
	}
	client.conn = newWSConn(conn, client.option.PendingWriteNum, client.option.MaxMsgLen)
	return nil
}

func (client *WSClient) Start(option *ClientOption) *Agent {
	client.option = option
	client.option.Check()
	err := client.dial()
	if err != nil {
		return nil
	}
	agent := &Agent{}
	agent.init(client.conn, client.option.MaxMsgLen)

	go func() {
		for {
			data, err := client.conn.ReadMsg()
			if err != nil {
				break
			}
			if len(agent.msgChan) == cap(agent.msgChan) {
				log.Error("Agent msgChan full")
				continue
			}
			agent.msgChan <- data
		}
		close(agent.msgChan)
		agent.OnClose()
	}()

	return agent
}
