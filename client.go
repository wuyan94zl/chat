package chart_server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

// 客户端
type Client struct {
	hub *Hub
	// websocket 连接
	conn *websocket.Conn
	// 发送
	send chan []byte
	// 客户端管道ID
	channelId string
	// 客户端ID
	Id uint64
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 转发客户端发送消息
func (c *Client) WriteMsg() {
	defer func() {
		_ = c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send: // 收到channel消息 执行发送
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{}) // 错误 关闭 channel
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)
			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(<-c.send)
			}
			//客户端关闭，退出
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

// 接受到客户端消息
func (c *Client) ReadMsg() {
	defer func() {
		// 注销
		c.hub.unregister <- c
		// 推送注销消息
		_ = c.conn.Close()
	}()
	for {
		// 读取到channel消息
		_, strByte, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message := Message{}
		err = json.Unmarshal(strByte, &message)
		message.SendTime = time.Now().Format("2006-01-02 15:04:05")
		if string(strByte) != "" {
			sendMessage(c.Id, message)
			c.hub.broadcast <- message // 转发读取到的channel消息
		}
	}
}

func RunWs(hub *Hub, w http.ResponseWriter, r *http.Request, channelId string, clientId uint64) {
	//升级get请求为webSocket协议
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("升级get请求错误", err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte), channelId: channelId, Id: clientId}
	//连接时休眠1秒  防止刷新页面 先连接后退出
	time.Sleep(time.Duration(1) * time.Second)
	client.hub.register <- client

	go client.ReadMsg()
	go client.WriteMsg()
}
