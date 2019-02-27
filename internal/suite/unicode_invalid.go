package suite

import (
	"math"
	. "unicode"
)

var UnicodeInvalid = map[string]*RangeTable{
	"Non-Character":          NonCharacter,
	"UTF8 Invalid Bytes":     UTF8InvalidBytes,
}

var (
	NonCharacter         = _NonCharacter         // Code points permanently reserved for internal use: https://www.unicode.org/faq/private_use.html#noncharacters
	UTF8InvalidBytes     = _UTF8InvalidBytes     // Bytes that must never appear in a valid UTF8 sequence
)

var _NonCharacter = func() *RangeTable {
	rt := RangeTable{}
	rt.R16 = append(rt.R16, Range16{Lo: 0xfdd0, Hi: 0xfdef, Stride: 1})
	for i := uint32(0); i <= 0x100000; i += 0x10000 {
		cp0, cp1 := i+0xfffe, i+0xffff
		if cp1 < math.MaxUint16 {
			r16 := Range16{Lo: uint16(cp0), Hi: uint16(cp1), Stride: 1}
			rt.R16 = append(rt.R16, r16)
		} else {
			r32 := Range32{Lo: cp0, Hi: cp1, Stride: 1}
			rt.R32 = append(rt.R32, r32)
		}
	}
	return &rt
}()

var _UTF8InvalidBytes = func() *RangeTable {
	return &RangeTable{
		R16: []Range16{
			{Lo: 0xc0, Hi: 0xc1, Stride: 1},
			{Lo: 0xf5, Hi: 0xfd, Stride: 1},
		},
		LatinOffset: 2,
	}
}()

