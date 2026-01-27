package dao

import (
	"errors"
	"log"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/User"
	"zhihu/utils/strings"
)

func ReadUser(req *User.CreateUserReq) error {
	var user User.User
	res := configs.Db.Where("name = ?", req.Name).First(&user)
	if res.Error != nil {
		log.Fatal("找不到用户：", res.Error)
		return res.Error
	}
	// 查询用户信息

	if strings.VerifyPassword(user.Password, req.Password) == false {
		return errors.New("密码错误！")
	}
	return nil
}
