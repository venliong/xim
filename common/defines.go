package common

// 业务逻辑类型, 对应nodenet的组
const (
	LOGIC_STATE     = "state"  // 用户状态管理
	LOGIC_TEMPGROUP = "tgroup" // 临时讨论组
)

/* 业务类型
const (
	LogicPushMessage = iota // 发消息

	// 临时讨论组
	LogicTGRecv = iota // 临时讨论组长连接接收
	LogicTGSend = iota // 临时讨论组发消息
)
*/

// 消息显示类型
const (
	MsgText = iota // 文本消息
)
