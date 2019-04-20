package network

import (
	"net"
	"sync"
)

type TcpServant struct {
	NetServant
	server   *TcpServer
	lsnr     *net.TCPListener
	lock     sync.RWMutex
	channels map[string]INetChannel
}

func NewTcpServant(server *TcpServer, port int) (servent *TcpServant) {
	server.serverLog.Info("New tcp servant %s:%d", server.host, port)
	servent = &TcpServant{server: server, lock: sync.RWMutex{}, channels: make(map[string]INetChannel)}
	servent.Setup(&server.NetServer, server.host, port)
	return
}

func (c *TcpServant) Channel(key *string) (channel INetChannel) {
	c.lock.RLock()
	channel = c.channels[*key]
	c.lock.RUnlock()
	return
}

func (c *TcpServant) doStop() {
	log := c.server.lsnrLog

	defer func() {
		if v := recover(); nil != v {
			log.Error("error on clean %v", v)
		}
		c.lock.RUnlock()
	}()

	c.lsnr.Close()
	c.lock.RLock()
	for _, ref := range c.channels {
		channel := ref.(*TcpChannel)
		channel.conn.Close()
	}
}

func (c *TcpServant) doStart() {
	log := c.server.lsnrLog

	addr := c.LocalAddr()
	log.Debug("doStart %s", addr)

	if laddr, e := net.ResolveTCPAddr("tcp4", addr); e != nil {
		panic(e)
	} else {
		log.Info("try listen on %s", c.Key())
		if lsnr, err := net.ListenTCP("tcp", laddr); nil != err {
			panic(err)
		} else {
			c.lsnr = lsnr
			go c.loopAccept()
			log.Info("listen on %s ok.", c.Key())
		}
	}
}

func (c *TcpServant) loopAccept() {
	log := c.server.lsnrLog

	for c.server.IsActive() {
		if acceptConn, err := c.lsnr.AcceptTCP(); nil != err {
			panic(err)
		} else {
			go func(conn *net.TCPConn) {
				channel := NewTcpChannel(&c.NetServant, acceptConn)
				go c.openChannel(channel)
			}(acceptConn)
		}
	}

	log.Info("---------loopAccept.end")
}

func (c *TcpServant) openChannel(channel *TcpChannel) {
	key := channel.key.String()

	c.lock.Lock()
	c.channels[key] = channel
	c.lock.Unlock()

	channel.Open()

	c.lock.Lock()
	delete(c.channels, key)
	c.lock.Unlock()

	channel.Clear()
}
