package mmysqlTool

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type ERROR_T struct {
	Number  int64
	Message string
}

type MySql_t struct {
	Mysqls []MySqlDetail_t
}

type MySqlDetail_t struct {
	NickName string // 详情别名
	Name     string // 数据库名称
	Host     string // 地址
	Port     int64  // 端口
	User     string // 用户
	Password string // 密码
	Tables   []Table_t
}

type Table_t struct {
	Name    string // 表名称
	Content string // 内容
}

var TestType = false
var Mysql = MySql_t{}
var dbs = gin.H{}

// 增
func AddMysql(dbName string, tableName string, keys []string, values []string) (int64, int, error) {
	wenValue := []string{}
	for i := 0; i < len(keys); i++ {
		wenValue = append(wenValue, "?")
	}
	dbString := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", tableName, strings.Join(keys, ","), strings.Join(wenValue, ","))
	args := make([]any, len(values))
	for i, v := range values {
		args[i] = v
	}
	return execute(dbName, tableName, dbString, args)
}

// 删
func DelectMysql(dbName string, tableName string, where string, whereValues []any) (int64, int, error) {
	if len(where) == 0 {
		return 0, 10003, errors.New("where 不能为空")
	}
	dbString := "DELETE FROM " + tableName + " WHERE " + where
	return execute(dbName, tableName, dbString, whereValues)
}

// 改
func UpdateMysql(dbName string, tableName string, content string, contentValues []any, where string, whereValues []any) (int64, int, error) {
	if len(where) == 0 {
		return 0, 10003, errors.New("缺少where条件")
	}
	dbString := "UPDATE " + tableName + " SET " + content + " WHERE " + where
	return execute(dbName, tableName, dbString, append(contentValues, whereValues...))
}

