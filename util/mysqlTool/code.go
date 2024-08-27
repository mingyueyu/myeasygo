package mysqlTool

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func stringFromCode(code int64) string{
    switch code {
    case 0: return "成功"
    case 1005: return "创建表失败"
    case 1006: return "创建数据库失败"
    case 1007: return "数据库已存在，创建数据库失败"
    case 1008: return "数据库不存在，删除数据库失败"
    case 1009: return "不能删除数据库文件导致删除数据库失败"
    case 1010: return "不能删除数据目录导致删除数据库失败"
    case 1011: return "删除数据库文件失败"
    case 1012: return "不能读取系统表中的记录"

    case 1016: return "文件无法打开，使用后台修复或者使用phpmyadmin进行修复。"

    case 1020: return "记录已被其他用户修改"
    case 1021: return "硬盘剩余空间不足，请加大硬盘可用空间"
    case 1022: return "关键字重复，更改记录失败"
    case 1023: return "关闭时发生错误"
    case 1024: return "读文件错误"
    case 1025: return "更改名字时发生错误"
    case 1026: return "写文件错误"

    case 1032: return "记录不存在"
    
    case 1036: return "数据表是只读的，不能对它进行修改"
    case 1037: return "系统内存不足，请重启数据库或重启服务器"
    case 1038: return "用于排序的内存不足，请增大排序缓冲区"

    case 1040: return "已到达数据库的最大连接数，请加大数据库可用连接数"
    case 1041: return "系统内存不足"
    case 1042: return "无效的主机名"
    case 1043: return "无效连接"
    case 1044: return "当前用户没有访问数据库的权限"
    case 1045: return "不能连接数据库，用户名或密码错误"

    case 1048: return "字段不能为空"
    case 1049: return "数据库不存在"
    case 1050: return "数据表已存在"
    case 1051: return "数据表不存在"

    case 1054: return "数据库缺少字段"

    case 1062: return "字段值重复，入库失败"

    case 1065: return "无效的SQL语句，SQL语句为空"

    case 1081: return "不能建立Socket连接"

    case 1114: return "数据表已满，不能容纳任何记录"

    case 1116: return "打开的数据表太多"

    case 1129: return "数据库出现异常，请重启数据库"
    case 1130: return "连接数据库失败，没有连接数据库的权限"

    case 1133: return "数据库用户不存在"

    case 1141: return "当前用户无权访问数据库"
    case 1142: return "当前用户无权访问数据表"
    case 1143: return "当前用户无权访问数据表中的字段"

    case 1146: return "数据表不存在"
    case 1147: return "未定义用户对数据表的访问权限"

    case 1149: return "SQL语句语法错误"

    case 1158: return "网络错误，出现读错误，请检查网络连接状况"
    case 1159: return "网络错误，读超时，请检查网络连接状况"
    case 1160: return "网络错误，出现写错误，请检查网络连接状况"
    case 1161: return "网络错误，写超时，请检查网络连接状况"

    case 1169: return "字段值重复，更新记录失败"

    case 1177: return "打开数据表失败"

    case 1180: return "提交事务失败"
    case 1181: return "回滚事务失败"

    case 1203: return "当前用户和数据库建立的连接已到达数据库的最大连接数，请增大可用的数据库连接数或重启数据库"
    case 1205: return "加锁超时"
    case 1211: return "当前用户没有创建用户的权限"
    case 1216: return "外键约束检查失败，更新子表记录失败"
    case 1217: return "外键约束检查失败，删除或修改主表记录失败"
    case 1226: return "当前用户使用的资源已超过所允许的资源，请重启数据库或重启服务器"
    case 1227: return "权限不足，您无权进行此操作"
    case 1235: return "MySQL版本过低，不具有本功能"
    case 1250: return "客户端不支持服务器要求的认证协议，请考虑升级客户端。"

    case 2002: return "服务器端口不对，请咨询空间商正确的端口。"
    case 2003: return "mysql服务没有启动，请启动该服务"

    // 自定义code，参数相关
    case 10000: return "参数错误"
    case 10001: return "参数格式错误"
    case 10002: return "参数不能为空"
    case 10003: return "缺少where条件"
    case 10004: return "缺少参数"
    // 数据库相关
    case 10010: return "没配置数据库"
    case 10011: return "数据不存在，不能删除"
    case 10012: return "数据不存在，不能更新"
    case 10013: return "数据库获取计数失败"
    // 数据相关
    case 10020: return "数据不存在"
    case 10021: return "数据已存在"
    case 10022: return "上传失败"


    default: return fmt.Sprintf("未知错误（%d）", code)
    }
}

func ReturnFail(code int64, data interface{}) gin.H {
    return gin.H{
        "code": code,
        "data": data,
        "msg":  stringFromCode(code),
    }
}

func ReturnSuccess(data interface{}) gin.H {
    code := int64(0)
    return gin.H{
        "code": code,
        "data": data,
        "msg":  stringFromCode(code),
    }
}
