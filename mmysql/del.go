package mmysql

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingyueyu/myeasygo/util/mysqlTool"
)

func Delete(r *gin.Engine, relativePath string, dbName string, tableName string) {
	DeletePlus(r, relativePath, dbName, tableName, nil, nil)
}

func DeletePlus(r *gin.Engine, relativePath string, dbName string, tableName string, funcParam func(c *gin.Context, param gin.H) (gin.H, int, error), funcResult func(c *gin.Context, result int64) (int64, int, error)) {
	r.POST(relativePath, func(c *gin.Context) {
		param, err := ParamToGinH(c)
		if err != nil {
			if TestType {
				panic(err)
			}
			c.JSON(http.StatusOK, mysqlTool.ReturnFail(10001, err.Error()))
			return
		}
		// 处理参数
		if funcParam != nil {
			tparam, tcode, err := funcParam(c, param)
			if err != nil {
				if TestType {
					panic(err)
				}
				c.JSON(http.StatusOK, mysqlTool.ReturnFail(tcode, err.Error()))
				return
			}
			param = tparam
		}
		re, tcode, err := MysqlDel(param, dbName, tableName)
		if err != nil {
			if TestType {
				panic(err)
			}
			c.JSON(http.StatusOK, mysqlTool.ReturnFail(tcode, err.Error()))
		} else {
			// 处理返回值
			if funcResult != nil {
				tresult, tcode, err := funcResult(c, re)
				if err != nil {
					if TestType {
						panic(err)
					}
					c.JSON(http.StatusOK, mysqlTool.ReturnFail(tcode, err.Error()))
					return
				}
				re = tresult
			}
			c.JSON(http.StatusOK, mysqlTool.ReturnSuccess(re))
		}
	})
}

func MysqlDel(param gin.H, dbName string, tableName string) (int64, int, error) {
	table := tableNameFromeParam(param, tableName)
	where, whereValues := whereString(param, nil)
	count, tcode, err := mysqlTool.DelectMysql(dbName, table, where, whereValues)
	// if tcode == 10010 {
	// 	dealwithMysql()
	// 	where, whereValues := whereString(param, nil)
	// 	return mysqlTool.DelectMysql(dbName, table, where, whereValues)
	// }
	return count, tcode, err
}
