package lib

import (
	"github.com/astaxie/beego/config"
	"sync"
)

var (
	once     sync.Once
	TCConfig Config
)

type Config struct {
	DBConnectStr string
	IpWhiteList  string
}

func InitConfig() {
	once.Do(innerInit)
}

func innerInit() {
	cf, err := config.NewConfig("ini", "config.ini")
	if err != nil {
		panic(err)
	}
	dbConnectStr := cf.DefaultString("dbconnect", "")
	ipWhiteList := cf.DefaultString("ipwhitelist", "")
	TCConfig = Config{DBConnectStr: dbConnectStr, IpWhiteList: ipWhiteList}
}
