package demo

import (
	"fmt"

	"github.com/mingyueyu/myeasygo/mredis"
	"github.com/mingyueyu/myeasygo/util"
)

func RedisInit(){
	err := mredis.SetValue("test", "redistest", 0)
	if err != nil {
		fmt.Println("设置redis值错误：",err)
	}
	mredis.SetValue("test2", 122, 0)
	v, err := mredis.GetValue("test")
	if err != nil {
		fmt.Println("获取redis值错误：",err)
	}else {
		fmt.Println(v)
	}
	v, err = mredis.GetValue("test2")
	if err != nil {
		fmt.Println("获取redis值错误：",err)
	}else {
		fmt.Println(util.JsonString(v))
	}
	
}