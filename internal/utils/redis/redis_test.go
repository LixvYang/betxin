package redis

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/lixvyang/betxin/internal/utils"
	"github.com/lixvyang/betxin/pkg/convert"

	"gopkg.in/ini.v1"
)

func TestZADD(t *testing.T) {
	f, err := ini.Load("../../../configs/config.ini")
	if err != nil {
		log.Printf("配置文件读取错误:%s", err)
	}

	utils.LoadRedis(f)

	type Person struct {
		Id    int
		Phone int
	}

	// BatchDel("test")
	p1 := Person{Id: 111, Phone: 1111}
	p2 := Person{Id: 2222, Phone: 3333}
	p3 := Person{Id: 111, Phone: 1444111}
	p4 := Person{Id: 555111, Phone: 1111}
	NewRedisClient(context.Background())
	members := []*redis.Z{
		{
			Score:  100,
			Member: convert.Marshal(p1),
		},
		{
			Score:  200,
			Member: convert.Marshal(p2),
		},
		{
			Score:  100,
			Member: convert.Marshal(p3),
		},
		{
			Score:  50,
			Member: convert.Marshal(p4),
		},
	}
	ZADD("test", members...)
	vv, _ := ZREVRANGE("test", 0, -1)
	for k, v := range vv {
		if vv[k] != v {
			t.Errorf("err Val")
		}
	}
}

// func TestGet(t *testing.T) {
// 	f, err := ini.Load("../../../configs/config.ini")
// 	if err != nil {
// 		log.Printf("配置文件读取错误:%s", err)
// 	}

// 	utils.LoadRedis(f)

// 	NewRedisClient(context.Background())
// 	key, value := "name", "lixv"
// 	Set(key, value, 10*time.Second)
// 	vv := Get(key).Val()
// 	fmt.Println("vv: ", vv)
// 	if vv != value {
// 		t.Errorf("err Val")
// 	}
// }

func TestSADD(t *testing.T) {
	f, err := ini.Load("../../../configs/config.ini")
	if err != nil {
		log.Printf("配置文件读取错误:%s", err)
	}

	utils.LoadRedis(f)
	key := "SADDTEST"
	value := "tid_123132131"
	SADD(key, value)

	vv := SISMEMBER(key, value)
	fmt.Println(vv)

	if !SREM(key, value) {
		t.Error("error ", err)
	}
}
