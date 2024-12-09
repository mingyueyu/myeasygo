package qywxTool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Qywx_t struct {
	Qywxs  []QywxDetail_t
	CorpId string
}

type QywxDetail_t struct {
	NickName   string // 详情别名
	Host       string
	Agentid    string
	Corpsecret string
}

var Qywx Qywx_t

var defaultQywxValue QywxDetail_t
var access_token = ""

func QywxIdGinH() gin.H {
	return gin.H{
		"appid":   Qywx.CorpId,
		"agentid": defaultQywxValue.Agentid,
	}
}

// 获取access_token
func QywxAccessToken(corpsecret string) string {
	if len(corpsecret) == 0 {
		corpsecret = defaultQywxValue.Corpsecret
	}
	valueMap := QywxHttpGet(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", Qywx.CorpId, corpsecret))
	if nil == valueMap || nil == valueMap["access_token"] {
		return ""
	}
	// fmt.Println("\naccess_token to map ", valueMap)
	return valueMap["access_token"].(string)
}

// 获取成员id
func QywxUserInfoId(code string) string {
	valueMap := qywxRequestGet("https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo", "&code="+code)
	if nil == valueMap || nil == valueMap["UserId"] {
		return ""
	}
	return valueMap["UserId"].(string)
}

// 获取部门
func QywxDepartment() map[string]interface{} {
	return qywxRequestGet("https://qyapi.weixin.qq.com/cgi-bin/department/list", "")
}

// 获取子部门
func QywxSubDepartment(subDepartmentId int) map[string]interface{} {
	return qywxRequestGet("https://qyapi.weixin.qq.com/cgi-bin/department/list", "&id="+strconv.Itoa(subDepartmentId))
}

// 获取子部门ID
func QywxSubDepartmentID(subDepartmentId int) map[string]interface{} {
	return qywxRequestGet("https://qyapi.weixin.qq.com/cgi-bin/department/simplelist", "&id="+strconv.Itoa(subDepartmentId))
}

// 获取部门详情
func QywxSubDepartmentDetail(department_id int) map[string]interface{} {
	return qywxRequestGet("https://qyapi.weixin.qq.com/cgi-bin/department/get", "&department_id="+strconv.Itoa(department_id))
}

// 获取部门成员
func QywxUserSimplelist(department_id int) map[string]interface{} {
	return qywxRequestGet("https://qyapi.weixin.qq.com/cgi-bin/user/simplelist", "&department_id="+strconv.Itoa(department_id))
}

// 获取部门成员详情
func QywxUserList(department_id int) map[string]interface{} {
	return qywxRequestGet("https://qyapi.weixin.qq.com/cgi-bin/user/list", "&department_id="+strconv.Itoa(department_id))
}

// 获取成员详细信息
func QywxUserInfo(qywxUserId string) map[string]interface{} {
	return qywxRequestGet("https://qyapi.weixin.qq.com/cgi-bin/user/get", "&userid="+qywxUserId)
}

// 手机号获取userid
func QywxUserFromPhone(mobile string) map[string]interface{} {
	param := make(map[string]interface{})
	param["mobile"] = mobile
	return qywxRequestPost("https://qyapi.weixin.qq.com/cgi-bin/user/getuserid", param)
}

// 邮箱号获取userid
func QywxUserFromEmail(email string) map[string]interface{} {
	param := make(map[string]interface{})
	param["email"] = email
	param["email_type"] = 1 // 邮箱类型：1-企业邮箱（默认）；2-个人邮箱
	return qywxRequestPost("https://qyapi.weixin.qq.com/cgi-bin/user/getuserid", param)
}

// 发送通知
func QywxSendMessage(param map[string]interface{}) map[string]interface{} {
	param["agentid"] = defaultQywxValue.Agentid
	param["save"] = 0                        // 表示是否是保密消息，0表示可对外分享，1表示不能分享且内容显示水印，默认为0
	param["enable_id_trans"] = 0             // 表示是否开启id转译，0表示否，1表示是，默认0。仅第三方应用需要用到，企业自建应用可以忽略。
	param["enable_duplicate_check"] = 0      // 表示是否开启重复消息检查，0表示否，1表示是，默认0
	param["duplicate_check_interval"] = 1800 // 表示是否重复消息检查的时间间隔，默认1800s，最大不超过4小时
	return qywxRequestPost("https://qyapi.weixin.qq.com/cgi-bin/message/send", param)
}

// 发起审核
func qywxApplyEvent(param map[string]interface{}) map[string]interface{} {
	// system.Bot(fmt.Sprintf("企业微信参数：%v", JsonString(param)))
	return qywxRequestPost("https://qyapi.weixin.qq.com/cgi-bin/oa/applyevent", param)
}

// 审核详情
func qywxApprovalDetail(param map[string]interface{}) map[string]interface{} {
	return qywxRequestPost("https://qyapi.weixin.qq.com/cgi-bin/oa/getapprovaldetail", param)
}

// token过期将更新并重新请求
func qywxRequestGet(url string, urlEnd string) map[string]interface{} {
	if len(access_token) == 0 {
		access_token = QywxAccessToken(defaultQywxValue.Corpsecret)
	}
	result := QywxHttpGet(fmt.Sprintf("%s?access_token=%s%s", url, access_token, urlEnd))
	switch result["errcode"].(float64) {
	case 0:
		return result
	case 40014:
		{
			access_token = QywxAccessToken(defaultQywxValue.Corpsecret)
			result = QywxHttpGet(fmt.Sprintf("%s?access_token=%s%s", url, access_token, urlEnd))
			// system.BotAlert(fmt.Sprintf("qywxRequestGet不合法重新获取\n%v", JsonString(result)))
			break
		}
	case 42001:
		{
			access_token = QywxAccessToken(defaultQywxValue.Corpsecret)
			result = QywxHttpGet(fmt.Sprintf("%s?access_token=%s%s", url, access_token, urlEnd))
			// system.BotAlert(fmt.Sprintf("qywxRequestGet过期重新获取\n%v", JsonString(result)))
			break
		}
	default:
		{
			// system.BotAlert(fmt.Sprintf("qywxRequestGet\n%v\n%s", JsonString(result), access_token))
			break
		}
	}
	return result
}

func qywxRequestPost(url string, param map[string]interface{}) map[string]interface{} {
	if len(access_token) == 0 {
		access_token = QywxAccessToken(defaultQywxValue.Corpsecret)
	}
	result := QywxHttpPost(fmt.Sprintf("%s?access_token=%s", url, access_token), param)
	if result["errcode"] != nil {
		switch result["errcode"].(float64) {
		case 0:
			return result
		case 40014:
			{
				access_token = QywxAccessToken(defaultQywxValue.Corpsecret)
				result = QywxHttpPost(fmt.Sprintf("%s?access_token=%s", url, access_token), param)
				// system.BotAlert(fmt.Sprintf("qywxRequestPost不合法重新获取\n%v", JsonString(result)))
				break
			}
		case 42001:
			{
				access_token = QywxAccessToken(defaultQywxValue.Corpsecret)
				result = QywxHttpPost(fmt.Sprintf("%s?access_token=%s", url, access_token), param)
				// system.BotAlert(fmt.Sprintf("qywxRequestPost过期重新获取\n%v", JsonString(result)))
				break
			}
		default:
			{
				// system.BotAlert(fmt.Sprintf("qywxRequestPost\n%v", JsonString(result)))
			}
		}
		return result
	} else {
		return result
	}
}

// 网络请求
func QywxHttpGet(url string) map[string]interface{} {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// fmt.Println(string(body))
	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// fmt.Printf("json to map %v", result)
	return result
}

func QywxHttpPost(url string, param map[string]interface{}) map[string]interface{} {
	bytesData, err := json.Marshal(param)
	if err != nil {
		fmt.Printf("%v", err)
		return nil
	}
	resp, err := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer([]byte(bytesData)))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// fmt.Println(string(body))
	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// fmt.Printf("json to map %v", result)
	return result
}
