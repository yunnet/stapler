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

func (c *Config) LogFileName(ymd, hms *string) string {
	return fmt.Sprintf("%s/%s/%s.%s.log", c.RootPath, *ymd, c.AppName, *hms)
}

func (c *Config) copyForm(cfg *Config) {
	if nil == cfg {
		fmt.Println("*logger.Config.copyFrom: warn : config is null.")
		return
	}

	if len(cfg.AppName) > 0 {
		c.AppName = cfg.AppName
	}
	root := cfg.RootPath
	cnt := len(root)
	if cnt > 0 {
		if root[cnt-1:] == "/" {
			root = root[0 : cnt-1]
		}
	} else {
		root, _ = os.Getwd()
	}

	c.RootPath = root
	c.Level = cfg.Level
	c.Source = cfg.Source
	c.Console = cfg.Console
	c.GenFile = cfg.GenFile
}
