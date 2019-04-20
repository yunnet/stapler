package network

import (
	"github.com/yunnet/stapler/logger"
	"testing"
	"time"
)

func TestTcpServer(t *testing.T) {
	logger.Setup(&logger.Config{AppName: "demo", RootPath: "z:/", Source: false, Console: true})

	s := NewTcpServer()

	option := &KOption{BufSize: 1024}
	s.Setup("127.0.0.1", []int{6000}, option)

	s.Start()

	log := logger.Logger("demo")
	for {
		log.Info("links=%d, recv=%d, sents:%d", s.ChannelCount(), s.AllRecvBytes(), s.AllSentBytes())
		time.Sleep(time.Second * 3)
	}
}
