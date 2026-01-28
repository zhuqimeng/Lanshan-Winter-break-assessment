package configs

import (
	"zhihu/app/api/internal/model/User"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB(dsn string) error {
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		Logger.Fatal("InitDb", zap.Error(err))
		return err
	}
	err = Db.AutoMigrate(&User.User{})
	if err != nil {
		Logger.Fatal("InitDb", zap.Error(err))
		return err
	}
	return nil
}
