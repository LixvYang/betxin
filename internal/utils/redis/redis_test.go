package redis

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/lixvyang/betxin/internal/utils"

	"github.com/go-redis/redis/v8"
	"gopkg.in/ini.v1"
)

func TestZADD(t *testing.T) {
	f, err := ini.Load("../../../configs/config.ini")
	if err != nil {
		log.Printf("配置文件读取错误:%s", err)
	}

	utils.LoadRedis(f)

	NewRedisClient(context.Background())
	members := []*redis.Z{
		{
			Score:  100,
			Member: "EEWDWDW",
		},
		{
			Score:  200,
			Member: "nihao",
		},
		{
			Score:  100,
			Member: "test",
		},
	}
	ZADD("test", members...)
	vv := ZRANGE("test")
	for k, v := range vv {
		if vv[k] != v {
			t.Errorf("err Val")
		}
	}
}

func TestGet(t *testing.T) {
	f, err := ini.Load("../../../configs/config.ini")
	if err != nil {
		log.Printf("配置文件读取错误:%s", err)
	}

	utils.LoadRedis(f)

	NewRedisClient(context.Background())
	key, value := "name", "lixv"
	Set(key, value, 10*time.Second)
	vv := Get(key).Val()
	fmt.Println("vv: ", vv)
	if vv != value {
		t.Errorf("err Val")
	}
}
