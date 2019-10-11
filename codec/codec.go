package codec

import "github.com/playnb/util"

type Codec interface {
	Encode(util.BuffData) util.BuffData
	Decode(util.BuffData) util.BuffData
	Clone() Codec //每次实例是否复制
}
