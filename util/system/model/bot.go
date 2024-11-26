package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Bot(content string) {
	botFromUrl(content, SettingData.BotUrl)
}

func BotErr(code int64, errString string) {
	botFromUrl(fmt.Sprintf("(code %d)%s", code, errString), SettingData.BotUrl)
}

func BotAlert(content string) {
	botFromUrl(content, SettingData.BotUrlAlert)
}

func botFromUrl(content string, urlString string) {
	// fmt.Println(content)
	//newdata为xml的字符串数据
	client := &http.Client{}
	senddata := gin.H{
		"msgtype": "text",
		"text": gin.H{
			"content": content,
		},
	}
	bs, _ := json.Marshal(senddata) //senddata为结构体数据或者json数据
	var out bytes.Buffer
	json.Indent(&out, bs, "", "\t")
	//common.PrintDebug(out.String())

	req, err := http.NewRequest("POST", urlString, strings.NewReader(out.String()))
	if nil != err {
		// fmt.Println("botFromUrlErr:", err)
		return
	}
	req.Header.Set("Content-Type", "application/xml")

	resp, err := client.Do(req)
	if nil != err {
		// fmt.Println("botFromUrlErr:", err)
		return
	}
	defer resp.Body.Close()
}
