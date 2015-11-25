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
	Sig  string
	Conf Config // 系统配置信息
	GID  *gocommon.GlobalID

	users    *session.SessionManager
	mynode   *nodenet.Component
	passport *client.Passport
)

type User struct {
	ID string
	ch chan string
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
	if mynode == nil {
		return fmt.Errorf("No node: ", Conf.NodeName)
	}

	mynode.RegisterHandler(xim.MessagePushMsg{}, dealPushMsg)
	go mynode.Run()

	return nil
}

func SendMsgToUser(fromuserid, touserid, message string) error {
	cMsg := nodenet.NewMessage(GID.ID(), Conf.NodeName, nodenet.GetGraphByName("send"), xim.MessagePushMsg{From: fromuserid, To: touserid, Content: message})
	log.Infoln(cMsg)

	return nodenet.SendMsgToNext(cMsg)
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
	users.SetPrepireRelease(AccessPrepireRelease)

	passport = &client.Passport{ServAddr: Conf.Passport}

	GID = &gocommon.GlobalID{Type: Conf.NodeName}

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

func AccessPrepireRelease(ss session.SessionStore) {
	if ss != nil {
		info := ss.Get("info")
		if info != nil {
			select {
			case info.(*User).ch <- "TIMEOUT":
			default:
			}
			close(info.(*User).ch)
		}
	}
}

func dealPushMsg(data interface{}) (result interface{}, err error) {
	msg := data.(xim.MessagePushMsg)

	sess, _ := users.GetSessionById(&msg.To)
	user := sess.Get("info")
	if user == nil {
		log.Errorln("No such session: ", msg.To)
		return
	}
	log.Infoln("processPushMessage:", user, msg)

	bytemsg, _ := json.Marshal(msg)
	user.(*User).ch <- string(bytemsg)

	return nil, nil
}
