package network

import (
	"net"
	"fmt"
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

func NewTcpChannel(servant *NetServant, conn *net.TCPConn)(*TcpChannel) {
	server := servant.server
	option := server.option

	channel := &TcpChannel{owner:servant,
	                       chSend:make(chan []byte, option.SndQLen),
	                       chRecv:make(chan []byte, option.RcvQLen),
	                       conn:conn,
	                       }

	key_str := fmt.Sprintf("T%.4x-%s/%s", server.ChannelSeq(), servant.LocalAddr(), conn.RemoteAddr().String())
	channel.Setup(servant.server, MakeAddrByString(&key_str))
	return channel
}

func (this *TcpChannel)Open(){
	this.server.Connect(this)

	wg := &sync.WaitGroup{}
	wg.Add(3)

	go this.loopSend(wg)
	go this.loopRecv(wg)
	go this.loopData(wg)

	wg.Wait()
	this.server.Disconnect(this)
}

func (this *TcpChannel)Close()  {
	if atomic.CompareAndSwapInt32(&this.closed, 0, 1){
		close(this.chSend)
		close(this.chRecv)
		this.conn = nil
	}
}

func (this *TcpChannel)loopData(wg *sync.WaitGroup)  {
	log := this.server.channelLog

	defer func() {
		if v := recover(); nil != v{
			log.Error("error on data (%v)", v)
		}
		wg.Add(-1)
		this.Close()
	}()

	for data := range this.chRecv{
		this.server.Data(this, data)
	}
}

func (this *TcpChannel)loopSend(workgroup *sync.WaitGroup)  {
	log := this.server.channelLog
	defer func() {
		if v := recover(); nil != v{
			log.Error("error on send(%v)", v)
		}
		workgroup.Add(-1)
		this.Close()
	}()

	for data := range this.chSend{
		this.doSend(data)
	}
}

func (this *TcpChannel)doSend(data []byte){
	log := this.server.channelLog
	for{
		if nw, e := this.conn.Write(data); nil != e{
			log.Error("%s error on send %s.", this.key, e.Error())
			return
		}else if nw == 0{
			log.Warn("%s error on send Zero.", this.key)
			return
		}else{
			this.server.AddSentBytes(nw)
			data = data[nw:]
			if len(data) == 0{
				break
			}
		}
	}
}

func (this *TcpChannel)loopRecv(wg *sync.WaitGroup)  {
	log := this.server.channelLog

	server := this.server
	bufobj := server.bufPool.Pop()

	defer func() {
		if v := recover(); nil != v{
			log.Error("error on recv(%v)", v)
		}
		server.bufPool.Push(bufobj)
		this.Close()
		wg.Add(-1)
	}()

	buf := bufobj.Data()
	for{
		if nr, e := this.conn.Read(buf); nil != e{
			log.Error("read error: (%s)", e.Error())
			break
		}else if 0 == nr{
			log.Warn("[%s] read Zero", this.key)
			break
		}else{
			this.server.AddRecvBytes(nr)

			data := make([]byte, nr)
			copy(data, buf[:nr])
			this.chRecv <- data
		}
	}
}

func (this *TcpChannel)Write(data []byte){
	this.chSend <- data
}