package utils

type KBufPool struct {
	KStack
	bufSize int
}

func NewKBufPool(_pool_size, _buf_size int) (pool *KBufPool) {
	println("New KBufPool(", _pool_size, ", bufsize: ", _buf_size, ")")

	pool = &KBufPool{bufSize: _buf_size}
	for i := 0; i < _pool_size; i++ {
		pool.Push(NewKBufObj(_buf_size))
	}
	return
}

func (this *KBufPool) Push(buf *KBufObj) {
	if nil != buf {
		buf.Clear()
		this.KStack.Push(buf)
	}
}

func (this *KBufPool) Pop() (buf *KBufObj) {
	if obj, ok := this.KStack.Pop(); ok {
		buf = obj.(*KBufObj)
	} else {
		buf = NewKBufObj(this.bufSize)
	}
	return
}
