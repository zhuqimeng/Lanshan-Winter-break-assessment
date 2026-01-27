package configs

import (
	"log"
	"zhihu/app/api/internal/model/User"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB(dsn string) error {
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败：", err)
		return err
	}
	err = Db.AutoMigrate(&User.User{})
	if err != nil {
		log.Fatal("自动迁移失败： ", err)
		return err
	}
	return nil
}
