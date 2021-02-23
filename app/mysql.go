package app

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var MysqlDb *gorm.DB

//连接数据库
func InitMysql() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // 禁用彩色打印
		},
	)
	mysqlDb, err := gorm.Open(mysql.Open(Config.MysqlDbDns), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Panic("mysql connect error ", err)
	}
	MysqlDb = mysqlDb
}
