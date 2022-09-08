package cache

// 缓存接口
type Cache[K comparable, V any] interface {
	// 添加或更新元素
	Put(key K, value V)
	// 获取元素
	Get(key K) (V, bool)
}

// 淘汰时触发
type OnEvict[K comparable, V any] func(entry *Entry[K, V])

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}
