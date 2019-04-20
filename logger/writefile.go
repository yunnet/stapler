package logger

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type FileWriter struct {
	file   *os.File
	writer *bufio.Writer
	ymd    string
	lines  int
}

func fileExists(filename string) bool {
	exists := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exists = false
	}
	return exists
}

func (c *FileWriter) Recv(line string) {
	if !manager.config.GenFile {
		return
	}

	if 0 == len(c.ymd) {
		c.Open()
	}

	if 0 < len(c.ymd) {
		_, _ = c.writer.WriteString(line)
		c.lines++
	}
}

func (c *FileWriter) Open() {
	filedir := fmt.Sprintf("%s/logs", manager.config.RootPath)

	if !fileExists(filedir) {
		fmt.Printf("new dir %s \r\n", filedir)
		os.MkdirAll(filedir, 0777)
	}

	ymd := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s/%s_%s.log", filedir, manager.config.AppName, ymd)

	var err error
	if fileExists(filename) {
		if c.file, err = os.OpenFile(filename, os.O_APPEND, 0666); nil == err {
			fmt.Printf("open log file %s ok!\r\n", filename)
		}
	} else {
		if c.file, err = os.Create(filename); err == nil {
			fmt.Printf("create log file %s ok! \r\n", filename)
		}
	}

	if nil != err {
		fmt.Printf("err: %s\r\n", err.Error())
	} else {
		c.writer = bufio.NewWriter(c.file)
		c.ymd = ymd
		c.file.WriteString("== open ==\r\n")
	}
}

func (c *FileWriter) Close() {
	if nil != c.writer {
		fmt.Printf("close %s\r\n", c.file.Name())
		c.writer.Flush()
		c.file.Close()

		c.ymd = ""
		c.file = nil
		c.writer = nil
	}
}

func (c *FileWriter) Flush() {
	if c.lines > 0 {
		c.writer.Flush()
		c.lines = 0
	}

	ymd := time.Now().Format("2006-01-02")
	if ymd != c.ymd {
		if len(c.ymd) > 0 {
			fmt.Printf("close log of today: %s\r\n", ymd)
			c.Close()
		}
	}
}
