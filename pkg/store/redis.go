package redis

import (
	config "UserStore/pkg/config"
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/rbcervilla/redisstore/v8"
)

//InitRedisConn return a redis conn
func InitOrDie(conf *config.RedisConfig) *redisstore.RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", conf.Uri, conf.Port),
	})
	// New default RedisStore
	store, err := redisstore.NewRedisStore(context.Background(), client)
	if err != nil {
		log.Fatalf("%s", err)
	}
	return store
}
