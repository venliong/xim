package main

import (
	"fmt"

	"github.com/liuhengloveyou/nodenet"
	"github.com/liuhengloveyou/xim"

	log "github.com/golang/glog"
)

func TGroutRecv(uid, gid string) error {
	g := nodenet.GetGraphByName(xim.API_TEMPGROUP)
	if len(g) < 1 {
		return fmt.Errorf("graph nil:", xim.API_TEMPGROUP)
	}

	iMsg := xim.Message{xim.LogicTGRecv, &xim.Message_TGLogin{Uid: uid, Gid: gid, Access: Conf.NodeName}}
	log.Infoln(iMsg, uid, gid, Conf.NodeName)

	cMsg, e := nodenet.NewMessage(Conf.NodeName, g, iMsg)
	if e != nil {
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

func TGroutSend(uid, gid, message string) error {
	g := nodenet.GetGraphByName(xim.API_TEMPGROUP)
	if len(g) < 1 {
		return fmt.Errorf("graph nil:", xim.API_TEMPGROUP)
	}

	cMsg, e := nodenet.NewMessage(Conf.NodeName, g, &xim.Message{xim.LogicTGSend, &xim.Message_PushMsg{From: uid, To: gid, Content: message}})
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
