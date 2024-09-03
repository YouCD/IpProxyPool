package database

import (
	"IpProxyPool/middleware/config"
	"database/sql"
	"fmt"
	"github.com/youcd/toolkit/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	sdkLog "log"
	"net/url"
	"os"
	"sync"
	"time"
)

var dbPingInterval = 90 * time.Second
var (
	db   *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	return db
}

func InitDB(setting *config.Database) *gorm.DB {
	var err error
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", url.QueryEscape(setting.Username), setting.Password, setting.Host, setting.Port) // 连接数据库

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// * 解决中文字符问题：Error 1366
		db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")

		sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s  CHARACTER SET utf8mb4 ", setting.DbName)
		// 创建数据库
		err = db.Exec(sql).Error
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&allowNativePasswords=true&parseTime=True&loc=Local",
			// 连接数据库的用户名
			url.QueryEscape(setting.Username),
			// 连接数据库的密码
			setting.Password,
			// 连接数据库的地址
			setting.Host,
			// 连接数据库的端口号
			setting.Port,
			// 连接数据库的具体数据库名称
			setting.DbName,
		)
		newLogger := logger.New(
			sdkLog.New(os.Stdout, "\r\n", sdkLog.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second,   // 慢 SQL 阈值
				LogLevel:      logger.Silent, // Log level
				Colorful:      false,         // 禁用彩色打印
			},
		)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			//NamingStrategy: schema.NamingStrategy{
			//	TablePrefix:   setting.Prefix, // 表名前缀，`User` 的表名应该是 `t_users`
			//	SingularTable: true,           // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
			//},
			PrepareStmt:            true, // 执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
			DisableAutomaticPing:   false,
			SkipDefaultTransaction: true, // 对于写操作（创建、更新、删除），为了确保数据的完整性，GORM 会将它们封装在事务内运行。但这会降低性能，你可以在初始化时禁用这种方式
			Logger:                 newLogger,
			AllowGlobalUpdate:      false,
		})
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		sqlDb, dbErr := db.DB()
		if dbErr != nil {
			log.Errorf("fail to connect database: %v\n", dbErr)
			os.Exit(-1)
		}
		// 设置连接池
		// 用于设置连接池中空闲连接的最大数量。
		sqlDb.SetMaxIdleConns(10)
		// 设置打开数据库连接的最大数量
		sqlDb.SetMaxOpenConns(100)

		go KeepAlivedDb(sqlDb)

		err = db.AutoMigrate(&IP{})
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})
	return db
}

func KeepAlivedDb(engine *sql.DB) {
	t := time.Tick(dbPingInterval)
	var err error
	for {
		<-t
		err = engine.Ping()
		if err != nil {
			log.Errorf("database ping error: %v\n", err.Error())
		}
	}
}
