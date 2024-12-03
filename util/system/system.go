package system

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mingyueyu/myeasygo/util/system/model"
)

func Init(aes string, tea [4]uint32) error{
	return FullInit(aes, tea, "", "")
}

func FullInit(aes string, tea [4]uint32, path string, fileName string) error {
	model.AesInit(aes)
	model.TeaInit(tea)
	if len(path) > 0 && len(fileName) > 0 {
		// 设置配置(正式环境不需要，只要被加密的数据文件)
		model.SetSetting(path, fileName)
		// if err != nil {
			// fmt.Println("设置配置失败：", err)
			// return
		// }
	}
	setting, err := model.ReadDefaultSetting()
	if err != nil {
		fmt.Println("读取默认配置失败：", err)
		return err
	} else {
		fmt.Println("读取配置成功：", model.JsonString(setting))
	}
	return nil
}

func ReturnFail(code int, msg string) gin.H {
	return model.ReturnFail(code, msg)
}

func ReturnSuccess(data interface{}) gin.H {
	return model.ReturnSuccess(data)
}

func JsonString(mapData interface{}) string {
	return model.JsonString(mapData)
}
