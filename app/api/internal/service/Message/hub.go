package Message

import (
	"encoding/json"
	"errors"
	"fmt"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/User"

	"go.uber.org/zap"
)

// 发送消息给指定用户
func (h *Hub) sendToUser(username string, msg User.WSMessage) bool {
	h.mutex.RLock()
	client, ok := h.clients[username]
	h.mutex.RUnlock()
	if !ok {
		return false
	}
	data, err := json.Marshal(msg)
	if err != nil {
		configs.Logger.Error("消息序列化失败", zap.String("Receiver", username), zap.Error(err))
		return false
	}
	select {
	case client.Send <- data:
		return true
	default:
		// 发送通道满，直接关闭连接
		h.mutex.Lock()
		close(client.Send)
		delete(h.clients, client.Username)
		h.mutex.Unlock()
		return false
	}
}

// 推送离线消息
func (h *Hub) pushOfflineMSG(client *Client) {
	var messages []User.Message
	err := configs.Db.Where("receiver_name = ? AND delivered = ?", client.Username, false).Order("created_at ASC").Find(&messages).Error
	if err != nil {
		configs.Logger.Error("获取离线消息失败", zap.String("Receiver", client.Username), zap.Error(err))
		return
	}
	if len(messages) == 0 {
		return
	}
	configs.Logger.Info("开始推送离线消息", zap.String("Receiver", client.Username), zap.Int("Count", len(messages)))
	// 转换为聊天消息格式
	for _, msg := range messages {
		chatMSG := User.ChatMessage{
			ID:           msg.ID,
			SenderName:   msg.SenderName,
			ReceiverName: msg.ReceiverName,
			Content:      msg.Content,
			CreatedAt:    msg.CreatedAt,
		}
		success := h.sendToUser(client.Username, User.WSMessage{
			Type: "chat",
			Data: chatMSG,
		})
		if success {
			// 标记为已发送
			configs.Db.Model(User.Message{}).Where("id = ?", msg.ID).Update("delivered", true)
		}
	}
	configs.Logger.Info("推送离线消息成功", zap.String("Receiver", client.Username))
}

func (h *Hub) handleRegister(client *Client) {
	h.mutex.Lock()
	// 如果用户已经在线，先踢掉旧的连接
	if oldClient, ok := h.clients[client.Username]; ok {
		close(oldClient.Send)
		delete(h.clients, client.Username)
	}
	h.clients[client.Username] = client
	h.mutex.Unlock()
	configs.Logger.Info("用户上线", zap.String("username", client.Username))
	go h.pushOfflineMSG(client)
}

func (h *Hub) handleUnregister(client *Client) {
	h.mutex.Lock()
	if _, ok := h.clients[client.Username]; ok {
		delete(h.clients, client.Username)
		close(client.Send)
	}
	h.mutex.Unlock()
	configs.Logger.Info("用户下线", zap.String("username", client.Username))
}

func (h *Hub) handleMessage(message *User.Message) {
	// 1.保存到数据库
	if err := configs.Db.Create(message).Error; err != nil {
		configs.Logger.Error("保存消息失败", zap.String("Receiver", message.ReceiverName), zap.Error(err))
		return
	}
	// 2.转化消息格式
	chatMSG := User.ChatMessage{
		ID:           message.ID,
		SenderName:   message.SenderName,
		ReceiverName: message.ReceiverName,
		Content:      message.Content,
		CreatedAt:    message.CreatedAt,
	}
	// 3.根据不同情况发送给接收者
	h.mutex.RLock()
	_, online := h.clients[message.ReceiverName]
	h.mutex.RUnlock()
	if online {
		success := h.sendToUser(message.ReceiverName, User.WSMessage{
			Type: "chat",
			Data: chatMSG,
		})
		if success {
			// 4. 推送成功，标记为已送达
			configs.Db.Model(&User.Message{}).Where("id = ?", message.ID).Update("delivered", true)
		}
		configs.Sugar.Info(fmt.Sprintf("实时消息： %s -> %s", message.SenderName, message.ReceiverName))
	} else {
		configs.Sugar.Info(fmt.Sprintf("离线消息（已存库）： %s -> %s", message.SenderName, message.ReceiverName))
	}
}

// Run 持续运行 Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.handleRegister(client)
		case client := <-h.unregister:
			h.handleUnregister(client)
		case message := <-h.broadcast:
			h.handleMessage(message)
		}
	}
}

// BroadcastMsg 广播消息
func (h *Hub) BroadcastMsg(senderName, receiverName, content string) error {
	if senderName == receiverName {
		return errors.New("不要自言自语了，快去找人聊天吧！")
	}
	// 检查关注关系
	var count int64
	err := configs.Db.Model(&User.Relation{}).Where("(username = ? AND follower = ?) OR (follower = ? AND username = ?)", senderName, receiverName, senderName, receiverName).Count(&count).Error
	if err != nil {
		return err
	}
	if count != 2 {
		return errors.New("只有相互关注的用户才能发送消息")
	}
	// 创建消息
	msg := &User.Message{
		SenderName:   senderName,
		ReceiverName: receiverName,
		Content:      content,
	}
	// 广播消息
	select {
	case h.broadcast <- msg:
		return nil
	default:
		return errors.New("系统繁忙，请稍后重试")
	}
}
