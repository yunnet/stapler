package logger

import (
	"fmt"
	"os"
	"sync/atomic"
)

var (
	manager = NewManager()
	writers = NewWrites()
)

type ILogWriter interface {
	Open()
	Recv(line string)
	Flush()
	Close()
}

type ILogger interface {
	Trace(string, ...interface{})
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
	Critical(string, ...interface{})
	Fatal(string, ...interface{})
}

func Setup(cfg *Config) {
	if atomic.CompareAndSwapInt32(&manager.active, 0, 1) {
		manager.config.copyForm(cfg)
		writers.Setup()
		go manager.run()
	} else {
		fmt.Println("logger setup already.")
	}
}

func Close() {
	if atomic.CompareAndSwapInt32(&manager.active, 1, 0) {
		done := &LogItem{itemType: itDONE}
		manager.items.Offer(done)
	}
}

func Logger(name string) ILogger {
	if atomic.LoadInt32(&manager.active) < 1 {
		println("please call logger.Init(config) first.")
		os.Exit(0)
	}
	return manager.Logger(name)
}
