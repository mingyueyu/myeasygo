package mmysql

import (
	"errors"
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
				}
			} else {
				delete(param, "IP")
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
	delete(param, "createTime")
	delete(param, "modifyTime")
	if param["content"] == nil {
		return nil, 10004, errors.New("缺少参数 content")
	}
	content := param["content"].(gin.H)
	if content["infoId"] == nil {
		content["infoId"] = mysqlTool.GetTimeLongName()
	}
	table := tableName
	if param["table"] != nil {
		table = fmt.Sprintf("%s_%v", table, param["table"])
	}
	if withYear {
		t := time.Now()
		if param["year"] != nil {
			table = fmt.Sprintf("%s_%v", table, param["year"])
		} else {
			table = fmt.Sprintf("%s_%d", tableName, t.Year())
		}
		// 有年才有月
		if withMouth {
			if param["mouth"] != nil {
				table = fmt.Sprintf("%s%02v", table, param["mouth"])
			} else {
				table = fmt.Sprintf("%s%02d", table, t.Month())
			}
		}
	}
	keyStr, valueStr := sqlKeyValuesFromMap(content)
	num, tcode, err := mysqlTool.AddMysql(dbName, table, keyStr, valueStr)
	if tcode == 10010 {
		dealwithMysql()
		num, tcode, err = mysqlTool.AddMysql(dbName, table, keyStr, valueStr)
	}
	if err != nil {
		if TestType {
			panic(err)
		}
		return nil, tcode, err
	}
	// IP 不返回
	delete(content, "IP")
	content["ID"] = num
	return content, 0, nil
}
