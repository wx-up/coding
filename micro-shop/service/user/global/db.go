package global

import (
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/gorm/schema"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

var initDBOnce sync.Once

func DB() *gorm.DB {
	return db
}

func InitDB() {
	initDBOnce.Do(func() {
		initDB()
	})
}

func initDB() {
	dsn := "root:123456@tcp(localhost:3306)/user_srv?charset=utf-8&parseTime=True&multiStatements=true&loc=Local"
	conf := mysql.New(mysql.Config{
		DSN: dsn,
	})
	res, err := gorm.Open(conf, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: false,
		},
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond, // 慢 sql 的预值
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
			LogLevel:                  logger.Info,
		}),
	})
	if err != nil {
		panic(err)
	}
	db = res
}
