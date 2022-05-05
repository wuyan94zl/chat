package chart

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	send      chan []byte
	channelId []string
	Id        uint64
	Detail    interface{}
}

var clients map[uint64]*Client

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *Client) WriteMsg() {
	defer func() {
		_ = c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{}) // 错误 关闭 channel
				msgStore.ErrorLogServer(fmt.Errorf("系统错误：未知"))
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadMsg() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()
	for {
		_, strByte, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				msgStore.ErrorLogServer(fmt.Errorf("系统错误：%v", err))
			}
			break
		}
		if clients[c.Id] != c {
			msgStore.ErrorLogServer(fmt.Errorf("用户`%d`未登录，不能发送消息", c.Id))
			continue
		}

		message := Message{}
		err = json.Unmarshal(strByte, &message)
		message.SendTime = time.Now().Format("2006-01-02 15:04:05")
		message.UserId = c.Id
		message.Detail = c.Detail
		if message.ToUserId > 0 {
			if _, ok := clients[message.ToUserId]; !ok {
				continue
			}
		} else if _, ok := c.hub.clients[message.ChannelId][c]; !ok {
			msgStore.ErrorLogServer(fmt.Errorf("用户`%d`未监听`%s`频道，不能发送消息", c.Id, message.ChannelId))
			continue
		}

		if string(strByte) != "" {
			if msgStore != nil {
				msgStore.SendMessage(message)
			}
			c.hub.broadcast <- message // 转发读取到的channel消息
		}
	}
}
