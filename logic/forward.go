/*
 * 消息路由
 */

package logic

import (
	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/session"
	"github.com/liuhengloveyou/xim/common"

	log "github.com/golang/glog"
)

func init() {
	nodenet.RegisterWorker("ForwardMessage", common.MessageForward{}, ForwardMessage)
}

func ForwardMessage(data interface{}) (result interface{}, err error) {
	var msg = data.(common.MessageForward)
	log.Infoln("ForwardMessage:", msg)

	sess, err := session.GetSessionById(msg.ToUserid)
	if err != nil {
		return nil, err
	}
	log.Info("tosession:", msg.ToUserid, sess)

	cMsg := nodenet.NewMessage("", "", make([]string, 0), nil)

	keys := sess.Keys()
	for i := 0; i < len(keys); i++ {
		if sess.Get(keys[i]) == nil {
			log.Errorln("stat nil:", msg.ToUserid, keys[i])
			continue
		}
		stat := sess.Get(keys[i]).(*common.MessageLogin)
		log.Infoln("ForwardMessage:", keys[i], stat, msg)

		msg.ToAccess = stat.AccessName
		msg.ToSession = stat.AccessSession

		cMsg.Payload = msg
		log.Infoln("ForwardMessage to:", keys[i], stat, msg)
		if err = nodenet.SendMsgToComponent(stat.AccessName, cMsg); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
