package logger

import "fmt"

type ConsoleWriter struct {
	enable bool
}

func (this *ConsoleWriter)Recv(msg string)  {
	if this.enable{
		fmt.Print(msg)
	}
}

func (this *ConsoleWriter)Open(){

}

func (this *ConsoleWriter)Flush(){

}

func (this *ConsoleWriter)Close(){

}
