package interval

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/constraints"
)

// 区间，双闭
type Interval[Endpoint constraints.Ordered, T any] struct {
	start Endpoint
	end   Endpoint
	value T
}

func (i *Interval[Endpoint, T]) SetValue(value T) {
	i.value = value
}

func (i *Interval[Endpoint, T]) String() string {
	return fmt.Sprintf("|%v|%v|%v|", i.start, i.end, i.value)
}

func NewInterval[Endpoint constraints.Ordered, T any](start, end Endpoint) *Interval[Endpoint, T] {
	return &Interval[Endpoint, T]{
		start: start,
		end:   end,
	}
}

func NewIntervalWithValue[Endpoint constraints.Ordered, T any](start, end Endpoint, value T) *Interval[Endpoint, T] {
	return &Interval[Endpoint, T]{
		start: start,
		end:   end,
		value: value,
	}
}

// 节点
type Node[Endpoint constraints.Ordered, T any] struct {
	interval *Interval[Endpoint, T]
	maxEnd   Endpoint
	left     *Node[Endpoint, T]
	right    *Node[Endpoint, T]
}

func NewNode[Endpoint constraints.Ordered, T any](interval *Interval[Endpoint, T]) *Node[Endpoint, T] {
	return &Node[Endpoint, T]{
		interval: interval,
		maxEnd:   interval.end,
	}
}

// 区间树
// 支持区间作为节点
type IntervalTree[Endpoint constraints.Ordered, T any] struct {
	root *Node[Endpoint, T]
}

func NewIntervalTree[Endpoint constraints.Ordered, T any](intervals ...*Interval[Endpoint, T]) *IntervalTree[Endpoint, T] {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].start < intervals[j].start
	})
	return &IntervalTree[Endpoint, T]{root: buildTree(intervals)}
}

func (t *IntervalTree[Endpoint, T]) Query(value Endpoint) []T {
	var result []T
	t.query(t.root, value, &result)
	return result
}

func (t *IntervalTree[Endpoint, T]) Print(spaceLen int) {
	t.print(t.root, "", true, true, spaceLen)
}

func (t *IntervalTree[Endpoint, T]) query(node *Node[Endpoint, T], value Endpoint, result *[]T) {
	if node == nil {
		return
	}

	if value >= node.interval.start && value <= node.interval.end {
		*result = append(*result, node.interval.value)
	}

	if node.left != nil && node.left.maxEnd >= value {
		t.query(node.left, value, result)
	}

	if node.right != nil && node.right.interval.start <= value {
		t.query(node.right, value, result)
	}
}

func (t *IntervalTree[Endpoint, T]) print(root *Node[Endpoint, T], prefix string, isTail, isFirst bool, spaceLen int) {
	if root == nil {
		return
	}

	fmt.Printf("%s", prefix)
	if isTail {
		if !isFirst {
			fmt.Printf("└%s ", strings.Repeat("─", spaceLen))
			prefix += "  " + strings.Repeat(" ", spaceLen)
		}
	} else {
		if root.right != nil {
			fmt.Printf("├%s ", strings.Repeat("─", spaceLen))
		} else {
			fmt.Printf("└%s ", strings.Repeat("─", spaceLen))
		}
		prefix += fmt.Sprintf("│%s ", strings.Repeat(" ", spaceLen))
	}
	fmt.Printf("|%v,%v|%v|%v|\n", root.interval.start, root.interval.end, root.maxEnd, root.interval.value)

	t.print(root.left, prefix, false, false, spaceLen)
	t.print(root.right, prefix, true, false, spaceLen)
}

func buildTree[Endpoint constraints.Ordered, T any](intervals []*Interval[Endpoint, T]) *Node[Endpoint, T] {
	if len(intervals) == 0 {
		return nil
	}

	mid := (len(intervals) - 1) / 2
	root := NewNode(intervals[mid])

	root.left = buildTree(intervals[:mid])
	root.right = buildTree(intervals[mid+1:])

	if root.left != nil && root.left.maxEnd > root.maxEnd {
		root.maxEnd = root.left.maxEnd
	}
	if root.right != nil && root.right.maxEnd > root.maxEnd {
		root.maxEnd = root.right.maxEnd
	}

	return root
}
