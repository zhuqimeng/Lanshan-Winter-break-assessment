package main

import (
	"fmt"
	"zhihu/app/api/configs"
	"zhihu/app/api/router"
)

func main() {
	fmt.Println("Welcome to use my project ZhiHu.")
	configs.InitLogger()
	defer func() {
		if err := configs.Logger.Sync(); err != nil {
			fmt.Println(err)
		}
	}()
	configs.InitDB()
	configs.InitLlm()
	hub := router.InitHub()
	configs.Sugar.Info("init success")
	router.Router(hub)
}
