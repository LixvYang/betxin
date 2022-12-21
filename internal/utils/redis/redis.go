package redis

import (
	"context"
	"time"

	"github.com/lixvyang/betxin/internal/utils"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	ctx         context.Context
	redisClient *redis.Client
}

var r *RedisClient

func NewRedisClient(ctx context.Context) {
	r = &RedisClient{
		ctx: ctx,
		redisClient: redis.NewClient(&redis.Options{
			Addr:     utils.RedisHost + ":" + utils.RedisPort,
			Password: utils.RedisPassword, // no password set
			DB:       0,                   // use default DB
		},
		),
	}
}

func Get(key string) *redis.StringCmd {
	return r.redisClient.Get(r.ctx, key)
}
func Exists(key string) bool {
	return r.redisClient.Exists(r.ctx, key).Val() != 0
}

func Set(key string, value interface{}, expiration time.Duration) {
	if Exists(key) {
		r.redisClient.Expire(r.ctx, key, expiration)
		return
	}
	r.redisClient.Set(r.ctx, key, value, expiration)
}

func DelKeys(keys ...string) {
	for i := 0; i < len(keys); i++ {
		Del(keys[i])
	}
}

func Del(key string) {
	if !Exists(key) {
		return
	}
	r.redisClient.Del(r.ctx, key)
}

// 批量删除
func BatchDel(key string) {
	iter := r.redisClient.Scan(r.ctx, 0, key+"*", 0).Iterator()
	for iter.Next(r.ctx) {
		Del(iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
}

func ZADD(key string, members ...*redis.Z) {
	r.redisClient.ZAdd(r.ctx, key, members...)
}

func ZRANGE(key string) []string {
	return r.redisClient.ZRange(r.ctx, key, 0, -1).Val()
}
