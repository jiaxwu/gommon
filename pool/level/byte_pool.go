package level

import (
	"sort"
)

type BytePool struct {
	p        *Pool[[]byte]
	levels   []int
	perLevel int
	maxLevel int
}

// 比如level=5，则等级为perLevel、perLevel*2、perLevel*4、perLevel*8、perLevel*16和更大
func NewBytePool(perLevel, maxLevel int) *BytePool {
	levels := make([]int, maxLevel)
	levels[0] = perLevel
	for i := 1; i < maxLevel; i++ {
		levels[i] = levels[i-1] * 2
	}
	newFunc := func(l int) []byte {
		if l == maxLevel {
			return nil
		}
		return make([]byte, 0, levels[l])
	}
	clearFunc := func(b []byte) []byte {
		return b[:0]
	}
	p := &BytePool{
		perLevel: perLevel,
		maxLevel: maxLevel,
		levels:   levels,
	}
	levelFunc := func(b []byte) int {
		return p.level(cap(b))
	}
	p.p = New(newFunc, clearFunc, levelFunc, maxLevel)
	return p
}

func (p *BytePool) level(capacity int) int {
	return sort.SearchInts(p.levels, capacity)
}

// 获取字节数组
func (p *BytePool) Get(capacity int) []byte {
	// 计算等级
	l := p.level(capacity)
	// 尝试获取
	b := p.p.Get(l)
	// 获取不到则新创建一个
	if cap(b) < capacity {
		return make([]byte, capacity)
	}
	return b
}

// 归还字节数组
func (p *BytePool) Put(b []byte) {
	p.p.Put(b)
}
