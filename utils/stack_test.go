package utils

import (
	"testing"
	"time"
)

func TestKStack(t *testing.T) {
	println("test kstatk")
	stack := NewStack()

	type node struct {
		val int
	}

	println("setup")

	for i := 0; i < 100; i++ {
		stack.Push(&node{val: i})
	}

	println("go 1")
	go func() {
		for {
			if val, ok := stack.Pop(); ok {
				nd := val.(*node)
				println("1.pop=", nd.val)
			}
			time.Sleep(time.Microsecond)
		}
	}()

	println("go 2")
	go func() {
		for {
			if val, ok := stack.Pop(); ok {
				nd := val.(*node)
				println("2.pop=", nd.val)
			}
			time.Sleep(time.Microsecond)
		}
	}()

	tm := time.NewTimer(time.Second * 5)
	println("start timer")
	for {
		select {
		case <-tm.C:
			println("time")
			tm.Reset(time.Second * 5)
		}
	}
}
