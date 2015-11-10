package main

import (
	"fmt"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/xim"

	log "github.com/golang/glog"
)

func TGroutRecv(uid, gid string) (user *User, e error) {
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
		return info.(*User), nil
	}

	g := nodenet.GetGraphByName(xim.API_TEMPGROUP)
	if len(g) < 1 {
		return nil, fmt.Errorf("graph nil:", xim.API_TEMPGROUP)
	}

	iMsg := xim.Message{xim.LogicTGRecv, &xim.Message_TGLogin{Uid: uid, Gid: gid, Access: Conf.NodeName}}
	log.Infoln(iMsg, uid, gid, Conf.NodeName)

	cMsg, e := nodenet.NewMessage(GID.ID(), Conf.NodeName, g, iMsg)
	if e != nil {
		return nil, e
	}
	cMsg.DispenseKey = gid
	log.Infoln(cMsg)

	if e = nodenet.SendMsgToNext(cMsg); e != nil {
		log.Errorln("SendMsgToNext ERR:", e.Error())
		return nil, e
	}

	log.Infoln("tgroup userlogin:", gid, uid)
	user = &User{ID: fmt.Sprintf("%s.%s", gid, uid), ch: make(chan string)}
	sess.Set("info", user)

	return user, nil
}

func TGroutSend(uid, gid, message string) error {
	g := nodenet.GetGraphByName(xim.API_TEMPGROUP)
	if len(g) < 1 {
		return fmt.Errorf("graph nil:", xim.API_TEMPGROUP)
	}

	cMsg, e := nodenet.NewMessage(GID.ID(), Conf.NodeName, g, &xim.Message{xim.LogicTGSend, &xim.Message_PushMsg{From: uid, To: gid, Content: message}})
	if e != nil {
		log.Errorln(e)
		return e
	}
	cMsg.DispenseKey = gid
	log.Infoln(cMsg)

	if e := nodenet.SendMsgToNext(cMsg); e != nil {
		log.Exitln("SendMsgToNext ERR:", e.Error())
		return e
	}

	return nil
}
