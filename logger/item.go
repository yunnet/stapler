package logger

import (
	"fmt"
	"sync/atomic"
	"time"
)

type ItemType int

const (
	itDATA ItemType = iota
	itDONE
)

type LogItem struct {
	seq      uint64
	from     *LogInstance
	itemType ItemType
	key      string
	level    LogLevel
	time     time.Time
	msg      string
	file     string
	lines    int
}

var g_log_seq uint64

func NewLogItem(_from *LogInstance, _level LogLevel, _msg, _file string, _lines int) *LogItem {
	atomic.CompareAndSwapUint64(&g_log_seq, 0xEFFFFFFF, 0)
	return &LogItem{
		seq:      atomic.AddUint64(&g_log_seq, 1),
		from:     _from,
		itemType: itDATA,
		key:      _from.key,
		level:    _level,
		time:     time.Now(),
		msg:      _msg,
		file:     _file,
		lines:    _lines,
	}
}

func (c *LogItem) String() string {
	time_str := c.time.Format("2006-01-02 15:04:05")
	if c.from.config.Source {
		return fmt.Sprintf("%.8X %s [%s] %s (%s:%d) : %s\r\n",
			c.seq, time_str, c_logLevelNames[c.level], c.key, c.file, c.lines, c.msg)
	} else {
		return fmt.Sprintf("%.8X %s [%s] %s : %s\r\n",
			c.seq, time_str, c_logLevelNames[c.level], c.key, c.msg)
	}
}
