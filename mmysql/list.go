package mmysql

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mingyueyu/myeasygo/util/system"
	"github.com/mingyueyu/myeasygo/util/system/model"
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
				c.JSON(http.StatusOK, system.ReturnFail(tcode, err.Error()))
				return
			}
			param = tparam
		}
		re, count, tcode, err := MysqlList(param, dbName, tableName, searchTargets)
		if err != nil {
			if TestType {
				panic(err)
			}
			c.JSON(http.StatusOK, system.ReturnFail(tcode, err.Error()))
		} else {
			target := system.ReturnSuccess(re)
			target["count"] = count
			// 处理返回值
			if funcResult != nil {
				tresult, tcount, tcode, err := funcResult(c, re, count)
				if err != nil {
					if TestType {
						panic(err)
					}
					c.JSON(http.StatusOK, system.ReturnFail(tcode, err.Error()))
					return
				}
				target = system.ReturnSuccess(tresult)
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
			c.JSON(http.StatusOK, system.ReturnFail(10001, err.Error()))
		} else {
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
			re, count, tcode, err := MysqlList(param, dbName, tableName, searchTargets)
			if err != nil {
				if TestType {
					panic(err)
				}
				c.JSON(http.StatusOK, system.ReturnFail(tcode, err.Error()))
			} else {
				target := system.ReturnSuccess(re)
				target["count"] = count
				// 处理返回值
				if funcResult != nil {
					tresult, tcount, tcode, err := funcResult(c, re, count)
					if err != nil {
						if TestType {
							panic(err)
						}
						c.JSON(http.StatusOK, system.ReturnFail(tcode, err.Error()))
						return
					}
					target = system.ReturnSuccess(tresult)
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
	sortString, sortValues := getSort(param["sort"])
	page := paramInt(param["page"], 1) - 1
	limit := paramInt(param["limit"], 20)
	where, whereValues := whereString(param, searchTargets)
	list, count, tcode, err := model.ListMysql(dbName, table, where, whereValues, sortString, sortValues, page, limit)
	// if tcode == 10010 {
	// 	dealwithMysql()
	// 	where, whereValues := whereString(param, searchTargets)
	// 	return system.ListMysql(dbName, table, where, whereValues, sortString, sortValues, page, limit)
	// }
	return list, count, tcode, err
}

// 获取排序数据
func getSort(sort any) (string, []any) {
	var sortString string
	var sortValues []any
	if reflect.TypeOf(sort) == reflect.TypeOf(gin.H{}) {
		sort := sort.(gin.H)
		if sort["field"] != nil && len(sort["field"].(string)) > 0 && sort["type"] != nil && len(sort["type"].(string)) > 0 {
			sortType := strings.ToUpper(sort["type"].(string))
			sortField := sort["field"].(string)
			if strings.Compare(sortType, "ASC") != 0 && strings.Compare(sortType, "DESC") != 0 {
				return sortString, sortValues
			}
			if strings.Contains(sortField, ";") || strings.Contains(sortField, ",") {
				return sortString, sortValues
			}
			// sortString = "? " + sortType
			// sortValues = append(sortValues, sortField)
			sortString = sortField + " " + sortType
		}
	} else if reflect.TypeOf(sort) == reflect.TypeOf([]gin.H{}) {
		array := sort.([]gin.H)
		sorts := []string{}
		for i := 0; i < len(array); i++ {
			item := array[i]
			if item != nil {
				if item["field"] != nil && len(item["field"].(string)) > 0 && item["type"] != nil && len(item["type"].(string)) > 0 {
					sortType := strings.ToUpper(item["type"].(string))
					sortField := item["field"].(string)
					if strings.Compare(sortType, "ASC") != 0 && strings.Compare(sortType, "DESC") != 0 {
						continue
					}
					if strings.Contains(sortField, ";") || strings.Contains(sortField, ",") {
						continue
					}
					// sorts = append(sorts, "? " + sortType)
					// sortValues = append(sortValues, sortField)
					sorts = append(sorts, sortField+" "+sortType)
				}
			}
		}
		if len(sorts) > 0 {
			sortString = strings.Join(sorts, ",")
		}
	}
	fmt.Println("sortString:", sortString, " - sortValues:", system.JsonString(sortValues))
	return sortString, sortValues
}
