package common

// 接口类型, 对应nodenet的组
const (
	API_TEMPGROUP = "tgroup" // 临时讨论组
	API_CHAT      = "chat"   // 1对1聊天
)

// 业务类型
const (
	LogicPushMessage = iota // 发消息

	// 临时讨论组
	LogicTGRecv = iota // 临时讨论组长连接接收
	LogicTGSend = iota // 临时讨论组发消息
)

// 消息类型
const (
	MsgText = iota // 文本消息
)

// 用户信息
type UserInfo struct {
	Id        int64  `validate:"-" json:"id,omitempty"`
	Cellphone string `validate:"noneor,cellphone" json:"cellphone,omitempty"`
	Email     string `validate:"noneor,email" json:"email,omitempty"`
	Nickname  string `validate:"noneor,max=20" json:"nickname,omitempty"`
	Password  string `validate:"nonone,min=6,max=64" json:"password,omitempty"`
}
