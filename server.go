package chart

import (
	"fmt"
	"net/http"
	"time"
)

func NewServer(w http.ResponseWriter, r *http.Request, clientId uint64, detail interface{}, store messageInterface) {
	if globHub == nil {
		globHub = newServer()
		go globHub.run()
	}
	if clients == nil {
		clients = make(map[uint64]*Client)
	}
	if msgStore == nil {
		msgStore = store
	}
	runWs(globHub, w, r, clientId, detail)
}

func runWs(hub *Hub, w http.ResponseWriter, r *http.Request, clientId uint64, detail interface{}) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("升级get请求错误", err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte), Id: clientId, Detail: detail}
	//连接时休眠1秒  防止刷新页面 先连接后退出
	time.Sleep(time.Duration(1) * time.Second)
	client.hub.register <- client
	go client.ReadMsg()
	go client.WriteMsg()
}
