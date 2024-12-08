package main

// Tree定义一个二叉树
type Tree[T any] struct {
	cmp  func(T, T) int
	root *node[T]
}

// 二叉树的一个节点
type node[T any] struct {
	left, right *node[T]
	val         T
}

// find查找值
func (bt *Tree[T]) find(val T) **node[T] {
	pl := &bt.root
	for *pl != nil {
		switch cmp := bt.cmp(val, (*pl).val); {
		case cmp < 0:
			pl = &(*pl).left
		case cmp > 0:
			pl = &(*pl).right
		default:
			return pl
		}
	}
	return pl
}
