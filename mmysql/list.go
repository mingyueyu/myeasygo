package mmysql

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mingyueyu/myeasygo/util/mysqlTool"
)

func List(r *gin.Engine, relativePath string, dbName string, tableName string, searchTargets []string) {
	ListPlus(r, relativePath, dbName, tableName, searchTargets, nil, nil)
}

func ListPlus(r *gin.Engine, relativePath string, dbName string, tableName string, searchTargets []string, funcParam func(c *gin.Context, param gin.H) (gin.H, int, error), funcResult func(c *gin.Context, result []gin.H, count int64) ([]gin.H, int64, int, error)) {
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
		re, count, tcode, err := MysqlList(param, dbName, tableName, searchTargets)
		if err != nil {
			if TestType {
				panic(err)
			}
			c.JSON(http.StatusOK, mysqlTool.ReturnFail(tcode, err.Error()))
		} else {
			target := mysqlTool.ReturnSuccess(re)
			target["count"] = count
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
				target["count"] = tcount

			}
			c.JSON(http.StatusOK, target)
		}
	})
	r.POST(relativePath, func(c *gin.Context) {
		param, err := ParamToGinH(c)
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
			re, count, tcode, err := MysqlList(param, dbName, tableName, searchTargets)
			if err != nil {
				if TestType {
					panic(err)
				}
				c.JSON(http.StatusOK, mysqlTool.ReturnFail(tcode, err.Error()))
			} else {
				target := mysqlTool.ReturnSuccess(re)
				target["count"] = count
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
					target["count"] = tcount

				}
				c.JSON(http.StatusOK, target)
			}
		}
	})
}

func MysqlList(param gin.H, dbName string, tableName string, searchTargets []string) ([]gin.H, int64, int, error) {
	// 处理列表数据
	table := tableNameFromeParam(param, tableName)
	sortString := fmt.Sprintf("%s %s", "createTime", "DESC")
	if param["sort"] != nil {
		array := param["sort"].([]gin.H)
		sorts := []string{}
		for i := 0; i < len(array); i++ {
			sort := array[i]
			if sort != nil {
				if sort["field"] != nil && len(sort["field"].(string)) > 0 && sort["type"] != nil && len(sort["type"].(string)) > 0 {
					sorts = append(sorts, fmt.Sprintf("%s %s", sort["field"], sort["type"]))
				}
			}
		}
		if len(sorts) > 0 {
			sortString = strings.Join(sorts, ",")
		}
	}
	page := paramInt(param["page"], 1) - 1
	limit := paramInt(param["limit"], 20)
	list, count, tcode, err := mysqlTool.ListMysql(dbName, table, whereString(param, searchTargets), sortString, page, limit)
	if tcode == 10010 {
		dealwithMysql()
		return mysqlTool.ListMysql(dbName, table, whereString(param, searchTargets), sortString, page, limit)
	}
	return list, count, tcode, err
}
