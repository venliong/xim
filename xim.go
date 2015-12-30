package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/liuhengloveyou/xim/face"
)

var Sig string

var (
	service = flag.String("service", "", "启动什么服务? [access | logic]")
)

func main() {
	flag.Parse()

	sigHandler()

	switch *service {
	case "access":
		face.AccessMain()
	default:
		panic("Error service type: " + *service)
	}
}

func sigHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		s := <-c
		Sig = "service is suspend ..."
		fmt.Println("Got signal:", s)
	}()
}
