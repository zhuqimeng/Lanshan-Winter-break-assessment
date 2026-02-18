package Message

import (
	"sync"
	"zhihu/app/api/internal/model/User"

	"github.com/gorilla/websocket"
)

// Client 客户端连接
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	Username string
}

// Hub 连接中心
type Hub struct {
	clients    map[string]*Client // key: username
	broadcast  chan *User.Message // 消息广播通道
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan *User.Message, 100),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}
