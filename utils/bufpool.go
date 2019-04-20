package utils

type KBufPool struct {
	KStack
	bufSize int
}

func NewKBufPool(poolSize, bufSize int) (pool *KBufPool) {
	println("New KBufPool(", poolSize, ", bufsize: ", bufSize, ")")

	pool = &KBufPool{bufSize: bufSize}
	for i := 0; i < poolSize; i++ {
		pool.Push(NewKBufObj(bufSize))
	}
	return
}

func (c *KBufPool) Push(buf *KBufObj) {
	if nil != buf {
		buf.Clear()
		c.KStack.Push(buf)
	}
}

func (c *KBufPool) Pop() (buf *KBufObj) {
	if obj, ok := c.KStack.Pop(); ok {
		buf = obj.(*KBufObj)
	} else {
		buf = NewKBufObj(c.bufSize)
	}
	return
}
