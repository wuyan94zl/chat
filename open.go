package chart

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func JoinChannelIds(uid uint64, channelIds ...string) error {
	if cli, ok := clients[uid]; ok {
		for _, channelId := range channelIds {
			if cli.hub.clients[channelId] == nil {
				cli.hub.clients[channelId] = map[*Client]bool{}
			}
			cli.hub.clients[channelId][cli] = true
		}
		cli.channelId = channelIds
		return nil
	} else {
		msgStore.ErrorLogServer(fmt.Errorf("用户`%d`未登录", uid))
		return fmt.Errorf("用户`%d`未登录", uid)
	}
}

func UnJoinChannelIds(uid uint64, channelIds ...string) error {
	if cli, ok := clients[uid]; ok {
		delChannel := make(map[string]bool)
		for _, channelId := range channelIds {
			if cli.hub.clients[channelId] == nil {
				return fmt.Errorf("用户`%d`为监听：`%s` channel", uid, channelId)
			}

			if _, ok := cli.hub.clients[channelId][cli]; ok {
				delete(cli.hub.clients[channelId], cli)
				delChannel[channelId] = true
			}
		}
		var hasChannelIds []string
		for _, v := range cli.channelId {
			if _, ok := delChannel[v]; !ok {
				hasChannelIds = append(hasChannelIds, v)
			}
		}
		cli.channelId = hasChannelIds
		return nil
	} else {
		msgStore.ErrorLogServer(fmt.Errorf("用户`%d`未登录", uid))
		return fmt.Errorf("用户`%d`未登录", uid)
	}
}

func SendMessageToUid(uid, toUId uint64, msg string) {
	if cli, ok := clients[uid]; ok {
		message := Message{
			Content:  msg,
			UserId:   uid,
			ToUserId: toUId,
			Type:     msgTypeSendText,
			SendTime: time.Now().Format("2006-01-02 15:04:05"),
		}
		sendMessage(cli, message)
	}
}

func SendMessageToChannelIds(uid uint64, msg string, channelIds ...string) map[string]bool {
	resp := make(map[string]bool)
	if cli, ok := clients[uid]; ok {
		message := Message{
			Content:  msg,
			UserId:   uid,
			Type:     msgTypeSendText,
			SendTime: time.Now().Format("2006-01-02 15:04:05"),
		}

		for _, channelId := range channelIds {
			if !cli.hub.clients[channelId][cli] {
				resp[channelId] = false
				continue
			}
			message.ChannelId = channelId
			if sendMessage(cli, message) == nil {
				resp[channelId] = true
			} else {
				resp[channelId] = false
			}
		}
	}
	return resp

}

func sendMessage(c *Client, msg Message) error {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		msgStore.ErrorLogServer(fmt.Errorf("发送消息错误：%v", err))
		return err
	}
	byteMessage, err := json.Marshal(msg)
	if err != nil {
		msgStore.ErrorLogServer(fmt.Errorf("发送消息错误：%v", err))
		return err
	}
	_, _ = w.Write(byteMessage)
	n := len(c.send)
	for i := 0; i < n; i++ {
		_, _ = w.Write(<-c.send)
	}
	if err := w.Close(); err != nil {
		msgStore.ErrorLogServer(fmt.Errorf("发送消息错误：%v", err))
		return err
	}
	return nil
}
