package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

var TestType = false
var aesSecretKey = ""

func AesInit(key string) {
	aesSecretKey = key
}

func GetSetting(path string, fileName string) ([]byte, error) {
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
		if len(fileName) == 0 {
			fileName = "system"
		}
		fileName = fileName + ".conff"
		path = fmt.Sprintf("%s/conf/%s", tPath, fileName)
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
	unEncData, err := aesDecrypt(string(data), []byte(aesSecretKey))
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("解密失败：", err)
		return nil, err
	}
	setting := gin.H{}
	err = json.Unmarshal(unEncData, &setting)
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("error:", err)
		return nil, err
	}
	codeString := setting["code"].(string)
	delete(setting, "code")
	setting = passwordDecrypt(setting)
	j, err := json.MarshalIndent(setting, "", " ")
	if err != nil {
		if TestType {
			panic(err)
		}
		return nil, err
	}

	if checkConvert(codeString, "feasycom", string(j)) {
		return unEncData, nil
	} else {
		// fmt.Println("校验失败")
		return nil, errors.New("校验失败")
	}
}

func SetSetting(codePath string, fileName string, dealwithParam func(param gin.H) gin.H) error {
	if len(fileName) == 0 {
		fileName = "system"
	}
	path := fmt.Sprintf("%s/%s.json", codePath, fileName)
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
	setting := gin.H{}
	err = json.Unmarshal([]byte(data), &setting)
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("error:", err)
		return err
	}
	// fmt.Printf("读取源文件：%s", string(data))
	// 将数据结构转换为JSON格式
	jsonData, err := json.MarshalIndent(setting, "", " ")
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("set.Data json格式化失败：", err)
		return err
	}
	setting["code"] = convert("feasycom", string(jsonData))
	if dealwithParam != nil {
		setting = dealwithParam(setting)
	}
	// 密码加密
	setting = passwordEncrypt(setting)
	jsonData, err = json.MarshalIndent(setting, "", " ")
	if err != nil {
		if TestType {
			panic(err)
		}
		return err
	}
	// 加密结果
	value, err := aesEncrypt(jsonData, []byte(aesSecretKey))
	if err != nil {
		if TestType {
			panic(err)
		}
		// fmt.Println("加密失败：", err)
		return err
	}
	saveToExePath(value, codePath+"/conf", fileName)
	return saveToExePath(value, "", fileName)
}

// 密码加密
func passwordEncrypt(param map[string]interface {}) map[string]interface {} {
	for k, v := range param {
		if strings.Contains(strings.ToUpper(k), "PASSWORD") {
			param[k] = TeaEncryptString(v.(string))
		}
		if reflect.TypeOf(v).String() == "map[string]interface {}" {
			param[k] = passwordEncrypt(v.(map[string]interface {}))
		} else if reflect.TypeOf(v).String() == "[]interface {}" {
			for i := 0; i < len(v.([]interface {})); i++ {
				item := v.([]interface{})[i]
				if reflect.TypeOf(item).String() == "map[string]interface {}" {
					v.([]interface{})[i] = passwordEncrypt(item.(map[string]interface {}))
				}
			}
			param[k] = v
		}
	}
	return param
}

// 密码解密
func passwordDecrypt(param map[string]interface {}) map[string]interface {} {
	for k, v := range param {
		if strings.Contains(strings.ToUpper(k), "PASSWORD") {
			param[k] = TeaDecryptString(v.(string))
		}
		if reflect.TypeOf(v).String() == "map[string]interface {}" {
			param[k] = passwordDecrypt(v.(map[string]interface {}))
		} else if reflect.TypeOf(v).String() == "[]interface {}" {
			for i := 0; i < len(v.([]interface {})); i++ {
				item := v.([]interface{})[i]
				if reflect.TypeOf(item).String() == "map[string]interface {}" {
					v.([]interface{})[i] = passwordDecrypt(item.(map[string]interface {}))
				}
			}
			param[k] = v
		}
	}
	return param
}

// 加密函数
func aesEncrypt(plaintext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// PKCS7Padding填充
	plaintext = pKCS7Padding(plaintext, block.BlockSize())

	// 初始化向量（IV）
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// 加密模式（这里使用CBC模式）
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// Base64编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// PKCS7Padding填充
func pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 解密函数
func aesDecrypt(ciphertext string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)

	// 去除PKCS7Padding
	data = pKCS7UnPadding(data, block.BlockSize())

	return data, nil
}

// PKCS7UnPadding去除填充
func pKCS7UnPadding(data []byte, blockSize int) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

// 保存到指定路径
func saveToExePath(value string, path string, fileName string) error {
	fileName = fmt.Sprintf("%s.conff", fileName)
	if len(path) == 0 {
		tPath, err := myPath()
		if err != nil {
			if TestType {
				panic(err)
			}
			return err
		}
		path = fmt.Sprintf("%s/conf", tPath)
	}
	// 判断文件夹是否存在
	isHave, err := pathExists(path)
	if err != nil {
		if TestType {
			panic(err)
		}
		return err
	}
	if !isHave {
		// 文件夹不存在就创建文件夹
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			if TestType {
				panic(err)
			}
			return err
		}
	}
	path = fmt.Sprintf("%s/%s", path, fileName)
	err = os.WriteFile(path, []byte(value), 0644) // 指定文件路径、数据和权限（这里为0644）
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

// 判断文件是否存在
// path：要判断的文件路径
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	//当为空文件或文件夹存在
	if err == nil {
		return true, nil
	}
	//os.IsNotExist(err)为true，文件或文件夹不存在
	if os.IsNotExist(err) {
		return false, nil
	}
	//其它类型，不确定是否存在
	return false, err
}
