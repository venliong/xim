/*
接入层
*/

package access

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
	"github.com/liuhengloveyou/passport/client"
	"github.com/liuhengloveyou/xim"
)

const (
	HEARTBEAT = time.Duration(3) * time.Second
)

var (
	ConfJson map[string]interface{} // 系统配置信息

	mynode   *nodenet.Component
	Sig      string
	users    *xim.Session // 所有在线用户会话
	passport *client.Passport
)

type User struct {
	ID  string
	ch  chan []byte
	act int64
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
	nodenet.BuildFromConfig("../example/nodenet.conf")

	mynode = nodenet.GetComponentByName(ConfJson["nodeName"].(string))
	if mynode != nil {
		mynode.SetHandler(accessWork)
		go mynode.Run()
	}
}

func accessWork(msg interface{}) (result interface{}, err error) {
	fmt.Println(mynode.Name, msg)

	iMsg := msg.(map[string]interface{})["content"].(map[string]interface{})

	users.Get(iMsg["to"].(string)).(*User).ch <- iMsg["msg"].([]byte)

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

var proto = flag.String("proto", "http", "tcp or http?")

func StartService() {
	flag.Parse()

	sigHandler()

	initNodenet()

	passport = &client.Passport{ServAddr: ConfJson["passport"].(string)}

	switch *proto {
	case "tcp":
		TcpAccess()
	case "http":
		HttpAccess()
	default:
		panic("Error proto: " + *proto)
	}
}
