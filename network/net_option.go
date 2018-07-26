package network

import "stapler/utils"

type KOption struct {
	RcvQLen int
	SndQLen int
	BufPool *utils.KBufPool
	PoolSize int
	BufSize int
	Timeout int // 5分钟没有收到数据认为掉线
}

func NewKOption()(*KOption)  {
	return &KOption{
		RcvQLen:8,
		SndQLen:8,
		PoolSize:8192,
		BufSize:1024,
		Timeout:300,
	}
}

func (this KOption)copyFrom(src *KOption)  {
	if nil != src{
		if src.RcvQLen > 0{
			this.RcvQLen = src.RcvQLen
		}

		if src.SndQLen > 0{
			this.SndQLen = src.SndQLen
		}

		if src.PoolSize > 0{
			this.PoolSize = src.PoolSize
		}

		if src.BufSize > 0{
			this.BufSize = src.BufSize
		}

		if src.BufPool != nil{
			this.BufPool = src.BufPool
		}
	}
}

func (this KOption) GetBufPool()(*utils.KBufPool)  {
	if nil == this.BufPool{
		return utils.NewKBufPool(this.PoolSize, this.BufSize)
	}else{
		return this.BufPool
	}
}