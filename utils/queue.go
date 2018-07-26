package utils

import (
	"unsafe"
	"sync/atomic"
	"time"
)

type KQueue struct {
	head   unsafe.Pointer
	tail   unsafe.Pointer
	length int64
}

func NewQueue() (*KQueue) {
	node := &KNode{}
	return &KQueue{
		head: unsafe.Pointer(node),
		tail: unsafe.Pointer(node),
	}
}

func (this *KQueue) Offer(_val interface{}) (int64) {
	node := &KNode{val: _val}
	for {
		tail := this.tail
		next := (*KNode)(tail).next

		if tail == this.tail {
			if nil == next {
				if atomic.CompareAndSwapPointer(&(*KNode)(this.tail).next, next, unsafe.Pointer(node)) {
					atomic.CompareAndSwapPointer(&this.tail, tail, unsafe.Pointer(node))
					atomic.AddInt64(&this.length, 1)
					break
				}
			} else {
				atomic.CompareAndSwapPointer(&this.tail, tail, next)
			}
		}
	}
	return this.Length()
}

func (this *KQueue) Poll(timeout time.Duration) (val interface{}, ok bool) {
	val, ok = this.doPoll()
	for !ok && timeout > 0 {
		// println("doPoll timeout ", timeout)
		time.Sleep(time.Microsecond)
		timeout -= time.Microsecond
		val, ok = this.doPoll()
	}
	return
}

func (this *KQueue) doPoll() (interface{}, bool) {
	for {
		head := this.head
		tail := this.tail
		next := (*KNode)(head).next

		if head == this.head {
			if head == tail {
				if nil == next {
					return nil, false
				}
				atomic.CompareAndSwapPointer(&this.tail, tail, next)
			} else {
				val := (*KNode)(next).val
				if atomic.CompareAndSwapPointer(&this.head, head, next) {
					atomic.AddInt64(&this.length, -1)
					return val, true
				}
			}
		}
	}
}

func (this *KQueue) IsEmpty() (bool) {
	return 0 == this.Length()
}

func (this *KQueue) Length() (int64) {
	return atomic.LoadInt64(&this.length)
}