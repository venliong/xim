package face

import (
	"fmt"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/xim/common"

	log "github.com/golang/glog"
)

func TGroutRecv(uid, gid string) (user *common.UserMessage, e error) {
	if uid == "" {
		return nil, fmt.Errorf("userid nil")
	}
	if gid == "" {
		return nil, fmt.Errorf("groupid nil")
	}

	sid := fmt.Sprintf("%s.%s", gid, uid)
	sess, e := users.GetSessionById(&sid)
	if e != nil {
		return nil, e
	}
	info := sess.Get("info")
	if info != nil {
		log.Infoln("tgroup userlogined:", gid, uid)
		return info.(*common.UserMessage), nil
	}

	log.Infoln("tgroup userlogin:", gid, uid)
	sess.Set("info", common.NewUserMessage(fmt.Sprintf("%s.%s", gid, uid)))

	g := nodenet.GetGraphByName(common.API_TEMPGROUP)
	if len(g) < 1 {
		return nil, fmt.Errorf("graph nil:", common.API_TEMPGROUP)
	}

	cMsg := nodenet.NewMessage(GID.ID(), common.AccessConf.NodeName, g, common.MessageTGLogin{Uid: uid, Gid: gid, Access: common.AccessConf.NodeName})
	cMsg.DispenseKey = gid

	if e = nodenet.SendMsgToNext(cMsg); e != nil {
		log.Errorln("SendMsgToNext ERR:", e.Error())
		return nil, e
	}

	return user, nil
}

func TGroutSend(uid, gid, message string) error {
	cMsg := nodenet.NewMessage(GID.ID(), common.AccessConf.NodeName, nodenet.GetGraphByName(common.API_TEMPGROUP), &common.MessagePushMsg{From: uid, To: gid, Content: message})
	cMsg.DispenseKey = gid
	log.Infoln(cMsg)

	if e := nodenet.SendMsgToNext(cMsg); e != nil {
		log.Errorln("SendMsgToNext ERR:", e.Error())
		return e
	}

	return nil
}
