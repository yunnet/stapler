package network

type TcpServer struct {
	NetServer
	servants map[int]*TcpServant
}

func NewTcpServer() INetServer {
	server := &TcpServer{servants:make(map[int]*TcpServant)}
	server.init(TUN_TCP, server.doPorts, server.doStart, server.doStop)
	return server
}

func (this *TcpServer)WriteTo(addr *NetAddr, data []byte) bool {
	channel := this.Channel(addr)
	if nil != channel{
		channel.Write(data)
	}
	return nil != channel
}

func (this *TcpServer)Channel(addr *NetAddr) INetChannel {
	servant := this.servants[addr.GetLocalPort()]
	if nil == servant{
		return nil
	}else{
		key := addr.String()
		return servant.Channel(&key)
	}
}

func (this *TcpServer)doPorts(ports []int){
	for _,port := range ports{
		this.servants[port] = NewTcpServant(this, port)
	}
}

func (this *TcpServer)doStart(){
	this.serverLog.Info("dostart.")
	for _, servant := range this.servants{
		servant.doStart()
	}
}

func (this *TcpServer)doStop(){
	for _, servant := range this.servants{
		servant.doStop()
	}
}