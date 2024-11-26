package mmysql

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingyueyu/myeasygo/util/system"
	"github.com/mingyueyu/myeasygo/util/system/model"
)

func Update(r *gin.Engine, relativePath string, dbName string, tableName string) {
	UpdatePlus(r, relativePath, dbName, tableName, false, nil, nil)
}

func UpdatePro(r *gin.Engine, relativePath string, dbName string, tableName string, wihtIp bool) {
	UpdatePlus(r, relativePath, dbName, tableName, wihtIp, nil, nil)
}

func UpdatePlus(r *gin.Engine, relativePath string, dbName string, tableName string, wihtIp bool, funcParam func(c *gin.Context, param gin.H) (gin.H, int, error), funcResult func(c *gin.Context, result gin.H) (gin.H, int, error)) {
	r.POST(relativePath, func(c *gin.Context) {
		param, err := ParamToGinH(c)
		if err != nil {
			if TestType {
				panic(err)
			}
			c.JSON(http.StatusOK, system.ReturnFail(10001, err.Error()))
			return
		}
		// 处理参数
		if funcParam != nil {
			tparam, tcode, err := funcParam(c, param)
			if err != nil {
				if TestType {
					panic(err)
				}
				c.JSON(http.StatusOK, system.ReturnFail(tcode, err.Error()))
				return
			}
			param = tparam
		}
		if wihtIp {
			if param["content"] != nil {
				param["content"].(gin.H)["IP"] = c.ClientIP()
				param["content"].(gin.H)["userAgent"] = c.Request.UserAgent()
			}
		} else {
			if param["content"] != nil {
				delete(param["content"].(gin.H), "IP")
				delete(param["content"].(gin.H), "userAgent")
			}
		}
		re, tcode, err := MysqlUpdate(param, dbName, tableName)
		if err != nil {
			if TestType {
				panic(err)
			}
			c.JSON(http.StatusOK, system.ReturnFail(tcode, err.Error()))
		} else {
			// 处理结果
			if funcResult != nil {
				tresult, tcode, err := funcResult(c, re)
				if err != nil {
					if TestType {
						panic(err)
					}
					c.JSON(http.StatusOK, system.ReturnFail(tcode, err.Error()))
					return
				}
				re = tresult
			}
			c.JSON(http.StatusOK, system.ReturnSuccess(re))
		}

	})
}

func MysqlUpdate(param gin.H, dbName string, tableName string) (gin.H, int, error) {
	delete(param, "createTime")
	delete(param, "modifyTime")
	table := tableNameFromeParam(param, tableName)
	content, contentValues := sqlKeyValues(param["content"].(gin.H), ",")
	where, whereValues := whereString(param, nil)
	count, tcode, err := model.UpdateMysql(dbName, table, content, contentValues, where, whereValues)
	// if tcode == 10010 {
	// 	dealwithMysql()
	// 	content, contentValues = sqlKeyValues(param["content"].(gin.H), ",")
	// 	where, whereValues := whereString(param, nil)
	// 	count, tcode, err = system.UpdateMysql(dbName, table, content, contentValues, where, whereValues)
	// }
	if err != nil {
		if TestType {
			panic(err)
		}
		return nil, tcode, err
	} else {
		return gin.H{"count": count}, 0, nil
	}
}
