package network

import "fmt"

type NetServant struct {
	server *NetServer
	host string
	port int
}

func (this *NetServant)Setup(server *NetServer, host string, port int)  {
	this.server = server
	this.host = host
	this.port = port
}

func (this *NetServant) LocalAddr()(string){
	return fmt.Sprintf("%s:%d", this.host, this.port)
}

func (this *NetServant) Key()(string){
	return fmt.Sprintf("tcp://%s:%d", this.host, this.port)
}
