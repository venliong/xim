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
	"syscall"
	"time"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/client"
	"github.com/liuhengloveyou/xim/common"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
)

const (
	HEARTBEAT = time.Duration(3) * time.Second
)

type Config struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	NodeName string `json:"nodeName"`
	NodeConf string `json:"nodeConf"`
	Passport string `json:"passport"`
}

var (
	Conf     Config // 系统配置信息
	mynode   *nodenet.Component
	Sig      string
	users    *common.Session // 所有在线用户会话
	passport *client.Passport
)

type User struct {
	ID  string
	ch  chan []byte
	act int64
}

func init() {
	if e := gocommon.LoadJsonConfig("access.conf", &Conf); e != nil {
		panic(e)
	}

	if e := initNodenet(Conf.NodeConf); e != nil {
		panic(e)
	}

	users = common.NewSession()

}

func initNodenet(fn string) error {
	if e := nodenet.BuildFromConfig(fn); e != nil {
		return e
	}

	mynode = nodenet.GetComponentByName(Conf.NodeName)
	if mynode != nil {
		mynode.SetHandler(accessWork)
		go mynode.Run()
	}

	return nil
}

func accessWork(msg interface{}) (result interface{}, err error) {
	log.Infoln(Conf.NodeName, "msg: ", msg)

	b, _ := json.Marshal((msg.(map[string]interface{}))["content"])

	log.Infoln(string(b))
	var sm common.Message_SendMsg

	json.Unmarshal(b, &sm)
	log.Infoln("sm:", sm)

	user := users.Get(sm.ToUser)
	if user == nil {
		log.Errorln("No such user: ", sm.ToUser)
		return nil, nil
	}

	user.(*User).ch <- []byte(sm.Msg)

	return nil, nil
}

func SendMsgToUser(fromuserid, touserid, message string) error {
	iMsg := common.Message{common.MSG_SENDMSG, common.Message_SendMsg{fromuserid, touserid, message}}
	fmt.Println("iMsg: ", iMsg)

	g := nodenet.GetGraphByName("send")
	cMsg, _ := nodenet.NewMessage(Conf.NodeName, g, iMsg)
	fmt.Println("cMsg: ", cMsg)

	err := nodenet.SendMsgToNext(cMsg)
	fmt.Println(err)

	return nil
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

func main() {
	flag.Parse()

	sigHandler()

	passport = &client.Passport{ServAddr: Conf.Passport}

	switch *proto {
	case "tcp":
		TcpAccess()
	case "http":
		HttpAccess()
	default:
		panic("Error proto: " + *proto)
	}
}
