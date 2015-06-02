/*
消息路由
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	Sig      string
	ConfJson AppConf            // 系统配置信息
	accesss  map[string]*Access // 所有接入服务器
)

type AppConf struct {
	Accesss []string `json:"accesss"`
	Userids []int64  `json:"userids"`
}

type Access struct {
	url  string
	ch   chan []byte
	conn net.Conn
}

func init() {
	runtime.GOMAXPROCS(8)

	r, err := os.Open("./forward.conf")
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

func connToAccess(addr string) {
	return
}

func main() {
	flag.Parse()

	sigHandler()

	for _, one := range ConfJson.Accesss {
		fmt.Println(one)
	}
}
