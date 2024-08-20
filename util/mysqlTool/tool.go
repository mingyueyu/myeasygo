package mysqlTool

import (
	"fmt"
	"sync"
	"time"
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
    nsecond := ToMoreShort(t.Nanosecond()/1000000)
    if len(nsecond) < 2 {
		nsecond = fmt.Sprintf("%02s", nsecond)
	}
    target := fmt.Sprintf("%s%s%s%s%s%s%s", year, month, day, hour, minute, second, nsecond)
    return target
}