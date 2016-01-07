/*
 * 接入层
 */

package face

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/session"
	"github.com/liuhengloveyou/xim/common"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
)

var (
	Sig string
	GID *gocommon.GlobalID

	users  *session.SessionManager
	mynode *nodenet.Component
)

var (
	confile = flag.String("access_conf", "example/access.conf.sample", "接入服务配置文件路径.")
	proto   = flag.String("access_proto", "http", "接入服务网络协议.")
)

func AccessMain() {
	if e := common.InitAccessServ(*confile); e != nil {
		panic(e)
	}

	if e := initNodenet(common.AccessConf.NodeConf); e != nil {
		panic(e)
	}

	users = session.NewSessionManager(common.AccessConf.Session)
	users.SetPrepireRelease(AccessPrepireRelease)

	GID = &gocommon.GlobalID{Type: common.AccessConf.NodeName}

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

	mynode = nodenet.GetComponentByName(common.AccessConf.NodeName)
	if mynode == nil {
		return fmt.Errorf("No node: ", common.AccessConf.NodeName)
	}

	mynode.RegisterHandler(common.MessagePushMsg{}, dealPushMsg)
	go mynode.Run()

	return nil
}

func SendMsgToUser(fromuserid, touserid, message string) error {
	cMsg := nodenet.NewMessage(GID.ID(), common.AccessConf.NodeName, nodenet.GetGraphByName("send"), common.MessagePushMsg{From: fromuserid, To: touserid, Content: message})
	log.Infoln(cMsg)

	return nodenet.SendMsgToNext(cMsg)
}

func AccessPrepireRelease(ss session.SessionStore) {
	if ss != nil {
		user := ss.Get("info")
		if user != nil {
			user.(*common.UserMessage).Destroy()
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

	user.(*common.UserMessage).PushMessage(string(bytemsg))

	return nil, nil
}
