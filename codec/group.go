package codec

import "github.com/playnb/util"

/*
编码器组合
*/
type Group struct {
	group     []Codec
	replicate bool
}

func (g *Group) Encode(data util.BuffData) util.BuffData {
	if g == nil {
		return data
	}
	for i := 0; i < len(g.group); i++ {
		data = g.group[i].Encode(data)
	}
	return data
}

func (g *Group) Decode(data util.BuffData) util.BuffData {
	if g == nil {
		return data
	}
	for i := len(g.group); i > 0; i-- {
		data = g.group[i-1].Encode(data)
	}
	return data
}

func (g *Group) Clone() Codec {
	if g.replicate {
		return NewGroup(g.replicate, g.group...)
	} else {
		return g
	}
}

func (g *Group) AddCodec(c Codec) {
	g.group = append(g.group, c)
}

func NewGroup(replicate bool, c ...Codec) Codec {
	g := &Group{}
	g.replicate = replicate
	for _, v := range c {
		g.AddCodec(v.Clone())
	}
	return g
}
