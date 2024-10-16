package demo

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mingyueyu/myeasygo/mmysql"
	"github.com/mingyueyu/myeasygo/util/system"
)

func MysqlInit(r *gin.Engine) {
	// 数据库相关
	mmysql.Add(r, "/api/add", "testdb", "testTable")
	mmysql.AddPlus(r, "/api/addPlus", "testdb", "testTable", true, true, true, func(c *gin.Context, param gin.H) (gin.H, int, error) {
		fmt.Println("addPlus参数：", system.JsonString(param))
		return param, 0, nil
	}, func(c *gin.Context, result gin.H) (gin.H, int, error) {
		fmt.Println("addPlus结果：", system.JsonString(result))
		return result, 0, nil
	})
	mmysql.Delete(r, "/api/delete", "testdb", "testTable")
	mmysql.Update(r, "/api/update", "testdb", "testTable")
	mmysql.List(r, "/api/list", "testdb", "testTable", []string{"name", "age"})
	mmysql.Detail(r, "/api/detail", "testdb", "testTable")
	mmysql.Dif(r, "/api/dif", "testdb", "testTable")
}
