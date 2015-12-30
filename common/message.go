package common

import (
	"github.com/liuhengloveyou/nodenet"
)

func init() {
	nodenet.RegisterMessageType(MessagePushMsg{})
	nodenet.RegisterMessageType(MessageTGLogin{})
}

type MessagePushMsg struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Group   string `json:"group,omitempty"`
	Content string `json:"ctx"`
}

type MessageTGLogin struct {
	Gid    string `json:"gid"`
	Uid    string `json:"uid"`
	Access string `json:"access"`
}
