package lib

import (
	"github.com/astaxie/beego/logs"
)

var ULogger *logs.BeeLogger

func InitLog() {
	ULogger := logs.NewLogger(10000)
	ULogger.SetLogger("file", `{"filename":"tcexam.log"}`)
	ULogger.SetLogger("console", "")
	ULogger.SetLevel(6)
}

//package lib

//import (
//	"fmt"

//	seelog "github.com/cihub/seelog"
//)

//var ULogger seelog.LoggerInterface

//func loadAppConfig() {
//	appConfig := `
//<seelog type="adaptive" mininterval="2000000" maxinterval="100000000" critmsgcount="500" minlevel="debug">
//    <exceptions>
//        <exception filepattern="test*" minlevel="error"/>
//    </exceptions>
//    <outputs formatid="all">
//        <file path="all.log"/>
//        <filter levels="info">
//          <console formatid="fmtinfo"/>
//        </filter>
//        <filter levels="error,critical" formatid="fmterror">
//          <console/>
//          <file path="errors.log"/>
//        </filter>
//    </outputs>
//    <formats>
//        <format id="fmtinfo" format="[%Level] [%Date(01-02) %Time] %Msg%n"/>
//        <format id="fmterror" format="[%LEVEL] [%Date(01-02) %Time] [%FuncShort @ %File.%Line] %Msg%n"/>
//        <format id="all" format="[%Level] [%Date(01-02) %Time] [@ %File.%Line] %Msg%n"/>
//        <format id="criticalemail" format="Critical error on our server!\n    %Time %Date %RelFile %Func %Msg \nSent by Seelog"/>
//    </formats>
//</seelog>
//`
//	logger, err := seelog.LoggerFromConfigAsBytes([]byte(appConfig))
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	UseLogger(logger)
//}

//func init() {
//	DisableLog()
//	loadAppConfig()
//}

//// DisableLog disables all library log output
//func DisableLog() {
//	ULogger = seelog.Disabled

//}

//// UseLogger uses a specified seelog.LoggerInterface to output library log.
//// Use this func if you are using Seelog logging system in your app.
//func UseLogger(newLogger seelog.LoggerInterface) {
//	ULogger = newLogger
//}