// 查
func ListMysql(dbName string, tableName string, where string, whereValues []any, sort string, sortValues []any, pageNumber int64, pageSize int64) ([]gin.H, int64, int, error) {
	db, tcode, err := dbFromName(dbName)
	if err != nil {
		// if TestType {
		// 	panic(err)
		// }
		return nil, 0, tcode, err
	}
	// 处理参数
	whereString := ""
	if len(where) > 0 {
		whereString = fmt.Sprintf("WHERE %s", where)
	}
	orderByString := ""
	if len(sort) > 0 {
		orderByString = fmt.Sprintf("ORDER BY %s", sort)
	}
	if pageNumber < 0 {
		pageNumber = 0
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	dbString := fmt.Sprintf("SELECT * FROM %s %s %s LIMIT %d,%d;", tableName, whereString, orderByString, pageNumber*pageSize, pageSize)
	argsList := whereValues
	if sortValues != nil {
		argsList = append(argsList, sortValues...)
	}
	// fmt.Println("dbString:", dbString, " - whereValues:", JsonString(argsList))
	rows, err := db.Query(dbString, argsList...)
	if err != nil {
		fmt.Println("sql err:", err.Error())
		errcode := errorCode(err)
		if errcode != -1 && errcode == 1146 {
			// fmt.Printf("数据表%s不存在，尝试创建数据表", tableName)
			sqlStr, err := sqlCeateFromName(dbName, tableName)
			if err != nil {
				// 没有数据库
				if TestType {
					panic(err)
				}
				return nil, 0, 10010, err
			}
			_, err = db.Query(sqlStr)
			if err != nil {
				if TestType {
					panic(err)
				}
				return nil, 0, errorCode(err), err
			} else {
				// fmt.Printf("数据表%s创建成功", tableName)
				// fmt.Printf("\n==Insert-dbString:%s\n", dbString)
				rows, err = db.Query(dbString)
				if err != nil {
					if TestType {
						panic(err)
					}
					return nil, 0, errorCode(err), err
				}
			}
		} else {
			return nil, 0, errorCode(err), err
		}
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	//var result []map[string]interface{}
	var result = []gin.H{}
	//isContentTargetOrderId := false
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		if err != nil {
			if TestType {
				panic(err)
			}
			return nil, 0, errorCode(err), err
		}
		record := make(map[string]interface{})
		for i, col := range values {
			if col != nil {
				var val interface{}
				switch v := col.(type) {
				case int64, float64, time.Time: // 根据实际数据库列类型添加更多的case
					// fmt.Println("类型: 数字")
					val = v
				case []byte:
					// fmt.Println("类型: []byte")
					val = string(v)
					var v1 interface{}
					json.Unmarshal(v, &v1)
					if v1 != nil {
						if reflect.TypeOf(v1).Name() == "" {
							val = v1
						}
					}

				default:
					// fmt.Println("类型: 其他")
					val = fmt.Sprintf("%v", v)
				}
				record[columns[i]] = val
			}
		}
		result = append(result, record)
		//isContentTargetOrderId = true
	}
	count, tcode, err := CheckCount(dbName, tableName, where, whereValues)
	if err != nil {
		if TestType {
			panic(err)
		}
		return nil, 0, tcode, err
	}
	return result, count, 0, nil
}

func dbFromName(dbName string) (*sql.DB, int, error) {
	if dbs[dbName] != nil {
		// 尝试从存储中获取已存在的数据库连接。
		target := dbs[dbName].(*sql.DB)
		if target != nil {
			var err error = nil
			// 检查数据库连接是否可用。
			err = target.Ping()
			if err != nil {
				// if TestType {
				// 	panic(err)
				// }
				// 如果数据库连接不可用，打印错误信息并关闭连接。
				// fmt.Println("数据库没连接:", err)
				// fmt.Printf("\n数据库close前连接数OpenConnections：%v\nInUse:%v\n", target.Stats().OpenConnections, target.Stats().InUse)
				target.Close()
				// 如果数据库不存在, 创建数据库
				if errorCode(err) == 1049 {
					isOK, tcode, err := createDb(dbName, "")
					if isOK {
						return dbFromName(dbName)
					} else {
						return nil, tcode, err
					}
				}
			} else {
				// 如果数据库连接可用，直接返回连接对象。
				return target, 0, nil
			}
		}

	}
	sqlString, tcode, err := targetSqlString(dbName)
	if err != nil {
		if TestType {
			panic(err)
		}
		return nil, tcode, err
	}
	// fmt.Println(sqlString)
	// 如果不存在可用的数据库连接，尝试创建新的数据库连接。
	target, err := sql.Open("mysql", sqlString)
	if err != nil {
		if TestType {
			panic(err)
		}
		return nil, errorCode(err), err
	} else {
		// 将新创建的数据库连接存储，并返回连接对象。
		dbs[dbName] = target
		return target, 0, nil
	}
}

func targetSqlString(name string) (string, int, error) {
	// 先检查是否有别名对应
	for i := 0; i < len(Mysql.Mysqls); i++ {
		item := Mysql.Mysqls[i]
		if strings.Compare(item.NickName, name) == 0 {
			return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", item.User, item.Password, item.Host, item.Port, name), 0, nil
		}
	}
	// 别名没对应的再找出第一个name对应的数据库
	for i := 0; i < len(Mysql.Mysqls); i++ {
		item := Mysql.Mysqls[i]
		if strings.Compare(item.Name, name) == 0 {
			return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", item.User, item.Password, item.Host, item.Port, name), 0, nil
		}
	}
	return "", 10010, errors.New("没配置数据库")
}

func targetSqlStringWithNoDbName(name string) string {
	for i := 0; i < len(Mysql.Mysqls); i++ {
		item := Mysql.Mysqls[i]
		if strings.Compare(item.Name, name) == 0 {
			return fmt.Sprintf("%s:%s@tcp(%s:%d)/", item.User, item.Password, item.Host, item.Port)
		}
	}
	return ""
}

// 数据库字段类型
func typeNameFromTable(dbName, table string) []gin.H {
	db, _, err := dbFromName(dbName)
	if err != nil {
		if TestType {
			panic(err)
		}
		return nil
	}
	// 查询数据库元数据获取表字段类型
	query := `
        SELECT COLUMN_NAME, DATA_TYPE
        FROM INFORMATION_SCHEMA.COLUMNS
        WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;
    `
	rows, err := db.Query(query, dbName, table)
	if err != nil {
		fmt.Println("Error querying database:", err)
		return nil
	}
	defer rows.Close()
	target := []gin.H{}
	for rows.Next() {
		var columnName, dataType string
		if err := rows.Scan(&columnName, &dataType); err != nil {
			fmt.Println("Error scanning row:", err)
		}
		fmt.Printf("Column: %s, Type: %s\n", columnName, dataType)
		if dataType == "int" || dataType == "bigint" || dataType == "tinyint" {
			target = append(target, gin.H{
				columnName: "int64",
			})
		} else if dataType == "float" || dataType == "double" {
			target = append(target, gin.H{
				columnName: "float64",
			})
		} else if dataType == "json" {
			target = append(target, gin.H{
				columnName: "json",
			})
		} else {
			target = append(target, gin.H{
				columnName: "string",
			})
		}
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over rows:", err)
		return nil
	}
	return target
}

func DetailMysql(dbName string, table string, where string, whereValues []any) (gin.H, int, error) {
	res, count, tcode, err := ListMysql(dbName, table, where, whereValues, "", nil, 0, 1)
	if err != nil {
		if TestType {
			panic(err)
		}
		return nil, tcode, err
	} else if count == 1 && len(res) == 1 {
		// 成功
		return res[0], 0, nil
	} else {
		return nil, 10020, errors.New("not found")
	}
}

func DifMysql(dbName string, tableName string, field string, where string, whereValues []any) ([]gin.H, int, error) {
	db, tcode, err := dbFromName(dbName)
	if err != nil {
		if TestType {
			panic(err)
		}
		return nil, tcode, err
	}
	orderByString := fmt.Sprintf("ORDER BY %s DESC", field)
	// 处理参数
	whereString := ""
	if len(where) > 0 {
		whereString = fmt.Sprintf("WHERE %s", where)
	}
	dbString := fmt.Sprintf("SELECT DISTINCT %s, COUNT(*) FROM %s %s GROUP BY %s %s;", field, tableName, whereString, field, orderByString)
	rows, err := db.Query(dbString, whereValues...)
	if err != nil {
		errcode := errorCode(err)
		if errcode != -1 && errcode == 1146 {
			// fmt.Printf("数据表%s不存在，尝试创建数据表", tableName)
			sqlStr, err := sqlCeateFromName(dbName, tableName)
			if err != nil {
				// 没有数据库
				if TestType {
					panic(err)
				}
				return nil, 10010, err
			}
			_, err = db.Query(sqlStr)
			if err != nil {
				if TestType {
					panic(err)
				}
				return nil, errorCode(err), err
			} else {
				// fmt.Printf("数据表%s创建成功", tableName)
				// fmt.Printf("\n==Insert-dbString:%s\n", dbString)
				rows, err = db.Query(dbString)
				if err != nil {
					if TestType {
						panic(err)
					}
					return nil, errorCode(err), err
				}
			}
		} else {
			return nil, errorCode(err), err
		}
	}
	defer rows.Close()
	//字典类型
	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var result = []gin.H{}
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, 10020, errors.New("数据不存在")
		}
		record := make(map[string]interface{})
		for i, col := range values {
			if col != nil {
				var val interface{}
				switch v := col.(type) {
				case []byte:
					val = string(v)
				case int64, float64, time.Time: // 根据实际数据库列类型添加更多的case
					val = v
				default:
					// 处理其他类型或抛出错误
					val = fmt.Sprintf("%v", v)
				}
				record[columns[i]] = val
			}
		}
		if record["COUNT(*)"] == nil {
			return nil, 10013, errors.New("数据库获取计数失败")
		}
		resultMac := gin.H{
			"count": record["COUNT(*)"].(int64),
			"value": record[field],
		}
		result = append(result, resultMac)
	}
	return result, 0, nil
}


