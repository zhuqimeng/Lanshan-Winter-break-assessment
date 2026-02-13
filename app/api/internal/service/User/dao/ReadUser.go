package dao

import (
	"errors"
	"zhihu/app/api/configs"
	"zhihu/app/api/internal/model/User"
	"zhihu/utils/Strings"

	"go.uber.org/zap"
)

func ReadUser(req *User.CreateUserReq) error {
	var user User.User
	res := configs.Db.Where("name = ?", req.Name).First(&user)
	if res.Error != nil {
		configs.Logger.Error("ReadUser", zap.Error(res.Error))
		return res.Error
	}
	// 查询用户信息

	if Strings.VerifyPassword(user.Password, req.Password) == false {
		return errors.New("密码错误！")
	}
	return nil
}
