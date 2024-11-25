package demo

import (
	"fmt"

	"github.com/mingyueyu/myeasygo/util/mysqlTool"
	"github.com/mingyueyu/myeasygo/util/system"
)

func SettingInit() {
	// 设置加密参数，随意16字符值
	system.AesInit("1234567890123456")
	// 设置加密参数，随意4字节4长度uint32数组
	system.TeaInit([4]uint32{0x01234567, 0x89abcdef, 0xfedcba98, 0x76543210})

	// 设置配置(正式环境不需要，只要被加密的数据文件)
	err := system.SetSetting("/Users/mingyueyu/Desktop/Me/git/myeasygo/demo", "setting.json")
	if err != nil {
		// fmt.Println("设置配置失败：", err)
		return
	}
	// 读取默认配置
	setting, err := system.ReadDefaultSetting()
	if err != nil {
		fmt.Println("读取默认配置失败：", err)
		return
	} else {
		fmt.Println("读取配置成功：", system.JsonString(setting))
		// 更新数据库配置
		mysqls := system.MySqls
		body := []mysqlTool.MySql_t{}
		for i := 0; i < len(mysqls); i++ {
			item_mysql := mysqls[i]
			tables := []mysqlTool.Table_t{}
			for j := 0; j < len(item_mysql.Tables); j++ {
				item_table := item_mysql.Tables[j]
				tables = append(tables, mysqlTool.Table_t{
					Content: item_table.Content,
					Name:    item_table.Name,
				})
			}
			body = append(body, mysqlTool.MySql_t{
				Host:   item_mysql.Host,
				Name:   item_mysql.Name,
				Port:   item_mysql.Port,
				User:   item_mysql.User,
				Pwd:    item_mysql.Pwd,
				Tables: tables,
			})
		}
		mysqlTool.MysqlToolInit(body)
	}
	fmt.Println("自定义数据是：", system.JsonString(setting.Custom.(map[string]interface{})["detail"]))

}
