package network

import "fmt"

type NetServant struct {
	server *NetServer
	host   string
	port   int
}

func (c *NetServant) Setup(server *NetServer, host string, port int) {
	c.server = server
	c.host = host
	c.port = port
}

func (c *NetServant) LocalAddr() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}

func (c *NetServant) Key() string {
	return fmt.Sprintf("tcp://%s:%d", c.host, c.port)
}
