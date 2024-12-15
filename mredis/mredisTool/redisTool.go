package mredisTool

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis_t struct {
	Host     string // 地址
	Port     int    // 端口
	Password string // 密码
	Db       int    // 数据库
}

var TestType = false
var MyRedis = Redis_t{}
var ctx = context.Background()

var rdb *redis.Client

func SetValue(key string, value interface{}, expiration time.Duration) error {
	if MyRedis.Host == "" || MyRedis.Password == "" {
		return fmt.Errorf("redis not init")
	}
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", MyRedis.Host, MyRedis.Port),
			Password: MyRedis.Password, // no password set
			DB:       MyRedis.Db,  // use default DB
		})
	}
	err := rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println(err)
		return err
	}
	return nil
}

func GetValue(key string) (interface{}, error) {
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", MyRedis.Host, MyRedis.Port),
			Password: MyRedis.Password, // no password set
			DB:       MyRedis.Db,  // use default DB
		})
	}
	return rdb.Get(ctx, "testKey").Result()
}
