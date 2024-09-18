package main

import (
	"fmt"
	"log"
	"middleware/config"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	defaultDB *gorm.DB
)

type MySQL struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	CharSet         string        `mapstructure:"charset"`
	Loc             string        `mapstructure:"loc"`
	MaxOpenConns    int           `mapstructure:"maxopenconns"`
	MaxIdleConns    int           `mapstructure:"maxidleconns"`
	ConnMaxLifeTime time.Duration `mapstructure:"connmaxlifetime"`
}

func InitMySQL() {
	var mysqlCfg MySQL
	config.CfgMysql.Unmarshal(&mysqlCfg)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=%s",
		mysqlCfg.User,
		mysqlCfg.Password,
		mysqlCfg.Host,
		mysqlCfg.Port,
		mysqlCfg.Database,
		mysqlCfg.CharSet,
		url.QueryEscape(mysqlCfg.Loc),
	)

	// 连接mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Print(err.Error())
	}

	// 自动迁移表
	// if !db.Migrator().HasTable(&TableName{}) {
	// 	err = db.Set("gorm:table_options", "COMMENT='用户信息表'").AutoMigrate(&TableName{})
	// 	if err != nil {
	// 		log.Print(err.Error())
	// 	}
	// }

	// 设置最大连接数
	sqlDB, err := db.DB()
	if err != nil {
		log.Print(err.Error())
	}
	// 设置最大连接数
	sqlDB.SetMaxOpenConns(mysqlCfg.MaxOpenConns)
	// 设置空闲连接池最大连接数
	sqlDB.SetMaxIdleConns(mysqlCfg.MaxIdleConns)
	// 设置连接的最大生命周期
	sqlDB.SetConnMaxLifetime(mysqlCfg.ConnMaxLifeTime * time.Second)

	defaultDB = db
	if CheckDBConnection(&gin.Context{}) {
		log.Print("mysql connection success")
	} else {
		log.Print("mysql connection failed")
	}
}

func GetMySQLClient() *gorm.DB {
	if defaultDB == nil {
		InitMySQL()
	}
	return defaultDB
}

// 检查 MySQL 连接是否正常
func CheckDBConnection(ctx *gin.Context) bool {
	sqlDB, err := defaultDB.DB()
	if err != nil {
		return false
	}
	err = sqlDB.Ping()
	return err == nil
}
