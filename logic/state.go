/*
* 用户长连接在线状态更新
 */
package logic

import (
	"time"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/xim/common"

	log "github.com/golang/glog"
	"github.com/liuhengloveyou/passport/session"
)

func init() {
	nodenet.RegisterWorker("UerLogin", common.MessageLogin{}, UerLogin)
	nodenet.RegisterWorker("UerLogout", common.MessageLogout{}, UerLogout)
}

func UerLogin(data interface{}) (result interface{}, err error) {
	var msg = data.(common.MessageLogin)
	log.Infoln(msg)

	if msg.ClientType == "" {
		msg.ClientType = "XIM"
	}

	sess, err := session.GetSessionById(msg.Userid)
	if err != nil {
		return nil, err
	}
	log.Infoln("current session:", msg, sess)

	msg.UpdateTime = time.Now().Unix()
	if err = sess.Set(msg.ClientType, &msg); err != nil {
		return nil, err
	}
	log.Infoln("UserLogin OK:", msg, sess)

	return nil, nil
}

func UerLogout(data interface{}) (result interface{}, err error) {
	//	var msg = data.(*common.MessageLogin)
	return nil, nil
}
