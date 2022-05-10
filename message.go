package chart

const (
	msgTypeSendText  = 1
	msgTypeSendImage = 2
	msgTypeSendAudio = 3
	msgTypeSendVideo = 4
)

type Message struct {
	ChannelId    string      `json:"channel_id"`    // 管道ID
	ChannelTitle string      `json:"channel_title"` // 管道标题
	UserId       uint64      `json:"user_id"`
	Detail       interface{} `json:"detail"`
	ToUserId     uint64      `json:"to_user_id"`
	Type         uint8       `json:"type"`    // 消息类型
	Content      string      `json:"content"` // 消息内容
	SendTime     string      `json:"send_time"`
}

type messageInterface interface {
	SendMessage(msg Message)
	DelaySendMessage(channelId string, msg Message, uids []uint64)
	LoginServer(uid uint64)
	LogoutServer(uid uint64)
	ErrorLogServer(err error)
}

var msgStore messageInterface
