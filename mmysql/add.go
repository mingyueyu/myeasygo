package mmysql

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mingyueyu/myeasygo/util/mysqlTool"
)

func Add(r *gin.Engine, relativePath string, dbName string, tableName string) {
	AddPlus(r, relativePath, dbName, tableName, false, false, false, nil, nil)
}

func AddPro(r *gin.Engine, relativePath string, dbName string, tableName string, withYear bool, withMouth bool, wihtIp bool) {
	AddPlus(r, relativePath, dbName, tableName, withYear, withMouth, wihtIp, nil, nil)
}

func AddPlus(r *gin.Engine, relativePath string, dbName string, tableName string, withYear bool, withMouth bool, wihtIp bool, funcParam func(c *gin.Context, param gin.H) (gin.H, int, error), funcResult func(c *gin.Context, result gin.H) (gin.H, int, error)) {
	r.POST(relativePath, func(c *gin.Context) {
		param, err := ParamToGinH(c)
		if err != nil {
			if TestType {
				panic(err)
			}
			c.JSON(http.StatusOK, mysqlTool.ReturnFail(10001, err.Error()))
			return
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
			if wihtIp {
				if param["content"] != nil {
					param["content"].(gin.H)["IP"] = c.ClientIP()
					param["content"].(gin.H)["userAgent"] = c.Request.UserAgent()
				}else {
					param["IP"] = c.ClientIP()
					param["userAgent"] = c.Request.UserAgent()
				}
			} else {
				if param["content"] != nil {
					delete(param["content"].(gin.H), "IP")
					delete(param["content"].(gin.H), "userAgent")
				}else {
					delete(param, "IP")
					delete(param, "userAgent")
				}
			}
			re, tcode, err := MysqlAdd(param, dbName, tableName, withYear, withMouth)
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

func MysqlAdd(param gin.H, dbName string, tableName string, withYear bool, withMouth bool) (gin.H, int, error) {
	var content gin.H
	var table = tableName
	t := time.Now()
	var year = fmt.Sprintf("%d", t.Year())
	var month = fmt.Sprintf("%02d", t.Month())
	if param["content"] != nil {
		content = param["content"].(gin.H)
		if param["year"] != nil {
			year = fmt.Sprintf("%v", param["year"])
		}
		if param["mouth"] != nil {
			month = fmt.Sprintf("%02v", param["mouth"])
		}
	}else {
		content = param
		if param["table"] != nil {
			table = fmt.Sprintf("%s_%v", table, param["table"])
		}
	}
	delete(content, "createTime")
	delete(content, "modifyTime")
	content["infoId"] = mysqlTool.GetTimeLongName()
	if withYear {
		table = fmt.Sprintf("%s_%v", table, year)
		// 有年才有月
		if withMouth {
			table = fmt.Sprintf("%s%02v", table, month)
		}
	}
	keys, values := sqlKeyValuesFromMap(content)
	num, tcode, err := mysqlTool.AddMysql(dbName, table, keys, values)
	// if tcode == 10010 {
	// 	dealwithMysql()
	// 	num, tcode, err = mysqlTool.AddMysql(dbName, table, keys, values)
	// }
	// if tcode == 1062 { // 有重复字段
	// 	// 失败换infoId再试一次
	// 	content["infoId"] = mysqlTool.GetTimeLongName()
	// 	keys, values = sqlKeyValuesFromMap(content)
	// 	num, tcode, err = mysqlTool.AddMysql(dbName, table, keys, values)
	// }
	if err != nil {
		if TestType {
			panic(err)
		}
		return nil, tcode, err
	}
	// IP, userAgent 不返回
	delete(content, "IP")
	delete(content, "userAgent")
	content["ID"] = num
	return content, 0, nil
}
