package main

import (
	"fmt"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/session"
	"github.com/liuhengloveyou/xim"

	log "github.com/golang/glog"
)

func init() {
	nodenet.RegisterWorker("tempGroupLogin", xim.MessageTGLogin{}, TempGroupLogin)
	nodenet.RegisterWorker("tempGroupSend", xim.MessagePushMsg{}, TempGroupSend)
}

func TempGroupLogin(data interface{}) (result interface{}, err error) {
	var msg = data.(xim.MessageTGLogin)
	log.Infoln("tgroupLogin:", msg)

	sess, err := session.GetSessionById(&msg.Gid)
	if err != nil {
		return nil, err
	}

	if err := sess.Set(msg.Uid, msg.Access); err != nil {
		return nil, err
	}

	log.Infoln("tgroupLogin OK:", sess)

	return nil, nil
}

func TempGroupSend(data interface{}) (result interface{}, err error) {
	var msg = data.(xim.MessagePushMsg)
	log.Infoln(msg)

	sess, err := session.GetSessionById(&msg.To)
	if err != nil {
		return nil, err
	}

	keys := sess.Keys()
	log.Infoln("tempGroupSend:", msg.To, keys)
	for i := 0; i < len(keys); i++ {
		stat := sess.Get(keys[i])
		if stat == nil || msg.From == keys[i] {
			log.Errorln("tgroup skip:", keys[i], stat)
			continue
		}

		cMsg := nodenet.NewMessage("", "", nil, xim.MessagePushMsg{From: msg.From, To: fmt.Sprintf("%v.%v", msg.To, keys[i]), Group: msg.To, Content: msg.Content})
		log.Infoln("tgroup pushmsg: ", stat.(string), cMsg)

		if err = nodenet.SendMsgToComponent(stat.(string), cMsg); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
