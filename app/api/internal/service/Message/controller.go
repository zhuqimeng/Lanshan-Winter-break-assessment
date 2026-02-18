package Message

import (
	"net/http"
	"zhihu/app/api/configs"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 调试用，生产环境需要限制
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *Hub) HandleWebSocket(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username is empty"})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		configs.Sugar.Error("WebSocket升级失败", err)
		return
	}
	client := &Client{
		Hub:      h,
		Conn:     conn,
		Username: username,
		Send:     make(chan []byte, 256),
	}
	h.register <- client
	go client.WritePump()
	go client.ReadPump()
}
