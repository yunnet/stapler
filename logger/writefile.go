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

func fileExists(filename string) (bool) {
	exists := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exists = false
	}
	return exists
}

func (this *FileWriter) Recv(line string) {
	if !manager.config.GenFile {
		return
	}

	if 0 == len(this.ymd) {
		this.Open()
	}

	if 0 < len(this.ymd) {
		this.writer.WriteString(line)
		this.lines++
	}
}

func (this *FileWriter) Open() {
	filedir := fmt.Sprintf("%s/logs", manager.config.RootPath)

	if !fileExists(filedir) {
		fmt.Printf("new dir %s \r\n", filedir)
		os.MkdirAll(filedir, 0777)
	}

	ymd := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s/%s_%s.log", filedir, manager.config.AppName, ymd)

	var err error
	if fileExists(filename) {
		if this.file, err = os.OpenFile(filename, os.O_APPEND, 0666); nil == err {
			fmt.Printf("open log file %s ok!\r\n", filename)
		}
	} else {
		if this.file, err = os.Create(filename); err == nil {
			fmt.Printf("create log file &s ok! \r\n", filename)
		}
	}

	if nil != err {
		fmt.Printf("err: %s\r\n", err.Error())
	} else {
		this.writer = bufio.NewWriter(this.file)
		this.ymd = ymd
		this.file.WriteString("== open ==\r\n")
	}
}

func (this *FileWriter) Close() {
	if nil != this.writer {
		fmt.Printf("close %s\r\n", this.file.Name())
		this.writer.Flush()
		this.file.Close()

		this.ymd = ""
		this.file = nil
		this.writer = nil
	}
}

func (this *FileWriter) Flush() {
	if this.lines > 0 {
		this.writer.Flush()
		this.lines = 0
	}

	ymd := time.Now().Format("2006-01-02")
	if ymd != this.ymd {
		if len(this.ymd) > 0 {
			fmt.Printf("close log of today: %s\r\n", ymd)
			this.Close()
		}
	}
}
