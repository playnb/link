package pack

import (
	"encoding/binary"
	"errors"
	"github.com/playnb/util"
	"io"
)

// 处理黏包
// --------------------
// | lenMsgLen | data |
// --------------------

type PackParser interface {
	Init(littleEndian bool, lenMsgLen int, minMsgLen int, maxMsgLen int)
	Read(conn io.Reader) (util.BuffData, error)
	WriteAll(args ...[]byte) (util.BuffData, error)
	Write(args util.BuffData) (util.BuffData, error)
}

type packParser struct {
	lenMsgLen    int //表示长度的字节位数
	minMsgLen    int
	maxMsgLen    int
	littleEndian bool //大小端字节序
}

func NewMsgParser() PackParser {
	p := &packParser{}
	p.Init(false, 2, 1, 4096)
	return p
}

func (mp *packParser) Init(littleEndian bool, lenMsgLen int, minMsgLen int, maxMsgLen int) {
	mp.littleEndian = littleEndian
	mp.lenMsgLen = lenMsgLen
	mp.minMsgLen = minMsgLen
	mp.maxMsgLen = maxMsgLen
}

func (mp *packParser) Read(conn io.Reader) (util.BuffData, error) {
	var b [4]byte
	bufMsgLen := b[:mp.lenMsgLen]

	// read len
	if _, err := io.ReadFull(conn, bufMsgLen); err != nil {
		return nil, err
	}

	var msgLen int
	switch mp.lenMsgLen {
	case 1:
		msgLen = int(bufMsgLen[0])
	case 2:
		if mp.littleEndian {
			msgLen = int(binary.LittleEndian.Uint16(bufMsgLen))
		} else {
			msgLen = int(binary.BigEndian.Uint16(bufMsgLen))
		}
	case 4:
		if mp.littleEndian {
			msgLen = int(binary.LittleEndian.Uint32(bufMsgLen))
		} else {
			msgLen = int(binary.BigEndian.Uint32(bufMsgLen))
		}
	}

	if msgLen > mp.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < mp.minMsgLen {
		return nil, errors.New("message too short")
	}

	msgData := util.DefaultPool().Get(msgLen)
	if _, err := io.ReadFull(conn, msgData.GetPayload()); err != nil {
		msgData.Release()
		return nil, err
	}
	return msgData, nil
}

func (mp *packParser) WriteAll(args ...[]byte) (util.BuffData, error) {
	var msgLen int
	for i := 0; i < len(args); i++ {
		msgLen += len(args[i])
	}

	if msgLen > mp.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < mp.minMsgLen {
		return nil, errors.New("message too short")
	}
	//msgLen += mp.lenMsgLen

	msgData := util.DefaultPool().Get(msgLen + mp.lenMsgLen)

	switch mp.lenMsgLen {
	case 1:
		msgData.GetPayload()[0] = byte(msgLen)
	case 2:
		if mp.littleEndian {
			binary.LittleEndian.PutUint16(msgData.GetPayload(), uint16(msgLen))
		} else {
			binary.BigEndian.PutUint16(msgData.GetPayload(), uint16(msgLen))
		}
	case 4:
		if mp.littleEndian {
			binary.LittleEndian.PutUint32(msgData.GetPayload(), uint32(msgLen))
		} else {
			binary.BigEndian.PutUint32(msgData.GetPayload(), uint32(msgLen))
		}
	}

	l := mp.lenMsgLen
	for i := 0; i < len(args); i++ {
		copy(msgData.GetPayload()[l:], args[i])
		l += len(args[i])
	}

	return msgData, nil
}
func (mp *packParser) Write(args util.BuffData) (util.BuffData, error) {
	//利用util.BuffData减少复制和GC
	return mp.WriteAll(args.GetPayload())
}
