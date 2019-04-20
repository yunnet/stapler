package utils

import (
	"sync/atomic"
	"time"
	"unsafe"
)

type KQueue struct {
	head   unsafe.Pointer
	tail   unsafe.Pointer
	length int64
}

func NewQueue() *KQueue {
	node := &KNode{}
	return &KQueue{
		head: unsafe.Pointer(node),
		tail: unsafe.Pointer(node),
	}
}

func (c *KQueue) Offer(_val interface{}) int64 {
	node := &KNode{val: _val}
	for {
		tail := c.tail
		next := (*KNode)(tail).next

		if tail == c.tail {
			if nil == next {
				if atomic.CompareAndSwapPointer(&(*KNode)(c.tail).next, next, unsafe.Pointer(node)) {
					atomic.CompareAndSwapPointer(&c.tail, tail, unsafe.Pointer(node))
					atomic.AddInt64(&c.length, 1)
					break
				}
			} else {
				atomic.CompareAndSwapPointer(&c.tail, tail, next)
			}
		}
	}
	return c.Length()
}

func (c *KQueue) Poll(timeout time.Duration) (val interface{}, ok bool) {
	val, ok = c.doPoll()
	for !ok && timeout > 0 {
		// println("doPoll timeout ", timeout)
		time.Sleep(time.Microsecond)
		timeout -= time.Microsecond
		val, ok = c.doPoll()
	}
	return
}

func (c *KQueue) doPoll() (interface{}, bool) {
	for {
		head := c.head
		tail := c.tail
		next := (*KNode)(head).next

		if head == c.head {
			if head == tail {
				if nil == next {
					return nil, false
				}
				atomic.CompareAndSwapPointer(&c.tail, tail, next)
			} else {
				val := (*KNode)(next).val
				if atomic.CompareAndSwapPointer(&c.head, head, next) {
					atomic.AddInt64(&c.length, -1)
					return val, true
				}
			}
		}
	}
}

func (c *KQueue) IsEmpty() bool {
	return 0 == c.Length()
}

func (c *KQueue) Length() int64 {
	return atomic.LoadInt64(&c.length)
}
