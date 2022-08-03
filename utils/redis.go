package utils

import (
	"GoStatusServer/config"
	"github.com/go-redis/redis"
)

var Redisdb *redis.Client

func RedisInit() {
	Redisdb = redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisServer.Address,
		Password: config.Config.RedisServer.Password,
		DB:       config.Config.RedisServer.DB,
	})
}
