package logger

import (
	"testing"
	"time"
)

func Test_log(t *testing.T) {
	//println("hello world test.")
	Setup(&Config{AppName: "demo", RootPath: "g:/temp/", Source: true, Console: true, GenFile: true})

	logger := Logger("demo")
	var seq = 0
	for {
		logger.Info("Hello %d, %s", seq, "log")
		seq++
		time.Sleep(time.Nanosecond * 1)

		if seq > 1000 {
			break
		}
	}

	Close()
}
