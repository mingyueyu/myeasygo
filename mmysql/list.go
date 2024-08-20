package mmysql

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingyueyu/myeasygo/util/mysqlTool"
)

func List(r *gin.Engine, relativePath string, dbName string, tableName string, searchTargets []string) {
	ListPlus(r, relativePath, dbName, tableName, searchTargets, nil, nil)
}

func ListPlus(r *gin.Engine, relativePath string, dbName string, tableName string, searchTargets []string, funcParam func(c *gin.Context, param gin.H) (gin.H, int64, error), funcResult func(c *gin.Context, result []gin.H, count int64) ([]gin.H, int64, int64, error)) {
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
		re, count, tcode, err := list(param, dbName, tableName, searchTargets)
		if err != nil {
			if TestType {
				panic(err)
			}
			c.JSON(http.StatusOK, mysqlTool.ReturnFail(tcode, err.Error()))
		} else {
			target := mysqlTool.ReturnSuccess(re)
			target["cont"] = count
			// 处理返回值
			if funcResult != nil {
				tresult, tcount, tcode, err := funcResult(c, re, count)
				if err != nil {
					if TestType {
						panic(err)
					}
					c.JSON(http.StatusOK, mysqlTool.ReturnFail(tcode, err.Error()))
					return
				}
				target = mysqlTool.ReturnSuccess(tresult)
				target["cont"] = tcount

			}
			c.JSON(http.StatusOK, target)
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
			re, count, tcode, err := list(param, dbName, tableName, searchTargets)
			if err != nil {
				if TestType {
					panic(err)
				}
				c.JSON(http.StatusOK, mysqlTool.ReturnFail(tcode, err.Error()))
			} else {
				target := mysqlTool.ReturnSuccess(re)
				target["cont"] = count
				// 处理返回值
				if funcResult != nil {
					tresult, tcount, tcode, err := funcResult(c, re, count)
					if err != nil {
						if TestType {
							panic(err)
						}
						c.JSON(http.StatusOK, mysqlTool.ReturnFail(tcode, err.Error()))
						return
					}
					target = mysqlTool.ReturnSuccess(tresult)
					target["cont"] = tcount

				}
				c.JSON(http.StatusOK, target)
			}
		}
	})
}

func list(param gin.H, dbName string, tableName string, searchTargets []string) ([]gin.H, int64, int64, error) {
	// 处理列表数据
	table := tableNameFromeParam(param, tableName)
	sortString := fmt.Sprintf("%s %s", "createTime", "DESC")
	sort := paramGinH(param["sort"])
	if sort != nil {
		if sort["key"] != nil && len(sort["key"].(string)) > 0 && sort["value"] != nil && len(sort["value"].(string)) > 0 {
			sortString = fmt.Sprintf("%s %s", sort["key"], sort["value"])
		}
	}
	page := paramInt(param["page"], 1) - 1
	limit := paramInt(param["limit"], 20)
	return mysqlTool.ListMysql(dbName, table, whereString(param, searchTargets), sortString, page, limit)
}
