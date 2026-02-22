package configs

import (
	"fmt"
	"os"
	"strconv"
	"zhihu/app/api/internal/model/Document"
	"zhihu/app/api/internal/model/User"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Db  *gorm.DB
	Cli *redis.Client
)

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

	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Username = user
	}
	// 4. 设置密码
	cfg.Password = password

	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Host = host
	}

	if portStr := os.Getenv("DB_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err != nil {
			cfg.Port = port
		}
	}

	return &cfg, nil
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.DBName)
}

func createChineseFullTextIndex() error {
	if !Db.Migrator().HasIndex(&Document.Article{}, "idx_articles_search") {
		sql := `CREATE FULLTEXT INDEX idx_articles_search 
            ON articles(title, summary) 
            WITH PARSER ngram`
		return Db.Exec(sql).Error
	}
	return nil
}

func InitDB() {
	MyDbConfig, err := LoadDBConfig()
	if err != nil {
		Logger.Fatal("InitDb", zap.Error(err))
	}
	Db, err = gorm.Open(mysql.Open(MyDbConfig.DSN()), &gorm.Config{})
	if err != nil {
		Logger.Fatal("InitDb", zap.Error(err))
	}
	err = Db.AutoMigrate(&User.User{}, &User.LikeUrlUser{}, &User.Relation{},
		&User.FeedItem{}, &User.Message{},
		&Document.Article{}, &Document.Question{}, &Document.Answer{}, &Document.Comment{})
	if err != nil {
		Logger.Fatal("InitDb", zap.Error(err))
	}
	err = createChineseFullTextIndex()
	if err != nil {
		Logger.Fatal("InitDb", zap.Error(err))
	}
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "127.0.0.1"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	Cli = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})
	initBf(Cli)
}
