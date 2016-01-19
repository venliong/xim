package common

import (
	"github.com/liuhengloveyou/nodenet"
)

func init() {
	nodenet.RegisterMessageType(MessageLogin{})
	nodenet.RegisterMessageType(MessageLogout{})
	nodenet.RegisterMessageType(MessageForward{})
	nodenet.RegisterMessageType(MessageTGLogin{})
}

// 业务逻辑类型, 对应nodenet的组
const (
	LOGIC_STATE     = "state"   // 用户状态管理
	LOGIC_FORWARD   = "forward" // 消息推送
	LOGIC_TEMPGROUP = "tgroup"  // 临时讨论组
)

// 消息显示类型
const (
	MSG_TEXT  = "txt"   // 文本消息
	MSG_ICON  = "icon"  // 表情图标
	MSG_IMSG  = "img"   // 图片
	MSG_SOUND = "sound" // 语音
)

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

// 消息路由
type MessageForward struct {
	FromUserid  string `json:"fromuser"`          // 消息发送方用户ID
	FromeAccess string `json:"access,omitempty"`  // 发送用户所在接入点
	ToUserid    string `json:"touser"`            // 消息接收方ID
	ToAccess    string `json:"toaccess"`          // 消息接收方所在接入节点名
	ToSession   string `json:"tosession"`         // 消息接收方在接入节点上的会话ID
	ToGroupId   string `json:"togroup,omitempty"` // 群组ID或空
	ShowType    string `json:"type"`              // 消息显示类型
	Content     string `json:"ctx"`               // 消息内容
}

// 随时讨论组登录
type MessageTGLogin struct {
	Gid    string `json:"gid"`
	Uid    string `json:"uid"`
	Access string `json:"access"`
}
