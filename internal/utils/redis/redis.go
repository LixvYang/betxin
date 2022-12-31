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

func ZCARD(key string) int {
	return int(r.redisClient.ZCard(r.ctx, key).Val())
}

// 按分数正序返回
func ZRANGE(key string, offset, limit int) ([]string, error) {
	return r.redisClient.ZRange(r.ctx, key, int64(offset), int64(limit-1)).Result()
}

// 按分数倒序返回
func ZREVRANGE(key string, offset, limit int) ([]string, error) {
	return r.redisClient.ZRevRange(r.ctx, key, int64(offset), int64(limit-1)).Result()
}

// 点赞Incr
func Incr(key string) int64 {
	return r.redisClient.Incr(r.ctx, key).Val()
}

func SADD(key string, members ...any) {
	r.redisClient.SAdd(r.ctx, key, members...)
}

func SISMEMBER(key string, member any) bool {
	return r.redisClient.SIsMember(r.ctx, key, member).Val()
}

func SREM(key string, member any) bool {
	return r.redisClient.SRem(r.ctx, key, member).Val() == 1
}
