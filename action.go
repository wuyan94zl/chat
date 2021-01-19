package chart_server

func sendMessage(userId uint64, msg Message) {
	send := ChatMessage{}
	send.Type = msg.Type
	send.ChannelId = msg.ChannelId
	send.SendUserId = userId
	send.Content = msg.Content
	DB.Create(&send)
}
