package mredis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mingyueyu/myeasygo/mredis/mredisTool"
	"github.com/mingyueyu/myeasygo/util"
)

func RefreshSetting(data []byte) {
	err := json.Unmarshal([]byte(data), &mredisTool.MyRedis)
	if err != nil {
		fmt.Println("更新redis配置失败:", err)
	}else {
		fmt.Println("更新redis配置成功:", util.JsonString(mredisTool.MyRedis))
	}
}

func SetValue(key string, value interface{}, expiration time.Duration) error {
	return mredisTool.SetValue(key, value, expiration)
}

func GetValue(key string) (interface{}, error) {
	return mredisTool.GetValue(key)
}

