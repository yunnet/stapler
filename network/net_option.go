package network

import "stapler/utils"

type KOption struct {
	RcvQLen  int
	SndQLen  int
	BufPool  *utils.KBufPool
	PoolSize int
	BufSize  int
	Timeout  int // 5分钟没有收到数据认为掉线
}

func NewKOption() *KOption {
	return &KOption{
		RcvQLen:  8,
		SndQLen:  8,
		PoolSize: 8192,
		BufSize:  1024,
		Timeout:  300,
	}
}

func (c KOption) copyFrom(src *KOption) {
	if nil != src {
		if src.RcvQLen > 0 {
			c.RcvQLen = src.RcvQLen
		}

		if src.SndQLen > 0 {
			c.SndQLen = src.SndQLen
		}

		if src.PoolSize > 0 {
			c.PoolSize = src.PoolSize
		}

		if src.BufSize > 0 {
			c.BufSize = src.BufSize
		}

		if src.BufPool != nil {
			c.BufPool = src.BufPool
		}
	}
}

func (c KOption) GetBufPool() *utils.KBufPool {
	if nil == c.BufPool {
		return utils.NewKBufPool(c.PoolSize, c.BufSize)
	} else {
		return c.BufPool
	}
}
