package system

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	// "time"

	"github.com/wumansgy/goEncrypt/aes"
)

var TestType = false

var aesSecretKey = "quxingdong202408"

func ReadSetting() bool {
	path := fmt.Sprintf("%s/setting.ini", myPath())
	f, err := os.Open(path)
	if err != nil {
		if TestType {
			panic(err)
		}
		// 打开文件失败
		// fmt.Println(err)
		return false
	}
	defer f.Close() // 确保文件会被关闭
	// t1 := time.Now()
	// 读取文件内容
	data, err := io.ReadAll(f)
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("读取文件时出错:", err)
		return false
	}
	// fmt.Printf("读取文件：%s", string(data))
	unEncData, err := aes.AesCtrDecryptByBase64(string(data), []byte(aesSecretKey), nil)
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("解密失败：", err)
		return false
	}
	err = json.Unmarshal([]byte(unEncData), &Setting)
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("error:", err)
		return false
	}
	if unConvert(Setting.Code, Setting.Name, jsonString(Setting.Data)) {
		RefreshSetting(Setting)
	} else {
		// fmt.Println("校验失败")
	}
	// t2 := time.Now()
	// 计算读取文件的耗时
	// fmt.Printf("读取文件耗时：%s\n", t2.Sub(t1))
	return true
}

func SetSetting(codePath string, inPath string) bool {
	path := fmt.Sprintf("%s/setting%s.json", codePath, inPath)
	f, err := os.Open(path)
	if err != nil {
		if TestType {
			panic(err)
		}
		// 打开文件失败
		// fmt.Println(err)
		return false
	}
	defer f.Close() // 确保文件会被关闭
	// 读取文件内容
	data, err := io.ReadAll(f)
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("读取文件时出错:", err)
		return false
	}
	set := Setting_t{}
	err = json.Unmarshal([]byte(data), &set)
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("error:", err)
		return false
	}
	// fmt.Printf("读取源文件：%s", string(data))
	jsonData, err := json.MarshalIndent(set.Data, "", " ")
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("set.Data json格式化失败：", err)
		return false
	}
	set.Code = convert(set.Name, string(jsonData))
	jsonData, err = json.MarshalIndent(set, "", " ")
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("set json格式化失败：", err)
		return false
	}
	// 加密结果
	value, err := aes.AesCtrEncryptBase64(jsonData, []byte(aesSecretKey), nil)
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("加密失败：", err)
		return false
	}
	// fmt.Printf("设置：%s\n加密结果:%s", jsonString(set), value)
	// 保存到指定路径
	saveToExePath(value)
	return true
}

func saveToExePath(value string) {
	path := fmt.Sprintf("%s/setting.ini", myPath())

	err := os.WriteFile(path, []byte(value), 0644) // 指定文件路径、数据和权限（这里为0644）
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println(err)
	} else {
		// fmt.Println("数据已成功保存到文件！")
	}
}

func myPath() string {
	ex, err := os.Executable()
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("获取程序路径失败：", err)
		return ""
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func jsonString(mapData interface{}) string {
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

func convert(name string, targetJson string) string {
	return fmt.Sprintf("%X", md5.Sum([]byte(targetJson+"2020"+name+aesSecretKey)))
}

func unConvert(value string, name string, targetJson string) bool {
	if len(value) != 0 {
		if strings.Compare(value, fmt.Sprintf("%X", md5.Sum([]byte(targetJson+"2020"+name+aesSecretKey)))) == 0 {
			return true
		}
	}
	return false
}
