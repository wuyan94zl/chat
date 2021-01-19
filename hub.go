package chart_server

import (
	"encoding/json"
	"strings"
)

var clients map[uint64]*Client

// 管道消息
type Message struct {
	ChannelId string `json:"channel_id"` // 管道ID
	Type      uint8  `json:"type"`       // 消息类型 0：文字消息，1：font表情，2：图片表情
	Content   string `json:"content"`    // 消息内容
	SendTime  string `json:"send_time"`
}

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

func NewServer() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register: // 注册（加入聊天室）
			if v, ok := clients[client.Id]; ok {
				delIds := strings.Split(v.channelId, ",")
				for _, c := range delIds {
					delete(h.clients[c], clients[client.Id])
				}
			}
			clients[client.Id] = client
			channelIds := strings.Split(client.channelId, ",")
			for _, sk := range channelIds {
				if v, ok := h.clients[sk]; ok {
					v[client] = true
				} else {
					v := make(map[*Client]bool)
					v[client] = true
					h.clients[sk] = v
				}
			}
		case client := <-h.unregister: // 注销（退出聊天室）
			channelIds := strings.Split(client.channelId, ",")
			for _, sk := range channelIds {
				if v, ok := h.clients[sk]; ok {
					delete(v, client)
				}
			}
		case message := <-h.broadcast: // 接受消息
			sk := message.ChannelId
			for cli, _ := range h.clients[sk] {
				send, _ := json.Marshal(message)
				cli.send <- send // 向客户端发消息
			}
		}
	}
}
