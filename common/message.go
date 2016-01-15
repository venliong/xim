package common

import (
	"github.com/liuhengloveyou/nodenet"
)

func init() {
	nodenet.RegisterMessageType(MessageLogin{})
	nodenet.RegisterMessageType(MessageLogout{})
	nodenet.RegisterMessageType(MessagePushMsg{})
	nodenet.RegisterMessageType(MessageTGLogin{})
}

// 长连接登入
type MessageLogin struct {
	Userid        string // 用户ID
	ClientType    string // 客户端类型
	AccessName    string // 接入节点名
	AccessSession string // 接入节点会话ID
	UpdateTime    int64  // 状态更新时间
}

// 长连接登出
type MessageLogout struct {
	Userid        string // 用户ID
	ClientType    string // 客户端类型
	AccessName    string // 接入节点名
	AccessSession string // 接入节点会话ID
}

type MessageForward struct {
	ToSession string // 接入节点会话ID
	Content   string `json:"ctx"` // 消息内容
}

// 消息推送
type MessagePushMsg struct {
	From    string `json:"from"`            // 消息发送方用户ID
	To      string `json:"to"`              // 消息接收方ID
	Group   string `json:"group,omitempty"` // 群组ID或空
	Content string `json:"ctx"`             // 消息内容
}

// 随时讨论组登录
type MessageTGLogin struct {
	Gid    string `json:"gid"`
	Uid    string `json:"uid"`
	Access string `json:"access"`
}
