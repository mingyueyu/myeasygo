package mmysql

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingyueyu/myeasygo/util/mysqlTool"
)

func Detail(r *gin.Engine, relativePath string, dbName string, tableName string) {
	DetailPlus(r, relativePath, dbName, tableName, nil, nil)
}

func DetailPlus(r *gin.Engine, relativePath string, dbName string, tableName string, funcParam func(c *gin.Context, param gin.H) (gin.H, int64, error), funcResult func(c *gin.Context, result gin.H) (gin.H, int64, error)) {
	r.GET(relativePath, func(c *gin.Context) {
		param := paramFromGet(c)
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
		re, tcode, err := detail(param, dbName, tableName)
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
	r.POST(relativePath, func(c *gin.Context) {
		param, err := paramToGinH(c)
		if err != nil {
			if TestType {
				panic(err)
			}
			c.JSON(http.StatusOK, mysqlTool.ReturnFail(10001, err.Error()))
		} else {
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
			re, tcode, err := detail(param, dbName, tableName)
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
		}
	})
}

func detail(param gin.H, dbName string, tableName string) (gin.H, int64, error) {
	table := tableNameFromeParam(param, tableName)
	return mysqlTool.DetailMysql(dbName, table, whereString(param, nil))
}
