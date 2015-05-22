/*
HTTP长连接接入
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var Sig string
var ConfJson map[string]interface{} // 系统配置信息

var users map[string]chan []byte

func init() {
	runtime.GOMAXPROCS(8)

	users = make(map[string]chan []byte, 1000000)

	r, err := os.Open("./access.conf")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&ConfJson); err != nil {
		panic(err)
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

var proto = flag.String("proto", "http", "access proto. tcp or http")

func main() {
	flag.Parse()

	sigHandler()

	switch *proto {
	case "tcp":
		TcpAccess()
	case "http":
		HttpAccess()
	default:
		panic("Error proto: " + *proto)
	}
}
