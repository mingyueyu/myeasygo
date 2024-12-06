// 全局通用
package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var mutex = &sync.Mutex{}
var key = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

// 10进制转62进制
func ToMoreShort(ten int) string {
	if ten == 0 {
		return "0"
	}
	keyCount := len(key)
	tempTen := ten
	targetString := ""
	for tempTen > 0 {
		tempNumber := tempTen % keyCount
		targetString = fmt.Sprintf("%s%s", key[tempNumber], targetString)
		tempTen -= tempNumber
		tempTen = tempTen / keyCount
	}
	return targetString
}

// 62进制转10进制
func MoreShortToTen(str string) int {
	targetTen := 0
	strCount := len(str)
	keyCount := len(key)
	for i := 0; i < strCount; i++ {
		tempStr := str[strCount-i-1 : strCount-i]
		for j := 0; j < keyCount; j++ {
			if tempStr == key[j] {
				// 找到对应位置
				mi := 1
				for k := 0; k < i; k++ {
					mi *= keyCount
				}
				targetTen += mi * j
				break
			}
		}
	}
	return targetTen
}

func GetTimeLongName() string {
	mutex.Lock()
	defer mutex.Unlock()
	t := time.Now()
	year := ToMoreShort(t.Year())
	month := ToMoreShort(int(t.Month()))
	day := ToMoreShort(t.Day())
	hour := ToMoreShort(t.Hour())
	minute := ToMoreShort(t.Minute())
	second := ToMoreShort(t.Second())
	nsecond := ToMoreShort(t.Nanosecond() / 1000)
	if len(nsecond) < 4 {
		nsecond = fmt.Sprintf("%04s", nsecond)
	}
	target := fmt.Sprintf("%s%s%s%s%s%s%s", year, month, day, hour, minute, second, nsecond)
	return target
}

// mac地址加任意数
func NextMac(mac string, count int64) (string, error) {
	mac = strings.Replace(mac, ":", "", -1)
	mac = strings.ToUpper(mac)
	if len(mac) != 12 {
		return "", errors.New("mac地址错误")
	}
	decimal, err := strconv.ParseInt(mac, 16, 64)
	if err != nil {
		return "", errors.New("mac地址内容错误")
	}
	decimal += count
	for decimal > 0xffffffffffff {
		decimal -= 0xffffffffffff + 1
	}
	return fmt.Sprintf("%012X", decimal), nil
}

// 两个Mac相差多少
func MacToMacCount(mac1 string, mac2 string) (int64, error) {
	decimal1, err := strconv.ParseInt(mac1, 16, 64)
	if err != nil {
		return 0, errors.New("mac1地址内容错误")
	}
	decimal2, err := strconv.ParseInt(mac2, 16, 64)
	if err != nil {
		return 0, errors.New("mac2地址内容错误")
	}
	return decimal2 - decimal1, nil
}

func MacInsert(mac string, connect string) string {
	if len(mac)%2 != 0 {
		mac = "0" + mac
	}
	target := []string{}
	for i := 0; i < len(mac); i += 2 {
		target = append(target, mac[i:i+2])
	}
	return strings.Join(target, connect)
}

// json格式化字符串
func JsonString(mapData interface{}) string {
	j, err := json.MarshalIndent(mapData, "", " ")
	if err != nil {
		return ""
	}
	return string(j)
}

func MapToGinH(value map[string]interface{}) gin.H {
	result := gin.H{}
	for k, v := range value {
		result[k] = v
		typName := fmt.Sprintf("%v", reflect.TypeOf(v))
		if strings.Compare(typName, "map[string]interface {}") == 0 {
			result[k] = MapToGinH(v.(map[string]interface{}))
		} else if strings.Compare(typName, "[]interface {}") == 0 {
			list := []gin.H{}
			for i := 0; i < len(v.([]interface{})); i++ {
				item := v.([]interface{})[i]
				if strings.Compare(reflect.TypeOf(item).String(), "map[string]interface {}") == 0 {
					list = append(list, MapToGinH(item.(map[string]interface{})))
				}
			}
			result[k] = list
		}
	}
	return result
}