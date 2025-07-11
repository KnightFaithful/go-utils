package gormhelper

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
)

type GormConnectionConfig struct {
	Host     string
	Port     int64
	UserName string
	Password string
	DbName   string
	Options  []string
}

//GetGormConnection
/*
DSN
用户名:密码@tcp(ip:port)/数据库?charset=utf8mb4&parseTime=True&loc=Local
*/
func GetGormConnection(config GormConnectionConfig) (*gorm.DB, error) {
	option := ""
	if len(config.Options) > 0 {
		option = "?" + strings.Join(config.Options, "&")
	}
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v%v", config.UserName, config.Password, config.Host, config.Port, config.DbName, option)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
