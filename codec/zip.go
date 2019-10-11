package codec

import (
	"bytes"
	"compress/zlib"
	"github.com/playnb/util"
	"io/ioutil"
)

/*
zlib压缩
TODO: 后续可以使用单一的压缩缓冲
*/

type Zip struct {
}

func (e *Zip) Decode(data util.BuffData) util.BuffData {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(data.GetPayload())
	w.Close()
	data.Release()

	return util.MakeBuffDataBySlice(in.Bytes(), 0)
}

func (e *Zip) Encode(data util.BuffData) util.BuffData {
	r, _ := zlib.NewReader(bytes.NewReader(data.GetPayload()))
	data.Release()
	b, _ := ioutil.ReadAll(r)
	return util.MakeBuffDataBySlice(b, 0)
}

func (e *Zip) Clone() Codec {
	return e
}

func NewZipCodec() Codec {
	z := &Zip{}
	return z
}
