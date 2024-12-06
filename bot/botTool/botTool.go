package botTool

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
    Default = iota
    Alert 
    Error
)

type Bot_t struct {
	Bots []BotDetail_t
}

type BotDetail_t struct {
	NickName string // 详情别名
	Url      string
	AlertUrl string
	ErrorUrl string
}

var Bot Bot_t

func Note(nickName string, content string, noteType int) error {
	var tbot BotDetail_t
	if nickName == "" {
        tbot = Bot.Bots[0]
    }else {
		for _, v := range Bot.Bots {
			if strings.Compare(nickName, v.NickName) == 0 {
				tbot = v
			}
		}
	}
	if tbot == (BotDetail_t{}) {
		return errors.New("bot is not initialized")
	}
	switch noteType {
	case Alert:
		return note(content, tbot.AlertUrl)
	case Error:
		return note(content, tbot.ErrorUrl)
	default:
		return note(content, tbot.Url)
	}
}

func note(content string, urlString string) error {
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
		return err
	}
	req.Header.Set("Content-Type", "application/xml")

	resp, err := client.Do(req)
	if nil != err {
		// fmt.Println("botFromUrlErr:", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}

