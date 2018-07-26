package utils

import "encoding/hex"

type KBufObj struct {
	data []byte
	slice []byte
}

func NewKBufObj(size int)(*KBufObj)  {
	bytes := make([]byte, size)
	return &KBufObj{data: bytes, slice:bytes[0:0]}
}

func (this *KBufObj)Write(buf ...byte)  {
	this.slice = append(this.slice, buf...)
}

func (this *KBufObj)WriteBuf(buf *KBufObj)  {
	this.slice = append(this.slice, buf.slice...)
}

func (this *KBufObj)WriteString(str string)  {
	this.Write([]byte(str)...)
}

func (this *KBufObj)Clear(){
	this.slice = this.data[0:0]
}

func (this *KBufObj)Size()(int){
	return len(this.slice)
}

func (this *KBufObj)ByteAt(idx int)(byte){
	return this.slice[idx]
}

func (this *KBufObj)GetMem(size int)([]byte){
	return this.data[0:size]
}

func (this *KBufObj)Read(size int, dst *[]byte){
	if size > this.Size(){
		size = this.Size()
	}

	if nil != dst{
		*dst = append(*dst, this.slice[0:size]...)
	}

	this.slice = this.slice[size:]
}

func (this *KBufObj)IndexOf(src byte, offset int)(int){
	for i:=offset; i<len(this.slice) ; i++ {
		if this.slice[i] == src{
			return i
		}
	}
	return -1
}

func (this *KBufObj)ToHex()(string){
	return hex.EncodeToString(this.slice)
}

func (this *KBufObj)Slice()([]byte){
	return this.slice
}

func (this *KBufObj)Data()([]byte){
	return this.data
}

func (this *KBufObj)String()(string){
	return string(this.slice)
}