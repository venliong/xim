package face

import (
	"fmt"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/session"
	"github.com/liuhengloveyou/xim/common"
	"github.com/liuhengloveyou/xim/service"

	log "github.com/golang/glog"
)

func StateUpdate(token string) (info *UserSession, e error) {
	if token == "" {
		return nil, fmt.Errorf("token nil")
	}

	sess, e := session.GetSessionById(token)
	if e != nil {
		log.Errorln("session ERR:", token, e.Error())
		return nil, fmt.Errorf("会话错误.")
	}

	/*
		create := sess.Get("sync")
		if (time.Now().Unix() - create) > 30 {
			log.Infoln("state need not update:", sess)
			return nil
		}
	*/

	var user *service.User
	if sess.Get("user") == nil {
		log.Errorln("session ERR:", token, "no user in sess.")
		return nil, fmt.Errorf("会话错误.")
	}
	user = sess.Get("user").(*service.User)

	if sess.Get("info") != nil {
		log.Infoln("state need not update:", sess)
		return sess.Get("info").(*UserSession), nil
	}

	g := nodenet.GetGraphByName(common.LOGIC_STATE)
	if len(g) < 1 {
		return nil, fmt.Errorf("graph not config:", common.LOGIC_STATE)
	}

	cMsg := nodenet.NewMessage(GID.ID(), common.AccessConf.NodeName, g,
		common.MessageLogin{
			Userid:        user.Userid,
			ClientType:    user.Client,
			AccessName:    common.AccessConf.NodeName,
			AccessSession: sess.Id("")})
	cMsg.DispenseKey = user.Userid

	if e = nodenet.SendMsgToNext(cMsg); e != nil {
		log.Errorln("SendMsgToNext ERR:", e.Error())
		return nil, e
	}

	log.Infoln("userlogin:", sess)
	info = NewUserSession(token)
	sess.Set("info", info)

	return info, nil
}
