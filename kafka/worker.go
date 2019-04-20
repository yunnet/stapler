package kafka

import (
	"github.com/yunnet/stapler/logger"
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

func (this *KWorker) init(name string, on_start, on_open, on_stop func()) {
	this.name = name
	this.lg = logger.Logger(name)
}

func (this *KWorker) getStatus() wsStatusType {
	return wsStatusType(atomic.LoadInt32(&this.status))
}

func (this *KWorker) setStatus(value wsStatusType) {
	atomic.StoreInt32(&this.status, int32(value))
}

func (this *KWorker) Setup(zk_addrs string) {
	this.zkAddrs = strings.Split(zk_addrs, ",")
}

func (this *KWorker) Start() {
	if atomic.CompareAndSwapInt32(&this.active, 0, 1) {
		this.lg.Info("start")

		this.closing = make(chan bool)

		if nil != this.onStart {
			this.lg.Info("onStart")
			this.onStart()
		}

		if nil != this.onOpen {
			this.lg.Debug("check open loop")

			go func() {
				for this.isActive() {
					switch this.getStatus() {
					case wsClosed:
						if nil != this.onOpen {
							this.onOpen()
						}

					case wsConning:
						this.lg.Info("connecting to %s ...", this.zkAddrs)

					case wsConnected:
						//this.lg.Debug("connected[%s] ok, total pkgs: %d", this.name, this.allPkgs)
					}
					time.Sleep(time.Second * 5)
				}
			}()
		}
	}
}

func (this *KWorker) Stop() {
	if atomic.CompareAndSwapInt32(&this.active, 1, 0) {
		this.lg.Info("stop")

		close(this.closing)
		if nil != this.onStop {
			this.onStop()
		}
	}
}

func (this *KWorker) isActive() bool {
	return atomic.LoadInt32(&this.active) > 0
}
