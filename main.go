package main

import (
	"fmt"
	"log"
	"zhihu/app/api/configs"
	"zhihu/app/api/router"
)

func main() {
	fmt.Println("welcome to use my project")
	if err := configs.InitDB("root:Cyzhu8899312_@tcp(127.0.0.1:3306)/zhihu?charset=utf8mb4&parseTime=True&loc=Local"); err != nil {
		log.Fatal(err)
	}
	router.Router()
}
