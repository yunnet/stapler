package utils

import "encoding/hex"

type KBufObj struct {
	data  []byte
	slice []byte
}

func NewKBufObj(size int) *KBufObj {
	bytes := make([]byte, size)
	return &KBufObj{data: bytes, slice: bytes[0:0]}
}

func (c *KBufObj) Write(buf ...byte) {
	c.slice = append(c.slice, buf...)
}

func (c *KBufObj) WriteBuf(buf *KBufObj) {
	c.slice = append(c.slice, buf.slice...)
}

func (c *KBufObj) WriteString(str string) {
	c.Write([]byte(str)...)
}

func (c *KBufObj) Clear() {
	c.slice = c.data[0:0]
}

func (c *KBufObj) Size() int {
	return len(c.slice)
}

func (c *KBufObj) ByteAt(idx int) byte {
	return c.slice[idx]
}

func (c *KBufObj) GetMem(size int) []byte {
	return c.data[0:size]
}

func (c *KBufObj) Read(size int, dst *[]byte) {
	if size > c.Size() {
		size = c.Size()
	}

	if nil != dst {
		*dst = append(*dst, c.slice[0:size]...)
	}

	c.slice = c.slice[size:]
}

func (c *KBufObj) IndexOf(src byte, offset int) int {
	for i := offset; i < len(c.slice); i++ {
		if c.slice[i] == src {
			return i
		}
	}
	return -1
}

func (c *KBufObj) ToHex() string {
	return hex.EncodeToString(c.slice)
}

func (c *KBufObj) Slice() []byte {
	return c.slice
}

func (c *KBufObj) Data() []byte {
	return c.data
}

func (c *KBufObj) String() string {
	return string(c.slice)
}
