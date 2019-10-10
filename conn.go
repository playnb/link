package link

import (
	"github.com/playnb/util"
	"net"
)

var _UniqueID = uint64(1000)

type Conn interface {
	GetUniqueID() uint64
	CloseChan() chan bool

	ReadMsg() (util.BuffData, error)
	WriteMsg(data util.BuffData) error
	Close()

	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}
