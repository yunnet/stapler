package utils

import (
	"sync/atomic"
	"unsafe"
)

type KNode struct {
	val  interface{}
	next unsafe.Pointer
}

type KStack struct {
	top    unsafe.Pointer
	length int32
}

func NewStack() (stack *KStack) {
	return &KStack{}
}

func (c *KStack) Length() int32 {
	return atomic.LoadInt32(&c.length)
}

func (c *KStack) Push(_val interface{}) {
	if nil != _val {
		node := &KNode{val: _val}
		for {
			node.next = c.top
			if atomic.CompareAndSwapPointer(&c.top, node.next, unsafe.Pointer(node)) {
				atomic.AddInt32(&c.length, 1)
				break
			}
		}
	}
}

func (c *KStack) Pop() (interface{}, bool) {
	for {
		top := c.top
		if nil == top {
			return nil, false
		} else {
			node := (*KNode)(top)
			if atomic.CompareAndSwapPointer(&c.top, top, node.next) {
				atomic.AddInt32(&c.length, -1)
				return node.val, true
			}
		}
	}
	return nil, false
}
