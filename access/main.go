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
	"time"
)

const (
	HEARTBEAT = time.Duration(3) * time.Second
)

var (
	Sig      string
	ConfJson map[string]interface{} // 系统配置信息
	users    map[string]*User       // 所有在线用户

)

type User struct {
	ch   chan []byte
	beat int64
}

func init() {
	runtime.GOMAXPROCS(8)

	users = make(map[string]*User, 100000)

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
