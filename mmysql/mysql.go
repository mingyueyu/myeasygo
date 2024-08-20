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
)

var TestType = false

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
	return gin.H{
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
		"sort":   c.QueryMap("sort"),
	}
}

func whereString(param gin.H, searchTargets []string) string {
	whereStrings := []string{}
	and := paramGinH(param["and"])
	if and != nil {
		if whereString := sqlContentKeyValues(and, false); len(whereString) > 0 {
			whereStrings = append(whereStrings, whereString)
		}
	}
	or := paramGinH(param["or"])
	if or != nil {

		if whereString := sqlContentKeyValues(or, true); len(whereString) > 0 {
			whereStrings = append(whereStrings, whereString)
		}
	}
	before := paramGinH(param["before"])
	if before != nil {
		if whereString := sqlLikeKeyValues(before, true, false); len(whereString) > 0 {
			whereStrings = append(whereStrings, whereString)
		}
	}
	mid := paramGinH(param["mid"])
	if mid != nil {
		if whereString := sqlLikeKeyValues(mid, true, true); len(whereString) > 0 {
			whereStrings = append(whereStrings, whereString)
		}
	}
	after := paramGinH(param["after"])
	if after != nil {
		if whereString := sqlLikeKeyValues(after, false, true); len(whereString) > 0 {
			whereStrings = append(whereStrings, whereString)
		}
	}
	if param["search"] != nil && len(param["search"].(string)) > 0 && searchTargets != nil {
		if whereString := sqlSearchValues(param["search"].(string), searchTargets); len(whereString) > 0 {
			whereStrings = append(whereStrings, whereString)
		}
	}
	return strings.Join(whereStrings, " AND ")
}

func sqlKeyValuesFromMap(param gin.H) (string, string) {
	keys, values := keysValuesFromParam(param)
	return strings.Join(keys, ","), strings.Join(values, ",")
}

func sqlContentKeyValues(content gin.H, isOr bool) string {
	keys, values := keysValuesFromParam(content)
	wheres := []string{}
	for i := 0; i < len(keys); i++ {
		k := keys[i]
		value := values[i]
		wheres = append(wheres, k+"="+value)
	}
	if len(wheres) > 0 {
		if isOr {
			return fmt.Sprintf("(%s)", strings.Join(wheres, " OR "))
		} else {
			return fmt.Sprintf("(%s)", strings.Join(wheres, " AND "))
		}
	} else {
		return ""
	}
}

func sqlContentValue(content gin.H) string {
	keys, values := keysValuesFromParam(content)
	contents := []string{}
	for i := 0; i < len(keys); i++ {
		k := keys[i]
		value := values[i]
		contents = append(contents, k+"="+value)
	}
	if len(contents) > 0 {
		return strings.Join(contents, ",")
	} else {
		return ""
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

		wheres = append(wheres, fmt.Sprintf("%s LIKE \"%s\"", k, v))
	}
	if len(wheres) > 0 {
		return fmt.Sprintf("(%s)", strings.Join(wheres, " OR "))
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
				instrs = append(instrs, fmt.Sprintf("INSTR(%s,\"%s\")", item, search))
			}
		}
	}
	if len(instrs) > 0 {
		return fmt.Sprintf("(%s)", strings.Join(instrs, " OR "))
	} else {
		return ""
	}
}

func keysValuesFromParam(scene gin.H) ([]string, []string) {
	keys := []string{}
	values := []string{}
	for k, v := range scene {
		if v == nil {
			continue
		}
		value := ""
		typeName := reflect.TypeOf(v).Name()
		if strings.Compare(typeName, "string") == 0 {
			value = fmt.Sprintf("\"%s\"", v)
		} else {
			value = fmt.Sprintf("%v", v)
		}
		keys = append(keys, k)
		values = append(values, value)
	}
	return keys, values
}

func paramToGinH(c *gin.Context) (gin.H, error) {

	// 这样读取字节流之后，整个c.request.body就已经读空啦。再次无法读到数据。
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		if TestType { panic(err)}
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(data, &m)
	if err != nil {
		if TestType { panic(err)}
		return nil, err
	}

	// 再次读取数据（复制字节流）
	c.Request.Body = io.NopCloser(bytes.NewReader(data))
	return mapToGinH(m), nil
}

func mapToGinH(value map[string]interface{}) gin.H {
	result := gin.H{}
	for k, v := range value {
		result[k] = v
		typName := fmt.Sprintf("%v", reflect.TypeOf(v))
		if strings.Compare(typName, "map[string]interface {}") == 0 {
			result[k] = mapToGinH(v.(map[string]interface{}))
		}
	}
	return result
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
		i, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
		if TestType { panic(err)}
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
		if TestType { panic(err)}
			return nil
		} else {
			return gin.H(m)
		}
	default:
		return nil
	}
}

// func wrap(handler HandlerFunc) func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		code, data, err := handler(c)
// 		if err != nil {
// 			c.JSON(http.StatusOK, gin.H{
// 				"code": code,         //状态
// 				"msg":  mysqlTool.StringFromCode(code), //描述信息
// 				"data": err.Error(), // 错误信息
// 			})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{
// 			"code": 0, //状态
// 			"msg":  "Success",     //描述信息
// 			"data": data,           //数据
// 		})
// 	}
// }
