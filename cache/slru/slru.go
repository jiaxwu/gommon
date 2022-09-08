package slru

import (
	"github.com/jiaxwu/gommon/container/list"
)

// 段类型
type SegmentType uint8

const (
	// 试用段
	SegmentTypeProbation = 1
	// 保护段
	SegmentTypeProtected = 2
)

// 保护段比例
const ProtectedPercentage = 0.8

// 淘汰时触发
type OnEvict[K comparable, V any] func(entry *Entry[K, V])

type Entry[K comparable, V any] struct {
	Key         K
	Value       V
	SegmentType SegmentType
}

// 分段最近最少使用
// 一开始Put()会写入第一个段
// 然后被Get()才写入第二个段
// 避免某些很少读取的值把一直读取的值给淘汰了
// 优点：稳定淘汰，避免大量失效
// 非线程安全，请根据业务加锁
type Cache[K comparable, V any] struct {
	entries      map[K]*list.Element[*Entry[K, V]]
	probation    *list.List[*Entry[K, V]] // 试用段，容易淘汰
	protected    *list.List[*Entry[K, V]] // 保护段，不容易淘汰
	probationCap int
	protectedCap int
	onEvict      OnEvict[K, V]
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	protectedCap := int(float64(capacity) * ProtectedPercentage)
	probationCap := capacity - protectedCap
	if probationCap < 1 {
		panic("too small probation capacity")
	}
	if protectedCap < 1 {
		panic("too small protected capacity")
	}
	return &Cache[K, V]{
		probation:    list.New[*Entry[K, V]](),
		protected:    list.New[*Entry[K, V]](),
		probationCap: probationCap,
		protectedCap: probationCap,
	}
}

// 设置 OnEvict
func (c *Cache[K, V]) SetOnEvict(onEvict OnEvict[K, V]) {
	c.onEvict = onEvict
}

// 添加或更新元素
func (c *Cache[K, V]) Put(key K, value V) {
	// 如果 key 已经存在，直接Get()更新元素状态
	if elem, ok := c.entries[key]; ok {
		elem.Value.Value = value
		c.access(elem)
		return
	}

	// 如果已经到达最大尺寸，先剔除probation段的一个元素
	if c.Full() {
		c.evict(c.probation)
	}

	// 添加元素到probation段
	c.put(c.probation, &Entry[K, V]{
		Key:         key,
		Value:       value,
		SegmentType: SegmentTypeProbation,
	})
}

// 获取元素
func (c *Cache[K, V]) Get(key K) (V, bool) {
	// 如果存在则更新元素状态
	if elem, ok := c.entries[key]; ok {
		c.updateElement(elem)
		return elem.Value.Value, true
	}

	// 不存在返回空值和false
	var value V
	return value, false
}

// 获取元素，不更新状态
func (c *Cache[K, V]) Peek(key K) (V, bool) {
	// 如果存在
	if elem, ok := c.entries[key]; ok {
		return elem.Value.Value, true
	}

	// 不存在返回空值和false
	var value V
	return value, false
}

// 是否包含元素，不更新状态
func (c *Cache[K, V]) Contains(key K) bool {
	_, ok := c.entries[key]
	return ok
}

// 获取缓存的Keys
func (c *Cache[K, V]) Keys() []K {
	keys := make([]K, c.Len())
	for elem, i := c.evictList.Back(), 0; elem != nil; elem, i = elem.Prev(), i+1 {
		keys[i] = elem.Value.Key
	}
	return keys
}

// 获取缓存的Values
func (c *Cache[K, V]) Values() []V {
	values := make([]V, c.Len())
	for elem, i := c.evictList.Back(), 0; elem != nil; elem, i = elem.Prev(), i+1 {
		values[i] = elem.Value.Value
	}
	return values
}

// 获取缓存的Entries
func (c *Cache[K, V]) Entries() []*Entry[K, V] {
	entries := make([]*Entry[K, V], c.Len())
	for elem, i := c.evictList.Back(), 0; elem != nil; elem, i = elem.Prev(), i+1 {
		entries[i] = elem.Value
	}
	return entries
}

// 移除元素
func (c *Cache[K, V]) Remove(key K) bool {
	if elem, ok := c.entries[key]; ok {
		c.removeElement(elem)
		return true
	}
	return false
}

// 清空缓存
func (c *Cache[K, V]) Clear(needOnEvict bool) {
	// 触发回调
	if needOnEvict && c.onEvict != nil {
		for elem, i := c.evictList.Back(), 0; elem != nil; elem, i = elem.Prev(), i+1 {
			c.onEvict(elem.Value)
		}
	}

	// 清空
	c.entries = make(map[K]*list.Element[*Entry[K, V]])
	c.evictList.Init()
}

// 改变容量
func (c *Cache[K, V]) Resize(capacity int, needOnEvict bool) {
	diff := c.Len() - capacity
	if diff < 0 {
		diff = 0
	}
	for i := 0; i < diff; i++ {
		c.Evict()
	}
	c.capacity = capacity
}

// 元素个数
func (c *Cache[K, V]) Len() int {
	return len(c.entries)
}

// 容量
func (c *Cache[K, V]) Cap() int {
	return c.capacity
}

// 缓存满了
func (c *Cache[K, V]) Full() bool {
	return c.probation.Full() && c.protected.Full()
}



// 添加元素
func (c *Cache[K, V]) put(evictList *list.List[*Entry[K, V]], entry *Entry[K, V]) {
	elem := evictList.PushFront(entry)
	c.entries[entry.Key] = elem
}

// 移除节点
func (c *Cache[K, V]) removeElement(evictList *list.List[*Entry[K, V]], elem *list.Element[*Entry[K, V]]) {
	evictList.Remove(elem)
	entry := elem.Value
	delete(c.entries, entry.Key)
}

// 更新元素
func (c *Cache[K, V]) access(elem *list.Element[*Entry[K, V]]) {
	entry := elem.Value
	// 在保护段则直接移动到头部即可
	if entry.SegmentType == SegmentTypeProtected {
		c.protected.MoveToFront(elem)
		return
	}

	// 否则一定在试用段，需要提升到保护段
	// 如果保护段满了，则把保护段的一个元素移动到试用段
	if c.protected.Len() == c.protectedCap {
		// 从保护段淘汰一个元素
		entry := c.protected.RemoveBack()
		// 更新段类型
		entry.SegmentType = SegmentTypeProbation
		// 添加到试用段
		c.entries[entry.Key] = c.probation.PushFront(entry)
	}
	
	// 从试用段移除
	c.removeElement(c.probation, elem)

	// 提升到
	if 
}

// 淘汰元素
func (c *Cache[K, V]) evict(evictList *list.List[*Entry[K, V]]) {
	elem := evictList.Back()
	if elem != nil {
		c.removeElement(evictList, elem)
		// 回调
		if c.onEvict != nil {
			c.onEvict(elem.Value)
		}
	}
}
