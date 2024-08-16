package mysqlTool

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"strings"
)

type MySql_t struct {
	Name string // 数据库名称
	Host string // 地址
	Port int64  // 端口
	User string // 用户
	Pwd  string // 密码
	Table []Table_t
}

type Table_t struct {
	Name    string // 表名称
	Content string // 内容
}

var mysqls = []MySql_t{}
var dbs = gin.H{}

func MysqlToolInit(params []MySql_t) {
	mysqls = params
}

func dbFromName(dbName string) (*sql.DB, error) {
	if dbs[dbName] != nil {
		// 尝试从存储中获取已存在的数据库连接。
		target := dbs[dbName].(*sql.DB)
		if target != nil {
			var err error = nil
			// 检查数据库连接是否可用。
			err = target.Ping()
			if err != nil {
				// 如果数据库连接不可用，打印错误信息并关闭连接。
				fmt.Println("数据库没连接:", err)
				fmt.Printf("\n数据库close前连接数OpenConnections：%v\nInUse:%v\n", target.Stats().OpenConnections, target.Stats().InUse)
				target.Close()
			} else {
				// 如果数据库连接可用，直接返回连接对象。
				return target, nil
			}
		}

	}
	sqlString, err := targetSqlString(dbName)
	if err != nil {
		return nil, err
	}
	// 如果不存在可用的数据库连接，尝试创建新的数据库连接。
	target, err := sql.Open("mysql", sqlString)
	if err != nil {
		return nil, err
	} else {
		// 将新创建的数据库连接存储，并返回连接对象。
		dbs[dbName] = target
		return target, nil
	}
}

func targetSqlString(name string) (string, error) {
	for i := 0; i < len(mysqls); i++ {
		item := mysqls[i]
		if strings.Compare(item.Name, name) == 0 {
			return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", item.User, item.Pwd, item.Host, item.Port, name), nil
		}
	}
	return "", errors.New("数据库没配置")
}

// 增
func AddMysql(dbName string, tableName string, keys string, values string) (int64, error) {
	db, err := dbFromName(dbName)
	if err != nil {
		return -1, err
	}
	defer db.Close()
	dbString := fmt.Sprintf("insert into %s(%s) values(%s)", tableName, keys, values)
	fmt.Printf("数据库新增：%s", dbString)
	//	// 2. exec
	ret, err := db.Exec(dbString) //exec执行（Python中的exec就是执行字符串代码的，返回值是None，eval有返回值）
	if err != nil {
		log.Println("err:", err)
		log.Println("error:", err.Error())
		if strings.Contains(err.Error(), "Error 1146") {
			log.Printf("数据表%s不存在，尝试创建数据表", tableName)
			_, err = db.Query(sqlCeateFromName(dbName, tableName))
			if err != nil {
				return -1, err
			} else {
				log.Printf("数据表%s创建成功", tableName)
				log.Printf("\n==Insert-dbString:%s\n", dbString)
				ret, err = db.Exec(dbString)
				if err != nil {
					return -1, err
				}
			}
		} else {
			return -1, err
		}
	}
	// 增
	// 如果是插入数据的操作，能够拿到插入数据的id
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get insert id fail,err:%v\n", err)
		return -1, err
	}
	fmt.Println("insert id:", id)
	return id, nil
}

// 删
func DelectMysql(dbName string, table string, where string) (gin.H, error) {
	db, err := dbFromName(dbName)
	if err != nil {
		return nil, err
	}
	if len(where) == 0 {
		return ReturnFaile(1149, "where 不能为空"), errors.New("where 不能为空")
	}
	dbString := "DELETE FROM " + table + " WHERE " + where
	result, err := db.Exec(dbString)
	if err != nil {
		return nil, err
	}
	rowNum, err := result.RowsAffected()
	if err != nil {
		return nil, err
	} else if rowNum == 0 {
		return nil, errors.New("没有找到需修改的数据")
	}
	return gin.H{
		"code": 0,
		"msg":  "OK",
		"data": gin.H{
			"count": rowNum,
		},
	}, nil
}

// 改
func UpdateMysql(dbName string, table string, content string, where string) (gin.H, error) {
	db, err := dbFromName(dbName)
	if err != nil {
		return nil, err
	}
	if len(where) == 0 {
		return ReturnFaile(1149, "where 不能为空"), errors.New("where 不能为空")
	}
	dbString := "UPDATE " + table + " SET " + content + " WHERE " + where
	result, err := db.Exec(dbString)
	if err != nil {
		return ReturnFaile(1149, err), err
	}
	rowsCount, _ := result.RowsAffected()
	fmt.Printf("update success, affected rows:[%d]\n", rowsCount)

	return gin.H{
		"code": 0,
		"data": nil,
		"msg":  "Success",
	}, nil
}

// 查
func ListMysql(dbName string, table string, where string, sort string, pageNumber int64, pageSize int64) ([]gin.H, int64, error) {
	db, err := dbFromName(dbName)
	if err != nil {
		return nil, 0, err
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
	dbString := fmt.Sprintf("SELECT * FROM %s %s %s LIMIT %d,%d;", table, whereString, orderByString, pageNumber*pageSize, pageSize)
	fmt.Println("List-dbString", dbString)
	rows, err := db.Query(dbString)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, err
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
	count, err := CheckCount(dbName, table, where)
	if err != nil {
		return nil, 0, err
	}
	return result, count, nil
}

func DetailMysql(dbName string, table string, where string) (gin.H, error) {
	res, count, err := ListMysql(dbName, table, where, "", 0, 1)
	if err != nil {
		return ReturnFaile(1149, err), err
	} else if count == 1 && len(res) == 1 {
		// 成功
		return gin.H{
			"code": 0,
			"data": res[0],
			"msg":  "Success",
		}, nil
	} else {
		return ReturnFaile(1149, "not found"), errors.New("not found")
	}
}

func DifMysql(dbName string, table string, field string, where string) (gin.H, error) {
	db, err := dbFromName(dbName)
	if err != nil {
		return nil, err
	}
	orderByString := fmt.Sprintf("ORDER BY %s DESC", field)
	// 处理参数
	whereString := ""
	if len(where) > 0 {
		whereString = fmt.Sprintf("WHERE %s", where)
	}
	dbString := fmt.Sprintf("SELECT DISTINCT %s, COUNT(*) FROM %s %s GROUP BY %s %s;", field, table, whereString, field, orderByString)
	log.Println("dbString:", dbString)
	rows, err := db.Query(dbString)
	if err != nil {
		return nil, err
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
	var result []gin.H
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
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
		count, err := strconv.Atoi(record["COUNT(*)"].(string))
		if err != nil {
			return nil, err
		}
		resultMac := gin.H{
			"count": count,
			"value": record[field],
		}
		result = append(result, resultMac)
	}
	return gin.H{
		"code":  0,
		"data":  result,
		"count": len(result),
		"msg":   "Success",
	}, nil
}

// 检查数量
func CheckCount(dbName string, table string, where string) (int64, error) {
	db, err := dbFromName(dbName)
	if err != nil {
		return 0, err
	}
	whereString := ""
	if len(where) > 0 {
		whereString = fmt.Sprintf("WHERE %s", where)
	}
	dbString := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", table, whereString)
	fmt.Printf("CheckCount-dbString:%s", dbString)
	rows, err := db.Query(dbString)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	total := 0
	for rows.Next() {
		err := rows.Scan(
			&total,
		)
		if err != nil {
			fmt.Println("GetKnowledgePointListTotal error", err)
			continue
		}
	}
	return int64(total), nil
}
