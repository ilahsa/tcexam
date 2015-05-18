package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/funny/link"
	"github.com/ilahsa/tcexam/lib"
	"runtime"
)

var (
	benchmark  = flag.Bool("bench", false, "is for b1enchmark, will disable print")
	buffersize = flag.Int("buffer", 1024, "session read buffer size")
)

func log(v ...interface{}) {
	if !*benchmark {
		fmt.Println(v...)

	}
}

// This is an echo server demo work with the echo_client.
// usage:
//     go run echo_server/main.go
func main() {
	lib.InitConfig()
	//initlog
	lib.ULogger = logs.NewLogger(10000)
	lib.ULogger.SetLogger("file", `{"filename":"tcexam.log"}`)

	lib.ULogger.SetLogger("console", "")

	flag.Parse()

	runtime.GOMAXPROCS(8)
	//	link.DefaultConnBufferSize = *buffersize
	//	link.DefaultProtocol = lib.TCProtocol

	link.DefaultConfig.InBufferSize = 1024
	link.DefaultConfig.OutBufferSize = 1024
	link.DefaultConfig.SendChanSize = 10000
	pool := link.NewMemPool(10, 1, 10)

	server, err := link.Listen("tcp", "0.0.0.0:10010", lib.Protocol, pool)
	if err != nil {
		panic(err)
	}
	/// 记录系统的开始时间
	lib.Exec(`insert into user_activities(user_id,active_time,active_type,user_type,other_info) values('system',now(),'start','system','')`)
	lib.ULogger.Info("server start %s", server.Listener().Addr().String())

	server.Serve(func(session *link.Session) {
		lib.ULogger.Info("client %s %s", session.Conn().RemoteAddr().String(), "in")

		session.Process(func(msg *link.Buffer) error {
			var receiveTmp []byte
			if len(msg.Data) > 100 {
				receiveTmp = msg.Data[0:100]
			} else {
				receiveTmp = msg.Data
			}

			lib.ULogger.Info("client %s %s %s", session.Conn().RemoteAddr().String(), "say:", string(receiveTmp))

			var dat map[string]string

			err := json.Unmarshal(msg.Data, &dat)
			if err != nil {
				lib.ULogger.Error("bad request,req is %s", string(msg.Data))
				return errors.New("bad request")
			} else {
				tmpMap := map[string]string{}
				for k, v := range dat {
					if k != "file" {
						tmpMap[k] = v
					}
				}
				lib.ULogger.Info("receive %s %s", session.Conn().RemoteAddr().String(), tmpMap)
				er := lib.Process(session, dat)
				//ULogger.Infof("tttt %v\n", er)
				if er != nil {
					lib.ULogger.Error("Error: panic ", err)
				}
				return nil
			}

		})

		lib.ULogger.Info("client %s %s", session.Conn().RemoteAddr().String(), "close")
		if session.State != nil {
			u := session.State.(*lib.User)
			userId := u.Id
			add := session.Conn().RemoteAddr().String()
			if u.UserType == "P" {
				lib.VFMapInstance.DelSessionByP(session)
				lib.Exec(`insert into user_activities(user_id,active_time,active_type,user_type,other_info) values(?,now(),'end','production',?)`, userId, add)
			} else if u.UserType == "C" {
				lib.Exec(`insert into user_activities(user_id,active_time,active_type,user_type,other_info) values(?,now(),'end','customer',?)`, userId, add)
				lib.VFMapInstance.DelSessionByC(session)
			}
		}
	})
}
