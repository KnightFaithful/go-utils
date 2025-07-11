package mysqlclient

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type MySQLConfig struct {
	Host         string `yaml:"host"`
	Port         int64  `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"dbname"`
	Charset      string `yaml:"charset"`
	ParseTime    bool   `yaml:"parse_time"`
	Loc          string `yaml:"loc"`
	MaxIdleConns int64  `yaml:"max_idle_conns"`
	MaxOpenConns int64  `yaml:"max_open_conns"`
}

func InitDB(config *MySQLConfig) *gorm.DB {
	// 构建 DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.Charset,
		config.ParseTime,
		config.Loc,
	)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败: %v", err))
	}

	// 获取通用数据库对象并配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("获取数据库实例失败: %v", err))
	}

	sqlDB.SetMaxIdleConns(int(config.MaxIdleConns))
	sqlDB.SetMaxOpenConns(int(config.MaxOpenConns))
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}
