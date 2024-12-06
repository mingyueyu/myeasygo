package emailTool

import (
	"errors"
	"strings"

	"github.com/go-gomail/gomail"
)

type Email_t struct {
	Emails []EmailDetail_t
}

type EmailDetail_t struct {
	NickName   string // 详情别名
	Name       string
	Email      string
	Password   string
	ServerHost string
	ServerPort int64
}

type EmailParam struct {
	FromEmail EmailDetail_t
	// Toers 接收者邮件，如有多个，则以英文逗号(“,”)隔开，不能为空
	Toers string
	// CCers 抄送者邮件，如有多个，则以英文逗号(“,”)隔开，可以为空
	CCers string
}

// 全局变量，因为发件人账号、密码，需要在发送时才指定
var emailDetail EmailDetail_t

var m *gomail.Message

func Email(email EmailDetail_t, toEmail string, toCers string, subject string, body string) (int, error) {
	if email.Email == "" || email.Password == "" || email.ServerHost == "" {
		return -1, errors.New("email is not initialized")
	}
	// 结构体赋值
	myEmail := EmailParam{
		FromEmail: email,
		Toers:     toEmail,
		CCers:     toCers,
	}
	initEmail(myEmail)
	return sendEmail(subject, body)
}

func initEmail(ep EmailParam) {
	emails := []string{}
	emailDetail = ep.FromEmail
	m = gomail.NewMessage()
	if len(ep.Toers) == 0 {
		return
	}
	for _, tmp := range strings.Split(ep.Toers, ",") {
		emails = append(emails, strings.TrimSpace(tmp))
	}
	// 收件人可以有多个，故用此方式
	m.SetHeader("To", emails...)
	toers := []string{}
	//抄送列表
	if len(ep.CCers) != 0 {
		for _, tmp := range strings.Split(ep.CCers, ",") {
			toers = append(toers, strings.TrimSpace(tmp))
		}
		m.SetHeader("Cc", toers...)
	}
	// 发件人
	// 第三个参数为发件人别名，如"李大锤"，可以为空（此时则为邮箱名称）
	m.SetAddressHeader("From", emailDetail.Email, emailDetail.Name)
}

// SendEmail body支持html格式字符串
func sendEmail(subject string, body string) (int, error) {
	// fmt.Printf("\n发送邮件：【主题】%s\n【内容】%s\n", subject, body)
	// 主题
	m.SetHeader("Subject", subject)
	// 正文
	m.SetBody("text/html", body)
	d := gomail.NewDialer(emailDetail.ServerHost, int(emailDetail.ServerPort), emailDetail.Email, emailDetail.Password)
	// 发送
	err := d.DialAndSend(m)
	if err != nil {
		return -1, err
	}
	return 0, nil
}
