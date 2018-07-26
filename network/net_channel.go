package network

import (
	"stapler/logger"
	"sync/atomic"
)

type INetChannel interface {
	Logger() (logger.ILogger)
	Addr() *NetAddr
	Attr() interface{}
	SetAttr(interface{})
	Write([]byte)
	String() (string)
}

type NetChannel struct {
	server *NetServer
	key    *NetAddr
	attr   interface{}

	closed   int32
	lastRecv int64
	lastSent int64
}

const NIL  = ""

func (this *NetChannel) Setup(server *NetServer, key *NetAddr) {
	atomic.StoreInt32(&this.closed, 0)
	this.server = server
	this.key = key
}

func (this *NetChannel) Clear() {
	this.key.Clear()
	this.lastRecv = 0
	this.lastSent = 0
	this.attr = nil
}

func (this *NetChannel) Logger() (logger.ILogger) {
	return this.server.serverLog
}

func (this *NetChannel) String() (string) {
	return this.key.String()
}

func (this *NetChannel) Attr() (interface{}) {
	return this.attr
}

func (this *NetChannel) SetAttr(attr interface{}) {
	this.attr = attr
}

func (this *NetChannel) Addr() (*NetAddr) {
	return this.key
}

func (this *NetChannel) LastRecv() (int64) {
	return this.lastRecv
}

func (this *NetChannel) LastSent() (int64) {
	return this.lastSent
}
