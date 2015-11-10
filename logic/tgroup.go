package main

import (
	"encoding/json"
	"fmt"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/passport/session"
	"github.com/liuhengloveyou/xim"

	log "github.com/golang/glog"
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

	switch msg.LogicType {
	case xim.LogicTGRecv:
		err = tempGroupLogin(msg.Content)
		if err != nil {
			log.Errorln("tempGroupLogin ERR:", err.Error())
		}
	case xim.LogicTGSend:
		err = tempGroupSend(msg.Content)
		if err != nil {
			log.Errorln("tempGroupLogin ERR:", err.Error())
		}
	default:
		log.Errorf("末知的业务类型: [%v]", msg.LogicType)
	}

	return nil, nil
}

func tempGroupSend(data interface{}) error {
	msg := &xim.Message_PushMsg{}
	b, e := json.Marshal(data)
	if e != nil {
		return e
	}
	if e := json.Unmarshal(b, msg); e != nil {
		return e
	}
	log.Infoln("tgroupsend:", msg)

	sess, err := session.GetSessionById(&msg.To)
	if err != nil {
		return err
	}

	keys := sess.Keys()
	log.Infoln("tempGroupSend:", msg.To, keys)
	for i := 0; i < len(keys); i++ {
		stat := sess.Get(keys[i])
		if stat == nil || msg.From == keys[i] {
			log.Errorln("skip:", keys[i], stat)
			continue
		}

		cMsg, _ := nodenet.NewMessage("", "", nil, xim.Message{xim.LogicPushMessage, &xim.Message_PushMsg{From: msg.From, To: fmt.Sprintf("%v.%v", msg.To, keys[i]), Group: msg.To, Content: msg.Content}})
		log.Infoln("tgroup pushmsg: ", cMsg)
		nodenet.SendMsgToComponent(stat.(string), cMsg)
	}

	return nil
}

func tempGroupLogin(data interface{}) error {
	var msg xim.Message_TGLogin
	b, e := json.Marshal(data)
	if e != nil {
		return e
	}
	if e := json.Unmarshal(b, &msg); e != nil {
		return e
	}
	log.Infoln("tgroupLogin:", msg)

	sess, err := session.GetSessionById(&msg.Gid)
	if err != nil {
		return err
	}

	if err := sess.Set(msg.Uid, msg.Access); err != nil {
		return err
	}

	log.Infoln("tgroupLogin OK:", sess)

	return nil
}
