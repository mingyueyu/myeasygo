package mysqlTool

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MySql_t struct {
	Name   string // 数据库名称
	Host   string // 地址
	Port   int64  // 端口
	User   string // 用户
	Pwd    string // 密码
	Tables []Table_t
}

type Table_t struct {
	Name    string // 表名称
	Content string // 内容
}

type ERROR_T struct {
	Number  int64
	Message string
}

var TestType = false
var CreateDbWhenNoDb = false
var mysqls = []MySql_t{}
var dbs = gin.H{}

func MysqlToolInit(params []MySql_t) {
	mysqls = params
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
				if TestType {
					panic(err)
				}
				// 如果数据库连接不可用，打印错误信息并关闭连接。
				// fmt.Println("数据库没连接:", err)
				// fmt.Printf("\n数据库close前连接数OpenConnections：%v\nInUse:%v\n", target.Stats().OpenConnections, target.Stats().InUse)
				target.Close()
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
	for i := 0; i < len(mysqls); i++ {
		item := mysqls[i]
		if strings.Compare(item.Name, name) == 0 {
			return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", item.User, item.Pwd, item.Host, item.Port, name), 0, nil
		}
	}
	return "", 10010, errors.New("没配置数据库")
}

func targetSqlStringWithNoDbName(name string) string {
	for i := 0; i < len(mysqls); i++ {
		item := mysqls[i]
		if strings.Compare(item.Name, name) == 0 {
			return fmt.Sprintf("%s:%s@tcp(%s:%d)/", item.User, item.Pwd, item.Host, item.Port)
		}
	}
	return ""
}

// 增
func AddMysql(dbName string, tableName string, keys string, values string) (int64, int, error) {
	db, tcode, err := dbFromName(dbName)
	if err != nil {
		if TestType {
			panic(err)
		}
		return -1, tcode, err
	}
	defer db.Close()
	dbString := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", tableName, keys, values)
	// fmt.Printf("数据库新增：%s", dbString)
	//	// 2. exec
	ret, err := db.Exec(dbString) //exec执行（Python中的exec就是执行字符串代码的，返回值是None，eval有返回值）
	if err != nil {
		errcode := errorCode(err)
		if errcode != -1 && errcode == 1146 {
			isOK, tcode, err := createTable(dbName, tableName)
			if isOK {
				ret, err = db.Exec(dbString)
				if err != nil {
					if TestType {
						panic(err)
					}
					return -1, errorCode(err), err
				}
			} else {
				return -1, tcode, err
			}
		} else if errcode == 1049 {
			isOK, tcode, err := createDb(dbName, tableName)
			if isOK {
				ret, err = db.Exec(dbString)
				if err != nil {
					if TestType {
						panic(err)
					}
					return -1, errorCode(err), err
				}
			} else {
				return -1, tcode, err
			}
		} else {
			// fmt.Printf("get insert id fail,err:%v\n", err)
			return -1, errorCode(err), err
		}
	}
	// 增
	// 如果是插入数据的操作，能够拿到插入数据的id
	id, err := ret.LastInsertId()
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Printf("get insert id fail,err:%v\n", err)
		return -1, errorCode(err), err
	}
	// fmt.Println("insert id:", id)
	return id, 0, nil
}

// 删
func DelectMysql(dbName string, tableName string, where string) (int64, int, error) {
	db, tcode, err := dbFromName(dbName)
	if err != nil {
		if TestType {
			panic(err)
		}
		return 0, tcode, err
	}
	if len(where) == 0 {
		return 0, 10003, errors.New("where 不能为空")
	}
	dbString := "DELETE FROM " + tableName + " WHERE " + where
	result, err := db.Exec(dbString)
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
				result, err = db.Exec(dbString)
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
	rowNum, err := result.RowsAffected()
	if err != nil {
		if TestType {
			panic(err)
		}
		return 0, errorCode(err), err
	} else if rowNum == 0 {
		return 0, 10011, errors.New("数据不存在，不能删除")
	}
	return rowNum, 0, nil
}

// 改
func UpdateMysql(dbName string, tableName string, content string, where string) (int64, int, error) {
	db, tcode, err := dbFromName(dbName)
	if err != nil {
		if TestType {
			panic(err)
		}
		return 0, tcode, err
	}
	if len(where) == 0 {
		return 0, 10003, errors.New("缺少where条件")
	}
	dbString := "UPDATE " + tableName + " SET " + content + " WHERE " + where
	result, err := db.Exec(dbString)
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
				result, err = db.Exec(dbString)
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
	rowsCount, _ := result.RowsAffected()
	// fmt.Printf("update success, affected rows:[%d]\n", rowsCount)

	return rowsCount, 0, nil
}

// 查
func ListMysql(dbName string, tableName string, where string, sort string, pageNumber int64, pageSize int64) ([]gin.H, int64, int, error) {
	db, tcode, err := dbFromName(dbName)
	if err != nil {
		if TestType {
			panic(err)
		}
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
	// fmt.Println("List-dbString", dbString)
	rows, err := db.Query(dbString)
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
		result = append(result, record)
		//isContentTargetOrderId = true
	}
	count, tcode, err := CheckCount(dbName, tableName, where)
	if err != nil {
		if TestType {
			panic(err)
		}
		return nil, 0, tcode, err
	}
	return result, count, 0, nil
}

func DetailMysql(dbName string, table string, where string) (gin.H, int, error) {
	res, count, tcode, err := ListMysql(dbName, table, where, "", 0, 1)
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

func DifMysql(dbName string, tableName string, field string, where string) ([]gin.H, int, error) {
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
	rows, err := db.Query(dbString)
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

// 检查数量
func CheckCount(dbName string, table string, where string) (int64, int, error) {
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
	rows, err := db.Query(dbString)
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
	fmt.Println("数据库",dbName,"不存在，尝试创建数据库")
	dsn := targetSqlStringWithNoDbName(dbName)
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        fmt.Println(err)
    }
    defer db.Close()
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName))
	if err != nil {
		fmt.Println("数据库",dbName,"创建失败：", err)
		if TestType {
			panic(err)
		}
		return false, errorCode(err), err
	} else {
		fmt.Println("数据库",dbName,"创建成功")
		return createTable(dbName, tableName)
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
