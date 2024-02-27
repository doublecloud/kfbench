package value

import (
	"strconv"
)

type BytesGenerator struct {
	size int64
}

func NewBytesGenerator(size int64) *BytesGenerator {
	return &BytesGenerator{size: size}
}

func (g *BytesGenerator) Generate(num int64) []byte {
	return NewBytes(num, g.size)
}

type StaticBytes struct {
	v []byte
}

func NewStaticBytes(size int64) *StaticBytes {
	return &StaticBytes{v: NewBytes(0, size)}
}

func (g *StaticBytes) Generate(_ int64) []byte {
	return g.v
}

func NewBytes(num int64, size int64) []byte {
	v := make([]byte, size)
	Format(v, num)
	return v
}

func Format(v []byte, num int64) {
	var buf [20]byte // max int64 takes 19 bytes, then we add a space
	b := strconv.AppendInt(buf[:0], num, 10)
	b = append(b, ' ')

	n := copy(v, b)
	for n != len(v) {
		n += copy(v[n:], b)
	}
}
