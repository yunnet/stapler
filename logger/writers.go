package logger

type Logwriters struct {
	manager *LogManager
	items map[string]ILogWriter
}

func NewWrites()*Logwriters {
	return &Logwriters{items:make(map[string]ILogWriter)}
}

func (this *Logwriters)Setup()  {
	this.items["file"] = &FileWriter{file:nil, writer:nil}
	this.items["cons"] = &ConsoleWriter{enable: manager.config.Console}
}

func (this *Logwriters)Close()  {
	for _, w := range this.items{
		w.Close()
	}
}

func (this *Logwriters)Recv(item *LogItem)  {
	line := item.String()
	for _, w := range this.items{
		w.Recv(line)
	}
}

func (this *Logwriters)Flush()  {
	for _, w := range this.items{
		w.Flush()
	}
}


