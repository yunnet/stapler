package logger

import (
	"fmt"
	"stapler/utils"
	"sync"
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

func (c *LogManager) run() {
	defer func() {
		writers.Close()
		for k := range c.insts {
			delete(c.insts, k)
		}
	}()

	terminated := false
	needFlush := false
	lastFlush := time.Now()

	var item *LogItem
	const itv = time.Second * 1

	for !terminated {
		if ref, ok := c.items.Poll(itv); ok {
			needFlush = true
			item = ref.(*LogItem)
			switch item.itemType {
			case itDATA:
				writers.Recv(item)
			case itDONE:
				terminated = true
			}
		} else {
			if needFlush && lastFlush.Add(itv).Before(time.Now()) {
				writers.Flush()
				lastFlush = time.Now()
				needFlush = false
			}
		}
	}
	fmt.Println("end of loop no write log")
}

func (c *LogManager) Logger(name string) ILogger {
	defer c.locks.Unlock()
	c.locks.Lock()

	l := c.insts[name]
	if nil == l {
		cfg := NewConfig()
		cfg.copyForm(c.config)

		l = &LogInstance{config: cfg, key: name, items: c.items}
		c.insts[name] = l
	}
	return l
}
