package logger

import (
	"sync"
	"stapler/utils"
	"time"
)

type LogManager struct {
	writers *Logwriters
	locks   *sync.Mutex
	insts   map[string]ILogger
	items   *utils.KQueue
	config  *Config
	active  int32
}

func NewManager() *LogManager {
	return &LogManager{
		locks:  &sync.Mutex{},
		insts:  make(map[string]ILogger),
		config: NewConfig(),
		items:  utils.NewQueue(),
		active: 0,
	}
}

func (this *LogManager) run() {
	defer func() {
		writers.Close()
		for k, _ := range this.insts {
			delete(this.insts, k)
		}
	}()

	terminated := false
	need_flush := false
	last_flush := time.Now()

	var item *LogItem
	const itv = time.Second * 1

	for !terminated {
		if ref, ok := this.items.Poll(itv); ok {
			need_flush = true
			item = ref.(*LogItem)
			switch item.itemType {
			case itDATA:
				writers.Recv(item)
			case itDONE:
				terminated = true
			}
		} else {
			if need_flush && last_flush.Add(itv).Before(time.Now()) {
				writers.Flush()
				last_flush = time.Now()
				need_flush = false
			}
		}
	}
	println("end of loop no write log")
}

func (this *LogManager) Logger(name string) ILogger {
	defer this.locks.Unlock()
	this.locks.Lock()

	l := this.insts[name]
	if nil == l {
		cfg := NewConfig()
		cfg.copyForm(this.config)

		l = &LogInstance{config: cfg, key: name, items: this.items}
		this.insts[name] = l
	}
	return l
}
