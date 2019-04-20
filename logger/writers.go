package logger

type Logwriters struct {
	manager *LogManager
	items   map[string]ILogWriter
}

func NewWrites() *Logwriters {
	return &Logwriters{items: make(map[string]ILogWriter)}
}

func (c *Logwriters) Setup() {
	c.items["file"] = &FileWriter{file: nil, writer: nil}
	c.items["cons"] = &ConsoleWriter{enable: manager.config.Console}
}

func (c *Logwriters) Close() {
	for _, w := range c.items {
		w.Close()
	}
}

func (c *Logwriters) Recv(item *LogItem) {
	line := item.String()
	for _, w := range c.items {
		w.Recv(line)
	}
}

func (c *Logwriters) Flush() {
	for _, w := range c.items {
		w.Flush()
	}
}
