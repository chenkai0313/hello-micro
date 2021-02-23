package app

import (
	"github.com/go-redis/redis"
	"log"
)

var RedisDB *redis.Client

func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "" + Config.Redis.HostName + ":" + Config.Redis.Port + "",
		Password: "" + Config.Redis.Password + "",
		DB:       Config.Redis.DB,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Panic("redis connect error ", err)
	}
	RedisDB = client
}
