package fix

import (
	"bytes"
)

type BufferPool struct {
	p *Pool[*bytes.Buffer]
}

func NewBufferPool(cacheSize, size, cap int) *BufferPool {
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
	return &BufferPool{
		p: NewPool(cacheSize, newFunc, clearFunc),
	}
}

func (p *BufferPool) Get() *bytes.Buffer {
	return p.p.Get()
}

func (p *BufferPool) Put(b *bytes.Buffer) {
	p.p.Put(b)
}
