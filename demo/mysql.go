package demo

import (
	"github.com/gin-gonic/gin"
	"github.com/mingyueyu/myeasygo/mmysql"
	"github.com/mingyueyu/myeasygo/util/mysqlTool"
)

func MysqlInit(r *gin.Engine){
	// 设置数据库配置
	setupMysql()
	// 数据库相关
	mmysql.Add(r, "/api/add", "test", "test", false, false, false)
	mmysql.Delete(r, "/api/delete", "test", "test")
	mmysql.Update(r, "/api/update", "test", "test")
	mmysql.List(r, "/api/list", "test", "test", []string{"name", "age"})
	mmysql.Detail(r, "/api/detail", "test", "test")
	mmysql.Dif(r, "/api/dif", "test", "test")
}


func setupMysql(){
	mysqlTool.MysqlToolInit([]mysqlTool.MySql_t{
		{
			Host:   "localhost",
			Name:   "test",
			Port:   3306,
			Pwd:    "12345678",
			User:   "root",
			Tables: []mysqlTool.Table_t{
				{
					Name:    "test",
					Content: "name varchar(255), age varchar(255)",
				},
			},
		},
	})
}