func SumMysql(dbName string, tableName string, field string, where string, whereValues []any) (int64, int, error) {
	db, tcode, err := dbFromName(dbName)
	if err != nil {
		if TestType {
			panic(err)
		}
		return 0, tcode, err
	}
	// 处理参数
	whereString := ""
	if len(where) > 0 {
		whereString = fmt.Sprintf("WHERE %s", where)
	}
	dbString := fmt.Sprintf("SELECT SUM(%s) AS count FROM %s %s;", field, tableName, whereString)
	rows, err := db.Query(dbString, whereValues...)
	if err != nil {
		errcode := errorCode(err)
		if errcode != -1 && errcode == 1146 {
			// fmt.Printf("数据表%s不存在，尝试创建数据表", tableName)
			sqlStr, err := sqlCeateFromName(dbName, tableName)
			if err != nil {
				// 没有数据库
				if TestType {
					panic(err)
				}
				return 0, 10010, err
			}
			_, err = db.Query(sqlStr)
			if err != nil {
				if TestType {
					panic(err)
				}
				return 0, errorCode(err), err
			} else {
				// fmt.Printf("数据表%s创建成功", tableName)
				// fmt.Printf("\n==Insert-dbString:%s\n", dbString)
				rows, err = db.Query(dbString)
				if err != nil {
					if TestType {
						panic(err)
					}
					return 0, errorCode(err), err
				}
			}
		} else {
			return 0, errorCode(err), err
		}
	}
	defer rows.Close()
	//字典类型
	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// var result = []gin.H{}
	var count = int64(0)
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		if err != nil {
			return 0, 10020, errors.New("数据不存在")
		}
		record := make(map[string]interface{})
		for i, col := range values {
			if col != nil {
				var val interface{}
				switch v := col.(type) {
				case []byte:
					val = string(v)
				case int64, float64, time.Time: // 根据实际数据库列类型添加更多的case
					val = v
				default:
					// 处理其他类型或抛出错误
					val = fmt.Sprintf("%v", v)
				}
				record[columns[i]] = val
			}
		}
		if record["count"] == nil {
			return 0, 10013, errors.New("数据库获取计数失败")
		}
		str := record["count"].(string)
		count, err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			return 0, 10013, errors.New("数据库获取计数失败")
		}
		// resultMac := gin.H{
		// 	"count": record["count"].(int64),
		// }
		// result = append(result, resultMac)
	}
	return count, 0, nil
}

