package email

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mingyueyu/myeasygo/email/emailTool"
	"github.com/mingyueyu/myeasygo/util"
)

var EmailInfo emailTool.Email_t

func RefreshSetting(data []byte) {
	err := json.Unmarshal([]byte(data), &EmailInfo)
	if err != nil {
		fmt.Println("更新email配置失败:", err)
	}else {
		fmt.Println("更新email配置成功", util.JsonString(EmailInfo))
	}
}

func Email(toEmail string, toCers string, subject string, body string) (int, error) {
	return OtherEmail("", toEmail, toCers, subject, body)
}

func OtherEmail(name string, toEmail string, toCers string, subject string, body string) (int, error) {
	if  len(EmailInfo.Emails) == 0  {
        return -1, errors.New("email is not initialized")
    }
	if len(name) == 0{
		return emailTool.Email(EmailInfo.Emails[0], toEmail, toCers, subject, body)
	}else {
		for _, v := range EmailInfo.Emails {
			if strings.Compare(v.Name, name) == 0 {
				return emailTool.Email(v, toEmail, toCers, subject, body)
			}
		}
	}
	return -1, errors.New("have no email for this name")
}
