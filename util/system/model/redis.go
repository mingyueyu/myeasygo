package model

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var myredis = Redis_t{}
var ctx = context.Background()

var rdb *redis.Client

func RedisInit(host string, port int, pwd string, db int) {
	myredis.Host = host
	myredis.Port = port
	myredis.Pwd = pwd
	myredis.Db = db
}

func ExampleClient() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", myredis.Host, myredis.Port),
		Password: myredis.Pwd, // no password set
		DB:       myredis.Db,  // use default DB
	})

	// val, err := rdb.Get(ctx, "testKey").Result()
	_, err := rdb.Get(ctx, "testKey").Result()
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println(err)
	}
	// fmt.Println("Before testKey:", val)

	err = rdb.Set(ctx, "testKey", "testValue3", 0).Err()
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println(err)
	}

	// val, err = rdb.Get(ctx, "testKey").Result()
	_, err = rdb.Get(ctx, "testKey").Result()
	if err == redis.Nil {
		// fmt.Println("testKey does not exist")
	} else if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Printf("rdb error:%v", err)
	} else {
		// fmt.Println("testKey:", val)
	}
	// Output: key value
	// key2 does not exist
}

func SetRedisValue(key string, value interface{}, expiration time.Duration) bool {
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", myredis.Host, myredis.Port),
			Password: myredis.Pwd, // no password set
			DB:       myredis.Db,  // use default DB
		})
	}
	err := rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println(err)
		return false
	}
	return true
}

func GetRedisValue(key string) (interface{}, error) {
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", myredis.Host, myredis.Port),
			Password: myredis.Pwd, // no password set
			DB:       myredis.Db,  // use default DB
		})
	}
	return rdb.Get(ctx, "testKey").Result()
}
