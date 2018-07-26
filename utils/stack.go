package utils

import (
	"unsafe"
	"sync/atomic"
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

func (this *KStack) Length() int32 {
	return atomic.LoadInt32(&this.length)
}

func (this *KStack) Push(_val interface{}) {
	if nil != _val {
		node := &KNode{val: _val}
		for {
			node.next = this.top
			if atomic.CompareAndSwapPointer(&this.top, node.next, unsafe.Pointer(node)) {
				atomic.AddInt32(&this.length, 1)
				break
			}
		}
	}
}

func (this *KStack) Pop() (interface{}, bool) {
	for {
		top := this.top
		if nil == top {
			return nil, false
		} else {
			node := (*KNode)(top)
			if atomic.CompareAndSwapPointer(&this.top, top, node.next) {
				atomic.AddInt32(&this.length, -1)
				return node.val, true
			}
		}
	}
	return nil, false
}
