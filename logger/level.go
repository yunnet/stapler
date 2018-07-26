package logger

type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	CRITICAL
	FATAL
	NONE
)

var (
	c_logLevelNames = [...]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "CRITICAL", "FATAL", "NONE"}
)
