package pool

import (
	"bytes"
)

type BufferFixPool struct {
	p *FixPool[*bytes.Buffer]
}

func NewBufferFixPool(cacheSize, size, cap int) *BufferFixPool {
	if size > cap {
		panic("size must be less then cap")
	}
	newFunc := func() *bytes.Buffer {
		var b []byte
		if cap > 0 {
			b = make([]byte, size, cap)
		}
		return bytes.NewBuffer(b)
	}
	clearFunc := func(b *bytes.Buffer) *bytes.Buffer {
		b.Reset()
		return b
	}
	return &BufferFixPool{
		p: NewFixPool(cacheSize, newFunc, clearFunc),
	}
}

func (p *BufferFixPool) Get() *bytes.Buffer {
	return p.p.Get()
}

func (p *BufferFixPool) Put(b *bytes.Buffer) {
	p.p.Put(b)
}
