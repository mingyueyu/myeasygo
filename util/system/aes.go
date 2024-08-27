package system

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	// "time"

	"github.com/wumansgy/goEncrypt/aes"
)

var TestType = true

var aesSecretKey = ""

// 设置aesSecretKey
func AesInit(key string) {
	aesSecretKey = key
}

// 读取默认配置
func ReadDefaultSetting() (*Setting_t, error) {
	return ReadSetting("", "setting.ini")
}

// 读取配置
func ReadSetting(path string, fileName string) (*Setting_t, error) {
	if len(aesSecretKey) == 0 {
		return nil, errors.New("aesSecretKey为空,请先设置aesSecretKey")
	}
	if len(path) == 0 {
		tPath, err := myPath()
		if err != nil {
			if TestType {
				panic(err)
			}
			return nil, err
		}
		path = fmt.Sprintf("%s/%s", tPath, fileName)
	}
	f, err := os.Open(path)
	if err != nil {
		if TestType {
			panic(err)
		}
		return nil, err
	}
	defer f.Close() // 确保文件会被关闭
	// t1 := time.Now()
	// 读取文件内容
	data, err := io.ReadAll(f)
	if err != nil {
		if TestType {
			panic(err)
		}
		return nil, err
	}
	// fmt.Printf("读取文件：%s", string(data))
	unEncData, err := aes.AesCtrDecryptByBase64(string(data), []byte(aesSecretKey), nil)
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("解密失败：", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(unEncData), &Setting)
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("error:", err)
		return nil, err
	}
	// 邮箱密码解密
	Setting.Data.EmailPwd = TeaDecryptString(Setting.Data.EmailPwd)
	// redis密码解密
	Setting.Data.Redis.Pwd = TeaDecryptString(Setting.Data.Redis.Pwd)
	// 数据库密码解密
	for i := 0; i < len(Setting.Data.MySqls); i++ {
		Setting.Data.MySqls[i].Pwd = TeaDecryptString(Setting.Data.MySqls[i].Pwd)
	}

	j, err := json.MarshalIndent(Setting.Data, "", " ")
	if err != nil {
		if TestType {
			panic(err)
		}
		return nil, err
	}
	if checkConvert(Setting.Code, Setting.Name, string(j)) {
		RefreshSetting(Setting)
	} else {
		// fmt.Println("校验失败")
		return nil, errors.New("校验失败")
	}
	// t2 := time.Now()
	// 计算读取文件的耗时
	// fmt.Printf("读取文件耗时：%s\n", t2.Sub(t1))
	return &Setting, nil
}

// 设置配置
func SetSetting(codePath string, fileName string) error {
	path := fmt.Sprintf("%s/%s", codePath, fileName)
	f, err := os.Open(path)
	if err != nil {
		if TestType {
			panic(err)
		}
		// 打开文件失败
		// fmt.Println(err)
		return err
	}
	defer f.Close() // 确保文件会被关闭
	// 读取文件内容
	data, err := io.ReadAll(f)
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("读取文件时出错:", err)
		return err
	}
	set := Setting_t{}
	err = json.Unmarshal([]byte(data), &set)
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("error:", err)
		return err
	}
	// fmt.Printf("读取源文件：%s", string(data))
	jsonData, err := json.MarshalIndent(set.Data, "", " ")
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("set.Data json格式化失败：", err)
		return err
	}
	set.Code = convert(set.Name, string(jsonData))
	// 设置完成code后再进行密码加密
	// 邮箱密码加密
	set.Data.EmailPwd = TeaEncryptString(set.Data.EmailPwd)
	// redis密码加密
	set.Data.Redis.Pwd = TeaEncryptString(set.Data.Redis.Pwd)
	
	// 数据库密码加密
	for i := 0; i < len(set.Data.MySqls); i++ {
		set.Data.MySqls[i].Pwd = TeaEncryptString(set.Data.MySqls[i].Pwd)
	}
	fmt.Printf("加密前:%s", JsonString(set))
	jsonData, err = json.MarshalIndent(set, "", " ")
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("set json格式化失败：", err)
		return err
	}
	
	// 加密结果
	value, err := aes.AesCtrEncryptBase64(jsonData, []byte(aesSecretKey), nil)
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("加密失败：", err)
		return err
	}
	// fmt.Printf("设置：%s\n加密结果:%s", jsonString(set), value)
	// 保存到指定路径
	saveToExePath(value, codePath, "setting.ini")
	return saveToExeDefaultPath(value)
}

// 保存到默认路径，即程序路径
func saveToExeDefaultPath(value string) error {
	return saveToExePath(value, "", "setting.ini")
}

// 保存到指定路径
func saveToExePath(value string, path string, fileName string) error {
	if len(path) == 0 {
		tPath, err := myPath()
		if err != nil {
			if TestType {
				panic(err)
			}
			return err
		}
		path = fmt.Sprintf("%s/%s", tPath, fileName)
	} else {
		path = fmt.Sprintf("%s/%s", path, fileName)
	}

	err := os.WriteFile(path, []byte(value), 0644) // 指定文件路径、数据和权限（这里为0644）
	if err != nil {
		if TestType {
			panic(err)
		}
		return err
		// fmt.Println(err)
	}
	return nil
}

// 获取程序路径
func myPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("获取程序路径失败：", err)
		return "", err
	}
	exPath := filepath.Dir(ex)
	return exPath, nil
}

// json格式化字符串
func JsonString(mapData interface{}) string {
	j, err := json.MarshalIndent(mapData, "", " ")
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("json格式化失败：", err)
		return ""
	}
	return string(j)
}

// 转md5
func convert(name string, targetJson string) string {
	return fmt.Sprintf("%X", md5.Sum([]byte(targetJson+"2020"+name+aesSecretKey)))
}

// md5校验
func checkConvert(value string, name string, targetJson string) bool {
	if len(value) != 0 {
		if strings.Compare(value, fmt.Sprintf("%X", md5.Sum([]byte(targetJson+"2020"+name+aesSecretKey)))) == 0 {
			return true
		}
	}
	return false
}
