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

// func sqlBanner(tableName string) string {
// 	return sqlDefaultContent(tableName, "`title` varchar(128) COMMENT '标题',"+
// 		"`url` varchar(256) COMMENT '网址'")
// }

// func sqlCheckin(tableName string) string {
// 	return sqlDefaultContent(tableName, "`userId` varchar(128) NOT NULL COMMENT '签到者ID',"+
// 		"`score` int NOT NULL COMMENT '签到分数',"+
// 		"`type` varchar(128) NOT NULL COMMENT '类型（普签，补签，加签，福利）'")
// }

// func sqlCreateAppFeedback(tableName string) string {
// 	return sqlDefaultContent(tableName, "`platform` varchar(16) NOT NULL COMMENT '平台（如iOS，Android）',"+
// 		"`DUID` varchar(512) NOT NULL COMMENT '设备唯一标识',"+
// 		"`appName` varchar(128) NOT NULL COMMENT 'app名称',"+
// 		"`appVersion` varchar(128) NOT NULL COMMENT 'app版本号',"+
// 		"`deviceModel` varchar(128) NOT NULL COMMENT '设备型号',"+
// 		"`type` varchar(128) NOT NULL COMMENT '反馈类型（建议、功能异常、界面异常、合作、联系）',"+
// 		"`title` varchar(128) COMMENT '标题',"+
// 		"`content` text COMMENT '内容',"+
// 		"`tel` varchar(128) COMMENT '座机',"+
// 		"`phone` varchar(128) COMMENT '手机',"+
// 		"`email` varchar(256) COMMENT '邮箱',"+
// 		"`wechat` varchar(128) COMMENT '微信',"+
// 		"`qq` varchar(128) COMMENT 'QQ'")
// }

// func sqlCreateAppFeedbackReply(tableName string) string {
// 	return sqlDefaultContent(tableName, "`feedbackId` varchar(16) NOT NULL COMMENT '反馈infoId',"+
// 		"`userName` varchar(512) NOT NULL COMMENT '回复者名称',"+
// 		"`userId` varchar(128) NOT NULL COMMENT '回复者ID',"+
// 		"`ownerReply` varchar(128) NOT NULL DEFAULT '0' COMMENT '提交者回复（0为回复者回复，1为提交者回复）',"+
// 		"`content` text COMMENT '内容'")
// }

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
		") ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8;", name, field)
}
