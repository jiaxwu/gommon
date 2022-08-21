package pool

type BytePool struct {
	p *Pool[[]byte]
}

func NewBytePool(size, cap int) *BytePool {
	if size > cap {
		panic("size must be less then cap")
	}
	newFunc := func() []byte {
		return make([]byte, size, cap)
	}
	clearFunc := func(b []byte) []byte {
		return b[:0]
	}
	return &BytePool{
		p: New(newFunc, clearFunc),
	}
}

// 获取字节数组
func (p *BytePool) Get() []byte {
	return p.p.Get()
}

// 归还字节数组
func (p *BytePool) Put(b []byte) {
	p.p.Put(b)
}
