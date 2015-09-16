package xim

import (
//"encoding/json"

)

// 消息类型
const (
	MSG_SENDMSG = iota // 发消息

	// 临时讨论组
	MSG_TG_LOGIN = iota // 临时讨论组登录
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

type Message_TGLogin struct {
	Uid    string `json:"uid"`
	Gid    string `json:"gid"`
	Access string `json:"access"`
}
