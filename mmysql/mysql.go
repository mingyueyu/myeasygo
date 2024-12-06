package mmysql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mingyueyu/myeasygo/mmysql/mmysqlTool"
	"github.com/mingyueyu/myeasygo/util"
)

var TestType = false

func RefreshSetting(data []byte) {
	err := json.Unmarshal([]byte(data), &mmysqlTool.Mysql)
	if err != nil {
		fmt.Println("更新mysql配置失败:", err)
	}else {
		fmt.Println("更新mysql配置成功")
	}
}

func ParamToGinH(c *gin.Context) (gin.H, error) {

	// 这样读取字节流之后，整个c.request.body就已经读空啦。再次无法读到数据。
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		if TestType {
			panic(err)
		}
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(data, &m)
	if err != nil {
		if TestType {
			panic(err)
		}
		return nil, err
	}

	// 再次读取数据（复制字节流）
	c.Request.Body = io.NopCloser(bytes.NewReader(data))
	return util.MapToGinH(m), nil
}

// type HandlerFunc func(c *gin.Context) (code int64, data interface{}, err error)
func tableNameFromeParam(param gin.H, tableName string) string {
	table := tableName
	if param["table"] != nil {
		table = fmt.Sprintf("%s_%v", table, param["table"])
	}
	if paramInt(param["year"], -1) > -1 {
		table = fmt.Sprintf("%s_%v", table, param["year"])
		if paramInt(param["month"], -1) > -1 {
			table = fmt.Sprintf("%s%02v", table, param["month"])

		}
	}
	return table
}

func paramFromGet(c *gin.Context) gin.H {
	re := gin.H{
		"search": c.Query("search"),
		"page":   c.Query("page"),
		"limit":  c.Query("limit"),
		"year":   c.Query("year"),
		"month":  c.Query("month"),
		"and":    c.QueryMap("and"),
		"or":     c.QueryMap("or"),
		"before": c.QueryMap("before"),
		"mid":    c.QueryMap("mid"),
		"after":  c.QueryMap("after"),
		"sort":   c.QueryArray("sort"),
	}
	if re["sort"] != nil {
		sorts := re["sort"].([]string)
		target := []gin.H{}
		for i := 0; i < len(sorts); i++ {
			item := sorts[i]
			list := strings.Split(item, ":")
			target = append(target, gin.H{"field": list[0], "type": list[1]})
		}
		re["sort"] = target
	}
	return re
}

func whereString(param gin.H, searchTargets []string) (string, []any) {
	whereStrings := []string{}
	whereValuelist := []any{}
	and := paramGinH(param["and"])
	if and != nil {
		if whereString, whereValues := sqlKeyValues(and, "AND"); len(whereString) > 0 {
			whereStrings = append(whereStrings, fmt.Sprintf("(%s)", whereString))
			whereValuelist = append(whereValuelist, whereValues...)
		}
	}
	if param["or"] != nil {
		orType := reflect.TypeOf(param["or"]).String()
		if strings.Compare(orType, "gin.H") == 0 {
			or := paramGinH(param["or"])
			if or != nil {
				if whereString, whereValues := sqlKeyValues(or, "OR"); len(whereString) > 0 {
					whereStrings = append(whereStrings, fmt.Sprintf("(%s)", whereString))
					whereValuelist = append(whereValuelist, whereValues...)
				}
			}
		} else if strings.Compare(orType, "[]gin.H") == 0 {
			list := param["or"].([]gin.H)
			targetList := []string{}
			for i := 0; i < len(list); i++ {
				or := list[i]
				if whereString, whereValues := sqlKeyValues(or, "AND"); len(whereString) > 0 {
					targetList = append(targetList, whereString)
					whereValuelist = append(whereValuelist, whereValues...)
				}
			}
			whereStrings = append(whereStrings, fmt.Sprintf("(%s)", strings.Join(targetList, " OR ")))
		}
	}
	
	before := paramGinH(param["before"])
	if before != nil {
		if whereString := sqlLikeKeyValues(before, true, false); len(whereString) > 0 {
			whereStrings = append(whereStrings, fmt.Sprintf("(%s)", whereString))
		}
	}
	mid := paramGinH(param["mid"])
	if mid != nil {
		if whereString := sqlLikeKeyValues(mid, true, true); len(whereString) > 0 {
			whereStrings = append(whereStrings, fmt.Sprintf("(%s)", whereString))
		}
	}
	after := paramGinH(param["after"])
	if after != nil {
		if whereString := sqlLikeKeyValues(after, false, true); len(whereString) > 0 {
			whereStrings = append(whereStrings, fmt.Sprintf("(%s)", whereString))
		}
	}
	
	if param["search"] != nil && len(param["search"].(string)) > 0 && searchTargets != nil {
		if whereString := sqlSearchValues(param["search"].(string), searchTargets); len(whereString) > 0 {
			whereStrings = append(whereStrings, fmt.Sprintf("(%s)", whereString))
		}
	}
	return strings.Join(whereStrings, " AND "), whereValuelist
}

