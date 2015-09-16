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

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/client"
	"github.com/liuhengloveyou/passport/session"
	"github.com/liuhengloveyou/xim"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
)

type Config struct {
	Addr     string      `json:"addr"`
	Port     int         `json:"port"`
	NodeName string      `json:"nodeName"`
	NodeConf string      `json:"nodeConf"`
	Passport string      `json:"passport"`
	Session  interface{} `json:"session"`
}

var (
	Conf     Config // 系统配置信息
	users    *session.SessionManager
	mynode   *nodenet.Component
	Sig      string
	passport *client.Passport
)

type User struct {
	ID  string
	ch  chan string
	act int64
}

var (
	confile = flag.String("c", "access.conf.sample", "配置文件路径.")
	proto   = flag.String("p", "http", "接入网络协议.")
)

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
	var sm xim.Message_SendMsg

	json.Unmarshal(b, &sm)
	log.Infoln("sm:", sm)

	sess, _ := users.GetSessionById(sm.ToUser)
	user := sess.Get("info")
	if user == nil {
		log.Errorln("No such user: ", sm.ToUser)
		return nil, nil
	}

	user.(*User).ch <- sm.Msg

	return nil, nil
}

func SendMsgToUser(fromuserid, touserid, message string) error {
	iMsg := xim.Message{xim.MSG_SENDMSG, xim.Message_SendMsg{fromuserid, touserid, message}}
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

func main() {
	flag.Parse()

	if e := gocommon.LoadJsonConfig(*confile, &Conf); e != nil {
		panic(e)
	}

	if e := initNodenet(Conf.NodeConf); e != nil {
		panic(e)
	}

	users = session.NewSessionManager(Conf.Session)
	passport = &client.Passport{ServAddr: Conf.Passport}

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
