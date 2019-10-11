package connect

import (
	"github.com/playnb/util"
	"net"
	"sync/atomic"
)

var _UniqueID = uint64(1000)

func GetConnUniqueID() uint64 {
	return atomic.AddUint64(&_UniqueID, 1)
}

type Conn interface {
	GetUniqueID() uint64
	CloseChan() chan bool

	ReadMsg() (util.BuffData, error)
	WriteMsg(data util.BuffData) error
	Close()

	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}
