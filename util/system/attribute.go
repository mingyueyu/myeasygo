package system

type Setting_t struct {
	Name string
	Data SettingData_t
	Test SettingTest_t
	Custom interface{}
	Code string // 校验码
}

type SettingData_t struct {
	Url          string    // 网络地址
	Port         int       // 程序端口
	Origin       string    // 跨域支持的网站
	LoginTimeout int64     // 登录过期分钟数
	Email        string    // 发送邮箱号
	EmailPwd     string    // 邮箱密码
	ImagePath    string    // 本地图片路径
	ImageUrl     string    // 图片地址
	FilePath     string    // 本地文件路径
	FileUrl      string    // 文件地址
	BotUrl       string    // 机器人地址
	BotUrlAlert  string    // 提示机器人地址
	MySqls        []MySql_t // 数据库
	Redis        Redis_t   // redis
}

type SettingTest_t struct {
	Asker              string  // 请求者
}

type ModelType_t struct {
	Id       string
	Module   string
	Finished string
	Antenna  string
}

type MySql_t struct {
	Name  string // 数据库名称
	Host  string // 地址
	Port  int64    // 端口
	User  string // 用户
	Pwd   string // 密码
	Tables []Table_t
}

type Table_t struct {
	Name    string // 表名称
	Content string // 内容
}

type Redis_t struct {
	Host string // 地址
	Port int    // 端口
	Pwd  string // 密码
	Db   int    // 数据库
}

var Setting = Setting_t{}
var SettingData = SettingData_t{}
var SettingTest = SettingTest_t{}
var MySqls = []MySql_t{}
var Redis = Redis_t{}

func RefreshSetting(set Setting_t) {
	Setting = set
	SettingData = set.Data
	SettingTest = set.Test
	MySqls = SettingData.MySqls
	Redis = SettingData.Redis
}
