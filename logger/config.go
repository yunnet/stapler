package logger

import (
	"fmt"
	"os"
)

type Config struct {
	AppName  string
	RootPath string
	Level    LogLevel
	Source   bool
	Console  bool
	GenFile  bool
}

func NewConfig() *Config {
	return &Config{AppName: "app", RootPath: ".", Source: true, Level: INFO, Console: true}
}

func (this *Config) LogFileName(ymd, hms *string) string {
	return fmt.Sprintf("%s/%s/%s.%s.log", this.RootPath, *ymd, this.AppName, *hms)
}

func(this *Config) copyForm(cfg *Config){
	if nil == cfg{
		println("*logger.Config.copyFrom: warn : config is null.")
		return
	}

	if len(cfg.AppName) > 0{
		this.AppName = cfg.AppName
	}
	root := cfg.RootPath
	cnt := len(root)
	if cnt > 0{
		if root[cnt - 1:] == "/"{
			root = root[0:cnt - 1]
		}
	}else{
		root, _ = os.Getwd()
	}

	this.RootPath = root
	this.Level = cfg.Level
	this.Source = cfg.Source
	this.Console = cfg.Console
	this.GenFile = cfg.GenFile
}