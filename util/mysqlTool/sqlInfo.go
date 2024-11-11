package mysqlTool

import (
	"errors"
	"fmt"
	"strings"
)

func sqlCreateDbFromName(dbName string) string {
	return fmt.Sprintf("CREATE DATABASE %s;", dbName)
}

func sqlCeateFromName(dbName string, tableName string) (string, error) {
	defaultTableName := sqlTableName(tableName)
	// 先取相应的数据库数据
	for i := 0; i < len(mysqls); i++ {
		item := mysqls[i]
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
