package network

type TcpServer struct {
	NetServer
	servants map[int]*TcpServant
}

func NewTcpServer() INetServer {
	server := &TcpServer{servants: make(map[int]*TcpServant)}
	server.init(TUN_TCP, server.doPorts, server.doStart, server.doStop)
	return server
}

func (c *TcpServer) WriteTo(addr *NetAddr, data []byte) bool {
	channel := c.Channel(addr)
	if nil != channel {
		channel.Write(data)
	}
	return nil != channel
}

func (c *TcpServer) Channel(addr *NetAddr) INetChannel {
	servant := c.servants[addr.GetLocalPort()]
	if nil == servant {
		return nil
	} else {
		key := addr.String()
		return servant.Channel(&key)
	}
}

func (c *TcpServer) doPorts(ports []int) {
	for _, port := range ports {
		c.servants[port] = NewTcpServant(c, port)
	}
}

func (c *TcpServer) doStart() {
	c.serverLog.Info("dostart.")
	for _, servant := range c.servants {
		servant.doStart()
	}
}

func (c *TcpServer) doStop() {
	for _, servant := range c.servants {
		servant.doStop()
	}
}
