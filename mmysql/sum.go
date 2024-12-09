package mmysql

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingyueyu/myeasygo/mmysql/mmysqlTool"
	"github.com/mingyueyu/myeasygo/util"
)

func Sum(r *gin.Engine, relativePath string, dbName string, tableName string) {
	SumPlus(r, relativePath, dbName, tableName, nil, nil)
}

func SumPlus(r *gin.Engine, relativePath string, dbName string, tableName string, funcParam func(c *gin.Context, param gin.H) (gin.H, int, error), funcResult func(c *gin.Context, result int64) (int64, int, error)) {
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
			re, tcode, err := MysqlSum(param, dbName, tableName)
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

func MysqlSum(param gin.H, dbName string, tableName string) (int64, int, error) {
	table := tableNameFromeParam(param, tableName)
	if param["field"] == nil || len(param["field"].(string)) == 0 {
		return 0, 10004, errors.New("field参数不能为空")
	}
	where, whereValues := whereString(param, nil)
	re, tcode, err := mmysqlTool.SumMysql(dbName, table, param["field"].(string), where, whereValues)
	// if tcode == 10010 {
	// 	dealwithMysql()
	// 	where, whereValues := whereString(param, nil)
	// 	list, tcode, err = system.DifMysql(dbName, table, param["field"].(string), where, whereValues)
	// }
	return re, tcode, err
}
