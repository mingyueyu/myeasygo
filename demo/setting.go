package demo

import "github.com/mingyueyu/myeasygo/system"

func SettingInit() {
	system.FullInit("1234567890123456", [4]uint32{0x01234567, 0x89abcdef, 0xfedcba98, 0x76543210},"/Users/mingyueyu/Desktop/Me/git/myeasygo/demo/settingfile","")
}
