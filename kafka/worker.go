package kafka

import (
	"stapler/logger"
	"strings"
	"sync/atomic"
	"time"
)

type IKWorker interface {
	Start()
	Stop()
}

type KWorker struct {
	name    string
	active  int32
	status  int32
	zkAddrs []string
	lg      logger.ILogger

	closing chan bool
	onStart func()
	onOpen  func()
	onStop  func()
}

type wsStatusType int32

const (
	wsClosed wsStatusType = iota
	wsConning
	wsConnected
)

func (c *KWorker) init(name string, onStart, onOpen, onStop func()) {
	c.name = name
	c.lg = logger.Logger(name)
}

func (c *KWorker) getStatus() wsStatusType {
	return wsStatusType(atomic.LoadInt32(&c.status))
}

func (c *KWorker) setStatus(value wsStatusType) {
	atomic.StoreInt32(&c.status, int32(value))
}

func (c *KWorker) Setup(zkAddrs string) {
	c.zkAddrs = strings.Split(zkAddrs, ",")
}

func (c *KWorker) Start() {
	if atomic.CompareAndSwapInt32(&c.active, 0, 1) {
		c.lg.Info("start")

		c.closing = make(chan bool)

		if nil != c.onStart {
			c.lg.Info("onStart")
			c.onStart()
		}

		if nil != c.onOpen {
			c.lg.Debug("check open loop")

			go func() {
				for c.isActive() {
					switch c.getStatus() {
					case wsClosed:
						if nil != c.onOpen {
							c.onOpen()
						}

					case wsConning:
						c.lg.Info("connecting to %s ...", c.zkAddrs)

					case wsConnected:
						//c.lg.Debug("connected[%s] ok, total pkgs: %d", c.name, c.allPkgs)
					}
					time.Sleep(time.Second * 5)
				}
			}()
		}
	}
}

func (c *KWorker) Stop() {
	if atomic.CompareAndSwapInt32(&c.active, 1, 0) {
		c.lg.Info("stop")

		close(c.closing)
		if nil != c.onStop {
			c.onStop()
		}
	}
}

func (c *KWorker) isActive() bool {
	return atomic.LoadInt32(&c.active) > 0
}
