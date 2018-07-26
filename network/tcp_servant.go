package network

import (
	"net"
	"sync"
)

type TcpServant struct {
	NetServant
	server *TcpServer
	lsnr *net.TCPListener
	lock sync.RWMutex
	channels map[string] INetChannel
}

func NewTcpServant(server *TcpServer, port int)(servent *TcpServant)  {
	server.serverLog.Info("New tcp servant %s:%d", server.host, port)
	servent = &TcpServant{server:server, lock:sync.RWMutex{}, channels:make(map[string]INetChannel)}
	servent.Setup(&server.NetServer, server.host, port)
	return
}

func (this *TcpServant) Channel(key *string)(channel INetChannel)  {
	this.lock.RLock()
	channel = this.channels[*key]
	this.lock.RUnlock()
	return
}

func (this *TcpServant)doStop()  {
	log := this.server.lsnrLog

	defer func() {
		if v:= recover(); nil != v{
			log.Error("error on clean %v", v)
		}
		this.lock.RUnlock()
	}()

	this.lsnr.Close()
	this.lock.RLock()
	for _, ref := range this.channels{
		channel := ref.(*TcpChannel)
		channel.conn.Close()
	}
}

func (this *TcpServant)doStart(){
	log := this.server.lsnrLog

	addr := this.LocalAddr()
	log.Debug("doStart %s", addr)

	if laddr, e := net.ResolveTCPAddr("tcp4", addr); e != nil{
		panic(e)
	}else{
		log.Info("try listen on %s", this.Key())
		if lsnr, err := net.ListenTCP("tcp", laddr); nil != err{
			panic(err)
		}else{
			this.lsnr = lsnr
			go this.loopAccept()
			log.Info("listen on %s ok.", this.Key())
		}
	}
}

func (this *TcpServant) loopAccept()  {
	log := this.server.lsnrLog

	for this.server.IsActive(){
		if acceptConn, err := this.lsnr.AcceptTCP(); nil != err{
			panic(err)
		}else{
			go func(conn *net.TCPConn) {
				channel := NewTcpChannel(&this.NetServant, acceptConn)
				go this.openChannel(channel)
			}(acceptConn)
		}
	}

	log.Info("---------loopAccept.end")
}

func (this *TcpServant)openChannel(channel *TcpChannel){
	key := channel.key.String()

	this.lock.Lock()
	this.channels[key] = channel
	this.lock.Unlock()

	channel.Open()

	this.lock.Lock()
	delete(this.channels, key)
	this.lock.Unlock()

	channel.Clear()
}
