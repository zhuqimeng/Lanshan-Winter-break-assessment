package Message

import (
	"encoding/json"
	"time"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/User"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024
)

func (c *Client) handleMessage(wsMsg User.WSMessage) {
	switch wsMsg.Type {
	case "chat":
		var chatMsg User.ChatMessage
		data, _ := json.Marshal(wsMsg.Data)
		err := json.Unmarshal(data, &chatMsg)
		if err != nil {
			configs.Sugar.Error(err)
			return
		}
		err = c.Hub.BroadcastMsg(c.Username, chatMsg.ReceiverName, chatMsg.Content)
		if err != nil {
			configs.Sugar.Error(err)
			c.Hub.sendToUser(c.Username, User.WSMessage{
				Type: "error",
				Data: User.ErrorMessage{
					Message: err.Error(),
					Code:    "SEND_FAILED",
				},
			})
			return
		}
	case "ping":
		// 心跳，不处理
	}
}

// ReadPump 读取客户端消息
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		if err := c.Conn.Close(); err != nil {
			configs.Sugar.Error(err)
		}
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		configs.Sugar.Error(err)
	}
	c.Conn.SetPongHandler(func(string) error {
		err = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			configs.Sugar.Error(err)
		}
		return nil
	})
	for {
		var wsMsg User.WSMessage
		err = c.Conn.ReadJSON(&wsMsg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				configs.Sugar.Error("WebSocket读取错误", err.Error())
			}
			break
		}
		c.handleMessage(wsMsg)
		// 处理客户端发来的消息
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := c.Conn.Close(); err != nil {
			configs.Sugar.Error(err)
		}
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					configs.Sugar.Error(err)
				}
				return
			}
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				configs.Sugar.Error(err)
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				configs.Sugar.Error(err)
				return
			}
		case <-ticker.C:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				configs.Sugar.Error(err)
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				configs.Sugar.Error(err)
				return
			}
		}
	}
}
