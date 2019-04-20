package network

import (
	"log"
	"sync/atomic"
)

type ServerCallback struct {
	Connect func(INetChannel)
	Disconnect func(INetChannel)
	Data func(INetChannel, []byte)
}

func NewServerCallback(onConnect, onDisconnect func(INetChannel), onData func(INetChannel, []byte)) *ServerCallback {
	return &ServerCallback{Connect:onConnect, Disconnect:onDisconnect, Data:onData}
}

func (c *ServerCallback)copyFrom(src *ServerCallback)  {
	if nil != src{
		if nil != src.Connect{
			c.Connect = src.Connect
		}

		if nil != src.Disconnect{
			c.Disconnect = src.Disconnect
		}

		if nil != src.Data{
			c.Data = src.Data
		}
	}
}

var defLinkCountS int64

func OnConnect(channel INetChannel)  {
	log.Printf("++++++ OnConnect %s %d \r\n", channel.Addr(), atomic.AddInt64(&defLinkCountS, 1))
}

func OnDisconnect(channel INetChannel)  {
	log.Printf("------ OnDisconect %s \r\n", channel.Addr())
}

func OnData(channel INetChannel, data []byte)  {
	channel.Write(data)
}
