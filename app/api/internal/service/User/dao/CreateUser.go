package dao

import (
	"errors"
	"log"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/User"
	"zhihu/utils/randoms"
	"zhihu/utils/strings"
)

func CreateUser(req *User.CreateUserReq) error {
	var count int64
	if err := configs.Db.Model(&User.User{}).Where("name = ?", req.Name).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("已存在的用户名")
	}
	// 检查用户是否存在

	user := User.User{
		Name:     req.Name,
		Password: req.Password,
	}
	salt, err := randoms.GenerateSalt()
	if err != nil {
		return err
	}
	hash := strings.HashPassword(req.Password + salt)
	user.Password = salt + "_" + hash
	// 创建用户信息

	if err := configs.Db.Create(&user).Error; err != nil {
		log.Fatal("用户创建失败：", err)
		return err
	}
	// 写入数据库

	return nil
}
