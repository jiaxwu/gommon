/*
Copyright 2013 Google Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package consistenthash provides an implementation of a ring hash.

// 使用 https://github.com/golang/groupcache/blob/master/consistenthash/consistenthash_test.go 的实现，添加Reset方法

package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 哈希函数
type Hash func(data []byte) uint32

// 哈希环
// 注意，非线程安全，业务需要自行加锁
type HashRing struct {
	hash Hash
	// 每个真实节点的虚拟节点数量
	replicas int
	// 哈希环，按照节点哈希值排序
	ring []int
	// 节点哈希值到真实节点字符串，哈希映射的逆过程
	nodes map[int]string
}

func New(replicas int, fn Hash) *HashRing {
	r := &HashRing{
		replicas: replicas,
		hash:     fn,
		nodes:    make(map[int]string),
	}
	if r.hash == nil {
		r.hash = crc32.ChecksumIEEE
	}
	return r
}

// 哈希环上是否有节点
func (r *HashRing) Empty() bool {
	return len(r.ring) == 0
}

// 添加新节点到哈希环
// 注意，如果加入的节点已经存在，会导致哈希环上面重复，如果不确定是否存在请使用Reset
func (m *HashRing) Add(nodes ...string) {
	for _, node := range nodes {
		// 每个节点创建多个虚拟节点
		for i := 0; i < m.replicas; i++ {
			// 每个虚拟节点计算哈希值
			hash := int(m.hash([]byte(strconv.Itoa(i) + node)))
			// 加入哈希环
			m.ring = append(m.ring, hash)
			// 哈希值到真实节点字符串映射
			m.nodes[hash] = node
		}
	}
	// 哈希环排序
	sort.Ints(m.ring)
}

// 先清空哈希环再设置
func (r *HashRing) Reset(nodes ...string) {
	// 先清空
	r.ring = nil
	r.nodes = map[int]string{}
	// 再重置
	r.Add(nodes...)
}

// 获取Key对应的节点
func (r *HashRing) Get(key string) string {
	// 如果哈希环位空，则直接返回
	if r.Empty() {
		return ""
	}

	// 计算Key哈希值
	hash := int(r.hash([]byte(key)))

	// 二分查找第一个大于等于Key哈希值的节点
	idx := sort.Search(len(r.ring), func(i int) bool { return r.ring[i] >= hash })

	// 这里是特殊情况，也就是数组没有大于等于Key哈希值的节点
	// 但是逻辑上这是一个环，因此第一个节点就是目标节点
	if idx == len(r.ring) {
		idx = 0
	}

	// 返回哈希值对应的真实节点字符串
	return r.nodes[r.ring[idx]]
}
