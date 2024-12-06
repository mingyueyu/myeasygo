package system

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mingyueyu/myeasygo/bot"
	"github.com/mingyueyu/myeasygo/email"
	"github.com/mingyueyu/myeasygo/mmysql"
	"github.com/mingyueyu/myeasygo/mredis"
	"github.com/mingyueyu/myeasygo/util"
)

type SettingData_t struct {
	Url          string // 网络地址
	Port         int    // 程序端口
	Origin       string // 跨域支持的网站
	LoginTimeout int64  // 登录过期分钟数
	ImagePath    string // 本地图片路径
	ImageUrl     string // 图片地址
	FilePath     string // 本地文件路径
	FileUrl      string // 文件地址
	Other        gin.H  // 其他
}

var Setting SettingData_t

func RefreshSetting(data []byte) {
	err := json.Unmarshal([]byte(data), &Setting)
	if err != nil {
		fmt.Println("更新system配置失败:", err)
	}else {
		fmt.Println("更新system配置成功")
	}
}

func Init(aes string, tea [4]uint32) error {
	return FullInit(aes, tea, "", "")
}

func FullInit(aes string, tea [4]uint32, setPath string, getPath string) error {
	util.AesInit(aes)
	util.TeaInit(tea)
	if len(setPath) > 0 {
		// 设置配置(正式环境不需要，只要被加密的数据文件)
		setSetting(setPath)
	}
	err := getSetting(getPath)
	if err != nil {
		fmt.Println("读取配置失败：", err)
		return err
	}
	return nil
}

// 设置配置
func setSetting(path string) error {
	errStrings := []string{}
	if err := util.SetSetting(path, "", nil); err != nil {
		errStrings = append(errStrings, fmt.Sprintln("设置配置失败：", err))
	}
	if err := util.SetSetting(path, "mysql", nil); err != nil {
		errStrings = append(errStrings, fmt.Sprintln("设置mysql配置失败：", err))
	}
	if err := util.SetSetting(path, "redis", nil); err != nil {
		errStrings = append(errStrings, fmt.Sprintln("设置redis配置失败：", err))
	}
	if err := util.SetSetting(path, "email", nil); err != nil {
		errStrings = append(errStrings, fmt.Sprintln("设置email配置失败：", err))
	}
	if err := util.SetSetting(path, "bot", nil); err != nil {
		errStrings = append(errStrings, fmt.Sprintln("设置bot配置失败：", err))
	}
	if len(errStrings) > 0 {
		fmt.Println(strings.Join(errStrings, "\n"))
		return fmt.Errorf(strings.Join(errStrings, "\n"))
	}
	return nil
}

func getSetting(path string) error {
	errStrings := []string{}
	// 获取系统配置
	re, err := util.GetSetting(path, "")
	if err != nil {
		errStrings = append(errStrings, fmt.Sprintln("获取系统配置失败：", err))
	} else {
		RefreshSetting(re)
	}

	re, err = util.GetSetting(path, "mysql")
	if err != nil {
		errStrings = append(errStrings, fmt.Sprintln("获取MySql配置失败：", err))
	} else {
		mmysql.RefreshSetting(re)
	}
	re, err = util.GetSetting(path, "redis")
	if err != nil {
		errStrings = append(errStrings, fmt.Sprintln("获取redis配置失败：", err))
	} else {
		mredis.RefreshSetting(re)
	}
	re, err = util.GetSetting(path, "email")
	if err != nil {
		errStrings = append(errStrings, fmt.Sprintln("获取email配置失败：", err))
	} else {
		email.RefreshSetting(re)
	}
	re, err = util.GetSetting(path, "bot")
	if err != nil {
		errStrings = append(errStrings, fmt.Sprintln("获取bot配置失败：", err))
	} else {
		bot.RefreshSetting(re)
	}
	if len(errStrings) > 0 {
		return fmt.Errorf(strings.Join(errStrings, "\n"))
	}
	return nil
}
