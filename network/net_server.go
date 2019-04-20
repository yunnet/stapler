package network

import (
	"github.com/yunnet/stapler/logger"
	"github.com/yunnet/stapler/utils"
	"sync"
	"sync/atomic"
)

type INetServer interface {
	Logger() logger.ILogger
	Setup(host string, ports []int, option *KOption)
	Start()
	Stop()
	IsActive() bool
	TunnelType() TunnelType
	SetCallback(callback *ServerCallback)
	Channel(*NetAddr) INetChannel

	ChannelCount() int32
	AllRecvBytes() int64
	AllSentBytes() int64

	WriteTo(dst *NetAddr, msg []byte) bool
}

type NetServer struct {
	active    int32
	timeout   int64
	bufPool   *utils.KBufPool
	waitGroup *sync.WaitGroup

	serverLog  logger.ILogger
	lsnrLog    logger.ILogger
	channelLog logger.ILogger

	channelSeq   int64
	channelCount int32
	allRecvBytes int64
	allSentBytes int64

	tunnelType TunnelType
	host       string
	option     *KOption
	onPorts    func([]int)
	onStart    func()
	onStop     func()
	onCallback *ServerCallback
}

func (c *NetServer) OnPorts() func([]int) {
	return c.onPorts
}

func (c *NetServer) SetOnPorts(onPorts func([]int)) {
	c.onPorts = onPorts
}

func (c *NetServer) init(tunnel TunnelType, onPorts func([]int), onStart func(), onStop func()) {
	c.tunnelType = tunnel
	name := string(tunnel)

	c.serverLog = logger.Logger(name + ".server")
	c.lsnrLog = logger.Logger(name + ".lsnr")
	c.channelLog = logger.Logger(name + "channel")

	c.waitGroup = &sync.WaitGroup{}
	c.onCallback = &ServerCallback{Connect: OnConnect, Disconnect: OnDisconnect, Data: OnData}

	c.onPorts = onPorts
	c.onStart = onStart
	c.onStop = onStop
}

func (c *NetServer) Setup(host string, ports []int, option *KOption) {
	if c.IsActive() {
		panic("can not setup on active status")
	} else {
		c.host = host

		c.option = NewKOption()
		c.option.copyFrom(option)

		c.bufPool = c.option.GetBufPool()

		c.onPorts(ports)
	}
}

func (c *NetServer) Start() {
	if atomic.CompareAndSwapInt32(&c.active, 0, 1) {
		c.onStart()
	} else {
		c.serverLog.Warn("already started.")
	}
}

func (c *NetServer) Stop() {
	if atomic.CompareAndSwapInt32(&c.active, 1, 0) {
		c.onStop()
	} else {
		c.serverLog.Warn("not started.")
	}
}

func (c *NetServer) TunnelType() TunnelType {
	return c.tunnelType
}

func (c *NetServer) Logger() logger.ILogger {
	return c.serverLog
}

func (c *NetServer) SetCallback(callback *ServerCallback) {
	c.onCallback.copyFrom(callback)
}

func (c *NetServer) IsActive() bool {
	return atomic.LoadInt32(&c.active) == 1
}

func (c *NetServer) Connect(channel INetChannel) {
	atomic.AddInt32(&c.channelCount, 1)
	c.onCallback.Connect(channel)
}

func (c *NetServer) Disconnect(channel INetChannel) {
	atomic.AddInt32(&c.channelCount, -1)
	c.onCallback.Disconnect(channel)
}

func (c *NetServer) Data(channel INetChannel, data []byte) {
	c.onCallback.Data(channel, data)
}

func (c *NetServer) ChannelCount() int32 {
	return atomic.LoadInt32(&c.channelCount)
}

func (c *NetServer) AddRecvBytes(size int) {
	atomic.AddInt64(&c.allRecvBytes, int64(size))
}

func (c *NetServer) AllRecvBytes() int64 {
	return atomic.LoadInt64(&c.allRecvBytes)
}

func (c *NetServer) AddSentBytes(size int) {
	atomic.AddInt64(&c.allSentBytes, int64(size))
}

func (c *NetServer) AllSentBytes() int64 {
	return atomic.LoadInt64(&c.allSentBytes)
}

func (c *NetServer) ChannelSeq() int64 {
	return atomic.AddInt64(&c.channelSeq, 1) & 0xFFFF
}