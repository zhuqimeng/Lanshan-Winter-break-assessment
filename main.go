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
	if err := configs.InitDB("root:Cyzhu8899312_@tcp(127.0.0.1:3306)/ZhiHu?charset=utf8mb4&parseTime=True&loc=Local"); err != nil {
		configs.Sugar.Errorf("init db err:%v", err)
	}
	configs.Sugar.Info("init success")
	router.Router()
}
