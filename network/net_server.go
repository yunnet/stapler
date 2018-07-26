package network

import (
	"stapler/logger"
	"stapler/utils"
	"sync"
	"sync/atomic"
)

type INetServer interface {
	Logger() (logger.ILogger)
	Setup(host string, ports []int, option *KOption)
	Start()
	Stop()
	IsActive() (bool)
	TunnelType() (TunnelType)
	SetCallback(callback *ServerCallback)
	Channel(*NetAddr) (INetChannel)

	ChannelCount() (int32)
	AllRecvBytes() (int64)
	AllSentBytes() (int64)

	WriteTo(dst *NetAddr, msg []byte) (bool)
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

func (n *NetServer) OnPorts() func([]int) {
	return n.onPorts
}

func (n *NetServer) SetOnPorts(onPorts func([]int)) {
	n.onPorts = onPorts
}

func (this *NetServer) init(tunnel TunnelType, onPorts func([]int), onStart func(), onStop func()) {
	this.tunnelType = tunnel
	name := string(tunnel)

	this.serverLog = logger.Logger(name + ".server")
	this.lsnrLog = logger.Logger(name + ".lsnr")
	this.channelLog = logger.Logger(name + "channel")

	this.waitGroup = &sync.WaitGroup{}
	this.onCallback = &ServerCallback{Connect: OnConnect, Disconnect: OnDisconnect, Data: OnData}

	this.onPorts = onPorts
	this.onStart = onStart
	this.onStop = onStop
}

func (this *NetServer) Setup(host string, ports []int, option *KOption) {
	if this.IsActive() {
		panic("can not setup on active status")
	} else {
		this.host = host

		this.option = NewKOption()
		this.option.copyFrom(option)

		this.bufPool = this.option.GetBufPool()

		this.onPorts(ports)
	}
}

func (this *NetServer) Start() {
	if atomic.CompareAndSwapInt32(&this.active, 0, 1) {
		this.onStart()
	} else {
		this.serverLog.Warn("already started.")
	}
}

func (this *NetServer) Stop() {
	if atomic.CompareAndSwapInt32(&this.active, 1, 0) {
		this.onStop()
	} else {
		this.serverLog.Warn("not started.")
	}
}

func (this *NetServer) TunnelType() (TunnelType) {
	return this.tunnelType
}

func (this *NetServer) Logger() (logger.ILogger) {
	return this.serverLog
}

func (this *NetServer) SetCallback(callback *ServerCallback) {
	this.onCallback.copyFrom(callback)
}

func (this *NetServer) IsActive() (bool) {
	return atomic.LoadInt32(&this.active) == 1
}

func (this *NetServer) Connect(channel INetChannel) {
	atomic.AddInt32(&this.channelCount, 1)
	this.onCallback.Connect(channel)
}

func (this *NetServer) Disconnect(channel INetChannel) {
	atomic.AddInt32(&this.channelCount, -1)
	this.onCallback.Disconnect(channel)
}

func (this *NetServer) Data(channel INetChannel, data []byte) {
	this.onCallback.Data(channel, data)
}

func (this *NetServer) ChannelCount() (int32) {
	return atomic.LoadInt32(&this.channelCount)
}

func (this *NetServer) AddRecvBytes(size int) {
	atomic.AddInt64(&this.allRecvBytes, int64(size))
}

func (this *NetServer) AllRecvBytes() (int64) {
	return atomic.LoadInt64(&this.allRecvBytes)
}

func (this *NetServer) AddSentBytes(size int) {
	atomic.AddInt64(&this.allSentBytes, int64(size))
}

func (this *NetServer) AllSentBytes() (int64) {
	return atomic.LoadInt64(&this.allSentBytes)
}

func (this *NetServer) ChannelSeq() (int64) {
	return atomic.AddInt64(&this.channelSeq, 1) & 0xFFFF
}
