/*
 * 用户长连接在线状态
 */
package service

import (
	"fmt"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/session"
	"github.com/liuhengloveyou/xim/common"

	log "github.com/golang/glog"
)

func StateUpdate(sess session.SessionStore) (info *UserSession, e error) {
	var user *User
	if sess.Get("user") == nil {
		log.Errorln("session ERR: ", sess)
		return nil, fmt.Errorf("会话错误.")
	}
	user = sess.Get("user").(*User)

	if sess.Get("info") != nil {
		log.Infoln("state need not update:", sess)
		return sess.Get("info").(*UserSession), nil
	}

	g := nodenet.GetGraphByName(common.LOGIC_STATE)
	if len(g) < 1 {
		return nil, fmt.Errorf("graph not config:", common.LOGIC_STATE)
	}

	cMsg := nodenet.NewMessage(common.GID.ID(), common.AccessConf.NodeName, g,
		common.MessageLogin{
			Userid:        user.Userid,
			ClientType:    user.Client,
			AccessName:    common.AccessConf.NodeName,
			AccessSession: sess.Id("")})
	cMsg.DispenseKey = user.Userid

	log.Infoln("StateUpdate:", cMsg, cMsg.Payload)
	if e = nodenet.SendMsgToNext(cMsg); e != nil {
		log.Errorln("SendMsgToNext ERR:", e.Error())
		return nil, e
	}

	log.Infoln("userlogin:", sess)
	info = NewUserSession(sess.Id(""))
	sess.Set("info", info)

	return info, nil
}
