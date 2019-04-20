package logger

import (
	"fmt"
	"github.com/yunnet/stapler/utils"
	"runtime"
	"strings"
)

type LogInstance struct {
	key string
	items *utils.KQueue
	config *Config
}

func (this *LogInstance)doLog(level LogLevel, _fmt string, args...interface{}) {
	if level < this.config.Level{
		println("logger skip", level, _fmt, args)
		return
	}

	var(
		file = ""
		rows = 0
		ok   = false
	)

	if this.config.Source{
		if _, file, rows, ok = runtime.Caller(2); ok{
			if idx := strings.LastIndex(file, "/"); idx > 0{
				file = file[idx + 1:]
			}
		}
	}

	msg := _fmt
	if len(args) > 0{
		msg = fmt.Sprintf(_fmt, args...)
	}

	item := NewLogItem(this, level, msg, file, rows)
	this.items.Offer(item)
}

func(this LogInstance)Trace(_fmt string, args ...interface{}){
	this.doLog(TRACE, _fmt, args...)
}

func(this LogInstance)Debug(_fmt string, args ...interface{}){
	this.doLog(DEBUG, _fmt, args...)
}

func(this LogInstance)Info(_fmt string, args ...interface{}){
	this.doLog(INFO, _fmt, args...)
}

func(this LogInstance)Warn(_fmt string, args ...interface{}){
	this.doLog(WARN, _fmt, args...)
}

func(this LogInstance)Error(_fmt string, args ...interface{}){
	this.doLog(ERROR, _fmt, args...)
}

func(this LogInstance)Critical(_fmt string, args ...interface{}){
	this.doLog(CRITICAL, _fmt, args...)
}

func(this LogInstance)Fatal(_fmt string, args ...interface{}){
	this.doLog(FATAL, _fmt, args...)
}
