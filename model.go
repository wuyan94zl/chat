package chart

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

// 聊天记录
type ChatMessage struct {
	BaseModel
	ChannelId  string `json:"channel_id"validate:"required,numeric"`   // 聊天室code
	SendUserId uint64 `json:"send_user_id"validate:"required,numeric"` // 发送用户
	Type       uint8  `json:"type"validate:"required,numeric"`         // 消息类型
	Content    string `json:"content"validate:"required"`              // 消息内容
}
