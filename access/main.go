/*
接入层
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

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/xim"
)

const (
	HEARTBEAT = time.Duration(3) * time.Second
)

var (
	mynode   *nodenet.Component
	Sig      string
	ConfJson map[string]interface{} // 系统配置信息
	users    *xim.Session           // 所有在线用户会话

)

type User struct {
	ch   chan string
	beat int64
}

func init() {
	runtime.GOMAXPROCS(8)

	users = xim.NewSession()

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

func initNodenet() {
	nodenet.BuildFromConfig("./nodenet.conf")

	mynode = nodenet.GetComponentByName(ConfJson["nodeName"].(string))
	if mynode != nil {
		mynode.SetHandler(accessWork)
		go mynode.Run()
	}
}

func accessWork(msg interface{}) (result interface{}, err error) {
	fmt.Println(mynode.Name, msg)

	iMsg := msg.(map[string]interface{})["content"].(map[string]interface{})

	users.Get(iMsg["to"].(string)).(*User).ch <- iMsg["msg"].(string)

	return nil, nil
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

var proto = flag.String("proto", "http", "tcp || http")

func main() {
	flag.Parse()

	sigHandler()

	initNodenet()

	switch *proto {
	case "tcp":
		TcpAccess()
	case "http":
		HttpAccess()
	default:
		panic("Error proto: " + *proto)
	}
}
