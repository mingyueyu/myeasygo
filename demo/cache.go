package demo

import (
	"fmt"
	"time"

	"github.com/mingyueyu/myeasygo/util/cache"
)

func CacheInit(){
	// 设置缓存
	cache.Set("test", "testValue")
	// 获取缓存
	v, is := cache.Get("test")
	if is {
		fmt.Println("获取缓存值为：", v)
	} else {
		fmt.Println("没有缓存值")
	}

	// 设置缓存，设定有效期秒
	cache.SetDuration("testDuration", "testDurationValue", 2)
	vDuration, isDuration := cache.Get("testDuration")
	if isDuration {
		fmt.Println("获取缓存值为：", vDuration)
	} else {
		fmt.Println("没有缓存值")
	}
	
	// 创建一个在 2 秒后触发的定时器
    timer := time.After(time.Duration(1) * time.Second)
    // 等待定时器触发
    <-timer
	vDuration, isDuration = cache.Get("testDuration")
	if isDuration {
		fmt.Println("获取缓存值为：", vDuration)
	} else {
		fmt.Println("没有缓存值")
	}
	// 创建一个在 2 秒后触发的定时器
    timer = time.After(time.Duration(2) * time.Second)
    // 等待定时器触发
    <-timer

	vDuration, isDuration = cache.Get("testDuration")
	if isDuration {
		fmt.Println("获取缓存值为：", vDuration)
	} else {
		fmt.Println("没有缓存值")
	}

	cache.Cleanup()

	cache.Delete("test")

	v, is = cache.Get("test")
	if is {
		fmt.Println("获取缓存值为：", v)
	} else {
		fmt.Println("没有缓存值")
	}
}