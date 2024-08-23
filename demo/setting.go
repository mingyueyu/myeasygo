package demo

import (
	"fmt"

	"github.com/mingyueyu/myeasygo/util/system"
)

func SettingInit() {
	system.AesInit("1234567890123456")
	system.TeaInit([4]uint32{0x01234567, 0x89abcdef, 0xfedcba98, 0x76543210})

	err := system.SetSetting("/Users/mingyueyu/Desktop/Me/git/myeasygo/demo", "setting.json")
	if err != nil {
		fmt.Println("设置配置失败：", err)
		return
	}
	// 读取默认配置
	setting, err := system.ReadDefaultSetting()
	if err != nil {
		fmt.Println("读取默认配置失败：", err)
		return
	}
	fmt.Println("读取配置成功：", system.JsonString(setting))
}
