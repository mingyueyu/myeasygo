package demo

import (
	"github.com/gin-gonic/gin"
	"github.com/mingyueyu/myeasygo/mmysql"
)

func MysqlInit(r *gin.Engine){
	// 数据库相关
	mmysql.Add(r, "/api/add", "testDB", "testTable")
	mmysql.Delete(r, "/api/delete", "testDB", "testTable")
	mmysql.Update(r, "/api/update", "testDB", "testTable")
	mmysql.List(r, "/api/list", "testDB", "testTable", []string{"name", "age"})
	mmysql.Detail(r, "/api/detail", "testDB", "testTable")
	mmysql.Dif(r, "/api/dif", "testDB", "testTable")
}