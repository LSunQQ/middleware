package main

import (
	"fmt"
	"log"
	"middleware/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	defaultRedis *redis.Client
)

// Redis 绑定数据库配置
type Redis struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Database       int    `mapstructure:"database"`
	MaxIdleConns   int    `mapstructure:"maxidleconns"`
	MaxActiveConns int    `mapstructure:"maxactiveconns"`
}

// InitRedis 初始化redis
func InitRedis() {
	var redisCfg Redis
	config.CfgRedis.Unmarshal(&redisCfg)

	rc := redis.NewClient(&redis.Options{
		Addr:           fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Username:       redisCfg.User,
		Password:       redisCfg.Password,
		DB:             redisCfg.Database,
		MaxIdleConns:   redisCfg.MaxIdleConns,   // 最大空闲连接数
		MaxActiveConns: redisCfg.MaxActiveConns, // 最大活跃连接数
	})

	defaultRedis = rc

	if CheckRedisConnection(&gin.Context{}) {
		log.Print("redis connection success")
	} else {
		log.Print("redis connection failed")
	}
}

// GetRedisClient 获取redis客户端
func GetRedisClient() *redis.Client {
	if defaultRedis == nil {
		InitRedis()
	}
	return defaultRedis
}

// GetFromRedis 从redis中获取数据
func GetFromRedis(ctx *gin.Context, key string) (string, error) {
	data, err := defaultRedis.Get(ctx, key).Result()
	if err == redis.Nil {
		// 缓存不存在
		return "", nil
	} else if err != nil {
		// 其他错误
		return "", fmt.Errorf("failed to get data from redis: %v", err)
	}
	return data, nil
}

// SetToRedis 将数据写入redis
func SetToRedis(ctx *gin.Context, key string, value string, expiration time.Duration) error {
	return defaultRedis.Set(ctx, key, value, expiration).Err()
}

// 检查 Redis 连接是否正常
func CheckRedisConnection(ctx *gin.Context) bool {
	_, err := defaultRedis.Ping(ctx).Result()
	return err == nil
}
