package initialize

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"os"
	"server/global"
)

// ConnectRedis 连接到redis
func ConnectRedis() redis.Client {
	redisConfig := global.Config.Redis

	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
	_, err := client.Ping().Result()
	if err != nil {
		global.Log.Error("Failed to connect to redis:", zap.Error(err))
		os.Exit(1)
	}

	return *client
}
