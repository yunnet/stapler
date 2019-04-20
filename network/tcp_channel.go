package network

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
)

type TcpChannel struct {
	NetChannel
	owner *NetServant
	conn net.Conn
	chSend chan []byte
	chRecv chan []byte
}

func NewTcpChannel(servant *NetServant, conn *net.TCPConn) *TcpChannel {
	server := servant.server
	option := server.option

	channel := &TcpChannel{owner:servant,
	                       chSend:make(chan []byte, option.SndQLen),
	                       chRecv:make(chan []byte, option.RcvQLen),
	                       conn:conn,
	                       }

	keyStr := fmt.Sprintf("T%.4x-%s/%s", server.ChannelSeq(), servant.LocalAddr(), conn.RemoteAddr().String())
	channel.Setup(servant.server, MakeAddrByString(&keyStr))
	return channel
}

func (c *TcpChannel)Open(){
	c.server.Connect(c)

	wg := &sync.WaitGroup{}
	wg.Add(3)

	go c.loopSend(wg)
	go c.loopRecv(wg)
	go c.loopData(wg)

	wg.Wait()
	c.server.Disconnect(c)
}

func (c *TcpChannel)Close()  {
	if atomic.CompareAndSwapInt32(&c.closed, 0, 1){
		close(c.chSend)
		close(c.chRecv)
		c.conn = nil
	}
}

func (c *TcpChannel)loopData(wg *sync.WaitGroup)  {
	log := c.server.channelLog

	defer func() {
		if v := recover(); nil != v{
			log.Error("error on data (%v)", v)
		}
		wg.Add(-1)
		c.Close()
	}()

	for data := range c.chRecv{
		c.server.Data(c, data)
	}
}

func (c *TcpChannel)loopSend(workgroup *sync.WaitGroup)  {
	log := c.server.channelLog
	defer func() {
		if v := recover(); nil != v{
			log.Error("error on send(%v)", v)
		}
		workgroup.Add(-1)
		c.Close()
	}()

	for data := range c.chSend{
		c.doSend(data)
	}
}

func (c *TcpChannel)doSend(data []byte){
	log := c.server.channelLog
	for{
		if nw, e := c.conn.Write(data); nil != e{
			log.Error("%s error on send %s.", c.key, e.Error())
			return
		}else if nw == 0{
			log.Warn("%s error on send Zero.", c.key)
			return
		}else{
			c.server.AddSentBytes(nw)
			data = data[nw:]
			if len(data) == 0{
				break
			}
		}
	}
}

func (c *TcpChannel)loopRecv(wg *sync.WaitGroup)  {
	log := c.server.channelLog

	server := c.server
	bufobj := server.bufPool.Pop()

	defer func() {
		if v := recover(); nil != v{
			log.Error("error on recv(%v)", v)
		}
		server.bufPool.Push(bufobj)
		c.Close()
		wg.Add(-1)
	}()

	buf := bufobj.Data()
	for{
		if nr, e := c.conn.Read(buf); nil != e{
			log.Error("read error: (%s)", e.Error())
			break
		}else if 0 == nr{
			log.Warn("[%s] read Zero", c.key)
			break
		}else{
			c.server.AddRecvBytes(nr)

			data := make([]byte, nr)
			copy(data, buf[:nr])
			c.chRecv <- data
		}
	}
}

func (c *TcpChannel)Write(data []byte){
	c.chSend <- data
}