// 检查数量
func CheckCount(dbName string, table string, where string, whereValues []any) (int64, int, error) {
	db, tcode, err := dbFromName(dbName)
	if err != nil {
		if TestType {
			panic(err)
		}
		return 0, tcode, err
	}
	whereString := ""
	if len(where) > 0 {
		whereString = fmt.Sprintf("WHERE %s", where)
	}
	dbString := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", table, whereString)
	// fmt.Printf("CheckCount-dbString:%s", dbString)
	rows, err := db.Query(dbString, whereValues...)
	if err != nil {
		if TestType {
			panic(err)
		}
		return 0, errorCode(err), err
	}
	defer rows.Close()
	total := int64(0)
	for rows.Next() {
		err := rows.Scan(
			&total,
		)
		if err != nil {
			if TestType {
				panic(err)
			}
			// fmt.Println("GetKnowledgePointListTotal error", err)
			continue
		}
	}
	return total, 0, nil
}

// 没有数据库的创建数据库，并创建对应数据表
func createDb(dbName string, tableName string) (bool, int, error) {
	fmt.Println("数据库", dbName, "不存在，尝试创建数据库")
	dsn := targetSqlStringWithNoDbName(dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName))
	if err != nil {
		fmt.Println("数据库", dbName, "创建失败：", err)
		if TestType {
			panic(err)
		}
		return false, errorCode(err), err
	} else {
		fmt.Println("数据库", dbName, "创建成功")
		if len(tableName) > 0 {
			return createTable(dbName, tableName)
		} else {
			return true, 0, nil
		}
	}
}

func createTable(dbName string, tableName string) (bool, int, error) {
	db, _, _ := dbFromName(dbName)
	sqlStr, err := sqlCeateFromName(dbName, tableName)
	if err != nil {
		// 没有数据库
		if TestType {
			panic(err)
		}
		return false, 10010, err
	}

	_, err = db.Query(sqlStr)
	if err != nil {
		if TestType {
			panic(err)
		}
		return false, errorCode(err), err
	} else {
		return true, 0, nil
	}
}

func errorCodeMsg(e error) (int, string) {
	target := ERROR_T{}
	j, err := json.MarshalIndent(e, "", " ")
	if err != nil {
		return -1, "error转json格式化失败"
	}
	err = json.Unmarshal([]byte(j), &target)
	if err != nil || target.Number == 0 {
		return -1, "error绑定失败"
	}
	// 强制转int
	return int(target.Number), target.Message
}

