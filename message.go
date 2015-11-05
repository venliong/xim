package xim

import (
//"encoding/json"

)

// 消息类型
const (
	MSG_PUSHMSG = iota // 推送消息

	// 临时讨论组
	MSG_TG_LOGIN = iota // 临时讨论组登录
)

type Message struct {
	MsgType int         `json:"type"`
	Content interface{} `json:"content"`
}

type Message_PushMsg struct {
	From string `json:"from"`
	To   string `json:"to"`
	Msg  string `json:"msg"`
}

type Message_TGLogin struct {
	Uid    string `json:"uid"`
	Gid    string `json:"gid"`
	Access string `json:"access"`
}
