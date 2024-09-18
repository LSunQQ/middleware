package config

import (
	"strings"

	"github.com/spf13/viper"
)

var (
	CfgMysql *viper.Viper // MySQL配置
	CfgRedis *viper.Viper // Redis配置
	CfgKafka *viper.Viper // Kafka配置
)

// InitConfig 初始化配置文件
func InitConfig(filepath string) {

	viper.AutomaticEnv()                      // 开启自动从环境变量中读取配置
	replacer := strings.NewReplacer(".", "_") // 环境变量.改_
	viper.SetEnvKeyReplacer(replacer)

	viper.SetConfigFile(filepath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// 读取 MySQL 配置
	CfgMysql = viper.Sub("MySQL")
	if CfgMysql == nil {
		panic("config not found MySQL")
	}

	// 读取 Redis 配置
	CfgRedis = viper.Sub("Redis")
	if CfgRedis == nil {
		panic("config not found Redis")
	}

	// 读取 Kafka 配置
	CfgKafka = viper.Sub("Kafka")
	if CfgKafka == nil {
		panic("config not found Kafka")
	}
}
