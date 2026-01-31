package configs

import (
	"fmt"
	"os"
	"zhihu/app/api/internal/model/User"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string // 从环境变量获取，不写进配置文件
	DBName   string `mapstructure:"dbname"`
}

func LoadDBConfig() (*DBConfig, error) {
	// 1. 读取基本配置
	viper.SetConfigFile("app/api/configs/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 2. 从环境变量获取密码（更安全）
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		// 或者从 viper 读取
		password = viper.GetString("database.password")
		if password == "" {
			return nil, fmt.Errorf("数据库密码未设置，请设置 DB_PASSWORD 环境变量")
		}
	}

	// 3. 解析配置
	var cfg DBConfig
	if err := viper.UnmarshalKey("database", &cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 4. 设置密码
	cfg.Password = password

	return &cfg, nil
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.DBName)
}

func InitDB() error {
	MyDbConfig, err := LoadDBConfig()
	if err != nil {
		Logger.Fatal("InitDb", zap.Error(err))
		return err
	}
	Db, err = gorm.Open(mysql.Open(MyDbConfig.DSN()), &gorm.Config{})
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