func sqlKeyValuesFromMap(param gin.H) ([]string, []string) {
	keys, values := keysValuesFromParam(param)
	return keys, values
}

func sqlKeyValues(content gin.H, spliceStrig string) (string, []any) {
	keys, values := keysValuesFromParam(content)
	wheres := []string{}
	whereValues := []any{}
	for i := 0; i < len(keys); i++ {
		k := keys[i]
		value := values[i]
		if strings.Compare(value, "'IS NOT NULL'") == 0 {
			wheres = append(wheres, k+" IS NOT NULL")
		} else if strings.Compare(value, "'IS NULL'") == 0 || strings.Compare(value, "IS NULL") == 0 {
			wheres = append(wheres, k+" IS NULL")
		} else {
			if strings.LastIndex(k, "-") == len(k)-1 {
				k = k[:len(k)-1]
				value = k + "-" + value
				wheres = append(wheres, k+"=?")
				whereValues = append(whereValues, value)
			} else if strings.LastIndex(k, "+") == len(k)-1 {
				k = k[:len(k)-1]
				value = k + "+" + value
				wheres = append(wheres, k+"=?")
				whereValues = append(whereValues, value)
			} else if strings.LastIndex(k, ">=") == len(k)-2 || strings.LastIndex(k, "<=") == len(k)-2 || strings.LastIndex(k, ">") == len(k)-1 || strings.LastIndex(k, "<") == len(k)-1  {
				wheres = append(wheres, k+"?")
				whereValues = append(whereValues, value)
			} else {
				wheres = append(wheres, k+"=?")
				whereValues = append(whereValues, value)
			}
		}
	}
	if len(wheres) > 0 {
		return strings.Join(wheres, " "+spliceStrig+" "), whereValues
	} else {
		return "", whereValues
	}
}

func sqlLikeKeyValues(like gin.H, isBefore bool, isAfter bool) string {
	wheres := []string{}
	for k, v := range like {
		if v == nil {
			continue
		}
		if isBefore {
			v = fmt.Sprintf("%s%%", v)
		}
		if isAfter {
			v = fmt.Sprintf("%%%s", v)
		}

		wheres = append(wheres, fmt.Sprintf("%s LIKE '%s'", k, v))
	}
	if len(wheres) > 0 {
		return strings.Join(wheres, " OR ")
	} else {
		return ""
	}
}

func sqlSearchValues(search string, targets []string) string {
	instrs := []string{}
	if len(search) > 0 {
		for i := 0; i < len(targets); i++ {
			item := targets[i]
			if len(item) > 0 {
				instrs = append(instrs, fmt.Sprintf("INSTR(%s,'%s')", item, search))
			}
		}
	}
	if len(instrs) > 0 {
		return strings.Join(instrs, " OR ")
	} else {
		return ""
	}
}

func keysValuesFromParam(scene gin.H) ([]string, []string) {
	keys := []string{}
	values := []string{}
	for k, v := range scene {
		if v == nil {
			keys = append(keys, k)
			values = append(values, "IS NULL")
			continue
		}
		value := ""
		typeName := reflect.TypeOf(v).Name()
		if strings.Compare(typeName, "string") == 0 {
			value = fmt.Sprintf("%s", v)
		} else if strings.Compare(typeName, "H") == 0 || strings.Compare(typeName, "") == 0 {
			value = util.JsonString(v)
		} else {
			value = fmt.Sprintf("%v", v)
		}
		keys = append(keys, k)
		values = append(values, value)
	}
	return keys, values
}

func paramInt(value interface{}, defaultValue int64) int64 {
	if value == nil {
		return defaultValue
	}
	switch reflect.TypeOf(value).Name() {
	case "int":
		return int64(value.(int))
	case "int64":
		return value.(int64)
	case "float64":
		return int64(value.(float64))
	case "string":
		// if len(value.(string)) == 0{
		// 	return defaultValue
		// }
		i, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			if TestType {
				panic(err)
			}
			return defaultValue
		} else {
			return i
		}
	default:
		return defaultValue
	}
}

func paramGinH(value interface{}) gin.H {
	if value == nil {
		return nil
	}
	switch fmt.Sprintf("%v", reflect.TypeOf(value)) {
	case "gin.H":
		return value.(gin.H)
	case "map[string]string":
		m := make(map[string]interface{})
		for k, v := range value.(map[string]string) {
			m[k] = v
		}
		return gin.H(m)
	case "map[string]interface {}":
		return gin.H(value.(map[string]interface{}))
	case "string":
		m := make(map[string]interface{})
		err := json.Unmarshal([]byte(value.(string)), &m)
		if err != nil {
			if TestType {
				panic(err)
			}
			return nil
		} else {
			return gin.H(m)
		}
	default:
		return nil
	}
}
