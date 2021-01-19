package chart

// 写入聊天记录信息
func sendMessage(userId uint64, msg Message) {
	if DB != nil {
		send := ChatMessage{}
		send.Type = msg.Type
		send.ChannelId = msg.ChannelId
		send.SendUserId = userId
		send.Content = msg.Content
		DB.Create(&send)
	}
}
