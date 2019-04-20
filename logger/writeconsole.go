package logger

import "fmt"

type ConsoleWriter struct {
	enable bool
}

func (c *ConsoleWriter)Recv(msg string)  {
	if c.enable{
		fmt.Print(msg)
	}
}

func (c *ConsoleWriter)Open(){

}

func (c *ConsoleWriter)Flush(){

}

func (c *ConsoleWriter)Close(){

}
