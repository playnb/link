package link

import (
	"errors"
	"github.com/playnb/util"
	"github.com/playnb/util/log"
	"net"
	"sync"
	"sync/atomic"
)

type ImpConn struct {
	sync.Mutex
	conn      net.Conn
	writeChan chan util.BuffData //发送缓冲
	closeFlag bool
	uniqueID  uint64

	doDestroy func()
	closeChan chan bool
}

//PendingWriteNum: 发送缓冲区大小
//msgParser: 黏包的解析器
func newImpConn(conn net.Conn, pendingWriteNum int) *ImpConn {
	iConn := &ImpConn{}
	iConn.conn = conn
	iConn.closeChan = make(chan bool)
	iConn.writeChan = make(chan util.BuffData, pendingWriteNum)
	iConn.closeFlag = false
	iConn.uniqueID = atomic.AddUint64(&_UniqueID, 1)

	//发送线程
	go func() {
		for b := range iConn.writeChan {
			if b == nil {
				//发送nil标识断开连接
				break
			}
			buf := b.GetPayload()
			//log.Debug("%d 发送数据 %s", iConn.GetUniqueID(), buf)
			_, err := conn.Write(buf)
			b.Release()
			if err != nil {
				log.Error(err.Error())
				break
			}
		}

		log.Trace("%d 结束发送线程", iConn.GetUniqueID())
		conn.Close()
		iConn.Lock()
		iConn.closeFlag = true
		iConn.Unlock()
	}()

	return iConn
}

func (iConn *ImpConn) GetUniqueID() uint64 {
	return iConn.uniqueID
}

func (iConn *ImpConn) CloseChan() chan bool {
	return iConn.closeChan
}

func (iConn *ImpConn) doClose() {
	//iConn.conn.SetLinger(0)
	//iConn.conn.Close()
	iConn.doDestroy()

	if !iConn.closeFlag {
		close(iConn.writeChan)
		close(iConn.closeChan)
		iConn.closeFlag = true
	}
}

func (iConn *ImpConn) Close() {
	iConn.Lock()
	defer func() {
		iConn.closeFlag = true
		iConn.Unlock()
	}()

	if iConn.closeFlag {
		return
	}
	iConn.doClose()
	iConn.doWrite(nil)
}

func (iConn *ImpConn) doWrite(b util.BuffData) error {
	if iConn.closeFlag || b == nil {
		return errors.New("写入已关闭的连接")
	}

	if len(iConn.writeChan) == cap(iConn.writeChan) {
		log.Error("发送缓冲到达上限,丢弃消息 %v", b)
		if b == nil {
			iConn.doClose()
		}
		return errors.New("发送缓冲到达上限,丢弃消息")
	}

	iConn.writeChan <- b
	return nil
}

func (iConn *ImpConn) LocalAddr() net.Addr {
	return iConn.conn.LocalAddr()
}

func (iConn *ImpConn) RemoteAddr() net.Addr {
	return iConn.conn.RemoteAddr()
}
