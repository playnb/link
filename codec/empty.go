package codec

import "github.com/playnb/util"

/*
空解码器
*/
type Empty struct {
}

func (e *Empty) Decode(data util.BuffData) util.BuffData {
	return data
}

func (e *Empty) Encode(data util.BuffData) util.BuffData {
	return data
}

func (e *Empty) Clone() Codec {
	return e
}
