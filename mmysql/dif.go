package mmysql

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingyueyu/myeasygo/mmysql/mmysqlTool"
	"github.com/mingyueyu/myeasygo/util"
)

func Dif(r *gin.Engine, relativePath string, dbName string, tableName string) {
	DifPlus(r, relativePath, dbName, tableName, nil, nil)
}

func DifPlus(r *gin.Engine, relativePath string, dbName string, tableName string, funcParam func(c *gin.Context, param gin.H) (gin.H, int, error), funcResult func(c *gin.Context, result []gin.H) ([]gin.H, int, error)) {
	r.POST(relativePath, func(c *gin.Context) {
		param, err := ParamToGinH(c)
		if err != nil {
			if TestType {
				panic(err)
			}
			c.JSON(http.StatusOK, util.ReturnFail(10001, err.Error()))
		} else {
			// 处理参数
			if funcParam != nil {
				tparam, tcode, err := funcParam(c, param)
				if err != nil {
					if TestType {
						panic(err)
					}
					c.JSON(http.StatusOK, util.ReturnFail(tcode, err.Error()))
					return
				}
				param = tparam
			}
			re, tcode, err := MysqlDif(param, dbName, tableName)
			if err != nil {
				if TestType {
					panic(err)
				}
				c.JSON(http.StatusOK, util.ReturnFail(tcode, err.Error()))
			} else {
				// 处理返回值
				if funcResult != nil {
					tresult, tcode, err := funcResult(c, re)
					if err != nil {
						if TestType {
							panic(err)
						}
						c.JSON(http.StatusOK, util.ReturnFail(tcode, err.Error()))
						return
					}
					re = tresult
				}
				c.JSON(http.StatusOK, util.ReturnSuccess(re))
			}

		}
	})
}

func MysqlDif(param gin.H, dbName string, tableName string) ([]gin.H, int, error) {
	table := tableNameFromeParam(param, tableName)
	if param["field"] == nil || len(param["field"].(string)) == 0 {
		return nil, 10004, errors.New("field参数不能为空")
	}
	where, whereValues := whereString(param, nil)
	list, tcode, err := mmysqlTool.DifMysql(dbName, table, param["field"].(string), where, whereValues)
	// if tcode == 10010 {
	// 	dealwithMysql()
	// 	where, whereValues := whereString(param, nil)
	// 	list, tcode, err = system.DifMysql(dbName, table, param["field"].(string), where, whereValues)
	// }
	return list, tcode, err
}
