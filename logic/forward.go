/*
 * 消息路由
 */

package logic

import (
	"time"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/session"
	"github.com/liuhengloveyou/xim/common"

	log "github.com/golang/glog"
)

func init() {
	nodenet.RegisterWorker("ForwardMessage", common.MessageForward{}, ForwardMessage)
	nodenet.RegisterWorker("ConfirmMessage", common.MessageConfirm{}, ConfirmMessage)
}

func ConfirmMessage(data interface{}) (result interface{}, err error) {
	var msg = data.(common.MessageConfirm)
	log.Infof("ConfirmMessage: %#v\n", msg)

	sess, err := session.GetSessionById(msg.FromUserid)
	if err != nil {
		return nil, err
	}

	if sess.Get("info") != nil {
		log.Errorf("ConfirmMessage ERR: no info.")
		return nil, nil
	}
	info := sess.Get("info").(*StateSession)

	if info.Confirm > msg.ConfirmMessage {
		log.Warningf("ConfirmMessage old: %#v. %#v", info, msg)
		return nil, nil
	}

	info.Confirm = msg.ConfirmMessage

	return nil, nil
}

func ForwardMessage(data interface{}) (result interface{}, err error) {
	var msg = data.(common.MessageForward)
	log.Infof("ForwardMessage: %#v\n", msg)
	newid := common.GID.LogicClock(msg.MsgId)
	log.Infof("logicClock: %d->%d\n", msg.MsgId, newid)
	msg.MsgId = newid // 全局排序

	sess, err := session.GetSessionById(msg.ToUserid)
	if err != nil {
		return nil, err
	}
	log.Infof("ForwardMessage tosession: %v, %#v\n", msg.ToUserid, sess)

	if sess.Get("info") != nil {
		log.Errorf("ForwardMessage ERR: no info.")
		return nil, nil
	}
	info := sess.Get("info").(*StateSession)

	msg.FromeAccess = ""
	msg.Time = time.Now().Unix()
	setOfflineMessage(info, &msg) // 先放到离线消息队列

	if info.Alive <= 0 {
		log.Infof("offline: %v. \n", msg.ToUserid)
		return nil, nil
	}

	dealOfflineMessage(sess, info) // 处理离线消息

	return nil, nil
}

func dealOfflineMessage(sess session.SessionStore, info *StateSession) {
	for info.Messages.Len() > 0 {
		info.lock.Lock()

		e := info.Messages.Front()
		one := e.Value.(*common.MessageForward)
		if one.MsgId <= info.Confirm {
			info.Messages.Remove(e)
			info.lock.Unlock()
			continue
		}

		if one.MsgId <= info.Pushed {
			info.lock.Unlock()
			continue
		}
		info.lock.Unlock()

		if pushMessage(sess, one) == true {
			info.Pushed = one.MsgId // 推送成功
		}
	}

}

func pushMessage(sess session.SessionStore, message *common.MessageForward) (ok bool) {
	cMsg := nodenet.NewMessage("", "", make([]string, 0), nil)
	keys := sess.Keys()

	for i := 0; i < len(keys); i++ {
		if keys[i] == "info" {
			continue
		}
		if sess.Get(keys[i]) == nil {
			continue
		}
		stat := sess.Get(keys[i]).(*common.MessageLogin)
		if stat.UpdateTime <= 0 {
			continue // 已退出
		}

		message.ToAccess = stat.AccessName
		message.ToSession = stat.AccessSession
		cMsg.Payload = message
		log.Infof("ForwardMessage: %s, %#v, %#v", keys[i], stat, message)

		if err := nodenet.SendMsgToComponent(stat.AccessName, cMsg); err != nil {
			log.Infof("ForwardMessage ERR: %s, %#v, %#v; %v", keys[i], stat, message, err)
		} else {
			ok = true
		}
	}

	return
}

func setOfflineMessage(info *StateSession, message *common.MessageForward) {
	info.lock.Lock()
	defer info.lock.Unlock()

	info.Messages.PushBack(message)
}
