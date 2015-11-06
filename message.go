package xim

import (
//"encoding/json"

)

type Message struct {
	LogicType int         `json:"type"`
	Content   interface{} `json:"content"`
}

type Message_PushMsg struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Group   string `json:"group,omitempty"`
	Content string `json:"ctx"`
}

type Message_TGLogin struct {
	Gid    string `json:"gid"`
	Uid    string `json:"uid"`
	Access string `json:"access"`
}
