package chart_server

import (
	"gorm.io/gorm"
	"net/http"
)

var GlobHub *Hub

// Independent operation
//func RunServer(addr string) {
//	GlobHub = NewServer()
//	go GlobHub.Run()
//	mux := http.NewServeMux()
//	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
//		RunWs(GlobHub, w, r,"asd",123)
//	})
//	go http.ListenAndServe(addr, mux)
//}

// Framework-based operation
func Server(w http.ResponseWriter, r *http.Request, db *gorm.DB, channelId string, clientId uint64) {
	DB = db
	RunWs(GlobHub, w, r, channelId, clientId)
}

func ChannelAllUserId(channelId string) []uint64 {
	var ids []uint64
	if list, ok := GlobHub.clients[channelId]; ok {
		for c, _ := range list {
			ids = append(ids, c.Id)
		}
	}
	return ids
}

func AddUserNum() int {
	num := 0
	for _, v := range GlobHub.clients {
		num += len(v)
	}
	return num
}
