package xim

import (
//"encoding/json"
)

const (
	MSG_SENDMSG = iota // 发消息
)

type Message struct {
	MsgType int         `json:"type"`
	Content interface{} `json:"content"`
}

type Message_SendMsg struct {
	FromUser string `json:"from"`
	ToUser   string `json:"to"`
	Msg      string `json:"msg"`
}
