package chart_server

import (
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

type BaseModel struct {
	Id        uint64    `gorm:"column:id;primaryKey;autoIncrement;not null"json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;index"json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"json:"updated_at"`
}

type ImUser struct {
	BaseModel
	Nickname string     `json:"name"validate:"required"`
	ImUsers  []*ImUser  `gorm:"-"json:"im_users"`  // 用户聊天好友
	ImGroups []*ImGroup `gorm:"-"json:"im_groups"` // 用户聊天群
	ImChats  []*ImChat  `gorm:"-"json:"im_chats"`  // 用户聊天室
}

type ImGroup struct {
	BaseModel
	Name    string    `json:"name"`
	ImUsers []*ImUser `gorm:"-"json:"im_users"` // 聊天群全部用户
}

// 聊天室
type ImChat struct {
	BaseModel
	Type   uint8     `json:"type"validate:"required,numeric"` // 聊天室类型 0：单聊，1：群聊
	Code   string    `json:"code"validate:"required"`         // 聊天室code，群聊更加GroupId生成唯一值，单聊根据2个用户id生成唯一值
	Name   string    `json:"name"validate:"required"`
	ImUser []*ImUser `gorm:"-"json:"room_id"validate:"required"`
}

// 聊天记录
type ChatMessage struct {
	BaseModel
	ChannelId  string `json:"channel_id"validate:"required,numeric"`   // 聊天室code
	SendUserId uint64 `json:"send_user_id"validate:"required,numeric"` // 发送用户
	Type       uint8  `json:"type"validate:"required,numeric"`         // 消息类型
	Content    string `json:"content"validate:"required"`              // 消息内容
}
