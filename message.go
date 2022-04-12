package chart

const (
	msgTypeLogin     = 1
	msgTypeLogout    = 2
	msgTypeSendText  = 3
	msgTypeSendImage = 4
)

type Message struct {
	ChannelId string `json:"channel_id"` // 管道ID
	UserId    uint64 `json:"user_id"`
	Type      uint8  `json:"type"`    // 消息类型
	Content   string `json:"content"` // 消息内容
	SendTime  string `json:"send_time"`
}

type messageInterface interface {
	SendMessage(msg Message)
	LoginServer(uid uint64)
	LogoutServer(uid uint64)
	ErrorLogServer(err error)
}

var msgStore messageInterface
