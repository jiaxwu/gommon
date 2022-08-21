package pool

type ByteFixPool struct {
	p *FixPool[[]byte]
}

// cacheSize: 字节池缓存长度
// size: 字节数组长度
// cap: 字节数组容量
func NewByteFixPool(cacheSize, size, cap int) *ByteFixPool {
	if size > cap {
		panic("size must be less then cap")
	}
	newFunc := func() []byte {
		return make([]byte, size, cap)
	}
	clearFunc := func(b []byte) []byte {
		return b[:0]
	}
	return &ByteFixPool{
		p: NewFixPool(cacheSize, newFunc, clearFunc),
	}
}

func (p *ByteFixPool) Get() []byte {
	return p.p.Get()
}

func (p *ByteFixPool) Put(b []byte) {
	p.p.Put(b)
}
