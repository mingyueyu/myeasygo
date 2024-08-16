package mmysql

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mingyueyu/myeasygo/mmysql"
)

func Add(r *gin.Engine, relativePath string, dbName string, tableName string, withYear bool, withMouth bool, wihtIp bool) {
	r.POST(relativePath, func(c *gin.Context) {
		param, err := paramToGinH(c)
		if err != nil {
			c.JSON(http.StatusOK, mmysql.ReturnFaile(10000, err.Error()))
			return
		}else{
			add(c, param, dbName, tableName, withYear, withMouth, wihtIp)
		}
	})
}

func Delete(r *gin.Engine, relativePath string, dbName string, tableName string) {
	r.POST(relativePath, func(c *gin.Context) {
		param, err := paramToGinH(c)
		if err != nil {
			c.JSON(http.StatusOK, mmysql.ReturnFaile(10000, err.Error()))
		} else {
			del(c, param, dbName, tableName)
		}
	})
}

func Update(r *gin.Engine, relativePath string, dbName string, tableName string) {
	r.POST(relativePath, func(c *gin.Context) {
		param, err := paramToGinH(c)
		if err != nil {
			c.JSON(http.StatusOK, mmysql.ReturnFaile(10000, err.Error()))
		} else {
			update(c, param, dbName, tableName)
		}
	})
}

func List(r *gin.Engine, relativePath string, dbName string, tableName string, searchTargets []string) {
	r.GET(relativePath, func(c *gin.Context) {
		list(c, paramFromGet(c), dbName, tableName, searchTargets)
	})
	r.POST(relativePath, func(c *gin.Context) {
		param, err := paramToGinH(c)
		if err != nil {
			c.JSON(http.StatusOK, mmysql.ReturnFaile(10000, err.Error()))
		} else {
			list(c, param, dbName, tableName, searchTargets)
		}
	})
}

func Detail(r *gin.Engine, relativePath string, dbName string, tableName string) {
	r.GET(relativePath, func(c *gin.Context) {
		detail(c, paramFromGet(c), dbName, tableName)
	})
	r.POST(relativePath, func(c *gin.Context) {
		param, err := paramToGinH(c)
		if err != nil {
			c.JSON(http.StatusOK, mmysql.ReturnFaile(10000, err.Error()))
		} else {
			detail(c, param, dbName, tableName)
		}
	})
}

func Dif(r *gin.Engine, relativePath string, dbName string, tableName string) {
	r.POST(relativePath, func(c *gin.Context) {
		param, err := paramToGinH(c)
		if err != nil {
			c.JSON(http.StatusOK, mmysql.ReturnFaile(10000, err.Error()))
		} else {
			dif(c, param, dbName, tableName)
		}
	})
}

func add(c *gin.Context, param gin.H, dbName string, tableName string, withYear bool, withMouth bool, wihtIp bool) {
	if param["content"] == nil {
		c.JSON(http.StatusOK, mmysql.ReturnFaile(1002, errors.New("content is nil")))
		return
	}
	content := param["content"].(gin.H)
	if content["infoId"] == nil {
		content["infoId"] = mmysql.GetTimeLongName()
	}
	if wihtIp {
		content["IP"] = c.ClientIP()
	}
	content["IP"] = nil
	keyStr, valueStr := sqlKeyValuesFromMap(content)
	db, err := mmysql.DbFromName(dbName)
	if err != nil {
		c.JSON(http.StatusOK, mmysql.ReturnFaile(10000, err.Error()))
		return
	}
	table := tableName
	if param["table"] != nil {
		table = fmt.Sprintf("%s%v", table, param["table"])
	}
	if withYear {
		t := time.Now()
		year := fmt.Sprintf("%d", t.Year())
		month := fmt.Sprintf("%02d", t.Month())
		if param["year"] != nil {
			year = fmt.Sprintf("%v", param["year"])
		}
		if param["month"] != nil {
			year = fmt.Sprintf("%02v", param["month"])
		}
		table = fmt.Sprintf("%s_%v", table, year)
		if withMouth {
			table = fmt.Sprintf("%s%02v", table, month)
		}
	}
	num, err := mmysql.AddMysql(db, table, keyStr, valueStr)
	if err != nil || num < 0 {
		c.JSON(http.StatusOK, mmysql.ReturnFaile(10002, err.Error()))
	} else {
		// IP 不返回
		delete(content, "IP")
		c.JSON(http.StatusOK, mmysql.RequestResult(content, nil))
	}

}

