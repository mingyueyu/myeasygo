package demo

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mingyueyu/myeasygo/mmysql"
	"github.com/mingyueyu/myeasygo/util"
)

func MysqlInit(r *gin.Engine) {
	// 数据库相关
	mmysql.Add(r, "/api/add", "testdb", "testTable")
	mmysql.AddPro(r, "/api/addPro", "testdb", "testTable", true, true, true)
	mmysql.AddPlus(r, "/api/addPlus", "testdb", "testTable", true, true, true, func(c *gin.Context, param gin.H) (gin.H, int, error) {
		fmt.Println("addPlus参数：", util.JsonString(param))
		return param, 0, nil
	}, func(c *gin.Context, result gin.H) (gin.H, int, error) {
		fmt.Println("addPlus结果：", util.JsonString(result))
		return result, 0, nil
	})
	mmysql.Sum(r, "/api/sum", "testdb", "testTable")
	
	mmysql.Delete(r, "/api/delete", "testdb", "testTable")
	mmysql.DeletePlus(r, "/api/deletePlus", "testdb", "testTable", func(c *gin.Context, param gin.H) (gin.H, int, error) {
		fmt.Println("deletePlus参数：", util.JsonString(param))
		return param, 0, nil
	},func(c *gin.Context, result int64) (int64, int, error) {
		fmt.Println("deletePlus结果：", result)
		return result, 0, nil
	})
	mmysql.Update(r, "/api/update", "testdb", "testTable")
	mmysql.UpdatePro(r, "/api/updatePro", "testdb", "testTable", true)
	mmysql.UpdatePlus(r, "/api/updatePlus", "testdb", "testTable", true, func(c *gin.Context, param gin.H) (gin.H, int, error) {
		fmt.Println("updatePlus参数：", util.JsonString(param))
		return param, 0, nil
	}, func(c *gin.Context, result gin.H) (gin.H, int, error) {
		fmt.Println("updatePlus结果：", util.JsonString(result))
		return result, 0, nil
	})
	mmysql.List(r, "/api/list", "testdb", "testTable", []string{"name", "age"})
	mmysql.ListPlus(r, "/api/listPlus", "testdb", "testTable", []string{"name", "age"}, func(c *gin.Context, param gin.H) (gin.H, int, error) {
		fmt.Println("listPlus参数：", util.JsonString(param))
		return param, 0, nil
	}, func(c *gin.Context, result []gin.H, count int64) ([]gin.H, int64, int, error) {
		fmt.Println("listPlus结果(", count, "个)：", util.JsonString(result))
		return result, count, 0, nil
	})
	mmysql.Detail(r, "/api/detail", "testdb", "testTable")
	mmysql.DetailPlus(r, "/api/detailPlus", "testdb", "testTable", func(c *gin.Context, param gin.H) (gin.H, int, error) {
		fmt.Println("detailPlus参数：", util.JsonString(param))
		return param, 0, nil
	}, func(c *gin.Context, result gin.H) (gin.H, int, error) {
		fmt.Println("detailPlus结果：", util.JsonString(result))
		return result, 0, nil
	})
	mmysql.Dif(r, "/api/dif", "testdb", "testTable")
	mmysql.DifPlus(r, "/api/difPlus", "testdb", "testTable", func(c *gin.Context, param gin.H) (gin.H, int, error) {
		fmt.Println("difPlus参数：", util.JsonString(param))
		return param, 0, nil
	}, func(c *gin.Context, result []gin.H) ([]gin.H, int, error) {
		fmt.Println("difPlus结果：", util.JsonString(result))
		return result, 0, nil
	})
}
