/*
接入层
*/

package face

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/liuhengloveyou/xim/common"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/client"
	"github.com/liuhengloveyou/passport/session"

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

var (
	confile = flag.String("access_conf", "example/access.conf.sample", "接入服务配置文件路径.")
	proto   = flag.String("access_proto", "http", "接入服务网络协议.")
)

func AccessMain() {
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

	switch *proto {
	case "tcp":
		TcpAccess()
	case "http":
		HttpAccess()
	default:
		panic("Error proto: " + *proto)
	}
}

func initNodenet(fn string) error {
	if e := nodenet.BuildFromConfig(fn); e != nil {
		return e
	}

	mynode = nodenet.GetComponentByName(Conf.NodeName)
	if mynode == nil {
		return fmt.Errorf("No node: ", Conf.NodeName)
	}

	mynode.RegisterHandler(common.MessagePushMsg{}, dealPushMsg)
	go mynode.Run()

	return nil
}

func SendMsgToUser(fromuserid, touserid, message string) error {
	cMsg := nodenet.NewMessage(GID.ID(), Conf.NodeName, nodenet.GetGraphByName("send"), common.MessagePushMsg{From: fromuserid, To: touserid, Content: message})
	log.Infoln(cMsg)

	return nodenet.SendMsgToNext(cMsg)
}

func AccessPrepireRelease(ss session.SessionStore) {
	if ss != nil {
		user := ss.Get("info")
		if user != nil {
			user.(*common.User).Destroy()
		}
	}
}

func dealPushMsg(data interface{}) (result interface{}, err error) {
	msg := data.(common.MessagePushMsg)

	sess, _ := users.GetSessionById(&msg.To)
	user := sess.Get("info")
	if user == nil {
		log.Errorln("No such session: ", msg.To)
		return
	}

	bytemsg, _ := json.Marshal(msg)
	log.Infoln("processPushMessage:", user, string(bytemsg))

	user.(*common.User).PushMessage(string(bytemsg))

	return nil, nil
}
