package router

import "zhihu/app/api/internal/service/Message"

func InitHub() *Message.Hub {
	hub := Message.NewHub()
	go hub.Run()
	return hub
}
