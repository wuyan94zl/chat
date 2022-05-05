package chart

import (
	"encoding/json"
)

type Hub struct {
	// 所有通道客户端
	clients map[string]map[*Client]bool
	// 发送消息
	broadcast chan Message
	// 注册客户端
	register chan *Client
	// 注销客户端
	unregister chan *Client
}

var globHub *Hub

func newServer() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register: // 登录
			if _, ok := clients[client.Id]; !ok {
				clients[client.Id] = client
				msgStore.LoginServer(client.Id)
			}
		case client := <-h.unregister: // 注销 / 退出
			// 退出所有聊天室
			for _, sk := range client.channelId {
				if v, ok := h.clients[sk]; ok {
					delete(v, client)
				}
			}
			// 退出聊天服务
			delete(clients, client.Id)
			msgStore.LogoutServer(client.Id)
		case message := <-h.broadcast: // 接受消息
			send, _ := json.Marshal(message)
			var sent []uint64
			if message.ToUserId > 0 {
				if cli, ok := clients[message.ToUserId]; ok {
					cli.send <- send // to one user
					sent = append(sent, cli.Id)
				} else {
					if msgStore != nil {
						msgStore.DelaySendMessage(message.ChannelId, message, sent)
					}
				}
			} else {
				if channel, ok := h.clients[message.ChannelId]; ok {
					for cli, _ := range channel {
						cli.send <- send // to all users
						sent = append(sent, cli.Id)
					}
				}
				if msgStore != nil {
					msgStore.DelaySendMessage(message.ChannelId, message, sent)
				}
			}
		}
	}
}