func list(c *gin.Context, param gin.H, dbName string, tableName string, searchTargets []string) {
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
	db, err := mmysql.DbFromName(dbName)
	if err != nil {
		c.JSON(http.StatusOK, mmysql.ReturnFaile(10000, err.Error()))
		return
	}
	list, count, err := mmysql.ListMysql(db, table, whereString(param, searchTargets), sortString, page, limit)
	if err != nil {
		c.JSON(http.StatusOK, mmysql.ReturnFaile(10003, err.Error()))
		return
	} else {
		result := mmysql.RequestResult(list, nil)
		result["count"] = count
		c.JSON(http.StatusOK, result)
	}
}

func del(c *gin.Context, param gin.H, dbName string, tableName string) {
	table := tableNameFromeParam(param, tableName)
	db, err := mmysql.DbFromName(dbName)
	if err != nil {
		c.JSON(http.StatusOK, mmysql.ReturnFaile(10000, err.Error()))
		return
	}
	re, err := mmysql.DelectMysql(db, table, whereString(param, nil))
	if err != nil {
		c.JSON(http.StatusOK, mmysql.ReturnFaile(10003, err.Error()))
		return
	} else {
		c.JSON(http.StatusOK, re)
	}
}

func update(c *gin.Context, param gin.H, dbName string, tableName string) {
	table := tableNameFromeParam(param, tableName)
	db, err := mmysql.DbFromName(dbName)
	if err != nil {
		c.JSON(http.StatusOK, mmysql.ReturnFaile(10000, err.Error()))
		return
	}
	mmysql.UpdateMysql(db, table, sqlContentValue(param["content"].(gin.H)), whereString(param, nil))
}

func detail(c *gin.Context, param gin.H, dbName string, tableName string) {
	table := tableNameFromeParam(param, tableName)
	db, err := mmysql.DbFromName(dbName)
	if err != nil {
		c.JSON(http.StatusOK, mmysql.ReturnFaile(10000, err.Error()))
		return
	}
	res, err := mmysql.DetailMysql(db, table, whereString(param, nil))
	if err != nil {
		c.JSON(http.StatusOK, mmysql.ReturnFaile(10003, err.Error()))
		return
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func dif(c *gin.Context, param gin.H, dbName string, tableName string) {
	table := tableNameFromeParam(param, tableName)
	db, err := mmysql.DbFromName(dbName)
	if err != nil {
		c.JSON(http.StatusOK, mmysql.ReturnFaile(10000, err.Error()))
		return
	}
	if param["field"] == nil || len(param["field"].(string)) == 0 {
		c.JSON(http.StatusOK, mmysql.ReturnFaile(1002, errors.New("field参数不能为空")))
		return
	}
	result, err := mmysql.DifMysql(db, table, param["field"].(string), whereString(param, nil))
	if err != nil {
		c.JSON(http.StatusOK, mmysql.ReturnFaile(10003, err.Error()))
		return
	}
	c.JSON(http.StatusOK, result)
}

func tableNameFromeParam(param gin.H, tableName string) string {
	table := tableName
	if param["table"] != nil {
		table = fmt.Sprintf("%s%v", table, param["table"])
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
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(data, &m)
	if err != nil {
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
			return nil
		} else {
			return gin.H(m)
		}
	default:
		return nil
	}
}
