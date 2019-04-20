package logger

import (
	"fmt"
	"runtime"
	"stapler/utils"
	"strings"
)

type LogInstance struct {
	key    string
	items  *utils.KQueue
	config *Config
}

func (c *LogInstance) doLog(level LogLevel, _fmt string, args ...interface{}) {
	if level < c.config.Level {
		println("logger skip", level, _fmt, args)
		return
	}

	var (
		file = ""
		rows = 0
		ok   = false
	)

	if c.config.Source {
		if _, file, rows, ok = runtime.Caller(2); ok {
			if idx := strings.LastIndex(file, "/"); idx > 0 {
				file = file[idx+1:]
			}
		}
	}

	msg := _fmt
	if len(args) > 0 {
		msg = fmt.Sprintf(_fmt, args...)
	}

	item := NewLogItem(c, level, msg, file, rows)
	c.items.Offer(item)
}

func (c LogInstance) Trace(_fmt string, args ...interface{}) {
	c.doLog(TRACE, _fmt, args...)
}

func (c LogInstance) Debug(_fmt string, args ...interface{}) {
	c.doLog(DEBUG, _fmt, args...)
}

func (c LogInstance) Info(_fmt string, args ...interface{}) {
	c.doLog(INFO, _fmt, args...)
}

func (c LogInstance) Warn(_fmt string, args ...interface{}) {
	c.doLog(WARN, _fmt, args...)
}

func (c LogInstance) Error(_fmt string, args ...interface{}) {
	c.doLog(ERROR, _fmt, args...)
}

func (c LogInstance) Critical(_fmt string, args ...interface{}) {
	c.doLog(CRITICAL, _fmt, args...)
}

func (c LogInstance) Fatal(_fmt string, args ...interface{}) {
	c.doLog(FATAL, _fmt, args...)
}
