package User

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	SenderName   string `json:"sender_name" gorm:"not null"`
	ReceiverName string `json:"receiver_name" gorm:"not null"`
	Content      string `json:"content" gorm:"type:text;not null"`
}

// WSMessage WebSocket消息格式
type WSMessage struct {
	Type string      `json:"type"` // 消息类型: chat/offline/history/error
	Data interface{} `json:"data"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	ID           uint      `json:"id,omitempty"`
	SenderName   string    `json:"sender_name"`
	ReceiverName string    `json:"receiver_name"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	IsSelf       bool      `json:"is_self,omitempty"` // 是否是自己的消息（前端用）
}

// ErrorMessage 错误消息
type ErrorMessage struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
