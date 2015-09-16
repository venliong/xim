package main

import (
	"encoding/json"
	"fmt"

	log "github.com/golang/glog"
	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/session"
	"github.com/liuhengloveyou/xim"
)

func init() {
	nodenet.SetWorker(xim.API_TEMPGROUP, tempGroupWorker)
}

func tempGroupWorker(data interface{}) (result interface{}, err error) {
	b, e := json.Marshal(data)
	if e != nil {
		return nil, e
	}

	var msg xim.Message
	if e := json.Unmarshal(b, &msg); e != nil {
		return nil, e
	}

	switch msg.MsgType {
	case xim.MSG_TG_LOGIN:
		err = tempGroupLogin(msg.Content)
		if err != nil {
			log.Errorln("tempGroupLogin ERR:", err.Error())
		}
		log.Infoln("tempGroupLogin OK:", msg.Content)
		return nil, err
	default:
		return nil, fmt.Errorf("末知的消息类型: [%v]", msg.MsgType)
	}

	return nil, nil
}

func tempGroupLogin(data interface{}) error {
	b, e := json.Marshal(data)
	if e != nil {
		return e
	}

	var msg xim.Message_TGLogin
	if e := json.Unmarshal(b, &msg); e != nil {
		return e
	}
	fmt.Println(msg)

	sess, err := session.GetSessionById(msg.Gid)
	if err != nil {
		return err
	}

	if err := sess.Set(msg.Uid, msg.Access); err != nil {
		return err
	}

	fmt.Println(">>>", sess)

	return nil
}
