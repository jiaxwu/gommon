package level

import (
	"bytes"
	"sort"
)

type BufferPool struct {
	p        *Pool[*bytes.Buffer]
	levels   []int
	perLevel int
	maxLevel int
}

// 比如level=5，则等级为perLevel、perLevel*2、perLevel*4、perLevel*8、perLevel*16和更大
func NewBufferPool(perLevel, maxLevel int) *BufferPool {
	levels := make([]int, maxLevel)
	levels[0] = perLevel
	for i := 1; i < maxLevel; i++ {
		levels[i] = levels[i-1] * 2
	}
	newFunc := func(l int) *bytes.Buffer {
		if l == maxLevel {
			return nil
		}
		return bytes.NewBuffer(make([]byte, 0, levels[l]))
	}
	clearFunc := func(b *bytes.Buffer) *bytes.Buffer {
		b.Reset()
		return b
	}
	p := &BufferPool{
		perLevel: perLevel,
		maxLevel: maxLevel,
		levels:   levels,
	}
	levelFunc := func(b *bytes.Buffer) int {
		return p.level(b.Cap())
	}
	p.p = New(newFunc, clearFunc, levelFunc, maxLevel)
	return p
}

func (p *BufferPool) level(capacity int) int {
	return sort.SearchInts(p.levels, capacity)
}

// 获取字节Buffer
func (p *BufferPool) Get(capacity int) *bytes.Buffer {
	// 计算等级
	l := p.level(capacity)
	// 尝试获取
	b := p.p.Get(l)
	// 获取不到则新创建一个
	if b.Cap() < capacity {
		return bytes.NewBuffer(make([]byte, capacity))
	}
	return b
}

// 归还字节Buffer
func (p *BufferPool) Put(b *bytes.Buffer) {
	p.p.Put(b)
}
