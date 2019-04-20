package network

import (
	"stapler/logger"
	"sync/atomic"
)

type INetChannel interface {
	Logger() logger.ILogger
	Addr() *NetAddr
	Attr() interface{}
	SetAttr(interface{})
	Write([]byte)
	String() string
}

type NetChannel struct {
	server *NetServer
	key    *NetAddr
	attr   interface{}

	closed   int32
	lastRecv int64
	lastSent int64
}

const NIL = ""

func (c *NetChannel) Setup(server *NetServer, key *NetAddr) {
	atomic.StoreInt32(&c.closed, 0)
	c.server = server
	c.key = key
}

func (c *NetChannel) Clear() {
	c.key.Clear()
	c.lastRecv = 0
	c.lastSent = 0
	c.attr = nil
}

func (c *NetChannel) Logger() logger.ILogger {
	return c.server.serverLog
}

func (c *NetChannel) String() string {
	return c.key.String()
}

func (c *NetChannel) Attr() interface{} {
	return c.attr
}

func (c *NetChannel) SetAttr(attr interface{}) {
	c.attr = attr
}

func (c *NetChannel) Addr() *NetAddr {
	return c.key
}

func (c *NetChannel) LastRecv() int64 {
	return c.lastRecv
}

func (c *NetChannel) LastSent() int64 {
	return c.lastSent
}
