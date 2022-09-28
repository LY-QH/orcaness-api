package infra

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB Object
var _db map[string]*gorm.DB

func Db(resource ...string) *gorm.DB {
	rs := "default"
	if len(resource) == 1 {
		rs = resource[0]
	}
	if _db == nil {
		_db = make(map[string]*gorm.DB)
	}

	if _db[rs] == nil {
		// 从配置文件中获取参数
		prefix := "database." + rs + "."
		host := viper.GetString(prefix + "host")
		port := viper.GetString(prefix + "port")
		database := viper.GetString(prefix + "name")
		username := viper.GetString(prefix + "username")
		password := viper.GetString(prefix + "password")
		charset := viper.GetString(prefix + "charset")
		loc := viper.GetString(prefix + "loc")
		// 字符串拼接
		sqlStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
			username,
			password,
			host,
			port,
			database,
			charset,
			url.QueryEscape(loc),
		)
		fmt.Println("数据库连接:", sqlStr)
		// 配置日志输出
		newLogger := logger.New(
			log.New(os.Stdout, "\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,   // 缓存日志时间
				LogLevel:                  logger.Silent, // 日志级别
				IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,         // Disable color
			},
		)
		db, err := gorm.Open(mysql.Open(sqlStr), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			fmt.Println("打开数据库失败", err)
			panic("打开数据库失败" + err.Error())
		}

		// debug mode
		mode := os.Getenv("GIN_MODE")
		if mode == "debug" {
			db = db.Debug()
		}

		_db[rs] = db
	}

	return _db[rs]
}
