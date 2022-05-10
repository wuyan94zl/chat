package chart

import (
	"encoding/json"
	"fmt"
	"time"
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

func SendMessageToUid(uid, toUId uint64, msg string, tp uint8) {
	if cli, ok := clients[uid]; ok {
		message := Message{}
		if uid == toUId {
			json.Unmarshal([]byte(msg), &message)
			message.ToUserId = toUId
		} else {
			message = Message{
				Content:  msg,
				UserId:   cli.Id,
				Detail:   cli.Detail,
				ToUserId: toUId,
				Type:     tp,
				SendTime: time.Now().Format("2006-01-02 15:04:05"),
			}
		}
		sendMessage(cli, message)
	}
}

func SendMessageToUids(uid uint64, msg string, tp uint8, toUIds ...uint64) {
	if cli, ok := clients[uid]; ok {
		message := Message{
			Content:  msg,
			UserId:   cli.Id,
			Detail:   cli.Detail,
			Type:     tp,
			SendTime: time.Now().Format("2006-01-02 15:04:05"),
		}
		for _, uid := range toUIds {
			message.ToUserId = uid
			sendMessage(cli, message)
		}
	}
}

func SendMessageToChannelIds(uid uint64, msg string, tp uint8, channelIds ...string) map[string]bool {
	resp := make(map[string]bool)
	if cli, ok := clients[uid]; ok {
		message := Message{
			Content:  msg,
			UserId:   uid,
			Detail:   cli.Detail,
			Type:     tp,
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
	c.hub.broadcast <- msg
	return nil
}