func errorCode(e error) int {
	code, _ := errorCodeMsg(e)
	return code
}

// 执行MySQL
func execute(dbName string, tableName string, dbString string, params []any) (int64, int, error) {
	db, tcode, err := dbFromName(dbName)
	if err != nil {
		if TestType {
			panic(err)
		}
		return 0, tcode, err
	}
	stmt, err := db.Prepare(dbString)
	if err != nil {
		errcode := errorCode(err)
		if errcode == 1146 {
			isOK, tcode, err := createTable(dbName, tableName)
			if isOK {
				stmt, err = db.Prepare(dbString)
				if err != nil {
					if TestType {
						panic(err)
					}
					return 0, errorCode(err), err
				}
			} else {
				return 0, tcode, err
			}
		} else if errcode == 1049 {
			isOK, tcode, err := createDb(dbName, tableName)
			if isOK {
				stmt, err = db.Prepare(dbString)
				if err != nil {
					if TestType {
						panic(err)
					}
					return 0, errorCode(err), err
				}
			} else {
				return 0, tcode, err
			}
		} else {
			return 0, errorCode(err), err
		}
	}
	defer stmt.Close()
	result, err := stmt.Exec(params...)
	if err != nil {
		if TestType {
			panic(err)
		}
		fmt.Println("执行失败", dbString)
		return 0, errorCode(err), err
	}
	if strings.HasPrefix(dbString, "UPDATE") {
		count, err := result.RowsAffected()
		if err != nil {
			return 0, errorCode(err), err
		} else {
			if count == 0 {
				return 0, 10023, errors.New("没有需要更新的数据")
			}
			return count, 0, nil
		}
	} else {
		lid, err := result.LastInsertId()
		if err != nil {
			if TestType {
				panic(err)
			}
			return 0, errorCode(err), err
		} else {
			return lid, 0, nil
		}
	}
}

func sqlCreateDbFromName(dbName string) string {
	return fmt.Sprintf("CREATE DATABASE %s;", dbName)
}

func sqlCeateFromName(dbName string, tableName string) (string, error) {
	defaultTableName := sqlTableName(tableName)
	// 先取相应的数据库数据
	for i := 0; i < len(Mysql.Mysqls); i++ {
		item := Mysql.Mysqls[i]
		if strings.Compare(item.Name, dbName) == 0 {
			for j := 0; j < len(item.Tables); j++ {
				if strings.Compare(item.Tables[j].Name, defaultTableName) == 0 {
					return sqlDefaultContent(tableName, item.Tables[j].Content), nil
				}
			}
		}
	}
	return "", errors.New("未找到数据库" + dbName + "表" + tableName + "内容")
}

// 截取表名前缀
func sqlTableName(tableName string) string {
	return strings.Split(tableName, "_")[0]
}

// 假设有一个安全函数，用于检查并清洗model参数，避免SQL注入风险
func sanitizeModel(model string) (string, error) {
	// 这里应实现对model的检查，移除潜在的危险字符等
	// 若发现model不合法，返回错误
	targer := strings.ReplaceAll(strings.ToUpper(model), " ", "")
	if len(targer) == 0 || strings.Contains(targer, ";DROPTABLE") {
		return "", fmt.Errorf("invalid model")
	}
	return model, nil
}
func sqlDefaultContent(name string, field string) string {
	// 先对model进行清洗和验证
	_, err := sanitizeModel(name)
	if err != nil {
		if TestType {
			panic(err)
		}
		return "" // 返回错误空值，让调用方处理
	}
	_, err = sanitizeModel(field)
	if err != nil {
		if TestType {
			panic(err)
		}
		return "" // 返回错误空值，让调用方处理
	}
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` ("+
		"`ID` bigint NOT NULL AUTO_INCREMENT,"+
		"`infoId` varchar(16) NOT NULL COMMENT '自定义唯一标记Id',%s,"+
		"`createTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',"+
		"`modifyTime` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改日期',"+
		"PRIMARY KEY (`ID`),"+
		"UNIQUE KEY `infoId` (`infoId`)"+
		") ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8 COLLATE = utf8_bin;", name, field)
}
