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

	log.Infoln("tgroup userlogin:", gid, uid)
	user = &User{ID: fmt.Sprintf("%s.%s", gid, uid), ch: make(chan string)}
	sess.Set("info", user)

	g := nodenet.GetGraphByName(xim.API_TEMPGROUP)
	if len(g) < 1 {
		return nil, fmt.Errorf("graph nil:", xim.API_TEMPGROUP)
	}

	iMsg := &xim.MessageTGLogin{Uid: uid, Gid: gid, Access: Conf.NodeName}
	log.Infoln(iMsg, uid, gid, Conf.NodeName)

	cMsg := nodenet.NewMessage(GID.ID(), Conf.NodeName, g, iMsg)
	cMsg.DispenseKey = gid

	if e = nodenet.SendMsgToNext(cMsg); e != nil {
		log.Errorln("SendMsgToNext ERR:", e.Error())
		return nil, e
	}

	return user, nil
}

func TGroutSend(uid, gid, message string) error {
	cMsg := nodenet.NewMessage(GID.ID(), Conf.NodeName, nodenet.GetGraphByName(xim.API_TEMPGROUP), &xim.MessagePushMsg{From: uid, To: gid, Content: message})
	cMsg.DispenseKey = gid
	log.Infoln(cMsg)

	if e := nodenet.SendMsgToNext(cMsg); e != nil {
		log.Errorln("SendMsgToNext ERR:", e.Error())
		return e
	}

	return nil
}
