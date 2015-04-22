package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/funny/link"
	"github.com/ilahsa/tcexam/lib"
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
	flag.Parse()
	lib.InitDbConfig()
	link.DefaultConnBufferSize = *buffersize
	link.DefaultProtocol = lib.TCProtocol

	server, err := link.Listen("tcp", "0.0.0.0:10010")
	if err != nil {
		panic(err)
	}

	lib.ULogger.Info("server start:", server.Listener().Addr().String())

	server.Serve(func(session *link.Session) {
		lib.ULogger.Info("client", session.Conn().RemoteAddr().String(), "in")

		session.Process(func(msg *link.InBuffer) error {
			lib.ULogger.Info("client", session.Conn().RemoteAddr().String(), "say:", string(msg.Data))

			var dat map[string]string

			err := json.Unmarshal(msg.Data, &dat)
			if err != nil {
				lib.ULogger.Errorf("bad request,req is %s", string(msg.Data))
				return errors.New("bad request")
			} else {
				er := lib.Process(session, dat)
				//ULogger.Infof("tttt %v\n", er)
				if er != nil {
					panic(er)
				}
				return er
			}

		})

		lib.ULogger.Info("client", session.Conn().RemoteAddr().String(), "close")
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
