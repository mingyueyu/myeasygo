package demo

import (
	"fmt"

	"github.com/mingyueyu/myeasygo/util/system"
	"github.com/mingyueyu/myeasygo/util/system/model"
)

func SettingInit() {
	system.FullInit("1234567890123456", [4]uint32{0x01234567, 0x89abcdef, 0xfedcba98, 0x76543210},"/Users/mingyueyu/Desktop/Me/git/myeasygo/demo", "setting.json")
	fmt.Println("自定义数据是：", system.JsonString(model.Setting.Custom.(map[string]interface{})["detail"]))

}
