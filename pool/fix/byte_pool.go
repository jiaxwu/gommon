package fix

type BytePool struct {
	p *Pool[[]byte]
}

// cacheSize: 字节池缓存长度
// size: 字节数组长度
// cap: 字节数组容量
func NewBytePool(cacheSize, size, cap int) *BytePool {
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
		p: NewPool(cacheSize, newFunc, clearFunc),
	}
}

func (p *BytePool) Get() []byte {
	return p.p.Get()
}

func (p *BytePool) Put(b []byte) {
	p.p.Put(b)
}
